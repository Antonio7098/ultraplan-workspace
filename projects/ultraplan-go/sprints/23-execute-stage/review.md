# Sprint Review: 23-execute-stage

> Project: `ultraplan-go`
> Sprint: `23-execute-stage`
> Inputs used: `sprint-index.md`, `requirements.md`, `technical-handbook.md`, `reasoning.md`, `plan.md`, `../ultraplan-go` implementation, `system/contracts/core/*`, `system/contracts/runtime/*`, `system/contracts/surfaces/cli.md`, `system/protocols/architecture-review-protocol.md`, `system/protocols/sprint-review-protocol.md`

This review was produced by the Review Sprint Protocol. Independent subagents reviewed the implementation against every contract selected in `sprint-index.md` plus the technical handbook. Findings are synthesized below.

## 1. Sprint Status Summary

| Field | Value |
| --- | --- |
| Sprint Goal | Execute validated `plan.md` implementation tasks through the generic runtime boundary with durable task state, safe target-repository boundaries, resumability, and clear diagnostics. |
| Implementation State | All eleven plan tasks are marked `[x]` in `plan.md`; the implementation in `../ultraplan-go` matches the contracted surface for stage wiring, run state, target safety, prompts, model resolution, and summary generation. The consolidated `internal/sprint/execute_test.go` and `internal/app/sprint_execute_commands_test.go` files required by `requirements.md` lines 31-32 remain absent. |
| Verification Commands | `go test ./...` PASS · `go test -race ./...` PASS · `go build ./cmd/ultraplan` PASS |
| Selected Contracts | Architecture, CLI Surface, Configuration, Documentation, Errors, LLM Evaluation / Cost / Safety, LLM Runtime, Observability, Performance, Persistence And Migrations, Security, Testing, Workflows |
| Overall Outcome | **Accepted with follow-ups**. Core execute semantics, atomic `.run-state.json`, target containment, evidence-gated completion, stale-running recovery, prompt rendering, summary generation, flow integration, and CLI wiring for `execute`, `validate execute`, `prompt execute`, `flow --to execute`, `--task`, `--dry-run`, `--resume`, and `--model` are implemented and pass deterministic tests. Required Fixes from the prior review remain open: the two contract-required consolidated test files (`internal/sprint/execute_test.go`, `internal/app/sprint_execute_commands_test.go`) are still missing, the help/status text is stale, `redactPlanning` does not cover `ExecuteModel`/`ExecuteVariant`, `Service.executeModelSelection` silently falls back to a hard-coded literal instead of failing fast, `mapSprintError` does not classify the `ErrExecuteRunState*` sentinels, `renderSprintStatus` and `StatusSummary` do not surface execute readiness/progress, `PlanningStages()` does not include `StageExecute`, and `internal/platform/logging.Logger` is never instantiated. The unaddressed Required Fixes are bounded: they are documentation/test/redaction polish and do not block the architectural conformance that the sprint was about. Future Follow-ups F-1 through F-8 also remain. |

## 2. Verification

| Check | Command | Result |
| --- | --- | --- |
| Full test suite | `GOCACHE=/tmp/ultraplan-go-build-cache go test ./...` from `../ultraplan-go` | All packages pass (app cached, sprint 0.030s, study cached, platform/runtime cached, project cached, workspace cached, codeextract cached, platform/config cached, platform/logging cached) |
| Sprint/app/config focused | `go test ./internal/sprint ./internal/app ./internal/platform/config` | All packages pass |
| Race suite | `GOCACHE=/tmp/ultraplan-go-build-cache go test -race ./...` from `../ultraplan-go` | All packages pass (sprint 1.134s, app cached, others cached) |
| CLI build | `GOCACHE=/tmp/ultraplan-go-build-cache GOMODCACHE=/tmp/ultraplan-go-mod-cache go build ./cmd/ultraplan` | Builds successfully with no errors |

The race suite passes for `internal/sprint` and `internal/app`, which is the relevant scope for execute behavior.

## 3. Contract Conformance

Each contract was reviewed by an independent subagent against the requirements in the contract file, the sprint's `requirements.md`, `reasoning.md` decisions, and the implementation in `../ultraplan-go`. Only `direct` and `partial` applicability items that produced a concrete finding are listed here.

### 3.1 Architecture — `pass_with_issues` (was `pass_with_issues`)

| Requirement | Applicability | Status | Evidence |
| --- | --- | --- | --- |
| Module boundaries explicit | direct | satisfied | `internal/sprint/execute.go:15` declares `StageExecute`; per-concern files at `execute_state.go`, `execute_plan.go`, `execute_target.go`, `execute_model.go`, and `execute.go`; `internal/platform/runtime/*` mtimes pre-date 2026-07-09 (untouched for Sprint 23) |
| Dependency direction inward | direct | satisfied | `internal/sprint/execute.go:11-12` imports only `pruntime` (`internal/platform/runtime`) and `workspace`; `grep internal/study` under `internal/sprint` returns no matches; `go list -deps ./internal/platform/runtime` only reaches `agentwrap`, `agentwrap/opencode`, `internal/platform/config` |
| Domain layer pure | direct | satisfied | `internal/sprint/domain.go:58-129` is pure data types (`ExecuteTaskStatus`, `ExecuteTaskIdentity`, `ExecuteTaskRecord`, `ExecuteRunState`, `ExecuteTargetRef`) with no transport/SDK imports inside type definitions |
| Use cases depend on ports | partial | **partially satisfied** | `Runtime` interface at `internal/sprint/flow.go:12-14` is narrow (one method); but `internal/sprint/runtime_validation.go:8` directly imports `github.com/Antonio7098/agentwrap` and constructs `agentwrap.ValidationSpec` types. Pre-existing, not introduced by Sprint 23 |
| Entrypoint adapters stay thin | partial | **partially satisfied** | `runSprint` execute case at `internal/app/sprint_commands.go:158-184` is ~25 lines and calls `Service.Execute` after `parseSprintExecuteArgs` (line 368); but the `flow --to execute` path at `sprint_commands.go:237-239` does stage sequencing in `internal/app` rather than through a sprint-owned facade |
| Module public APIs small and stable | partial | **partially satisfied** | The package godoc at `internal/sprint/doc.go:3-7` still says the package "models only the governed planning chain through plan.md" and lists "execute implementation work" as forbidden at lines 12-13, contradicting the new sprint-owned execute semantics |
| Sprint owns execute semantics | direct | satisfied | Task extraction and IDs in `execute_plan.go`; run-state schema and atomic persistence in `execute_state.go`; target resolution and workdir containment in `execute_target.go`; model precedence in `execute_model.go` and `execute.go:330-344`; prompt rendering and summary generation in `execute.go:196-254`; reconciliation in `execute.go:267-285` |
| Target path containment | direct | satisfied | `internal/sprint/execute_target.go:13-52` rejects missing/relative/alternate paths and enforces lexical containment |

**Issues**:
- `internal/sprint/doc.go:3-7` still claims "This package models only the governed planning chain through plan.md" and forbids "execute implementation work" (line 12-13), contradicting the new sprint-owned execute semantics, `execute.md`, and `.run-state.json` ownership.
- `internal/sprint/runtime_validation.go:8` directly imports `github.com/Antonio7098/agentwrap` (pre-existing) and constructs `agentwrap.ValidationSpec` types — platform/runtime exposes `*agentwrap.ValidationSpec` through its public `Request` API, leaking the agentwrap SDK through the product boundary.
- `flow --to execute` orchestration lives in `internal/app/sprint_commands.go:237-239` rather than a single sprint-owned `FlowTo`/`FlowExecute` facade.

### 3.2 CLI Surface — `fail` (was `pass_with_issues`)

