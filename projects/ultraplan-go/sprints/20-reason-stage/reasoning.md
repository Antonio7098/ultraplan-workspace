# Sprint Reasoning: Reason Stage

> Project: `ultraplan-go`
> Sprint: `20-reason-stage`
> Output: `projects/ultraplan-go/sprints/20-reason-stage/reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/20-reason-stage/requirements.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/sprints/20-reason-stage/sprint-index.md`, `projects/ultraplan-go/sprints/20-reason-stage/technical-handbook.md`, `templates/sprint-reasoning.md`

This document decides. It synthesizes selected context, handbook evidence, selected contracts, and sprint requirements into final sprint decisions for implementing the reason stage.

It does not replace `sprint-index.md`, `technical-handbook.md`, or any `reasoning/*.md` area artifact.

## Sprint Purpose

- **Goal:** Implement the planning reason stage so `area-reasoning` and final `reasoning` can be validated, previewed, flowed, generated through the generic runtime boundary, and recorded in `flow-state.json` from valid `requirements.md`, `sprint-index.md`, and `technical-handbook.md`.
- **Non-Goals:** Do not generate or validate `plan.md`; do not create task breakdowns or implementation checklists in reasoning artifacts; do not implement sprint execution, smoke investigation, review automation, issues, Git mutation, `.run-state.json`, `smoke.md`, `smoke.json`, generated `review.md`, `issues.md`, or `issues.json`; do not reselect context or use unselected contracts, reports, templates, protocols, docs, or prior reviews as authoritative decision sources.
- **Depends On:** Valid Sprint 17 flow-state foundations, Sprint 18 selected-context behavior, Sprint 19 technical-handbook behavior, valid project catalog data in `project-index.md`, Phase 2 planning scope in PRD/TRD/Architecture docs, and the selected evidence distilled in `technical-handbook.md`.

## Selected Context And Pre-Reasoning Artifacts

| Artifact | Path | How It Was Used |
| --- | --- | --- |
| Project Index | `projects/ultraplan-go/project-index.md` | Confirmed project scope, target implementation directory, Phase 2 planning workflow, selected catalog pool, available evidence reports, available Architecture reasoning template, and review protocols. |
| Requirements | `projects/ultraplan-go/sprints/20-reason-stage/requirements.md` | Provided the sprint contract, required outputs, acceptance criteria, non-goals, constraints, dependencies, and review expectations. |
| PRD | `projects/ultraplan-go/docs/PRD.md` | Confirmed product behavior: planning runs through `plan.md`, generated artifacts remain editable, runtime success is not product success, and execution/smoke/review/issues/Git mutation are deferred. |
| TRD | `projects/ultraplan-go/docs/TRD.md` | Confirmed `internal/sprint` ownership, planning stage model, validators, prompt rendering, generic runtime boundary, strict `flow-state.json`, exit codes, diagnostics, and testing expectations. |
| Architecture Doc | `projects/ultraplan-go/docs/ARCHITECTURE.md` | Confirmed module-driven ownership, allowed dependencies, product/platform separation, reuse boundaries, and prohibition on extracting study semantics into planning. |
| Sprint Index | `projects/ultraplan-go/sprints/20-reason-stage/sprint-index.md` | Selected contracts, evidence reports, Architecture reasoning template, required protocols, excluded context, and non-goal boundaries for this sprint. |
| Technical Handbook | `projects/ultraplan-go/sprints/20-reason-stage/technical-handbook.md` | Provided distilled study evidence, concrete repository/source references, relevant patterns, trade-offs, anti-patterns, design pressures, and open questions for final decisions. |
| Sprint Reasoning Template | `templates/sprint-reasoning.md` | Provided the structure for this final decision artifact. |

## Area-Specific Reasoning Inputs

The sprint index selects the Architecture reasoning template with expected output `projects/ultraplan-go/sprints/20-reason-stage/reasoning/architecture.md`. No area-specific reasoning documents were present under `projects/ultraplan-go/sprints/20-reason-stage/reasoning/*.md` when this document was created.

| Area | Reasoning Document | Key Conclusion | Evidence Basis | Impact On Final Decision |
| --- | --- | --- | --- | --- |
| Architecture | `projects/ultraplan-go/sprints/20-reason-stage/reasoning/architecture.md` | No area artifact was available as an input. Architecture decisions in this final reasoning are therefore made directly from selected requirements, project docs, sprint index, and handbook evidence. | No area-specific artifact evidence. The selected handbook and project architecture docs provide the available architecture grounding. | The implementation must still support Architecture area reasoning as a selected template. Product validation and flow must require the area artifact before final reasoning completion when the selected template is present. |

## Sprint Technical Handbook Summary

- **Relevant Patterns:** Thin CLI entrypoints delegate into owned services; product behavior stays in `internal/` modules; manual composition roots and explicit seams are preferred over framework DI; stdout/stderr, filesystem, and runtime boundaries should be injectable for tests; state completion follows successful generation plus validation; workspace paths and selected artifacts are trust boundaries.
- **Important Trade-Offs:** Keeping reason-stage behavior inside `internal/sprint` avoids premature generic packages but may duplicate small Markdown/validation mechanics; strict validators prevent false success but must avoid brittle prose matching; fake-runtime and command-level tests add setup but keep default verification offline; deterministic prompt previews constrain output richness but improve reviewability.
- **Warnings / Anti-Patterns:** Do not let command handlers own reasoning business rules; do not introduce global mutable runtime/config/prompt/validation state; do not directly invoke shell, Git, OpenCode, provider APIs, or agent adapters from sprint logic; do not treat runtime success as artifact success; do not allow unselected context to leak into generated reasoning decisions.
- **Evidence Confidence:** High for package boundaries, command delegation, error classification, IO seams, security, and testing because multiple selected reports cite concrete mature Go CLI implementations. Medium for terminal UX, state/context, and philosophy where the evidence is more pattern-oriented but still consistent with the sprint scope.

