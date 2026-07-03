# Sprint Plan: 07 Prompt Composition Generation

> Project: `ultraplan-go`
> Sprint: `07-prompt-composition-generation`
> Source: `projects/ultraplan-go/sprints/07-prompt-composition-generation/reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/sprints/07-prompt-composition-generation/requirements.md`, `projects/ultraplan-go/sprints/07-prompt-composition-generation/reasoning.md`, `projects/ultraplan-go/sprints/07-prompt-composition-generation/sprint-index.md`, `projects/ultraplan-go/sprints/07-prompt-composition-generation/technical-handbook.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/project-index.md`, `templates/sprint-plan.md`; no `projects/ultraplan-go/sprints/07-prompt-composition-generation/reasoning/*.md` files were present.

This plan executes `reasoning.md`. It must not invent architecture, scope, or decisions.

## Reasoning Source

- **Sprint Reasoning:** `projects/ultraplan-go/sprints/07-prompt-composition-generation/reasoning.md`
- **Sprint Index:** `projects/ultraplan-go/sprints/07-prompt-composition-generation/sprint-index.md`
- **Technical Handbook:** `projects/ultraplan-go/sprints/07-prompt-composition-generation/technical-handbook.md`
- **Area Reasoning:** none present; `sprint-index.md` selected `reasoning/prompt-composition.md`, but no `reasoning/*.md` files were found.

## Sprint Status

- **Status:** complete
- **Owner:** implementation agent
- **Start Date:** 2026-05-31
- **Completion Date:** 2026-05-31

## Decisions To Execute

| Decision | Source Section | Execution Implication |
| --- | --- | --- |
| Study owns prompt composition | `reasoning.md#decision-1-study-owns-prompt-composition` | Add prompt request/result/manifest concepts in `internal/study/domain.go`; add deterministic prompt builders in `internal/study/prompts.go`; keep command handlers thin in `internal/app/study_prompt_commands.go`. |
| Deterministic manifest-backed prompt results | `reasoning.md#decision-2-deterministic-manifest-backed-prompt-results` | Builders must return prompt text plus manifest with stable fields and deterministic workspace-relative paths where available. |
| Directory analysis prompts preserve repository isolation and code citations | `reasoning.md#decision-3-directory-analysis-prompts-preserve-repository-isolation-and-code-citations` | Directory prompts must include base prompt, dimension Markdown, `templates/repo-analysis.md`, metadata, expected output path, source-isolation rules, and file-path-line citation rules unless disabled by dimension policy. |
| Markdown document prompts are embedded-document only | `reasoning.md#decision-4-markdown-document-prompts-are-embedded-document-only` | Markdown prompts must use stripped document content, exclude YAML frontmatter, forbid external file/repository/code inspection, and default away from code citation requirements. |
| Synthesis uses applicable report manifests only | `reasoning.md#decision-5-synthesis-uses-applicable-report-manifests-only` | Synthesis prompt composition must include only applicable source reports, sort them deterministically, fail on missing applicable reports, and exclude inapplicable Markdown pairs. |
| Preview command is runtime-free rendering only | `reasoning.md#decision-6-preview-command-is-runtime-free-rendering-only` | Add preview/dry-run CLI wiring that resolves refs, calls study prompt builders, writes stdout or explicit output path, and never invokes runtime, agentwrap, OpenCode, subprocesses, providers, or network calls. |
| Tests focus on determinism, source-kind semantics, errors, and no-runtime boundaries | `reasoning.md#decision-7-tests-focus-on-determinism-source-kind-semantics-errors-and-no-runtime-boundaries` | Add unit tests for prompt builders/manifests and command tests for preview success/failure/no-runtime behavior; verify with `go test ./...` and `go build ./cmd/ultraplan`. |

## Requirements / Contracts To Satisfy

