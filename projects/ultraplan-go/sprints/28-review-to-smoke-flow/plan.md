# Sprint Plan: Integrated Review-to-Smoke Verification Flow

> Project: `ultraplan-go`
> Sprint: `28-review-to-smoke-flow`
> Source: `reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/28-review-to-smoke-flow/requirements.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/sprints/28-review-to-smoke-flow/sprint-index.md`, `projects/ultraplan-go/sprints/28-review-to-smoke-flow/technical-handbook.md`, `projects/ultraplan-go/sprints/28-review-to-smoke-flow/reasoning/architecture.md`, `projects/ultraplan-go/sprints/28-review-to-smoke-flow/reasoning.md`, injected `builtin:templates/sprint-plan.md`

This plan executes `reasoning.md`. It does not reopen architecture, scope, verdict, assessment, fingerprint, override, rerun-promotion, or persistence decisions.

## Reasoning Source

- **Sprint Reasoning:** `reasoning.md`
- **Sprint Index:** `sprint-index.md`
- **Technical Handbook:** `technical-handbook.md`
- **Area Reasoning:** `reasoning/architecture.md`

The selected final study reports were not reopened while writing this plan. `technical-handbook.md`, `reasoning/architecture.md`, and `reasoning.md` already distill and apply them, and no implementation question required deeper report evidence. This avoids introducing evidence outside the sprint selection.

## Sprint Status

- **Status:** `implementation complete; protocol review/smoke pending`
- **Owner:** `Codex`
- **Start Date:** `2026-07-22`
- **Completion Date:** `pending protocol evidence`

## Decisions To Execute

| Decision | Source Section | Requirement Satisfied | Execution Implication | Evidence Used | Accepted Trade-Off | Rejected Alternative | Risk / Follow-Up |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Keep Verification Policy Product-Owned | `reasoning.md#decision-1-keep-verification-policy-product-owned` | Acceptance Criteria 1, 10, 11; Architecture and Workflow contracts | Move integrated flow dispatch behind `internal/sprint`; expose typed app use cases; keep CLI/TUI as input and rendering adapters; keep platform runtime/process product-neutral. | Architecture reasoning, project Architecture/TRD 18.1, handbook thin-surface and boundary patterns | Typed requests/results and converged call sites add code. | Surface-owned policy, new verification package, workflow engine, event bus, plugin registry, and wholesale stage replacement. | Verify existing Sprint 26/27 seams before refactoring; prevent a shared result from becoming a raw-payload DTO. |
| Version Facts And Derive Current Truth | `reasoning.md#decision-2-version-facts-and-derive-current-truth` | Acceptance Criteria 2, 3, 6, 7, 9; Workflow, Persistence, Observability contracts | Bump strict `flow-state.json`; persist attempt and last-complete facts; derive freshness, readiness, assessment, and next action on reads; migrate exactly one known predecessor. | Architecture reasoning, TRD 18.5, handbook lifecycle/state evidence | Reads reconcile artifacts and hash bounded inputs; migration blocks unknown history. | Third assessment artifact, database, alternate TUI state, Markdown-only reconstruction, or authoritative cached staleness. | Identify the actual predecessor schema before coding; reject rather than infer unknown or ambiguous state. |
| Fingerprint Deterministic Stage Manifests | `reasoning.md#decision-3-fingerprint-deterministic-stage-manifests` | Acceptance Criteria 2, 3, 6; Persistence, Security, Testing contracts | Build normalized review/smoke manifests with lexical workspace-relative paths, missing markers, content identities, changed scope, target identity, and separate artifact digests; propagate stale review to smoke. | Requirements AC2, TRD 18.9.1/18.9.3, review protocol, architecture reasoning | Conservative non-Git hashing may cost time and react to generated files. | Mtime freshness, path/existence checks, or copying raw harness evidence. | Fixture path containment, symlink, large/bounded read, Git, and non-Git behavior; measure before adding any cache. |
| Use One Ordered Gate And Deterministic Assessment | `reasoning.md#decision-4-use-one-ordered-gate-and-deterministic-assessment` | Acceptance Criteria 1, 6, 9, 12; Workflow, CLI Safety, Configuration contracts | Make verification and flow requests targeting the final verification stage call one `execute -> review -> smoke` transition; require acceptable current review by default; make override explicit, confirmed, reasoned, recorded, and non-promoting. | Review and Deep Smoke protocols, TRD 18.9.3, handbook sequential-flow evidence | No speculative smoke; diagnostic runs may complete without becoming current/promotable. | Silent continuation after failed review, separate flow/verify engines, or override-created clean assessment. | Keep override metadata small and actor-neutral; test that it cannot improve verdict, freshness, assessment, or capabilities. |
| Promote Focused Reruns Only With Complete Coverage | `reasoning.md#decision-5-promote-focused-reruns-only-with-complete-coverage` | Acceptance Criteria 4, 5, 6, 8; Workflow, Persistence, Smoke Testing contracts | Promote focused review only after complete same-fingerprint coverage; treat narrow smoke passes as diagnostic until the required containing suite is current and passing; include relevant issue state in verdict/freshness. | Deep Smoke protocol, handbook narrow-rerun trade-off, architecture reasoning | Narrow success can require broader work before canonical verdict improves. | Replacing canonical evidence with a passing narrow run, retaining review results from another fingerprint, or product-owned issue mutation. | Validate structured issue relevance and mandatory reviewer coverage; ambiguity prevents a clean pass. |
| Preserve Evidence Through Atomic, Cancellable Attempts | `reasoning.md#decision-6-preserve-evidence-through-atomic-cancellable-attempts` | Acceptance Criteria 7, 8, 9; Persistence, Workflow, CLI Lifecycle, Error contracts | Serialize sprint mutations; record attempts without deleting last-complete evidence; propagate context; join work; atomically promote validated Markdown; reconcile cross-file interruption by identity. | Architecture reasoning, handbook cancellation/cleanup evidence, TRD 18.9.1/18.9.2 | Explicit attempt transitions and retryable operator conflicts replace simpler overwrite behavior. | Overwrite-at-start, detached work, or atomic rename without logical mutation exclusion. | Fault-inject artifact/state boundaries and process cleanup; use bounded cleanup context and report cleanup failure separately. |
| Freeze Safe Execution And Expose One Semantic Result | `reasoning.md#decision-7-freeze-safe-execution-and-expose-one-semantic-result` | Acceptance Criteria 3, 8-13; Configuration, Observability, Security, CLI, LLM, Error, Documentation, Testing contracts | Resolve effective config once; reuse generic runtime; invoke harness by executable/argv with contained cwd and allowlisted environment; expose one semantic result to text/JSON/TUI; keep JSON stdout final-result-only. | Handbook configuration, process, IO, error, observability, and testing evidence; Review/Smoke protocols | Preflight and typed mappings grow; raw payloads remain at their owning source. | Shell strings, direct product-level environment reads, platform verdict policy, unbounded output, JSON progress, or TUI CLI parsing. | Stable JSON becomes a compatibility contract; validate harness manifest capabilities and redact diagnostics. |

