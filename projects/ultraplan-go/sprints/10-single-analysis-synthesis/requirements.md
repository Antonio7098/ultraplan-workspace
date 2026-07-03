# Sprint Requirements: 10 Single Analysis and Synthesis

> Project: `ultraplan-go`
> Sprint: `10-single-analysis-synthesis`
> Purpose: the authoritative, human-readable sprint contract. All other sprint artifacts must satisfy these requirements.

## Sprint Goal

Implement `ultraplan study <study> run <dimension-ref> <source-ref>` and `ultraplan study <study> synthesize <dimension-ref>` so one analysis task or one synthesis task can execute through the agentwrap-backed runtime and be accepted only after the expected report artifact exists and passes study-owned validation.

## Required Outputs

| Output | Path | Description |
| --- | --- | --- |
| Sprint requirements | `projects/ultraplan-go/sprints/10-single-analysis-synthesis/requirements.md` | Defines the binding goal, scope, acceptance criteria, constraints, dependencies, and review expectations for this sprint. |
| Sprint index | `projects/ultraplan-go/sprints/10-single-analysis-synthesis/sprint-index.md` | Selects the roadmap entries, source docs, evidence reports, contracts, and review protocols that govern this execution sprint. |
| Technical handbook | `projects/ultraplan-go/sprints/10-single-analysis-synthesis/technical-handbook.md` | Gives implementation guidance for single analysis execution, single synthesis execution, runtime request construction, validation gating, diagnostics, and fake-runtime tests. |
| Execution reasoning | `projects/ultraplan-go/sprints/10-single-analysis-synthesis/reasoning/single-analysis-synthesis.md` | Records reasoning for task ownership, runtime request mapping, validation boundaries, skip behavior, CLI output, and deferred batch/run-loop behavior. |
| Consolidated reasoning | `projects/ultraplan-go/sprints/10-single-analysis-synthesis/reasoning.md` | Summarizes accepted sprint decisions, tradeoffs, risks, and requirement mappings. |
| Implementation plan | `projects/ultraplan-go/sprints/10-single-analysis-synthesis/plan.md` | Provides the ordered task plan and verification checklist for implementation. |
| Sprint review | `projects/ultraplan-go/sprints/10-single-analysis-synthesis/review.md` | Reviews completed implementation against these requirements and records carry-forward decisions. |
| Study analysis execution | `/home/antonioborgerees/coding/ultraplan-go/internal/study/run.go` | Adds study-owned single analysis orchestration: resolve prompt inputs, skip inapplicable Markdown pairs, call the platform runtime, and validate the per-source report. |
| Study synthesis execution | `/home/antonioborgerees/coding/ultraplan-go/internal/study/synthesize.go` | Adds study-owned single-dimension synthesis orchestration: require applicable valid source reports, call the platform runtime, and validate the final report. |
| Study service wiring | `/home/antonioborgerees/coding/ultraplan-go/internal/study/service.go` | Wires the runtime dependency and exposes service methods for single analysis and synthesis without moving workflow logic into the CLI. |
| Study execution domain types | `/home/antonioborgerees/coding/ultraplan-go/internal/study/domain.go` | Adds only the request/result/status fields needed for single execution, validation summaries, skip results, and safe diagnostics. |
| CLI command wiring | `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_commands.go` | Implements `study run` and `study synthesize` command parsing, output rendering, exit status mapping, and help text through thin app-layer adapters. |
| App composition update | `/home/antonioborgerees/coding/ultraplan-go/internal/app/app.go` | Provides the study service with the configured platform runtime while preserving explicit composition and test seams. |
| Study execution tests | `/home/antonioborgerees/coding/ultraplan-go/internal/study/run_test.go` | Tests single analysis and synthesis orchestration with fake runtimes, validation success/failure, missing outputs, runtime failures, and inapplicable Markdown skips. |
| CLI execution tests | `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_run_commands_test.go` | Tests command behavior, stdout/stderr separation, exit statuses, skip messages, validation diagnostics, and fake-runtime execution paths. |

