# Sprint Requirements: 23-execute-stage

> Project: `ultraplan-go`
> Sprint: `23-execute-stage`
> Purpose: the authoritative, human-readable sprint contract. All other sprint artifacts must satisfy these requirements.

## Sprint Goal

Execute validated `plan.md` implementation tasks through the generic runtime boundary with durable task state, safe target-repository boundaries, resumability, and clear diagnostics.

## Required Outputs

| Output | Path | Description |
| --- | --- | --- |
| Requirements | `projects/ultraplan-go/sprints/23-execute-stage/requirements.md` | This sprint contract derived from the roadmap, project index, project docs, and prior sprint reviews. |
| Sprint index | `projects/ultraplan-go/sprints/23-execute-stage/sprint-index.md` | Selected contracts, evidence reports, reasoning templates, review protocols, and explicitly excluded context for execute-stage work. |
| Technical handbook | `projects/ultraplan-go/sprints/23-execute-stage/technical-handbook.md` | Evidence-backed implementation guidance for the execute stage without final implementation decisions. |
| Architecture reasoning | `projects/ultraplan-go/sprints/23-execute-stage/reasoning/architecture.md` | Area reasoning for module ownership, dependency direction, runtime boundary use, persistence, and target-repository safety when the Architecture template is selected. |
| Final reasoning | `projects/ultraplan-go/sprints/23-execute-stage/reasoning.md` | Final decisions, assumptions, risks, and expected evidence for execute-stage implementation. |
| Sprint plan | `projects/ultraplan-go/sprints/23-execute-stage/plan.md` | Ordered implementation plan whose tasks trace to `reasoning.md` decisions and evidence. |
| Execute summary artifact | `projects/ultraplan-go/sprints/23-execute-stage/execute.md` | Editable Markdown summary of execute task counts, terminal states, completion evidence, and failed-task diagnostics after execute runs. |
| Execute run state artifact | `projects/ultraplan-go/sprints/23-execute-stage/.run-state.json` | Versioned, atomic, resumable execute task state written only when execute begins or resumes. |
| Flow state artifact | `projects/ultraplan-go/sprints/23-execute-stage/flow-state.json` | Versioned sprint flow state that can represent execute readiness, progress, success, failure, and prerequisite failures. |
| Execute domain implementation | `../ultraplan-go/internal/sprint/execute.go` | Sprint-owned execute task extraction, task validation, deterministic task IDs, execution orchestration, and summary generation. |
| Execute persistence implementation | `../ultraplan-go/internal/sprint/execute_state.go` | Sprint-owned `.run-state.json` schema, strict load validation, atomic writes, stale running-task recovery, and resume reconciliation. |
| Sprint flow implementation updates | `../ultraplan-go/internal/sprint/flow.go` | Flow-stage updates that allow `execute` after valid prerequisites through `plan` and reject deferred stages. |
| Sprint prompt implementation updates | `../ultraplan-go/internal/sprint/prompts.go` | Deterministic execute prompt rendering and prompt preview data using workspace-relative paths and safe target scope. |
| Sprint service implementation updates | `../ultraplan-go/internal/sprint/service.go` | Service use cases for execute validation, prompt preview, flow through execute, task execution, resume, and status integration. |
| Sprint validation implementation updates | `../ultraplan-go/internal/sprint/validation.go` | Execute prerequisite validation, target-path safety checks, plan task extractability checks, and `execute.md` validation. |
| App command implementation updates | `../ultraplan-go/internal/app/sprint_commands.go` | Thin CLI wiring for `validate execute`, `prompt execute`, `flow --to execute`, `execute`, task selection, dry-run, and status rendering. |
| Execute package tests | `../ultraplan-go/internal/sprint/execute_test.go` | Unit tests for task extraction, deterministic IDs, prerequisites, prompt rendering, state transitions, resume, evidence/diagnostic gating, and non-goals. |
| Execute command tests | `../ultraplan-go/internal/app/sprint_execute_commands_test.go` | Command tests for help, execute invocation, dry-run/prompt behavior, invalid prerequisites, task selection, exit codes, and status output. |

## Acceptance Criteria

