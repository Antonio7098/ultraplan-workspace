# Sprint Reasoning: 05 Study Listing Inspection

> Project: `ultraplan-go`
> Sprint: `05-study-listing-inspection`
> Output: `projects/ultraplan-go/sprints/05-study-listing-inspection/reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/05-study-listing-inspection/requirements.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/sprints/05-study-listing-inspection/sprint-index.md`, `projects/ultraplan-go/sprints/05-study-listing-inspection/technical-handbook.md`, `templates/sprint-reasoning.md`

This document decides. It synthesizes selected context, handbook evidence, and contracts into final sprint decisions.

It does not replace `sprint-index.md`, `technical-handbook.md`, or any area-specific reasoning artifact.

## Sprint Purpose

- **Goal:** Add top-level Markdown document source discovery, frontmatter-based dimension applicability, and visible study listing output for source kind and applicability while preserving deterministic study inspection behavior.
- **Non-Goals:** Runtime execution for Markdown document sources, prompt composition changes, report validation changes, summary generation changes, run/run-all/run-loop/status/synthesis scheduling changes beyond compile-time compatibility, source cloning, assisted initialization, YAML schema expansion, JSON schema stabilization, target workflows, sprint planning, and sprint execution.
- **Depends On:** Prior CLI shell, workspace discovery, study listing/resolution, and study initialization layout from earlier sprints; project source documents `PRD.md`, `TRD.md`, and `ARCHITECTURE.md`; selected contracts from `sprint-index.md`.

## Selected Context And Pre-Reasoning Artifacts

| Artifact | Path | How It Was Used |
| --- | --- | --- |
| Sprint Requirements | `projects/ultraplan-go/sprints/05-study-listing-inspection/requirements.md` | Binding sprint scope, acceptance criteria, non-goals, constraints, dependencies, output paths, and review expectations. |
| Project Index | `projects/ultraplan-go/project-index.md` | Confirmed available source documents, selected study evidence pool, active contracts, and absence of prior decision artifacts. |
| Product Requirements | `projects/ultraplan-go/docs/PRD.md` | Established product need for source listing, Markdown document sources, `applicable_dimensions`, stable script-friendly output, and deferred runtime/prompt/report behavior. |
| Technical Requirements | `projects/ultraplan-go/docs/TRD.md` | Defined source model shape, direct-child source discovery, frontmatter parsing/stripping, normalization, applicability helper behavior, and future flows that must eventually respect applicability. |
| Architecture | `projects/ultraplan-go/docs/ARCHITECTURE.md` | Constrained ownership to `internal/study`, preserved dependency direction, and rejected global technical-layer packages for study-specific behavior. |
| Sprint Index | `projects/ultraplan-go/sprints/05-study-listing-inspection/sprint-index.md` | Selected contracts, evidence reports, review protocols, non-goals, and excluded contexts. It also selected area reasoning for source applicability, but no such artifact exists in the workspace at the time of this document. |
| Technical Handbook | `projects/ultraplan-go/sprints/05-study-listing-inspection/technical-handbook.md` | Supplied evidence from Go CLI studies for thin CLI design, command factories, IO test seams, error wrapping, deterministic fixture/output tests, direct-child scanning, trust boundaries, and stable text output. |

## Area-Specific Reasoning Inputs

The sprint index selected an Architecture reasoning artifact at `projects/ultraplan-go/sprints/05-study-listing-inspection/reasoning/source-applicability.md`, but the `reasoning/` directory and area-specific reasoning files are not present. No area-specific conclusions were available to summarize.

| Area | Reasoning Document | Key Conclusion | Evidence Basis | Impact On Final Decision |
| --- | --- | --- | --- | --- |
| Source applicability | `projects/ultraplan-go/sprints/05-study-listing-inspection/reasoning/source-applicability.md` | Not available in the workspace. Final decisions in this document resolve source kind, frontmatter, applicability, malformed metadata, ordering, and listing semantics directly from requirements, docs, sprint index, and technical handbook evidence. | `requirements.md`, `PRD.md`, `TRD.md`, `ARCHITECTURE.md`, `technical-handbook.md` | `plan.md` can proceed from the final decisions below without reopening architecture. |

## Sprint Technical Handbook Summary

