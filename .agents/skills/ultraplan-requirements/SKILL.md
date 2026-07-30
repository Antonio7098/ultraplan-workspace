---
name: ultraplan-requirements
description: Manually run the UltraPlan requirements stage for a selected project sprint. Use only when the user explicitly invokes $ultraplan-requirements or directly asks to run this exact UltraPlan stage; do not invoke implicitly.
---

# UltraPlan Requirements

Run this stage interactively while preserving UltraPlan's governed artifact chain.

## Operating contract

1. Locate the workspace root and resolve the project and sprint. If either is ambiguous, ask the user instead of guessing.
2. Run `ultraplan project <project> status` and `ultraplan sprint <project> <sprint> status --json`. Treat files and fresh CLI status as authoritative; never hand-edit flow-state JSON.
3. Check these prerequisites:

- project index
- roadmap and relevant project docs

4. Validate every prerequisite that has a sprint validation command. If anything is missing, invalid, stale, or internally inconsistent, show the exact gaps and ask whether to fill them. Do not fill prerequisite gaps until the user agrees. If they agree, run the corresponding earlier UltraPlan skills in canonical order, then return to this stage.
5. If the target is already complete and valid, summarize that state and ask before regenerating or materially changing it.
6. If the user explicitly asks for a proposal, analysis, or discussion only, inspect all relevant evidence and return that without writing artifacts or advancing state. Otherwise, do the stage now; do not stop at a proposal.
7. Use the current effective prompt and concrete paths:

    ultraplan sprint <project> <sprint> prompt requirements

   The resolved prompt can include workspace or project overrides and therefore takes precedence over the canonical prompt below.
8. Complete the stage-specific workflow, preserve unrelated user edits, and do not cross declared mutation boundaries.
9. Run `ultraplan sprint <project> <sprint> validate requirements` when supported. Fix validation findings within this stage rather than declaring success early.
10. Run `ultraplan sprint <project> <sprint> status --json` after writes or governed execution so flow-state, freshness, and artifact status are reconciled. Re-run project status if the project or sprint index changed.
11. Inspect downstream artifacts for references made stale by this change. Update directly coupled indexes/references when safe; otherwise report the exact dependent stage that must be revisited. Never delete or silently rewrite downstream decisions.
12. Finish with the artifact/result paths, validation outcome, state transition, and any remaining blocker or dependent stage.

## Stage workflow

Create or revise the exact requirements artifact from the resolved prompt.
If prior sprint reviews exist, carry forward only still-applicable decisions. Do not silently broaden the roadmap scope.

## Canonical stage prompt

The current resolved CLI prompt wins over this embedded baseline.

# Create Sprint Requirements

> **Inputs:** project-index.md, roadmap.md, docs/*.md, injected requirements template, prior sprint reviews

---

Use this prompt to create the sprint requirements document for a new sprint.

## Required Inputs

Read these files first:

1. Project index: `.ultra/projects/{project}/project-index.md`
2. Project roadmap: `.ultra/projects/{project}/roadmap.md`
3. Project docs: all markdown files in `.ultra/projects/{project}/docs/*.md`
4. Injected requirements template section in this prompt
5. Prior sprint decisions (if any): `.ultra/projects/{project}/sprints/*/review.md`

## Output

Write to: `.ultra/projects/{project}/sprints/{sprint-slug}/requirements.md`

## Instructions

1. Read the project index and roadmap to understand the overall project scope and sprint sequence.
2. Read all project docs in `docs/` for requirements context.
3. If this is not the first sprint, read prior sprint review(s) for carry-forward decisions.
4. Fill the injected requirements template section with:
   - **Sprint Goal**: one clear sentence describing what must be achieved.
   - **Required Outputs**: enumerate every deliverable file with its path and a one-line description.
   - **Acceptance Criteria**: checkable criteria that prove the sprint succeeded.
   - **Non-Goals**: explicitly what is NOT included.
   - **Constraints**: architectural, technical, or process constraints.
   - **Dependencies**: prior sprints or external inputs this sprint needs.
   - **Review Expectations**: what will be checked at review and how.

## Quality Bar

- The sprint goal must be achievable in one sprint.
- All required outputs must be specific (file paths, not vague descriptions).
- Acceptance criteria must be objectively verifiable.
- Non-goals must be specific enough to prevent scope creep.
- Constraints must be binding (not aspirational).

## Skip Criteria

Skip creating requirements if:

- `requirements.md` already exists and is complete
- Required inputs (roadmap, project-index) are missing
