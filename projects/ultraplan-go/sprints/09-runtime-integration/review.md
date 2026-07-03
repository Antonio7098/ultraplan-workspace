# Sprint Review: 09 Runtime Integration

> Project: `ultraplan-go`
> Sprint: `09-runtime-integration`
> Review Date: `2026-06-01`
> Status: `accepted with follow-ups`

## Review Inputs

- `sprint-index.md` — selected contracts, evidence reports, reasoning templates, excluded context
- `technical-handbook.md` — distilled patterns, trade-offs, anti-patterns, design pressures
- `reasoning.md` and `reasoning/runtime-integration.md` — sprint decisions and trade-offs
- `plan.md` — ordered task execution
- Architecture Review Protocol and Sprint Review Protocol

---

## Decision Conformance

| Decision From `reasoning.md` | Implemented? | Evidence | Notes |
| --- | --- | --- | --- |
| Decision 1: Generic Platform Runtime Boundary | yes | `internal/platform/runtime/runtime.go:15` — Request/Result/Event types; no product imports via `go list -deps` | Thin mapping layer stays generic |
| Decision 2: Agentwrap/OpenCode Composition | yes | `internal/platform/runtime/opencode.go:41` — `ObservingRuntime -> ValidatingRuntime -> PolicyRunner -> opencode.Runtime` | Required wrapper order preserved |
| Decision 3: Minimal Runtime Config Mapping, Validation, And Redaction | yes | `internal/platform/config/config.go:262` — Validate rejects unknown health/capability names, empty executable, invalid timeouts; `redaction.go` — Sensitive/RedactValue | Centralized preflight validation |
| Decision 4: Runtime Health Through Platform Boundary Only | yes | `internal/app/health_commands.go:61` — Appends runtime checks after config validation; no study workflow invocation | Thin health query path |
| Decision 5: Policy, Fallback, Sandbox, And Permission Mapping | yes | `agentwrap.go:97-103` — Path rules fail unless best_effort; `opencode.go:20-24` — BasicPolicy with MaxAttempts; `runtime_test.go:66-85` | Fail-fast permission enforcement |
| Decision 6: Event, Metadata, Diagnostics, Usage, And Error Mapping | yes | `runtime.go:267-284` — mapEvent preserves Kind via `string(event.Kind())`; `mapUsage` preserves Known flags; `mapSDKError` uses `errors.As` not string matching; `runtime_test.go` covers every canonical event kind | Canonical event kinds preserved |
| Decision 7: Deterministic Fake-First Verification | yes | `go test ./...` passes without OpenCode/credentials; `runtime_test.go` uses fakeRuntime/fakeRun | Fake-first default confirmed |

---

## Contract Conformance

| Contract / Requirement | Satisfied? | Evidence | Deviations |
| --- | --- | --- | --- |
| **Architecture** | pass | No product imports in `internal/platform/runtime`; wrapper order explicit; thin health command | None |
| **Errors** | pass | SDKError preserved via `errors.As` (`runtime.go:304-312`); cause chains via `classedError.Unwrap()`; no string matching | None |
| **Configuration** | partial | CFG-COMPAT-001: version rejection only, no migration path; CFG-ENV-001: no explicit environment/production guard for debug log level | Documented in follow-ups |
| **Observability** | partial | OBS-CORR-001: RunID/SessionID preserved but LLM prompt reference not in correlation metadata; OBS-DIAG-001: no canonical unified diagnostics surface; OBS-ALERT-001: no metrics/alerting integration. Runtime event-kind mapping now has exhaustive canonical table coverage. | Documented in follow-ups |
| **Security** | pass | Permission path rules fail-closed; redaction applied to env/model/diagnostics; no hardcoded secrets; no unsafe deserialization | None |
| **Testing** | partial | Runtime fake tests now cover success, SDK error classification, cancellation, timeout, malformed event, retry/fallback metadata, validation metadata, permission unsupported behavior, usage knownness, raw omission facts, and all canonical event kinds. TEST-SMOKE-002 real-dependency smoke and TEST-E2E-001 dedicated e2e suite remain deferred. | Documented in follow-ups |
| **CLI Surface** | pass | Stable exit codes (ExitRuntime, ExitValidation); JSON/text output separation; --json flag; --help docs | None |
| **LLM Runtime** | partial | LLM-PROMPT-001: prompts lack explicit version identifiers and effective version not propagated into runtime metadata | Documented in follow-ups |
| **LLM Evaluation / Cost / Safety** | pass | EVAL-SAFETY-001: permission guardrails with fail-before-launch; EVAL-COST-001: Usage with Known flags, EstimatedCost, latency via StartedAt/FinishedAt | None |
| **Workflows** | pass | WF-BOUNDARY-001: engine encapsulated behind Adapter; WF-STATE-001: Result exposes explicit terminal state; WF-RETRY-001: PolicySummary exposes retry/fallback decisions | WF-SCOPE-001, WF-IDEMPOTENCY-001, WF-COMP-001, WF-VERSION-001 not_triggered (infrastructure only) |

