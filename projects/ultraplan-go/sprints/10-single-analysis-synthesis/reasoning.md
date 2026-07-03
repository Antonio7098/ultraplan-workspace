# Sprint Reasoning: 10 Single Analysis and Synthesis

> Project: `ultraplan-go`
> Sprint: `10-single-analysis-synthesis`
> Output: `projects/ultraplan-go/sprints/10-single-analysis-synthesis/reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/10-single-analysis-synthesis/requirements.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/sprints/10-single-analysis-synthesis/sprint-index.md`, `projects/ultraplan-go/sprints/10-single-analysis-synthesis/technical-handbook.md`, `projects/ultraplan-go/sprints/10-single-analysis-synthesis/reasoning/single-analysis-synthesis.md`, `templates/sprint-reasoning.md`

This document decides. It synthesizes selected context, handbook evidence, area-specific reasoning, and contracts into final sprint decisions.

It does not replace `sprint-index.md`, `technical-handbook.md`, or `reasoning/*.md`.

## Sprint Purpose

- **Goal:** Implement `ultraplan study <study> run <dimension-ref> <source-ref>` and `ultraplan study <study> synthesize <dimension-ref>` so one analysis task or one synthesis task executes through the agentwrap-backed runtime and succeeds only after the expected report artifact exists and passes study-owned validation.
- **Non-Goals:** Do not implement `study run-all`, task matrix execution, worker pools, multi-task scheduling, `study run-loop`, resumable batch mutation, stale running recovery, cross-process retry scheduling, per-study locks, automatic synthesis after analysis, summary generation, code reference extraction, product-owned repair loops, durable agentwrap run-store persistence, stable public execution JSON schemas, target workflows, sprint planning, sprint execution, or direct OpenCode supervision by UltraPlan.
- **Depends On:** Sprints 1-9 provide buildable CLI/app wiring, workspace/config/logging/health foundations, study/dimension/source resolution, study artifact layout, Markdown source discovery and applicability, report validators and rating parsing, deterministic prompt builders, task/state vocabulary compatibility, and agentwrap/OpenCode platform runtime integration.

## Selected Context And Pre-Reasoning Artifacts

| Artifact | Path | How It Was Used |
| --- | --- | --- |
| Project Index | `projects/ultraplan-go/project-index.md` | Established the authoritative docs, selected contract pool, selected evidence report pool, repository target, and non-goal that target/sprint workflows remain deferred. |
| Sprint Requirements | `projects/ultraplan-go/sprints/10-single-analysis-synthesis/requirements.md` | Bound the sprint goal, required outputs, acceptance criteria, constraints, dependencies, non-goals, and review expectations. |
| Architecture Doc | `projects/ultraplan-go/docs/ARCHITECTURE.md` | Governed package ownership: `internal/study` owns study workflows and validation, `internal/app` stays thin, and `internal/platform/runtime` remains generic. |
| PRD | `projects/ultraplan-go/docs/PRD.md` | Grounded product behavior: single analysis, synthesis, Markdown applicability, validation as product success gate, deterministic artifacts, and deferred target/sprint scope. |
| TRD | `projects/ultraplan-go/docs/TRD.md` | Grounded technical behavior: agentwrap runtime use, run request mapping, event handling, error classification, validation, permissions, testing, build, and no direct OpenCode supervision. |
| Sprint Index | `projects/ultraplan-go/sprints/10-single-analysis-synthesis/sprint-index.md` | Selected the governing contracts, evidence reports, reasoning artifact, review protocols, scope, non-scope, and excluded context. |
| Technical Handbook | `projects/ultraplan-go/sprints/10-single-analysis-synthesis/technical-handbook.md` | Supplied evidence from studied Go CLIs for thin commands, explicit composition, fake runtimes, context propagation, safe diagnostics, validation gates, and avoiding scheduler/plugin overreach. |
| Area Reasoning | `projects/ultraplan-go/sprints/10-single-analysis-synthesis/reasoning/single-analysis-synthesis.md` | Resolved the main architecture choice: study-owned single execution over the generic runtime boundary, with thin CLI, study-owned validation, Markdown skip, synthesis preflight, fake-runtime tests, and no batch/workflow expansion. |
| Sprint Reasoning Template | `templates/sprint-reasoning.md` | Provided the required structure for this consolidated reasoning document. |

## Area-Specific Reasoning Inputs

| Area | Reasoning Document | Key Conclusion | Evidence Basis | Impact On Final Decision |
| --- | --- | --- | --- | --- |
| Single Analysis and Synthesis Architecture | `projects/ultraplan-go/sprints/10-single-analysis-synthesis/reasoning/single-analysis-synthesis.md` | Implement single analysis and synthesis in `internal/study`, call the existing generic `internal/platform/runtime` boundary, keep CLI code thin, treat study-owned artifact validation as the decisive product success gate, and avoid batch/workflow expansion. | Cites `ARCHITECTURE.md` ownership rules, `requirements.md` acceptance criteria, `PRD.md` runtime-success-is-not-product-success principle, `TRD.md` agentwrap/run-request/error/validation rules, and `technical-handbook.md` evidence for thin commands, explicit composition, fake external systems, context-aware single execution, safe diagnostics, and bounded abstractions. | Accepted as the controlling sprint architecture. Final decisions below refine validation placement, runtime request mapping, result/error output, and testing evidence without overriding the area reasoning. |

## Sprint Technical Handbook Summary

