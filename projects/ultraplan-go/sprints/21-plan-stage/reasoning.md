# Sprint Reasoning: Plan Stage

> Project: `ultraplan-go`
> Sprint: `21-plan-stage`
> Output: `projects/ultraplan-go/sprints/21-plan-stage/reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/sprints/21-plan-stage/requirements.md`, `projects/ultraplan-go/sprints/21-plan-stage/sprint-index.md`, `projects/ultraplan-go/sprints/21-plan-stage/technical-handbook.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/project-index.md`, `templates/sprint-reasoning.md`
> **Area Reasoning Inputs:** none present; `projects/ultraplan-go/sprints/21-plan-stage/reasoning/*.md` does not exist in the workspace.

This document decides. It synthesizes selected context, handbook evidence, available area-specific reasoning, and selected contracts into final sprint decisions for implementing the plan stage.

It does not replace `sprint-index.md`, `technical-handbook.md`, or `reasoning/*.md`.

## Sprint Purpose

- **Goal:** Implement the planning `plan` stage so validated `reasoning.md` can be previewed, flowed, generated, and validated into `plan.md`, with `flow-state.json` updated through `plan` and with no implementation, smoke, review, issue, or Git behavior.
- **Non-Goals:** The sprint will not execute implementation tasks, run smoke investigations, generate or validate product conformance reviews, create or manage issues, mutate Git state, introduce `.run-state.json`, `smoke.md`, `smoke.json`, product-generated `review.md`, `issues.md`, or `issues.json`, alter select/distill/reason contracts beyond consuming valid prerequisites, reuse study semantics, or add global workflow/prompt/validation/planning packages.
- **Depends On:** Sprint 20 reason-stage outputs, Sprint 19 technical-handbook behavior, Sprint 18 selected-context behavior, Sprint 17 planning flow-state and atomic write behavior, Sprint 16 project catalog boundaries, Sprint 14 diagnostics and exit-code discipline, Sprint 9 generic runtime boundary, and the selected project docs and study reports listed in `sprint-index.md`.

## Selected Context And Pre-Reasoning Artifacts

| Artifact | Path | How It Was Used |
| --- | --- | --- |
| Sprint Requirements | `projects/ultraplan-go/sprints/21-plan-stage/requirements.md` | Authoritative sprint contract for scope, outputs, acceptance criteria, non-goals, constraints, dependencies, and review expectations. Requirement IDs in this document use `AC-01` through `AC-33` for the acceptance criteria in order, `NG-01` onward for non-goals, and `C-01` onward for constraints. |
| Project Index | `projects/ultraplan-go/project-index.md` | Confirmed selected contracts, evidence reports, the Architecture reasoning template, project docs, and review protocols are available in the project catalog. Also confirmed Phase 2 stops at `plan.md` and planning modules may reuse infrastructure but not study semantics. |
| PRD | `projects/ultraplan-go/docs/PRD.md` | Established product scope: governed project/sprint planning through `plan.md`, editable artifacts, runtime success not being product success, calm CLI output, deterministic artifacts, and deferral of execution, smoke, review automation, issues, and Git mutation. |
| TRD | `projects/ultraplan-go/docs/TRD.md` | Established technical scope: `internal/sprint` ownership, planning stages through `plan`, `flow-state.json` schema/statuses, plan validator requirements, prompt rendering inputs, runtime boundary, command shapes, exit codes, path safety, atomic writes, diagnostics, and testing expectations. |
| Architecture Doc | `projects/ultraplan-go/docs/ARCHITECTURE.md` | Established module-driven layout, allowed dependencies, `sprint -> project/workspace/platform/runtime`, platform product-agnostic boundaries, reuse of infrastructure only, and prohibition on current implementation/smoke/review/issue/Git workflow stages. |
| Sprint Index | `projects/ultraplan-go/sprints/21-plan-stage/sprint-index.md` | Selected Architecture, Errors, Configuration, Observability, Security, Testing, Documentation, CLI Surface, LLM Runtime, LLM Evaluation / Cost / Safety, Workflows, and Persistence And Migrations contracts; selected 12 Go CLI study reports; selected Architecture area reasoning; selected Architecture Review and Sprint Review protocols; excluded concurrency, extensibility, performance, smoke, deferred stages, unselected context, and study service reuse. |
| Technical Handbook | `projects/ultraplan-go/sprints/21-plan-stage/technical-handbook.md` | Provided the source-evidence basis for thin CLI, module-owned behavior, constructor/runtime/store seams, writer injection, validation as completion gate, context-aware flow, safe diagnostics, fake-first tests, and deliberate non-goals. |
| Sprint Reasoning Template | `templates/sprint-reasoning.md` | Supplied the structure for final decisions, trade-off and debt analysis, source evidence, assumptions, risks, implementation constraints, and plan handoff. |

## Area-Specific Reasoning Inputs

| Area | Reasoning Document | Key Conclusion | Evidence Basis | Impact On Final Decision |
| --- | --- | --- | --- | --- |
| Architecture | `projects/ultraplan-go/sprints/21-plan-stage/reasoning/architecture.md` | No actual area-specific reasoning document exists in the workspace at the time this final reasoning was authored. The sprint index selects this artifact, so implementation and validation must treat the selected area reasoning file as a prerequisite for runtime-backed plan flow. | No area document evidence was available. Final reasoning uses `requirements.md`, `sprint-index.md`, `technical-handbook.md`, PRD, TRD, and Architecture docs instead of inventing an area-specific conclusion. | The final decisions below settle the architecture directly enough for `plan.md` to execute, and they preserve a risk item requiring the missing selected area reasoning artifact to be produced or the sprint index to be revised before the implemented flow can pass its own prerequisite checks. |