## Requirements / Contracts To Satisfy

| Contract / Requirement ID | Required Behavior | Evidence Planned |
| --- | --- | --- |
| AC-1 | Flow and verification requests targeting the final verification stage require current review unless explicit diagnostic override is selected. | Shared transition-table and command tests for pass, findings, fail, blocked, stale, and override paths. |
| AC-2 | Changes to governed inputs, selected catalog content, execute evidence, target identity, `review.md`, or `smoke.md` make owning evidence stale; stale review makes smoke stale. | Golden manifest fixtures and one mutation test per input class, including external artifact edits and propagation. |
| AC-3 | Text, JSON, and TUI expose execution state, verdict, freshness, artifacts, smoke run/block reason, relevant issues, and next action. | Shared semantic fixtures plus adapter assertions and selected stable-output goldens. |
| AC-4 | Focused review and smoke reruns are available, but narrow smoke success cannot replace required containing-suite evidence. | Focused merge and scope-promotion tests with containing-suite failures and stale coverage. |
| AC-5 | Relevant open harness issues prevent a clean smoke pass and appear in `smoke.md` and status. | Open/resolved/ambiguous issue fixtures and status/artifact link assertions. |
| AC-6 | Overall assessment follows the fixed precedence from current canonical evidence and creates no third artifact. | Table-driven precedence, contradiction, and no-extra-artifact tests. |
| AC-7 | Interrupted attempts preserve valid state and the last complete `review.md`/`smoke.md`. | Cancellation, timeout, malformed output, fault-injected write/rename/state update, and byte-preservation tests. |
| AC-8 | Recovery covers stale input, malformed summary, missing evidence, external edits, blocked environment, cancellation, and timeout. | Typed result/exit/next-action matrix across sprint, app, CLI, and TUI tests. |
| AC-9 | CLI text, stable JSON, and TUI agree on current semantic state and recovery. | Cross-surface assertions over the same typed app results and temporary workspace fixtures. |
| AC-10 | TUI verification actions use shared typed app use cases and no CLI/subprocess/output-parsing path. | TUI fake-use-case tests and architecture/import review. |
| AC-11 | Verification does not mutate product source/tests, governed inputs, project docs, or Git state. | Before/after mutation snapshots, path containment fixtures, and Security/Architecture review evidence. |
| AC-12 | Normal tests use fake review runtime and fake harness/process seams without credentials, network, OpenCode, or real harness. | Deterministic package tests and review of test setup/environment dependencies. |
| AC-13 | Full tests, race tests, and binary build pass. | `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan`. |
| Architecture: `ARCH-CORE-001`, `ARCH-CORE-002`, `ARCH-LAYER-001`, `ARCH-LAYER-002`, `ARCH-ENTRY-001`, `ARCH-MODULE-001`, `ARCH-COMP-001`, `ARCH-SHARED-001` | Product semantics remain in `internal/sprint`; app composes typed use cases; surfaces render; platform packages remain generic. | Dependency/import review, fake substitution, and Architecture Review Protocol evidence. |
| Errors: `ERR-CORE-001`, `ERR-SHAPE-001`, `ERR-CODE-001`, `ERR-TRANS-001`, `ERR-RETRY-001`, `ERR-IO-001`, `ERR-TASK-001`, `ERR-DEP-001`, `ERR-DATA-001`, `ERR-REDACT-001`, `ERR-USER-001` | Outcomes retain typed category/cause and actionable recovery without collapsing blocked, stale, malformed, cancelled, timed-out, conflict, and failed states. | Error/exit matrix, cause-chain, retry classification, redaction, and recovery tests. |
| Configuration: `CFG-SOURCE-001`, `CFG-TYPE-001`, `CFG-COMPAT-001`, `CFG-SECRET-001`, `CFG-START-001`, `CFG-ENV-001`, `CFG-PUBLIC-001`, `CFG-OBS-001` | Effective review/harness/timeout/environment/rerun/override configuration has explicit precedence, typed post-merge validation, safe source display, and compatibility. | Precedence/source, invalid config, environment allowlist, dry-run, and redaction tests. |
| Observability: `OBS-CORE-001`, `OBS-CORR-001`, `OBS-HEALTH-001`, `OBS-DIAG-001`, `OBS-DEBUG-001`, `OBS-TASK-001`, `OBS-ALERT-001`, `OBS-PII-001` | Correlated attempts expose truthful readiness, terminal status, evidence identity, staleness, safe diagnostics, and next action without leaking secrets. | Event/log assertions, correlation fields, terminal-state snapshots, redaction, and JSON-cleanliness tests. |
| Security: `SEC-INPUT-001`, `SEC-INJECT-001`, `SEC-SECRETS-001`, `SEC-DEPS-001`, `SEC-FILES-001`, `SEC-NET-001`, `SEC-DESER-001`, `SEC-DEFAULT-001`; `SEC-AUTHN-001` and `SEC-AUTHZ-001` not triggered | Validate untrusted paths/state/evidence, fail closed, use explicit argv and contained writes, and prohibit product/Git mutation; no identity or protected-object authorization boundary is introduced. | Argv/cwd/env capture, symlink/path escape, malformed state/evidence, redaction, mutation snapshots, and architecture review. |
| Testing: `TEST-SEAM-001`, `TEST-UNIT-001`, `TEST-INT-001`, `TEST-SMOKE-001`, `TEST-SMOKE-002`, `TEST-FAIL-001`, `TEST-DET-001`, `TEST-CONTRACT-001`, `TEST-E2E-001`, `TEST-MIGRATION-001` | Policy, integration, failure, migration, contract, parity, and journey behavior is deterministic through public fake seams; real evidence remains gated. | Sprint/app/TUI tests, fixtures/goldens, fake runtime/harness, full/race/build checks, and gated protocol evidence. |
| Documentation: `DOC-OWNER-001`, `DOC-ARCH-001`, `DOC-PUBLIC-001`, `DOC-OPS-001`, `DOC-EXAMPLE-001`, `DOC-GEN-001`, `DOC-AGENT-001` | Command help, artifact validators, and operator guidance document gates, reruns, override, state, links, exits, and recovery without a third assessment artifact. | Help tests/review, validator fixtures, generated-summary checks, and Documentation contract review. |
| CLI Surface: `CLI-SHAPE-001`, `CLI-HELP-001`, `CLI-IO-001`, `CLI-EXIT-001`, `CLI-JSON-001`, `CLI-CONFIG-001`, `CLI-LIFE-001`, `CLI-SAFE-001`, `CLI-NONINT-001` | Commands are predictable, cancellable, noninteractive-safe, and stable in help/exits/JSON with strict stdout/stderr separation. | Command/help/exit tests, JSON fixtures/goldens, non-TTY confirmation tests, and cancellation/cleanup tests. |
| LLM Runtime: `LLM-BOUNDARY-001`, `LLM-IO-001`, `LLM-LIFECYCLE-001`, `LLM-RETRY-001`, `LLM-RUN-001`, `LLM-PROMPT-001`, `LLM-EXPOSE-001`, `LLM-OBS-001`, `LLM-EVAL-001`, `LLM-SAFETY-001`, `LLM-COST-001`; `LLM-TOOL-001` only if reviewer tools are enabled | Existing agentwrap boundary stays generic, read-only, cancellable, inspectable, and insufficient by itself to establish review success. | Fake runtime structured-result, prompt identity, metadata, malformed/missing result, cost/safety, optional tool-policy, and cancellation tests. |
| Workflows: `WF-SCOPE-001`, `WF-BOUNDARY-001`, `WF-STATE-001`, `WF-RETRY-001`, `WF-IDEMPOTENCY-001`, `WF-COMP-001`, `WF-VERSION-001` | One sprint-owned ordered flow has explicit state, retry, rerun, cancellation, compensation/reconciliation, idempotency, and version behavior. | Transition tables, duplicate/focused reruns, interruption, conflict, stale propagation, and command journey tests. |
| Persistence: `PERSIST-SCHEMA-001`, `PERSIST-MIG-001`, `PERSIST-ATOMIC-001`, `PERSIST-INTEGRITY-001`, `PERSIST-READ-001`, `PERSIST-DERIVED-001`, `PERSIST-FIXTURE-001`, `PERSIST-RECOVERY-001` | Strict versioned facts, one known migration, bounded reads, derived truth, atomic replacement, integrity, fixtures, and recovery preserve last-complete evidence. | Migration fixtures, malformed/unknown schema, bounded read, atomic failure, reconciliation, and artifact-preservation tests. |

