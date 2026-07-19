# Sprint Plan: Automated Sprint Review

> Project: `ultraplan-go`
> Sprint: `26-review-stage`
> Source: `reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/26-review-stage/requirements.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/sprints/26-review-stage/sprint-index.md`, `projects/ultraplan-go/sprints/26-review-stage/technical-handbook.md`, `projects/ultraplan-go/sprints/26-review-stage/reasoning/architecture.md`, `projects/ultraplan-go/sprints/26-review-stage/reasoning.md`, `builtin:templates/sprint-plan.md`

This plan executes `reasoning.md`. It does not reopen architecture, scope, concurrency semantics, verdict policy, persistence semantics, or interface ownership. The selected handbook was sufficient for implementation planning, so linked final study reports were not reopened and no selected evidence was omitted.

## Reasoning Source

- **Sprint Reasoning:** `reasoning.md`
- **Sprint Index:** `sprint-index.md`
- **Technical Handbook:** `technical-handbook.md`
- **Area Reasoning:** `reasoning/architecture.md`

## Sprint Status

- **Status:** `complete`
- **Owner:** `Codex`
- **Start Date:** `2026-07-19`
- **Completion Date:** `2026-07-19`

## Decisions To Execute

| Decision | Source Section | Execution Implication |
| --- | --- | --- |
| Keep Review Product-Owned And Share Typed Use Cases | `reasoning.md#decision-1-keep-review-product-owned-and-share-typed-use-cases` | Keep review policy and durable review behavior in `internal/sprint`; expose terminal-independent requests, progress, and results through `internal/app`; keep CLI and TUI as adapters. |
| Resolve, Validate, Freeze, And Fingerprint One Review Manifest | `reasoning.md#decision-2-resolve-validate-freeze-and-fingerprint-one-review-manifest` | Resolve catalog selections and execute prerequisites before runtime work, reject invalid or escaping paths, capture one immutable invocation scope, fingerprint every governed input, and reject drift before persistence. |
| Run Independent Structured Reviewers With Bounded Collect-All Fan-Out | `reasoning.md#decision-3-run-independent-structured-reviewers-with-bounded-collect-all-fan-out` | Issue exactly one structured request per selected contract and one for the handbook; bound concurrency, continue healthy siblings after unit failures, and cancel/await all started work for invocation-wide termination. |
| Validate Evidence And Compute The Verdict Deterministically | `reasoning.md#decision-4-validate-evidence-and-compute-the-verdict-deterministically` | Validate versioned reviewer results and citations, run product checks independently of model prose, sort by stable keys, and let product code alone produce an allowed verdict. |
| Persist Only Complete Current Evidence And Reconcile Freshness | `reasoning.md#decision-5-persist-only-complete-current-evidence-and-reconcile-freshness` | Stage and validate complete Markdown before atomic replacement, update review flow state atomically afterward, preserve prior complete evidence on failure, and reconcile artifact/state fingerprints in status. |
| Expose One Truthful CLI, JSON, Status, Flow, And TUI Surface | `reasoning.md#decision-6-expose-one-truthful-cli-json-status-flow-and-tui-surface` | Drive review, validation, preview, dry-run, status, flow, JSON, and TUI from shared typed facts; keep machine stdout uncontaminated and non-execution paths runtime-free. |
| Enforce Read-Only Review And Safe Diagnostics Mechanically | `reasoning.md#decision-7-enforce-read-only-review-and-safe-diagnostics-mechanically` | Require enforceable read-only capabilities, canonical contained paths, narrow permitted writes, central redaction, and omission of raw native payloads and unrestricted runtime output. |
| Ship Embedded Assets And Verify Behavior At Product Boundaries | `reasoning.md#decision-8-ship-embedded-assets-and-verify-behavior-at-product-boundaries` | Extend the existing embedded-default mechanism for review prompt/template behavior, preserve explicit override/materialization rules, and prove the feature through pure, fake-runtime, persistence, command, TUI, race, build, and protocol checks. |

## Requirements / Contracts To Satisfy

Acceptance criteria IDs follow their order in `requirements.md:42-62`. Protocol IDs follow `sprint-index.md:78-81`.

