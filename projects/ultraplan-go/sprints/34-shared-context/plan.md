# Sprint Plan: Shared Context Integration and Grounded-Planning Release

> Project: `ultraplan-go`
> Sprint: `34-shared-context`
> Source: `projects/ultraplan-go/sprints/34-shared-context/reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/roadmap.md`, `projects/ultraplan-go/sprints/34-shared-context/requirements.md`, `projects/ultraplan-go/sprints/34-shared-context/code-context.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/sprints/34-shared-context/sprint-index.md`, `projects/ultraplan-go/sprints/34-shared-context/technical-handbook.md`, `projects/ultraplan-go/sprints/34-shared-context/reasoning/architecture.md`, `projects/ultraplan-go/sprints/34-shared-context/reasoning.md`, `../ultraplan-go/docs/plans/integrated-roadmap.md`, `../ultraplan-go/docs/plans/sprint-code-context-stage.md`, `system/protocols/architecture-review-protocol.md`, `system/protocols/review-sprint-protocol.md`, `builtin:prompts/plan-sprint.md`, `builtin:templates/sprint-plan.md`

This plan executes `reasoning.md`. It does not reopen architecture, scope, byte limits, failure semantics, or integration mechanisms.

## Reasoning Source

- **Sprint Reasoning:** `projects/ultraplan-go/sprints/34-shared-context/reasoning.md`
- **Sprint Index:** `projects/ultraplan-go/sprints/34-shared-context/sprint-index.md`
- **Technical Handbook:** `projects/ultraplan-go/sprints/34-shared-context/technical-handbook.md`
- **Area Reasoning:** `projects/ultraplan-go/sprints/34-shared-context/reasoning/architecture.md`

## Sprint Status

- **Status:** `implementation complete; ready for governed review`
- **Owner:** `UltraPlan execute stage`
- **Start Date:** `2026-08-20`
- **Completion Date:** `2026-08-20`

## Decisions To Execute

| Decision | Source Section | Requirement / Evidence Basis | Execution Implication | Accepted Trade-Off / Rejected Alternative | Risk / Follow-Up |
| --- | --- | --- | --- | --- | --- |
| One Sprint-Owned, Operation-Scoped Byte Contract | `reasoning.md#decision-1-one-sprint-owned-operation-scoped-byte-contract` | `AC-02`-`AC-06`, `AC-12`, `AC-14`, `AC-16`; `C-01`, `C-02`, `C-04`, `C-11`; Architecture and Testing contracts | Add one concrete renderer in `internal/sprint/prompt_context.go`; keep stable framing and the single suffix boundary in `internal/sprint/prompts.go`; preserve raw requirements and code-context slices exactly; build once per top-level operation. | Accept one bounded buffered string and explicit callers. Reject stage-local assembly, runtime middleware, a sprint-aware runtime request, a DI framework, mutable shared bytes, and a persisted prefix cache. | Fixed framing is compatibility-sensitive and every future agent route must opt in and be tested. |
| Deterministic, Contained, Bounded Transient Source Resolution | `reasoning.md#decision-2-deterministic-contained-bounded-transient-source-resolution` | `AC-01`-`AC-03`, `AC-06`, `AC-12`, `AC-16`; `C-02`-`C-04`, `C-06`, `C-08`, `C-09`, `C-11`; Security, Errors, Performance contracts | Parse only validated `Path`, `Lines`, optional `Symbol`, and `Rationale`; resolve sequentially in authored order; preserve duplicates, overlaps, bytes, and line endings; fail closed on unsafe or changed input; cap the whole prefix at 256 KiB and reserve 64 KiB for suffixes. | Accept sequential reads, strict symlink rejection, full failure for one bad reference, and bounded buffering. Reject recursive scans, concurrency without measurement, truncation, omission, prioritization, manifests, indexes, retrieval, and persisted excerpts. | Portable hard-link ancestry remains residual; the byte limit may not match every provider token limit and requires dogfood evidence. |
| Explicit Integration At Every Agent Boundary | `reasoning.md#decision-3-explicit-integration-at-every-agent-boundary` | `AC-03`-`AC-06`, `AC-09`, `AC-12`; `C-01`, `C-03`-`C-06`, `C-08`, `C-09`; Architecture, LLM Runtime, Workflows contracts | Compose through `service.go`, `execute.go`, `review.go`, and `smoke_author.go`; put all stage, output, task, coverage, model, run, attempt, session, timestamp, and diagnostic data after the boundary; preserve one context lineage. | Accept visible call-site wiring and inventory maintenance. Reject hidden middleware, globals, context service locators, interface-layer composition, and eager work for agent-free branches. | A missed route or per-reviewer rebuild would violate prefix equality; captured-request coverage is mandatory. |
| Preserve Workflow Authority, Exact-Once Flow, And Canonical Skill Delegation | `reasoning.md#decision-4-preserve-workflow-authority-exact-once-flow-and-canonical-skill-delegation` | `AC-07`-`AC-11`; `C-05`-`C-07`, `C-09`, `C-11`; CLI Surface, Workflows, Persistence contracts | Keep `flow.go` as the only cumulative ordering authority; run canonical code-context once after requirements; preserve Sprint 33 state/rerun behavior; add one manual-only skill that delegates to `ultraplan sprint "$PROJECT" "$SPRINT" flow --to code-context`. | Accept a non-customizable mechanics path in the materialized skill. Reject skill-owned source selection, duplicate transitions, adapter-specific semantics, and a new workflow engine. | Generated skill text/metadata can drift; flow changes can regress compatibility or duplicate execution. |
| Release Through Layered Evidence, Documentation, And Truthful Dogfood | `reasoning.md#decision-5-release-through-layered-evidence-documentation-and-truthful-dogfood` | All selected contracts; `AC-03`-`AC-16`; `C-03`-`C-11`; PRD Phase 5 KPI and roadmap Sprint 34 release gate | Combine exact-byte unit evidence, real request construction with fake runtimes, flow/state/interface/skill regressions, all named docs, release commands, selected reviews, and disposable gated requirements-to-plan dogfood. | Accept broad fixtures, golden maintenance, release-check cost, and a possible blocked dogfood outcome. Reject substring-only proof, fake-as-real claims, unsent-request claims, production-workspace dogfood, and provider cache-hit claims. | Credentials/runtime/network may block dogfood; a blocker must name the exact unavailable prerequisite. |

