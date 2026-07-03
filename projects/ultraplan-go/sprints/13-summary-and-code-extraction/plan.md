# Sprint Plan: 13 Summary And Code Extraction

> Project: `ultraplan-go`
> Sprint: `13-summary-and-code-extraction`
> Source: `projects/ultraplan-go/sprints/13-summary-and-code-extraction/reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/13-summary-and-code-extraction/requirements.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/sprints/13-summary-and-code-extraction/sprint-index.md`, `projects/ultraplan-go/sprints/13-summary-and-code-extraction/technical-handbook.md`, `projects/ultraplan-go/sprints/13-summary-and-code-extraction/reasoning.md`, `templates/sprint-plan.md`, `system/protocols/architecture-review-protocol.md`, `system/protocols/review-sprint-protocol.md`. Area reasoning omitted because no files matched `projects/ultraplan-go/sprints/13-summary-and-code-extraction/reasoning/*.md`.

This plan executes `reasoning.md`. It must not invent architecture, scope, or decisions beyond that document.

## Reasoning Source

- **Sprint Reasoning:** `projects/ultraplan-go/sprints/13-summary-and-code-extraction/reasoning.md`
- **Sprint Index:** `projects/ultraplan-go/sprints/13-summary-and-code-extraction/sprint-index.md`
- **Technical Handbook:** `projects/ultraplan-go/sprints/13-summary-and-code-extraction/technical-handbook.md`
- **Area Reasoning:** none present; consolidated reasoning records this absence.

## Sprint Status

- **Status:** complete
- **Owner:** implementation agent
- **Start Date:** 2026-06-12
- **Completion Date:** 2026-06-12

## Decisions To Execute

| Decision | Source Section | Execution Implication |
| --- | --- | --- |
| Study-Owned Summary Matrix Semantics | `reasoning.md#decision-1-study-owned-summary-matrix-semantics` | Implement deterministic `summary.csv` generation in `internal/study`, using discovered dimension order, valid parsed ratings, empty invalid/missing cells, `N/A` inapplicable cells, valid-only totals, and total-descending/source-ascending row sorting. |
| Standalone Summary Command And Loud Writes | `reasoning.md#decision-2-standalone-summary-command-and-loud-writes` | Add `ultraplan study <study> summary` as thin `internal/app` wiring over a study service method; write summaries atomically and fail non-zero on write or command failures without invoking runtime work. |
| Code Extraction Owned By `internal/codeextract` | `reasoning.md#decision-3-code-extraction-owned-by-internalcodeextract` | Add `internal/codeextract` domain, parser, resolver, and service files; keep CLI parsing/rendering in `internal/app`; do not introduce global parser, resolver, report, validation, or summary packages. |
| Source Table, Citation, And Resolution Rules | `reasoning.md#decision-4-source-table-citation-and-resolution-rules` | Parse required inline code reference forms, source tables, workspace/report-relative roots, source-root paths, source-prefixed paths, and basename fallback with unresolved diagnostics for malformed, missing, ambiguous, escape, and out-of-range references. |
| Deterministic Extraction Output And Exit Outcomes | `reasoning.md#decision-5-deterministic-extraction-output-and-exit-outcomes` | Render deterministic text and sprint-local JSON output; process reports in argument order; map all-resolved, partial unresolved, validation, usage, and filesystem outcomes to existing app exit conventions. |
| Local-Only Safe Path Handling And Bounded Fallback Search | `reasoning.md#decision-6-local-only-safe-path-handling-and-bounded-fallback-search` | Restrict file reads to parsed source roots; reject or mark path escapes unresolved; cache basename lookup per invocation; skip `.git`, `.hg`, `.svn`, `node_modules`, and `.ultraplan`. |
| Verification Scope And Architecture Review Gates | `reasoning.md#decision-7-verification-scope-and-architecture-review-gates` | Add deterministic unit, fixture, command-level, and exact-output tests; verify with `go test ./...`, `go build ./cmd/ultraplan`, changed-path review, import review, and safe-output review. |

## Requirements / Contracts To Satisfy

