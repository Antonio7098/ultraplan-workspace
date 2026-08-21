# Sprint Reasoning: Code-Context Stage Vertical Slice

> Project: `ultraplan-go`
> Sprint: `33-code-context-stage`
> Output: `projects/ultraplan-go/sprints/33-code-context-stage/reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/33-code-context-stage/requirements.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/sprints/33-code-context-stage/sprint-index.md`, `projects/ultraplan-go/sprints/33-code-context-stage/technical-handbook.md`, `projects/ultraplan-go/sprints/33-code-context-stage/reasoning/api-design.md`, `projects/ultraplan-go/sprints/33-code-context-stage/reasoning/architecture.md`, `projects/ultraplan-go/sprints/33-code-context-stage/reasoning/frontend.md`

This document decides. It synthesizes the selected context, handbook evidence, area-specific reasoning, and contracts into final Sprint 33 decisions. It does not replace the input artifacts and does not authorize implementation beyond this sprint.

## Sprint Purpose

- **Goal:** Make `code-context` a fully operational sprint stage immediately after requirements. It will inspect the resolved implementation repository without mutating it, produce and structurally validate one authoritative sprint-root `code-context.md`, preserve the last valid artifact across failed reruns, and expose truthful state and recovery through existing CLI, application, and browser boundaries.
- **Non-Goals:** Downstream requirements/code-context prompt-prefix injection, a manual code-context skill, broad documentation, and real-runtime release dogfood are Sprint 34 scope. Repository indexing, retrieval, embeddings, caches, a parallel JSON manifest, automatic staleness or amendment, expanded QA/repair, alternate persistence, graphs, cloud/Aren integration, source or Git mutation, and new generic stage/plugin/workflow/web frameworks are excluded.
- **Depends On:** Sprint 32's shared application and web operation baseline; existing sprint flow, state, artifact, runtime, configuration, defaults, mutation-lock, and atomic-write behavior; the project-owned implementation target; and the four implementation plans named in `requirements.md` as design inputs.

## Selected Context And Pre-Reasoning Artifacts

| Artifact | Path | How It Was Used |
| --- | --- | --- |
| Requirements | `projects/ultraplan-go/sprints/33-code-context-stage/requirements.md` | Defined the vertical-slice outputs, acceptance criteria, mutation boundary, compatibility behavior, evidence gates, and Sprint 34 exclusions. Acceptance bullets are referenced below as `AC-01` through `AC-18` in document order; constraints are `C-01` through `C-10` in document order. |
| Project Index | `projects/ultraplan-go/project-index.md` | Confirmed module ownership, available contracts, selected study evidence, review protocols, target implementation repository, and the absence of cataloged prior decisions. |
| Project Architecture | `projects/ultraplan-go/docs/ARCHITECTURE.md` | Fixed inward dependency direction, `internal/sprint` ownership, generic runtime and process boundaries, shared app use cases, web presentation ownership, and the evidence-gated architecture sequence. |
| Product Requirements | `projects/ultraplan-go/docs/PRD.md` | Established the Phase 5 product behavior, local-first filesystem authority, code-context artifact purpose, shared-interface parity, and explicit exclusions. |
| Technical Requirements | `projects/ultraplan-go/docs/TRD.md` | Supplied runtime, configuration, validation, atomic persistence, cancellation, web operation, security, and deterministic testing constraints, including the Sprint 33/Sprint 34 split. |
| Sprint Index | `projects/ultraplan-go/sprints/33-code-context-stage/sprint-index.md` | Limited this reasoning to the 12 selected contracts, 12 evidence reports, three area templates, two review protocols, and explicit excluded context. |
| Technical Handbook | `projects/ultraplan-go/sprints/33-code-context-stage/technical-handbook.md` | Distilled concrete Go CLI evidence for thin adapters, explicit composition, source-aware config, typed errors, caller-owned cancellation, bounded observability, narrow IO seams, layered tests, and fixed additive extension. |
| Prior Decision | Not applicable | `project-index.md` catalogs no prior decisions. Sprint 32 constraints are carried through `requirements.md`, not invented as a prior-decision artifact. |

## Area-Specific Reasoning Inputs

| Area | Reasoning Document | Key Conclusion | Evidence Basis | Impact On Final Decision |
| --- | --- | --- | --- | --- |
| API Design | `projects/ultraplan-go/sprints/33-code-context-stage/reasoning/api-design.md` | Extend the existing stage-valued CLI/app contract and generic guarded operation API; add no top-level command, route family, operation kind, durable web session, or caller-controlled target/output path. | Selected command, configuration, error, context, observability, testing, extensibility, and security reports; existing `/api/v1` operation model. | Fixes additive CLI/JSON behavior, synchronous read operations, asynchronous generation, error projection, conflict, confirmation, and compatibility semantics. |
| Architecture | `projects/ultraplan-go/sprints/33-code-context-stage/reasoning/architecture.md` | Keep the complete stage policy in focused `internal/sprint` code, use generic runtime and existing mechanical seams, generate to an isolated candidate, validate, atomically promote, interpret legacy state without read-time writes, and remain sequential. | All selected reports plus project module ownership and Phase 5 boundaries. | Fixes ownership, dependency direction, state compatibility, mutation scope, service flow, persistence ordering, and concurrency policy. |
| Frontend | `projects/ultraplan-go/sprints/33-code-context-stage/reasoning/frontend.md` | Extend the existing sprint page and typed view model with generic status, finding, preview, confirmation, progress, cancellation, and recovery components; preserve a complete no-JavaScript path. | Thin-adapter, typed-error, IO, cancellation, observability, testing, extensibility, and security evidence; Sprint 32 template hierarchy. | Fixes canonical DOM placement, truthful preserved-artifact/latest-attempt presentation, progressive enhancement, bounded rendering, and accessibility checks. |

All three area conclusions are adopted. None is overridden.

## Sprint Technical Handbook Summary

