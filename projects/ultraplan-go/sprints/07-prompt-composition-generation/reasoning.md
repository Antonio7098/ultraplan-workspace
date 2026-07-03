# Sprint Reasoning: 07 Prompt Composition Generation

> Project: `ultraplan-go`
> Sprint: `07-prompt-composition-generation`
> Output: `projects/ultraplan-go/sprints/07-prompt-composition-generation/reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/sprints/07-prompt-composition-generation/requirements.md`, `projects/ultraplan-go/sprints/07-prompt-composition-generation/reasoning.md`, `projects/ultraplan-go/sprints/07-prompt-composition-generation/sprint-index.md`, `projects/ultraplan-go/sprints/07-prompt-composition-generation/technical-handbook.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/project-index.md`, `templates/sprint-reasoning.md`; no `projects/ultraplan-go/sprints/07-prompt-composition-generation/reasoning/*.md` files were present.

This document decides. It synthesizes selected context, handbook evidence, area-specific reasoning, and contracts into final sprint decisions.

It does not replace `sprint-index.md`, `technical-handbook.md`, or `reasoning/*.md`.

## Sprint Purpose

- **Goal:** Implement deterministic study prompt composition for directory analysis, Markdown document analysis, and synthesis, including manifest-backed prompt previews, without invoking agentwrap, OpenCode, provider credentials, network access, or subprocess runtime execution.
- **Non-Goals:** Runtime-backed analysis or synthesis execution, agentwrap/OpenCode integration, runtime health checks, permission policy mapping, event/log mapping, retry/fallback/repair/cancellation/provider cost metadata, `run-all`, `run-loop`, durable task state, status calculations, worker pools, resumability, report validation changes beyond consuming prior path concepts, summary generation, code reference extraction, assisted study initialization, source cloning, Markdown source discovery changes, target workflows, sprint planning, and sprint execution.
- **Depends On:** Prior sprint capabilities for CLI/app composition, workspace discovery and path safety, study/source/dimension resolution, generated study layout, Markdown source kind/frontmatter/applicability behavior, and deterministic report path/validation concepts.

## Selected Context And Pre-Reasoning Artifacts

| Artifact | Path | How It Was Used |
| --- | --- | --- |
| Project Index | `projects/ultraplan-go/project-index.md` | Confirmed repository target, available docs, active contract pool, selected study evidence pool, and project non-goals. |
| Sprint Requirements | `projects/ultraplan-go/sprints/07-prompt-composition-generation/requirements.md` | Treated as the binding contract for outputs, acceptance criteria, constraints, dependencies, and review expectations. |
| Product Requirements | `projects/ultraplan-go/docs/PRD.md` | Grounded product behavior for prompt composition, dry-run previews, directory vs Markdown source semantics, artifact readability, deterministic paths, and deferred target/sprint scope. |
| Technical Requirements | `projects/ultraplan-go/docs/TRD.md` | Grounded prompt input requirements, deterministic prompt builder expectations, manifest/dry-run behavior, Markdown applicability/frontmatter rules, test requirements, path handling, and no-runtime test isolation. |
| Architecture | `projects/ultraplan-go/docs/ARCHITECTURE.md` | Governed package ownership: study behavior stays in `internal/study`, CLI wiring stays in `internal/app`, workspace remains path-owner, and platform/runtime stays generic. |
| Sprint Index | `projects/ultraplan-go/sprints/07-prompt-composition-generation/sprint-index.md` | Selected applicable contracts, evidence reports, non-goals, review protocols, excluded context, and the intended prompt-composition reasoning focus. |
| Technical Handbook | `projects/ultraplan-go/sprints/07-prompt-composition-generation/technical-handbook.md` | Supplied study evidence and trade-offs for thin CLI wiring, domain-owned prompt logic, DI, IO abstraction, actionable errors, deterministic tests, security boundaries, and avoiding speculative performance machinery. |
| Sprint Reasoning Template | `templates/sprint-reasoning.md` | Provided the required structure for final decisions, evidence, risks, constraints, and plan handoff. |

## Area-Specific Reasoning Inputs

The sprint index selected an Architecture reasoning artifact at `projects/ultraplan-go/sprints/07-prompt-composition-generation/reasoning/prompt-composition.md`, but no `reasoning/*.md` files were present. Final decisions therefore rely directly on `requirements.md`, project docs, `sprint-index.md`, and `technical-handbook.md`.

| Area | Reasoning Document | Key Conclusion | Evidence Basis | Impact On Final Decision |
| --- | --- | --- | --- | --- |
| Prompt composition architecture | Not present | No area-specific conclusion file exists to summarize. The sprint still has enough selected context and handbook evidence to decide implementation boundaries. | `sprint-index.md` selected the area; `technical-handbook.md` and project docs provide the usable evidence. | Decisions explicitly resolve prompt ownership, template boundaries, source-kind behavior, synthesis manifest shape, preview boundaries, and test strategy so `plan.md` does not reopen architecture. |

## Sprint Technical Handbook Summary

