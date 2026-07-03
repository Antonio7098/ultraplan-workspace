# Sprint Requirements: 06 Study Run Execution

> Project: `ultraplan-go`
> Sprint: `06-study-run-execution`
> Purpose: the authoritative, human-readable sprint contract. All other sprint artifacts must satisfy these requirements.

## Sprint Goal

Make study reports first-class validated artifacts by implementing per-source report validation, final report validation, rating parsing, and actionable validation diagnostics without introducing runtime execution.

## Required Outputs

| Output | Path | Description |
| --- | --- | --- |
| Sprint requirements | `projects/ultraplan-go/sprints/06-study-run-execution/requirements.md` | Defines the binding goal, scope, acceptance criteria, constraints, dependencies, and review expectations for this sprint. |
| Sprint index | `projects/ultraplan-go/sprints/06-study-run-execution/sprint-index.md` | Selects the roadmap entries, source docs, evidence reports, contracts, and review protocols that govern this sprint. |
| Technical handbook | `projects/ultraplan-go/sprints/06-study-run-execution/technical-handbook.md` | Gives implementation guidance for study-owned report validation, rating parsing, diagnostics, Markdown source rules, and tests. |
| Report validation reasoning | `projects/ultraplan-go/sprints/06-study-run-execution/reasoning/report-validation.md` | Records design reasoning for validator ownership, required sections, rating ambiguity, citation requirements, and diagnostic shape. |
| Consolidated reasoning | `projects/ultraplan-go/sprints/06-study-run-execution/reasoning.md` | Summarizes accepted sprint decisions, tradeoffs, risks, and requirement mappings. |
| Implementation plan | `projects/ultraplan-go/sprints/06-study-run-execution/plan.md` | Provides the ordered task plan and verification checklist for implementation. |
| Sprint review | `projects/ultraplan-go/sprints/06-study-run-execution/review.md` | Reviews completed implementation against these requirements and records carry-forward decisions. |
| Study report domain updates | `/home/antonioborgerees/coding/ultraplan-go/internal/study/domain.go` | Defines or extends report, validation status, validation result, validation check, and rating concepts needed by study report validation. |
| Study report path helpers | `/home/antonioborgerees/coding/ultraplan-go/internal/study/reports.go` | Provides deterministic per-source and final report path helpers and report metadata construction for validation. |
| Study report validation implementation | `/home/antonioborgerees/coding/ultraplan-go/internal/study/validation.go` | Implements per-source and final report validators with source-kind-aware citation rules and actionable diagnostics. |
| Rating parser implementation | `/home/antonioborgerees/coding/ultraplan-go/internal/study/rating.go` | Parses supported report rating formats, returns missing or ambiguous rating diagnostics, and does not invent scores. |
| Study report validation tests | `/home/antonioborgerees/coding/ultraplan-go/internal/study/validation_test.go` | Tests per-source validation, final report validation, Markdown-source citation behavior, missing sections, and diagnostics. |
| Rating parser tests | `/home/antonioborgerees/coding/ultraplan-go/internal/study/rating_test.go` | Tests supported rating formats, missing ratings, ambiguous ratings, invalid values, and score normalization. |

## Acceptance Criteria

- [ ] Per-source report validation fails when the expected report file is missing, empty, or unreadable, and diagnostics include the report path.
- [ ] Per-source report validation requires a heading, source information section, summary section, rating section, parseable rating, and required question or answer content.
- [ ] Directory source report validation enforces file-path-and-line citation shape unless the selected dimension explicitly disables code citation requirements.
- [ ] Markdown document source report validation does not require code citations by default and records the source kind in diagnostics.
- [ ] Markdown document source report validation still requires a rating, summary, and answers to dimension questions.
- [ ] Final report validation fails when the expected final report file is missing, empty, or unreadable, and diagnostics include the report path.
- [ ] Final report validation requires study parameters or equivalent study context, sources studied table, executive summary, rating summary, pattern or synthesis content, and open questions or notable absences.
- [ ] Rating parser accepts `**8 / 10**`, `8/10`, and `Rating: 8` formats.
- [ ] Rating parser rejects or warns on ambiguous ratings without inventing a score.
- [ ] Rating parser distinguishes missing ratings from invalid or ambiguous ratings in validation diagnostics.
- [ ] Validation results expose pass/fail status, check names, severity, path, expected value, observed value, and corrective guidance where safe.
- [ ] Validation preserves error cause chains for filesystem and parsing failures.
- [ ] Validation APIs are usable by later single-run, batch, synthesis, summary, and validate commands without depending on runtime execution.
- [ ] Inapplicable Markdown source/dimension pairs from Sprint 5 remain skipped or excluded by callers; this sprint must not turn them into validation failures.
- [ ] Unit tests cover valid reports, missing files, empty files, missing sections, missing ratings, ambiguous ratings, directory citation failures, Markdown-source citation exemptions, and final report failures.
- [ ] `go test ./...` passes in `/home/antonioborgerees/coding/ultraplan-go`.