## Contracts Applied

| Contract / Requirement ID | Constraint | Decision Impact | Expected Evidence |
| --- | --- | --- | --- |
| Architecture contract; `requirements.md` lines 57-60, 92-99, 102 | Keep reason-stage behavior in `internal/sprint`; keep CLI adapters thin; keep platform product-agnostic; do not import study semantics. | Decisions 1, 2, 6, and 7 keep ownership in `internal/sprint`, use `internal/app` only for wiring, and reject global planning packages. | Import review, package review, `go test ./...`, command tests proving handlers delegate. |
| Errors contract; `requirements.md` lines 61-65, 127-135 | Actionable diagnostics, cause preservation, deterministic ordering, and exit-code mapping are required. | Decisions 3, 5, and 7 require structured validation findings and stage-specific exit codes. | Unit tests for validation errors, unsupported stages, runtime failures, stdout/stderr separation, and exit codes `0`, `2`, `5`, runtime failure code. |
| Configuration contract; `requirements.md` lines 42-45, 100-101 | Prompt/runtime diagnostics must avoid secrets and unnecessary absolute paths while respecting runtime-facing config. | Decisions 4, 5, and 7 require workspace-relative prompt content by default and redaction of sensitive values. | Prompt preview tests, command output tests, review of rendered prompt fixtures. |
| Observability contract; `requirements.md` lines 53-55, 61, 131-135 | Failures must be visible, safe, and recorded in flow state. | Decisions 5 and 6 require failed stage recording, safe messages, and no later-stage completion. | Flow tests for runtime failure, validation failure, state write failure, and safe diagnostics. |
| Security contract; `requirements.md` lines 40, 45-47, 56-58, 94-101 | Validate workspace-safe paths, reject selected path escapes, prohibit direct shell/Git/OpenCode/provider calls from sprint code, redact secrets. | Decisions 1, 2, 4, 5, and 7 enforce selected-context loading and generic runtime-only execution. | Path escape tests, import review, fake runtime tests, prompt no-mutation rule tests. |
| Testing contract; `requirements.md` lines 66-75, 126-136 | Deterministic offline tests must cover selected-template detection, validators, prompts, flow, commands, race, and build. | All decisions require explicit unit, fixture, command, fake-runtime, race, and build evidence. | `go test ./...`, `go test -race ./...`, `go build ./cmd/ultraplan`, targeted tests in `internal/sprint` and `internal/app`. |
| Documentation contract; `requirements.md` lines 35-36, 42-45, 125 | Help and prompt previews must be maintainable and user-facing artifacts remain editable Markdown. | Decisions 3, 4, and 7 require help output and deterministic editable Markdown prompt targets. | Help output tests and review of generated prompt content. |
| CLI Surface contract; `requirements.md` lines 31-38, 60-65, 71 | `validate`, `prompt`, and `flow` command behavior must be scriptable, no ANSI by default, stdout/stderr separated, and correctly exit-coded. | Decision 7 keeps CLI wiring thin and maps service results to command output/exit codes. | `internal/app/sprint_commands_test.go`, help tests, no-runtime validate/prompt tests. |
| LLM Runtime contract; `requirements.md` lines 38, 40-41, 65, 71 | Non-dry-run flow uses only the generic platform runtime boundary; runtime success alone is insufficient. | Decision 5 requires fake-runtime seam and post-generation artifact validation before state completion. | Fake-runtime flow tests, import review, artifact-missing and artifact-invalid tests. |
| LLM Evaluation / Cost / Safety contract; `requirements.md` lines 41, 45, 100-101 | Safety boundaries and validation-after-runtime success are required; unsafe prompt/runtime diagnostics are excluded. | Decisions 4 and 5 include no-mutation prompt rules and safe diagnostics. | Prompt rendering tests and runtime failure diagnostic tests. |
| Workflows contract; `requirements.md` lines 37, 39, 51-56 | Flow supports Planning Phase 2 stages through `reasoning`; area reasoning is skipped only when no templates are selected; unsupported future stages are rejected. | Decision 6 defines stage ordering and flow-state representation. | Flow tests for dry-run, selected Architecture area requirement, no-template skip, unsupported stage rejection, no plan completion. |
| Persistence And Migrations contract; `requirements.md` lines 41, 54-56, 64 | `flow-state.json` remains strict, versioned, atomic, and preserves last valid state on write failure. | Decision 6 requires atomic writes and strict stage set. | State load/write tests, write failure tests, malformed/unsupported state tests. |

## Repos Studied / Source Evidence Used