| Contract / Requirement ID | Required Behavior | Evidence Planned |
| --- | --- | --- |
| AC-01; Security; Workflows | `review` requires valid prerequisites through execute and rejects missing, duplicate, unknown, unreadable, or escaping selected paths before runtime work. | Preflight fixture matrix, including canonical/symlink escape cases and zero runtime calls on failure. |
| AC-02; Architecture; Security | Resolve every selected contract and protocol dynamically through `project-index.md`; keep review semantics in `internal/sprint`. | Catalog resolution tests plus RP-01 dependency/ownership inspection. |
| AC-03; Persistence And Migrations; Workflows | Fingerprint governed sprint inputs, selected scope, target identity, changed paths, `plan.md`, `execute.md`, and `.run-state.json` when present. | Fingerprint table tests proving each input changes identity and absent/present run-state handling is explicit. |
| AC-04; Persistence And Migrations; Workflows | Freeze one manifest for workers and detect governed or target drift before persistence. | Mutation-during-run fixture yields stale failure and preserves prior complete evidence. |
| AC-05; Architecture; LLM Runtime; Performance; Workflows | Run one independent selected-contract request plus one handbook request with bounded collect-all fan-out owned by UltraPlan. | Fake-runtime request keys/count, peak concurrency, sibling collection, and no-subagent assertions. |
| AC-06; Configuration; LLM Runtime; Security | Apply central model/config precedence, context, bounded retries, read-only permissions, path containment, and safe diagnostics. | Config-source, request-policy, capability, retry, timeout, cancellation, and redaction tests. |
| AC-07; Errors; LLM Evaluation / Cost / Safety; Observability | Missing, malformed, failed, timed-out, cancelled, or runtime-failed reviewer units remain classified and cannot become pass. | Terminal-class matrix across sprint result, app DTO, text, JSON, status, and TUI. |
| AC-08; Errors; LLM Evaluation / Cost / Safety | Deterministic checks cover decisions, plan execution, approved verification evidence, coverage, citations, deviations, and missing evidence. | Pure check tests with positive, missing, stale, malformed, and deviation fixtures. |
| AC-09; Documentation; Persistence And Migrations | Completion requires valid sprint-root `review.md` with required sections, no placeholders, complete coverage, valid citations, and an allowed verdict. | Markdown validator tests and deterministic fixture/golden output. |
| AC-10; LLM Evaluation / Cost / Safety | Product code alone computes `pass`, `pass_with_findings`, `fail`, or `blocked`; applicable blocker/high findings force `fail`. | Verdict truth-table tests and shuffled-input stability tests. |
| AC-11; Persistence And Migrations; Security | Rerun atomically replaces only `review.md` and review flow-state fields without unrelated artifact, source, test, config, or Git mutation. | Failure-injected writes and before/after filesystem snapshots. |
| AC-12; Persistence And Migrations; Observability | Failed or incomplete attempts preserve the last complete `review.md` while exposing the current failed attempt and safe diagnostics. | Stage/write/flush/close/rename/state-update failure fixtures and status assertions. |
| AC-13; Persistence And Migrations; Workflows | Add `review` after `execute` and record execution status, verdict, artifact, timestamp, fingerprint, staleness, and diagnostics separately. | Strict state schema/load/transition/reconciliation tests, including legacy and mismatch fixtures. |
| AC-14; Errors; Observability; Workflows | Status, review validation, review prompt preview, review-target flow integration, and review JSON output agree on readiness, state, verdict, freshness, and next action. | Shared fixture assertions across all required review command surfaces. |
| AC-15; CLI Surface; Errors; Observability | Text and JSON output truthfully map all review outcomes and keep execution status distinct from conformance verdict. | Command/exit mapping tests and JSON structural assertions. |
| AC-16; CLI Surface; Configuration; Performance | Preview and dry-run expose model source, count, bound, scope, selections, writes, and mutation policy without runtime initialization or writes. | Lazy-factory spy and filesystem snapshot tests for CLI and TUI confirmation data. |
| AC-17; Architecture; CLI Surface; Observability | TUI uses shared typed app use cases for readiness, preview, dry-run, confirmation, progress, cancellation, result, findings, verdict, and artifact preview. | Deterministic TUI model/update/view tests using the same result fixtures as CLI. |
| AC-18; Architecture; CLI Surface | TUI does not call CLI handlers, parse terminal/native-provider output, or own alternate review persistence. | RP-01 dependency inspection and fake-use-case TUI tests. |
| AC-19; Documentation | Review prompt and output template work as embedded defaults; intentional overrides and `defaults install` remain explicit. | Missing-override fallback, valid override, invalid/unreadable override, and materialization tests. |
| AC-20; Security; Testing | Fake-runtime tests cover success, findings, fail/blocked cases, malformed/missing results, citation failures, stale input, timeout, cancellation, redaction, and atomic-write preservation. | Required scenario matrix in sprint/app/TUI tests; real OpenCode remains optional and gated. |
| AC-21; Testing; Verification gate | Repository tests, race tests, and CLI build pass from `../ultraplan-go`. | Recorded results for `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan`. |
| RP-01 Architecture Review | Verify module ownership, dependency direction, generic runtime boundary, shared app operations, and interface-only TUI behavior. | Completed `system/protocols/architecture-review-protocol.md` evidence referenced by the generated review. |
| RP-02 Sprint Review | Verify all selected contracts and handbook guidance, implementation evidence, deterministic verdict, and current valid `review.md`. | Completed `system/protocols/review-sprint-protocol.md` evidence and validated current artifact. |

