# Sprint Plan: Project Domain and Project Index

> Project: `ultraplan-go`
> Sprint: `16-project-domain-and-index`
> Source: `reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/roadmap.md`, `projects/ultraplan-go/sprints/16-project-domain-and-index/requirements.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/sprints/16-project-domain-and-index/sprint-index.md`, `projects/ultraplan-go/sprints/16-project-domain-and-index/technical-handbook.md`, `projects/ultraplan-go/sprints/16-project-domain-and-index/reasoning.md`, `templates/sprint-plan.md`

This plan executes `reasoning.md`. It must not invent architecture, scope, or decisions.

## Reasoning Source

- **Sprint Reasoning:** `projects/ultraplan-go/sprints/16-project-domain-and-index/reasoning.md`
- **Sprint Index:** `projects/ultraplan-go/sprints/16-project-domain-and-index/sprint-index.md`
- **Technical Handbook:** `projects/ultraplan-go/sprints/16-project-domain-and-index/technical-handbook.md`
- **Area Reasoning:** none present under `projects/ultraplan-go/sprints/16-project-domain-and-index/reasoning/*.md`

## Sprint Status

- **Status:** complete
- **Owner:** implementation agent
- **Start Date:** 2026-06-14
- **Completion Date:** 2026-06-14

## Decisions To Execute

| Decision | Source Section | Execution Implication |
| --- | --- | --- |
| Own project behavior in `internal/project` | `reasoning.md#decision-1-own-project-behavior-in-internalproject` | Add focused project package files for domain, docs, discovery, store, index parsing, validation, and service use cases; keep project business rules out of `internal/app`, platform, study, and sprint modules. |
| Discover and resolve projects by direct safe names only | `reasoning.md#decision-2-discover-and-resolve-projects-by-direct-safe-names-only` | Scan only direct non-hidden directories under `projects/`, validate project names as safe path segments, sort results, and return actionable missing, invalid, or ambiguous reference diagnostics. |
| Read project files through a workspace-safe filesystem store | `reasoning.md#decision-3-read-project-files-through-a-workspace-safe-filesystem-store` | Read `docs/*.md`, `roadmap.md`, `project-index.md`, and direct sprint directory metadata through workspace-safe path helpers; keep status and validation bounded and read-only. |
| Parse current project-index catalog tables narrowly and validate findings structurally | `reasoning.md#decision-4-parse-current-project-index-catalog-tables-narrowly-and-validate-findings-structurally` | Implement a narrow parser for the current catalog sections and structured validation findings with section, entry name, path, problem, cause, severity, and suggested action. |
| Expose project commands as thin, deterministic, runtime-free CLI wiring | `reasoning.md#decision-5-expose-project-commands-as-thin-deterministic-runtime-free-cli-wiring` | Add `ultraplan project list`, `ultraplan project <project> status`, and `ultraplan project <project> validate`; handlers only parse, delegate, render stable text, and map exit codes. |
| Verify with fixture-first, command-level, review-protocol evidence | `reasoning.md#decision-6-verification-will-be-fixture-first-command-level-and-review-protocol-driven` | Add project package tests and app command tests, then record test, build, architecture, and sprint review evidence in `review.md` after implementation. |

## Requirements / Contracts To Satisfy

