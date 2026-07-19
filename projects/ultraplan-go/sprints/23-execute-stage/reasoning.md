# Sprint Reasoning: Execute Stage

> Project: `ultraplan-go`
> Sprint: `23-execute-stage`
> Output: `projects/ultraplan-go/sprints/23-execute-stage/reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/23-execute-stage/requirements.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/sprints/23-execute-stage/sprint-index.md`, `projects/ultraplan-go/sprints/23-execute-stage/technical-handbook.md`, `projects/ultraplan-go/sprints/23-execute-stage/reasoning/architecture.md`

This document decides. It synthesizes selected context, handbook evidence, area-specific reasoning, and contracts into final sprint decisions.

It does not replace `sprint-index.md`, `technical-handbook.md`, or `reasoning/*.md`.

## Sprint Purpose

- **Goal:** Execute validated `plan.md` implementation tasks through the generic runtime boundary with durable `.run-state.json`, safe target-repository constraints, resumability, status visibility, and actionable diagnostics.
- **Non-Goals:** Smoke investigation, automated review artifacts, issue tracking, Git state mutation, hosted/browser UI, TUI behavior, cross-sprint scheduling, a generic workflow engine, and reuse of study services or study-specific models for execute behavior.
- **Depends On:** Valid prerequisites through `plan.md`, selected sprint context in `sprint-index.md`, handbook evidence, completed architecture reasoning, existing project/sprint/workspace modules, existing generic runtime integration, and the explicit target implementation repository `../ultraplan-go`.

## Selected Context And Pre-Reasoning Artifacts

| Artifact | Path | How It Was Used |
| --- | --- | --- |
| Requirements | `projects/ultraplan-go/sprints/23-execute-stage/requirements.md` | Defined the sprint goal, required implementation outputs, acceptance criteria, non-goals, constraints, dependencies, and review expectations for execute-stage work. |
| Project Index | `projects/ultraplan-go/project-index.md` | Confirmed the active contract pool, selected study evidence catalog, module ownership rules, non-goals, target implementation directory, and required review protocols. |
| Project Docs | `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md` | Supplied product principles, execute data flow, module boundaries, runtime/agentwrap constraints, flow-state and run-state requirements, command requirements, validation rules, security, observability, and testing requirements. |
| Sprint Index | `projects/ultraplan-go/sprints/23-execute-stage/sprint-index.md` | Selected contracts, evidence reports, architecture reasoning, review protocols, and exclusions. It constrained reasoning to execute from validated `plan.md` and excluded smoke, review, issue, Git mutation, TUI, study semantics, and generic workflow engine behavior. |
| Technical Handbook | `projects/ultraplan-go/sprints/23-execute-stage/technical-handbook.md` | Provided high-confidence evidence for thin command boundaries, module-local ownership, explicit composition, config precedence, validation-gated completion, durable state, safe observability, path safety, bounded work, and fake-runtime tests. |
| Architecture Reasoning | `projects/ultraplan-go/sprints/23-execute-stage/reasoning/architecture.md` | Selected the sprint-owned execute architecture, rejected global scheduler/study reuse/platform-owned execute semantics, and resolved core architecture trade-offs before final decisions. |

## Area-Specific Reasoning Inputs

| Area | Reasoning Document | Key Conclusion | Evidence Basis | Impact On Final Decision |
| --- | --- | --- | --- | --- |
| Architecture | `projects/ultraplan-go/sprints/23-execute-stage/reasoning/architecture.md` | Implement execute as a sprint-owned use case with sprint-local task extraction, prompts, validation, `.run-state.json`, `execute.md`, flow/status integration, and runtime calls delegated to `internal/platform/runtime` and agentwrap. | `technical-handbook.md` patterns and warnings; `ARCHITECTURE.md` sprint/runtime ownership rules; `TRD.md` sections 18.1 through 18.8; selected review protocols. | Final decisions preserve `internal/sprint` ownership, keep runtime generic, use `.run-state.json` as source of truth, constrain target paths, default to conservative execution, and reject shared workflow abstractions. |

## Sprint Technical Handbook Summary

- **Relevant Patterns:** Thin command wrappers over product-owned behavior; module-local ownership; explicit composition roots and narrow volatile interfaces; deterministic config precedence; validation-gated success; atomic durable state; context cancellation; structured redacted diagnostics; workspace-safe execution; fake-runtime and golden/fixture tests.
- **Important Trade-Offs:** Thin CLI wrappers add internal request/result types; sprint-local execute state duplicates some study run-loop code shape; validation-gated completion requires more evidence rules; conservative execution reduces throughput; writing state after meaningful transitions increases file-write frequency; path/permission safety can block convenient local workflows.
- **Warnings / Anti-Patterns:** Do not put execute in a monolithic command handler; do not create global scheduler/workflow packages; do not rely on package globals for runtime/config; do not bypass IO/test seams; do not mark success from runtime exit alone; do not spawn unbounded runtime tasks; do not persist unsafe raw runtime payloads; do not add smoke, review, issue, or Git mutation behavior.
- **Evidence Confidence:** High for architecture, command, DI, config, errors, IO, state/context, concurrency, observability, testing, security, and performance patterns because the handbook cites concrete mature Go CLI examples. Medium for terminal UX and philosophy because those reports are more contextual, but still useful for status and non-goal discipline.

## Contracts Applied