- **Relevant Patterns:** Thin command delegation, explicit composition roots, generic runtime boundaries, output validation as a product gate, injectable IO and fake external systems, context-aware single execution without batch machinery, safe structured diagnostics, and config preflight before work launch.
- **Important Trade-Offs:** Keep orchestration in `internal/study` instead of `internal/app`; keep runtime generic instead of product-aware; use generic expected-output runtime validation plus definitive study-owned report validation; accept readable analysis/synthesis flow duplication instead of a premature scheduler; prefer deterministic text output over rich progress UX.
- **Warnings / Anti-Patterns:** Do not let CLI handlers become orchestration bodies, do not put study semantics in platform runtime, do not treat runtime exit success as product success, do not hardcode real runtime or stdout/stderr in tests, do not introduce speculative plugin/scheduler frameworks, do not parse human runtime text for control flow, do not let event channels block `Wait`, and do not leak secrets or unsafe native payloads.
- **Evidence Confidence:** High overall. Most relevant reports are marked high confidence in `technical-handbook.md`, especially project structure, command architecture, dependency injection, configuration, error handling, IO abstraction, state/context, observability, testing, security, and philosophy. Concurrency and terminal UX evidence are medium confidence but used only to support limited decisions about single-run event draining and deterministic output.

## Contracts Applied

| Contract / Requirement ID | Constraint | Decision Impact | Expected Evidence |
| --- | --- | --- | --- |
| Architecture contract; `requirements.md:50-51`, `requirements.md:75-80` | Study workflow behavior must stay in `internal/study`; `internal/platform/runtime` must stay generic and product-agnostic. | Final decisions place resolution, applicability, prompt selection, synthesis preflight, and report validation in `internal/study`; runtime receives only generic request fields and metadata. | Import review confirms no study imports in platform runtime; tests confirm study service owns execution rules. |
| Errors contract; `requirements.md:40-42`, `requirements.md:49`, `requirements.md:52-53`, TRD section 14.2 | Runtime, validation, usage, workspace/config, and cancellation outcomes need classified diagnostics and correct exit mapping without string parsing. | Final decisions require typed/classified errors, preserved cause chains, runtime-vs-validation distinction, and stdout/stderr separation. | CLI command tests cover exit statuses and diagnostics; code review checks `errors.As` or platform classifications rather than human text parsing. |
| Configuration contract; `requirements.md:38`, `requirements.md:47`, TRD section 6 | Runtime/provider/model/timeout/permission/health settings must be resolved and validated before launch. | Runtime request mapping includes resolved provider, model, timeout, permission policy, health checks, capability checks, and safe config summaries. | Fake-runtime request assertions verify config fields; config failure tests prove fail-before-launch behavior where feasible. |
| Observability contract; `requirements.md:39`, `requirements.md:47`, `requirements.md:52-53` | Runtime metadata, task metadata, safe diagnostics, and usage/cost fields where available must be preserved safely. | Runtime metadata keys are flattened strings and safe values; CLI diagnostics include report paths and failed checks without prompts, secrets, or unsafe raw payloads. | Request assertions inspect metadata; output tests check safe diagnostics; review checks no prompt/env/native raw payload leaks. |
| Security contract; `requirements.md:35-38`, `requirements.md:45-47`, `requirements.md:81-85`, TRD section 22 | Source/document isolation, path safety, permission posture, redaction, and no direct OpenCode process handling are required. | Directory analysis uses source directory workdir; Markdown analysis uses study/workspace safe workdir with embedded stripped document content; execution routes through platform runtime/agentwrap only. | Fake-runtime tests assert workdir choice; prompt/validation tests cover Markdown behavior; review confirms no direct OpenCode invocation. |
| Testing contract; `requirements.md:55-60`, TRD section 23 | Default tests must use fakes and avoid OpenCode, credentials, network, wall-clock sleeps, and real subprocesses. | Study and CLI tests use fake runtime and buffer-backed IO; real smoke tests, if added, are opt-in and skipped by default. | `go test ./...` passes without OpenCode/provider credentials; `go build ./cmd/ultraplan` passes. |
| CLI Surface contract; `requirements.md:31-33`, `requirements.md:37`, `requirements.md:43`, `requirements.md:50`, `requirements.md:53-54` | Command shape, help, argument handling, output determinism, stdout/stderr separation, and optional JSON discipline must be respected. | CLI handlers parse arguments, call service methods, render summaries/skips to stdout, render diagnostics to stderr, and do not add stable JSON unless documented and tested. | Command tests cover run/synthesize usage, success, skip, runtime failure, validation failure, and output streams. |
| LLM Runtime contract; `requirements.md:38-39`, `requirements.md:47`, `requirements.md:52`, TRD sections 11-12 | Runtime requests must go through platform runtime/agentwrap with prompt, workdir, provider/model/timeout, metadata, permissions, health/caps, validation, event consumption, and final result handling. | Study builds generic runtime requests and consumes final platform runtime results; platform/agentwrap owns OpenCode supervision, event projection, policy, health, permissions, and validation wrappers. | Fake-runtime request assertions and code review verify no duplicate runtime supervisor and no blocking event handling. |
| LLM Evaluation / Cost / Safety contract; `requirements.md:40`, `requirements.md:48`, `requirements.md:83`, PRD principle | Runtime success is not product success; validation and safe metadata handling are required. | Study-owned per-source and final report validators are definitive gates after runtime completion; usage/cost fields are optional and unknown values stay unknown. | Tests cover runtime success with missing/invalid outputs; review checks no unknown usage rendered as zero. |