| Source / Repo / Report | Concrete Reference | Relevant Finding | Why It Matters For This Sprint | Used In Decision(s) |
| --- | --- | --- | --- | --- |
| `01-project-structure` | `technical-handbook.md` lines 16, 31-33; chezmoi `main.go:26-34`; gh-cli `cmd/gh/main.go:6`; k9s `cmd/root.go:112-127`; restic `cmd/restic/main.go:37-114` | Mature Go CLIs keep entrypoints thin and route work inward. | Supports keeping business rules in `internal/sprint` and command wiring in `internal/app`. | Decisions 1 and 7 |
| `01-project-structure` | `technical-handbook.md` lines 16, 33; chezmoi `internal/chezmoi/chezmoi.go:1-2`; restic `internal/restic/repository.go:18` | `internal/` boundaries protect application-owned behavior and import direction. | Supports sprint-owned validators, prompts, flow-state behavior, and no platform reverse imports. | Decisions 1, 2, and 6 |
| `02-command-architecture` | `technical-handbook.md` lines 17, 31, 60; gh-cli `pkg/cmdutil/factory.go:16-43`; helm `pkg/cmd/install.go:132-145`; opencode `cmd/root.go:49-183` warning | Command factories delegate; large handlers are a maintainability smell. | Supports thin `sprint_commands.go` wiring and service methods for validate/prompt/flow. | Decision 7 |
| `03-dependency-injection` | `technical-handbook.md` lines 18, 35, 53; gh-cli `pkg/cmd/factory/default.go:26-46`; go-task `executor.go:22-24`; helm `pkg/action/action.go:118-146`; opencode `internal/app/app.go:42-81` | Manual composition roots, constructor injection, functional options, and small interfaces beat DI frameworks. | Supports explicit store/runtime/clock/IO seams without broad abstractions. | Decisions 2, 5, and 7 |
| `04-configuration-management` | `technical-handbook.md` lines 19, 113; chezmoi `internal/cmd/config.go:2253-2287`; go-task `internal/flags/flags.go:314-327`; opencode `internal/config/config.go:609-641` | Config precedence and validation must be explicit, with safe rendering. | Supports redacted prompt/runtime summaries and no secret leakage in previews. | Decision 4 |
| `05-error-handling` | `technical-handbook.md` lines 20, 39, 66; rclone `fs/fserrors/error.go:22-29`; go-task `errors/errors.go:47-50`; restic `internal/errors/fatal.go:10`; gh-cli `internal/ghcmd/cmd.go:44-49` | Typed/wrapped errors and classification support actionable user diagnostics and exit codes. | Supports validation result shapes, unsupported-stage errors, and runtime failure mapping. | Decisions 3, 5, and 7 |
| `06-io-abstraction` | `technical-handbook.md` lines 21, 37, 64; gh-cli `pkg/iostreams/iostreams.go:551-568`; go-task `executor.go:541-577`; restic `internal/ui/terminal.go:10-36` | Testable CLIs inject output and filesystem boundaries. | Supports command tests for stdout/stderr, prompt previews, and no-ANSI behavior. | Decisions 2 and 7 |
| `07-state-context` | `technical-handbook.md` lines 22, 41, 56; gh-cli `internal/ghcmd/cmd.go:142`; go-task `task.go:89`; helm `pkg/cmd/install.go:333-347` | Long/runtime-backed work should carry `context.Context`; deterministic parsing need not overuse context. | Supports separating runtime-backed flow from no-runtime validate/prompt paths. | Decision 5 |
| `09-terminal-ux` | `technical-handbook.md` lines 23, 55, 87; chezmoi `internal/cmd/prompt.go:20-256`; gh-cli `internal/prompter/` | Scriptable non-interactive output and restrained progress are preferred. | Supports deterministic prompt previews, calm diagnostics, and no interactive prompts. | Decisions 4 and 7 |
| `10-logging-observability` | `technical-handbook.md` lines 24, 43, 88; helm `internal/logging/logging.go:31-66`; k9s `internal/slogs/keys.go:6-231`; gh-cli `pkg/iostreams/iostreams.go:52-54` | Diagnostics should be structured and routed separately from user/data output. | Supports stderr/stdout separation and safe flow failure details. | Decisions 5, 6, and 7 |
| `11-testing-strategy` | `technical-handbook.md` lines 25, 39, 72, 91; chezmoi `internal/cmd/main_test.go:64-174`; helm `internal/test/test.go:43`; go-task `task_test.go:166-169` | Fixture, fake, command-level, and golden tests provide high confidence without brittle prose matching. | Supports validator, prompt, flow, and command tests with required variables rather than exact generated prose. | All decisions |
| `13-security` | `technical-handbook.md` lines 26, 45, 70, 74, 89; opencode `internal/permission/permission.go:44-108`; opencode `internal/llm/tools/bash.go:41-55`; restic `internal/options/secret_string.go:15-20`; helm `pkg/registry/transport.go:37-41`; k9s `internal/config/json/validator.go:146` | Trust boundaries, path safety, no stringly shell execution, and redaction are essential. | Supports workspace-safe selected paths, no direct shell/Git/OpenCode/provider calls, and redacted prompts/diagnostics. | Decisions 1, 2, 4, 5, and 7 |
| `15-philosophy` | `technical-handbook.md` lines 27, 51, 57; gh-cli `pkg/cmd/factory/default.go:26-46`; restic `internal/backend/backend.go:19-90`; opencode `internal/pubsub/broker.go:10-19` | Coherent projects accept targeted complexity and reject scope outside the core workflow. | Supports deferring generic workflow engines, plugins, plan execution, and real-runtime smoke by default. | Decisions 1, 6, and 7 |

## Trade-Off And Debt Analysis

### Accepted Trade-Offs