- **Relevant Patterns:** Thin adapters over an owning product service; manual composition and narrow volatile-boundary interfaces; command-local inputs with shared lifecycle wrappers; source-aware configuration; wrapped errors with stable classification; caller-owned context; independent repository-read and artifact-write boundaries; bounded asynchronous ownership; safe structured diagnostics; layered tests; centralized trust-boundary validation; and additive fixed extension rather than plugins.
- **Important Trade-Offs:** A central service improves consistency but grows `internal/sprint`; lazy runtime construction keeps reads cheap but delays execution preflight; typed failures aid recovery but can overgrow the compatibility surface; real temporary repositories prove containment and rename semantics but cost more than pure fakes; candidate promotion preserves prior output but requires cleanup and separate latest-attempt status; read-time compatibility preserves query purity but keeps interpretation logic alive.
- **Warnings / Anti-Patterns:** Do not create a large command handler, globals, hidden explicit-flag behavior, detached contexts, fire-and-forget work, global IO bypasses, runtime-exit-only success, raw diagnostic projection, shell commands from source text, a stage/plugin registry, broad filesystem abstraction, or tests coupled to incidental counts and dynamic output.
- **Evidence Confidence:** High for package boundaries, command delegation, dependency injection, configuration, errors, IO, context, observability, and layered testing because the handbook cites multiple mature repositories. Medium for concurrency, extensibility, and security where repository examples are useful but do not prove UltraPlan-specific policy. Sprint requirements and area reasoning close those policy gaps.

## Contracts Applied

| Contract / Requirement ID | Constraint | Decision Impact | Expected Evidence |
| --- | --- | --- | --- |
| Architecture; `AC-01`, `AC-02`; `C-01`, `C-02` | Sprint owns stage policy; platform remains generic; interfaces use app use cases. | One focused `internal/sprint` vertical slice with additive app/CLI/web wiring. | Architecture Review dependency and ownership checks; stage order and app/web contract tests. |
| CLI Surface; `AC-03`, `AC-04` | Existing prompt, validate, flow, help, status, dry-run, and JSON surfaces must accept the stage without semantic drift. | Extend existing stage-valued requests and projections; reads remain non-mutating. | Command help, prompt, dry-run, status, validation, error, and JSON fixtures. |
| Configuration; `AC-05`, `AC-12` | Explicit override, stage setting, and fallback sources must remain distinguishable and safely projected. | Add code-context model/variant keys to existing source-aware resolution and redaction paths. | Parsing, precedence, supplied-flag, validation, source tracking, fallback, and redaction tests. |
| Documentation; `AC-11` | Embedded defaults and executable help are in scope; broad guides are not. | Ship synchronized embedded prompt/template defaults and update help only. | Embedded/materialized default parity and help tests; scope review confirms broad docs remain deferred. |
| Errors; `AC-07`, `AC-08`, `AC-09`, `AC-14` | Invalid output, runtime failure, cancellation, persistence failure, and uncertainty must be distinguishable and actionable. | Reuse stable categories with wrapped causes; no success from runtime exit or stale artifact presence. | Error mapping, missing/invalid output, cancellation, cleanup uncertainty, and recovery tests. |
| LLM Evaluation / Cost / Safety; `AC-05`, `AC-08`, `AC-12` | Runtime use must be attributable, bounded, validated, and safely projected. | One explicit operation, no transport auto-retry or stage fan-out, safe metadata only, unknown usage remains unknown. | Fake-runtime metadata and validation tests; review of redaction and no raw prompt/provider payloads. |
| LLM Runtime; `AC-05`, `AC-08`; `C-08` | Use the generic runtime with `context.Context`, stage overrides, permissions, and expected output behavior. | Product code constructs a generic request; runtime package gains no sprint semantics or direct OpenCode handling. | Request construction, context propagation, cancellation, runtime success/failure, and dependency review. |
| Observability; `AC-12`, `AC-13`, `AC-14` | Status and progress must be truthful, bounded, correlated, and safe across surfaces. | Reuse operation events and expose stage, attempt, model/variant source, validation, terminal outcome, and next action. | CLI/app/web parity, bounded event/SSE, redaction, reconnect, shutdown, and preserved-artifact/latest-attempt tests. |
| Persistence And Migrations; `AC-09`, `AC-10`; `C-03`, `C-06` | Artifact and state writes are atomic; legacy state remains readable; read-only status does not migrate. | Candidate validation precedes atomic promotion; legacy state is interpreted on read and serialized only on a later mutation. | Rename/write failure tests, legacy fixtures, no-write status checks, and recovery tests. |
| Security; `AC-05`, `AC-06`, `AC-07`, `AC-12`; `C-04`, `C-05` | Repository reads are contained, only one sprint artifact may be written, paths are safe, and projections are redacted. | Resolve target internally, reject caller paths and escaping references, use allowlisted bounded preview, and verify source/Git immutability. | Temporary repository before/after and Git-state comparisons, path tables, hostile Markdown, preview allowlist, and redaction tests. |
| Testing; `AC-15`, `AC-16`; `C-09` | Normal tests are deterministic/offline and must cover the complete vertical slice and release gates. | Use validator tables, fakes, fixtures, temporary repositories, command/app/web contracts, race tests, and repository-wide commands. | Focused package tests plus `go test ./...`, `go test -race ./...`, `go vet ./...`, `go build ./cmd/ultraplan`, and `git diff --check`. |
| Workflows; `AC-01`, `AC-02`, `AC-03`, `AC-08`, `AC-09`, `AC-13`, `AC-14`; `C-07` | Stage order, prerequisites, cumulative dispatch, validation gating, conflict, rerun, cancellation, and recovery must use one workflow truth. | Add one canonical stage and reuse existing flow and generic operation semantics with single terminal arbitration. | Ordered-stage, cumulative-flow, mutation-conflict, exact-once, rerun, cancellation, and recovery tests. |
| Sprint scope; `AC-17`, `AC-18`; `C-10` | Deferred systems and Sprint 34 work must not enter this vertical slice. | No index/RAG/cache/manifest/staleness/amendment/prefix/skill/broad-doc/real-runtime work. | Diff/dependency review and both required review protocols explicitly confirm exclusions. |

## Repos Studied / Source Evidence Used

The source reports were used through `technical-handbook.md` and the completed area reasoning. The concrete references below are the handbook's selected evidence, not new unselected context.