- **Relevant Patterns:** Keep study prompt logic behind thin CLI wiring; use factory-style command construction with explicit dependencies; load deterministic local inputs lazily; inject IO for preview output and tests; use structured actionable error paths; protect deterministic text artifacts with fixture/golden-style tests; keep prompt preview as a no-runtime trust boundary.
- **Important Trade-Offs:** Keep prompt composition in `internal/study` instead of a global prompt package; use a thin preview command that delegates to study builders; prefer deterministic full-output tests where stable and focused substring/manifest tests where wording churn would be excessive; prefer workspace-relative manifest paths while preserving absolute paths only for diagnostics; fail fast on missing inputs instead of producing partial runnable-looking prompts; avoid reusing future runtime execution paths for preview.
- **Warnings / Anti-Patterns:** Do not let command handlers become prompt composition engines; do not introduce hidden globals for templates/config/IO; do not bypass injectable output streams; do not break error chains or rely on string matching; do not treat Markdown document sources like repositories; do not add speculative concurrency, caching, pooling, or runtime machinery; do not assert brittle private implementation details in tests.
- **Evidence Confidence:** High for package/command/DI/error/IO/testing decisions because selected reports cite multiple mature Go CLIs and concrete source examples. Medium for security/performance/philosophy decisions because those sources are broader and adapted to this no-runtime preview sprint rather than copied directly.

## Contracts Applied

Requirement IDs below are local aliases for this reasoning document because `requirements.md` lists acceptance criteria and constraints without explicit IDs.

| Contract / Requirement ID | Constraint | Decision Impact | Expected Evidence |
| --- | --- | --- | --- |
| Architecture | Product behavior stays module-owned; no global prompt/template/report package; platform packages must not import `study`. | Prompt request/result types, manifest construction, template loading, source-kind rules, and synthesis report selection stay in `internal/study`; command parsing/output stays in `internal/app`. | Import/package review, `go test ./...`, `go build ./cmd/ultraplan`, and architecture review checks. |
| Errors | Failures must be actionable and preserve cause chains. | Missing templates, documents, reports, unknown refs, ambiguous refs, and inapplicable single-source previews return contextual errors naming the failed path/source/dimension/report. | Unit and command tests asserting diagnostic content and non-zero command behavior. |
| Security | Preview must not invoke subprocesses, OpenCode, providers, network access, or repository-wide scans; generated artifacts must avoid secrets and unnecessary absolute paths. | Prompt builders only read known workspace-managed files and explicit source docs/reports; preview directly renders prompt text/manifest without runtime paths. | Command tests proving no runtime invocation path; code review for absent agentwrap/OpenCode/subprocess/network calls in preview. |
| Testing | Deterministic unit, fixture, and command coverage is required. | Tests cover directory prompts, Markdown prompts, frontmatter stripping, applicability, synthesis manifests, missing inputs, deterministic ordering, preview success/failure, and no-runtime behavior. | `internal/study/prompts_test.go`, `internal/app/study_prompt_commands_test.go`, `go test ./...`. |
| Documentation | Generated sprint artifacts and CLI help must be maintainable and reviewable. | Preview command help and sprint artifacts must describe prompt-only behavior and avoid implying execution. | Help/output review and sprint review checks. |
| CLI Surface | Preview command must be script-friendly with stdout or explicit output path behavior and meaningful failures. | Add a minimal prompt preview or dry-run command surface in `internal/app` without committing to future runtime command semantics beyond the sprint scope. | Command tests for stdout, explicit output path, usage/failure cases, and exit behavior. |
| AC-01, AC-02 | Builders load required workspace prompt/template files deterministically and fail before runtime-facing behavior when files are missing. | Implement required template reads from resolved workspace paths and validate all required files before returning prompts. | Missing template tests naming `prompts/base.md`, `prompts/synthesize.md`, `templates/repo-analysis.md`, or final report template path. |
| AC-03, AC-04 | Directory prompts include required content and hard source-isolation/citation rules. | Directory prompt composition uses base prompt, dimension Markdown, report template, study/source/output metadata, and explicit file-path-line citation instructions unless disabled by dimension. | Directory prompt unit/golden test. |
| AC-05, AC-06, AC-07, AC-08 | Markdown prompts embed stripped document content, exclude frontmatter, forbid external exploration, and default away from code citation requirements. | Markdown prompt composition uses the existing frontmatter stripping/applicability behavior and document-only prompt section. | Markdown prompt and frontmatter stripping tests. |
| AC-09 | Inapplicable Markdown source/dimension pairs are skipped or reported without producing runnable analysis prompts. | Single-source preview returns a clear non-runnable/inapplicable error or result; synthesis report manifests filter by applicability under AC-11. | Applicability tests for analysis preview and synthesis manifest. |
| AC-10, AC-11 | Synthesis prompts include synthesis instructions, dimension content, final template, deterministic applicable report manifest, and expected final output path. | Synthesis builder gathers applicable per-source report paths in stable order and fails on missing applicable reports without inventing them. | Synthesis prompt/manifest tests with sorted expected paths and missing report failure. |
| AC-12, AC-13 | Builders return deterministic prompt text plus input manifest with prompt kind, study, dimension, source/source kind where applicable, template paths, input paths, and expected output path. | Define manifest structs and stable path rendering in `internal/study/domain.go`; sort all list fields deterministically. | Manifest unit tests and deterministic repeated-build test. |
| AC-14, AC-15 | Preview writes stdout or explicit output path without runtime; errors are actionable. | Preview command accepts prompt selection refs and output option, writes rendered output only, and never calls run/synthesize runtime paths. | Command tests for success, failure, output file/stdout, and no-runtime behavior. |
| AC-16, AC-17, AC-18, AC-19 | Required unit/command test coverage, `go test ./...`, and `go build ./cmd/ultraplan`. | Plan must include both focused tests and whole-repo verification. | Passing test/build commands recorded in review. |

