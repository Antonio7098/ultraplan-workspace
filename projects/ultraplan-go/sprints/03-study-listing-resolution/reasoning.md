# Sprint Reasoning: Study Domain, Listing, and Resolution

> Project: `ultraplan-go`
> Sprint: `03-study-listing-resolution`
> Output: `projects/ultraplan-go/sprints/03-study-listing-resolution/reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/03-study-listing-resolution/requirements.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/sprints/03-study-listing-resolution/sprint-index.md`, `projects/ultraplan-go/sprints/03-study-listing-resolution/technical-handbook.md`, `projects/ultraplan-go/sprints/03-study-listing-resolution/reasoning/architecture.md`, `templates/sprint-reasoning.md`

This document decides. It synthesizes selected context, handbook evidence, area-specific reasoning, and contracts into final sprint decisions.

It does not replace `sprint-index.md`, `technical-handbook.md`, or `reasoning/*.md`.

## Sprint Purpose

- **Goal:** Implement the read-only study domain model, deterministic filesystem discovery, reference resolution, listing service, and CLI commands needed to inspect existing studies.
- **Non-Goals:** Do not implement study initialization, YAML parsing, Markdown document source discovery, frontmatter parsing, applicability filtering, prompt composition, runtime execution, synthesis, report validation, summaries, run state, retry, cancellation, code-reference extraction, target workflows, sprint workflows, OpenCode wiring, or non-trivial JSON listing output.
- **Depends On:** Sprint 1 command registration and buildable CLI entrypoint; Sprint 2 workspace discovery and global `--workspace`; PRD study listing requirements; TRD study/source/dimension domain rules; ARCHITECTURE package boundaries.

## Selected Context And Pre-Reasoning Artifacts

| Artifact | Path | How It Was Used |
| --- | --- | --- |
| Project Index | `projects/ultraplan-go/project-index.md` | Confirmed project scope, target repository, active contracts, selected study evidence pool, and available architecture reasoning template. |
| Sprint Requirements | `projects/ultraplan-go/sprints/03-study-listing-resolution/requirements.md` | Treated as the authoritative sprint contract for required outputs, acceptance criteria, non-goals, constraints, dependencies, and review expectations. |
| PRD | `projects/ultraplan-go/docs/PRD.md` | Grounded product behavior for study listing, stable output, source kind display, workspace-local artifacts, and deferred target/sprint workflows. |
| TRD | `projects/ultraplan-go/docs/TRD.md` | Grounded technical behavior for `Study`, `Source`, `SourceKind`, `Dimension`, deterministic sorting, lookup rules, workspace path safety, and testing expectations. |
| Architecture | `projects/ultraplan-go/docs/ARCHITECTURE.md` | Applied module ownership and dependency direction: `internal/study` owns study behavior, `internal/app` owns command wiring, `internal/workspace` owns workspace resolution, platform packages must not import product modules. |
| Sprint Index | `projects/ultraplan-go/sprints/03-study-listing-resolution/sprint-index.md` | Selected contracts, evidence reports, architecture reasoning, excluded context, and required review protocols. |
| Technical Handbook | `projects/ultraplan-go/sprints/03-study-listing-resolution/technical-handbook.md` | Supplied studied Go CLI evidence for thin command handlers, injected dependencies, actionable errors, shallow discovery, deterministic tests, and avoiding premature abstractions. |
| Architecture Reasoning | `projects/ultraplan-go/sprints/03-study-listing-resolution/reasoning/architecture.md` | Supplied the area-specific architecture conclusion adopted by this sprint: a read-only `internal/study` module with a narrow service consumed by thin `internal/app` commands. |

## Area-Specific Reasoning Inputs

| Area | Reasoning Document | Key Conclusion | Evidence Basis | Impact On Final Decision |
| --- | --- | --- | --- | --- |
| Architecture | `projects/ultraplan-go/sprints/03-study-listing-resolution/reasoning/architecture.md` | Proceed with a minimal module-owned implementation: `internal/study` owns domain types, shallow discovery, reference resolution, and listing service; `internal/app` owns CLI registration, arguments, workspace option consumption, and output rendering. | `ARCHITECTURE.md` module ownership rules; sprint constraints; handbook evidence from chezmoi, gh-cli, helm, yq, go-task, restic, age, and gdu; acceptance criteria for deterministic listing and actionable errors. | Adopted directly. Final decisions reject CLI-owned discovery, global technical packages, recursive discovery, and Sprint 3 Markdown document source handling. |

