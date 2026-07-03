# Sprint Requirements: 07 Prompt Composition Generation

> Project: `ultraplan-go`
> Sprint: `07-prompt-composition-generation`
> Purpose: the authoritative, human-readable sprint contract. All other sprint artifacts must satisfy these requirements.

## Sprint Goal

Implement deterministic study prompt composition for directory analysis, Markdown document analysis, and synthesis, including manifest-backed prompt previews, without invoking agentwrap or OpenCode.

## Required Outputs

| Output | Path | Description |
| --- | --- | --- |
| Sprint requirements | `projects/ultraplan-go/sprints/07-prompt-composition-generation/requirements.md` | Defines the binding goal, scope, acceptance criteria, constraints, dependencies, and review expectations for this sprint. |
| Sprint index | `projects/ultraplan-go/sprints/07-prompt-composition-generation/sprint-index.md` | Selects the roadmap entries, source docs, evidence reports, contracts, and review protocols that govern this sprint. |
| Technical handbook | `projects/ultraplan-go/sprints/07-prompt-composition-generation/technical-handbook.md` | Gives implementation guidance for prompt template loading, analysis prompt builders, synthesis prompt builders, manifests, previews, and tests. |
| Prompt composition reasoning | `projects/ultraplan-go/sprints/07-prompt-composition-generation/reasoning/prompt-composition.md` | Records reasoning for prompt ownership, template boundaries, directory-vs-Markdown source rules, synthesis manifest shape, and preview behavior. |
| Consolidated reasoning | `projects/ultraplan-go/sprints/07-prompt-composition-generation/reasoning.md` | Summarizes accepted sprint decisions, tradeoffs, risks, and requirement mappings. |
| Implementation plan | `projects/ultraplan-go/sprints/07-prompt-composition-generation/plan.md` | Provides the ordered task plan and verification checklist for implementation. |
| Sprint review | `projects/ultraplan-go/sprints/07-prompt-composition-generation/review.md` | Reviews completed implementation against these requirements and records carry-forward decisions. |
| Study prompt domain updates | `/home/antonioborgerees/coding/ultraplan-go/internal/study/domain.go` | Defines or extends prompt request, prompt kind, prompt manifest, template reference, and preview result concepts needed by study prompt composition. |
| Study prompt composition implementation | `/home/antonioborgerees/coding/ultraplan-go/internal/study/prompts.go` | Loads base/report/synthesis templates and builds deterministic directory analysis, Markdown document analysis, and synthesis prompts with input manifests. |
| Study prompt preview CLI wiring | `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_prompt_commands.go` | Adds a runtime-free prompt preview command or dry-run command wiring for rendering composed prompts through the existing CLI app. |
| Study prompt composition tests | `/home/antonioborgerees/coding/ultraplan-go/internal/study/prompts_test.go` | Tests template loading, deterministic prompt text, manifests, directory rules, Markdown stripping/isolation rules, synthesis manifests, and missing-input failures. |
| Study prompt preview command tests | `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_prompt_commands_test.go` | Tests prompt preview CLI behavior, output paths or stdout behavior, non-zero errors, and absence of runtime invocation. |

## Acceptance Criteria