## Repos Studied / Source Evidence Used

| Source / Repo / Report | Concrete Reference | Relevant Finding | Why It Matters For This Sprint | Used In Decision(s) |
| --- | --- | --- | --- | --- |
| `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md`; examples `cmd/age/age.go:105`, `cmd/gh/main.go:6`, `cmd/root.go:112-127`, `cmd/evaluate_sequence_command.go:7` | Mature Go CLIs keep entrypoints thin and business logic in domain/protected packages with one-way imports. | Supports `internal/study` ownership and app-only CLI wiring. | Decisions 1, 6 |
| `02-command-architecture` | `studies/go-cli-study/reports/final/02-command-architecture.md`; examples `pkg/cmd/install.go:132-145`, `pkg/cmdutil/factory.go:16-43`, `pkg/cmd/issue/list/list.go:47-118`, caution `cmd/root.go:49-183` | Command factories and thin handlers improve testability and prevent business logic accumulation in commands. | Supports a preview command that delegates prompt assembly to study builders. | Decisions 1, 6 |
| `03-dependency-injection` | `studies/go-cli-study/reports/final/03-dependency-injection.md`; examples `internal/ghcmd/cmd.go:52-132`, `pkg/cmd/factory/default.go:26-46`, `executor.go:22-24`, `fs/cache/cache.go:16-21` | Manual composition, explicit dependencies, and minimal globals are favored. | Supports passing workspace roots, stores, and IO explicitly without global prompt registries. | Decisions 1, 6 |
| `05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md`; examples `pkg/storage/driver/driver.go:27-36`, `go-task/errors/errors_task.go:13-32`, `gh-cli/internal/ghcmd/cmd.go:281-301` | Wrapped errors, typed/sentinel categories, and user-facing diagnostics make CLI failures actionable. | Supports fail-fast missing input checks and diagnostics naming failed paths/refs. | Decisions 2, 3, 4, 5, 6 |
| `06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md`; examples `pkg/iostreams/iostreams.go:551-568`, `executor.go:553-564`, `internal/chezmoi/system.go:25`, `internal/fs/interface.go:10-31` | Injectable readers/writers and filesystem seams make command output and file behavior testable. | Supports stdout/output-path preview tests without direct `os.Stdout` coupling. | Decision 6 |
| `11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md`; examples `helm/internal/test/test.go:43`, `go-task/task_test.go:166-169`, `chezmoi/internal/cmd/main_test.go:64-174` | Table-driven, fixture, golden, and command-level tests protect deterministic CLI artifacts. | Supports deterministic prompt/manifest tests and command tests. | Decision 7 |
| `13-security` | `studies/go-cli-study/reports/final/13-security.md`; examples `internal/llm/tools/bash.go:41-55`, `internal/permission/permission.go:44-108`, `pkg/registry/transport.go:37-41` | Explicit trust boundaries, command filtering, permission boundaries, and secret redaction matter around untrusted inputs and runtimes. | Supports strict no-runtime preview, workspace-relative manifests, and document-only Markdown rules. | Decisions 2, 3, 5, 6 |
| `14-performance` | `studies/go-cli-study/reports/final/14-performance.md`; examples `gh-cli/pkg/cmdutil/factory.go:27-42`, `age/internal/stream/stream.go:20`, `yq/pkg/yqlib/stream_evaluator.go:78-113` | Avoid unnecessary startup work, prefer lazy initialization, and avoid large accidental loads. | Supports deterministic targeted reads only, no repo-wide scans, and no speculative caching/pooling. | Decisions 2, 3, 5, 6 |
| `15-philosophy` | `studies/go-cli-study/reports/final/15-philosophy.md`; examples `pkg/cmd/factory/default.go:26-46`, `internal/chezmoi/system.go:25-45`, `VISION.md:97` | Complexity should be deliberate; abstractions should serve clear product needs. | Supports minimal structs/builders and deferring runtime/durable workflow concerns. | Decisions 1, 7 |

## Trade-Off And Debt Analysis

### Accepted Trade-Offs