## Sprint Technical Handbook Summary

- **Relevant Patterns:** Thin CLI/domain-owned behavior; command factories with injected services and writers; explicit composition over globals; actionable typed/sentinel-style errors; deterministic bounded discovery; fixture-based unit and command tests.
- **Important Trade-Offs:** A small `internal/study.Service` boundary adds setup but keeps domain rules out of command handlers; prefix resolution improves UX but requires exact-match precedence and ambiguity detection; shallow discovery is fast and safe but cannot infer nested repository metadata.
- **Warnings / Anti-Patterns:** Do not put discovery and resolution in `RunE`; do not add global `internal/validation`, `internal/reports`, or `internal/scheduler`; do not recursively scan source repositories for listing; do not collapse missing and ambiguous references into generic errors; do not overbuild JSON, TUI, runtime, or workflow abstractions.
- **Evidence Confidence:** High for command/domain boundaries, error handling, tests, and bounded discovery because multiple selected reports cite mature Go CLIs and concrete source references. Medium for terminal formatting and security details because they guide shape and constraints but are less specific to this exact listing feature.

## Contracts Applied

Current gate: Skeleton / Local CLI. The selected contracts apply to implemented read-only listing behavior only. Runtime/provider observability, durable workflow state, diagnostics, migration, and release-grade public error payload requirements are deferred until their roadmap gates enter scope.

| Contract / Requirement ID | Constraint | Decision Impact | Expected Evidence |
| --- | --- | --- | --- |
| REQ-AC-01 | `ultraplan study list` lists discovered studies in deterministic ascending order. | Discovery sorts non-hidden study directories under workspace `studies/`. | `internal/study/study_test.go`; `internal/app/study_commands_test.go`; manual fixture command if reviewed. |
| REQ-AC-02 | `ultraplan study <study> list` lists sources and dimensions in deterministic ascending order. | Service resolves the study then returns sorted directory sources and dimensions. | Unit tests for ordering; command output tests. |
| REQ-AC-03 | Hidden and non-directory entries at `studies/` are ignored. | Study discovery reads only direct non-hidden directories. | Fixture tests with hidden dirs/files and regular files. |
| REQ-AC-04 | Source discovery returns one source for each non-hidden directory directly under `sources/`. | Sprint 3 source discovery is direct-child directory only. | Unit tests; review confirms no recursive traversal. |
| REQ-AC-05 | Markdown document sources are deferred to Sprint 5. | Top-level `.md` files under `sources/` are ignored in Sprint 3. | Regression test proving `.md` source files are ignored; review confirms no frontmatter parsing. |
| REQ-AC-06 | Dimension discovery returns direct Markdown files under `dimensions/` whose names begin with a numeric prefix. | Dimension discovery parses filenames only and ignores nested or non-matching entries. | Unit tests for matching and ignored entries. |
| REQ-AC-07 | Dimension numbers are normalized to two digits. | Domain normalization produces canonical `NN` numbers. | Domain unit tests. |
| REQ-AC-08 | Dimension references resolve by number, slug, full filename, exact reference, or unambiguous prefix. | Resolution candidate set includes canonical number, slug, filename, and exact reference aliases, with exact match before prefix. | Unit tests for exact, prefix, ambiguous, and missing dimension refs. |
| REQ-AC-09 | Source references resolve by exact name or unambiguous prefix. | Source resolver uses source name candidates only and rejects ambiguous prefixes. | Unit tests for source resolution. |
| REQ-AC-10 | Ambiguous or missing references return actionable non-zero command errors. | Structured domain errors distinguish not found from ambiguity and command layer renders concise guidance. | Command tests assert non-zero error behavior and message content. |
| REQ-AC-11 | Listing commands use existing workspace discovery and global `--workspace`. | `internal/app` command wiring calls prior workspace resolution instead of reimplementing workspace lookup. | Command test using `--workspace`; import/review check. |
| REQ-AC-12 | Listing output includes source kind for discovered directory sources. | Output rows include `directory` for each discovered directory source. | Command output tests. |
| REQ-AC-13 | Listing commands do not recursively scan source repository contents. | Discovery uses direct child reads only; no walk over source directories. | Code review; unit fixture with nested content not discovered. |
| REQ-AC-14 | `go test ./...` passes. | Plan must include full test verification from target repo. | Command output from `/home/antonioborgerees/coding/ultraplan-go`. |
| REQ-AC-15 | `go build ./cmd/ultraplan` passes. | Plan must include build verification from target repo. | Command output from `/home/antonioborgerees/coding/ultraplan-go`. |
| Architecture Contract | Preserve module ownership, dependency direction, and thin entrypoints. | Keep domain discovery/resolution in `internal/study`; CLI wiring/output in `internal/app`; no platform dependency on `study`. | Architecture review; import review. |
| Errors Contract | Produce actionable diagnostics and non-zero failures. | Missing and ambiguous refs use typed/structured errors with requested ref and candidate context. | Unit and command tests for error classes and output. |
| Security Contract | Keep filesystem behavior bounded and workspace-safe. | Prefer workspace-relative paths; do not escape workspace; do not recursively inspect repositories. | Code review and fixture tests. |
| Testing Contract | Use deterministic unit, fixture, and command tests without external runtime dependencies. | Tests use temp workspaces and injected/captured output; no network, OpenCode, or credentials. | `go test ./...`; test file review. |
| CLI Surface Contract | Commands must be stable, script-friendly, and use meaningful exit behavior. | Implement `ultraplan study list` and `ultraplan study <study> list` with stable text output and existing help/workspace behavior. | Command tests and manual help/output review. |