| Contract / Requirement ID | Constraint | Decision Impact | Expected Evidence |
| --- | --- | --- | --- |
| Architecture | Preserve `sprint -> project`, `sprint -> workspace`, `sprint -> platform/runtime`, `project -> workspace`, and `platform/* -> no product modules`; keep execute semantics in `internal/sprint`. | Execute task extraction, prompt rendering, validation, run state, summary, status, and flow remain sprint-owned; runtime package remains product-agnostic. | Import review, architecture review, package-level tests, and no product imports from `internal/platform/runtime`. |
| CLI Surface | Add/update `validate execute`, `prompt execute`, `flow --to execute`, `execute`, task selection, dry-run, help, exit codes, text/status behavior. | CLI/app wrappers parse requests and delegate to sprint use cases; unsupported stages are rejected. | Command tests for help, invalid args, dry-run/prompt, prerequisite failure, selected task, exit codes, and status output. |
| Configuration | Support global/default sprint model plus stage-specific overrides; `execute` override wins; diagnostics must redact secrets. | Model resolution is deterministic and product-owned before runtime request construction. | Unit tests for precedence, unknown stage keys, empty models, command overrides, and safe diagnostics. |
| Documentation | Generated execute artifacts must be editable/reviewable; state JSON must be versioned. | `execute.md` is a human summary citing `plan.md`; `.run-state.json` is machine state and source of truth. | `execute.md` validator tests and documentation/help review. |
| Errors | Prerequisite, plan extraction, runtime, validation, stale state, atomic write, and target path failures need actionable diagnostics. | Errors are classified and wrapped with operation/path/task context; status and `execute.md` point to recovery information. | Unit and command tests for each failure class and exit-code mapping. |
| LLM Runtime | Use existing generic runtime boundary backed by agentwrap/OpenCode; do not parse OpenCode stdout/stderr in product code. | Sprint builds generic runtime requests and consumes generic results/events/metadata only. | Import review, fake runtime tests, and gated real-runtime checks through agentwrap when explicitly run. |
| LLM Evaluation / Cost / Safety | Runtime success is insufficient; safe metadata may be recorded; unsafe payloads and secrets are omitted/redacted. | Completion requires evidence or explicit diagnostic; metadata summaries are safe and bounded. | Fake-runtime tests for success without evidence, redaction checks, metadata summary assertions. |
| Observability | Status, logs, task counts, attempts, diagnostics, and runtime summaries must be inspectable without runtime calls. | `.run-state.json` records task state and safe metadata; status computes from state, not runtime. | Status command tests and structured/text output checks. |
| Performance | Avoid unbounded runtime fan-out, expensive target scans, repeated large reads, and status calls that invoke runtime. | Default execution is serial for target safety; status reads bounded sprint artifacts and run state only. | Tests/benchmarks for state/status on larger fixture task counts if needed; code review for no recursive target scans in status. |
| Persistence And Migrations | `.run-state.json` and `flow-state.json` require schema versions, strict load validation, atomic writes, and recovery. | New execute state schema is versioned; malformed/unsupported state fails; stale running tasks recover conservatively. | State-store tests for atomic writes, schema validation, malformed JSON, unsupported versions, stale recovery, and last-known-good preservation. |
| Security | Target writes must be constrained to `../ultraplan-go`; secrets and unsafe payloads must be redacted; no Git mutation. | Execute validates target path before runtime and excludes Git mutation behavior from prompts/commands. | Target path tests, redaction tests, permission policy review, and non-goal regression tests. |
| Testing | Normal verification must be deterministic, offline, fake-runtime based; `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan` must pass. | Tests prioritize unit, fixture, command, fake-runtime, and golden outputs; real runtime is optional/gated. | Required commands from `../ultraplan-go` plus focused package tests. |
| Workflows | Execute is a stateful workflow stage after valid `plan.md`; it must support pending/running/complete/failed/cancelled, attempts, resumability, cancellation, and summaries. | Flow includes execute only after valid prerequisites; complete tasks are not redone on resume unless explicitly forced later. | Workflow tests for flow gating, resume, cancellation, stale recovery, task selection, and terminal summaries. |
| REQ-23-40 | `execute` and `flow --to execute` require valid prerequisites through `plan.md`. | Execute validates prerequisites before task extraction or runtime launch. | Command/use-case tests for missing and invalid prerequisites. |
| REQ-23-41 | Runtime model resolution supports global/default and stage-specific overrides; stage-specific wins. | Implement deterministic stage model resolution for `execute` and planning stages. | Config/model-resolution tests. |
| REQ-23-42 | Task extraction reads only validated `plan.md`, produces deterministic IDs, and preserves traceability. | Task IDs derive from stable normalized plan task fields and include plan/reasoning references. | Task extraction fixture tests and duplicate/ambiguous task diagnostics. |
| REQ-23-43 | Runtime-backed execute uses only generic runtime boundary through sprint/app seams. | No direct OpenCode calls, shell runtime calls, native stream parsing, or study service imports. | Import/code review and fake-runtime tests. |
| REQ-23-44 | Target repository paths are explicit, workspace-safe, and constrained to `../ultraplan-go` unless future config permits otherwise. | Target path validation is mandatory before runtime launch. | Path containment and escape tests. |
| REQ-23-45 | Runtime success alone is insufficient. | Complete only with expected evidence or explicit safe diagnostic. | Fake-runtime tests for runtime success with/without evidence. |
| REQ-23-46 / REQ-23-47 | `.run-state.json` is versioned, atomic, strict, resumable, and recovers stale running tasks. | Execute state is the durable source of truth with conservative recovery. | State lifecycle tests. |
| REQ-23-48 / REQ-23-49 | Status and `execute.md` reflect readiness/progress and terminal states without runtime calls. | Status reads artifacts/state; `execute.md` is a regenerated summary citing `plan.md`. | Status and summary golden/command tests. |
| REQ-23-50 / REQ-23-51 | Verification passes with fake runtimes; forbidden behaviors remain absent. | Normal test suite excludes real runtime and forbids smoke/review/issues/Git mutation. | `go test ./...`, `go test -race ./...`, `go build ./cmd/ultraplan`, and non-goal regression checks. |

