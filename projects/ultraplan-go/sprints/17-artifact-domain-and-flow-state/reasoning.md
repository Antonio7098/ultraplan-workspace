# Sprint Reasoning: Sprint Artifact Domain and Flow State

> Project: `ultraplan-go`
> Sprint: `17-artifact-domain-and-flow-state`
> Output: `projects/ultraplan-go/sprints/17-artifact-domain-and-flow-state/reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/17-artifact-domain-and-flow-state/requirements.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/sprints/17-artifact-domain-and-flow-state/sprint-index.md`, `projects/ultraplan-go/sprints/17-artifact-domain-and-flow-state/technical-handbook.md`, `templates/sprint-reasoning.md`

This document decides. It synthesizes selected context, handbook evidence, area-specific reasoning, and contracts into final sprint decisions.

It does not replace `sprint-index.md`, `technical-handbook.md`, or `reasoning/*.md`.

## Sprint Purpose

- **Goal:** Model planning-stage sprint artifacts and durable `flow-state.json` state through `plan.md`, and expose runtime-free sprint status inspection through `ultraplan sprint <project> <sprint> status`.
- **Non-Goals:** Do not generate, repair, or content-validate planning artifacts beyond existence/non-empty status inspection; do not implement `sprint validate`, `sprint prompt`, `sprint flow`, runtime-backed generation, implementation execution, smoke investigation, automated review, issue tracking, Git mutation, JSON status output, hosted service, browser UI, TUI, local API server, metrics exporter, or project-management workflow.
- **Depends On:** Sprint 16 project domain and project index boundaries for project resolution; Sprint 14 diagnostics and exit-code discipline; Sprint 12 strict state and atomic-write precedent; Sprint 9 runtime boundary decisions; PRD, TRD, and Architecture docs for Phase 2 planning scope.

## Selected Context And Pre-Reasoning Artifacts

| Artifact | Path | How It Was Used |
| --- | --- | --- |
| Project Index | `projects/ultraplan-go/project-index.md` | Established Phase 2 ownership rules, selected evidence pool, contract pool, review protocols, and non-goals for execution/smoke/review/issues/Git mutation. |
| Sprint Requirements | `projects/ultraplan-go/sprints/17-artifact-domain-and-flow-state/requirements.md` | Treated as the authoritative sprint contract for required files, accepted stage/status values, strict flow-state behavior, command surface, tests, verification commands, and non-goals. |
| Product Requirements | `projects/ultraplan-go/docs/PRD.md` | Grounded the product scope: governed project/sprint planning through `plan.md`, local filesystem artifacts, explicit state/errors, and deferred implementation execution, smoke, review, issues, and Git mutation. |
| Technical Requirements | `projects/ultraplan-go/docs/TRD.md` | Grounded `internal/sprint` ownership, planning-stage model, `flow-state.json` fields/statuses, path safety, atomic writes, validators, command requirements, and test expectations. |
| Architecture | `projects/ultraplan-go/docs/ARCHITECTURE.md` | Grounded module-owned behavior, allowed `sprint -> project/workspace` dependency direction, platform/product separation, and avoidance of global technical-layer packages. |
| Sprint Index | `projects/ultraplan-go/sprints/17-artifact-domain-and-flow-state/sprint-index.md` | Selected contracts, selected evidence reports, excluded runtime/workflow/performance/config contexts, selected Architecture reasoning template, and required review protocols. |
| Technical Handbook | `projects/ultraplan-go/sprints/17-artifact-domain-and-flow-state/technical-handbook.md` | Supplied distilled evidence from selected Go CLI reports, relevant patterns, anti-patterns, trade-offs, concrete source references, and open questions resolved in this document. |
| Sprint Reasoning Template | `templates/sprint-reasoning.md` | Provided required reasoning structure and quality bar for decisions, evidence, trade-offs, assumptions, constraints, and plan handoff. |

## Area-Specific Reasoning Inputs

The sprint index selected an Architecture reasoning template at `projects/ultraplan-go/sprints/17-artifact-domain-and-flow-state/reasoning/architecture.md`, but no `reasoning/` directory or area reasoning file exists for this sprint. No area-specific reasoning conclusions were available to summarize. This final sprint reasoning therefore resolves the architecture questions directly and records the missing selected area reasoning as a risk and review point.

| Area | Reasoning Document | Key Conclusion | Evidence Basis | Impact On Final Decision |
| --- | --- | --- | --- | --- |
| Architecture | `projects/ultraplan-go/sprints/17-artifact-domain-and-flow-state/reasoning/architecture.md` | No file exists; no area-specific conclusion was available. | Sprint index selected Architecture reasoning; technical handbook provided architecture evidence from `01-project-structure`, `02-command-architecture`, `03-dependency-injection`, `06-io-abstraction`, `12-extensibility`, and `15-philosophy`. | Final decisions in this document explicitly decide `internal/sprint` ownership, dependency direction, runtime-free status behavior, flow-state persistence, path representation, and deferred stages so `plan.md` can proceed without reopening architecture. |

## Sprint Technical Handbook Summary