| Trade-Off | Benefit | Cost / Constraint Accepted | Why Acceptable Now | Revisit Trigger |
| --- | --- | --- | --- | --- |
| Keep all reason-stage product rules in `internal/sprint` instead of adding global planning packages | Preserves module ownership and prevents study/planning semantic coupling. | Some Markdown validation and prompt helper mechanics may be duplicated inside sprint. | Sprint scope is narrow and selected evidence warns against premature global abstractions. | A later sprint demonstrates the same mechanical helper is needed by multiple modules without product semantics. |
| Validate reasoning artifacts semantically but not by exact prose | Catches missing decisions, evidence, assumptions, risks, placeholders, and unselected context while avoiding brittle tests. | Validators must use clear structural heuristics, not subjective quality judgments. | The sprint requires editable Markdown generated by a runtime; strict exact prose would make flow fragile. | Reopen if invalid artifacts routinely pass or valid artifacts routinely fail due to under-specified checks. |
| Use fake runtime for default flow verification | Keeps tests deterministic, offline, and free of OpenCode/provider/network dependencies. | Real OpenCode behavior is not covered by default sprint verification. | Requirements explicitly exclude real-runtime smoke from default verification and select fake-first evidence. | Reviewer requests deep smoke or a later sprint selects smoke harness evidence. |
| Keep CLI output plain text and no ANSI by default | Makes automation and tests stable; supports stdout/stderr separation. | Runtime-backed flow output is less visually rich. | Sprint acceptance criteria require calm, scriptable, deterministic output. | Future UX sprint selects terminal UX or TUI scope. |
| Require selected Architecture area reasoning for product flow completion even if this manually authored document had no area artifact input | Aligns implementation with selected-template semantics and prevents skipped required reasoning. | The current created `reasoning.md` records an input gap and cannot be used as proof that the selected area artifact existed. | Requirements say selected templates require their area artifacts before final reasoning completes. | If sprint index changes to select no reasoning templates or the Architecture area artifact is created and validated. |
| Reject unsupported future stages at parsing/preflight rather than representing them as skipped states | Maintains strict Phase 2 state model and prevents hidden expansion into execution/smoke/review/issues. | Users get hard errors for familiar prototype stages until future scope is selected. | PRD/TRD and requirements explicitly defer those stages. | A later requirements revision adds those stages to Phase 2 or beyond. |

### Potential Technical Debt

| Debt / Shortcut | Why It Might Accrue | Current Mitigation | Owner / Follow-Up |
| --- | --- | --- | --- |
| Section-based Markdown validators may become inconsistent across planning stages | Sprint 19 and Sprint 20 both validate Markdown artifacts with stage-specific sections. | Keep validators local and small; extract only mechanical helpers after repeated stable needs. | Revisit during Sprint 21 `plan.md` if duplication becomes concrete. |
| Selected-context reference detection may start with conservative string/path checks | Detecting every unselected context citation semantically is hard in editable Markdown. | Enforce selected paths, selected names, placeholder detection, required sections, and obvious unselected catalog references. | Improve when validation false negatives are observed. |
| Flow-state schema evolution pressure | Additional planning stages after `plan` may need state expansion later. | Keep strict schema version, known stage set, and rejection of unsupported stages now. | Future execution/smoke/review/issue stage requirements. |
| Prompt rendering may accumulate long string builders | Required prompt variables are extensive for area and final reasoning prompts. | Keep prompt rendering in `internal/sprint/prompts.go` with table/manifest-oriented helpers and tests. | Extract only if prompt composition becomes hard to review in Sprint 21. |
| Command tests may duplicate setup fixtures across validate, prompt, and flow | Each command path needs similar sprint/project fixture data. | Use focused fixture helpers in tests only; do not add production abstractions for fixtures. | Revisit if command test setup obscures behavior. |

### Future Considerations

| Consideration | Deferred Until | Reason Deferred | What Should Be Preserved Now |
| --- | --- | --- | --- |
| `plan.md` generation and validation | Sprint 21 | Current sprint stops at final reasoning and must not require plan artifacts. | Stage order includes `plan` as supported future Planning Phase 2 stage but does not mark it complete during `--to reasoning`. |
| Real-runtime smoke evidence | Explicit deep-smoke selection or reviewer request | Default verification must be offline and fake-first. | Runtime boundary must be generic and fakeable so smoke can be added without changing architecture. |
| Stable JSON output for sprint validate/prompt/flow | Future CLI/API scope | Requirements exclude adding new stable JSON output unless already supported. | Keep service result structs deterministic enough that JSON can be added later without moving business logic into commands. |
| Generic workflow engine | Later product phase, if ever | PRD/TRD explicitly reject turning UltraPlan into a general workflow engine. | Keep stage model explicit and bounded to Planning Phase 2. |
| Shared planning validation/prompt package | Only after multiple modules need identical mechanical behavior | Architecture docs and handbook warn against global technical-layer packages. | Keep helper boundaries narrow and product-neutral if extracted later. |
| Execution, smoke, review automation, issues, and Git mutation | Later requirements revision | These are deferred product capabilities and excluded from this sprint. | Reject unsupported stages and do not create deferred artifacts in flow state. |

## Final Decisions

### Decision 1: Reason-Stage Ownership Stays In `internal/sprint`