## Repos Studied / Source Evidence Used

| Source / Repo / Report | Concrete Reference | Relevant Finding | Why It Matters For This Sprint | Used In Decision(s) |
| --- | --- | --- | --- | --- |
| `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md`; `cmd/gh/main.go:6`, `cmd/root.go:112-127`, `internal/restic/repository.go:18` cited by handbook | Mature Go CLIs keep entrypoints thin and use internal package boundaries. | Supports app/CLI delegation and sprint-owned execute semantics. | Decisions 1, 8 |
| `02-command-architecture` | `studies/go-cli-study/reports/final/02-command-architecture.md`; `pkg/cmd/install.go:132-145`, `pkg/cmdutil/factory.go:16-43` cited by handbook | Command factories/options keep workflow logic out of command handlers. | Supports thin execute, prompt, validate, flow, and status command wiring. | Decisions 1, 8 |
| `03-dependency-injection` | `studies/go-cli-study/reports/final/03-dependency-injection.md`; `internal/cmd/config.go:362`, `executor.go:22-24`, `internal/backend/backend.go:19-90` cited by handbook | Composition roots, constructor injection, and narrow interfaces improve testability. | Supports fake runtime, state store, clock, and output seams without global state. | Decisions 1, 2, 8 |
| `04-configuration-management` | `studies/go-cli-study/reports/final/04-configuration-management.md`; `internal/cmd/config.go:2253-2287`, `internal/flags/flags.go:314-327` cited by handbook | Config precedence and changed-flag tracking should be explicit and validated after merge. | Directly supports execute model override precedence and diagnostics. | Decision 7 |
| `05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md`; `fs/fserrors/error.go:22-29`, `internal/errors/fatal.go:10` cited by handbook | Wrapped/classified errors and exit mapping improve actionable CLI behavior. | Shapes prerequisite, state, runtime, evidence, and path diagnostics. | Decisions 2, 3, 4, 8 |
| `06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md`; `pkg/iostreams/iostreams.go:551-568`, `internal/fs/interface.go:10-31` cited by handbook | IO/filesystem seams enable deterministic command and persistence tests. | Supports state-store tests, status output capture, and summary golden tests. | Decisions 3, 6, 8 |
| `07-state-context` | `studies/go-cli-study/reports/final/07-state-context.md`; `task.go:89`, `internal/session/session.go:12-23`, `internal/restic/lock.go:105` cited by handbook | Context propagation, cancellation, and explicit persistent task/session state are required for long operations. | Shapes execute cancellation, attempts, stale running recovery, and resumability. | Decisions 3, 6 |
| `08-concurrency` | `studies/go-cli-study/reports/final/08-concurrency.md`; `task.go:87`, `gdu/pkg/analyze/parallel.go:13`, `opencode/cmd/root.go:130,252` cited by handbook | Bounded structured concurrency avoids leaks and unbounded task growth. | Supports serial default and bounded future execution. | Decisions 6, 9 |
| `09-terminal-ux` | `studies/go-cli-study/reports/final/09-terminal-ux.md`; `internal/cmd/prompt.go:20-256`, `signals.go:11-31` cited by handbook | Long operations need calm progress, cancellation, and non-TTY/script behavior. | Shapes status and execute summary requirements. | Decision 5, 8 |
| `10-logging-observability` | `studies/go-cli-study/reports/final/10-logging-observability.md`; `internal/logging/logging.go:31-66`, `internal/slogs/keys.go:6-231` cited by handbook | Structured fields, stdout/stderr separation, and debug controls support recovery. | Shapes run-state metadata summaries and safe diagnostics. | Decisions 3, 5 |
| `11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md`; `internal/cmd/main_test.go:64-174`, `pkg/httpmock/stub.go:35-199`, `task_test.go:166-169` cited by handbook | Table, command, fake/mock, fixture, and golden tests are the right confidence strategy. | Shapes normal verification and fake-runtime requirements. | Decision 10 |
| `13-security` | `studies/go-cli-study/reports/final/13-security.md`; `internal/permission/permission.go:44-108`, `internal/llm/tools/bash.go:41-55`, `pkg/registry/transport.go:37-41` cited by handbook | Trust boundaries, permissions, path safety, and redaction must be explicit. | Shapes target path validation, permission policy, and no Git mutation. | Decisions 4, 5, 9 |
| `14-performance` | `studies/go-cli-study/reports/final/14-performance.md`; `pkg/cmdutil/factory.go:27-42`, `age/internal/stream/stream.go:20`, `gdu/pkg/analyze/parallel.go:13` cited by handbook | Defer expensive work, stream/bound data, and cap concurrency. | Shapes status without runtime calls, no target scans, and conservative execution. | Decisions 6, 8, 9 |
| `15-philosophy` | `studies/go-cli-study/reports/final/15-philosophy.md`; `VISION.md:97`, `pkg/cmd/factory/default.go:26-46`, `internal/chezmoi/system.go:25-45` cited by handbook | Coherent CLIs keep non-goals explicit and accept complexity deliberately. | Supports rejection of workflow engine, smoke/review/issues/Git mutation, and premature abstractions. | Decisions 1, 9 |
| Architecture reasoning | `projects/ultraplan-go/sprints/23-execute-stage/reasoning/architecture.md` | Chose sprint-owned execute with generic runtime boundary and rejected shared workflow abstractions. | Converts evidence into sprint-specific architectural conclusions. | All decisions |