## Non-Goals

- Runtime execution for `ultraplan study <study> run <dimension-ref> <source-ref>` is not included despite the sprint folder name.
- Prompt composition for directory analysis, Markdown document analysis, or synthesis is not included.
- Agentwrap/OpenCode integration, runtime health checks, policy runner wiring, retry, fallback, and event mapping are not included.
- `run-all`, `run-loop`, `status`, cancellation, durable task state, and bounded worker pools are not included.
- Summary generation and `summary.csv` writing are not included.
- Code reference extraction is not included.
- A public `ultraplan study <study> validate` command and stable JSON validation output are not included unless already present from prior work.
- Target workflows, sprint planning, and sprint execution remain out of scope.

## Constraints

- The authoritative roadmap scope for Sprint 6 is report validation and rating parsing; the directory slug must not expand this sprint into runtime execution.
- Keep study report validation behavior in `internal/study`; do not create global `validation`, `reports`, or `parsing` packages for study-owned behavior.
- Preserve the dependency direction from `docs/ARCHITECTURE.md`: `study` may use `workspace` and platform helpers, but platform packages must not import `study`.
- Apply the roadmap Skeleton / Local CLI Gate only; Runtime / Provider and Batch / Durable Workflow gates do not apply because runtime execution and durable run-loop behavior are non-goals.
- Runtime success must not be treated as product success in the API design: validation helpers must be able to prove expected artifacts exist and satisfy report rules.
- Validation diagnostics must be deterministic and suitable for unit tests.
- Validators must not require real OpenCode, `agentwrap`, network access, provider credentials, or generated model output beyond local fixtures.
- Markdown document source reports must not require code citations by default, matching Sprint 5 source-kind and applicability behavior.
- Directory source citation validation must check citation shape only; resolving cited files and extracting snippets is deferred to the code extraction sprint.
- Rating parsing must not coerce ambiguous or missing values into `0`.
- Generated or persisted machine-readable schemas introduced in this sprint must include a clear version or be kept internal/test-only.

## Dependencies

| Prior Sprint / Output | Required For | Notes |
| --- | --- | --- |
| Sprint 1: Go Module, CLI Shell, and App Composition | Buildable CLI and package structure | Required for package placement, test execution, and `cmd/ultraplan` build continuity. |
| Sprint 2: Workspace, Config, Logging, and Health Skeleton | Workspace path resolution and filesystem safety | Required for report paths and diagnostics that remain workspace-aware. |
| Sprint 3: Study Domain, Listing, and Resolution | Study, source, dimension, and lookup model | Required so validators can reason about dimensions, sources, source names, and deterministic study metadata. |
| Sprint 4: Study Initialization From YAML | Generated study directory layout | Required for expected `studies/<study>/reports/source/` and `studies/<study>/reports/final/` layouts used by validation fixtures. |
| Sprint 5: Markdown Document Sources and Applicability | Source kind and applicability behavior | Required so Markdown report validation can avoid default code citation requirements and inapplicable pairs remain skipped. |
| `projects/ultraplan-go/docs/PRD.md` | Product report validation behavior | Governs per-source success, synthesis success, Markdown document source validation, and runtime-success-is-not-product-success principles. |
| `projects/ultraplan-go/docs/TRD.md` | Technical validation and rating requirements | Governs report validator checks, validation result mapping, rating parser formats, and diagnostics. |
| `projects/ultraplan-go/docs/ARCHITECTURE.md` | Package ownership and dependency direction | Governs keeping validation and report behavior in `internal/study`. |
| `projects/ultraplan-go/sprints/05-study-listing-inspection/review.md` | Carry-forward decisions | Confirms Sprint 5 source-kind/applicability implementation is complete and runtime/prompt/report validation remained deferred. |

## Review Expectations

| What | How Verified |
| --- | --- |
| Scope alignment | Review `requirements.md`, `sprint-index.md`, `technical-handbook.md`, `reasoning.md`, and `plan.md` against roadmap Sprint 6 and confirm no runtime, prompt composition, run-loop, summary, code extraction, target, or sprint-execution scope was added. |
| Per-source validation behavior | Run or inspect tests proving missing, empty, malformed, and valid per-source reports produce deterministic validation results and actionable diagnostics. |
| Final report validation behavior | Run or inspect tests proving missing, empty, malformed, and valid final reports produce deterministic validation results and actionable diagnostics. |
| Rating parser behavior | Run or inspect tests for `**8 / 10**`, `8/10`, `Rating: 8`, missing ratings, ambiguous ratings, invalid ratings, and non-invented scores. |
| Source-kind citation rules | Run or inspect tests proving directory reports require citation shape by default while Markdown document reports do not require code citations by default. |
| Architecture conformance | Review implementation paths and imports to confirm study-owned validation remains in `internal/study`, platform packages do not import product modules, and no global validation/report package was introduced. |
| Verification commands | Run `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go` and record the result in `review.md`. |
