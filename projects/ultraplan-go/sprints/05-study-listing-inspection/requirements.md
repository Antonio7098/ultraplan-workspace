# Sprint Requirements: 05 Study Listing Inspection

> Project: `ultraplan-go`
> Sprint: `05-study-listing-inspection`
> Purpose: the authoritative, human-readable sprint contract. All other sprint artifacts must satisfy these requirements.

## Sprint Goal

Add top-level Markdown document source discovery, frontmatter-based dimension applicability, and visible study listing/inspection output for source kind and applicability without introducing runtime execution behavior.

## Required Outputs

| Output | Path | Description |
| --- | --- | --- |
| Sprint requirements | `projects/ultraplan-go/sprints/05-study-listing-inspection/requirements.md` | Defines the binding goal, scope, acceptance criteria, constraints, dependencies, and review expectations for this sprint. |
| Sprint index | `projects/ultraplan-go/sprints/05-study-listing-inspection/sprint-index.md` | Selects the roadmap entries, source docs, evidence reports, and contracts that govern this sprint. |
| Technical handbook | `projects/ultraplan-go/sprints/05-study-listing-inspection/technical-handbook.md` | Gives implementation guidance for Markdown source discovery, frontmatter parsing, applicability filtering, CLI output, and tests. |
| Applicability reasoning | `projects/ultraplan-go/sprints/05-study-listing-inspection/reasoning/source-applicability.md` | Records design reasoning for source kinds, frontmatter normalization, malformed frontmatter behavior, and list output semantics. |
| Consolidated reasoning | `projects/ultraplan-go/sprints/05-study-listing-inspection/reasoning.md` | Summarizes accepted sprint decisions, tradeoffs, risks, and requirement mappings. |
| Implementation plan | `projects/ultraplan-go/sprints/05-study-listing-inspection/plan.md` | Provides the ordered task plan and verification checklist for implementation. |
| Sprint review | `projects/ultraplan-go/sprints/05-study-listing-inspection/review.md` | Reviews completed implementation against these requirements and records carry-forward decisions. |
| Study source model implementation | `/home/antonioborgerees/coding/ultraplan-go/internal/study/domain.go` | Extends or confirms the `Source` model with source kind, normalized applicable dimensions, and frontmatter metadata needed for listing and filtering. |
| Study source discovery implementation | `/home/antonioborgerees/coding/ultraplan-go/internal/study/discovery.go` | Discovers non-hidden source directories and top-level Markdown document sources under `studies/<study>/sources/` deterministically. |
| Markdown frontmatter implementation | `/home/antonioborgerees/coding/ultraplan-go/internal/study/markdown.go` | Parses leading YAML frontmatter, strips leading frontmatter, normalizes `applicable_dimensions`, and reports invalid values with source path context. |
| Applicability helper implementation | `/home/antonioborgerees/coding/ultraplan-go/internal/study/applicability.go` | Provides `GetApplicableSources(sources []Source, dimension Dimension) []Source` or an equivalent exported study helper. |
| Study listing CLI implementation | `/home/antonioborgerees/coding/ultraplan-go/internal/study/cli.go` | Updates `ultraplan study <study> list` output to show source kind and Markdown applicability filters in deterministic text output. |
| Study source tests | `/home/antonioborgerees/coding/ultraplan-go/internal/study/source_test.go` | Tests discovery, source kind assignment, deterministic ordering, ignored entries, frontmatter parsing/stripping, normalization, and applicability filtering. |
| Study listing command tests | `/home/antonioborgerees/coding/ultraplan-go/internal/study/cli_test.go` | Tests human-facing listing output for directory sources, Markdown document sources, unfiltered Markdown sources, and filtered Markdown sources. |

## Acceptance Criteria

- [ ] `ultraplan study <study> list` discovers and displays both source directories and top-level `.md` document sources under `studies/<study>/sources/`.
- [ ] Source listing identifies each source as `directory` or `markdown` in deterministic order.
- [ ] Hidden entries, non-directory non-Markdown files, and nested Markdown files inside directory sources are ignored by source discovery.
- [ ] Markdown frontmatter is parsed only when it is the leading document block delimited by `---`.
- [ ] Markdown files without leading frontmatter return an empty frontmatter map and preserve the original body for stripping behavior.
- [ ] `stripFrontmatter` removes only leading frontmatter and returns the document body without the frontmatter block.
- [ ] `applicable_dimensions` accepts numeric and string values and normalizes `1`, `"1"`, and `"01"` to dimension `01`.
- [ ] Missing or empty `applicable_dimensions` means a Markdown source applies to every dimension.
- [ ] Invalid `applicable_dimensions` values are reported with the source path and offending value.
- [ ] `GetApplicableSources` or the equivalent study helper always includes directory sources.
- [ ] `GetApplicableSources` or the equivalent study helper includes unfiltered Markdown sources for every dimension.
- [ ] `GetApplicableSources` or the equivalent study helper includes filtered Markdown sources only when the selected normalized dimension number matches.
- [ ] Inapplicable source/dimension pairs are represented as skipped or excluded from listing-derived applicability checks, not as failures.
- [ ] Existing `study list` and `study <study> list` behavior from Sprint 3 remains deterministic and script-friendly.
- [ ] Unit and command tests cover valid frontmatter, absent frontmatter, malformed frontmatter, invalid applicability values, hidden entries, nested Markdown files, and mixed directory/Markdown source ordering.
- [ ] `go test ./...` passes in `/home/antonioborgerees/coding/ultraplan-go`.

