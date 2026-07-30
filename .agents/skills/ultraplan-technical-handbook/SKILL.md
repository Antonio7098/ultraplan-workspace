---
name: ultraplan-technical-handbook
description: Manually run the UltraPlan technical-handbook stage for a selected project sprint. Use only when the user explicitly invokes $ultraplan-technical-handbook or directly asks to run this exact UltraPlan stage; do not invoke implicitly.
---

# UltraPlan Technical Handbook

Run this stage interactively while preserving UltraPlan's governed artifact chain.

## Operating contract

1. Locate the workspace root and resolve the project and sprint. If either is ambiguous, ask the user instead of guessing.
2. Run `ultraplan project <project> status` and `ultraplan sprint <project> <sprint> status --json`. Treat files and fresh CLI status as authoritative; never hand-edit flow-state JSON.
3. Check these prerequisites:

- requirements
- sprint-index

4. Validate every prerequisite that has a sprint validation command. If anything is missing, invalid, stale, or internally inconsistent, show the exact gaps and ask whether to fill them. Do not fill prerequisite gaps until the user agrees. If they agree, run the corresponding earlier UltraPlan skills in canonical order, then return to this stage.
5. If the target is already complete and valid, summarize that state and ask before regenerating or materially changing it.
6. If the user explicitly asks for a proposal, analysis, or discussion only, inspect all relevant evidence and return that without writing artifacts or advancing state. Otherwise, do the stage now; do not stop at a proposal.
7. Use the current effective prompt and concrete paths:

    ultraplan sprint <project> <sprint> prompt technical-handbook

   The resolved prompt can include workspace or project overrides and therefore takes precedence over the canonical prompt below.
8. Complete the stage-specific workflow, preserve unrelated user edits, and do not cross declared mutation boundaries.
9. Run `ultraplan sprint <project> <sprint> validate technical-handbook` when supported. Fix validation findings within this stage rather than declaring success early.
10. Run `ultraplan sprint <project> <sprint> status --json` after writes or governed execution so flow-state, freshness, and artifact status are reconciled. Re-run project status if the project or sprint index changed.
11. Inspect downstream artifacts for references made stale by this change. Update directly coupled indexes/references when safe; otherwise report the exact dependent stage that must be revisited. Never delete or silently rewrite downstream decisions.
12. Finish with the artifact/result paths, validation outcome, state transition, and any remaining blocker or dependent stage.

## Stage workflow

Create or revise the technical handbook from only the evidence selected by the sprint index.
Distil patterns, trade-offs, cautions, and open questions. Preserve the boundary between evidence and final design decisions.

## Canonical stage prompt

The current resolved CLI prompt wins over this embedded baseline.

# Create Technical Handbook

> **Inputs:** `sprint-index.md`, `project-index.md`, `requirements.md`, `docs/*.md`, reports selected in `sprint-index.md` -> "Selected Evidence Reports" section, injected technical-handbook template

---

Use this prompt to create the technical handbook for a sprint.

## Required Inputs

Read these files first:

1. Sprint index: `.ultra/projects/{project}/sprints/{sprint-slug}/sprint-index.md` — contains the **Selected Evidence Reports** section that tells you which reports to read
2. Injected technical-handbook template section in this prompt

Then read the evidence reports listed in the sprint index's "Selected Evidence Reports" section. The project index's "Available Evidence Reports" table is the authoritative source for those report paths.

Do NOT read all evidence reports indiscriminately. Only read the ones the sprint index selects.

## Output

Write to: `.ultra/projects/{project}/sprints/{sprint-slug}/technical-handbook.md`

Use the injected technical-handbook template section. Fill every section. The handbook distills selected studies and reports into sprint-relevant patterns, trade-offs, cautions, and questions. It does NOT decide implementation.

## Instructions

1. Read the sprint index first to find which evidence reports to read.
2. Read those evidence reports directly. Evidence reports cite specific source files and line numbers from real study repos (e.g., `mitchellh-cli/cli.go:70`, `go-task/cmd/task/task.go:23-46`). These are your evidence base.
3. Read `requirements.md` for sprint scope context.
4. Extract from the evidence reports:
   - Relevant patterns that apply to this sprint's scope
   - Trade-offs observed across sources
   - Warnings or anti-patterns to avoid
   - Open questions that sprint reasoning must resolve
   - Design pressures visible in the study evidence
5. At the top of `technical-handbook.md`, write `> **Inputs Used:**` listing the exact files used.
6. Do not use contracts, PRD, TRD, or other governance documents as handbook evidence. They may provide scope context only.
7. Do not hardcode patterns — ground them in the evidence from the reports.
8. Do not make final implementation decisions — defer those to sprint reasoning.
9. Cite specific source files and line numbers from the reports (e.g., `mitchellh-cli/cli.go:70`).

## Skip Criteria

Skip creating technical-handbook if:

- `technical-handbook.md` already exists and is complete
- Sprint index is missing or has no "Selected Evidence Reports" section
- Selected evidence reports cannot be located

## Quality Bar

The technical handbook must:

- Cite specific evidence reports and their findings
- Identify at least 3 relevant patterns with evidence basis (file paths + line numbers from evidence reports)
- Document at least 2 trade-offs with benefit/cost analysis
- Include at least 2 warnings or anti-patterns
- List open questions for sprint reasoning
- Point to specific evidence (file paths, sections) to inspect