- **Decision:** Implement selected-template detection, area reasoning manifests, final reasoning input assembly, reasoning validators, prompt rendering, and flow-state transitions in `internal/sprint`, using focused files such as `reasoning.go`, `prompts.go`, `flow.go`, `state.go`, `service.go`, and `store_fs.go` rather than new global planning, workflow, validation, prompts, reports, or catalog packages.
- **Rationale:** The sprint changes sprint planning behavior and state. The Architecture doc says `internal/sprint` owns planning artifacts, flow state, prompt rendering, and validation through `plan.md`; the TRD says planning validators are product behavior owned by `internal/sprint`. Keeping this behavior local avoids premature abstractions and preserves product context.
- **Study / Source Grounding:** `technical-handbook.md` lines 31-35, 51, 60-64; `01-project-structure` evidence from chezmoi `main.go:26-34`, gh-cli `cmd/gh/main.go:6`, k9s `cmd/root.go:112-127`, chezmoi `internal/chezmoi/chezmoi.go:1-2`, restic `internal/restic/repository.go:18`; `15-philosophy` evidence from gh-cli `pkg/cmd/factory/default.go:26-46` and restic `internal/backend/backend.go:19-90`.
- **Trade-Offs Accepted:** Local sprint code may duplicate small Markdown/path helper mechanics, but this is preferable to global packages before stable cross-module reuse exists.
- **Technical Debt / Future Impact:** Any extraction must be mechanical, product-neutral, and justified by repeated needs. Planning semantics must not be moved into platform or study packages.
- **Alternatives Rejected:** A global `internal/planning` or `internal/workflow` package is rejected because it contradicts selected architecture constraints and expands scope. A global `internal/validation` or `internal/prompts` package is rejected because it separates product rules from sprint state and repeats anti-patterns called out in project architecture. Reusing `internal/study` report validators is rejected because sprint reasoning artifacts are not study reports and the TRD forbids study semantic reuse.
- **Contracts Satisfied:** Architecture, Security, Testing; `requirements.md` lines 57-60, 92-99, 102; TRD sections 18.1, 18.2, 18.6.
- **Evidence Required:** Import review proves `internal/sprint` does not import `internal/study` and platform packages do not import product modules; code review proves no global planning/workflow/validation/prompt package was added; `go test ./...` passes.

### Decision 2: Selected Reasoning Templates Become A Deterministic Manifest

- **Decision:** Represent selected reasoning templates as a deterministic manifest derived from `sprint-index.md` and validated against `project-index.md`. Each manifest entry includes template name, selected template path, deterministic output path under `reasoning/`, selected contracts/evidence/protocol summaries needed for prompts, and workspace-safety/readability validation results. Ordering is stable by template name or selected output path.
- **Rationale:** The reason stage needs a single source of truth for `validate area-reasoning`, `prompt area-reasoning`, dry-run flow, runtime-backed area generation, and final reasoning prerequisites. A manifest avoids re-parsing ad hoc in commands while preserving selected-context-only behavior.
- **Study / Source Grounding:** `technical-handbook.md` lines 95-109 and 117-133 identify selected templates, conditional area reasoning, path trust boundaries, and validation ordering as design pressures/open questions. `13-security` evidence from opencode permission boundaries and k9s config schema validation supports explicit trust-boundary validation. `03-dependency-injection` evidence supports constructing this through explicit store/service seams rather than globals.
- **Trade-Offs Accepted:** A manifest adds a small internal model, but it centralizes deterministic behavior and diagnostics. It is not a generic catalog model and must remain scoped to sprint reasoning.
- **Technical Debt / Future Impact:** If later sprints add more reasoning templates, the manifest can support multiple entries without reopening architecture. It must not grow into a generic project catalog replacement.
- **Alternatives Rejected:** Inferring outputs by scanning `reasoning/*.md` is rejected because selected templates define required outputs, not existing files. Hardcoding `architecture.md` is rejected because the project index already models templates and later selections may differ. Allowing unselected template paths is rejected because it violates selected-context constraints and path safety.
- **Contracts Satisfied:** Architecture, Security, Workflows, Testing; `requirements.md` lines 39, 42, 46-49, 66-68, 94-95; sprint-index selected Architecture template lines 59-65.
- **Evidence Required:** Unit tests for no templates, one template, missing template path, unreadable path, unselected template reference, path escape, deterministic ordering, and selected output path; prompt and flow tests consume the same manifest.

### Decision 3: Reasoning Validators Enforce Artifact Contracts Without Brittle Prose Matching

- **Decision:** Implement separate validators for selected area reasoning files and final `reasoning.md`. Area validation requires file existence, non-empty content, no placeholders, required area-decision and trade-off sections, and no obvious references to unselected contracts/evidence/templates/protocols as decision sources. Final reasoning validation requires file existence, non-empty content, no placeholders, final decisions, expected evidence, assumptions and risks, selected requirement/contract/evidence mapping, and required selected-area references when templates are selected.
- **Rationale:** The product principle is that runtime success is not product success, but reasoning artifacts are editable Markdown and should not be rejected for harmless wording differences. Validators should check structural contract and selected-context discipline, not exact prose.
- **Study / Source Grounding:** `technical-handbook.md` lines 39, 54, 68, 72, 97-99, 121; `11-testing-strategy` evidence from chezmoi command tests, helm test fixtures, and go-task tests supports fixture/golden coverage without brittle assertions; `05-error-handling` evidence supports classified, actionable validation failures.
- **Trade-Offs Accepted:** Validators may miss subtle unsupported claims if they are not expressed as paths/names. This is acceptable because the sprint’s validator goal is enforceable artifact contract, not semantic truth certification.
- **Technical Debt / Future Impact:** The unselected-context check may need refinement after real artifacts reveal false positives or false negatives. The validator must report file, section, and selected-context findings clearly enough for repair.
- **Alternatives Rejected:** Existence-only validation is rejected because it would allow empty or placeholder artifacts to complete. Exact heading/prose validation is rejected because it would make generated editable Markdown brittle. Requiring `plan.md`, task checklists, smoke evidence, review evidence, issue artifacts, or Git state is rejected because those are non-goals for reasoning validation.
- **Contracts Satisfied:** Errors, Testing, Workflows, Security; `requirements.md` lines 31-36, 47-50, 66-69; TRD section 18.6.
- **Evidence Required:** `internal/sprint/reasoning_test.go` covers valid content, missing file, empty file, placeholders, missing decisions, missing trade-offs, missing expected evidence, missing assumptions/risks, missing required area references, unselected-context references, and workspace path escapes.