All 13 selected contracts are represented above: Architecture, CLI Surface, Configuration, Documentation, Errors, LLM Evaluation / Cost / Safety, LLM Runtime, Observability, Performance, Persistence And Migrations, Security, Testing, and Workflows.

## Tasks

- [x] **Task 1: Extend Review Domain And Flow State**

  > Executes: Decisions 1 and 5; AC-09, AC-10, AC-13, AC-14; Architecture; Persistence And Migrations

  - [x] Extend `internal/sprint/domain.go` with the `review` stage and typed review execution status, allowed verdicts, applicability, severity, coverage, finding, diagnostic, fingerprint, artifact, staleness, progress, and result facts required by all surfaces.
  - [x] Keep execution status separate from verdict in domain objects and versioned `flow-state.json`; do not encode a failing verdict as failed execution or a successful runtime as completed review.
  - [x] Extend stage ordering, artifact paths, strict loading, status initialization, and flow transitions so `review` follows valid `execute` and no smoke/verify/issues behavior enters this sprint.
  - [x] Define explicit supported legacy-state handling and unsupported-version diagnostics using existing state conventions; do not silently ignore new or malformed review fields.
  - [x] Add table tests for stage order, allowed values, strict state loading, legacy/manual review non-current behavior, malformed state, and execution-status/verdict combinations.

- [x] **Task 2: Resolve And Freeze The Review Manifest**

  > Executes: Decision 2; AC-01 through AC-04, AC-16; Configuration; Security; Workflows

  - [x] Add sprint-owned readiness/preflight behavior that validates governed artifacts through execute and resolves selected contracts and protocols only through the existing project catalog API.
  - [x] Reject missing, duplicate, unknown, unreadable, non-canonical, symlink-escaping, or otherwise out-of-policy selected entries before constructing the runtime.
  - [x] Capture immutable workspace-relative paths/content identities for requirements, sprint index, handbook, selected area reasoning, final reasoning, plan, execute summary/state, selected contracts/protocols, target implementation identity, and explicit changed-path scope.
  - [x] Resolve and validate effective review model source, timeout, positive per-invocation concurrency, agentwrap retry policy, required capabilities, permission policy, reviewer count, and permitted writes before scheduling; do not choose a new numeric concurrency default in this sprint plan.
  - [x] Compute one deterministic fingerprint from canonical manifest identities and contents, and provide reviewer-specific immutable projections without worker rereads becoming authoritative.
  - [x] Recheck governed and target identity immediately before persistence; classify drift as stale/non-current and preserve prior evidence.
  - [x] Add deterministic fingerprint and preflight fixtures showing every governed input and changed-path component affects identity and all invalid scope classes cause zero runtime calls and zero writes.

- [x] **Task 3: Orchestrate Independent Read-Only Reviewers**

  > Executes: Decisions 3 and 7; AC-05 through AC-07, AC-11, AC-20; LLM Runtime; Performance; Security

  - [x] Implement sprint-owned review orchestration in focused `internal/sprint/review.go` behavior using the existing generic agentwrap-backed runtime boundary.
  - [x] Render exactly one versioned structured request for each selected contract and one for the technical handbook; use selected protocols as constraints, not extra reviewer units.
  - [x] Set contained work directories, model/timeout metadata, structured validation, required capabilities, explicit read-only permissions, and only the read/search/list operations required by review.
  - [x] Enforce a configured positive per-invocation concurrency bound; give each started run one owner that drains events, calls `Wait`, and returns an isolated terminal coverage result.
  - [x] Continue healthy siblings after an ordinary reviewer failure, panic, missing result, or malformed result; convert worker panic to a safe failed coverage outcome rather than process termination.
  - [x] On user cancellation, deadline, or invocation-wide infrastructure invalidation, stop scheduling, cancel active runs, await bounded cleanup, and preserve cleanup facts separately from the primary outcome.
  - [x] Use agentwrap-owned bounded retries and OpenCode adaptation; do not add runtime subagents, direct OpenCode execution, custom retry loops, global throttles, or durable partial-result caches.
  - [x] Add fake-runtime tests for exact request count/keys, peak activity, sibling continuation, retry, timeout, cancellation, panic, event draining, awaited completion, unsupported safety capability, and race/leak behavior.

