# Sprint Requirements: 13 Summary And Code Extraction

> Project: `ultraplan-go`
> Sprint: `13-summary-and-code-extraction`
> Purpose: the authoritative, human-readable sprint contract. All other sprint artifacts must satisfy these requirements.

## Sprint Goal

Implement standalone study summary generation and code-reference extraction so completed study outputs can be reviewed, compared, and audited from deterministic text, CSV, and JSON surfaces.

## Required Outputs

| Output | Path | Description |
| ---------- | -------- | --------------- |
| Requirements | `projects/ultraplan-go/sprints/13-summary-and-code-extraction/requirements.md` | This sprint contract defining scope, outputs, acceptance criteria, constraints, dependencies, and review expectations. |
| Sprint index | `projects/ultraplan-go/sprints/13-summary-and-code-extraction/sprint-index.md` | Selected contracts, evidence reports, docs, protocols, and exclusions for summary generation and code extraction. |
| Technical handbook | `projects/ultraplan-go/sprints/13-summary-and-code-extraction/technical-handbook.md` | Sprint-specific implementation guidance for summary semantics, rating warnings, citation parsing, source resolution, snippet output, and safe diagnostics. |
| Area reasoning: summary generation | `projects/ultraplan-go/sprints/13-summary-and-code-extraction/reasoning/summary-generation.md` | Detailed reasoning for summary command behavior, matrix shape, missing/inapplicable semantics, warnings, sorting, and carry-forward decisions. |
| Area reasoning: code extraction | `projects/ultraplan-go/sprints/13-summary-and-code-extraction/reasoning/code-extraction.md` | Detailed reasoning for citation syntax, source table parsing, reference resolution, unresolved reporting, output formats, and module boundaries. |
| Consolidated reasoning | `projects/ultraplan-go/sprints/13-summary-and-code-extraction/reasoning.md` | Decision log summarizing the chosen design and mapping decisions to requirements and selected evidence. |
| Implementation plan | `projects/ultraplan-go/sprints/13-summary-and-code-extraction/plan.md` | Ordered implementation checklist with verification commands and review evidence expectations. |
| Sprint review | `projects/ultraplan-go/sprints/13-summary-and-code-extraction/review.md` | Completed review recording acceptance evidence, deviations, residual risks, and follow-ups. |
| Study summary implementation | `/home/antonioborgerees/coding/ultraplan-go/internal/study/summary.go` | Study-owned deterministic `summary.csv` generation, score matrix creation, warning collection, total sorting, and atomic write behavior. |
| Study summary domain updates | `/home/antonioborgerees/coding/ultraplan-go/internal/study/domain.go` | Minimal request/result/warning types required for standalone summary behavior and reviewable command output. |
| Study summary service wiring | `/home/antonioborgerees/coding/ultraplan-go/internal/study/service.go` | Public study service method required to regenerate summaries without launching runtime work. |
| Study summary CLI wiring | `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_commands.go` | Thin CLI handling for `ultraplan study <study> summary`, including help text, output rendering, and exit-code mapping. |
| Study summary tests | `/home/antonioborgerees/coding/ultraplan-go/internal/study/summary_test.go` | Unit and fixture tests for matrix generation, totals, ordering, warnings, missing reports, ambiguous ratings, inapplicable pairs, and atomic writes. |
| Study summary command tests | `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_summary_commands_test.go` | Command-level tests for `study summary` help, success, warnings, failures, output paths, and non-runtime behavior. |
| Code extraction domain | `/home/antonioborgerees/coding/ultraplan-go/internal/codeextract/domain.go` | Code extraction request/result/source/reference/snippet/unresolved types owned by the `codeextract` module. |
| Source table parser | `/home/antonioborgerees/coding/ultraplan-go/internal/codeextract/parser.go` | Parser for report source tables and supported inline code citation forms. |
| Reference resolver | `/home/antonioborgerees/coding/ultraplan-go/internal/codeextract/resolver.go` | Resolver for source-relative, source-prefixed, and basename fallback citation paths with ignored-directory exclusions. |
| Extraction service | `/home/antonioborgerees/coding/ultraplan-go/internal/codeextract/service.go` | Orchestrates report loading, citation parsing, source resolution, line extraction, unresolved reporting, and text/JSON result construction. |
| Code extraction CLI wiring | `/home/antonioborgerees/coding/ultraplan-go/internal/app/code_commands.go` | Thin top-level command handling for `ultraplan code <report>...`, including `--json` and optional output-file behavior. |
| App command dispatch updates | `/home/antonioborgerees/coding/ultraplan-go/internal/app/app.go` | Top-level dispatch/help updates needed to expose the `code` command without moving product behavior into `internal/app`. |
| Code extraction tests | `/home/antonioborgerees/coding/ultraplan-go/internal/codeextract/codeextract_test.go` | Fixture tests for source table parsing, citation parsing, line/range/list extraction, path resolution, unresolved references, and text/JSON output. |
| Code command tests | `/home/antonioborgerees/coding/ultraplan-go/internal/app/code_commands_test.go` | Command-level tests for arguments, multiple reports, JSON/text output, output-file writes, unresolved reporting, and exit statuses. |