| Contract / Requirement ID | Required Behavior | Evidence Planned |
| --- | --- | --- |
| AC-01 through AC-11 | Standalone summary command, deterministic CSV shape, missing/inapplicable semantics, warnings, helper reuse, atomic write, deterministic output, no runtime invocation. | `internal/study/summary_test.go`, `internal/app/study_summary_commands_test.go`, exact CSV/output assertions, write-failure tests, import/runtime boundary review. |
| AC-12 through AC-25 | Top-level code command, source table parsing, citation parsing, resolution, snippets, unresolved diagnostics, text/JSON output, output-file writes, multi-report processing, exit outcomes. | `internal/codeextract/codeextract_test.go`, `internal/app/code_commands_test.go`, exact text/JSON assertions, resolver fixtures, output write failure tests. |
| AC-26 through AC-28 | Summary remains in `internal/study`; extraction remains in `internal/codeextract`; `internal/app` stays thin; `internal/platform/runtime` remains product-free. | Changed-file review, import review, absence of new global product parser/resolver packages, architecture review protocol. |
| AC-29 | Human and JSON output avoid secrets, sensitive env, prompt bodies, unrelated Markdown bodies, and unsafe runtime payloads. | Exact output review, fixture tests with sensitive-looking unrelated content, no runtime diagnostics in summary/code output. |
| AC-30 through AC-31 | Normal tests are offline/local-fixture based and the CLI builds. | `go test ./...` and `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`. |
| Architecture | Product behavior belongs to owning modules; runtime is generic; CLI adapters are thin. | Architecture review against `system/protocols/architecture-review-protocol.md`. |
| Errors / CLI Surface | Outcomes are actionable and scriptable with meaningful exit-code categories. | Command tests for success, usage, validation, partial, workspace/filesystem, and write failures. |
| Security / Performance | Extraction is local-only, contained, bounded, deterministic, and avoids unbounded scans or accidental leakage. | Resolver tests for containment, ignored dirs, path escapes, ambiguous basenames, deterministic traversal, and safe diagnostics. |
| Testing / Documentation / Persistence And Migrations | Observable outputs are fixture-backed, help text is covered, and generated writes are explicit and safe. | Golden/exact output tests, help-output tests, atomic write tests, optional output-file write tests, review evidence in `review.md`. |

## Tasks

- [x] **Task 1: Inventory Existing Study And App Surfaces**
  > Executes: `Decision 1`, `Decision 2`, `Decision 7`, `AC-01 through AC-11`, `AC-26`, `AC-28`
  - [x] Inspect existing study discovery, source applicability, report path, validation, and rating parser helpers in `/home/antonioborgerees/coding/ultraplan-go/internal/study/`.
  - [x] Inspect existing app command dispatch, help rendering, exit-code mapping, workspace resolution, and command-test patterns in `/home/antonioborgerees/coding/ultraplan-go/internal/app/`.
  - [x] Record any missing or incompatible helper as an implementation note before adding new code; prefer minimal exported helper exposure over duplicated parsing or path logic.
  - [x] Confirm no implementation step needs runtime, agentwrap, OpenCode, network access, subprocess execution, run-loop mutation, prompt changes, target commands, or sprint commands.

- [x] **Task 2: Implement Study Summary Domain And Service Behavior**
  > Executes: `Decision 1`, `Decision 2`, `AC-01 through AC-11`, `AC-26`, `AC-30`
  - [x] Add minimal summary request/result/warning types in `internal/study/domain.go` only if existing types cannot represent command output and warning details.
  - [x] Implement or harden `internal/study/summary.go` to discover dimensions and sources deterministically, reuse applicability filtering, reuse report path helpers, and reuse the existing rating parser.
  - [x] Represent valid ratings as numeric cells, missing expected reports as empty cells with warnings, missing ratings as empty cells with warnings, ambiguous ratings as empty cells with warnings, and inapplicable Markdown source/dimension pairs as `N/A` without missing-report warnings.
  - [x] Compute totals from valid rating values only and sort rows by total descending with source name ascending as the deterministic tie-breaker.
  - [x] Write `studies/<study>/summary.csv` atomically using an explicit path and safe permissions; surface write failures as command-failing errors, not warnings.
  - [x] Expose one focused `study.Service` summary method in `internal/study/service.go` for CLI use without launching or depending on runtime execution.

- [x] **Task 3: Test Summary Semantics And Writes**
  > Executes: `Decision 1`, `Decision 2`, `Decision 7`, `AC-02 through AC-11`, `AC-30`
  - [x] Add `internal/study/summary_test.go` fixture/unit tests for exact CSV header, dimension ordering, source ordering, total sorting, tie sorting, valid rating totals, repeated-run byte identity, and warning ordering.
  - [x] Cover missing expected reports, missing ratings, ambiguous ratings, and inapplicable Markdown pairs as distinct cases in both implementation result data and expected output.
  - [x] Test that inapplicable pairs do not produce missing-report warnings and do not fail summary generation.
  - [x] Test atomic write success and simulated write failure, including preservation of prior content where applicable.
  - [x] Include a test or review assertion that existing rating parser and report path helpers are used rather than duplicated.

