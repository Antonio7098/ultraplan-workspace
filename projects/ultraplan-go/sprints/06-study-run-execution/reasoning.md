# Sprint Reasoning: 06 Study Run Execution

> Project: `ultraplan-go`
> Sprint: `06-study-run-execution`
> Output: `projects/ultraplan-go/sprints/06-study-run-execution/reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/sprints/06-study-run-execution/requirements.md`, `projects/ultraplan-go/sprints/06-study-run-execution/reasoning.md`, `projects/ultraplan-go/sprints/06-study-run-execution/sprint-index.md`, `projects/ultraplan-go/sprints/06-study-run-execution/technical-handbook.md`, `projects/ultraplan-go/sprints/06-study-run-execution/reasoning/*.md` (not present), `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/project-index.md`, `templates/sprint-reasoning.md`

This document decides. It synthesizes selected context, handbook evidence, and contracts into final sprint decisions.

It does not replace `sprint-index.md`, `technical-handbook.md`, or future `reasoning/*.md` artifacts.

## Sprint Purpose

- **Goal:** Make study reports first-class validated artifacts by implementing study-owned per-source report validation, final report validation, rating parsing, deterministic report path helpers, and actionable validation diagnostics.
- **Non-Goals:** Runtime execution, prompt composition, agentwrap/OpenCode wiring, runtime health checks, policy runner wiring, retry/fallback behavior, event mapping, run-all/run-loop/status/cancellation behavior, durable task state, summary generation, `summary.csv` writing, code reference extraction, target workflows, sprint planning, and sprint execution.
- **Depends On:** Prior study-side sprints for package structure, workspace path resolution, study/source/dimension domain, generated study layout, and Sprint 5 Markdown document source kind and applicability behavior.

## Selected Context And Pre-Reasoning Artifacts

| Artifact | Path | How It Was Used |
| --- | --- | --- |
| Sprint Requirements | `projects/ultraplan-go/sprints/06-study-run-execution/requirements.md` | Binding scope, acceptance criteria, non-goals, implementation paths, dependencies, and review expectations. Requirement IDs below use `AC-01` through `AC-17` in the order listed under acceptance criteria. |
| Project Index | `projects/ultraplan-go/project-index.md` | Confirmed target implementation directory, active contract pool, available evidence reports, review protocols, and project-level non-goals. |
| PRD | `projects/ultraplan-go/docs/PRD.md` | Confirmed product principles: generated reports are durable artifacts, runtime success is not product success, Markdown document sources do not require code exploration, and target/sprint workflows are deferred. |
| TRD | `projects/ultraplan-go/docs/TRD.md` | Confirmed report validator checks, rating parser formats, validation result mapping expectations, source-kind behavior, package ownership, testing requirements, and `go test ./...` acceptance. |
| Architecture | `projects/ultraplan-go/docs/ARCHITECTURE.md` | Confirmed module-driven ownership: `internal/study` owns reports, validation, parsing, summary-facing helpers, and study artifact paths; platform packages must not import study. |
| Sprint Index | `projects/ultraplan-go/sprints/06-study-run-execution/sprint-index.md` | Selected contracts are Architecture, Errors, Security, Testing, and Performance; selected evidence reports are `01`, `05`, `06`, `10`, `11`, `13`, `14`, and `15`; excluded runtime, workflow, persistence, CLI-surface, and documentation context is not sprint scope. |
| Technical Handbook | `projects/ultraplan-go/sprints/06-study-run-execution/technical-handbook.md` | Evidence basis for study-owned validation, typed diagnostics, preserved error causes, source-kind-aware rules, bounded artifact processing, safe diagnostics, and fixture-first tests. |
| Sprint Reasoning Template | `templates/sprint-reasoning.md` | Structure for final decisions, evidence mapping, trade-offs, risks, constraints, and plan handoff. |

## Area-Specific Reasoning Inputs