## Repos Studied / Source Evidence Used

| Source / Repo / Report | Concrete Reference | Relevant Finding | Why It Matters For This Sprint | Used In Decision(s) |
| --- | --- | --- | --- | --- |
| `go-cli-study` / `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md`; chezmoi `main.go:26-34`, gh-cli `cmd/gh/main.go:6`, yq `cmd/root.go:9` | Mature CLIs keep entrypoints thin and place behavior in internal/domain packages. | Supports `internal/study` ownership and thin `internal/app` command handlers. | Decisions 1, 5 |
| `go-cli-study` / `02-command-architecture` | `studies/go-cli-study/reports/final/02-command-architecture.md`; gh-cli `pkg/cmd/issue/list/list.go:47-118`, helm `pkg/cmd/install.go:132-145`, restic `cmd/restic/cmd_backup.go:84-115` | Command factory/delegate patterns keep command trees maintainable. | Supports command wiring that delegates to a study service. | Decisions 1, 5 |
| `go-cli-study` / `03-dependency-injection` | `studies/go-cli-study/reports/final/03-dependency-injection.md`; gh-cli `internal/ghcmd/cmd.go:52-132`, `pkg/cmdutil/factory.go:16-43`, go-task `executor.go:93-114` | Explicit composition and constructor/factory injection avoid global state and improve tests. | Supports constructing service/output dependencies at the app edge. | Decisions 1, 5, 6 |
| `go-cli-study` / `04-configuration-management` | `studies/go-cli-study/reports/final/04-configuration-management.md`; restic `internal/global/global.go:139,147`, go-task `internal/flags/flags.go:314-327` | Flag/default precedence should be centralized and explicit. | Supports reusing Sprint 2 workspace/global flag behavior rather than inventing new config handling. | Decision 5 |
| `go-cli-study` / `05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md`; go-task `errors/errors_task.go:13-32`, age `cmd/age/tui.go:37-54`, gh-cli `internal/ghcmd/cmd.go:281-301` | Actionable errors benefit from typed context and hints. | Supports distinct missing vs ambiguous reference errors with candidate context. | Decision 4 |
| `go-cli-study` / `06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md`; gh-cli `pkg/iostreams/iostreams.go:551-568`, go-task `executor.go:553-564`, mitchellh-cli `ui_mock.go:27-33` | Injectable output streams make command output testable. | Supports command tests capturing listing output and errors. | Decisions 5, 6 |
| `go-cli-study` / `09-terminal-ux` | `studies/go-cli-study/reports/final/09-terminal-ux.md`; chezmoi `internal/cmd/prompt.go:20-256`, gh-cli `internal/prompter/` | CLI output should be calm, script-friendly, and avoid unnecessary terminal complexity. | Supports minimal stable text listing instead of rich/TUI output. | Decision 5 |
| `go-cli-study` / `11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md`; chezmoi `internal/cmd/main_test.go:64-174`, go-task `task_test.go:166-169`, gdu `cmd/gdu/app/app_test.go:682` | Mature CLIs combine domain tests, command tests, fixtures, and output assertions. | Supports required unit and command test split. | Decision 6 |
| `go-cli-study` / `13-security` | `studies/go-cli-study/reports/final/13-security.md`; k9s `internal/config/json/validator.go:146`, gdu `internal/common/ignore.go:16-37`, gh-cli `pkg/surveyext/editor_manual.go:23` | Path/input trust boundaries require deliberate validation and safe file operations. | Supports workspace-relative output, no escaping workspace, and shallow discovery. | Decisions 2, 5 |
| `go-cli-study` / `14-performance` | `studies/go-cli-study/reports/final/14-performance.md`; gh-cli `pkg/cmdutil/factory.go:27-42`, yq `pkg/yqlib/stream_evaluator.go:78-113`, gdu `pkg/analyze/parallel.go:13` | Fast CLIs defer expensive work and avoid unbounded scans. | Supports direct-child discovery only and no recursive source scans. | Decisions 2, 3 |
| `go-cli-study` / `15-philosophy` | `studies/go-cli-study/reports/final/15-philosophy.md`; age `age.go:18`, gh-cli `pkg/cmd/factory/default.go:26-46`, restic `internal/backend/backend.go:19-90` | Coherent tools avoid complexity that is not tied to product purpose. | Supports rejecting global technical layers and broad workflow/runtime abstractions in this sprint. | Decisions 1, 5 |