## Requirements / Contracts To Satisfy

The `AC-*` and `C-*` labels are defined in `reasoning.md#contracts-applied` and map in order to `requirements.md:44-59` and `requirements.md:73-83`.

| Contract / Requirement ID | Required Behavior | Evidence Planned |
| --- | --- | --- |
| `AC-01`; `C-02`, `C-03` | Keep `code-context.md` reference-only and unchanged; consume exact relative paths/ranges, optional symbols, and rationale without persisting source excerpts. | Parser/validator agreement tests, direct artifact-byte comparisons, and mutation-scope assertions. |
| `AC-02`; `C-01`, `C-04` | Render stable instructions, sprint identity, exact requirements, exact code-context, transient evidence, other stable context, and one final stage boundary in that order. | Full renderer golden, raw framed-slice assertions, order assertions, and exactly-one-boundary checks. |
| `AC-03`, `AC-04`; `C-04` | Produce an identical prefix for unchanged governed inputs and source bytes; exclude all dynamic and stage-specific data. | Cross-stage raw-byte comparisons while varying stage, output, task, reviewer, timestamp, run, attempt, and session facts. |
| `AC-05`, `AC-06`; `C-01`, `C-05`, `C-08` | Prefix every agent-backed sprint-index, handbook, area/final reasoning, plan, execute, independent review, and smoke-author request; label snippets transient/untrusted and permit additional live reads. | Captured fake-runtime requests for every route plus agent-free negative cases and runtime dependency inspection. |
| `AC-07`; `C-05` | Preserve `requirements -> code-context -> sprint-index -> technical-handbook -> area-reasoning -> reasoning -> plan`, with code-context exactly once. | Ordered fake-runtime call log and cumulative-flow state assertions. |
| `AC-08`, `AC-09`; `C-05`, `C-06` | Preserve explicit rerun, failure, cancellation, atomic replacement, last-valid artifact, compatibility, and CLI/TUI/browser agreement. | Sprint 33 regression suites plus app, TUI, and browser fixtures over the same durable state. |
| `AC-10`, `AC-11`; `C-07` | Materialize a manual-only skill that delegates to the canonical flow and preserves dry-run, custom files, force restoration, all-skill inclusion, and synchronized content/metadata. | `internal/workspace/skills_test.go` and command-surface skill tests. |
| `AC-12`; `C-09` | Keep normal evidence offline and deterministic while proving exact bytes, failures, call coverage, flow, cancellation, recovery, and browser behavior. | Temporary repositories, fake runtimes, focused package tests, full normal tests, and race tests. |
| `AC-13`; `C-09`, `C-10` | Run a real implementation-repository requirements-to-plan flow only in a disposable workspace; report pass only after a sent runtime request and validated artifacts/call order/prefix/mutation evidence. | Gated dogfood record with command log, prerequisites, runtime/model, prompt captures, state/artifact validation, before/after identities, and `pass` or exact `blocked`. |
| `AC-14`; `C-11` | Align every named guide with actual order, shared-context behavior, limits, recovery, manual skill, and no cache-hit guarantee. | Path-by-path documentation review against CLI help, tests, renderer framing, and dogfood procedure. |
| `AC-15` | Pass normal, race, vet, build, and whitespace release gates in `../ultraplan-go`. | Complete command output retained for review. |
| `AC-16`; `C-01`-`C-11` | Introduce none of the excluded cache, retrieval, index, manifest, staleness, provenance, QA/repair, persistence, framework, hosted, issue, or Git-mutation capabilities. | Dependency and diff review under Architecture Review and Sprint Review. |

