# Sprint Plan: 09 Runtime Integration

> Project: `ultraplan-go`
> Sprint: `09-runtime-integration`
> Source: `reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/09-runtime-integration/requirements.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/sprints/09-runtime-integration/sprint-index.md`, `projects/ultraplan-go/sprints/09-runtime-integration/technical-handbook.md`, `projects/ultraplan-go/sprints/09-runtime-integration/reasoning/runtime-integration.md`, `projects/ultraplan-go/sprints/09-runtime-integration/reasoning.md`, `templates/sprint-plan.md`, `/home/antonioborgerees/coding/agentwrap/docs/README.md`, `/home/antonioborgerees/coding/agentwrap/docs/architecture.md`, `/home/antonioborgerees/coding/agentwrap/docs/core-api.md`, `/home/antonioborgerees/coding/agentwrap/docs/events-and-metadata.md`, `/home/antonioborgerees/coding/agentwrap/docs/errors-config-health.md`, `/home/antonioborgerees/coding/agentwrap/docs/wrappers.md`, `/home/antonioborgerees/coding/agentwrap/docs/opencode-adapter.md`, `/home/antonioborgerees/coding/agentwrap/docs/integration-guide.md`, `/home/antonioborgerees/coding/agentwrap/runtime.go`, `/home/antonioborgerees/coding/agentwrap/events.go`, `/home/antonioborgerees/coding/agentwrap/metadata.go`, `/home/antonioborgerees/coding/agentwrap/config.go`, `/home/antonioborgerees/coding/agentwrap/health.go`, `/home/antonioborgerees/coding/agentwrap/policy.go`, `/home/antonioborgerees/coding/agentwrap/permissions.go`, `/home/antonioborgerees/coding/agentwrap/validation.go`, `/home/antonioborgerees/coding/agentwrap/observability.go`, `/home/antonioborgerees/coding/agentwrap/errors.go`, `/home/antonioborgerees/coding/agentwrap/opencode/options.go`, `/home/antonioborgerees/coding/ultraplan-go/go.mod`, `/home/antonioborgerees/coding/ultraplan-go/internal/platform/config/config.go`, `/home/antonioborgerees/coding/ultraplan-go/internal/platform/config/config_test.go`, `/home/antonioborgerees/coding/ultraplan-go/internal/platform/config/redaction.go`, `/home/antonioborgerees/coding/ultraplan-go/internal/platform/runtime/doc.go`, `/home/antonioborgerees/coding/ultraplan-go/internal/app/app.go`, `/home/antonioborgerees/coding/ultraplan-go/internal/app/app_test.go`, `/home/antonioborgerees/coding/ultraplan-go/internal/app/health_commands.go`

This plan executes `reasoning.md`. It must not invent architecture, scope, or decisions.

## Reasoning Source

- **Sprint Reasoning:** `projects/ultraplan-go/sprints/09-runtime-integration/reasoning.md`
- **Sprint Index:** `projects/ultraplan-go/sprints/09-runtime-integration/sprint-index.md`
- **Technical Handbook:** `projects/ultraplan-go/sprints/09-runtime-integration/technical-handbook.md`
- **Area Reasoning:** `projects/ultraplan-go/sprints/09-runtime-integration/reasoning/runtime-integration.md`
- **Supplemental Runtime Reference:** `/home/antonioborgerees/coding/agentwrap/docs` and public Agentwrap source files listed in the inputs line.

## Sprint Status

- **Status:** `complete with documented follow-ups`
- **Owner:** `implementation agent`
- **Start Date:** `2026-06-01`
- **Completion Date:** `2026-06-01`

## Decisions To Execute

