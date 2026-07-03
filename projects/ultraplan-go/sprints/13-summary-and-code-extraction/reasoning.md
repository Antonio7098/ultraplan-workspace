# Sprint Reasoning: 13 Summary And Code Extraction

> Project: `ultraplan-go`
> Sprint: `13-summary-and-code-extraction`
> Output: `projects/ultraplan-go/sprints/13-summary-and-code-extraction/reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/13-summary-and-code-extraction/requirements.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/sprints/13-summary-and-code-extraction/sprint-index.md`, `projects/ultraplan-go/sprints/13-summary-and-code-extraction/technical-handbook.md`, `templates/sprint-reasoning.md`, `studies/go-cli-study/reports/final/01-project-structure.md`, `studies/go-cli-study/reports/final/02-command-architecture.md`, `studies/go-cli-study/reports/final/03-dependency-injection.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/06-io-abstraction.md`, `studies/go-cli-study/reports/final/09-terminal-ux.md`, `studies/go-cli-study/reports/final/10-logging-observability.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/13-security.md`, `studies/go-cli-study/reports/final/14-performance.md`, `studies/go-cli-study/reports/final/15-philosophy.md`

This document decides. It synthesizes selected context, handbook evidence, available area-specific reasoning, and contracts into final sprint decisions.

It does not replace `sprint-index.md`, `technical-handbook.md`, or `reasoning/*.md`.

Requirement IDs used below are assigned from the ordered acceptance criteria in `requirements.md`: `AC-01` is the first acceptance criterion at line 40, and `AC-31` is the final acceptance criterion at line 70.

## Sprint Purpose

- **Goal:** Implement standalone deterministic study summary generation and code-reference extraction so completed study outputs can be reviewed, compared, and audited from CSV, text, and JSON surfaces.
- **Non-Goals:** Do not implement runtime execution, prompt composition changes, report template changes, run-loop orchestration, stale task recovery, target/sprint commands, remote source fetching, hosted/UI workflows, plugin systems, or release-wide stable public JSON schema guarantees.
- **Depends On:** Sprint requirements identify dependencies on prior study-side work for Markdown source discovery and applicability, report path helpers, validation and rating parsing, prompt/report artifact conventions, existing Sprint 11 summary behavior, and compatibility with completed run-loop outputs. No prior decision artifacts are listed in `project-index.md`.

## Selected Context And Pre-Reasoning Artifacts

| Artifact | Path | How It Was Used |
| --- | --- | --- |
| Project Index | `projects/ultraplan-go/project-index.md` | Established project scope, target implementation directory, available contracts, selected evidence report pool, and absence of prior decisions. |
| Sprint Requirements | `projects/ultraplan-go/sprints/13-summary-and-code-extraction/requirements.md` | Provided authoritative sprint goal, required outputs, acceptance criteria, constraints, dependencies, and review expectations. |
| PRD | `projects/ultraplan-go/docs/PRD.md` | Grounded user-facing need for citation auditing, deterministic summaries, local-first artifacts, deferred target/sprint workflows, and study-side-only scope. |
| TRD | `projects/ultraplan-go/docs/TRD.md` | Grounded module layout, CLI behavior, exit-code classes, source applicability rules, citation syntax, source table parsing, resolution algorithm, summary semantics, persistence, security, and testing requirements. |
| Architecture | `projects/ultraplan-go/docs/ARCHITECTURE.md` | Applied module ownership rule: `internal/study` owns summary, `internal/codeextract` owns extraction, `internal/app` stays thin, and `internal/platform/runtime` remains product-free. |
| Sprint Index | `projects/ultraplan-go/sprints/13-summary-and-code-extraction/sprint-index.md` | Selected the binding contracts, evidence reports, review protocols, exclusions, and next-artifact chain for this sprint. |
| Technical Handbook | `projects/ultraplan-go/sprints/13-summary-and-code-extraction/technical-handbook.md` | Supplied study evidence, relevant patterns, anti-patterns, trade-offs, and open questions for final decisions. |
| Sprint Reasoning Template | `templates/sprint-reasoning.md` | Provided this document structure and required sections. |
| Selected Evidence Reports | `studies/go-cli-study/reports/final/*.md` listed in the inputs line | Used concrete studied repo references for CLI boundaries, error handling, IO seams, output testing, path safety, bounded scanning, and engineering trade-offs. |
| Prior Decision | N/A | No prior decision artifacts are listed in `project-index.md`; carry-forward constraints come from requirements, docs, and selected evidence. |

## Area-Specific Reasoning Inputs

The sprint index selected two area reasoning outputs, but no `reasoning/*.md` files currently exist for this sprint. Therefore no area-specific conclusions were available to summarize or override. This consolidated reasoning makes final decisions directly from `requirements.md`, `sprint-index.md`, `technical-handbook.md`, project docs, selected contracts, and selected evidence reports.

