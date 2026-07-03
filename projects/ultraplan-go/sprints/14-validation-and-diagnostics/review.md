# Sprint Review: Validation Command, Diagnostics, and JSON Stability

> Project: `ultraplan-go`
> Sprint: `14-validation-and-diagnostics`
> Review Date: `2026-06-13`
> Status: `accepted`

## Review Inputs

- `projects/ultraplan-go/sprints/14-validation-and-diagnostics/sprint-index.md`
- `projects/ultraplan-go/sprints/14-validation-and-diagnostics/technical-handbook.md`
- `projects/ultraplan-go/sprints/14-validation-and-diagnostics/reasoning.md`
- `projects/ultraplan-go/sprints/14-validation-and-diagnostics/plan.md`
- `projects/ultraplan-go/sprints/14-validation-and-diagnostics/requirements.md`
- `system/protocols/architecture-review-protocol.md`
- `system/protocols/sprint-review-protocol.md`
- Implementation in `/home/antonioborgerees/coding/ultraplan-go`

## Decision Conformance

| Decision From `reasoning.md` | Implemented? | Evidence | Notes |
| --- | --- | --- | --- |
| Decision 1: Study-Owned Validation Result Model | yes | `internal/study/domain.go:112-123` (ValidationCheck, ValidationResult, StudyValidationResult); `internal/study/validation_command.go:20-108` | Typed result structures with schema_version, checks, severity, path, expected, observed, guidance. Study-owned. |
| Decision 2: Runtime-Free Study Validation Command | yes | `internal/app/study_commands.go:655-687`; `internal/study/validation_command.go:12-18`; tests `study_validate_commands_test.go:77-94` assert studyRuntimeFactory is never called | Validate/status never invoke runtime, network, source cloning, prompts, synthesis, repair, retry, fallback, or subprocess. |
| Decision 3: Stable JSON Envelope In `internal/app` | yes | `internal/app/json_output.go:9-32` (jsonEnvelope with schema_version=1, command, workspace, status, generated_at, result); payload-agnostic | Shared envelope used by all three --json surfaces. No study/runtime imports in envelope. |
| Decision 4: Status JSON From Durable Local State | yes | `internal/app/study_commands.go:755-889` (statusJSON, statusCounts, statusTaskJSON, knownValue[T]); `internal/app/study_status_commands_test.go:51-135` | Runtime-free. run_id, counts, tasks, locks, run_metadata, validation summaries, retry/fallback metadata, unknown usage/cost as known=false. |
| Decision 5: Health JSON Through Existing Runtime Health Boundary | yes | `internal/app/health_commands.go:30-110`; `internal/app/app_test.go:254-261` (stubRuntimeHealth replaces runtimeHealthChecks) | Uses existing platform runtime/agentwrap health boundary. Platform runtime remains product-agnostic. |
| Decision 6: Redaction Before Text Or JSON Rendering | yes | `internal/platform/config/redaction.go:21-41` (Sensitive, RedactValue); `internal/app/study_commands.go:729-753` (redactedStudyValidationResult); `internal/app/health_commands.go:158-174` (sanitizeHealthChecks) | Applied before both text and JSON rendering. Tests seed sk-test/Bearer values and assert absence. |
| Decision 7: Fake-First Schema And Behavior Tests | yes | `internal/study/validation_command_test.go`; `internal/app/study_validate_commands_test.go`; `internal/app/study_status_commands_test.go`; `internal/app/app_test.go:225-252` | Deterministic offline tests. JSON schema assertions parse and verify envelope/result fields. |

## Contract Conformance

