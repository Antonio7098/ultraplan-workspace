> **Inputs Used:** `projects/ultraplan-go/sprints/03-study-listing-resolution/technical-handbook.md`, `projects/ultraplan-go/sprints/03-study-listing-resolution/requirements.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/PRD.md`, `system/reasoning/architecture_reasoning_template.md`

# Architecture Reasoning: Study Domain, Listing, and Resolution

## 0. Feature Summary

### Feature name

Study domain discovery, reference resolution, and listing commands.

### User/product goal

Users need to inspect existing studies from an UltraPlan workspace without initializing studies, running analysis, invoking OpenCode, or scanning source repositories recursively. The required user-facing commands are `ultraplan study list` and `ultraplan study <study> list`.

### Current task

Implement the sprint 3 study-side read model:

- `internal/study/domain.go` for `Study`, `Source`, `SourceKind`, and `Dimension`.
- `internal/study/discovery.go` for shallow deterministic discovery under `studies/`, `sources/`, and `dimensions/`.
- `internal/study/resolve.go` for exact and unambiguous-prefix reference resolution.
- `internal/study/service.go` for listing use cases.
- `internal/app/study_commands.go` for command wiring and output.
- Unit and command tests for domain normalization, discovery, ordering, resolution, output, and actionable errors.

### Non-goals

- Do not implement study initialization, YAML parsing, cloning, or assisted completion.
- Do not implement Markdown document source discovery, frontmatter parsing, or applicability filtering in this sprint, even though the PRD/TRD define it for the broader product.
- Do not implement prompt composition, runtime execution, synthesis, report validation, summary generation, run state, retry, cancellation, target workflows, or sprint workflows.
- Do not add non-trivial JSON listing output unless existing CLI output support makes it trivial.

## 1. First-Principles Breakdown

### Core behaviour

The feature reads the resolved workspace filesystem, derives a deterministic study listing model, resolves user-supplied study/source/dimension references, and prints stable human-readable lists. It does not create or update study artifacts.

### Inputs

- Resolved workspace root from existing workspace discovery and global `--workspace` behavior.
- Filesystem entries under `studies/`, `studies/<study>/sources/`, and `studies/<study>/dimensions/`.
- User command arguments: optional study reference for `ultraplan study <study> list`.
- Existing command IO streams from `internal/app`.

### Outputs

- Deterministically sorted study names for `ultraplan study list`.
- Deterministically sorted source and dimension rows for `ultraplan study <study> list`.
- Actionable non-zero errors for missing and ambiguous study, source, or dimension references.

### Durable state

- Created: none.
- Read: workspace `studies/` tree, direct child entries of each selected study's `sources/` and `dimensions/` directories.
- Updated: none.
- Deleted: none.

### Ephemeral state

- In-memory `Study`, `Source`, and `Dimension` values derived from shallow filesystem reads.
- Candidate reference sets used to detect exact, prefix, missing, and ambiguous matches.
- Command output buffers in tests.

### Derived state

- Study names are derived from non-hidden directories directly under `studies/`; source of truth is the workspace filesystem; computed on demand; staleness risk is low.
- Directory sources are derived from non-hidden directories directly under `sources/`; source of truth is the study filesystem; computed on demand; staleness risk is low.
- Dimension number and slug are derived from Markdown filenames beginning with a numeric prefix; source of truth is filenames; computed on demand; staleness risk is low.

### Side effects

- Database write: no.
- File write: no.
- Network call: no.
- Queue/event emission: no.
- Email/notification: no.
- Logs/metrics/traces: no new observability surface required for this sprint.
- Terminal output: yes, command stdout and stderr through existing app IO paths.

## 2. Existing Architecture Fit

### Contract applicability and phase gate

Current gate: Skeleton / Local CLI.

Applies now:

- ARCH-CORE-001: keep `internal/study`, `internal/workspace`, and `internal/app` ownership explicit.
- ARCH-ENTRY-001: keep CLI handlers thin and delegate study behavior to `internal/study`.
- ERR-TRANS-001: preserve cause chains when command errors wrap study/workspace failures.
- TEST-FAIL-001: cover missing and ambiguous reference failures.
- CLI Surface: preserve command shape, output streams, and exit statuses.