- [x] **Task 4: Wire Standalone Summary CLI**
  > Executes: `Decision 2`, `Decision 7`, `AC-01`, `AC-08 through AC-11`, `AC-26 through AC-31`
  - [x] Update `internal/app/study_commands.go` so `ultraplan study <study> summary` appears in help and dispatches to the study service summary method.
  - [x] Keep command handling limited to argument parsing, workspace/study resolution, service invocation, output rendering, and exit mapping.
  - [x] Render deterministic human output that includes summary path and warnings with source, dimension, report path when known, and reason.
  - [x] Add `internal/app/study_summary_commands_test.go` for help output, success, warnings, missing/invalid study or usage failures, write failure, output path rendering, deterministic output, and no runtime invocation.

- [x] **Task 5: Implement Code Extraction Domain, Parser, And Resolver**
  > Executes: `Decision 3`, `Decision 4`, `Decision 6`, `AC-12 through AC-19`, `AC-24`, `AC-27 through AC-31`
  - [x] Add `internal/codeextract/domain.go` with focused request, result, report, source, reference, snippet, unresolved diagnostic, and status/outcome types needed by service and renderers.
  - [x] Add `internal/codeextract/parser.go` to parse report source tables with index/source/path rows, source names, workspace-relative paths, report-relative paths, supported inline code citations, single lines, closed ranges, explicit line lists, and legacy dash range normalization.
  - [x] Treat malformed line specs and malformed source rows as actionable diagnostics rather than silently dropping references.
  - [x] Add `internal/codeextract/resolver.go` to normalize and contain source roots, try source-root-relative paths, try source-prefixed path stripping, then perform deterministic basename fallback inside parsed roots only.
  - [x] Cache basename lookup work per extraction invocation and skip `.git`, `.hg`, `.svn`, `node_modules`, and `.ultraplan` during fallback traversal.
  - [x] Report unresolved references for missing files, ambiguous basename matches, path escape attempts, out-of-range lines, and malformed line specs with report/source/citation context.

- [x] **Task 6: Implement Code Extraction Service And Output Rendering**
  > Executes: `Decision 3`, `Decision 4`, `Decision 5`, `Decision 6`, `AC-12 through AC-25`, `AC-29 through AC-31`
  - [x] Add `internal/codeextract/service.go` to orchestrate report loading, citation parsing, source table parsing, source root resolution, snippet extraction, unresolved aggregation, and deterministic result construction.
  - [x] Process multiple report inputs in argument order and continue inspecting later reports even when earlier reports contain unresolved references.
  - [x] Reject reports that contain supported citations requiring resolution but lack a parseable sources table with a validation-style outcome; allow reports with no supported code references to produce an empty successful inspection result unless existing app conventions require a clearer no-reference status.
  - [x] Render default text output with report path, source name, cited path, requested normalized line spec, resolved path when found, line-numbered snippets, and unresolved summary.
  - [x] Render `--json` output with deterministic sprint-local structured fields for reports, sources, references, resolution status, snippets, and unresolved diagnostics without claiming release-wide schema stability.
  - [x] Ensure output includes only requested snippets and safe context, not prompt bodies, sensitive environment values, unsafe runtime payloads, or unrelated embedded Markdown document bodies.

- [x] **Task 7: Wire Top-Level Code CLI**
  > Executes: `Decision 3`, `Decision 5`, `Decision 7`, `AC-12`, `AC-20 through AC-31`
  - [x] Update `internal/app/app.go` and add `internal/app/code_commands.go` so top-level help exposes `ultraplan code <report>...` with `--json` and `--output <path>`.
  - [x] Keep command handling limited to argument parsing, workspace path setup, service invocation, renderer selection, optional output-file write, and exit mapping.
  - [x] Map outcomes to existing app conventions while preserving reasoning categories: success when all references resolve, partial completion for unresolved references, validation for malformed report/source-table preconditions, usage for invalid arguments, and workspace/filesystem for unreadable inputs or output write failures.
  - [x] Write optional extraction output to the requested explicit path with safe permissions and non-zero filesystem failure handling.
  - [x] Add `internal/app/code_commands_test.go` for help, no reports usage error, multiple reports in argument order, all-resolved success, unresolved partial outcome, malformed/missing table validation outcome, JSON output, text output, output-file writes, write failure, unreadable input, and safe diagnostics.