| Area | Reasoning Document | Key Conclusion | Evidence Basis | Impact On Final Decision |
| --- | --- | --- | --- | --- |
| Summary generation | `reasoning/summary-generation.md` | No file exists for this sprint. | No area reasoning input was present. Requirements, TRD summary rules, architecture docs, and handbook evidence are used directly. | Final decisions specify summary semantics without reopening area reasoning. |
| Code extraction | `reasoning/code-extraction.md` | No file exists for this sprint. | No area reasoning input was present. Requirements, TRD code extraction rules, architecture docs, and handbook evidence are used directly. | Final decisions specify extraction parsing, resolution, output, exit behavior, and safety constraints. |

## Sprint Technical Handbook Summary

- **Relevant Patterns:** Thin CLI adapters, product-owned modules, request/result types at command seams, injectable IO, deterministic renderers, behavior-focused errors, fixture/golden output tests, local-only path resolution, bounded source lookup, and deliberate scope boundaries.
- **Important Trade-Offs:** Thin CLI vs command-owned logic, exact/golden output tests vs flexible assertions, basename fallback usefulness vs ambiguity and scan risk, sprint-local JSON determinism vs release-wide schema commitment, warnings as non-fatal results vs visible partial exits, and minimal filesystem seams vs broad abstractions.
- **Warnings / Anti-Patterns:** Do not put parsing/resolution/totals in `internal/app`; do not duplicate existing study helpers; do not silently downgrade write failures or unresolved refs; do not hardcode output to `os.Stdout`; do not trust path strings; do not allow unbounded fallback scans; do not test only internal state.
- **Evidence Confidence:** High for CLI boundaries, testing, error handling, IO injection, security, and performance because selected reports cite repeated patterns across 16 Go CLI repositories. Medium for terminal UX and philosophy because those reports are broader and less directly tied to this non-interactive sprint, but still useful for output discipline and scope control.

## Contracts Applied

| Contract / Requirement ID | Constraint | Decision Impact | Expected Evidence |
| --- | --- | --- | --- |
| Architecture | Product logic must live in owning modules; CLI remains thin; platform runtime remains generic. | Summary stays in `internal/study`; extraction stays in `internal/codeextract`; `internal/app` parses, invokes, renders, maps exits only. | Import/path review; changed-file review; no `internal/platform/runtime` product imports; command tests exercise public surfaces. |
| Errors | Failures and partial outcomes need classification, wrapping, actionable diagnostics, and scriptable exits. | Summary write failures are non-zero; extraction distinguishes usage, validation, filesystem, and partial unresolved outcomes. | Unit and command tests for warning output, missing table, unresolved refs, write failure, unreadable inputs, and exit categories. |
| Security | Paths, source roots, diagnostics, snippets, and writes must not leak secrets or traverse outside intended roots. | Extraction normalizes and contains paths under parsed roots; ignores dangerous fallback dirs; rejects escapes; output prints only requested snippets and safe context. | Tests for path escapes, missing files, ambiguous basename, ignored dirs, output redaction review, and no runtime/env/prompt payload output. |
| Testing | Normal tests must be deterministic, fixture-backed, and offline. | Add unit, fixture, command-level, and exact/golden-style tests for summary CSV, text, JSON, warnings, and failures. | `go test ./...`; exact expected CSV/text/JSON fixtures; command tests with fake/local files only. |
| Documentation | CLI help and generated artifacts must be reviewable and maintainable. | Help exposes `study <study> summary` and top-level `code`; output text is deterministic and repair-oriented. | Help-output command tests; review of command usage text; generated `summary.csv` and extraction output examples in tests. |
| CLI Surface | Commands need stable shape, flags, help, text/JSON output, and exit-code handling. | Implement `ultraplan study <study> summary`; implement `ultraplan code <report>...` with `--json` and `--output <path>`. | Command tests for help, success, warning/partial, usage errors, JSON output, and output-file writes. |
| Performance | Avoid unbounded repository scans and large accidental memory use. | Basename fallback is scoped to parsed source roots, uses deterministic bounded traversal, skips ignored dirs, and caches lookup work per command invocation. | Resolver tests for ignored dirs and deterministic ambiguity; code review of traversal boundaries and cache scope. |
| Persistence And Migrations | Generated writes need explicit paths, safe permissions, and atomicity where required. | `summary.csv` writes atomically; optional extraction output writes use explicit destination and fail loudly on filesystem errors. | Atomic write tests; simulated write failure tests; file permission/path review. |
| AC-01 through AC-11 | Summary CLI, semantics, warnings, sorting, determinism, helper reuse, and atomic write behavior. | Decide exact summary matrix semantics and service/CLI shape. | `internal/study/summary_test.go`; `internal/app/study_summary_commands_test.go`; exact CSV and warning output checks. |
| AC-12 through AC-25 | Code CLI, source parsing, citation parsing, resolution, output, multi-report behavior, exit mapping, and module ownership. | Decide exact extraction module, parser, resolver, service, rendering, and exit policy. | `internal/codeextract/codeextract_test.go`; `internal/app/code_commands_test.go`; text/JSON golden or exact-output tests. |
| AC-26 through AC-31 | Module boundaries, runtime separation, safe diagnostics, offline tests, and build. | Preserve product/platform separation and non-runtime scope. | Import review; `go test ./...`; `go build ./cmd/ultraplan`; absence of OpenCode/provider/network dependencies in normal tests. |

## Repos Studied / Source Evidence Used