## Trade-Off And Debt Analysis

### Accepted Trade-Offs

| Trade-Off | Benefit | Cost / Constraint Accepted | Why Acceptable Now | Revisit Trigger |
| --- | --- | --- | --- | --- |
| Add a small `internal/study.Service` boundary for read-only listing. | Keeps domain behavior out of CLI handlers and creates a stable use-case seam. | Slightly more structure than direct command-handler filesystem reads. | The service boundary is narrow and required to preserve architecture/testability. | If service grows into a generic manager with unrelated responsibilities. |
| Keep discovery shallow and filesystem-derived on each command. | Fast, safe, deterministic, and simple. | No caching or inferred repository metadata. | Listing is read-only and source directories may be huge. | If listing latency becomes unacceptable with very large numbers of direct children. |
| Ignore top-level `.md` source files in Sprint 3. | Keeps sprint aligned with explicit requirements and avoids frontmatter/applicability scope. | Broader PRD/TRD Markdown source behavior remains temporarily incomplete. | Sprint requirements explicitly defer Markdown document sources to Sprint 5. | Sprint 5 or revised requirements introduce Markdown document source discovery. |
| Use filename-derived dimensions without parsing Markdown content. | Satisfies number/slug/file discovery with minimal IO. | Listing title/purpose may be unavailable or derived rather than content-rich. | Sprint 3 acceptance focuses on direct Markdown filename discovery and normalized numbers. | Requirements demand dimension content parsing, title extraction, or prompt metadata. |
| Exact match first, then unambiguous prefix resolution. | Supports ergonomic refs while avoiding surprising ambiguous matches. | Resolver must maintain candidate sets and ambiguity errors. | Required by acceptance criteria and supported by CLI error evidence. | Users need fuzzy matching, typo suggestions, or alias ranking. |
| Minimal stable text output only. | Avoids output complexity and keeps tests simple. | JSON/table/color polish is deferred. | Non-trivial JSON is a sprint non-goal and text listing is sufficient for review. | Existing CLI output foundation makes JSON trivial or automation requirements change. |

### Potential Technical Debt

| Debt / Shortcut | Why It Might Accrue | Current Mitigation | Owner / Follow-Up |
| --- | --- | --- | --- |
| Sprint 3 `SourceKind` may include Markdown as a domain concept while discovery ignores Markdown files. | PRD/TRD already include Markdown sources, but sprint requirements defer discovery. | Tests must explicitly prove `.md` source files are ignored for Sprint 3. | Sprint 5 Markdown document source implementation. |
| Dimension titles may be filename-derived instead of parsed from content. | Future listings may want richer dimension metadata. | Preserve `Dimension` fields and keep content parsing out until required. | Future dimension metadata or prompt-composition sprint. |
| Reference matching helper could become duplicated between sources and dimensions. | Source and dimension candidate rules differ but share exact/prefix mechanics. | Keep any helper private to `internal/study` and rule-specific candidate construction explicit. | Refactor only if duplication becomes error-prone after more commands use resolution. |
| Error rendering could become split between domain and app layers. | Domain needs structured errors; app needs user-facing text. | Domain returns typed/structured facts; app maps to command output once. | Revisit when more commands need shared error rendering. |

