--- manifest ---
{
  "kind": "directory_analysis",
  "study": "agent-harness-study",
  "dimension": "24.01-public-api-surface",
  "source": "langgraph",
  "source_kind": "directory",
  "templates": [
    "prompts/base.md",
    "templates/repo-analysis.md"
  ],
  "dimension_path": "studies/agent-harness-study/dimensions/24.01-public-api-surface.md",
  "expected_output_path": "studies/agent-harness-study/reports/source/langgraph-24.01-public-api-surface.md"
}
--- prompt ---
## Base Prompt

# Base Prompt

## Dimension

# Dimension 24.01: Public API Surface

## Purpose

Study the explicit public API exposed to application developers, operators, and extension authors.

## Steps

1. Identify the stable import paths, client objects, CLI commands, service endpoints, and documented entry points.
2. Separate public API from examples, internals, generated code, and accidental exports.
3. Inspect naming, grouping, lifecycle ownership, and discoverability of the API surface.
4. Check whether the public API is documented with minimal runnable examples.
5. Assess whether public API choices expose runtime internals or preserve abstraction boundaries.

## Evidence to Capture

- Public packages, modules, clients, command groups, or HTTP/RPC routes
- Documentation and example coverage
- API stability markers or internal/experimental labels
- Import/export boundaries
- Evidence of accidental public surface area

## Questions

1. What is the intended public API surface?
2. Is the stable API easy to distinguish from internal implementation details?
3. Does the API expose the right level of abstraction for agent harness users?
4. Are examples sufficient to use the API correctly without reading internals?

## Rating

| Score | Meaning |
|-------|---------|
| 1-3 | Absent, implicit, ad-hoc, or unsafe |
| 4-6 | Present but inconsistent, weakly documented, or fragile |
| 7-8 | Clear model with tests, explicit interfaces, and operational safeguards |
| 9-10 | Mature, durable, observable, extensible, and proven under failure or scale |

> Can a new integration identify and use the stable API without depending on implementation details?

## Output

Write findings to the exact `Expected output` path provided by UltraPlan prompt metadata. Do not write to `reports/repo`.

## Repository Analysis Template

# Repository Analysis Task

You are producing one source report for an UltraPlan study.

## Hard Requirements

- Write the final report to the exact `Expected output` path shown in the prompt metadata.
- Create parent directories before writing the report.
- Ignore any conflicting output path mentioned inside dimension text.
- Keep the investigation bounded: inspect at most 8 files unless the answer is impossible.
- Prefer documentation and public interface files. Include at least one citation to a `.md`, `.go`, or `.txt` file with a line number, such as `README.md:1`.
- Do not run tests, package installs, network operations, or long recursive scans.
- The report must be self-contained markdown and must satisfy the structure below.

## Required Report Shape

```markdown
# <source> - <dimension>

## Source Info

- Source: <source name>
- Dimension: <dimension id and title>
- Repository/path inspected: <path>

## Summary

<Short factual summary of what this source shows for the dimension.>

## Evidence

- `<file.md:line>` <what the evidence shows>
- `<file.go:line>` <what the evidence shows, if applicable>

## Questions and Answers

Q: <study question or dimension question>
A: <answer grounded in the evidence above>

Q: <another relevant question>
A: <answer grounded in the evidence above>

## Rating

Rating: <integer 1-10>/10

<One sentence explaining the rating.>

## Gaps

- <Important uncertainty or missing evidence.>
```

## Validation Notes

UltraPlan validation requires a top-level `#` heading, source information, summary, question/answer content, and a parseable rating line like `Rating: 7/10`.

## Metadata

Study: agent-harness-study
Source: langgraph
Source kind: directory
Source path: studies/agent-harness-study/sources/langgraph
Expected output: studies/agent-harness-study/reports/source/langgraph-24.01-public-api-surface.md

## Source Isolation Rules

Inspect only the selected source directory. Do not inspect unrelated workspace files, sibling sources, provider configuration, or generated reports except the explicit template and dimension inputs in this prompt.

## Citation Rules

Cite code claims with workspace-relative file paths and line numbers from the selected source directory.

