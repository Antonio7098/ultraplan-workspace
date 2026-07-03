# Sprint Requirements: 09 Runtime Integration

> Project: `ultraplan-go`
> Sprint: `09-runtime-integration`
> Purpose: the authoritative, human-readable sprint contract. All other sprint artifacts must satisfy these requirements.

## Sprint Goal

Integrate UltraPlan's generic runtime boundary with `github.com/Antonio7098/agentwrap` and `agentwrap/opencode` so later study execution can use OpenCode through agentwrap wrappers, health checks, policies, validation, permissions, and events without UltraPlan reimplementing runtime supervision.

## Required Outputs

| Output | Path | Description |
| --- | --- | --- |
| Sprint requirements | `projects/ultraplan-go/sprints/09-runtime-integration/requirements.md` | Defines the binding goal, scope, acceptance criteria, constraints, dependencies, and review expectations for this sprint. |
| Sprint index | `projects/ultraplan-go/sprints/09-runtime-integration/sprint-index.md` | Selects the roadmap entries, source docs, evidence reports, contracts, and review protocols that govern this runtime integration sprint. |
| Technical handbook | `projects/ultraplan-go/sprints/09-runtime-integration/technical-handbook.md` | Gives implementation guidance for the platform runtime boundary, agentwrap composition, OpenCode adapter construction, health checks, policy mapping, validation, events, and tests. |
| Runtime integration reasoning | `projects/ultraplan-go/sprints/09-runtime-integration/reasoning/runtime-integration.md` | Records reasoning for runtime ownership, wrapper order, config mapping, permission posture, event/log mapping, test seams, and deferred study orchestration. |
| Consolidated reasoning | `projects/ultraplan-go/sprints/09-runtime-integration/reasoning.md` | Summarizes accepted sprint decisions, tradeoffs, risks, and requirement mappings. |
| Implementation plan | `projects/ultraplan-go/sprints/09-runtime-integration/plan.md` | Provides the ordered task plan and verification checklist for implementation. |
| Sprint review | `projects/ultraplan-go/sprints/09-runtime-integration/review.md` | Reviews completed implementation against these requirements and records carry-forward decisions. |
| Agentwrap dependency update | `/home/antonioborgerees/coding/ultraplan-go/go.mod` | Adds `github.com/Antonio7098/agentwrap` as the runtime SDK dependency, using a local replace only if required for development. |
| Platform runtime API | `/home/antonioborgerees/coding/ultraplan-go/internal/platform/runtime/runtime.go` | Defines UltraPlan's thin internal runtime request, result, health, validation, permission, and event-facing types without study semantics. |
| Agentwrap runtime composition | `/home/antonioborgerees/coding/ultraplan-go/internal/platform/runtime/agentwrap.go` | Builds the default agentwrap runtime stack using `ObservingRuntime`, `ValidatingRuntime`, `PolicyRunner`, and a concrete runtime. |
| OpenCode adapter construction | `/home/antonioborgerees/coding/ultraplan-go/internal/platform/runtime/opencode.go` | Constructs the `agentwrap/opencode` runtime from resolved UltraPlan config without direct OpenCode process handling. |
| Runtime health mapping | `/home/antonioborgerees/coding/ultraplan-go/internal/platform/runtime/health.go` | Maps configured health checks to agentwrap health/capability checks and exposes safe health diagnostics. |
| Runtime policy and permission mapping | `/home/antonioborgerees/coding/ultraplan-go/internal/platform/runtime/policy.go` | Maps retries, fallback targets, sandbox, and permission posture into agentwrap policy and permission structures. |
| Runtime event and diagnostics mapping | `/home/antonioborgerees/coding/ultraplan-go/internal/platform/runtime/events.go` | Maps agentwrap events, metadata, safe diagnostics, validation summaries, retry/fallback records, and usage/cost signals into UltraPlan runtime results/log fields. |
| Runtime config extensions | `/home/antonioborgerees/coding/ultraplan-go/internal/platform/config/config.go` | Extends config only as needed for runtime adapter options, policy, permissions, health, and redacted diagnostics. |
| Health command runtime wiring | `/home/antonioborgerees/coding/ultraplan-go/internal/app/health_commands.go` | Extends `ultraplan health` to report runtime health through the platform runtime boundary when runtime checks are requested or configured. |
| Runtime unit and fake tests | `/home/antonioborgerees/coding/ultraplan-go/internal/platform/runtime/runtime_test.go` | Tests agentwrap request mapping, wrapper composition, health mapping, policy/permission mapping, event handling, validation failure handling, retry/fallback behavior, and safe diagnostics with fakes. |
| Runtime config tests | `/home/antonioborgerees/coding/ultraplan-go/internal/platform/config/config_test.go` | Tests runtime config defaults, validation, environment/CLI precedence where applicable, required health names, and redaction of sensitive runtime values. |
| Health command tests | `/home/antonioborgerees/coding/ultraplan-go/internal/app/health_commands_test.go` | Tests runtime health output, failure diagnostics, exit status mapping, and no direct OpenCode/process invocation from the CLI command. |