### Future Considerations

| Consideration | Deferred Until | Reason Deferred | What Should Be Preserved Now |
| --- | --- | --- | --- |
| Markdown document source discovery and applicability filtering. | Sprint 5 or revised requirements. | Explicit Sprint 3 non-goal. | `SourceKind` should not make directory-only assumptions impossible. |
| Study initialization and YAML parsing. | Study initialization sprint. | Sprint 3 inspects existing studies only. | Discovery should tolerate existing filesystem structures without requiring generated metadata. |
| Runtime execution, prompt composition, report validation, and run state. | Execution/run-loop sprints. | Listing does not run analysis or validate generated reports. | `internal/study` should remain the product module for future study behavior without adding runtime dependencies now. |
| JSON output for listing. | When CLI output foundation or automation requirements justify it. | Non-trivial JSON listing is a non-goal. | Service should return structured values so JSON can be added later without changing discovery. |
| Caching or indexing study metadata. | If direct-child discovery becomes measurably slow. | No evidence that simple listings need caching. | Keep discovery deterministic and side-effect free. |

## Final Decisions

### Decision 1: Create A Minimal `internal/study` Read Model

- **Decision:** Implement `internal/study/domain.go`, `discovery.go`, `resolve.go`, and `service.go` with `Study`, `Source`, `SourceKind`, `Dimension`, shallow discovery, reference resolution, and read-only listing use cases.
- **Rationale:** Study concepts are first-class PRD/TRD domain objects, and sprint constraints require discovery/resolution outside CLI handlers. A cohesive `internal/study` package is the smallest design that preserves module ownership and supports later study behavior.
- **Study / Source Grounding:** `technical-handbook.md` cites thin CLI/domain-owned evidence from chezmoi `main.go:26-34`, gh-cli `cmd/gh/main.go:6`, helm `pkg/cmd/install.go:132-145`, and yq `cmd/root.go:9`; `reasoning/architecture.md` adopts the same conclusion; `ARCHITECTURE.md` states product modules own product behavior and `study` owns sources and dimensions.
- **Trade-Offs Accepted:** Accept a small service boundary and several focused files now instead of a single command-handler implementation.
- **Technical Debt / Future Impact:** Avoids CLI/domain coupling; potential future growth must stay use-case-oriented and not become a generic manager.
- **Alternatives Rejected:** CLI-owned discovery/resolution because it violates sprint constraints and handbook anti-patterns; global technical packages such as `internal/discovery` or `internal/validation` because architecture requires module-owned behavior; clean-architecture subpackages because the feature is too small to justify them.
- **Contracts Satisfied:** Architecture Contract; Testing Contract; REQ-AC-01 through REQ-AC-10; sprint required outputs for `/internal/study/domain.go`, `discovery.go`, `resolve.go`, and `service.go`.
- **Evidence Required:** Unit tests for domain normalization, discovery, sorting, ignored entries, and reference resolution; architecture import review confirming platform packages do not import `internal/study`.

### Decision 2: Use Shallow Deterministic Filesystem Discovery

- **Decision:** Discover studies from direct non-hidden directories under workspace `studies/`; discover sources from direct non-hidden directories under `studies/<study>/sources/`; discover dimensions from direct Markdown files under `studies/<study>/dimensions/` whose filenames begin with numeric prefixes; sort all returned collections ascending deterministically.
- **Rationale:** Listing must be fast, safe, stable, and read-only. The sprint explicitly forbids recursive source scans and defers Markdown document source discovery.
- **Study / Source Grounding:** `technical-handbook.md` cites bounded traversal/performance evidence from gdu `pkg/analyze/parallel.go:13`, yq `pkg/yqlib/stream_evaluator.go:78-113`, and gh-cli lazy factory fields `pkg/cmdutil/factory.go:27-42`; security evidence warns about path/input trust boundaries including gdu `internal/common/ignore.go:16-37`; PRD performance requirement says listing studies should avoid recursive scans of source repositories.
- **Trade-Offs Accepted:** No caching, recursive metadata inference, or nested repository inspection. Missing `sources/` or `dimensions/` under a valid study may produce empty sections rather than fatal errors unless a real filesystem error occurs.
- **Technical Debt / Future Impact:** Future Markdown and metadata-rich listing will need to extend discovery, but current direct-child rules preserve a safe base.
- **Alternatives Rejected:** Recursive source scanning because it violates sprint acceptance and performance/security evidence; treating top-level `.md` files as sources because Sprint 3 explicitly defers Markdown document sources; requiring generated YAML metadata because listing must inspect existing studies without initialization workflows.
- **Contracts Satisfied:** Security Contract; Testing Contract; REQ-AC-01 through REQ-AC-07, REQ-AC-12, REQ-AC-13.
- **Evidence Required:** Unit tests for hidden entries, non-directory entries, ignored `.md` source files, nested content ignored, deterministic sorting, numeric dimension prefix detection, and two-digit normalization; code review confirming no recursive walk over source directories.