- **Relevant Patterns:** Use thin CLI command wiring; keep domain, artifact rules, flow-state persistence, and status derivation inside `internal/sprint`; construct dependencies explicitly; use IO/filesystem seams only where they make command and write-failure tests practical; render deterministic plain text; use behavior-focused tests.
- **Important Trade-Offs:** Accept sprint-owned behavior over a generic workflow package; accept strict flow-state rejection over lenient migration; accept plain text over rich terminal formatting; accept explicit command enumeration over registry/plugin mechanisms; accept targeted filesystem seams over a broad abstraction layer.
- **Warnings / Anti-Patterns:** Do not put sprint business rules in command handlers; do not create global `planning`, `workflow`, `stages`, `validation`, `reports`, or `prompts` packages; do not use package-level mutable state; do not flatten error causes; do not trust persisted paths by string convention; do not invoke runtime, shell, study services, source cloning, prompt generation, or project-index mutation for status.
- **Evidence Confidence:** High for thin CLI, module ownership, command/IO testing, strict errors, path safety, and behavior tests because the handbook cites multiple mature Go CLIs and selected evidence reports. Medium for terminal UX and extensibility because those findings are more contextual, but they still align with this sprint's small, runtime-free command scope.

## Contracts Applied

| Contract / Requirement ID | Constraint | Decision Impact | Expected Evidence |
| --- | --- | --- | --- |
| Architecture contract | Product behavior must remain module-owned with clear dependency direction and thin entrypoints. | `internal/sprint` owns sprint stages, artifact paths, state, status derivation, and validation findings; `internal/app` only wires and renders. | Import review, `internal/sprint/doc.go`, absence of global workflow/planning packages, command tests. |
| Errors contract | Diagnostics must be actionable, causes preserved internally, and exit codes meaningful. | Flow-state validation failures return exit code `5`; usage/reference issues use established usage/workspace codes; errors wrap causes and preserve state paths internally. | Tests for invalid refs, malformed JSON, unsafe paths, unsupported stages/statuses, write failure diagnostics, stdout/stderr separation. |
| Observability contract | Status output must truthfully report state, paths, failures, and diagnostics without runtime side effects. | Plain text status includes project, sprint, sprint root, flow-state path, every stage, status, artifact path, and safe errors. | Golden or stable command-output tests, no ANSI assertions, no runtime invocation assertions. |
| Security contract | User paths and persisted paths must not escape workspace or sprint root; diagnostics must redact unsafe details. | Slug validation, workspace-safe sprint resolution, sprint-root-contained artifact paths, strict persisted path validation, safe error detail only. | Path escape tests, invalid slug tests, unsafe flow-state path tests, import/runtime review. |
| Testing contract | Tests must be deterministic, offline, fake-first, and behavior-focused. | Sprint tests cover discovery, stage ordering, strict loading, atomic writes, status derivation, and command behavior without OpenCode/network/Git. | `go test ./...`, `go test -race ./...`, command fixture tests. |
| Documentation contract | Package docs and CLI help must document supported behavior and deferred behavior. | Add `internal/sprint/doc.go`, top-level help, and sprint command help that describe planning-stage-only scope. | Documentation review and command help tests. |
| CLI Surface contract | Commands must be script-friendly, support help, separate stdout/stderr, and return stable exit codes. | Implement only `ultraplan sprint <project> <sprint> status` in this sprint; no `validate`, `prompt`, or `flow`. | Command tests for help, malformed args, exit codes, stdout/stderr, deterministic output. |
| Persistence And Migrations contract | Durable machine state must be versioned, strict, atomically written, and reject unsupported schemas. | `flow-state.json` uses schema version `1`, strict loading, no compatibility paths, same-directory temp write, flush/close, rename, best-effort parent sync. | State load/write tests, write failure tests preserving last valid state. |
| `requirements.md:30-36` | `internal/sprint` must own sprint planning artifacts and exactly six stages/artifacts. | Domain model defines only `requirements`, `sprint-index`, `technical-handbook`, `area-reasoning`, `reasoning`, and `plan`; supported artifacts are fixed. | Domain tests for stage order and unsupported implementation/review stages. |
| `requirements.md:37-42` | `flow-state.json` fields, statuses, strict loading, and atomic writes are required. | Implement strict schema fields, one entry per stage, allowed statuses only, safe paths, project/sprint matching, and atomic persistence. | Unit tests for malformed JSON, schema version, duplicate/unknown stages, unknown statuses, escapes, mismatches, atomic write failure. |
| `requirements.md:43-49` | Status derivation, `area-reasoning` skip rules, status output, exit codes, and runtime-free inspection are required. | Status refresh derives from artifacts without runtime calls; recorded failures remain failed; `area-reasoning` skips only with explicit no-template selection. | Fixture tests for empty/partial/complete artifacts, skipped area reasoning, failed preservation, exit code `0`/`5`, runtime-free assertions. |
| `requirements.md:50-52` | Sprint must not import study or create broad global packages; platform remains product-agnostic. | Reuse `internal/project` only for project boundary and `internal/workspace` for path safety; no `internal/study` dependency; platform packages do not import sprint/project/app. | Import review and package layout review. |
| `requirements.md:53-58` | Required sprint and command tests plus verification commands must pass. | Plan must include tests for valid/invalid state, references, artifacts, command behavior, and build/race verification. | `go test ./...`, `go test -race ./...`, `go build ./cmd/ultraplan`. |

