# Sprint Reasoning: Distill Stage

> Project: `ultraplan-go`
> Sprint: `19-distill-stage`
> Output: `projects/ultraplan-go/sprints/19-distill-stage/reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/sprints/19-distill-stage/requirements.md`, `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/sprints/19-distill-stage/sprint-index.md`, `projects/ultraplan-go/sprints/19-distill-stage/technical-handbook.md`, `templates/sprint-reasoning.md`; no area reasoning files were present under `projects/ultraplan-go/sprints/19-distill-stage/reasoning/`.

This document decides. It synthesizes selected context, handbook evidence, area-specific reasoning availability, and contracts into final sprint decisions.

It does not replace `sprint-index.md`, `technical-handbook.md`, or `reasoning/*.md`.

## Sprint Purpose

- **Goal:** Implement the planning distill stage so `technical-handbook.md` can be generated, previewed, flowed, and validated from valid `sprint-index.md` selected evidence only, with flow-state updates through `technical-handbook` and without making implementation decisions in the handbook artifact.
- **Non-Goals:** Do not generate or validate `reasoning/*.md`, final `reasoning.md`, or `plan.md`; do not make final sprint decisions inside `technical-handbook.md`; do not read or cite unselected evidence reports; do not implement sprint execution, smoke investigation, review automation, issue tracking, Git mutation, study-service reuse, global workflow/prompt/validation packages, new stable JSON output, or real-runtime smoke requirements.
- **Depends On:** Sprint 18 select-stage behavior and `sprint-index.md`; Sprint 17 sprint artifact domain, supported Planning Phase 2 stages, and atomic `flow-state.json`; Sprint 16 project catalog boundary; Sprint 14 diagnostics, redaction, and exit-code conventions; Sprint 9 generic runtime boundary; selected `go-cli-study` evidence reports cataloged by `project-index.md` and selected by `sprint-index.md`.

## Selected Context And Pre-Reasoning Artifacts

| Artifact | Path | How It Was Used |
| --- | --- | --- |
| Project Index | `projects/ultraplan-go/project-index.md` | Confirmed project scope, Phase 2 planning workflow, cataloged source docs, active contracts, selected-evidence pool, reasoning template pool, and review protocols. Established that `internal/project` owns catalog behavior and `internal/sprint` owns planning artifacts and flow state through `plan.md`. |
| Sprint Requirements | `projects/ultraplan-go/sprints/19-distill-stage/requirements.md` | Primary sprint contract. Defined required outputs, acceptance criteria, constraints, dependencies, non-goals, review expectations, expected files, exit codes, and verification commands. Requirement IDs in this document refer to the ordered acceptance criteria as `REQ-AC-01` through `REQ-AC-40`. |
| PRD | `projects/ultraplan-go/docs/PRD.md` | Grounded product principles: runtime success is not product success, generated files stay editable, Phase 2 supports planning through `plan.md`, and execution/smoke/review/issues/Git mutation are deferred. |
| TRD | `projects/ultraplan-go/docs/TRD.md` | Grounded technical rules for module ownership, Planning Phase 2 stages, flow state, validators, prompt rendering, runtime boundary, command shape, exit codes, diagnostics, path safety, atomic writes, and testing. |
| Architecture | `projects/ultraplan-go/docs/ARCHITECTURE.md` | Grounded dependency direction, module-owned behavior, platform/product separation, `sprint -> project/workspace/platform/runtime`, and prohibition on extracting study semantics into planning. |
| Sprint Index | `projects/ultraplan-go/sprints/19-distill-stage/sprint-index.md` | Established selected contracts, selected evidence reports, selected Architecture reasoning template, required review protocols, and explicit excluded context. Confirmed that the distill stage must consume selected evidence and stop before reasoning and plan generation. |
| Technical Handbook | `projects/ultraplan-go/sprints/19-distill-stage/technical-handbook.md` | Supplied study-backed patterns, trade-offs, anti-patterns, examples, design pressures, and open questions. Used as the main evidence distillation input for final sprint decisions. |
| Sprint Reasoning Template | `templates/sprint-reasoning.md` | Provided required document structure and final-decision expectations. |

## Area-Specific Reasoning Inputs

The sprint index selects the Architecture reasoning template with output path `projects/ultraplan-go/sprints/19-distill-stage/reasoning/architecture.md`. No area-specific reasoning files currently exist for this sprint.

| Area | Reasoning Document | Key Conclusion | Evidence Basis | Impact On Final Decision |
| --- | --- | --- | --- | --- |
| Architecture | Not present | Final sprint reasoning must resolve architecture decisions directly rather than summarizing an area file. | `requirements.md`, `sprint-index.md`, `technical-handbook.md`, `TRD.md` section 18, and `ARCHITECTURE.md` module ownership rules. | Decisions below explicitly settle package ownership, dependency direction, runtime boundary, and rejected architecture alternatives so `plan.md` can execute without reopening architecture. |

## Sprint Technical Handbook Summary

- **Relevant Patterns:** Keep CLI entrypoints thin; keep product behavior module-owned; use explicit composition roots and injectable collaborators; validate artifacts after runtime success; use selected manifests and path validation before runtime; keep IO captureable and output scriptable; propagate context through runtime-backed flow; use fake-first and command-level tests.
- **Important Trade-Offs:** Keep distill behavior inside `internal/sprint` even if Markdown parsing/validation code duplicates some local mechanics; validate selected evidence strictly rather than discovering files permissively; use explicit service/store/runtime seams rather than direct runtime/filesystem calls in command handlers; prefer structured diagnostics and deterministic output over richer terminal UX.
- **Warnings / Anti-Patterns:** Do not put selected-evidence rules in CLI handlers; do not read unselected evidence opportunistically; do not treat runtime success as stage success; do not bypass IO abstractions; do not add global mutable config or global validation/prompt/workflow packages; do not use decision language in `technical-handbook.md`.
- **Evidence Confidence:** High for module ownership, command thinness, path safety, validation-after-runtime, fakes, and test strategy because multiple selected reports agree and cite concrete mature Go CLI source references. Medium for terminal UX richness because this sprint intentionally favors scriptable output and does not implement interactive TUI behavior.