| Source / Repo / Report | Concrete Reference | Relevant Finding | Why It Matters For This Sprint | Used In Decision(s) |
| --- | --- | --- | --- | --- |
| Project structure report | `studies/go-cli-study/reports/final/01-project-structure.md`; chezmoi `main.go:16`; gh-cli `cmd/gh/main.go:6`; k9s `cmd/root.go:112-127` | Mature Go CLIs keep entrypoints thin and delegate product behavior to internal/domain packages. | Supports keeping summary and extraction behavior out of `internal/app`. | Decisions 1, 2, 3, 7 |
| Command architecture report | `studies/go-cli-study/reports/final/02-command-architecture.md`; gh-cli `pkg/cmdutil/factory.go:16-43`; helm `pkg/cmd/install.go:132-145`; opencode `cmd/root.go:49-183` as caution | Command factories and thin handlers are testable; long `RunE` functions accumulate product logic. | Supports command adapters that only parse args, call services, render outputs, and map exits. | Decisions 2, 5, 7 |
| Dependency injection report | `studies/go-cli-study/reports/final/03-dependency-injection.md`; gh-cli `pkg/cmd/factory/default.go:26-46`; go-task `executor.go:22-24`; restic `internal/restic/repository.go:18-66` | Manual DI, visible composition roots, narrow seams, and functional options dominate mature Go CLIs. | Supports minimal request/result types and injectable output/filesystem seams without a DI framework. | Decisions 2, 5, 7 |
| Error handling report | `studies/go-cli-study/reports/final/05-error-handling.md`; go-task `errors/errors.go:47-50`; gh-cli `internal/ghcmd/cmd.go:44-49`; age `cmd/age/tui.go:47-54`; rclone `fs/fserrors/error.go:22-29` | Errors should preserve chains, classify outcomes, map exits, and provide hints. | Supports visible summary write failures, partial extraction exits, and actionable unresolved diagnostics. | Decisions 1, 2, 5 |
| IO abstraction report | `studies/go-cli-study/reports/final/06-io-abstraction.md`; gh-cli `pkg/iostreams/iostreams.go:551-568`; go-task `executor.go:553-564`; restic `internal/ui/mock.go:10-53`; restic `internal/fs/interface.go:10-31` | Injectable IO and filesystem seams enable deterministic command tests and simulated failures. | Supports buffer-backed text/JSON tests and write-failure tests without broad abstractions. | Decisions 2, 5, 7 |
| Terminal UX report | `studies/go-cli-study/reports/final/09-terminal-ux.md`; gh-cli `internal/prompter/`; yq `pkg/yqlib/color_print.go:7-9` | CLI-first tools should keep non-interactive output concise, deterministic, and script-friendly. | Supports stable text output rather than rich TUI/progress behavior for this sprint. | Decisions 2, 5 |
| Logging and observability report | `studies/go-cli-study/reports/final/10-logging-observability.md`; gh-cli `pkg/iostreams/iostreams.go:52-54`; helm `internal/logging/logging.go:31-66`; k9s `internal/slogs/keys.go:6-231` | Diagnostics should stay separate from data output and use consistent safe context. | Supports not mixing secrets/prompts/runtime payloads into human/JSON extraction output. | Decisions 5, 6 |
| Testing strategy report | `studies/go-cli-study/reports/final/11-testing-strategy.md`; chezmoi `internal/cmd/main_test.go:64-174`; helm `internal/test/test.go:43`; go-task `task_test.go:166-169`; rclone `cmd/bisync/bisync_test.go:1435-1479` | Fixture-backed command tests and golden/exact output tests catch CLI artifact regressions. | Supports exact CSV, text, JSON, warning, and help-output tests. | Decisions 1, 2, 5, 7 |
| Security report | `studies/go-cli-study/reports/final/13-security.md`; go-task `internal/execext/exec.go:59-66`; opencode `internal/permission/permission.go:44-108`; restic `internal/options/secret_string.go:15-20`; helm `pkg/registry/transport.go:37-41`; gh-cli `pkg/surveyext/editor_manual.go:23` | Trust boundaries, argument/path discipline, and redaction must be explicit. | Supports local-only extraction, path containment, ignored dirs, no runtime calls, and safe diagnostics. | Decisions 3, 4, 5, 6 |
| Performance report | `studies/go-cli-study/reports/final/14-performance.md`; gh-cli `pkg/cmdutil/factory.go:27-42`; yq `pkg/yqlib/stream_evaluator.go:78-113`; gdu `pkg/analyze/parallel.go:13`; rclone `lib/pool/pool.go:17-24` | Large local processing should be lazy, streaming/bounded, and cached where useful. | Supports bounded basename fallback and deterministic per-command lookup caching. | Decisions 4, 6 |
| Philosophy report | `studies/go-cli-study/reports/final/15-philosophy.md`; gh-cli `pkg/cmd/factory/default.go:26-46`; opencode `internal/pubsub/broker.go:10-19`; lazygit `VISION.md:97` | Strong projects accept complexity deliberately and reject scope that does not serve the goal. | Supports not adding generic parser packages, plugins, workflow engines, or public schema guarantees. | Decisions 3, 5, 7 |

## Trade-Off And Debt Analysis

### Accepted Trade-Offs