### Decision 4: Prompt Rendering Is Deterministic, Selected-Context-Only, And Non-Mutating Except Expected Output

- **Decision:** Add deterministic prompt rendering for `area-reasoning` and final `reasoning` stages in `internal/sprint/prompts.go`. Prompt previews include project slug, sprint slug, sprint path, requirements path, sprint-index path, technical-handbook path, selected reasoning template path and output path for area prompts, required area reasoning artifact paths for final prompts, final output path, selected contracts/evidence summary, required sections, selected-context-only rules, workspace-relative paths by default, and explicit no-mutation rules for all non-output artifacts.
- **Rationale:** The runtime needs enough context to generate the expected Markdown artifact, but prompt previews must be safe, deterministic, reviewable, and non-mutating. Prompt rendering is also a validation/debug surface and must not invoke runtime.
- **Study / Source Grounding:** `technical-handbook.md` lines 43-45, 55, 87-90, 101, 107, 113; `04-configuration-management` supports explicit config and redaction; `09-terminal-ux` supports non-interactive scriptable previews; `13-security` supports path and secret boundaries.
- **Trade-Offs Accepted:** Prompts may be verbose because they carry selected context and no-mutation constraints. This is acceptable because previews are deterministic and reviewable.
- **Technical Debt / Future Impact:** Prompt builders may grow long. Keep them table/manifest-oriented and covered by tests, but do not extract a global prompt package in this sprint.
- **Alternatives Rejected:** Letting runtime infer paths from workspace state is rejected because it weakens determinism and path safety. Embedding full unselected project docs or reports is rejected because reasoning must be selected-context-only. Allowing prompt previews to write normal output artifacts is rejected because preview commands must not generate `reasoning/*.md` or `reasoning.md` except an explicit preview path if such a flag already exists.
- **Contracts Satisfied:** Documentation, CLI Surface, Configuration, Security, LLM Evaluation / Cost / Safety; `requirements.md` lines 35-36, 42-45, 61, 69, 100-101.
- **Evidence Required:** Prompt tests assert required variables, selected template manifest, selected context summary, output path, required sections, workspace-relative paths, selected-only rules, no-mutation rules, no ANSI output, and no runtime invocation.

### Decision 5: Runtime-Backed Flow Uses Only The Generic Runtime Boundary And Validates After Generation

- **Decision:** Non-dry-run `flow --to reasoning` invokes only the generic platform runtime boundary from `internal/sprint`, first for required selected area reasoning artifacts and then for final reasoning. A stage is complete only when prerequisites validate, runtime execution succeeds, expected artifact exists, artifact validation passes, and `flow-state.json` is atomically updated. Validate, prompt, and dry-run flow paths never invoke runtime.
- **Rationale:** PRD and TRD establish that runtime success is not product success. The sprint must add runtime-backed generation while preserving safety, fake-first tests, and product-owned validation.
- **Study / Source Grounding:** `technical-handbook.md` lines 39-41, 68, 95, 103, 105, 111; `07-state-context` evidence supports context for runtime-backed work; `05-error-handling` supports classification and cause chains; `10-logging-observability` supports safe diagnostics; `13-security` supports no direct shell/OpenCode/provider/Git calls.
- **Trade-Offs Accepted:** Flow implementation has more gates than direct prompt execution, but this prevents false success and unsafe state completion.
- **Technical Debt / Future Impact:** Runtime diagnostics must be safe and may initially be concise. Future observability can add richer agentwrap metadata without changing reason-stage completion rules.
- **Alternatives Rejected:** Directly invoking OpenCode, agentwrap adapters, shell commands, provider APIs, or Git from `internal/sprint` is rejected because runtime execution must go through platform runtime only. Marking completion after runtime exit without validation is rejected because it violates product principles. Running final reasoning before required selected area reasoning validates is rejected because selected templates are prerequisites.
- **Contracts Satisfied:** LLM Runtime, LLM Evaluation / Cost / Safety, Security, Errors, Observability, Testing; `requirements.md` lines 31-41, 50, 65, 70-72, 97-101.
- **Evidence Required:** Flow tests cover dry-run no mutation, fake-runtime successful area generation, fake-runtime successful final generation, runtime success with missing artifact, runtime success with invalid artifact, runtime failure, validation failure after generation, selected-template preflight failure, and validate/prompt no-runtime command paths.

### Decision 6: Flow-State Remains Strict, Atomic, And Limited To Planning Phase 2 Stages