## Acceptance Criteria

- [ ] `ultraplan study <study> run <dimension-ref> <source-ref>` resolves the workspace, study, dimension, and source using existing study resolution behavior.
- [ ] Single analysis uses the existing deterministic analysis prompt builder rather than duplicating prompt composition in the command layer or runtime package.
- [ ] Directory source analysis runtime requests use the source directory as the execution work directory and preserve source-isolation and file-line citation prompt requirements.
- [ ] Markdown document analysis runtime requests use the study or workspace directory as the safe work directory, embed stripped document content through the existing prompt builder, and do not request external code or filesystem exploration.
- [ ] If the selected Markdown source is not applicable to the selected dimension, `study run` exits successfully without invoking the runtime and prints a clear skipped/inapplicable message.
- [ ] Single analysis constructs an agentwrap-backed platform runtime request with prompt text, work directory, provider/model/timeout config, task metadata, permission policy, required health/capability checks, and expected output validation.
- [ ] Single analysis metadata includes task kind, study, dimension number/slug, source name, source kind, output path, runtime, provider, and model without requiring `internal/platform/runtime` to import study types.
- [ ] Runtime success alone is insufficient for analysis success: the expected per-source report must exist, be non-empty, and pass the existing study-owned per-source report validator.
- [ ] Analysis validation failures return the validation exit status and include the failed check names and report path in user-facing diagnostics.
- [ ] Runtime failures return the runtime exit status and preserve classified runtime diagnostics without hiding the original cause chain.
- [ ] `ultraplan study <study> synthesize <dimension-ref>` resolves the workspace, study, and dimension using existing study resolution behavior.
- [ ] Synthesis uses the existing deterministic synthesis prompt builder rather than duplicating report manifest construction in the command layer or runtime package.
- [ ] Synthesis considers only applicable source/dimension pairs and excludes inapplicable Markdown sources from required input checks.
- [ ] Synthesis fails before runtime execution when any applicable per-source report for the selected dimension is missing or invalid, and diagnostics identify each blocking report path.
- [ ] Synthesis runtime requests include prompt text, safe work directory, provider/model/timeout config, task metadata, permission policy, required health/capability checks, and expected output validation.
- [ ] Runtime success alone is insufficient for synthesis success: the expected final report must exist, be non-empty, and pass the existing study-owned final report validator.
- [ ] Synthesis validation failures return the validation exit status and include the failed check names and final report path in user-facing diagnostics.
- [ ] Study execution logic stays in `internal/study`; CLI handlers only parse arguments, call service methods, render output, and map errors to exit statuses.
- [ ] `internal/platform/runtime` remains generic and does not learn study semantics, report semantics, synthesis gating rules, or source applicability rules.
- [ ] The implementation consumes runtime events without allowing event channels to block and always obtains the final runtime result through the platform runtime abstraction.
- [ ] User-facing text output is deterministic enough for command tests and separates normal summaries on stdout from diagnostics on stderr.
- [ ] Any JSON output added by this sprint is explicitly documented and tested; stable public JSON schemas beyond existing command behavior are not required.
- [ ] Fake-runtime tests cover analysis success, synthesis success, runtime failure, successful runtime with missing output, successful runtime with invalid report, inapplicable Markdown skip, and synthesis blocked by invalid or missing source reports.
- [ ] Normal tests do not require OpenCode, provider credentials, network access, or real subprocess execution.
- [ ] Any real OpenCode smoke test added by this sprint is gated by explicit environment variables and skipped by default; it must go through `agentwrap/opencode` via the platform runtime boundary.
- [ ] Existing study discovery, applicability, report validation, prompt composition, run-state, runtime integration, health, config, and status tests continue to pass.
- [ ] `go test ./...` passes in `/home/antonioborgerees/coding/ultraplan-go` without requiring OpenCode or provider credentials.
- [ ] `go build ./cmd/ultraplan` passes in `/home/antonioborgerees/coding/ultraplan-go`.