### Decision 3: Normalize Dimension Identity From Filenames

- **Decision:** Derive dimension `Number`, `Slug`, and `File` from direct Markdown filenames with numeric prefixes, normalizing numbers to two digits for discovered dimensions and reference matching.
- **Rationale:** Sprint 3 needs deterministic dimension listing and resolution, not full Markdown prompt metadata parsing. Filename-derived identity matches generated dimension conventions and keeps IO bounded.
- **Study / Source Grounding:** `TRD.md` defines `Dimension.Number`, `Slug`, and `File`, requires zero-padded numbers, and requires lookup by number, slug, filename, or unambiguous prefix. No external source repo provides a more relevant domain-specific filename rule; the decision is primarily grounded in project TRD and sprint requirements.
- **Trade-Offs Accepted:** Dimension title/purpose may remain empty or derived from slug rather than parsed from Markdown content.
- **Technical Debt / Future Impact:** Later prompt composition or rich listing may parse dimension content; preserving `Dimension` fields leaves room for that without changing command shape.
- **Alternatives Rejected:** Parsing Markdown headings now because it is not required and would expand scope; accepting non-numeric Markdown dimensions because acceptance requires filenames beginning with numeric prefixes; storing dimensions only as strings because resolution and future study behavior need structured domain values.
- **Contracts Satisfied:** REQ-AC-06, REQ-AC-07, REQ-AC-08; Architecture Contract; Testing Contract.
- **Evidence Required:** Unit tests for filenames like `1-foo.md`, `01-foo.md`, ignored `foo.md`, ignored nested files, normalized number `01`, slug extraction, and filename reference matching.

### Decision 4: Implement Exact-First, Unambiguous-Prefix Resolution With Structured Errors

- **Decision:** Resolve study, source, and dimension references by exact canonical matches first, then by unambiguous prefixes. Study and source refs match names; dimension refs match number, slug, full filename, exact reference aliases, and unambiguous prefixes. Missing and ambiguous results return structured domain errors that the CLI maps to actionable non-zero command errors.
- **Rationale:** Prefix resolution is required for usability, but exact-first matching and ambiguity rejection avoid surprising command behavior. Distinct error classes are required by acceptance criteria.
- **Study / Source Grounding:** `technical-handbook.md` cites actionable error evidence from go-task `errors/errors_task.go:13-32`, age `cmd/age/tui.go:37-54`, gh-cli `internal/ghcmd/cmd.go:281-301`, helm `pkg/storage/driver/driver.go:27-36`, and go-task `errors/errors.go:47-50`. `reasoning/architecture.md` specifically recommends typed/structured errors for not found and ambiguity.
- **Trade-Offs Accepted:** Resolver complexity increases slightly to maintain candidate sets and ambiguity lists.
- **Technical Debt / Future Impact:** Fuzzy matching and did-you-mean suggestions are deferred; candidate lists should stay canonical and concise to avoid noisy errors.
- **Alternatives Rejected:** Prefix-first matching because it could override exact matches; generic `not found` errors because they fail acceptance and error-contract expectations; fuzzy/typo matching because it is not required and could make resolution non-deterministic.
- **Contracts Satisfied:** Errors Contract; CLI Surface Contract; REQ-AC-08, REQ-AC-09, REQ-AC-10.
- **Evidence Required:** Unit tests for exact vs prefix precedence, missing refs, ambiguous refs, dimension number/slug/filename refs, and source prefix refs; command tests asserting non-zero actionable errors.

### Decision 5: Wire Thin Study Listing Commands In `internal/app`