| Contract / Requirement ID | Required Behavior | Evidence Planned |
| --- | --- | --- |
| Architecture | Product behavior stays module-owned; `internal/project` owns project catalog behavior; `internal/app` stays thin; platform packages stay product-agnostic. | Import review, package placement review, command tests that exercise service delegation, and absence of prohibited packages. |
| Errors | Failures preserve causes internally and render safe actionable diagnostics with meaningful exit codes. | Unit tests for finding/error classification and command tests for stderr and exit code `5`. |
| Security | Workspace-managed paths are normalized; catalog path escapes fail closed; diagnostics avoid unsafe absolute-path leakage where possible. | Tests for `..` escapes, absolute escapes, missing paths, explicit external handling, and safe output. |
| Testing | Deterministic offline unit, fixture, and command tests cover domain and CLI behavior. | `internal/project/project_test.go`, `internal/app/project_commands_test.go`, `go test ./...`, and `go test -race ./...`. |
| Documentation | Package docs and CLI help describe ownership, commands, and catalog-only semantics. | `internal/project/doc.go`, top-level help tests, project help tests, and command-specific help tests. |
| CLI Surface | Project commands are stable, script-friendly, documented, and map validation failures to exit code `5`. | Command tests for help, output, stdout/stderr separation, deterministic ordering, and exit codes. |
| Performance | Discovery and status are bounded and avoid recursive source/evidence tree scans. | Tests with hidden, nested, file, and irrelevant entries plus code review for bounded scans. |
| AC-01, AC-02, AC-03 | Project list discovers direct project roots, ignores hidden/non-directory/nested entries, validates names, and sorts by name. | Discovery tests and `project list` command tests. |
| AC-04 | Project status reports docs, Markdown docs, `roadmap.md`, `project-index.md`, `sprints/`, and sprint directory health. | Project status unit tests and command output tests. |
| AC-05, AC-06 | Project validate exits `0` when valid and exit code `5` when required files, rows, paths, or path safety checks fail. | Validation unit tests and command tests for success/failure exit codes. |
| AC-07, AC-08, AC-09, AC-10 | Project-index parsing recognizes source documents, contracts, evidence reports, reasoning templates, and review protocols; paths are workspace-safe; project index is catalog-only. | Parser tests based on the real `project-index.md`, validation tests for missing/escaping paths, and package doc/help review. |
| AC-11, AC-12, AC-15, AC-18 | Project commands do not use study internals, runtime, agentwrap, network, workflow execution, or global catalog/validation/planning packages. | Import review, package review, command tests requiring no runtime fakes, and diff review for prohibited package creation. |
| AC-16, AC-17 | Tests cover named fixtures and help documents the project command family. | Fixture tests for valid/missing/malformed/path cases and help tests for top-level/project subcommands. |
| AC-19, AC-20, AC-21 | Offline tests, race tests, and CLI build pass from `/home/antonioborgerees/coding/ultraplan-go`. | `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan`. |

## Tasks

- [x] **Task 1: Define Project Package Contract**
  > Executes: Decision 1, Architecture, Documentation, AC-10, AC-11, AC-12, AC-18
  - [x] Create `/home/antonioborgerees/coding/ultraplan-go/internal/project/doc.go` documenting `internal/project` ownership, boundaries, and catalog-only semantics.
  - [x] Create `/home/antonioborgerees/coding/ultraplan-go/internal/project/domain.go` with `Project`, `ProjectIndex`, catalog entry types, project status categories, validation finding/result types, and stable status constants.
  - [x] Keep domain names project-specific and avoid study, sprint-stage, report, rating, scheduler, runtime, or workflow terms unless required by Sprint 16.
  - [x] Confirm no new global `internal/catalog`, `internal/validation`, `internal/reports`, `internal/prompts`, `internal/planning`, or workflow-engine package is introduced.

- [x] **Task 2: Implement Discovery And Reference Resolution**
  > Executes: Decision 2, Security, Performance, AC-01, AC-02, AC-03, AC-13, AC-16
  - [x] Add `/home/antonioborgerees/coding/ultraplan-go/internal/project/discovery.go` for direct-child discovery under `projects/` only.
  - [x] Ignore hidden entries, non-direct children, and non-directory entries without recursively scanning project contents.
  - [x] Validate project names as filesystem-safe single path segments before resolving workspace paths.
  - [x] Sort discovered project names and missing/ambiguous reference candidates deterministically.
  - [x] Return actionable diagnostics for invalid, missing, and ambiguous project references.

- [x] **Task 3: Implement Workspace-Safe Project Store**
  > Executes: Decision 3, Security, Performance, AC-04, AC-08, AC-12, AC-15, AC-16
  - [x] Add `/home/antonioborgerees/coding/ultraplan-go/internal/project/store_fs.go` for read-only access to project docs, `roadmap.md`, `project-index.md`, and direct sprint directory metadata.
  - [x] Normalize all project-managed paths through existing workspace-safe path helpers before filesystem access.
  - [x] Discover only Markdown docs under `docs/*.md` and direct sprint directories under `sprints/`.
  - [x] Return structured presence/health inputs for status and validation instead of letting CLI handlers inspect the filesystem.
  - [x] Keep the store free of runtime, network, Git, study service, and sprint flow behavior.