No area-specific reasoning files exist under `projects/ultraplan-go/sprints/06-study-run-execution/reasoning/*.md` at the time this document was created. The sprint index selected an architecture reasoning output path for report validation, but no file is present, so this consolidated reasoning makes the final decisions directly from the requirements, project docs, sprint index, and technical handbook.

## Sprint Technical Handbook Summary

- **Relevant Patterns:** Study-owned report validation; validation results as domain values; actionable diagnostics with preserved causes; source-kind-aware validation rules; small filesystem seams; fixture-first validator tests; bounded local artifact processing; safe diagnostics and redaction discipline.
- **Important Trade-Offs:** Keep validators in `internal/study` rather than a global validation package; use structured diagnostics rather than plain error strings; validate citation shape now but defer resolution; keep filesystem seams minimal; parse ratings strictly rather than coercing uncertainty to `0`.
- **Warnings / Anti-Patterns:** Do not parse diagnostics by string; do not panic on malformed report input; do not let direct IO or globals leak into validation core where it harms deterministic tests; do not treat runtime/file-existence success as product success; do not scan whole source repositories during report validation; do not turn inapplicable Markdown pairs into failures.
- **Evidence Confidence:** High for package ownership, typed errors, IO seams, testing, and security because selected reports cite concrete mature Go CLI implementations. Medium for observability, performance, and philosophy because they influence constraints and future suitability but this sprint does not implement runtime logs, status, or batch execution.

## Contracts Applied

| Contract / Requirement ID | Constraint | Decision Impact | Expected Evidence |
| --- | --- | --- | --- |
| Architecture | Product behavior stays with the module that owns the state; no global `validation`, `reports`, or `parsing` packages. | Validation, report paths, rating parsing, diagnostics, and report metadata live in `internal/study`. | Architecture review confirms files are under `internal/study` and platform packages do not import study. |
| Errors | Failures must preserve cause chains and expose actionable diagnostics. | Filesystem and parsing failures are wrapped; validation checks expose name, severity, path, expected, observed, and guidance. | Unit tests inspect diagnostics and error wrapping for missing/unreadable/malformed artifacts. |
| Security | Paths and diagnostics must be safe and avoid leaking secrets or large raw content. | Diagnostics may include report paths and concise observed facts, but not full report dumps, environment, credentials, or runtime payloads. | Review checks diagnostic fields and tests cover path-bearing failures. |
| Testing | Behavior must be deterministic and fixture-testable without real runtime dependencies. | Validators and rating parser are tested with local fixtures/temp files and table-driven cases. | `go test ./...` passes; tests cover valid, missing, empty, malformed, rating, citation, Markdown exemption, and final report cases. |
| Performance | Validation must avoid unnecessary scans and unbounded work. | Validators read expected report files and inspect content directly; directory citation checks validate shape only. | Review confirms no repository traversal or citation resolution is added. |
| AC-01, AC-02 | Per-source reports fail when missing, empty, unreadable, or structurally incomplete. | Per-source validator performs required file and section checks and includes report path in diagnostics. | `validation_test.go` cases for missing, empty, unreadable if feasible, missing sections, valid report. |
| AC-03, AC-04, AC-05, AC-14 | Directory and Markdown document sources have different citation obligations, and inapplicable Markdown pairs remain caller-skipped. | Directory sources require file-line citation shape unless dimension disables it; Markdown sources do not require code citations by default and record source kind. | `validation_test.go` cases for directory citation failure, Markdown exemption, Markdown required rating/summary/answers, and no validation failure for skipped pairs. |
| AC-06, AC-07 | Final reports must fail when missing, empty, unreadable, or lacking required synthesis sections. | Final report validator checks expected artifact path and required final report sections. | `validation_test.go` cases for valid final report, missing final report, empty final report, and missing required final sections. |
| AC-08, AC-09, AC-10 | Rating parsing accepts specified formats and distinguishes missing, invalid, and ambiguous values without invented scores. | Rating parser returns explicit result state and diagnostics rather than coercing to zero. | `rating_test.go` cases for `**8 / 10**`, `8/10`, `Rating: 8`, missing, invalid, ambiguous, normalization. |
| AC-11, AC-12, AC-13 | Validation APIs expose structured checks, preserve causes, and remain usable by later run/synthesis/status/summary/validate commands without runtime execution. | API returns domain validation results and checks independent of runtime. | Tests inspect result fields; review confirms no dependency on agentwrap/OpenCode runtime execution. |
| AC-15, AC-16, AC-17 | Required unit coverage and full test suite pass. | Plan must include focused tests and `go test ./...` verification in `/home/antonioborgerees/coding/ultraplan-go`. | Passing test output recorded in review. |