- [ ] Prompt builders load required workspace prompt and report template files from the resolved workspace using deterministic paths.
- [ ] Missing `prompts/base.md`, `prompts/synthesize.md`, `templates/repo-analysis.md`, or final report template inputs fail before any runtime-facing behavior, with diagnostics naming the missing path.
- [ ] Directory analysis prompt composition includes the shared base prompt, selected dimension Markdown, report template, study name, source name, source path, expected output path, and hard source-isolation rules.
- [ ] Directory analysis prompts explicitly require file-path-and-line citations for code claims unless the selected dimension disables code citation requirements.
- [ ] Markdown document analysis prompt composition includes the shared base prompt, selected dimension Markdown, report template, study name, source name, expected output path, and stripped Markdown document content.
- [ ] Markdown document prompts do not include YAML frontmatter from the source document.
- [ ] Markdown document prompts explicitly state that all source material is embedded in the prompt and that the runtime must not inspect external files, repositories, or code.
- [ ] Markdown document prompts explicitly state that code citation requirements do not apply by default unless the dimension states otherwise.
- [ ] Prompt composition skips or reports inapplicable Markdown source/dimension pairs without producing a runnable analysis prompt.
- [ ] Synthesis prompt composition includes the synthesis prompt, selected dimension Markdown, final report template, deterministic manifest of applicable per-source report paths, and expected final output path.
- [ ] Synthesis manifests include only applicable source reports and exclude inapplicable Markdown source/dimension pairs from required inputs.
- [ ] Prompt builders return both composed prompt text and an input manifest identifying prompt kind, study, dimension, source when applicable, source kind when applicable, template paths, input document/report paths, and expected output path.
- [ ] Prompt text and manifest output are deterministic for identical filesystem inputs.
- [ ] Prompt preview can print to stdout or write to an explicit output path without invoking agentwrap, OpenCode, provider credentials, network access, or subprocess runtime execution.
- [ ] Prompt preview errors are actionable for unknown studies, unknown dimensions, unknown sources, ambiguous references, inapplicable Markdown pairs, missing templates, missing Markdown document files, and missing applicable per-source reports for synthesis.
- [ ] Unit tests cover directory analysis prompts, Markdown document prompts, frontmatter stripping, inapplicable Markdown pairs, synthesis prompts, manifest contents, missing template failures, and deterministic ordering.
- [ ] Command tests cover prompt preview success, prompt preview failure, and no-runtime behavior.
- [ ] `go test ./...` passes in `/home/antonioborgerees/coding/ultraplan-go`.
- [ ] `go build ./cmd/ultraplan` passes in `/home/antonioborgerees/coding/ultraplan-go`.

## Non-Goals

- Agentwrap/OpenCode runtime integration is not included.
- Actual execution of `ultraplan study <study> run <dimension-ref> <source-ref>` or `ultraplan study <study> synthesize <dimension-ref>` is not included unless limited to already-existing help/placeholder behavior.
- Runtime health checks, permission policy mapping, event/log mapping, retry, fallback, repair, cancellation, and provider cost metadata are not included.
- `run-all`, `run-loop`, durable task state, status calculations, worker pools, and resumability are not included.
- Report validation changes beyond consuming Sprint 6 validation/path concepts are not included.
- Summary generation and `summary.csv` writing are not included.
- Code reference extraction is not included.
- Assisted study initialization, source cloning changes, and Markdown source discovery changes are not included.
- Target workflows, sprint planning, and sprint execution remain out of scope.

## Constraints

- Prompt composition behavior must live in `internal/study`; do not create global `internal/prompts`, `internal/templates`, `internal/reports`, or `internal/validation` packages for study-owned behavior.
- CLI parsing and preview rendering must stay in `internal/app`; prompt assembly, validation of prompt inputs, and manifest construction must stay outside command handlers.
- Preserve the dependency direction from `docs/ARCHITECTURE.md`: `study` may use `workspace` and platform helpers, but platform packages must not import `study`.
- Apply the roadmap Skeleton / Local CLI Gate only; Runtime / Provider and Batch / Durable Workflow gates do not apply because runtime execution and durable orchestration are non-goals.
- Prompt builders must not invoke subprocesses, OpenCode, agentwrap runtimes, provider APIs, network calls, or repository-wide scans.
- Markdown document source prompts must use the existing Sprint 5 frontmatter stripping and applicability behavior.
- Directory prompt source isolation rules must not be weakened by shared base prompt content.
- Synthesis prompt construction must use deterministic ordering for source reports and must not invent missing reports.
- Prompt manifests must avoid absolute paths in user-facing output unless needed for local diagnostics; prefer workspace-relative paths where available.
- Tests must use local fixtures or temporary directories and must not require OpenCode, agentwrap runtime execution, provider credentials, network access, or real source repositories.
- Errors must preserve cause chains for filesystem and parsing failures and must identify the failing template, dimension, source, report, or output path.