| Trade-Off | Benefit | Cost / Constraint Accepted | Why Acceptable Now | Revisit Trigger |
| --- | --- | --- | --- | --- |
| Keep prompt composition in `internal/study` instead of creating `internal/prompts` or `internal/templates` | Preserves product ownership and dependency direction from architecture docs and study evidence. | `internal/study` gains more focused files and local template/report concepts. | Prompt behavior is study-specific and not shared with another current product module. | A second product module needs the same prompt/template behavior for non-study state. |
| Use a thin preview command that calls study prompt builders | Keeps command tests simple and avoids prompt business logic in `internal/app`. | Requires request/result structs crossing from app to study. | The indirection is small and aligns with existing Go CLI evidence. | Preview command grows into a broad UX surface requiring separate app service orchestration. |
| Fail fast on missing templates/documents/reports instead of partial prompts | Prevents misleading runnable-looking prompts and protects deterministic manifests. | Users cannot inspect partial output when one input is absent. | The sprint goal is correctness and deterministic preview, not best-effort explanation generation. | Product adds an explicit diagnostic-only mode for incomplete workspaces. |
| Use workspace-relative paths in manifests by default | Keeps artifacts relocatable, reviewable, and safer to share. | Some diagnostics may need resolved absolute paths for local troubleshooting. | Requirements allow absolute paths only where needed for diagnostics. | Users need stable machine-readable absolute paths for automation and explicitly opt in. |
| Treat Markdown document prompts as a separate source-kind branch | Enforces document-only trust boundary and avoids repository/citation leakage. | Some prompt builder code branches by source kind instead of one fully generic path. | Directory and Markdown sources have materially different rules and validation expectations. | More source kinds are added and a shared source-kind strategy abstraction becomes worthwhile. |
| Require deterministic manifest ordering for synthesis | Makes prompt previews and tests stable; prevents missing report ambiguity. | Sorting and applicability filtering must be explicit in builders. | Synthesis correctness depends on knowing exactly which reports are required. | Report path rules change or synthesis begins accepting dynamic optional inputs. |
| Avoid runtime/durable workflow reuse in preview | Proves no-runtime behavior and avoids side effects. | Later runtime commands may duplicate some reference resolution before sharing prompt builders. | Runtime execution is explicitly a non-goal and preview must be safe without credentials. | Runtime-backed `run --dry-run` enters scope and needs unified preflight behavior. |

### Potential Technical Debt

| Debt / Shortcut | Why It Might Accrue | Current Mitigation | Owner / Follow-Up |
| --- | --- | --- | --- |
| Manifest schema may become de facto stable before runtime manifests are designed | Tests and users may start depending on field names and path formats. | Keep manifest minimal, explicit, workspace-relative, and directly tied to acceptance criteria. | Future runtime execution sprint should either adopt or version this manifest shape. |
| Citation-disable detection may be wording-based if dimensions lack structured metadata | Requirements mention disabling citation requirements but current dimension model only has Markdown fields. | Define a narrow helper/rule in study code and cover it with tests; do not spread string checks through command code. | Future dimension schema/versioning work should make citation policy explicit if needed. |
| Preview CLI shape may need adjustment when actual `run`/`synthesize` execution is implemented | A standalone preview command could diverge from future dry-run semantics. | Keep preview command minimal and backed by reusable study builders, not runtime command internals. | Runtime execution sprint should decide whether preview remains separate or becomes shared dry-run plumbing. |
| Prompt text tests can be noisy if templates evolve | Full golden assertions may require broad updates for wording-only changes. | Combine stable manifest assertions with focused required-text checks and golden tests where whole text is intentionally stable. | Test maintenance during template/content revisions. |
| Error categories may be less formal than future stable JSON errors | Current sprint needs actionable errors, not a public structured error schema. | Preserve causes and expose enough typed/sentinel handling for command tests. | Stable JSON/exit-code sprint or public automation contract. |

### Future Considerations

| Consideration | Deferred Until | Reason Deferred | What Should Be Preserved Now |
| --- | --- | --- | --- |
| Agentwrap/OpenCode runtime request construction | Runtime execution sprint | This sprint composes and previews prompts only. | Prompt builders return text and manifest that future runtime request construction can consume. |
| Runtime validation and repair prompts | Runtime/validation sprint | Repair and validation wrapping are non-goals here. | Prompt outputs must include expected report paths and source kind metadata. |
| Durable run state, run-loop, and resumability | Batch/durable workflow sprint | Run-loop state is explicitly out of scope. | Deterministic applicable source/report ordering and manifest fields should align with later task state. |
| Summary generation and inapplicable-cell representation | Summary sprint | Summary CSV writing is non-goal. | Applicability filtering must be consistent so later summary logic can reuse it. |
| Code reference extraction | Code extraction sprint | This sprint only requires prompts to request citations where applicable. | Directory prompts should consistently request file path and line citations. |
| Additional prompt/template package | Only if multiple product modules need shared behavior | Current prompt behavior is study-specific. | Keep helpers cohesive and avoid leaking product behavior into platform packages. |
| Stable JSON output for preview manifests | CLI automation contract sprint | Current acceptance criteria require manifests but not a stable public JSON schema. | Use deterministic struct fields and tests so future JSON support can be added without changing semantics. |

## Final Decisions

### Decision 1: Study Owns Prompt Composition