## Repos Studied / Source Evidence Used

| Source / Repo / Report | Concrete Reference | Relevant Finding | Why It Matters For This Sprint | Used In Decision(s) |
| --- | --- | --- | --- | --- |
| `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md`; chezmoi `main.go:26-34`; gh-cli `cmd/gh/main.go:6`; yq `cmd/root.go:9` | Mature Go CLIs keep entrypoints thin and protected behavior under internal/domain packages with one-way dependencies. | Supports `internal/sprint` as the owning package and `internal/app` as wiring only. | Decisions 1, 6 |
| `02-command-architecture` | `studies/go-cli-study/reports/final/02-command-architecture.md`; gh-cli `pkg/cmdutil/factory.go:16-43`; helm `pkg/cmd/install.go:132-145`; restic `cmd/restic/main.go:37-114` | Command layers construct commands, parse flags, delegate, and render; large `RunE` functions are a warning. | Supports a thin `sprint status` adapter and explicit command registration. | Decisions 5, 6 |
| `03-dependency-injection` | `studies/go-cli-study/reports/final/03-dependency-injection.md`; chezmoi `internal/cmd/config.go:362`; opencode `internal/app/app.go:42-81`; go-task `executor.go:22-24` | Explicit composition roots and constructor/function option patterns improve testability without global state. | Supports injected clock/filesystem/IO seams and avoids package-level mutable state. | Decisions 1, 4, 5 |
| `05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md`; gh-cli `internal/ghcmd/cmd.go:44-49`; go-task `errors/errors.go:47-50`; restic `internal/errors/fatal.go:10` | Typed/classified errors and cause chains support actionable diagnostics and exit-code mapping. | Supports strict flow-state errors, state path retention, and exit code `5` for validation failures. | Decisions 3, 5 |
| `06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md`; gh-cli `pkg/iostreams/iostreams.go:551-568`; restic `internal/fs/interface.go:10-31` | IO and filesystem seams are valuable when they enable command and persistence tests. | Supports focused seams for output rendering, artifact inspection, atomic write failure injection, and no hardcoded terminal output. | Decisions 4, 5 |
| `07-state-context` | `studies/go-cli-study/reports/final/07-state-context.md`; gh-cli `internal/ghcmd/cmd.go:142`; restic `internal/global/global.go:46-89` | Durable state and operation context should be explicit and inspectable. | Supports explicit `FlowState`, `StageState`, state loading, refresh, and context-aware service methods without globals. | Decisions 1, 3, 4 |
| `09-terminal-ux` | `studies/go-cli-study/reports/final/09-terminal-ux.md`; chezmoi `internal/cmd/prompt.go:20-256`; gh-cli `internal/prompter/prompter.go:16-81`; yq `pkg/yqlib/color_print.go:7-9` | CLI-first status commands should stay plain and script-friendly unless interaction needs justify richer terminal behavior. | Supports deterministic non-ANSI status text and help output. | Decision 5 |
| `10-logging-observability` | `studies/go-cli-study/reports/final/10-logging-observability.md`; gh-cli `pkg/iostreams/iostreams.go:52-54`; helm `internal/logging/logging.go:71` | Useful observability depends on stdout/stderr separation, stable fields, and safe diagnostics. | Supports status output shape, safe stage errors, and command tests for stdout/stderr separation. | Decisions 3, 5 |
| `11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md`; chezmoi `internal/cmd/main_test.go:64-174`; helm `internal/test/test.go:43`; restic `internal/backend/mock/backend.go:14-26` | Strong CLIs use table-driven, command-level, golden-compatible, and mock-backed behavior tests. | Supports the test suite shape and expected evidence. | Decision 7 |
| `12-extensibility` | `studies/go-cli-study/reports/final/12-extensibility.md`; restic `cmd/restic/main.go:77-106`; rclone `fs/registry.go:22`; helm `internal/plugin/runtime_subprocess.go:65-79` | Extension points should be deliberate; closed static command lists are often better for bounded command surfaces. | Supports avoiding plugin/registry/stage-engine mechanisms for one status command and six fixed stages. | Decisions 1, 2, 6 |
| `13-security` | `studies/go-cli-study/reports/final/13-security.md`; gh-cli `pkg/surveyext/editor_manual.go:23`; restic `internal/options/secret_string.go:15-20`; k9s `internal/config/json/validator.go:146`; gdu `internal/common/ignore.go:16-37` | Path, command, and secret boundaries must be enforced in code, not convention. | Supports slug validation, path containment, safe persisted paths, no runtime/shell invocation, and redacted diagnostics. | Decisions 2, 3, 5 |
| `15-philosophy` | `studies/go-cli-study/reports/final/15-philosophy.md`; age `age.go:18`; gh-cli `pkg/cmd/factory/default.go:26-46`; restic `internal/backend/backend.go:19-90` | Accept complexity only when it serves current scope; use cohesive module-oriented abstractions. | Supports minimal sprint-owned design and rejecting speculative workflow/platform abstractions. | Decisions 1, 2, 6 |

## Trade-Off And Debt Analysis

### Accepted Trade-Offs