## Sprint Technical Handbook Summary

- **Relevant Patterns:** Thin command-to-domain delegation, module-owned product behavior, explicit constructor/runtime/store seams, IO/writer injection for command tests, validation as the completion gate, context-aware flow, safe diagnostics, fake-first tests, and deliberate non-goal discipline.
- **Important Trade-Offs:** Keep plan behavior in `internal/sprint` instead of global packages; use explicit service/store/runtime injection rather than direct filesystem/runtime calls; make validation strict enough to gate success without becoming a prose linter; combine focused fixture/substr assertions with limited golden tests; keep `flow-state.json` explicit and atomic; use the generic runtime boundary rather than direct OpenCode or shell execution.
- **Warnings / Anti-Patterns:** Do not let command handlers accumulate plan rules, do not hide state behind globals, do not bypass output and filesystem seams, do not accept contexts and ignore them, do not add shell/Git/provider/OpenCode calls to `internal/sprint`, do not leak secrets in diagnostics, and do not assert brittle prompt prose unrelated to required behavior.
- **Evidence Confidence:** High for package/CLI/test/runtime-boundary decisions because multiple selected reports and concrete repos support them; medium for terminal UX, context, config, and philosophy decisions because the findings are consistent but less directly plan-stage-specific.

## Contracts Applied

| Contract / Requirement ID | Constraint | Decision Impact | Expected Evidence |
| --- | --- | --- | --- |
| Architecture, C-01, C-05, C-06, C-12, AC-18, AC-19, AC-20, AC-21 | Product behavior stays module-owned; `internal/sprint` owns plan behavior; `internal/app` is thin; platform packages stay product-agnostic; study semantics are not reused. | Plan validator, prompt renderer, prerequisite loading, task/evidence trace checks, flow-state transitions, and service use cases will live in `internal/sprint`; command wiring only parses, calls services, renders output, and maps exit codes. | Import review, Architecture Review protocol, `go test ./...`, and command tests proving behavior remains behind sprint services. |
| Errors, Observability, CLI Surface, AC-22, AC-23, AC-24, AC-25, AC-26 | Diagnostics are deterministic, scriptable, separated between stdout/stderr, redacted, and mapped to exit codes. | Plan validation, prompt, dry-run, flow, runtime failure, invalid prerequisite, invalid plan, and unsupported-stage paths will return the required exit codes and safe messages. | Command tests in `internal/app/sprint_commands_test.go`; validation tests asserting stable findings and no ANSI output. |
| Security, Configuration, LLM Runtime, LLM Evaluation / Cost / Safety, C-04, C-07, C-08, C-09, AC-05, AC-10, AC-11, AC-12 | Runtime execution goes only through generic platform runtime; prompts use safe selected paths and no-mutation rules; diagnostics redact secrets and avoid unsafe raw payloads. | `flow --to plan` will build a generic platform runtime request; `internal/sprint` will not call OpenCode, agentwrap adapters, shell commands, provider APIs, or Git directly. | Fake-runtime tests; import/code review; prompt preview tests asserting selected paths, no-mutation instructions, and no absolute paths except local diagnostics when justified. |
| Workflows, Persistence And Migrations, AC-03, AC-04, AC-06, AC-13, AC-14, AC-15, AC-16, AC-17 | Planning flow supports only `requirements`, `sprint-index`, `technical-handbook`, `area-reasoning`, `reasoning`, and `plan`; `flow-state.json` is strict, versioned, and atomically written. | Flow will validate prerequisites before runtime, complete only after generated `plan.md` validates and state writes, mark failures safely, and reject unsupported future stages and artifacts. | Flow tests for dry-run, success, missing/invalid generated artifact, runtime failure, invalid prerequisites, validation failure, write failure, unsupported stages, and state contents. |
| Testing, Documentation, AC-27, AC-28, AC-29, AC-30, AC-31, AC-32, AC-33 | Verification must be deterministic, offline, fake-first, and reviewable; help and artifact docs must describe behavior. | Implement focused unit and command tests, fake runtime seams, prompt invariant checks, plan validation fixtures, and final review evidence. | `go test ./...`, `go test -race ./...`, `go build ./cmd/ultraplan`, help-output tests, and `projects/ultraplan-go/sprints/21-plan-stage/review.md`. |
| PRD/TRD Phase 2 Scope, NG-01 through NG-11, C-02, C-03, C-10, C-11 | Phase 2 ends at `plan.md`; `plan.md` is editable Markdown and traces to `reasoning.md`; no deferred implementation/smoke/review/issue/Git behavior is modeled. | Plan validation and prompt rendering will reject/instruct against future-stage execution semantics, mutation, and unselected authoritative context. | Plan validation tests for forbidden content; prompt tests for no-execution/no-mutation rules; review checks for absent deferred stage models. |

## Repos Studied / Source Evidence Used