## Tasks

- [x] **Task 1: Establish The Shared Prefix Byte Contract**
  > Executes: `Decision 1`; `AC-02`-`AC-04`, `AC-06`, `AC-12`; `C-01`, `C-02`, `C-04`
  - [x] Add `../ultraplan-go/internal/sprint/prompt_context.go` with one concrete, context-aware renderer returning `string, error`; reuse the exact requirements and code-context values from the existing planning-input storage boundary rather than creating persistence or a second artifact model.
  - [x] Add stable shared instructions, fixed external artifact frames, transient-evidence framing, and one constant stage-specific boundary marker to `../ultraplan-go/internal/sprint/prompts.go`; ensure the boundary is the final prefix bytes.
  - [x] Compose in the decided order and append raw artifact strings without trimming, newline insertion, line-ending normalization, reserialization, or regeneration.
  - [x] Keep project and sprint identity stable in the prefix; keep stage/template/output/model/run/task/reviewer/smoke facts out of it.
  - [x] Add `../ultraplan-go/internal/sprint/prompt_context_test.go` fixtures for LF, CRLF, missing final newlines, exact framed slices, order, exactly one boundary, unchanged-input determinism, operation-scoped reuse, and dynamic-data exclusion.
  - [x] Stop and record a deviation before continuing if exact stored bytes cannot pass unchanged through the existing planning-input API or if a platform/runtime contract change appears necessary. No deviation was required; `internal/platform/runtime` remains unchanged.

- [x] **Task 2: Resolve Transient Source Evidence Safely And Deterministically**
  > Executes: `Decision 2`; `AC-01`-`AC-03`, `AC-06`, `AC-12`, `AC-16`; `C-02`-`C-04`, `C-06`, `C-08`, `C-09`
  - [x] Parse only the Sprint 33 validated reference fields from `code-context.md`; test agreement with `ValidateCodeContextContent` and preserve document order, duplicates, overlaps, selected bytes, and line endings.
  - [x] Resolve only non-empty repository-relative paths under the supplied implementation root; reject absolute, volume-qualified, lexical escape, canonical escape, symlink-component, non-regular, missing, invalid UTF-8, malformed/zero/descending/out-of-range, changed-read, and uncertain-containment cases.
  - [x] Open and verify the regular file handle, use bounded range reads with checks before/during/between reads, and repeat identity/canonical containment checks before accepting bytes; start no goroutines and preserve the caller context.
  - [x] Frame each excerpt with stable path/range/symbol/rationale metadata and explicit language that it is untrusted transient prepared evidence, is not stored in `code-context.md`, and does not limit further repository inspection.
  - [x] Enforce checked arithmetic, a 256 KiB maximum for the complete prefix, and a 64 KiB stage-suffix reserve; count fixed text, both exact artifacts, every duplicate/overlap, and excerpt framing; fail before runtime launch without truncating, omitting, sorting, summarizing, or prioritizing entries.
  - [x] Preserve `%w` causes and cancellation identity while exposing safe sprint-owned categories for path, containment, file kind, missing source, range, changed read, encoding, and budget failures; diagnostics may include safe relative path/range and byte counts but no source, full prompt, secret, or default absolute path.
  - [x] Extend `prompt_context_test.go` with temporary-repository cases for valid/invalid path and range classes, authored order, duplicates/overlaps, cancellation, symlinks/non-regular files, deterministic file replacement detection, exact budget fit and one-byte overflow, no recursion, and no mutation.
  - [x] Stop before runtime invocation on any unresolved reference or budget uncertainty; do not weaken containment to accommodate symlinked or changing sources.