| Trade-Off | Benefit | Cost / Constraint Accepted | Why Acceptable Now | Revisit Trigger |
| --- | --- | --- | --- | --- |
| Summary uses `N/A` only for inapplicable Markdown pairs, with empty cells for missing/invalid ratings. | Makes inapplicable distinct from missing expected reports while preserving simple CSV. | Consumers must interpret empty cells by consulting warnings. | Requirements explicitly call for distinct inapplicable representation and warnings for missing/ambiguous cases. | If users need machine-readable summary diagnostics beyond CSV. |
| Extraction supports only TRD inline code forms. | Keeps parser deterministic and auditable. | Some prose citations or Markdown link styles will not be extracted. | Requirements exclude citation forms outside single-line, range, and explicit line-list forms. | If report templates change to emit additional citation syntax. |
| Basename fallback remains enabled but bounded and ambiguity-producing. | Recovers legacy or weak citations without modifying reports. | Adds traversal cost and ambiguity risk. | Requirements require fallback; security/performance constraints bound it. | If large source roots make fallback too slow or ambiguous in practice. |
| JSON output is deterministic but sprint-local. | Enables automation and exact tests now. | Later releases may change schema after Sprint 14 work. | Requirements explicitly defer release-wide stable JSON schema compatibility. | Sprint 14 or release planning defines stable public JSON schemas. |
| Warnings and unresolved refs do not abort all processing. | Users can inspect useful results even when some reports are incomplete. | Exit-code mapping and output summaries must be very clear. | Requirements demand multi-report processing and partial/unresolved reporting. | If scripts need fail-fast mode later. |
| Minimal local file seams instead of broad filesystem framework. | Keeps implementation small and focused. | Some filesystem behavior may still use concrete `os` calls. | Evidence supports narrow seams at volatile/tested boundaries, not abstractions everywhere. | If multiple modules begin duplicating write, read, or path containment behavior. |

### Potential Technical Debt

| Debt / Shortcut | Why It Might Accrue | Current Mitigation | Owner / Follow-Up |
| --- | --- | --- | --- |
| Sprint-local JSON shape for `code --json`. | Automation users may treat it as stable before release-wide schema decisions. | Document and test deterministic fields without promising release-wide compatibility. | Sprint 14 schema compatibility work. |
| Fixed ignored-directory set for basename fallback. | Different ecosystems may need custom ignores or may store relevant code in ignored names. | Start with a documented conservative set and exact tests. | Future extraction configuration or performance tuning decision. |
| Empty summary cells rely on warnings for reason detail. | CSV alone cannot distinguish missing report, missing rating, and ambiguous rating. | Human output and `SummaryWarning` results carry source, dimension, report path, and reason. | Future structured summary artifact if needed. |
| Source table parser may encode assumptions about report table headings. | Existing or hand-written reports may vary slightly. | Support the required row shape and tests for malformed/actionable diagnostics. | Future parser compatibility work if real reports show common variations. |
| Basename lookup cache scope is per command only. | Repeated invocations repeat scans. | Keeps persistence out of scope and avoids cache invalidation complexity. | Add persistent cache only if profiling shows repeated extraction pain. |

### Future Considerations

| Consideration | Deferred Until | Reason Deferred | What Should Be Preserved Now |
| --- | --- | --- | --- |
| Release-wide stable JSON schemas | Sprint 14 / release hardening | Current sprint requires deterministic tested JSON for `code` only. | Keep JSON fields intentional, documented in tests, and versionable later. |
| Public reusable extraction API | Explicit product request | Current scope is CLI product behavior, not SDK exposure. | Keep `internal/codeextract` request/result types clean but internal. |
| Configurable extraction ignore rules | Performance or ecosystem need appears | No config changes are in selected sprint context. | Keep ignored-dir list centralized in `internal/codeextract`. |
| Structured summary JSON | Requirement change | Sprint requires CSV and deterministic human output only for summary. | Keep summary result/warning structs expressive enough to render other formats later. |
| Remote source resolution | Requirement change | Sprint explicitly excludes cloning/fetching/remote URL resolution. | Keep resolver local-root-only and error messages explicit. |
| Run-loop integration beyond compatibility | Future workflow sprint | Current sprint inspects completed outputs and must not mutate orchestration state. | Summary and extraction should not depend on run state. |

## Final Decisions

### Decision 1: Study-Owned Summary Matrix Semantics