## Dependencies

| Prior Sprint / Output | Required For | Notes |
| --- | --- | --- |
| Sprint 1: Go Module, CLI Shell, and App Composition | Buildable CLI, command wiring, and package structure | Required for preview command wiring, tests, and `cmd/ultraplan` build continuity. |
| Sprint 2: Workspace, Config, Logging, and Health Skeleton | Workspace discovery, template path resolution, and path safety | Required so prompt builders locate workspace-managed `prompts/`, `templates/`, studies, sources, and reports. |
| Sprint 3: Study Domain, Listing, and Resolution | Study, source, dimension, and lookup model | Required for resolving selected studies, dimensions, and sources by supported references. |
| Sprint 4: Study Initialization From YAML | Generated study directory layout | Required for expected `dimensions/`, `sources/`, `reports/source/`, and `reports/final/` structures used by prompt inputs and fixtures. |
| Sprint 5: Markdown Document Sources and Applicability | Markdown source kind, frontmatter stripping, and applicability filtering | Required for Markdown document prompts and inapplicable pair handling. |
| Sprint 6: Report Validation and Rating Parsing | Deterministic report paths and validation concepts | Required for synthesis prompt manifests and expected output/report path helpers. |
| `projects/ultraplan-go/docs/PRD.md` | Product prompt composition behavior | Governs analysis prompts, synthesis prompts, dry-run preview, Markdown document source isolation, and deferred target/sprint scope. |
| `projects/ultraplan-go/docs/TRD.md` | Technical prompt builder behavior | Governs prompt inputs, manifest requirements, missing-file failure behavior, and tests. |
| `projects/ultraplan-go/docs/ARCHITECTURE.md` | Package ownership and dependency direction | Governs keeping prompt behavior in `internal/study` and runtime-generic behavior out of product modules. |
| `projects/ultraplan-go/sprints/05-study-listing-inspection/review.md` | Carry-forward decisions | Confirms Markdown source discovery, frontmatter stripping, and applicability behavior are complete and runtime/prompt/report execution remained deferred. |
| `projects/ultraplan-go/sprints/06-study-run-execution/review.md` | Carry-forward decisions | Confirms report validation/path/rating concepts are complete and prompt composition/runtime execution remained deferred. |

## Review Expectations

| What | How Verified |
| --- | --- |
| Scope alignment | Review `requirements.md`, `sprint-index.md`, `technical-handbook.md`, `reasoning.md`, and `plan.md` against the roadmap prompt-composition milestone and confirm no runtime execution, run-loop, summary, code extraction, target, or sprint-execution scope was added. |
| Required outputs exist | Check each path listed in Required Outputs and confirm implementation files are present or intentionally replaced by equivalent files recorded in `review.md`. |
| Directory prompt behavior | Run or inspect tests proving directory prompts include base/dimension/template content, source/output metadata, source isolation, and file-line citation rules. |
| Markdown prompt behavior | Run or inspect tests proving Markdown prompts embed stripped document content, exclude frontmatter, forbid external filesystem/code exploration, and handle inapplicable pairs as non-runnable. |
| Synthesis prompt behavior | Run or inspect tests proving synthesis prompts include synthesis instructions, dimension content, final template, deterministic applicable report manifest, and final output path. |
| Manifest behavior | Run or inspect tests proving prompt manifests contain the required fields and use deterministic workspace-relative paths where available. |
| Preview behavior | Run or inspect command tests proving preview prints or writes composed prompts and does not invoke agentwrap, OpenCode, network calls, or provider credentials. |
| Architecture conformance | Review implementation paths and imports to confirm study-owned prompt behavior remains in `internal/study`, CLI wiring remains in `internal/app`, platform packages do not import product modules, and no global prompt/template package was introduced. |
| Verification commands | Run `go test ./...` and `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go` and record the results in `review.md`. |
