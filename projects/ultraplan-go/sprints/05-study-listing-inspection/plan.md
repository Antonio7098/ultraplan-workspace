# Sprint Plan: 05 Study Listing Inspection

> Project: `ultraplan-go`
> Sprint: `05-study-listing-inspection`
> Source: `projects/ultraplan-go/sprints/05-study-listing-inspection/reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/05-study-listing-inspection/requirements.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/sprints/05-study-listing-inspection/sprint-index.md`, `projects/ultraplan-go/sprints/05-study-listing-inspection/technical-handbook.md`, `projects/ultraplan-go/sprints/05-study-listing-inspection/reasoning.md`, `templates/sprint-plan.md`, `/home/antonioborgerees/coding/ultraplan-go/internal/study/domain.go`, `/home/antonioborgerees/coding/ultraplan-go/internal/study/discovery.go`, `/home/antonioborgerees/coding/ultraplan-go/internal/study/resolve.go`, `/home/antonioborgerees/coding/ultraplan-go/internal/study/service.go`, `/home/antonioborgerees/coding/ultraplan-go/internal/study/study_test.go`, `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_commands.go`, `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_commands_test.go`, `/home/antonioborgerees/coding/ultraplan-go/go.mod`

This plan executes `projects/ultraplan-go/sprints/05-study-listing-inspection/reasoning.md`. It must not invent architecture, scope, or decisions.

## Reasoning Source

- **Sprint Reasoning:** `projects/ultraplan-go/sprints/05-study-listing-inspection/reasoning.md`
- **Sprint Index:** `projects/ultraplan-go/sprints/05-study-listing-inspection/sprint-index.md`
- **Technical Handbook:** `projects/ultraplan-go/sprints/05-study-listing-inspection/technical-handbook.md`
- **Area Reasoning:** none present; the selected `reasoning/source-applicability.md` artifact is absent and `reasoning.md` records that absence.

## Sprint Status

- **Status:** complete
- **Owner:** implementation agent
- **Start Date:** 2026-05-31
- **Completion Date:** 2026-05-31

## Decisions To Execute

| Decision | Source Section | Execution Implication |
| --- | --- | --- |
| Discover only direct child directory and Markdown sources | `reasoning.md#decision-1-discover-only-direct-child-directory-and-markdown-sources` | Extend `DiscoverSources` to include non-hidden top-level `.md` files, keep direct-child-only scanning, ignore nested Markdown, and preserve deterministic source ordering. |
| Model source kind, frontmatter, and normalized applicability in study domain | `reasoning.md#decision-2-model-source-kind-frontmatter-and-normalized-applicability-in-study-domain` | Extend `Source` with `ApplicableDimensions []string` and `Frontmatter map[string]any` or equivalent, preserving existing `Kind` values. |
| Parse and strip only leading YAML frontmatter | `reasoning.md#decision-3-parse-and-strip-only-leading-yaml-frontmatter` | Add `internal/study/markdown.go` with leading-only frontmatter parsing and stripping helpers, using existing `gopkg.in/yaml.v3`. |
| Normalize `applicable_dimensions` strictly and deduplicate deterministically | `reasoning.md#decision-4-normalize-applicable_dimensions-strictly-and-deduplicate-deterministically` | Normalize numeric/string YAML values to two-digit dimension strings, treat missing/empty as all dimensions, and fail invalid values with path/value context. |
| Provide a study-owned applicability helper | `reasoning.md#decision-5-provide-a-study-owned-applicability-helper` | Add `internal/study/applicability.go` with `GetApplicableSources` or equivalent exported helper that excludes inapplicable Markdown sources without treating them as failures. |
| Keep listing plain, deterministic, and capturable | `reasoning.md#decision-6-keep-listing-plain-deterministic-and-capturable` | Update list output so directory rows remain `name directory`, unfiltered Markdown rows render `name markdown all`, and filtered Markdown rows render `name markdown 01,03`. |
| Keep errors actionable and cause-preserving | `reasoning.md#decision-7-keep-errors-actionable-and-cause-preserving` | Wrap filesystem/YAML/metadata errors with source path context and preserve causes for CLI classification. |
| Verify through public behavior and local fixtures | `reasoning.md#decision-8-verify-through-public-behavior-and-local-fixtures` | Add or update `internal/study/source_test.go` or current study tests and `internal/app/study_commands_test.go`/`internal/study/cli_test.go` with temp-directory fixtures and output assertions. |

## Requirements / Contracts To Satisfy