- **Decision:** Implement summary generation in `internal/study` as a deterministic score matrix with header `source,<dimension...>,total`; dimension columns in discovered deterministic dimension order; rows sorted by total descending and source name ascending; valid parsed ratings as numeric cells; empty cells for missing expected report, missing rating, or ambiguous rating; `N/A` for inapplicable Markdown source/dimension pairs; totals as the sum of valid rating values only.
- **Rationale:** This directly satisfies the sprint's distinction between missing and inapplicable cells without inventing a new artifact format. It also preserves CSV reviewability and makes re-runs byte-identical when inputs do not change.
- **Study / Source Grounding:** `technical-handbook.md` design pressures require deterministic CSV semantics, warning collection, `N/A` cells, total-descending order, and helper reuse. TRD section 17 requires deterministic discovery, rating parsing, total-descending sorting, missing ratings as empty cells, and distinct inapplicable representation. The testing report supports exact/golden CSV checks via helm `internal/test/test.go:43` and go-task `task_test.go:166-169`.
- **Trade-Offs Accepted:** Empty cells intentionally group missing report, missing rating, and ambiguous rating in CSV; the detailed reason moves to warnings and command output. `N/A` is a human-readable sentinel rather than a numeric value.
- **Technical Debt / Future Impact:** If future automation needs per-cell reason codes, add a structured summary artifact rather than overloading CSV. Current result/warning types should keep enough information to add that later.
- **Alternatives Rejected:** Using `0` for missing or inapplicable cells is rejected because it would incorrectly affect totals and hide absence. Using empty cells for inapplicable pairs is rejected because it would not distinguish inapplicable from missing expected reports. Adding a global `internal/summary` package is rejected because architecture assigns study summary behavior to `internal/study`.
- **Contracts Satisfied:** Architecture, Errors, Testing, Persistence And Migrations, AC-01 through AC-11, AC-26, AC-30.
- **Evidence Required:** Unit tests comparing exact CSV for valid ratings, missing reports, missing ratings, ambiguous ratings, inapplicable pairs, totals, tie sorting, and repeated-run determinism; tests proving existing rating parser and report path helpers are reused; review of changed paths to confirm summary remains in `internal/study`.

### Decision 2: Standalone Summary Command And Loud Writes

- **Decision:** Expose `ultraplan study <study> summary` through thin `internal/app` wiring that resolves the workspace/study, calls a `study.Service` summary method, renders deterministic human output with warnings, writes `studies/<study>/summary.csv` atomically, and returns non-zero on write or command failures. The command must not invoke runtime, agentwrap, OpenCode, network access, or subprocess execution.
- **Rationale:** Summary regeneration is an artifact-inspection operation over completed study outputs. It should be runnable independently from batch execution and should fail loudly if it cannot update the durable CSV artifact.
- **Study / Source Grounding:** `technical-handbook.md` cites gh-cli `pkg/cmdutil/factory.go:16-43`, helm `pkg/cmd/install.go:132-145`, and IOStreams examples as support for thin command layers with injected output. The error-handling report warns against silent failure logging, citing lazygit `pkg/commands/git_commands/file_loader.go:52`. TRD section 21 requires explicit file writes and atomic state-style write discipline; requirements AC-09 requires atomic summary writes and non-zero failure.
- **Trade-Offs Accepted:** Command output remains concise rather than exposing all internal matrix detail; tests, CSV, and warnings carry the full review surface. Atomic write implementation adds a small local write helper or focused function instead of a broad filesystem abstraction.
- **Technical Debt / Future Impact:** If more generated artifacts require atomic writes, the implementation may later extract a platform filesystem helper. This sprint should avoid creating a generic persistence framework prematurely.
- **Alternatives Rejected:** Regenerating summary only as a side effect of `run-all` is rejected because users need standalone review and repair. Treating write failure as a warning is rejected because it creates invalid success. Putting command behavior directly in `RunE` is rejected because selected evidence flags monolithic command handlers as an anti-pattern.
- **Contracts Satisfied:** Architecture, CLI Surface, Errors, Persistence And Migrations, Security, Testing, AC-01, AC-08 through AC-11, AC-26 through AC-31.
- **Evidence Required:** Command tests for help, success, warning rendering, missing study/usage failures, write failure, deterministic human output, and non-runtime behavior; atomic write tests preserving previous content on failure where applicable; import review proving no runtime/agentwrap/OpenCode calls.

### Decision 3: Code Extraction Owned By `internal/codeextract`

- **Decision:** Add `internal/codeextract` with focused files for domain types, source-table/citation parser, resolver, and service orchestration. `internal/app` will only parse `ultraplan code <report>...` arguments and flags, call the extraction service, render text or JSON, write optional output, and map exit outcomes.
- **Rationale:** Code-reference extraction is a separate product capability with its own state transformations: report loading, source table parsing, citation parsing, source-root resolution, snippet extraction, and unresolved diagnostics. Keeping it in one owning module prevents a global parser/resolver layer from accumulating cross-module coupling.
- **Study / Source Grounding:** `ARCHITECTURE.md` explicitly states `internal/codeextract` owns parsing, resolution, extraction, and validation behavior. The project-structure report's thin CLI pattern and command-architecture report's helm/gh-cli examples support command delegation. The philosophy report supports accepting only deliberate complexity and rejecting generic extension surfaces that are not in scope.
- **Trade-Offs Accepted:** `internal/codeextract` may duplicate tiny generic-looking helpers such as citation line-spec parsing rather than creating global `internal/parser` or `internal/resolver` packages. This keeps ownership clear at the cost of less reuse for hypothetical future modules.
- **Technical Debt / Future Impact:** If study validation later needs the same parser, a shared internal helper can be extracted with a concrete second consumer. Until then, code extraction owns the implementation.
- **Alternatives Rejected:** Adding `internal/reports`, `internal/parser`, `internal/resolver`, or `internal/validation` is rejected because it violates the module-driven architecture unless reuse is proven. Putting extraction under `internal/study` is rejected because extraction can inspect arbitrary report paths and is a distinct product capability. Placing extraction logic in CLI commands is rejected because it blocks unit testing and violates selected evidence.
- **Contracts Satisfied:** Architecture, CLI Surface, Testing, Security, Performance, AC-12, AC-14 through AC-25, AC-27, AC-28, AC-31.
- **Evidence Required:** Changed-file review showing extraction behavior in `internal/codeextract`; command-level tests proving `internal/app` is only an adapter; import review showing no product dependency from `internal/platform/runtime` and no global parser/resolver package creation.

