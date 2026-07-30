---
name: ultraplan-area-reasoning
description: Manually run the UltraPlan area-reasoning stage for a selected project sprint. Use only when the user explicitly invokes $ultraplan-area-reasoning or directly asks to run this exact UltraPlan stage; do not invoke implicitly.
---

# UltraPlan Area Reasoning

Run this stage interactively while preserving UltraPlan's governed artifact chain.

## Operating contract

1. Locate the workspace root and resolve the project and sprint. If either is ambiguous, ask the user instead of guessing.
2. Run `ultraplan project <project> status` and `ultraplan sprint <project> <sprint> status --json`. Treat files and fresh CLI status as authoritative; never hand-edit flow-state JSON.
3. Check these prerequisites:

- requirements
- sprint-index
- technical-handbook

4. Validate every prerequisite that has a sprint validation command. If anything is missing, invalid, stale, or internally inconsistent, show the exact gaps and ask whether to fill them. Do not fill prerequisite gaps until the user agrees. If they agree, run the corresponding earlier UltraPlan skills in canonical order, then return to this stage.
5. If the target is already complete and valid, summarize that state and ask before regenerating or materially changing it.
6. If the user explicitly asks for a proposal, analysis, or discussion only, inspect all relevant evidence and return that without writing artifacts or advancing state. Otherwise, do the stage now; do not stop at a proposal.
7. Use the current effective prompt and concrete paths:

    ultraplan sprint <project> <sprint> prompt area-reasoning

   The resolved prompt can include workspace or project overrides and therefore takes precedence over the canonical prompt below.
8. Complete the stage-specific workflow, preserve unrelated user edits, and do not cross declared mutation boundaries.
9. Run `ultraplan sprint <project> <sprint> validate area-reasoning` when supported. Fix validation findings within this stage rather than declaring success early.
10. Run `ultraplan sprint <project> <sprint> status --json` after writes or governed execution so flow-state, freshness, and artifact status are reconciled. Re-run project status if the project or sprint index changed.
11. Inspect downstream artifacts for references made stale by this change. Update directly coupled indexes/references when safe; otherwise report the exact dependent stage that must be revisited. Never delete or silently rewrite downstream decisions.
12. Finish with the artifact/result paths, validation outcome, state transition, and any remaining blocker or dependent stage.

## Stage workflow

Create or revise every area-reasoning document selected by the sprint index, and no others.
When the user requests a deep dive, treat the work as an interactive design discussion: surface design pressures, alternatives, trade-offs, risks, and evidence; resolve one meaningful decision at a time; then record the conclusions unless the user asked for discussion or a proposal only.

## Canonical stage prompt

The current resolved CLI prompt wins over this embedded baseline.

# Create Area Reasoning

> **Inputs:** technical-handbook.md, requirements.md, docs/*.md, injected selected reasoning template

---

Use this prompt to create optional area-specific reasoning documents for a sprint.

## Required Inputs

Read these files first:

1. Technical handbook: `.ultra/projects/{project}/sprints/{sprint-slug}/technical-handbook.md`
2. Sprint requirements: `.ultra/projects/{project}/sprints/{sprint-slug}/requirements.md`
3. Project docs: all markdown files in `.ultra/projects/{project}/docs/*.md`
4. Injected selected reasoning template section in this prompt

## Output

For each area selected by sprint-index, write to:
`.ultra/projects/{project}/sprints/{sprint-slug}/reasoning/<area>.md`

Only create files for areas explicitly selected in sprint-index. Do NOT create area reasoning documents for ceremony.

## Instructions

1. The flow has already determined which area files to create. Do not infer extra areas.
2. Read the technical handbook for evidence context.
3. Read `requirements.md` and all project docs in `docs/` for sprint-specific scope and constraints.
4. For each selected area:
   - Use the injected selected reasoning template section as source material to reason through, not as the literal output structure
   - At the very top of the file, add an `> **Inputs Used:**` line listing the exact files used for that document
   - Include these exact required `##` sections with concrete content:
     - `## Area Decisions`
     - `## Trade-Offs`
     - `## Evidence`
     - `## Risks`
   - Ground decisions in technical handbook evidence
   - Record the key conclusion and evidence basis
   - Note any open questions or risks
5. Do not duplicate content from technical-handbook — synthesize into area-specific conclusions.
6. Ensure each area reasoning document is self-contained and can be understood without reading other area documents.

## Skip Criteria

Skip creating area reasoning if:

- No areas are selected in sprint-index
- Area reasoning files already exist and are complete
- Contains placeholders

## Quality Bar

Each area reasoning document must:

- Have a clear area name and scope
- Cite technical handbook evidence
- Make final area-specific decisions (no more "TBD")
- Record rejected alternatives
- Note risks, assumptions, and open questions
- Be referenced in sprint-reasoning.md