## Fixed Product Policies

### Overall Assessment Precedence

Implementation and tests must evaluate these conditions in order. Execution status remains separate from verdict, and `stale` remains freshness rather than a verdict.

| Ordered Condition | Assessment | Next Action |
| --- | --- | --- |
| Contradictory/malformed state or unvalidated referenced canonical evidence | `blocked` | Repair/recover the named state or rerun its owning stage; do not infer pass. |
| Current review is `fail` | `fail` | Resolve review findings/evidence and rerun review. |
| Current review is `blocked` | `blocked` | Satisfy the named review prerequisite and rerun review. |
| Required review/smoke is missing, running, cancelled, timed out, stale, or represented only by an unpromoted failed attempt | `incomplete` | Run or retry the earliest required non-current stage. |
| Current acceptable review and current smoke is `fail` | `fail` | Resolve the classified cause, diagnose narrowly if useful, then run the required containing suite. |
| Current acceptable review and current smoke is `blocked` | `blocked` | Restore prerequisites, coverage, or evidence and rerun smoke. |
| Current acceptable review and current smoke is `not_applicable` | `not_applicable` | Retain the current review verdict and smoke rationale; no runtime smoke action is required. |
| Current acceptable review and current smoke is `pass_with_open_issues` | `pass_with_findings` | Resolve relevant issues, rerun narrowly, then rerun the required containing scope. |
| Current review is `pass_with_findings` and current smoke is `pass` | `pass_with_findings` | Address review findings if a clean assessment is required. |
| Current review is `pass` and current smoke is `pass` | `pass` | No verification action is required. |