Deferred:

- ARCH-LAYER-002 broad ports for runtime/process/persistence dependencies; listing uses narrow local filesystem reads and does not cross a volatile provider boundary.
- ARCH-COMP-001 full registrar/container structure; introduce before runtime/provider or durable workflow wiring grows command composition.
- OBS-CORR-001, OBS-DIAG-001, OBS-ALERT-001; listing does not introduce runtime logs, traces, diagnostics, or alertable background work.
- PERF-CONC-001 and workflow/persistence gates; listing has no worker pool, run-loop, durable task state, or migration behavior.

### Where does this behaviour naturally belong?

Domain discovery, sorting, normalization, and resolution belong in `internal/study`. Command registration, argument parsing, workspace option consumption, and output rendering belong in `internal/app`. Workspace root resolution remains in `internal/workspace`.

This follows the project architecture rule that product modules own product behavior and workspace owns location/path rules. `ARCHITECTURE.md` states that `internal/study` owns sources, dimensions, filesystem persistence, validation, reports, and CLI workflows, while `internal/workspace` owns root discovery and path safety. The sprint narrows that broader ownership to read-only discovery and listing.

### Existing workflow affected

Current flow before this sprint:

1. `cmd/ultraplan` delegates to `internal/app`.
2. `internal/app` registers existing commands and global `--workspace` behavior from prior sprints.
3. Workspace discovery resolves a workspace root or returns a workspace error.
4. No study-domain package or study listing commands exist yet for this sprint scope.

### Proposed new flow

Proposed flow:

1. `internal/app` registers `study list` and `study <study> list` using existing command factory patterns.
2. Command execution asks existing workspace behavior for the resolved workspace root.
3. `internal/app` constructs or calls an `internal/study.Service` with the workspace root.
4. `internal/study` performs shallow deterministic discovery and reference resolution.
5. `internal/app` formats the returned domain values to stdout and maps domain errors to actionable command errors.

### Does the current architecture still fit?

[x] Yes - the feature fits cleanly into the existing shape.

The technical handbook's strongest evidence supports thin CLI delegates and domain-owned behavior. It cites mature Go CLIs such as chezmoi `main.go:26-34`, gh-cli `cmd/gh/main.go:6`, helm `pkg/cmd/install.go:132-145`, and yq `cmd/root.go:9` as evidence for keeping business behavior out of command handlers. The current project architecture already expects this split.

### Refactor-before-feature decision

Decision: implement directly with a small new `internal/study` package and app command wiring.

Reason: the required behavior is cohesive, read-only, and aligned with documented module boundaries. A larger refactor or new global technical layer would add complexity not justified by sprint 3.

## 3. Design Options Considered

### Option A: Minimal module-owned implementation

Description:
Create `internal/study` with concrete domain types, shallow discovery functions, resolution helpers, and a small service used by `internal/app` commands.

Pros:
- Matches `ARCHITECTURE.md` and sprint constraints that study behavior must live in `internal/study`.
- Keeps command handlers thin and testable, matching technical handbook evidence from command architecture and IO abstraction reports.
- Avoids global technical packages such as `internal/validation`, `internal/reports`, or `internal/scheduler`.
- Supports direct unit tests for domain behavior and command tests for output.

Cons:
- Requires a small service boundary even though the first use cases are simple listings.
- Some logic may later grow when study initialization and Markdown document sources arrive.

Risks:
- If the service surface is too broad, it may become a generic study manager. Keep it use-case-oriented: list studies and list one study.

### Option B: Implement discovery and resolution inside command handlers

Description:
Add all filesystem reads, sorting, prefix matching, and output generation directly in `internal/app/study_commands.go`.

Pros:
- Fewest files in the short term.
- Easy to see command output and command behavior in one place.