## Trade-Off And Debt Analysis

Analyze the consequences of the chosen direction before listing final decisions.

### Accepted Trade-Offs

| Trade-Off | Benefit | Cost / Constraint Accepted | Why Acceptable Now | Revisit Trigger |
| --- | --- | --- | --- | --- |
| Sprint-local execute orchestration instead of shared scheduler | Preserves product ownership and avoids premature abstraction. | Duplicates some task-state code shape from study workflows. | TRD and architecture require reuse of infrastructure, not study semantics. | A second non-study execute-like workflow proves identical mechanics and change reasons. |
| Runtime success plus evidence gate | Prevents false completion and aligns with product principles. | More validation rules, diagnostics, and failed states. | Requirements explicitly state runtime success is insufficient. | Evidence rules become too manual or block valid tasks frequently. |
| `.run-state.json` as source of truth and `execute.md` as summary | Enables reliable resume/status while keeping reviewable Markdown. | Markdown summary can become stale between writes. | Status can read state directly and summary can be refreshed after execute transitions. | Users need authoritative offline review from `execute.md` alone. |
| Serial default execution | Reduces conflicting agent edits in one target repository and simplifies diagnostics. | Lower throughput for many tasks. | Sprint scope favors safety, resumability, and target containment over maximum speed. | Task conflicts are controlled by explicit file scopes or future parallel safety policy. |
| Strict target path checks | Prevents accidental workspace/host writes and enforces sprint scope. | Some legitimate local layouts require future configuration. | Requirements name `../ultraplan-go` as the only allowed target for this sprint. | A later sprint adds explicit safe target configuration. |
| Persisting safe runtime metadata summaries only | Reduces leakage risk and state bloat. | Some debug detail may require separate runtime logs or gated debug retention. | Security and observability contracts prioritize safe diagnostics. | Operators need durable deep runtime inspection beyond safe summaries. |

### Potential Technical Debt

| Debt / Shortcut | Why It Might Accrue | Current Mitigation | Owner / Follow-Up |
| --- | --- | --- | --- |
| Duplicate state-transition mechanics between study and sprint | Both domains have task states, attempts, validation, and persistence. | Keep sprint semantics local; extract only mechanical atomic write helpers if already product-neutral. | Future architecture decision after another concrete reuse case. |
| Plan task syntax parser may be narrow | Sprint 23 must extract executable tasks deterministically from current `plan.md` structure. | Reject unsupported or ambiguous tasks with diagnostics instead of guessing. | Plan/execute refinement sprint if task syntax expands. |
| `execute.md` can lag behind `.run-state.json` | Summary Markdown is not the durable state source. | Status reads `.run-state.json`; execute refreshes summary after command completion and meaningful terminal transitions. | Documentation/status refinement if users rely on summary during live runs. |
| Run metadata summaries may omit useful native details | Raw payloads are unsafe by default. | Store run/session IDs, safe categories, validation summaries, and omission reasons. | Observability hardening sprint if debug retention is needed. |
| No general lock decision for execute state | Concurrent execute commands could race on `.run-state.json`. | Atomic writes and conservative state validation; tests should expose race-sensitive paths. | Add explicit sprint execute lock if concurrent mutation becomes likely or race tests reveal risk. |

### Future Considerations

| Consideration | Deferred Until | Reason Deferred | What Should Be Preserved Now |
| --- | --- | --- | --- |
| Parallel execute tasks | Later sprint with conflict controls | Multiple coding-agent tasks can edit the same repository unpredictably. | Task IDs, target path, expected evidence, and status fields that can support concurrency later. |
| Configurable target repositories | Later target-safety sprint | Current requirements permit only explicit `../ultraplan-go`. | Target resolver seam and diagnostics that name approved target source. |
| Durable agentwrap run store | Observability/release hardening | Sprint needs safe summaries, not full run inspection persistence. | Store run/session IDs and safe metadata summaries. |
| Stable public JSON schema | Release hardening | Current sprint needs deterministic command output, but public compatibility is broader. | Keep JSON structures simple, named, versionable, and covered by tests. |
| Execute locks | Reliability hardening if concurrent execute is used | Not explicitly required, but may become important for state races. | Keep run-state store boundary narrow enough to add lock acquisition. |
| Smoke/review/issue/Git workflows | Future explicit roadmap scope | Requirements and selected context exclude them. | Maintain explicit rejected-stage handling and non-goal regression tests. |

## Decisions

This section records the final sprint decisions for execute-stage implementation. These decisions are binding for `plan.md` and must not be reopened during implementation unless a later sprint explicitly changes the architecture or requirements.

### Decision 1: Sprint Owns Execute Semantics