## Repos Studied / Source Evidence Used

| Source / Repo / Report | Concrete Reference | Relevant Finding | Why It Matters For This Sprint | Used In Decision(s) |
| --- | --- | --- | --- | --- |
| `go-cli-study / 01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md`; `cmd/root.go:112-127`, `internal/restic/repository.go:18`, `internal/chezmoi/chezmoi.go:1-2` | Mature Go CLIs keep entrypoints thin and product behavior inside protected internals with one-way dependency flow. | Supports keeping report validation in `internal/study` and avoiding global technical-layer packages. | Decision 1 |
| `go-cli-study / 05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md`; `rclone/vfs/zip.go:21`, `gh-cli/pkg/ssh/ssh_keys.go:64`, `helm/pkg/storage/driver/driver.go:27-48`, `age/cmd/age/tui.go:47-54` | Mature CLIs wrap causes, use typed/sentinel checks, and provide actionable guidance. | Supports structured validation checks and preserved filesystem/parsing causes. | Decisions 3, 4 |
| `go-cli-study / 06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md`; `pkg/iostreams/iostreams.go:551-568`, `internal/chezmoi/system.go:25`, `internal/fs/interface.go:10-31` | IO seams should be focused around volatile boundaries and test needs. | Supports minimal report file reading seams or temp-file fixtures rather than a broad filesystem abstraction. | Decisions 1, 5 |
| `go-cli-study / 10-logging-observability` | `studies/go-cli-study/reports/final/10-logging-observability.md`; `internal/logging/logging.go:31-66`, `internal/slogs/keys.go:6-231`, `pkg/iostreams/iostreams.go:52-54` | Stable structured fields make diagnostics and status machine-usable later. | Supports validation check fields that future status/validate commands can surface without parsing text. | Decision 3 |
| `go-cli-study / 11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md`; `chezmoi/internal/cmd/main_test.go:64-174`, `go-task/task_test.go:166-169`, `helm/internal/test/test.go:43` | Table-driven tests, fixtures, and golden-style checks are common in robust CLIs. | Supports fixture/table-driven validator and rating parser tests. | Decision 5 |
| `go-cli-study / 13-security` | `studies/go-cli-study/reports/final/13-security.md`; `internal/options/secret_string.go:15-20`, `pkg/registry/transport.go:37-41`, `internal/config/json/validator.go:146` | Explicit trust boundaries, path validation, schema validation, and redaction reduce leakage and unsafe operations. | Supports safe diagnostic contents and source-kind trust boundary for Markdown document sources. | Decisions 2, 3 |
| `go-cli-study / 14-performance` | `studies/go-cli-study/reports/final/14-performance.md`; `gh-cli/pkg/cmdutil/factory.go:27-42`, `yq/pkg/yqlib/stream_evaluator.go:78-113`, `gdu/pkg/analyze/parallel.go:13` | Mature tools avoid eager scans and unbounded processing. | Supports validating expected report artifacts only and deferring citation resolution/code extraction. | Decisions 2, 5 |
| `go-cli-study / 15-philosophy` | `studies/go-cli-study/reports/final/15-philosophy.md`; `pkg/cmd/factory/default.go:26-46`, `internal/chezmoi/system.go:25-45`, `VISION.md:97` | Good Go CLIs accept deliberate complexity only where it preserves maintainability and explicit non-goals. | Supports a small but typed validation surface and strict scope boundaries. | Decisions 1, 3, 5 |