## Repos Studied / Source Evidence Used

| Source / Repo / Report | Concrete Reference | Relevant Finding | Why It Matters For This Sprint | Used In Decision(s) |
| --- | --- | --- | --- | --- |
| `01-project-structure` report | `studies/go-cli-study/reports/final/01-project-structure.md`; chezmoi `main.go:26-34`, gh-cli `cmd/gh/main.go:6`, yq `cmd/root.go:9`, restic `internal/restic/repository.go:18` | Mature Go CLIs keep entrypoints thin and imports unidirectional from CLI to domain/internal packages. | Supports keeping `internal/app` thin and avoiding product semantics in platform runtime. | Decisions 1, 4 |
| `02-command-architecture` report | `studies/go-cli-study/reports/final/02-command-architecture.md`; helm `pkg/cmd/install.go:132-145`, gh-cli `pkg/cmdutil/factory.go:16-43` | Command constructors and handlers delegate to reusable action/service logic. | Supports CLI parse/render/error-mapping only, not prompt/runtime/validation orchestration. | Decisions 1, 4 |
| `03-dependency-injection` report | `studies/go-cli-study/reports/final/03-dependency-injection.md`; opencode `internal/app/app.go:42-81`, gh-cli `pkg/cmd/factory/default.go:26-46`, go-task `executor.go:22-24` | Manual composition roots, constructor injection, and focused interfaces are common and testable. | Supports wiring the platform runtime into the study service through app composition and fakeable seams. | Decisions 1, 3, 6 |
| `04-configuration-management` report | `studies/go-cli-study/reports/final/04-configuration-management.md`; chezmoi `internal/cmd/config.go:2253-2287`, go-task `internal/flags/flags.go:314-327`, opencode `internal/config/config.go:609-641` | Effective config should be merged and validated before expensive work starts. | Supports fail-before-launch runtime config handling and request fields for provider/model/timeout/health/permissions. | Decision 3 |
| `05-error-handling` report | `studies/go-cli-study/reports/final/05-error-handling.md`; rclone `fs/fserrors/error.go:22-29`, go-task `errors/errors.go:47-50`, restic `internal/errors/fatal.go:10` | High-quality CLIs preserve cause chains and classify errors for user-facing behavior. | Supports runtime vs validation vs usage/config/cancellation distinction and exit mapping. | Decisions 4, 5 |
| `06-io-abstraction` report | `studies/go-cli-study/reports/final/06-io-abstraction.md`; gh-cli `pkg/iostreams/iostreams.go:551-568`, go-task `executor.go:553-564`, restic `internal/ui/mock.go:10-53` | Commands should inject IO and external systems for deterministic tests. | Supports buffer-backed CLI tests and fake-runtime execution tests. | Decisions 4, 6 |
| `07-state-context` report | `studies/go-cli-study/reports/final/07-state-context.md`; gh-cli `internal/ghcmd/cmd.go:142`, helm `pkg/cmd/install.go:333-347`, go-task `task.go:89` | External or long-running work should propagate `context.Context` and cancellation. | Supports context-aware runtime calls without adding batch state machinery. | Decisions 3, 5 |
| `08-concurrency` report | `studies/go-cli-study/reports/final/08-concurrency.md`; opencode `cmd/root.go:130,252`, lazygit `pkg/tasks/tasks.go:338-357`, restic `internal/repository/repository.go:567` | Cleanup and event/channel handling matter even for bounded units; general concurrency is only justified for fan-out. | Supports draining/consuming events and waiting for final results, while rejecting worker pools and schedulers. | Decisions 3, 7 |
| `09-terminal-ux` report | `studies/go-cli-study/reports/final/09-terminal-ux.md`; gh-cli `pkg/iostreams/iostreams.go:116-130`, restic `internal/ui/terminal.go:10-34` | Output should match interaction density and remain deterministic when needed. | Supports concise success/skip summaries and deterministic diagnostics instead of rich progress UX. | Decision 4 |
| `10-logging-observability` report | `studies/go-cli-study/reports/final/10-logging-observability.md`; helm `internal/logging/logging.go:31-66`, k9s `internal/slogs/keys.go:6-231`, gh-cli `pkg/iostreams/iostreams.go:52-54` | Operators need structured diagnostics, stable fields, and stdout/stderr separation. | Supports metadata and diagnostics requirements without leaking unsafe payloads. | Decisions 3, 4 |
| `11-testing-strategy` report | `studies/go-cli-study/reports/final/11-testing-strategy.md`; chezmoi `internal/cmd/main_test.go:64-174`, gh-cli `pkg/httpmock/stub.go:35-199`, lazygit `pkg/commands/oscommands/fake_cmd_obj_runner.go:17-26` | Strong CLIs use fakes, fixtures, table tests, and command-level tests for external systems. | Supports fake-runtime service tests and CLI tests as default evidence. | Decision 6 |
| `12-extensibility` report | `studies/go-cli-study/reports/final/12-extensibility.md`; helm `internal/plugin/runtime_subprocess.go:65-79`, rclone `fs/registry.go:407`, go-task `executor.go:20-24` | Extension points and registries add cost and should be bounded by actual need. | Supports rejecting plugin systems, broad registries, and general workflow engines for this sprint. | Decision 7 |
| `13-security` report | `studies/go-cli-study/reports/final/13-security.md`; opencode `internal/permission/permission.go:44-108`, restic `internal/options/secret_string.go:15-20`, helm `pkg/registry/transport.go:37-41` | Runtime execution needs explicit trust boundaries, permission gates, argument-safe process handling, and redaction. | Supports source/document isolation, permission policy mapping, and safe runtime diagnostics through agentwrap. | Decisions 2, 3, 4 |
| `15-philosophy` report | `studies/go-cli-study/reports/final/15-philosophy.md`; lazygit `VISION.md:97`, rclone `fs/types.go:16-59`, restic `internal/backend/backend.go:19-90` | Coherent Go CLIs accept complexity only where it serves product goals and keep non-goals explicit. | Supports the minimal single-task design and rejection of scheduler/repair/JSON/batch scope. | Decisions 1, 7 |