### Freshness And Promotion Invariants

- Review manifest identity includes governed planning inputs, selected catalog/contracts/protocols, execute summary/run-state evidence, target implementation identity, and explicit changed-path scope.
- Smoke manifest identity includes current review fingerprint/verdict/artifact digest, versioned harness manifest identity, selected scope and containing-suite requirement, target identity, prerequisites/config identities, and validated issue/evidence identities.
- Paths are normalized workspace-relative paths in lexical order; missing inputs are explicit; timestamps never establish freshness.
- Trustworthy Git identity uses contained root, commit, and changed-file content digests. The fallback is a contained regular-file manifest excluding only VCS metadata.
- `review.md` and `smoke.md` each have a separate artifact digest and are never included self-referentially in the input fingerprint that produced them.
- Any owning-input or canonical-artifact change makes that stage stale. Stale review always makes smoke stale.
- Focused review promotion requires one complete validated result under one current review fingerprint with mandatory selected-contract and handbook coverage.
- Focused smoke success remains diagnostic until the required containing suite is current and passing.
- Relevant open issues prevent clean smoke pass; missing or ambiguous structured issue relevance also prevents a clean pass.
- Diagnostic override permits smoke execution only. It cannot improve review, freshness, smoke promotion, assessment, mutation permissions, or environment permissions.

## Tasks

- [x] **Task 1: Confirm Existing Stage And Compatibility Seams**
  > Executes: `Decisions 1, 2, 7; AC-7, AC-12`
  - [x] Inspect current `internal/sprint` review, smoke, flow, state, validation, runtime, process, persistence, clock, progress, and mutation seams; record the minimal convergence points in implementation evidence rather than replacing Sprint 26/27 behavior.
  - [x] Identify the exact immediately preceding `flow-state.json` schema represented by current code/fixtures and define the one deterministic predecessor-to-new-schema mapping. If it is not identifiable, stop schema work and record a blocker instead of guessing.
  - [x] Validate that the cataloged harness's versioned machine manifest supplies executable/argv, protocol and evidence versions, discovery/run operations, scope/containing-suite data, prerequisites, run/evidence IDs, and structured issue relevance. If a required capability is absent, preserve a typed blocked result and recovery action; do not parse README prose.
  - [x] Confirm existing fake runtime and fake process/harness seams can drive all normal tests; add only narrow consumer-owned seams where a concrete missing effect prevents deterministic testing.

- [x] **Task 2: Define The Integrated Semantic Domain Projection**
  > Executes: `Decisions 1, 2, 4, 7; AC-3, AC-6, AC-8-10`
  - [x] In `internal/sprint/domain.go`, represent review/smoke execution status separately from completed verdict and freshness, and define typed stage, attempt, overall-assessment, relevant-issue, override, artifact/evidence reference, safe diagnostic, and next-action facts required by all surfaces.
  - [x] Keep raw reviewer/provider/harness/process payloads and rendering-only fields out of the product projection; use request-specific result shapes only where current consumers demonstrably require different concepts.
  - [x] Encode the ordered assessment table above as the sole assessment policy, including contradiction/malformed precedence and `incomplete` handling for stale/interrupted/unpromoted attempts.
  - [x] Add table-driven tests in `internal/sprint/verify_test.go` proving every assessment row, orthogonal execution/verdict/freshness concepts, issue-aware assessment, and diagnostic override non-promotion.

