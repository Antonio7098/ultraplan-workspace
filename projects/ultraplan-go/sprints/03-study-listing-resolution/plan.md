# Sprint Plan: Study Domain, Listing, and Resolution

> Project: `ultraplan-go`
> Sprint: `03-study-listing-resolution`
> Source: `reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/03-study-listing-resolution/requirements.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/sprints/03-study-listing-resolution/sprint-index.md`, `projects/ultraplan-go/sprints/03-study-listing-resolution/technical-handbook.md`, `projects/ultraplan-go/sprints/03-study-listing-resolution/reasoning/architecture.md`, `projects/ultraplan-go/sprints/03-study-listing-resolution/reasoning.md`, `templates/sprint-plan.md`

This plan executes `reasoning.md`. It must not invent architecture, scope, or decisions.

## Reasoning Source

- **Sprint Reasoning:** `projects/ultraplan-go/sprints/03-study-listing-resolution/reasoning.md`
- **Sprint Index:** `projects/ultraplan-go/sprints/03-study-listing-resolution/sprint-index.md`
- **Technical Handbook:** `projects/ultraplan-go/sprints/03-study-listing-resolution/technical-handbook.md`
- **Area Reasoning:** `projects/ultraplan-go/sprints/03-study-listing-resolution/reasoning/architecture.md`

## Sprint Status

- **Status:** `complete`
- **Owner:** `implementation agent`
- **Start Date:** `2026-05-30`
- **Completion Date:** `2026-05-30`

## Decisions To Execute

| Decision | Source Section | Execution Implication |
| --- | --- | --- |
| Create a minimal `internal/study` read model | `reasoning.md#decision-1-create-a-minimal-internalstudy-read-model` | Add `domain.go`, `discovery.go`, `resolve.go`, and `service.go` under `/home/antonioborgerees/coding/ultraplan-go/internal/study`; keep discovery, sorting, normalization, resolution, and listing use cases in this package. |
| Use shallow deterministic filesystem discovery | `reasoning.md#decision-2-use-shallow-deterministic-filesystem-discovery` | Read only direct children of workspace `studies/`, study `sources/`, and study `dimensions/`; ignore hidden entries where required; sort results deterministically; do not walk source repository contents. |
| Normalize dimension identity from filenames | `reasoning.md#decision-3-normalize-dimension-identity-from-filenames` | Build `Dimension` values from direct `.md` filenames beginning with numeric prefixes; normalize discovered numbers to two digits; derive slug and file identity without parsing Markdown body content. |
| Implement exact-first, unambiguous-prefix resolution with structured errors | `reasoning.md#decision-4-implement-exact-first-unambiguous-prefix-resolution-with-structured-errors` | Resolve study/source/dimension refs by exact canonical matches before prefix matches; distinguish not-found and ambiguous errors; include actionable candidate/context data for CLI rendering. |
| Wire thin study listing commands in `internal/app` | `reasoning.md#decision-5-wire-thin-study-listing-commands-in-internalapp` | Add `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_commands.go`; reuse existing workspace discovery/global `--workspace`; delegate domain work to `internal/study`; render stable text output including source kind. |
| Verify with focused domain and command tests only | `reasoning.md#decision-6-verify-with-focused-domain-and-command-tests-only` | Add `/home/antonioborgerees/coding/ultraplan-go/internal/study/study_test.go` and `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_commands_test.go`; use temp fixture workspaces; run `go test ./...` and `go build ./cmd/ultraplan`. |

## Requirements / Contracts To Satisfy