| Source / Repo / Report | Concrete Reference | Relevant Finding | Why It Matters For This Sprint | Used In Decision(s) |
| --- | --- | --- | --- | --- |
| restic / `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md`; `restic/cmd/restic/main.go:37-114`, `restic/internal/restic/repository.go:18` | Thin entrypoints and inward, acyclic product boundaries. | Supports sprint-owned behavior behind app/interface adapters. | 1, 2, 5 |
| Helm / `01-project-structure`, `02-command-architecture` | `helm/pkg/cmd/install.go:132-145,333-347`, `helm/pkg/action/install.go:73-140` | Commands delegate to action behavior and propagate caller cancellation. | Supports shared use cases and one operation context instead of CLI/web workflow duplication. | 2, 4, 5 |
| gh-cli / `02-command-architecture`, `03-dependency-injection` | `gh-cli/pkg/cmdutil/factory.go:16-43`, `gh-cli/internal/ghcmd/cmd.go:52-132` | Focused factories and explicit manual dependency construction. | Supports additive stage wiring and narrow existing seams without DI/plugin machinery. | 2, 4 |
| rclone / `02-command-architecture`, `12-extensibility` | `rclone/cmd/cmd.go:240-340`, `rclone/fs/rc/registry.go:41-48` | Shared lifecycle wrappers are valuable; dynamic registries add collision and initialization costs. | Supports generic operation reuse and rejects a stage registry. | 1, 4, 5 |
| go-task and restic / `04-configuration-management` | `go-task/internal/flags/flags.go:314-327`, `restic/internal/global/global.go:139,147` | Effective configuration must retain precedence and whether a flag was explicitly supplied. | Prevents omitted CLI flags from masking stage/global model and variant fallbacks. | 4 |
| k9s / `04-configuration-management`, `13-security` | `k9s/internal/config/k9s.go:423-451`, `k9s/internal/config/json/validator.go:146` | Validate after source merge and centralize structural trust-boundary checks. | Supports effective config validation and one sprint-owned structural validator. | 3, 4 |
| Helm and gh-cli / `05-error-handling` | `helm/pkg/storage/driver/driver.go:27-48`, `gh-cli/internal/ghcmd/cmd.go:44-49,281-301` | Preserve error identity and render it appropriately at boundaries. | Supports stable shared failure classes and safe CLI/web recovery guidance. | 5 |
| restic and gh-cli / `06-io-abstraction` | `restic/internal/fs/interface.go:10-31`, `gh-cli/pkg/iostreams/iostreams.go:551-568` | External-system seams and injectable streams aid testing, but broad IO abstraction has cost. | Supports existing narrow runtime/atomic seams plus real filesystem tests rather than a virtual filesystem. | 2, 7 |
| restic and opencode / `07-state-context`, `08-concurrency` | `restic/internal/restic/lock.go:290-305`, `opencode/root.go:261-279` | Cleanup may need a bounded context; cancellation must be followed by explicit bounded waiting. | Supports canonical cancellation, bounded reconciliation, and explicit uncertainty. | 2, 5 |
| lazygit / `08-concurrency`, `13-security` | `lazygit/pkg/gui/background.go:35,46,123`, `lazygit/cmd_obj_builder.go:38` | Asynchronous work needs ownership and bounds; explicit argv is safer than shell construction. | Supports no stage-local fan-out and preservation of generic runtime/process trust boundaries. | 2, 5 |
| Helm, k9s, go-task / `10-logging-observability` | `helm/internal/logging/logging.go:31-71`, `k9s/internal/slogs/keys.go:6-231`, `go-task/internal/output/output.go:12-14` | Structured diagnostics must stay separate from machine-readable result output. | Supports bounded safe fields and prevents JSON/SSE contamination by raw runtime output. | 4, 5 |
| chezmoi, Helm, restic / `11-testing-strategy` | `chezmoi/internal/cmd/main_test.go:64-174`, `helm/internal/test/test.go:43`, `restic/cmd/restic/integration_helpers_test.go:188-235` | Behavior, golden, fixture, fake, and isolated integration tests prove different contracts. | Supports the layered verification matrix and selective use of goldens. | 7 |
| go-task / `12-extensibility` | `go-task/executor.go:20-24,91-122` | Additive options and internal factories evolve fixed applications at lower cost than plugins. | Supports explicit canonical stage registration and minimum composition changes. | 1, 2 |
| restic and Helm / `13-security` | `restic/internal/options/secret_string.go:15-20`, `helm/pkg/registry/transport.go:37-41` | Dedicated safe formatting and scrubbed diagnostics are necessary at trust boundaries. | Supports projection-time allowlists and exclusion of prompts, excerpts, unsafe paths, and raw payloads. | 3, 4, 5 |

## Trade-Off And Debt Analysis

### Accepted Trade-Offs

| Trade-Off | Benefit | Cost / Constraint Accepted | Why Acceptable Now | Revisit Trigger |
| --- | --- | --- | --- | --- |
| Grow `internal/sprint` with focused files | Keeps stage policy, state, validation, and artifact semantics under one owner. | The package becomes larger. | Code-context is sprint workflow state, not an independent product domain. | Concrete readability or dependency pressure demonstrates a stable subpackage boundary. |
| Explicit fixed stage registration | Deterministic order and compile-time discoverability. | Canonical projections can drift if separately maintained. | The workflow is closed and product-defined. | Multiple externally contributed stages require versioned lifecycle semantics. |
| Lazy runtime construction | Prompt, validate, status, and dry-run stay cheap and non-mutating. | Runtime/preflight errors occur at execution time. | Read paths must not invoke expensive or side-effectful infrastructure. | Users need a separate explicit runtime preflight for this stage. |
| Isolated candidate then atomic promotion | Failed, cancelled, or invalid reruns preserve the last valid artifact. | Candidate cleanup and artifact/latest-attempt distinction add complexity. | Preserving authoritative content is more important than direct-write simplicity. | Existing atomic helpers cannot provide same-directory replacement and cleanup guarantees. |
| Structural, not semantic, validation | Deterministically rejects malformed and unsafe artifacts. | A valid artifact may still omit important repository context. | Semantic completeness cannot be proven mechanically in Sprint 33. | A later evidence-backed QA or amendment workflow defines semantic review. |
| Read-time compatibility without migration writes | Legacy workspaces remain usable and read-only status remains pure. | Compatibility interpretation stays in the load path. | Hidden migration would violate query semantics and recovery expectations. | An explicit versioned migration policy removes support for the legacy representation. |
| Sequential stage orchestration | Simplifies cancellation, cost attribution, mutation scope, and terminal arbitration. | No stage-internal repository-inspection speedup. | No measured independent work units justify fan-out. | Profiling shows unacceptable latency and a bounded partitioning model exists. |
| Real temporary repository tests for filesystem guarantees | Proves containment, source immutability, permissions, and rename behavior. | Slower and more environment-sensitive than pure fakes. | These guarantees depend on real local filesystem semantics. | A proven contract suite can supply equivalent cross-platform guarantees. |
| Reuse generic browser operations | Preserves confirmation, conflict, progress, cancellation, and recovery semantics. | Generic request/view models may need small additive fields. | A route-specific workflow would duplicate product behavior and durable truth. | The existing operation contract cannot represent a sprint stage after narrow additive extension. |

