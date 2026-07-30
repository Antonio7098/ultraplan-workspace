---
name: ultraplan-plan
description: Manually run the UltraPlan plan stage for a selected project sprint. Use only when the user explicitly invokes $ultraplan-plan or directly asks to run this exact UltraPlan stage; do not invoke implicitly.
---

# UltraPlan Sprint Plan

Run this stage interactively while preserving UltraPlan's governed artifact chain.

## Operating contract

1. Locate the workspace root and resolve the project and sprint. If either is ambiguous, ask the user instead of guessing.
2. Run `ultraplan project <project> status` and `ultraplan sprint <project> <sprint> status --json`. Treat files and fresh CLI status as authoritative; never hand-edit flow-state JSON.
3. Check these prerequisites:

- requirements
- sprint-index
- technical-handbook
- reasoning

4. Validate every prerequisite that has a sprint validation command. If anything is missing, invalid, stale, or internally inconsistent, show the exact gaps and ask whether to fill them. Do not fill prerequisite gaps until the user agrees. If they agree, run the corresponding earlier UltraPlan skills in canonical order, then return to this stage.
5. If the target is already complete and valid, summarize that state and ask before regenerating or materially changing it.
6. If the user explicitly asks for a proposal, analysis, or discussion only, inspect all relevant evidence and return that without writing artifacts or advancing state. Otherwise, do the stage now; do not stop at a proposal.
7. Use the current effective prompt and concrete paths:

    ultraplan sprint <project> <sprint> prompt plan

   The resolved prompt can include workspace or project overrides and therefore takes precedence over the canonical prompt below.
8. Complete the stage-specific workflow, preserve unrelated user edits, and do not cross declared mutation boundaries.
9. Run `ultraplan sprint <project> <sprint> validate plan` when supported. Fix validation findings within this stage rather than declaring success early.
10. Run `ultraplan sprint <project> <sprint> status --json` after writes or governed execution so flow-state, freshness, and artifact status are reconciled. Re-run project status if the project or sprint index changed.
11. Inspect downstream artifacts for references made stale by this change. Update directly coupled indexes/references when safe; otherwise report the exact dependent stage that must be revisited. Never delete or silently rewrite downstream decisions.
12. Finish with the artifact/result paths, validation outcome, state transition, and any remaining blocker or dependent stage.

## Stage workflow

Create or revise plan.md from the resolved prompt.
Carry decisions forward rather than reopening them. Make tasks ordered, bounded, testable, traceable to requirements, and explicit about files, checks, evidence, dependencies, and stop conditions. Do not implement the plan in this stage.

## Canonical stage prompt

The current resolved CLI prompt wins over this embedded baseline.

# Sprint Planning - Evidence-Grounded Implementation Plan

> **Inputs:** project-index.md, requirements.md, docs/*.md, sprint-index.md, technical-handbook.md, area reasoning (if any), reasoning.md, injected sprint-plan template

---

Use this prompt to plan one implementation sprint for an Ultra project.

The outputs are a sprint reasoning document and then a sprint plan. Do not implement code. Do not make ungrounded architecture decisions. Every important decision must trace to project requirements, roadmap scope, evolved study evidence, or an explicitly named open question.

## Required Inputs

Load these files first:

1. Project index: `.ultra/projects/{project}/project-index.md`
2. Sprint requirements: `.ultra/projects/{project}/sprints/{sprint-slug}/requirements.md`
3. Project docs: all markdown files in `.ultra/projects/{project}/docs/*.md`
5. Sprint index: `.ultra/projects/{project}/sprints/{sprint-slug}/sprint-index.md`
6. Technical handbook: `.ultra/projects/{project}/sprints/{sprint-slug}/technical-handbook.md`
7. Area reasoning: `.ultra/projects/{project}/sprints/{sprint-slug}/reasoning/*.md` (if present)
8. Sprint reasoning: `.ultra/projects/{project}/sprints/{sprint-slug}/reasoning.md`
9. Injected sprint-plan template section in this prompt

## Evidence Loading Order

Use this order so the plan stays grounded:

1. Read `requirements.md` first.
2. Read the sprint reasoning to understand what was decided and why.
3. Read the sprint index to understand selected context.
4. Read the technical handbook for evidence.
5. Read area reasoning for area-specific decisions.
6. Read the PRD and TRD sections relevant to this sprint.
7. Read the project index for traceability back to selected studies/contracts if needed.
8. Open linked final reports only when a specific decision needs deeper evidence.
9. Resolve code references only for specific implementation questions.
10. At the top of `plan.md`, write an `> **Inputs Used:**` line that lists the exact files used for this document.

If context is too large:

1. Keep sprint reasoning, PRD, TRD, and roadmap sprint section in context.
2. Load evidence packs tied to current decisions.
3. Load final reports only when evidence packs are insufficient.
4. Record omitted evidence and why it was omitted.

## Planning Rules

- Start from fundamentals. Do not pull later-sprint complexity into earlier sprints.
- Respect the roadmap scope. If evidence suggests scope change, record as recommendation.
- Separate requirements from design decisions.
- Use sprint reasoning to justify design choices.
- Prefer small, testable increments over broad abstractions.
- Record tradeoffs, rejected alternatives, and anti-patterns.
- Do not implement code while planning.

## Decision Discipline

For every major decision, capture:

- Decision made
- Requirement it satisfies
- Evidence used
- Tradeoff accepted
- Alternative rejected
- Risk or follow-up

If evidence is insufficient, write an open question instead of guessing.

## Sprint Reasoning Output

Write the sprint reasoning to:
`.ultra/projects/{project}/sprints/{sprint-slug}/reasoning.md`

Use the injected sprint-reasoning template section when this prompt is used to create or repair sprint reasoning. Fill every section.

## Sprint Plan Output

Write the sprint plan to:
`.ultra/projects/{project}/sprints/{sprint-slug}/plan.md`

Use the injected sprint-plan template section. Fill every section. The plan must cite `reasoning.md` and carry forward its decisions, expected evidence, risks, assumptions, and open questions.

## Quality Bar

A good sprint plan is specific enough that implementation can proceed without rereading every study report, but evidence-grounded enough that decisions can be audited later.

Avoid generic phrases like "clean architecture" unless named with evidence basis.