| Contract / Requirement ID | Required Behavior | Evidence Planned |
| --- | --- | --- |
| REQ-AC-01 | `ultraplan study list` lists discovered studies under workspace `studies/` in deterministic ascending order. | Study unit tests for sorted discovery; command test for output order. |
| REQ-AC-02 | `ultraplan study <study> list` lists sources and dimensions in deterministic ascending order. | Study service/discovery tests; command output test. |
| REQ-AC-03 | Study discovery ignores hidden entries and non-directory entries at the `studies/` level. | Fixture test with hidden directories, hidden files, and regular files. |
| REQ-AC-04 | Source discovery returns one source for each non-hidden directory directly under `studies/<study>/sources/`. | Fixture test with multiple direct source dirs plus ignored entries. |
| REQ-AC-05 | Sprint 3 source discovery does not implement Markdown document sources. | Regression test proving top-level `.md` files under `sources/` are ignored. |
| REQ-AC-06 | Dimension discovery returns direct Markdown files under `dimensions/` whose names begin with a numeric dimension prefix. | Fixture test for matching `.md` files and ignored non-matching, nested, hidden, or non-Markdown entries. |
| REQ-AC-07 | Dimension numbers are normalized to two digits. | Domain test for `1-foo.md` and `01-foo.md` normalization. |
| REQ-AC-08 | Dimension references resolve by number, slug, full filename, exact reference, or unambiguous prefix. | Resolver table tests for exact, prefix, ambiguous, and missing dimension refs. |
| REQ-AC-09 | Source references resolve by exact name or unambiguous prefix. | Resolver table tests for exact, prefix, ambiguous, and missing source refs. |
| REQ-AC-10 | Ambiguous or missing study, source, and dimension references return actionable non-zero command errors. | Command tests asserting returned errors and useful message content. |
| REQ-AC-11 | Listing commands use existing workspace discovery and global `--workspace`. | Command test using `--workspace`; code review of `internal/app` wiring. |
| REQ-AC-12 | Listing output includes source kind for discovered directory sources. | Command output test expecting `directory` for source rows. |
| REQ-AC-13 | Listing commands do not recursively scan source repository contents. | Fixture with nested repository-like files that are not surfaced; code review confirms no recursive walk. |
| REQ-AC-14 | `go test ./...` passes from `/home/antonioborgerees/coding/ultraplan-go`. | Verification command output. |
| REQ-AC-15 | `go build ./cmd/ultraplan` passes from `/home/antonioborgerees/coding/ultraplan-go`. | Verification command output. |
| Architecture Contract | Preserve module ownership and dependency direction. | Import review: `internal/study` owns domain behavior, `internal/app` owns CLI wiring/output, platform packages do not import `internal/study`. |
| Errors Contract | Missing and ambiguous references produce actionable diagnostics and non-zero failures. | Structured domain errors plus command tests for error rendering. |
| Security Contract | Keep filesystem behavior bounded and workspace-safe. | Shallow discovery tests and review for workspace-relative output where practical. |
| Testing Contract | Use deterministic tests without OpenCode, network, cloned repositories, or provider credentials. | Temp workspace tests; `go test ./...`. |
| CLI Surface Contract | Commands are stable, script-friendly, and respect existing workspace/help behavior. | Command tests for output and errors; optional manual help/output review. |

## Tasks

- [x] **Task 1: Inspect Existing CLI And Workspace Seams**
  > Executes: `Decision 5`, `REQ-AC-11`, `Architecture Contract`
  - [x] Inspect existing `internal/app` command construction, output writer conventions, and test helpers.
  - [x] Inspect existing `internal/workspace` APIs and global `--workspace` behavior from Sprint 2.
  - [x] Record any mismatch with `reasoning.md` assumptions before editing implementation files.

- [x] **Task 2: Add Study Domain Types**
  > Executes: `Decision 1`, `Decision 3`, `REQ-AC-07`, `Architecture Contract`
  - [x] Create `/home/antonioborgerees/coding/ultraplan-go/internal/study/domain.go` with `Study`, `Source`, `SourceKind`, and `Dimension` fields required by the sprint.
  - [x] Include `SourceKindDirectory` and keep `SourceKindMarkdown` only as a domain constant if needed by PRD/TRD shape while discovery remains directory-only.
  - [x] Add private normalization helpers for dimension numbers and filename-derived dimension identity.
  - [x] Keep domain types independent of CLI parsing and output formatting.

- [x] **Task 3: Add Shallow Discovery**
  > Executes: `Decision 2`, `Decision 3`, `REQ-AC-01` through `REQ-AC-07`, `REQ-AC-13`, `Security Contract`
  - [x] Create `/home/antonioborgerees/coding/ultraplan-go/internal/study/discovery.go`.
  - [x] Discover studies from direct non-hidden directories under `<workspace>/studies` and sort ascending by canonical name.
  - [x] Discover sources from direct non-hidden directories under `<study>/sources`; ignore files, hidden entries, top-level `.md` files, and nested repository contents.
  - [x] Discover dimensions from direct `.md` files under `<study>/dimensions` whose filenames begin with a numeric prefix; normalize numbers to two digits and sort deterministically.
  - [x] Treat absent `studies/`, `sources/`, or `dimensions/` as empty only if existing workspace validation does not require those directories; surface real filesystem read/stat errors.

- [x] **Task 4: Add Reference Resolution**
  > Executes: `Decision 4`, `REQ-AC-08`, `REQ-AC-09`, `REQ-AC-10`, `Errors Contract`
  - [x] Create `/home/antonioborgerees/coding/ultraplan-go/internal/study/resolve.go`.
  - [x] Resolve studies and sources by exact name first, then unambiguous prefix.
  - [x] Resolve dimensions by exact number, slug, full filename, and exact reference aliases before unambiguous prefix matching across allowed candidates.
  - [x] Add structured errors for study/source/dimension not found and ambiguous references, with requested ref and canonical candidate names.
  - [x] Keep shared matching helpers private to `internal/study` if introduced; do not create global matching/validation packages.

