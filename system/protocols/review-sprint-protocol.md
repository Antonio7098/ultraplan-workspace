# Review Sprint Protocol

> Role: post-implementation sprint review orchestrator.
> Used by: `review.md` when selected by `sprint-index.md`.
> Reads contracts dynamically from the sprint index — do not hardcode contract names.

This protocol orchestrates a distributed review of sprint implementation against selected contracts and technical handbook guidance.

Reviews must account for sprint scope. A selected contract or handbook recommendation may apply only partially. Reviewers must distinguish real implementation problems from guidance that is intentionally out of scope, deferred, or not triggered by the change.

## Purpose

After implementation is complete, use this protocol to coordinate review of:
- Each selected contract's conformance
- Implementation against technical handbook guidance
- Final sprint review synthesis

## How It Works

1. Read `sprint-index.md` to discover selected contracts dynamically.
2. For each selected contract, spawn an independent subagent to review implementation against that contract.
3. Spawn a subagent to review the implementation against the technical-handbook.md guidance.
4. Collect all findings and synthesize into `review.md`.

---

## 1. Discover Sprint Context

### Required identifiers

```text
Project: <project>
Sprint: <sprint-slug>
Sprint directory: .ultra/projects/<project>/sprints/<sprint-slug>/
```

### Required inputs

Read these files to discover what must be reviewed:

```text
.ultra/projects/<project>/sprints/<sprint-slug>/sprint-index.md
.ultra/projects/<project>/sprints/<sprint-slug>/requirements.md
.ultra/projects/<project>/sprints/<sprint-slug>/technical-handbook.md
.ultra/projects/<project>/sprints/<sprint-slug>/reasoning.md
.ultra/projects/<project>/sprints/<sprint-slug>/plan.md
```

Do not hardcode contract names. Parse the "Selected Contracts" table from `sprint-index.md` to get the actual contracts for this sprint.

---

## 2. Parse Selected Contracts

From `sprint-index.md`, extract the "Selected Contracts" table:

```markdown
| Contract | Why Selected |
| -------- | ------------ |
| Architecture | ... |
| Errors | ... |
| Security | ... |
| Testing | ... |
| Performance | ... |
```

Each contract name maps to a contract file under `.ultra/system/contracts/`:

| Contract Name | Contract Path |
| ------------- | ------------- |
| Architecture | `.ultra/system/contracts/core/architecture.md` |
| Errors | `.ultra/system/contracts/core/errors.md` |
| Security | `.ultra/system/contracts/core/security.md` |
| Testing | `.ultra/system/contracts/core/testing.md` |
| Performance | `.ultra/system/contracts/runtime/performance.md` |

---

## 3. Spawn Contract Review Subagents

For each contract in the sprint index's "Selected Contracts" table, spawn a subagent with this prompt:

```
## Contract Review: {contract_name}

### Context
- Project: {project}
- Sprint: {sprint-slug}
- Sprint directory: .ultra/projects/{project}/sprints/{sprint-slug}/
- Contract: {contract_name}
- Contract path: {contract_path}
- Implementation directory: /home/antonioborgerees/coding/{project}/

### Your Task
Review the implementation in /home/antonioborgerees/coding/{project}/ against the requirements in {contract_path}.

Use the Architecture Review Protocol at .ultra/system/protocols/architecture-review-protocol.md as the review framework. Apply only the sections relevant to {contract_name}'s scope.

### Review Steps
1. Read {contract_path} to understand the contract requirements.
2. Read requirements.md, sprint-index.md, reasoning.md, and plan.md to understand this sprint's scope, non-goals, selected evidence, and accepted trade-offs.
3. Explore the implementation in /home/antonioborgerees/coding/{project}/internal/study/ (or the relevant module for this sprint).
4. For each contract requirement, first judge applicability to this sprint:
   - `direct`: the sprint implements behavior this requirement governs.
   - `partial`: the requirement matters, but only a subset is triggered by this sprint.
   - `not_triggered`: the requirement is selected for awareness but the implementation does not touch the relevant surface.
   - `explicitly_deferred`: sprint requirements, reasoning, or plan make this requirement a non-goal or future concern.
5. Only mark a requirement `not_satisfied` when the applicable portion is violated by the implemented scope.
6. Do not create findings for full-contract conformance when the sprint scope does not warrant it. Record those as `not_applicable` with the scope reason.
7. Return a structured finding report.

### Return Format
Return a JSON object with this structure:
{{
  "contract": "{contract_name}",
  "project": "{project}",
  "sprint": "{sprint-slug}",
  "requirements_reviewed": [
    {{
      "requirement_id": "REQ-ID",
      "title": "requirement title",
      "applicability": "direct|partial|not_triggered|explicitly_deferred",
      "scope_reason": "why this requirement does or does not apply to this sprint",
      "status": "satisfied|not_satisfied|not_applicable",
      "severity": "blocker|high|medium|low|info",
      "evidence": "file:line or test result",
      "issue": "description of problem if not satisfied",
      "recommended_action": "specific fix, follow-up, or none"
    }}
  ],
  "overall_status": "pass|pass_with_issues|fail",
  "issues": ["list of blocking issues"],
  "recommendations": ["list of non-blocking recommendations"],
  "findings_summary": "brief summary of main findings"
}}
```