## Acceptance Criteria

- [ ] `ultraplan study <study> summary` is available in help output and regenerates `studies/<study>/summary.csv` without invoking runtime, agentwrap, OpenCode, network access, or subprocess execution.
- [ ] The summary command discovers dimensions and sources using existing study discovery behavior and uses existing applicability rules so inapplicable Markdown document source/dimension pairs are represented distinctly and are never counted as missing expected reports.
- [ ] `summary.csv` has a deterministic header containing `source`, one column per discovered dimension in deterministic dimension order, and `total`.
- [ ] Summary cells use valid parsed rating values for valid reports, empty cells for missing expected reports or missing/ambiguous ratings, and `N/A` for inapplicable Markdown source/dimension pairs.
- [ ] Total score is the sum of valid rating values only; missing reports, missing ratings, ambiguous ratings, and inapplicable cells do not contribute to totals.
- [ ] Summary rows are sorted by total score descending, with source name ascending as the deterministic tie-breaker.
- [ ] Missing expected reports, missing ratings, and ambiguous ratings produce warnings that identify the source, dimension, report path when known, and reason.
- [ ] Inapplicable Markdown source/dimension pairs do not produce missing-report warnings and do not make the summary command fail.
- [ ] Summary generation writes `summary.csv` atomically and a write failure is reported as a non-zero command outcome rather than a success with only a silent or low-visibility warning.
- [ ] Re-running summary generation with identical inputs produces byte-identical CSV output and deterministic human output.
- [ ] The summary implementation reuses the existing study rating parser and report path helpers instead of adding duplicate rating or path parsing logic.
- [ ] `ultraplan code <report>...` is available in top-level help output and accepts one or more Markdown report paths.
- [ ] Code extraction parses a report sources table containing source names and source paths, including rows shaped like index `1`, source `source-name`, and path `sources/source-name`.
- [ ] Code extraction rejects reports that require source resolution but do not contain a parseable sources table, returning a non-zero validation-style outcome with actionable diagnostics.
- [ ] Citation parsing supports inline code references shaped as `path/to/file.go:42`, `path/to/file.go:42-58`, and `path/to/file.go:42,47,53`.
- [ ] Citation parsing tolerates the legacy range dash form documented in the TRD and normalizes extracted output to the standard hyphen range form.
- [ ] Reference resolution first tries citation paths relative to each parsed source root, then source-prefixed path stripping, then basename search inside source roots.
- [ ] Basename fallback search excludes `.git`, `node_modules`, and other explicitly documented ignored directories.
- [ ] Line extraction supports single lines, closed ranges, and explicit line lists, preserving requested line numbers in output.
- [ ] Out-of-range line requests, missing files, ambiguous basename matches, and malformed line specs are reported as unresolved references with source report context and do not silently disappear.
- [ ] Text output includes report path, source name, cited path, requested line spec, resolved path when found, rendered line-numbered snippets, and an unresolved reference summary.
- [ ] JSON output is supported for `ultraplan code <report>... --json` and exposes deterministic structured fields for reports, sources, references, resolution status, snippets, and unresolved diagnostics.
- [ ] `ultraplan code <report>... --output <path>` writes the extraction result to the requested file and reports filesystem failures with a non-zero outcome.
- [ ] Multiple report inputs are processed deterministically in argument order, and one report's unresolved references do not prevent other reports from being inspected.
- [ ] The code extraction command returns success when all references resolve, a distinct non-zero validation/partial outcome when one or more references are unresolved, and a usage/workspace/filesystem outcome for invalid arguments or unreadable inputs according to existing app exit-code conventions.
- [ ] Code extraction behavior is owned by `internal/codeextract`; `internal/app` only parses command arguments, invokes the module, renders output, and maps exit codes.
- [ ] Study summary behavior remains owned by `internal/study`; no global `internal/summary`, `internal/reports`, `internal/parser`, `internal/resolver`, or `internal/validation` package is introduced.
- [ ] `internal/platform/runtime` remains generic and has no dependency on `internal/study`, `internal/codeextract`, or CLI command packages.
- [ ] Human and JSON output do not print secrets, sensitive environment values, prompt bodies, embedded Markdown document bodies unrelated to extracted snippets, or unsafe native runtime payloads.
- [ ] Normal tests use local fixtures only; `go test ./...` passes without OpenCode, provider credentials, network access, or real runtime execution.
- [ ] The CLI builds with `go build ./cmd/ultraplan`.