## Acceptance Criteria

- [ ] `go.mod` declares `github.com/Antonio7098/agentwrap` as the runtime SDK dependency and runtime integration imports agentwrap root package APIs rather than defining a competing public runtime SDK.
- [ ] `internal/platform/runtime` owns generic runtime execution concerns only: prompt, work directory, provider, model, timeout, metadata, health requirements, capability requirements, permissions, validation expectations, events, and results.
- [ ] `internal/platform/runtime` does not import `internal/study`, `internal/codeextract`, or other product modules.
- [ ] The default production composition is `agentwrap.ObservingRuntime` -> `agentwrap.ValidatingRuntime` -> `agentwrap.PolicyRunner` -> `opencode.Runtime`, unless reasoning records an explicit, reviewed wrapper-order exception.
- [ ] OpenCode execution is constructed through `agentwrap/opencode` options such as executable, extra args, environment, and stderr limits; UltraPlan code does not launch `opencode` directly or parse OpenCode stdout/stderr directly.
- [ ] UltraPlan config maps runtime executable, provider/model, timeout, retry count, fallback model/provider when configured, required health checks, required capabilities, sandbox, and permission posture into agentwrap request/config structures.
- [ ] Config validation rejects unknown runtime names, unsupported required health check names, empty runtime executable values, invalid timeouts, invalid retry counts, and unsupported required permission policy features before runtime launch.
- [ ] Sensitive runtime environment or config values are redacted from config output, logs, health diagnostics, and runtime error diagnostics.
- [ ] Runtime requests carry `context.Context` and preserve cancellation semantics through agentwrap calls.
- [ ] Runtime request metadata supports correlation fields needed by later study execution: UltraPlan run ID, task ID, task kind, study, dimension, source, source kind, output path, runtime, provider, and model, without making platform runtime depend on study types.
- [ ] Agentwrap `RunRequest.Validation` can express required output file validation and custom validation hooks without embedding study report semantics in `internal/platform/runtime`.
- [ ] Validation failures from `agentwrap.ValidatingRuntime` make the logical runtime result fail even when the underlying runtime exits successfully.
- [ ] Retry, rate-limit, fallback, and exhaustion behavior is delegated to `agentwrap.PolicyRunner` and surfaced through result metadata and events; UltraPlan does not implement retry/fallback inside the OpenCode adapter.
- [ ] Permission policy mapping uses agentwrap permission structures and fails before launch when required policy features cannot be represented by the OpenCode adapter.
- [ ] Runtime health checks use agentwrap health/capability APIs where supported and return actionable diagnostics for missing executable, missing required capability, unsupported health check, authentication/provider/model failures, and runtime unavailable states.
- [ ] `ultraplan health` can include runtime health information without invoking any study workflow, prompt execution, batch scheduler, or report validation logic.
- [ ] Agentwrap canonical event kinds are accepted and mapped without dropping the event kind: lifecycle, session, message, progress, tool, artifact, permission, blocking, usage, warning, fatal_error, rate_limit, validation, retry, fallback, final_result, and native_extension.
- [ ] Unknown or future native events projected as `native_extension` do not fail a run by themselves.
- [ ] Unsafe raw native payload bytes are not persisted or logged by default; diagnostics record raw payload presence, source, encoding, safety, and omission reason where available.
- [ ] Runtime results expose safe metadata for attempts, policy decisions, validation, repair, sessions, permissions, cleanup, artifacts, usage, estimated cost when available, warnings, and classified errors.
- [ ] Unknown token or usage values remain unknown and are not converted to zero.
- [ ] Runtime-facing errors preserve agentwrap `SDKError` classifications using `errors.As` or agentwrap helpers rather than string matching.
- [ ] Runtime error diagnostics include safe category, operation, user detail, provider, model, runtime kind, exit code/signal when available, retry-after value when available, and redacted metadata.
- [ ] Fake-runtime tests cover successful run, successful runtime exit with validation failure, runtime failure, timeout, cancellation, malformed event, rate limit, validation failure, retry, fallback, permission failure, and observability metadata.
- [ ] Health command tests cover successful health, missing executable or unavailable runtime, unsupported health requirement, redacted diagnostics, and stable exit status for runtime failure.
- [ ] Runtime tests do not require OpenCode, provider credentials, network access, or real subprocess execution.
- [ ] Any real OpenCode smoke test added by this sprint is gated by explicit environment variables and skipped by default; it must go through `agentwrap/opencode`.
- [ ] Existing prompt composition, run-state persistence, status, report validation, study listing, and study initialization tests continue to pass.
- [ ] `go test ./...` passes in `/home/antonioborgerees/coding/ultraplan-go` without requiring OpenCode or provider credentials.
- [ ] `go build ./cmd/ultraplan` passes in `/home/antonioborgerees/coding/ultraplan-go`.