Cons:
- Violates sprint constraint that domain discovery and resolution logic must stay outside CLI command handlers.
- Conflicts with technical handbook anti-pattern warning: do not put domain rules in `RunE`.
- Makes resolution and discovery harder to unit test without command scaffolding.

Risks:
- Prefix-matching and missing/ambiguous error behavior could diverge across future commands.

### Option C: Introduce generic filesystem/repository abstractions first

Description:
Create broad interfaces or platform-level packages for discovery, validation, reports, and listing, then implement study behavior through those abstractions.

Pros:
- Could support future reuse if many modules need identical discovery semantics.
- Makes mocking explicit if integration tests need a fake filesystem.

Cons:
- `ARCHITECTURE.md` explicitly warns against global technical-layer packages and large clean-architecture trees before there is a concrete need.
- The technical handbook warns against over-broad abstractions for a small listing sprint.
- This sprint only needs shallow local reads and can be tested with temp fixture workspaces.

Risks:
- Premature abstractions could obscure the simple read-only flow and create package cycles or false reuse.

### Chosen option

Chosen option: Option A.

Reason: it is the simplest honest design that satisfies the sprint's architecture constraints, acceptance criteria, and evidence. It keeps study rules near study state while preserving thin command handlers.

## 4. Abstraction Check

### Are we adding a new abstraction?

[x] Yes - service/component.
[x] Yes - data structure/DTO/config object.

The service is earned because `internal/app` needs a narrow use-case boundary into `internal/study` without knowing discovery internals. Domain data structures are earned because `Study`, `Source`, `SourceKind`, and `Dimension` are explicit PRD/TRD concepts and required sprint outputs.

### Why is it earned?

- It protects stable domain concepts from being represented as loose strings and filesystem paths in command handlers.
- It creates a clear boundary between CLI/output policy and study discovery/resolution mechanism.
- It makes testing simpler: domain unit tests can exercise discovery and resolution without Cobra command setup, while command tests can focus on output and exit behavior.
- It follows technical handbook evidence favoring explicit service surfaces over global state, including gh-cli factory/service examples and go-task constructor injection examples.

### Bad abstraction smell check

- Generic manager naming: reject names like `StudyManager`; prefer `Service` only inside the `study` package or explicit functions with use-case names.
- Single implementation: acceptable for the service because it is not polymorphic; do not introduce an interface unless tests or app construction need one.
- Boolean mode switches: avoid. Separate methods such as `ListStudies` and `ListStudy` are clearer than a mode parameter.
- Flow traceability: keep discovery and resolution in focused files rather than nested packages.

## 5. DRY and Duplication Check

### Is any duplication being introduced?

[x] No duplicated business/domain knowledge should be introduced.

Single sources of truth:

- Study discovery rules live in `internal/study/discovery.go`.
- Source reference matching rules live in `internal/study/resolve.go`.
- Dimension reference matching rules live in `internal/study/resolve.go`.
- Command output formatting lives in `internal/app/study_commands.go`.

Some code-shape similarity between source and dimension resolution is acceptable only if it keeps each rule explicit. If shared candidate matching is introduced, it should be a small internal helper in `internal/study`, not a global matching package.

## 6. Coupling Check

### Dependencies required

- `internal/workspace`: command flow needs the existing resolved workspace root and workspace path behavior from sprint 2.
- Go standard library filesystem/path packages: `internal/study` needs shallow directory reads and workspace-relative path handling.
- Existing app command framework: `internal/app` owns CLI registration, arguments, help, and output.

### Coupling risks

Global coupling: none introduced. Avoid package-level mutable workspace or output state.

Content coupling: none introduced if `internal/app` consumes `study.Service` results rather than recreating discovery internals.

Stamp coupling: low. Passing whole `Study` values to listing output is appropriate because output needs study, sources, and dimensions. Internal helper functions should receive narrower slices or values where possible.

### Narrow dependency check

Each function should receive only what it needs. Discovery functions should receive a workspace root or study root, not the full app object. Resolution functions should receive candidate slices and the reference string, not command objects.

### Dependency injection check