| Decision | Source Section | Execution Implication |
| --- | --- | --- |
| Thin generic platform runtime boundary | `reasoning.md#decision-1-thin-generic-platform-runtime-boundary` | Add `internal/platform/runtime` types and mapping helpers without importing `internal/study`, `internal/codeextract`, or other product modules. |
| Agentwrap/OpenCode composition and ownership | `reasoning.md#decision-2-agentwrapopencode-composition-and-ownership` | Use `github.com/Antonio7098/agentwrap` and `agentwrap/opencode`; compose `ObservingRuntime -> ValidatingRuntime -> PolicyRunner -> opencode.Runtime`; do not launch or decode OpenCode directly. |
| Minimal runtime config mapping, validation, and redaction | `reasoning.md#decision-3-minimal-runtime-config-mapping-validation-and-redaction` | Extend existing config minimally, validate before launch/health, and redact sensitive runtime values. |
| Runtime health through platform boundary only | `reasoning.md#decision-4-runtime-health-through-platform-boundary-only` | Replace the current `runtime.opencode: skipped` health placeholder with platform runtime health checks that do not run studies or prompts. |
| Policy, fallback, sandbox, and permission mapping | `reasoning.md#decision-5-policy-fallback-sandbox-and-permission-mapping` | Map retries, fallback provider/model, sandbox, permission mode, and structured permissions into Agentwrap policy and permission types. |
| Event, metadata, diagnostics, usage, and error mapping | `reasoning.md#decision-6-event-metadata-diagnostics-usage-and-error-mapping` | Preserve canonical event kinds, raw payload omission facts, nil usage semantics, safe metadata, and `SDKError` classifications. |
| Deterministic fake-first verification | `reasoning.md#decision-7-deterministic-fake-first-verification` | Cover normal behavior through fakes; any real OpenCode smoke path is opt-in and skipped by default. |

## Requirements / Contracts To Satisfy

| Contract / Requirement ID | Required Behavior | Evidence Planned |
| --- | --- | --- |
| Architecture, AC-02, AC-03 | Runtime logic stays in `internal/platform/runtime`; no product imports. | Import review; runtime package tests; Architecture Review Protocol. |
| LLM Runtime, AC-01, AC-04, AC-05, AC-11, AC-12, AC-13 | Use Agentwrap root APIs, `agentwrap/opencode`, validation, policy, and required wrapper order. | `go.mod`; wrapper composition tests; adapter construction tests; no direct `opencode` process code. |
| Configuration, AC-06, AC-07, AC-08 | Runtime config maps supported fields, rejects invalid/unsupported values, and redacts secrets. | Config tests for defaults, env/CLI precedence where present, health/capability names, timeouts, retries, permissions, and redaction. |
| CLI Surface, AC-15, AC-16, AC-25 | `ultraplan health` reports runtime health without invoking study workflows. | Health command tests for text/JSON, failure exit status, redacted diagnostics, and no prompt/study/report calls. |
| Observability and Errors, AC-17 through AC-23 | Event kinds, raw payload omission, metadata, usage knownness, warnings, and `SDKError` classifications are preserved safely. | Event/result/error mapping tests and review of `events.go`. |
| Security, AC-08, AC-14, AC-19 | Sandbox and permissions map through Agentwrap; unsupported required policy features fail before launch; sensitive data is redacted. | Permission mapping tests; unsupported path-rule tests; diagnostics redaction tests. |
| Workflows and Testing, AC-09, AC-24, AC-26 through AC-30 | Context/cancellation propagate, retries/fallback delegate to Agentwrap, default tests are fake-only, build and full tests pass. | Fake runtime tests; `go test ./...`; `go build ./cmd/ultraplan`; optional smoke skip evidence if added. |

## Tasks

- [x] **Task 1: Add Agentwrap Dependency And Runtime Package Baseline**
  > Executes: `Decision 1`, `Decision 2`, `AC-01`, `AC-02`, `AC-03`
  - [x] Add `github.com/Antonio7098/agentwrap` to `/home/antonioborgerees/coding/ultraplan-go/go.mod`; use a local `replace` only if required for sibling-checkout development.
  - [x] Replace the placeholder runtime package doc with runtime package files listed by requirements: `runtime.go`, `agentwrap.go`, `opencode.go`, `health.go`, `policy.go`, and `events.go`.
  - [x] Keep public dependencies in `internal/platform/runtime` limited to standard library, platform-safe packages, Agentwrap root APIs, and `agentwrap/opencode`.
  - [x] Add compile-time interface assertions for internal fakes or adapters only where they clarify conformance without creating a public SDK.