- [x] **Task 4: Validate Evidence And Synthesize The Verdict**

  > Executes: Decision 4; AC-07 through AC-10, AC-20; Errors; LLM Evaluation / Cost / Safety; Testing

  - [x] Define and validate a versioned structured reviewer result containing coverage identity, applicability (`direct`, `partial`, `not_triggered`, or `deferred`), findings, severities, evidence citations, and safe metadata.
  - [x] Canonicalize every citation to a workspace- or target-relative contained path and validate that each cited file and line/range exists before accepting the finding.
  - [x] Run deterministic product checks for all final decisions, plan-task execution evidence, the approved verification commands, complete reviewer coverage, deviations, missing evidence, and citation validity independently of reviewer prose.
  - [x] Treat missing, failed, malformed, or invalid coverage as terminal non-passing evidence; do not downgrade incomplete coverage to a warning.
  - [x] Normalize and sort coverage, checks, findings, deviations, and diagnostics by stable product keys after concurrent collection.
  - [x] Implement the product verdict truth table exactly: only `pass`, `pass_with_findings`, `fail`, or `blocked`; applicable blocker/high findings force `fail`, and only medium/low/info findings may yield `pass_with_findings`.
  - [x] Add pure table tests for every severity/verdict combination, applicability classes, missing evidence, deviations, verification failure, citation escape/invalid line, duplicate identity, and byte-stable normalized results from shuffled inputs.

- [x] **Task 5: Render, Validate, Persist, And Reconcile Current Evidence**

  > Executes: Decisions 4 and 5; AC-09, AC-11 through AC-15; Documentation; Persistence And Migrations

  - [x] Implement `internal/sprint/review_validation.go` checks for project/sprint/target identity, fingerprint/date, selected contracts/protocols, all required conformance sections, complete coverage, contained valid citations, no placeholders, and allowed deterministic verdicts.
  - [x] Render canonical `review.md` only from validated aggregate product data; reviewers must not write files or author canonical sections directly.
  - [x] Stage, flush, close, validate, and atomically replace sprint-root `review.md`; never delete or truncate the prior complete artifact before replacement succeeds.
  - [x] After artifact replacement, atomically update only review fields in `flow-state.json` with execution status, verdict, artifact path, timestamp, input fingerprint, staleness, and safe diagnostics.
  - [x] Reconcile state and artifact validity/fingerprints in status; classify interruption or mismatch as stale/non-current rather than inferring completion from file existence.
  - [x] Limit successful and failed review attempts to the canonical/staged review artifact and review flow-state fields; persist no partial reviewer cache or alternate state.
  - [x] Add failure injection for staging, write, flush, close, validation, rename, and flow-state update, proving the prior artifact stays byte-identical when replacement fails and unrelated files remain unchanged.
  - [x] Add rerun, legacy/manual artifact, artifact/state mismatch, stale-input, and successful replacement fixtures with normalized timestamps for deterministic output comparison.

- [x] **Task 6: Expose Shared App Use Cases And CLI Surfaces**

  > Executes: Decisions 1 and 6; AC-14 through AC-18; Architecture; CLI Surface; Observability

  - [x] Extend `internal/app` composition and typed sprint use cases for review readiness/status, validation, dry-run, prompt preview, execution, progress, cancellation, and final result using lazy runtime construction.
  - [x] Ensure shared DTOs carry scope, fingerprint, model source, reviewer count/concurrency, permitted writes, execution status, progress, verdict, staleness, findings, safe diagnostics, artifact path, and next action without terminal or TUI types.
  - [x] Wire review execution, review validation, review prompt preview, and review-target flow integration into `internal/app/sprint_commands.go` as thin parse/delegate/render adapters and extend status/help consistently.
  - [x] Map semantic preflight, blocked, runtime, timeout, cancellation, validation, persistence, and failing-verdict outcomes through the existing app exit/error boundary without string parsing.
  - [x] Keep final JSON as the only content on stdout in JSON mode; route progress and diagnostics through injected non-data channels and omit ANSI, raw payloads, and unsafe absolute paths.
  - [x] Make prompt preview and dry-run expose selected model source, reviewer count/bound, target/scope, contracts/protocols, expected writes, and mutation policy while proving they construct no runtime and write no review state.
  - [x] Add command/use-case tests for every required command, help, text, JSON structure, exit class, cancellation, no-write preview/dry-run, output separation, freshness/recovery, and shared DTO parity.