---

## Plan Execution

| Plan Task | Status | Evidence |
| --- | --- | --- |
| Task 1: Agentwrap Dependency And Runtime Package Baseline | done | `go.mod` declares agentwrap; `internal/platform/runtime/*.go` present |
| Task 2: Define Generic Runtime Request, Result, Health, And Diagnostic Types | done | `runtime.go:15-140` — Request, Result, Event, Artifact, Usage, Error types |
| Task 3: Extend Runtime Config, Validation, And Redaction | done | `config.go:262-318` — Validate; `redaction.go` — Sensitive/RedactValue |
| Task 4: Construct OpenCode Adapter And Wrapper Stack | done | `opencode.go:41-50` — ObservingRuntime -> ValidatingRuntime -> PolicyRunner -> opencode.Runtime |
| Task 5: Implement Health And Capability Mapping Plus CLI Health Wiring | done | `health_commands.go:61` — runtime health appended; `health.go:43-103` — capability/health checks |
| Task 6: Implement Policy, Fallback, Sandbox, And Permission Mapping | done | `agentwrap.go:75-103` — mapPermissionPolicy with fail-closed path rules |
| Task 7: Implement Event, Error, Metadata, And Safe Diagnostic Mapping | done | `runtime.go:267-331` — mapEvent, mapUsage, mapSDKError; `events.go` — canonical kinds and raw omission facts |
| Task 8: Add Fake-First Runtime, Config, And Health Tests | done | `runtime_test.go`, `config_test.go`, and app health tests cover fake-first runtime/config/health behavior; real OpenCode smoke remains optional and deferred |
| Task 9: Verify, Review, And Record Handoff Evidence | done | `go test ./...` and `go build ./cmd/ultraplan` pass |

---

## Verification Evidence

| Check | Command / Method | Result | Notes |
| --- | --- | --- | --- |
| Full test suite | `cd /home/antonioborgerees/coding/ultraplan-go && go test ./...` | pass | Rerun after follow-up tests on 2026-06-01T13:28:08+01:00; no OpenCode/credentials required |
| CLI build | `cd /home/antonioborgerees/coding/ultraplan-go && go build ./cmd/ultraplan` | pass | Rerun after follow-up tests on 2026-06-01T13:28:08+01:00; builds successfully |
| Import boundary | `go list -deps ./internal/platform/runtime` | pass, no product imports | Verified no `internal/study`, `internal/codeextract`, `internal/app` imports |
| Direct OpenCode guard | `rg "exec\.Command\|opencode run" internal/platform/runtime internal/app` | reviewed | No UltraPlan-owned OpenCode process launch; matches are `opencode.WithStderrLimit` option and app IO fields |
| SDKError preservation | `runtime_test.go:49-63` — `TestAdapterPreservesSDKErrorClassification` | pass | `errors.As` used, not string matching |
| Permission path rules | `runtime_test.go:66-85` — `TestPermissionPathRulesFailUnlessBestEffort` | pass | Fail-fast without best_effort; pass with best_effort |
| Usage knownness | `runtime_test.go:21-46` — `TestAdapterStartRunMapsEventsUsageAndError` | pass | `InputTokensKnown` preserved as true when non-nil |
| Canonical event kinds | `runtime_test.go` — `TestAdapterMapsAllCanonicalEventKinds` | pass | All Agentwrap canonical event kinds preserve kind, type, and payload facts |
| Runtime failure matrix | `runtime_test.go` — cancellation, timeout, malformed event, retry/fallback metadata, validation metadata tests | pass | Fake-first coverage expanded after initial review |
| Health command | `app_test.go:212-224` — `TestRunHealthRuntimeFailureExitsRuntime` | pass | ExitRuntime on runtime failure |

---

## Deviations