### Potential Technical Debt

| Debt / Shortcut | Why It Might Accrue | Current Mitigation | Owner / Follow-Up |
| --- | --- | --- | --- |
| Duplicated stage projections | Order, help, artifacts, config, status, and templates may each encode stage membership. | Drive from the current canonical list where possible and add exact-order/parity tests; do not add a dynamic registry. | Sprint 33 implementation and Architecture Review. |
| Legacy compatibility branch | Supporting pre-code-context state adds load-path complexity. | Isolate deterministic interpretation, fixture it, and serialize current state only on mutation. | Future explicit flow-state migration decision. |
| Candidate cleanup residue | Process death or cleanup timeout may leave temporary files. | Use contained candidate paths, bounded cleanup, and truthful `interrupted`/`cleanup_uncertain` outcomes. | Sprint 33 recovery tests; revisit with durable worker design only if needed. |
| Structural validator format coupling | Markdown parser rules may become brittle as users edit the template. | Validate semantic headings/blocks and path/range grammar, not incidental whitespace or prose. | Future template-version decision if format evolution proves necessary. |
| App composition growth | More stage dependencies can enlarge the central app container. | Inject the existing cohesive sprint service/lazy runtime only; keep helpers in `internal/sprint`. | Future architecture review after repeated additions. |
| Stable JSON expansion | New stage and optional metadata can accidentally change nullability or leak internal types. | Additive DTO fields only, compatibility fixtures, and explicit projection/redaction tests. | Sprint 33 API tests and Sprint Review. |
| UI distinction complexity | A valid preserved artifact and failed latest rerun require multiple state dimensions. | Add explicit typed fields and semantic labels; never infer this in templates. | Sprint 33 app/web view-model tests. |

### Future Considerations

| Consideration | Deferred Until | Reason Deferred | What Should Be Preserved Now |
| --- | --- | --- | --- |
| Shared exact prompt prefix | Sprint 34 | Current scope ends at the standalone vertical slice. | Stable authoritative `requirements.md` and `code-context.md` bytes and deterministic prompt inputs. |
| Manual code-context skill | Sprint 34 | Explicit roadmap split. | Stage command/app behavior must be callable without a special manual abstraction. |
| Broad docs and real-runtime dogfood | Sprint 34 release gate | Sprint 33 requires executable help/defaults and deterministic fake coverage only. | Accurate help, safe metadata, fake-first test seams, and no false real-runtime claim. |
| Automatic freshness/amendments | Later explicit sprint | Source identity and amendment semantics are not yet defined. | Explicit rerun, truthful latest outcome, and no automatic freshness claim. |
| Content identity, QA/repair, retrieval, persistence, graph, cloud/Aren | Evidence-gated post-Phase-5 work | Each requires prior dogfood and measured value. | Markdown/filesystem authority, product-owned state, narrow seams, and no shadow index/store. |
| Stage-internal concurrency | Measured performance trigger | Current evidence supports ownership and bounds, not parallelization. | One caller-owned context and generic event handling that can remain bounded. |

## Decisions

The authoritative decisions are specified in full under **Final Decisions** below:

1. Insert one canonical `code-context` stage and interpret legacy flow state without read-time writes.
2. Implement a sequential sprint-owned runtime service with isolated candidate generation and atomic promotion.
3. Keep one authoritative Markdown artifact under a strict structural and containment validator.
4. Extend source-aware configuration and the existing stage-valued CLI/application contract.
5. Reuse generic browser operations and present artifact validity separately from the latest attempt.
6. Propagate canonical cancellation, stable failure identity, bounded cleanup, and safe observability.
7. Prove the vertical slice with layered offline tests, real temporary repositories where needed, release commands, and required reviews.

## Final Decisions

### Decision 1: Canonical Stage Order And Legacy State Compatibility

- **Decision:** Add `StageCodeContext` exactly once to the canonical planning sequence immediately after `requirements` and before `sprint-index`. All order, prerequisite, cumulative-flow, status, artifact, and stage-count behavior must derive from that definition. Valid completed requirements make the stage ready; only a successful run with a structurally valid authoritative artifact makes sprint-index ready. Pre-code-context flow state is interpreted deterministically without a read-time write, preserves prior outcomes, and is emitted in current form only during a later explicit mutation.
- **Rationale:** One workflow vocabulary prevents divergent behavior across flow, CLI, app, and web, while side-effect-free compatibility keeps old workspaces trustworthy.
- **Study / Source Grounding:** `technical-handbook.md` sections "Relevant Patterns," "Trade-Offs," and "Design Pressures" cite explicit root registration in Helm, shared lifecycle behavior in rclone, and lower-cost additive extension in go-task (`02-command-architecture`, `12-extensibility`). The Architecture area reasoning applies that evidence to the existing sprint domain and rejects read-time migration.
- **Trade-Offs Accepted:** Fixed registration requires parity tests; compatibility logic remains in the load path until a future explicit migration policy.
- **Technical Debt / Future Impact:** A legacy interpretation branch and possible duplicated projections remain. Preserve a canonical source and compatibility fixtures rather than introducing a plugin registry or second state format.
- **Alternatives Rejected:** A dynamic stage/plugin registry adds collision, lifecycle, and versioning complexity; a parallel workflow/state file creates split truth; read-time migration hides writes; treating artifact presence alone as completion violates validation gating.
- **Contracts Satisfied:** Architecture, Persistence And Migrations, Workflows, CLI Surface; `AC-01`, `AC-02`, `AC-03`, `AC-08`, `AC-10`; `C-01`, `C-03`, `C-06`.
- **Evidence Required:** Exact order and exact-once membership tests; prerequisite and cumulative `flow --to` tests; stage count/artifact mapping/status JSON tests; legacy fixtures for representative prior outcomes; byte/mtime or injected-store proof that status does not write; atomic current-state write on a later mutation.