- [x] **Task 7: Add TUI Review Operation Parity**

  > Executes: Decisions 1 and 6; AC-17, AC-18; Architecture; CLI Surface; Observability

  - [x] Extend the existing TUI operational controls to display review readiness and scope, dry-run/prompt preview, model source, reviewer count/bound, selected contracts/protocols, permitted writes, and confirmation before runtime work.
  - [x] Start and cancel review through the same typed app use cases as CLI and translate typed reviewer lifecycle facts into live progress without interpreting native provider payloads.
  - [x] Present validated findings with navigation, execution status, deterministic verdict, staleness, safe diagnostics, next action, and read-only `review.md` preview.
  - [x] Preserve terminal interaction only in `internal/tui`; do not invoke CLI handlers, parse stdout, recompute readiness/verdict/freshness, or persist alternate review state.
  - [x] Add deterministic model/update/view tests for readiness, confirmation, progress, cancellation, failed/malformed coverage, findings navigation, verdict/result, stale recovery, preview, foreground/background behavior, and narrow-terminal rendering where existing patterns require it.
  - [x] Run CLI/JSON/TUI assertions against common semantic fixtures so each surface reports the same scope, state, verdict, freshness, diagnostics, and next action.

- [x] **Task 8: Complete Embedded Review Assets And Documentation**

  > Executes: Decision 8; AC-09, AC-16, AC-19; Documentation; Security

  - [x] Update the existing embedded `prompts/review.md` default to request the versioned structured contract/handbook result and enforce read-only, contained citation behavior.
  - [x] Update the existing embedded `templates/review.md` default to cover every validator-required review section and deterministic verdict field without placeholders in rendered output.
  - [x] Preserve workspace asset precedence: use a readable valid intentional override when present, embedded fallback only when absent, and fail for an existing unreadable or invalid override.
  - [x] Preserve explicit materialization through `ultraplan defaults install`; do not make workspace `prompts/` or `templates/` copies prerequisites or export them during workspace initialization.
  - [x] Update command help and relevant user-facing recovery guidance for prerequisites, dry-run/preview, cancellation, stale evidence, failed attempts, verdict meaning, and explicit non-goals.
  - [x] Add embedded fallback, override, invalid/unreadable override, materialization, prompt fixture, Markdown fixture/golden, and help tests.

- [x] **Task 9: Prove Product Boundaries And Complete Review Protocols**

  > Executes: Decision 8; AC-20, AC-21; RP-01, RP-02; Testing

  - [x] Complete the required fake-runtime matrix: pass, pass with findings, blocker/high fail, blocked preflight, missing result, malformed result, runtime failure, timeout, cancellation, retry, sibling collection, bounded concurrency, panic, stale input, citation escape, invalid line, redaction, and failed atomic replacement.
  - [x] Complete persistence snapshots proving no mutation to product source/tests, governed inputs, workspace config, unrelated sprint artifacts, or Git state and proving only complete current evidence is accepted.
  - [x] Complete app/command and TUI parity tests, including clean JSON stdout, lazy non-execution paths, progress/cancellation, status reconciliation, and artifact preview.
  - [x] Run targeted sprint, app, workspace/default, and TUI tests before the repository-wide gates; fix behavior rather than weakening required assertions.
  - [x] Run `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan` from `../ultraplan-go` and record command, result, and relevant diagnostics in `execute.md`.
  - [x] Apply RP-01 Architecture Review and record ownership/import/runtime-boundary evidence.
  - [x] Apply RP-02 Sprint Review and confirm every selected contract plus the handbook has valid coverage, the verdict is product-computed, and current `review.md` validates.
  - [x] Keep any real OpenCode run gated and optional; normal completion evidence must remain deterministic and credential-independent.

## Execution Order And Checkpoints

| Checkpoint | Tasks | Required Before Continuing |
| --- | --- | --- |
| Domain foundation | 1 | Review stage/status/verdict facts and strict state behavior are testable without runtime work. |
| Stable scope | 2 | Preflight, catalog resolution, immutable manifest, and fingerprint tests pass before reviewer scheduling is added. |
| Runtime collection | 3 | Exact coverage units, safety policy, bounds, cancellation, and awaited cleanup are proven with fakes. |
| Deterministic product result | 4 | Structured validation, deterministic checks, ordering, and verdict truth table pass independently of rendering. |
| Durable current evidence | 5 | Complete artifact validation, atomic replacement, preservation, and freshness reconciliation pass under failure injection. |
| Local surfaces | 6-8 | Shared DTO semantics, CLI/JSON/TUI parity, runtime-free previews, and embedded assets are complete. |
| Release gate | 9 | Full scenario matrix, required commands, and both selected review protocols have recorded evidence. |