- [x] **Task 2: Define Generic Runtime Request, Result, Health, And Diagnostic Types**
  > Executes: `Decision 1`, `Decision 6`, `AC-02`, `AC-10`, `AC-20`, `AC-21`, `AC-23`
  - [x] Define narrow UltraPlan internal types for request prompt, workdir, provider, model, timeout, metadata, validation, health requirements, capability requirements, permissions, sandbox, events, diagnostics, and result summaries.
  - [x] Include generic correlation metadata keys for run ID, task ID, task kind, study, dimension, source, source kind, output path, runtime, provider, and model without importing study types.
  - [x] Represent usage and estimated cost with known/unknown semantics so missing token/cost values are not converted to zero.
  - [x] Map Agentwrap `RunResult`, `RunMetadata`, `ArtifactRef`, warnings, validation, repair, cleanup, attempts, policy, permissions, and session facts into safe UltraPlan result fields.

- [x] **Task 3: Extend Runtime Config, Validation, And Redaction**
  > Executes: `Decision 3`, `Decision 5`, `AC-06`, `AC-07`, `AC-08`, `AC-14`
  - [x] Extend `internal/platform/config.Config` only for required runtime adapter options: executable, extra args, safe env additions, stderr/diagnostic limit, required health, required capabilities, sandbox, permission mode/policy posture, retry count, and fallback provider/model where not already represented.
  - [x] Keep existing precedence order: defaults, workspace config, environment, command flags where applicable.
  - [x] Map runtime config into Agentwrap concepts: `RunRequest.Provider`, `RunRequest.Model`, `Timeout`, `RequireHealth`, `RequireCaps`, `Sandbox`, `Permissions`, `PermissionPolicy`, `Validation`, and Agentwrap `EffectiveConfig` or config layers where practical.
  - [x] Validate known runtime names, non-empty executable, positive timeout, non-negative retries, supported health IDs, supported capability IDs, permission actions/tools, unsupported required path rules, and fallback provider/model shape before launch.
  - [x] Redact sensitive env/config/model/runtime values in config output, logs, health diagnostics, and error diagnostics.

- [x] **Task 4: Construct OpenCode Adapter And Wrapper Stack**
  > Executes: `Decision 2`, `Decision 7`, `AC-04`, `AC-05`, `AC-11`, `AC-12`, `AC-13`
  - [x] Build OpenCode only with `opencode.NewRuntime`, `opencode.WithExecutable`, `opencode.WithExtraArgs`, `opencode.WithEnv`, and `opencode.WithStderrLimit`.
  - [x] Compose production runtime as `agentwrap.ObservingRuntime{Runtime: agentwrap.ValidatingRuntime{Runtime: agentwrap.PolicyRunner{Runtime: opencodeRuntime}}}` with explicit options for store, sinks, policy, validation, and persistence policy.
  - [x] Use `agentwrap.NewMemoryRunStore` only for tests/local inspection in this sprint; do not introduce durable runtime records.
  - [x] Preserve caller context through `StartRun`, event draining, `Wait`, health checks, and cancellation.
  - [x] Ensure validation failure after runtime success fails the logical result through `agentwrap.ValidatingRuntime`.

- [x] **Task 5: Implement Health And Capability Mapping Plus CLI Health Wiring**
  > Executes: `Decision 4`, `AC-15`, `AC-16`, `AC-25`
  - [x] Map configured health names to Agentwrap health IDs: `runtime_available`, `structured_output`, `workdir`, `config`, `provider`, `model`, `authentication`, and `runtime_paths` as supported by this sprint's config.
  - [x] Map configured capability names to Agentwrap capabilities including `structured_events`, `raw_payloads`, `cancellation`, `artifacts`, `permissions`, `usage`, `validation_events`, and session capability variants.
  - [x] Call `HealthChecker.CheckHealth` when the configured runtime supports it; use `Capabilities` for required capability checks.
  - [x] Convert health reports to existing `healthCheck` text/JSON output shape with actionable, redacted messages and stable runtime failure exit behavior.
  - [x] Prove health does not invoke study workflows, prompt composition, schedulers, report validation, or direct OpenCode process code.

- [x] **Task 6: Implement Policy, Fallback, Sandbox, And Permission Mapping**
  > Executes: `Decision 5`, `AC-06`, `AC-13`, `AC-14`, `AC-20`, `AC-24`
  - [x] Map retry count and backoff defaults into `agentwrap.BasicPolicy` with explicit bounds.
  - [x] Map backup provider/model into `agentwrap.FallbackAlternative` request overrides while preserving prompt, workdir, timeout, sandbox, permissions, permission policy, validation, and metadata.
  - [x] Surface policy attempts, retry, wait, fallback, rate-limit, and exhaustion facts from `RunMetadata.Policy` and canonical events.
  - [x] Map permission mode and structured posture into `agentwrap.PermissionPolicy`, `PermissionTool`, and `PermissionAction`.
  - [x] Reject unsupported required permission features, especially path-level rules for the current OpenCode subprocess adapter, unless `PermissionUnsupportedBestEffort` is explicitly configured.