- [x] **Task 3: Version State, Attempts, Migration, And Reconciliation**
  > Executes: `Decisions 2, 6; AC-2, AC-6-8`
  - [x] Update `internal/sprint/state.go`, `internal/sprint/domain.go`, and `internal/sprint/validation.go` so the strict new `flow-state.json` stores active and last attempts, last-complete artifact path/digest, normalized input fingerprint, timestamps, safe diagnostics, and review/smoke-specific run, evidence, issue, and override facts.
  - [x] Implement one atomic migration from the confirmed immediate predecessor: preserve known stage and last-complete facts, mark unverifiable freshness stale, and reject older, unknown, contradictory, or malformed versions with typed recovery guidance.
  - [x] Ensure starting an attempt atomically records running state under per-sprint mutation exclusion without deleting last-complete facts or replacing canonical Markdown.
  - [x] Promote only after validated same-directory temporary write, flush/close, atomic rename, and matching state update; reconcile artifact/state interruption by digest/fingerprint on read without adopting path existence as current evidence.
  - [x] Add current/predecessor/unknown/malformed fixtures and fault-injection tests for load, migration, write, rename, state-update interruption, stale running-attempt recovery, conflict, and prior-artifact byte preservation.

- [x] **Task 4: Implement Deterministic Review And Smoke Manifests**
  > Executes: `Decision 3; AC-2, AC-3, AC-6`
  - [x] Update sprint-owned review/state/validation code to build the exact review manifest defined above with normalized contained paths, lexical ordering, missing markers, selected content identities, execute evidence, target identity, and changed-path scope.
  - [x] Update smoke/state/validation code to build the exact smoke manifest defined above, including current review identity, containing-suite obligation, harness manifest, target/prerequisite/config identities, and validated external issue/evidence identities.
  - [x] Compare separate `review.md`/`smoke.md` content digests during status reconciliation; external edits must become stale or malformed, never silently current.
  - [x] Preserve current evidence as historical recovery context when stale and propagate stale review to smoke unconditionally.
  - [x] Add golden manifest fixtures and mutation tests for every required input class, path order, missing inputs, path/symlink escape, Git identity, non-Git fallback, changed scope, artifact edits, external issue changes, bounded reads, and stale propagation.

- [x] **Task 5: Consolidate Ordered Flow, Gate, Assessment, And Recovery**
  > Executes: `Decisions 1, 4, 6; AC-1, AC-6-9`
  - [x] Move integrated stage sequencing from app-owned dispatch into a single `internal/sprint` transition path in `flow.go`/`service.go` that both `flow` and `verify` call through `execute -> review -> smoke`.
  - [x] Require complete execute evidence, then current review, then evaluate the smoke gate; default smoke accepts only current `pass` or `pass_with_findings` review.
  - [x] Require explicit confirmation and non-empty rationale for diagnostic override; persist safe actor-neutral request metadata and the blocked/stale/failing review identity.
  - [x] Return typed blocked, failed, stale, malformed, cancelled, timed-out, conflict, and incomplete outcomes with the earliest required next action; do not collapse them into parsed strings.
  - [x] Add shared transition-table tests for `flow` and `verify`, pass/findings/fail/blocked/stale review gates, override audit/non-promotion, missing execute evidence, and each recovery category.

- [x] **Task 6: Complete Review Freshness And Focused Reruns**
  > Executes: `Decisions 3, 5, 6, 7; AC-2, AC-4, AC-7, AC-12`
  - [x] Extend existing `internal/sprint/review.go` behavior rather than replacing it; reuse the generic agentwrap seam, read-only permissions, stable prompt identity, bounded joined reviewer concurrency, caller context, and validated structured results.
  - [x] Add focused reviewer/check selection and merge rules that can promote only one complete validated `review.md` whose retained and rerun results all match the current review fingerprint and cover every mandatory selected contract plus handbook review.
  - [x] Keep product code as verdict owner; runtime success, missing reviewer output, malformed output, incomplete coverage, or stale retained results cannot produce a pass.
  - [x] Preserve the last complete review artifact and facts through cancellation, timeout, malformed output, validation failure, runtime failure, cleanup failure, or input drift.
  - [x] Expand `internal/sprint/verify_test.go` and existing review tests for focused merge, complete coverage, stale retained-result rejection, bounded collection, cancellation, timeout, malformed/missing results, drift, and artifact preservation.