| Requirement | Applicability | Status | Evidence |
| --- | --- | --- | --- |
| `validate execute` wired | direct | satisfied | `internal/app/sprint_commands.go:86-87` calls `service.ValidateExecute` for `StageExecute` |
| `prompt execute` wired | direct | satisfied | `internal/app/sprint_commands.go:118-119` calls `service.PromptExecute` for `StageExecute` |
| `flow --to execute` wired | direct | satisfied | `internal/app/sprint_commands.go:204-205` adds execute to the stage chain for `flow`; `:237-239` invokes `service.Execute` with `Resume: true`; `parseSprintFlowArgs` at `:362` accepts `--to execute` |
| `execute` subcommand | direct | satisfied | `internal/app/sprint_commands.go:158-184` dispatches `execute` to `execService.Execute` after `parseSprintExecuteArgs` |
| `--task` selection | direct | satisfied | `parseSprintExecuteArgs` at `internal/app/sprint_commands.go:368`; `ExecuteRequest.TaskID` propagated to `execute.go:94-96, 146-148` |
| Dry-run | direct | satisfied | `parseSprintExecuteArgs` accepts `--dry-run`; `execute.go:66-79` short-circuits and returns the rendered prompt |
| Help | partial | **partially satisfied** | `sprintExecuteHelp()` does not exist and there is no `--help` handling for the `execute` subcommand (only `--help` at lines 28-31 covers status); `sprintHelp()` at line 549-550 still says "Supports planning stages through plan.md only. It does not execute implementation…"; `sprintValidateHelp`, `sprintPromptHelp`, `sprintFlowHelp` at lines 564-606 do not list `execute` |
| Exit codes | partial | **partially satisfied** | `sprint_commands.go:174-180` maps prerequisite → `ExitValidation`, failed-task result → `ExitPartial`, runtime error → `ExitRuntime`; uses fragile `strings.Contains(err.Error(), "failed tasks"\|"runtime")` matching; `mapSprintError` at `:328-341` does not classify `ErrExecuteRunStateMissing/Malformed/Unsupported` (fall through to `ExitWorkspace`) |
| `--json` output | direct | **not satisfied** | No `--json` flag for `sprint <p> <s> status` or `sprint <p> <s> execute`; only `config show` has `--json` (uses `jsonEnvelope` at `json_output.go:11-32`); `status_json.go` is study-only |
| Status output reflects execute | direct | **not satisfied** | `renderSprintStatus` (`sprint_commands.go:395-408`) iterates only planning stages; `PlanningStages()` (`domain.go:165-174`) excludes `StageExecute`; `StatusSummary` has no `Execute` field — REQ-23-48 unmet |
| `execute` command tests | direct | **not satisfied** | `internal/app/sprint_execute_commands_test.go` does not exist; `sprint_commands_test.go` has zero execute-related cases |

**Issues**:
- `internal/app/sprint_execute_commands_test.go` required by `requirements.md` line 32 is absent.
- `sprintHelp()` (`sprint_commands.go:519-552`) contradicts the wired `execute` command; `sprintValidateHelp`, `sprintPromptHelp`, `sprintFlowHelp` (lines 564-606) do not list `execute`.
- No `sprintExecuteHelp()` and no `--help` dispatch for the `execute` subcommand at `sprint_commands.go:32-44`.
- `renderSprintStatus` (`sprint_commands.go:395-408`) and `StatusSummary` do not surface execute readiness/progress.
- `mapSprintError` (`sprint_commands.go:328-341`) does not classify `ErrExecuteRunStateMissing/Malformed/Unsupported`.
- `renderSprintExecute` (`sprint_commands.go:478-508`) prints per-task ID/status/attempts but not task counts by terminal state or evidence/diagnostic lines that `WriteExecuteSummary` writes to `execute.md`.

### 3.3 Configuration — `fail` (was `pass_with_issues`)

| Requirement | Applicability | Status | Evidence |
| --- | --- | --- | --- |
| Config from explicit sources | direct | satisfied | `internal/platform/config/config.go:111-134` documents precedence; sources map at `config.go:113-128` records `default`/`workspace`/`env`/`cli` tiers |
| Typed/parsed/validated | direct | satisfied | `Config` is fully typed; `Validate()` at `config.go:321-378` rejects unknown top-level fields and missing required runtime fields |
| Public config allowlisted | direct | satisfied | `Redacted` struct at `internal/platform/config/redaction.go:7-16`; `setField` rejects unknown config fields |
| Secrets separated from public config | direct | **not satisfied** | `redactPlanning` at `internal/platform/config/redaction.go:48-56` redacts `RequirementsModel`, `SprintIndexModel`, `TechnicalHandbookModel`, `AreaReasoningModel`, `ReasoningModel`, `PlanModel` but **does NOT** redact `ExecuteModel` or `ExecuteVariant` (or any `*Variant` field) |
| Effective config inspectable | direct | **not satisfied** | `internal/app/config_commands.go:67-68` prints `planning.execute_model` and `planning.execute_variant` unredacted (because they are not passed through `RedactValue`); the resolved execute model source/value is not surfaced via `config show` |
| Required config fails before work starts | partial | **partially satisfied** | `ResolveExecuteModel` at `execute_model.go:15-30` returns an explicit error; but `Service.executeModelSelection` at `internal/sprint/execute.go:330-344` silently returns `"provider/model"` with `Source: "default"` instead of failing fast and uses a precedence that drops `models.primary` |
| REQ-23-41 stage model precedence | direct | **partially satisfied** | `ResolveExecuteModel` (helper, tested) implements the full chain; `Service.executeModelSelection` (in-service, not tested) duplicates with drift |

**Issues**:
- `redactPlanning` does not redact `ExecuteModel`/`ExecuteVariant` (R-5 from prior review remains unresolved; the redaction test in `config_test.go:69-83` does not cover these fields).
- `Service.executeModelSelection` (`execute.go:330-344`) hard-codes `"provider/model"` fallback and skips the `models.primary` tier. The two resolvers have drifted; only the helper is exercised by tests, so the production code path is not regression-protected.
- Sources map (`config.go:113-128`) does not surface the resolved stage-tier (`planning.execute_model` vs `models.primary` vs `models.default`); field-path diagnostics on empty values are absent (Decision 7 expects source visibility).

### 3.4 Documentation — `fail` (was `pass_with_issues`)

| Requirement | Applicability | Status | Evidence |
| --- | --- | --- | --- |
| Execute summary generator | direct | satisfied | `internal/sprint/execute.go:223-254` `WriteExecuteSummary` writes `execute.md` with plan citation, run-state path, task counts by terminal state, per-task evidence/diagnostic lines |
| Summary relpath helper | direct | satisfied | `internal/sprint/artifacts.go:26-27` adds `StageExecute` case to `ArtifactRelPath` returning `execute.md` |
| Editable Markdown + versioned JSON | direct | satisfied | `execute.md` is editable Markdown; `.run-state.json` is versioned (`ExecuteRunStateSchemaVersion = 1`) |
| Agent/human help text | direct | **not satisfied** | `internal/app/sprint_commands.go:519-606` help text remains stale; `sprintHelp()` at line 549-550 still says "Supports planning stages through plan.md only. It does not execute implementation…"; per-command help blocks (`sprintValidateHelp`, `sprintPromptHelp`, `sprintFlowHelp`) do not list `execute` |
| Project-wide docs | partial | **not satisfied** | `docs/cli-reference.md:223-268` and `docs/user-guide.md:215-232` do not list `execute`, validate/prompt/flow variants, `--task/--dry-run/--resume/--model` flags, or status output that surfaces execute; `internal/sprint/doc.go` still says "models only the governed planning chain through plan.md" |
| Recovery runbook | partial | **partially satisfied** | Stale-running recovery and `.run-state.json` schema policy are not documented in `docs/recovery.md` |

**Issues**:
- `sprintHelp()`/`sprintValidateHelp()`/`sprintPromptHelp()`/`sprintFlowHelp()` are stale; R-2 from prior review remains open.
- `docs/cli-reference.md` and `docs/user-guide.md` are not updated; R-3 (status content) overlaps because status docs do not mention execute readiness/progress.
- `internal/sprint/doc.go` is stale (Architecture finding mirrors).
- No golden/fixture test for `execute.md` (F-8 still open).

### 3.5 Errors — `pass_with_issues` (unchanged)

| Requirement | Applicability | Status | Evidence |
| --- | --- | --- | --- |
| Sentinel errors for execute run state | direct | satisfied | `ErrExecuteRunStateMissing`, `ErrExecuteRunStateMalformed`, `ErrExecuteRunStateUnsupported` at `internal/sprint/domain.go:18-20` |
| Wrapped errors with operation/path | direct | satisfied | `LoadExecuteRunState` wraps `fs.ErrNotExist` (`execute_state.go:42-46`), `json.Unmarshal` errors (`:48-50`), `ValidateExecuteRunState` results (`:51-53`); `SaveExecuteRunState` wraps `os.CreateTemp`, `os.Rename`, etc. |
| Actionable diagnostics | direct | satisfied | `ValidationFinding{Problem, Cause, Suggestion}` carried into `execute_target.go:13-32`, `execute_plan.go:23-107`, `execute.go:188-191` |
| Exit-code mapping | direct | **partially satisfied** | `sprint_commands.go:174-180` maps prerequisite/partial/runtime categories for `execute`; `mapSprintError` at `:328-341` matches `ErrFlowStateMalformed/Unsupported` but not `ErrExecuteRunState*`; uses `strings.Contains` rather than `errors.Is`/`errors.As` |
| Runtime-error classification & redaction | direct | **partially satisfied** | `ExecuteRuntimeSummary` at `domain.go:88-97` stores only safe fields with explicit `OmissionReason: "raw runtime payloads omitted"`; `safeError` at `flow.go:139-151` strips control characters and truncates at 180 chars but does **not** run secret-pattern redaction; diagnostic messages are persisted unredacted |
| Persistence errors classified distinctly | direct | **partially satisfied** | Read failures wrap `ErrExecuteRunStateMissing/Malformed/Unsupported` correctly; `SaveExecuteRunState` write failures (mkdir, write, rename) return raw wrapped errors without a `data.*` sentinel and collapse to `ExitWorkspace` via `mapSprintError` |
| Cancellation maps to inspectable state | direct | satisfied | `execute.go:125-127` maps `ctx.Err()` to `ExecuteTaskCancelled` with `cancelled` diagnostic |
| Cancellation exit code | direct | **partially satisfied** | `ExitCancel` is defined in `app.go:23` but is never emitted for execute cancellation paths |