## Contracts Applied

| Contract / Requirement ID | Constraint | Decision Impact | Expected Evidence |
| --- | --- | --- | --- |
| Architecture | Product behavior stays module-owned; platform packages stay product-agnostic; CLI handlers remain thin. | Distill behavior lives in `internal/sprint`; command parsing lives in `internal/app`; platform runtime does not learn sprint semantics. | Import review, architecture review, tests proving command handlers delegate to sprint service. |
| Errors | Diagnostics must be actionable, cause-preserving, deterministic, and mapped to stable exit codes. | Validation, selected evidence, usage, runtime, and flow-state failures use classified errors and safe messages. | Command tests for exit codes `0`, `2`, `5`, runtime failure mapping, stdout/stderr separation. |
| Configuration | Runtime configuration and redaction stay at established boundaries. | Prompt and flow use existing config/runtime resolution; prompts and diagnostics avoid secrets and unnecessary absolute paths. | Prompt tests and command tests for redaction and workspace-relative paths. |
| Observability | Diagnostics and status must be calm, scriptable, and safe. | Flow reports dry-run, selected evidence, generation success/failure, and state transitions without ANSI by default. | Command tests and review checks for output discipline. |
| Security | Workspace-safe paths, selected-only evidence, no shell/OpenCode/Git direct invocation, secret redaction. | Selected evidence loader validates catalog membership, sprint selection, readability, and path safety before reads or runtime. | Loader tests for missing, unreadable, uncataloged, unselected, and escaping paths; import review for no direct shell/OpenCode/Git calls. |
| Testing | Default verification must be deterministic, offline, fake-first, and broad enough to catch side effects. | Unit, fixture, fake-runtime, and command tests are required before review. | `go test ./...`, `go test -race ./...`, `go build ./cmd/ultraplan`; targeted tests in `handbook_test.go`, `flow_test.go`, `sprint_commands_test.go`. |
| Documentation | CLI help and editable Markdown artifact behavior must remain understandable. | Help exposes `technical-handbook` for validate, prompt, and flow; prompt instructs editable Markdown output. | Help-output tests and review of generated prompt content. |
| CLI Surface | Command shape, stage support, usage errors, stdout/stderr separation, and exit-code behavior are constrained. | Add thin wiring for `validate technical-handbook`, `prompt technical-handbook`, and `flow --to technical-handbook`; reject unsupported stages. | Command tests for help, unsupported stages, output streams, no ANSI, and exit codes. |
| LLM Runtime | Non-dry-run flow goes only through generic runtime boundary. | `internal/sprint` constructs sprint-owned prompt/request data and calls platform runtime; it does not invoke OpenCode, provider APIs, shell commands, or Git. | Fake-runtime tests and import review. |
| Workflows | Flow dry-run, runtime flow, failures, and next-stage readiness must be explicit. | Flow validates prerequisites before runtime, validates generated handbook after runtime, atomically updates `flow-state.json`, and stops before later reasoning/plan completion. | Flow tests for dry-run, success, missing artifact, invalid artifact, runtime failure, validation failure, write failure, next-stage readiness. |
| Persistence And Migrations | `flow-state.json` is strict, versioned, and atomically written; unsupported legacy stages are not supported. | State model includes only Planning Phase 2 stages and rejects `implementation`, `execute`, `smoke`, `review`, and `issues` as current supported stages. | State tests, fixture tests with legacy unsupported stages, atomic-write failure tests. |
| `REQ-AC-01` to `REQ-AC-03` | Validate, prompt, and dry-run paths must not invoke runtime or mutate artifacts/state as completion. | Service methods separate validation, prompt preview, dry-run, and non-dry-run flow. | Command and flow tests with runtime spy/fake proving no calls and no writes where forbidden. |
| `REQ-AC-04` to `REQ-AC-18` | Runtime-backed flow and handbook validation requirements. | Prerequisite validation, selected-evidence loading, prompt manifest, post-generation validation, no-decision checks, and completion criteria are implemented as sprint behavior. | Handbook and flow tests for every validation rule and state transition. |
| `REQ-AC-19` to `REQ-AC-24` | Flow-state and stage support constraints. | State transitions stop at Planning Phase 2 and keep area reasoning ready or skipped correctly. | Flow-state tests and review of supported stage constants. |
| `REQ-AC-25` to `REQ-AC-31` | Package boundaries, CLI discipline, outputs, and exit codes. | `internal/sprint` owns rules; `internal/app` remains thin; platform stays generic; outputs are deterministic and safe. | Import review and command tests. |
| `REQ-AC-32` to `REQ-AC-40` | Test coverage and verification commands. | Implementation plan must create targeted unit, flow, command, race, and build evidence. | `go test ./...`, `go test -race ./...`, `go build ./cmd/ultraplan`, plus sprint review evidence. |

## Repos Studied / Source Evidence Used