## Non-Goals

- Implementing new runtime execution, prompt composition, report generation, repair, retry, fallback, or provider behavior.
- Implementing `ultraplan study <study> run-loop`, stale running task recovery, per-study lock files, force unlock, or durable retry/cancellation logic if not already completed by Sprint 12.
- Changing the report templates or requiring existing generated reports to be rewritten before summary or code extraction can inspect them.
- Guaranteeing semantic correctness of generated prose beyond rating parsing, source table parsing, citation syntax parsing, and reference resolution rules.
- Supporting citation forms outside the TRD-supported single-line, range, and explicit line-list forms.
- Resolving remote repository URLs, fetching missing sources, cloning repositories, or extracting snippets from non-local sources.
- Adding a browser UI, TUI dashboard, hosted service, remote artifact store, plugin system, or workflow engine.
- Implementing target scaffolding, sprint planning, sprint execution, or any target/sprint CLI commands.
- Declaring release-wide stable JSON schema compatibility; this sprint requires deterministic tested JSON for the `code` command only, while broader public JSON stability remains Sprint 14 scope.
- Adding real OpenCode smoke tests to the default test suite.

## Constraints

- Summary behavior must stay in `internal/study`; code extraction behavior must stay in `internal/codeextract`; CLI adapters must stay thin and must not own product parsing or resolution logic.
- `internal/platform/runtime` must not import product modules and this sprint must not call agentwrap, OpenCode, or provider runtimes.
- Existing study discovery, source applicability, report path, validation, and rating parsing behavior must be reused rather than duplicated.
- Missing expected reports and inapplicable Markdown pairs must be represented differently in both implementation and tests.
- Summary writes and optional extraction output writes must use explicit paths, safe permissions, and workspace/path containment where applicable; generated summary writes must be atomic.
- All user-facing diagnostics must preserve enough context to repair the problem while avoiding prompt, secret, environment, and unsafe runtime payload leakage.
- Reference resolution must only read local files under parsed source roots or explicitly accepted report-relative/workspace-relative paths; it must not traverse outside intended roots through `..` escapes.
- Basename fallback must be bounded to parsed source roots and must skip ignored directories to avoid scanning unrelated dependency or VCS content.
- Tests must be deterministic, fixture-based, and offline by default.
- Any changes to existing summary behavior from Sprint 11, including total-score sorting and loud summary write failures, must be documented in reasoning and covered by tests.

## Dependencies