- **Relevant Patterns:** Keep CLI wrappers thin and put behavior in `internal/study`; use command construction and output seams that remain testable; wrap errors with `%w` and source path context; rely on local fixtures and command-output assertions; inspect only direct children; treat Markdown/frontmatter as untrusted user-authored input; keep terminal output plain and stable.
- **Important Trade-Offs:** Concrete study-owned filesystem discovery is acceptable for this narrow local scope; parsing Markdown frontmatter during discovery creates a larger listing error surface but enables consistent source metadata; deterministic text output is favored over a premature JSON schema; directory and Markdown sources use asymmetric applicability rules; leading-only frontmatter parsing is strict but predictable.
- **Warnings / Anti-Patterns:** Do not make `study <study> list` a large command handler; do not recursively scan source repositories; do not drop underlying YAML/filesystem causes; do not hardcode uncapturable stdout/stderr paths; do not add runtime, prompt, scheduler, report-validation, or summary behavior; do not trust malformed metadata.
- **Evidence Confidence:** High. The handbook selected directly relevant Go CLI reports for package structure, command architecture, IO seams, error handling, testing, security, performance, and maintainability. Terminal UX confidence is medium but adequate because the sprint only needs stable non-interactive listing output.

## Contracts Applied

| Contract / Requirement ID | Constraint | Decision Impact | Expected Evidence |
| --- | --- | --- | --- |
| Architecture contract; `ARCHITECTURE.md` module ownership | Study behavior belongs in `internal/study`; platform packages must not import product modules; avoid global `validation`, `scheduler`, `reports`, or `prompts` packages. | Discovery, frontmatter, applicability, and listing decisions are implemented as focused files in `internal/study`. | Architecture review confirms import direction and file placement. |
| Errors contract; `requirements.md` AC-09, C-10 | Invalid applicability and malformed frontmatter diagnostics must include source path/offending value and preserve causes. | Markdown parsing and normalization failures fail listing/discovery with wrapped errors. | Unit/command tests assert path/value context; review checks `%w` cause preservation. |
| Security contract; `requirements.md` C-11, C-12 | User-authored source files are untrusted; workspace path safety must be preserved; no runtime behavior is introduced. | Only direct children under `sources/` are inspected; hidden entries and nested Markdown are ignored; no external execution or recursive traversal. | Tests for ignored entries; architecture review confirms no runtime calls. |
| Testing contract; `requirements.md` AC-16, AC-17 | Deterministic unit and command tests must cover discovery, frontmatter, filtering, listing, and `go test ./...`. | Public behavior tests are required in `internal/study/source_test.go` and `internal/study/cli_test.go`. | `go test ./...` passes from `/home/antonioborgerees/coding/ultraplan-go`. |
| CLI Surface contract; `requirements.md` AC-01, AC-02, AC-15 | Listing remains deterministic, script-friendly, and human-readable. | Text output includes source kind and Markdown applicability annotations without stabilizing a new JSON schema. | Command tests compare deterministic output. |
| `PRD.md` FR-4, FR-19; `TRD.md` 9A | Markdown document sources must be direct top-level `.md` files; `applicable_dimensions` normalizes to two digits; missing/empty filter means all dimensions. | Source model gains kind, frontmatter metadata, and normalized applicability. | Unit tests for numeric/string normalization and default applicability. |
| `requirements.md` Non-Goals | Runtime, prompt, validation, summary, and workflow behavior are out of scope. | Helper shape may support future flows, but this sprint does not schedule, execute, validate reports, or synthesize. | Review confirms no added runtime/scheduler/report behavior beyond compile compatibility. |

## Repos Studied / Source Evidence Used