## Trade-Off And Debt Analysis

### Accepted Trade-Offs

| Trade-Off | Benefit | Cost / Constraint Accepted | Why Acceptable Now | Revisit Trigger |
| --- | --- | --- | --- | --- |
| Keep all report validation behavior in `internal/study` | Matches architecture and report semantics stay near study domain state. | Similar validation result shapes may be duplicated later if other modules need validators. | This sprint validates study reports only; target/sprint validators are non-goals. | Another product module needs the same validation primitives with concrete duplication. |
| Structured validation results instead of returning only `error` | Later run, synthesis, status, summary, and validate flows can inspect check fields. | More domain types and more precise tests are required. | Acceptance criteria require deterministic diagnostics and future API reuse. | Public stable JSON validation output requires schema versioning or different field names. |
| Citation shape validation only | Satisfies directory-source evidence requirement without repository traversal. | Cannot prove cited files exist, line numbers are valid, or snippets match claims. | Code-reference extraction is explicitly deferred and performance constraints prohibit broad scans. | Code extraction sprint begins or citation resolution becomes validation scope. |
| Markdown sources exempt from code citations by default | Preserves Sprint 5 document-source semantics and avoids invalid repository assumptions. | Document reports may contain weaker evidence unless dimensions specify document citations. | PRD/TRD explicitly distinguish Markdown document sources from directory sources. | Dimension schema gains explicit document citation requirements. |
| Strict rating ambiguity handling | Prevents false numeric scores and protects later summary generation. | Some generated reports will fail or warn even when a human could infer intent. | Missing, invalid, ambiguous, and valid scores are distinct states in requirements. | A repair/suggestion sprint adds bounded automatic clarification. |
| Minimal filesystem seams and local fixtures | Keeps implementation small and deterministic. | Some tests may use real temporary files instead of pure in-memory filesystems. | Validation reads known local artifacts only and does not need a broad storage abstraction. | Validation begins reading remote stores, virtual filesystems, or many generated artifact types. |

### Potential Technical Debt

| Debt / Shortcut | Why It Might Accrue | Current Mitigation | Owner / Follow-Up |
| --- | --- | --- | --- |
| Section heading matching may be heuristic | Generated Markdown headings can vary across reports. | Define a small explicit set of accepted heading variants in tests and avoid `TBD` behavior. | Future report template/version sprint if formats stabilize. |
| Dimension disables code citations may be minimal | Current dimension model may not have a dedicated boolean for citation policy. | Use existing dimension citation content or a narrowly scoped helper without broad schema changes. | Prompt/report template or dimension schema sprint. |
| Validation result fields may become public JSON later | Internal names may not be ideal for stable external APIs. | Keep any schema internal/test-only unless explicitly versioned. | Public `study validate --json` sprint. |
| Rating parser supports only required formats | Other plausible rating formats may appear in generated reports. | Reject/diagnose unsupported or ambiguous values rather than inventing scores. | Summary generation or repair sprint can add formats with fixtures. |
| Final report validation checks presence rather than semantic completeness | Local validation cannot verify truth or synthesis quality. | Require core sections and rating summary while keeping semantic correctness outside validation scope. | Human review workflow or evaluation sprint. |

### Future Considerations

| Consideration | Deferred Until | Reason Deferred | What Should Be Preserved Now |
| --- | --- | --- | --- |
| Runtime execution integration | Single-run execution sprint | Runtime and agentwrap wiring are explicit non-goals. | Validation APIs must be callable without runtime and later embeddable as agentwrap validators. |
| Stable JSON validation output | Public validate command sprint | CLI surface/stable output contract is excluded. | Validation results should be structured and deterministic internally. |
| Citation resolution and snippet extraction | Code extraction sprint | Current sprint checks citation shape only. | Keep citation-shape regex/parser narrow and not coupled to filesystem resolution. |
| Summary CSV generation | Summary sprint | This sprint only parses ratings and validates reports. | Rating parser must preserve missing/invalid/ambiguous states for later summary code. |
| Repair prompts | Repair/runtime sprint | Prompt composition and runtime repair are excluded. | Diagnostics should include safe corrective guidance and expected/observed facts. |
| Persisted validation schemas | Durable state or public API sprint | Persistence/migration contract is not selected. | Do not persist unversioned machine-readable schemas as product artifacts. |