| Contract / Requirement ID | Required Behavior | Evidence Planned |
| --- | --- | --- |
| `requirements.md` AC-01, AC-02 | `ultraplan study <study> list` discovers and displays directory and Markdown sources with source kind. | Command test for mixed sources and deterministic output. |
| `requirements.md` AC-03 | Hidden entries, non-directory non-Markdown files, and nested Markdown files are ignored. | Unit test for direct-child discovery fixture. |
| `requirements.md` AC-04 through AC-07 | Frontmatter parsing and stripping are leading-only; absent frontmatter preserves body. | Unit tests for parse/strip helpers. |
| `requirements.md` AC-08 through AC-10 | `applicable_dimensions` accepts numeric/string values, normalizes to two digits, defaults to all dimensions when missing/empty, and reports invalid values with path/value context. | Unit tests for normalization and invalid diagnostics. |
| `requirements.md` AC-11 through AC-14 | Applicability helper always includes directories, includes unfiltered Markdown, includes filtered Markdown only on normalized match, and excludes inapplicable pairs without failure. | Unit tests for helper behavior and ordering. |
| `requirements.md` AC-15 | Existing study listing remains deterministic and script-friendly. | Updated command tests pin text output. |
| `requirements.md` AC-16 | Tests cover valid/absent/malformed frontmatter, invalid applicability, ignored entries, nested Markdown, and mixed ordering. | Study and command test cases. |
| `requirements.md` AC-17 | Full Go test suite passes. | `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go`. |
| Architecture contract, `ARCHITECTURE.md` | Study behavior remains in `internal/study`; platform packages do not import `study`; no global validation/scheduler/reports/prompts packages. | Diff review and import review. |
| Errors contract | Errors preserve causes and include actionable source context. | Error-focused unit/command tests and review of `%w` wrapping. |
| Security contract | Source discovery uses direct child inspection only and treats Markdown metadata as untrusted input. | Fixture tests and scope review. |
| Testing contract | Behavior is covered with local filesystem fixtures and no runtime/network dependencies. | Normal test run without OpenCode, `agentwrap`, network, or credentials. |
| CLI Surface contract | Listing output remains plain, deterministic, and capturable. | Command output tests using existing app test seams. |

## Tasks

- [x] **Task 1: Inspect Current Study Seams**
  > Executes: Decisions 1, 2, 6, 8; Architecture contract
  - [x] Confirm current `internal/study/domain.go`, `discovery.go`, `resolve.go`, and `service.go` shapes before editing.
  - [x] Confirm whether command output stays in `internal/app/study_commands.go` for this sprint or whether an existing `internal/study/cli.go` seam exists by implementation time.
  - [x] Preserve existing source/dimension resolution behavior and tests unless a requirement explicitly changes it.

- [x] **Task 2: Extend Source Domain Metadata**
  > Executes: Decision 2; AC-02, AC-08, AC-11 through AC-13; `TRD.md` 8.2 and 9A
  - [x] Add `ApplicableDimensions []string` and `Frontmatter map[string]any` or equivalent fields to `Source` in `internal/study/domain.go`.
  - [x] Keep `SourceKindDirectory` and `SourceKindMarkdown` values unchanged.
  - [x] Ensure zero-value `ApplicableDimensions` means all dimensions for Markdown sources and no filter for directory sources.

- [x] **Task 3: Add Markdown Frontmatter Helpers**
  > Executes: Decisions 3, 4, 7; AC-04 through AC-10, AC-16; Errors and Security contracts
  - [x] Create `internal/study/markdown.go` using `gopkg.in/yaml.v3`, already present in `go.mod`.
  - [x] Implement leading-only `parseFrontmatter(content string)` or equivalent helper.
  - [x] Implement `stripFrontmatter(content string) string` that removes only the leading block and leaves non-frontmatter content unchanged.
  - [x] Normalize `applicable_dimensions` values from YAML numeric and string values to two-digit strings using the existing dimension normalization behavior where suitable.
  - [x] Treat missing or empty `applicable_dimensions` as all dimensions by returning an empty normalized filter.
  - [x] Return errors for malformed leading frontmatter and invalid applicability values; call sites must add source path and offending value context while preserving causes.
  - [x] Deduplicate normalized applicability values deterministically.