- **Decision:** Implement execute-stage product behavior in `internal/sprint`: plan task extraction, deterministic task IDs, execute prompt rendering, task validation, task state transitions, `.run-state.json`, `execute.md`, execute status, and flow integration. App/CLI code remains thin wiring; `internal/platform/runtime` remains generic execution infrastructure.
- **Rationale:** Execute transforms sprint planning artifacts and sprint task state. The architecture assigns these responsibilities to `internal/sprint`, while the runtime boundary must not understand sprint stages or plan semantics.
- **Study / Source Grounding:** `technical-handbook.md` relevant patterns for thin command boundaries and module-local ownership; `01-project-structure`, `02-command-architecture`, and `03-dependency-injection` reports; `ARCHITECTURE.md` module ownership and runtime boundary rules; `reasoning/architecture.md` chosen Option A.
- **Trade-Offs Accepted:** Accept sprint-local task-state mechanics instead of a shared scheduler. This duplicates code shape but protects module boundaries.
- **Technical Debt / Future Impact:** Future extraction is allowed only for product-neutral mechanical helpers such as atomic JSON writes, not sprint/study task semantics.
- **Alternatives Rejected:** Global workflow engine rejected because cross-sprint scheduling and generic workflow behavior are non-goals. Reusing study run-loop services rejected because TRD says to reuse infrastructure, not study semantics. CLI-owned execute logic rejected because it would become a monolithic command handler.
- **Contracts Satisfied:** Architecture, CLI Surface, Workflows, Testing, REQ-23-42, REQ-23-43.
- **Evidence Required:** Import review showing allowed dependency direction; unit tests for task extraction and IDs; command tests showing CLI delegates behavior; architecture review protocol checks.

### Decision 2: Execute Only From Validated Plan Tasks

- **Decision:** Execute validates prerequisites through `plan.md` before task extraction or runtime launch. It reads only validated `plan.md` entries, rejects unsupported/ambiguous task syntax, assigns deterministic task IDs from stable normalized plan task fields, and stores traceability back to the plan task and reasoning decision/requirement where available.
- **Rationale:** The execute stage must be governed. Running ad hoc prompts or inferred work would bypass the planning contract and make resume/status unreliable.
- **Study / Source Grounding:** `technical-handbook.md` design pressures for deterministic IDs and validation; `TRD.md` sections 18.5.1, 18.6, and 18.7; `11-testing-strategy` fixture/fake guidance.
- **Trade-Offs Accepted:** A strict parser may reject human-readable but under-specified tasks. This is acceptable because unsafe execution is worse than a repairable plan validation error.
- **Technical Debt / Future Impact:** The task syntax may need expansion later; preserving clear diagnostics and plan fingerprints keeps future migration possible.
- **Alternatives Rejected:** Inferring executable tasks from arbitrary Markdown rejected because it is non-deterministic. Running an operator-supplied free-form execute prompt rejected because it breaks traceability. Auto-repairing `plan.md` during execute rejected because execute should not reopen planning architecture.
- **Contracts Satisfied:** Workflows, Errors, Persistence And Migrations, Testing, REQ-23-40, REQ-23-42, REQ-23-45.
- **Evidence Required:** Fixture tests for valid tasks, duplicate/ambiguous tasks, unsupported syntax, stable IDs across identical validated content, and traceability fields in run state.

### Decision 3: Versioned Run State Is The Durable Source Of Truth

- **Decision:** Create `projects/ultraplan-go/sprints/23-execute-stage/.run-state.json` when execute begins or resumes. It is versioned, strictly loaded, written atomically in the same directory, and stores project, sprint, target path, plan path/fingerprint, deterministic tasks, statuses `pending`, `running`, `complete`, `failed`, `cancelled`, attempts, timestamps, safe diagnostics, and runtime metadata summaries where available.
- **Rationale:** Execute must be resumable after interruption and inspectable without runtime calls. Durable state, not Markdown or runtime memory, must govern resume and status.
- **Study / Source Grounding:** `technical-handbook.md` durable state pattern; `07-state-context`, `10-logging-observability`, and `06-io-abstraction` reports; `TRD.md` sections 18.5.1, 20, and 21; `reasoning/architecture.md` state decision.
- **Trade-Offs Accepted:** Writing state after meaningful transitions increases file writes and schema care. This is acceptable because recovery and interruption safety are core requirements.
- **Technical Debt / Future Impact:** First schema has no migration path beyond strict rejection of unsupported versions. Future migrations can be added once additional versions exist.
- **Alternatives Rejected:** Using `execute.md` as state rejected because Markdown is reviewable summary, not reliable machine state. In-memory-only task state rejected because it is not resumable. Leniently accepting malformed state rejected because it risks corrupt resume decisions.
- **Contracts Satisfied:** Persistence And Migrations, Observability, Errors, Workflows, REQ-23-46, REQ-23-47, REQ-23-48.
- **Evidence Required:** State tests for schema versioning, malformed JSON, unsupported versions, atomic write behavior, stale running recovery, resume without redoing complete tasks, cancellation, and status counts from state.

### Decision 4: Completion Requires Evidence Or Explicit Safe Diagnostic