## Non-Goals

- `ultraplan study <study> run-all`, task matrix execution, bounded worker pools, and multi-task scheduling are not included.
- `ultraplan study <study> run-loop`, resumable batch mutation, stale running task recovery, retry scheduling across process restarts, cancellation persistence, and per-study lock files are not included.
- Automatic synthesis after analysis is not included; `study run` and `study synthesize` remain separate commands in this sprint.
- Summary generation and `summary.csv` updates are not included.
- Code reference extraction is not included.
- Report repair prompts and bounded repair loops are not required unless already supplied by the configured runtime validation wrapper without new product workflow logic.
- Durable agentwrap run-store persistence across process exits is not required for this sprint.
- Stable public JSON schemas for execution results are not required unless explicitly selected in the sprint index.
- Target workflows, sprint planning, and sprint execution remain out of scope.
- UltraPlan must not implement a new OpenCode subprocess supervisor, retry engine, permission translator, event decoder, or validation wrapper that duplicates agentwrap behavior.

## Constraints

- Preserve the dependency direction from `docs/ARCHITECTURE.md`: `internal/study -> internal/platform/runtime -> agentwrap/opencode`; `internal/platform/runtime` must not import product modules.
- Apply the roadmap Runtime / Provider Gate for runtime-capable services: pass `context.Context`, use the platform runtime boundary, preserve correlation metadata, surface safe usage/cost fields when available, and verify with fake-runtime tests.
- Do not apply full Batch / Durable Workflow Gate behavior beyond compatibility with existing task IDs and validation summaries; worker pools, run-loop resume execution, lock files, and long-lived cancellation state are later scope.
- Runtime output validation must be product-owned by `internal/study` and exposed to agentwrap through generic validation hooks or post-run checks; report semantics must not move into `internal/platform/runtime`.
- Runtime requests must fail before launch when required prompt inputs, required source reports for synthesis, config values, health checks, capabilities, or permission policy features are invalid.
- Error translation must preserve cause chains and runtime classifications; do not classify runtime failures by parsing human-readable strings.
- Sensitive config, environment, prompt diagnostics, native runtime payloads, and runtime metadata must be redacted or omitted from user-facing diagnostics unless already marked safe.
- Tests must use fake runtimes for default coverage and must not depend on OpenCode, provider credentials, network access, wall-clock sleeps, or long-running subprocesses.
- Any optional real-runtime smoke coverage must be opt-in, environment-gated, skipped by default in `go test ./...`, and routed through `agentwrap/opencode` via `internal/platform/runtime`.
- Implementation must remain minimal and reviewable; avoid speculative plugin systems, broad registries, a general workflow engine, or new global technical-layer packages.

## Dependencies

