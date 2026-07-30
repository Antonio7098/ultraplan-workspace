---
name: ultraplan-execute
description: Manually run the UltraPlan execute stage for a selected project sprint. Use only when the user explicitly invokes $ultraplan-execute or directly asks to run this exact UltraPlan stage; do not invoke implicitly.
---

# UltraPlan Execute

Run this stage interactively while preserving UltraPlan's governed artifact chain.

## Operating contract

1. Locate the workspace root and resolve the project and sprint. If either is ambiguous, ask the user instead of guessing.
2. Run `ultraplan project <project> status` and `ultraplan sprint <project> <sprint> status --json`. Treat files and fresh CLI status as authoritative; never hand-edit flow-state JSON.
3. Check these prerequisites:

- validated plan and all planning artifacts
- resolvable target implementation directory

4. Validate every prerequisite that has a sprint validation command. If anything is missing, invalid, stale, or internally inconsistent, show the exact gaps and ask whether to fill them. Do not fill prerequisite gaps until the user agrees. If they agree, run the corresponding earlier UltraPlan skills in canonical order, then return to this stage.
5. If the target is already complete and valid, summarize that state and ask before regenerating or materially changing it.
6. If the user explicitly asks for a proposal, analysis, or discussion only, inspect all relevant evidence and return that without writing artifacts or advancing state. Otherwise, do the stage now; do not stop at a proposal.
7. Use the current effective prompt and concrete paths:

    ultraplan sprint <project> <sprint> prompt execute

   The resolved prompt can include workspace or project overrides and therefore takes precedence over the canonical prompt below.
8. Complete the stage-specific workflow, preserve unrelated user edits, and do not cross declared mutation boundaries.
9. Run `ultraplan sprint <project> <sprint> validate execute` when supported. Fix validation findings within this stage rather than declaring success early.
10. Run `ultraplan sprint <project> <sprint> status --json` after writes or governed execution so flow-state, freshness, and artifact status are reconciled. Re-run project status if the project or sprint index changed.
11. Inspect downstream artifacts for references made stale by this change. Update directly coupled indexes/references when safe; otherwise report the exact dependent stage that must be revisited. Never delete or silently rewrite downstream decisions.
12. Finish with the artifact/result paths, validation outcome, state transition, and any remaining blocker or dependent stage.

## Stage workflow

Run a dry-run first:

    ultraplan sprint <project> <sprint> execute --dry-run

Then execute or resume the governed plan:

    ultraplan sprint <project> <sprint> execute --resume

Follow progress, inspect partial failures, and continue until the execution state is complete or a genuine blocker requires the user. Do not replace this governed execution with untracked ad-hoc implementation.

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