| Deviation | Reason | Risk | Follow-Up |
| --- | --- | --- | --- |
| Used local `replace github.com/Antonio7098/agentwrap => ../agentwrap` | Sprint reasoning allowed sibling-checkout development | Replace must be removed before module release | Use proper versioned agentwrap dependency before publishing |
| Configuration contract partial: no config migration path | Only version rejection exists; migration not implemented for this sprint | Users with old config versions get rejection without upgrade path | Add migration or clear rejection behavior for version mismatches |
| No explicit environment/production concept | Debug log level accepted without production safety guard | Unsafe dev settings could pass in production | Consider adding ULTRAPLAN_ENV with production guards for unsafe settings |
| LLM prompt version not propagated to runtime metadata | Prompt versioning infrastructure not in scope for this sprint | Runtime events cannot correlate to specific prompt revision | Add prompt version tracking and runtime metadata propagation in future execution sprint |
| No canonical unified diagnostics surface | Health and config commands exist separately; no single operator path | Operators must check multiple commands for full state | Consider consolidated diagnostics command in future |
| No metrics/alerting integration | Observability signals use structured logs but no external metrics export | High-severity signals lack machine-usable alert streams | Define alert codes and structured event export when metrics integration is needed |

---

## Follow-Ups Completed After Initial Review

| Follow-Up | Status | Evidence |
| --- | --- | --- |
| Expand runtime fake tests for cancellation, timeout, malformed event, retry/fallback metadata, validation metadata, permission failure, and safe diagnostics | done | `/home/antonioborgerees/coding/ultraplan-go/internal/platform/runtime/runtime_test.go` now includes fake-first coverage for cancellation/timeout failures, malformed native extension events, policy retry/fallback metadata, validation metadata, permission unsupported behavior, SDK classifications, raw omission facts, and usage knownness. |
| Add exhaustive canonical event kind table coverage | done | `TestAdapterMapsAllCanonicalEventKinds` covers lifecycle, session, message, progress, tool, artifact, permission, blocking, usage, warning, fatal_error, rate_limit, validation, retry, fallback, final_result, and native_extension. |

---

## Technical Handbook Conformance

All 10 patterns from `technical-handbook.md` are followed:
- Thin Runtime Boundary Behind Thin Commands — `health_commands.go:61` delegates without owning runtime
- Manual Composition Root With Explicit Wrapper Stack — `opencode.go:41-50` visible composition
- Decorator/Wrapper For Cross-Cutting Runtime Concerns — wrappers via composition, not duplication
- Merged Config Then Validate Once — `config.go:94-113` merge then validate
- Typed Error Classification Over String Matching — `runtime.go:304-312` uses `errors.As`
- Injectable IO And Fake Runtime Seams — `runtime_test.go` fakeRuntime/fakeRun
- Context As Lifecycle, Not Service Locator — ctx propagated through all runtime calls
- Structured Events And Logs For Operator Diagnostics — canonical event kinds preserved
- Permission And Trust Boundaries As First-Class Runtime Inputs — fail-closed path rules
- Behavioral Tests Over Implementation Detail Tests — assertions on observable behavior only

All 9 anti-patterns avoided:
- No runtime supervisor reimplementation — OpenCode via agentwrap options only
- No string matching for error classification
- No global runtime singletons
- No CLI health invoking study workflows
- No dropped unknown events or unsafe payload facts
- No unknown usage/cost converted to zero
- No context.Background in work paths
- No private wrapper implementation detail assertions
- No silent permission best-effort when requirements are strict

---

## Review Protocol Results

| Protocol | Applied? | Findings |
| --- | --- | --- |
| Architecture Review Protocol | yes | Runtime package has no product imports; wrapper order explicit; health remains thin query path |
| Sprint Review Protocol | yes | All tasks completed or explicitly deferred; verification commands pass; deviation documented |

---

## Final Assessment

- **Status:** `accepted with minor deferred follow-ups`
- **Residual Risks:**
  - Config migration path gap (CFG-COMPAT-001)
  - No explicit environment/production guard for unsafe settings (CFG-ENV-001)
  - LLM prompt version not propagated to runtime metadata (LLM-PROMPT-001)
  - No canonical unified diagnostics surface (OBS-DIAG-001)
  - No metrics/alerting integration (OBS-ALERT-001)
  - Optional real-dependency smoke tests and e2e suite gaps (TEST-SMOKE-002, TEST-E2E-001)
  - Local `replace` directive for agentwrap must be resolved before release
- **Required Follow-Ups:**
  - Implement config migration path or documented rejection behavior for version mismatches
  - Add explicit environment concept with production guards for unsafe debug settings
  - Add prompt version identifiers and effective version propagation into runtime metadata
  - Add canonical unified diagnostics surface (or document existing multi-command path)
  - Define stable alert codes and structured event export for metrics/alerting integration
  - Consider dedicated e2e test suite and optional real-dependency smoke tests for health command
  - Remove local `replace` directive and use proper versioned agentwrap before module release