**Issues**:
- `mapSprintError` does not classify `ErrExecuteRunStateMissing/Malformed/Unsupported` (R-4 from prior review remains unresolved).
- Execute exit-code map uses `strings.Contains` heuristics (`sprint_commands.go:176-181`) — fragile against wording changes; should be `errors.Is` checks against stable sentinels.
- `safeError` does not redact secret patterns (API keys, bearer tokens, `sk-`/`ghp-`/`gho_`/`xox*-` prefixes).
- `SaveExecuteRunState` write failures lack a `data.*` sentinel and fall through to `ExitWorkspace`.

### 3.6 LLM Evaluation / Cost / Safety — `pass_with_issues` (unchanged)

| Requirement | Applicability | Status | Evidence |
| --- | --- | --- | --- |
| Runtime success insufficient for completion | direct | satisfied | `internal/sprint/execute.go:124-142` explicit switch: `ctx.Err` → `cancelled`; `runErr` → `failed` with `runtime-failed`; `run.Artifacts > 0` → `complete` with evidence; `hasDiagnosticOnlyCompletion(run)` → `complete` with `diagnostic-only-completion`; default → `failed` with `missing-evidence` |
| Safe metadata persistence | direct | satisfied | `ExecuteRuntimeSummary` at `domain.go:88-97` records only `RunID`, `SessionID`, `Model`, `ModelSource`, `PermissionSummary`, `ValidationSummary`, and `OmissionReason: "raw runtime payloads omitted"` (built at `execute.go:287-289`) |
| Omission reason recorded | direct | satisfied | Every persisted runtime summary sets `OmissionReason` |
| Eval-Model-001 reviewable | direct | **partially satisfied** | `executeModelSelection` (`execute.go:330-344`) precedence is reviewable but uses a different fallback ("provider/model") than the documented Decision 7 chain; only `ResolveExecuteModel` in `execute_model.go:15-30` is exercised by tests |
| Fake-runtime lifecycle tests | direct | **not satisfied** | `internal/sprint/execute_test.go` is absent (verified by glob); `sprint_index_test.go:170-198` `fakeRuntime` is wired only to `FlowSprintIndex`, not `Service.Execute`; the five completion arms in `execute.go:124-142` are not exercised by any test |
| Safe metadata redaction in config diagnostics | direct | **partially satisfied** | `redactPlanning` covers `RequirementsModel` through `PlanModel`; `ExecuteModel`/`ExecuteVariant` are not redacted (same gap as the Configuration and Security reviews) |
| Eval-Cost-001 (token/duration/cost) | direct | **partially satisfied** | `ExecuteRuntimeSummary` declares a `UsageSummary` field but `runtimeSummary` at `execute.go:287-289` never populates it; `pruntime.Result` exposes `Usage`, `EstimatedCost`, `StartedAt/FinishedAt`, all discarded |

**Issues**:
- The contract's `Eval-REG-001` requirement that "Critical AI Behaviours Need Regression Examples" is unmet for the completion switch — five arms are not exercised by tests.
- `redactPlanning` does not cover `ExecuteModel`/`ExecuteVariant`.
- `runtimeSummary` discards `Usage`/`EstimatedCost`/latency that the platform Result carries.

### 3.7 LLM Runtime — `pass_with_issues` (was `pass_with_issues`)

| Requirement | Applicability | Status | Evidence |
| --- | --- | --- | --- |
| Generic `internal/platform/runtime` boundary | direct | satisfied | `internal/platform/runtime/*` was not modified for Sprint 23 (mtimes pre-date 2026-07-09); the `Runtime` interface in `internal/sprint/flow.go:12-14` is narrow (one method); `Service.Execute` (`execute.go:118`) calls through the seam |
| Sprint uses only the seam | direct | **partially satisfied** | `internal/sprint/runtime_validation.go:8` directly imports `github.com/Antonio7098/agentwrap` and constructs `agentwrap.ValidationSpec` types. **Pre-existing**, not introduced by Sprint 23 but not yet refactored. Enabled by `internal/platform/runtime/runtime.go:30` exposing the agentwrap type directly through the platform boundary |
| No native OpenCode stream parsing | direct | satisfied | `internal/sprint/*` has no `os/exec`, no stream parsing, and no OpenCode adapter import |
| Tool execution boundary | direct | satisfied | `execute.go:117` adds explicit `pruntime.PermissionPathRule{Path: target.Path, Action: "allow"}`; no hardcoded tool selection |
| Lifecycle controls explicit | direct | satisfied | `ExecuteTaskStatus` enum (`pending`/`running`/`complete`/`failed`/`cancelled`) at `domain.go:60-66`; `execute.go:92-149` terminal transitions; `reconcileExecuteState` (`execute.go:267-285`) recovers stale running tasks |
| Cost/latency observability | direct | **partially satisfied** | `runtimeSummary` (`execute.go:287-289`) keeps only the safe meta set; `pruntime.Result` already carries `Usage`, `EstimatedCost`, `StartedAt/FinishedAt` (`runtime.go:47-67`) which are discarded |
| Prompt version propagation | partial | **not satisfied** | `execute.go:115` attaches metadata `{project, sprint, stage, task, model_source}` but no `prompt_id`/`prompt_version`; pre-existing pattern gap (upstream planning stages also do not propagate) |
| Platform `policy.go` stays neutral | direct | **partially satisfied** | `internal/platform/runtime/policy.go:4-14` exports `MetadataStudy`, `MetadataDimension`, `MetadataTaskKind`, `MetadataSource`, `MetadataSourceKind`, `MetadataOutputPath`; grep shows zero consumers in the repo (dead, mode-specific content in shared substrate). Pre-existing. |

**Issues**:
- `internal/sprint/runtime_validation.go` directly imports `github.com/Antonio7098/agentwrap` and constructs `agentwrap.ValidationSpec` values (pre-existing, F-5).
- `internal/platform/runtime/policy.go` exports unused domain-metadata constants (F-6).
- `ExecuteRuntimeSummary` discards `Usage`/`EstimatedCost`/latency from `pruntime.Result`.

### 3.8 Observability — `fail` (was `pass_with_issues`)

| Requirement | Applicability | Status | Evidence |
| --- | --- | --- | --- |
| Status reflects execute readiness/progress without runtime calls | direct | **not satisfied** | `Service.Execute` writes `.run-state.json` at each transition (`execute.go:113, 143`) and `execute.md` on completion (`:150`); but `Service.Status` (in `internal/sprint/service.go`) and `renderSprintStatus` (`sprint_commands.go:395-408`) do not load `.run-state.json` or `execute.md`; `PlanningStages()` (`domain.go:165-174`) does not include `StageExecute`. REQ-23-48 is unimplemented at the operator surface |
| Structured diagnostics | direct | **partially satisfied** | `ExecuteDiagnostic{Code, Message, At}` (`domain.go:76-80`) is structured and validated; `ExecuteRuntimeSummary` carries safe metadata; but no `internal/platform/logging.Logger` is invoked from any execute path — the canonical structured logging surface is dead code (verified by repo-wide grep) |
| Counts, attempts, diagnostics inspectable | direct | **partially satisfied** | Persisted in `.run-state.json`, but `renderSprintStatus` and `StatusSummary` do not surface them |
| Sensitive data redaction in diagnostics | direct | **not satisfied** | `ExecuteDiagnostic.Message` and `ExecuteEvidence.Summary` are persisted unredacted; `safeError` strips CR/LF/NUL but does not run `config.RedactValue` over messages and does not scrub secret-bearing patterns |