Side-effectful dependencies are created at the edge and passed inward where useful. For this sprint, concrete filesystem reads are acceptable because the behavior is local, shallow, and testable with temporary directories. Do not introduce a filesystem interface unless existing platform filesystem abstractions already require it.

## 7. State and Mutation Check

### Mutation points

- `internal/study`: no durable mutation; only in-memory derived structs.
- `internal/app`: no durable mutation; only command output writes.
- Tests: create fixture workspaces in temporary directories.

### Is mutation explicit from names and flow?

[x] Yes.

Listing methods should be named as queries, for example `ListStudies` and `ListStudy`. No method should imply initialization, validation repair, or run state mutation.

### Are queries and commands separated where practical?

[x] Yes.

Study service methods are read-only queries. CLI commands perform output side effects after receiving results.

### Could hidden state make this hard to debug?

[x] No, if package globals are avoided.

## 8. Function/Class/Module Shape

### Primary unit of behaviour

[x] Module.
[x] Service/component.

### Why this unit?

The behavior is cohesive around study filesystem state and reference rules. A single `internal/study` package split into focused files matches `ARCHITECTURE.md`: start with one package per module and multiple focused files, then split only when there is a concrete readability or dependency benefit.

### If using composition, what components are combined?

- `internal/app`: command construction, argument parsing, workspace option handling, output rendering.
- `internal/workspace`: workspace root discovery and path rules.
- `internal/study`: study domain values, shallow discovery, reference resolution, listing use cases.

## 9. Error Handling Design

### Expected failures

- Missing workspace or invalid workspace: reuse existing workspace command error behavior.
- Missing `studies/` directory: treat as an empty study listing unless existing workspace validation considers it invalid; this keeps read-only listing calm and avoids inventing mutation/initialization behavior.
- Missing selected study: return a study not found error that names the requested reference and suggests `ultraplan study list`.
- Ambiguous selected study: return an ambiguity error that lists matching study names.
- Missing source or dimension reference: return source/dimension not found errors naming the requested reference and available context.
- Ambiguous source or dimension reference: return ambiguity errors listing conflicting candidates.
- Filesystem read failure: wrap with the path and operation, preserving workspace-relative output where practical.

### Unexpected failures

- Permission denied, broken symlinks, or stat errors should surface as filesystem/workspace errors with safe path context.
- Malformed dimension filenames should be ignored unless they begin like a dimension and cannot normalize safely; tests should define this edge. The acceptance criteria only require Markdown files directly under `dimensions/` whose names begin with a numeric prefix.

### Error taxonomy

Error types/classes introduced or reused:

- Study not found: requested study reference matched no discovered study.
- Source not found: requested source reference matched no source in a resolved study.
- Dimension not found: requested dimension reference matched no dimension in a resolved study.
- Ambiguous reference: requested prefix matched more than one candidate.
- Filesystem/workspace error: shallow discovery could not read required workspace paths.

The technical handbook supports structured missing/ambiguous reference errors with actionable hints, citing go-task `TaskNotFoundError` with `DidYouMean`, age `errorWithHint()`, and gh-cli `printError()`.

### Retry / recovery behaviour

- Retry safe: yes, commands are read-only.
- Partial progress possible: no durable partial progress.
- Compensation needed: no.

## 10. Observability Design

### Events/logs/metrics/traces needed

No new logging, metrics, traces, event records, run IDs, or observability stores are required for this sprint. The sprint-index explicitly excludes the Observability contract and runtime diagnostics.

### Minimum useful fields

For user-facing errors, include:

- Operation or command context.
- Requested reference.
- Candidate matches for ambiguity.
- Workspace-relative path where practical.
- Suggested next command, such as `ultraplan study list`.

### Sensitive data check

Could this expose PII/secrets/user content? Yes, if absolute local paths are printed unnecessarily.

Mitigation: prefer workspace-relative paths in output and errors where practical, satisfying sprint constraints and PRD artifact guidance to avoid absolute paths unless necessary for diagnostics.

## 11. Testing Strategy

### Unit tests