## Trade-Off And Debt Analysis

### Accepted Trade-Offs

| Trade-Off | Benefit | Cost / Constraint Accepted | Why Acceptable Now | Revisit Trigger |
| --- | --- | --- | --- | --- |
| Study-owned orchestration instead of command-contained execution | Preserves product ownership, improves service-level testing, and supports later batch reuse without CLI duplication. | Requires service request/result types and explicit CLI rendering. | Requirements explicitly require study execution logic in `internal/study` and thin CLI handlers. | Revisit only if a future non-CLI API requires a different service boundary. |
| Generic runtime boundary instead of study-aware runtime helpers | Keeps platform runtime reusable and prevents reverse product coupling. | Study code must flatten metadata and own validation/report semantics. | Architecture and TRD explicitly prohibit runtime from understanding studies, reports, or synthesis gates. | Revisit only with a reviewed architecture exception if platform runtime cannot express a generic needed behavior. |
| Generic runtime expected-output validation plus definitive study post-run validation | Runtime wrappers can observe missing outputs, while study validators remain product source of truth. | Risk of duplicate diagnostics if not coordinated. | Requirements require expected output validation and study-owned validators; both are useful when clearly separated. | Revisit if users see duplicate or conflicting validation messages. |
| Readable separate analysis and synthesis flows instead of a generic task executor | Keeps distinct preconditions clear: Markdown skip for analysis and source-report preflight for synthesis. | Some shared code shape remains. | The sprint executes one task at a time and excludes schedulers/workflow engines. | Revisit when `run-all` or `run-loop` needs shared scheduling logic. |
| Deterministic text output instead of rich event/progress streaming | Makes command tests stable and stdout/stderr behavior script-friendly. | Less interactive feedback during long runtime execution. | Single execution scope needs clear summaries and diagnostics more than a TUI/progress layer. | Revisit when batch status, logs, or TUI/watch mode enters scope. |
| Fake-runtime default tests instead of real OpenCode coverage by default | Keeps normal verification deterministic and credential-free. | Real adapter drift is not fully covered by default tests. | Requirements prohibit default tests from needing OpenCode, credentials, network, or subprocesses. | Revisit with gated smoke tests when runtime behavior needs real-adapter confidence. |

### Potential Technical Debt

| Debt / Shortcut | Why It Might Accrue | Current Mitigation | Owner / Follow-Up |
| --- | --- | --- | --- |
| Minimal execution result model may need expansion for batch/run-loop | Future durable runs need attempts, retry timing, event records, and state reconciliation beyond single-command results. | Keep fields compatible with existing task/status vocabulary and do not persist a new schema now. | Future `run-all` or `run-loop` sprint. |
| Validation diagnostics may be split between runtime validation and study validators | Users may see overlapping missing-output or invalid-report information. | Treat study validator summary as definitive product gate and render one coherent diagnostic. | Sprint 10 implementation review; future diagnostics polish if needed. |
| Permission policy defaults may be too broad or too restrictive for Markdown document analysis | Markdown tasks need to write reports but should not explore external files. | Use existing permission policy capabilities, document-specific prompts, safe workdir, and fake-runtime request checks. | Sprint 10 implementation review; future security hardening. |
| Fake runtime may not model every agentwrap/OpenCode edge case | Default tests may miss adapter-specific event or metadata behavior. | Keep fake tests focused on UltraPlan-owned mapping/gating and allow optional environment-gated smoke tests. | Runtime integration follow-up or release hardening. |
| CLI JSON output remains unspecified if not implemented | Automation consumers may later request stable execution result schemas. | Do not add JSON output unless documented and tested in this sprint. | Future CLI automation sprint. |

### Future Considerations

| Consideration | Deferred Until | Reason Deferred | What Should Be Preserved Now |
| --- | --- | --- | --- |
| `study run-all` task matrix and bounded workers | Batch execution sprint | Explicit non-goal; single-task semantics should be proven first. | Service methods and result fields should be reusable without depending on CLI state. |
| Stateful `study run-loop`, locks, restart recovery, and durable agentwrap run-store | Durable workflow sprint | Requires persistence, concurrency, stale task recovery, lock policy, and resume reconciliation outside this sprint. | Do not mutate batch state; preserve task IDs/status vocabulary compatibility. |
| Automatic synthesis after analysis | Batch/completion workflow sprint | This sprint keeps `study run` and `study synthesize` separate. | Synthesis method should be callable independently after valid reports exist. |
| Product repair prompts and bounded repair loops | Repair workflow sprint | Requirements defer repair unless already supplied by configured runtime wrappers. | Keep validation errors actionable and preserve failed artifact paths/checks. |
| Summary generation and code-reference extraction | Summary/code extraction sprints | Explicit non-goals for single execution. | Do not couple execution results to summary or citation extraction packages. |
| Stable public JSON schemas | CLI automation/release stabilization | Sprint index excludes stable public execution JSON. | Avoid undocumented JSON or test/document it if added. |
| Target and sprint workflows | Future product phase after PRD/TRD revision | Current project scope is study-side only. | Do not add target/sprint commands, validators, task kinds, or product semantics. |