**Issues**:
- `Service.Status` / `renderSprintStatus` do not read `.run-state.json`; operators must open the JSON file directly to inspect execute progress (R-3 from prior review).
- No structured log events from execute paths; `internal/platform/logging.Logger` is dead code.
- `ExecuteDiagnostic` redaction is not applied before persistence.

### 3.9 Performance — `pass_with_issues` (unchanged)

| Requirement | Applicability | Status | Evidence |
| --- | --- | --- | --- |
| Status bounded and deterministic | direct | satisfied | `Service.Status` is single-shot; execute artifacts are not opened by status today |
| No unbounded goroutines / runtime fan-out | direct | satisfied | `internal/sprint/execute.go` has zero `go` statements; tasks processed serially (`:92-149`); `Complete` tasks skipped on resume (`:97-99`) |
| Status without runtime calls | direct | satisfied | Status does not invoke runtime |
| No recursive target scans | direct | satisfied | `execute_target.go:13-34` reads only the project index's `Target Implementation Directory` line and runs one `os.Stat` on the resolved path; status does not open the target repo |
| Save-after-meaningful-transition budget | partial | **partially satisfied** | `Service.Execute` calls `SaveExecuteRunState` twice per task (`:113, 143`) plus once at `Resume` reconciliation (`:89`); the intermediate save at `:113` discards the error with `_` |
| `PlanFingerprint` cached | partial | **partially satisfied** | `PlanFingerprint` is persisted in `.run-state.json` (compute at `execute.go:85`) but never compared to a recomputed value at load or resume; `mustReadPlan` (`execute.go:318-321`) silently swallows read errors, so the fingerprint can become a hash of an empty string |

**Issues**:
- `mustReadPlan` silently swallows read errors; computes fingerprint possibly over empty content.
- Intermediate `SaveExecuteRunState` at `execute.go:113` discards the error with `_`.

### 3.10 Persistence And Migrations — `pass_with_issues` (unchanged)

| Requirement | Applicability | Status | Evidence |
| --- | --- | --- | --- |
| Versioned schema | direct | satisfied | `ExecuteRunStateSchemaVersion = 1` at `domain.go:11`; `ValidateExecuteRunState` rejects mismatched `SchemaVersion` at `execute_state.go:117-119` |
| Strict load validation | direct | satisfied | `LoadExecuteRunState` at `execute_state.go:35-55` unmarshals then runs `ValidateExecuteRunState` (`:113-183`); rejects project/sprint mismatch, unsafe `planPath`, duplicate task IDs, invalid statuses, missing timestamps, unsafe diagnostics/evidence paths |
| Same-directory atomic writes | direct | **partially satisfied** | `SaveExecuteRunState` at `execute_state.go:57-111` uses `os.CreateTemp` in same directory, `fsync`, `os.Rename`, `syncDir`; `TestExecuteRunStateStrictLoadingAndAtomicWritePreservesPrior` exercises the `BeforeRename` failure and confirms last-known-good preservation; `WriteExecuteSummary` (`execute.go:250-253`) uses raw `os.WriteFile`, not atomic — could tear mid-write |
| Resumability | direct | **partially satisfied** | `reconcileExecuteState` (`execute.go:267-285`) reconciles existing state with newly extracted tasks; stale `running` tasks become `failed` with `stale-running` diagnostic; `Complete` tasks preserved (`byID` at `:278`) |
| Stale-running recovery | direct | satisfied | `execute.go:267-285` (in-memory reconciliation) and `execute.go:100-104` (per-task recovery at start of run) |
| Validator vs reconciler divergence | direct | **partially satisfied** | Strict validator rejects `Running` tasks without `StartedAt` (`execute_state.go:162`); `reconcileExecuteState` is the recovery path that can repair them. Intentional per REQ-23-47 but a future reload-then-validate code path (status) must use `reconcileExecuteState` |
| Plan fingerprint drift detection | direct | **not satisfied** | `PlanFingerprint` is recorded and validated by the schema (`execute_state.go:128-137`) but is not compared against a recomputed fingerprint at load/resume; drift between runs is silent (F-2) |
| No migration path | direct | satisfied (intentional) | First version is `1`; future versions will need migration logic |
| Atomic write (summary) | direct | **not satisfied** | `WriteExecuteSummary` uses `os.WriteFile` directly at `execute.go:253` |

**Issues**:
- `ValidateExecuteRunState` rejects a state file with a stale `running` task whose `startedAt` is set but the task is otherwise valid (validator vs reconciler divergence, partially intentional).
- No explicit plan-fingerprint mismatch detection (F-2).
- `WriteExecuteSummary` is non-atomic.

### 3.11 Security — `pass_with_issues` (was `pass_with_issues`)

| Requirement | Applicability | Status | Evidence |
| --- | --- | --- | --- |
| Path containment / write scoping | direct | **partially satisfied** | `ApprovedExecuteTargetPath` hard-pinned to `/home/antonioborgerees/coding/ultraplan-go` (`execute_target.go:11`); `ResolveExecuteTarget` (`execute_target.go:13-34`) requires an absolute path matching the approved target and an existing directory; `ValidateExecuteWorkdir` (`execute_target.go:36-52`) rejects empty/relative/outside-target workdirs. **Symlink resolution is not used** (no `filepath.EvalSymlinks` in either function); a symlink inside the target root pointing outside would pass lexical containment |
| Secrets redacted in config, diagnostics, status | direct | **partially satisfied** | `internal/platform/config/redaction.go:48-56` covers `RequirementsModel` through `PlanModel` but not `ExecuteModel`/`ExecuteVariant`; `safeError` strips control characters but does not redact secret patterns; `ExecuteDiagnostic.Message` and `ExecuteEvidence.Summary` are persisted unredacted |
| No shell/cmd construction | direct | satisfied | `internal/sprint/*` has no `os/exec` and no shell-out |
| Target authorization allowlist | direct | satisfied | Exact-match allowlist via `ApprovedExecuteTargetPath` |
| Workspace-relative diagnostics | direct | satisfied | `Service.Execute` and related paths use `workspace.Rel(s.root, statePath)` (`execute.go:154`) and `ArtifactRelPath` for summary paths |
| No Git mutation | direct | satisfied | No `os/exec` to git; `execute_target.go:54-61` `ExecuteSafetyInstructions` explicitly forbids Git mutation; `plan.go:166-167` rejects `git commit`/`git push` mentions in plan content |
| Smoke/review/issues/hosted/TUI excluded | direct | satisfied | `execute_target.go:54-61` lists forbidden artifacts; `RenderExecutePrompt` (`execute.go:215-218`) embeds them |
| In-memory payload retention | direct | **partially satisfied** | `ExecuteResult.Runtime []pruntime.Result` (`execute.go:33, 119`) retains the full runtime result including `Events.Payload`, `Attempts`, `Memory`; persisted run-state is safe via `runtimeSummary` (`execute.go:287-289`); the field is exported and unannotated (no `json:"-"`), so any future marshal would leak |
| Git mutation denial enforcement | partial | **partially satisfied** | Rejection is layered: (1) `plan.go:166-167` text scanning for `git commit`/`git push`; (2) `execute_target.go:54-61` prompt-level instructions; (3) `execute.go:117` single allow rule for `target.Path`. Enforcement is delegated to the runtime adapter; no code-level scanner or explicit deny rule for `.git` inside the target |

**Issues**:
- `ValidateExecuteWorkdir` is lexical-only; does not call `filepath.EvalSymlinks` (F-1).
- `Planning.ExecuteModel`/`ExecuteVariant` are not redacted in `config show` (R-5).
- `ExecuteResult.Runtime` retains the full runtime payload at the service boundary.
- No defense-in-depth for Git mutation (`git commit`/`git push` tokens are not scanned at the execute stage; `.git` paths are not explicitly denied in `runtimeReq.Policy.PathRules`).

### 3.12 Testing — `fail` (was `fail` — R-1 still unresolved)