- [x] **Task 7: Implement Event, Error, Metadata, And Safe Diagnostic Mapping**
  > Executes: `Decision 6`, `AC-17`, `AC-18`, `AC-19`, `AC-20`, `AC-21`, `AC-22`, `AC-23`
  - [x] Preserve all canonical event kinds: lifecycle, session, message, progress, tool, artifact, permission, blocking, usage, warning, fatal_error, rate_limit, validation, retry, fallback, final_result, and native_extension.
  - [x] Treat `native_extension` and unsupported future native projections as non-fatal unless Agentwrap returns a terminal classified error.
  - [x] Preserve raw payload presence, source, encoding, safety, omission reason, and default unsafe-byte omission through Agentwrap observing records or equivalent safe mapping.
  - [x] Use `errors.As` or `agentwrap.ErrorAs` to preserve `SDKError` category, operation, user detail, provider, model, runtime kind, exit code/signal, retry-after, and redacted metadata.
  - [x] Keep stdout/stderr/native payloads out of logs unless explicitly safe and bounded by Agentwrap diagnostics.

- [x] **Task 8: Add Fake-First Runtime, Config, And Health Tests**
  > Executes: `Decision 7`, `AC-24` through `AC-30`
  - [x] Add `internal/platform/runtime/runtime_test.go` with Agentwrap-compatible fake runtime/run coverage for success, validation failure after runtime success, runtime failure, timeout, cancellation, malformed event, rate limit, retry, fallback, permission failure, event drain, safe diagnostics, and observability metadata.
  - [x] Extend `internal/platform/config/config_test.go` for runtime defaults, env/file/CLI precedence where applicable, validation failures, health/capability mapping, fallback mapping, permissions, and redaction.
  - [x] Add or extend `internal/app/health_commands_test.go` or existing app tests for runtime health success/failure, unsupported health requirement, missing executable/unavailable runtime, redacted diagnostics, JSON/text shape, and exit status.
  - [x] If a real OpenCode smoke test is added, gate it behind explicit environment variables and skip by default in `go test ./...`.

- [x] **Task 9: Verify, Review, And Record Handoff Evidence**
  > Executes: `Review Expectations`, `AC-28`, `AC-29`, `AC-30`
  - [x] Run `go test ./...` in `/home/antonioborgerees/coding/ultraplan-go`.
  - [x] Run `go build ./cmd/ultraplan` in `/home/antonioborgerees/coding/ultraplan-go`.
  - [x] Review imports and code paths for no product imports from `internal/platform/runtime`, no direct `opencode` process launch, no direct stdout/stderr native parsing, and no public study execution scope.
  - [x] Produce `review.md` using the selected Architecture Review and Sprint Review protocols after implementation.

## Evidence Checklist

- [x] `go.mod` declares Agentwrap and implementation uses Agentwrap root APIs plus `agentwrap/opencode`.
- [x] Runtime behavior lives in `internal/platform/runtime` with no product-module imports.
- [x] Wrapper composition is explicit and covered by tests.
- [x] OpenCode is constructed only through `agentwrap/opencode` options.
- [x] Config validation rejects invalid runtime, health, capability, timeout, retry, permission, and executable values before launch.
- [x] Health output reports runtime health/capability status without study workflow execution.
- [x] Events preserve canonical kind and safe raw payload omission facts.
- [x] Usage/cost unknown values remain unknown.
- [x] Runtime errors preserve `SDKError` classifications without string matching.
- [x] Default tests require no OpenCode, credentials, network, or real subprocesses.
- [x] Required verification commands pass or failures are recorded in `review.md`.

## Verification Commands