- **Decision:** Define prompt request, prompt kind, prompt manifest, template reference, and preview result concepts in `internal/study/domain.go`, and implement prompt loading/composition in `internal/study/prompts.go`. CLI code in `internal/app` will only resolve command arguments, call study prompt builders, and render stdout/file output.
- **Rationale:** Prompt composition transforms study state: study, dimension, source, source kind, report paths, applicability, and template inputs. The project architecture explicitly says study owns prompts and reports, while the sprint constraints forbid global `internal/prompts`, `internal/templates`, `internal/reports`, or `internal/validation` packages for study-owned behavior.
- **Study / Source Grounding:** `technical-handbook.md` patterns for study-owned prompt logic and thin CLI wiring; `01-project-structure` evidence from `cmd/age/age.go:105`, `cmd/gh/main.go:6`, and one-way imports; `02-command-architecture` evidence from `pkg/cmd/install.go:132-145` and `pkg/cmdutil/factory.go:16-43`; `15-philosophy` evidence that deliberate abstractions should serve product needs.
- **Trade-Offs Accepted:** `internal/study` grows a prompt-specific implementation file, but ownership stays clear and avoids global technical-layer fragmentation.
- **Technical Debt / Future Impact:** Future runtime execution should reuse the study prompt builders rather than reimplementing prompt text assembly. If another module needs shared prompt mechanics, that later sprint can extract only genuinely common infrastructure.
- **Alternatives Rejected:** A global `internal/prompts` package is rejected because it violates architecture guidance and would separate prompt behavior from study state. Building prompt strings inside `internal/app` command handlers is rejected because it creates long command logic and makes domain tests weaker. Putting prompt behavior under `platform/runtime` is rejected because runtime must not know studies, dimensions, sources, reports, or synthesis semantics.
- **Contracts Satisfied:** Architecture, Testing, Documentation, AC-01 through AC-13, constraints from `requirements.md` lines 64-67.
- **Evidence Required:** Package/import review, `internal/study/prompts_test.go`, `internal/app/study_prompt_commands_test.go`, `go test ./...`, `go build ./cmd/ultraplan`.

### Decision 2: Deterministic Manifest-Backed Prompt Results

- **Decision:** Every prompt builder returns a result containing composed prompt text and a manifest. The manifest records prompt kind, study, dimension, source and source kind when applicable, template paths, input document/report paths, expected output path, and applicable source/report entries in deterministic order using workspace-relative paths where possible.
- **Rationale:** Requirements demand deterministic prompt text and manifest output for identical filesystem inputs. The manifest also becomes the reviewable evidence that prompt composition used the correct templates, source material, report inputs, and output path without invoking runtime behavior.
- **Study / Source Grounding:** `technical-handbook.md` design pressure says manifests are both user-facing preview artifacts and test oracles. `13-security` evidence supports avoiding unnecessary absolute local details. `14-performance` evidence supports targeted deterministic local reads and avoiding unnecessary scans.
- **Trade-Offs Accepted:** Manifest fields may become a stable surface earlier than the full runtime/task-state schema, but keeping them minimal and requirement-driven limits premature commitment.
- **Technical Debt / Future Impact:** A later runtime sprint should decide whether to version this manifest for public JSON automation or keep it as an internal preview structure. Path rendering must remain consistent with workspace path helpers.
- **Alternatives Rejected:** Returning only prompt text is rejected because it fails manifest acceptance criteria and weakens reviewability. Returning absolute paths everywhere is rejected because project docs prefer workspace-relative generated artifacts. Returning unordered maps/slices is rejected because it would make outputs non-deterministic.
- **Contracts Satisfied:** Architecture, Security, Testing, AC-10 through AC-13, constraints from `requirements.md` lines 71-72.
- **Evidence Required:** Manifest field tests, repeated-build determinism tests, synthesis ordering tests, review confirming workspace-relative paths except local diagnostics.

### Decision 3: Directory Analysis Prompts Preserve Repository Isolation And Code Citations

- **Decision:** Directory analysis prompt composition includes the shared base prompt, selected dimension Markdown, `templates/repo-analysis.md`, study name, source name, source path, expected per-source output path, hard source-isolation rules, and file-path-and-line citation requirements unless the selected dimension explicitly disables code citation requirements.
- **Rationale:** Directory sources are code/repository inputs and must preserve the product principle that architectural claims should cite concrete source references where applicable. Shared base prompt content cannot weaken source isolation, so directory-specific hard rules must be explicit in the composed prompt.
- **Study / Source Grounding:** PRD prompt composition requirements at lines 205-210 and product principle at line 47; TRD prompt inputs at lines 552-586 and prompt builder requirements at lines 915-925; `technical-handbook.md` warnings against weakening source isolation and treating source repository content as trusted. `13-security` supports explicit trust boundaries.
- **Trade-Offs Accepted:** Directory prompts carry explicit repeated guardrails even if the base prompt already contains related language. This duplication is acceptable because source isolation is safety-critical.
- **Technical Debt / Future Impact:** Citation-disable detection may be limited by existing dimension metadata. Implementation must keep any detection helper narrow and tested, and future dimension schema work may formalize citation policy.
- **Alternatives Rejected:** Relying only on `prompts/base.md` for source-isolation and citation rules is rejected because base content may be too broad or later edited. Using Markdown document citation defaults for directory sources is rejected because it would weaken code auditability. Recursively scanning a source repository during prompt composition is rejected because this sprint only composes prompts and must not perform repo-wide scans.
- **Contracts Satisfied:** Security, Testing, AC-03, AC-04, AC-14, constraints from `requirements.md` lines 68 and 70.
- **Evidence Required:** Directory prompt test asserting base/dimension/template/metadata/output/rules inclusion, citation-rule test including disabled-citation dimension case if current dimension fixtures can represent it, and code review confirming no repository scan/runtime invocation.