| Source / Repo / Report | Concrete Reference | Relevant Finding | Why It Matters For This Sprint | Used In Decision(s) |
| --- | --- | --- | --- | --- |
| `go-cli-study / 01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md`; `cmd/root.go:112-127`, `cmd/restic/main.go:37-114`, `internal/restic/repository.go:18` | Mature Go CLIs keep entrypoints thin and business logic in `internal/` domain packages. | Supports keeping discovery, parsing, filtering, and list formatting in `internal/study`, not `cmd` or global packages. | Decisions 1, 2, 5, 6 |
| `go-cli-study / 02-command-architecture` | `studies/go-cli-study/reports/final/02-command-architecture.md`; `pkg/cmd/issue/list/list.go:47-118`, `pkg/cmd/install.go:132-145`, `opencode/cmd/root.go:49-183` | Command factories and option structs keep command handlers from absorbing business logic; large `RunE` handlers are a warning. | Shapes `study <study> list` as orchestration/presentation only. | Decision 6 |
| `go-cli-study / 05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md`; `rclone/vfs/zip.go:21`, `helm/pkg/storage/driver/secrets.go:76`, `age/cmd/age/tui.go:37-54` | Wrapped errors and actionable diagnostics preserve cause chains while giving users context. | Supports failing invalid/malformed frontmatter with source path and offending value instead of silently ignoring it. | Decisions 3, 4, 7 |
| `go-cli-study / 06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md`; `pkg/iostreams/iostreams.go:551-568`, `internal/ui/terminal.go:10-36`, `ui.go:19-43` | Capturable stdout/stderr and filesystem seams make CLI output testable. | Listing output must remain observable in command tests. | Decisions 6, 8 |
| `go-cli-study / 09-terminal-ux` | `studies/go-cli-study/reports/final/09-terminal-ux.md`; `internal/cmd/prompt.go:20-256`, `cmd/evaluate_sequence_command.go:67` | Stable non-interactive output is preferable for CLI-first inspection commands. | Supports plain deterministic text annotations rather than TUI/progress behavior. | Decision 6 |
| `go-cli-study / 11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md`; `chezmoi/internal/cmd/main_test.go:64-174`, `helm/internal/test/test.go:43`, `go-task/task_test.go:166-169` | Fixture-driven and golden/output tests provide high confidence for CLI behavior. | Requires local filesystem fixtures and command output assertions. | Decision 8 |
| `go-cli-study / 13-security` | `studies/go-cli-study/reports/final/13-security.md`; `internal/config/json/validator.go:146`, `pkg/registry/transport.go:37-41`, `pkg/surveyext/editor_manual.go:23` | Treat paths and user-authored data as trust boundaries; validate inputs; avoid unsafe execution. | Supports direct-child-only discovery and strict frontmatter validation. | Decisions 1, 3, 4, 7 |
| `go-cli-study / 14-performance` | `studies/go-cli-study/reports/final/14-performance.md`; `pkg/cmdutil/factory.go:27-42`, `yq/pkg/yqlib/stream_evaluator.go:78-113`, `gdu/pkg/analyze/parallel.go:13` | Fast CLIs avoid recursive/unbounded scans and defer expensive work. | Supports only inspecting direct children of `sources/` and ignoring nested Markdown. | Decision 1 |
| `go-cli-study / 15-philosophy` | `studies/go-cli-study/reports/final/15-philosophy.md`; `VISION.md:97`, `pkg/cmd/factory/default.go:26-46`, `internal/chezmoi/system.go:25-45` | Non-goals and interfaces should be deliberate; complexity should not be accumulated early. | Supports avoiding runtime/prompt/validation/scheduler expansion and avoiding unnecessary abstraction. | Decisions 1, 5, 6, 8 |

## Trade-Off And Debt Analysis

### Accepted Trade-Offs

| Trade-Off | Benefit | Cost / Constraint Accepted | Why Acceptable Now | Revisit Trigger |
| --- | --- | --- | --- | --- |
| Parse frontmatter during source discovery | Source list exposes kind/applicability consistently and invalid metadata is caught before later flows depend on it. | `study <study> list` can fail on malformed Markdown metadata. | The sprint explicitly requires list-visible applicability and invalid value diagnostics. | Future validation mode distinguishes warnings from hard errors per command context. |
| Use direct-child filesystem inspection only | Fast, deterministic, safe, and aligned with performance/security evidence. | Nested Markdown files cannot become standalone Markdown sources. | Requirements explicitly say nested Markdown belongs to directory source content. | Product adds recursive source discovery or explicit source manifests. |
| Keep directory sources always applicable | Preserves existing repository source behavior and avoids surprise filtering. | Source applicability logic is asymmetric by source kind. | PRD/TRD define applicability only for Markdown document sources. | Future source kinds need applicability metadata. |
| Fail invalid `applicable_dimensions` values instead of ignoring them | Prevents hidden task-matrix mistakes and surfaces bad user metadata. | A single bad Markdown file can block source listing. | Requirements demand invalid values be reported with path and offending value. | A future `validate` command introduces warning-only inspection modes. |
| Stabilize text output, not JSON schema | Delivers script-friendly visible behavior without premature machine-readable compatibility promises. | Scripts parsing text still depend on formatting. | JSON output schema stabilization is a non-goal. | A future sprint selects the CLI Surface contract for stable JSON list output. |
| Use concrete study helpers instead of broad interfaces | Keeps implementation small and consistent with current local filesystem scope. | Some filesystem failure modes are tested through temp dirs rather than fakes. | Architecture permits concrete local collaborators for narrow side effects. | New volatile storage, remote sources, or reusable FS abstraction enters scope. |

### Potential Technical Debt

