# Sprint Plan: 06 Study Run Execution

> Project: `ultraplan-go`
> Sprint: `06-study-run-execution`
> Source: `projects/ultraplan-go/sprints/06-study-run-execution/reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/sprints/06-study-run-execution/requirements.md`, `projects/ultraplan-go/sprints/06-study-run-execution/reasoning.md`, `projects/ultraplan-go/sprints/06-study-run-execution/sprint-index.md`, `projects/ultraplan-go/sprints/06-study-run-execution/technical-handbook.md`, `projects/ultraplan-go/sprints/06-study-run-execution/reasoning/*.md` (not present), `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/project-index.md`, `templates/sprint-plan.md`, `/home/antonioborgerees/coding/ultraplan-go/internal/study/domain.go`, `/home/antonioborgerees/coding/ultraplan-go/internal/study/discovery.go`, `/home/antonioborgerees/coding/ultraplan-go/internal/study/applicability.go`, `/home/antonioborgerees/coding/ultraplan-go/internal/study/resolve.go`, `/home/antonioborgerees/coding/ultraplan-go/internal/study/study_test.go`

This plan executes `projects/ultraplan-go/sprints/06-study-run-execution/reasoning.md`. It must not invent architecture, scope, or decisions.

## Reasoning Source

- **Sprint Reasoning:** `projects/ultraplan-go/sprints/06-study-run-execution/reasoning.md`
- **Sprint Index:** `projects/ultraplan-go/sprints/06-study-run-execution/sprint-index.md`
- **Technical Handbook:** `projects/ultraplan-go/sprints/06-study-run-execution/technical-handbook.md`
- **Area Reasoning:** none present; `projects/ultraplan-go/sprints/06-study-run-execution/reasoning/` does not exist at planning time

## Sprint Status

- **Status:** not started
- **Owner:** implementation agent
- **Start Date:** pending
- **Completion Date:** pending

## Decisions To Execute

| Decision | Source Section | Execution Implication |
| --- | --- | --- |
| Keep validation, report paths, ratings, and diagnostics study-owned | `reasoning.md#decision-1-keep-validation-report-paths-ratings-and-diagnostics-study-owned` | Add domain concepts in `internal/study/domain.go`, deterministic helpers in `internal/study/reports.go`, validators in `internal/study/validation.go`, and parser in `internal/study/rating.go`; do not create global `validation`, `reports`, or `parsing` packages. |
| Validate per-source and final reports as local artifacts with source-kind-aware rules | `reasoning.md#decision-2-validate-per-source-and-final-reports-as-local-artifacts-with-source-kind-aware-rules` | Validate only expected report paths, require required sections and rating, enforce directory citation shape by default, exempt Markdown document sources from code citation requirements by default, and require final report context/table/summary/synthesis/open-question sections. |
| Return structured validation results with deterministic actionable diagnostics | `reasoning.md#decision-3-return-structured-validation-results-with-deterministic-actionable-diagnostics` | Expose validation status, checks, severity, path, expected, observed, source kind where relevant, and safe guidance; wrap filesystem and parsing causes. |
| Parse ratings strictly and preserve uncertainty | `reasoning.md#decision-4-parse-ratings-strictly-and-preserve-uncertainty` | Accept only required rating formats for this sprint and distinguish valid, missing, invalid, and ambiguous ratings without inventing scores. |
| Keep validation bounded, deterministic, and runtime-free | `reasoning.md#decision-5-keep-validation-bounded-deterministic-and-runtime-free` | Do not add runtime execution, prompt composition, agentwrap/OpenCode wiring, run-loop behavior, summary writing, code extraction, repository scanning, or public stable validation JSON. |

## Requirements / Contracts To Satisfy