- [ ] `projects/ultraplan-go/sprints/23-execute-stage/sprint-index.md` validates as a subset of `projects/ultraplan-go/project-index.md` and explicitly excludes smoke, review automation, issue tracking, Git mutation, and cross-sprint scheduling.
- [ ] `projects/ultraplan-go/sprints/23-execute-stage/technical-handbook.md` cites only selected evidence and does not make final implementation decisions.
- [ ] `projects/ultraplan-go/sprints/23-execute-stage/reasoning.md` records execute-stage decisions, assumptions, risks, rejected alternatives, and expected verification evidence.
- [ ] `projects/ultraplan-go/sprints/23-execute-stage/plan.md` cites `reasoning.md`, includes executable tasks with expected evidence, and does not invoke implementation work before the execute stage.
- [ ] `ultraplan sprint <project> <sprint> execute` and `ultraplan sprint <project> <sprint> flow --to execute` require valid prerequisites through `plan.md` before runtime-backed implementation begins.
- [ ] Runtime model resolution supports a global/default sprint model and stage-specific overrides for `sprint-index`, `technical-handbook`, `area-reasoning`, `reasoning`, `plan`, and `execute`; stage-specific values win over global/default values.
- [ ] Execute task extraction reads only validated `plan.md` entries, produces deterministic task IDs for the same validated plan content, and preserves traceability back to the plan task and reasoning decision.
- [ ] Runtime-backed execute work uses only the generic `internal/platform/runtime` boundary through sprint/app seams; `internal/sprint` does not directly invoke OpenCode, shell out to runtime tools, parse native runtime streams, or import study services.
- [ ] Target repository paths are explicit, workspace-safe, and constrained to `../ultraplan-go` unless a future explicit configuration permits another safe target.
- [ ] Runtime success alone is insufficient: a task reaches `complete` only when expected evidence is present or an explicit safe diagnostic explains why evidence cannot be machine-validated.
- [ ] `projects/ultraplan-go/sprints/23-execute-stage/.run-state.json` is versioned, written atomically, loaded strictly, resumable, and records `pending`, `running`, `complete`, `failed`, and `cancelled` task states with attempts, timestamps, diagnostics, and runtime metadata summaries where available.
- [ ] Interrupted or stale `running` execute tasks are recovered to a retryable or failed state with a diagnostic, without redoing `complete` tasks unless explicitly forced.
- [ ] Sprint status reflects execute readiness and progress from `flow-state.json`, `.run-state.json`, `execute.md`, and prerequisite validation without requiring a runtime call.
- [ ] `projects/ultraplan-go/sprints/23-execute-stage/execute.md` cites `plan.md`, summarizes task counts and terminal states, records completion evidence, and records actionable failed-task diagnostics or pointers to `.run-state.json`.
- [ ] Normal verification passes with fake runtimes: `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan` from `../ultraplan-go`.
- [ ] No execute path invokes smoke investigation, automated review, local issue tracking, Git add/commit/push/branch/PR mutation, or cross-sprint scheduling behavior.

## Non-Goals

- Running smoke investigations or producing `smoke.md` / `smoke.json`.
- Producing automated conformance `review.md` artifacts or implementing review automation.
- Creating or managing issues, issue files, issue JSON, assignments, scheduling, or project-management workflow.
- Automatically running Git add, commit, push, branch creation, PR creation, checkout, reset, or any other Git state mutation.
- Building a generic workflow engine, cross-sprint scheduler, hosted service, browser UI, or TUI surface.
- Reusing study services, source/dimension models, study report validation, rating parsing, summary generation, or study run-loop scheduling for execute behavior.
- Replacing `github.com/Antonio7098/agentwrap` / `agentwrap/opencode` with a competing runtime contract.

## Constraints

- Use workspace-relative paths in generated sprint artifacts and diagnostics.
- Keep generated sprint artifacts editable Markdown, except versioned JSON state files `flow-state.json` and `.run-state.json`.
- `internal/sprint` owns execute task extraction, prompt rendering, execute validation, execute run state, execute summary generation, and flow through execute.
- Preserve dependency direction: `sprint -> project`, `sprint -> workspace`, `sprint -> platform/runtime`, `project -> workspace`, `platform/* -> no product modules`, and `study -> no project or sprint modules`.
- Keep `internal/platform/runtime` generic; it may know prompts, working directories, models, timeouts, permissions, expected outputs, events, and results, but not project catalogs, sprint stages, plan semantics, execute task states, or artifact validators.
- Use `github.com/Antonio7098/agentwrap` and `agentwrap/opencode` through the existing platform runtime integration; do not parse OpenCode stdout/stderr directly in product code.
- Treat runtime success as necessary but never sufficient; artifact/state validation and evidence/diagnostic checks are product-owned gates.
- Persist `.run-state.json` with same-directory atomic write behavior and reject malformed or unsupported schema versions with actionable diagnostics.
- Redact secrets and unsafe native runtime payloads from logs, diagnostics, `execute.md`, and JSON/text status output.
- Normal tests must be deterministic, offline, fake-runtime based, and must not require OpenCode, provider credentials, network access, or mutation of Git state.