- [x] **Task 4: Parse `project-index.md` Catalog Tables**
  > Executes: Decision 4, Architecture, Security, Documentation, AC-07, AC-08, AC-10, AC-18
  - [x] Add `/home/antonioborgerees/coding/ultraplan-go/internal/project/index.go` with a narrow Markdown table parser for the current catalog sections in `projects/ultraplan-go/project-index.md`.
  - [x] Recognize source documents, active contracts, available evidence reports, available reasoning templates, and review protocols.
  - [x] Preserve section name, entry name, referenced path, and descriptive fields needed for validation diagnostics.
  - [x] Treat unknown sections as outside Sprint 16 unless a requirements-backed rule is added.
  - [x] Do not infer sprint task plans, selected sprint context, planning-stage state, or study semantics from catalog rows.

- [x] **Task 5: Implement Validation And Status Service**
  > Executes: Decisions 1, 3, and 4, Errors, Security, AC-04, AC-05, AC-06, AC-08, AC-09, AC-14, AC-16
  - [x] Add `/home/antonioborgerees/coding/ultraplan-go/internal/project/validation.go` for required file checks, docs checks, catalog row shape checks, and catalog path resolution checks.
  - [x] Resolve catalog paths relative to the workspace and fail closed on missing files or workspace escapes unless an entry is explicitly modeled as external.
  - [x] Produce sorted structured findings containing severity, catalog section, entry name, referenced path, problem/cause, and suggested action.
  - [x] Render or expose status summaries for docs, Markdown docs, roadmap, project index, sprints directory, sprint directories, and catalog health.
  - [x] Preserve cause chains internally while ensuring user-facing diagnostics stay safe and actionable.

- [x] **Task 6: Expose Project Service API**
  > Executes: Decisions 1, 2, 3, 4, and 5, Architecture, CLI Surface, AC-01 through AC-15
  - [x] Add `/home/antonioborgerees/coding/ultraplan-go/internal/project/service.go` with use-case methods for list, status, and validate.
  - [x] Hide parser and filesystem internals behind service results used by `internal/app`.
  - [x] Keep service construction explicit and small, using workspace/path infrastructure and project store dependencies only.
  - [x] Keep validation result structs structured enough for future JSON output without adding a Sprint 16 JSON schema unless existing app conventions require it.
  - [x] Ensure project service does not import `internal/study`, `internal/sprint`, `internal/platform/runtime`, agentwrap, OpenCode, or workflow packages.

- [x] **Task 7: Wire CLI Commands And Help**
  > Executes: Decision 5, CLI Surface, Documentation, Errors, AC-01, AC-04, AC-05, AC-06, AC-13, AC-14, AC-17
  - [x] Add `/home/antonioborgerees/coding/ultraplan-go/internal/app/project_commands.go` with thin wiring for `ultraplan project list`, `ultraplan project <project> status`, and `ultraplan project <project> validate`.
  - [x] Register the project command family in `/home/antonioborgerees/coding/ultraplan-go/internal/app/app.go` without disturbing existing commands.
  - [x] Follow existing app conventions for argument parsing, output streams, help text, and exit code mapping.
  - [x] Map validation findings to exit code `5` and usage/reference failures to the established app error categories.
  - [x] Render deterministic, concise, script-friendly text with no prompts, spinners, ANSI color, runtime health checks, network calls, source cloning, study execution, or sprint flow execution.

- [x] **Task 8: Add Project Package Tests**
  > Executes: Decision 6, Testing, Security, Performance, AC-01 through AC-16
  - [x] Add `/home/antonioborgerees/coding/ultraplan-go/internal/project/project_test.go` with fixture-first unit tests.
  - [x] Cover valid fixture project, direct-child discovery, hidden/file/nested exclusion, deterministic sorting, invalid names, missing references, and ambiguous references.
  - [x] Cover docs present, empty docs, missing `roadmap.md`, missing `project-index.md`, sprint directory discovery, and status summaries.
  - [x] Cover catalog parsing for source documents, active contracts, evidence reports, reasoning templates, and review protocols using a fixture based on the real project index.
  - [x] Cover malformed catalog rows, unresolved catalog references, workspace path escapes, explicit external handling when implemented, sorted findings, and safe diagnostic fields.