### Decision 2: Sprint-Owned Sequential Runtime Service And Atomic Promotion

- **Decision:** Implement one cohesive `internal/sprint` code-context service. It validates prerequisites, resolves the implementation target through existing project/execute mechanisms, resolves prompt/template and effective model/variant, constructs a generic runtime request with the caller context and read-only posture, generates to a contained isolated candidate, requires runtime success plus candidate existence and structural validity, then atomically promotes only `code-context.md` and persists the truthful stage transition. The stage itself remains sequential and uses existing mutation conflict and runtime/event infrastructure.
- **Rationale:** This is the narrowest design that keeps product semantics with sprint state, runtime infrastructure generic, source and Git immutable, and reruns safe.
- **Study / Source Grounding:** `technical-handbook.md` cites restic/Helm inward boundaries (`01-project-structure`), gh-cli explicit factories (`03-dependency-injection`), restic IO seams (`06-io-abstraction`), Helm/restic context and cleanup behavior (`07-state-context`), and bounded concurrency examples (`08-concurrency`). `architecture.md` resolves these pressures in favor of a focused sprint service, real filesystem guarantees, and sequential orchestration.
- **Trade-Offs Accepted:** The sprint package grows; lazy runtime preflight occurs later; candidate cleanup is required; real filesystem tests cost more; no parallel speedup is attempted.
- **Technical Debt / Future Impact:** Candidate residue and app composition growth are possible. Keep candidate paths contained, cleanup bounded, and dependencies narrow. Sprint 34 can reuse the authoritative artifact without changing this service boundary.
- **Alternatives Rejected:** A new `internal/codecontext` domain splits sprint-owned state; runtime-owned validation/transitions reverse dependency direction; direct generation over the authoritative artifact risks corruption; a broad virtual filesystem obscures real rename/containment semantics; parallel repository workers add unearned race and cancellation complexity; direct OpenCode/process handling duplicates agentwrap.
- **Contracts Satisfied:** Architecture, LLM Runtime, LLM Evaluation / Cost / Safety, Security, Persistence And Migrations, Workflows, Errors; `AC-04`, `AC-05`, `AC-06`, `AC-08`, `AC-09`, `AC-14`; `C-01`, `C-04`, `C-05`, `C-06`, `C-08`.
- **Evidence Required:** Fake-runtime success, runtime failure, missing output, invalid output, timeout, cancellation, and override tests; temporary implementation repository before/after and Git-state comparisons; mutation conflict tests; failed write/rename and failed rerun preservation tests; proof that only sprint-root `code-context.md` changes; no goroutine leak/race findings.

### Decision 3: Single Markdown Artifact And Structural Validation Contract

- **Decision:** Store exactly one authoritative output at `projects/<project>/sprints/<sprint>/code-context.md`. Ship an embedded editable template covering sprint scope, inspected repository areas, selected source excerpts, relationships, constraints, and open questions. Central sprint validation requires all sections, at least one exact language-tagged fenced source excerpt, a repository-relative contained path, a rationale for each excerpt, and well-formed optional line ranges and symbols. It rejects absent/empty output, template placeholders, missing data, absolute or escaping paths, malformed ranges, and missing fences. Validation explicitly does not claim semantic completeness.
- **Rationale:** Human-editable Markdown is consistent with filesystem authority and provides durable prepared evidence without creating an index or machine-readable shadow artifact.
- **Study / Source Grounding:** `technical-handbook.md` uses k9s centralized validation and restic/Helm redaction evidence (`13-security`) to support one trust-boundary validator, while warning that source evidence cannot establish semantic completeness. The Architecture and Frontend area reasoning require bounded, allowlisted handling of the resulting artifact.
- **Trade-Offs Accepted:** Structural validity cannot prove that the best or all relevant code was selected; Markdown grammar must remain stable enough for deterministic parsing.
- **Technical Debt / Future Impact:** Validator/template coupling may grow. Validate durable semantic structure rather than incidental prose or whitespace, and defer versioned semantic identity or amendment behavior to a later gated sprint.
- **Alternatives Rejected:** A parallel JSON manifest creates a second source of truth; a repository index/RAG/cache exceeds scope; semantic completeness scoring would overclaim evidence; arbitrary hard excerpt limits lack a requirement; source mutation or generated test output violates the stage boundary.
- **Contracts Satisfied:** Security, Documentation, Testing, Persistence And Migrations, LLM Evaluation / Cost / Safety; `AC-06`, `AC-07`, `AC-08`, `AC-09`, `AC-11`, `AC-18`; `C-03`, `C-04`, `C-05`, `C-07`.
- **Evidence Required:** Table-driven valid/invalid fixtures for every section, excerpt, path, rationale, fence, language, range, symbol, empty content, and placeholder rule; containment tests; representative artifact review; embedded/template materialization parity; no parallel manifest/index/cache dependency or file in the diff.

### Decision 4: Source-Aware Configuration And Additive CLI/App Contract