- [x] **Task 5: Add Read-Only Listing Service**
  > Executes: `Decision 1`, `Decision 2`, `Decision 4`, `Architecture Contract`
  - [x] Create `/home/antonioborgerees/coding/ultraplan-go/internal/study/service.go`.
  - [x] Expose narrow query methods such as `ListStudies` and `ListStudy` or equivalent use-case names.
  - [x] Ensure `ListStudy` resolves the study reference and returns the resolved study with sorted sources and dimensions.
  - [x] Do not add init, prompt, runtime, scheduler, validation, report, summary, target, or sprint workflow methods in this sprint.

- [x] **Task 6: Wire Study Listing Commands**
  > Executes: `Decision 5`, `REQ-AC-01`, `REQ-AC-02`, `REQ-AC-10`, `REQ-AC-11`, `REQ-AC-12`, `CLI Surface Contract`
  - [x] Create `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_commands.go` using existing command factory and output conventions.
  - [x] Register `ultraplan study list` and `ultraplan study <study> list` without implementing other study commands.
  - [x] Resolve workspace through existing global `--workspace` behavior rather than reimplementing workspace lookup.
  - [x] Render stable text output for studies, sources, and dimensions; include `directory` source kind for directory sources.
  - [x] Map structured study errors to actionable command errors that distinguish missing and ambiguous references.

- [x] **Task 7: Add Study Package Tests**
  > Executes: `Decision 6`, `REQ-AC-01` through `REQ-AC-09`, `REQ-AC-13`, `Testing Contract`
  - [x] Create `/home/antonioborgerees/coding/ultraplan-go/internal/study/study_test.go`.
  - [x] Cover domain normalization for dimension numbers and filename-derived slug/file identity.
  - [x] Cover deterministic study, source, and dimension ordering.
  - [x] Cover ignored hidden entries, non-directory study entries, non-directory source entries, top-level `.md` source files, nested source contents, and non-matching dimension files.
  - [x] Cover exact, prefix, ambiguous, and missing resolution for studies, sources, and dimensions.

- [x] **Task 8: Add Command Tests**
  > Executes: `Decision 6`, `REQ-AC-01`, `REQ-AC-02`, `REQ-AC-10`, `REQ-AC-11`, `REQ-AC-12`, `CLI Surface Contract`
  - [x] Create `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_commands_test.go`.
  - [x] Use temp fixture workspaces and existing command test helpers; do not require OpenCode, network access, cloned repositories, or credentials.
  - [x] Assert `study list` output order and empty-list behavior chosen from existing workspace validation.
  - [x] Assert `study <study> list` output includes sorted sources, sorted dimensions, and source kind.
  - [x] Assert missing and ambiguous study refs return non-zero actionable errors.
  - [x] Assert `--workspace` routes discovery to the explicit fixture workspace.