| Trade-Off | Benefit | Cost / Constraint Accepted | Why Acceptable Now | Revisit Trigger |
| --- | --- | --- | --- | --- |
| Sprint-owned stage/state model instead of a generic workflow engine | Keeps product semantics close to sprint artifacts and prevents study/workflow leakage. | Some mechanical state or validation concepts may be repeated later. | This sprint has one bounded status use case and six fixed planning stages. | `sprint flow` or another product module needs the same proven mechanical helper without product semantics. |
| Strict flow-state loading instead of lenient compatibility | Keeps durable state trustworthy and makes corrupt or unsafe state loud. | Manually edited or older state files can fail status with exit code `5`. | Requirements explicitly demand strict version/stage/status/path/project/sprint checks and no migration path. | A future migration requirement defines supported older state shapes. |
| Status refresh from artifacts after valid or missing state | Lets status reflect human-editable Markdown artifacts and recover from absent state. | Refresh logic can mask stale state if strict loading and failure preservation are not separated. | Refresh is limited to existence/non-empty artifact inspection and only runs after no invalid persisted state is accepted. | Future validators add content semantics or runtime-backed stage generation. |
| Plain deterministic text instead of rich terminal UI | Easier to test, script-friendly, no ANSI surprises, and aligned with status-only scope. | Less visual polish for humans. | Sprint requires deterministic status and no TUI/browser/API behavior. | Product requirements add interactive sprint dashboard or structured JSON status. |
| Targeted filesystem seams instead of a broad filesystem interface | Keeps implementation smaller while still testing write failures and unsafe paths. | Some path helpers remain concrete and may need refactoring later. | Current volatile boundaries are artifact inspection, flow-state load/write, clock, and command IO. | Multiple packages need the same mechanical filesystem abstraction or tests become brittle. |

### Potential Technical Debt

| Debt / Shortcut | Why It Might Accrue | Current Mitigation | Owner / Follow-Up |
| --- | --- | --- | --- |
| Minimal `sprint-index.md` parsing for `area-reasoning` skip detection | Sprint 17 must not implement full sprint-index validation, but skip behavior needs evidence that no templates were selected. | Implement only a narrow detector for an explicit no-template selection; otherwise do not skip missing area reasoning. | Sprint 18 validation should own full sprint-index selection checks. |
| Persisted `failed` stages cannot be cleared by a dedicated command in this sprint | Status preserves recorded failures, while generation/validation commands are deferred. | Document that recorded failures remain failed until a future stage action or manual valid state change; tests assert no silent clearing. | Sprint 18-21 planning flow commands should define failure clearing semantics. |
| No JSON status output | PRD/TRD list JSON as broader desirable status behavior, but requirements exclude it unless already supported without scope expansion. | Keep renderer isolated enough that future JSON can use the same status summary without changing domain rules. | Future CLI surface sprint defining stable status JSON schema. |
| Missing selected Architecture area reasoning file | Sprint index selected an area reasoning artifact that does not exist. | This reasoning document resolves architecture decisions directly and records the gap as a risk for review. | Sprint review should confirm whether direct final reasoning is acceptable or request backfill. |

### Future Considerations

| Consideration | Deferred Until | Reason Deferred | What Should Be Preserved Now |
| --- | --- | --- | --- |
| `sprint validate [stage]` | Sprint 18 or later | Current sprint inspects status only and does not content-validate planning files. | Keep validation findings/domain types narrow and avoid command assumptions that block validators. |
| `sprint prompt <stage>` and `sprint flow --to <stage>` | Sprint 19-21 planning generation work | Runtime-backed generation and prompt preview are explicit non-goals. | Keep `internal/sprint` stage/order/artifact APIs stable enough for future prompt and flow services. |
| Runtime integration for planning stages | Later runtime-backed planning sprint | Sprint status must not call agentwrap/OpenCode/network. | Keep platform/runtime product-agnostic and do not import runtime from status path. |
| Implementation execution, smoke, review, issues, and Git mutation | Later requirements revision beyond Phase 2 planning | PRD, TRD, project index, sprint index, and requirements all defer these. | Do not model them as supported stages or accept them in strict flow-state loading. |
| Shared atomic-write helper | When repeated by sprint/project/study without product semantics | This sprint can reuse existing generic helper if present, but should not extract prematurely. | Keep atomic write behavior mechanical and product-agnostic if extracted. |

## Final Decisions

### Decision 1: `internal/sprint` Owns Sprint Planning State