- [x] **Task 3: Integrate Planning Stages Through One Explicit Boundary**
  > Executes: `Decisions 1 and 3`; `AC-03`-`AC-06`, `AC-12`; `C-01`, `C-04`, `C-05`, `C-08`, `C-09`
  - [x] Update `../ultraplan-go/internal/sprint/service.go` so sprint-index, technical-handbook, area reasoning, final reasoning, and plan agent-bound prompt construction prepares the common prefix through the one renderer and appends the existing stage renderer output only after the boundary.
  - [x] Include previews/dry runs of agent-bound prompts in the same composition path while keeping genuinely skipped area reasoning and other agent-free branches runtime-free and free of unnecessary source preparation.
  - [x] Resolve the existing implementation target through current project/execute mechanisms and carry the same operation context from source preparation into runtime invocation.
  - [x] Preserve workspace/project prompt overrides, templates, mutation contracts, output validation, stage model selection, prompt identity, and runtime metadata after the boundary.
  - [x] Cover the planning preview/runtime composition boundary through shared exact-prefix fixtures and the existing handbook/reasoning/plan request-construction suites; additional-live-read permission is asserted in the shared renderer fixture.
  - [x] Stop and inventory the missing call path if any planning agent request cannot use the common composition operation; no missing planning path was found and no runtime decorator was added.

- [x] **Task 4: Integrate Execute, Conformance Review, And Smoke Authoring Without Orchestration Changes**
  > Executes: `Decisions 1 and 3`; `AC-03`-`AC-06`, `AC-08`, `AC-12`; `C-01`, `C-03`-`C-06`, `C-08`, `C-09`
  - [x] Update `../ultraplan-go/internal/sprint/execute.go` to prepare one prefix per execute operation and place every task prompt, resume/session preamble, task ID, attempt, output/evidence contract, and model/runtime detail after the boundary.
  - [x] Update `../ultraplan-go/internal/sprint/review.go` to prepare one immutable prefix before independent reviewer fan-out and reuse it for every contract and handbook request; retain frozen scope, coverage sorting, concurrency, structured-result repair, citation rules, aggregation, and product-owned verdicts.
  - [x] Update `../ultraplan-go/internal/sprint/smoke_author.go` to add the same prefix only when smoke suite authoring invokes an agent; retain harness-only write policy, protected snapshots, authoring selection, cancellation, and recovery.
  - [x] Keep `../ultraplan-go/internal/platform/runtime` and agentwrap unchanged and sprint-neutral; pass the complete composed prompt through the existing generic request.
  - [x] Extend cross-route, review, and smoke fixtures with exact-prefix comparisons, fan-out reuse, dynamic-suffix placement, cancellation behavior, and agent-free negative behavior without weakening existing assertions.
  - [x] Stop and record a reasoning deviation if integration requires changing review fan-out/verdict authority, smoke scope/mutation authority, execute task/session semantics, or the generic runtime request contract. No such change or deviation was required.