| Source / Repo / Report | Concrete Reference | Relevant Finding | Why It Matters For This Sprint | Used In Decision(s) |
| --- | --- | --- | --- | --- |
| `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md`; `cmd/gh/main.go:6`, `main.go:16`, `cmd/root.go:112-127`, `internal/chezmoi/chezmoi.go:1-2` | Mature Go CLIs keep entrypoints thin and route work inward with unidirectional dependencies. | Supports keeping distill rules in `internal/sprint` and CLI wiring in `internal/app`. | Decisions 1, 6 |
| `02-command-architecture` | `studies/go-cli-study/reports/final/02-command-architecture.md`; `pkg/cmdutil/factory.go:16-43`, `pkg/cmd/install.go:132-145`, `cmd/restic/cmd_backup.go:84-115` | Command factories and injected dependencies separate command wiring from behavior. | Supports command tests around thin handlers and service seams. | Decisions 1, 6 |
| `03-dependency-injection` | `studies/go-cli-study/reports/final/03-dependency-injection.md`; `internal/ghcmd/cmd.go:52-132`, `pkg/cmd/factory/default.go:26-46`, `executor.go:93-114`, `internal/app/app.go:42-81` | Explicit collaborators and composition roots improve testability without DI frameworks. | Supports service/store/runtime seams and fake-runtime testing. | Decisions 1, 5, 7 |
| `04-configuration-management` | `studies/go-cli-study/reports/final/04-configuration-management.md`; `internal/cmd/config.go:2253-2287`, `internal/flags/flags.go:314-327`, `internal/config/config.go:609-641` | Config precedence and validation should be explicit after merging. | Supports using existing platform configuration/redaction boundaries and failing before runtime when prerequisites are invalid. | Decisions 4, 5 |
| `05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md`; `cmd/age/tui.go:37-54`, `cmd/restic/main.go:199-209`, `internal/ghcmd/cmd.go:44-49`, `pkg/action/uninstall.go:232-254` | User diagnostics, cause preservation, exit-code mapping, and aggregated failures are important in CLI flows. | Supports validation diagnostics, selected-evidence finding aggregation, and command exit-code mapping. | Decisions 2, 3, 5, 6 |
| `06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md`; `pkg/iostreams/iostreams.go:551-568`, `internal/ui/mock.go:10-53`, `cmd/ls/ls.go:42`, `cli.go:47` | Captureable stdout/stderr and filesystem seams make command behavior testable. | Supports no direct stdout/stderr in business rules and command-level output tests. | Decisions 6, 7 |
| `07-state-context` | `studies/go-cli-study/reports/final/07-state-context.md`; `internal/ghcmd/cmd.go:142`, `task.go:89`, `pkg/cmd/install.go:333-347`, `cmd/restic/cleanup.go:24-38` | Context propagation and explicit state are needed for runtime-backed operations. | Supports context-aware flow and safe failed-state recording. | Decisions 5, 7 |
| `09-terminal-ux` | `studies/go-cli-study/reports/final/09-terminal-ux.md`; `internal/cmd/prompt.go:124-137`, `internal/prompter/prompter.go:16-81`, `signals.go:11-31` | Scriptable non-TTY behavior and cancellation-aware output beat rich terminal UX for default automation. | Supports calm, deterministic, no-ANSI output and dry-run previews. | Decisions 4, 6 |
| `10-logging-observability` | `studies/go-cli-study/reports/final/10-logging-observability.md`; `fs/log/slog.go:21`, `internal/logging/logging.go:31-66`, `internal/slogs/keys.go:6-231`, `pkg/iostreams/iostreams.go:52-54` | Structured diagnostics should be separated from user output. | Supports stdout/stderr separation and safe flow diagnostics. | Decisions 5, 6 |
| `11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md`; `internal/cmd/main_test.go:64-174`, `acceptance/acceptance_test.go:26-29`, `pkg/httpmock/stub.go:35-199`, `internal/test/test.go:43` | Table-driven, fixture, fake, command, and golden tests provide CLI confidence. | Supports targeted tests for validators, prompt previews, flow, and command behavior. | Decision 7 |
| `13-security` | `studies/go-cli-study/reports/final/13-security.md`; `internal/execext/exec.go:59-66`, `internal/permission/permission.go:44-108`, `internal/options/secret_string.go:15-20`, `pkg/registry/transport.go:37-41`, `internal/config/json/validator.go:146` | Trust boundaries around paths, shell execution, permissions, and secrets must be explicit. | Supports selected-only evidence loading, workspace safety, redaction, and no direct shell/OpenCode/Git calls. | Decisions 2, 4, 5 |
| `15-philosophy` | `studies/go-cli-study/reports/final/15-philosophy.md`; `pkg/cmd/factory/default.go:26-46`, `internal/chezmoi/system.go:25-45`, `internal/pubsub/broker.go:10-19`, `website/src/docs/experiments/index.md:17-21` | Sustainable Go CLI design accepts complexity only under real product pressure. | Supports minimal local sprint implementation instead of broad abstractions. | Decisions 1, 7 |

## Trade-Off And Debt Analysis

### Accepted Trade-Offs

| Trade-Off | Benefit | Cost / Constraint Accepted | Why Acceptable Now | Revisit Trigger |
| --- | --- | --- | --- | --- |
| Keep all distill rules in `internal/sprint` instead of shared planning packages | Preserves module ownership and avoids premature abstractions. | Some Markdown heading, placeholder, and diagnostic mechanics may be duplicated locally. | Sprint scope is one planning stage and project architecture explicitly prefers local ownership first. | A later stage repeats identical mechanical behavior in multiple modules and extraction can stay product-neutral. |
| Strict selected-evidence validation instead of permissive discovery | Prevents unselected or uncataloged evidence from entering the handbook. | Requires more parsing, ordering, and diagnostics. | The sprint goal is selected-evidence-only distillation; correctness matters more than convenience. | Sprint index format changes or catalog explicitly supports external evidence manifests. |
| Validation after runtime success instead of trusting runtime completion | Prevents invalid handbook artifacts from being marked complete. | Adds post-generation failure states and tests. | PRD and TRD both state runtime success is not product success. | None for this product principle; only validator details may evolve. |
| Focused prompt assertions instead of full-prompt golden lock-in | Keeps required prompt content stable while allowing prose improvements. | Some non-required prompt wording can change without test failure. | Requirements warn against brittle prose assertions beyond required variables and constraints. | Prompt regressions occur that focused assertions fail to catch. |
| Fake-first runtime tests instead of real OpenCode smoke | Keeps default verification offline, deterministic, and credential-free. | Does not prove real provider behavior in this sprint. | Requirements explicitly exclude real-runtime smoke from default verification. | A future completed sprint selects deep smoke evidence. |
| Aggregate selected-evidence validation findings where safe | Gives users a fuller repair list. | Slightly more validation plumbing than fail-fast. | Evidence report `05-error-handling` supports multi-error patterns and requirements ask for actionable diagnostics. | Diagnostics become noisy or expensive for very large manifests. |