- **Decision:** A task reaches `complete` only when it was extracted from valid `plan.md`, runtime execution succeeds when invoked, expected evidence is present, or an explicit safe diagnostic explains why the expected evidence cannot be machine-validated. Runtime success without evidence marks the task failed, not complete.
- **Rationale:** Product principles and sprint requirements state that runtime success is not product success. Execute must avoid false positives that make later review or resume unsafe.
- **Study / Source Grounding:** `technical-handbook.md` validation-gated completion pattern; `05-error-handling` and `11-testing-strategy` reports; `PRD.md` runtime-success principle; `TRD.md` section 18.7.
- **Trade-Offs Accepted:** Some valid implementation work may require a diagnostic-only completion path. That path must be explicit, safe, and reviewable.
- **Technical Debt / Future Impact:** Evidence classification may start simple and need refinement as plan tasks diversify. Diagnostics should preserve enough context for future stricter validators.
- **Alternatives Rejected:** Runtime-exit-only completion rejected by requirements. Manual operator marking without diagnostic rejected because it lacks review evidence. Treating missing evidence as warning rejected because it undermines status truthfulness.
- **Contracts Satisfied:** LLM Evaluation / Cost / Safety, Errors, Testing, Workflows, REQ-23-45, REQ-23-49.
- **Evidence Required:** Fake-runtime tests for success with evidence, success without evidence, runtime failure, validation failure, and diagnostic-only completion; review checks that `execute.md` records completion evidence or diagnostics.

### Decision 5: Target Repository And Runtime Diagnostics Are Safety-Gated

- **Decision:** Execute targets only the explicit implementation repository `../ultraplan-go` for this sprint. Target paths are resolved and validated before runtime launch, prompts and permission policy constrain work to that target, unsafe path escapes fail, automatic Git mutation remains excluded, and persisted diagnostics redact secrets and omit unsafe raw runtime payloads by default.
- **Rationale:** Execute allows runtime-backed writes in a sibling repository, so path safety, permission posture, and redaction are central safety requirements.
- **Study / Source Grounding:** `technical-handbook.md` workspace-safe execution and redaction warnings; `13-security` report; `TRD.md` sections 18.5.1, 19, 20, and 22; `requirements.md` target path and no-Git constraints.
- **Trade-Offs Accepted:** Strict path validation can block alternate local layouts. This is acceptable because no future target configuration has been selected.
- **Technical Debt / Future Impact:** Target configuration should be introduced later only with explicit safe configuration, migration, and tests.
- **Alternatives Rejected:** Allowing arbitrary target paths rejected due to escape risk. Letting runtime decide workdir scope rejected because product owns target safety. Adding Git add/commit/push hooks rejected because requirements explicitly exclude Git mutation.
- **Contracts Satisfied:** Security, LLM Runtime, Observability, Errors, REQ-23-43, REQ-23-44, REQ-23-51.
- **Evidence Required:** Target containment tests, path escape rejection tests, redaction tests, permission/request mapping review, and regression checks that no smoke/review/issues/Git mutation behavior is added.

### Decision 6: Execute Runs Conservatively And Resumes Safely

- **Decision:** Execute runs one task at a time by default for this sprint, with context cancellation, no unbounded goroutines, state saved after meaningful transitions, stale `running` tasks recovered to retryable or failed state with diagnostics, and completed tasks skipped on resume unless a future explicit force behavior is designed.
- **Rationale:** Multiple coding-agent tasks in one target repository can conflict. Serial default execution improves safety, simplifies evidence gates, and still satisfies execute requirements.
- **Study / Source Grounding:** `technical-handbook.md` bounded concurrency and state patterns; `07-state-context`, `08-concurrency`, and `14-performance` reports; `TRD.md` sections 20 and 18.5.1; architecture reasoning assumptions.
- **Trade-Offs Accepted:** Lower throughput is accepted in exchange for safer target edits and clearer diagnostics.
- **Technical Debt / Future Impact:** Parallelism can be added later if task scopes, locks, or conflict controls are designed. Current state schema should not preclude future active-task counts.
- **Alternatives Rejected:** Unbounded parallel runtime tasks rejected for resource and edit-conflict risk. Bounded parallelism in this sprint rejected because no conflict policy is selected. Redoing complete tasks on resume rejected because it violates resumability and can create duplicate edits.
- **Contracts Satisfied:** Workflows, Performance, Persistence And Migrations, Testing, REQ-23-46, REQ-23-47, REQ-23-50.
- **Evidence Required:** Workflow tests for resume, stale recovery, cancellation, no redo of complete tasks, no goroutine leaks in fake-runtime tests, and normal `go test -race ./...`.

### Decision 7: Stage Model Resolution Is Deterministic And Visible

- **Decision:** Execute uses the sprint stage model resolution order: explicit command override when supported, configured `execute` stage model, configured global sprint/planning model, `models.primary`, then `models.default`. Unknown stage keys and empty model values fail validation. Diagnostics and prompt previews show the selected model source without secrets.
- **Rationale:** Execute and planning stages need predictable runtime configuration, and stage-specific execute overrides must win over global/default values.
- **Study / Source Grounding:** `technical-handbook.md` config precedence pattern; `04-configuration-management` report; `TRD.md` sections 6 and 18.7.
- **Trade-Offs Accepted:** Source-aware resolution adds implementation detail but prevents ambiguous runtime behavior.
- **Technical Debt / Future Impact:** If future stages are added, the supported key list must be updated intentionally and covered by tests.
- **Alternatives Rejected:** Ad hoc lookup inside execute rejected because it would diverge from planning stages. Silent fallback from unknown keys rejected because it hides config mistakes. Printing full runtime config rejected because of secret leakage risk.
- **Contracts Satisfied:** Configuration, LLM Runtime, Observability, Security, REQ-23-41.
- **Evidence Required:** Unit tests for precedence, stage-specific override, fallback, unknown keys, empty values, command override, and redacted diagnostics.