- **Decision:** Add `ultraplan study list` and `ultraplan study <study> list` in `internal/app/study_commands.go`, using existing workspace discovery/global `--workspace`, delegating domain work to `internal/study`, and rendering minimal stable human-readable output including source kind.
- **Rationale:** `internal/app` owns CLI composition, argument parsing, help, workspace option consumption, and output. The command layer should not own study discovery or resolution.
- **Study / Source Grounding:** `technical-handbook.md` cites command factory/delegate evidence from gh-cli `pkg/cmd/issue/list/list.go:47-118`, restic `cmd/restic/main.go:37-114`, helm `pkg/cmd/root.go:105`, and injected output evidence from gh-cli `pkg/iostreams/iostreams.go:551-568`, go-task `executor.go:553-564`, and mitchellh-cli `ui_mock.go:27-33`. Terminal UX evidence supports calm, script-friendly listing output.
- **Trade-Offs Accepted:** Output remains minimal text rather than rich tables or JSON. The app layer formats results while the study layer remains output-agnostic.
- **Technical Debt / Future Impact:** JSON output can be added later because service returns structured values. Avoid hardcoding output paths or global writers that would make tests brittle.
- **Alternatives Rejected:** Rich TUI/color/table output because it expands scope; non-trivial JSON listing because it is a non-goal; command handlers that read directories directly because they violate architecture and testability constraints.
- **Contracts Satisfied:** CLI Surface Contract; Architecture Contract; REQ-AC-01, REQ-AC-02, REQ-AC-10, REQ-AC-11, REQ-AC-12; sprint required output for `/internal/app/study_commands.go`.
- **Evidence Required:** Command tests for `study list`, `study <study> list`, `--workspace`, output ordering, source kind display, and missing/ambiguous study errors; help/exit behavior review if command framework requires it.

### Decision 6: Verify With Focused Domain And Command Tests Only

- **Decision:** Add `internal/study/study_test.go` for domain/discovery/resolution behavior and `internal/app/study_commands_test.go` for CLI output/error behavior, using fixture/temp workspaces and no OpenCode, network, cloned repositories, or credentials.
- **Rationale:** Sprint behavior is deterministic and read-only, so normal tests should be fast, local, and fixture-based. Full runtime or integration dependencies would violate sprint constraints.
- **Study / Source Grounding:** `technical-handbook.md` cites test evidence from chezmoi `internal/cmd/main_test.go:64-174`, go-task `task_test.go:166-169`, gdu `cmd/gdu/app/app_test.go:682`, and gh-cli output test seams. Requirements explicitly say normal unit tests must not require OpenCode, network, cloned repositories, or AI provider credentials.
- **Trade-Offs Accepted:** No real-runtime or large-repository integration tests in this sprint.
- **Technical Debt / Future Impact:** Later runtime and run-loop sprints will need fake runtime and gated OpenCode tests, but listing does not.
- **Alternatives Rejected:** Manual-only verification because acceptance requires deterministic test coverage; real repository scans because they are slow and out of scope; adding filesystem interfaces solely for tests because temp directories are sufficient for shallow local reads.
- **Contracts Satisfied:** Testing Contract; Security Contract; REQ-AC-14, REQ-AC-15 plus all behavior-specific acceptance criteria.
- **Evidence Required:** `go test ./...` and `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`; review of test cases for hidden entries, ignored files, deferred Markdown sources, deterministic ordering, and errors.

## Expected Evidence

| Evidence Type | Required Evidence | Source / Command / Review Check |
| --- | --- | --- |
| Tests | Domain tests cover normalization, discovery ordering, ignored hidden/non-directory entries, ignored `.md` source files, shallow source discovery, dimension discovery, and source/dimension/study resolution. | `/home/antonioborgerees/coding/ultraplan-go/internal/study/study_test.go`; `go test ./...`. |
| Tests | Command tests cover `ultraplan study list`, `ultraplan study <study> list`, global `--workspace`, deterministic output, source kind display, and non-zero actionable missing/ambiguous errors. | `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_commands_test.go`; `go test ./...`. |
| Build | CLI builds successfully. | `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`. |
| Runtime | Optional manual fixture run demonstrates listing output for a workspace with multiple studies, sources, and dimensions. | Built binary or `go run ./cmd/ultraplan --workspace <fixture> study list`; not a substitute for tests. |
| Review | Architecture review confirms `internal/study` owns domain discovery/resolution, `internal/app` owns CLI wiring/output, workspace discovery is reused, and platform packages do not import `internal/study`. | `system/protocols/architecture-review-protocol.md`; import review. |
| Review | Sprint review confirms required outputs exist, deferred scope is excluded, and all acceptance criteria are mapped to tests or review checks. | `system/protocols/sprint-review-protocol.md`. |
| Documentation | No standalone docs are required in this sprint beyond command help/output already covered by CLI surface. | Review command help/output if existing app test conventions include help assertions. |