| Requirement | Applicability | Status | Evidence |
| --- | --- | --- | --- |
| Unit tests for business logic | direct | **not satisfied** | `internal/sprint/execute_test.go` (required by `requirements.md` line 31) does not exist as a consolidated runner test file. `Service.Execute`, `WriteExecuteSummary`, `RenderExecutePrompt`, `reconcileExecuteState`, `runtimeSummary`, `hasDiagnosticOnlyCompletion`, `safeArtifactPath`, `hasFailedExecuteTask`, `executeModelSelection` are all untested at the function level |
| Wiring/persistence integration coverage | direct | **partially satisfied** | `Service.ValidateExecute` is exercised in `execute_plan_test.go:71-100`; `execute_state_test.go:11-56` covers atomic write + last-known-good preservation; `execute_state_test.go:58-127` covers 12 validation cases; `execute_plan_test.go` and `execute_target_test.go` cover extraction/target. But no full integration test seeding `.run-state.json` with stale running tasks and running `Service.Execute` end-to-end |
| Failure paths tested explicitly | direct | **partially satisfied** | Helper-level paths (`execute_state_test.go:58-127`, `execute_target_test.go:18-31`, `execute_plan_test.go:54-69`, `execute_model_test.go:45-52`, `execute_target_test.go:33-47`) cover malformed/missing/unsupported/duplicate/unsupported-syntax/missing-model/escape well. **Runner-level failure paths (runtime-failed, missing-evidence, cancelled, diagnostic-only completion) are untested** because `execute_test.go` is missing |
| Determinism | direct | satisfied | All current tests use fixed times or `time.Now` injection; no real-runtime calls |
| Compatibility tests for CLI/contracts | direct | **not satisfied** | `internal/app/sprint_execute_commands_test.go` (required by `requirements.md` line 32) does not exist. No command-level tests cover `execute`, `validate execute`, `prompt execute`, `flow --to execute`, flags, exit codes, or text/JSON output |
| End-to-end scenarios | direct | **not satisfied** | No scenario test for the full execute journey: prepare → dry-run → run → evidence gate → summary → status → resume |
| Golden tests | direct | **partially satisfied** | No golden tests for `RenderExecutePrompt` or `WriteExecuteSummary` (F-8); state-level golden-like assertions exist for atomic write |

**Issues**:
- `internal/sprint/execute_test.go` is absent (R-1).
- `internal/app/sprint_execute_commands_test.go` is absent (R-1).
- The five completion arms in `execute.go:124-142` are not exercised by any test.
- No end-to-end scenario test exists for the full execute journey.

### 3.13 Workflows — `pass_with_issues` (was `pass_with_issues`)

| Requirement | Applicability | Status | Evidence |
| --- | --- | --- | --- |
| Execute is a stateful workflow stage | direct | **partially satisfied** | `StageExecute` at `internal/sprint/execute.go:15`. `validateFlowTarget` (`flow.go:154`) accepts `StageExecute`. `parseSprintFlowArgs` (`sprint_commands.go:362`) accepts `--to execute`. The `flow` switch dispatches to `service.Execute` for `StageExecute` (`sprint_commands.go:237-239`). **However, `PlanningStages()` (`domain.go:165-174`) does not include `StageExecute`**, so flow chains and `StatusSummary` cannot surface execute readiness/progress end to end |
| Prerequisite gating through `plan.md` | direct | satisfied | `Service.prepareExecute` (`execute.go:164-194`) reads the project index, validates target, reads `plan.md`, extracts tasks, validates target workdir; any failure populates `findings` and `Service.Execute` returns an error before runtime launch |
| Task states with attempts, resumability, cancellation | direct | satisfied | `ExecuteTaskStatus` enum (`pending`, `running`, `complete`, `failed`, `cancelled`) at `domain.go:60-66`. `Service.Execute` (`execute.go:92-149`) tracks `Attempts`, `StartedAt`, `CompletedAt`, and supports `ctx.Err()` cancellation |
| Stale running recovery | direct | satisfied | `reconcileExecuteState` (`execute.go:267-285`) transitions stale `running` to `failed` with `stale-running` diagnostic on resume |
| Skip `complete` on resume | direct | satisfied | `execute.go:97-99` skips tasks already in `ExecuteTaskComplete` |
| Terminal-state flow summaries | direct | **partially satisfied** | `WriteExecuteSummary` writes a Markdown summary with task counts and per-task status; but `flow-state.json`'s `stages` array (written by `flow.go` builders) does not currently include a flow-stage for execute; no `flowExecuteSuccessStages` / `flowExecuteFailedStages` builders in `flow.go` |
| Deferred stages rejected | direct | **partially satisfied** | `validateFlowTarget` (`flow.go:153-157`) accepts `StageExecute` and rejects anything else with "supports requirements, sprint-index, technical-handbook, area-reasoning, reasoning, plan, and execute". `smoke`/`review`/`issues` are not explicitly named as deferred/non-goal stages here, but `plan.go:166-167` and `ExecuteSafetyInstructions` (`execute_target.go:54-61`) forbid them in content |
| Idempotency for retried tasks | direct | **not satisfied** | `ExecuteTaskRecord` and `ExecuteRequest` carry no `IdempotencyKey`/`RunKey`/`DedupeToken`; the runtime is not told that a retry is a retry; re-running a failed task can produce duplicate external side effects |
| Compensation / reconciliation | direct | **not satisfied** | On `ctx` cancellation or per-task failure, the run-state records the status but the workdir may contain partial writes; no rollback, undo, cleanup, or reconciliation hook. `hasFailedExecuteTask` only reports failure; no compensation surface |

**Issues**:
- `PlanningStages()` (`internal/sprint/domain.go:165-174`) does not include `StageExecute` (F-3).
- `flowExecuteSuccessStages`/`flowExecuteFailedStages` builders are not added in `flow.go` (F-4).
- No idempotency key for retried tasks.
- No compensation or reconciliation for partial external side effects.
- `Service.FlowExecute` (or extension of `Service.Flow`) does not exist as a discrete method; flow through execute is a switch case in `internal/app/sprint_commands.go:237-239`.

### 3.14 Technical Handbook — `pass_with_issues` (was `pass_with_issues`)

Nine patterns and eight anti-patterns were reviewed. Status summary:

| Pattern | Applicability | Status |
| --- | --- | --- |
| Thin command boundary over product-owned execute behavior | direct | partial (stale help text + status summary gap) |
| Module-local ownership with enforced dependency direction | direct | followed |
| Explicit composition root and narrow volatile interfaces | direct | followed |
| Config precedence with stage-specific override visibility | direct | partial (redaction gap + duplicated model selection) |
| Runtime success paired with artifact validation | direct | followed (in code) |
| Atomic durable state and resumable task lifecycle | direct | followed |
| Inspectable diagnostics with safe observability | direct | partial (re-export of latent gaps; no structured logging; in-memory payload retention) |
| Workspace-safe and permission-aware execution boundary | direct | followed (with symlink caveat) |
| Behavior-focused test harnesses with fakes and golden fixtures | direct | partial (no execute_test.go, no sprint_execute_commands_test.go, no execute.md golden) |

| Anti-pattern | Status |
| --- | --- |
| Execute must not become a monolithic `RunE` | avoided (Service.Execute delegates; command handler is ~25 lines) |
| No global workflow/scheduler package | avoided |
| No package-level mutable config or runtime state | avoided |
| No bypassing IO abstractions | avoided |
| No success without artifact/evidence validation | avoided |
| No unbounded goroutines or runtime fan-out | avoided |
| No persistence of unsafe raw runtime payloads or secrets | partial (`redactPlanning` gap for `ExecuteModel`/`ExecuteVariant`; `ExecuteResult.Runtime` retains full payload in memory) |
| No smoke, review, issue tracking, Git mutation | avoided (safety instructions embedded in prompt; plan validator rejects; doc.go stale) |

## 4. Handbook Applicability Table

This table records the patterns and anti-patterns from the technical handbook that were classified as `partial`, `not_triggered`, or `explicitly_deferred` for this sprint, with the scope reason.

| Pattern / Anti-pattern | Classification | Scope reason |
| --- | --- | --- |
| Thin command boundary over product-owned execute behavior | partial | Command handlers stay small but `sprintHelp`/`sprintValidateHelp`/`sprintPromptHelp`/`sprintFlowHelp` are stale; `renderSprintStatus` does not surface execute; `renderSprintExecute` is missing task counts and evidence/diagnostic lines |
| Config precedence with stage-specific override visibility | partial | `redactPlanning` covers all stage models except `ExecuteModel`/`ExecuteVariant`; `Service.executeModelSelection` duplicates `ResolveExecuteModel` with a different fallback; precedence is otherwise correct |
| Inspectable diagnostics with safe observability | partial | `ExecuteDiagnostic` is structured and redacted re control characters, but `internal/platform/logging.Logger` is never invoked from execute paths; status does not load `.run-state.json`; `ExecuteResult.Runtime` retains the full payload at the service boundary |
| Workspace-safe and permission-aware execution boundary | partial | Lexical containment is enforced; `filepath.EvalSymlinks` is not used (a hardening item per F-1); Git-mutation denial is layered but not defended at the code level |
| Behavior-focused test harnesses with fakes and golden fixtures | partial | Unit tests cover extraction/state/target/model; no `internal/sprint/execute_test.go` exists; no `internal/app/sprint_execute_commands_test.go` exists; no fixture/golden test for `execute.md` |
| Atomic durable state and resumable task lifecycle | direct | Followed for everything in scope |
| Module-local ownership with enforced dependency direction | direct | Followed; `internal/platform/runtime/*` is product-agnostic |
| Explicit composition root and narrow volatile interfaces | direct | Followed through `Service` + `WithRuntime`/`WithStageRuntime` |
| Runtime success must be paired with artifact validation | direct | Followed via the five-arm completion switch in `execute.go:124-142` |