### Potential Technical Debt

| Debt / Shortcut | Why It Might Accrue | Current Mitigation | Owner / Follow-Up |
| --- | --- | --- | --- |
| Local Markdown validation heuristics in `internal/sprint` | Required-section and no-decision checks may grow across future planning stages. | Keep helpers small, unexported, and stage-specific; extract only mechanical helpers after repeated use. | Future reasoning/plan sprints. |
| No-decision wording checks may be imperfect | Overly broad checks can reject valid trade-off prose or miss subtle decisions. | Use targeted banned decision sections/phrases and tests for allowed evidence language. | Sprint 20 reasoning validator can refine shared policy if needed. |
| Selected-evidence parsing depends on current `sprint-index.md` tables | Table heading or formatting changes could break loader behavior. | Validate with deterministic diagnostics and fixtures; keep sprint-index validator authoritative. | Project/sprint planning maintenance if index format evolves. |
| Flow-state schema strictness may reject legacy sprint states | Existing Sprint 19 `flow-state.json` may contain unsupported legacy stages. | Reject unsupported stages as current supported state and preserve last valid state on write failures. | Future migration sprint if old state compatibility becomes a product need. |
| Prompt preview output is text-first, not stable JSON | Automation may later want machine-readable manifests. | Expose deterministic text and internal manifest; do not add stable JSON unless requirements select it. | Future CLI surface sprint. |

### Future Considerations

| Consideration | Deferred Until | Reason Deferred | What Should Be Preserved Now |
| --- | --- | --- | --- |
| Area reasoning generation and validation | Sprint 20 | Sprint 19 stops at distillation and next-stage readiness. | Stage constants and flow transitions for `area-reasoning`, `reasoning`, and `plan` remain explicit but not completed by this sprint. |
| Final plan generation and validation | Sprint 21 | Plan belongs after final reasoning decisions. | `technical-handbook` prompt and validation must not create task plans. |
| Real OpenCode smoke verification | Future selected smoke review | Requirements require offline fake-first verification only. | Runtime calls go through generic platform runtime so later smoke can exercise real adapter without sprint code changes. |
| Stable JSON output for sprint commands | Future CLI contract revision | Explicitly out of scope unless existing app infrastructure already supports it without expansion. | Deterministic internal result structs and text output discipline. |
| Shared Markdown validation package | After multiple product modules prove identical mechanical needs | Current architecture rejects global validation packages by default. | Keep logic cohesive and tests clear enough to extract later without semantic coupling. |
| Planning execution, smoke, review automation, issues, and Git mutation | Later product phase | PRD/TRD defer these capabilities. | Do not model unsupported stages as current Planning Phase 2 state. |

## Final Decisions

### Decision 1: Distill Behavior Lives In `internal/sprint`

- **Decision:** Implement technical-handbook input manifests, selected-evidence loading, handbook validation, prompt rendering, flow behavior, and flow-state transitions as `internal/sprint` product behavior. Use focused files matching the required outputs: `handbook.go`, `prompts.go`, `flow.go`, `service.go`, `store_fs.go`, and `state.go`. `internal/app` only parses command arguments, delegates to sprint services, renders output, and maps exit codes.
- **Rationale:** The distill stage transforms sprint planning state and sprint artifacts. `ARCHITECTURE.md` and `TRD.md` section 18 assign sprint planning artifacts, validators, prompts, and flow state to `internal/sprint`. Keeping behavior in the owning module prevents global-layer coupling and prevents CLI handlers from becoming the product workflow.
- **Study / Source Grounding:** `technical-handbook.md` selected patterns cite `01-project-structure`, `02-command-architecture`, `03-dependency-injection`, and `15-philosophy`. Concrete references include `cmd/gh/main.go:6`, `cmd/root.go:112-127`, `pkg/cmdutil/factory.go:16-43`, `internal/ghcmd/cmd.go:52-132`, and `pkg/cmd/factory/default.go:26-46` showing thin entrypoints, command factories, and explicit collaborators.
- **Trade-Offs Accepted:** Accept some local sprint-specific Markdown and validation code instead of a global package. Accept a slightly larger sprint package to preserve ownership and dependency clarity.
- **Technical Debt / Future Impact:** Local helpers may be candidates for extraction later, but only as mechanical platform helpers after repeated concrete use. This decision avoids the larger debt of global `validation`, `prompts`, `workflow`, or `planning` packages.
- **Alternatives Rejected:** Rejected `internal/planning`, `internal/workflow`, `internal/validation`, or `internal/prompts` packages because sprint requirements and architecture prohibit global packages for this scope. Rejected implementing selected-evidence and validator rules in `internal/app` because command layers must remain thin and testable. Rejected importing `internal/study` because Phase 2 may reference study outputs as catalog paths but must not reuse study semantics.
- **Contracts Satisfied:** Architecture, CLI Surface, Testing, Documentation, `REQ-AC-25`, `REQ-AC-26`, `REQ-AC-27`, `REQ-AC-28`.
- **Evidence Required:** Import review confirms `internal/sprint` does not import `internal/study`; platform packages do not import product modules; `internal/app` command handlers delegate to service methods. Command tests prove CLI handlers do not own validation or runtime behavior.