| Source / Repo / Report | Concrete Reference | Relevant Finding | Why It Matters For This Sprint | Used In Decision(s) |
| --- | --- | --- | --- | --- |
| `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md`; k9s `cmd/root.go:112-127`; restic `cmd/restic/main.go:37-114`; yq `cmd/root.go:9` | Mature Go CLIs keep entrypoints thin and dependencies one-way into internal/domain packages. | Supports keeping plan-stage behavior in `internal/sprint` and `internal/app` as command wiring only. | Decisions 1, 5 |
| `02-command-architecture` | `studies/go-cli-study/reports/final/02-command-architecture.md`; helm `pkg/cmd/install.go:132-145`; gh-cli `pkg/cmdutil/factory.go:16-43`; chezmoi `internal/cmd/config.go:1833-1845` | Command factories and thin `RunE` delegates scale better than handlers containing business rules. | Supports thin `validate plan`, `prompt plan`, and `flow --to plan` CLI wiring. | Decisions 1, 5 |
| `03-dependency-injection` | `studies/go-cli-study/reports/final/03-dependency-injection.md`; gh-cli `pkg/cmd/factory/default.go:26-46`; go-task `executor.go:22-24`; restic `internal/backend/backend.go:19-90` | Explicit composition roots, constructor injection, and narrow seams improve testability. | Supports fake store/runtime/clock seams without creating a generic DI framework. | Decisions 1, 4, 5 |
| `04-configuration-management` | `studies/go-cli-study/reports/final/04-configuration-management.md`; chezmoi `internal/cmd/config.go:2253-2287`; go-task `internal/flags/flags.go:314-327`; restic `internal/global/global.go:139,147` | Clear config precedence and post-merge validation prevent surprising behavior. | Supports using resolved runtime config and redacted diagnostics for plan flow rather than ad hoc environment access. | Decisions 3, 4 |
| `05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md`; go-task `errors/errors.go:47-50`; gh-cli `internal/ghcmd/cmd.go:44-49`; restic `internal/errors/fatal.go:10` | First-class error values, wrapping, sentinels, and exit-code mapping support scriptable diagnostics. | Supports deterministic validation findings and required exit-code mapping. | Decisions 2, 4, 5 |
| `06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md`; gh-cli `pkg/iostreams/iostreams.go:551-568`; go-task `executor.go:553-564`; restic `internal/ui/terminal.go:10-36` | Injectable stdout/stderr/filesystem/runtime seams make command behavior testable. | Supports command tests for no runtime invocation, stdout/stderr separation, prompt previews, and no ANSI output. | Decisions 3, 5, 6 |
| `07-state-context` | `studies/go-cli-study/reports/final/07-state-context.md`; gh-cli `internal/ghcmd/cmd.go:142`; helm `pkg/cmd/install.go:333-347`; restic `cmd/restic/cleanup.go:24-38` | Context propagation and centralized explicit state matter for cancellable operations. | Supports context-aware non-dry-run plan flow and explicit stage failure/complete state. | Decision 4 |
| `09-terminal-ux` | `studies/go-cli-study/reports/final/09-terminal-ux.md`; chezmoi `internal/cmd/prompt.go:124-137`; rclone `lib/terminal/terminal.go:82-86` | Scriptable CLI commands need calm output, non-TTY awareness, and no unnecessary TUI behavior. | Supports stable text output, no ANSI by default, and prompt/validation output suitable for scripts. | Decisions 3, 5 |
| `10-logging-observability` | `studies/go-cli-study/reports/final/10-logging-observability.md`; helm `internal/logging/logging.go:31-66`; k9s `internal/slogs/keys.go:6-231`; gh-cli `pkg/iostreams/iostreams.go:52-54` | Structured diagnostics and stdout/stderr separation are operationally important. | Supports safe runtime failure summaries, validation findings, and flow-state failure visibility. | Decisions 4, 5 |
| `11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md`; chezmoi `internal/cmd/main_test.go:64-174`; helm `internal/test/test.go:43`; rclone `cmd/bisync/bisync_test.go:1435-1479` | Mature CLIs combine table tests, command integration tests, fixtures, goldens, and fakes. | Supports plan validation fixtures, prompt invariant tests, fake-runtime flow tests, and command tests. | Decision 6 |
| `13-security` | `studies/go-cli-study/reports/final/13-security.md`; restic `internal/options/secret_string.go:15-20`; helm `pkg/registry/transport.go:37-41`; opencode `internal/llm/tools/bash.go:41-55` | Trust boundaries, redaction, path validation, and constrained execution reduce risk. | Supports workspace-safe artifact access, redacted diagnostics, no shell/Git/provider calls from `internal/sprint`, and prompt no-mutation rules. | Decisions 3, 4, 5 |
| `15-philosophy` | `studies/go-cli-study/reports/final/15-philosophy.md`; age `age.go:18`; lazygit `VISION.md:97`; go-task `website/src/docs/experiments/index.md:17-21` | Strong projects keep non-goals visible and accept complexity deliberately. | Supports stopping Phase 2 at `plan`, rejecting deferred stages, and avoiding speculative abstractions. | Decisions 1, 2, 4 |

## Trade-Off And Debt Analysis

### Accepted Trade-Offs