| Debt / Shortcut | Why It Might Accrue | Current Mitigation | Owner / Follow-Up |
| --- | --- | --- | --- |
| Text listing format may become an implicit API | Scripts may rely on exact columns/phrasing before JSON is stabilized. | Command tests lock current deterministic output; docs/review note JSON is not stabilized. | Future CLI output/versioning sprint. |
| Frontmatter parser/helper names may remain unexported or narrow | Future prompt composition needs stripped body and future schedulers need applicability. | Implement helpers in `internal/study` with clear tests and an exported applicability helper equivalent to `GetApplicableSources`. | Future prompt/runtime sprint reuses or promotes helpers if needed. |
| Duplicate normalized dimension values could hide user mistakes | Silently deduplicating duplicates may mask redundant metadata. | Decide to deduplicate to deterministic normalized lists while preserving invalid-value errors. | Future validation sprint may add warnings for duplicates. |
| Malformed frontmatter hard-fails listing | Users may want to inspect other valid sources despite one bad file. | Error includes source path and cause; validation mode can later report multiple issues. | Future study validation command refinement. |

### Future Considerations

| Consideration | Deferred Until | Reason Deferred | What Should Be Preserved Now |
| --- | --- | --- | --- |
| Markdown prompt composition with stripped document body | Future prompt/runtime sprint | Runtime execution and prompt composition are explicit non-goals. | `stripFrontmatter` behavior and tests. |
| Run/run-all/run-loop matrix filtering | Future execution/scheduler sprint | Scheduling and durable workflows are out of scope. | `GetApplicableSources` helper semantics. |
| Report validation differences for Markdown sources | Future validation sprint | Report validation changes are non-goals. | Source kind stored on `Source`. |
| Summary representation for inapplicable pairs | Future summary sprint | Summary generation changes are non-goals. | Applicability helper returns exclusion/skipped semantics. |
| Stable JSON output for listing | Future CLI/API compatibility sprint | JSON schema stabilization is a non-goal. | Deterministic source model and text output tests. |

## Final Decisions

### Decision 1: Discover Only Direct Child Directory And Markdown Sources

- **Decision:** Extend source discovery in `internal/study` to return one `directory` source for each non-hidden direct child directory and one `markdown` source for each non-hidden top-level `.md` file directly under `studies/<study>/sources/`. Ignore hidden entries, non-directory non-Markdown files, and nested Markdown files inside source directories. Directory source names remain the directory basename. Markdown source names use the Markdown filename including `.md` so listing and prefix lookup preserve the exact direct child and avoid collisions with same-basename directories. Sort discovered sources deterministically by source name, with kind/path as deterministic tie-breakers if the current model permits duplicates during discovery.
- **Rationale:** This directly satisfies the sprint goal while preserving Sprint 3 listing behavior and avoiding recursive repository scans. Markdown document sources are first-class sources only when they are top-level files.
- **Study / Source Grounding:** `technical-handbook.md` relevant patterns for study-owned behavior and direct-child discovery; `go-cli-study / 14-performance` warns against recursive/unbounded scans; `go-cli-study / 13-security` supports explicit path trust boundaries; `TRD.md` 9A.1 defines direct-child directory and Markdown discovery.
- **Trade-Offs Accepted:** Nested Markdown files are not surfaced independently. Discovery remains concrete and local rather than a generalized scanning framework. Markdown names include the file extension, which is slightly more verbose but avoids ambiguity with directory sources.
- **Technical Debt / Future Impact:** Future recursive or manifest-based source discovery would need a new decision because this sprint intentionally makes direct-child discovery the invariant.
- **Alternatives Rejected:** Recursive discovery of all Markdown files was rejected because it violates requirements, risks scanning large repositories, and blurs document sources with repository content. Treating only directories as sources was rejected because PRD/TRD and acceptance criteria require top-level Markdown document sources. Stripping `.md` from Markdown source names was rejected because the existing source name rule is direct-child basename and stripping can collide with a same-named directory.
- **Contracts Satisfied:** Architecture, Security, Testing, CLI Surface; `requirements.md` AC-01, AC-02, AC-03, AC-15, AC-16; `TRD.md` 9A.1; `PRD.md` FR-4 and FR-19.
- **Evidence Required:** Unit tests in `internal/study/source_test.go` for mixed entries, hidden entries, non-Markdown files, nested Markdown files, kind assignment, and ordering; command tests showing listed directory and Markdown sources.

### Decision 2: Model Source Kind, Frontmatter, And Normalized Applicability In Study Domain