## Evidence Checklist

- [x] Tests prove all AC-01 through AC-21 behaviors and all 13 selected contracts have traceable evidence.
- [x] Fake-runtime records prove exactly one request per selected contract plus one handbook request, bounded concurrency, collect-all behavior, read-only policy, and awaited cleanup.
- [x] Pure tests prove deterministic schema validation, citations, checks, ordering, and verdict synthesis.
- [x] Persistence evidence proves atomic replacement, prior-artifact preservation, strict state loading, fingerprint reconciliation, and permitted mutation scope.
- [x] CLI/app evidence covers review, validation, preview, flow, status, text/JSON output, semantic exits, progress, cancellation, and runtime-free dry-run behavior.
- [x] TUI evidence covers readiness, confirmation, progress, cancellation, finding navigation, verdict, staleness/recovery, and artifact preview through shared use cases.
- [x] Security evidence proves path/citation containment, required capability enforcement, redaction, raw-payload omission, and no forbidden mutation.
- [x] Embedded asset evidence proves fallback, intentional override, invalid override failure, and explicit materialization.
- [x] Documentation updates cover help, status fields, artifact sections, verdicts, failure/recovery, and non-goals.
- [x] Required Go test, race, and build commands pass and are recorded in `execute.md`.
- [x] RP-01 Architecture Review evidence is complete.
- [x] RP-02 Sprint Review evidence is complete and the current generated `review.md` validates.
- [x] Deviations from `reasoning.md` are recorded and approved before implementation continues.

## Verification Commands

Run all commands from `../ultraplan-go`.

| Check | Command | Expected Result |
| --- | --- | --- |
| Sprint review/domain tests | `go test ./internal/sprint/...` | Manifest, orchestration, deterministic checks/verdict, persistence, validation, flow, and freshness tests pass. |
| App and CLI review tests | `go test ./internal/app/...` | Required commands, typed use cases, JSON/text output, exit mapping, dry-run, and cancellation tests pass. |
| TUI review tests | `go test ./internal/tui/...` | Deterministic readiness, confirmation, progress, cancellation, findings, verdict, recovery, and preview tests pass. |
| Embedded asset tests | `go test ./internal/workspace/...` | Embedded review assets, override validation, and explicit materialization tests pass. |
| Full repository tests | `go test ./...` | All repository tests pass without OpenCode, credentials, or an interactive terminal. |
| Race gate | `go test -race ./...` | All tests pass with no reviewer, event-drain, progress, capture, cancellation, or persistence races. |
| CLI build gate | `go build ./cmd/ultraplan` | The UltraPlan CLI builds successfully. |

## Risks And Blockers

| Risk / Blocker | Source | Mitigation | Status |
| --- | --- | --- | --- |
| Existing project catalog API may not expose every canonical duplicate/path check needed by review. | `reasoning.md#assumptions-and-risks` | Extend only the existing catalog boundary and prove every preflight class; do not duplicate catalog parsing in review. | mitigated |
| Existing runtime boundary may not enforce all required read-only capabilities. | `reasoning.md#assumptions-and-risks` | Block before launch rather than degrade to prompt-only policy; adapt only through the generic agentwrap-backed boundary. | mitigated |
| Execute evidence may not form a trustworthy Git-independent target identity. | `reasoning.md#assumptions-and-risks` | Combine contained execute target identity, explicit changed paths, and fingerprints; block and record the open implementation question if trustworthy identity cannot be formed. | mitigated |
| Existing app/TUI progress seam may require a narrow extension. | `reasoning.md#assumptions-and-risks` | Reuse the narrowest typed callback/subscription seam; keep policy out of TUI and test parity/cancellation. | mitigated |
| Interruption between artifact and flow-state replacement can create cross-file mismatch. | `reasoning.md#decision-5-persist-only-complete-current-evidence-and-reconcile-freshness` | Fingerprint both, reconcile status, mark mismatch non-current, and require a fresh idempotent rerun. | mitigated |
| Governed or target inputs may drift while reviewers run. | `reasoning.md#decision-2-resolve-validate-freeze-and-fingerprint-one-review-manifest` | Freeze captured scope, recheck identity before persistence, fail stale, and preserve prior complete evidence. | mitigated |
| Cancellation or panic may leak reviewer/runtime work. | `reasoning.md#decision-3-run-independent-structured-reviewers-with-bounded-collect-all-fan-out` | Use one owner per started run, drain events, cancel, await bounded cleanup, convert panic to failed coverage, and run race/leak tests. | mitigated |
| Collect-all review can consume cost and time after a unit fails. | `reasoning.md#trade-off-and-debt-analysis` | Bound per-invocation concurrency and retries, expose progress/usage where safe, and support cancellation; add no cross-run throttle without measurements. | accepted |
| Malicious or malformed repository/model content may escape citations or disclose secrets. | `reasoning.md#assumptions-and-risks` | Validate structured schema, canonical paths and lines, redact diagnostics, omit raw payloads, and prohibit shell/source execution. | mitigated |
| Flow-state schema changes may break legacy readers or hide unsupported data. | `reasoning.md#assumptions-and-risks` | Version and strictly load state, define explicit supported/unsupported behavior, and retain legacy fixtures. | mitigated |
| CLI, JSON, and TUI renderers may drift semantically. | `reasoning.md#decision-6-expose-one-truthful-cli-json-status-flow-and-tui-surface` | Use one typed result model and common parity fixtures; prohibit surface-owned readiness, verdict, and freshness rules. | mitigated |
| An existing review asset override may be invalid or unreadable. | `reasoning.md#assumptions-and-risks` | Fail validation/preflight; use embedded fallback only when the override is absent. | mitigated |
| Verification command evidence may be absent or stale. | `reasoning.md#assumptions-and-risks` | Deterministically check exact approved commands and prevent pass when required evidence is missing/currentness cannot be established. | mitigated |
| `internal/sprint` may become harder to navigate. | `reasoning.md#potential-technical-debt` | Use focused review files and concrete pure functions; reassess extraction only after Sprint 27/28 demonstrates pressure. | accepted |