- [x] **Task 9: Add Command Tests**
  > Executes: Decisions 5 and 6, CLI Surface, Documentation, Errors, Testing, AC-01, AC-04, AC-05, AC-06, AC-13, AC-14, AC-17
  - [x] Add `/home/antonioborgerees/coding/ultraplan-go/internal/app/project_commands_test.go` with command-level tests using controlled fixture workspaces and captured output streams.
  - [x] Cover top-level help listing the `project` family and project-specific help for `project`, `list`, `status`, and `validate`.
  - [x] Cover `project list`, `project <project> status`, and `project <project> validate` success output with deterministic ordering.
  - [x] Cover validation failure output, stderr/stdout separation, exit code `5`, missing/ambiguous project diagnostics, and no ANSI escape sequences.
  - [x] Assert command tests do not need runtime fakes, OpenCode, provider credentials, network access, Git commands, or long sleeps.

- [x] **Task 10: Verify And Record Review Evidence**
  > Executes: Decision 6, Review Expectations, AC-19, AC-20, AC-21
  - [x] Run `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go` and record the result.
  - [x] Run `go test -race ./...` from `/home/antonioborgerees/coding/ultraplan-go` and record the result.
  - [x] Run `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go` and record the result.
  - [x] Perform Architecture Review checks from `system/protocols/architecture-review-protocol.md` for package ownership, dependency direction, runtime/study separation, and prohibited package absence.
  - [x] Create `projects/ultraplan-go/sprints/16-project-domain-and-index/review.md` with implementation evidence, deviations, verification results, and review conclusions after implementation.

## Evidence Checklist

- [x] Tests prove project discovery, reference resolution, safe name validation, bounded scans, deterministic ordering, status summaries, catalog parsing, validation findings, and command output.
- [x] Diagnostic evidence proves validation failures include section, entry name, referenced path, problem/cause, suggested action, and safe user-facing rendering.
- [x] Documentation updates are complete in `internal/project/doc.go` and CLI help.
- [x] Deviations from `reasoning.md` are recorded before implementation continues.
- [x] Architecture Review and Sprint Review protocols have evidence in `review.md`.
- [x] Non-goal exclusions are checked: no sprint flow state, planning-stage generation, implementation execution, smoke/review automation, issue tracking, Git mutation, study service consumption, runtime invocation, or network behavior.

## Verification Commands

| Check | Command | Expected Result |
| --- | --- | --- |
| Offline tests | `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go` | Passes without OpenCode, provider credentials, network access, Git commands, or long sleeps. |
| Race tests | `go test -race ./...` from `/home/antonioborgerees/coding/ultraplan-go` | Passes without data races. |
| CLI build | `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go` | Builds the CLI binary successfully. |
| Help smoke | `ultraplan --help`, `ultraplan project --help`, `ultraplan project list --help`, `ultraplan project <project> status --help`, `ultraplan project <project> validate --help` | Help lists the project command family and subcommands. |
| Import review | Inspect imports for `internal/project`, `internal/app`, and platform packages | `internal/project` has no study/sprint/runtime/agentwrap/OpenCode imports; platform packages do not import product modules. |
| Package review | Inspect created packages | No global `internal/catalog`, `internal/validation`, `internal/reports`, `internal/prompts`, `internal/planning`, or workflow-engine package exists. |

## Risks And Blockers

| Risk / Blocker | Source | Mitigation | Status |
| --- | --- | --- | --- |
| Current `project-index.md` table shapes may not cover future variants. | `reasoning.md` assumption | Kept parser narrow, tested recognized current-format sections, and left future schema variants out of scope. | mitigated |
| Workspace path helpers may be missing or too study-specific. | `reasoning.md` assumption | Existing `workspace.ResolveInside` and `workspace.Rel` were sufficient; no new helper required. | closed |
| Project reference resolution behavior could be ambiguous if exact and prefix matching conventions differ from existing app patterns. | `reasoning.md` assumption | Matched existing study exact-then-prefix resolution style with deterministic missing/ambiguous diagnostics. | closed |
| External catalog entry modeling is underspecified. | `reasoning.md` risk | Implemented URL-like external detection only; workspace-relative and escaping paths fail closed. | mitigated |
| Status and validate could diverge if they duplicate catalog health logic. | `reasoning.md` risk | Status derives catalog health from the same validation path as validate. | closed |
| Text output may become automation-sensitive before JSON output is scoped. | `reasoning.md` risk | Text output is concise and covered by command tests; JSON schema remains deferred. | mitigated |
| Project commands could accidentally pull in study or runtime dependencies through app composition. | `reasoning.md` risk | Import review found no study, sprint, runtime, agentwrap, or OpenCode imports in `internal/project`. | closed |
| Diagnostics may leak absolute local paths. | `reasoning.md` risk | User-facing finding paths are catalog/workspace-relative; escape causes can include existing workspace helper detail. | mitigated |