| Contract / Requirement ID | Required Behavior | Evidence Planned |
| --- | --- | --- |
| Architecture | Prompt behavior stays in `internal/study`; CLI rendering stays in `internal/app`; no global prompt/template/report package; platform packages do not import `study`. | Import/package review, unit tests, command tests, `go test ./...`, `go build ./cmd/ultraplan`. |
| Errors | Missing templates, documents, reports, unknown refs, ambiguous refs, and inapplicable single-source previews produce actionable diagnostics that preserve causes. | Unit and command tests assert diagnostic content and non-zero failures. |
| Security | Preview does not invoke runtime, subprocesses, OpenCode, agentwrap runtime execution, providers, network calls, or repo-wide scans; artifacts avoid unnecessary absolute paths. | Command tests and code review of preview path. |
| Testing | Deterministic local tests cover builders, manifests, source-kind behavior, errors, and command rendering. | `internal/study/prompts_test.go`, `internal/app/study_prompt_commands_test.go`, `go test ./...`. |
| Documentation | Help and generated sprint artifacts accurately describe prompt-only preview behavior. | Command help review and sprint review. |
| CLI Surface | Preview supports stdout or explicit output path and script-friendly failures. | Command tests for stdout, output file, usage errors, and failure exit behavior. |
| AC-01, AC-02 | Required prompt/template files are loaded from deterministic workspace paths and missing files fail before runtime-facing behavior. | Missing-template unit tests naming `prompts/base.md`, `prompts/synthesize.md`, `templates/repo-analysis.md`, and final report template path. |
| AC-03, AC-04 | Directory analysis prompts include required content, source isolation, and citation rules. | Directory prompt tests with required substring/manifest assertions. |
| AC-05, AC-06, AC-07, AC-08 | Markdown prompts embed stripped content, exclude frontmatter, forbid external exploration, and default away from code citation requirements. | Markdown prompt and frontmatter stripping tests. |
| AC-09 | Inapplicable Markdown source/dimension pairs do not produce runnable analysis prompts. | Applicability tests for single-source preview behavior. |
| AC-10, AC-11 | Synthesis prompts include synthesis inputs and deterministic applicable report manifest, excluding inapplicable Markdown pairs. | Synthesis prompt tests, manifest-ordering tests, missing applicable report tests. |
| AC-12, AC-13 | Builders return deterministic prompt text plus manifests with required fields. | Manifest exact-field tests and repeated-build determinism tests. |
| AC-14, AC-15 | Preview prints or writes prompt output without runtime and returns actionable errors for all required failure cases. | Command success/failure/no-runtime tests. |
| AC-16, AC-17, AC-18, AC-19 | Unit/command tests exist, `go test ./...` passes, and `go build ./cmd/ultraplan` passes. | Verification command results recorded in sprint review. |

## Tasks

- [x] **Task 1: Inspect Existing Study And Workspace APIs**
  > Executes: `Decision 1`, `Decision 2`, Architecture, AC-01 through AC-13
  - [x] Inspect `/home/antonioborgerees/coding/ultraplan-go/internal/study/domain.go` for existing `Study`, `Source`, `Dimension`, source-kind, applicability, frontmatter, and report path concepts.
  - [x] Inspect existing study stores/services/resolution helpers to reuse study, dimension, and source lookup behavior instead of duplicating it in command handlers.
  - [x] Inspect workspace path helpers for resolving `prompts/base.md`, `prompts/synthesize.md`, `templates/repo-analysis.md`, `templates/report.md`, study dimensions, sources, per-source reports, and final reports.
  - [x] Confirm the actual final report template path; default to `templates/report.md` only if the current workspace and TRD structure support that path.
  - [x] Record any missing prerequisite API as an implementation blocker or minimal local helper, not as a new global package.

- [x] **Task 2: Define Prompt Domain Types**
  > Executes: `Decision 1`, `Decision 2`, AC-12, AC-13
  - [x] Add or extend prompt kind values for directory analysis, Markdown document analysis, and synthesis in `internal/study/domain.go`.
  - [x] Add prompt request types or one minimal request shape that can express study, dimension, optional source, source kind, and expected output path without runtime settings.
  - [x] Add manifest types for prompt kind, study, dimension, source/source kind when applicable, template references, input document path, input report paths, source report manifest entries, and expected output path.
  - [x] Add prompt result or preview result type containing prompt text and manifest.
  - [x] Keep manifest fields requirement-driven and avoid durable runtime-state or public JSON versioning unless already present.