- **Decision:** Add `internal/sprint` as the product module that owns sprint discovery, sprint reference resolution, planning stage/domain types, artifact path rules, `flow-state.json` strict loading/writing, artifact-derived status refresh, validation-shaped findings, and status service behavior. `internal/sprint` may import `internal/project` for project-root/catalog boundaries and `internal/workspace` for path safety; it must not import `internal/study` or make platform packages aware of sprint semantics.
- **Rationale:** The sprint requirements and architecture docs require product behavior to remain with the state it transforms. Sprint planning artifacts and `flow-state.json` are sprint-owned state, not study behavior and not generic platform behavior.
- **Study / Source Grounding:** `technical-handbook.md` relevant patterns on thin CLI and module-owned domain/persistence; `01-project-structure` with chezmoi `main.go:26-34`, gh-cli `cmd/gh/main.go:6`, and yq `cmd/root.go:9`; `15-philosophy` with restic `internal/backend/backend.go:19-90`; Architecture docs lines defining `internal/sprint` ownership and dependency direction.
- **Trade-Offs Accepted:** Accept some local sprint-specific helpers over premature shared abstractions. This keeps behavior understandable now and avoids study semantic leakage.
- **Technical Debt / Future Impact:** Future prompt/flow/validation stages can extend `internal/sprint`; if mechanical helpers repeat across modules, extract only product-neutral helpers later.
- **Alternatives Rejected:** Rejected `internal/workflow` or `internal/stages` because six planning stages do not justify a global workflow abstraction. Rejected reusing `internal/study` run-loop types because study task states, reports, sources, dimensions, ratings, and scheduling are explicitly out of scope.
- **Contracts Satisfied:** Architecture contract; Security contract; `requirements.md:30,49-52`; PRD Phase 2 planning scope; TRD 18.1 and 18.2.
- **Evidence Required:** Import review confirms `internal/sprint` has no `internal/study` dependency and platform packages import no product packages; package docs explain ownership and deferred behavior; tests exercise sprint behavior through public package/service surfaces.

### Decision 2: Model Exactly Six Planning Stages And Seven Planning Artifacts

- **Decision:** Define supported stages as exactly `requirements`, `sprint-index`, `technical-handbook`, `area-reasoning`, `reasoning`, and `plan`, in that order. Define supported planning artifacts as exactly `requirements.md`, `sprint-index.md`, `technical-handbook.md`, `reasoning/`, `reasoning.md`, `plan.md`, and `flow-state.json`. Strict loading rejects `implementation`, `execute`, `smoke`, `review`, `issues`, and any unknown stage.
- **Rationale:** The sprint goal is to model planning through `plan.md` only. Treating execution/review/smoke/issues as unsupported is necessary to prevent Phase 2 planning status from becoming a project-management or implementation workflow engine.
- **Study / Source Grounding:** `technical-handbook.md` design pressures and warnings against generic stage engines; `12-extensibility` evidence including restic `cmd/restic/main.go:77-106` supporting closed command/stage lists for bounded surfaces; `15-philosophy` simplicity guidance; PRD sections deferring sprint execution/smoke/review/issues.
- **Trade-Offs Accepted:** Closed enumeration is less flexible than registry-based stages, but it makes strict state validation, output ordering, and non-goal enforcement reviewable.
- **Technical Debt / Future Impact:** Future phases that add execution or smoke must define a new schema and explicit migration or separate state, not silently reuse this planning schema.
- **Alternatives Rejected:** Rejected stage plugin/registry design because there is no current extension requirement. Rejected accepting unknown stages with warnings because the requirements demand strict rejection and because unknown stages could smuggle deferred workflows into current status.
- **Contracts Satisfied:** Architecture, CLI Surface, Persistence And Migrations; `requirements.md:34-36,40,52-53`; PRD/TRD deferred-scope requirements.
- **Evidence Required:** Domain tests assert exact stage order and artifact paths; strict load tests reject unsupported `implementation` and `review`; review confirms no implementation/smoke/review/issues commands or modeled stages were added.

### Decision 3: Store Strict Versioned Flow State With Workspace-Relative Safe Paths

- **Decision:** Persist `flow-state.json` with schema version `1`, project name, sprint slug, updated timestamp, and one state entry per supported planning stage. Each `StageState` records stage, status, workspace-relative artifact path, optional last-run timestamp, and optional safe error detail. Allowed statuses are exactly `missing`, `ready`, `complete`, `failed`, and `skipped`. Loading rejects malformed JSON, unsupported schema versions, missing fields, unknown/duplicate stages, unknown statuses, path escapes, project mismatch, and sprint mismatch.
- **Rationale:** `flow-state.json` is durable machine state, not a best-effort cache. Workspace-relative paths keep artifacts reviewable and relocatable while path normalization and sprint-root containment prevent unsafe persisted paths.
- **Study / Source Grounding:** `technical-handbook.md` strict state and security patterns; `05-error-handling` with go-task `errors/errors.go:47-50` and restic `internal/errors/fatal.go:10`; `13-security` with k9s `internal/config/json/validator.go:146` and gdu path caution `internal/common/ignore.go:16-37`; `07-state-context` explicit durable state examples.
- **Trade-Offs Accepted:** Strict rejection may inconvenience users who hand-edit state incorrectly. The benefit is that corrupted or unsafe state cannot be silently accepted.
- **Technical Debt / Future Impact:** No backward-compatible loader is introduced. A future migration sprint must define concrete old shapes before compatibility code is added.
- **Alternatives Rejected:** Rejected sprint-root-relative persisted paths because requirements require workspace-relative artifact paths and project status output benefits from workspace-level consistency. Rejected absolute persisted paths because they are less reviewable and more host-specific. Rejected lenient unknown-field/stage/status acceptance because it undermines strict validation.
- **Contracts Satisfied:** Persistence And Migrations, Security, Errors, Observability; `requirements.md:37-42,47`; TRD 18.5.
- **Evidence Required:** Unit tests for all strict rejection cases, path escape attempts, project/sprint mismatches, unknown statuses, duplicate stages, unsupported schema versions, and safe error preservation; command tests assert invalid state exits `5` and reports the state path safely.

