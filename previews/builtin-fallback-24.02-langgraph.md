--- manifest ---
{
  "kind": "directory_analysis",
  "study": "agent-harness-study",
  "dimension": "24.02-interface-contract-design",
  "source": "langgraph",
  "source_kind": "directory",
  "templates": [
    "prompts/base.md",
    "templates/repo-analysis.md"
  ],
  "dimension_path": "studies/agent-harness-study/dimensions/24.02-interface-contract-design.md",
  "expected_output_path": "studies/agent-harness-study/reports/source/langgraph-24.02-interface-contract-design.md"
}
--- prompt ---
## Base Prompt

# Base Dimension — Execution Instructions

This file defines the shared workflow for every study dimension. Read it first, then read the selected dimension file for the study-specific purpose, steps, evidence, and questions.

## Hard Rules

These rules are NOT optional. Violations invalidate the study.

1. **NO cross-source filesystem access.** When studying a source, you may ONLY access files inside that source's directory. Accessing files from another source (e.g., reading `../other-source/`) is BANNED. Each source is studied in isolation.
2. **EVERY code mention MUST include a file path.** Whenever you reference a class, function, type, config key, test, or any code element, you MUST include the file path and line number (e.g., `src/core/loop.ts:42`). Line numbers are highly encouraged even for non-code mentions.
3. **Cite evidence, not vibes.** Every claim about architecture, patterns, or tradeoffs must trace back to a specific file path. If you cannot find evidence, state "No evidence found" and describe what you searched.

Violations of rules 1 or 2 require a rewrite before the study can be accepted.

## Execution Workflow

1. **Read the dimension file**
   - Read `../../prompts/base.md` for shared execution rules.
   - Read the selected `../../dimensions/{NN}-{name}.md` for the study content.

2. **Analyze each reference source**
   - For every source in the selected group, inspect the code following the selected dimension.
   - Prefer implementation, tests, configuration, and public interfaces over README-level claims.
   - For each dimension, assign a rating score (1–10) based on the rubric in the dimension file. Include the score and rationale in the output.
    - Write findings to `reports/source/{NN}-{dimension-name}/{source-name}.md` using `../../templates/repo-analysis.md`.

## Quality Bar

Each study should:

- Cite concrete evidence for major findings: file paths AND line numbers, symbols, config keys, tests, docs, or observed behavior.
- Distinguish implemented behavior from inferred intent.
- Capture tradeoffs and failure modes, not just feature presence.
- Call out missing evidence explicitly when a question cannot be answered.
- Keep recommendations specific enough to become engineering work.
- Format every evidence citation as `path/to/file.ts:NN` — not just a filename alone.

## Template Usage

- Use `../../templates/repo-analysis.md` for each per-source analysis.
- Fill every `{{placeholder}}`.
- Replace placeholders with concrete prose, tables, or bullet lists as appropriate.
- Do not leave empty sections. If there is no finding, write `No clear evidence found` and explain the search boundary.

## Output Structure

```
reports/source/{NN}-{dimension-name}/
├── {source-1}.md
├── {source-2}.md
└── ...
```

## Evidence Guidelines

Every piece of evidence MUST include a file path. Include line numbers whenever possible.

Format: `path/to/file.ts:NN` (e.g., `src/core/loop.ts:42`)

Useful evidence includes:

- Source files that implement the behavior under study — with line numbers pointing to the relevant symbols.
- Public APIs, type definitions, interfaces, schemas, or decorators — with line numbers.
- Tests that show intended behavior or edge cases — with test name and line number.
- Runtime configuration, policy files, workflow definitions, or plugin manifests — with line numbers.
- Documentation only when it is tied back to implementation or accepted as a stated design goal — include the file path.

**Bad**: "The agent loop uses an event-driven pattern."
**Good**: "The agent loop uses an event-driven pattern (`src/core/loop.ts:42-58`), dispatching events through a central bus (`src/core/bus.ts:12`)."

Avoid unsupported claims such as "robust", "enterprise-grade", "flexible", or "production-ready" unless the analysis explains the concrete mechanism that earns the label.

## Dimension

# Dimension 24.02: Interface Contract Design

## Purpose

Study how the system defines, narrows, validates, and evolves contracts between runtimes, agents, tools, providers, stores, and host applications.

## Steps

1. Identify central interfaces, protocols, abstract base classes, trait objects, schemas, or service contracts.
2. Inspect interface size, dependency direction, error contracts, context propagation, cancellation, and lifecycle methods.
3. Check whether implementations are substitutable without hidden assumptions.
4. Look for compile-time, schema-time, or runtime contract validation.
5. Assess whether contracts encode semantic guarantees or only structural shape.

## Evidence to Capture

- Interface/protocol definitions
- Adapter implementations
- Contract tests or conformance suites
- Error, cancellation, streaming, and lifecycle semantics
- Validation logic for implementations or schemas

## Questions

1. Are interfaces small, coherent, and owned by the consumer side?
2. Do contracts specify behavior, not just method signatures?
3. Can providers, tools, stores, and runtimes be replaced safely?
4. Are compatibility failures caught early by tests or validation?

## Rating

| Score | Meaning |
|-------|---------|
| 1-3 | Absent, implicit, ad-hoc, or unsafe |
| 4-6 | Present but inconsistent, weakly documented, or fragile |
| 7-8 | Clear model with tests, explicit interfaces, and operational safeguards |
| 9-10 | Mature, durable, observable, extensible, and proven under failure or scale |

> Can two independent implementations satisfy the same contract without relying on undocumented behavior?

## Output

Write findings to the exact `Expected output` path provided by UltraPlan prompt metadata. Do not write to `reports/repo`.

## Repository Analysis Template

# Source Analysis: {{source_name}}

## {{dimension_title}}

### Source Info

| Field | Value |
|-------|-------|
| Name | {{source_name}} |
| Path | `{{source_path}}` |
| Language / Stack | {{language}} |
| Analyzed | {{date}} |

## Summary

{{approach_summary}}

## Rating

{{rating}}

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
{{evidence_rows}}

## Answers to Dimension Questions

{{dimension_answers}}

## Architectural Decisions

{{key_decisions}}

## Notable Patterns

{{patterns}}

## Tradeoffs

{{tradeoffs}}

## Failure Modes / Edge Cases

{{failure_modes}}

## Future Considerations

{{future_considerations}}

## Questions / Gaps

{{gaps}}

---

Generated by `{{dimension_file}}` against `{{source_name}}`.

## Metadata

Study: agent-harness-study
Source: langgraph
Source kind: directory
Source path: studies/agent-harness-study/sources/langgraph
Expected output: studies/agent-harness-study/reports/source/langgraph-24.02-interface-contract-design.md

## Source Isolation Rules

Inspect only the selected source directory. Do not inspect unrelated workspace files, sibling sources, provider configuration, or generated reports except the explicit template and dimension inputs in this prompt.

## Citation Rules

Cite code claims with workspace-relative file paths and line numbers from the selected source directory.