### Decision 4: Markdown Document Prompts Are Embedded-Document Only

- **Decision:** Markdown document analysis prompt composition includes the shared base prompt scoped for document analysis, selected dimension Markdown, report template, study/source/output metadata, and the stripped Markdown document body. It excludes YAML frontmatter, states all source material is embedded, forbids external files/repositories/code inspection, and states code citation requirements do not apply by default unless the dimension says otherwise.
- **Rationale:** Markdown document sources are not repositories. The prior Markdown source behavior provides frontmatter stripping and applicability filtering, and this sprint must consume those capabilities rather than redefining discovery or scheduling.
- **Study / Source Grounding:** PRD Markdown source requirements at lines 279-287 and 460-495; TRD source discovery/frontmatter/applicability requirements at lines 479-550 and prompt requirements at lines 892-904; `technical-handbook.md` warning not to treat Markdown document sources like repositories; `13-security` trust-boundary evidence.
- **Trade-Offs Accepted:** Markdown prompts branch from directory prompt behavior and include explicit document-only instructions even when base prompt language may discuss code analysis. This avoids accidental filesystem exploration.
- **Technical Debt / Future Impact:** If future dimensions require document section citations or other non-code citation rules, prompt builder policy may need structured dimension metadata. For now, it must not require code citations by default.
- **Alternatives Rejected:** Reusing directory prompt rules for Markdown files is rejected because it would invite external exploration and incorrect code citation requirements. Embedding frontmatter in the prompt is rejected because frontmatter is scheduling metadata, not analysis content. Treating an inapplicable Markdown pair as a runnable prompt with a warning is rejected because product rules say inapplicable pairs are skipped, not failed or executed.
- **Contracts Satisfied:** Security, Errors, Testing, AC-05 through AC-10, AC-14, constraints from `requirements.md` lines 69 and 73-74.
- **Evidence Required:** Markdown prompt tests asserting stripped body inclusion, frontmatter exclusion, document-only instruction text, no-code-citation default text, inapplicable pair behavior, and missing Markdown file diagnostics.

### Decision 5: Synthesis Uses Applicable Report Manifests Only

- **Decision:** Synthesis prompt composition includes `prompts/synthesize.md`, selected dimension Markdown, the final report template, expected final output path, and a deterministic manifest of existing applicable per-source report paths. Directory sources are always applicable; Markdown sources are included only when applicable to the selected dimension. Missing applicable reports fail composition; inapplicable Markdown pairs are excluded and do not create missing-report failures.
- **Rationale:** Synthesis should compare only reports that should exist for the selected dimension. Inventing missing reports or requiring inapplicable Markdown reports would violate prior applicability behavior and produce misleading synthesis prompts.
- **Study / Source Grounding:** Requirements acceptance criteria for synthesis manifests; PRD synthesis scenario and inapplicable Markdown rules at lines 92-94, 223-235, and 492-495; TRD applicability and prompt input rules at lines 523-550 and 905-912; `technical-handbook.md` design pressure around deterministic applicable report manifests.
- **Trade-Offs Accepted:** Synthesis prompt preview becomes strict: one missing applicable input blocks prompt generation. This is preferable to generating a prompt that cannot produce a valid final report.
- **Technical Debt / Future Impact:** Later runtime synthesis may distinguish invalid reports from missing reports, but this sprint only needs existing report-path inputs and missing-input failures. The deterministic ordering should carry forward into run-loop synthesis gating.
- **Alternatives Rejected:** Including all sources regardless of applicability is rejected because it would require nonexistent Markdown reports. Silently omitting missing applicable reports is rejected because synthesis would hide incomplete analysis. Sorting by filesystem discovery order is rejected because it can vary and break deterministic output.
- **Contracts Satisfied:** Errors, Security, Testing, AC-10 through AC-13, constraints from `requirements.md` lines 71 and 74.
- **Evidence Required:** Synthesis prompt tests with directory and Markdown sources, deterministic sorted manifest assertions, missing applicable report failure tests, and inapplicable Markdown exclusion tests.

### Decision 6: Preview Command Is Runtime-Free Rendering Only