| Trade-Off | Benefit | Cost / Constraint Accepted | Why Acceptable Now | Revisit Trigger |
| --- | --- | --- | --- | --- |
| Keep plan behavior inside `internal/sprint` instead of shared planning/prompt/validation/workflow packages | Preserves module ownership and avoids premature cross-module coupling. | Some Markdown validation and prompt code may resemble study code without reuse. | The sprint has one concrete planning module owner and explicit requirements forbidding study semantic reuse. | A later sprint proves an identical mechanical helper is needed by multiple modules and can be extracted without product semantics. |
| Validate `plan.md` structurally rather than semantically proving implementation quality | Prevents invalid success and catches missing traceability, sections, placeholders, and forbidden future-stage behavior. | Generated prose may pass validation while still needing human review for quality. | Phase 2 requires editable artifacts and reviewable evidence, not semantic correctness guarantees. | Repeated invalid plans pass validation because checks are too shallow. |
| Use prompt invariant tests instead of full exact-prose tests everywhere | Keeps tests stable while proving required variables, selected inputs, paths, traceability, and no-mutation rules. | Some harmless wording regressions will not be caught. | Requirements care about deterministic required content more than exact prose. | Prompt wording becomes part of a public compatibility contract or regressions slip through invariant checks. |
| Require post-runtime artifact validation before marking `plan` complete | Prevents runtime success from bypassing product success criteria. | Adds failure modes where runtime succeeds but flow fails due missing/invalid artifact or state write failure. | PRD/TRD and Sprint 21 acceptance criteria make validation the completion gate. | None for this sprint; this is a core product invariant. |
| Keep runtime-backed plan flow sequential through `plan` only | Simpler state transitions, clearer diagnostics, and no scheduler complexity. | Does not prepare a generic execution DAG for later phases. | The selected sprint excludes concurrency evidence, execution, smoke, review, issues, and Git mutation. | A future requirements revision adds implementation execution or multi-stage run loops. |
| Treat the selected but missing Architecture area reasoning file as an input gap rather than inventing its conclusions | Keeps the document truthful and avoids fabricating evidence. | The flow implementation must still require selected area reasoning before plan runtime execution, so current artifacts are not fully prereq-complete. | The user requested final reasoning now and area files are optional inputs only if present. | The sprint index continues to select Architecture reasoning and validation requires the file; create the area reasoning artifact or revise selection before executing plan flow. |

### Potential Technical Debt

| Debt / Shortcut | Why It Might Accrue | Current Mitigation | Owner / Follow-Up |
| --- | --- | --- | --- |
| Local Markdown section checks in `internal/sprint` | Plan, reasoning, and handbook validators may grow repeated heading/placeholder helpers. | Keep helpers private and focused; extract only mechanical helpers after repeated stable use. | Future planning validation maintenance sprint if duplication becomes costly. |
| Large `sprint.Service` surface | Adding validate, prompt, flow, status, and prerequisite use cases can make one service too broad. | Use request/result types and focused files while keeping one module-owned service boundary. | Revisit after plan stage if service methods become hard to test or reason about. |
| Prompt rendering text drift | Invariant tests may not catch all accidental wording shifts. | Assert required paths, sections, traceability rules, no-mutation rules, no-execution rules, and deterministic ordering. | Add golden fixtures for stable blocks if prompt regressions occur. |
| Flow-state schema growth | Adding plan failure details and skipped area reasoning can expand state complexity. | Keep supported stages fixed, reject unknown future stages, and preserve atomic strict loading/writing. | Revisit only when a future phase intentionally adds new planning stages. |
| Missing selected area reasoning artifact | The final reasoning can proceed, but implemented flow should fail prerequisites until the selected artifact exists. | Record the gap as a risk and require validation tests for selected area reasoning prerequisites. | Sprint artifact authoring before runtime-backed plan flow is used. |

### Future Considerations

| Consideration | Deferred Until | Reason Deferred | What Should Be Preserved Now |
| --- | --- | --- | --- |
| Sprint implementation execution from `plan.md` | Future requirements revision after Planning Phase 2 | Current PRD, TRD, and requirements stop at `plan.md`. | Reject `execute`/`implementation` stages and avoid modeling `.run-state.json`. |
| Smoke investigation automation | Future smoke-specific sprint | Smoke artifacts and real-runtime smoke harness are excluded from default plan-stage verification. | Keep `smoke.md` and `smoke.json` out of current validators and flow state. |
| Product-generated conformance review automation | Future review-specific sprint | Current review evidence is a human/post-implementation artifact, not a product-generated planning stage. | Do not require or generate product `review.md` in Phase 2 flow. |
| Issue tracking and Git mutation | Future explicitly scoped product phase | PRD and sprint requirements defer issues and Git mutation. | Reject issue/Git stages and keep prompts from instructing mutation. |
| Shared validation/prompt packages | Only after repeated stable use across modules | Global packages would violate selected architecture unless proven mechanical and product-neutral. | Keep plan logic local and keep helper extraction narrow. |
| Real OpenCode smoke evidence for plan flow | Optional deep-smoke review request | Default verification is offline and fake-runtime based. | Keep runtime seam compatible with generic platform runtime and fake-first tests. |

## Final Decisions

### Decision 1: Own Plan Behavior In `internal/sprint`