---

## 4. Spawn Handbook Review Subagent

Spawn a subagent to review implementation against the technical-handbook.md guidance:

```
## Technical Handbook Implementation Review

### Context
- Project: {project}
- Sprint: {sprint-slug}
- Sprint directory: .ultra/projects/{project}/sprints/{sprint-slug}/
- Handbook: .ultra/projects/{project}/sprints/{sprint-slug}/technical-handbook.md
- Implementation directory: /home/antonioborgerees/coding/{project}/

### Your Task
Review the implementation in /home/antonioborgerees/coding/{project}/internal/study/ against the guidance in the technical-handbook.md.

The technical handbook distills evidence from studied reports into guidance for implementation. Your job is to check whether the implementation follows that guidance.

### Review Steps
1. Read the technical-handbook.md.
2. Read requirements.md, sprint-index.md, reasoning.md, and plan.md to understand this sprint's scope, non-goals, selected evidence, and accepted trade-offs.
3. Extract the "Relevant Patterns", "Trade-Offs", "Anti-Patterns And Warnings", "Design Pressures", and "Open Questions For Reasoning" sections.
4. For each pattern or guidance statement, first judge scope applicability:
   - `direct`: the implementation should follow this guidance in the current sprint.
   - `partial`: only part of the guidance applies because the sprint is narrower than the evidence context.
   - `not_triggered`: the implementation does not touch the behavior the guidance addresses.
   - `explicitly_deferred`: sprint requirements, reasoning, or plan defer this guidance to a later sprint.
5. Check whether the implementation follows the applicable portion of each pattern or guidance statement.
6. Check that anti-patterns are avoided where the implementation touches the relevant area.
7. Do not create findings for handbook guidance that is intentionally broader than current sprint scope. Record those as `not_applicable` with the scope reason.
8. Return a structured finding report.

### Return Format
Return a JSON object with this structure:
{{
  "contract": "technical_handbook",
  "project": "{project}",
  "sprint": "{sprint-slug}",
  "patterns_reviewed": [
    {{
      "pattern": "pattern name from handbook",
      "applicability": "direct|partial|not_triggered|explicitly_deferred",
      "scope_reason": "why this guidance does or does not apply to this sprint",
      "status": "followed|not_followed|not_applicable",
      "severity": "blocker|high|medium|low|info",
      "evidence": "file:line or test result",
      "issue": "description if not followed",
      "recommended_action": "specific fix, follow-up, or none"
    }}
  ],
  "antipatterns_checked": [
    {{
      "antipattern": "antipattern name from handbook",
      "applicability": "direct|partial|not_triggered|explicitly_deferred",
      "scope_reason": "why this anti-pattern check does or does not apply to this sprint",
      "status": "avoided|present|not_applicable",
      "severity": "blocker|high|medium|low|info",
      "evidence": "file:line if present",
      "issue": "description if present",
      "recommended_action": "specific fix, follow-up, or none"
    }}
  ],
  "overall_status": "pass|pass_with_issues|fail",
  "issues": ["list of blocking issues"],
  "recommendations": ["list of non-blocking recommendations"],
  "findings_summary": "brief summary of main findings"
}}
```

---

## 5. Collect And Synthesize

After all subagents return:

1. Collect all contract review reports.
2. Collect the handbook review report.
3. Synthesize into `review.md` using the template at `.ultra/system/templates/review.md`.

### Synthesis Rules

- **Contract conformance table**: Map each contract's requirements to implementation evidence.
- **Handbook conformance**: Map handbook patterns to implementation evidence.
- **Applicability table**: Record any selected contract requirements or handbook guidance that were `partial`, `not_triggered`, or `explicitly_deferred`, including the scope reason.
- **Issue classification**: Treat only violated applicable requirements as issues. Do not count justified `not_applicable` items as failures.
- **Follow-up classification**: Separate required fixes from future follow-ups. Future follow-ups must cite the sprint non-goal, reasoning decision, or accepted trade-off that justifies deferral.
- **Decision conformance**: Check `reasoning.md` decisions against what was implemented.
- **Plan execution**: Check each task in `plan.md` for completion evidence.
- **Verification**: Record `go test ./...` output and any other verification commands.
- **Deviations**: Document any deviations from `reasoning.md` or `plan.md` with reason and risk.
- **Final assessment**: Determine overall sprint status (accepted/accepted with follow-ups/blocked/rejected).

---

## 6. Stop Conditions

Stop and report a blocker if:
- `sprint-index.md` is missing or has no "Selected Contracts" table.
- Any selected contract file cannot be found.
- A subagent returns a "fail" status for a blocking issue.
- Implementation is missing entirely.

---

## 7. Output

Write the final review to:

```text
.ultra/projects/<project>/sprints/<sprint-slug>/review.md
```

The review must:
- Cite specific file paths and line numbers as evidence.
- Record deviations with reason and risk.
- Not introduce new scope or architecture.
- Be actionable — recommendations must be specific.