- [x] **Task 8: Test Code Extraction Module Thoroughly**
  > Executes: `Decision 4`, `Decision 5`, `Decision 6`, `Decision 7`, `AC-13 through AC-25`, `AC-29 through AC-31`
  - [x] Add `internal/codeextract/codeextract_test.go` fixture tests for valid source rows, workspace-relative source roots, report-relative source roots, missing source table when citations exist, reports without supported citations, malformed rows, duplicate or ambiguous sources, and deterministic source ordering.
  - [x] Cover single-line, range, explicit line-list, legacy dash range, malformed line specs, duplicate citations, and prose or Markdown-link references that should not be parsed.
  - [x] Cover source-root-relative resolution, source-prefixed stripping, basename fallback, ignored directories, missing files, ambiguous basename matches, path escapes, symlink or cleaned-path containment where supported by current helpers, and out-of-range line requests.
  - [x] Cover line-numbered snippet rendering that preserves requested line numbers for single lines, closed ranges, and explicit lists.
  - [x] Cover exact or golden-style text and JSON output ordering, unresolved summaries, and absence of unsafe unrelated content.

- [x] **Task 9: Architecture, Boundary, And Non-Goal Review Before Final Verification**
  > Executes: `Decision 3`, `Decision 6`, `Decision 7`, `AC-26 through AC-31`
  - [x] Inspect changed files and imports to confirm summary behavior is in `internal/study`, extraction behavior is in `internal/codeextract`, and `internal/app` only adapts CLI concerns.
  - [x] Confirm no new `internal/summary`, `internal/reports`, `internal/parser`, `internal/resolver`, `internal/validation`, plugin, workflow engine, target, sprint, or generic future-horizon package was introduced.
  - [x] Confirm `internal/platform/runtime` has no dependency on `internal/study`, `internal/codeextract`, or command packages.
  - [x] Search/review summary and code commands for absence of runtime, agentwrap/OpenCode launches, network access, subprocess execution, prompt composition changes, run-loop orchestration, and default real-runtime test requirements.
  - [x] Review text and JSON output for secret, environment, prompt, raw native payload, and unrelated Markdown leakage.

- [x] **Task 10: Final Verification And Review Handoff**
  > Executes: `Decision 7`, `AC-30`, `AC-31`, review expectations in `requirements.md`
  - [x] Run `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go` and record the result for `review.md`.
  - [x] Run `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go` and record the result for `review.md`.
  - [x] If verification fails, fix applicable implementation or test issues before review; if a failure is unrelated or blocked by pre-existing state, document the evidence and blocker clearly.
  - [x] Prepare review evidence for architecture review, sprint review, exact-output tests, command tests, import review, safe-output review, and any deviations from `reasoning.md` or this plan.

## Evidence Checklist

- [x] Summary matrix tests prove exact CSV shape, deterministic ordering, rating totals, missing/invalid/inapplicable semantics, warnings, atomic writes, and repeated-run determinism.
- [x] Summary command tests prove help, success, warnings, usage/workspace/write failures, output path reporting, and non-runtime behavior.
- [x] Code extraction module tests prove source table parsing, citation parsing, line extraction, resolution order, ignored dirs, containment, unresolved diagnostics, text rendering, and JSON rendering.
- [x] Code command tests prove help, arguments, multiple reports, `--json`, `--output`, output write failure, unresolved partial outcome, validation outcome, and safe diagnostics.
- [x] Architecture review evidence proves product behavior stayed in owning modules and no forbidden global packages or runtime dependencies were introduced.
- [x] Verification evidence includes `go test ./...` and `go build ./cmd/ultraplan` results.
- [x] Deviations from `reasoning.md` or this plan are recorded before implementation continues.
- [x] Required review protocols have usable evidence, with the stale sprint-review protocol path noted.

## Verification Commands

| Check | Command | Expected Result |
| --- | --- | --- |
| Full Go test suite | `go test ./...` | Passes without OpenCode, provider credentials, network access, or real runtime execution. |
| CLI build | `go build ./cmd/ultraplan` | Builds the `ultraplan` CLI successfully. |
| Architecture import review | Inspect changed imports and package paths | `internal/study` owns summary, `internal/codeextract` owns extraction, `internal/app` is thin, and `internal/platform/runtime` has no product imports. |
| Non-goal review | Inspect changed files and command surfaces | No runtime execution, prompt changes, run-loop orchestration, target/sprint commands, remote fetching, plugin/workflow engine, or release-wide JSON schema commitment. |
| Safe output review | Inspect exact text/JSON output tests and renderers | Outputs include requested audit context only and omit secrets, env values, prompt bodies, unsafe runtime payloads, and unrelated Markdown bodies. |