### Decision 8: Status, Prompt, Validate, Flow, And Execute Commands Stay Thin

- **Decision:** CLI/app wiring exposes `validate execute`, `prompt execute`, `flow --to execute`, `execute`, task selection, dry-run/prompt preview, help, status, and exit-code mapping by calling sprint/app use cases. Status reads `flow-state.json`, `.run-state.json`, `execute.md`, and prerequisite validation without calling runtime.
- **Rationale:** Command handlers should remain scriptable, testable adapters. Execute workflow logic belongs in sprint services and app use cases.
- **Study / Source Grounding:** `technical-handbook.md` command and IO patterns; `02-command-architecture`, `06-io-abstraction`, `09-terminal-ux`, and `14-performance` reports; `TRD.md` command requirements.
- **Trade-Offs Accepted:** More request/result structs and output adapters are needed. This is acceptable for command testability and future TUI reuse.
- **Technical Debt / Future Impact:** JSON status schema may need later stabilization; keep output deterministic and covered by tests now.
- **Alternatives Rejected:** Monolithic `RunE` execute implementation rejected by handbook warnings. Runtime-backed status rejected because status must be inspectable offline. TUI behavior rejected because not selected for this sprint.
- **Contracts Satisfied:** CLI Surface, Observability, Documentation, Errors, Testing, REQ-23-40, REQ-23-48.
- **Evidence Required:** Command tests for help, prompt preview, dry-run, status without runtime, invalid prerequisites, invalid stages, task selection, and exit codes.

### Decision 9: Deferred Behaviors Are Explicitly Rejected

- **Decision:** Execute must not introduce smoke investigation, `smoke.md`/`smoke.json`, automated `review.md`, issue files/JSON, assignment/scheduling, project-management workflow, Git add/commit/push/branch/PR/checkout/reset, hosted/browser UI, TUI, or cross-sprint scheduling.
- **Rationale:** The sprint index, requirements, PRD, TRD, and project index all exclude these behaviors. Keeping them out protects scope and trust boundaries.
- **Study / Source Grounding:** `technical-handbook.md` anti-pattern warning; `15-philosophy` report on explicit non-goals; `13-security` report on trust boundaries; selected sprint index excluded context.
- **Trade-Offs Accepted:** Execute is narrower than the prototype's broader workflow suite. This is acceptable because Phase 2 scope is controlled execution only.
- **Technical Debt / Future Impact:** Future workflows must be added by explicit requirements and reasoning, not incidental execute hooks.
- **Alternatives Rejected:** Porting prototype smoke/review/issues now rejected as out of scope. Adding opt-in Git mutation rejected because no safe contract is selected. Building generic cross-sprint execution rejected because it is a non-goal.
- **Contracts Satisfied:** Security, Workflows, Architecture, CLI Surface, REQ-23-51.
- **Evidence Required:** Regression tests or code review checks for rejected stages/commands/artifacts and absence of Git mutation behavior.

### Decision 10: Verification Uses Deterministic Fake-Runtime Evidence

- **Decision:** Normal verification relies on deterministic offline unit, fixture, command, fake-runtime, golden, race, and build checks. Real OpenCode/runtime smoke is optional and gated, not required for sprint acceptance.
- **Rationale:** Agentic runtime behavior is nondeterministic and may require credentials/network. Sprint acceptance requires reliable local verification with fake runtimes.
- **Study / Source Grounding:** `technical-handbook.md` testing pattern; `11-testing-strategy` report; `TRD.md` testing sections 23 and 26; requirements verification gates.
- **Trade-Offs Accepted:** Fake-runtime tests do not prove provider behavior. They do prove UltraPlan's product-owned logic, state, validation, prompts, and command behavior.
- **Technical Debt / Future Impact:** Gated real-runtime smoke can be added or run later through external harnesses when selected.
- **Alternatives Rejected:** Requiring OpenCode/provider credentials for normal tests rejected. Manual-only verification rejected because durable state and validators are core. Skipping race/build checks rejected by requirements.
- **Contracts Satisfied:** Testing, LLM Runtime, Workflows, Performance, REQ-23-50.
- **Evidence Required:** From `../ultraplan-go`: `go test ./...`, `go test -race ./...`, `go build ./cmd/ultraplan`; plus focused fake-runtime and command tests.

## Expected Evidence