## Dependencies

| Prior Sprint / Output | Required For | Notes |
| --- | --- | --- |
| `projects/ultraplan-go/project-index.md` | Context selection and sprint-index validation | Execute-stage selections must come from the project catalog. |
| `projects/ultraplan-go/roadmap.md` | Sprint scope and non-goals | Sprint 23 is the first Planning Phase 2 sprint that includes controlled implementation execution. |
| `projects/ultraplan-go/docs/PRD.md` | Product behavior and non-goals | Requires controlled execute from validated `plan.md`; defers smoke, reviews, issues, and Git mutation. |
| `projects/ultraplan-go/docs/TRD.md` | Technical behavior, state, runtime, and command requirements | Defines execute run state, model resolution, validation gates, runtime boundary, and command surface. |
| `projects/ultraplan-go/docs/ARCHITECTURE.md` | Module ownership and dependency rules | Requires execute behavior in `internal/sprint` and generic runtime behavior in `internal/platform/runtime`. |
| `projects/ultraplan-go/sprints/16-project-domain-and-index/review.md` | Project catalog prerequisite | Confirms project discovery and project-index validation are accepted. |
| `projects/ultraplan-go/sprints/17-artifact-domain-and-flow-state/review.md` | Sprint artifact and flow-state prerequisite | Confirms sprint discovery, artifact paths, strict `flow-state.json`, and status refresh are accepted. |
| `projects/ultraplan-go/sprints/18-select-stage/review.md` | Sprint-index prerequisite | Confirms selected-context validation and flow through `sprint-index` are accepted. |
| `projects/ultraplan-go/sprints/19-distill-stage/review.md` | Technical-handbook prerequisite | Confirms selected-evidence distillation and flow through `technical-handbook` are accepted. |
| `projects/ultraplan-go/sprints/20-reason-stage/review.md` | Reasoning prerequisite | Confirms area/final reasoning generation, validation, and flow through `reasoning` are accepted. |
| `projects/ultraplan-go/sprints/21-plan-stage/review.md` | Plan prerequisite | Confirms `plan.md` validation, prompt rendering, and flow through `plan` are accepted. |
| Sprint 22 planning documentation release gate | Documentation and release baseline | Roadmap marks Sprint 22 complete; no review file exists in this workspace, so Sprint 23 must not assume additional carry-forward decisions beyond the roadmap status. |
| `../ultraplan-go` | Target implementation directory | Execute must constrain implementation work to this explicit target repository unless later configuration safely overrides it. |

## Review Expectations

| What | How Verified |
| --- | --- |
| Planning artifact chain | Validate `requirements.md`, `sprint-index.md`, `technical-handbook.md`, `reasoning.md`, and `plan.md`; inspect traceability from requirements to plan tasks. |
| Execute prerequisites | Tests and command checks show execute refuses missing, invalid, or unvalidated prerequisites through `plan.md` before runtime-backed work starts. |
| Runtime boundary | Import review and tests show sprint execute uses only the generic runtime seam and does not call OpenCode, shell runtime commands, agentwrap adapters, or study services directly. |
| Model resolution | Unit tests cover global/default model fallback, stage-specific overrides, execute-specific override precedence, unknown stage keys, and empty model diagnostics. |
| Task extraction and traceability | Unit tests cover deterministic IDs, duplicate/ambiguous task handling, plan-task references, evidence expectations, and stable ordering. |
| Target path safety | Unit tests cover target path containment, explicit target directory use, path escape rejection, and workspace-relative diagnostics. |
| Execute run-state durability | Unit tests cover schema versioning, atomic write preservation, malformed/unsupported state rejection, status transitions, stale running-task recovery, resume without redoing complete tasks, and cancellation state. |
| Evidence/diagnostic completion gate | Fake-runtime tests cover runtime success with evidence, runtime success without evidence, runtime failure, validation failure, and explicit diagnostic-only completion cases. |
| Status and summary output | Command tests verify status is runtime-free, reflects execute progress, and `execute.md` cites `plan.md` with counts, terminal states, evidence, and diagnostics. |
| Deferred behavior exclusion | Tests or code review confirm no smoke, review automation, issue tracking, Git mutation, TUI, hosted service, or cross-sprint scheduler commands or stages are introduced. |
| Verification gates | Run `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan` from `../ultraplan-go`; gated real-runtime smoke is optional and not required for normal acceptance. |