## Final Decisions

### Decision 1: Keep Validation, Report Paths, Ratings, And Diagnostics Study-Owned

- **Decision:** Implement report metadata/domain extensions in `internal/study/domain.go`, deterministic report path helpers in `internal/study/reports.go`, report validators in `internal/study/validation.go`, and rating parsing in `internal/study/rating.go`. Do not create global `internal/validation`, `internal/reports`, or `internal/parsing` packages.
- **Rationale:** Study reports are study-owned durable artifacts. The validators need study concepts such as source kind, dimensions, source names, report folders, applicability, and expected report paths, so moving them to global technical layers would fracture product context.
- **Study / Source Grounding:** `technical-handbook.md` relevant patterns and trade-offs; `go-cli-study / 01-project-structure` evidence for thin entrypoints and internal product behavior (`cmd/root.go:112-127`, `internal/restic/repository.go:18`, `internal/chezmoi/chezmoi.go:1-2`); `ARCHITECTURE.md` lines 75-87 and 139-158 identify `internal/study` as owner of reports and validation.
- **Trade-Offs Accepted:** Future modules may need similar validation shapes, but no other module has this concrete need in Sprint 6.
- **Technical Debt / Future Impact:** If codeextract or future target/sprint modules need shared result primitives, extract only proven generic pieces later with a clear owner and versioning story.
- **Alternatives Rejected:** A global `internal/validation` package was rejected because it violates sprint constraints and would separate report semantics from study state. A broad `internal/reports` package was rejected because report path and metadata rules are currently study-specific. Adding agentwrap validators as the only API was rejected because this sprint must not depend on runtime execution.
- **Contracts Satisfied:** Architecture; AC-01 through AC-17 where implementation paths and API ownership are relevant; requirements constraints on package placement and runtime independence.
- **Evidence Required:** Architecture review confirms implementation files are under `internal/study`; import review confirms platform packages do not import study; `go test ./...` passes.

### Decision 2: Validate Per-Source And Final Reports As Local Artifacts With Source-Kind-Aware Rules

- **Decision:** Add per-source validators that check expected file existence/readability, non-empty content, top-level heading, source information, summary, rating section, parseable rating, question/answer content, and citation shape when required. Add final report validators that check expected file existence/readability, non-empty content, study parameters or equivalent context, sources studied table, executive summary, rating summary, pattern/synthesis content, and open questions or notable absences.
- **Rationale:** Product success requires valid report artifacts, not just file existence or future runtime success. Directory and Markdown document sources have different evidence obligations, so citation validation must be source-kind-aware.
- **Study / Source Grounding:** `technical-handbook.md` source-kind-aware validation rules and bounded artifact processing; PRD report validation lines 259-263 and Markdown source requirements lines 279-287 and 492-496; TRD report validators lines 1299-1334; security evidence from `go-cli-study / 13-security`; performance evidence from `go-cli-study / 14-performance`.
- **Trade-Offs Accepted:** Citation validation checks shape only, not file existence or snippet correctness. Final report validation checks required structure, not semantic quality.
- **Technical Debt / Future Impact:** Citation resolution remains deferred to code extraction; semantic review remains human/product review. Final report section matching must be explicit enough to avoid unpredictable generated-heading failures.
- **Alternatives Rejected:** Validating only file existence was rejected because PRD/TRD require required sections, rating, and citation discipline. Requiring code citations for Markdown document sources by default was rejected because Sprint 5 and PRD/TRD distinguish document sources from repositories. Resolving all citations during validation was rejected because code extraction is deferred and performance constraints prohibit repository traversal.
- **Contracts Satisfied:** Architecture, Security, Performance, Testing; AC-01 through AC-07, AC-14, AC-15.
- **Evidence Required:** `internal/study/validation_test.go` covers valid per-source reports, missing files, empty files, missing sections, missing ratings, directory citation failures, Markdown citation exemptions, valid final reports, missing/empty final reports, and missing final sections.