| Contract | Satisfied? | Evidence | Deviations |
| --- | --- | --- | --- |
| Architecture | yes | Study validates in `internal/study`; app wires/renders in `internal/app`; platform/runtime untouched; no forbidden global packages; dependency direction app->study->platform preserved. Import grep confirms no reverse edges. | none |
| Errors | yes | Exit code 0/5 mapped correctly. Cause chains preserved via %w. Typed ValidationCheck with name/severity/path/expected/observed/guidance. ERR-RETRY-001 and ERR-TASK-001 not triggered (inspection-only). | none |
| Configuration | yes | Health JSON uses merged config from loadEffectiveConfig. RedactValue applied to all check messages. CFG-COMPAT-001 and CFG-ENV-001 not triggered. | none |
| Observability | yes | Stable versioned envelope for all three surfaces. Status exposes counts, task states, locks, run metadata, validation summaries, retry/fallback, unknown usage/cost. Health exposes stable check IDs. All ANSI-free. | none |
| Security | yes | Secrets never exposed. Redaction applied before rendering. Paths workspace-relative. Tests seed secret-like values and assert absence. SEC-AUTHN-001, SEC-AUTHZ-001, SEC-INJECT-001, SEC-DEPS-001, SEC-NET-001 not triggered. | none |
| Testing | yes | Deterministic offline tests. Fake-first command tests with replaceable runtime seams. JSON schema assertions. go test ./..., go test -race ./..., go build ./cmd/ultraplan all pass. TEST-SMOKE-002 (real OpenCode) explicitly deferred. | none |
| Documentation | yes | Help text updated for study validate, validate --json, status --json, health --json. JSON envelope schema versioned. DOC-OPS-001, DOC-EXAMPLE-001, DOC-AGENT-001 not triggered. | Minor: no dedicated test for top-level study --help mentioning both --json variants. |
| CLI Surface | yes | Command shape correct. --json flags with stable envelope. Exit codes 0/5/6 mapped. Text mode backward-compatible. JSON ANSI-free with no mixed human text. --workspace global. | none |
| LLM Runtime | yes | Platform/runtime stays product-agnostic (only imports agentwrap + config). Validate/status runtime-free (tests prove no invocation). Health routes through existing Adapter.Health boundary. | none |
| LLM Evaluation / Cost / Safety | yes | EVAL-COST-001 satisfied: knownValue[T] preserves unknown usage/cost as Known=false. Tests assert known=false when absent. No prompt/model/tool changes. | none |
| Workflows | yes | Status exposes run-loop counts, task state summaries, retry/fallback/cancellation metadata, lock diagnostics, validation summaries, run metadata. Runtime-free. Schema versioning exposed. WF-IDEMPOTENCY-001, WF-COMP-001 explicitly deferred. | none |
| Performance | yes | Validation bounded to known workspace/study artifacts. No recursive source scans. No new caches or concurrency. Validate/status runtime-free (no eager runtime setup). | none |
| Persistence And Migrations | yes | Strict run-state.json schema handling. Typed sentinels: ErrRunStateMissing/Malformed/Unsupported. Distinct diagnostics. No silent migration. Lock diagnostics surfaced. PERSIST-RECOVERY-001 explicitly deferred. | none |

## Plan Execution

| Plan Task | Status | Evidence |
| --- | --- | --- |
| Task 1: Baseline Current Implementation | done | All inspection steps completed before implementation. |
| Task 2: Add Study Validation Result Model | done | `internal/study/domain.go:112-123` with ValidationCheck, ValidationResult, StudyValidationResult, ValidationCounts. |
| Task 3: Implement Study Validation Service | done | `internal/study/validation_command.go:20-108` with bounded orchestration over study artifacts. |
| Task 4: Wire Validate Command And Help | done | `internal/app/study_commands.go:655-697` (runStudyValidate + studyValidateHelp). |
| Task 5: Add Shared JSON Envelope Renderer | done | `internal/app/json_output.go:9-32` (payload-agnostic envelope). |
| Task 6: Extend Status JSON | done | `internal/app/study_commands.go:755-889` (statusJSON with counts, tasks, locks, run metadata, validation summaries, retry/fallback, knownValue[T]). |
| Task 7: Extend Health JSON | done | `internal/app/health_commands.go:30-110` (health --json with stable check IDs). |
| Task 8: Extend Redaction Policy | done | `internal/platform/config/redaction.go:21-41` (Sensitive, RedactValue with token-like markers). |
| Task 9: Add Validation Service Tests | done | `internal/study/validation_command_test.go` covers valid, missing, invalid, malformed, inapplicable, run-state diagnostics. |
| Task 10: Add Validate Command Tests | done | `internal/app/study_validate_commands_test.go` covers text/JSON, exit codes, redaction, no runtime invocation. |
| Task 11: Extend Status And Health Tests | done | `internal/app/study_status_commands_test.go` and `internal/app/app_test.go:225-252` cover status JSON shape and health JSON. |
| Task 12: Run Verification And Prepare Review | done | go test ./..., go test -race ./..., go build ./cmd/ultraplan all passed. Review recorded. |

## Verification Evidence

| Check | Command / Method | Result | Notes |
| --- | --- | --- | --- |
| Offline unit and command tests | `go test ./...` | pass | All packages green. No OpenCode, network, or long sleeps. |
| Race verification | `go test -race ./...` | pass | No data races detected. |
| CLI build | `go build ./cmd/ultraplan` | pass | Builds successfully. |

## Applicability Summary