- [x] **Task 4: Extend Source Discovery**
  > Executes: Decisions 1, 2, 3, 4, 7; AC-01 through AC-03, AC-08 through AC-10, AC-16
  - [x] Update `DiscoverSources` in `internal/study/discovery.go` to inspect only direct children of `studies/<study>/sources/`.
  - [x] Keep non-hidden directories as `SourceKindDirectory` with directory basename names.
  - [x] Add non-hidden top-level `.md` files as `SourceKindMarkdown` with filename-including-extension names.
  - [x] Read Markdown files only enough to parse leading frontmatter metadata required for listing/applicability.
  - [x] Ignore hidden entries, unrelated files, and nested Markdown files.
  - [x] Sort by source name, with deterministic kind/path tie-breakers if needed.
  - [x] Wrap read/parse errors with study/source path context.

- [x] **Task 5: Add Applicability Helper**
  > Executes: Decision 5; AC-11 through AC-14; `TRD.md` 9A.3
  - [x] Create `internal/study/applicability.go` with `GetApplicableSources(sources []Source, dimension Dimension) []Source` or an equivalent exported study helper.
  - [x] Always include directory sources.
  - [x] Include Markdown sources with empty filters for every dimension.
  - [x] Include filtered Markdown sources only when `dimension.Number` matches a normalized applicability value.
  - [x] Preserve input ordering in the returned slice.
  - [x] Represent inapplicable pairs by exclusion, not by an error.

- [x] **Task 6: Update Study Listing Output**
  > Executes: Decision 6; AC-01, AC-02, AC-15, AC-16; CLI Surface contract
  - [x] Update `ultraplan study <study> list` rendering in the existing command seam.
  - [x] Keep directory rows as `  <name> directory`.
  - [x] Render unfiltered Markdown rows as `  <name> markdown all`.
  - [x] Render filtered Markdown rows as `  <name> markdown <comma-separated-normalized-dimensions>`.
  - [x] Keep dimensions output deterministic and preserve existing Sprint 3 behavior.
  - [x] Do not add or stabilize JSON output in this sprint.

- [x] **Task 7: Add Study Unit Tests**
  > Executes: Decisions 1 through 5, 7, 8; AC-03 through AC-14, AC-16
  - [x] Update existing `internal/study/study_test.go` or add `internal/study/source_test.go` for mixed directory/Markdown discovery and ordering.
  - [x] Test hidden entries, unrelated files, and nested Markdown files are ignored.
  - [x] Test source kind assignment and Markdown source names with `.md` extension.
  - [x] Test leading frontmatter, absent frontmatter, non-leading delimiter content, malformed frontmatter, and stripping behavior.
  - [x] Test `applicable_dimensions` numeric/string/zero-padded normalization, empty/missing defaults, duplicate handling, and invalid values with source path/offending value diagnostics.
  - [x] Test applicability helper inclusion/exclusion and ordering.

- [x] **Task 8: Add Command Output Tests**
  > Executes: Decisions 6, 7, 8; AC-01, AC-02, AC-15, AC-16; CLI Surface and Testing contracts
  - [x] Update `internal/app/study_commands_test.go` or add `internal/study/cli_test.go` according to the current command seam.
  - [x] Cover directory source output, unfiltered Markdown source output, filtered Markdown source output, mixed deterministic ordering, and unchanged dimension listing.
  - [x] Cover malformed frontmatter or invalid applicability command failure if the error surfaces through `study <study> list`.
  - [x] Assert output is captured through existing test buffers, not hardcoded terminal streams.

- [x] **Task 9: Scope And Architecture Review Before Verification**
  > Executes: Non-goals; Architecture, Security, CLI Surface contracts
  - [x] Confirm no runtime execution, prompt composition, report validation, summary generation, run-loop/status/synthesis scheduling, target, or sprint workflow behavior was added.
  - [x] Confirm no global `validation`, `scheduler`, `reports`, or `prompts` packages were introduced.
  - [x] Confirm platform packages do not import `internal/study`.

- [x] **Task 10: Run Verification**
  > Executes: AC-17; Testing contract
  - [x] Run `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go`.
  - [x] If tests fail, fix only sprint-related failures or record unrelated blockers in `review.md`.
  - [x] Record final verification evidence for sprint review.

## Evidence Checklist

- [x] Discovery tests prove direct child directories and top-level Markdown files are discovered and hidden/unrelated/nested entries are ignored.
- [x] Source model tests prove kind, normalized applicability, and frontmatter metadata are populated as required.
- [x] Frontmatter tests prove leading-only parse/strip behavior, absent frontmatter behavior, malformed frontmatter handling, normalization, deduplication, and invalid diagnostics.
- [x] Applicability tests prove directory inclusion, unfiltered Markdown inclusion, filtered matching, and inapplicable exclusion.
- [x] Command tests prove deterministic listing output for directory, unfiltered Markdown, filtered Markdown, and mixed sources.
- [x] Error tests or assertions prove source path and offending value are present for invalid `applicable_dimensions` values and cause chains are preserved.
- [x] Architecture review confirms implementation remains in `internal/study` plus existing command wiring, with no runtime/prompt/report/scheduler expansion.
- [x] `go test ./...` passes from `/home/antonioborgerees/coding/ultraplan-go`.
- [x] Deviations from `reasoning.md` are recorded before implementation continues.
- [x] Required review protocols from `sprint-index.md` have evidence for `review.md`.