## Non-Goals

- Public execution of `ultraplan study <study> run <dimension-ref> <source-ref>` is not included beyond preserving existing help or placeholders.
- Public execution of `ultraplan study <study> synthesize <dimension-ref>` is not included beyond preserving existing help or placeholders.
- `ultraplan study <study> run-all`, bounded worker pools, full batch scheduling, and synthesis-after-analysis orchestration are not included.
- `ultraplan study <study> run-loop` resume execution, per-study locks, stale-running recovery persistence, and active task cancellation are not included.
- Persisting agentwrap run records durably across process exits is not required unless explicitly selected in the sprint index; in-memory or fake stores are sufficient for this sprint's tests.
- Study report semantics must not move into `internal/platform/runtime`; study-specific validators and task construction remain owned by `internal/study` for later execution sprints.
- Summary generation and `summary.csv` writing are not included.
- Code reference extraction is not included.
- Stable public JSON schemas for runtime status or health are not required unless already present; any JSON output added must be explicitly documented and tested.
- Target workflows, sprint planning, and sprint execution remain out of scope.
- UltraPlan must not implement a new OpenCode subprocess supervisor, native event decoder, retry engine, permission translator, or validation wrapper that duplicates agentwrap behavior.

## Constraints

- Runtime infrastructure must live in `internal/platform/runtime`; do not create runtime behavior inside `internal/study`, `internal/app`, or global product packages.
- Platform runtime may expose thin UltraPlan-internal seams for dependency injection and tests, but the external runtime integration boundary must be agentwrap's root package API.
- Preserve the dependency direction from `docs/ARCHITECTURE.md`: `study -> platform/runtime -> agentwrap/opencode`; `platform/runtime` must not know study semantics or import product modules.
- Apply the roadmap Runtime / Provider Gate: `context.Context` propagation, runtime/process execution ports, bounded retry ownership through agentwrap, correlation metadata, cost/usage drivers in operational signals, fake-runtime tests, and gated real-provider smoke coverage when available.
- Do not apply full Batch / Durable Workflow Gate scope except where runtime metadata must be compatible with existing run-state placeholders; worker pools, run-loop mutation, lock files, and durable run stores are later scope.
- Runtime config and diagnostics must use redaction for sensitive values and must not log full environment variables.
- Health checks and runtime requests must fail fast on invalid config or unsupported required capabilities before launching OpenCode.
- Tests must use agentwrap-compatible fakes for normal coverage and must not depend on OpenCode, provider credentials, network access, or long-running subprocesses.
- Any optional OpenCode smoke test must be opt-in, environment-gated, and skipped by default in `go test ./...`.
- Errors must preserve cause chains and agentwrap classifications; do not classify runtime failures by parsing human-readable error strings.
- Event handling must drain or consume agentwrap run events in tests so event channels cannot block.
- Runtime metadata and event mapping must not convert unknown usage or cost fields to zero values.
- Implementation must remain reviewable and minimal; avoid speculative plugin systems, broad registries, or a general workflow engine.

## Dependencies