### Decision 4: Source Table, Citation, And Resolution Rules

- **Decision:** Parse supported inline-code citations matching `path:line`, `path:start-end`, `path:line,line,line`, and the documented legacy range dash form, normalizing ranges to standard hyphen output. A report requires a parseable sources table only when supported code citations are present or source resolution is otherwise requested. Source table paths resolve deterministically by trying workspace-relative paths first, then report-directory-relative paths. Each citation resolves by trying source-root-relative paths, source-prefixed path stripping, then basename fallback inside parsed source roots. Basename fallback returns unresolved when no match or more than one match is found.
- **Rationale:** This follows TRD resolution order while making ambiguous path bases deterministic and safe. Requiring a source table only when resolution is needed allows reports with no supported code refs to be inspected without a false validation failure, while still rejecting report inputs that cannot resolve actual references.
- **Study / Source Grounding:** TRD section 16 defines citation syntax, source table parsing, resolution order, and output fields. `technical-handbook.md` highlights source table parsing as a trust boundary and basename fallback as useful but dangerous. The performance report supports bounded lookup via gdu `pkg/analyze/parallel.go:13` and yq streaming/bounded processing references. The security report supports explicit path discipline and trust boundaries.
- **Trade-Offs Accepted:** Workspace-relative source paths take precedence over report-relative paths to match generated tables like `sources/source-name`; report-relative fallback supports moved or copied reports. Basename fallback improves legacy compatibility but may mark ambiguous references unresolved.
- **Technical Debt / Future Impact:** Parser compatibility may need expansion if generated reports vary table headings or use Markdown links. This sprint deliberately limits support to required forms.
- **Alternatives Rejected:** Treating missing source tables as success when supported citations exist is rejected because references would silently disappear. Searching the whole workspace by basename is rejected because it risks unrelated matches and large scans. Resolving remote URLs or cloning missing sources is rejected because sprint scope is local-only.
- **Contracts Satisfied:** Security, Performance, Errors, Testing, AC-13 through AC-19, AC-24, AC-29, AC-31.
- **Evidence Required:** Parser and resolver tests for valid source rows, workspace-relative paths, report-relative paths, missing table when refs exist, no-ref reports, single-line/range/list citations, legacy dash normalization, source-prefixed stripping, basename fallback, missing files, ambiguous basename, malformed line specs, path escapes, and out-of-range lines.

### Decision 5: Deterministic Extraction Output And Exit Outcomes

- **Decision:** Default text output must include report path, source name, cited path, requested normalized line spec, resolved path when found, line-numbered snippets, and an unresolved summary. `--json` must emit deterministic sprint-local structured fields for reports, sources, references, status, snippets, and unresolved diagnostics. Multiple reports are processed in argument order. Exit mapping is: success when all references resolve; partial completion exit code `8` when inspection succeeds but one or more references remain unresolved; validation exit code `5` for malformed report inputs such as required-but-missing source table or malformed source rows; usage exit code `2` for invalid arguments; workspace/filesystem exit code `4` for unreadable inputs or output write failures.
- **Rationale:** This preserves useful extraction results while still making unresolved references visible to scripts. JSON determinism supports automation and tests without making a release-wide compatibility promise.
- **Study / Source Grounding:** TRD section 7.2 defines suggested exit classes, including `5` validation failure and `8` partial completion. `technical-handbook.md` identifies unresolved references as validation-style outcomes and asks whether partial completion or validation should represent unresolved refs. The error-handling report supports distinct exit mapping with gh-cli `internal/ghcmd/cmd.go:44-49` and rclone `cmd/cmd.go:497-516`. The IO and testing reports support injectable output and exact text/JSON tests.
- **Trade-Offs Accepted:** Unresolved references use partial completion rather than validation failure when other report inspection completed. Missing source tables use validation failure because the report cannot satisfy a required resolution precondition.
- **Technical Debt / Future Impact:** Scripts may initially depend on sprint-local JSON field names. Future stable schema work must either preserve these fields or version changes deliberately.
- **Alternatives Rejected:** Returning success with unresolved refs is rejected because it hides audit failures. Failing fast on the first unresolved ref is rejected because it prevents inspection of later references and reports. Treating all extraction problems as exit `1` is rejected because requirements and TRD call for scriptable outcome classes.
- **Contracts Satisfied:** CLI Surface, Errors, Documentation, Testing, Security, AC-12, AC-20 through AC-25, AC-29 through AC-31.
- **Evidence Required:** Command tests for help, no reports usage error, multiple report argument order, all-resolved success, unresolved partial exit, malformed/missing table validation exit, JSON exact output, text exact output, `--output` writes, output write failure, and safe diagnostic content review.