- **Decision:** Extend or confirm `Source` in `internal/study/domain.go` with `Kind`, normalized `ApplicableDimensions []string`, and `Frontmatter map[string]any` or equivalent fields. Directory sources use `SourceKindDirectory` and have no applicability filter. Markdown sources use `SourceKindMarkdown`; empty `ApplicableDimensions` means all dimensions.
- **Rationale:** The source model must carry enough metadata for list output now and for future scheduling/prompt/report flows without re-reading or rediscovering metadata inconsistently.
- **Study / Source Grounding:** `TRD.md` 8.2 and 9A define `SourceKind`, `ApplicableDimensions`, `Frontmatter`, and normalized two-digit dimensions. `technical-handbook.md` supports keeping this in `internal/study`; `go-cli-study / 01-project-structure` supports domain-owned behavior.
- **Trade-Offs Accepted:** Directory and Markdown sources have asymmetric metadata. This is acceptable because applicability is only defined for Markdown document sources.
- **Technical Debt / Future Impact:** Future source kinds may need a broader applicability model, but this sprint should not introduce abstraction before it is needed.
- **Alternatives Rejected:** A separate global metadata registry was rejected because architecture requires behavior near study state and no global product packages. Applying frontmatter to directory sources was rejected because directory sources must always apply and should not depend on metadata files.
- **Contracts Satisfied:** Architecture, Security, Testing; `requirements.md` AC-02, AC-08, AC-11, AC-12, AC-13; `TRD.md` 8.2, 9A.2, 9A.3; `PRD.md` FR-19.
- **Evidence Required:** Unit tests assert kind values, empty Markdown applicability meaning all dimensions, directory applicability regardless of metadata, and deterministic model values returned by discovery.

### Decision 3: Parse And Strip Only Leading YAML Frontmatter

- **Decision:** Implement `parseFrontmatter(content string) (map[string]any, error)` or equivalent and `stripFrontmatter(content string) string` in `internal/study/markdown.go`. Parse only a leading `---` block delimited by a closing `---`; documents without leading frontmatter return an empty map and unchanged body; stripping removes only that leading block and returns the body. Malformed leading frontmatter returns an error with source path context at call sites.
- **Rationale:** Leading-only frontmatter avoids treating examples or body content as metadata and gives prompt composition a predictable future stripping behavior without implementing prompt composition now.
- **Study / Source Grounding:** `TRD.md` 9A.2 defines helper behavior; `requirements.md` AC-04 through AC-07 define leading-only parsing and stripping; `technical-handbook.md` trade-off on leading-only parsing supports strict, predictable body behavior; `go-cli-study / 13-security` supports treating user-authored files as untrusted.
- **Trade-Offs Accepted:** Non-leading metadata-like blocks are ignored even if users intended them as metadata. Malformed leading blocks fail rather than becoming partial metadata.
- **Technical Debt / Future Impact:** Future validation may need richer diagnostics or multiple-error accumulation, but the strict helper behavior is stable for prompts and scheduling.
- **Alternatives Rejected:** Scanning for any `---` block anywhere in the document was rejected because it could misclassify examples/content and breaks predictable body stripping. Ignoring malformed frontmatter was rejected because it would hide bad applicability metadata and conflict with error requirements.
- **Contracts Satisfied:** Errors, Security, Testing; `requirements.md` AC-04, AC-05, AC-06, AC-07, AC-09, AC-16; `TRD.md` 9A.2.
- **Evidence Required:** Unit tests for leading frontmatter, absent frontmatter, malformed leading frontmatter, non-leading delimiter content, and stripping behavior.

### Decision 4: Normalize `applicable_dimensions` Strictly And Deduplicate Deterministically