## Final Decisions

### Decision 1: Study-Owned Single Execution

- **Decision:** Implement single analysis and single synthesis as `internal/study` service behavior. `internal/study/run.go` owns one analysis run, `internal/study/synthesize.go` owns one synthesis run, and `internal/study/service.go` exposes the service methods with a runtime dependency. CLI code only parses, delegates, renders, and maps exits.
- **Rationale:** The selected behavior is study workflow behavior: resolving studies/dimensions/sources, applying Markdown applicability, building prompts, computing report paths, validating reports, and enforcing synthesis gates. Putting this in `internal/study` matches the module-driven architecture and keeps future non-CLI callers possible.
- **Study / Source Grounding:** `technical-handbook.md` patterns for thin CLI with study-owned behavior and explicit composition cite chezmoi `main.go:26-34`, gh-cli `cmd/gh/main.go:6`, helm `pkg/cmd/install.go:132-145`, gh-cli `pkg/cmdutil/factory.go:16-43`, opencode `internal/app/app.go:42-81`, and gh-cli `pkg/cmd/factory/default.go:26-46`.
- **Trade-Offs Accepted:** Accept service request/result types and some orchestration code in `internal/study` instead of quick inline command code.
- **Technical Debt / Future Impact:** The result model may grow later for batch/run-loop, but no durable schema is introduced now.
- **Alternatives Rejected:** Rejected command-contained execution because it duplicates study rules in `internal/app`; rejected a general scheduler/workflow engine because `run-all`, worker pools, locks, retries across process restarts, and run-loop are non-goals.
- **Contracts Satisfied:** Architecture, CLI Surface, Testing; `requirements.md:31-33`, `requirements.md:43-44`, `requirements.md:50-51`, `requirements.md:75-80`.
- **Evidence Required:** Service tests prove analysis and synthesis orchestration through fakes; command tests prove thin parsing/rendering and exit mapping; architecture review confirms no product orchestration moved into `internal/app`.

### Decision 2: Existing Prompt Builders And Applicability Rules Are The Only Prompt/Skip Sources Of Truth

- **Decision:** Single analysis must call the existing deterministic analysis prompt builder, and synthesis must call the existing deterministic synthesis prompt builder. Markdown source inapplicability is decided before runtime launch by existing applicability behavior. Directory analysis uses the source directory as workdir. Markdown document analysis embeds stripped document content through the prompt builder and uses a safe study/workspace workdir.
- **Rationale:** Requirements and PRD/TRD already define prompt behavior, Markdown document isolation, and applicability. This sprint connects those existing capabilities to runtime execution without duplicating prompt composition or frontmatter rules.
- **Study / Source Grounding:** `technical-handbook.md` supports source/document isolation through security evidence from opencode `internal/permission/permission.go:44-108`, restic `internal/options/secret_string.go:15-20`, and the handbook warning not to move workflow logic into CLI or runtime. No additional external repo source is needed for Markdown applicability itself because it is product-specific and governed by `PRD.md`, `TRD.md`, and prior sprint carry-forward requirements.
- **Trade-Offs Accepted:** Accept separate directory and Markdown workdir/prompt handling because source isolation rules differ by source kind.
- **Technical Debt / Future Impact:** Permission policy for Markdown document analysis may need tightening later, but the workdir and prompt constraints keep this sprint safe and reviewable.
- **Alternatives Rejected:** Rejected duplicating prompt text in command/runtime code; rejected invoking runtime for inapplicable Markdown pairs and expecting validation to fail; rejected using the Markdown file directory as a repository-like workdir that invites external exploration.
- **Contracts Satisfied:** Architecture, Security, CLI Surface, LLM Runtime; `requirements.md:33-38`, `requirements.md:43-47`, `requirements.md:50-51`, PRD sections 2.3 and 2.5.1, TRD sections 9A and 10.
- **Evidence Required:** Tests assert prompt builders are used, directory and Markdown workdirs differ correctly, inapplicable Markdown skips return success without runtime invocation, and synthesis excludes inapplicable Markdown reports from required inputs.

### Decision 3: Generic Agentwrap-Backed Runtime Request Mapping