### Decision 4: Refresh Status From Artifacts Only After State Is Missing Or Valid

- **Decision:** The status service first resolves project and sprint boundaries, then strictly loads `flow-state.json` if it exists. If the file is malformed or invalid, status fails with exit code `5` and does not repair it. If state is missing or valid, the service inspects planning artifacts without invoking runtime behavior, derives current statuses, preserves recorded `failed` stages and safe errors, writes refreshed state atomically, and returns a status summary.
- **Rationale:** This order satisfies both strict-state and artifact-as-state requirements. It prevents corrupt state from being masked by refresh while still letting status reflect human-created or edited planning files when state is absent or valid but stale.
- **Study / Source Grounding:** `technical-handbook.md` trade-off on refresh from artifact inspection versus strict load; `07-state-context` explicit state roots; `05-error-handling` cause-preserving errors; `10-logging-observability` status/diagnostics discipline. No runtime study sources are relevant because this decision intentionally excludes runtime execution.
- **Trade-Offs Accepted:** Recorded failures may remain even when a file appears present and non-empty; this is intentional because requirements say recorded stage errors remain failed and no repair/validate command is in scope.
- **Technical Debt / Future Impact:** Future planning generation or validation commands need explicit rules for clearing `failed` states after successful rerun or validation.
- **Alternatives Rejected:** Rejected always deriving from artifacts before loading state because it could hide invalid JSON, unsafe paths, or unsupported stages. Rejected never writing refreshed state during status because requirements require status refresh through the command path and durable `flow-state.json` inspection.
- **Contracts Satisfied:** Persistence And Migrations, Observability, Errors, Security; `requirements.md:41-49,61-68`.
- **Evidence Required:** Tests for missing flow-state refresh, valid stale refresh, invalid state failing before refresh, failed state preservation, atomic write success/failure, and absence of runtime/study/Git/network invocations.

### Decision 5: Derive Stage Status Deterministically With Narrow Artifact Rules

- **Decision:** Artifact inspection checks only expected paths under the selected sprint root. Present non-empty files mark file-backed stages `complete`. Missing stages blocked by earlier incomplete stages are `missing`. The first missing stage after completed predecessors is `ready`. `area-reasoning` is `complete` when `reasoning/` contains at least one readable non-empty Markdown file, `skipped` only when a readable `sprint-index.md` explicitly selects no reasoning templates, and otherwise `missing` or `ready` according to stage order. Recorded failed stages remain `failed` with safe error detail.
- **Rationale:** This implements runtime-free planning status without content validation. It honors the special optional nature of area reasoning only when explicit no-template evidence exists, while avoiding silent skips when templates might be selected or `sprint-index.md` cannot be read.
- **Study / Source Grounding:** `technical-handbook.md` open questions on `area-reasoning` skip evidence and deterministic artifact status; `13-security` path validation references; `11-testing-strategy` behavior-focused fixture tests. No external repo source decides the exact `area-reasoning` semantics because this is project-specific product behavior derived from requirements.
- **Trade-Offs Accepted:** The no-template detector is intentionally narrow and not a full `sprint-index.md` validator. This avoids implementing Sprint 18 validation scope while still enforcing the Sprint 17 skip rule.
- **Technical Debt / Future Impact:** Full mapping of selected reasoning templates to required `reasoning/*.md` files is deferred to later stage validation. Current status can say area reasoning is complete if at least one non-empty reasoning file exists, even if a later validator finds the selection incomplete.
- **Alternatives Rejected:** Rejected always skipping `area-reasoning` when `reasoning/` is absent because requirements prohibit silent skips. Rejected validating `sprint-index.md` selections against `project-index.md` or exact reasoning file counts because that is a non-goal for this sprint and belongs to later validation.
- **Contracts Satisfied:** Observability, Testing, Security; `requirements.md:34,36,43-45,49,53,61-66`.
- **Evidence Required:** Fixture tests for no artifacts, partial artifacts, all artifacts, missing `reasoning/`, selected reasoning templates, explicit no-template selection, unreadable or missing `sprint-index.md`, deterministic reasoning file sorting, and failure preservation.

### Decision 6: Add Only Thin CLI Wiring For `sprint status`