### Decision 3: Return Structured Validation Results With Deterministic Actionable Diagnostics

- **Decision:** Define validation status/result/check concepts that expose pass/fail status, check name, severity, artifact path, expected value, observed value, source kind where relevant, and safe corrective guidance. Filesystem and parsing failures must preserve error cause chains with wrapping while still producing validation diagnostics.
- **Rationale:** Later single-run, batch, synthesis, summary, status, and validate flows need inspectable validation results. Tests need deterministic fields, not fragile string matching.
- **Study / Source Grounding:** `technical-handbook.md` patterns for validation results as domain values and actionable diagnostics with preserved causes; `go-cli-study / 05-error-handling` references for `%w`, typed/sentinel errors, and user guidance; `go-cli-study / 10-logging-observability` references for stable structured fields; TRD validation failure behavior lines 1355-1366.
- **Trade-Offs Accepted:** More domain types and assertions are required now, but this avoids later rewriting string errors into structured diagnostics.
- **Technical Debt / Future Impact:** The result shape is internal for now. If exposed as stable JSON later, add explicit schema/versioning rather than treating current Go structs as a public contract by accident.
- **Alternatives Rejected:** Plain `error` strings were rejected because callers would need to scrape text and could not distinguish missing, invalid, ambiguous, or skipped checks. Panicking on malformed report input was rejected because malformed generated content is expected operational input. Dumping raw report excerpts into diagnostics was rejected for security/noise reasons.
- **Contracts Satisfied:** Errors, Security, Testing; AC-01, AC-04, AC-10, AC-11, AC-12, AC-13.
- **Evidence Required:** Tests assert check names, severities, paths, expected/observed values, guidance where safe, source kind for Markdown diagnostics, and preserved wrapped errors for filesystem/parsing failures.

### Decision 4: Parse Ratings Strictly And Preserve Uncertainty

- **Decision:** Implement a rating parser that accepts `**8 / 10**`, `8/10`, and `Rating: 8`, normalizes valid numeric scores, and returns distinct states for valid, missing, invalid, and ambiguous ratings. Ambiguous ratings must produce diagnostics and must not invent or coerce a score, including to `0`.
- **Rationale:** Later summary generation depends on distinguishing missing values from real numeric scores. Ambiguity should trigger repair or human correction, not silently degrade product metrics.
- **Study / Source Grounding:** `technical-handbook.md` strict rating ambiguity handling; TRD rating parser requirements lines 1438-1445; PRD summary requirements lines 272-278 distinguish missing scores from zero. Error-handling evidence supports typed parse outcomes instead of string-only failures.
- **Trade-Offs Accepted:** Some reports that a human could interpret may fail validation or warn until repaired. This is preferable to corrupting downstream summaries.
- **Technical Debt / Future Impact:** Additional rating formats can be added later with fixtures when observed in real reports. The parser should stay independent of summary CSV writing.
- **Alternatives Rejected:** Coercing missing or ambiguous ratings to `0` was rejected because it violates requirements and misrepresents product quality. Selecting the first rating found was rejected because multiple scores can be ambiguous. Making rating parsing part of summary generation was rejected because validators need it now and summary generation is out of scope.
- **Contracts Satisfied:** Errors, Testing; AC-02, AC-05, AC-08, AC-09, AC-10, AC-15.
- **Evidence Required:** `internal/study/rating_test.go` covers all accepted formats, score normalization, missing ratings, invalid ratings, ambiguous ratings, and no invented score.