### Decision 2: Selected Evidence Is Loaded Through A Strict Catalog-And-Selection Manifest

- **Decision:** Build a deterministic selected-evidence manifest from valid `requirements.md`, valid `project-index.md`, and valid `sprint-index.md`. Load only reports selected in `sprint-index.md` that also exist in `project-index.md`, resolve inside the workspace unless explicitly cataloged external, verify readability, sort deterministically, and fail with actionable diagnostics on missing, unreadable, uncataloged, unselected, duplicate, or escaping paths.
- **Rationale:** The sprint goal is selected-evidence-only distillation. The project index is a catalog, not a sprint plan; the sprint index is the selected context. Any broader discovery would violate the sprint contract and make the handbook unverifiable.
- **Study / Source Grounding:** `technical-handbook.md` selected patterns cite `13-security`, `04-configuration-management`, and `05-error-handling`. Concrete references include `internal/config/json/validator.go:146`, `internal/execext/exec.go:59-66`, `internal/permission/permission.go:44-108`, and `pkg/action/uninstall.go:232-254`, supporting explicit validation, path trust boundaries, and aggregated actionable diagnostics.
- **Trade-Offs Accepted:** Accept more validation code and fixtures to avoid accidental evidence expansion. Prefer deterministic manifest order over filesystem discovery convenience.
- **Technical Debt / Future Impact:** Manifest parsing is tied to current project-index and sprint-index Markdown table shape. Future schema changes should update project/sprint validators and fixtures rather than adding permissive fallback discovery.
- **Alternatives Rejected:** Rejected scanning `studies/**/reports/final/*.md` because it would include unselected reports. Rejected trusting `sprint-index.md` paths without project-index membership because it would bypass the project catalog boundary. Rejected reading all project docs as handbook evidence because selected evidence reports, not docs, are the handbook evidence source; docs remain constraints and context.
- **Contracts Satisfied:** Security, Errors, Configuration, Testing, `REQ-AC-04`, `REQ-AC-11`, `REQ-AC-12`, `REQ-AC-15`, `REQ-AC-16`, `REQ-AC-17`, `REQ-AC-34`.
- **Evidence Required:** `handbook_test.go` covers multiple selected reports, missing files, unreadable files, project-index mismatches, unselected evidence references, escaping paths, duplicate ordering, and deterministic manifest diagnostics. Flow tests prove selected-evidence failures happen before runtime.

### Decision 3: Handbook Validation Is Semantic Enough To Gate The Distill Stage

- **Decision:** Validate `technical-handbook.md` for existence, non-empty content, template placeholders, required sections, selected studies/reports, relevant patterns, trade-offs, anti-patterns, open questions, evidence pointers, selected evidence traces, unselected evidence references, and implementation-decision language. Allow evidence-backed observations, warnings, open questions, and trade-off framing when they do not decide architecture or create implementation tasks.
- **Rationale:** The handbook is editable Markdown and evidence distillation, not final reasoning or a plan. Validation must be stronger than file existence because runtime success alone is insufficient and because later reasoning depends on a constrained evidence artifact.
- **Study / Source Grounding:** `technical-handbook.md` uses `05-error-handling`, `11-testing-strategy`, and `15-philosophy` to support validation gates, regression tests, and deliberate complexity. Concrete references include `errors/errors.go:47-50`, `cmd/restic/main.go:199-209`, `internal/cmd/main_test.go:64-174`, `internal/test/test.go:43`, and `website/src/docs/experiments/index.md:17-21`.
- **Trade-Offs Accepted:** Accept heuristic no-decision checks even though natural language cannot be perfectly classified. Tests must include both prohibited decision sections and allowed trade-off/evidence wording to reduce false positives.
- **Technical Debt / Future Impact:** No-decision checks may require refinement when Sprint 20 adds final reasoning validation. Keep rules local and named so future refinements are safe.
- **Alternatives Rejected:** Rejected existence-only validation because it would allow empty or irrelevant artifacts to pass. Rejected requiring exact prose or full golden content because the handbook is editable Markdown. Rejected allowing uncataloged citations because it would break selected-evidence traceability.
- **Contracts Satisfied:** Errors, Security, Testing, Documentation, `REQ-AC-13`, `REQ-AC-14`, `REQ-AC-15`, `REQ-AC-16`, `REQ-AC-17`, `REQ-AC-18`, `REQ-AC-33`.
- **Evidence Required:** `handbook_test.go` covers valid handbook content, missing file, empty file, placeholders, each missing required section, missing selected studies/reports, missing patterns, missing trade-offs, missing anti-patterns, missing open questions, missing evidence pointers, implementation-decision wording, allowed evidence wording, selected trace success, and unselected evidence rejection.

### Decision 4: Prompt Rendering Is Deterministic, Previewable, Selected-Only, And Non-Mutating

