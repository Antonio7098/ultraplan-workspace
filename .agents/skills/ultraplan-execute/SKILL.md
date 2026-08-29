---
name: ultraplan-execute
description: Manually run the UltraPlan execute stage when given a project sprint path or project/sprint references. Use only when the user explicitly invokes $ultraplan-execute or directly asks to run this exact UltraPlan stage; do not invoke implicitly.
---

# UltraPlan Execute

Run this stage interactively while preserving UltraPlan's governed artifact chain.

## Operating contract

1. Treat a supplied sprint path as UltraPlan stage input, not as a Git target. For an input such as `projects/<project>/sprints/<sprint>/` or `.ultra/projects/<project>/sprints/<sprint>/`, find the workspace root, derive `<project>` and `<sprint>` from the path, and read the matching `project-index.md`. The sprint directory contains governed stage artifacts. Read `<sprint>/.workspace.json` and use its absolute `path` as the implementation target. Verify that the record names the expected source repository and branch, that `path` exists, that Git reports it as a worktree of the recorded `sourceRoot` on the recorded `branch`, and that the worktree is clean. Confirm from the code and Git history that the prior sprint's implementation is present in substance; do not require an exact prior commit, branch, worktree, artifact fingerprint, or evidence identity. Run every implementation edit, source inspection, test, formatter, build, and Git command in that worktree. Never implement in `Target Implementation Directory`, `Repository`, `sourceRoot`, the workspace root, or the checkout from which UltraPlan was launched. If the workspace record is absent, malformed, stale, dirty, missing the prior implementation, or otherwise fails verification, stop and report the exact problem. Do not guess a worktree or fall back to another checkout. This worktree rule overrides target-directory wording in the resolved or canonical execution prompt. Do not search nested source repositories for a similarly named skill, and do not ask what target to use merely because the supplied input is a directory.
2. If no sprint path was supplied, locate the workspace root and resolve the project and sprint from explicit references and the current location. Ask only when the project index is missing, a required implementation target cannot be resolved, or more than one project/sprint remains possible.
3. Run all UltraPlan commands from the resolved workspace root. Run `ultraplan project <project> status` and `ultraplan sprint <project> <sprint> status --json`. Treat files, the project index, and fresh CLI status as authoritative; never hand-edit flow-state JSON.
4. Check these prerequisites:

- validated plan and all planning artifacts
- resolvable target implementation directory

5. Inspect the selected sprint's planning artifacts directly, but use only the selected sprint's recorded worktree as the execution admission gate. The worktree must exist, validate as the recorded Git worktree, be clean, and contain the prior sprint's implementation in substance. The prior implementation need not come from an exact commit, branch, worktree, artifact fingerprint, or evidence identity. Prior-sprint review, smoke, QA, verification, dogfood, worktree, and promotion state are context only and must not block source changes for the selected sprint. If the selected worktree fails an admission requirement, show the exact gaps and ask whether to fill them. Do not fill prerequisite gaps until the user agrees.
6. If the target is already complete and valid, summarize that state and ask before regenerating or materially changing it.
7. If the user explicitly asks for a proposal, analysis, or discussion only, inspect all relevant evidence and return that without writing artifacts or advancing state. Otherwise, do the stage now; do not stop at a proposal.
8. Act as the execution agent and perform the entire stage manually with your own file-editing and command tools. Use the UltraPlan CLI only for `project <project> status`, `sprint <project> <sprint> status --json`, `sprint <project> <sprint> prompt execute`, and the final `sprint <project> <sprint> review --dry-run --json` readiness check. Do not use any other UltraPlan CLI command during execution, including execution dry-run, validate, execute, verify, smoke, or flow commands.
9. Use the current effective prompt and concrete paths:

    ultraplan sprint <project> <sprint> prompt execute

   The resolved prompt can include workspace or project overrides and therefore takes precedence over the canonical prompt below.
10. Complete the stage-specific workflow, preserve unrelated user edits, and do not cross declared mutation boundaries.
11. Inspect the completed implementation and governed execution artifacts yourself and run the plan's required checks directly. Do not use an UltraPlan CLI validation command; use sprint status only to reconcile the state after the manual work is recorded.
12. Run `ultraplan sprint <project> <sprint> status --json` after writes or governed review execution so flow-state, freshness, and artifact status are reconciled. Re-run project status if the project or sprint index changed.
13. Inspect downstream artifacts for references made stale by this change. Update directly coupled indexes/references when safe; otherwise report the exact dependent stage that must be revisited. Never delete or silently rewrite downstream decisions.
14. Finish with the artifact/result paths, validation outcome, state transition, and any remaining blocker or dependent stage.

## Stage workflow

Use the resolved execution prompt and approved plan to perform the implementation yourself in the target implementation directory.

Act as the execution agent: work through incomplete plan tasks in order, inspect the code, edit the implementation, run the required checks directly, and maintain plan checkboxes, execution evidence, and the execution artifacts required by the prompt. Continue until the plan is complete or a genuine blocker requires the user.

For this stage, use the UltraPlan CLI only to check project or sprint state, materialise the effective execution prompt, and run the review readiness check described below. Do not use CLI execution dry runs, validation, execute, verify, smoke, or flow commands to perform, preview, validate, or complete the sprint. Inspect and verify the implementation and governed artifacts yourself, then use sprint status to reconcile the resulting state.