### Decision 5: Keep Validation Bounded, Deterministic, And Runtime-Free

- **Decision:** Validators read known expected report paths and inspect local Markdown content only. They do not invoke OpenCode, agentwrap, network access, provider credentials, prompt composition, run state, worker pools, source repository scans, citation resolution, summary writing, or public CLI output behavior.
- **Rationale:** The sprint folder name includes execution, but the authoritative Sprint 6 scope is report validation and rating parsing. Keeping validation runtime-free makes it reusable by later runtime and workflow commands without reopening architecture.
- **Study / Source Grounding:** `technical-handbook.md` design pressures and anti-patterns; `go-cli-study / 14-performance` bounded processing references; `go-cli-study / 11-testing-strategy` fixture-first tests; PRD principle that runtime success is not product success; sprint-index excluded context lines 73-95.
- **Trade-Offs Accepted:** This sprint will not prove runtime integration, repair, status display, persisted validation records, or full citation correctness.
- **Technical Debt / Future Impact:** Later runtime integration must adapt these validators into agentwrap validation expectations or task-state validation without changing the validation semantics.
- **Alternatives Rejected:** Implementing `ultraplan study <study> run` was rejected because runtime execution is a non-goal. Adding a public stable `study validate --json` command was rejected because CLI Surface and Documentation contracts are excluded unless already present. Scanning all reports/repositories to discover validation work was rejected because expected paths are deterministic and batch discovery belongs to later workflow code.
- **Contracts Satisfied:** Architecture, Performance, Testing, Security; AC-13, AC-14, AC-15, AC-16, AC-17.
- **Evidence Required:** Tests run without OpenCode, agentwrap runtime execution, network, or provider credentials; review confirms no run-loop, worker pool, prompt, code extraction, summary, or public JSON validation scope was added.

## Expected Evidence

| Evidence Type | Required Evidence | Source / Command / Review Check |
| --- | --- | --- |
| Tests | Per-source validator tests for valid reports, missing files, empty files, unreadable files if portable, missing required sections, missing question/answer content, missing rating, ambiguous rating, directory citation failure, Markdown citation exemption, and Markdown rating/summary/answers requirements. | `/home/antonioborgerees/coding/ultraplan-go/internal/study/validation_test.go` |
| Tests | Final report validator tests for valid final report, missing file, empty file, unreadable file if portable, missing study context, missing sources table, missing executive summary, missing rating summary, missing synthesis content, and missing open questions/notable absences. | `/home/antonioborgerees/coding/ultraplan-go/internal/study/validation_test.go` |
| Tests | Rating parser tests for `**8 / 10**`, `8/10`, `Rating: 8`, missing rating, invalid value, ambiguous values, and normalized score output. | `/home/antonioborgerees/coding/ultraplan-go/internal/study/rating_test.go` |
| Tests | Full Go test suite passes. | Run `go test ./...` in `/home/antonioborgerees/coding/ultraplan-go`. |
| Review | Architecture boundaries remain valid: validation/rating/report helpers are in `internal/study`; no global validation/report/parsing package; platform packages do not import study. | Architecture Review protocol selected by `sprint-index.md`. |
| Review | Scope boundaries remain valid: no runtime execution, prompt composition, run-loop, summary generation, code extraction, target/sprint workflow, stable public validate JSON, or durable workflow behavior added. | Sprint Review protocol and requirements non-goals. |
| Review | Diagnostics are deterministic and safe: include paths/checks/expected/observed/guidance, preserve causes, avoid raw secrets or large report dumps. | Errors, Security, Testing contract checks. |
| Documentation | Sprint artifacts reflect final decisions and plan can execute without reopening architecture. | `reasoning.md`, `plan.md`, `review.md`. |

## Assumptions And Risks