- **Decision:** Render the `technical-handbook` prompt deterministically from the sprint requirements, validated sprint index, selected evidence manifest, required handbook sections, no-decision instructions, selected-evidence-only constraints, output path, and no-mutation instructions. Use workspace-relative paths in prompt content unless absolute paths are required for local diagnostics. `prompt technical-handbook` writes no handbook artifact and invokes no runtime; if an existing preview-output path feature is present, it may write only to that explicit preview path.
- **Rationale:** Prompt preview is a planning safety feature. It must show exactly what would be sent to runtime without side effects, and it must constrain the runtime to generate editable Markdown at the handbook output path without mutating catalogs, docs, evidence reports, source repositories, Git state, implementation files, `sprint-index.md`, reasoning artifacts, or `plan.md`.
- **Study / Source Grounding:** `technical-handbook.md` cites `04-configuration-management`, `09-terminal-ux`, `10-logging-observability`, and `13-security`. Concrete references include `internal/cmd/config.go:2253-2287`, `internal/cmd/prompt.go:124-137`, `pkg/iostreams/iostreams.go:52-54`, `internal/options/secret_string.go:15-20`, and `pkg/registry/transport.go:37-41`.
- **Trade-Offs Accepted:** Use focused assertions for required prompt variables and constraints rather than freezing the whole prompt prose. Accept text-first preview output and defer stable JSON preview unless a pre-existing app path supports it without scope expansion.
- **Technical Debt / Future Impact:** Prompt prose can evolve; required manifest fields and no-mutation constraints must remain stable. Future stages can reuse the pattern inside `internal/sprint` without extracting a global prompt package.
- **Alternatives Rejected:** Rejected runtime-backed prompt preview because preview must be side-effect-free. Rejected prompt content with broad absolute paths because generated artifacts should remain relocatable and avoid leaking local details. Rejected embedding unselected evidence report content because the sprint index is the selection boundary.
- **Contracts Satisfied:** Configuration, Security, CLI Surface, Documentation, Testing, `REQ-AC-02`, `REQ-AC-08`, `REQ-AC-09`, `REQ-AC-10`, `REQ-AC-35`.
- **Evidence Required:** Prompt tests assert project slug, sprint slug, sprint path, requirements path, sprint-index path, output path, selected evidence manifest, required sections, no-decision instructions, selected-evidence-only instructions, no-mutation instructions, deterministic ordering, workspace-relative paths, no runtime invocation, and no handbook write.

### Decision 5: Flow Through `technical-handbook` Uses Runtime Only After Prerequisite Validation And Completes Only After Post-Generation Validation And Atomic State Update

- **Decision:** `flow --to technical-handbook` validates `requirements.md`, `sprint-index.md`, and selected evidence before runtime. Dry-run reports stages and selected evidence without runtime, handbook writes, or completion marks. Non-dry-run invokes only the generic platform runtime boundary, then requires `technical-handbook.md` to exist and pass handbook validation before atomically updating `flow-state.json`. Runtime failure, missing artifact, invalid artifact, selected-evidence failure, validation failure, or state-write failure returns non-zero and records or preserves safe state according to the failure point.
- **Rationale:** The product principle is that runtime success is not product success. Flow-state is durable product state and must represent validated artifact completion, not just attempted generation.
- **Study / Source Grounding:** `technical-handbook.md` cites `03-dependency-injection`, `05-error-handling`, `07-state-context`, `10-logging-observability`, `11-testing-strategy`, and `13-security`. Concrete references include `executor.go:93-114`, `internal/ghcmd/cmd.go:44-49`, `internal/ghcmd/cmd.go:142`, `task.go:89`, `internal/logging/logging.go:31-66`, and `internal/execext/exec.go:59-66`.
- **Trade-Offs Accepted:** Accept explicit state-transition and failure-case tests. Keep flow orchestration local to `internal/sprint` rather than adding a workflow engine.
- **Technical Debt / Future Impact:** Flow code may grow as later stages are added; keep stage transition helpers explicit and Planning Phase 2-specific so future additions do not require architecture changes. Unsupported legacy stages remain rejected rather than preserved as supported behavior.
- **Alternatives Rejected:** Rejected marking complete immediately after runtime exit because generated content still needs validation. Rejected direct OpenCode, shell, provider API, or Git invocation from `internal/sprint` because runtime execution belongs behind the platform runtime boundary. Rejected adding a generic workflow engine because this sprint only needs staged planning flow through `technical-handbook`.
- **Contracts Satisfied:** LLM Runtime, Workflows, Persistence And Migrations, Security, Errors, Observability, `REQ-AC-03`, `REQ-AC-04`, `REQ-AC-05`, `REQ-AC-06`, `REQ-AC-07`, `REQ-AC-18`, `REQ-AC-19`, `REQ-AC-20`, `REQ-AC-21`, `REQ-AC-22`, `REQ-AC-23`, `REQ-AC-24`, `REQ-AC-36`.
- **Evidence Required:** `flow_test.go` covers dry-run, fake-runtime success, runtime success with missing artifact, runtime success with invalid artifact, runtime failure, selected-evidence failure before runtime, validation failure after generation, flow-state write failure, safe failure messages, no later-stage completion, and atomic write preservation.

### Decision 6: CLI Wiring Is Thin, Scriptable, And Strict About Supported Stages