- **Decision:** Implement plan-stage domain, prerequisite loading, plan validation, plan prompt rendering, task/evidence trace checks, and flow-state transitions in `internal/sprint`, using focused files such as `plan.go`, `prompts.go`, `flow.go`, `state.go`, `service.go`, and `store_fs.go`; keep `internal/app/sprint_commands.go` as thin CLI wiring.
- **Rationale:** The PRD, TRD, Architecture doc, requirements, and sprint index all make sprint planning a product-owned module through `plan.md`. Keeping plan behavior local preserves the module-driven architecture and avoids premature global abstractions or study-semantic reuse.
- **Study / Source Grounding:** `technical-handbook.md` cites `01-project-structure`, k9s `cmd/root.go:112-127`, restic `cmd/restic/main.go:37-114`, yq `cmd/root.go:9`, `02-command-architecture`, helm `pkg/cmd/install.go:132-145`, gh-cli `pkg/cmdutil/factory.go:16-43`, and `15-philosophy` for thin entrypoints, one-way dependencies, and non-goal discipline.
- **Trade-Offs Accepted:** Some local plan validation and prompt helpers may resemble study helpers, but this avoids extracting unstable shared planning abstractions.
- **Technical Debt / Future Impact:** A larger `internal/sprint` module may need file-level organization, but no subpackage or global package should be introduced until there is a concrete readability or dependency benefit.
- **Alternatives Rejected:** A global `internal/planning`, `internal/workflow`, `internal/validation`, or `internal/prompts` package is rejected because it violates selected architecture and requirements. Putting plan rules in CLI handlers is rejected because it makes command tests brittle and hides business behavior outside the owning module. Reusing `internal/study` services or models is rejected because study semantics are explicitly out of scope for planning.
- **Contracts Satisfied:** Architecture, Workflows, Testing, Documentation, C-01, C-05, C-06, C-12, AC-18, AC-19, AC-20, AC-21.
- **Evidence Required:** Import review proving `internal/sprint` does not import `internal/study`, platform packages do not import product modules, and no global planning/workflow/validation/prompt package was added; command tests proving CLI delegates to sprint services; `go test ./...`.

### Decision 2: Validate `plan.md` As Traceable Editable Markdown

- **Decision:** Add a plan validator that fails when `plan.md` is missing, empty, has placeholders, omits a `reasoning.md` citation or explicit reference, omits decisions to execute, omits a task checklist, omits an evidence checklist, omits risks/blockers, omits success criteria, contains forbidden deferred-stage behavior, or has implementation tasks without clear traceability to reasoning decisions and verification evidence.
- **Rationale:** `plan.md` is the executable planning artifact for humans and future implementation work, but runtime success cannot be trusted unless the artifact exists and passes structural and traceability checks. The validator must be strict about required evidence and deferred-stage prohibitions without becoming a brittle prose linter.
- **Study / Source Grounding:** `technical-handbook.md` cites `05-error-handling`, gh-cli `internal/ghcmd/cmd.go:44-49`, go-task `errors/errors.go:47-50`, restic `internal/errors/fatal.go:10`, and `11-testing-strategy` for deterministic validation diagnostics and behavior-focused tests. `15-philosophy` supports keeping deferred behavior visible and out of scope.
- **Trade-Offs Accepted:** Structural checks may not prove that a plan is strategically excellent, but they prevent known invalid successes and keep the artifact editable.
- **Technical Debt / Future Impact:** If plan prose evolves, validators should continue checking stable headings, references, checklists, and forbidden behavior rather than exact prose. Repeated false positives should drive rule refinement, not validator removal.
- **Alternatives Rejected:** Accepting any non-empty Markdown is rejected because it would allow invalid success. Requiring exact template text is rejected because it would make human editing fragile. Semantically executing or simulating tasks from `plan.md` is rejected because implementation execution is out of scope.
- **Contracts Satisfied:** Errors, Testing, Workflows, Documentation, PRD/TRD planning validation, AC-01, AC-07, AC-08, AC-09, AC-25, AC-27, C-02, C-03.
- **Evidence Required:** `internal/sprint/plan_test.go` cases for valid content, missing file, empty file, placeholders, missing `reasoning.md` reference, missing decisions, missing tasks, missing evidence checklist, missing risks/blockers, missing success criteria, untraced tasks, and forbidden deferred-stage content.

### Decision 3: Render Deterministic Plan Prompts Without Runtime Or Mutation

- **Decision:** Implement `prompt plan` so it renders a deterministic preview containing project slug, sprint slug, sprint path, requirements path, sprint-index path, technical-handbook path, selected area reasoning paths, final `reasoning.md` path, `plan.md` output path, required sections, traceability rules, and no-execution/no-mutation rules. It writes to stdout by default; if the existing command surface supports an explicit preview-output path, only that path may be written.
- **Rationale:** Users need to inspect the exact plan prompt before runtime execution. The prompt must guide the runtime to create editable Markdown only at `plan.md` and not mutate prerequisite artifacts, docs, source repositories, workspace config, implementation files, or Git state.
- **Study / Source Grounding:** `technical-handbook.md` cites `06-io-abstraction`, gh-cli `pkg/iostreams/iostreams.go:551-568`, `09-terminal-ux`, chezmoi `internal/cmd/prompt.go:124-137`, `13-security`, restic `internal/options/secret_string.go:15-20`, helm `pkg/registry/transport.go:37-41`, and `04-configuration-management` for testable output seams, calm scriptable output, redaction, and resolved config discipline.
- **Trade-Offs Accepted:** Prompt tests should assert required variables, selected paths, ordering, and rules rather than all prose, accepting limited wording drift.
- **Technical Debt / Future Impact:** If prompt wording becomes a compatibility concern, stable prompt blocks may need golden fixtures. For this sprint, invariant tests are enough.
- **Alternatives Rejected:** Invoking runtime from `prompt plan` is rejected because preview must be runtime-free. Writing `plan.md` from prompt preview is rejected because prompt preview must not generate artifacts. Embedding unselected evidence or absolute local paths by default is rejected because selected-context and portability rules require workspace-relative prompt content unless local diagnostics need otherwise.
- **Contracts Satisfied:** Security, Configuration, CLI Surface, Documentation, LLM Runtime, AC-02, AC-10, AC-11, AC-12, AC-22, AC-23, AC-28, AC-30, C-07, C-08, C-09.
- **Evidence Required:** Prompt-rendering tests for required variables, selected reasoning input paths including selected Architecture area path, output path, required sections, traceability rules, no-mutation/no-execution rules, workspace-relative paths, deterministic ordering, no runtime calls, and no `plan.md` writes.