After the implementation and execution evidence are complete, run `ultraplan sprint <project> <sprint> review --dry-run --json`. Require the top-level status and `result.execution_status` to be `ready` with no blocking diagnostics. If review is not ready, inspect its diagnostics, fix every in-scope execution artifact or evidence problem directly, reconcile with sprint status, and repeat the review dry-run until it is ready or a genuine blocker must be reported. Do not launch the actual review from this skill.

## Canonical stage prompt

The current resolved CLI prompt wins over this embedded baseline.

# Sprint Execution - Plan-Driven Implementation

> **Inputs:** plan.md, reasoning.md, sprint-index.md, technical-handbook.md, docs/*.md, roadmap.md

---

Use this prompt to execute one approved sprint plan.

The goal is implementation, verification, and evidence capture. Do not redesign unless the plan is blocked, contradicted by the codebase, or missing a decision required for safe implementation.

## Required Inputs

Load these files first:

1. Sprint plan: `.ultra/projects/{project}/sprints/{sprint-slug}/plan.md`
2. Sprint reasoning: `.ultra/projects/{project}/sprints/{sprint-slug}/reasoning.md`
3. Sprint index: `.ultra/projects/{project}/sprints/{sprint-slug}/sprint-index.md`
4. Technical handbook: `.ultra/projects/{project}/sprints/{sprint-slug}/technical-handbook.md`
5. Project docs: all markdown files in `.ultra/projects/{project}/docs/*.md`
6. Project roadmap: `.ultra/projects/{project}/roadmap.md`
7. Target implementation directory: read from `.ultra/projects/{project}/project-index.md` under `Project Scope` as `Target Implementation Directory`

The **target implementation directory** is a per-project setting stored in `.ultra/projects/{project}/project-index.md`. Write all implementation source files (`.go`, `.mod`, `.sum`, test files, config, etc.) into that target implementation directory, NOT into the sprint directory. The sprint directory (`.ultra/projects/{project}/sprints/{sprint-slug}/`) is for sprint artifacts only (sprint-index.md, reasoning.md, plan.md, review.md).

If `Target Implementation Directory` is missing, empty, or ambiguous, pause before editing implementation files and update the project index with the correct per-project directory.

## Execution Rules

- Follow the sprint plan before inventing new work.
- Use sprint reasoning to understand why the plan made its design choices.
- Keep work inside the sprint scope.
- Do not pull non-scope or later-sprint behavior into the implementation.
- If the plan is wrong, incomplete, or unsafe, pause and update the plan with the reason.
- Prefer small, reviewable edits that match the existing codebase.
- Preserve product-specific boundaries from the project docs in `docs/` and the roadmap.
- Do not silently skip checklist items, tests, risks, or quality gates.
- Record explicit deferrals with reason, impact, and follow-up.

## Implementation Workflow

1. Read the sprint plan and identify the first incomplete task.
2. Inspect the existing codebase in the project's target implementation directory before editing.
3. Implement one coherent task or sub-task at a time.
4. Run the smallest useful verification after each meaningful change.
5. Update the sprint plan checklist and `Execution Evidence` as work progresses.
6. Repeat until the sprint is complete or blocked.
7. Run the sprint's full verification set before marking it complete.

## Writing Implementation Files

Write all source files (Go, YAML, etc.) to the project's **target implementation directory**.

Example paths when the project's target implementation directory is `/path/to/go-todo`:

- `/path/to/go-todo/go.mod`
- `/path/to/go-todo/cmd/todo/main.go`
- `/path/to/go-todo/internal/model/task.go`
- `/path/to/go-todo/internal/store/jsonstore.go`
- `/path/to/go-todo/internal/config/paths.go`
- `/path/to/go-todo/internal/store/jsonstore_test.go`

## Marking Implementation Complete

When all tasks are done (or blocked), write a `.run-state.json` file to the sprint directory:

**File:** `.ultra/projects/{project}/sprints/{sprint-slug}/.run-state.json`

**Content:**

```json
{
  "status": "complete",
  "completedAt": "<ISO-8601 timestamp>",
  "files": ["<list of files created/changed>"],
  "testsRun": ["<list of tests/checks run and their results>"],
  "blockers": ["<list of blockers if any>"]
}
```

This file is the implementation artifact that proves the sprint was executed.

## Verification Rules

Run the tests and checks named in the sprint plan. If a named check cannot run:

- record the command or check that could not run
- record why it could not run
- record the residual risk
- add a follow-up or blocker if the risk is material

## Updating The Sprint Plan

Keep `plan.md` current during execution. Update:

- `Status`
- task and sub-task checkboxes
- testing and documentation checklist
- risks and blockers
- execution evidence

## Completion Rules

Before marking the sprint complete:

1. All in-scope checklist items are complete or explicitly deferred.
2. Required tests and checks have passed or have justified deferrals.
3. Risks are closed, mitigated, or carried forward.
4. `Execution Evidence` includes important commands, decisions, and deferrals.
5. `Review And Sign-Off` has accurate sprint status and completion date.
6. `.run-state.json` written to sprint directory as implementation artifact.

## Final Response

When execution is done or blocked, report concisely:

- implementation summary
- files changed
- tests/checks run
- remaining blockers, risks, or deferrals

The sprint plan is the durable record — keep it current.