- **Decision:** Study execution constructs generic platform runtime requests with prompt text, workdir, provider, model, timeout, task metadata, permission policy, required health checks, required capabilities, expected output validation, and context propagation. Metadata is flattened into safe key/value fields and includes task kind, study, dimension number/slug, source name, source kind, output path, runtime, provider, and model where applicable.
- **Rationale:** The runtime boundary is the volatile external system seam and must remain product-agnostic. The study module can translate product facts into generic request fields without making `internal/platform/runtime` import or understand study types.
- **Study / Source Grounding:** `technical-handbook.md` explicit composition and runtime-boundary evidence cites opencode `internal/app/app.go:42-81`, gh-cli `pkg/cmd/factory/default.go:26-46`, go-task `executor.go:22-24`, config precedence examples from chezmoi `internal/cmd/config.go:2253-2287`, go-task `internal/flags/flags.go:314-327`, and opencode `internal/config/config.go:609-641`, plus security evidence from opencode `internal/permission/permission.go:44-108`.
- **Trade-Offs Accepted:** Accept explicit request mapping in `internal/study` and app composition wiring instead of convenience product-aware runtime helpers.
- **Technical Debt / Future Impact:** Future batch execution may add attempt/run-state metadata, but current metadata must stay generic and safe.
- **Alternatives Rejected:** Rejected study-aware runtime helpers; rejected direct OpenCode subprocess handling; rejected parsing OpenCode stdout/stderr or human diagnostics for control flow; rejected hidden globals for runtime/config.
- **Contracts Satisfied:** Architecture, Configuration, Observability, Security, LLM Runtime, LLM Evaluation / Cost / Safety; `requirements.md:38-39`, `requirements.md:47`, `requirements.md:52`, `requirements.md:75-85`, TRD sections 6, 11, 12, 19, 20, and 22.
- **Evidence Required:** Fake-runtime tests assert request fields; review confirms `internal/platform/runtime` has no study imports; tests or code review prove event channels are consumed/drained and final result is obtained through the platform abstraction.

### Decision 4: Study-Owned Validation Is The Product Success Gate

- **Decision:** Runtime success alone is never enough. Analysis succeeds only after the expected per-source report exists, is non-empty, and passes the existing study-owned per-source report validator. Synthesis succeeds only after all applicable per-source reports pass preflight, runtime completes, and the expected final report exists, is non-empty, and passes the existing final report validator. Runtime expected-output validation may be attached generically, but study validators remain definitive.
- **Rationale:** The PRD principle is explicit: runtime success is not product success. Validation rules are study semantics and belong in `internal/study`, not `internal/platform/runtime`.
- **Study / Source Grounding:** `technical-handbook.md` validation-as-product-gate evidence cites go-task `errors/errors.go:47-50`, restic `internal/errors/fatal.go:10`, helm `pkg/storage/driver/driver.go:27-36`, opencode `internal/permission/permission.go:44-108`, restic `internal/options/secret_string.go:15-20`, and validation/provenance evidence from helm `pkg/provenance/sign.go:281-288` and restic `internal/crypto/crypto.go:309-311`.
- **Trade-Offs Accepted:** Accept a two-layer validation boundary: generic expected-output checks can inform runtime validation, while study-owned validators decide product success.
- **Technical Debt / Future Impact:** Duplicate validation messages are possible unless rendered carefully. The implementation must produce one coherent user-facing diagnostic with report path and failed checks.
- **Alternatives Rejected:** Rejected treating runtime exit success as product success; rejected moving report validators into platform runtime; rejected validating only in CLI; rejected synthesizing when any applicable per-source report is missing or invalid.
- **Contracts Satisfied:** Architecture, Errors, LLM Runtime, LLM Evaluation / Cost / Safety, Testing; `requirements.md:40-41`, `requirements.md:45-49`, `requirements.md:55`, PRD sections 1.4 and 2.6, TRD section 15.
- **Evidence Required:** Tests cover runtime success with missing output, runtime success with invalid per-source report, runtime success with invalid final report, synthesis blocked by missing reports, synthesis blocked by invalid reports, and diagnostics listing failed check names and report paths.

### Decision 5: Typed Outcomes, Safe Diagnostics, And Deterministic CLI Output

- **Decision:** Execution returns structured statuses for completed, skipped, runtime failed, validation failed, preflight blocked, and cancelled outcomes. CLI summaries and skip messages go to stdout. Runtime, validation, preflight, usage/config/workspace, and cancellation diagnostics go to stderr. Runtime classifications and cause chains are preserved through platform/agentwrap error types rather than parsing human-readable text.
- **Rationale:** Users need actionable output and scripts need deterministic behavior. The sprint must distinguish skipped inapplicable work from failure, runtime failure from validation failure, and preflight blockers from launched runtime failures.
- **Study / Source Grounding:** `technical-handbook.md` error and observability evidence cites rclone `fs/fserrors/error.go:22-29`, go-task `errors/errors.go:47-50`, restic `internal/errors/fatal.go:10`, helm `internal/logging/logging.go:31-66`, k9s `internal/slogs/keys.go:6-231`, gh-cli `pkg/iostreams/iostreams.go:52-54`, and terminal UX examples from gh-cli `pkg/iostreams/iostreams.go:116-130` and restic `internal/ui/terminal.go:10-34`.
- **Trade-Offs Accepted:** Accept less rich progress output during single execution in exchange for stable command tests and clear stream separation.
- **Technical Debt / Future Impact:** Rich runtime event presentation and stable JSON result schemas remain future work unless explicitly documented and tested in this sprint.
- **Alternatives Rejected:** Rejected streaming all runtime events to stdout; rejected returning generic error text without classification; rejected parsing runtime strings; rejected adding undocumented JSON output.
- **Contracts Satisfied:** Errors, Observability, CLI Surface, Security; `requirements.md:37`, `requirements.md:41-42`, `requirements.md:49`, `requirements.md:52-54`, `requirements.md:82-83`, TRD sections 7, 14, and 19.
- **Evidence Required:** CLI tests assert exit codes, stdout/stderr separation, skip success, validation diagnostics, runtime diagnostics, and no sensitive prompt/env/native payloads in normal output.

### Decision 6: Fake-Runtime First Verification