| Prior Sprint / Output | Required For | Notes |
| --- | --- | --- |
| Sprint 1: Go Module, CLI Shell, and App Composition | Buildable CLI and app wiring | Required for command integration and `cmd/ultraplan` build continuity. |
| Sprint 2: Workspace, Config, Logging, and Health Skeleton | Workspace/config resolution and redacted diagnostics | Required before runtime execution can resolve settings and emit safe failures. |
| Sprint 3: Study Domain, Listing, and Resolution | Study, dimension, and source lookup | Required for resolving command arguments and stable source/dimension metadata. |
| Sprint 4: Study Initialization From YAML | Study directory and artifact layout | Required for expected source, dimension, report, prompt, and template paths. |
| Sprint 5: Markdown Document Sources and Applicability | Source kind and applicability behavior | Required for inapplicable Markdown skip behavior and synthesis gating. |
| Sprint 6: Report Validation and Rating Parsing | Product artifact validation | Required to prove runtime success is not product success for analysis and synthesis outputs. |
| Sprint 7: Prompt Composition Generation | Deterministic analysis and synthesis prompts | Required to feed runtime requests without duplicating prompt logic. |
| Sprint 8: Run State Persistence and Status | Stable task IDs and state vocabulary | Required as compatibility context; this sprint does not implement batch resume execution. |
| Sprint 9: Agentwrap/OpenCode Runtime Integration | Platform runtime boundary | Required to execute prompts through agentwrap/OpenCode wrappers without direct OpenCode supervision. |
| `projects/ultraplan-go/docs/PRD.md` | Product execution behavior | Governs single analysis, synthesis, validation gates, diagnostics, Markdown source rules, and deferred target/sprint scope. |
| `projects/ultraplan-go/docs/TRD.md` | Technical execution behavior | Governs agentwrap request mapping, validation, policy, permission, event handling, errors, fake-runtime tests, and build requirements. |
| `projects/ultraplan-go/docs/ARCHITECTURE.md` | Package ownership and dependency direction | Governs keeping study workflow behavior in `internal/study` and generic runtime behavior in `internal/platform/runtime`. |
| `projects/ultraplan-go/sprints/05-study-listing-inspection/review.md` | Carry-forward decisions | Confirms Markdown source discovery/applicability are complete. |
| `projects/ultraplan-go/sprints/06-study-run-execution/review.md` | Carry-forward decisions | Confirms report validators and rating parser are complete and runtime execution was deferred. |
| `projects/ultraplan-go/sprints/07-prompt-composition-generation/review.md` | Carry-forward decisions | Confirms prompt builders are deterministic, runtime-free, and ready to feed execution. |
| `projects/ultraplan-go/sprints/08-run-state-persistence/review.md` | Carry-forward decisions | Confirms task/state vocabulary exists while run-loop mutation remains deferred. |
| `projects/ultraplan-go/sprints/09-runtime-integration/review.md` | Carry-forward decisions | Confirms agentwrap runtime integration, wrapper stack, runtime health, event mapping, and fake-runtime coverage are complete; config migration, prompt version metadata, unified diagnostics, metrics/alerts, optional real smoke tests, and local `replace` cleanup remain deferred. |

## Review Expectations

| What | How Verified |
| --- | --- |
| Scope alignment | Review `requirements.md`, `sprint-index.md`, `technical-handbook.md`, `reasoning.md`, and `plan.md` to confirm only single analysis and single synthesis execution were added, with no batch, run-loop, summary, code extraction, target, or sprint-execution scope. |
| Required outputs exist | Check each path listed in Required Outputs and confirm implementation files are present or intentionally replaced by equivalent files recorded in `review.md`. |
| Architecture conformance | Review implementation paths and imports to confirm study execution behavior lives in `internal/study`, CLI wiring stays thin in `internal/app`, platform runtime remains generic, and no global validation/scheduler/prompts/runtime workflow packages were introduced. |
| Analysis command behavior | Run or inspect tests proving `study run` resolves inputs, skips inapplicable Markdown pairs without runtime invocation, invokes the runtime for applicable pairs, validates the expected per-source report, and maps runtime vs validation failures to the correct exits. |
| Synthesis command behavior | Run or inspect tests proving `study synthesize` requires valid applicable per-source reports before runtime execution, excludes inapplicable Markdown pairs, validates the expected final report, and maps runtime vs validation failures to the correct exits. |
| Runtime request mapping | Inspect code and tests proving prompt text, work directory, model/provider/timeout, metadata, health/capability requirements, permission policy, and output validation are passed through the platform runtime boundary. |
| Validation gating | Run or inspect tests proving successful runtime completion does not mark analysis or synthesis successful when the expected report is missing or invalid. |
| Diagnostics and safety | Review text output, errors, and tests for deterministic stdout/stderr separation, safe redaction/omission of sensitive runtime diagnostics, preserved cause chains, and actionable report paths/check names. |
| Runtime-free default tests | Run `go test ./...` without OpenCode/provider environment and confirm execution tests use fakes or skip gated smoke paths by default. |
| Verification commands | Run `go test ./...` and `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go` and record the results in `review.md`. |