- [x] **Task 9: Verify And Review**
  > Executes: `REQ-AC-14`, `REQ-AC-15`, selected review protocols`
  - [x] Run `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go`.
  - [x] Run `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`.
  - [x] Review imports and package boundaries for Architecture Review evidence.
  - [x] Confirm deferred scope remains absent: no study init, YAML parsing, Markdown source frontmatter/applicability, prompt/runtime/run-loop/synthesis/report/summary/target/sprint workflows, or non-trivial JSON listing.
  - [x] Capture verification evidence for `review.md`.

## Evidence Checklist

- [x] Tests prove domain normalization, deterministic discovery, ignored entries, shallow discovery, and exact/prefix resolution.
- [x] Command tests prove listing output, source kind display, `--workspace`, and actionable non-zero errors.
- [x] Build evidence exists for `go build ./cmd/ultraplan`.
- [x] Full test evidence exists for `go test ./...`.
- [x] Deviations from `reasoning.md` are recorded before implementation continues.
- [x] Architecture Review evidence confirms package ownership and dependency direction.
- [x] Sprint Review evidence confirms required outputs exist and deferred scope remains excluded.
- [x] No documentation artifact is required beyond command help/output unless implementation reveals an existing docs convention that must be updated.

## Verification Commands

| Check | Command | Expected Result |
| --- | --- | --- |
| Unit and command test suite | `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go` | Passes without OpenCode, network access, cloned repositories, or AI provider credentials. |
| CLI build | `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go` | Builds the CLI binary successfully. |
| Optional manual study list smoke | `go run ./cmd/ultraplan --workspace <fixture-workspace> study list` from `/home/antonioborgerees/coding/ultraplan-go` | Prints discovered studies in deterministic ascending order. |
| Optional manual study detail smoke | `go run ./cmd/ultraplan --workspace <fixture-workspace> study <study> list` from `/home/antonioborgerees/coding/ultraplan-go` | Prints sorted sources with `directory` kind and sorted dimensions. |

## Risks And Blockers

| Risk / Blocker | Source | Mitigation | Status |
| --- | --- | --- | --- |
| Existing Sprint 2 workspace discovery or global `--workspace` API is absent or incompatible. | `reasoning.md#assumptions-and-risks` | Inspected existing app/workspace seams first; implementation uses `discoverWorkspace(deps)` and the existing global flag parser. | closed |
| Missing `studies/` behavior may conflict with existing workspace validation. | `reasoning.md#assumptions-and-risks` | Existing workspace initialization creates `studies/`; discovery treats absent `studies/` as empty and command tests cover empty valid workspaces. | closed |
| Missing `sources/` or `dimensions/` behavior may be unclear for hand-edited studies. | `reasoning.md#assumptions-and-risks` | Discovery treats absent directories as empty sections unless real read errors occur; unit tests cover absent directories. | closed |
| Prefix resolution can surprise users if exact and prefix matching are not ordered consistently. | `reasoning.md#assumptions-and-risks` | Exact match wins; prefix match must be unique; unit tests cover exact, prefix, ambiguous, and missing cases. | closed |
| Error output could leak absolute local paths or become noisy. | `reasoning.md#assumptions-and-risks` | Reference errors include requested ref and canonical candidates, not source filesystem paths; study list still prints workspace root consistent with existing app conventions. | closed |
| Future Markdown document source support may conflict with Sprint 3 ignored `.md` files. | `reasoning.md#trade-off-and-debt-analysis` | `SourceKindMarkdown` remains available as a domain constant while discovery ignores `.md` source files; tests lock Sprint 3 behavior. | carried forward |
| Dimension filename edge cases may be ignored if they do not start with numeric prefixes. | `reasoning.md#assumptions-and-risks` | Acceptance requires numeric prefixes; unit tests cover ignored non-conforming filenames. | closed |

## Review Inputs

Review should use:

- `projects/ultraplan-go/sprints/03-study-listing-resolution/sprint-index.md`
- `projects/ultraplan-go/sprints/03-study-listing-resolution/technical-handbook.md`
- `projects/ultraplan-go/sprints/03-study-listing-resolution/reasoning/architecture.md`
- `projects/ultraplan-go/sprints/03-study-listing-resolution/reasoning.md`
- `projects/ultraplan-go/sprints/03-study-listing-resolution/plan.md`
- implementation diff in `/home/antonioborgerees/coding/ultraplan-go`
- verification evidence from `go test ./...` and `go build ./cmd/ultraplan`
- `system/protocols/architecture-review-protocol.md`
- `system/protocols/sprint-review-protocol.md`

## Execution Log

| Date / Step | Action | Evidence / Notes |
| --- | --- | --- |
| 2026-05-30 / planning | Created sprint implementation plan from `reasoning.md`. | No implementation performed. Plan is grounded in requirements, PRD/TRD/ARCHITECTURE, sprint index, technical handbook, and architecture reasoning. |
| 2026-05-30 / inspection | Inspected `internal/app` command routing, output/error helpers, `runForTest`, and `internal/workspace` discovery/path APIs. | Existing `discoverWorkspace(deps)` and global `--workspace` behavior were compatible with sprint assumptions. |
| 2026-05-30 / implementation | Added `internal/study` domain, shallow discovery, exact-first reference resolution, and read-only listing service. | Files: `internal/study/domain.go`, `internal/study/discovery.go`, `internal/study/resolve.go`, `internal/study/service.go`; updated `internal/study/doc.go`. |
| 2026-05-30 / implementation | Added `ultraplan study list` and `ultraplan study <study> list` command routing and stable text rendering. | Files: `internal/app/app.go`, `internal/app/study_commands.go`; help now includes `study`. |
| 2026-05-30 / tests | Added domain and command tests for deterministic listing, ignored entries, resolution, `--workspace`, source kind output, and actionable errors. | Files: `internal/study/study_test.go`, `internal/app/study_commands_test.go`, `internal/app/app_test.go`. |
| 2026-05-30 / verification | Ran `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go`. | Passed. |
| 2026-05-30 / verification | Ran `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`. | Passed; transient `ultraplan` build artifact removed. |

## Completion Criteria

- [x] All tasks are complete or explicitly deferred.
- [x] Required implementation files exist at the paths listed in `requirements.md`.
- [x] Verification commands were run or deferrals are documented.
- [x] Evidence satisfies the expectations from `reasoning.md`.
- [x] Missing and ambiguous reference behavior is covered by tests and command output is actionable.
- [x] Deferred scope remains excluded from implementation.
- [x] `review.md` can evaluate conformance without guessing intent.