### Decision 4: Flow Sequentially Through `plan` With Post-Generation Gates

- **Decision:** Extend sprint flow so `flow --to plan --dry-run` reports required stages, validated reasoning inputs, expected `plan.md`, and flow-state impact without runtime or mutation; `flow --to plan` validates all prerequisites, invokes only the generic platform runtime boundary for plan generation, requires `plan.md` to exist and pass validation after runtime success, and atomically updates `flow-state.json` only after validation passes.
- **Rationale:** The flow must preserve the Phase 2 stage chain while adding plan generation. Runtime success alone is not product success, and failed generation or validation must leave safe, inspectable state without introducing later execution/smoke/review/issue/Git stages.
- **Study / Source Grounding:** `technical-handbook.md` cites `07-state-context`, gh-cli `internal/ghcmd/cmd.go:142`, helm `pkg/cmd/install.go:333-347`, restic `cmd/restic/cleanup.go:24-38`, `10-logging-observability`, helm `internal/logging/logging.go:31-66`, and `13-security`, opencode `internal/llm/tools/bash.go:41-55`, for explicit state, context propagation, safe diagnostics, and constrained runtime execution.
- **Trade-Offs Accepted:** The flow remains a sequential planning-stage orchestrator rather than a generic workflow engine. This keeps scope small but means future implementation execution will need a separately reasoned model.
- **Technical Debt / Future Impact:** Flow-state handling will grow one more stage and failure mode. Strict stage validation and tests should prevent accidental future-stage creep.
- **Alternatives Rejected:** Marking `plan` complete immediately after runtime success is rejected because it violates validation-as-completion. Directly invoking OpenCode, shell commands, provider APIs, agentwrap adapters, or Git from `internal/sprint` is rejected because the platform runtime boundary and security constraints prohibit it. Adding `implementation`, `execute`, `smoke`, `review`, `issues`, `.run-state.json`, smoke artifacts, review artifacts, or issue artifacts as supported stages is rejected because Phase 2 stops at `plan`.
- **Contracts Satisfied:** Workflows, Persistence And Migrations, Security, LLM Runtime, LLM Evaluation / Cost / Safety, AC-03, AC-04, AC-05, AC-06, AC-13, AC-14, AC-15, AC-16, AC-17, AC-24, AC-25, AC-26, AC-29, C-04, C-07, C-10, C-11.
- **Evidence Required:** Flow tests for dry-run non-mutation, successful fake-runtime generation, runtime success with missing artifact, runtime success with invalid artifact, runtime failure, invalid prerequisites, invalid selected area reasoning, post-generation validation failure, atomic state write failure, successful stage completion through `plan`, failed state recording, and unsupported-stage rejection.

### Decision 5: Keep Service, Store, Runtime, Output, And Exit-Code Seams Narrow

- **Decision:** Add narrow sprint service methods for validating plan, rendering plan prompts, and flowing to plan; update the filesystem store to read `reasoning.md`, selected `reasoning/*.md`, and `plan.md` through workspace-safe paths; inject the generic runtime seam only into flow; and keep command handlers responsible only for argument parsing, service calls, stdout/stderr rendering, and exit-code mapping.
- **Rationale:** Validate and prompt commands must be provably runtime-free, while flow needs fake-runtime tests and generic runtime execution. Narrow seams provide testability without turning sprint into a generic framework.
- **Study / Source Grounding:** `technical-handbook.md` cites `03-dependency-injection`, gh-cli `pkg/cmd/factory/default.go:26-46`, go-task `executor.go:22-24`, restic `internal/backend/backend.go:19-90`, `06-io-abstraction`, and `05-error-handling` for explicit construction, fakes, writer seams, and exit-code mapping.
- **Trade-Offs Accepted:** More request/result types and constructor wiring are acceptable because they directly support command tests, fake runtime, and safe output checks.
- **Technical Debt / Future Impact:** If service construction becomes too broad, split request handling internally by file before introducing interfaces or subpackages. Preserve cause chains and typed validation findings for future JSON output.
- **Alternatives Rejected:** Direct filesystem/runtime access in command handlers is rejected because it bypasses module ownership and test seams. A broad service locator or context-carried dependencies are rejected because they hide side effects and make tests order-dependent. Direct `os.Stdout`/`os.Stderr` writes in domain logic are rejected because command tests require output capture.
- **Contracts Satisfied:** Architecture, Errors, Observability, CLI Surface, Security, Testing, AC-01, AC-02, AC-21, AC-22, AC-23, AC-24, AC-25, AC-26, AC-30, C-01, C-08, C-12.
- **Evidence Required:** Command tests proving `validate plan` and `prompt plan` do not call runtime, `flow --to plan` uses a fake runtime seam, stdout/stderr are separated, no ANSI escapes are emitted, exit codes map to requirements, safe diagnostics are redacted, and unsupported stages return usage exit code `2`.

### Decision 6: Verify Offline With Focused Plan And Command Tests