- [x] **Task 3: Implement Deterministic Template Loading And Path Rendering**
  > Executes: `Decision 2`, AC-01, AC-02, AC-12, AC-13
  - [x] Implement template loading in `internal/study/prompts.go` using deterministic workspace paths for `prompts/base.md`, `prompts/synthesize.md`, `templates/repo-analysis.md`, and the chosen final report template.
  - [x] Return wrapped errors for missing or unreadable required inputs and name the failing path in the diagnostic.
  - [x] Render manifest paths as workspace-relative where helpers allow it; use absolute paths only in local diagnostics when required.
  - [x] Ensure all slices in manifests are sorted deterministically and no maps are emitted in nondeterministic order.
  - [x] Avoid any subprocess, OpenCode, agentwrap runtime, provider, network, or repository-wide scan call.

- [x] **Task 4: Implement Directory Analysis Prompt Builder**
  > Executes: `Decision 3`, AC-03, AC-04, AC-12, AC-13
  - [x] Compose directory prompts from base prompt, selected dimension Markdown, `templates/repo-analysis.md`, study name, source name, source path, expected output path, and directory-specific source-isolation rules.
  - [x] Add explicit instruction to explore only the selected source directory and not unrelated workspace or sibling source content.
  - [x] Add file-path-and-line citation requirements for code claims unless the selected dimension disables code citation requirements.
  - [x] If citation-disable behavior depends on dimension text, centralize the detection rule in one narrow helper and test supported wording.
  - [x] Populate analysis manifest fields for prompt kind, study, dimension, source, source kind, template paths, and expected output path.

- [x] **Task 5: Implement Markdown Document Analysis Prompt Builder**
  > Executes: `Decision 4`, AC-05, AC-06, AC-07, AC-08, AC-09, AC-12, AC-13
  - [x] Reuse existing Sprint 5 frontmatter stripping and applicability behavior for Markdown document sources.
  - [x] Fail with an actionable missing-document error when the Markdown file does not exist or cannot be read.
  - [x] Compose Markdown prompts from base prompt, selected dimension Markdown, report template, study/source/output metadata, and stripped document body.
  - [x] Exclude YAML frontmatter from the embedded document content.
  - [x] Add explicit instructions that all source material is embedded, the runtime must not inspect external files/repositories/code, and code citation requirements do not apply by default unless the dimension states otherwise.
  - [x] Return a clear non-runnable error or skipped result for inapplicable Markdown source/dimension pairs without producing a runnable prompt.

- [x] **Task 6: Implement Synthesis Prompt Builder**
  > Executes: `Decision 5`, AC-10, AC-11, AC-12, AC-13
  - [x] Compose synthesis prompts from `prompts/synthesize.md`, selected dimension Markdown, final report template, deterministic applicable per-source report manifest, and expected final output path.
  - [x] Use applicability filtering so directory sources are included and Markdown sources are included only when applicable to the selected dimension.
  - [x] Sort applicable report inputs deterministically by the existing source ordering or a documented stable source key.
  - [x] Fail on each missing applicable per-source report and name the source/report path in the diagnostic.
  - [x] Exclude inapplicable Markdown source/dimension pairs from the manifest and from missing-report failures.

- [x] **Task 7: Add Prompt Builder Unit Tests**
  > Executes: `Decision 7`, Testing, AC-01 through AC-13, AC-16
  - [x] Add `internal/study/prompts_test.go` with local fixtures or temporary workspaces only.
  - [x] Test directory prompt content, manifest fields, source-isolation rules, citation rules, and disabled-citation behavior where supported.
  - [x] Test Markdown prompt stripped body inclusion, frontmatter exclusion, document-only rules, no-code-citation default text, missing document failures, and inapplicable pair behavior.
  - [x] Test synthesis prompt content, deterministic applicable report ordering, inapplicable Markdown exclusion, missing applicable report failures, and final output path.
  - [x] Test missing template failures for `prompts/base.md`, `prompts/synthesize.md`, `templates/repo-analysis.md`, and final report template input.
  - [x] Test repeated builder calls produce identical prompt text and manifest output for identical filesystem inputs.