## 5. Issue Classification

The following items are categorized as **Required Fixes** (must be resolved before the next sprint) or **Future Follow-ups** (deferred with justification). R-1 through R-7 and F-1 through F-8 from the prior review are merged into this section.

### 5.1 Required Fixes

| ID | Severity | Issue | Evidence | Recommended action |
| --- | --- | --- | --- | --- |
| R-1 | blocker | `internal/sprint/execute_test.go` and `internal/app/sprint_execute_commands_test.go` are absent. `requirements.md` lines 31-32 require both as consolidated runner and command test files; `sprint_commands_test.go` has zero execute cases. | `requirements.md` lines 31-32; `ls internal/sprint/execute_test.go` → not found; `ls internal/app/sprint_execute_commands_test.go` → not found | Create `internal/sprint/execute_test.go` covering fake-runtime success-with-evidence → `complete`, success-without-evidence → `failed` with `missing-evidence`, runtime error → `failed` with `runtime-failed`, cancellation → `cancelled`, diagnostic-only completion → `complete` with `diagnostic-only-completion`, stale-running recovery on resume, resume skipping `complete` tasks, dry-run prompt preview, and `executeModelSelection` precedence including empty-model rejection. Create `internal/app/sprint_execute_commands_test.go` covering help, `--task`, `--dry-run`, `--resume`, `--model`, invalid prerequisites, invalid stage, invalid `--to execute`, exit-code mapping (`ExitValidation`/`ExitPartial`/`ExitRuntime`), and text rendering |
| R-2 | high | `sprintHelp()` text at `internal/app/sprint_commands.go:519-552` line 549-550 still says "Supports planning stages through plan.md only. It does not execute implementation…". `sprintValidateHelp`, `sprintPromptHelp`, `sprintFlowHelp` (lines 564-606) do not list `execute`. No `sprintExecuteHelp()` and no `--help` dispatch for the `execute` subcommand at `sprint_commands.go:32-44`. | R-2 from prior review; Documentation review; CLI Surface review | Update `sprintHelp()`, `sprintValidateHelp()`, `sprintPromptHelp()`, `sprintFlowHelp()`. Add `sprintExecuteHelp()` and dispatch `execute --help` at `sprint_commands.go:32-44`. Add an `Execute` section to `docs/cli-reference.md` and an `Execute Sprint Tasks` section to `docs/user-guide.md` documenting `--task`/`--dry-run`/`--prompt`/`--resume`/`--model` flags and exit-code semantics |
| R-3 | high | `renderSprintStatus` (`sprint_commands.go:395-408`) and `StatusSummary` do not surface execute readiness/progress. `PlanningStages()` (`internal/sprint/domain.go:165-174`) excludes `StageExecute`. | R-3 from prior review; Observability and Workflows reviews | Extend `PlanningStages()` (or add a parallel `ExecuteStages()`) with `StageExecute`. Extend `DeriveStages` (`service.go:1046-1088`) to read `.run-state.json` and `execute.md` without invoking runtime and to emit a derived `StageExecute` row. Update `renderSprintStatus` to print counts by `ExecuteTaskStatus`, last diagnostic per task, and the `.run-state.json` pointer. Status must not invoke runtime |
| R-4 | high | `mapSprintError` (`internal/app/sprint_commands.go:328-341`) does not classify `ErrExecuteRunStateMissing/Malformed/Unsupported`; they fall through to `ExitWorkspace`. | R-4 from prior review; Errors review | Extend `mapSprintError` with `errors.Is` cases for the three sentinels and map them to `ExitValidation` (or a dedicated `ExitData` if/when the taxonomy grows). Replace the `strings.Contains` heuristics in the execute switch at `sprint_commands.go:176-181` with `errors.Is` checks against new `ErrExecuteTaskFailure`/`ErrExecuteRuntimeFailure` sentinels |
| R-5 | high | `redactPlanning` (`internal/platform/config/redaction.go:48-56`) does not redact `ExecuteModel`/`ExecuteVariant` (or any `*Variant` field). The redaction test at `config_test.go:69-83` does not cover these fields. | R-5 from prior review; Configuration, LLM Evaluation, Security reviews | Add `p.ExecuteModel = RedactValue("planning.execute_model", p.ExecuteModel)` and `p.ExecuteVariant = RedactValue("planning.execute_variant", p.ExecuteVariant)`, plus the six other `*Variant` fields. Extend `TestRedactSensitiveValues` with marker-bearing execute_model/execute_variant assertions |
| R-6 | high | No structured log events emitted at execute task transitions. `internal/platform/logging.Logger` is defined but never imported or instantiated; `grep -r 'platform/logging' .` returns zero matches. | Observability review (`OBS-CORE-001` blocker) | Inject a `Logger` into `Service` (e.g. `Service.WithLogger`), instantiate it in `runSprint`/`sprintRuntimeService` using `deps.logging`, and emit `log.Info` at task start, completion, failure, cancellation, stale-running recovery. Fields should include `project`, `sprint`, `stage=execute`, `task_id`, `attempt`, `run_id`, `model_source`. Add `--verbose` to elevate to debug |
| R-7 | high | `Service.executeModelSelection` (`internal/sprint/execute.go:330-344`) duplicates the precedence with a different fallback (`"provider/model"` literal, drops `models.primary`) and silently returns it; `PrepareExecute` only emits a `missing execute model` finding when the empty model is reached. The public `ResolveExecuteModel` in `execute_model.go:15-30` is the only path covered by tests. | Configuration review (REQ-23-41); Technical Handbook review | Delete the duplicated chain in `execute.go:330-344` and route `Service.executeModelSelection` through `ResolveExecuteModel` (or extract a single helper). Add a `validateExecuteModel` flow that rejects unknown stage keys and empty values with field-path diagnostics (planning.execute_model → planning.plan_model → models.primary → models.default) so operators see exactly which field is misconfigured without leaking values |
| R-8 | high | `ExecuteDiagnostic.Message` and `ExecuteEvidence.Summary` are persisted unredacted. `safeError` (`internal/sprint/flow.go:139-151`) strips CR/LF/NUL and truncates at 180 chars but does not redact secret patterns. `config.RedactValue` is never invoked for diagnostics/evidence. | Observability review (`OBS-PII-001`); Errors review (`ERR-REDACT-001`); Security review (`SEC-INPUT-001` adjacent) | Introduce a redactor helper that runs `config.RedactValue` and common secret patterns (e.g. `(?i)(sk-\|gho-\|ghp-\|xox[abp]-\|AKIA[0-9A-Z]{16}\|bearer\s+[A-Za-z0-9._\-]+)`) before persisting. Apply at `execute.go:103, 127, 130, 134, 141, 273` and inside `WriteExecuteSummary`. Consider removing the 180-char truncation for persisted diagnostics |
| R-9 | medium | `internal/sprint/doc.go:3-13` still declares "only the governed planning chain through plan.md" and forbids "execute implementation work", contradicting the new sprint-owned execute semantics, `execute.md`, and `.run-state.json` artifacts. | Architecture and Documentation reviews | Update `internal/sprint/doc.go` to include `execute.md` and `.run-state.json` in the governed artifact list, document `execute` as a sprint-owned stage, and replace "execute implementation work" with the narrower "Git mutation / smoke / review / issue tracking / TUI / hosted/browser" exclusions from `execute_target.go:54-61` |
| R-10 | medium | `docs/cli-reference.md:223-268`, `docs/user-guide.md:215-232`, `docs/recovery.md`, and `docs/release-checklist.md:84` (which already enumerates `ultraplan sprint --help`) are not updated for the new execute stage. | Documentation review (`DOC-PUBLIC-001` / `DOC-OPS-001`) | Add an `Execute` section to `docs/cli-reference.md`, an `Execute Sprint Tasks` section to `docs/user-guide.md`, an `Execute Recovery` subsection to `docs/recovery.md` covering `.run-state.json` inspection and stale-running recovery, and ensure `docs/release-checklist.md` requires a help regenerator pass after sprint command surface changes |