| Requirement Category | Direct | Partial | Not Triggered | Deferred |
| --- | --- | --- | --- | --- |
| Architecture | 3 | 4 | 0 | 0 |
| Errors | 7 | 2 | 2 | 0 |
| Configuration | 3 | 2 | 2 | 0 |
| Observability | 4 | 4 | 0 | 0 |
| Security | 2 | 3 | 4 | 0 |
| Testing | 7 | 3 | 0 | 1 |
| Documentation | 1 | 2 | 4 | 0 |
| CLI Surface | 13 | 1 | 1 | 0 |
| LLM Runtime | 3 | 2 | 8 | 0 |
| LLM Eval / Cost / Safety | 1 | 0 | 6 | 0 |
| Workflows | 2 | 2 | 0 | 2 |
| Performance | 1 | 2 | 4 | 0 |
| Persistence And Migrations | 3 | 2 | 3 | 1 |
| Technical Handbook | 6 patterns + 7 antipatterns all followed/avoided | 0 | 0 | 0 |

## Deviations

| Deviation | Reason | Risk | Follow-Up |
| --- | --- | --- | --- |
| Health JSON result includes workspace root in result payload | Existing health text already exposes workspace context; diagnostics inside checks prefer relative/safe strings. | Low -- workspace root is not a secret. | none required |
| Status JSON shaped in `internal/app` from durable `RunState` | No separate study status DTO introduced; existing study state model exposes needed durable data. | Low -- app layer is correct owner of rendering shape. | none required |
| No dedicated health_commands_test.go file | Health JSON tests are in app_test.go alongside other health tests. | Low -- tests exist and cover the surface. | Consider dedicated file if test density grows. |
| No top-level study --help test for both --json variants | Help text is correct per study_commands.go:132-133 but not directly tested at top-level. | Low -- subcommand help is tested. | Optional: add TestStudyTopLevelHelpMentionsValidateAndStatusJSON. |

## Technical Handbook Conformance

| Pattern | Status | Evidence |
| --- | --- | --- |
| Thin command wiring, study-owned behavior | followed | `internal/app/study_commands.go:655-687` |
| Stable machine output as a separate rendering path | followed | `internal/app/json_output.go:1-32` |
| Typed diagnostics with preserved error chains | followed | `internal/study/domain.go:112-123`, `internal/study/validation_command.go:139-172` |
| Redaction as a type/policy boundary before logs or JSON | followed | `internal/platform/config/redaction.go:21-41`, `internal/app/study_commands.go:745-753` |
| Fake-first inspection commands | followed | `internal/app/study_validate_commands_test.go:77-94`, `internal/app/study_status_commands_test.go:51-135` |
| Bounded state/artifact inspection | followed | `internal/study/validation_command.go:20-108` |

| Anti-Pattern | Status | Evidence |
| --- | --- | --- |
| Business logic in RunE | avoided | `internal/app/study_commands.go:655-687` (thin wiring) |
| Direct os.Stdout/ANSI/progress in machine output | avoided | `internal/app/json_output.go:20-32`, tests assert no ANSI |
| Config/runtime globals hidden behind validation | avoided | `internal/app/study_commands.go:20-22` (replaceable seam), tests prove no invocation |
| Panic or string-only parse failures | avoided | `internal/study/validation_command.go:139-172` (typed findings) |
| Secret leakage through logs, JSON, or observed details | avoided | Tests seed sk-test/Bearer values and assert absence |
| Unbounded recursive scans or eager expensive setup | avoided | Validation bounded to known artifacts, no filepath.WalkDir |
| Implementation-detail tests instead of behavior/schema tests | avoided | Tests assert JSON field names, envelope shape, exit codes |

## Review Protocol Results

| Protocol | Applied? | Findings |
| --- | --- | --- |
| Architecture Review | yes | All 8 architecture requirements satisfied. Study validation owned by internal/study. CLI wiring/rendering in internal/app. Platform runtime product-agnostic. No forbidden global packages. Dependency direction preserved. |
| Sprint Review | yes | All 13 contracts satisfied. All 7 reasoning decisions implemented. All 12 plan tasks complete. Verification commands pass. |

## Final Assessment

- **Status:** `accepted`
- **Residual Risks:** Safe diagnostics may omit details useful for manual debugging (carried forward per reasoning.md); unsupported run-state versions produce guidance but no automatic migration (by design); JSON schema versioning may need finer granularity if payloads evolve at different rates.
- **Required Follow-Ups:** None blocking. Optional low-priority items: add TestStudyTopLevelHelpMentionsValidateAndStatusJSON; add service-level secret-redaction assertion in validation_command_test.go; consider dedicated health_commands_test.go if test density grows; extract inline status JSON types from study_commands.go into named types for discoverability.