| Contract / Requirement ID | Required Behavior | Evidence Planned |
| --- | --- | --- |
| Architecture | Study-owned validation and report behavior remains under `internal/study`; platform packages do not import `study`; no global report/validation/parsing package is introduced. | File placement review, import review, and `go test ./...`. |
| Errors | Filesystem and parsing errors preserve causes while diagnostics remain actionable. | Unit tests using `errors.Is` or `errors.As` where practical and assertions on validation check fields. |
| Security | Diagnostics include safe paths and concise observed facts, not raw large report content, secrets, runtime payloads, or environment values. | Test/review checks for diagnostic fields and absence of raw content dumps. |
| Testing | Behavior is deterministic and fixture/table-testable without OpenCode, agentwrap runtime execution, network, provider credentials, or generated model output. | `internal/study/validation_test.go`, `internal/study/rating_test.go`, and `go test ./...`. |
| Performance | Validation reads expected report artifacts only and validates citation shape without resolving files or scanning repositories. | Review confirms no repository traversal or citation resolution is added. |
| AC-01 | Missing, empty, and unreadable per-source report failures include report path. | Per-source validation tests with missing and empty files; unreadable-file coverage via a portable read-failure seam if natural, otherwise documented review note. |
| AC-02 | Per-source report validation requires heading, source information, summary, rating section, parseable rating, and required question/answer content. | Per-source valid/missing-section/missing-answer tests. |
| AC-03 | Directory source validation enforces file-path-and-line citation shape unless code citations are explicitly disabled for the dimension. | Directory citation failure and disabled-citation tests. |
| AC-04 | Markdown document source validation does not require code citations by default and records source kind in diagnostics. | Markdown source validation tests assert pass/fail behavior and diagnostic `source_kind`. |
| AC-05 | Markdown document reports still require rating, summary, and answers. | Markdown source missing rating/summary/answers tests. |
| AC-06 | Missing, empty, and unreadable final report failures include report path. | Final report validation tests with missing and empty files; unreadable-file coverage handled as for AC-01. |
| AC-07 | Final reports require study context, sources studied table, executive summary, rating summary, synthesis/pattern content, and open questions/notable absences. | Final report valid and missing-section tests. |
| AC-08 | Rating parser accepts `**8 / 10**`, `8/10`, and `Rating: 8`. | Rating parser table tests. |
| AC-09 | Rating parser rejects or warns on ambiguous ratings without inventing a score. | Ambiguous rating tests assert no score is returned. |
| AC-10 | Rating diagnostics distinguish missing from invalid or ambiguous ratings. | Rating parser tests and validator tests for diagnostic classification. |
| AC-11 | Validation results expose pass/fail status, check names, severity, path, expected, observed, and guidance where safe. | Validation result unit assertions. |
| AC-12 | Validation preserves error cause chains for filesystem and parsing failures. | Wrapped error assertions. |
| AC-13 | Validation APIs are usable by later run, batch, synthesis, summary, and validate commands without runtime execution. | API review confirms plain study-domain inputs and no runtime dependencies. |
| AC-14 | Inapplicable Markdown source/dimension pairs remain skipped or excluded by callers. | Tests or review confirm validators do not discover skipped pairs and `GetApplicableSources` behavior remains intact. |
| AC-15 | Unit tests cover valid reports, missing files, empty files, missing sections, missing ratings, ambiguous ratings, directory citation failures, Markdown citation exemptions, and final report failures. | `internal/study/validation_test.go` and `internal/study/rating_test.go`. |
| AC-16 | `go test ./...` passes in `/home/antonioborgerees/coding/ultraplan-go`. | Recorded command output in sprint review. |
| AC-17 | Full sprint verification is reviewable. | `review.md` can cite test output and architecture/scope review evidence. |

## Tasks

- [ ] **Task 1: Confirm Current Study Domain Shape**
  > Executes: Decision 1, AC-13, Architecture
  - [ ] Inspect `internal/study/domain.go`, `discovery.go`, `applicability.go`, and `resolve.go` before editing.
  - [ ] Identify the smallest domain additions needed for report kind, validation status, validation result, validation check, severity, and rating concepts.
  - [ ] Preserve existing `SourceKind`, `Source`, `Dimension`, `Study`, and `GetApplicableSources` behavior unless a requirement directly needs an additive change.

- [ ] **Task 2: Add Study-Owned Report And Validation Domain Types**
  > Executes: Decisions 1 and 3, AC-11, AC-13
  - [ ] Extend `internal/study/domain.go` with report, validation, diagnostic, and rating result concepts needed by validators.
  - [ ] Keep result structures internal-domain-friendly and deterministic; do not create public JSON schema promises unless versioned.
  - [ ] Include fields for status, check name, severity, path, expected, observed, source kind where relevant, corrective guidance, and underlying cause where needed.
  - [ ] Ensure any error wrapping supports `errors.Is` or `errors.As` for filesystem and parser failures.

- [ ] **Task 3: Add Deterministic Report Path Helpers**
  > Executes: Decision 1, AC-01, AC-06, AC-13
  - [ ] Create `internal/study/reports.go` for expected per-source and final report paths under `studies/<study>/reports/source/` and `studies/<study>/reports/final/`.
  - [ ] Use existing study, source, and dimension identity values rather than inventing a separate report registry.
  - [ ] Add report metadata construction helpers only where validators or later callers need them.
  - [ ] Keep paths deterministic and local; do not scan report directories to discover work in this sprint.