## Assumptions And Open Questions

| Item | Planning Treatment | Required Implementation Action |
| --- | --- | --- |
| Catalog resolution can be extended without a second parser. | Assumption from `reasoning.md:239`. | Verify the existing project API first; if insufficient, extend that API and retain project ownership. |
| Agentwrap-backed runtime can express structured validation, cancellation, required capabilities, and read-only policy. | Assumption from `reasoning.md:240`. | Verify capability support during Task 3; unsupported mandatory safety blocks execution. |
| Execute target plus changed paths can produce trustworthy implementation identity. | Assumption from `reasoning.md:241` and TRD open question. | Reuse concrete execute evidence; do not choose timestamps or a weak fallback. Record a decision deviation/open question if no trustworthy contained identity exists. |
| Existing app/TUI plumbing can carry typed reviewer lifecycle events. | Assumption from `reasoning.md:242`. | Reuse or narrowly extend the existing operation progress seam without terminal/provider types crossing into app or sprint. |
| Existing stores may already offer narrower atomic-write/failure seams than named files imply. | Open verification question from `reasoning/architecture.md:117`. | Inspect and reuse conforming seams; do not create a parallel filesystem abstraction solely to match anticipated file names. |
| Exact review concurrency default is not decided by reasoning. | Deliberately unresolved configuration detail. | Use established central configuration/default conventions and require a positive effective value; record any new default as a deviation before implementation. |

No unresolved question may be answered by weakening read-only enforcement, fingerprint completeness, reviewer coverage, deterministic verdicts, atomic preservation, or CLI/TUI parity.

## Guardrails And Rejected Alternatives

| Rejected Alternative | Required Guardrail |
| --- | --- |
| Global `internal/review`, generic workflow engine, plugin registry, or premature sprint subpackage tree | Keep focused review behavior in `internal/sprint` unless a concrete dependency/readability problem is recorded before deviation. |
| Review orchestration in `internal/app`, CLI, or TUI | Keep product policy in sprint; app composes typed operations and interfaces only render/interact. |
| Sprint semantics in `internal/platform/runtime` or direct OpenCode execution | Keep runtime generic and use agentwrap/OpenCode through the existing boundary. |
| Hardcoded selected contracts/protocols or per-worker live scope reads | Resolve through `project-index.md`, validate once, freeze one manifest, and recheck identity before persistence. |
| Fail-fast on an ordinary reviewer-unit failure | Continue healthy siblings and preserve all possible coverage; reserve invocation-wide cancellation for user/deadline/infrastructure invalidation. |
| Sequential or unbounded reviewer execution | Use configurable bounded per-invocation fan-out with deterministic sequential post-processing. |
| Model-authored canonical Markdown or model-selected verdict | Accept only validated structured evidence; product code renders and computes verdict. |
| Presence or timestamps as freshness proof | Reconcile valid artifact/state content fingerprints against current governed and target identity. |
| Cross-file transaction claim, database, or journal for two files | Use independent atomic replacement plus explicit fingerprint reconciliation. |
| Durable partial reviewer cache or mixed-fingerprint resume | Recover with a fresh manifest and complete rerun; use bounded agentwrap retry only within requests. |
| Prompt-only safety, banned-command lists, or best-effort required capabilities | Enforce permissions/capabilities and path policy mechanically or block before launch. |
| TUI-to-CLI calls, stdout parsing, provider-payload interpretation, or alternate TUI state | Use typed app operations and canonical sprint persistence only. |
| Progress on JSON stdout or unvalidated streamed findings as final data | Separate data/progress channels and expose final findings only after validation and deterministic sorting. |
| Required workspace prompt/template copies or implicit materialization | Use embedded defaults, optional intentional overrides, and explicit `defaults install`. |
| Smoke, `verify`, issue tracking, automatic fixes, Git mutation, remote review, or cross-run scheduling | Keep these out of Sprint 26 implementation and evidence except as explicit non-goals. |