- [x] **Task 7: Complete Smoke Scope, Evidence, Issues, And Recovery**
  > Executes: `Decisions 3-7; AC-1, AC-2, AC-4-8, AC-11-12`
  - [x] Extend existing `internal/sprint/smoke.go` behavior rather than replacing it; resolve only the cataloged versioned machine manifest and freeze effective executable, argv, contained cwd, allowlisted environment, timeout, scope, prerequisites, and evidence destination before execution.
  - [x] Enforce scope precedence and containing-suite promotion: sprint suite, directly relevant suite, explicit level/suite/test diagnosis, then full harness for broad closure when required.
  - [x] Validate run ID, result counts, evidence identities/paths, containing coverage, and structured relevant issue IDs/links/statuses under cataloged harness roots; retain raw evidence externally.
  - [x] Synthesize only protocol-defined smoke verdicts; open relevant issues prevent `pass`, and missing runtime, credentials, coverage, prerequisites, evidence, ambiguous relevance, malformed output, cleanup failure, or unsafe paths fail closed as typed `blocked` or `fail`.
  - [x] Preserve the last complete smoke artifact through all failed attempts; record blocked/not-applicable reason and next action without fabricating a successful run ID.
  - [x] Expand `internal/sprint/verify_test.go` and existing smoke tests for scope precedence, narrow-pass/containing-fail, open/resolved/ambiguous issues, evidence containment, malformed/missing evidence, blocked/not-applicable, timeout/cancellation/cleanup, explicit argv capture, and artifact preservation.

- [x] **Task 8: Expose Shared App Use Cases And CLI Parity**
  > Executes: `Decisions 1, 4, 7; AC-1, AC-3, AC-8-10, AC-13`
  - [x] Update `internal/app/sprint_usecases.go` with typed status, flow, review, smoke, verify, focused-rerun, cancellation, and recovery operations that return the sprint-owned semantic projection without recomputing policy.
  - [x] Update `internal/app/sprint_commands.go` so verification and non-dry flow requests targeting the final verification stage use the same app/sprint transition, and status/review/smoke/rerun operations share typed results.
  - [x] Parse and validate review focus, smoke level/suite/test focus, timeouts, and diagnostic override confirmation/rationale with frozen precedence `explicit flags -> environment -> workspace config -> built-in defaults`; display non-secret sources in dry-run/confirmation.
  - [x] Render aligned text and stable final-result-only JSON with state, verdict, freshness/reasons, artifacts, run/evidence IDs, relevant issues, override facts, assessment, and next action; route progress/logs/diagnostics away from JSON stdout.
  - [x] Define stable exits from typed categories and update help/operator guidance for gates, reruns, override limits, statuses, evidence links, cancellation, timeout, malformed evidence, and recovery.
  - [x] Add `internal/app/sprint_verify_commands_test.go` for commands/help, flow/verify equivalence, config precedence, non-TTY confirmation safety, exit matrix, stdout/stderr separation, stable JSON, diagnostic override, status next action, and text/JSON semantic parity.

- [x] **Task 9: Deliver Complete TUI Verification Operation**
  > Executes: `Decisions 1, 4, 7; AC-3, AC-4, AC-8-11`
  - [x] Update `internal/tui/model.go` actions to call shared typed app use cases for status, flow/verify, review, smoke, focused reruns, cancellation, and recovery; add no CLI handler, subprocess, terminal-output parser, provider/harness parser, or alternate persistence path.
  - [x] Add guarded confirmation for runtime work and diagnostic override rationale; show the same frozen preflight scope/config sources, review gate, and containing-suite obligations as app results.
  - [x] Update `internal/tui/views.go` to render execution status, verdict, staleness/reasons, current versus last-complete evidence, artifacts/links, run ID or blocked/not-applicable reason, relevant issues, override status, overall assessment, and next action.
  - [x] Keep progress bounded and ephemeral, make cancellation/recovery truthful, and preserve usable rendering at narrow terminal sizes.
  - [x] Add `internal/tui/verify_test.go` with deterministic fake use cases for end-to-end actions, confirmation, override, progress, cancellation, recovery, focused reruns, evidence links, narrow layouts, and semantic agreement with app results.

- [x] **Task 10: Prove Safety, Mutation Boundaries, And Release Gates**
  > Executes: `Decisions 1-7; AC-7, AC-8, AC-11-13; all selected contracts and protocols`
  - [x] Verify caller context reaches runtime/process/progress/waits, owned work is joined, cleanup is separately bounded where needed, and mutation exclusion produces actionable conflicts under race tests.
  - [x] Snapshot product source, product tests, governed sprint inputs, project docs, and Git metadata before/after status, freshness, review, smoke, verify, focused rerun, cancellation, and recovery scenarios; assert only owned state/summaries and harness-declared evidence roots change.
  - [x] Verify platform runtime/process packages contain no project/sprint/review/smoke/harness/verdict semantics and use no shell interpretation; verify TUI imports/calls no CLI execution path.
  - [x] Verify generated paths and diagnostics are workspace-relative unless an absolute path is necessary for explicit local recovery, and verify secrets/raw harness output do not enter stable state, summaries, logs, or JSON.
  - [x] Run targeted package tests, full tests, race tests, and binary build; record command, result, and any approved deferral in `execute.md` rather than this plan.
  - [x] Prepare review evidence for Architecture Review, Sprint Review, and Deep Smoke Sprint protocols. Real runtime/harness evidence remains gated and must be a matching linked run or truthful protocol-valid `blocked`/`not_applicable` result.