- [x] **Task 5: Lock In Exact-Once Flow And Sprint 33 State Guarantees**
  > Executes: `Decision 4`; `AC-07`-`AC-09`, `AC-12`; `C-05`, `C-06`, `C-09`
  - [x] Audit `../ultraplan-go/internal/sprint/flow.go`; no implementation adjustment was required because the canonical order already contains code-context exactly once and dispatches the existing operation.
  - [x] Add an explicit ordered `flowStages(StagePlan)` assertion in `sprint_index_test.go`; existing direct/cumulative, state, and request-construction fixtures retain operation equivalence and prove prompt previews do not advance state.
  - [x] Preserve and rerun `code_context_test.go` and `verify_test.go` coverage for validation/runtime failure, cancellation/deadline, candidate cleanup, atomic promotion, last-valid-artifact restoration, state-write failure, recovery, and pre-code-context compatibility.
  - [x] Run the shared app, TUI, and browser interface fixtures so CLI/JSON/TUI/browser continue to agree without adapter-owned prompt logic.
  - [x] Stop on any interface divergence or changed state transition; no divergence or transition regression was observed.

- [x] **Task 6: Materialize The Manual Canonical Code-Context Skill**
  > Executes: `Decision 4`; `AC-10`, `AC-11`; `C-05`, `C-07`, `C-11`
  - [x] Register `ultraplan-code-context` in `../ultraplan-go/internal/workspace/skills.go` using existing `StageSkill` and materialization conventions; keep it manual-only.
  - [x] Make the generated workflow delegate exactly to `ultraplan sprint "$PROJECT" "$SPRINT" flow --to code-context`, with the selected project and sprint substituted by the materialized skill workflow; do not duplicate target resolution, repository inspection, prompt composition, validation, state transitions, or artifact promotion.
  - [x] Preserve selection by stage and full name, dry-run non-mutation, deterministic output, customization preservation, `--force` restoration, all-skill inclusion, and synchronized `SKILL.md`/`agents/openai.yaml` content and metadata.
  - [x] Extend `../ultraplan-go/internal/workspace/skills_test.go` and retain the existing app skill command tests for resolution, count/all inclusion, manual-only metadata, exact delegation, dry run, normal materialization, custom files, force, isolation, and embedded/generated synchronization.
  - [x] Stop if the skill requires a second product implementation or direct TUI/web/CLI-internal invocation instead of canonical command delegation. It delegates only through the documented canonical CLI command.

- [x] **Task 7: Align Release Documentation With Executable Behavior**
  > Executes: `Decision 5`; `AC-06`, `AC-09`-`AC-11`, `AC-13`, `AC-14`; `C-05`, `C-07`, `C-08`, `C-10`, `C-11`
  - [x] Update `../ultraplan-go/README.md` with the grounded-planning order, reference-only artifact, shared-prefix guarantee, transient live evidence, live inspection permission, manual skill, and supported release boundary.
  - [x] Update `../ultraplan-go/docs/cli-reference.md` and `../ultraplan-go/docs/user-guide.md` with code-context commands, exact-once cumulative placement, status, generation/inspection/rerun, downstream reuse, and actionable failure behavior.
  - [x] Update `../ultraplan-go/docs/architecture.md` with `internal/sprint` ownership, exact rendering order and boundary, operation-scoped reuse, 256 KiB prefix/64 KiB suffix reserve, contained source resolution, and generic runtime separation.
  - [x] Update `../ultraplan-go/docs/recovery.md` with cancellation, failed generation/resolution, atomic rerun, restart, last-valid artifact, and browser recovery semantics.
  - [x] Update `../ultraplan-go/docs/planning-smoke.md` with deterministic fake coverage and the gated disposable real-repository requirements-to-plan procedure, exact prerequisites, mutation checks, observed prompt/call evidence, and truthful `pass`/`blocked` rules.
  - [x] Update `../ultraplan-go/docs/stage-skills.md`, `../ultraplan-go/docs/local-web.md`, and `../ultraplan-go/internal/workspace/scaffold/templates/README.md` with manual canonical delegation and shared interface readiness/progress/findings/artifact/rerun/cancellation/recovery behavior.
  - [x] Compare documentation with command definitions and executable tests; no provider cache-hit guarantee, measurement, ownership, or dependency is claimed.