- **Decision:** Add code-context model and variant settings to the existing fixed precedence, validation, source tracking, redaction, and `config show` projections. Preserve whether command overrides were explicitly supplied. Extend the existing sprint-stage requests for `prompt`, `validate`, `flow --to`, help, status, and stable JSON. Prompt preview, validation, status, readiness, config display, artifact preview, and dry-run remain synchronous and non-mutating; generation/rerun uses the existing long-running operation capability. Callers cannot supply target or output paths.
- **Rationale:** This preserves one scriptable command and typed application vocabulary while making effective runtime choice explainable and preventing path injection.
- **Study / Source Grounding:** `technical-handbook.md` cites go-task, restic, and k9s for explicit override tracking and post-merge validation (`04-configuration-management`), gh-cli/Helm for thin command wiring (`02-command-architecture`), and Helm/k9s/go-task for result/diagnostic separation (`10-logging-observability`). `api-design.md` adopts these patterns for additive stage-valued requests.
- **Trade-Offs Accepted:** Parsers and request types must carry supplied/not-supplied metadata; every stable projection needs drift tests; runtime preflight is delayed until execution.
- **Technical Debt / Future Impact:** Stable JSON fields and canonical stage projection can drift. Keep changes additive, map through explicit DTOs, and retain current envelope/nullability semantics.
- **Alternatives Rejected:** A top-level `code-context` command fragments sprint flow; caller-controlled paths weaken containment; bound default values masquerading as explicit flags break fallback truth; eager runtime construction makes read paths expensive; full prompts or raw diagnostics in status/JSON leak sensitive and unstable data.
- **Contracts Satisfied:** CLI Surface, Configuration, Errors, Observability, LLM Runtime, Security; `AC-03`, `AC-04`, `AC-05`, `AC-12`; `C-02`, `C-04`, `C-08`.
- **Evidence Required:** Help and stage parsing tests; deterministic prompt preview and dry-run no-mutation checks; model/variant precedence, explicit-source, fallback, validation, and redaction tables; stable JSON fixtures; typed app request/result tests; proof that preview/dry-run do not construct/invoke runtime or create candidate/artifact/state files.

### Decision 5: Generic Browser Operation And Truthful Sprint Presentation

- **Decision:** Extend the existing sprint page and typed view model to show code-context in canonical order using current generic status, finding, artifact-preview, confirmation, operation-progress, cancellation, error, and next-action components. Runtime generation and explicit rerun use the existing prepare/start/status/events/cancel contract and mutation conflict key. No route, operation kind, durable browser session, client store, or stage-specific JavaScript controller is added. The page must separately present authoritative artifact validity, stage readiness, and latest operation outcome, including a failed rerun with a preserved valid artifact. Server rendering remains complete without JavaScript; SSE only enhances progress.
- **Rationale:** The browser is a projection over shared app truth. Reusing the existing operation lifecycle preserves security, recovery, accessibility, and CLI/web parity without a second product workflow.
- **Study / Source Grounding:** `technical-handbook.md` thin-adapter and shared-lifecycle evidence from Helm/rclone (`01-project-structure`, `02-command-architecture`), caller-context evidence (`07-state-context`), and bounded diagnostics evidence (`10-logging-observability`, `13-security`) support this choice. `frontend.md` specifies canonical placement, no-JavaScript behavior, and preserved-artifact/latest-attempt separation.
- **Trade-Offs Accepted:** The view model becomes richer and generic components must remain expressive; explicit confirmation adds an operator step; bounded previews omit full artifact content.
- **Technical Debt / Future Impact:** The sprint page may expose hardcoded stage assumptions. Replace only those necessary with canonical typed projection; do not introduce a frontend framework or component registry.
- **Alternatives Rejected:** A dedicated route/page duplicates sprint context; a new operation kind moves product semantics into web; browser-owned durable state conflicts with filesystem authority; automatic rerun violates explicit intent and excluded staleness work; cancel-on-disconnect confuses subscription with work ownership; initial full artifact rendering weakens bounds; raw runtime logs are unsafe and inaccessible progress UI.
- **Contracts Satisfied:** Architecture, Security, Observability, Workflows, Testing, CLI Surface; `AC-12`, `AC-13`, `AC-14`; `C-02`, `C-03`, `C-04`.
- **Evidence Required:** App/web contract and `httptest` coverage for prepare/start/status/events/cancel reuse, no route-specific kind, confirmation staleness, conflict, disconnect, reconnect, slow subscribers, explicit cancellation, shutdown, and recovery; template tests for DOM order, all stage states, semantic controls/live regions, no-JavaScript flow, hostile Markdown, escaped/redacted values, bounded allowlisted preview, and preserved artifact plus failed latest attempt.

### Decision 6: Canonical Cancellation, Stable Failure Identity, And Safe Observability

- **Decision:** One caller-owned `context.Context` reaches target resolution, runtime start/events/wait, validation, and normal persistence. Explicit cancellation and server shutdown use the existing canonical cancellation function and single-terminal arbitration; browser disconnect only removes a subscription. Bounded reconciliation may outlive the work context under the existing cleanup policy, but timeout records `interrupted` or `cleanup_uncertain`, never success. Reuse existing typed/stable identities for prerequisite, config, conflict, runtime, missing output, invalid output, persistence, cancellation, interruption, and cleanup uncertainty, preserving wrapped causes internally. Public metadata is allowlisted and bounded to correlation, stage, attempt, runtime/model/variant source, timing, validation, cancellation, terminal outcome, and next action; unknown usage/cost remains unknown.
- **Rationale:** Operators need reliable recovery semantics across CLI and browser, and uncertainty must never be collapsed into success or inferred from artifact presence.
- **Study / Source Grounding:** `technical-handbook.md` cites Helm for signal context, restic for bounded cleanup, and opencode/lazygit for cancel-and-wait ownership (`07-state-context`, `08-concurrency`); Helm/gh-cli support typed boundary rendering (`05-error-handling`); Helm/k9s/restic support structured safe projection (`10-logging-observability`, `13-security`). Both Architecture and API area reasoning adopt these conclusions.
- **Trade-Offs Accepted:** Bounded cleanup can delay return; explicit uncertainty adds states and UI copy; stable categories require compatibility discipline.
- **Technical Debt / Future Impact:** Over-classification could expand the public surface. Add a new class only when existing identity plus wrapped detail cannot express recovery behavior. Preserve cleanup metadata separately from primary outcome.
- **Alternatives Rejected:** Detached contexts make cancellation decorative; fire-and-forget cleanup creates uncertain ownership; cancel-on-SSE-disconnect violates server operation ownership; raw errors/provider events disclose unstable data; string parsing is not a contract; runtime exit success or old artifact presence cannot establish stage success; transport auto-retry risks duplicate cost.
- **Contracts Satisfied:** Errors, Observability, Workflows, LLM Runtime, LLM Evaluation / Cost / Safety, Security; `AC-08`, `AC-09`, `AC-12`, `AC-13`, `AC-14`; `C-03`, `C-08`.
- **Evidence Required:** Context propagation and runtime cancel tests; event drain and wait assertions; idempotent explicit/shutdown cancellation; disconnect non-cancellation; terminal race arbitration; bounded cleanup timeout and restart recovery; safe error-code/exit mapping; redaction and event-bound tests; missing/invalid output after successful runtime remains failed.