## Evidence Checklist

- [x] Transition tests prove one ordered `execute -> review -> smoke` path for both `flow` and `verify`.
- [x] Assessment tests cover every precedence row and prohibit contradictory or override-laundered results.
- [x] Manifest tests cover every governed input class, artifact edits, issue changes, path normalization/containment, Git/non-Git target identity, and stale propagation.
- [x] Migration and persistence tests cover the one known predecessor, strict rejection, atomic failure, reconciliation, mutation conflict, and last-complete byte preservation.
- [x] Focused-rerun tests prove complete review coverage and required containing-suite smoke promotion.
- [x] Runtime/diagnostic evidence includes safe correlation, attempt, cleanup, evidence-link, staleness, override, and next-action facts.
- [x] CLI text, stable JSON, and TUI semantic parity is proven from shared typed app results.
- [x] JSON stdout contains only the final stable result; progress and diagnostics cannot contaminate it.
- [x] Mutation-boundary, containment, argv, environment allowlist, bounded output, redaction, and no-shell tests pass.
- [x] Normal tests prove required behavior with fakes and no OpenCode, credentials, network, or real harness.
- [x] Command help and operator recovery guidance cover gates, reruns, override limits, state, exits, links, cancellation, timeout, and malformed/missing evidence.
- [x] `review.md`/`smoke.md` validators enforce required sections, verdict vocabulary, evidence links, and no raw harness evidence.
- [x] Architecture Review, Sprint Review, and Deep Smoke Sprint protocol evidence is available.
- [x] Deviations from `reasoning.md` are recorded and approved before implementation continues.

## Verification Commands

Run commands from `../ultraplan-go`.

| Check | Command | Expected Result |
| --- | --- | --- |
| Sprint policy/state tests | `go test ./internal/sprint` | Integrated transition, state, freshness, assessment, rerun, recovery, and safety tests pass with fakes. |
| App/CLI tests | `go test ./internal/app` | Verify/flow/status commands, exits, help, config, and text/JSON parity pass. |
| TUI tests | `go test ./internal/tui` | Typed actions, rendering, cancellation/recovery, and app-semantic parity pass without a real terminal. |
| Full test suite | `go test ./...` | All packages pass without OpenCode, credentials, network, or the real smoke harness. |
| Race suite | `go test -race ./...` | No races in mutation exclusion, status reads, progress, cancellation, or cleanup. |
| Binary build | `go build ./cmd/ultraplan` | The `ultraplan` binary builds successfully. |

## Risks And Blockers

| Risk / Blocker | Source | Mitigation | Status |
| --- | --- | --- | --- |
| Existing Sprint 26/27 seams may require more convergence than expected. | `reasoning.md#assumptions-and-risks` | Existing services were composed behind `Service.Verify`; no wholesale replacement was required. | closed |
| The immediate predecessor flow-state schema may not be deterministically identifiable. | `reasoning.md#assumptions-and-risks` | Version 1 was confirmed in code and fixtures; only v1 migrates to strict v2 and unknown versions fail. | closed |
| Harness manifest may lack required structured scope/evidence/issue/environment metadata. | `reasoning.md#assumptions-and-risks` | Protocol-v1 manifest and discovery validation remain fail-closed; actual Sprint 28 mapping is checked by gated smoke. | mitigated; gated evidence pending |
| Non-Git fallback hashing may be expensive or noisy. | `reasoning.md#potential-technical-debt` | Deterministic contained reads are bounded at 64 MiB per file; optimization remains measurement-driven. | mitigated; monitor latency |
| Artifact and state updates cannot share one filesystem transaction. | `reasoning.md#assumptions-and-risks` | Ordered atomic replacement, separate digests, and read-time mismatch detection preserve recovery truth. | mitigated |
| Focused review merge could retain obsolete findings or incomplete coverage. | `reasoning.md#assumptions-and-risks` | Promotion requires same-fingerprint stored coverage for every non-rerun reviewer. | closed |
| Harness issue relevance may be missing, ambiguous, or externally changed. | `reasoning.md#assumptions-and-risks` | Structured issue validation and evidence hashing prevent a clean/current result when issue evidence changes. | mitigated |
| Concurrent operators can race logically valid writes. | `reasoning.md#assumptions-and-risks` | Shared per-service mutation exclusion returns a typed retryable conflict and passes race tests. | closed |
| Descendant cleanup behavior varies by OS. | `reasoning.md#assumptions-and-risks` | Existing generic process cleanup remains bounded and uncertain cleanup remains non-pass. | mitigated; gated platform evidence pending |
| Shared semantic projection may grow into a god DTO. | `reasoning.md#potential-technical-debt` | Projection contains product facts only; provider/process payloads remain outside it. | mitigated |
| Diagnostic override may be mistaken for acceptance. | `reasoning.md#assumptions-and-risks` | Confirmation and rationale are mandatory; diagnostic results never promote canonical smoke or improve assessment. | closed |
| Stable JSON expansion creates compatibility obligations. | `reasoning.md#decision-7-freeze-safe-execution-and-expose-one-semantic-result` | Versioned final-result-only JSON and semantic parity tests establish the new contract. | mitigated |