- **Decision:** Extend `flow-state.json` handling only for `requirements`, `sprint-index`, `technical-handbook`, `area-reasoning`, `reasoning`, and existing `plan` stage awareness. `area-reasoning` is `skipped` only when validated `sprint-index.md` selects no reasoning templates. When Architecture is selected, `area-reasoning` completes only after `projects/ultraplan-go/sprints/20-reason-stage/reasoning/architecture.md` exists and validates. Successful flow through `reasoning` marks validated prerequisites and reasoning stages complete, does not mark `plan` complete, and rejects unsupported future stages/state entries.
- **Rationale:** Flow state is durable product state. It must represent validated facts and current Planning Phase 2 scope, not runtime intention or future product phases.
- **Study / Source Grounding:** `technical-handbook.md` lines 95, 99, 105, 125; TRD sections 18.4 and 18.5; `10-logging-observability` evidence supports visible failed states; `15-philosophy` supports deliberate scope rejection.
- **Trade-Offs Accepted:** Unsupported stages fail hard rather than being treated as skipped or ignored. This preserves strict schema behavior and makes scope boundaries reviewable.
- **Technical Debt / Future Impact:** Later execution/smoke/review/issue stages will require an explicit schema and requirements revision rather than silent extension.
- **Alternatives Rejected:** Adding `implementation`, `execute`, `smoke`, `review`, `issues`, `.run-state.json`, `smoke.md`, or issue artifacts to current flow state is rejected because they are out of scope. Marking `area-reasoning` skipped when templates are selected but files are missing is rejected because selected templates require artifacts. Allowing state write failure to partially update state is rejected because prior sprint constraints require preserving the last valid state.
- **Contracts Satisfied:** Workflows, Persistence And Migrations, Observability, Errors; `requirements.md` lines 37, 39, 41, 51-56, 64, 70; TRD section 18.5.
- **Evidence Required:** State tests cover strict stage set, no-template skip, selected-template complete/fail, final reasoning complete, no plan completion, unsupported stage rejection, malformed/unsupported state rejection, atomic write success, and write failure preserving last valid state.

### Decision 7: CLI Wiring Is Thin, Scriptable, And Exit-Code Disciplined

- **Decision:** Extend `internal/app/sprint_commands.go` only to parse `validate area-reasoning`, `validate reasoning`, `prompt area-reasoning`, `prompt reasoning`, and `flow --to reasoning`, call `internal/sprint` services, render stdout/stderr, provide help, and map exit codes. Business rules, selected-context loading, validation, prompt construction, runtime invocation, and flow-state transitions remain in `internal/sprint`.
- **Rationale:** Command tests need to prove behavior at the surface, but command handlers should not own workflow rules. Keeping handlers thin preserves architecture and reduces test brittleness.
- **Study / Source Grounding:** `technical-handbook.md` lines 31, 37, 60, 83, 88; `02-command-architecture` evidence from gh-cli `pkg/cmdutil/factory.go:16-43`, helm `pkg/cmd/install.go:132-145`, and opencode `cmd/root.go:49-183` warning; `06-io-abstraction` evidence from gh-cli IOStreams; `09-terminal-ux` and `10-logging-observability` support scriptable output.
- **Trade-Offs Accepted:** More service request/result types may be needed, but they keep command code focused and testable.
- **Technical Debt / Future Impact:** If command output formats grow, the service result shape should remain stable enough to support JSON later without moving business logic into `internal/app`.
- **Alternatives Rejected:** Putting validation/prompt/flow logic directly in command handlers is rejected because large handlers are a known anti-pattern. Adding interactive prompts is rejected because commands must be scriptable. Adding stable new JSON output is rejected unless already supported without scope expansion.
- **Contracts Satisfied:** CLI Surface, Errors, Documentation, Testing, Observability; `requirements.md` lines 31-38, 60-65, 71, 125, 129, 135.
- **Evidence Required:** Command tests cover help output, valid/invalid stages, stdout/stderr separation, no ANSI output, exit codes `0`, `2`, `5`, runtime failure mapping, validate/prompt no-runtime behavior, and flow invoking only a fake generic runtime seam.

## Expected Evidence

| Evidence Type | Required Evidence | Source / Command / Review Check |
| --- | --- | --- |
| Unit Tests | Selected-template manifest detection, deterministic ordering, path safety, area/final validation, placeholder detection, selected-context rejection, flow-state strict load/write behavior, unsupported stage parsing. | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/reasoning_test.go`, `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/flow_test.go` |
| Prompt Tests | Prompt previews include required variables, selected context summary, output paths, required sections, selected-only rules, no-mutation rules, workspace-relative paths, and do not invoke runtime. | `go test ./internal/sprint` or equivalent package tests |
| Flow Tests | Dry-run does not mutate artifacts/state; fake runtime writes selected area and final artifacts; runtime success with missing/invalid artifact fails; runtime failure fails safely; validation failure after generation fails; no-template area skip; selected-template area requirement; no plan completion. | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/flow_test.go` |
| Command Tests | `validate area-reasoning`, `validate reasoning`, `prompt area-reasoning`, `prompt reasoning`, and `flow --to reasoning` command behavior, help text, exit codes, stdout/stderr separation, no ANSI output, no runtime on validate/prompt, fake runtime on non-dry-run flow. | `/home/antonioborgerees/coding/ultraplan-go/internal/app/sprint_commands_test.go` |
| Architecture Review | `internal/sprint` owns reason-stage behavior; `internal/app` stays thin; platform packages remain product-agnostic; `internal/study` is not imported; no global planning/workflow/validation/prompt package is added. | Architecture Review protocol selected in `sprint-index.md` |
| Security Review | Workspace-safe path handling, selected-template path validation, no direct shell/Git/OpenCode/provider/agent adapter calls from `internal/sprint`, redacted prompt/diagnostic output. | Code review and targeted tests for path escapes and diagnostics |
| Offline Verification | Full unit/fixture/command suite passes. | `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go` |
| Race Verification | Race suite passes. | `go test -race ./...` from `/home/antonioborgerees/coding/ultraplan-go` |
| Build Verification | CLI builds. | `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go` |
| Sprint Review | Implementation evidence, deviations, verification results, and architecture review conclusions are recorded after execution. | `projects/ultraplan-go/sprints/20-reason-stage/review.md` |