### Decision 6: Local-Only Safe Path Handling And Bounded Fallback Search

- **Decision:** Code extraction may read only local files under parsed and contained source roots after normalization. Path specs containing `..` escapes, absolute paths outside accepted roots, or symlink-resolved escapes are unresolved validation diagnostics, not silently followed. Basename fallback is scoped to parsed source roots, deterministic, cached within one extraction command invocation, and skips this fixed ignored set: `.git`, `.hg`, `.svn`, `node_modules`, and `.ultraplan`.
- **Rationale:** Reports are generated or hand-edited Markdown and must be treated as untrusted input. Local-only containment prevents a citation from reading arbitrary files, while a fixed ignored set keeps fallback lookup predictable and bounded.
- **Study / Source Grounding:** Requirements AC-18, AC-29, and constraints require local root containment and ignored directories. `technical-handbook.md` warns not to treat path strings as trusted and not to allow unbounded fallback scans. The security report cites explicit trust-boundary and path/command discipline, including opencode permission boundaries and restic secret redaction. The performance report supports bounded traversal and cache reuse.
- **Trade-Offs Accepted:** Some legitimate files inside ignored directories will not be found by basename fallback, though direct source-relative citations can still resolve if under the root and allowed. The ignore list is fixed rather than configurable to avoid new config scope.
- **Technical Debt / Future Impact:** Ecosystem-specific ignore configuration may become useful if false negatives appear. That should be a future requirements decision, not an implicit config addition now.
- **Alternatives Rejected:** Following absolute paths from reports is rejected because it violates source isolation. Scanning parent directories or the full workspace is rejected because it risks secret leakage and unbounded cost. Configurable ignore rules are rejected for this sprint because configuration changes were explicitly excluded by `sprint-index.md`.
- **Contracts Satisfied:** Security, Performance, Errors, Testing, AC-17 through AC-20, AC-24, AC-29, AC-31.
- **Evidence Required:** Resolver tests for path escape attempts, symlink/cleaned-path containment where supported by current path helpers, ignored directories, ambiguous matches, deterministic ordering, and no reads outside source roots; code review of filesystem traversal and diagnostics.

### Decision 7: Verification Scope And Architecture Review Gates

- **Decision:** Implementation must ship with deterministic unit, fixture, command-level, and exact-output tests for both capabilities; normal verification is `go test ./...` and `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`. Review must inspect package boundaries, runtime separation, safe diagnostics, output determinism, and non-goal leakage before the sprint is considered complete.
- **Rationale:** This sprint is primarily about audit surfaces. Regressions in output shape, warnings, exit outcomes, and path safety are product regressions even when internal tests pass.
- **Study / Source Grounding:** The testing report identifies testscript/golden/fixture-backed tests as top-tier CLI practice, citing chezmoi `internal/cmd/main_test.go:64-174`, helm `internal/test/test.go:43`, and go-task `task_test.go:166-169`. The architecture and command reports support import/path review. The security and logging reports support safe diagnostic review.
- **Trade-Offs Accepted:** Exact-output tests can make intentional formatting changes noisier. This is acceptable because deterministic output is a sprint requirement and generated artifacts must be reviewable.
- **Technical Debt / Future Impact:** If output grows complex, tests may need helper normalizers. Avoid broad normalization now because it could mask determinism bugs.
- **Alternatives Rejected:** Testing only internal state is rejected because the sprint outputs are CLI and file artifacts. Relying on real OpenCode or network smoke tests is rejected because normal tests must be offline. Skipping build verification is rejected because the CLI binary is a required output.
- **Contracts Satisfied:** Testing, Architecture, CLI Surface, Security, Documentation, AC-30, AC-31, review expectations in requirements lines 112-127.
- **Evidence Required:** `go test ./...`; `go build ./cmd/ultraplan`; exact-output fixtures or assertions for CSV/text/JSON/help; architecture review notes; sprint review evidence recorded in `review.md`.

## Expected Evidence

| Evidence Type | Required Evidence | Source / Command / Review Check |
| --- | --- | --- |
| Tests | Summary matrix exact CSV tests for ratings, totals, ordering, missing reports, missing ratings, ambiguous ratings, inapplicable `N/A`, warnings, atomic writes, and repeated-run determinism. | `/home/antonioborgerees/coding/ultraplan-go/internal/study/summary_test.go`; `go test ./...` |
| Tests | Summary command tests for help, success, warnings, missing study/usage failure, write failure, output path, and no runtime invocation. | `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_summary_commands_test.go`; `go test ./...` |
| Tests | Code extraction module tests for source tables, citation syntax, line specs, range normalization, path resolution, basename fallback, ignored dirs, unresolved diagnostics, snippets, JSON/text output. | `/home/antonioborgerees/coding/ultraplan-go/internal/codeextract/codeextract_test.go`; `go test ./...` |
| Tests | Code command tests for help, usage, multiple reports in argument order, `--json`, `--output`, unresolved partial exit, validation exit, filesystem failures, and safe output. | `/home/antonioborgerees/coding/ultraplan-go/internal/app/code_commands_test.go`; `go test ./...` |
| Build | CLI builds successfully. | `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go` |
| Runtime Boundary | No agentwrap, OpenCode, provider, network, or subprocess execution is used by summary or code extraction commands. | Import review, grep/review changed files, command tests with fakes/local fixtures only. |
| Architecture | Summary logic is in `internal/study`; extraction logic is in `internal/codeextract`; `internal/app` is thin; no global parser/resolver/report package is introduced; `internal/platform/runtime` remains product-free. | Architecture review protocol from `system/protocols/architecture-review-protocol.md`. |
| Security | Extraction cannot read outside parsed source roots; ignored dirs are skipped; path escapes and ambiguous basenames are unresolved; outputs omit secrets, env, prompt bodies, and unsafe runtime payloads. | Resolver tests, output exact tests, and code review. |
| Documentation | Help output exposes `ultraplan study <study> summary` and `ultraplan code <report>...`, including `--json` and `--output`. | Help-output command tests and CLI help review. |
| Review | Deviations, residual risks, commands run, and acceptance status are recorded. | `projects/ultraplan-go/sprints/13-summary-and-code-extraction/review.md` after implementation. |