## Open Questions To Resolve During Implementation

| Question | Source | Required Handling |
| --- | --- | --- |
| Are existing workspace path helpers sufficient for project and catalog path safety? | `reasoning.md` assumption | Inspect existing workspace APIs during implementation; add only generic helper behavior if needed and keep project semantics in `internal/project`. |
| What project reference resolution convention is already established by the CLI? | `reasoning.md` assumption | Follow existing exact/prefix behavior if present; otherwise implement exact plus deterministic unambiguous reference handling with clear missing/ambiguous diagnostics. |
| How should explicitly external catalog entries be identified? | `reasoning.md` risk | Support only clearly modeled external entries or URL-like references; fail closed for ambiguous workspace-escaping paths. |
| Should status reuse full catalog validation or a lighter shared health helper? | `reasoning.md` risk | Prefer one validation-backed source of truth so `status` and `validate` do not diverge. |

## Review Inputs

Review should use:

- `projects/ultraplan-go/sprints/16-project-domain-and-index/sprint-index.md`
- `projects/ultraplan-go/sprints/16-project-domain-and-index/technical-handbook.md`
- `projects/ultraplan-go/sprints/16-project-domain-and-index/reasoning.md`
- `projects/ultraplan-go/sprints/16-project-domain-and-index/plan.md`
- `projects/ultraplan-go/sprints/16-project-domain-and-index/requirements.md`
- `projects/ultraplan-go/project-index.md`
- `projects/ultraplan-go/roadmap.md`
- `projects/ultraplan-go/docs/PRD.md`
- `projects/ultraplan-go/docs/TRD.md`
- `projects/ultraplan-go/docs/ARCHITECTURE.md`
- `system/protocols/architecture-review-protocol.md`
- `system/protocols/sprint-review-protocol.md`
- implementation diff
- verification evidence

No area-specific reasoning files were present when this plan was written.

## Execution Log

| Date / Step | Action | Evidence / Notes |
| --- | --- | --- |
| 2026-06-14 / planning | Created implementation plan from `reasoning.md` and required inputs. | No code implementation performed during planning. |
| 2026-06-14 / implementation | Added `internal/project` domain, discovery, store, parser, validation, service, CLI wiring, package tests, and command tests. | Files changed under `/home/antonioborgerees/coding/ultraplan-go/internal/project`, `/home/antonioborgerees/coding/ultraplan-go/internal/app/project_commands.go`, `/home/antonioborgerees/coding/ultraplan-go/internal/app/project_commands_test.go`, and `/home/antonioborgerees/coding/ultraplan-go/internal/app/app.go`. |
| 2026-06-14 / verification | Ran required verification commands from `/home/antonioborgerees/coding/ultraplan-go`. | `go test ./...` passed; `go test -race ./...` passed; `go build ./cmd/ultraplan` passed. |
| 2026-06-14 / review | Performed architecture/package/import checks and wrote sprint review evidence. | `rg "internal/(study|sprint|platform/runtime)|agentwrap|opencode" internal/project internal/app/project_commands.go` returned no matches; prohibited global package search returned no matches. |

## Completion Criteria

- [x] All required outputs in `requirements.md` are implemented or explicitly deferred with requirements-backed justification.
- [x] All tasks are complete or explicitly deferred before review.
- [x] Verification commands were run or deferrals are documented with blockers.
- [x] Evidence satisfies `reasoning.md` expected evidence and `requirements.md` review expectations.
- [x] `review.md` can evaluate conformance without guessing intent.
- [x] Non-goals remain excluded, especially Sprint 17 flow state, runtime-backed planning generation, sprint execution, smoke/review automation, issue tracking, and Git mutation.