| Evidence Type | Required Evidence | Source / Command / Review Check |
| --- | --- | --- |
| Tests | Unit tests for task extraction, deterministic IDs, plan prerequisite validation, target path safety, model resolution, run-state schema/load/save, stale recovery, task transitions, evidence gates, summary rendering, and non-goals. | `../ultraplan-go/internal/sprint/execute_test.go` and related sprint tests. |
| Tests | Command tests for `validate execute`, `prompt execute`, `flow --to execute`, `execute`, `--task`, dry-run/prompt behavior, invalid prerequisites, invalid stages, exit codes, text/status output, and runtime-free status. | `../ultraplan-go/internal/app/sprint_execute_commands_test.go` or equivalent app/command test files. |
| Tests | Fake-runtime tests for success with evidence, success without evidence, runtime failure, cancellation, timeout, validation failure, retry/fallback metadata where applicable, and safe diagnostics. | Agentwrap-compatible fake runtime fixtures in `../ultraplan-go` tests. |
| Tests | Regression coverage that execute does not invoke smoke, review automation, issue tracking, Git mutation, TUI behavior, hosted/browser behavior, or study services. | Unit/import/code review checks and command tests. |
| Verification | Full normal verification passes. | Run from `../ultraplan-go`: `go test ./...`, `go test -race ./...`, `go build ./cmd/ultraplan`. |
| Runtime | Runtime requests use generic platform/runtime and agentwrap metadata only; no direct OpenCode stdout/stderr parsing. | Import review, fake-runtime request assertions, architecture review. |
| Persistence | `.run-state.json` and `flow-state.json` are versioned, strictly loaded, atomically written, resumable, and recover stale running tasks with diagnostics. | State-store tests and fixture state files. |
| Runtime | `execute.md` cites `plan.md`, summarizes task counts/terminal states, records completion evidence, and records failed-task diagnostics or `.run-state.json` pointers. | Execute summary golden tests and validator tests. |
| Review | Architecture and sprint review protocols are satisfied. | `system/protocols/architecture-review-protocol.md` and `system/protocols/sprint-review-protocol.md`. |
| Documentation | Help/status output and artifact text describe supported execute behavior and exclusions. | CLI help tests, status golden tests, and docs/help review. |

## Assumptions And Risks

| Item | Type | Impact | Mitigation / Follow-Up |
| --- | --- | --- | --- |
| Valid `plan.md` task structure is explicit enough to parse deterministically. | Assumption | If false, execute cannot safely assign stable IDs. | Reject unsupported/ambiguous tasks before runtime and require plan repair. |
| Existing platform/runtime can support fake-runtime injection through app/sprint seams. | Assumption | If false, tests may require refactoring composition before execute. | Keep runtime dependency injected at app/sprint boundary; do not construct runtime in helpers. |
| Target repository is `../ultraplan-go`. | Assumption | Alternate layouts are rejected this sprint. | Validate target source and path; defer configurable targets. |
| Runtime-created edits can conflict across tasks. | Risk | Parallel execution could corrupt or confuse target changes. | Run serially by default and defer parallelism until conflict controls exist. |
| `.run-state.json` can become inconsistent with target repo if runtime changes files but evidence is weak. | Risk | Status could overstate progress. | Require evidence or explicit diagnostic; persist failed/missing evidence diagnostics. |
| Diagnostics may leak secrets or unsafe native payloads. | Risk | Security and review failure. | Redact secrets, omit unsafe raw payloads, store safe summaries and omission facts. |
| Concurrent execute invocations may race on state. | Risk | Corrupt or stale task status. | Atomic writes and strict reload validation now; add explicit locks if implementation/review exposes risk. |
| `execute.md` may be stale relative to run state. | Risk | Human readers could inspect outdated summary. | Treat `.run-state.json` as source of truth; status reads state; refresh summary after execute transitions. |
| Gated real runtime may behave differently than fake runtime. | Risk | Integration issues may surface later. | Keep fake-runtime as acceptance; optional gated OpenCode checks can be run separately when selected. |

## Implementation Constraints

- Keep execute semantics in `internal/sprint`; do not move product behavior into `internal/platform/runtime`, `internal/study`, or a global workflow/scheduler package.
- Keep `internal/platform/runtime` generic: prompts, workdir, model, timeout, permissions, expected outputs, events, and results only; no project/sprint/plan/task semantics.
- Use agentwrap/OpenCode through the existing runtime integration; do not parse OpenCode stdout/stderr directly in sprint/app product code.
- Validate prerequisites through `plan.md` before execute task extraction or runtime launch.
- Extract executable tasks only from validated `plan.md`; assign deterministic IDs and preserve plan/reasoning traceability.
- Constrain target writes to `../ultraplan-go` unless a future explicit safe configuration is selected.
- Treat runtime success as necessary but insufficient; require evidence or explicit safe diagnostic before `complete`.
- Persist `.run-state.json` as versioned JSON with strict load validation and same-directory atomic writes.
- Recover stale `running` tasks conservatively with diagnostics; do not redo `complete` tasks on resume unless future force behavior is explicitly designed.
- Keep `execute.md` as editable Markdown summary, not source of truth.
- Keep status runtime-free and derived from prerequisites, `flow-state.json`, `.run-state.json`, and `execute.md` where relevant.
- Default execute task execution to serial for target-repository safety; do not add unbounded goroutines or runtime fan-out.
- Redact secrets and omit unsafe raw runtime payloads from logs, diagnostics, status, `execute.md`, and JSON/text output.
- Do not introduce smoke, review automation, issue tracking, Git mutation, TUI behavior, hosted/browser UI, or cross-sprint scheduling.
- Normal tests must be deterministic, offline, fake-runtime based, and must not require OpenCode, provider credentials, network access, or Git mutation.

## Plan Handoff

`plan.md` must execute these decisions. It must not invent architecture, scope, or decisions beyond this document.

The plan must carry forward:

- final decisions
- applicable contracts and requirement IDs
- expected evidence
- risks and assumptions
- required review protocols

## Phase Exit Criteria

- [x] Selected context was read and used.
- [x] Area-specific reasoning documents were completed where applicable or explicitly marked not applicable.
- [x] Area-specific reasoning conclusions are reflected or explicitly overridden in final decisions.
- [x] Contracts are explicitly mapped to decisions and expected evidence.
- [x] Final decisions are clear enough for `plan.md` to execute.
- [x] Expected evidence is specific and reviewable.
