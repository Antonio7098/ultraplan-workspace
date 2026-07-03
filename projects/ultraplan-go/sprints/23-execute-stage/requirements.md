# Sprint Requirements: 23-execute-stage

> Project: ultraplan-go
> Sprint: 23-execute-stage
> Purpose: implement the governed execute stage from the roadmap, PRD, TRD, and architecture docs.

## Sprint Goal

Execute validated `plan.md` implementation tasks through the generic runtime boundary with durable task state, safe target-repository boundaries, resumability, and clear diagnostics.

## Required Outputs

| Output | Path | Description |
| --- | --- | --- |
| Sprint index | projects/ultraplan-go/sprints/23-execute-stage/sprint-index.md | Selected contracts, evidence, reasoning templates, protocols, and excluded context. |
| Technical handbook | projects/ultraplan-go/sprints/23-execute-stage/technical-handbook.md | Evidence-backed technical guidance for the sprint. |
| Sprint reasoning | projects/ultraplan-go/sprints/23-execute-stage/reasoning.md | Final decisions, assumptions, risks, and evidence mapping. |
| Sprint plan | projects/ultraplan-go/sprints/23-execute-stage/plan.md | Ordered implementation plan derived from reasoning. |
| Execute prompt rendering | internal sprint package and CLI prompt paths | Render execute prompts from validated plan tasks without bypassing the runtime boundary. |
| Execute task extraction | internal sprint package | Extract deterministic executable task IDs from valid `plan.md` entries. |
| Execute run state | projects/ultraplan-go/sprints/23-execute-stage/.run-state.json | Versioned, atomic, resumable task state with attempts, diagnostics, and terminal statuses. |
| Execute summary | projects/ultraplan-go/sprints/23-execute-stage/execute.md | Summary of task counts, terminal states, and failed-task diagnostics. |
| Flow through execute | sprint flow/status commands | `flow --to execute`, `execute`, and sprint status reflect execute progress. |

## Acceptance Criteria

- [ ] The sprint index validates against the project catalog.
- [ ] The technical handbook cites selected evidence and avoids final implementation decisions.
- [ ] The reasoning artifact records decisions, assumptions, risks, and evidence.
- [ ] The plan traces tasks to reasoning decisions and completion evidence.
- [ ] `execute` requires valid prerequisites through `plan`.
- [ ] Runtime model selection supports a global/default model and stage-specific overrides for planning stages and execute; stage-specific values win over global/default values.
- [ ] Tasks trace to validated plan entries and preserve deterministic task IDs.
- [ ] Runtime-backed implementation uses only the generic platform runtime boundary.
- [ ] Target repository paths are explicit and workspace-safe.
- [ ] Runtime success alone is insufficient; task completion requires expected evidence or explicit diagnostics.
- [ ] `.run-state.json` is versioned, written atomically, resumable, and records pending, running, complete, failed, and cancelled task states.
- [ ] Sprint status shows execute progress.
- [ ] `execute.md` cites `plan.md`, summarizes task counts and terminal states, and records actionable failed-task diagnostics or pointers to `.run-state.json`.
- [ ] No smoke, automated review, issue tracking, Git add/commit/push, or cross-sprint scheduler behavior is invoked.

## Non-Goals

- Smoke investigation runs.
- Automated conformance review.
- Issue tracking or project-management workflow.
- Automatic Git add, commit, push, branch, or PR mutation.
- Cross-sprint scheduling or a general workflow engine.

## Constraints

- Use workspace-relative paths.
- Select only project-index catalog entries.
- Keep generated artifacts editable Markdown.
- Keep execute behavior owned by `internal/sprint`.
- Keep runtime execution behind `internal/platform/runtime`; do not add sprint semantics to the platform package.
- Constrain target repository writes to an explicit workspace-safe implementation directory unless explicitly configured otherwise.
- Preserve the dependency direction: `sprint -> project`, `sprint -> workspace`, `sprint -> platform/runtime`.

## Dependencies

| Prior Sprint / Output | Required For | Notes |
| --- | --- | --- |
| Project index | Sprint planning | The project catalog must validate before flow generation. |
| Valid `plan.md` | Execute task extraction | Execute may only run after prerequisites through plan validate. |
| Runtime model config | Execute and planning runtime requests | Supports global/default and stage-specific model selection. |
| Flow state | Status and resume visibility | Execute appears only after valid prerequisites through plan. |

## Review Expectations

| What | How Verified |
| --- | --- |
| Planning artifact chain | Run stage validation commands through plan. |
| Execute stage validation | Run execute command/status tests with fake runtime and invalid prerequisite cases. |
| Run-state durability | Unit tests for atomic write/load, malformed state rejection, resumability, stale running-task recovery, and deterministic task IDs. |
| Runtime boundary | Tests prove execute uses the generic runtime adapter and does not shell out through sprint code. |