## Risks And Blockers

| Risk / Blocker | Source | Mitigation | Status |
| --- | --- | --- | --- |
| Existing study helpers may be missing or not directly reusable. | `reasoning.md#assumptions-and-risks` | Existing discovery, applicability, report path, and rating helpers were reused; only warning metadata and service wiring were added. | closed |
| Existing app exit-code names may differ from reasoning's numeric categories. | `reasoning.md#decision-5-deterministic-extraction-output-and-exit-outcomes` | Existing app constants already include usage `2`, workspace `4`, validation `5`, and partial `8`; command tests cover partial and validation outcomes. | closed |
| Source table parser may reject real reports with slight heading or formatting variation. | `reasoning.md#assumptions-and-risks` | Parser supports required row shape and backtick paths; broader compatibility remains future scope. | mitigated |
| Basename fallback can be slow or ambiguous on large source roots. | `reasoning.md#decision-6-local-only-safe-path-handling-and-bounded-fallback-search` | Fallback is scoped to parsed roots, cached per invocation, skips ignored dirs, and reports ambiguity. | mitigated |
| Summary CSV empty cells are not self-describing. | `reasoning.md#trade-off-and-debt-analysis` | Detailed warning output includes source, dimension, path, and reason. | mitigated |
| Requested snippets may contain sensitive source content. | `reasoning.md#assumptions-and-risks` | Output is limited to explicitly requested lines under parsed source roots; safe-output review found no prompt/env/runtime payload rendering. | mitigated |
| No area-specific reasoning files exist despite being selected in sprint index and listed in requirements. | `reasoning.md#area-specific-reasoning-inputs` | Consolidated `reasoning.md` was used as the decision source. | closed |
| `sprint-index.md` previously named `system/protocols/sprint-review-protocol.md`, but the canonical file is `system/protocols/review-sprint-protocol.md`. | Input loading for this plan | `sprint-index.md` now references the canonical review protocol path. | closed |

## Review Inputs

Review should use:

- `projects/ultraplan-go/sprints/13-summary-and-code-extraction/requirements.md`
- `projects/ultraplan-go/sprints/13-summary-and-code-extraction/sprint-index.md`
- `projects/ultraplan-go/sprints/13-summary-and-code-extraction/technical-handbook.md`
- `projects/ultraplan-go/sprints/13-summary-and-code-extraction/reasoning.md`
- `projects/ultraplan-go/sprints/13-summary-and-code-extraction/plan.md`
- `projects/ultraplan-go/sprints/13-summary-and-code-extraction/reasoning/*.md` if created later, though none were present during planning
- `system/protocols/architecture-review-protocol.md`
- `system/protocols/review-sprint-protocol.md`
- implementation diff under `/home/antonioborgerees/coding/ultraplan-go`
- verification evidence from `go test ./...` and `go build ./cmd/ultraplan`

## Execution Log

| Date / Step | Action | Evidence / Notes |
| --- | --- | --- |
| 2026-06-02 / planning | Loaded required planning inputs and created sprint implementation plan. | `requirements.md`, `reasoning.md`, `sprint-index.md`, `technical-handbook.md`, PRD, TRD, architecture doc, project index, plan template, and available review protocols were used. No area reasoning files were present. |
| 2026-06-12 / implementation | Implemented standalone summary command, hardened summary ordering/warnings, added code extraction parser/resolver/service/renderers, wired top-level code command, and added fixture-backed tests. | Changed implementation under `/home/antonioborgerees/coding/ultraplan-go`: `internal/study`, `internal/codeextract`, and `internal/app`. |
| 2026-06-12 / verification | Ran required verification. | `go test ./...` passed. `go build ./cmd/ultraplan` passed. Boundary review found no new forbidden global packages and no `internal/codeextract` runtime dependency. |

## Completion Criteria

- [x] All tasks are complete or explicitly deferred with requirement-grounded reasons.
- [x] Verification commands were run or deferrals are documented with blocker evidence.
- [x] Evidence satisfies the expectations from `reasoning.md` and this plan.
- [x] `review.md` can evaluate conformance without guessing intent.
- [x] Review records any deviations, residual risks, stale protocol path issue, and follow-ups.