- **Decision:** Default study and command tests must use fake runtimes, temporary workspaces/fixtures, existing validators/prompt builders, and buffer-backed IO. Normal `go test ./...` must not require OpenCode, provider credentials, network, wall-clock sleeps, or real subprocess execution. Any real OpenCode smoke test is optional, environment-gated, skipped by default, and routed through `agentwrap/opencode` via the platform runtime boundary.
- **Rationale:** The sprint needs reviewable proof of UltraPlan-owned behavior: request mapping, skip behavior, validation gates, diagnostics, and exit mapping. Those behaviors are best verified deterministically with fakes.
- **Study / Source Grounding:** `technical-handbook.md` testing and IO evidence cites chezmoi `internal/cmd/main_test.go:64-174`, gh-cli `pkg/httpmock/stub.go:35-199`, lazygit `pkg/commands/oscommands/fake_cmd_obj_runner.go:17-26`, gh-cli `pkg/iostreams/iostreams.go:551-568`, go-task `executor.go:553-564`, and restic `internal/ui/mock.go:10-53`.
- **Trade-Offs Accepted:** Accept that default tests do not prove real OpenCode behavior end-to-end.
- **Technical Debt / Future Impact:** Real adapter drift can be covered by future gated smoke tests or release hardening.
- **Alternatives Rejected:** Rejected default real-runtime tests; rejected tests that require network/credentials/subprocesses; rejected hardcoded `os.Stdout`/`os.Stderr` or global runtime construction in command tests.
- **Contracts Satisfied:** Testing, CLI Surface, LLM Runtime; `requirements.md:55-60`, TRD section 23 and section 24.1.
- **Evidence Required:** `internal/study/run_test.go`, `internal/app/study_run_commands_test.go`, `go test ./...`, and `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`.

### Decision 7: Defer Batch, Workflow, Repair, Summary, Code Extraction, And Target/Sprint Scope

- **Decision:** Do not introduce worker pools, task matrices, durable run-loop mutation, stale task recovery, retry scheduling across process restarts, per-study locks, automatic synthesis, summary updates, code-reference extraction, product repair loops, durable agentwrap run-store persistence, stable public execution JSON schemas, target workflows, or sprint workflows. Keep single execution compatible with existing task/status vocabulary but do not build the later workflow engine now.
- **Rationale:** The sprint exists to connect one analysis and one synthesis to the runtime and validation gates. Adding batch/workflow scope would obscure acceptance criteria and contradict selected non-goals.
- **Study / Source Grounding:** `technical-handbook.md` bounded-abstraction and philosophy evidence cites lazygit `VISION.md:97`, rclone `fs/types.go:16-59`, restic `internal/backend/backend.go:19-90`, helm `internal/plugin/runtime_subprocess.go:65-79`, rclone `fs/registry.go:407`, and go-task `executor.go:20-24`. Concurrency evidence from restic `internal/repository/repository.go:567` and opencode `cmd/root.go:130,252` supports event cleanup without adding general scheduling machinery.
- **Trade-Offs Accepted:** Accept future integration work for `run-all`/`run-loop` rather than premature generalization now.
- **Technical Debt / Future Impact:** Later workflow sprints may refactor shared execution helpers or expand result/state models.
- **Alternatives Rejected:** Rejected general scheduler/workflow engine; rejected broad plugin/registry systems; rejected adding automatic synthesis after `study run`; rejected product repair prompts and loops; rejected stable JSON schema unless explicitly implemented and tested.
- **Contracts Satisfied:** Architecture, Testing, LLM Runtime, LLM Evaluation / Cost / Safety; `requirements.md:62-73`, `requirements.md:78-80`, `requirements.md:86`, sprint-index excluded context lines 81-99.
- **Evidence Required:** Code review confirms no scheduler/worker pool/run-loop/summary/codeextract/target/sprint implementation entered this sprint; tests focus on single analysis and single synthesis only.

## Expected Evidence

| Evidence Type | Required Evidence | Source / Command / Review Check |
| --- | --- | --- |
| Study service tests | Analysis success, synthesis success, runtime failure, runtime success with missing output, runtime success with invalid output, Markdown inapplicable skip with zero runtime invocations, synthesis blocked by missing reports, and synthesis blocked by invalid reports. | `/home/antonioborgerees/coding/ultraplan-go/internal/study/run_test.go` |
| Runtime request mapping tests | Prompt text, workdir, provider, model, timeout, metadata, permission policy, health/caps, expected output validation, and context path are passed through fake runtime. | `/home/antonioborgerees/coding/ultraplan-go/internal/study/run_test.go` and code review of `internal/study/run.go` and `internal/study/synthesize.go` |
| CLI command tests | `study run` and `study synthesize` argument behavior, help text, success output, skip output, validation diagnostics, runtime diagnostics, exit statuses, and stdout/stderr separation. | `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_run_commands_test.go` |
| Architecture review | Study execution behavior in `internal/study`; CLI thin in `internal/app`; platform runtime generic; no product imports in `internal/platform/runtime`; no direct OpenCode subprocess handling. | `system/protocols/architecture-review-protocol.md` checks during `review.md` |
| Runtime safety review | No prompt bodies, embedded Markdown bodies, secrets, full environment values, or unsafe native payloads in normal diagnostics; runtime classifications preserved through platform/agentwrap types. | Code review of diagnostics and error mapping; command tests for output content |
| Default verification | Normal tests and build pass without OpenCode/provider credentials/network/subprocess dependence. | `go test ./...` and `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go` |
| Optional smoke evidence | If real runtime smoke tests are added, they are environment-gated, skipped by default, and use `agentwrap/opencode` through `internal/platform/runtime`. | Test skip conditions and review notes in `review.md` |
| Sprint review | Acceptance criteria, constraints, non-goals, and carry-forward risks are checked and recorded. | `system/protocols/sprint-review-protocol.md` and `projects/ultraplan-go/sprints/10-single-analysis-synthesis/review.md` |