- [x] **Task 8: Add Runtime-Free Preview CLI Wiring**
  > Executes: `Decision 6`, AC-14, AC-15, Architecture, CLI Surface, Security
  - [x] Add `internal/app/study_prompt_commands.go` or extend the existing app command tree with a minimal prompt preview or dry-run command surface.
  - [x] Keep command handlers limited to argument parsing, workspace/study/dimension/source resolution, calling study prompt builders, and rendering output.
  - [x] Support stdout output and explicit output path writing.
  - [x] Decide manifest exposure using the minimal existing output conventions; if no convention exists, include manifest in a deterministic preview format without inventing a broad public JSON contract.
  - [x] Return non-zero actionable errors for unknown studies, dimensions, sources, ambiguous references, inapplicable Markdown pairs, missing templates, missing Markdown files, and missing synthesis reports.
  - [x] Do not call run/synthesize runtime paths, agentwrap, OpenCode, providers, network APIs, or subprocess execution.

- [x] **Task 9: Add Preview Command Tests**
  > Executes: `Decision 6`, `Decision 7`, AC-14, AC-15, AC-17
  - [x] Add `internal/app/study_prompt_commands_test.go` with fixture workspaces and capturable output streams.
  - [x] Test successful analysis preview to stdout.
  - [x] Test successful preview writing to an explicit output path.
  - [x] Test synthesis preview success with applicable report manifest.
  - [x] Test failure cases for unknown or ambiguous refs, inapplicable Markdown pairs, missing templates, missing Markdown files, and missing synthesis reports.
  - [x] Test no-runtime behavior by ensuring preview wiring uses prompt builders directly and cannot invoke runtime dependencies in the tested command path.

- [x] **Task 10: Verify Scope, Architecture, And Build**
  > Executes: Architecture Review, Sprint Review, AC-18, AC-19
  - [x] Review changed imports and files to confirm prompt behavior remains in `internal/study`, CLI rendering remains in `internal/app`, platform packages do not import `study`, and no global `internal/prompts`, `internal/templates`, `internal/reports`, or `internal/validation` package was introduced.
  - [x] Review diff for non-goal scope: no runtime execution, run-loop, summary generation, code extraction, target workflow, sprint workflow, provider config, network, subprocess, or durable workflow additions.
  - [x] Run `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go`.
  - [x] Run `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`.
  - [x] Record verification results and any carry-forward issues in `review.md` during sprint review.

## Evidence Checklist

- [x] Tests prove directory prompt content, source isolation, citation rules, manifest fields, and deterministic output.
- [x] Tests prove Markdown prompt frontmatter stripping, embedded-document-only behavior, inapplicable pair handling, missing document diagnostics, and no default code citation requirement.
- [x] Tests prove synthesis prompt content, applicable report manifest ordering, missing applicable report failures, and inapplicable Markdown exclusion.
- [x] Tests prove missing required template failures before runtime-facing behavior.
- [x] Command tests prove preview stdout/output-path success and actionable non-zero failures.
- [x] Command tests or code review prove no runtime, agentwrap, OpenCode, subprocess, provider, network, or repository-wide scan invocation.
- [x] Documentation/help text is accurate for prompt-only preview and does not imply execution.
- [x] Architecture review evidence confirms package ownership and dependency direction.
- [x] Sprint review evidence confirms no non-goal scope was added.
- [x] Deviations from `reasoning.md` are recorded before implementation continues.

## Verification Commands

| Check | Command | Expected Result |
| --- | --- | --- |
| Unit and command tests | `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go` | Passes without OpenCode, agentwrap runtime execution, provider credentials, network access, or real source repositories. |
| CLI build | `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go` | Builds the CLI binary successfully. |
| Architecture review | Inspect implementation diff and imports | Prompt behavior is study-owned, CLI preview is app-owned, platform packages do not import study, and no forbidden global packages are introduced. |
| Scope review | Inspect implementation diff and command behavior | No runtime execution, run-loop, summary generation, code extraction, target workflow, sprint workflow, provider config, network, or subprocess scope is added. |

## Risks And Blockers