| Prior Sprint / Output | Required For | Notes |
| ------------------------ | ------------------- | --------- |
| Sprint 5: Markdown source discovery and applicability | Summary missing vs inapplicable semantics | Summary generation must reuse directory-vs-Markdown source discovery and `GetApplicableSources` behavior. |
| Sprint 6: Report validation and rating parsing | Summary score extraction and warning behavior | Summary generation must reuse strict rating parsing and report path conventions. |
| Sprint 7: Prompt composition | Report artifact conventions | Code extraction and summary inspect reports produced by existing analysis/synthesis workflows; no prompt changes are required. |
| Sprint 8: Run-state persistence and status primitives | Optional context for completed study outputs | Summary and code extraction must not require direct state mutation; they may inspect existing artifacts independently. |
| Sprint 9: Agentwrap/OpenCode runtime integration | Runtime boundary preservation | This sprint must not bypass or extend runtime execution; platform runtime remains untouched except for import-boundary verification. |
| Sprint 10: Single analysis and synthesis | Valid report generation path | Summary and code extraction operate on reports produced and validated by single-run and synthesis commands. |
| Sprint 11: `run-all` batch execution | Existing summary implementation and carry-forward decisions | Sprint 13 must harden standalone summary behavior, total sorting, warnings, and loud write failures while preserving deterministic batch compatibility. |
| Sprint 12: durable `run-loop`, retry, and cancellation roadmap scope | Completed-output workflow context | If Sprint 12 is complete, summary/code commands must remain compatible with run-loop outputs; if Sprint 12 artifacts are absent, this sprint must not invent durable orchestration to compensate. |
| Project docs: PRD, TRD, ARCHITECTURE | Scope and architectural rules | Study-side-only scope, module ownership, code extraction ownership, validation, security, performance, and testing rules remain binding. |

## Review Expectations

| What | How Verified |
| ------------- | ----------------------- |
| Required sprint artifacts exist and match the standard artifact chain | Inspect `projects/ultraplan-go/sprints/13-summary-and-code-extraction/` for `requirements.md`, `sprint-index.md`, `technical-handbook.md`, `reasoning/`, `reasoning.md`, `plan.md`, and `review.md`. |
| Summary CLI surface is implemented | Run or test `ultraplan study <study> summary --help` and command-level tests for success, warning, usage failure, missing study, and write failure cases. |
| Summary CSV semantics are correct | Unit tests compare exact CSV output for valid ratings, missing reports, missing ratings, ambiguous ratings, inapplicable Markdown pairs, total calculation, total-descending sort, and repeated-run determinism. |
| Summary write failures are loud | Tests simulate write failure and verify the command returns a non-zero outcome and preserves prior valid summary content when applicable. |
| Code CLI surface is implemented | Run or test `ultraplan code --help`, `ultraplan code <report>...`, `ultraplan code <report>... --json`, and `ultraplan code <report>... --output <path>`. |
| Source table parsing is correct | Fixture tests cover valid source rows, workspace-relative paths, report-relative paths, missing source table, malformed rows, and duplicate/ambiguous sources. |
| Citation parsing is correct | Fixture tests cover single-line, range, list, legacy dash range, malformed line specs, duplicate citations, and citations that should not be parsed. |
| Reference resolution is correct | Fixture tests cover source-relative paths, source-prefixed paths, basename fallback, ignored directories, missing files, ambiguous basename matches, path escape attempts, and out-of-range lines. |
| Text and JSON outputs are deterministic and safe | Golden or exact-output tests verify output ordering, line-numbered snippets, unresolved summaries, JSON field stability for this command, and absence of secrets/prompts/unsafe runtime payloads. |
| Architecture boundaries are preserved | Inspect imports and changed paths; verify summary behavior lives in `internal/study`, extraction behavior lives in `internal/codeextract`, CLI wiring is thin, and `internal/platform/runtime` has no product imports. |
| No future-horizon or runtime scope leaked in | Search/review implementation for absence of new runtime execution, prompt execution, run-loop orchestration, target commands, sprint commands, hosted/UI features, and generic workflow abstractions. |
| Verification passes | Run `go test ./...` and `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`; record results in `review.md`. |