- Domain normalization for dimension numbers and slugs derived from filenames.
- Study discovery ignores hidden entries and non-directory files under `studies/`.
- Source discovery returns non-hidden direct child directories only and ignores files, hidden entries, and nested repository contents.
- Dimension discovery returns direct Markdown files with numeric prefixes and normalizes numbers to two digits.
- Deterministic sorting of studies, sources, and dimensions.
- Study/source/dimension reference resolution for exact, prefix, missing, and ambiguous cases.

### Use-case/workflow tests

- `ultraplan study list` prints discovered studies in deterministic order.
- `ultraplan study <study> list` prints sources with directory source kind and dimensions in deterministic order.
- Missing and ambiguous study references return non-zero actionable errors.
- Command behavior uses explicit `--workspace` through existing global workspace behavior.

### Integration tests

No real OpenCode, network, cloned repository, or provider-credential integration tests are needed. The requirements explicitly prohibit normal unit tests from depending on those systems.

### Regression tests

- Hidden files/directories are ignored.
- Non-directory source entries do not become Sprint 3 sources.
- Top-level `.md` source files remain deferred and ignored in Sprint 3, even though broader PRD/TRD source discovery includes Markdown document sources.
- Source repositories are not recursively scanned.

### Testability check

[x] Yes, the core behavior can be tested without real infrastructure.

The technical handbook supports behavior-focused unit and command tests using fixtures and output assertions, citing chezmoi testscript, go-task golden normalization, and gdu command helpers.

## 12. Performance and Scale Assumptions

### Known constraints

- Expected volume: dozens of studies, sources, and dimensions per workspace for listing; source directories may themselves be large repositories.
- Latency sensitivity: simple listing should feel fast and avoid expensive work.
- Memory sensitivity: low for shallow listings; high risk only if source repositories are scanned recursively, which is forbidden.
- Concurrency concerns: none for this sprint; discovery is shallow and read-only.

### Assumptions being made

- Direct children of `studies/`, `sources/`, and `dimensions/` are small enough to read and sort in memory.
- Absent `sources/` or `dimensions/` directories for an existing study can be represented as empty sections rather than fatal errors, unless filesystem read errors indicate a real problem.
- Dimension filenames follow the product convention `NN-slug.md`; files that do not begin with a numeric prefix are outside Sprint 3 dimension discovery.

### Measurement plan

[x] Not needed for this change.

### Optimization decision

Do not add caching, indexing, concurrency, or filesystem watchers. Performance evidence in the technical handbook favors bounded traversal and deferring expensive work; Sprint 3 satisfies that by reading only direct children and sorting deterministic slices.

## 13. Implementation Plan

### Files likely to change

- `/home/antonioborgerees/coding/ultraplan-go/internal/study/domain.go`: define study domain types and normalization helpers.
- `/home/antonioborgerees/coding/ultraplan-go/internal/study/discovery.go`: implement shallow workspace/study/source/dimension discovery.
- `/home/antonioborgerees/coding/ultraplan-go/internal/study/resolve.go`: implement exact and unambiguous-prefix resolution with typed errors.
- `/home/antonioborgerees/coding/ultraplan-go/internal/study/service.go`: expose listing use cases to `internal/app`.
- `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_commands.go`: register study listing commands and render output.
- `/home/antonioborgerees/coding/ultraplan-go/internal/study/study_test.go`: cover domain, discovery, sorting, ignored entries, and resolution.
- `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_commands_test.go`: cover command output and command errors.

### Step-by-step plan

1. Add domain types and normalization helpers in `internal/study`.
2. Add shallow discovery for studies, directory sources, and Markdown dimensions with deterministic sorting.
3. Add resolution helpers and typed/structured errors for missing and ambiguous references.
4. Add a small service with `ListStudies` and `ListStudy` use cases.
5. Wire app commands through existing workspace resolution and injected output streams.
6. Add unit tests for study package behavior.
7. Add command tests for output and non-zero actionable errors.
8. Verify with `go test ./...` and `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`.