### Decision 7: Layered Offline Verification And Scope Enforcement

- **Decision:** Build deterministic evidence in layers: domain/order and compatibility fixtures; validator tables; fake-runtime stage tests; real temporary repository/filesystem isolation and atomicity tests; embedded default and configuration tests; CLI/stable JSON tests; typed app and generic web operation tests; browser template/accessibility/security tests; cancellation/shutdown/recovery tests; and repository-wide test, race, vet, build, and diff checks. Use goldens only for stable help/default/JSON/Markdown shapes and structural assertions for dynamic paths, IDs, timestamps, runtime details, and progress. Real-runtime dogfood is explicitly deferred, not marked passed or simulated.
- **Rationale:** No single test style can prove stage semantics, filesystem isolation, compatibility, interface parity, and recovery. Layering provides reviewable evidence while keeping normal tests offline.
- **Study / Source Grounding:** `technical-handbook.md` cites chezmoi command tests, Helm golden support, and restic isolated integration environments (`11-testing-strategy`), plus gh-cli/restic IO seams (`06-io-abstraction`). The handbook warns against brittle incidental assertions. The area reasoning documents turn those findings into concrete test matrices.
- **Trade-Offs Accepted:** The suite is broader and temporary-repository/race tests cost more than helper-only tests; selective goldens require explicit review when updated.
- **Technical Debt / Future Impact:** Test duplication across CLI/app/web is possible. Share fixtures and compare projections of the same typed state, but retain boundary-specific assertions. Sprint 34 must add gated real-runtime and shared-prefix evidence rather than weakening this baseline.
- **Alternatives Rejected:** Helper-only unit tests miss integration contracts; end-to-end-only tests obscure causes and require live providers; blanket snapshots are brittle; in-memory-only filesystem tests cannot prove rename/containment; claiming unavailable real-runtime evidence as passing violates truthfulness; pulling Sprint 34 proof into this sprint breaks roadmap scope.
- **Contracts Satisfied:** Testing, Architecture, Security, Persistence And Migrations, Workflows, Documentation; `AC-01` through `AC-18`; `C-01` through `C-10` as applicable to their covered behavior.
- **Evidence Required:** Focused tests listed in `requirements.md` Required Outputs; `go test ./...`; `go test -race ./...`; `go vet ./...`; `go build ./cmd/ultraplan`; `git diff --check`; Architecture Review and Sprint Review confirmations; explicit diff/dependency check for every `AC-18` exclusion and Sprint 34 deferral.

## Expected Evidence

| Evidence Type | Required Evidence | Source / Command / Review Check |
| --- | --- | --- |
| Domain and flow tests | Exact canonical order, membership, stage count, prerequisites, cumulative dispatch, readiness, artifact mapping, and stage transitions. | `internal/sprint/sprint_test.go`, `internal/sprint/sprint_index_test.go`, focused `go test` for `internal/sprint`. |
| Compatibility and persistence tests | Pre-stage fixtures preserve outcomes; status performs no write; current state writes atomically on mutation; failed candidate/promotion preserves last valid artifact. | `internal/sprint/verify_test.go`, `internal/sprint/code_context_test.go`; file content/mtime and injected atomic-failure checks. |
| Validator and isolation tests | Every structural failure is actionable; paths/ranges are contained; source and Git state are unchanged; only `code-context.md` may change. | Table-driven `internal/sprint/code_context_test.go` with temporary repositories and before/after snapshots. |
| Runtime tests | Fake-runtime success/failure, missing/invalid output, model/variant overrides, context cancellation, event drain/wait, metadata, and rerun behavior. | `internal/sprint/code_context_test.go`; fake `agentwrap.Runtime`/`Run`; normal tests remain offline. |
| Defaults and configuration tests | Embedded prompt/template parity with materialized defaults; precedence, explicit source, fallback, validation, projection, and redaction. | `internal/workspace/workspace_test.go`, `internal/platform/config/config_test.go`, config command tests. |
| CLI and stable JSON tests | Help, accepted stage, prompt preview, validation, flow, status, dry-run non-mutation, error codes, and additive stable JSON. | `internal/app/sprint_commands_test.go`; structural JSON fixtures. |
| Shared app/web tests | Typed readiness/artifact/findings projection; generic operation prepare/run/fingerprint/progress/cancel/recovery; no route-specific operation kind. | `internal/app/web_usecases_test.go`, `internal/app/web_operations_test.go`, `internal/web/operations_contract_test.go`. |
| Browser tests | Canonical placement; all states; preserved artifact/latest attempt distinction; bounded safe preview; no-JavaScript use; accessibility; hostile content; confirmation and cancellation behavior. | `internal/web/templates_test.go` and `httptest` coverage. |
| Cancellation and recovery tests | Disconnect does not cancel; explicit and shutdown cancellation reach canonical context; bounded cleanup yields truthful terminal/uncertain outcomes; restart does not infer success. | App/web operation and shutdown tests under race detection. |
| Release commands | All package tests, race tests, vet, binary build, and whitespace checks pass in the implementation repository. | From `../ultraplan-go`: `go test ./...`, `go test -race ./...`, `go vet ./...`, `go build ./cmd/ultraplan`, `git diff --check`. |
| Review | Ownership, dependency direction, scope, mutation boundaries, contracts, evidence, and exclusions are confirmed. | `system/protocols/architecture-review-protocol.md` and `system/protocols/review-sprint-protocol.md`. |
| Documentation | Help and embedded/default-install prompt/template behavior are current; broad guides and release dogfood are explicitly deferred. | Help/default golden tests and Sprint Review scope check. |

## Assumptions And Risks