- **Decision:** Add command wiring for `ultraplan sprint <project> <sprint> validate technical-handbook`, `prompt technical-handbook`, and `flow --to technical-handbook`. Help output documents the stage consistently. Unsupported stages such as `execute`, `implementation`, `smoke`, `review`, and `issues` return usage exit code `2`. Validation and prompt commands return `0` on success and `5` on validation or selected-evidence failures. Runtime failures during non-dry-run flow map to the established runtime exit code. Text output separates stdout and stderr, remains deterministic, avoids ANSI by default, and redacts sensitive values.
- **Rationale:** The sprint command surface already exists; Sprint 19 extends it with the distill stage without expanding future execution behavior. Thin command wiring keeps business rules testable and consistent with the architecture.
- **Study / Source Grounding:** `technical-handbook.md` cites `01-project-structure`, `02-command-architecture`, `05-error-handling`, `06-io-abstraction`, `09-terminal-ux`, and `10-logging-observability`. Concrete references include `pkg/cmdutil/factory.go:16-43`, `pkg/iostreams/iostreams.go:551-568`, `cmd/restic/main.go:199-209`, `internal/ghcmd/cmd.go:44-49`, and `internal/cmd/prompt.go:124-137`.
- **Trade-Offs Accepted:** Keep output plain and scriptable rather than adding richer progress UI. Do not add stable new JSON output unless already supported by existing app infrastructure without scope expansion.
- **Technical Debt / Future Impact:** Help text and command tests must evolve as later planning stages are implemented. Unsupported future stages remain intentionally rejected until selected by requirements.
- **Alternatives Rejected:** Rejected adding `execute`, `smoke`, `review`, or `issues` aliases because they are deferred product phases. Rejected command handlers that perform file reads, validation, or runtime calls directly. Rejected ANSI-rich progress output because this sprint requires scriptable default output.
- **Contracts Satisfied:** CLI Surface, Errors, Observability, Documentation, Testing, `REQ-AC-01`, `REQ-AC-02`, `REQ-AC-03`, `REQ-AC-29`, `REQ-AC-30`, `REQ-AC-31`, `REQ-AC-37`.
- **Evidence Required:** `sprint_commands_test.go` covers validate, prompt, flow, help output, unsupported-stage exit `2`, validation exit `5`, success exit `0`, runtime-free validate/prompt, fake-runtime flow seam, stdout/stderr separation, no ANSI output, deterministic ordering, and safe redaction.

### Decision 7: Verification Is Offline, Fake-First, And Reviewable

- **Decision:** Implement targeted unit, fixture, fake-runtime, and command tests before review. Required verification is `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`. Sprint review evidence must record implementation results, deviations, command output expectations, import-boundary review, and acceptance-criteria conclusions. No real OpenCode, provider credentials, network, Git mutation, or smoke harness execution is required for this sprint.
- **Rationale:** Requirements explicitly demand deterministic offline verification and fake-runtime tests. The selected testing evidence supports command-level and fixture-heavy coverage for CLI workflows.
- **Study / Source Grounding:** `technical-handbook.md` cites `11-testing-strategy`, `03-dependency-injection`, `06-io-abstraction`, and `15-philosophy`. Concrete references include `internal/cmd/main_test.go:64-174`, `acceptance/acceptance_test.go:26-29`, `pkg/httpmock/stub.go:35-199`, `internal/test/test.go:43`, `pkg/iostreams/iostreams.go:551-568`, and `internal/chezmoitest/chezmoitest.go:86-92`.
- **Trade-Offs Accepted:** Fake tests do not prove actual provider behavior, but they prove product logic, side-effect boundaries, validation gates, and command behavior deterministically. Race verification is included even though the sprint should add little concurrency because runtime-backed flow and state writes touch shared abstractions.
- **Technical Debt / Future Impact:** Real-runtime confidence is deferred to an explicitly selected smoke review. The fake-runtime seam must remain representative of the platform runtime boundary.
- **Alternatives Rejected:** Rejected relying on manual CLI trials because acceptance criteria require repeatable tests. Rejected real OpenCode smoke as default verification because credentials/network are non-goals. Rejected only unit tests because command behavior, stdout/stderr, exit codes, and no-runtime guarantees are core requirements.
- **Contracts Satisfied:** Testing, Observability, CLI Surface, Security, `REQ-AC-32`, `REQ-AC-33`, `REQ-AC-34`, `REQ-AC-35`, `REQ-AC-36`, `REQ-AC-37`, `REQ-AC-38`, `REQ-AC-39`, `REQ-AC-40`.
- **Evidence Required:** Passing `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan`; review evidence in `projects/ultraplan-go/sprints/19-distill-stage/review.md`; targeted tests in `internal/sprint/handbook_test.go`, `internal/sprint/flow_test.go`, and `internal/app/sprint_commands_test.go`.

## Expected Evidence

| Evidence Type | Required Evidence | Source / Command / Review Check |
| --- | --- | --- |
| Tests | Handbook validator covers valid content, missing file, empty file, placeholders, missing required sections, evidence traces, unselected evidence, no-decision rejection, and allowed evidence language. | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/handbook_test.go`; `go test ./...` |
| Tests | Selected evidence loader covers multiple reports, missing/unreadable files, project-index mismatch, unselected evidence, path escape, deterministic ordering, and actionable diagnostics. | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/handbook_test.go`; `go test ./...` |
| Tests | Prompt rendering covers required variables, manifest, output path, required sections, no-decision instructions, selected-evidence-only instructions, no-mutation instructions, workspace-relative paths, and no runtime/write side effects. | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/handbook_test.go` or prompt-specific sprint tests; `go test ./...` |
| Tests | Flow covers dry-run, fake-runtime successful generation, runtime success with missing artifact, runtime success with invalid artifact, runtime failure, selected-evidence failure before runtime, validation failure after runtime, flow-state write failure, and next-stage readiness/skipping. | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/flow_test.go`; `go test ./...` |
| Tests | Command tests cover validate, prompt, flow, help output, exit codes, stdout/stderr separation, no ANSI output, safe redaction, unsupported-stage rejection, and runtime-free validate/prompt. | `/home/antonioborgerees/coding/ultraplan-go/internal/app/sprint_commands_test.go`; `go test ./...` |
| Runtime | Non-dry-run flow uses a fake runtime seam in tests and only the generic platform runtime boundary in production code. | Fake-runtime assertions; import review for no direct OpenCode, shell, provider API, or Git calls in `internal/sprint`. |
| Persistence | Flow-state writes are atomic and preserve last valid state on write failure. Unsupported legacy stages are not supported Planning Phase 2 stages. | State and flow tests; code review of `internal/sprint/state.go` and `internal/sprint/flow.go`. |
| Review | Architecture boundaries are preserved. | Architecture Review protocol: `internal/sprint` owns distill rules; `internal/app` remains thin; `internal/platform/*` imports no product packages; no global planning/workflow/validation/prompt package. |
| Review | Sprint acceptance criteria, deviations, verification commands, and residual risks are recorded. | Sprint Review protocol; `projects/ultraplan-go/sprints/19-distill-stage/review.md`. |
| Verification | Offline test suite passes. | `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go`. |
| Verification | Race suite passes. | `go test -race ./...` from `/home/antonioborgerees/coding/ultraplan-go`. |
| Verification | CLI builds. | `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`. |