- [x] **Task 8: Produce Layered Release And Review Evidence**
  > Executes: `Decision 5`; `AC-03`-`AC-16`; `C-03`-`C-11`
  - [x] Run focused renderer, planning, execute, review, smoke-author, flow, state, interface, and skill suites after their corresponding increments; retain failures and fixes as execution evidence rather than updating goldens blindly.
  - [x] Run the full normal, race, vet, build, and whitespace release commands from `../ultraplan-go`; all passed on the final tree.
  - [x] Evaluate the documented gated requirements-to-plan dogfood. It was not invoked because the active manual execute skill permits only project status, sprint status, and execute-prompt UltraPlan commands; invoking `flow --to plan` would exceed the granted execute-stage CLI permission boundary.
  - [x] Record the gated dogfood as `blocked` by that exact command-permission prerequisite. No fake-runtime, constructed-only, unsent, denied, or artifact-only evidence is represented as real-runtime success.
  - [x] Prepare the implementation diff and verification/blocker evidence for Architecture Review and Sprint Review; dependency and scope inspection found no excluded subsystem, new dependency, route-specific prompt semantics, runtime-layer sprint coupling, or automatic Git mutation.
  - [x] Record any deviation from `reasoning.md` before implementation continues. No implementation deviation is recorded; the truthful gated dogfood blocker is an evidence limitation, not a product-design deviation.

## Evidence Checklist

- [x] Exact renderer and focused raw-slice tests prove byte preservation, fixed order, one boundary, and LF/CRLF/no-final-newline behavior.
- [x] Temporary-repository tests prove parser agreement, deterministic ordered source resolution, duplicates/overlaps, containment, non-regular/replacement handling, UTF-8/range validation, cancellation, 256 KiB/64 KiB budget edges, no recursion, and no mutation.
- [x] Captured previews/requests prove identical prefix bytes for sprint-index, handbook, area/final reasoning, plan, execute, every independent review request, and agent-backed smoke authoring.
- [x] Captured requests prove stage/output/task/reviewer/smoke/model/run/session/attempt facts begin after the boundary or remain runtime metadata.
- [x] Agent-free conditional paths produce no fabricated runtime request or unnecessary context work.
- [x] Flow evidence proves canonical order and exact-once code-context scheduling/execution through the existing canonical operation.
- [x] Sprint 33 regression evidence proves cancellation, failure, atomic replacement, last-valid preservation, compatibility, state truthfulness, and browser recovery remain intact.
- [x] Skill evidence proves manual-only canonical delegation, dry run, customization preservation, force restoration, all-skill inclusion, and synchronized generated files.
- [x] CLI/JSON/TUI/browser fixtures agree on readiness, progress, findings, artifact, rerun, cancellation, recovery, and next action through shared app operations.
- [x] Documentation updates cover every required path, match executable behavior, and make no provider cache-hit claim.
- [x] Gated dogfood is recorded precisely as blocked by the active execute-stage UltraPlan CLI allowlist; fake or unsent work is not represented as real-runtime success.
- [x] Release commands pass and complete outcomes are recorded in `execute.md` for review.
- [x] Manual Architecture Review preflight covers ownership, dependency direction, explicit call-site completeness, containment, state/side effects, simplicity, errors, observability, testing, performance, and exclusions; it does not replace canonical review.
- [x] Manual Sprint Review preflight covers every selected contract, the handbook, all five decisions, plan execution, verification evidence, deviations, and truthful dogfood status; it does not create or replace `review.md`.
- [x] No deviation from `reasoning.md` was required.

## Verification Commands

Run repository commands from `../ultraplan-go` unless another working directory is stated.

| Check | Command | Expected Result |
| --- | --- | --- |
| Sprint renderer and route coverage | `go test ./internal/sprint` | Exact-byte, containment, planning, execute, review, smoke-author, flow, state, cancellation, and recovery tests pass offline. |
| Skill and interface regressions | `go test ./internal/workspace ./internal/app ./internal/tui ./internal/web` | Skill materialization and CLI/JSON/TUI/browser agreement pass without route-specific prompt behavior. |
| Full normal suite | `go test ./...` | All packages pass with deterministic normal-test dependencies. |
| Race suite | `go test -race ./...` | All packages pass under the race detector, including fan-out and cancellation paths. |
| Static analysis | `go vet ./...` | No vet findings. |
| Binary build | `go build ./cmd/ultraplan` | The production CLI builds successfully. |
| Whitespace and patch validity | `git diff --check` | No whitespace errors or malformed patch lines. |
| Gated real-repository dogfood | `go run ./cmd/ultraplan --workspace "$DOGFOOD_WORKSPACE" sprint "$DOGFOOD_PROJECT" "$DOGFOOD_SPRINT" flow --to plan` | In the disposable workspace defined by `docs/planning-smoke.md`, a real request is sent and valid context/plan, exact-once order, shared-prefix captures, and protected-location identity checks pass; otherwise the exact prerequisite is recorded as blocked. |
| Plan conformance review | `ultraplan sprint ultraplan-go 34-shared-context review` | The selected Architecture Review and Sprint Review coverage produces a current validated verdict from the implementation diff and retained evidence; no smoke claim is inferred. |