### Migration/backwards compatibility

- Schema migration needed: no.
- API compatibility affected: no public API; new CLI commands are additive within planned scope.
- Existing data affected: no writes.
- Feature flag needed: no.

### Rollback plan

Because the feature is additive and read-only, rollback is removing the new `internal/study` files, app command wiring, and tests. No persisted data cleanup is required.

## 14. Final Pre-Implementation Decision

Decision: Proceed.

Reason: Sprint 3 fits the existing module-driven architecture. The evidence supports `internal/study` owning discovery and resolution, `internal/app` owning command wiring and output, and `internal/workspace` remaining responsible for workspace root/path behavior.

Complexity introduced: low.

Complexity removed: none yet; this establishes the intended boundary for future study-side work.

Main trade-off: accept a small `study.Service` boundary now to prevent CLI handlers from accumulating domain rules. Reject broader abstractions, recursive discovery, Markdown document source handling, and runtime-related architecture until later sprints require them.

## Key Conclusion And Evidence Basis

The final area decision is to implement Sprint 3 as a read-only `internal/study` module with a narrow service consumed by thin `internal/app` commands. This is grounded in:

- Technical handbook evidence for thin CLI/domain-owned behavior from chezmoi, gh-cli, helm, and yq.
- Technical handbook evidence for explicit service/factory injection and testable IO from gh-cli, go-task, restic, and mitchellh-cli.
- Technical handbook evidence for actionable typed reference errors from go-task, age, and gh-cli.
- Technical handbook performance/security evidence favoring deterministic bounded discovery instead of recursive source scans.
- `ARCHITECTURE.md` dependency rules: product modules may depend on workspace/platform packages, platform packages must not depend on product modules, and study behavior should remain in `internal/study`.
- Sprint requirements that domain discovery and resolution must stay outside CLI command handlers and listing commands must use prior workspace behavior.

## Rejected Alternatives

- Rejected CLI-owned discovery and resolution because it violates the sprint constraint and technical handbook anti-patterns around large `RunE` command handlers.
- Rejected global technical packages such as `internal/discovery`, `internal/validation`, or `internal/reports` because the project architecture says module-owned behavior should stay with the owning product module.
- Rejected recursive source inspection because sprint requirements forbid recursive source repository scans and evidence emphasizes bounded traversal for performance and security.
- Rejected Sprint 3 Markdown document source support because requirements explicitly defer top-level `.md` source discovery, frontmatter parsing, and applicability filtering to a later sprint.

## Risks, Assumptions, And Open Questions

### Risks

- Prefix resolution can surprise users if exact matches and prefix matches are not ordered consistently. Mitigation: exact match wins; otherwise only one prefix candidate may match.
- Error output can become noisy if it prints absolute paths or too many candidates. Mitigation: prefer workspace-relative paths and concise candidate lists.
- Future Markdown document source support may require expanding `SourceKind` behavior. Mitigation: include `SourceKindMarkdown` only if required by the domain model, but keep Sprint 3 discovery limited to directories.

### Assumptions

- Missing `studies/` can be listed as empty if the workspace itself is valid.
- Missing `sources/` or `dimensions/` under a resolved study can be shown as empty sections rather than fatal errors.
- Existing workspace behavior already enforces the global `--workspace` precedence required by this sprint.

### Open questions

- Should ambiguous reference errors include all candidate aliases or only canonical names? Decision for implementation should prefer canonical names unless tests show aliases are needed.
- Should dimension title be parsed from Markdown content now or derived from filename? Sprint 3 requirements only require file discovery, number normalization, and listing; deriving from filename is sufficient unless existing docs/tests demand content parsing.
- Should absent `studies/` be an empty result or workspace validation error? This reasoning chooses empty listing for a valid workspace, but implementation should align with existing workspace validation from prior sprints.

## Sprint Reasoning Reference

This document should be referenced by `projects/ultraplan-go/sprints/03-study-listing-resolution/reasoning.md` when the sprint-level reasoning artifact is created. That artifact is not present yet in the current sprint directory.