## Non-Goals

- Runtime execution for Markdown document sources is not included.
- Prompt composition for Markdown document analysis is not included.
- Report validation changes for Markdown document sources are not included.
- Summary generation changes for inapplicable pairs are not included.
- `run`, `run-all`, `run-loop`, `status`, and synthesis scheduling changes are not included beyond preserving compile-time compatibility.
- Source cloning, assisted study initialization, and YAML schema expansion are not included.
- JSON output schema stabilization for listing commands is not included unless already present from prior work.
- Target, sprint planning, and sprint execution workflows remain out of scope.

## Constraints

- Keep study behavior in `internal/study`; do not create global `validation`, `scheduler`, `reports`, or `prompts` packages for this sprint.
- Preserve the dependency direction from `docs/ARCHITECTURE.md`: `study` may use `workspace` and platform helpers, but platform packages must not import `study`.
- Apply the roadmap Skeleton / Local CLI Gate only; runtime/provider and durable workflow gate requirements do not apply unless existing code paths already require them.
- Source discovery must not recursively scan source repositories for listing; it may inspect only direct children of `studies/<study>/sources/`.
- Markdown document sources must be top-level `.md` files directly under a study's `sources/` directory.
- Directory sources are always applicable and must not depend on frontmatter.
- Normalized dimension numbers must be two digits for this sprint's matching behavior.
- Listing output must remain deterministic and suitable for command tests.
- Errors must preserve cause chains when translated to CLI errors and must not swallow filesystem or malformed frontmatter failures.
- Tests must use local filesystem fixtures or temporary directories and must not require OpenCode, `agentwrap`, network access, or real provider credentials.

## Dependencies

| Prior Sprint / Output | Required For | Notes |
| --- | --- | --- |
| Sprint 1: Go Module, CLI Shell, and App Composition | Buildable CLI and package structure | Required for command wiring, test execution, and `cmd/ultraplan` build continuity. |
| Sprint 2: Workspace, Config, Logging, and Health Skeleton | Workspace discovery and path safety | Required so study listing can resolve the workspace root and operate within workspace-managed paths. |
| Sprint 3: Study Domain, Listing, and Resolution | Existing study model, source/dimension lookup, and listing commands | This sprint extends source discovery/listing rather than replacing it. |
| Sprint 4: Study Initialization From YAML | Generated study directory layout | Required for the expected `studies/<study>/sources/` and `dimensions/` structures used by fixtures and listing. |
| `projects/ultraplan-go/docs/PRD.md` | Product behavior for study listing and Markdown document sources | Governs user-facing behavior and non-goals. |
| `projects/ultraplan-go/docs/TRD.md` | Technical source model, discovery, frontmatter, and applicability rules | Governs helper behavior, normalization, and tests. |
| `projects/ultraplan-go/docs/ARCHITECTURE.md` | Package ownership and dependency direction | Governs where implementation belongs. |

## Review Expectations

| What | How Verified |
| --- | --- |
| Scope alignment | Review `requirements.md`, `sprint-index.md`, `technical-handbook.md`, `reasoning.md`, and `plan.md` against roadmap Sprint 5 and confirm no runtime, prompt, validation, summary, run-loop, target, or sprint-execution scope was added. |
| Source discovery behavior | Run or inspect tests proving direct child directories and top-level `.md` files are discovered while hidden entries, unrelated files, and nested Markdown files are ignored. |
| Frontmatter behavior | Run or inspect tests for leading-only frontmatter parsing, stripping, absent frontmatter, malformed frontmatter, and invalid `applicable_dimensions` diagnostics. |
| Applicability filtering | Run or inspect tests proving directory sources always apply, unfiltered Markdown applies to all dimensions, filtered Markdown applies only to matching normalized dimensions, and inapplicable pairs are skipped or excluded rather than failed. |
| CLI listing behavior | Run or inspect command tests for `ultraplan study <study> list` output containing source kind and applicability details in deterministic order. |
| Architecture conformance | Review implementation paths and imports to confirm study-owned behavior remains in `internal/study` and platform packages do not import product modules. |
| Verification commands | Run `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go` and record the result in `review.md`. |
