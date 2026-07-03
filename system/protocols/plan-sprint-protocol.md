# Plan Sprint Protocol

> Role: sprint planning protocol.
> Used by: agents preparing `.ultra/projects/<project>/sprints/<sprint>/plan.md`.
> Sources: `.ultra/README.md` and `.ultra/prompts/*.md`.

This protocol instructs the agent to work from sprint requirements and project index through the Ultra sprint planning artifacts, ending with an evidence-grounded sprint plan. It is a planning protocol only. Do not implement code while using it.

## Purpose

Use this protocol when asked to plan a sprint or produce the planning artifacts for one sprint.

The goal is to produce a sprint plan that can be executed without reopening architecture, while preserving traceability back to:

- sprint requirements
- project docs
- project index selections
- selected evidence reports
- technical handbook findings
- area reasoning, when selected
- sprint reasoning decisions

The final output is:

```text
.ultra/projects/<project>/sprints/<sprint-slug>/plan.md
```

Planning may also create or validate these prerequisite artifacts:

```text
.ultra/projects/<project>/sprints/<sprint-slug>/requirements.md
.ultra/projects/<project>/sprints/<sprint-slug>/sprint-index.md
.ultra/projects/<project>/sprints/<sprint-slug>/technical-handbook.md
.ultra/projects/<project>/sprints/<sprint-slug>/reasoning/*.md
.ultra/projects/<project>/sprints/<sprint-slug>/reasoning.md
```

---

## 1. Establish Sprint Context

### Required identifiers

Identify:

```text
Project: <project>
Sprint: <sprint-slug>
Sprint directory: .ultra/projects/<project>/sprints/<sprint-slug>/
```

### Required workspace sources

Read the workflow and prompt guidance before writing planning artifacts:

```text
.ultra/README.md
.ultra/prompts/create-requirements.md
.ultra/prompts/create-sprint-index.md
.ultra/prompts/create-technical-handbook.md
.ultra/prompts/create-area-reasoning.md
.ultra/prompts/create-sprint-reasoning.md
.ultra/prompts/plan-sprint.md
```

Use templates from:

```text
.ultra/system/templates/requirements.md
.ultra/system/templates/sprint-index.md
.ultra/system/templates/technical-handbook.md
.ultra/system/templates/sprint-reasoning.md
.ultra/system/templates/sprint-plan.md
```

If an artifact already exists, inspect it before deciding whether to keep it, update it, or stop for missing inputs. Do not overwrite complete artifacts casually.

---

## 2. Planning Flow Order

Follow the Ultra sprint flow order from `.ultra/README.md`:

```text
requirements
-> sprint-index
-> technical-handbook
-> area-reasoning, only if selected
-> reasoning
-> plan
```

Do not skip forward unless the prerequisite artifact already exists and is complete.

Do not create implementation files during this protocol.

---

## 3. Requirements Stage

### Inputs

Read:

```text
.ultra/projects/<project>/project-index.md
.ultra/projects/<project>/roadmap.md
.ultra/projects/<project>/docs/*.md
.ultra/system/templates/requirements.md
.ultra/projects/<project>/sprints/*/review.md, if prior sprint reviews exist
```

### Output

Write or validate:

```text
.ultra/projects/<project>/sprints/<sprint-slug>/requirements.md
```

### Rules

- Requirements are the sprint contract.
- The sprint goal must be achievable in one sprint.
- Acceptance criteria must be checkable.
- Non-goals and constraints must be explicit.
- Prior sprint review decisions must be carried forward when they constrain this sprint.
- If `requirements.md` already exists and is complete, treat it as authoritative.

Do not invent requirements from implementation preferences.

---

## 4. Sprint Index Stage

### Inputs

Read in this order:

```text
.ultra/projects/<project>/sprints/<sprint-slug>/requirements.md
.ultra/projects/<project>/project-index.md
.ultra/projects/<project>/docs/*.md
.ultra/system/templates/sprint-index.md
```

### Output

Write or validate:

```text
.ultra/projects/<project>/sprints/<sprint-slug>/sprint-index.md
```

### Rules

- The sprint index selects context. It does not make implementation decisions.
- Every selected contract, evidence report, protocol, and reasoning template must come from `project-index.md`.
- Contracts apply as flat whole documents. Do not create per-clause contract mappings in the sprint index.
- Copy selected evidence report rows from the project index rather than paraphrasing unavailable reports.
- Record excluded context and why it is excluded.
- Write an `> **Inputs Used:**` line listing exact files used.

If selected sections are empty, still contain placeholders, or reference items absent from `project-index.md`, the sprint index is not complete.

---

## 5. Technical Handbook Stage

### Inputs

Read:

```text
.ultra/projects/<project>/sprints/<sprint-slug>/sprint-index.md
.ultra/system/templates/technical-handbook.md
.ultra/projects/<project>/sprints/<sprint-slug>/requirements.md
```

Then read only the evidence reports selected in the sprint index.

### Output

Write or validate:

```text
.ultra/projects/<project>/sprints/<sprint-slug>/technical-handbook.md
```

### Rules

- The handbook distills selected evidence. It does not make final implementation decisions.
- Do not read every available report. Read only selected evidence reports.
- Do not use contracts, PRD, TRD, or other governance documents as evidence. They may provide scope context only.
- Extract patterns, trade-offs, warnings, anti-patterns, design pressures, and open questions.
- Cite specific evidence report paths and source references where available.
- Write an `> **Inputs Used:**` line listing exact files used.

If selected evidence reports cannot be located, stop and record the missing paths rather than inventing evidence.

---

## 6. Area Reasoning Stage

### When to run

Run this stage only when `sprint-index.md` explicitly selects area reasoning templates.