## Assumptions And Risks

| Item | Type | Impact | Mitigation / Follow-Up |
| --- | --- | --- | --- |
| Existing sprint-index parsing can be extended or reused to expose selected evidence deterministically. | Assumption | If not true, the implementation needs a focused parser update before handbook work. | Keep parsing in `internal/sprint`, add fixtures for current Sprint 19 index shape, and fail with diagnostics rather than guessing. |
| `project-index.md` catalog validation already provides enough structure to compare selected evidence paths. | Assumption | If project catalog APIs are too weak, selected evidence validation could duplicate catalog parsing. | Use `internal/project` exported catalog behavior where available; add only focused catalog read methods if necessary. |
| No area reasoning files are available for this sprint reasoning document. | Assumption | Final reasoning must settle architecture directly. | This document records final decisions with architecture grounding; plan must not require a missing area file to reopen decisions. |
| No-decision validation may reject legitimate wording. | Risk | Users could be blocked when handbook prose uses words like "prefer" or "decision" in an evidence context. | Test allowed trade-off and guidance wording; target section headers and final-decision phrases rather than broad word bans. |
| Markdown table parsing for selected evidence may be brittle. | Risk | Minor table formatting changes could produce false validation failures. | Use deterministic parser tests and actionable diagnostics; do not silently fall back to broader discovery. |
| Flow-state write failures can occur after artifact generation succeeds. | Risk | Artifact exists but stage is not complete, which may surprise users. | Return non-zero, preserve last valid state, record safe failure, and require rerun/revalidation to mark complete. |
| Existing Sprint 19 `flow-state.json` may contain unsupported legacy stages. | Risk | Loading strict state could fail. | Treat unsupported stages as invalid current Planning Phase 2 state per requirements; do not preserve them as supported stages. |
| Runtime result shape from platform runtime may not expose all desired diagnostics. | Risk | Flow failure messages could be sparse. | Use safe platform/runtime result fields and cause chains; do not parse OpenCode/native streams directly. |
| Help output tests can become brittle. | Risk | Small wording changes cause failures. | Assert required command names, supported stages, unsupported-stage behavior, and key usage lines rather than every sentence unless existing test style prefers goldens. |

## Implementation Constraints

- `internal/sprint` owns technical-handbook validation, selected-evidence loading, prompt rendering, and flow-state transitions for this distill stage.
- `internal/app` must only parse arguments, call sprint services, render output, and map exit codes.
- `internal/sprint` may depend on `internal/project`, `internal/workspace`, generic platform configuration/redaction, generic platform runtime, and generic atomic file helpers.
- `internal/sprint` must not import `internal/study` or reuse study source, dimension, report, rating, summary, scheduler, or run-loop behavior.
- `internal/platform/*` must not import `internal/sprint`, `internal/project`, `internal/study`, or `internal/app`.
- Do not introduce global `internal/catalog`, `internal/planning`, `internal/workflow`, `internal/stages`, `internal/validation`, `internal/reports`, or `internal/prompts` packages.
- Validate `requirements.md`, `sprint-index.md`, and selected evidence before runtime-backed `flow --to technical-handbook`.
- `validate technical-handbook`, `prompt technical-handbook`, and `flow --to technical-handbook --dry-run` must not invoke runtime.
- Runtime-backed flow must use only the generic platform runtime boundary; no direct OpenCode, agentwrap adapter, shell, provider API, or Git invocation from `internal/sprint`.
- Runtime success alone is insufficient; completion requires generated artifact existence, handbook validation success, and atomic flow-state update.
- Selected evidence loading must read only reports selected by `sprint-index.md` and cataloged by `project-index.md`.
- Handbook validation must fail on missing file, empty file, placeholders, missing required sections, missing selected evidence traces, unselected evidence references, and implementation-decision language.
- Prompt content must prefer workspace-relative paths, include selected evidence manifest and required constraints, instruct editable Markdown output at `projects/ultraplan-go/sprints/19-distill-stage/technical-handbook.md`, and forbid mutation of catalogs, docs, evidence, source repos, Git state, implementation files, `sprint-index.md`, reasoning artifacts, or `plan.md`.
- Flow-state support is limited to `requirements`, `sprint-index`, `technical-handbook`, `area-reasoning`, `reasoning`, and `plan`.
- Do not support `implementation`, `execute`, `smoke`, `review`, `issues`, `.run-state.json`, `smoke.md`, `smoke.json`, product-generated `review.md`, `issues.md`, or `issues.json` as current Planning Phase 2 stages.
- Text output must be calm, deterministic, no ANSI by default, stdout/stderr-separated, and redacted.
- Exit code `0` applies to successful validation, prompt rendering, dry-run, and completed flow; exit code `2` applies to malformed arguments and unsupported stages; exit code `5` applies to validation and selected-evidence failures; runtime failures use the established runtime exit code.
- Default verification must be offline, deterministic, fake-first, and free of real OpenCode, provider credentials, network calls, long sleeps, shell/Git mutation, and smoke harness requirements.

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
- [x] Area-specific reasoning documents were completed where applicable or explicitly marked not present.
- [x] Area-specific reasoning conclusions are reflected or explicitly overridden in final decisions.
- [x] Contracts are explicitly mapped to decisions and expected evidence.
- [x] Final decisions are clear enough for `plan.md` to execute.
- [x] Expected evidence is specific and reviewable.