| Check | Command | Expected Result |
| --- | --- | --- |
| Full tests | `cd /home/antonioborgerees/coding/ultraplan-go && go test ./...` | Passes without OpenCode, provider credentials, network access, or real subprocess execution. |
| CLI build | `cd /home/antonioborgerees/coding/ultraplan-go && go build ./cmd/ultraplan` | Builds successfully. |
| Import boundary | `cd /home/antonioborgerees/coding/ultraplan-go && go list -deps ./internal/platform/runtime` | Does not include `ultraplan-go/internal/study`, `ultraplan-go/internal/codeextract`, or other product modules. |
| Direct OpenCode guard | `cd /home/antonioborgerees/coding/ultraplan-go && rg "exec\\.Command|opencode run|--format json|Stdout|Stderr" internal/platform/runtime internal/app` | Any matches are reviewed; OpenCode process/stdout/stderr handling must be inside Agentwrap, not UltraPlan. |

## Risks And Blockers

| Risk / Blocker | Source | Mitigation | Status |
| --- | --- | --- | --- |
| Agentwrap sibling checkout changes before implementation | `reasoning.md` confirmed assumption | Compiled against sibling checkout with local `replace`; `go test ./...` and `go build ./cmd/ultraplan` pass. | mitigated |
| Runtime facade grows into a competing SDK | `reasoning.md` risk | Runtime types are narrow and map directly to agentwrap request/result/health/event/error concepts. | mitigated |
| Health wiring accidentally starts execution | `reasoning.md` risk | Health path calls Agentwrap/OpenCode health and capability checks only; app tests stub runtime health and do not invoke study workflows. | mitigated |
| Permission expectations exceed OpenCode subprocess support | `reasoning.md` risk; Agentwrap OpenCode docs | Path rules fail before launch unless `best_effort` is configured; covered by runtime test. | mitigated |
| Event drain or wait paths leak goroutines | `reasoning.md` risk | Adapter drains `Run.Events()` and then calls `Wait`; fake test covers event drain/result mapping. | mitigated with follow-up coverage for cancellation stress |
| Raw diagnostics leak secrets | `reasoning.md` risk | Config/env redaction extended and raw event bytes are omitted with omission facts. | mitigated |

## Review Inputs

Review should use:

- `projects/ultraplan-go/sprints/09-runtime-integration/sprint-index.md`
- `projects/ultraplan-go/sprints/09-runtime-integration/technical-handbook.md`
- `projects/ultraplan-go/sprints/09-runtime-integration/reasoning/runtime-integration.md`
- `projects/ultraplan-go/sprints/09-runtime-integration/reasoning.md`
- this `plan.md`
- `/home/antonioborgerees/coding/agentwrap/docs`
- implementation diff in `/home/antonioborgerees/coding/ultraplan-go`
- verification evidence
- `system/protocols/architecture-review-protocol.md`
- `system/protocols/sprint-review-protocol.md`

## Execution Log

| Date / Step | Action | Evidence / Notes |
| --- | --- | --- |
| 2026-06-01 | Created implementation plan | Plan grounded in sprint reasoning, Agentwrap docs/source, and current UltraPlan Go app/config/runtime placeholders. |
| 2026-06-01 | Implemented runtime boundary and Agentwrap/OpenCode composition | Added `github.com/Antonio7098/agentwrap` with local sibling `replace`; added `internal/platform/runtime` request/result/health/policy/event mapping and OpenCode stack construction. |
| 2026-06-01 | Extended config and health wiring | Added runtime config fields, validation, redaction, text output, and `ultraplan health` runtime health/capability reporting with fake-first app tests. |
| 2026-06-01 | Verified implementation | `go test ./...` passed; `go build ./cmd/ultraplan` passed; `go list -deps ./internal/platform/runtime` showed no product-module imports; direct OpenCode guard found no UltraPlan-owned process launch/stdout parsing. |
| 2026-06-01 | Recorded review follow-ups | `review.md` documents accepted follow-ups for exhaustive event matrix coverage, broader fake runtime scenarios, optional smoke testing, and future public schema decisions. |
| 2026-06-01 | Completed appropriate follow-ups | Expanded fake runtime tests for canonical event kinds, cancellation, timeout, malformed event, retry/fallback metadata, validation metadata, permission unsupported behavior, safe raw omission, and SDK error classifications; reran `go test ./...`, `go build ./cmd/ultraplan`, import boundary, and direct OpenCode guard. |

## Completion Criteria

- [x] All implementation tasks are complete or explicitly deferred in `review.md`.
- [x] Verification commands were run or deferrals are documented.
- [x] Evidence satisfies the expectations from `reasoning.md`.
- [x] `review.md` can evaluate conformance without guessing intent.