| Risk / Blocker | Source | Mitigation | Status |
| --- | --- | --- | --- |
| Existing prior-sprint APIs for study resolution, applicability, frontmatter stripping, or report paths may be incomplete. | `reasoning.md#assumptions-and-risks` | Inspect APIs first; use minimal local helpers only where needed; record prerequisite gaps in review if they block acceptance. | closed: existing APIs were sufficient. |
| Final report template path may differ from the assumed `templates/report.md`. | `reasoning.md#assumptions-and-risks`, TRD workspace structure | Confirm current workspace/template naming before implementing synthesis builder and tests. | closed: workspace scaffold and validation require `templates/report.md`. |
| Citation-disable detection may be brittle if dimensions lack structured metadata. | `reasoning.md#assumptions-and-risks` | Centralize the rule in one helper and test supported wording; leave structured policy as future schema work. | mitigated: centralized helper supports explicit disabled wording and is tested. |
| Preview command shape may need adjustment in future runtime dry-run sprint. | `reasoning.md#assumptions-and-risks` | Keep preview minimal and backed by reusable study prompt builders; avoid coupling to future runtime execution. | carried forward: preview is intentionally minimal. |
| Prompt text tests may become noisy when templates change. | `reasoning.md#assumptions-and-risks` | Assert exact manifests and ordering; use focused required-text checks unless full golden coverage is intentionally stable. | mitigated: tests use focused required-content checks and exact manifest/order assertions. |
| Manifest paths may leak absolute local details. | `reasoning.md#assumptions-and-risks` | Prefer workspace-relative paths in manifest output; reserve absolute paths for diagnostics only. | closed: manifest paths are workspace-relative. |
| Strict synthesis failures may block partial prompt exploration. | `reasoning.md#assumptions-and-risks` | Make errors identify missing applicable reports and explain synthesis input requirements; defer partial mode unless explicitly scoped later. | carried forward: strict missing-report failure remains by design. |

## Review Inputs

Review should use:

- `projects/ultraplan-go/sprints/07-prompt-composition-generation/requirements.md`
- `projects/ultraplan-go/sprints/07-prompt-composition-generation/sprint-index.md`
- `projects/ultraplan-go/sprints/07-prompt-composition-generation/technical-handbook.md`
- `projects/ultraplan-go/sprints/07-prompt-composition-generation/reasoning.md`
- `projects/ultraplan-go/sprints/07-prompt-composition-generation/plan.md`
- implementation diff in `/home/antonioborgerees/coding/ultraplan-go`
- verification evidence from `go test ./...` and `go build ./cmd/ultraplan`
- `system/protocols/architecture-review-protocol.md`
- `system/protocols/sprint-review-protocol.md`

## Execution Log

| Date / Step | Action | Evidence / Notes |
| --- | --- | --- |
| 2026-05-31 / Planning | Created sprint implementation plan from `reasoning.md`, requirements, project docs, sprint index, and technical handbook. | No implementation code changed. Area reasoning files were absent and recorded as absent. |
| 2026-05-31 / Implementation | Added study-owned prompt manifest/request/result types, deterministic directory/Markdown/synthesis builders, and runtime-free prompt preview command wiring. | Changed `internal/study/domain.go`, `internal/study/prompts.go`, and `internal/app/study_commands.go`; added unit and command tests. |
| 2026-05-31 / Verification | Ran required verification commands from `/home/antonioborgerees/coding/ultraplan-go`. | `go test ./...` passed. `go build ./cmd/ultraplan` passed. The build artifact was removed after verification. |
| 2026-05-31 / Review | Reviewed changed imports and scope. | Prompt behavior remains in `internal/study`; CLI rendering remains in `internal/app`; no runtime, network, provider, subprocess, run-loop, summary, target, or sprint workflow scope was added. |

## Completion Criteria

- [x] All tasks are complete or explicitly deferred in `review.md` with requirement impact stated.
- [x] `go test ./...` was run from `/home/antonioborgerees/coding/ultraplan-go`, or a blocker is documented.
- [x] `go build ./cmd/ultraplan` was run from `/home/antonioborgerees/coding/ultraplan-go`, or a blocker is documented.
- [x] Evidence satisfies the expectations from `reasoning.md`.
- [x] Architecture and sprint review protocols can evaluate conformance without guessing intent.
- [x] `review.md` records final acceptance status, verification output, deviations, and carry-forward decisions.