## Assumptions And Risks

| Item | Type | Impact | Mitigation / Follow-Up |
| --- | --- | --- | --- |
| Existing Sprint 2 workspace discovery and global `--workspace` behavior are available and reusable. | Assumption | Listing command wiring depends on this; reimplementing workspace lookup would violate scope. | Use existing app/workspace APIs; if absent or incompatible, stop and align with Sprint 2 behavior rather than inventing a second path. |
| Missing `studies/` in an otherwise valid workspace can be represented as an empty study list unless existing workspace validation treats it as invalid. | Assumption | Affects user-facing output and errors for empty workspaces. | Follow existing workspace validation first; otherwise return empty list calmly for read-only listing. |
| Missing `sources/` or `dimensions/` under a resolved study can be shown as empty sections unless filesystem errors indicate a real problem. | Assumption | Avoids failing on partially created or hand-edited studies. | Test absent directories if implementation chooses empty behavior; surface permission/read errors separately. |
| Prefix resolution could surprise users if exact matches and prefixes are not ordered consistently. | Risk | Users may get the wrong study/source/dimension or ambiguous results. | Exact match wins; otherwise require a single prefix candidate; test ambiguous cases. |
| Error output could leak absolute local paths or become too verbose. | Risk | Reduces security and usability. | Prefer workspace-relative paths and concise canonical candidate lists where practical. |
| Future Markdown source support may need to reinterpret ignored Sprint 3 `.md` files. | Risk | Users may expect PRD-wide Markdown behavior earlier than Sprint 5. | Regression test Sprint 3 deferred behavior and document in final review; implement Markdown support only when its sprint requires it. |
| Dimension filename parsing edge cases may vary across user-created files. | Risk | Some dimensions may be ignored if filenames do not start with numeric prefixes. | Acceptance criteria require numeric prefixes; ignored non-conforming files should be covered by tests and can be revisited with validation workflows later. |

## Implementation Constraints

- Keep study behavior in `internal/study`; do not create global `internal/validation`, `internal/reports`, `internal/scheduler`, `internal/discovery`, or similar technical-layer packages for Sprint 3.
- Keep CLI parsing, command registration, workspace option consumption, output rendering, and command error mapping in `internal/app`.
- Use existing workspace discovery and global `--workspace` behavior from prior sprints; do not invent a second workspace resolution mechanism.
- Do not allow listing to escape the resolved workspace for workspace-managed paths.
- Use shallow direct-child reads only for `studies/`, `sources/`, and `dimensions/`; do not recursively scan source repositories or nested source contents.
- Ignore hidden entries where required by acceptance criteria.
- Ignore top-level `.md` source files in Sprint 3; do not parse frontmatter or applicability filters.
- Normalize discovered dimension numbers to two digits.
- Resolve exact matches before prefixes and reject ambiguous prefixes with actionable errors.
- Prefer workspace-relative paths in user-facing output and errors where practical.
- Do not add OpenCode, `agentwrap`, prompt, runtime, scheduler, run state, report validation, summary, target, or sprint workflow code in this sprint.
- Tests must be deterministic and must not require network access, cloned repositories, OpenCode, AI provider credentials, or real runtime execution.

## Plan Handoff

`plan.md` must execute these decisions. It must not invent architecture, scope, or decisions beyond this document.

The plan must carry forward:

- final decisions
- applicable contracts and requirement IDs
- expected evidence
- risks and assumptions
- required review protocols

## Phase Exit Criteria

- [x] Selected context was read and used.
- [x] Area-specific reasoning documents were completed where applicable or explicitly marked not applicable.
- [x] Area-specific reasoning conclusions are reflected or explicitly overridden in final decisions.
- [x] Contracts are explicitly mapped to decisions and expected evidence.
- [x] Final decisions are clear enough for `plan.md` to execute.
- [x] Expected evidence is specific and reviewable.
