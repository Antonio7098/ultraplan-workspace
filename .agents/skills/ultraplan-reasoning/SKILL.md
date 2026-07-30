---
name: ultraplan-reasoning
description: Manually run the UltraPlan reasoning stage for a selected project sprint. Use only when the user explicitly invokes $ultraplan-reasoning or directly asks to run this exact UltraPlan stage; do not invoke implicitly.
---

# UltraPlan Sprint Reasoning

Run this stage interactively while preserving UltraPlan's governed artifact chain.

## Operating contract

1. Locate the workspace root and resolve the project and sprint. If either is ambiguous, ask the user instead of guessing.
2. Run `ultraplan project <project> status` and `ultraplan sprint <project> <sprint> status --json`. Treat files and fresh CLI status as authoritative; never hand-edit flow-state JSON.
3. Check these prerequisites:

- requirements
- sprint-index
- technical-handbook
- selected area-reasoning documents or an explicit none selection

4. Validate every prerequisite that has a sprint validation command. If anything is missing, invalid, stale, or internally inconsistent, show the exact gaps and ask whether to fill them. Do not fill prerequisite gaps until the user agrees. If they agree, run the corresponding earlier UltraPlan skills in canonical order, then return to this stage.
5. If the target is already complete and valid, summarize that state and ask before regenerating or materially changing it.
6. If the user explicitly asks for a proposal, analysis, or discussion only, inspect all relevant evidence and return that without writing artifacts or advancing state. Otherwise, do the stage now; do not stop at a proposal.
7. Use the current effective prompt and concrete paths:

    ultraplan sprint <project> <sprint> prompt reasoning

   The resolved prompt can include workspace or project overrides and therefore takes precedence over the canonical prompt below.
8. Complete the stage-specific workflow, preserve unrelated user edits, and do not cross declared mutation boundaries.
9. Run `ultraplan sprint <project> <sprint> validate reasoning` when supported. Fix validation findings within this stage rather than declaring success early.
10. Run `ultraplan sprint <project> <sprint> status --json` after writes or governed execution so flow-state, freshness, and artifact status are reconciled. Re-run project status if the project or sprint index changed.
11. Inspect downstream artifacts for references made stale by this change. Update directly coupled indexes/references when safe; otherwise report the exact dependent stage that must be revisited. Never delete or silently rewrite downstream decisions.
12. Finish with the artifact/result paths, validation outcome, state transition, and any remaining blocker or dependent stage.

## Stage workflow

Create or revise the final sprint reasoning document from the resolved prompt.
When the user requests a deep dive, discuss design pressures, competing approaches, accepted trade-offs, technical debt, future consequences, risks, and evidence before committing the conclusions to the artifact. Do not collapse a requested discussion into a shallow one-shot answer.

## Canonical stage prompt

The current resolved CLI prompt wins over this embedded baseline.

# Create Sprint Reasoning

> **Inputs:** project-index.md, requirements.md, docs/*.md, sprint-index.md, technical-handbook.md, area reasoning (if any), injected sprint-reasoning template

---

Use this prompt to create the sprint reasoning document for a sprint.

## Required Inputs

Read these files first:

1. Project index: `.ultra/projects/{project}/project-index.md`
2. Sprint requirements: `.ultra/projects/{project}/sprints/{sprint-slug}/requirements.md`
3. Project docs: all markdown files in `.ultra/projects/{project}/docs/*.md`
4. Sprint index: `.ultra/projects/{project}/sprints/{sprint-slug}/sprint-index.md`
5. Technical handbook: `.ultra/projects/{project}/sprints/{sprint-slug}/technical-handbook.md`
6. Area reasoning files: `.ultra/projects/{project}/sprints/{sprint-slug}/reasoning/*.md` (if any)
7. Injected sprint-reasoning template section in this prompt

## Output

Write to: `.ultra/projects/{project}/sprints/{sprint-slug}/reasoning.md`

Use the injected sprint-reasoning template section. Fill every section. This document makes the final sprint decisions — it synthesizes selected context, handbook evidence, area-specific reasoning, and contracts into final sprint decisions. It does NOT execute implementation.

## Instructions

1. Read `requirements.md` first.
2. Read the project index and all project docs in `docs/` for full supporting context.
3. Read the sprint index to understand selected context.
4. Read the technical handbook for study evidence.
5. Read area reasoning documents if present.
6. At the top of `reasoning.md`, write an `> **Inputs Used:**` line that lists the exact files used for this document.
7. If area reasoning files exist, summarize their actual conclusions. Do not say their output paths are empty, unwritten, or pending when files are present.
8. Synthesize all evidence into final sprint decisions.
9. For each decision:
   - State what will be done and why
   - Record the study/source grounding: cite `technical-handbook.md`, selected evidence reports, and concrete studied repos/source references where relevant
   - If no study sources are relevant to a decision, explicitly say so and justify why
   - Record trade-offs accepted, potential technical debt, and future impacts
   - Name rejected alternatives and why
   - Map to applicable contracts/requirement IDs
   - Define expected evidence (tests, logs, review checks)
10. Add a "Repos Studied / Source Evidence Used" section using the technical handbook and selected reports. Include why each repo/report mattered, and which decisions it influenced.
11. Add deeper analysis of accepted trade-offs, possible technical debt, and future considerations before final decisions.
12. Record assumptions and risks.
13. Define implementation constraints.
14. The plan must be able to execute these decisions without reopening architecture.

## Skip Criteria

Skip creating sprint-reasoning if:

- `reasoning.md` already exists and is complete
- Sprint index or technical handbook is missing
- Contains placeholders
- No final decisions are recorded

## Quality Bar

The sprint reasoning must:

- Reference sprint-index selected context
- Reference technical-handbook evidence
- Cite concrete study sources/repositories or explain why no sources were relevant
- Analyze accepted trade-offs, potential technical debt, and future considerations
- Make final, specific decisions (no "TBD" on core architecture)
- Record at least 2 rejected alternatives with rationale
- Map decisions to contracts and requirement IDs
- Define specific, reviewable expected evidence
- Record risks with mitigation or follow-up
- Be actionable by plan.md without reopening architecture