## Assumptions And Risks

| Item | Type | Impact | Mitigation / Follow-Up |
| --- | --- | --- | --- |
| Existing prompt builders and validators are ready to call from execution code. | Assumption | If APIs are not ready, implementation may need small adapter changes. | Inspect existing APIs before coding; keep changes local to `internal/study`; do not duplicate prompt or validation logic. |
| Platform runtime can accept generic metadata, expected output validation, and fakeable execution without study imports. | Assumption | If not true, request mapping may require platform runtime additions. | Add only generic runtime fields or narrow consumer-owned interfaces; reject product-aware runtime changes. |
| Existing exit-code conventions can represent validation, runtime, usage, config/workspace, and cancellation outcomes. | Assumption | If conventions are incomplete, CLI mapping may become inconsistent. | Reuse existing constants where possible; add minimal local mapping only if needed and cover with command tests. |
| Runtime expected-output validation and study post-run validation may double-report failures. | Risk | Users may see noisy or conflicting diagnostics. | Render study validation as definitive product result and include runtime validation only as safe supporting diagnostics when useful. |
| Synthesis preflight may accidentally require reports for inapplicable Markdown sources. | Risk | Valid synthesis could be blocked incorrectly. | Always compute applicable sources first and add targeted tests for inapplicable Markdown exclusion. |
| Markdown document analysis could accidentally encourage filesystem exploration. | Risk | Source isolation and security requirements could be violated. | Use existing Markdown prompt builder with stripped content, safe workdir, restrictive permission posture, and request-mapping tests. |
| Runtime event channel handling could block final result retrieval. | Risk | Single execution may hang or leak goroutines. | Consume or drain runtime events per platform abstraction and always obtain final result through runtime `Wait`/equivalent. |
| Fake runtime may diverge from real agentwrap/OpenCode behavior. | Risk | Default tests may miss real adapter issues. | Keep fake tests focused on UltraPlan-owned behavior; add optional environment-gated smoke tests if implementation needs real-runtime confidence. |
| Future batch/run-loop may need different state fields. | Risk | Some execution result code may be refactored later. | Keep current result fields minimal and compatible with existing task/status vocabulary; do not persist a new batch schema. |

## Implementation Constraints

- `internal/study` owns resolution use, source applicability, prompt builder calls, output path decisions, synthesis preflight, runtime request construction from study facts, and report validation gates.
- `internal/app` must remain a thin adapter: parse arguments, call study service methods, render deterministic output, and map errors/results to exit statuses.
- `internal/platform/runtime` must not import `internal/study` or learn study, dimension, source, report, Markdown applicability, synthesis gating, or validation semantics.
- Runtime execution must go through the existing platform runtime and agentwrap/OpenCode boundary; UltraPlan must not implement direct OpenCode subprocess supervision, native event decoding, permission translation, policy retry/fallback, or runtime string parsing.
- Directory source analysis uses the source directory as runtime workdir; Markdown document analysis uses a safe study/workspace workdir and embedded stripped document content; synthesis uses a safe study/workspace workdir.
- Inapplicable Markdown source/dimension analysis is a successful skipped result and must not invoke runtime.
- Synthesis must fail before runtime launch when any applicable per-source report is missing or invalid, and diagnostics must identify every blocking path and failed check where available.
- Runtime success is insufficient for product success; expected artifacts must exist, be non-empty, and pass the appropriate study-owned validator.
- Runtime request metadata must be safe flattened values, not raw study structs or unsafe native payloads.
- Diagnostics must preserve cause chains and runtime classifications through typed/platform errors and must not expose prompts, embedded document content, secrets, full environment values, or unsafe native payloads in normal output.
- Default tests must use fake runtimes and buffer-backed IO; they must not require OpenCode, provider credentials, network access, wall-clock sleeps, or real subprocesses.
- Do not add stable public JSON output unless it is explicitly documented and tested in this sprint.
- Do not implement batch/run-loop, summary, code extraction, repair workflow, target, or sprint execution scope in this sprint.

## Plan Handoff

`plan.md` must execute these decisions. It must not invent architecture, scope, or decisions beyond this document.

The plan must carry forward:

- Final decisions 1-7.
- Applicable contracts and requirement mappings in `Contracts Applied`.
- Expected evidence in `Expected Evidence`.
- Risks and assumptions in `Assumptions And Risks`.
- Required architecture and sprint review protocols selected by `sprint-index.md`.
- The constraint that implementation must be able to proceed without reopening package ownership, validation boundaries, runtime request ownership, CLI behavior, or deferred scope.

## Phase Exit Criteria

- [x] Selected context was read and used.
- [x] Area-specific reasoning documents were completed where applicable or explicitly marked not applicable.
- [x] Area-specific reasoning conclusions are reflected or explicitly overridden in final decisions.
- [x] Contracts are explicitly mapped to decisions and expected evidence.
- [x] Final decisions are clear enough for `plan.md` to execute.
- [x] Expected evidence is specific and reviewable.