| Item | Type | Impact | Mitigation / Follow-Up |
| --- | --- | --- | --- |
| Existing project/execute target resolution can provide the implementation repository without a new discovery service. | Assumption | If false, the stage cannot safely locate its read root. | Extend the existing project-owned typed resolution narrowly; do not accept raw caller paths or create an index. |
| Existing atomic helper supports same-directory replacement while preserving the old file on failure. | Assumption | If false, rerun safety is incomplete. | Verify with fault tests; add only the missing mechanical guarantee in the existing filesystem boundary. |
| Existing generic operation request and conflict key can represent one sprint stage mutation. | Assumption | If false, web reuse may be blocked. | Narrowly extend the typed generic app request/key; do not add a code-context route or operation kind. |
| Existing sprint/app view model can represent artifact validity separately from latest attempt outcome. | Assumption | If false, failed reruns may appear successful. | Add explicit typed projection fields before template rendering. |
| Runtime permission policy plus fixed work/output policy is enforceable by current agentwrap integration. | Assumption | Unsupported required restrictions could permit mutation risk. | Fail preflight when required policy is unsupported and verify repository/Git immutability in temporary-repository tests. |
| Canonical stage lists can drift across projections. | Risk | CLI, flow, status, and browser may disagree. | Derive from canonical order where possible and assert exact parity at each boundary. |
| Artifact promotion and flow-state persistence can partially succeed. | Risk | Artifact and state may disagree after failure. | Define and test transition ordering, retain recoverable diagnostics, and never infer completion from file presence. |
| A preserved old artifact can mask a failed rerun. | Risk | Operators may believe new context succeeded. | Keep latest operation outcome, stage state, and artifact validity separate across CLI/app/web. |
| Generated source references may escape or disclose unsafe paths/content. | Risk | Security boundary and browser safety could fail. | Central containment validation, repository-relative paths, allowlisted bounded preview, escaping, and hostile-content tests. |
| Cancellation races with completion or promotion. | Risk | Duplicate or false terminal outcomes are possible. | Single terminal arbitration, canonical cancellation, promotion only after validation, bounded reconciliation, and race tests. |
| Cleanup may time out. | Risk | Candidate residue or uncertain ownership remains. | Record `interrupted`/`cleanup_uncertain`, retain mutation lock until reconciled, and recover conservatively on restart. |
| Structural validation may be presented as semantic completeness. | Risk | Users may overtrust a partial context pack. | Use precise labels and help; empty findings mean no structural findings only; downstream source inspection remains allowed. |
| Full prompts, excerpts, paths, provider payloads, or stderr may leak through diagnostics. | Risk | Secrets or source content may be exposed. | Projection-time allowlists, redaction, bounds, stream separation, and negative tests. |
| Runtime-backed dogfood is unavailable in Sprint 33. | Risk | Deterministic tests cannot prove real provider behavior. | Record it as deferred to Sprint 34; do not claim or simulate a pass. |

## Implementation Constraints

- Keep all code-context stage semantics in focused files under `internal/sprint`; `internal/platform/runtime` remains product-neutral.
- Preserve dependency direction: `cmd/ultraplan -> internal/app -> internal/sprint`, `internal/web -> internal/app`, and platform packages import no product modules.
- Add `code-context` once, immediately after requirements, and make all projections agree with the canonical order.
- Require valid requirements before execution and a successful validated code-context operation before sprint-index readiness.
- Resolve the implementation target through existing project/execute mechanisms; do not accept a caller-provided target or output path.
- Permit repository reads but no implementation source, test, Git, governed-input, or unrelated sprint-artifact mutation.
- Write only the sprint-root `code-context.md`, using an isolated candidate and existing atomic replacement guarantees.
- Preserve the last valid artifact on runtime, validation, cancellation, cleanup, or persistence failure; separately report the latest attempt outcome.
- Keep prompt preview, validation, status, readiness, config display, artifact preview, and dry-run non-mutating and free of runtime invocation.
- Use embedded prompt/template defaults with optional intentional workspace overrides and synchronized `defaults install` copies.
- Keep validation structural and deterministic; never claim semantic completeness or exclusive repository coverage.
- Preserve explicit model/variant override-source information and apply existing precedence, validation, and redaction.
- Propagate one caller-owned context through real work; use canonical explicit/shutdown cancellation, bounded cleanup, and truthful uncertainty.
- Keep the stage sequential; add no worker pool, event broker, registry, plugin system, generic stage framework, or broad filesystem abstraction.
- Reuse existing typed app use cases and generic web operations; add no route-specific workflow, durable web state, or code-context operation kind.
- Keep browser output server-rendered, progressively enhanced, bounded, escaped, accessible, and fully usable without JavaScript.
- Keep public logs/JSON/SSE metadata bounded and allowlisted; exclude full prompts, excerpts, secrets, unsafe paths, raw payloads, and unrestricted stderr.
- Keep normal tests offline and deterministic; use real temporary repositories only where actual filesystem and Git guarantees are under proof.
- Do not implement Sprint 34 prefix injection, manual skill, broad documentation, or real-runtime dogfood in this sprint.
- Do not introduce repository indexing, RAG/embeddings, caches/cache keys, provider cache-control, parallel manifests, staleness/amendment, content identity, QA/repair, alternate persistence, graphs, cloud, Aren, source mutation, or Git mutation.
- Apply both `system/protocols/architecture-review-protocol.md` and `system/protocols/review-sprint-protocol.md` after implementation.

## Plan Handoff

`plan.md` must execute these decisions. It must not invent architecture, scope, or decisions beyond this document.

The plan must carry forward:

- all seven final decisions and their rejected alternatives
- the `AC-01` through `AC-18` and `C-01` through `C-10` requirement mapping
- all 12 selected contracts
- the expected evidence matrix and exact release commands
- the assumptions, risks, and mitigations
- the mutation, security, cancellation, compatibility, and scope constraints
- Architecture Review and Sprint Review protocols
- explicit Sprint 34 and post-Phase-5 deferrals

## Phase Exit Criteria

- [x] Selected context was read and used.
- [x] All three selected area-specific reasoning documents were completed and summarized.
- [x] Area-specific conclusions are reflected in the final decisions without override.
- [x] Contracts and requirement IDs are mapped to decisions and expected evidence.
- [x] Final decisions fix ownership, stage order, runtime flow, artifact/validation semantics, compatibility, APIs, browser behavior, cancellation, observability, and verification without requiring `plan.md` to reopen architecture.
- [x] Expected evidence is specific and reviewable.