- **Decision:** Add CLI wiring in `internal/app/study_prompt_commands.go` for a prompt preview or dry-run command that resolves study/dimension/source references, calls study prompt builders, and prints to stdout or writes to an explicit output path. It must not call analysis/synthesis runtime execution, agentwrap, OpenCode, provider config, network access, or subprocess execution.
- **Rationale:** The sprint goal is prompt composition and preview before runtime integration. A runtime-free preview command gives users and tests direct visibility into generated prompts and manifests while preserving security and avoiding credentials.
- **Study / Source Grounding:** PRD dry-run prompt preview requirement at lines 205-210 and UX requirement at line 560; TRD prompt builder requirements at lines 915-925; `technical-handbook.md` patterns for IO injection, thin commands, and no-runtime trust boundaries; `06-io-abstraction` examples for capturable output; `13-security` examples for explicit boundaries.
- **Trade-Offs Accepted:** Preview does not reuse future runtime command paths, so later execution commands may need to share resolution/preflight intentionally. This is acceptable because it proves no-runtime behavior and avoids accidental side effects now.
- **Technical Debt / Future Impact:** Future `run --dry-run` behavior may either call the same builders or wrap this preview logic. Keep the command surface minimal and avoid overcommitting to runtime semantics.
- **Alternatives Rejected:** Implementing preview by calling `study run` or `study synthesize` with a dry-run flag is rejected because those paths are runtime-facing and may later grow side effects. Writing only to files and not stdout is rejected because requirements require stdout or explicit output path. Invoking agentwrap/OpenCode in a validation/preflight step is rejected because runtime integration is a non-goal and preview must work without credentials.
- **Contracts Satisfied:** CLI Surface, Security, Errors, Testing, Documentation, AC-14, AC-15, AC-17, constraints from `requirements.md` lines 65 and 68.
- **Evidence Required:** Command tests for stdout success, explicit output path success, unknown/ambiguous/inapplicable/missing-input failures, non-zero errors, and a no-runtime assertion through code path/fake dependency review.

### Decision 7: Tests Focus On Determinism, Source-Kind Semantics, Errors, And No-Runtime Boundaries

- **Decision:** Implement table-driven unit tests in `internal/study/prompts_test.go` and command tests in `internal/app/study_prompt_commands_test.go`. Use local fixtures or temporary workspaces; assert deterministic text/manifest behavior, required prompt content, missing input diagnostics, Markdown frontmatter stripping, applicability handling, synthesis manifests, preview stdout/output-path behavior, and absence of runtime invocation. Run `go test ./...` and `go build ./cmd/ultraplan`.
- **Rationale:** Prompt composition is deterministic local IO, so tests should prove exact behavior without OpenCode, agentwrap runtime execution, provider credentials, network access, or real source repositories.
- **Study / Source Grounding:** `technical-handbook.md` testing pattern summary; `11-testing-strategy` golden/fixture evidence including `helm/internal/test/test.go:43`, `go-task/task_test.go:166-169`, and `chezmoi/internal/cmd/main_test.go:64-174`; TRD unit/fixture/golden test requirements at lines 1593-1687.
- **Trade-Offs Accepted:** Some prompt text assertions may be golden or full-text and can be noisy on wording changes, while focused assertions may miss incidental regressions. Use both where appropriate: manifests and ordering should be exact; prose checks should target required sections unless the full prompt is intentionally stable.
- **Technical Debt / Future Impact:** Tests should avoid private implementation coupling so future runtime integration can reuse prompt builders without rewriting all tests.
- **Alternatives Rejected:** Testing only through CLI command tests is rejected because it would obscure domain failures and make error cases harder to isolate. Testing only with substring assertions is rejected because deterministic manifest ordering could regress. Requiring OpenCode or provider credentials in tests is rejected because requirements prohibit runtime dependencies for this sprint.
- **Contracts Satisfied:** Testing, Security, Errors, CLI Surface, AC-16 through AC-19, constraints from `requirements.md` lines 73 and 104.
- **Evidence Required:** Passing focused unit tests, passing command tests, `go test ./...`, `go build ./cmd/ultraplan`, and review notes confirming no real runtime/network/source repository dependency.

## Expected Evidence

| Evidence Type | Required Evidence | Source / Command / Review Check |
| --- | --- | --- |
| Tests | Directory prompt includes base prompt, dimension Markdown, report template, study/source/output metadata, hard source isolation, and file-path-line citation requirements. | `internal/study/prompts_test.go` |
| Tests | Markdown prompt embeds stripped document body, excludes YAML frontmatter, forbids external file/repository/code inspection, and defaults away from code citation requirements. | `internal/study/prompts_test.go` |
| Tests | Inapplicable Markdown source/dimension pairs are skipped or reported without runnable prompt generation and excluded from synthesis required inputs. | `internal/study/prompts_test.go` |
| Tests | Synthesis prompt includes synthesis prompt, dimension Markdown, final template, expected final output path, deterministic applicable report manifest, and missing applicable report failure. | `internal/study/prompts_test.go` |
| Tests | Manifest fields include prompt kind, study, dimension, source/source kind where applicable, template paths, input document/report paths, and expected output path using deterministic workspace-relative paths where possible. | `internal/study/prompts_test.go` |
| Tests | Missing `prompts/base.md`, `prompts/synthesize.md`, `templates/repo-analysis.md`, final report template inputs, Markdown documents, and applicable report inputs fail with actionable diagnostics naming the path. | `internal/study/prompts_test.go`, `internal/app/study_prompt_commands_test.go` |
| Tests | Preview command prints to stdout or writes to explicit output path and returns non-zero errors for unknown studies, unknown dimensions, unknown sources, ambiguous references, inapplicable pairs, missing templates, missing documents, and missing synthesis reports. | `internal/app/study_prompt_commands_test.go` |
| Runtime | No runtime evidence is expected because runtime execution is out of scope; required evidence is absence of runtime invocation. | Code review: preview path does not call agentwrap, OpenCode, provider APIs, network calls, or subprocess execution. |
| Review | Architecture conformance: study-owned prompt behavior, app-owned CLI wiring, platform packages do not import study, no global prompt/template package. | Architecture review protocol selected in `sprint-index.md`. |
| Review | Scope conformance: no runtime execution, run-loop, summary, code extraction, target, or sprint workflow scope added. | Sprint review protocol selected in `sprint-index.md`. |
| Build | Whole-repo tests pass. | `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go` |
| Build | CLI builds. | `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go` |
| Documentation | Sprint artifacts remain consistent and preview help/output does not imply execution. | `requirements.md`, `sprint-index.md`, `technical-handbook.md`, `reasoning.md`, `plan.md`, command help review |