- [ ] **Task 4: Implement Strict Rating Parser**
  > Executes: Decision 4, AC-08, AC-09, AC-10
  - [ ] Create `internal/study/rating.go` with parser support for `**8 / 10**`, `8/10`, and `Rating: 8`.
  - [ ] Return explicit valid, missing, invalid, and ambiguous states.
  - [ ] Normalize valid scores without coupling to summary CSV writing.
  - [ ] Do not coerce missing, invalid, or ambiguous ratings into `0` or any invented score.

- [ ] **Task 5: Test Rating Parser First**
  > Executes: Decision 4, AC-08, AC-09, AC-10, AC-15
  - [ ] Create `internal/study/rating_test.go` with table-driven tests for accepted formats.
  - [ ] Cover missing ratings, invalid values, multiple conflicting ratings, repeated identical ratings if supported, and no-score behavior for ambiguous cases.
  - [ ] Assert classification separately from human-readable error text.

- [ ] **Task 6: Implement Per-Source Report Validation**
  > Executes: Decisions 2, 3, and 5, AC-01 through AC-05, AC-11 through AC-14
  - [ ] Create `internal/study/validation.go` with a per-source validator that reads the expected report file and returns structured validation results.
  - [ ] Check file existence/readability, non-empty content, report heading, source information section, summary section, rating section, parseable rating, and required question/answer content.
  - [ ] For directory sources, require file-path-and-line citation shape by default unless the dimension explicitly disables code citation requirements.
  - [ ] For Markdown document sources, do not require code citations by default; still require rating, summary, and answers and record `SourceKindMarkdown` in diagnostics.
  - [ ] Validate citation shape only; do not resolve cited files, extract snippets, or traverse source repositories.
  - [ ] Treat inapplicable Markdown source/dimension pairs as caller-skipped and do not add discovery logic that turns them into missing-report failures.

- [ ] **Task 7: Implement Final Report Validation**
  > Executes: Decisions 2, 3, and 5, AC-06, AC-07, AC-11 through AC-13
  - [ ] Add a final report validator in `internal/study/validation.go` using deterministic expected final report paths.
  - [ ] Check file existence/readability, non-empty content, study parameters or equivalent context, sources studied table, executive summary, rating summary, pattern or synthesis content, and open questions or notable absences.
  - [ ] Keep final report validation structural; do not attempt semantic quality scoring or generated-output repair.

- [ ] **Task 8: Test Report Validators**
  > Executes: Decisions 2, 3, and 5, AC-01 through AC-07, AC-11 through AC-15
  - [ ] Create `internal/study/validation_test.go` with local temp-file fixtures or inline Markdown fixtures, following existing `study_test.go` helper style where practical.
  - [ ] Cover valid per-source reports, missing files, empty files, missing heading/source/summary/rating/answers, missing ratings, ambiguous ratings, directory citation failures, directory citation disabled, Markdown citation exemption, and Markdown required content failures.
  - [ ] Cover valid final reports, missing files, empty files, missing study context, missing sources table, missing executive summary, missing rating summary, missing synthesis content, and missing open questions/notable absences.
  - [ ] Assert deterministic diagnostic fields instead of parsing diagnostic strings.
  - [ ] Add unreadable-file coverage only if it is portable in this codebase; otherwise document why missing/read-error coverage is the portable substitute in `review.md`.

- [ ] **Task 9: Scope And Architecture Review Before Final Test Run**
  > Executes: Decisions 1 and 5, Architecture, Security, Performance, AC-13, AC-14
  - [ ] Confirm all new implementation files are under `/home/antonioborgerees/coding/ultraplan-go/internal/study`.
  - [ ] Confirm no new global `internal/validation`, `internal/reports`, or `internal/parsing` package exists.
  - [ ] Confirm no runtime execution, prompt composition, agentwrap/OpenCode wiring, run-loop, worker pool, summary generation, code extraction, target, or sprint workflow behavior was added.
  - [ ] Confirm diagnostics avoid raw large report dumps, secrets, runtime payloads, environment values, and repository scans.

- [ ] **Task 10: Run Verification And Prepare Review Evidence**
  > Executes: AC-15, AC-16, AC-17, Testing
  - [ ] Run `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go`.
  - [ ] If tests fail, fix sprint-owned failures without broadening scope.
  - [ ] Record test command, result, architecture review notes, scope review notes, risks, and any deviations in `review.md`.

## Evidence Checklist