## Risks And Blockers

| Risk / Blocker | Source | Mitigation | Status |
| --- | --- | --- | --- |
| Renderer parser diverges from Sprint 33 validation grammar. | `reasoning.md#potential-technical-debt`; `code-context.md` | Parse only four validated fields, share canonical fixtures, and test parser/validator agreement; add no compatibility grammar. | `mitigated` |
| Trimming, newline normalization, or frame drift corrupts exact artifact reuse. | `reasoning.md#assumptions-and-risks` | Append raw stored strings between external frames; pair goldens with direct slice/order/boundary assertions. | `mitigated` |
| One agent-backed path bypasses composition or rebuilds during fan-out. | Decisions 1 and 3 | Enumerate all required routes, capture each request, assert one boundary, and prepare once at operation entry. | `mitigated` |
| Source changes during a read or containment remains uncertain. | Decision 2; Architecture reasoning | Use canonical pre/post checks, handle identity checks, bounded reads, symlink rejection, and fail closed before runtime. | `mitigated` |
| Hard links evade portable ancestry proof. | `reasoning.md#potential-technical-debt` | Keep reads contained, regular-file-only, bounded, read-only, and permission-limited; retain as an explicit residual risk. | `residual` |
| Prompt-like source text influences the agent. | `reasoning.md#assumptions-and-risks` | Label evidence untrusted, frame every excerpt, preserve least-privilege runtime policy, and never interpret source in product code. | `mitigated` |
| 256 KiB prefix plus 64 KiB reserve does not fit a configured provider or existing transport. | Decision 2; `code-context.md#open-questions` resolved by `reasoning.md` | Enforce the decided product bound, record observed sizes, dogfood representative suffixes, and change limits only through a reviewed decision/test/doc update. | `residual; dogfood blocked` |
| Local filesystem reads delay cancellation. | `reasoning.md#assumptions-and-risks` | Use bounded chunks/ranges and context checks before, during, between references, and before return; start no detached work. | `mitigated` |
| Golden updates bless accidental byte incompatibility. | Decision 5 | Require focused raw-slice/order/boundary/dynamic-exclusion assertions and deliberate fixture diff review. | `mitigated` |
| Skill text or documentation drifts from canonical CLI behavior. | Decisions 4 and 5 | Generate from one embedded definition, test synchronization, compare docs with command definitions and executable tests. | `mitigated` |
| Gated runtime prerequisites are unavailable. | Decision 5 | The active execute-stage UltraPlan CLI allowlist excludes `flow --to plan`; record the exact permission blocker and rerun the disposable procedure under an authorized gate. | `blocked` |
| Later roadmap work enters Sprint 34. | Requirements Non-Goals; roadmap stop point | Reject cache, index, retrieval, staleness, identity/provenance, QA/repair, persistence, graph, cloud, issue, and framework additions in diff/review. | `closed` |

## Evidence Boundaries And Reconciliation

- The current `requirements.md` and validated `reasoning.md` supersede the older copied-excerpt language in `../ultraplan-go/docs/plans/sprint-code-context-stage.md`: Sprint 34 keeps `code-context.md` reference-only and resolves selected source text transiently without rewriting or persisting it.
- `projects/ultraplan-go/roadmap.md` remains authoritative for Sprint 34 sequencing and scope. The implementation plans are design inputs only and do not authorize post-Sprint-34 identity, provenance, QA/repair, retrieval, persistence, graph, cloud, or Aren work.
- The 11 selected final study reports were not reopened for this plan because `technical-handbook.md`, `reasoning/architecture.md`, and validated `reasoning.md` contain sufficient cited comparative evidence and final sprint-specific decisions. Reopen a linked report only if implementation exposes a concrete unresolved question not decided by reasoning.
- Other implementation plans under `../ultraplan-go/docs/plans/` were inspected for overlap and omitted because they govern local-server history or explicitly deferred retrieval, QA/repair, persistence, and graph work; including them would risk scope expansion.
- Current implementation source was resolved only through the validated `code-context.md` references and a focused test-location inventory needed to name integration/regression files; no new architecture decision was inferred from source inspection.