| Item | Type | Impact | Mitigation / Follow-Up |
| --- | --- | --- | --- |
| Existing study domain types from prior sprints include enough source kind, dimension, and path data for validators. | Assumption | If not, small domain additions are needed. | Add only the fields required by validation and keep them in `internal/study/domain.go`. |
| Dimension-level disabling of directory citations can be represented without a broad schema redesign. | Assumption | If not, implementation could overreach into dimension schema or prompt design. | Use a narrow helper based on existing dimension citation expectations; defer richer schema to a later sprint. |
| Markdown heading variants can be kept explicit and deterministic. | Risk | Overly loose matching may pass weak reports; overly strict matching may fail reasonable generated reports. | Define accepted variants in tests and document them in validator code or fixtures. |
| Portable unreadable-file tests may be OS-dependent. | Risk | Permission behavior can vary across environments. | Prefer injectable read failure seam if already natural; otherwise cover missing/empty/read errors without brittle platform assumptions. |
| Validation checks could accidentally become user-facing stable JSON. | Risk | Later compatibility burden without schema review. | Keep machine-readable structures internal unless a public validate command explicitly versions them. |
| Area reasoning file expected by requirements is absent. | Risk | Some architecture nuance might not be pre-analyzed separately. | This consolidated reasoning records final decisions directly and plan should not require another reasoning pass. |
| Inapplicable Markdown source/dimension pairs are caller-skipped. | Assumption | Validators receiving an inapplicable expected path could incorrectly fail it as missing. | Preserve AC-14: callers exclude skipped pairs; validators should not create failures for pairs explicitly marked skipped. |

## Implementation Constraints

- Keep all new report validation, rating parsing, report path helper, and validation domain behavior under `/home/antonioborgerees/coding/ultraplan-go/internal/study`.
- Do not create global `internal/validation`, `internal/reports`, or `internal/parsing` packages for this sprint.
- Do not add runtime execution, OpenCode invocation, agentwrap runtime composition, prompt composition, run-loop, worker pool, status, repair, retry, fallback, summary generation, code extraction, target, or sprint execution behavior.
- Validate expected report paths only; do not recursively scan source repositories or resolve citation targets.
- Directory source reports require file-path-and-line citation shape by default unless a dimension explicitly disables code citation requirements.
- Markdown document source reports do not require code citations by default, must still require rating/summary/answers, and diagnostics must include source kind.
- Rating parsing must never coerce missing, invalid, or ambiguous ratings into `0` or any invented score.
- Validation diagnostics must be deterministic, testable, actionable, and safe to display.
- Filesystem and parsing errors must preserve cause chains where an underlying error exists.
- Inapplicable Markdown source/dimension pairs remain skipped or excluded by callers and must not become validation failures.
- Any schema-like machine-readable structure introduced in this sprint must remain internal/test-only unless explicitly versioned.
- Normal verification is local and runtime-free; tests must not require OpenCode, network, provider credentials, or generated model output.

## Plan Handoff

`plan.md` must execute these decisions. It must not invent architecture, scope, or decisions beyond this document.

The plan must carry forward:

- Final decisions for study-owned validators, source-kind-aware report checks, structured diagnostics, strict rating parsing, and runtime-free bounded validation.
- Applicable contracts and requirement IDs: Architecture, Errors, Security, Testing, Performance, AC-01 through AC-17.
- Expected evidence: focused validator/rating tests, architecture review checks, scope review checks, and `go test ./...`.
- Risks and assumptions around heading matching, citation-disable representation, unreadable-file test portability, and internal-only validation result schemas.
- Required review protocols: Architecture Review and Sprint Review from `sprint-index.md`.

## Phase Exit Criteria

- [x] Selected context was read and used.
- [x] Area-specific reasoning documents were completed where applicable or explicitly marked not applicable.
- [x] Area-specific reasoning conclusions are reflected or explicitly noted as absent.
- [x] Contracts are explicitly mapped to decisions and expected evidence.
- [x] Final decisions are clear enough for `plan.md` to execute.
- [x] Expected evidence is specific and reviewable.