- [ ] Rating parser tests prove accepted formats, missing classification, invalid classification, ambiguous classification, and no invented score.
- [ ] Per-source validation tests prove missing, empty, malformed, valid, rating, citation, Markdown exemption, and diagnostics behavior.
- [ ] Final report validation tests prove missing, empty, malformed, and valid final report behavior.
- [ ] Validation diagnostics expose deterministic fields required by AC-11 and preserve cause chains required by AC-12.
- [ ] Architecture evidence confirms study-owned files and dependency direction.
- [ ] Scope evidence confirms no runtime execution, prompt composition, run-loop, summary generation, code extraction, target, sprint execution, or public stable JSON validation scope was added.
- [ ] Security evidence confirms diagnostics are safe and bounded.
- [ ] Performance evidence confirms validation reads expected report artifacts only and does not scan source repositories.
- [ ] Deviations from `reasoning.md` are recorded before implementation continues.
- [ ] Required review protocols have evidence in `review.md`.

## Verification Commands

| Check | Command | Expected Result |
| --- | --- | --- |
| Full Go test suite | `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go` | Passes without OpenCode, network, provider credentials, or generated model output. |
| Optional focused study tests during implementation | `go test ./internal/study` from `/home/antonioborgerees/coding/ultraplan-go` | Passes and covers rating/report validation behavior. |
| Architecture file placement review | Inspect changed paths | New report validation, rating, and report helper code is under `internal/study`. |
| Scope review | Inspect changed paths and diffs | No runtime execution, prompt composition, agentwrap/OpenCode wiring, run-loop, summary, code extraction, target, or sprint workflow behavior added. |

## Risks And Blockers

| Risk / Blocker | Source | Mitigation | Status |
| --- | --- | --- | --- |
| Existing study domain may lack enough dimension metadata to decide citation-disabled behavior. | `reasoning.md#assumptions-and-risks` | Use the narrowest existing signal available; if no safe signal exists, implement default directory citation requirement and record an open question for richer dimension schema. | open |
| Heading matching may be too strict or too loose for generated Markdown. | `reasoning.md#assumptions-and-risks` | Define explicit accepted variants in tests and keep matching deterministic. | open |
| Portable unreadable-file tests may be OS-dependent. | `reasoning.md#assumptions-and-risks` | Prefer a natural small read-failure seam only if it keeps implementation minimal; otherwise cover missing/empty/read errors and document residual portability. | open |
| Validation structures could accidentally become stable public JSON. | `reasoning.md#assumptions-and-risks` | Keep schema-like structures internal/test-only unless a later public validate command versions them. | open |
| Area-specific reasoning file required by requirements is absent. | `reasoning.md#area-specific-reasoning-inputs` | Consolidated `reasoning.md` records final decisions directly; implementation should not wait for another reasoning artifact unless the user requests it. | mitigated |
| Inapplicable Markdown pairs are caller-skipped. | `requirements.md#acceptance-criteria`, `reasoning.md#implementation-constraints` | Preserve existing `GetApplicableSources` behavior and avoid validators that discover missing reports for skipped pairs. | open |

## Review Inputs

Review should use:

- `projects/ultraplan-go/sprints/06-study-run-execution/requirements.md`
- `projects/ultraplan-go/sprints/06-study-run-execution/sprint-index.md`
- `projects/ultraplan-go/sprints/06-study-run-execution/technical-handbook.md`
- `projects/ultraplan-go/sprints/06-study-run-execution/reasoning/*.md` where created; none were present at planning time
- `projects/ultraplan-go/sprints/06-study-run-execution/reasoning.md`
- `projects/ultraplan-go/sprints/06-study-run-execution/plan.md`
- Implementation diff in `/home/antonioborgerees/coding/ultraplan-go`
- Verification evidence from `go test ./...`
- Architecture Review protocol from `system/protocols/architecture-review-protocol.md`
- Sprint Review protocol from `system/protocols/sprint-review-protocol.md`

## Execution Log

| Date / Step | Action | Evidence / Notes |
| --- | --- | --- |
| 2026-05-31 / planning | Created implementation plan from sprint reasoning and selected evidence. | No implementation code changed by this planning step. |

## Completion Criteria

- [ ] All tasks are complete or explicitly deferred in `review.md`.
- [ ] Verification commands were run or deferrals are documented with reason.
- [ ] Evidence satisfies the expectations from `reasoning.md`.
- [ ] `review.md` can evaluate conformance without guessing intent.
- [ ] No architecture, scope, or validation behavior deviates from `reasoning.md` without an explicit recorded decision.