## Verification Commands

| Check | Command | Expected Result |
| --- | --- | --- |
| Full test suite | `go test ./...` | Passes from `/home/antonioborgerees/coding/ultraplan-go` without OpenCode, `agentwrap`, network, or provider credentials. |
| Architecture import review | `go list ./...` | Package graph loads; review confirms platform packages do not import `internal/study`. |

## Risks And Blockers

| Risk / Blocker | Source | Mitigation | Status |
| --- | --- | --- | --- |
| Area-specific reasoning artifact is absent despite sprint-index selection. | `reasoning.md#assumptions-and-risks` | Used final decisions in `reasoning.md`; no conflicting artifact appeared during implementation. | closed |
| Existing source tests currently assert directory-only discovery. | Current implementation inspection; `reasoning.md#assumptions-and-risks` | Updated tests to Sprint 5 acceptance criteria while preserving ignored nested Markdown behavior. | closed |
| Listing text may become an implicit API. | `reasoning.md#trade-off-and-debt-analysis` | Kept output simple and test-pinned; no JSON schema stability was claimed. | mitigated |
| Malformed frontmatter hard-fails listing. | `reasoning.md#assumptions-and-risks` | Implemented hard failure with source path, offending value where applicable, and cause-preserving wrapping. | mitigated |
| Dimension normalization above two digits may be ambiguous. | `reasoning.md#assumptions-and-risks` | Followed existing `normalizeDimensionNumber` behavior; no conflict found in tests. | closed |
| Future runtime/status/summary flows need identical applicability semantics. | `reasoning.md#assumptions-and-risks` | Centralized current semantics in `study.GetApplicableSources` and covered it with tests. | mitigated |

## Review Inputs

Review should use:

- `projects/ultraplan-go/sprints/05-study-listing-inspection/requirements.md`
- `projects/ultraplan-go/sprints/05-study-listing-inspection/sprint-index.md`
- `projects/ultraplan-go/sprints/05-study-listing-inspection/technical-handbook.md`
- `projects/ultraplan-go/sprints/05-study-listing-inspection/reasoning.md`
- this `plan.md`
- implementation diff
- verification evidence
- `system/protocols/architecture-review-protocol.md`
- `system/protocols/sprint-review-protocol.md` (selected by sprint index but absent from this workspace; recorded in `review.md`)

## Execution Log

| Date / Step | Action | Evidence / Notes |
| --- | --- | --- |
| 2026-05-30 / planning | Created implementation plan from requirements, reasoning, sprint index, technical handbook, PRD/TRD/Architecture docs, and current study/listing code seams. | No code implementation performed. |
| 2026-05-31 / implementation | Added Markdown source discovery, leading frontmatter parsing/stripping, normalized `applicable_dimensions`, and study-owned applicability filtering. | Changed `/home/antonioborgerees/coding/ultraplan-go/internal/study/domain.go`, `discovery.go`, `markdown.go`, `applicability.go`, and `study_test.go`. |
| 2026-05-31 / CLI output | Updated existing app command seam to render directory rows unchanged and Markdown rows as `markdown all` or comma-separated normalized dimensions. | Changed `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_commands.go` and `study_commands_test.go`. |
| 2026-05-31 / verification | Ran `go test ./internal/study ./internal/app`, `go test ./...`, and `go list ./...`. | All passed. `go list ./...` package graph loads and no platform package imports `internal/study`. |
| 2026-05-31 / review artifact | Applied available Architecture Review protocol and sprint completion checks in `review.md`. | Selected `system/protocols/sprint-review-protocol.md` is absent from the workspace; review records the missing protocol as an artifact deferral, not an implementation blocker. |

## Completion Criteria

- [x] All tasks are complete or explicitly deferred with requirement mapping.
- [x] Verification commands were run or deferrals are documented.
- [x] Evidence satisfies the expectations from `reasoning.md`.
- [x] `review.md` can evaluate conformance without guessing intent.
- [x] Non-goals remain excluded: no runtime execution, prompt composition, report validation, summary generation, run-loop/status/synthesis scheduling, target workflows, or sprint execution workflows.