- **Decision:** Add `ultraplan sprint <project> <sprint> status` and command help. `internal/app` parses arguments, resolves workspace/project dependencies through existing app composition, calls the sprint status service, renders deterministic plain text, and maps exit codes. Business rules, stage ordering, path rules, state validation, and artifact status derivation stay in `internal/sprint`.
- **Rationale:** This matches the selected CLI surface while preserving module boundaries. The command is an adapter, not the source of sprint behavior.
- **Study / Source Grounding:** `technical-handbook.md` thin CLI pattern; `02-command-architecture` with gh-cli `pkg/cmdutil/factory.go:16-43`, helm `pkg/cmd/install.go:132-145`, and restic `cmd/restic/main.go:37-114`; `06-io-abstraction` with gh-cli IOStreams; `09-terminal-ux` plain CLI-first output evidence.
- **Trade-Offs Accepted:** Explicit command wiring grows linearly as sprint commands are added, but it is clearer and safer than a registry for this small command family.
- **Technical Debt / Future Impact:** Future `validate`, `prompt`, and `flow` commands can reuse the sprint service/domain without moving behavior into app handlers.
- **Alternatives Rejected:** Rejected dynamic command registration or plugin-style stage commands because selected evidence warns against unnecessary extensibility and subprocess/plugin trust boundaries. Rejected JSON output in this sprint because requirements exclude it unless existing app infrastructure supports it without scope expansion.
- **Contracts Satisfied:** CLI Surface, Documentation, Observability, Errors; `requirements.md:21-23,45-48,54`.
- **Evidence Required:** Command tests for top-level help, `sprint --help`, `status --help`, malformed arguments, invalid/missing/ambiguous sprint refs, invalid flow state exit `5`, deterministic stdout, stderr diagnostics, no ANSI, and no runtime invocation.

### Decision 7: Keep Persistence Atomic And Failure-Loud

- **Decision:** Flow-state writes use a temp file in the same directory, JSON encoding, flush and close, rename over the prior file, and best-effort parent directory sync. Write or rename failures preserve the last valid `flow-state.json`, return a non-zero command result, and preserve the flow-state path and underlying cause internally.
- **Rationale:** Atomic state writes are required to prevent partial writes from corrupting durable state. The status command is mutating only in the narrow sense of refreshing `flow-state.json`, so write failure must be visible rather than silently ignored.
- **Study / Source Grounding:** `technical-handbook.md` persistence trade-offs; Sprint 12 precedent named in requirements; TRD 13.3 and 21.1 atomic write requirements; `05-error-handling` cause chains; `06-io-abstraction` filesystem seam evidence for testing write failures.
- **Trade-Offs Accepted:** Same-directory temp writes and parent sync add implementation complexity versus simple overwrite, but they protect durable state and match requirements.
- **Technical Debt / Future Impact:** If a generic atomic-write helper already exists and is product-neutral, reuse it; otherwise keep the implementation local until reuse is proven.
- **Alternatives Rejected:** Rejected direct overwrite because partial writes could corrupt the only state file. Rejected ignoring parent sync entirely because the requirements call for best-effort directory sync.
- **Contracts Satisfied:** Persistence And Migrations, Errors, Testing; `requirements.md:41-42,46,55-58`.
- **Evidence Required:** Unit tests inject temp write, flush/close, rename, and parent sync failures where feasible; tests assert previous valid content remains after failure and errors preserve state path and cause.

### Decision 8: Verify With Offline Behavior Tests And Review Checks

- **Decision:** Implement tests in `internal/sprint/sprint_test.go` and `internal/app/sprint_commands_test.go` covering domain order, discovery, reference resolution, path safety, status derivation, strict loading, atomic writes, unsupported stages/statuses, command rendering, exit codes, help, stdout/stderr separation, deterministic ordering, and runtime-free behavior. Final verification must run `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`.
- **Rationale:** The sprint's risk is not algorithmic complexity; it is boundary drift, unsafe paths, state corruption, and accidental scope expansion. Behavior tests and import/review checks are the right evidence.
- **Study / Source Grounding:** `technical-handbook.md` behavior-focused and golden-compatible test pattern; `11-testing-strategy` with chezmoi `internal/cmd/main_test.go:64-174`, helm `internal/test/test.go:43`, and restic `internal/backend/mock/backend.go:14-26`; `10-logging-observability` stdout/stderr separation evidence.
- **Trade-Offs Accepted:** Some tests may use fixture text/golden-style assertions, which require deliberate updates when output changes. The benefit is stable, reviewable CLI behavior.
- **Technical Debt / Future Impact:** Future JSON output or validators should add structured tests rather than weakening current text output assertions.
- **Alternatives Rejected:** Rejected private-helper-heavy tests because selected evidence warns that behavior tests age better. Rejected gated real-runtime smoke because the sprint is explicitly runtime-free and the smoke harness was excluded by the sprint index.
- **Contracts Satisfied:** Testing, Observability, CLI Surface, Architecture; `requirements.md:24-25,53-58,104-114`.
- **Evidence Required:** Passing test/build commands, review checklist confirming non-goals, import review, and `review.md` recording verification and deviations after implementation.

## Expected Evidence