- **Decision:** Accept `applicable_dimensions` as a YAML list containing numeric or string values. Normalize valid values such as `1`, `"1"`, and `"01"` to two-digit strings such as `01`. Missing or empty lists mean all dimensions. Invalid values fail with source path and offending value. Duplicate normalized values are deduplicated in first-seen or sorted deterministic order for stable output.
- **Rationale:** Normalized matching is required now for listing-derived applicability and future task matrix filtering. Deduplication avoids repeated display output while not changing the semantic set.
- **Study / Source Grounding:** `requirements.md` AC-08 through AC-10 define accepted value forms and invalid diagnostics; `TRD.md` 9A.2 defines normalization/default rules; `technical-handbook.md` highlights validation of untrusted metadata and open question resolution for duplicates.
- **Trade-Offs Accepted:** Duplicate values do not currently produce warnings. The sprint emphasizes matching behavior and deterministic output, not schema linting.
- **Technical Debt / Future Impact:** A future study validation command may add duplicate warnings or stricter schema reporting without changing the normalized applicability semantics.
- **Alternatives Rejected:** Preserving raw values was rejected because matching `1` and `01` would be unreliable. Treating missing applicability as applying to no dimensions was rejected because PRD/TRD require missing or empty values to mean all dimensions. Warning-only invalid values were rejected because the acceptance criteria require invalid values to be reported with path and value and because hidden invalid filters could create incorrect task matrices later.
- **Contracts Satisfied:** Errors, Security, Testing, CLI Surface; `requirements.md` AC-08, AC-09, AC-10, AC-12, AC-13, AC-16; `TRD.md` 8.2, 9A.2, 9A.3; `PRD.md` FR-19.
- **Evidence Required:** Unit tests for numeric values, string values, zero-padded strings, missing field, empty list, duplicate values, non-list values, invalid strings, negative/zero/non-integer values if parser encounters them, and error messages containing source path plus offending value.

### Decision 5: Provide A Study-Owned Applicability Helper

- **Decision:** Add `GetApplicableSources(sources []Source, dimension Dimension) []Source` or an equivalent exported study helper in `internal/study/applicability.go`. It always includes directory sources, includes unfiltered Markdown sources for every dimension, includes filtered Markdown sources only when `dimension.Number` matches a normalized applicability value, and represents inapplicable pairs by exclusion from the returned list.
- **Rationale:** The helper makes current listing-derived checks and future scheduling/validation/summary flows use the same rule without adding runtime behavior in this sprint.
- **Study / Source Grounding:** `TRD.md` 9A.3 specifies the helper and filtering semantics; `requirements.md` AC-11 through AC-14 require directory inclusion, unfiltered Markdown inclusion, filtered matching, and skipped/excluded inapplicable pairs. `technical-handbook.md` supports a concrete study-owned helper rather than new global packages.
- **Trade-Offs Accepted:** Exclusion does not record a durable skipped task in this sprint because task state and runtime flows are out of scope.
- **Technical Debt / Future Impact:** Future run-loop and summary work must call this helper or preserve identical semantics. If future flows need skip records, they can layer that behavior above the helper.
- **Alternatives Rejected:** Returning all sources with an `Applicable` flag was rejected for this sprint because requirements allow skipped/excluded listing-derived applicability checks and no task matrix is implemented now. Making the helper unexported only was rejected because future flows and tests need a stable study-level semantic surface.
- **Contracts Satisfied:** Architecture, Testing; `requirements.md` AC-11, AC-12, AC-13, AC-14; `TRD.md` 9A.3; `PRD.md` FR-19.
- **Evidence Required:** Unit tests for directory sources, unfiltered Markdown, filtered Markdown matching and non-matching dimensions, deterministic returned ordering, and exclusion of inapplicable pairs without failure.

### Decision 6: Keep Listing Plain, Deterministic, And Capturable

- **Decision:** Update `ultraplan study <study> list` text output to show each source's kind and Markdown applicability filter in deterministic order. Directory sources display as `directory` with no applicability annotation. Markdown sources display as `markdown all` when missing/empty `applicable_dimensions`, or `markdown 01,03` using comma-separated normalized dimensions when filtered. No new JSON schema is introduced.
- **Rationale:** Listing is the visible behavior for this sprint. It must remain human-readable, script-friendly, and testable without interactive output or schema expansion.
- **Study / Source Grounding:** `PRD.md` FR-4 requires source kind and Markdown applicability exposure; `requirements.md` AC-01, AC-02, and AC-15 require deterministic script-friendly listing; `technical-handbook.md` cites terminal UX evidence for calm stable output and IO abstraction/testing evidence for capturable output.
- **Trade-Offs Accepted:** Text formatting becomes test-pinned. The `all` marker is intentionally compact and script-friendly, but future user documentation may want more descriptive prose. JSON consumers must wait for a future schema-stabilization sprint.
- **Technical Debt / Future Impact:** Future JSON output should be added deliberately without breaking current text behavior unless a compatibility decision is made.
- **Alternatives Rejected:** Adding a TUI/progress-style output was rejected because listing is non-interactive and must be script-friendly. Adding/stabilizing JSON output now was rejected because the sprint explicitly says JSON schema stabilization is a non-goal.
- **Contracts Satisfied:** CLI Surface, Testing, Architecture; `requirements.md` AC-01, AC-02, AC-15, AC-16; `PRD.md` FR-4; `TRD.md` 7.1 and 9A.
- **Evidence Required:** Command tests in `internal/study/cli_test.go` for directory sources, Markdown sources with all dimensions, Markdown sources with filters, deterministic mixed ordering, and captured output.