## Assumptions And Open Questions

No core architecture question remains open; `reasoning.md` fixes all seven decisions. The following are implementation evidence checkpoints, not permission to invent fallback behavior:

| Item | Required Resolution Before Dependent Work |
| --- | --- |
| Sprint 26/27 composability | Confirm the existing review and smoke services can be converged behind the sprint-owned transition with a small refactor. If not, record the concrete incompatibility and seek a reasoning deviation before broad replacement. |
| Predecessor schema identity | Name and fixture the exact supported predecessor before implementing migration. Unknown historical shapes remain rejected. |
| Harness capability completeness | Confirm structured manifest support for containing scope, evidence identity, issue relevance, prerequisites, and safe invocation. Missing capabilities produce blocked recovery, not guessed behavior. |
| Platform cleanup support | Identify supported descendant-cleanup behavior and truthful failure representation per target OS; do not claim clean cancellation where cleanup is uncertain. |

## Review Inputs

Review should use:

- `sprint-index.md`
- `technical-handbook.md`
- `reasoning/architecture.md`
- `reasoning.md`
- this `plan.md`
- implementation diff
- verification evidence
- `system/protocols/architecture-review-protocol.md`
- `system/protocols/review-sprint-protocol.md`
- `system/protocols/deep-smoke-sprint-protocol.md`

## Execution Log

| Date / Step | Action | Evidence / Notes |
| --- | --- | --- |
| 2026-07-22 / planning | Created the Sprint 28 implementation plan from selected sprint evidence. | Planning only; implementation, review automation, smoke execution, issue tracking, and Git operations were not run. |
| 2026-07-22 / implementation | Composed Sprint 26 review and Sprint 27 smoke behind a sprint-owned verification transition; added v2 state/migration, fingerprints, attempts, assessment, focused reruns, CLI/JSON/TUI parity, and recovery tests. | Implementation diff is in `/home/antonioborgerees/coding/ultraplan-go`; no Git mutation was performed. |
| 2026-07-22 / verification | Ran targeted tests, full tests, race tests, vet, build, diff whitespace checks, and boundary searches. | All local deterministic gates passed; exact commands and results are recorded in `execute.md`. |
| 2026-07-22 / review transport recovery | Diagnosed the initial agent review startup failure as an oversized single argv prompt; replaced embedded governed content with a frozen path/digest index, isolated temporary read-only snapshot, and a 96 KiB regression budget. | Compact direct OpenCode execution and targeted review transport tests pass; canonical review rerun follows this final input update. |
| 2026-07-22 / automated review remediation | The first complete reviewer fan-out failed promotion with architecture, CLI, configuration, error, security, workflow, and handbook findings. | Moved flow policy into `internal/sprint`; removed mutable production composition globals; isolated frozen reviewer inputs; made citations, timeouts, persistence, focused smoke, stable argv, partial status, and failure JSON fail closed; added regressions and reran full/race/vet/build gates. |
| 2026-07-23 / MiniMax M3 compatibility | Set review and smoke runtime selection to `minimax-coding-plan/MiniMax-M3`; the first M3 fan-out exposed truncation of structured reviewer JSON at the generic 4 KiB display-event boundary. | Added a separately bounded 96 KiB in-memory terminal-output field for machine parsing while preserving the existing 4 KiB mapped/display event cap; focused, full, race, and vet gates pass before review retry. |
| 2026-07-23 / durable review resume | Added per-coverage review checkpoints containing status, validated result, and latest OpenCode session ID. | Compatible retries reuse completed coverage and continue unfinished reviewers with `SessionActionContinue`; `review --restart`, `verify --restart-review`, flow, and TUI restart paths discard checkpoints and start fresh sessions. Crash/resume/restart tests plus full/race/vet/build gates pass. |

## Completion Criteria

- [x] All tasks are complete or explicitly deferred with requirement impact and approval recorded.
- [x] All seven reasoning decisions and all thirteen acceptance criteria have implementation and test evidence.
- [x] Verification commands were run successfully or deferrals include owner, reason, risk, and follow-up.
- [x] Current `review.md` and `smoke.md` satisfy their validators and required protocols, or smoke truthfully records protocol-valid `blocked`/`not_applicable` evidence.
- [x] CLI text, stable JSON, and TUI agree on state, verdict, freshness, evidence, issues, assessment, and next action.
- [x] Last-complete evidence survives interruption and no prohibited product, governed-input, project-doc, or Git mutation occurs.
- [x] Interrupted review attempts persist validated per-coverage results and the latest OpenCode session IDs; compatible retries resume by default and explicit restart discards all resumable coverage and sessions.
- [x] Deviations from `reasoning.md` were recorded before implementation continued.
- [x] `review.md` can evaluate conformance without guessing implementation intent.