## Assumptions And Risks

| Item | Type | Impact | Mitigation / Follow-Up |
| --- | --- | --- | --- |
| No Architecture area reasoning artifact was available as an input to this document. | Risk | Product validation for final reasoning should fail when Architecture is selected but the area artifact is missing. This manually authored reasoning cannot claim area-specific conclusions. | Implementation must require selected area artifacts before final reasoning completion. Create or validate `projects/ultraplan-go/sprints/20-reason-stage/reasoning/architecture.md` before relying on automated flow completion. |
| `requirements.md` acceptance criteria have no explicit stable IDs. | Assumption | Decision mappings use file line ranges and acceptance-criteria descriptions instead of source-defined requirement IDs. | Plan/review should preserve these mappings or add stable IDs in a future requirements format if needed. |
| Existing Sprint 17-19 APIs provide enough hooks for flow-state and prerequisite validation. | Assumption | If current implementation shape differs, the plan may need small integration steps but must not reopen architecture. | Keep changes minimal and local; adapt names to existing code while preserving decisions. |
| Selected-context reference detection cannot prove semantic absence of all unselected ideas. | Risk | Generated reasoning could mention unsupported concepts without obvious path/name markers. | Enforce selected paths/names/templates/protocols, required sections, and review checks; treat semantic review as part of sprint review evidence. |
| Fake runtime may not reflect every real OpenCode behavior. | Risk | Default tests could miss integration issues in real runtime output/event behavior. | Keep runtime boundary generic and fakeable; defer real-runtime smoke to explicit deep-smoke selection or reviewer request. |
| Strict flow-state schema may reject older experimental state files. | Risk | Users with prototype or manually edited state may need repair. | Return actionable validation errors and preserve last valid state; do not silently migrate unsupported stage data in this sprint. |
| Prompt previews may expose local absolute paths if diagnostics require them. | Risk | Artifacts could become less portable or reveal local details. | Default to workspace-relative paths; use absolute paths only for local diagnostics where required and avoid secrets/full environment. |

## Implementation Constraints

- `internal/sprint` owns selected-template detection, area/final reasoning validation, prompt rendering, selected-context loading, runtime-backed reason-stage flow, and `flow-state.json` transitions.
- `internal/app` command handlers parse arguments, call sprint services, render output, and map exit codes only.
- `internal/sprint` may depend on `internal/project`, `internal/workspace`, generic platform config/redaction/runtime/filesystem helpers, and atomic file helpers where they already exist.
- `internal/sprint` must not import `internal/study`, reuse study source/dimension/report/rating/summary/scheduler/run-loop semantics, or extract those semantics into shared packages.
- `internal/platform/*` packages must not import `internal/sprint`, `internal/project`, `internal/study`, or `internal/app`.
- Non-dry-run flow must invoke only the generic platform runtime boundary and must not directly call shell commands, Git, OpenCode, provider APIs, or agentwrap concrete adapters from `internal/sprint`.
- Validate, prompt, and dry-run flow paths must not invoke runtime and must not create or overwrite reasoning artifacts.
- Prompt preview content must be deterministic, use workspace-relative paths by default, avoid secrets, include no-mutation rules, and write only to an explicit preview path if that behavior already exists.
- Area reasoning may be skipped only when validated `sprint-index.md` selects no reasoning templates.
- Final reasoning completion requires selected area reasoning artifacts when templates are selected.
- Flow through `reasoning` must not mark `plan` complete and must not create or validate implementation, smoke, review, issue, or Git artifacts.
- `flow-state.json` writes must be atomic; write failure preserves the last valid state and returns non-zero.
- Unsupported stages such as `implementation`, `execute`, `smoke`, `review`, `issues`, `.run-state.json`, `smoke.md`, `smoke.json`, generated `review.md`, `issues.md`, and `issues.json` must be rejected as current Planning Phase 2 behavior.
- Tests must be deterministic, offline, fake-first, and free of network calls, provider credentials, real OpenCode execution, Git mutation, and long sleeps by default.

## Plan Handoff

`plan.md` must execute these decisions. It must not invent architecture, scope, or decisions beyond this document.

The plan must carry forward:

- final decisions in Decisions 1-7
- selected contracts and `requirements.md` line mappings recorded above
- expected evidence for unit, prompt, flow, command, architecture, security, race, build, and review checks
- the missing Architecture area artifact risk and selected-template completion constraint
- strict no-plan/no-execution/no-smoke/no-review-automation/no-issues/no-Git-mutation boundaries
- required Architecture Review and Sprint Review protocols selected in `sprint-index.md`

## Phase Exit Criteria

- [x] Selected context was read and used.
- [x] Area-specific reasoning document availability was checked and the missing selected Architecture artifact was recorded as a risk.
- [x] Area-specific reasoning conclusions are not claimed because no area artifact was present.
- [x] Contracts are explicitly mapped to decisions and expected evidence.
- [x] Final decisions are clear enough for `plan.md` to execute without reopening architecture.
- [x] Expected evidence is specific and reviewable.