### Decision 7: Keep Errors Actionable And Cause-Preserving

- **Decision:** Filesystem, YAML/frontmatter, and applicability normalization errors are wrapped with operation/source context using cause-preserving errors. Invalid applicability diagnostics include the Markdown source path and offending value. CLI translation preserves the chain and returns non-zero failure for malformed metadata encountered during listing/discovery.
- **Rationale:** Bad frontmatter is user-authored input. Users need a precise path and value to repair it, while tests and future callers need preserved causes for classification.
- **Study / Source Grounding:** `technical-handbook.md` error-wrapping pattern and anti-pattern warnings; `go-cli-study / 05-error-handling` evidence for `%w`, typed/sentinel handling, and actionable diagnostics; `requirements.md` C-10 and AC-09.
- **Trade-Offs Accepted:** Listing can fail early rather than partially listing valid sources around one invalid Markdown file.
- **Technical Debt / Future Impact:** Future validation commands may aggregate multiple source errors, but this sprint only needs clear failure for command/listing behavior.
- **Alternatives Rejected:** String-only errors were rejected because they break cause preservation. Silent fallback to all-dimensions on invalid metadata was rejected because it can execute or list future applicability incorrectly.
- **Contracts Satisfied:** Errors, Security, CLI Surface, Testing; `requirements.md` AC-09, AC-16, C-10; `TRD.md` 9A.2 and 14.2.
- **Evidence Required:** Tests assert invalid values include path/value, malformed frontmatter preserves underlying parse error when applicable, and CLI output exposes actionable context.

### Decision 8: Verify Through Public Behavior And Local Fixtures

- **Decision:** Add or update tests in `internal/study/source_test.go` and `internal/study/cli_test.go` using local temporary directories/fixtures. Tests cover discovery, source kind assignment, ignored entries, frontmatter parsing/stripping, normalization, applicability filtering, command output, and full `go test ./...` compatibility. No tests require OpenCode, `agentwrap`, network access, or real credentials.
- **Rationale:** The sprint is local inspection behavior. Public behavior tests give confidence without introducing runtime integration requirements.
- **Study / Source Grounding:** `technical-handbook.md` testing pattern; `go-cli-study / 11-testing-strategy` examples for fixture-driven and output regression tests; `requirements.md` testing constraints and review expectations.
- **Trade-Offs Accepted:** Some private helper behavior is tested only where it is part of the study package contract; command output is pinned intentionally.
- **Technical Debt / Future Impact:** Future runtime and scheduler tests will need fake runtime coverage, but not in this sprint.
- **Alternatives Rejected:** Real-runtime or network-backed tests were rejected because scope and constraints prohibit OpenCode/network/credentials. Brittle tests of private implementation internals were rejected in favor of source lists, errors, helper results, and rendered output.
- **Contracts Satisfied:** Testing, Security, CLI Surface; `requirements.md` AC-16, AC-17, C-12; `TRD.md` 23.1 for Markdown discovery/frontmatter/applicability unit coverage.
- **Evidence Required:** `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go`; review of test files for required cases; no runtime/network dependency in normal test path.

## Expected Evidence

| Evidence Type | Required Evidence | Source / Command / Review Check |
| --- | --- | --- |
| Tests | Discovery covers direct directories, top-level Markdown files, hidden entries, non-Markdown files, nested Markdown files, deterministic ordering, and source kind assignment. | `/home/antonioborgerees/coding/ultraplan-go/internal/study/source_test.go` |
| Tests | Frontmatter parsing/stripping covers leading block, absent frontmatter, malformed frontmatter, non-leading delimiters, numeric/string normalization, missing/empty applicability, duplicates, and invalid values. | `/home/antonioborgerees/coding/ultraplan-go/internal/study/source_test.go` or dedicated Markdown tests in `internal/study` |
| Tests | Applicability helper includes directories, includes unfiltered Markdown for every dimension, includes filtered Markdown only on normalized match, and excludes inapplicable pairs without failure. | `/home/antonioborgerees/coding/ultraplan-go/internal/study/source_test.go` |
| Tests | `study <study> list` output shows source kind and Markdown applicability details for directory, unfiltered Markdown, and filtered Markdown sources in deterministic order. | `/home/antonioborgerees/coding/ultraplan-go/internal/study/cli_test.go` |
| Tests | Full test suite passes. | `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go` |
| Runtime | No runtime behavior is executed or required. | Review confirms no OpenCode, `agentwrap`, network, provider, scheduler, prompt, report-validation, or summary behavior was added beyond compile-time compatibility. |
| Review | Architecture ownership is preserved. | Architecture review confirms implementation remains in `internal/study`, `study` may use `workspace`, and platform packages do not import `study`. |
| Review | Error behavior is actionable and cause-preserving. | Review checks wrapped errors and diagnostics include source path/offending value for invalid `applicable_dimensions`. |
| Documentation | Sprint artifacts remain aligned. | `reasoning.md`, future `plan.md`, and future `review.md` map decisions to `requirements.md`, `sprint-index.md`, and `technical-handbook.md`. |

