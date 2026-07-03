# Sprint Review: 06 Study Run Execution

> Project: `ultraplan-go`
> Sprint: `06-study-run-execution`
> Plan: `plan.md`
> Reasoning: `reasoning.md`
> Reviewed: 2026-05-31

This review records implementation and verification evidence for the plan-driven sprint execution. It does not add runtime execution, prompt composition, run-loop behavior, summary generation, code extraction, target workflows, or public stable validation output.

## Review Inputs

- `requirements.md`
- `sprint-index.md`
- `technical-handbook.md`
- `reasoning.md`
- `plan.md`
- Architecture Review protocol from `system/protocols/architecture-review-protocol.md`
- Sprint Review protocol from `system/protocols/sprint-review-protocol.md`
- Implementation diff in `/home/antonioborgerees/coding/ultraplan-go/internal/study/`

## Implementation Evidence

| Requirement Area | Status | Evidence |
| --- | --- | --- |
| Study-owned report validation domain | pass | `internal/study/domain.go` defines report kind, validation status, severity, checks, results, rating states, and `Dimension.DisableCodeCitations`. |
| Deterministic report paths | pass | `internal/study/reports.go` defines source and final report paths under `reports/source/` and `reports/final/`. |
| Strict rating parsing | pass | `internal/study/rating.go` accepts `**8 / 10**`, `8/10`, and `Rating: 8`; distinguishes valid, missing, invalid, and ambiguous states without invented scores. |
| Per-source validation | pass | `internal/study/validation.go` validates expected report file, non-empty content, heading, source info, summary, rating, question/answer content, rating parseability, and citation shape when required. |
| Source-kind rules | pass | Directory sources require citation shape by default; `Dimension.DisableCodeCitations` disables that requirement explicitly; Markdown document sources skip code citations by default while still requiring report structure and ratings. |
| Final report validation | pass | `internal/study/validation.go` checks final report file, non-empty content, study context, sources studied table, executive summary, rating summary, synthesis/pattern content, and open questions/notable absences. |
| Deterministic diagnostics | pass | `ValidationCheck` exposes check name, status, severity, path, expected, observed, source kind, guidance, and preserved internal cause. File-read diagnostics use sanitized observed values while preserving wrapped errors for `errors.Is`. |

## Plan Execution

| Plan Task | Status | Evidence |
| --- | --- | --- |
| Task 1: Confirm current study domain shape | done | Existing `Study`, `Source`, `SourceKind`, `Dimension`, discovery, applicability, and resolution behavior preserved with additive domain fields only. |
| Task 2: Add study-owned report and validation domain types | done | `internal/study/domain.go`. |
| Task 3: Add deterministic report path helpers | done | `internal/study/reports.go`. |
| Task 4: Implement strict rating parser | done | `internal/study/rating.go`. |
| Task 5: Test rating parser first | done | `internal/study/rating_test.go`. |
| Task 6: Implement per-source report validation | done | `internal/study/validation.go`. |
| Task 7: Implement final report validation | done | `internal/study/validation.go`. |
| Task 8: Test report validators | done | `internal/study/validation_test.go`. |
| Task 9: Scope and architecture review before final test run | done | All new implementation is under `internal/study`; no global validation, reports, or parsing package was introduced. |
| Task 10: Run verification and prepare review evidence | done | `go test ./internal/study -v` and `go test ./...` passed. |

## Verification Evidence

| Check | Command / Method | Result | Notes |
| --- | --- | --- | --- |
| Full Go test suite | `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go` | pass | All packages pass. |
| Focused study tests | `go test ./internal/study -v` | pass | 25 study tests pass, including rating and report validation coverage. |
| Architecture: file placement | Inspect changed paths | pass | New files are `internal/study/reports.go`, `internal/study/rating.go`, `internal/study/validation.go`, and tests. |
| Architecture: no global packages | Inspect `internal/` and changed paths | pass | No `internal/validation`, `internal/reports`, or `internal/parsing` package created. |
| Scope: no runtime execution | Search changed implementation | pass | No agentwrap/OpenCode invocation, prompt composition, run-loop, worker pool, repair, retry, fallback, status, summary, or code extraction behavior added. |
| Performance: bounded validation | Inspect validator implementation | pass | Validators read expected report files only and validate citation shape with a regex; no repository traversal or citation resolution. |
| Security: safe diagnostics | Inspect `fileCheck` and rating diagnostics | pass | File-read observed values are sanitized; rating diagnostics report classifications/reasons rather than raw report content for non-valid states. |

## Acceptance Criteria Coverage

| Acceptance Criteria | Status | Evidence |
| --- | --- | --- |
| AC-01, AC-06 missing/empty/read failures include report path | pass | `validation_test.go` covers missing source report and empty source/final report; result path is deterministic. |
| AC-02 per-source required structure and rating | pass | `validation_test.go` covers valid reports, missing summary/rating/questions, missing rating parse, and ambiguous rating. |
| AC-03 directory citation shape unless disabled | pass | `validation_test.go` covers missing directory citation failure and `DisableCodeCitations` skip behavior. |
| AC-04 Markdown citation exemption and source kind diagnostics | pass | `validation_test.go` covers Markdown source citation skip with `SourceKindMarkdown`. |
| AC-05 Markdown reports still require rating, summary, and answers | pass | Markdown source validation still runs the same required-section and rating checks before skipping citations. |
| AC-07 final report required sections | pass | `validation_test.go` covers valid final report and missing context/table/rating/synthesis/open-question failures. |
| AC-08, AC-09, AC-10 rating parser behavior | pass | `rating_test.go` covers accepted formats, missing, invalid, ambiguous, and no invented score. |
| AC-11 structured validation results | pass | `ValidationResult` and `ValidationCheck` expose status, check names, severity, path, expected, observed, source kind, and guidance. |
| AC-12 preserved causes | pass | Read failures wrap causes; tests assert `errors.Is(res.Err, os.ErrNotExist)` and `errors.Is(res.Err, errValidationRead)`. |
| AC-13 runtime-free reusable APIs | pass | Public functions take study domain values and return study validation results; no runtime dependencies. |
| AC-14 skipped inapplicable Markdown pairs | pass | Applicability discovery remains unchanged; validators do not discover missing report work. |
| AC-15, AC-16, AC-17 tests and reviewability | pass | Focused tests and `go test ./...` pass; this review records evidence. |

## Review Protocol Results

| Protocol | Result | Evidence |
| --- | --- | --- |
| Architecture Review | pass | Implementation stays inside `internal/study`, preserves dependency direction, and introduces no global validation/report/parsing package. |
| Sprint Review | pass | Requirements, non-goals, selected contracts, tests, and verification evidence are satisfied. |

## Residual Risks

| Risk | Severity | Status |
| --- | --- | --- |
| Final report and per-source section matching is structural and intentionally heuristic. | low | Accepted for this sprint; future report template/versioning can tighten matching. |
| Validators read full expected report files into memory. | low | Accepted because validation is bounded to known report artifacts and no repository scan is performed. Add a size cap if batch validation later processes untrusted or very large reports. |
| `RatingResult.Raw` preserves raw text for valid ratings. | low | Internal-only today; clear or sanitize before exposing rating results through public JSON/logging. |

## Final Assessment

- **Status:** accepted.
- **Blocking Issues:** none.
- **Verification:** `go test ./internal/study -v` and `go test ./...` passed on 2026-05-31.