## Assumptions And Risks

| Item | Type | Impact | Mitigation / Follow-Up |
| --- | --- | --- | --- |
| Existing study helpers for discovery, applicability, report paths, and rating parsing are present and reusable. | Assumption | If absent or incompatible, summary implementation may need small refactors before feature work. | Plan must inspect current `internal/study` APIs first and reuse or minimally expose existing helpers. |
| Existing app exit-code conventions can represent usage `2`, filesystem/workspace `4`, validation `5`, and partial `8`. | Assumption | If current app differs, command behavior could become inconsistent. | Plan must inspect current exit mapping and adapt names while preserving chosen semantic categories. |
| Existing generated reports contain source tables close to the TRD row shape. | Assumption | Parser may reject real reports that vary headings or formatting. | Fixture tests should include required shape plus reasonable whitespace/backtick variations; broader compatibility remains future scope. |
| Basename fallback may be slow on very large source roots. | Risk | `ultraplan code` could become sluggish for huge repos. | Cache basename indexes per command invocation, skip ignored dirs, and keep traversal scoped to parsed roots. |
| Basename fallback may produce ambiguous matches frequently. | Risk | Users may need to improve report citations before snippets resolve. | Report ambiguity with source/report context and all candidate-safe relative paths where appropriate. |
| `summary.csv` empty cells are not self-describing. | Risk | Users reading CSV alone may not know whether a cell is missing report, missing rating, or ambiguous rating. | Human output warnings and tests must make reasons visible; future structured summary can add machine-readable reason codes. |
| Output snippets may include sensitive code if reports cite sensitive files. | Risk | The command intentionally prints requested source snippets, which could include secrets present in source code. | Do not print unrelated file content, prompt bodies, env, or runtime payloads; keep snippet extraction limited to requested lines and source roots. |
| No area-specific reasoning files exist. | Risk | Some nuanced decisions are made in this consolidated document instead of area artifacts. | This document records final decisions explicitly enough for `plan.md` to execute without reopening architecture. |

## Implementation Constraints

- Summary behavior must stay in `internal/study`; code extraction behavior must stay in `internal/codeextract`; CLI adapters must stay in `internal/app` and remain thin.
- Do not introduce `internal/summary`, `internal/reports`, `internal/parser`, `internal/resolver`, `internal/validation`, or similar global product-layer packages for this sprint.
- `internal/platform/runtime` must not import or depend on `internal/study`, `internal/codeextract`, or command packages.
- Summary generation must reuse existing study source discovery, applicability filtering, report path helpers, validation/rating parsing behavior, and deterministic ordering conventions.
- Summary must distinguish missing expected reports from inapplicable Markdown source/dimension pairs in code, output, and tests.
- Summary writes must be atomic and explicit; write failures must return non-zero outcomes and must not be downgraded to low-visibility warnings.
- Code extraction must not fetch, clone, or resolve remote repositories or URLs.
- Code extraction must normalize and contain source paths under parsed roots and must reject or mark unresolved any path escape.
- Basename fallback must be deterministic, scoped to parsed source roots, cached within one command invocation, and skip `.git`, `.hg`, `.svn`, `node_modules`, and `.ultraplan`.
- Text and JSON output must be deterministic and must not include secrets, sensitive environment values, prompt bodies, embedded Markdown document bodies unrelated to requested snippets, or unsafe native runtime payloads.
- Normal tests must use local fixtures only and must not require OpenCode, provider credentials, network access, or real runtime execution.
- `go test ./...` and `go build ./cmd/ultraplan` are required verification commands for the sprint review.
- `plan.md` must not reopen module ownership, summary semantics, extraction resolution order, ignored-directory set, exit outcome categories, or non-runtime scope.

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
- [x] Area-specific reasoning documents were checked; none exist for this sprint, and that absence is explicitly recorded.
- [x] Area-specific reasoning conclusions are reflected or explicitly unavailable.
- [x] Contracts are explicitly mapped to decisions and expected evidence.
- [x] Final decisions are clear enough for `plan.md` to execute.
- [x] Expected evidence is specific and reviewable.