### 5.2 Future Follow-ups (Deferred with Justification)

| ID | Severity | Issue | Justification |
| --- | --- | --- | --- |
| F-1 | medium | `ValidateExecuteWorkdir` does not resolve symlinks (`filepath.EvalSymlinks`). | REQ-23-44 requires containment; lexical containment is in place; symlink resolution is hardening. Future sprint should call `EvalSymlinks` on both target and workdir before the `inside()` check and recompute containment on the resolved real paths |
| F-2 | medium | Plan fingerprint drift is not detected on resume. | `PlanFingerprint` is recorded and validated by the schema (`execute_state.go:128-137`) but is not compared against a recomputed fingerprint at load time. The sprint explicitly accepts this as a follow-up (`reasoning.md` Decision 3). Future sprint should add explicit mismatch detection that fails closed unless a force flag is set, and compute the fingerprint only after a successful plan read (replacing `mustReadPlan` which silently swallows read errors) |
| F-3 | medium | `PlanningStages()` does not include `StageExecute`. | The decision to include execute in `PlanningStages()` is partially captured by R-3 above; for the flow-chain/diagnosis aspect (e.g. `flowExecuteSuccessStages`/`flowExecuteFailedStages` builders) the change is a follow-up |
| F-4 | medium | `flowExecuteSuccessStages`/`flowExecuteFailedStages` builders are not added to `flow.go`. | Same as F-3; per-stage flow summaries are deferred |
| F-5 | medium | `runtime_validation.go` directly imports `github.com/Antonio7098/agentwrap` and constructs `agentwrap.ValidationSpec` types. | Pre-existing (mtime predates Sprint 23), enabled by `internal/platform/runtime/runtime.go:30` exposing `agentwrap.ValidationSpec` directly. A future refactor should introduce a generic `pruntime.ValidationSpec` in the platform layer and migrate `runtime_validation.go` off agentwrap |
| F-6 | medium | `internal/platform/runtime/policy.go:4-14` exports `MetadataStudy`/`MetadataDimension` and other unused mode-specific metadata constants. | Pre-existing, unused outside the platform package. A cleanup sprint should move consumed metadata constants out of platform/runtime into their owning modules |
| F-7 | medium | Stage-key and empty-value validation for `planning.execute_model` is incomplete: unknown stage keys are not rejected explicitly with field-path diagnostics. | The current behavior falls back through `models.primary`/`models.default` in `ResolveExecuteModel` and via a hard-coded literal in `Service.executeModelSelection`. R-7 above addresses the duplication and silent fallback; full stage-key validation is a hardening follow-up |
| F-8 | medium | No golden tests for `execute.md` summary output or for rendered execute prompts. | Test infrastructure exists for state/plan/target/model; summary and prompt renderers can be covered in a future sprint once the surface stabilizes |
| F-9 | medium | `ExecuteResult.Runtime` retains the full runtime result at the service boundary. | Persisted state is safe via `runtimeSummary`; the service-boundary retention is not blocking. Future sprint should either drop the field or replace it with a redacted summary type (`json:"-"`) |
| F-10 | medium | `runtimeSummary` discards `pruntime.Result.Usage`/`EstimatedCost`/`StartedAt`/`FinishedAt`. | Persistence-and-cost observability for execute-stage runs is reduced; future sprint should add bounded `TokenUsage`/`EstimatedCost`/`DurationMS` fields populated from the platform Result and a round-trip test |
| F-11 | medium | `WriteExecuteSummary` uses `os.WriteFile` directly, bypassing the atomic temp+rename pattern. | Atomic summary writes can be added in a follow-up by extracting an `atomicWriteFile` helper in `state.go` |
| F-12 | medium | No idempotency key is passed to the runtime boundary; the runtime is not told a task is a retry. | A future execute-state safety refinement should derive an `IdempotencyKey` from `sprint`+`PlanFingerprint`+`task.ID` (or a prior run id) and pass it via `runtimeReq.Metadata`, with documentation of the idempotency stance |
| F-13 | medium | No compensation or reconciliation strategy for partial external side effects on cancellation/failure. | Decision 6 and Decision 5 already accepted this; a future sprint that adds conflict controls or a worktree-per-task strategy would also handle compensation |
| F-14 | low | `safeError` truncates at 180 chars but does not redact secret patterns. | R-8 above addresses redaction; keeping a separate short cap for operator-facing render output is acceptable |
| F-15 | low | Strict validator rejects `Running` tasks without `StartedAt`; the reconciler is the recovery path. Intentional per REQ-23-47 but worth documenting. | Document the validator/reconciler boundary in code comments; future status path should call `reconcileExecuteState` rather than the strict validator |

## 6. Decision Conformance

`reasoning.md` records 10 decisions. Each was checked against the implementation:

| Decision | Conformance |
| --- | --- |
| 1. Sprint Owns Execute Semantics | conformant. Task extraction/IDs in `execute_plan.go`; run-state in `execute_state.go`; target safety in `execute_target.go`; model precedence in `execute_model.go` and `execute.go:330-344`; prompt rendering and summary in `execute.go:196-254`; reconciliation in `execute.go:267-285`. Doc.go is stale (R-9) |
| 2. Execute Only From Validated Plan Tasks | conformant. `Service.prepareExecute` (`execute.go:164-194`) calls `ExtractExecutePlanTasks` after validating the target and reading `plan.md`; deterministic IDs from `deterministicExecuteTaskID`; duplicate/ambiguous/unsupported tasks rejected with findings |
| 3. Versioned Run State Is The Durable Source Of Truth | conformant with F-2 caveat. `ExecuteRunStateSchemaVersion = 1`; strict load validation; same-directory atomic writes; `reconcileExecuteState` recovers stale running tasks; complete tasks preserved on resume. Plan fingerprint drift not yet detected |
| 4. Completion Requires Evidence Or Explicit Safe Diagnostic | conformant in code (`execute.go:124-142`). Test coverage is missing (R-1) |
| 5. Target Repository And Runtime Diagnostics Are Safety-Gated | conformant with caveats. Target pinned to `ApprovedExecuteTargetPath`; `ValidateExecuteWorkdir` enforces lexical containment; `ExecuteSafetyInstructions` is invoked from `RenderExecutePrompt`. Symlink resolution missing (F-1); redaction incomplete (R-5, R-8); in-memory payload retention (F-9) |
| 6. Execute Runs Conservatively And Resumes Safely | conformant. Serial execution (`execute.go:92-149`); no goroutines; complete tasks skipped on resume; stale running recovered with diagnostic |
| 7. Stage Model Resolution Is Deterministic And Visible | conformant with caveats. `Service.executeModelSelection` implements `override → planning.execute_model → planning.plan_model → runtime.config → "provider/model"` literal fallback instead of `models.primary → models.default`; source surfaced via `Selection.Source`; field-path diagnostics partial (R-7, F-7) |
| 8. Status, Prompt, Validate, Flow, And Execute Commands Stay Thin | conformant in code with caveats. `sprint_commands.go:158-184` execute branch is thin (~25 lines); `parseSprintExecuteArgs` and `renderSprintExecute` are focused. `renderSprintStatus` is missing the execute summary (R-3); help text is stale (R-2) |
| 9. Deferred Behaviors Are Explicitly Rejected | conformant. `ExecuteSafetyInstructions` lists `smoke.md`, `review.md`, `issues.md`, `Git mutation`, `hosted/browser`, `TUI`, `cross-sprint scheduling` as forbidden; `RenderExecutePrompt` embeds them; `plan.go:166-167` rejects `git commit`/`git push` and `--to smoke/review/issues` mentions in plan content. Help text is stale (R-2); doc.go is stale (R-9) |
| 10. Verification Uses Deterministic Fake-Runtime Evidence | partially conformant. Verification commands pass. Fake-runtime lifecycle tests and command-level tests are missing (R-1) |

## 7. Plan Execution