- **Decision:** Implement deterministic offline tests centered on `internal/sprint/plan_test.go` and `internal/app/sprint_commands_test.go`, plus existing package tests, race tests, and CLI build. Use fixtures and fakes for runtime and filesystem behavior; do not require real OpenCode, network, provider credentials, smoke harnesses, or Git mutation for default verification.
- **Rationale:** Sprint 21 changes validation, prompt rendering, flow, state, runtime seams, and CLI behavior. The evidence must prove these behaviors at the module and command level without relying on external runtime availability.
- **Study / Source Grounding:** `technical-handbook.md` cites `11-testing-strategy`, chezmoi `internal/cmd/main_test.go:64-174`, helm `internal/test/test.go:43`, go-task `task_test.go:166-169`, and rclone `fstest/mockfs/mockfs.go:31-62` for command integration tests, fixtures, goldens, fakes, and avoiding brittle assertions.
- **Trade-Offs Accepted:** Default tests will not prove real OpenCode plan generation; they prove UltraPlan behavior around prompt construction, runtime requests, artifact validation, and state transitions with fakes.
- **Technical Debt / Future Impact:** Real-runtime smoke can be added later through the cataloged smoke harness if a reviewer requests deep smoke evidence. Keep fake-runtime seams compatible with that future path.
- **Alternatives Rejected:** Gating sprint acceptance on real OpenCode smoke is rejected because requirements call for deterministic offline verification. Testing only domain validators without command coverage is rejected because exit codes, stdout/stderr, help, and runtime-free paths are acceptance criteria. Full exact prompt-prose goldens for all output are rejected because they would be brittle; invariant assertions are preferred.
- **Contracts Satisfied:** Testing, Documentation, CLI Surface, Observability, Security, AC-27, AC-28, AC-29, AC-30, AC-31, AC-32, AC-33, NG-06, NG-11, C-09.
- **Evidence Required:** `go test ./...`, `go test -race ./...`, `go build ./cmd/ultraplan`, plan validator unit tests, prompt invariant tests, fake-runtime flow tests, command tests, help-output checks, and Architecture/Sprint Review protocol entries in `projects/ultraplan-go/sprints/21-plan-stage/review.md`.

## Expected Evidence

