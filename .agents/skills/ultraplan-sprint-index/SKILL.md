---
name: ultraplan-sprint-index
description: Manually run the UltraPlan sprint-index stage for a selected project sprint. Use only when the user explicitly invokes $ultraplan-sprint-index or directly asks to run this exact UltraPlan stage; do not invoke implicitly.
---

# UltraPlan Sprint Index

Run this stage interactively while preserving UltraPlan's governed artifact chain.

## Operating contract

1. Locate the workspace root and resolve the project and sprint. If either is ambiguous, ask the user instead of guessing.
2. Run `ultraplan project <project> status` and `ultraplan sprint <project> <sprint> status --json`. Treat files and fresh CLI status as authoritative; never hand-edit flow-state JSON.
3. Check these prerequisites:

- requirements

4. Validate every prerequisite that has a sprint validation command. If anything is missing, invalid, stale, or internally inconsistent, show the exact gaps and ask whether to fill them. Do not fill prerequisite gaps until the user agrees. If they agree, run the corresponding earlier UltraPlan skills in canonical order, then return to this stage.
5. If the target is already complete and valid, summarize that state and ask before regenerating or materially changing it.
6. If the user explicitly asks for a proposal, analysis, or discussion only, inspect all relevant evidence and return that without writing artifacts or advancing state. Otherwise, do the stage now; do not stop at a proposal.
7. Use the current effective prompt and concrete paths:

    ultraplan sprint <project> <sprint> prompt sprint-index

   The resolved prompt can include workspace or project overrides and therefore takes precedence over the canonical prompt below.
8. Complete the stage-specific workflow, preserve unrelated user edits, and do not cross declared mutation boundaries.
9. Run `ultraplan sprint <project> <sprint> validate sprint-index` when supported. Fix validation findings within this stage rather than declaring success early.
10. Run `ultraplan sprint <project> <sprint> status --json` after writes or governed execution so flow-state, freshness, and artifact status are reconciled. Re-run project status if the project or sprint index changed.
11. Inspect downstream artifacts for references made stale by this change. Update directly coupled indexes/references when safe; otherwise report the exact dependent stage that must be revisited. Never delete or silently rewrite downstream decisions.
12. Finish with the artifact/result paths, validation outcome, state transition, and any remaining blocker or dependent stage.

## Stage workflow

Create or revise the sprint index from the resolved prompt.
Keep it a selection document: update selected contracts, evidence, reasoning templates, carry-forward decisions, and exclusions without making implementation decisions.

## Canonical stage prompt

The current resolved CLI prompt wins over this embedded baseline.

# Create Sprint Index

> **Inputs:** `project-index.md`, `requirements.md`, `docs/*.md`, injected sprint-index template

---

Use this prompt to create the sprint index for a sprint.

## Required Inputs

Read these files first:

1. Project index: `.ultra/projects/{project}/project-index.md`
2. **Sprint requirements**: `.ultra/projects/{project}/sprints/{sprint-slug}/requirements.md`
3. Project docs: all markdown files in `.ultra/projects/{project}/docs/*.md`
4. Injected sprint-index template section in this prompt

## Output

Write to: `.ultra/projects/{project}/sprints/{sprint-slug}/sprint-index.md`

Use the injected sprint-index template section. Fill every section. The sprint index selects what must be read, distilled, reasoned through, or checked for this sprint. It does NOT make implementation decisions.

## Instructions

1. Read `requirements.md` first — it defines the sprint goal, required outputs, acceptance criteria, non-goals, and constraints.
2. Read the project index to understand the available pool of contracts, studies, evidence reports, reasoning templates, and protocols.
3. Read the PRD and TRD for supporting product and technical context.
4. Use the injected sprint-index template section. Fill every section.
5. At the top of `sprint-index.md`, write `> **Inputs Used:**` listing the exact files used.
6. **Contracts**: list each selected contract by simple name (e.g. "Architecture"). The contract applies as a flat whole — no requirement ID mappings or per-clause breakdowns. If a clause is particularly important, mention it in the Why Selected column.
7. **Selected Evidence Reports**: copy the relevant rows from the project index's "Available Evidence Reports" table. These tell the technical handbook which reports to read.
8. **Reasoning templates**: select which area reasoning templates apply and specify their output filenames.
9. **Excluded context**: record what is explicitly excluded and why. The table must include explicit rows covering these known deferred behaviors when they are out of scope: implementation execution, smoke investigation, review automation, issue tracking, and Git mutation.
10. Do not invent sections or add content outside the template.

## Skip Criteria

Skip creating sprint-index if:

- `sprint-index.md` already exists and is complete
- Required inputs (project-index, requirements.md) are missing
- Selected sections are empty or contain placeholders

## Quality Bar

The sprint index must:

- Name the sprint goal and planned output clearly
- List selected contracts by simple name only (flat, not mapped to IDs)
- Copy the relevant evidence report rows from the project index
- Record carry-forward decisions from prior sprints
- Explicitly exclude non-goals and non-relevant context
- Include explicit Excluded Context rows for implementation execution, smoke investigation, review automation, issue tracking, and Git mutation unless the sprint requirements explicitly bring one of those behaviors into scope
- Reference only items that appear in the project index (no new items invented)