## Assumptions And Risks

| Item | Type | Impact | Mitigation / Follow-Up |
| --- | --- | --- | --- |
| No area-specific reasoning file exists despite sprint-index selection. | Risk | A planned pre-reasoning artifact did not inform this document. | This document resolves the open questions directly and records the absence; `plan.md` should not wait for that artifact unless the team creates it later and intentionally revises this reasoning. |
| Existing code may already have source/domain/listing shapes from prior sprints. | Assumption | Implementation may need to extend rather than replace existing structs/functions. | Plan should inspect current `internal/study` files before editing and preserve existing behavior/tests. |
| YAML parsing dependency may already exist or may need to be introduced. | Risk | Adding a new dependency for frontmatter could be heavier than needed. | Prefer existing project YAML dependency if present; otherwise choose the smallest idiomatic Go YAML dependency already acceptable in the codebase. |
| Text listing format changes can break scripts. | Risk | Existing users/tests may rely on current output. | Keep output additive and deterministic where possible; command tests lock expected output. |
| Hard failure on malformed frontmatter may reduce partial visibility. | Risk | One invalid Markdown file blocks listing all sources. | Provide path/cause diagnostics; future validation can aggregate issues if needed. |
| Normalizing dimensions to two digits assumes sprint dimension numbering stays below 100 or uses two-digit matching for now. | Assumption | Values like `100` need clear behavior. | Implement two-digit normalization per sprint; tests should define invalid or accepted behavior according to existing dimension normalization helpers. |
| Future run-loop/status/summary flows need identical applicability semantics. | Risk | Divergent future implementations could treat skipped pairs differently. | Introduce a study-owned helper and document exclusion/skipped semantics in tests. |

## Implementation Constraints

- Keep all source discovery, frontmatter parsing, applicability filtering, and list output behavior in `internal/study`.
- Do not create global `validation`, `scheduler`, `reports`, or `prompts` packages for this sprint.
- Preserve dependency direction: `study` may depend on `workspace` and platform helpers, but platform packages must not import `study`.
- Inspect only direct children of `studies/<study>/sources/`; do not recursively scan source directories.
- Discover only non-hidden directories and top-level `.md` files as sources.
- Ignore hidden entries, non-directory non-Markdown files, and nested Markdown files inside source directories.
- Parse YAML frontmatter only when it is the leading document block delimited by `---`.
- Strip only leading frontmatter and preserve unchanged content when no leading frontmatter exists.
- Normalize `applicable_dimensions` values to two-digit dimension numbers; accept numeric and string values; treat missing/empty filters as all dimensions.
- Report invalid applicability values with source path and offending value; preserve underlying causes for filesystem/YAML errors.
- Directory sources are always applicable and must not depend on frontmatter.
- Inapplicable Markdown source/dimension pairs are excluded/skipped, not failures.
- Keep listing output deterministic, plain text, calm, and capturable in command tests.
- Do not add runtime execution, prompt composition, report validation, summary generation, durable workflow, target, or sprint-execution behavior.
- Tests must use local filesystem fixtures/temp dirs and must not require OpenCode, `agentwrap`, network access, or provider credentials.

## Plan Handoff

`plan.md` must execute these decisions. It must not invent architecture, scope, or decisions beyond this document.

The plan must carry forward:

- final decisions
- applicable contracts and requirement IDs
- expected evidence
- risks and assumptions
- required review protocols

## Phase Exit Criteria

- [ ] Selected context was read and used.
- [ ] Area-specific reasoning documents were completed where applicable or explicitly marked not available.
- [ ] Area-specific reasoning conclusions are reflected or explicitly overridden in final decisions.
- [ ] Contracts are explicitly mapped to decisions and expected evidence.
- [ ] Final decisions are clear enough for `plan.md` to execute.
- [ ] Expected evidence is specific and reviewable.