## Review Inputs

Review should use:

- `projects/ultraplan-go/project-index.md`
- `projects/ultraplan-go/sprints/34-shared-context/requirements.md`
- `projects/ultraplan-go/sprints/34-shared-context/code-context.md`
- `projects/ultraplan-go/sprints/34-shared-context/sprint-index.md`
- `projects/ultraplan-go/sprints/34-shared-context/technical-handbook.md`
- `projects/ultraplan-go/sprints/34-shared-context/reasoning/architecture.md`
- `projects/ultraplan-go/sprints/34-shared-context/reasoning.md`
- this `projects/ultraplan-go/sprints/34-shared-context/plan.md`
- `system/protocols/architecture-review-protocol.md`
- `system/protocols/review-sprint-protocol.md`
- implementation diff in `../ultraplan-go`
- focused and release-command output
- gated dogfood evidence or its exact blocker
- any recorded reasoning deviation or explicit deferral

## Execution Log

| Date / Step | Action | Evidence / Notes |
| --- | --- | --- |
| `2026-08-20 / planning` | Created the implementation plan from validated requirements, final/area reasoning, selected evidence, governing product documents, mandatory implementation plans, and selected review protocols. | Plan stage only; no implementation, review, smoke, issue, Git, or runtime work executed. |
| `2026-08-20 / Tasks 1-4` | Implemented the concrete shared renderer, fail-closed transient source resolution, and explicit planning/execute/review/smoke-author composition. | Exact artifact/source bytes, one final boundary, route equality, fan-out reuse, cancellation, replacement detection, budget limits, and legacy behavior are covered by `internal/sprint` tests. Generic runtime code was not changed. |
| `2026-08-20 / Tasks 5-7` | Audited canonical flow, added exact-order coverage, materialized the manual code-context skill, and aligned all required guides. | Flow retains code-context exactly once; skill delegates exactly to the canonical flow; workspace/app/TUI/web suites pass. |
| `2026-08-20 / release verification` | Ran the complete plan verification matrix from `../ultraplan-go`. | `go test ./internal/sprint`, interface/skill packages, `go test ./...`, `go test -race ./...`, `go vet ./...`, `go build ./cmd/ultraplan`, and `git diff --check` all passed. |
| `2026-08-20 / gated dogfood` | Evaluated the documented disposable real-runtime procedure without invoking it. | **Blocked:** the active manual execute skill authorizes only project status, sprint status, and execute-prompt UltraPlan commands; `flow --to plan` is outside that permission boundary. No unsent/fake evidence is claimed as a pass. |
| `2026-08-20 / Architecture Review preflight` | Applied the selected architecture protocol manually to the implementation and dependency diff. | **Approve with comments.** Sprint ownership, explicit call sites, generic runtime separation, failure boundaries, state neutrality, test coverage, and exclusions are sound. Residual hard-link ancestry and provider-capacity risks remain documented. This is preflight evidence, not canonical `review.md`. |
| `2026-08-20 / Sprint Review preflight` | Applied the selected sprint-review protocol manually against governed inputs, decisions, tasks, diff, and retained verification evidence. | **Pass with gated evidence blocked.** Implementation and deterministic evidence cover the selected contracts; no blocker/high implementation finding was identified. Governed reviewer fan-out and `review.md` remain the next stage. |

## Completion Criteria

- [x] All tasks are complete or explicitly blocked with requirement and reasoning impact recorded.
- [x] All five reasoning decisions and every `AC-*`/`C-*` mapping have implementation and evidence.
- [x] Verification commands were run or the exact gated-dogfood permission blocker is documented.
- [x] Exact-byte, containment, route, flow, state, interface, skill, documentation, and truthful blocked-dogfood evidence satisfies `reasoning.md`.
- [x] No excluded architecture or later-roadmap capability entered the implementation.
- [x] Architecture Review and Sprint Review can evaluate conformance without guessing intent.
- [x] `review.md` can distinguish implementation failure, missing evidence, and the unavailable gated permission prerequisite truthfully.