If no area reasoning templates are selected, record that area reasoning is skipped. Do not create reasoning files for ceremony.

### Inputs

Read:

```text
.ultra/projects/<project>/sprints/<sprint-slug>/technical-handbook.md
.ultra/projects/<project>/sprints/<sprint-slug>/requirements.md
.ultra/projects/<project>/docs/*.md
.ultra/system/reasoning/<selected-template>.md
```

### Output

For each selected area, write or validate:

```text
.ultra/projects/<project>/sprints/<sprint-slug>/reasoning/<area>.md
```

### Rules

- Create only the area reasoning documents selected by `sprint-index.md`.
- Use the selected reasoning template as the document format.
- Ground decisions in `technical-handbook.md` evidence.
- Record conclusions, rejected alternatives, risks, assumptions, and open questions.
- Keep each area reasoning document self-contained.
- Write an `> **Inputs Used:**` line listing exact files used.

Area reasoning may make area-specific decisions, but final sprint-wide decisions belong in `reasoning.md`.

---

## 7. Sprint Reasoning Stage

### Inputs

Read in this order:

```text
.ultra/projects/<project>/sprints/<sprint-slug>/requirements.md
.ultra/projects/<project>/project-index.md
.ultra/projects/<project>/docs/*.md
.ultra/projects/<project>/sprints/<sprint-slug>/sprint-index.md
.ultra/projects/<project>/sprints/<sprint-slug>/technical-handbook.md
.ultra/projects/<project>/sprints/<sprint-slug>/reasoning/*.md, if present
.ultra/system/templates/sprint-reasoning.md
```

### Output

Write or validate:

```text
.ultra/projects/<project>/sprints/<sprint-slug>/reasoning.md
```

### Rules

- Sprint reasoning is where final sprint decisions are made.
- Synthesize selected context, handbook evidence, area reasoning, and contracts.
- For each major decision, state:
  - decision made
  - requirement or contract it satisfies
  - evidence used
  - trade-off accepted
  - alternative rejected
  - technical debt or future impact
  - expected evidence
  - risk or follow-up
- Cite `technical-handbook.md`, selected evidence reports, and concrete studied repo/source references where relevant.
- If no study sources are relevant to a decision, explicitly say why.
- Include repos studied or source evidence used.
- Map decisions to applicable contracts and requirement IDs.
- Define implementation constraints.
- Write an `> **Inputs Used:**` line listing exact files used.

The reasoning must be specific enough that `plan.md` can execute without reopening architecture.

---

## 8. Sprint Plan Stage

### Inputs

Read in this order:

```text
.ultra/projects/<project>/sprints/<sprint-slug>/requirements.md
.ultra/projects/<project>/sprints/<sprint-slug>/reasoning.md
.ultra/projects/<project>/sprints/<sprint-slug>/sprint-index.md
.ultra/projects/<project>/sprints/<sprint-slug>/technical-handbook.md
.ultra/projects/<project>/sprints/<sprint-slug>/reasoning/*.md, if present
.ultra/projects/<project>/docs/*.md
.ultra/projects/<project>/project-index.md
.ultra/system/templates/sprint-plan.md
```

Open linked final reports only when a specific implementation decision needs deeper evidence. Resolve code references only for specific implementation questions.

### Output

Write or validate:

```text
.ultra/projects/<project>/sprints/<sprint-slug>/plan.md
```

### Rules

- The plan executes `reasoning.md`.
- Do not invent architecture, scope, or decisions in the plan.
- Carry forward final decisions, expected evidence, risks, assumptions, open questions, and required protocols.
- Break work into small, testable tasks.
- Separate implementation steps from verification steps.
- Include verification commands with expected results.
- Include review inputs and completion criteria.
- Write an `> **Inputs Used:**` line listing exact files used.

If the plan needs a decision not present in `reasoning.md`, return to the sprint reasoning stage instead of adding the decision only to `plan.md`.

---

## 9. Evidence And Context Discipline

When context is too large:

```text
1. Keep requirements, sprint reasoning, PRD/TRD sections relevant to the sprint, and roadmap sprint section in context.
2. Load sprint-index selected evidence before broad project materials.
3. Load final reports only when selected evidence is insufficient for a specific decision.
4. Record omitted evidence and why it was omitted.
```

Do not make ungrounded architecture decisions. Every important decision must trace to requirements, roadmap scope, selected evidence, existing project docs, or an explicitly named open question.

Avoid generic phrases such as "clean architecture" unless the plan names the concrete structure, requirement, and evidence basis.

---

## 10. Stop Conditions

Stop planning and report the blocker if:

- `project-index.md` is missing.
- `requirements.md` is missing and cannot be created from roadmap/docs.
- `sprint-index.md` selects evidence reports that cannot be found.
- `technical-handbook.md` cannot be grounded in selected evidence.
- `reasoning.md` contains placeholders or lacks final decisions.
- `plan.md` would need to invent decisions absent from `reasoning.md`.

Do not fill gaps by guessing.

---

## 11. Quality Gate

Before considering sprint planning complete, verify:

```text
- [ ] requirements.md states goal, acceptance criteria, non-goals, constraints, and dependencies.
- [ ] sprint-index.md selects only context available in project-index.md.
- [ ] technical-handbook.md cites selected evidence and defers final decisions.
- [ ] area reasoning exists only for selected areas, or is explicitly skipped.
- [ ] reasoning.md makes final sprint decisions with evidence, alternatives, risks, and expected evidence.
- [ ] plan.md executes reasoning.md without inventing new architecture or scope.
- [ ] plan.md has concrete tasks, verification commands, review inputs, and completion criteria.
- [ ] All generated artifacts include exact Inputs Used lines.
```

The sprint plan is ready only when implementation can begin from `plan.md` and `reasoning.md` without rereading every study report.