## Assumptions And Risks

| Item | Type | Impact | Mitigation / Follow-Up |
| --- | --- | --- | --- |
| Existing study resolution, source discovery, frontmatter stripping, applicability filtering, and report path helpers are available from prior sprints. | Assumption | If absent or incomplete, prompt composition may need prerequisite fixes before core builder work. | Plan should inspect existing APIs first and use prior helpers rather than reimplementing discovery. |
| `templates/report.md` or another existing final report template path is the final report template input, while `templates/repo-analysis.md` is per-source analysis. | Assumption | Wrong template path would make synthesis prompts fail or use the wrong format. | Implementation should confirm current workspace template naming and tests should name the chosen final template path. |
| Dimensions currently may not have structured metadata for disabling code citations. | Risk | Directory prompt citation-disable behavior could become brittle if inferred from prose. | Centralize the rule in one helper and test supported wording or add minimal structured handling only if existing domain supports it. |
| Area-specific reasoning file selected by `sprint-index.md` was not present. | Risk | Some design nuance may be missing from the normal artifact chain. | This reasoning document makes the final prompt composition decisions directly and records absence explicitly. |
| Preview CLI shape may conflict with future runtime command dry-run UX. | Risk | Later command migration may be needed. | Keep preview minimal and backed by reusable study builders; avoid overcommitting to runtime semantics. |
| Prompt text can become too brittle in tests. | Risk | Template wording changes may create noisy failures. | Assert exact manifests and ordering; use focused required-text assertions unless golden coverage is intentionally valuable. |
| Manifest paths may accidentally leak absolute local details. | Risk | Generated previews become less portable and less safe to share. | Prefer workspace-relative paths in manifests; reserve absolute paths for diagnostics only. |
| Missing applicable reports for synthesis may frustrate users who want partial previews. | Risk | Strict failure blocks exploration of incomplete studies. | Error messages should identify missing reports and explain that synthesis requires applicable inputs; partial synthesis can be considered in a later explicit mode. |

## Implementation Constraints

- Prompt composition behavior must live in `internal/study`; do not create global `internal/prompts`, `internal/templates`, `internal/reports`, or `internal/validation` packages for this sprint.
- CLI parsing and preview rendering must stay in `internal/app`; prompt assembly, validation of prompt inputs, source-kind rules, applicability handling, and manifest construction must stay outside command handlers.
- Preserve dependency direction: `study` may use `workspace` and platform helpers, but platform packages must not import `study`; `platform/runtime` must not understand study semantics.
- Prompt builders must not invoke subprocesses, OpenCode, agentwrap runtimes, provider APIs, network calls, runtime health checks, or repository-wide scans.
- Prompt preview must not call future or existing runtime-facing analysis/synthesis execution paths.
- Prompt builders must load required workspace prompt/template files from deterministic resolved paths and fail before returning runnable prompt text when required inputs are missing.
- Markdown document prompts must use existing frontmatter stripping and applicability behavior, exclude frontmatter from embedded content, and forbid external file/repository/code inspection.
- Directory prompt source isolation and citation rules must be explicit and must not rely only on shared base prompt text.
- Synthesis manifests must include only applicable source reports, use deterministic ordering, and fail on missing applicable reports without inventing inputs.
- Manifest output must prefer workspace-relative paths where available and avoid absolute paths unless required for local diagnostics.
- Errors must preserve cause chains for filesystem/parsing failures and identify the failing template, dimension, source, report, document, or output path.
- Tests must use local fixtures or temporary directories and must not require OpenCode, agentwrap runtime execution, provider credentials, network access, or real source repositories.
- The implementation must remain within the roadmap Skeleton / Local CLI Gate and must not add runtime/provider or batch/durable workflow behavior.

## Plan Handoff

`plan.md` must execute these decisions. It must not invent architecture, scope, or decisions beyond this document.

The plan must carry forward:

- final decisions
- applicable contracts and requirement IDs
- expected evidence
- risks and assumptions
- required review protocols

The plan should order implementation so existing study/workspace APIs are inspected first, domain/request/manifest types are defined second, prompt builders and tests are added third, CLI preview wiring and command tests are added fourth, and full verification commands are run last.

## Phase Exit Criteria

- [x] Selected context was read and used.
- [x] Area-specific reasoning documents were completed where applicable or explicitly marked not applicable.
- [x] Area-specific reasoning conclusions are reflected or explicitly overridden in final decisions.
- [x] Contracts are explicitly mapped to decisions and expected evidence.
- [x] Final decisions are clear enough for `plan.md` to execute.
- [x] Expected evidence is specific and reviewable.