| Task | Plan status | Implementation evidence | Conformance |
| --- | --- | --- | --- |
| 1 — Domain types & state schema | marked `[x]` | `domain.go:58-129`, `execute_state.go:14-183` | conformant |
| 2 — Prerequisite validation & plan extraction | marked `[x]` | `execute_plan.go:23-143`, `service.go:257 ValidateExecute` | conformant |
| 3 — Target resolver & safety | marked `[x]` | `execute_target.go:13-61`, `execute_target_test.go` | conformant for lexical; F-1 (symlinks) deferred |
| 4 — Stage model resolution | marked `[x]` (was `[ ]` partial) | `execute_model.go:15-30`, `execute.go:330-344`, `config.go:50-51, 113` | partial (R-5, R-7, F-7) |
| 5 — Execute prompts & runtime requests | marked `[x]` | `execute.go:196-221 RenderExecutePrompt`, `runtimeRequest` at `:115-117` | conformant |
| 6 — Conservative runner & durable transitions | marked `[x]` | `execute.go:57-162 Service.Execute` | conformant for code; test coverage missing (R-1) |
| 7 — Execute summary & runtime-free status | marked `[x]` | `execute.go:223-254 WriteExecuteSummary` | partial — summary implemented; status not wired (R-3); non-atomic summary write (F-11) |
| 8 — Execute flow & validation integration | marked `[x]` | `flow.go:153-157 validateFlowTarget`, `execute.go:15 StageExecute`, `sprint_commands.go:204-205, 237-239` | partial — flow validation accepts execute; `PlanningStages()` does not include `StageExecute` (F-3) |
| 9 — App/CLI commands | marked `[x]` | `sprint_commands.go:158-184, 86-87, 118-119, 368-393, 478-508` | conformant in code; stale help text (R-2); no `sprintExecuteHelp()`; `sprint_execute_commands_test.go` missing (R-1); mapSprintError gap (R-4) |
| 10 — Deferred behavior exclusions & dependency boundaries | marked `[x]` | `execute_target.go:54-61`, `execute.go:215-218`, `plan.go:166-167` | conformant; regression test for non-goal stages at CLI layer is part of the missing `sprint_execute_commands_test.go` (R-1) |
| 11 — Verification & review evidence | marked `[x]` | `go test ./...` PASS, `go test -race ./...` PASS, `go build ./cmd/ultraplan` PASS | conformant for build/test; review file produced |

The plan checklist is fully marked `[x]`; the implementation matches the contracted surface. The unaddressed Required Fixes are below the level of plan-task completion (test coverage, help text, redaction, status surfacing) and do not invalidate any plan task.

## 8. Deviations

| Deviation | Reason | Risk |
| --- | --- | --- |
| `internal/sprint/execute.go` is implemented but was not listed as a single file in the original plan task list (it was implicitly covered by Tasks 5, 6, 7). | The implementation split orchestration into `execute.go` for the use-case surface and `execute_state.go`/`execute_plan.go`/`execute_target.go`/`execute_model.go` for the cohesive helpers; the plan grouped state and orchestration separately. | Low. The file is owned by `internal/sprint` and respects all dependency rules. R-9 asks the doc.go to acknowledge `execute.go` as the orchestration home |
| `Service.executeModelSelection` (`execute.go:330-344`) hard-codes a precedence list that differs from `ResolveExecuteModel` in `execute_model.go:15-30` and uses a literal `"provider/model"` fallback. | The runner uses an in-service `stageRuntime` map for stage-specific overrides; the standalone helper uses `c.Planning.*` fields directly. The two resolvers have drifted. Both surface `Selection.Source`. | Medium. Only `ResolveExecuteModel` is exercised by tests; the in-service path is not regression-protected. R-7 asks to unify |
| `reconcileExecuteState` (`execute.go:267-285`) merges old task state into new task records by ID only, without explicit plan-fingerprint comparison. | The plan acknowledges plan fingerprint drift as a follow-up (`reasoning.md` Decision 3, F-2). | Medium. If `plan.md` is edited between runs in a way that changes task IDs, the merge silently drops old tasks and creates new pending ones. F-2 asks for explicit mismatch detection that fails closed unless a force flag is set |
| `WriteExecuteSummary` (`execute.go:223-254`) uses `os.WriteFile` directly rather than the atomic temp+rename pattern used for `.run-state.json` and `flow-state.json`. | Implementation chose simplicity for a reviewable Markdown summary; durable state is the JSON, summary is the artifact. | Low-to-medium. A torn write to `execute.md` could leave it half-written. F-11 asks for an atomic write helper |
| `ExecuteResult.Runtime` (`execute.go:33, 119`) retains the full `pruntime.Result` slice in memory at the service boundary. | The runner passed the runtime result through to `result.Runtime` for callers that want raw events; persisted state is safe via `runtimeSummary` (`:287-289`). | Medium. Today `renderSprintExecute` (`sprint_commands.go:478-508`) does not print `result.Runtime`, but the field is exported without `json:"-"` so any future marshal would leak full `Events.Payload`/`Attempts`/`Memory`. F-9 asks to replace it with a redacted summary type |
| Stale-running recovery at `execute.go:100-104` and `reconcileExecuteState` (`execute.go:267-285`) emit two distinct messages for conceptually the same recovery. | Both paths exist for safety; deduplicating them risks losing one. The reasoning (`reasoning.md`) does not require a single code path. | Low. Both messages end up in `Diagnostics[]` as separate entries with the same `stale-running` code. F-15 asks to document the divergence |
| `internal/sprint/runtime_validation.go` directly imports `github.com/Antonio7098/agentwrap` and constructs `agentwrap.ValidationSpec` types. | Pre-existing pattern enabled by `internal/platform/runtime/runtime.go:30` exposing the agentwrap type directly through the platform boundary. Sprint 23 did not introduce this; a future refactor should add a generic `pruntime.ValidationSpec` and migrate the caller. | Low for sprint behavior; pre-existing boundary concern surfaced by Architecture and LLM Runtime reviews (F-5) |
| `internal/platform/runtime/policy.go` exports `MetadataStudy`/`MetadataDimension` and other unused mode-specific metadata constants. | Pre-existing, unused outside the platform package. Sprint 23 inherits but does not worsen the boundary leak. | Low. Future cleanup (F-6) |

## 9. Final Assessment

**Outcome**: **Accepted with follow-ups**.

**Rationale**:

- The execute stage's core contract is implemented: sprint-owned semantics in `internal/sprint`; generic runtime boundary unchanged; `internal/platform/runtime` remains product-agnostic; `.run-state.json` is versioned, atomic, strictly validated, and recoverable; target paths are pinned to the approved root with lexical containment; runtime success is gated by evidence or explicit diagnostic; stale `running` tasks are recovered; complete tasks are skipped on resume; the `execute` subcommand, `validate execute`, `prompt execute`, `flow --to execute`, `--task`, `--dry-run`, `--resume`, and `--model` are wired; `execute.md` is generated with plan citation, task counts, evidence, and diagnostics.
- `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan` all pass.
- Ten required fixes (R-1 through R-10) and fifteen future follow-ups (F-1 through F-15) are documented. The required fixes are bounded: they are documentation/test/redaction polish, missing structured logging, missing status rendering, and missing runner/CLI test coverage; they do not require architectural change. The future follow-ups are explicit hardening items accepted as deferred per `reasoning.md`.
- No new scope or architecture is introduced; no forbidden behavior is present (smoke, review, issues, Git mutation, TUI, hosted/browser, cross-sprint scheduling are all excluded by `ExecuteSafetyInstructions`, plan validators, absence of related code paths, and prompt-level safety instructions).
- This is the second pass through the Review Sprint Protocol for sprint 23-execute-stage. The first pass (saved as `review.md` earlier) identified R-1 through R-7 and F-1 through F-8. Re-running the protocol produced independent subagent reviews that corroborated R-1 through R-7 and added R-8 through R-10 plus F-9 through F-15. The Required Fixes that were open in the prior review are still open in this review (R-1 through R-7 are unresolved). The new Required Fixes (R-8 through R-10) reflect deeper subagent findings on structured logging, diagnostic redaction, and project-wide doc propagation.

**Sign-off conditions for next sprint**:

- Land R-1 through R-10 before the next sprint begins.
- Address or re-scope at least R-1 (the missing `internal/sprint/execute_test.go` and `internal/app/sprint_execute_commands_test.go`) and R-3 (execute in `renderSprintStatus` and `PlanningStages()`) before the implementation opens a new sprint, because both materially affect verification coverage and operator visibility.
- Address R-5 (Extend `redactPlanning` to cover `ExecuteModel`/`ExecuteVariant`) before any real runtime exposes configured model values through `ultraplan config show`.
- Address R-6 (Wire `internal/platform/logging.Logger`) before any production-runtime smoke runs against the execute stage.
- Track F-1 through F-15 as hardening items; if any of F-1 (symlink resolution), F-2 (plan-fingerprint drift), F-9 (in-memory payload retention), or F-12 (idempotency key) becomes a blocker for a downstream sprint, promote it to a Required Fix.
- Update `projects/ultraplan-go/docs/ARCHITECTURE.md`, `docs/PRD.md`, and `docs/TRD.md` to mention the new execute stage, `.run-state.json`, and the execute/runtime seam boundaries (covered by R-10 indirectly).