## Review Inputs

Review should use:

- `sprint-index.md`
- `technical-handbook.md`
- `reasoning/architecture.md`
- `reasoning.md`
- this `plan.md`
- implementation diff
- verification evidence recorded in `execute.md` and `.run-state.json`
- `system/protocols/architecture-review-protocol.md`
- `system/protocols/review-sprint-protocol.md`

## Execution Log

| Date / Step | Action | Evidence / Notes |
| --- | --- | --- |
| 2026-07-19 / plan repair | Normalized the generated task checklist so the product plan validator recognized top-level tasks and nested subtasks. | `ultraplan ... validate plan` returned `Validation: ok`; no scope or decision changed. |
| 2026-07-19 / domain and manifest | Added typed review execution state/verdict facts, catalog-resolved frozen manifests, content/target fingerprints, changed-path identity, plan/execute preflight checks, and freshness reconciliation. | `internal/sprint/domain.go`, `review.go`, `state.go`, `service.go`; focused sprint tests pass. |
| 2026-07-19 / orchestration and persistence | Added exact contract-plus-handbook bounded collect-all fan-out, structured schema extraction, read-only runtime requests, citation/verdict validation, deterministic rendering, atomic replacement, and prior-artifact preservation. | `internal/sprint/review.go`, `review_test.go`; success, findings, blocker, malformed, escape, invalid-line, drift, cancellation, concurrency, and rename-failure fixtures pass. |
| 2026-07-19 / surfaces | Added review execution, validation, prompt preview, review-target flow integration, JSON/text/status output, configured review model/concurrency, shared app operations, and TUI review artifact/status/preview/dry-run/run/progress controls. | `internal/app`, `internal/tui`; targeted packages pass. |
| 2026-07-19 / embedded defaults | Replaced the prototype manual review prompt/template with embedded automated-review defaults while retaining optional explicit overrides and `defaults install`; workspace initialization still does not materialize overrides. | `internal/workspace`; embedded/default boundary test passes. |
| 2026-07-19 / verification | Ran repository gates from the implementation target. | `go test ./...` pass; `go test -race ./...` pass; `go build ./cmd/ultraplan` pass. The build emitted a non-fatal Go module stat-cache warning because the host module cache is read-only. |
| 2026-07-19 / race-gate repair | A final repeated race run exposed an existing unsynchronized dependency-state read in the study run loop. Captured dependency state under the existing mutex before evaluation and reran the study and repository race gates. | `internal/study/run_loop.go`; narrow gate-required repair, with no review semantics or API changes. |
| 2026-07-19 / implementation note | Kept review validation, rendering, and atomic write helpers in focused `internal/sprint/review.go` rather than creating anticipated `review_validation.go` and store wrappers. | Same package ownership and public behavior; avoided file-shape-only abstraction. No requirement or architecture deviation. |

## Completion Criteria

- [x] All nine tasks are complete or explicitly deferred with requirement impact and approval recorded.
- [x] Every AC-01 through AC-21 row has implementation and evidence pointers.
- [x] Every selected contract and both required protocols have current evidence.
- [x] Verification commands were run successfully or any deferral is documented as a blocker; required gates cannot be waived for completion.
- [x] Evidence satisfies the expected unit, fake-runtime, persistence, CLI/app, TUI, runtime, security, build/race, review, and documentation categories from `reasoning.md#expected-evidence`.
- [x] No deviation from `reasoning.md` was implemented before being recorded and approved.
- [x] The last complete `review.md` is preserved across failed attempts and the successful current artifact validates against its matching fingerprint/state.
- [x] `review.md` can evaluate conformance without guessing implementation intent.
