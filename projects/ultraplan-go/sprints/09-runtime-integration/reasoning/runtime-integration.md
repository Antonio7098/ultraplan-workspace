> **Inputs Used:** `projects/ultraplan-go/sprints/09-runtime-integration/technical-handbook.md`, `projects/ultraplan-go/sprints/09-runtime-integration/requirements.md`, `projects/ultraplan-go/sprints/09-runtime-integration/sprint-index.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `system/reasoning/architecture_reasoning_template.md`

# Runtime Integration Architecture Reasoning

## 0. Feature Summary

### Feature name

Runtime integration

### Area name and scope

This area covers UltraPlan Go's platform runtime boundary for sprint `09-runtime-integration`: ownership, agentwrap/OpenCode composition, runtime config mapping, health checks, validation, policy/retry/fallback delegation, permission posture, event and diagnostics mapping, error classification, and test seams.

This area does not cover ceremony. It also does not decide public study execution, batch scheduling, durable workflow orchestration, summaries, code extraction, target workflows, or sprint execution.

### User/product goal

UltraPlan needs a generic, testable runtime boundary that can run later study work through `github.com/Antonio7098/agentwrap` and `agentwrap/opencode` without UltraPlan owning OpenCode process supervision, native event decoding, retry/fallback engines, permission translation, or validation wrappers.

### Current task

Create the reasoning basis for the sprint implementation files listed in `requirements.md`: `internal/platform/runtime/runtime.go`, `agentwrap.go`, `opencode.go`, `health.go`, `policy.go`, `events.go`, runtime config extensions, health command wiring, and fake-runtime/config/health tests.

### Non-goals

- Do not implement public `study run`, `synthesize`, `run-all`, or `run-loop` execution in this sprint.
- Do not create worker pools, durable runtime record migrations, lock files, summary generation, code reference extraction, target workflows, or sprint execution.
- Do not create a public UltraPlan runtime SDK that competes with agentwrap.
- Do not parse OpenCode stdout/stderr, launch OpenCode directly, or reimplement agentwrap policy, health, validation, permission, event, or diagnostics behavior.
- Do not persist unsafe native payload bytes or full sensitive runtime environment values by default.

## 1. First-Principles Breakdown

### Core behaviour

The platform runtime receives a generic execution or health request, validates runtime-facing config and required capabilities, maps the request into agentwrap types, delegates execution/health/policy/validation/permissions/events to agentwrap and `agentwrap/opencode`, then maps safe results, classified errors, events, diagnostics, and metadata back into UltraPlan internal types.

### Inputs

- Resolved UltraPlan config for runtime name, executable, extra args, environment additions, provider, model, timeout, retries, fallback target, health checks, capabilities, sandbox, permissions, and diagnostics limits.
- Runtime request data: prompt, work directory, provider, model, timeout, metadata, health requirements, capability requirements, permission policy, sandbox, validation expectations, and optional session fields.
- Caller `context.Context`, including cancellation and timeout state.
- Agentwrap runtime outputs: `RunResult`, `RunMetadata`, canonical events, health/capability responses, validation results, policy decisions, permission metadata, warnings, usage/cost fields, and `SDKError` values.

### Outputs

- UltraPlan runtime result with success/failure state, classified safe error diagnostics, validation summary, retry/fallback metadata, permission/session/cleanup/artifact/usage summaries, warnings, and safe native metadata.
- UltraPlan runtime event records or log-ready event projections that preserve agentwrap canonical event kind and safe payload metadata.
- Health command runtime diagnostics that are actionable for missing executable, unsupported requirement, unavailable runtime, authentication/provider/model failures, and capability gaps.
- Testable fake-runtime behaviours for success, validation failure, runtime failure, timeout, cancellation, malformed events, rate limit, retry, fallback, permission failure, and observability metadata.

### Durable state

- Created: no new durable runtime store is required by this sprint by default.
- Read: workspace config and existing run/task metadata fields from prior sprint context may be read by callers before constructing runtime metadata.
- Updated: runtime result metadata must be shaped so later run-state persistence can store agentwrap run IDs, session IDs, attempts, validation, policy, permissions, cleanup, artifacts, usage, warnings, and errors.
- Deleted: none.

### Ephemeral state

- Runtime request mapping values.
- Agentwrap runtime stack instance and OpenCode adapter options.
- Event drain state while waiting for a run.
- Health check results and capability checks.
- Redacted diagnostic summaries.
- Fake runtime scripts/records in tests.

### Derived state

- Runtime health status is derived from configured required health/capability checks and agentwrap health/capability responses; it is computed on demand for this sprint.
- Runtime success is derived from underlying runtime result plus validation result; it must not be inferred from process exit alone.
- Attempt and fallback summaries are derived from agentwrap policy metadata and events.
- Usage and cost summaries are derived from agentwrap metadata; unknown values remain unknown, not zero.

### Side effects

- File write: no platform-runtime-owned generated artifact writes; underlying runtime may create expected output files, and validators inspect them through agentwrap validation hooks.
- Network call: possible through OpenCode/provider only inside `agentwrap/opencode`; normal tests must not require network.
- Queue/event emission: event fan-out through agentwrap `EventSink` and UltraPlan log/event projection.
- Logs/metrics/traces: structured logs and diagnostics with run/task/runtime/provider/model correlation and redaction.

## 2. Existing Architecture Fit

### Contract applicability and phase gate

Current gate: runtime/provider integration slice.

Applies now:

- Architecture: `internal/platform/runtime` owns generic runtime execution, and platform packages must not import product modules. `ARCHITECTURE.md` defines `study -> platform/runtime` and `platform/runtime -> no product modules`.
- LLM Runtime: agentwrap is the runtime SDK; production composition is `ObservingRuntime -> ValidatingRuntime -> PolicyRunner -> opencode.Runtime`.
- Configuration: runtime config must be resolved through defaults/config/env/flags, validated once, mapped into agentwrap/opencode options, and redacted.
- Errors: runtime-facing errors preserve agentwrap `SDKError` classifications using `errors.As`/helpers rather than string matching.
- Observability: canonical events, run/task correlation, safe diagnostics, and unknown usage/cost handling are in scope.
- Security: sandbox and permission policy are explicit trust boundaries; unsupported required policy features fail before launch.
- Testing: fake-runtime tests and gated real OpenCode smoke tests are in scope; default test runs must not require OpenCode or credentials.
- CLI Surface: `ultraplan health` may report runtime health but must not invoke study workflows.

Deferred:

- Full batch/durable workflow gate: worker pools, run-loop mutation, durable runtime record migrations, locks, stale-running recovery, and active task cancellation orchestration remain later scope.
- Stable public runtime JSON schema: no new stable public runtime status/health schema is required unless implementation reuses an existing output mode.
- Public study execution: single analysis, synthesis, run-all, and run-loop command execution remain out of this sprint except preserving existing help or stubs.

### Where does this behaviour naturally belong?

Runtime behaviour belongs in `internal/platform/runtime` because it is generic execution infrastructure. Config parsing/validation extensions belong in `internal/platform/config`. CLI health presentation belongs in `internal/app/health_commands.go` and must call the platform runtime health boundary without directly touching agentwrap or OpenCode process details.

### Existing workflow affected

Current flow:

1. Earlier sprints provide app wiring, workspace/config/logging/health foundations, study metadata, prompt composition, validators, and run-state metadata fields.
2. Runtime execution remains deferred; health exists as a skeleton rather than an agentwrap-backed runtime check.
3. Prompt composition and report validation are product-owned but not yet connected to a real runtime boundary.

### Proposed new flow

Proposed flow:

1. App/config resolves and validates runtime settings before launch or health checks.
2. `internal/platform/runtime` builds or receives a generic UltraPlan runtime request with context, prompt, workdir, provider/model, timeout, metadata, requirements, permissions, sandbox, and validation hooks.
3. Platform runtime maps that request to agentwrap `RunRequest` and uses the production wrapper stack `ObservingRuntime -> ValidatingRuntime -> PolicyRunner -> opencode.Runtime`.
4. Runtime events are drained and mapped into safe UltraPlan event/diagnostic records while preserving canonical event kinds.
5. `Run.Wait` returns the final agentwrap result and classified error; UltraPlan maps safe metadata, validation, policy, permissions, usage/cost, artifacts, warnings, and errors.
6. `ultraplan health` calls platform runtime health/capability checks only, then prints actionable redacted diagnostics without running prompts, studies, schedulers, or report validators.

### Does the current architecture still fit?

[x] Yes. The sprint fits the existing module-driven architecture because runtime is already defined as generic platform infrastructure and agentwrap is the selected external runtime boundary.

### If it does not fit, why?

Not applicable. The current architecture anticipates this sprint. The main risk is implementation drift: importing study semantics into `platform/runtime`, direct OpenCode supervision in app code, or hidden global runtime config would violate the existing shape.

### Refactor-before-feature decision

Decision: implement directly with small focused platform/config/app additions.

Reason: `ARCHITECTURE.md` already reserves `internal/platform/runtime` for generic execution, and the sprint requirements list focused files. No larger architectural refactor is needed before implementation.

## 3. Design Options Considered

### Option A: Thin UltraPlan runtime facade over agentwrap

Description:
Define narrow UltraPlan-internal request/result/health/event/config mapping types in `internal/platform/runtime`, compose agentwrap wrappers in an explicit constructor, use `agentwrap/opencode` for OpenCode, and expose only safe generic runtime information to app/study callers.

Pros:
- Matches `ARCHITECTURE.md` dependency direction and `requirements.md` acceptance criteria.
- Keeps OpenCode process handling, native decoding, retries, validation wrappers, health, permissions, and diagnostics with agentwrap.
- Supports fake-runtime tests without OpenCode, credentials, network, or subprocesses.
- Preserves later study execution compatibility without importing study concepts.

Cons:
- Requires careful type mapping to avoid losing agentwrap classifications or overexposing native details.
- Adds internal translation code even though agentwrap already has rich types.
- Future product callers must pass metadata as generic fields rather than typed study structs.

Risks:
- A too-thin facade could leak agentwrap implementation details throughout product code.
- A too-thick facade could become a competing SDK. Keep the facade limited to UltraPlan's internal needs.

### Option B: Product modules call agentwrap directly

Description:
Let future `internal/study` execution code build agentwrap requests, compose wrappers, map events, and handle policy/validation directly.

Pros:
- Fewer runtime translation structs initially.
- Product code can tailor metadata and validation without an intermediate platform type.

Cons:
- Violates the selected architecture pressure that runtime execution is generic platform infrastructure.
- Duplicates mapping and error/diagnostic decisions across commands or product modules.
- Makes `ultraplan health` and tests more likely to depend on concrete OpenCode/agentwrap details.

Risks:
- Study semantics and runtime mechanics become coupled, making later adapters, health checks, and tests harder.

### Option C: UltraPlan-owned runtime supervisor and retry engine

Description:
Implement OpenCode process launch, JSON event decoding, retry/fallback, health, validation, permissions, and diagnostics directly in UltraPlan.

Pros:
- Maximum local control over behaviour and output shape.
- Can tune every edge case without waiting for SDK support.

Cons:
- Directly contradicts the sprint goal and TRD requirement to use agentwrap and not reimplement SDK features.
- Recreates volatile subprocess, event, policy, and permission complexity.
- Raises security and maintenance risk substantially.

Risks:
- Diverges from agentwrap, creates duplicate bugs, and makes OpenCode changes more expensive.

### Option D: General runtime registry/plugin system now

Description:
Create a broad runtime registry with multiple adapters, plugin metadata, and factory selection beyond OpenCode.

Pros:
- Future adapters could plug in with less refactor.
- Runtime names and capabilities could become centrally discoverable.

Cons:
- Exceeds current sprint scope; only OpenCode through agentwrap is required.
- Adds abstraction and config surface before there is a second concrete adapter.
- Conflicts with the handbook's complexity discipline and warning against broad registries without need.

Risks:
- Turns a runtime integration sprint into plugin infrastructure and weakens reviewability.

### Chosen option

Chosen option: Option A.

Reason: This is the smallest honest design that satisfies the sprint. The technical handbook favors thin boundaries, manual composition roots, explicit wrapper stacks, typed errors, injectable IO/fakes, context propagation, structured events, and first-class trust boundaries. Option A applies those pressures without adding a speculative plugin system or duplicating agentwrap.

## 4. Abstraction Check

### Are we adding a new abstraction?

[x] Yes - service/component/module: `internal/platform/runtime` as the generic runtime integration package.
[x] Yes - data structure/DTO/config object: narrow UltraPlan-internal runtime request, result, health, event, diagnostic, and config mapping types.
[x] Yes - strategy/function parameter: validation hooks/custom validators passed through agentwrap validation specs.
[ ] Yes - plugin/registry/factory.

### If yes, why is it earned?

- It isolates an external dependency: agentwrap/OpenCode is the volatile runtime boundary.
- It protects a stable domain concept: UltraPlan needs generic execution without study semantics.
- It makes testing simpler: fakes can implement agentwrap-compatible runtime/run behaviour and assert mapping.
- It creates a clear boundary between policy and mechanism: agentwrap owns runtime policy mechanics; UltraPlan owns config/request/result mapping and product validators.
- It prevents direct OpenCode or agentwrap details from spreading into CLI health and future product modules.

### If no, why are we keeping it concrete?

Not applicable for the platform runtime boundary. However, a general plugin/registry abstraction is intentionally not added because only one production runtime adapter is selected and agentwrap already provides the external runtime interface.

### Bad abstraction smell check

- Generic names like `Manager`, `Handler`, `Processor`, `Common`, or `Helper` are rejected.
- A broad runtime plugin registry is rejected because there is no second concrete adapter in scope.
- Boolean mode switches for study vs synthesis vs sprint are rejected inside `platform/runtime`; callers pass generic metadata and validation hooks.
- Optional caller-specific behaviour should be represented by explicit request fields or omitted.

## 5. DRY and Duplication Check

### Is any duplication being introduced?

[x] Yes - duplicated code shape only: mapping between UltraPlan internal types and agentwrap types.
[ ] Yes - duplicated business/domain knowledge.
[ ] Yes - duplicated infrastructure mechanics.

### If duplicating business knowledge, stop and consolidate

No business knowledge should be duplicated. Study report semantics remain in `internal/study`; platform runtime only accepts generic validation expectations and custom validator hooks.

### If duplicating code shape, is that acceptable?

The mapping layer is acceptable because it prevents product modules and CLI health from depending directly on every agentwrap type while preserving agentwrap classifications and metadata. It must stay thin and field-for-field where possible.

### If removing duplication, is the abstraction safe?

Do not consolidate agentwrap wrapper mechanics into UltraPlan helpers beyond construction/mapping. The single source of truth for retry/fallback, validation wrapper behaviour, health, OpenCode process launch, native event projection, and permission translation remains agentwrap.

## 6. Coupling Check

### Dependencies required

- `github.com/Antonio7098/agentwrap`: required SDK boundary for runtime execution, run requests/results, capabilities, health, validation, policy, events, metadata, stores/sinks, and `SDKError` classification.
- `github.com/Antonio7098/agentwrap/opencode`: required concrete runtime adapter for OpenCode construction and process/native event ownership.
- `internal/platform/config`: required to resolve and validate runtime config before constructing the stack.
- `internal/platform/logging` or existing app logging path: required to emit redacted diagnostics and structured runtime event summaries.
- `context.Context`: required on health and run paths for cancellation propagation.

### Coupling risks

Global coupling: none introduced if runtime config and stack are passed through constructors rather than globals.

Content coupling: none introduced if platform runtime does not import `internal/study`, `internal/codeextract`, or product modules.

Stamp coupling: possible if runtime request accepts a large config object; mitigate by resolving a narrow runtime config struct before construction.

### Narrow dependency check

Each function should receive only what it needs. OpenCode construction receives resolved executable/args/env/stderr options. Run mapping receives request fields and validation hooks. Health receives required health/capability names and runtime config. Event mapping receives agentwrap events and redaction policy.

Large objects should be passed only when the full concept is needed. Do not pass full workspace config into the OpenCode adapter constructor if a narrower runtime config is enough.

### Dependency injection check

Side-effectful dependencies should be created at the edge and passed inward. The app composition root can call a platform runtime constructor, but tests should inject fake agentwrap-compatible runtimes, fake health checkers, fake sinks/stores, and buffer-backed IO for health output.

## 7. State and Mutation Check

### Mutation points

- `internal/platform/config/config.go`: adds and validates runtime config fields.
- `internal/platform/runtime/agentwrap.go`: constructs the runtime wrapper stack and any in-memory store/sink selected for this sprint.
- `internal/platform/runtime/events.go`: accumulates safe event/result diagnostics during a run.
- `internal/app/health_commands.go`: extends health output with runtime health diagnostics when requested/configured.
- Future product callers: may persist mapped runtime metadata into run state, but durable mutation is not required in this sprint.

### Is mutation explicit from names and flow?

[x] Yes. Constructor functions build runtime dependencies; run/health methods perform external checks or execution; mapping functions return derived results.

### Are queries and commands separated where practical?

[x] Yes. Health/capability checks are read/query paths; runtime execution is a command path. `ultraplan health` must not start runs or execute study workflows.

### Could hidden state make this hard to debug?

[x] No, if globals are avoided and agentwrap store/sink choices are explicit. Hidden state would become a risk if runtime config, stores, or active runs are global singletons.

## 8. Function/Class/Module Shape

### Primary unit of behaviour

[x] Module/package.
[x] Adapter.
[x] Service/component.

### Why this unit?

The runtime package groups one cohesive volatile boundary: mapping UltraPlan runtime requests to agentwrap and mapping agentwrap results back. It is not product workflow logic and not CLI presentation logic.

### If using a class/object, what justifies it?

Go structs are justified for values that hold cohesive dependencies such as an agentwrap runtime, health checker, event sink/store, redactor, and resolved config. They manage lifecycle-adjacent external resources but should not hide global state.

### If using inheritance, why not composition?

Not applicable. Go composition and interfaces at volatile boundaries are sufficient.

### If using composition, what components are combined?

- `opencode.Runtime`: owns OpenCode process execution and native event projection.
- `agentwrap.PolicyRunner`: owns retry, wait, rate-limit handling, and fallback.
- `agentwrap.ValidatingRuntime`: owns validation expectations and logical failure on invalid outputs.
- `agentwrap.ObservingRuntime`: owns event/run observation, active/completed run projection, store/sink integration, and merged metadata.
- UltraPlan runtime mapper: owns generic request/result/event/health/config translation and safe diagnostics.

## 9. Error Handling Design

### Expected failures

- Invalid runtime config: fail before launch with field path and corrective guidance.
- Unknown runtime name, unsupported health/capability name, empty executable, invalid timeout, invalid retry count: fail before launch.
- Unsupported required permission policy feature: fail before launch unless explicitly configured as best-effort.
- Missing executable or unavailable runtime: surface runtime health/config diagnostic without prompt execution.
- Authentication/provider/model failure: surface agentwrap classification and safe provider/model/runtime fields.
- Timeout/cancellation/rate limit/runtime exit/malformed event/validation/repair exhausted/cleanup failure: preserve agentwrap classification and safe metadata.
- Validation failure after runtime success: mark logical runtime result failed and include safe validation summary.

### Unexpected failures

- Unexpected agentwrap or adapter errors are wrapped with operation context while preserving cause chains.
- Unknown SDK categories are surfaced as classified unknown runtime errors with safe operation/runtime/provider/model metadata when available.
- Unknown native events projected as `native_extension` are recorded safely and do not fail a run by themselves.

### Error taxonomy

Reuse agentwrap `SDKError` categories for runtime failures: configuration, health, runtime unavailable, provider unavailable, model unavailable, authentication, permission, rate limit, timeout, cancellation, malformed event, runtime exit, validation, repair exhausted, cleanup, and unknown.

UltraPlan-specific product categories remain outside `platform/runtime` except where callers wrap the returned runtime result/error at product boundaries.

### Retry / recovery behaviour

- Retry safe: only through `agentwrap.PolicyRunner`; UltraPlan does not implement retry loops inside the OpenCode adapter.
- Partial progress possible: yes, underlying runtime may produce artifacts or events before failure; mapping must surface safe artifact/cleanup/validation metadata.
- Compensation needed: no platform-owned compensation in this sprint; cleanup metadata comes from agentwrap.

## 10. Observability Design

### Events/logs/metrics/traces needed

- Runtime health check started/completed/failed, with runtime, provider/model when relevant, required checks/capabilities, outcome, duration, and safe diagnostics.
- Runtime run started/completed/failed/cancelled, with UltraPlan run ID, task ID, task kind, runtime, provider, model, workdir summary, timeout, and agentwrap run/session IDs.
- Agentwrap canonical events: lifecycle, session, message, progress, tool, artifact, permission, blocking, usage, warning, fatal_error, rate_limit, validation, retry, fallback, final_result, and native_extension.
- Policy decisions: attempts, rate-limit wait, retry, fallback target, exhaustion.
- Validation summary: passed/skipped/failed counts, safe check names, artifact paths, and failure reasons.
- Raw payload omission facts: presence, source, encoding, safety, and omission reason; raw bytes are omitted by default.

### Minimum useful fields

- UltraPlan run ID and task ID.
- Task kind, source kind, study/dimension/source/output path as generic metadata strings supplied by callers.
- Runtime kind, provider, model, executable path or safe executable label, timeout, attempt, fallback target.
- Agentwrap run ID, session ID, event ID, event kind, event type, timestamp.
- Outcome, duration, validation status, policy decision, permission status, usage/cost knownness, classified error category.

### Sensitive data check

Could this expose secrets or user content? Yes.

Mitigation: redact sensitive config/env values, do not log full prompts by default, do not persist unsafe native payload bytes by default, record raw payload omission facts instead of bytes, and use agentwrap safe fields/redaction helpers where practical.

## 11. Testing Strategy

### Unit tests

- Config defaults, precedence where applicable, validation errors, required health names, capabilities, timeout/retry/fallback validation, empty executable rejection, unknown runtime rejection, and redaction.
- Request mapping from UltraPlan internal request into agentwrap `RunRequest`, including context, prompt, workdir, provider/model, timeout, metadata, health/cap requirements, permissions, sandbox, validation, and session fields.
- Event mapping for all selected canonical event kinds, unknown `native_extension`, raw payload omission facts, malformed event handling, usage/cost unknownness, and warnings.
- Error mapping using `errors.As`/agentwrap helpers with representative `SDKError` categories.

### Use-case/workflow tests

- Successful fake runtime run writes expected artifact and returns successful logical result.
- Underlying runtime success with validation failure returns failed logical result.
- Runtime failure, timeout, cancellation, malformed event, rate limit, permission failure, retry, fallback, and exhaustion paths are covered with agentwrap-compatible fakes.
- Health command reports success, missing executable, unavailable runtime, unsupported health requirement, redacted diagnostics, and stable runtime failure exit status.
- Health command does not invoke study workflows, prompt execution, batch scheduler, or report validation logic.

### Integration tests

- Optional real OpenCode smoke tests may be added only when gated by explicit environment variables and skipped by default.
- Any real smoke path must use `agentwrap/opencode`, not direct process invocation.

### Regression tests

- Unknown usage/cost values remain unknown rather than zero.
- Unknown native events do not fail runs by themselves.
- Unsafe native raw payload bytes are not persisted/logged by default.
- Unsupported required permission features fail before launch.

### Testability check

[x] Yes, the core behaviour can be tested without real infrastructure by injecting agentwrap-compatible fake runtimes, fake health checkers, fake events, fake stores/sinks, temp files, and buffer-backed IO.

## 12. Performance and Scale Assumptions

### Known constraints

- Expected volume for this sprint: single run/health paths and fake tests, not full batch execution.
- Latency sensitivity: health output should fail fast and be actionable; runtime execution latency is dominated by OpenCode/provider and governed by timeout/context.
- Memory sensitivity: do not retain unsafe raw payload bytes or full prompt/native event streams unnecessarily; store safe summaries.
- Concurrency concerns: event channels must be drained and `Run.Wait` must be called for every started run to avoid blocking/leaks.

### Assumptions being made

- Agentwrap APIs provide the selected wrapper, health, validation, policy, permission, event, metadata, store/sink, and error concepts described in PRD/TRD/requirements.
- A local in-memory store or fake store is sufficient for this sprint unless implementation needs local inspection beyond process lifetime.
- Future durable run-state persistence can consume this sprint's mapped metadata without changing the runtime boundary.

### Measurement plan

[x] No benchmark or load test is needed for this sprint.

Evidence of acceptable performance is default tests completing without real subprocesses, health failing fast on invalid config, event drains not blocking, and no unbounded goroutine/channel behaviour in fake tests.

### Optimization decision

Do not optimize beyond clear timeouts, bounded diagnostics, redaction, event draining, and no unnecessary raw payload retention. Broad throughput tuning belongs to later batch/run-loop work.

## 13. Implementation Plan

### Files likely to change

- `/home/antonioborgerees/coding/ultraplan-go/go.mod`: add `github.com/Antonio7098/agentwrap` dependency, with a local replace only if required for development.
- `/home/antonioborgerees/coding/ultraplan-go/internal/platform/runtime/runtime.go`: define thin UltraPlan internal runtime request, result, health, validation, permission, event, metadata, and diagnostic-facing types.
- `/home/antonioborgerees/coding/ultraplan-go/internal/platform/runtime/agentwrap.go`: build default wrapper stack and expose testable constructor seams.
- `/home/antonioborgerees/coding/ultraplan-go/internal/platform/runtime/opencode.go`: map resolved config into `agentwrap/opencode` options.
- `/home/antonioborgerees/coding/ultraplan-go/internal/platform/runtime/health.go`: map health/capability requirements and safe diagnostics.
- `/home/antonioborgerees/coding/ultraplan-go/internal/platform/runtime/policy.go`: map retry/fallback/sandbox/permission posture into agentwrap policy structures.
- `/home/antonioborgerees/coding/ultraplan-go/internal/platform/runtime/events.go`: map events, metadata, usage/cost knownness, raw payload omission, warnings, and errors.
- `/home/antonioborgerees/coding/ultraplan-go/internal/platform/config/config.go`: add minimal runtime config fields, validation, and redaction.
- `/home/antonioborgerees/coding/ultraplan-go/internal/app/health_commands.go`: wire runtime health through platform runtime only.
- `/home/antonioborgerees/coding/ultraplan-go/internal/platform/runtime/runtime_test.go`: add fake-runtime and mapping tests.
- `/home/antonioborgerees/coding/ultraplan-go/internal/platform/config/config_test.go`: add runtime config tests.
- `/home/antonioborgerees/coding/ultraplan-go/internal/app/health_commands_test.go`: add health command runtime tests.

### Step-by-step plan

1. Add the agentwrap dependency and inspect available APIs to align exact type names and options.
2. Define narrow platform runtime types and mapping helpers without importing product modules.
3. Add resolved runtime config fields, validation, defaults, and redaction.
4. Implement OpenCode adapter construction through `opencode.NewRuntime` and options only.
5. Implement explicit production wrapper composition: `ObservingRuntime -> ValidatingRuntime -> PolicyRunner -> opencode.Runtime`.
6. Implement health/capability mapping and CLI health wiring without invoking study workflows.
7. Implement policy, permission, event, metadata, diagnostics, usage/cost, and error mapping.
8. Add agentwrap-compatible fakes and tests for request mapping, wrapper composition, health, policy, permissions, events, validation, retry/fallback, cancellation, redaction, and error classification.
9. Run `go test ./...` and `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`.

### Migration/backwards compatibility

- Schema migration needed: no durable runtime schema migration is selected.
- API compatibility affected: no public runtime SDK should be introduced.
- Existing data affected: no existing study artifacts should be rewritten by platform runtime integration.
- Feature flag needed: no, except optional real OpenCode smoke tests must be environment-gated and skipped by default.

### Rollback plan

Revert runtime integration files, config extensions, and health wiring as a unit. Because no durable migration is selected, rollback should not require data repair. If agentwrap dependency causes build issues, remove the dependency and restore health to pre-runtime skeleton behaviour.

## 14. Final Pre-Implementation Decision

Decision: Proceed.

Reason: The selected sprint, architecture docs, TRD, PRD, requirements, and technical handbook all converge on the same shape: a thin generic platform runtime boundary over agentwrap/opencode with explicit wrapper composition, config validation, safe diagnostics, typed errors, fake tests, and no product workflow or OpenCode supervision leakage.

Complexity introduced: medium.

Complexity removed: high relative to implementing OpenCode supervision, retry, validation, permission, event decoding, and diagnostics directly in UltraPlan.

Main trade-off: UltraPlan accepts a careful internal mapping layer so product modules and CLI health stay decoupled from raw agentwrap/OpenCode details, while agentwrap remains the source of truth for runtime mechanics.

## Key Conclusion and Evidence Basis

The final decision is to implement runtime integration as a generic `internal/platform/runtime` boundary that composes agentwrap wrappers in the required order and delegates OpenCode supervision to `agentwrap/opencode`.

Evidence basis:

- `requirements.md` requires platform runtime ownership, no product imports, agentwrap root package usage, default `ObservingRuntime -> ValidatingRuntime -> PolicyRunner -> opencode.Runtime` composition, OpenCode construction through `agentwrap/opencode`, config validation, safe diagnostics, typed errors, health checks, event mapping, fake tests, and gated smoke tests.
- `ARCHITECTURE.md` defines runtime as generic execution infrastructure and states that `study -> platform/runtime` while `platform/runtime -> no product modules`.
- `TRD.md` requires agentwrap dependency, OpenCode adapter usage, wrapper composition, validation through agentwrap, policy through `PolicyRunner`, canonical events, SDKError classification, permission policy mapping, observability hooks, and default tests without OpenCode.
- `PRD.md` says runtime success is not product success, product workflows stay product-owned, runtime adapters execute agent work, health checks are required, unknown event types must not crash the system, and normal tests use fake runtimes.
- `technical-handbook.md` highlights thin runtime boundaries, manual composition roots, decorator/wrapper composition, merged config validation, typed errors, injectable IO/fake runtime seams, context propagation, structured events/logs, explicit trust boundaries, and behaviour-focused tests.

## Rejected Alternatives

- Direct product-module agentwrap usage is rejected because it couples study workflow code to runtime mechanics and weakens the platform boundary.
- UltraPlan-owned OpenCode subprocess supervision is rejected because agentwrap/opencode owns process launch, stdout/stderr handling, native event projection, health, permissions, retry/fallback, validation, and diagnostics.
- Validation failures participating in retry/fallback by default are rejected for this sprint because the selected wrapper order places `ValidatingRuntime` outside `PolicyRunner`; any exception needs explicit reviewed reasoning later.
- A broad plugin/registry system is rejected because only OpenCode through agentwrap is selected and premature registry complexity conflicts with the handbook's maintainability pressure.
- Persisting unsafe native payload bytes by default is rejected because security and observability requirements prefer safe diagnostic facts and redaction.

## Risks, Assumptions, and Open Questions

### Risks

- Agentwrap API shape may differ from the names assumed by sprint docs; implementation must inspect the real dependency and adapt while preserving the required behaviour.
- The mapping layer can accidentally become too thick and duplicate agentwrap concepts; keep it narrow and test observable behaviour.
- Health output can become a hidden execution path if it starts prompts or study validation; tests must prevent this.
- Permission mapping may expose adapter limitations; unsupported required features must fail before launch.
- Event drains can block or leak goroutines if tests do not consume events and call `Wait`.

### Assumptions

- Agentwrap exposes enough APIs for the selected wrapper composition, health/capability checks, validation specs/hooks, policy/fallback, permission structures, event sinks/stores, metadata, and `SDKError` classification.
- Existing app/config/health foundations from prior sprints can accept minimal additions without broad refactor.
- In-memory/fake stores are enough for this sprint unless implementation needs local run inspection beyond process lifetime.
- Later study execution can pass generic metadata fields without requiring `platform/runtime` to import study types.

### Open questions

- Which exact agentwrap health check and capability IDs are available and should be accepted by config validation for this sprint?
- Which permission policy features can `agentwrap/opencode` represent as required versus best-effort, especially path-level rules?
- Which agentwrap store/sink implementation should production use now: memory-only, workspace-backed, or adapter-provided defaults?
- How should optional repair configuration be represented now, if at all, without pulling later study execution scope into this sprint?
- What exact `ultraplan health` text/JSON shape is already present and should be extended without creating a new stable public runtime schema?

## Consolidation Reference

This area reasoning document is selected by `sprint-index.md` as `projects/ultraplan-go/sprints/09-runtime-integration/reasoning/runtime-integration.md` and should be summarized by the sprint's consolidated reasoning artifact at `projects/ultraplan-go/sprints/09-runtime-integration/reasoning.md`.
