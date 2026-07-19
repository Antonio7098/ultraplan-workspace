> **Inputs Used:** `projects/ultraplan-go/sprints/23-execute-stage/sprint-index.md`, `projects/ultraplan-go/sprints/23-execute-stage/requirements.md`, `projects/ultraplan-go/sprints/23-execute-stage/technical-handbook.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `system/protocols/architecture-review-protocol.md`, `system/protocols/sprint-review-protocol.md`, `system/reasoning/architecture_reasoning_template.md`

# Architecture Reasoning: Execute Stage

## Area Decisions

- Execute-stage product behavior belongs in `internal/sprint`, including plan task extraction, execute prompt construction, `.run-state.json` persistence, `execute.md` rendering, validation, flow-state updates, and sprint status integration.
- Runtime execution stays behind `internal/platform/runtime` and agentwrap. Sprint code must build generic runtime requests and must not parse provider-specific output or call OpenCode directly.
- Workspace-safe target resolution belongs at the workspace/config boundary. Execute must constrain implementation edits to the approved target repository and must not mutate Git state automatically.
- `.run-state.json` is the source of truth for task progress; `execute.md` is a reviewable summary that can be regenerated from durable task state.

## Trade-Offs

- Keeping execute orchestration local to `internal/sprint` avoids a premature global scheduler, but it requires focused file organization and tests so sprint logic does not become opaque.
- Reusing only generic runtime infrastructure duplicates some mechanical task-state patterns from study workflows, but preserves product boundaries and avoids coupling sprint execution to study semantics.
- Treating `execute.md` as a summary rather than state accepts possible staleness in Markdown while keeping status and resume behavior grounded in `.run-state.json`.
- Sequential or bounded task execution is less aggressive than unbounded fan-out, but it keeps diagnostics, retries, cost, and target-repository changes easier to reason about.

## Evidence

- `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/TRD.md`, and `projects/ultraplan-go/sprints/23-execute-stage/technical-handbook.md` all support sprint-owned product behavior with product-agnostic platform/runtime boundaries.
- `system/protocols/architecture-review-protocol.md` and `system/protocols/sprint-review-protocol.md` provide the review checks for dependency direction, thin CLI wiring, state ownership, and no forbidden workflow abstractions.
- The selected Architecture reasoning template, `system/reasoning/architecture_reasoning_template.md`, drove the deeper option analysis, abstraction check, state model, and implementation constraints below.

## Risks

- `internal/sprint` can become too broad unless execute behavior is split into focused helpers for task extraction, run state, prompts, validation, summary rendering, and status.
- Runtime/provider diagnostics can leak unsafe detail unless errors are categorized and redacted before persistence or CLI display.
- Resume semantics can become inconsistent if task IDs are not deterministic from `plan.md` or if unsupported `.run-state.json` schema versions are accepted.
- Markdown summaries can drift from durable state if status reads `execute.md`; status should compute from `.run-state.json`.

## 0. Feature Summary

### Feature name

`governed sprint execute stage`

### User/product goal

Execute implementation tasks from a validated `plan.md` through the generic runtime boundary, with durable task state, resumability, safe target-repository boundaries, status visibility, and actionable diagnostics.

### Current task

Plan the architecture for adding execute-stage behavior to UltraPlan Go. The implementation must keep execute semantics in `internal/sprint`, route agentic execution through `internal/platform/runtime` and agentwrap, persist `.run-state.json`, render `execute.md`, and expose progress through sprint status and flow commands.

### Non-goals

- Do not implement smoke investigation.
- Do not implement automated review or conformance review execution.
- Do not implement issue tracking, assignment, or project-management workflows.
- Do not mutate Git state automatically with add, commit, push, branch, or PR operations.
- Do not create a cross-sprint scheduler or general workflow engine.
- Do not move sprint task semantics into `internal/platform/runtime`, `internal/study`, or a global scheduler package.

## 1. First-Principles Breakdown

### Core behaviour

The execute stage reads a validated sprint `plan.md`, extracts deterministic implementation tasks, runs pending tasks through a generic runtime request, records task transitions in durable run state, validates completion using expected evidence or explicit diagnostics, writes an execute summary, and reflects progress in sprint flow/status.

### Inputs

- Valid sprint prerequisites through `plan.md`.
- `projects/ultraplan-go/sprints/23-execute-stage/plan.md` entries that are explicit enough to produce deterministic task IDs.
- Sprint, project, and workspace paths resolved through workspace-safe rules.
- Runtime configuration with global/default model and `execute` stage override precedence.
- Target implementation directory configured or resolved as an explicit workspace-safe path.
- Runtime execution result, events, metadata, validation result, and safe diagnostics from `platform/runtime` and agentwrap.

### Outputs

- `projects/ultraplan-go/sprints/23-execute-stage/.run-state.json` with schema version, task records, task statuses, attempts, diagnostics, timestamps, target path, and runtime metadata summaries where available.
- `projects/ultraplan-go/sprints/23-execute-stage/execute.md` summarizing task counts, terminal states, and failed-task diagnostics or pointers to `.run-state.json`.
- Updated `flow-state.json` showing execute stage progress once valid prerequisites through `plan` exist.
- Sprint status output showing execute progress.
- Runtime-created implementation changes inside the approved target repository boundary.

### Durable state

- Created: `.run-state.json` when execute begins, `execute.md` when a summary is produced.
- Read: `requirements.md`, `sprint-index.md`, `technical-handbook.md`, selected `reasoning/*.md`, `reasoning.md`, `plan.md`, `flow-state.json`, workspace/project config, and existing `.run-state.json` on resume.
- Updated: `.run-state.json` after every meaningful task transition, `flow-state.json` for execute progress, `execute.md` summary.
- Deleted: none as part of execute. Existing complete state is not deleted unless a later explicit force behavior is designed.

### Ephemeral state

- In-memory extracted task list and plan fingerprint.
- Active runtime run handles and event drains.
- Resolved runtime model, timeout, variant, permission policy, and target workdir.
- Current task attempt context, cancellation context, and transient validation result.
- Progress counters used for status rendering.

### Derived state

- Execute progress counts.
- Source of truth: `.run-state.json` task records and terminal statuses.
- Stored or computed: counts may be stored in `execute.md` for review but should be computed from `.run-state.json` for status.
- Staleness risk: medium for `execute.md` if task state changes after it is written; low for status if it reads current state.

### Side effects

- Database write: no.
- File write: yes, atomic writes for `.run-state.json` and `flow-state.json`, editable Markdown write for `execute.md`, and runtime-produced implementation edits constrained to the target directory.
- Network call: yes when the configured runtime/provider performs model execution through agentwrap/OpenCode.
- Queue/event emission: no external queue; runtime events are consumed and mapped to local diagnostics/logs.
- Email/notification: no.
- Logs/metrics/traces: yes, structured logs/status fields should include project, sprint, task ID, attempt, runtime/model, duration, outcome, and safe error category.

## 2. Existing Architecture Fit

### Contract applicability and phase gate

Current gate: `runtime/provider/batch workflow`

Applies now:
- Architecture: execute must preserve `sprint -> project`, `sprint -> workspace`, and `sprint -> platform/runtime`, with platform packages importing no product modules.
- CLI Surface: `flow --to execute`, `execute`, `status`, validation, prompt rendering, and scriptable outputs are in scope.
- Configuration: execute model resolution must support global/default plus stage-specific override, with stage-specific values winning.
- Documentation: `.run-state.json` is machine state, while `execute.md` remains editable, reviewable Markdown.
- Errors: prerequisite failures, malformed state, runtime failure, cancellation, validation failure, and missing evidence need actionable diagnostics.
- LLM Evaluation / Cost / Safety: runtime metadata and cost/usage summaries are recorded when available and safe.
- LLM Runtime: execution must use the generic runtime/agentwrap boundary and must not parse OpenCode directly from sprint code.
- Observability: task metadata, run IDs, attempt counts, diagnostics, and status visibility are in scope.
- Performance: task execution must avoid unbounded fan-out, expensive scans, and repeated large artifact reads.
- Persistence And Migrations: `.run-state.json` must be versioned, atomic, resumable, and reject unsupported schema versions.
- Security: target path safety, secret redaction, runtime permission policy, source isolation, and no automatic Git mutation are in scope.
- Testing: fake runtime, fixture state files, validation, prompt, and command/status tests are required.
- Workflows: execute joins the sprint flow only after valid prerequisites through `plan` and must support retries, cancellation, and resumability.

Deferred:
- Smoke/review/issue workflow contracts: deferred until a later roadmap gate explicitly adds those stages.
- Stable public JSON compatibility: deferred until release hardening defines the minimum stable JSON schema, although current structures should be deterministic and testable.
- Cross-sprint scheduling: deferred because requirements exclude a general workflow engine.

### Where does this behaviour naturally belong?

Execute-stage product behavior belongs in `internal/sprint` because that module owns sprint artifacts, stage validation, flow state, prompt rendering, plan task extraction, deterministic task IDs, run-state persistence, execute status, and the `execute` lifecycle. Runtime adapter mechanics belong in `internal/platform/runtime` and agentwrap. Workspace-safe path resolution belongs in `internal/workspace`. Composition, command parsing, and output mapping belong in `internal/app` and CLI glue.

### Existing workflow affected

Current flow:
1. Sprint commands discover workspace, project, and sprint.
2. Sprint flow validates planning prerequisites through selected stages.
3. Planning stages render prompts, invoke runtime when needed, validate expected artifacts, and update `flow-state.json` through `plan`.

### Proposed new flow

Proposed flow:
1. CLI/app resolves workspace, project, sprint, config, runtime, and target path.
2. `internal/sprint` validates prerequisites through `plan` and rejects execute if `plan.md` is missing, invalid, or invokes excluded behavior.
3. `internal/sprint` extracts deterministic executable tasks from validated `plan.md` and initializes or resumes `.run-state.json`.
4. `internal/sprint` selects pending or requested tasks and builds execute prompts plus generic runtime requests.
5. `internal/platform/runtime` and agentwrap execute each task with configured model, permissions, workdir, timeout, observability, validation, retry, and cancellation behavior.
6. `internal/sprint` records transitions, attempts, safe diagnostics, validation/evidence results, and terminal statuses atomically.
7. `internal/sprint` writes or refreshes `execute.md` and updates `flow-state.json`/status from the current run state.

### Does the current architecture still fit?

```text
[x] Yes - the feature fits cleanly into the existing shape.
[ ] Partially - small refactor needed before adding the feature.
[ ] No - the feature reveals that the existing architecture is the wrong shape.
```

### If it does not fit, why?

```text
- [ ] Requires awkward boolean flags
- [ ] Requires caller-specific branching
- [ ] Duplicates business rules
- [ ] Forces unrelated responsibilities into one unit
- [ ] Requires reaching into internals
- [ ] Requires global/shared mutable state
- [ ] Makes tests much harder
- [ ] Hides important side effects
- [ ] Other: not applicable
```

### Refactor-before-feature decision

Decision: implement directly with focused sprint-owned files and existing platform/runtime boundary.

Reason: the project architecture already names `internal/sprint` as owner of execute prompts, task extraction, `.run-state.json`, stage validation, flow state, and status. The technical handbook evidence supports thin command wrappers, module-local ownership, explicit composition roots, runtime boundaries, atomic state, and fakeable runtime seams. No large cross-module refactor is needed before adding execute.

## 3. Design Options Considered

### Option A: Minimal direct implementation

Description:
Add focused execute-stage behavior inside `internal/sprint` using concrete helpers for plan task extraction, prompt rendering, run-state load/save, task transitions, evidence checks, and execute summary rendering. Use existing app composition to inject workspace/config/runtime/output dependencies and route runtime calls through the generic platform/runtime boundary.

Pros:
- Preserves module ownership stated in `ARCHITECTURE.md` and TRD section 18.1.
- Keeps runtime generic and product semantics in `sprint`.
- Avoids premature global scheduler/workflow abstraction.
- Enables fake-runtime and fixture tests through existing runtime and filesystem seams.
- Smallest design that satisfies the sprint requirements.

Cons:
- Some task-state mechanics may resemble study run-loop mechanics without immediate shared abstraction.
- `internal/sprint` grows in responsibility and needs focused file organization to remain readable.
- Future cross-module task-state reuse may require later extraction of mechanical atomic-write helpers.

Risks:
- Duplicated state-machine shape could tempt later inconsistent behavior between study and sprint if not tested.
- Execute summary can become stale unless status treats `.run-state.json` as source of truth.

### Option B: Refactor then implement

Description:
Extract shared task scheduling, run-state, validation, prompt, or workflow abstractions from study/planning code before adding sprint execute, then build execute on those abstractions.

Pros:
- Could reduce repeated atomic write or state transition mechanics if the existing implementations are truly identical.
- May create a reusable test harness for long-running workflows.

Cons:
- Conflicts with TRD and architecture guidance to reuse infrastructure, not study semantics.
- Risks creating `internal/scheduler`, `internal/validation`, or `internal/workflow` abstractions before repeated stable product needs exist.
- May force sprint semantics into generic APIs or require boolean flags for study-vs-sprint behavior.
- Delays the sprint goal without evidence that a broad refactor is required.

Risks:
- Unearned abstraction makes future review harder and may hide side effects.
- Extracting study task concepts could couple `sprint` to study run-loop semantics that are explicitly out of bounds.

### Option C: New workflow engine / scheduler component

Description:
Create a new product-neutral workflow engine or scheduler package to own stage execution, task scheduling, run state, retries, and resumability across study and sprint workflows.

Pros:
- Central place for cross-workflow orchestration if UltraPlan later supports many workflow types.
- Could offer uniform state transitions and status structures in a future release.

Cons:
- Directly violates the sprint non-goal of avoiding a cross-sprint scheduler or general workflow engine.
- The current behavior has one concrete sprint execute use case and an existing study run-loop with different product semantics.
- Increases coupling and likely introduces generic names, mode switches, and product leakage.
- Duplicates agentwrap's policy/runtime responsibilities instead of consuming them.

Risks:
- Platform/runtime or scheduler could start understanding sprint stages, which violates architecture dependency rules.
- The workflow engine could become a dumping ground for product behavior.

### Chosen option

Chosen option: A.

Reason:
The simplest honest design is to implement execute as sprint-owned product behavior, reuse only generic infrastructure, and keep runtime mechanics behind agentwrap/platform. The technical handbook cites thin command wrappers, module-local ownership, explicit dependency injection, atomic durable state, safe observability, and bounded concurrency as high-confidence evidence. Option A satisfies those pressures without inventing a global workflow abstraction.

## 4. Abstraction Check

### Are we adding a new abstraction?

```text
[ ] No
[ ] Yes - interface/protocol/trait
[ ] Yes - base class
[x] Yes - service/component/module
[ ] Yes - strategy/function parameter
[x] Yes - data structure/DTO/config object
[ ] Yes - plugin/registry/factory
```

### If yes, why is it earned?

```text
- [ ] Multiple real implementations exist now
- [ ] A second implementation is very likely soon
- [x] It isolates an external dependency
- [ ] It removes branching rather than adding it
- [x] It protects a stable domain concept
- [x] It makes testing simpler
- [x] It creates a clear boundary between policy and mechanism
- [x] Other: execute task/run-state data structures are required durable product state
```

The earned abstractions are sprint-local use-case/service functions plus execute-specific request/result/state structs. Runtime itself should reuse the existing agentwrap/platform boundary rather than a new runtime contract. Filesystem, runtime, clock, and output seams are justified at side-effect boundaries for deterministic tests.

### If no, why are we keeping it concrete?

Internal helpers for parsing plan entries, normalizing task IDs, computing counts, and rendering `execute.md` should stay concrete functions because they have one sprint-local implementation and no current volatility. A separate interface for each helper would be premature.

### Bad abstraction smell check

```text
- [ ] Generic name like Manager/Handler/Processor/Common/Helper
- [ ] Only one implementation with no real volatility
- [ ] Boolean flags or mode switches
- [ ] Optional parameters for caller-specific behaviour
- [ ] Interface mirrors a concrete class rather than consumer needs
- [ ] Abstraction exists only to reduce line count
- [ ] Abstraction makes the flow harder to trace
```

If these smells appear during implementation, the design should be reduced back to sprint-local concrete helpers or split by workflow step.

## 5. DRY and Duplication Check

### Is any duplication being introduced?

```text
[ ] No
[ ] Yes - duplicated business/domain knowledge
[x] Yes - duplicated code shape only
[x] Yes - duplicated infrastructure mechanics
```

### If duplicating business knowledge, stop and consolidate

Shared rule/knowledge:
No study business knowledge should be duplicated. Sprint execute completion rules, plan task extraction, and flow stages are separate sprint product rules.

Single source of truth should be:
`internal/sprint` for sprint execute semantics; `internal/workspace` for path safety; `internal/platform/runtime`/agentwrap for runtime execution; config package/app wiring for configuration precedence.

### If duplicating code shape, is that acceptable?

Reason duplication is acceptable:
Study run-loop state and sprint execute run state both look like task lifecycle state, but they change for different product reasons. The TRD explicitly says to keep study service, study prompt builders, study report validation, rating parsing, summary generation, and study task scheduling local to study. Mechanical duplication is acceptable initially for sprint-specific state transitions and validators.

### If removing duplication, is the abstraction safe?

```text
- [x] The duplicated code represents the same knowledge
- [x] Call sites change for the same reason
- [x] The new abstraction has a meaningful name
- [x] It does not introduce flags or caller-specific branches
- [x] It makes tests simpler
```

Only mechanical helpers that meet all checks should be extracted, such as atomic JSON write/read helpers or small diagnostic/result structs without product semantics. Do not extract a shared scheduler or validator for sprint and study during this sprint.

## 6. Coupling Check

### Dependencies required

- `internal/project`: validates project catalog context and resolves selected sprint/project metadata.
- `internal/workspace`: resolves workspace root, normalizes paths, and enforces target path safety.
- `internal/platform/runtime`: executes generic prompt requests and returns generic results/events/metadata.
- Agentwrap through platform/runtime: supervises OpenCode, policy, validation, observability, permissions, cancellation, and health checks.
- Config/app wiring: resolves runtime, model, timeout, variant, permission policy, and output modes.
- Filesystem abstraction or concrete store at sprint boundary: loads/writes plan, state, flow, and summaries with atomic persistence where required.
- Clock/time seam: records deterministic timestamps in tests.

### Coupling risks

```text
Global coupling:
[x] None
[ ] Introduced
[ ] Existing but unchanged
[ ] Existing and reduced

Content coupling:
[x] None
[ ] Introduced
[ ] Existing but unchanged
[ ] Existing and reduced

Stamp coupling:
[ ] None
[ ] Introduced
[x] Existing but unchanged
[ ] Existing and reduced
```

Stamp coupling risk exists if broad app/config objects are passed deep into sprint helpers. Keep helper parameters narrow, while allowing cohesive request structs at use-case boundaries.

### Narrow dependency check

Does each function receive only what it needs?
```text
[x] Yes
[ ] No - reason: not applicable
```

Are large objects passed only when the full concept is needed?
```text
[x] Yes
[ ] No - reason: not applicable
```

### Dependency injection check

Are side-effectful dependencies created at the edge and passed inward?
```text
[x] Yes
[ ] No
[ ] Not applicable
```

If no, why is internal construction acceptable here?
Not applicable. Runtime, filesystem/store, output, config, and clock should be composed at app/command edges or sprint service construction.

## 7. State and Mutation Check

### Mutation points

- `internal/sprint` execute run-state store: creates/updates `.run-state.json` atomically.
- `internal/sprint` flow-state store: updates execute stage in `flow-state.json` atomically.
- `internal/sprint` execute summary renderer/store: writes `execute.md`.
- `internal/platform/runtime`/agentwrap invocation: may mutate target repository files through runtime tools, constrained by target path and permission policy.
- Logging/observability sink: records safe runtime/task diagnostics and events.

### Is mutation explicit from names and flow?

```text
[x] Yes
[ ] No - improvement needed: not applicable
```

Mutating functions should be named around actions such as `RunExecute`, `SaveExecuteRunState`, `MarkTaskRunning`, `MarkTaskComplete`, `MarkTaskFailed`, `WriteExecuteSummary`, and `UpdateFlowState`.

### Are queries and commands separated where practical?

```text
[x] Yes
[ ] No - reason: not applicable
```

Status reads should not mutate state except for an explicit refresh/reconcile command path. Execute and flow operations own mutations.

### Could hidden state make this hard to debug?

```text
[ ] No
[x] Yes - risk: runtime-created implementation edits and agentwrap run records can be outside `.run-state.json` unless summarized with run IDs and diagnostics
```

Mitigation: persist task IDs, agentwrap run/session IDs, attempt metadata summaries, validation outcomes, and target path in `.run-state.json`; avoid persisting unsafe raw runtime payloads by default.

## 8. Function/Class/Module Shape

### Primary unit of behaviour

```text
[ ] Function
[ ] Class/object
[x] Module
[x] Service/component
[ ] Pipeline/workflow
[x] Adapter
[ ] Other: not applicable
```

### Why this unit?

The behavior spans validation, prompt rendering, runtime execution, persistence, and status, so a sprint-owned service/use-case surface is the right unit at the module boundary. Internally, pure parsing/rendering/task-ID helpers should be functions. Runtime remains an adapter boundary in platform/agentwrap, not a sprint-owned process abstraction.

### If using a class/object, what justifies it?

```text
- [x] Maintains meaningful state
- [x] Protects invariants
- [ ] Implements polymorphism
- [x] Manages lifecycle
- [ ] Required by framework
- [x] Groups cohesive behaviour
- [ ] Other: not applicable
```

The service object may hold injected dependencies such as stores, runtime, path resolver, clock, logger, and config resolver. Durable product state remains in files, not in long-lived mutable globals.

### If using inheritance, why not composition?

Inheritance is not used. Go composition and interfaces at external boundaries are sufficient.

### If using composition, what components are combined?

- Sprint service/use case: coordinates execute validation, task extraction, runtime calls, state transitions, and summary output.
- Sprint state store: loads and atomically writes `.run-state.json` and flow state.
- Runtime boundary: accepts generic prompt execution requests and returns generic run results/events.
- Workspace path resolver: verifies sprint paths and target implementation directory safety.
- Config/model resolver: applies stage-specific execute model precedence.
- Output/log adapter: renders text/JSON status and safe diagnostics.

## 9. Error Handling Design

### Expected failures

- Missing or invalid prerequisites through `plan`: reject execute with validation failure, cite required artifact and next validation command.
- `plan.md` has no executable tasks or unsupported task syntax: fail before runtime with actionable task-extraction diagnostics.
- Existing `.run-state.json` malformed or unsupported schema: reject load with state-file path and schema/version details.
- Target path escapes workspace or is not explicitly allowed: fail before runtime with workspace/path diagnostic.
- Runtime health/config/model failure: surface agentwrap/platform runtime category and selected model source where safe.
- Runtime failure, timeout, rate limit, permission denial, or cancellation: mark task failed/cancelled or retryable according to policy and persist safe details.
- Runtime success without required evidence: mark task failed unless explicit diagnostic explains why evidence cannot be machine-validated.
- Atomic state write failure: fail the command and preserve the last known good state if rename did not complete.

### Unexpected failures

- Unknown runtime event shape: do not crash by default; record safe native-extension/unknown event facts where provided by agentwrap.
- Unexpected filesystem error: wrap with operation and workspace-relative path when safe.
- Internal invariant violation such as duplicate task IDs: fail before runtime and point to plan/task extraction diagnostics.

### Error taxonomy

Error types/codes/classes/results introduced or reused:
- Validation failure: invalid stage prerequisite, invalid plan, missing evidence, invalid execute summary.
- Workspace/path failure: unsafe target or unresolved workspace/project/sprint path.
- State failure: malformed `.run-state.json`, unsupported schema, atomic write failure.
- Runtime failure: wraps or maps agentwrap `SDKError` categories such as configuration, health, runtime_unavailable, model_unavailable, authentication, permission, rate_limit, timeout, cancellation, malformed_event, runtime_exit, validation, repair_exhausted, cleanup, and unknown.
- Cancellation: user-initiated or context/timeout cancellation mapped to cancelled task state and cancellation exit behavior.

### Retry / recovery behaviour

- Retry safe? yes for failed/pending execute tasks when policy allows and task state remains non-terminal or explicitly re-run; no automatic retry for completed tasks unless forced.
- Partial progress possible? yes; completed tasks remain complete and are not redone on resume unless explicitly forced.
- Compensation needed? no automatic compensation. The product must record diagnostics and preserve state; it must not revert implementation edits or mutate Git state.

## 10. Observability Design

### Events/logs/metrics/traces needed

- Execute command started/completed/failed: emitted around the command/use-case boundary.
- Task extracted: emitted when deterministic task IDs are produced from `plan.md`.
- Task queued/started/completed/failed/cancelled: emitted at each state transition.
- Runtime run started/completed/failed: mapped from platform/runtime and agentwrap metadata.
- Validation/evidence check passed/failed: emitted before marking task complete or failed.
- State write succeeded/failed: emitted for `.run-state.json` and flow-state persistence failures.
- Summary written: emitted when `execute.md` is updated.

### Minimum useful fields

- project
- sprint
- command
- workspace-relative sprint root
- target repository path or target reference
- task ID
- plan source path and fingerprint when safe
- agentwrap run ID and session ID when available
- runtime/provider/model with selected model source
- attempt number
- operation name
- duration
- outcome/status
- error category/type
- validation/evidence result
- state file path

### Sensitive data check

Could this expose PII/secrets/user content?
```text
[ ] No
[x] Yes - mitigation: redact config/env values, omit unsafe raw runtime payloads by default, store paths workspace-relative where possible, summarize user content rather than logging full prompts or file diffs, and use agentwrap safe diagnostic fields
```

## 11. Testing Strategy

### Unit tests

- Plan task extraction from valid `plan.md` fixtures with deterministic task IDs.
- Duplicate task ID detection and unsupported task syntax diagnostics.
- Execute run-state schema validation, malformed JSON rejection, unsupported version rejection, and atomic write/load behavior.
- Task transition invariants for pending, running, complete, failed, and cancelled.
- Stale running task recovery policy with diagnostic.
- Stage-specific execute model resolution precedence and unknown/empty stage-key validation.
- Target repository path safety and workspace escape rejection.
- Execute summary count generation from run state.

### Use-case/workflow tests

- Execute refuses to run when prerequisites through `plan` are invalid.
- Execute initializes `.run-state.json`, runs a pending task through fake runtime, validates evidence, marks complete, and updates status.
- Runtime success without evidence marks task failed with diagnostic.
- Runtime failure persists failed state with safe agentwrap category/metadata.
- Cancellation saves state and marks active task cancelled or retryable according to policy.
- Resume skips completed tasks and recovers stale running tasks conservatively.
- `flow --to execute` stops at execute and rejects smoke/review/issues targets.

### Integration tests

- Gated OpenCode/agentwrap integration may verify execute request mapping and status/event projection when required environment is present.
- Normal test runs must use agentwrap-compatible fake runtime and must not require OpenCode, provider credentials, network access, or target-repository Git mutation.

### Regression tests

- Runtime exit success cannot mark task complete without evidence or explicit diagnostic.
- `.run-state.json` writes are atomic and preserve last known good state after write failure.
- Secrets and full environment values are not present in logs, config diagnostics, or run-state diagnostics.
- Existing historical `review.md` files do not cause Phase 2 to generate or validate automated review artifacts.

### Testability check

Can the core behaviour be tested without real infrastructure?
```text
[x] Yes
[ ] No - why: not applicable
```

The core can be tested with fixture plans, fake workspace/project roots, local temp filesystem stores, fake clock, fake output writers, and fake agentwrap-compatible runtime.

## 12. Performance and Scale Assumptions

### Known constraints

- Expected volume: dozens of plan tasks per sprint initially; design should remain reasonable for hundreds.
- Latency sensitivity: startup/status should stay fast; runtime execution dominates long operations.
- Memory sensitivity: avoid reading large implementation repositories or full runtime payloads into state/logs.
- Concurrency concerns: default should be conservative because multiple coding-agent tasks can conflict in one target repository.

### Assumptions being made

- Execute should run one task at a time by default for target-repository safety unless a later decision enables bounded parallelism with conflict controls.
- `plan.md` is small enough to parse in memory.
- `.run-state.json` can be loaded in memory for expected sprint task counts.
- Status should compute counts from `.run-state.json` without scanning the target repository.
- Agentwrap owns subprocess waiting, cleanup, retry/policy behavior, and event projection.

### Measurement plan

```text
- [ ] Not needed for this change
- [x] Unit benchmark
- [ ] Profiling
- [ ] Load test
- [ ] Production metric
```

Benchmark state load/status count computation only if fixture task counts approach thousands or status becomes noticeably slow. Runtime execution speed should not be optimized inside sprint code.

### Optimization decision

Do not optimize runtime throughput now. The important performance decision is bounded work: no recursive target-repository scans for status, no unbounded goroutines, no excessive state rewrites beyond meaningful transitions, and no raw event payload persistence by default.

## 13. Implementation Plan

### Files likely to change

- `internal/sprint/*execute*.go`: execute use case, request/result types, prompt rendering, task extraction, state transitions, and summary generation.
- `internal/sprint/*state*.go` or sprint store files: `.run-state.json` load/save and atomic persistence.
- `internal/sprint/*validation*.go`: execute prerequisite, plan task, run-state, and execute summary validation.
- `internal/sprint/*flow*.go`: include execute stage after valid plan prerequisites and reject deferred stages.
- `internal/sprint/*status*.go`: show execute progress from `.run-state.json`.
- `internal/app/*`: wire execute use case, config/model resolution, runtime dependency, output adapters, and command exit mapping.
- `cmd/ultraplan/*` or command registration files: expose the sprint execute command, `flow --to execute`, `prompt execute`, and validation/status options through thin CLI glue.
- Tests and fixtures under relevant package test directories: plan fixtures, fake runtime, state files, golden execute summary/status output.

### Step-by-step plan

1. Add sprint-local execute domain/state types and validators for schema version, task statuses, attempts, diagnostics, target path, and plan fingerprint.
2. Implement deterministic task extraction from validated `plan.md`, with clear diagnostics for unsupported or ambiguous tasks.
3. Implement run-state load/init/resume with stale running-task recovery and atomic writes.
4. Implement execute prompt rendering and generic runtime request construction without adding sprint semantics to platform/runtime.
5. Implement execute task runner using injected runtime, context cancellation, safe event/metadata mapping, and validation/evidence completion gates.
6. Implement `execute.md` summary rendering and status/flow integration.
7. Wire app/CLI commands as thin wrappers and reject deferred stages or Git mutation behavior.
8. Add unit, workflow, fake-runtime, fixture, and golden tests for success, failure, cancellation, resume, validation, redaction, and status output.

### Migration/backwards compatibility

- Schema migration needed? no migration for existing execute state because this is the first execute-stage run-state schema; include schema version and reject unsupported versions.
- API compatibility affected? yes for CLI surface by adding execute commands/stage support; no breaking changes are required.
- Existing data affected? existing planning artifacts are read and validated; no existing source/state data should be rewritten except sprint flow state and execute artifacts for this sprint.
- Feature flag needed? no, execute is selected sprint scope; real-runtime integration remains governed by config and fake-runtime tests.

### Rollback plan

If execute behavior breaks, disable or revert the CLI/app wiring for `execute` and `flow --to execute` while leaving planning-through-plan behavior intact. Because `.run-state.json` and `execute.md` are additive sprint artifacts and Git mutation is not automatic, rollback does not require reverting target repository changes automatically; users inspect run state and target diffs manually.

## 14. Final Pre-Implementation Decision

Decision: Proceed.

Reason:
The execute stage fits the existing module-driven architecture. `internal/sprint` owns execute semantics and durable state, `internal/platform/runtime` remains generic, `internal/workspace` owns path safety, and app/CLI remain composition and surface layers. Technical handbook evidence supports this shape through thin command boundaries, module-local ownership, explicit dependency injection, atomic state, bounded concurrency, safe observability, and fakeable tests.

Complexity introduced:
medium.

Complexity removed:
low.

Main trade-off:
Accept sprint-local task-state mechanics instead of prematurely extracting a shared scheduler. This may duplicate some code shape from study run-loop behavior, but it preserves product boundaries and avoids a global workflow abstraction before stable shared semantics exist.

## Evidence Basis And Rejected Alternatives

Key conclusion:
Implement execute as a sprint-owned use case with sprint-local state and validation, runtime calls delegated to platform/runtime and agentwrap, and no shared workflow engine.

Evidence basis:
- `technical-handbook.md` lines 31-35 emphasize thin command boundaries, module-local ownership, and explicit composition roots.
- `technical-handbook.md` lines 39-47 emphasize runtime success plus artifact validation, atomic durable state, safe diagnostics, workspace-safe execution, and fake-runtime tests.
- `technical-handbook.md` lines 63-79 warn against monolithic `RunE`, global workflow/scheduler packages, package-level mutable runtime state, unbounded runtime runs, unsafe raw payload persistence, and out-of-scope smoke/review/issues/Git mutation.
- `ARCHITECTURE.md` lines 208-245 assigns sprint ownership of planning and execute artifacts, task extraction, deterministic task IDs, flow state, execute run-state persistence, and status while excluding smoke, review automation, issue tracking, and Git mutation.
- `ARCHITECTURE.md` lines 287-329 requires platform/runtime to know generic execution only and not understand sprint semantics.
- `ARCHITECTURE.md` lines 483-512 says to reuse infrastructure, not study semantics, and explicitly keeps execute task semantics in `internal/sprint` rather than a global scheduler/workflow package.
- `TRD.md` lines 1520-1586 assigns execute behavior to `internal/sprint`, preserves dependency direction, and allows reuse only of generic infrastructure.
- `TRD.md` lines 1724-1761 define execute `.run-state.json` requirements, statuses, atomic writes, resumability, stable task IDs, target path constraints, and non-goals.
- `TRD.md` lines 1817-1854 define sprint prompt/runtime model resolution and require runtime success plus artifact/evidence validation before completion.
- `PRD.md` lines 46-55 state runtime success is not product success, resumability by default, explicit state/errors, and product-owned workflows.

Rejected alternatives:
- Shared workflow engine: rejected because cross-sprint scheduling/general workflow is a non-goal and would likely move product semantics out of `internal/sprint`.
- Extract study run-loop scheduler/state into shared packages first: rejected because the architecture and TRD explicitly say to reuse infrastructure, not study semantics.
- Put execute semantics in platform/runtime or agentwrap request metadata: rejected because runtime must stay generic and must not understand project/sprint stages.
- Let the CLI command own execute logic: rejected because command handlers should stay thin and testable, with workflow behavior in sprint/app use cases.
- Treat runtime exit success as task completion: rejected because PRD, TRD, requirements, and handbook all require validation/evidence or explicit diagnostics.

## Risks, Assumptions, And Open Questions

Risks:
- Plan task syntax may be too loose for deterministic task extraction; mitigation is to reject unsupported or ambiguous tasks before runtime.
- Runtime-created target edits may conflict if tasks run concurrently; mitigation is one-at-a-time default execution until conflict controls are designed.
- `.run-state.json` can drift from target repository reality if runtime succeeds but evidence is weak; mitigation is evidence validation and explicit diagnostics before completion.
- `execute.md` can become stale after further task attempts; mitigation is to treat `.run-state.json` as source of truth and refresh summary after execute transitions.
- Persisted diagnostics may leak user content or secrets if raw payloads are stored; mitigation is safe summaries, redaction, and omission of unsafe raw runtime payloads by default.

Assumptions:
- The final `reasoning.md` will reference this area document as the selected Architecture reasoning input.
- The first implementation can run execute tasks serially by default without violating sprint acceptance criteria.
- Existing platform/runtime and app composition can inject or adapt agentwrap-compatible fake runtimes for tests.
- The target implementation directory can be represented as a workspace-safe path or explicit configured target reference.

Open questions to resolve in final sprint reasoning or plan:
- Which exact `plan.md` task markers and fields are considered executable for deterministic task ID generation?
- What concrete task evidence qualifies as machine-validatable for this sprint, and when is an explicit diagnostic sufficient?
- Should `.run-state.json` store only run metadata summaries, or also workspace-local links to agentwrap run records when durable run stores are available?
- Should `execute.md` be written after every task transition or only after execute command completion?
- What minimum text and JSON status fields are required for execute progress in this sprint?