| Evidence Type | Required Evidence | Source / Command / Review Check |
| --- | --- | --- |
| Tests | Sprint domain tests cover stage order, exact supported statuses, artifact paths, discovery, hidden/non-directory ignores, invalid slugs, exact/prefix/ambiguous references, status derivation, skipped area reasoning, strict state loading, path escapes, project/sprint mismatch, atomic writes, and write failure preservation. | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/sprint_test.go`; `go test ./...` |
| Tests | Command tests cover help, argument errors, status output, stdout/stderr separation, exit codes, deterministic ordering, invalid flow-state diagnostics, invalid/ambiguous sprint refs, no ANSI, and no runtime invocation. | `/home/antonioborgerees/coding/ultraplan-go/internal/app/sprint_commands_test.go`; `go test ./...` |
| Build | CLI package builds successfully. | `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go` |
| Race | Race test suite passes. | `go test -race ./...` from `/home/antonioborgerees/coding/ultraplan-go` |
| Runtime | No runtime evidence is required because this sprint is explicitly runtime-free; evidence must instead show no calls to agentwrap/OpenCode/network/shell/study run-loop/project-index mutation during status. | Import review, command tests with fake/no-op runtime absence, review checklist. |
| Review | Architecture review confirms `internal/sprint` ownership, allowed imports, product/platform separation, absence of global workflow abstractions, strict planning-stage-only scope, and no deferred workflow modeling. | `system/protocols/architecture-review-protocol.md`; review notes in sprint `review.md` after implementation. |
| Review | Sprint review confirms required outputs, acceptance criteria, verification results, deviations, and non-goals. | `system/protocols/sprint-review-protocol.md`; `projects/ultraplan-go/sprints/17-artifact-domain-and-flow-state/review.md` |
| Documentation | Package docs and command help describe sprint ownership, planning-stage scope through `plan.md`, runtime-free status, and deferred execution/smoke/review/issues/Git behavior. | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/doc.go`; command help tests. |

## Assumptions And Risks

| Item | Type | Impact | Mitigation / Follow-Up |
| --- | --- | --- | --- |
| Existing Sprint 16 project APIs can resolve project roots under the workspace without duplicating catalog ownership. | Assumption | If unavailable or incompatible, sprint code may be tempted to duplicate project discovery. | Use only exported project boundaries; if missing, add the smallest project-facing method needed without copying project semantics. |
| Existing workspace path helpers can normalize and verify containment for sprint artifacts. | Assumption | If inadequate, unsafe path validation could be inconsistent. | Add or reuse product-neutral workspace path helpers; keep sprint-specific artifact semantics in `internal/sprint`. |
| Missing selected Architecture area reasoning file is acceptable because this document directly resolves architecture decisions. | Risk | Review may require backfilling `reasoning/architecture.md` before plan execution. | Record the gap here; sprint review or planning owner can request a backfill if process compliance requires it. |
| Narrow no-template detection for `area-reasoning` may not understand every valid future `sprint-index.md` shape. | Risk | Status may mark area reasoning `ready` or `missing` instead of `skipped` for uncommon explicit no-template wording. | Keep detection conservative; only skip on clear explicit no-template selection; full validation belongs to Sprint 18. |
| Persisted failed stage preservation can surprise users after they create the missing artifact manually. | Risk | Status may continue showing `failed` until state is cleared by a future command or manual state edit. | Document behavior in package docs/status diagnostics; future flow/validation commands define failure clearing. |
| Atomic parent directory sync behavior can vary by platform/filesystem. | Risk | Best-effort sync may be untestable or unsupported in some environments. | Treat parent sync failure as best-effort warning only if file rename succeeded; preserve core atomic write behavior. |
| The command mutates `flow-state.json` during status refresh. | Risk | A user may expect a read-only status command. | Requirements explicitly require refresh from command path; keep mutation limited to `flow-state.json`, atomic, deterministic, and documented. |

## Implementation Constraints

- `internal/sprint` must not import `internal/study`, reuse study source/dimension/report/run-loop models, or depend on study prompt/validation/rating/summary behavior.
- `internal/app` must not contain sprint stage ordering, artifact path rules, flow-state schema validation, status derivation, or atomic-write logic.
- Platform packages under `internal/platform/*` must remain product-agnostic and must not import `internal/sprint`, `internal/project`, `internal/study`, or `internal/app`.
- Do not add `internal/planning`, `internal/workflow`, `internal/stages`, `internal/validation`, `internal/reports`, `internal/prompts`, or another speculative global workflow package.
- Persisted stage artifact paths must be workspace-relative and must normalize to paths contained by the selected sprint root before filesystem access or state acceptance.
- Strict flow-state load failures must not be repaired or overwritten by artifact refresh in the same status command.
- Existing recorded `failed` stage state and safe error detail must remain failed during artifact-derived refresh.
- `area-reasoning` may be marked `skipped` only when readable `sprint-index.md` explicitly selects no reasoning templates; otherwise absent reasoning files are not skipped silently.
- Status inspection must not invoke agentwrap, OpenCode, runtime health checks, prompt generation, network calls, source cloning, shell execution, Git commands, project-index mutation, planning-stage generation, or study run-loop behavior.
- Status output must be deterministic plain text with no ANSI escape sequences and clear stdout/stderr separation.
- Normal test runs must be offline, deterministic, and not require OpenCode, provider credentials, network access, Git commands, subprocess execution, or long sleeps.
- Do not add compatibility code for unsupported pre-sprint `flow-state.json` shapes unless a future migration requirement explicitly defines those shapes.
- Do not modify unrelated dirty files or revert user changes in the target repository.

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
- [x] Area-specific reasoning documents were completed where applicable or explicitly marked not available for this sprint.
- [x] Area-specific reasoning conclusions are reflected or explicitly overridden in final decisions; no conclusions existed, so this document directly resolves the selected Architecture questions.
- [x] Contracts are explicitly mapped to decisions and expected evidence.
- [x] Final decisions are clear enough for `plan.md` to execute.
- [x] Expected evidence is specific and reviewable.