| Evidence Type | Required Evidence | Source / Command / Review Check |
| --- | --- | --- |
| Tests | Plan validation covers valid content, missing file, empty file, placeholders, missing `reasoning.md` citation/reference, missing decisions, missing tasks, missing evidence checklist, missing risks/blockers, missing success criteria, untraced tasks, and forbidden deferred-stage content. | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/plan_test.go` |
| Tests | Prompt preview covers required variables, selected context paths, selected Architecture area reasoning path, output path, required sections, traceability rules, no-mutation/no-execution rules, deterministic ordering, workspace-relative paths, no runtime calls, and no `plan.md` writes. | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/plan_test.go` and `/home/antonioborgerees/coding/ultraplan-go/internal/app/sprint_commands_test.go` |
| Tests | Flow covers dry-run non-mutation, fake-runtime success, missing generated artifact, invalid generated artifact, runtime failure, invalid prerequisites, selected area reasoning prerequisite failure, validation failure after generation, flow-state write failure, stage completion through `plan`, failed-stage recording, and unsupported-stage rejection. | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/plan_test.go` and related sprint flow tests |
| Tests | Command tests cover `validate plan`, `prompt plan`, `flow --to plan`, help output, malformed arguments, unsupported stages, exit codes `0`, `2`, and `5`, runtime failure mapping, stdout/stderr separation, no ANSI output, redaction, runtime-free validate/prompt paths, and fake-runtime flow. | `/home/antonioborgerees/coding/ultraplan-go/internal/app/sprint_commands_test.go` |
| Runtime | Runtime-backed flow uses only the generic platform runtime boundary, and fake runtime proves success and failure paths; no default test requires OpenCode, provider credentials, network, smoke harnesses, shell commands, Git mutation, or real runtime execution. | Fake runtime test seam and import review |
| State | `flow-state.json` remains strict, versioned, atomically written, preserves the last valid state on write failure, supports only Planning Phase 2 stages, records `plan` complete only after validation, and records failed plan flow with safe diagnostics. | Sprint flow/state tests and code review |
| Review | Architecture boundaries are preserved: `internal/sprint` owns plan behavior, `internal/app` remains thin, platform packages are product-agnostic, `internal/study` is not imported by sprint, and no global planning/workflow/validation/prompt package is introduced. | `system/protocols/architecture-review-protocol.md` applied in `projects/ultraplan-go/sprints/21-plan-stage/review.md` |
| Review | Sprint review records implementation evidence, deviations, verification results, unsupported-stage exclusions, selected-context compliance, and conclusions. | `system/protocols/sprint-review-protocol.md` applied in `projects/ultraplan-go/sprints/21-plan-stage/review.md` |
| Verification | Offline verification passes. | `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go` |
| Verification | Race verification passes. | `go test -race ./...` from `/home/antonioborgerees/coding/ultraplan-go` |
| Verification | CLI build passes. | `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go` |

## Assumptions And Risks

| Item | Type | Impact | Mitigation / Follow-Up |
| --- | --- | --- | --- |
| The selected Architecture area reasoning artifact is absent. | Risk | Runtime-backed plan flow should fail prerequisite validation until the selected area reasoning file exists, even though this final reasoning makes direct decisions from other selected inputs. | Produce `projects/ultraplan-go/sprints/21-plan-stage/reasoning/architecture.md` or revise `sprint-index.md` before relying on implemented flow to pass prerequisites. Keep tests for selected area reasoning requirements. |
| Existing sprint command infrastructure already supports stage parsing and common exit-code/output conventions from prior sprints. | Assumption | If names or seams differ, implementation may need minimal adaptation while preserving decisions. | Inspect existing `internal/sprint` and `internal/app` code before editing; keep changes minimal and module-owned. |
| Existing platform runtime abstraction can support a prompt, working directory, config, timeout, and expected output path without product semantics. | Assumption | If the seam lacks an expected-output concept, sprint flow must still validate `plan.md` after runtime success locally. | Use only generic runtime request fields and keep product validation in `internal/sprint`. |
| Plan validation may produce false positives if it overfits exact prose. | Risk | Human-edited valid plans could fail unnecessarily. | Check structural elements, explicit references, checklist/task/evidence traceability, placeholders, and forbidden content rather than exact wording. |
| Prompt content could accidentally include unsafe or unselected context. | Risk | Generated plans could cite unselected evidence or leak local diagnostics/secrets. | Prompt tests must assert selected paths only, workspace-relative paths where possible, redaction, and no-mutation rules. |
| Unsupported future-stage language may appear in legitimate explanatory contexts. | Risk | Validator could reject plans that mention deferred work only as non-goals. | Validator should reject invocation, automation, modeling, or tasking of deferred stages as current Phase 2 behavior, while allowing explicit non-goal/risk statements. |
| Flow-state write failure after valid generation can leave `plan.md` present but stage not complete. | Risk | User may see a generated valid plan but failed state. | Return non-zero, preserve last valid state, emit safe diagnostic, and require rerun or validation to complete state transition. |
| Real runtime behavior is not exercised by default verification. | Risk | Integration issues with OpenCode/agentwrap may appear after offline tests pass. | Keep generic runtime boundary clean and use fake-runtime tests by default; run optional deep smoke only when selected by review. |

## Implementation Constraints

- `internal/sprint` owns plan validation, plan prompt rendering, reasoning prerequisite checks, task/evidence trace checks, service use cases, and flow-state transitions for `plan`.
- `internal/app` CLI handlers may parse arguments, call sprint services, render stdout/stderr, and map exit codes only.
- `internal/sprint` must not import `internal/study` or reuse study services, source/dimension models, report/rating logic, summary generation, scheduling, or run-loop state.
- `internal/platform/*` packages must not import `internal/sprint`, `internal/project`, `internal/study`, or `internal/app`.
- Runtime-backed plan flow must use only the generic platform runtime boundary; `internal/sprint` must not invoke OpenCode, agentwrap adapters, provider APIs, shell commands, or Git directly.
- `validate plan`, `prompt plan`, and `flow --to plan --dry-run` must not invoke runtime or create/overwrite `plan.md`.
- Runtime success is insufficient; `plan` completes only after runtime success, `plan.md` exists, `plan.md` validates, and `flow-state.json` is atomically updated.
- Prompt content must use workspace-relative paths unless an absolute path is required for local diagnostics, must include selected context and no-mutation rules, and must not include secrets or unselected authoritative inputs.
- `flow-state.json` must remain strict, versioned, atomically written, and limited to `requirements`, `sprint-index`, `technical-handbook`, `area-reasoning`, `reasoning`, and `plan`.
- Unsupported stages and artifacts such as `implementation`, `execute`, `smoke`, `review`, `issues`, `.run-state.json`, `smoke.md`, `smoke.json`, product-generated `review.md`, `issues.md`, and `issues.json` must be rejected or ignored as out of current Planning Phase 2 scope.
- Plan validation must allow editable Markdown but reject missing structure, placeholders, missing `reasoning.md` traceability, untraced tasks, missing evidence, missing risks/blockers, missing success criteria, and forbidden deferred-stage behavior.
- Tests must be deterministic, offline, fake-first, and free of long sleeps, network calls, provider credentials, real OpenCode execution, smoke harness dependency, and Git mutation.

## Plan Handoff

`plan.md` must execute these decisions. It must not invent architecture, scope, authoritative evidence, or decisions beyond this document.

The plan must carry forward:

- Final decisions 1 through 6.
- Applicable contracts and requirement IDs: Architecture, Errors, Configuration, Observability, Security, Testing, Documentation, CLI Surface, LLM Runtime, LLM Evaluation / Cost / Safety, Workflows, Persistence And Migrations, AC-01 through AC-33, NG-01 through NG-11, and C-01 through C-12 as applicable to tasks.
- Expected evidence covering plan validation, prompt rendering, flow, command behavior, runtime seam, flow-state writes, import boundaries, review protocols, `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan`.
- Risks and assumptions, especially the missing selected Architecture area reasoning artifact and selected-context discipline.
- Required Architecture Review and Sprint Review protocols.

## Phase Exit Criteria

- [x] Selected context was read and used.
- [x] Area-specific reasoning documents were checked; the selected Architecture area reasoning document was not present, and this gap is explicitly recorded as a risk rather than treated as a completed input.
- [x] Area-specific reasoning conclusions are not fabricated; final decisions are derived from requirements, project docs, sprint index, and technical handbook evidence.
- [x] Contracts are explicitly mapped to decisions and expected evidence.
- [x] Final decisions are clear enough for `plan.md` to execute.
- [x] Expected evidence is specific and reviewable.