| Prior Sprint / Output | Required For | Notes |
| --- | --- | --- |
| Sprint 1: Go Module, CLI Shell, and App Composition | Buildable CLI, package structure, and app wiring | Required for adding runtime package files, health command wiring, tests, and `cmd/ultraplan` build continuity. |
| Sprint 2: Workspace, Config, Logging, and Health Skeleton | Config loading, redaction, logging, and health command foundation | Required for runtime config mapping, secret-safe diagnostics, and `ultraplan health` runtime checks. |
| Sprint 3: Study Domain, Listing, and Resolution | Study metadata concepts for future runtime metadata | Required only as context; platform runtime must not import study. |
| Sprint 4: Study Initialization From YAML | Workspace study layout and prompt/template locations | Required as context for future runtime work directories and output expectations. |
| Sprint 5: Markdown Document Sources and Applicability | Source kind and applicability behavior | Required as context for later task metadata and permission posture; no new applicability behavior is required here. |
| Sprint 6: Report Validation and Rating Parsing | Existing product validators | Required so runtime validation can support custom validator hooks without moving report semantics into platform runtime. |
| Sprint 7: Prompt Composition Generation | Deterministic prompt builders and manifests | Required because runtime requests will later consume composed prompts and expected output paths. |
| Sprint 8: Run State Persistence | Run IDs, task IDs, task metadata placeholders, status summaries | Required so runtime result metadata can align with persisted task state in later execution sprints. |
| `projects/ultraplan-go/docs/PRD.md` | Product runtime behavior | Governs runtime success not being product success, OpenCode runtime adapter expectations, health checks, logging, diagnostics, cancellation, and deferred target/sprint scope. |
| `projects/ultraplan-go/docs/TRD.md` | Technical runtime behavior | Governs agentwrap dependency, wrapper composition, request mapping, health, events, validation, policy, permissions, errors, tests, and build requirements. |
| `projects/ultraplan-go/docs/ARCHITECTURE.md` | Package ownership and dependency direction | Governs keeping runtime generic in `internal/platform/runtime` and product workflow behavior outside platform packages. |
| `projects/ultraplan-go/sprints/05-study-listing-inspection/review.md` | Carry-forward decisions | Confirms Markdown source discovery/applicability are complete and runtime execution remained deferred. |
| `projects/ultraplan-go/sprints/06-study-run-execution/review.md` | Carry-forward decisions | Confirms report validation/path/rating behavior is complete and runtime execution remained deferred. |
| `projects/ultraplan-go/sprints/07-prompt-composition-generation/review.md` | Carry-forward decisions | Confirms prompt composition and preview are runtime-free and ready to feed later execution. |
| `projects/ultraplan-go/sprints/08-run-state-persistence/review.md` | Carry-forward decisions | Confirms durable task state exists, runtime execution remains deferred, and runtime sprint should align metadata with existing state placeholders. |

## Review Expectations

| What | How Verified |
| --- | --- |
| Scope alignment | Review `requirements.md`, `sprint-index.md`, `technical-handbook.md`, `reasoning.md`, and `plan.md` against runtime/provider integration scope and confirm no public study run execution, worker pool, run-loop mutation, summary, code extraction, target, or sprint-execution scope was added. |
| Required outputs exist | Check each path listed in Required Outputs and confirm implementation files are present or intentionally replaced by equivalent files recorded in `review.md`. |
| Architecture conformance | Review implementation paths and imports to confirm generic runtime behavior remains in `internal/platform/runtime`, CLI health wiring remains in `internal/app`, platform packages do not import product modules, and no direct OpenCode process supervision was introduced. |
| Agentwrap composition | Inspect code and tests proving the default runtime stack uses agentwrap `ObservingRuntime`, `ValidatingRuntime`, `PolicyRunner`, and `opencode.Runtime` in the required order or documents a reviewed exception. |
| Config and redaction behavior | Run or inspect tests proving runtime config maps to agentwrap/opencode options, invalid config fails before launch, and sensitive environment/config values are redacted. |
| Health behavior | Run or inspect tests proving runtime health checks use agentwrap health/capability APIs and `ultraplan health` reports actionable success/failure without study workflow execution. |
| Validation behavior | Run or inspect fake-runtime tests proving validation failures fail the logical run even when the underlying runtime succeeds. |
| Policy behavior | Run or inspect fake-runtime tests proving retry, rate-limit, fallback, and exhaustion are delegated to agentwrap policy and surfaced in metadata/events. |
| Permission behavior | Run or inspect tests proving permission policies map through agentwrap structures and unsupported required policy features fail before launch. |
| Event and diagnostics behavior | Run or inspect tests proving canonical events, unknown native extensions, safe raw-payload omission, usage/cost metadata, and `SDKError` classifications are preserved without unsafe logging. |
| Runtime-free default tests | Run `go test ./...` without OpenCode/provider environment and confirm all runtime tests use fakes or skip gated smoke paths by default. |
| Verification commands | Run `go test ./...` and `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go` and record the results in `review.md`. |
