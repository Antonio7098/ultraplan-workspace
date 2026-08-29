# Sprint Plan: Evidence-Producing QA and Smoke Integration

> Project: `ultraplan-go`
> Sprint: `37-evidence-qa-smoke`
> Source: `reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/roadmap.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/sprints/37-evidence-qa-smoke/requirements.md`, `projects/ultraplan-go/sprints/37-evidence-qa-smoke/code-context.md`, `projects/ultraplan-go/sprints/37-evidence-qa-smoke/sprint-index.md`, `projects/ultraplan-go/sprints/37-evidence-qa-smoke/technical-handbook.md`, `projects/ultraplan-go/sprints/37-evidence-qa-smoke/reasoning/api-design.md`, `projects/ultraplan-go/sprints/37-evidence-qa-smoke/reasoning/architecture.md`, `projects/ultraplan-go/sprints/37-evidence-qa-smoke/reasoning/frontend.md`, `projects/ultraplan-go/sprints/37-evidence-qa-smoke/reasoning.md`

This plan executes `reasoning.md`. It does not reopen ownership, isolation, quorum, persistence, smoke reuse, or adapter-authority decisions.

## Reasoning Source

- **Sprint Reasoning:** `reasoning.md`
- **Sprint Index:** `sprint-index.md`
- **Technical Handbook:** `technical-handbook.md`
- **Area Reasoning:** `reasoning/architecture.md`, `reasoning/api-design.md`, `reasoning/frontend.md`

## Sprint Status

- **Status:** implemented; post-execution dogfood and reviews deferred
- **Owner:** implementation agent under UltraPlan execute control
- **Start Date:** 2026-08-25
- **Completion Date:** 2026-08-25 (offline implementation boundary)

## Decisions To Execute

| Decision | Source Section | Execution Implication |
| --- | --- | --- |
| D-01: Keep QA policy in `internal/sprint`, generic isolation in `internal/platform/process`, and adapters thin | `reasoning.md#d-01-ownership-and-dependency-shape` | Add product-neutral copy, process, cancellation, and cleanup facts below the sprint package. Keep evidence meaning, promotion, assessment, and canonical reports in `internal/sprint`. |
| D-02: Freeze an evidence map and typed checks from current governed inputs | `reasoning.md#d-02-evidence-map-and-completeness` | Evidence plans must bind requirements, expectations, target and map identities, approved paths, explicit argv, environment names, conditions, timeouts, output limits, and cleanup before execution. |
| D-03: Give each writable worker a fresh copy workspace and schedule writable work sequentially | `reasoning.md#d-03-isolation-scheduling-and-target-immutability` | No writable child starts unless source and target identity, containment, symlink policy, native write denial, process cleanup, and removal capability are proven. Children receive copy-relative paths, not the original target path. |
| D-04: Use direct evidence, three-call semantic analysis where required, and three fresh failed-shard evaluators | `reasoning.md#d-04-verdicts-and-failed-shard-adjudication` | Validate every model result locally. Incomplete calls are blocked. Two of three valid evaluator passes may adjudicate an eligible failed shard, but no evaluator may override integrity, freshness, permission, cleanup, or coverage gates. |
| D-05: Keep run control authoritative and write private QA schema v2 with read-only v1 compatibility | `reasoning.md#d-05-run-control-and-private-state` | Every publication checks the Sprint 35 fence. New evidence records use v2 and new deterministic ID kinds. V1 records remain readable but cannot satisfy Sprint 37 passing evidence. |
| D-06: Add the closed QA suite `smoke` and reuse canonical smoke execution | `reasoning.md#d-06-canonical-smoke-suite-integration` | `qa --suite smoke` must use existing authoring, discovery, containing-suite selection, invocation, evidence validation, `smoke.md`, flow, and compatibility behavior. It creates no second smoke runner or nested durable operation. |
| D-07: Keep app and JSON v1 additive and render one bounded projection | `reasoning.md#d-07-app-cli-tui-and-browser-contracts` | Add focused evidence, adjudication, issue, assessment, and smoke-suite queries behind `internal/app`. CLI, TUI, and browser render app facts and never read private QA files directly. |
| D-08: Make failures structured, redacted, and non-pass | `reasoning.md#d-08-and-d-09-errors-observability-configuration-and-bounds` | Distinguish blocked setup, assertion failure, stale input, permission denial, timeout, cancellation, cleanup uncertainty, persistence failure, and invalid evidence across state, errors, exits, and views. |
| D-09: Validate finite budgets before work | `reasoning.md#d-08-and-d-09-errors-observability-configuration-and-bounds` | Add explicit defaults and hard maxima for copy size, files, plans, checks, calls, patches, outputs, duration, retries, issues, storage, history, progress, and response pages. |
| D-10: Require deterministic tests, race checks, reviews, and hostile real-runtime smoke | `reasoning.md#d-10-acceptance-and-documentation` | The sprint cannot complete on fake-only evidence. Gated dogfood must prove target identity preservation, workspace cleanup, adjudication audit, and smoke parity. |

## Requirements / Contracts To Satisfy

The trace IDs below are plan-local names for groups of acceptance criteria in `requirements.md`. They do not replace the authoritative requirement text.

| Trace ID | Required Behavior | Evidence Planned |
| --- | --- | --- |
| ENTRY | Sprint 36 has a current acceptable Conformance Review and required smoke evidence before writable QA admission. | Admission fixtures for missing, stale, malformed, blocked, and current Sprint 36 evidence; a gated real-workspace preflight record. |
| ISO | Each writable shard attempt gets one fresh contained workspace. Copy isolation supports dirty and non-Git targets and rejects special files, unsafe links, escapes, races, and over-budget trees. | `internal/platform/process/isolation_test.go`, `internal/sprint/qa_investigation_test.go`, Linux/macOS capability tests where supported. |
| IMMUTABLE | Target source, tests, Git state, governed inputs, historical theory records, and undeclared harness paths remain unchanged. Drift blocks promotion and assessment. | Before/after identities, denied-write fixtures, Git-state fixtures, original-path omission checks, and dogfood identity comparison. |
| PLAN | Every generated check has a frozen pre-execution plan with expectation grounding, conditions, paths, explicit argv, bounded environment, timeout, output cap, and cleanup requirements. | Schema and strict-validation tables, normalized plan golden records, stale and mismatched identity tests. |
| EXEC | Commands avoid shell interpolation, use contained cwd and bounded environment, propagate cancellation, terminate descendants, cap output, and record exit, timeout, cancellation, truncation, and cleanup truthfully. | Process fault injection, descendant tests, cancellation and timeout tests, command-result goldens. |
| EVIDENCE | Evidence binds current input, implementation, map, shard, workspace, command, patch, and external identities. Patches and bounded outputs survive workspace cleanup and are never applied. | Evidence digest tests, generated patch retention tests, workspace-removal tests, and explicit absence of apply paths. |
| BOUNDS | Writable work remains sequential. Every concurrency, attempt, call, output, storage, duration, retry, and follow-up limit is finite and validated. | Configuration boundary tests, over-limit fixtures, sequential scheduling assertions, race tests around cancellation and publication. |
| ADJ | Only the global adjudicator promotes issues, marks repair eligibility, classifies regression candidates, groups root causes, and retains rejected evidence. | `qa_adjudication_test.go` tables for promotion, rejection, freshness, setup, containment, repeatability, grouping, dissent, and hostile model output. |
| ASSESS | Product code derives canonical QA assessment from current Conformance Review, accepted evidence, adjudication, blockers, and containing smoke evidence. It cannot change the review verdict. | Assessment matrix covering incomplete, blocked, fail, not applicable, pass with findings, and pass; model and narrow-smoke non-authority tests. |
| REPORT | `qa.md` is atomic, deterministic, current-state-derived, complete, and last-complete preserving. Current failed work remains visible separately. | Markdown golden tests, injected write failures, pointer and digest checks, recovery tests, current-failure versus last-complete fixtures. |
| STATE | Detailed QA state is private, versioned, bounded, contained, digest-linked, dependency-ordered, and fenced. `flow-state.json` remains a bounded projection. | V1 compatibility and v2 write tests, unknown-major rejection, strict JSON, path/symlink rejection, stale writers, partial publication, retention, and recovery. |
| SMOKE | `qa --suite smoke` uses the existing manifest-driven smoke implementation and preserves review gating, authoring scope, selection, containing suite, argv, environment, timeout, cancellation, cleanup, evidence, diagnostic-only behavior, `smoke.md`, and external authority. | Paired parity fixtures through `smoke` and QA, plus gated real-harness parity evidence. |
| RUN | Runtime-backed QA and smoke-suite work uses Sprint 35 acceptance, claims, fencing, heartbeats, replay, cancellation, retention, reconciliation, and terminal arbitration. | Durable run tests proving one accepted QA run, no nested smoke run, stale-owner rejection, replay, reconnect, cancellation, and terminal races. |
| SURFACES | CLI text/JSON, app DTOs, TUI, browser HTML/JSON, durable run detail, `qa.md`, `smoke.md`, and verification state agree on authority-bearing facts. | One shared parity fixture, JSON and HTML goldens, hostile-content tests, keyboard tests, no-JavaScript snapshots, reconnect and restart tests. |
| DOCS | CLI, architecture, user, browser, recovery, schema, and release documents describe writable QA, smoke compatibility, limits, failures, cancellation, and recovery accurately. | Documentation review against executable fixtures and command help. |
| DOGFOOD | A gated real-repository run produces contained evidence and one rejection or promotion audit, leaves the target unchanged, cleans its workspace, and proves smoke parity. | Gated run IDs, target identity pair, cleanup facts, adjudication record, smoke run identity, and review artifacts. |
| SCOPE | No production repair, patch application, permanent test promotion, planning-stage expansion, issue tracker, retrieval, hosted execution, alternate persistence, or Git mutation enters the sprint. | Diff and dependency review plus negative API and command tests. |

## Tasks

- [x] **Task 1: Enforce the Sprint 36 writable-admission gate and freeze policy limits**
  > Executes: `D-02`, `D-08`, `D-09`; `ENTRY`, `BOUNDS`, `SCOPE`
  - [ ] Extend `internal/sprint/qa_types.go`, `internal/sprint/qa_map.go`, and `internal/sprint/service.go` with typed Sprint 37 settings, limits, admission facts, and stable error categories. Keep `VerificationPhase` separate from `PlanningStage`.
  - [ ] Require current validated Sprint 36 Conformance Review, deterministic map coverage, read-only investigation evidence, cancellation/resume/invalidation/synthesis proof, and required containing smoke evidence before enabling a writable evidence plan.
  - [ ] Make missing, stale, diagnostic-only, narrow-only, malformed, or blocked prerequisites return `blocked` before workspace creation, runtime construction, or child execution.
  - [ ] Set writable investigator concurrency to one regardless of the older read-only investigator setting. Reject any configuration that attempts to raise it in Sprint 37.
  - [ ] Define and validate practical defaults and hard maxima for tree files and bytes, file size, generated checks, commands, output, patches, evidence records, issues, model calls, retries, timeouts, cleanup, retained attempts, and state size.
  - [ ] Add admission and budget tests to `internal/sprint/qa_test.go`, `internal/sprint/qa_state_test.go`, and configuration/app fixtures.
  - [ ] Stop condition: no later task may enable writable child work until all admission and limit tests pass.

- [x] **Task 2: Add product-neutral local isolation mechanics**
  > Executes: `D-01`, `D-03`, `D-08`, `D-09`; `ISO`, `IMMUTABLE`, `EXEC`, `BOUNDS`
  - [ ] Add `internal/platform/process/isolation.go` with request and result types for private workspace creation, bounded local copy, source identity facts, protected roots, contained cwd, explicit process execution, change capture, descendant cleanup, and workspace removal.
  - [ ] Keep the package free of sprint types, requirement IDs, evidence verdicts, issue semantics, and smoke selection.
  - [ ] Copy the exact local tree, including dirty and untracked files, without requiring Git. Reject devices, sockets, FIFOs, unsafe hard links, absolute links, escaping links, path races, and entries outside declared limits.
  - [ ] Create one `0700` workspace per shard attempt. Commands in that attempt use only its assigned workspace. Never reuse the workspace across shard attempts, retries, analyzers, or evaluators.
  - [ ] Run commands through explicit executable and argv fields, a contained cwd, an exact allowlisted environment, finite timeout, bounded streams, process-group ownership, context cancellation, and bounded TERM-to-KILL cleanup.
  - [ ] Expose capability facts for native protected-path write denial, process containment, descendant cleanup, and workspace removal. If the running platform cannot prove a required capability, return an unsupported or incomplete fact so sprint policy can block writable work.
  - [ ] Preserve the primary process outcome separately from cleanup and removal outcomes. Never report cleanup complete merely because the group leader exited.
  - [ ] Add `internal/platform/process/isolation_test.go` with injected creation, traversal, special-file, copy, command, timeout, cancellation, descendant, change-capture, cleanup, and removal failures.
  - [ ] Add race coverage for simultaneous cancellation and process exit. Gate platform-specific proof where the operating system lacks the required primitive.
  - [ ] Stop condition: unsupported native write denial, descendant cleanup, or removal proof must keep writable QA blocked on that platform.

- [x] **Task 3: Define private QA v2 evidence, patch, adjudication, issue, and assessment records**
  > Executes: `D-02`, `D-04`, `D-05`, `D-09`; `PLAN`, `EVIDENCE`, `ADJ`, `ASSESS`, `STATE`
  - [ ] Extend `internal/sprint/qa_types.go` and `internal/sprint/qa_state.go` with schema v2 records for frozen plans, generated checks, workspace and target identities, command results, evidence, patches, repeatability, containment, cleanup, adjudication, root-cause groups, promoted issues, regression candidates, assessment, canonical report references, and current failure.
  - [ ] Add deterministic scoped IDs for evidence, patch, adjudication, issue, and assessment without changing existing `qa-v1` attempt, map, shard, theory, challenge, and synthesis identities or claiming global content identity.
  - [ ] Add contained paths for `evidence/<evidence-id>.json`, `patches/<patch-id>.patch`, `adjudication.json`, and `issues.json` under the current attempt.
  - [ ] Normalize generated patches, cap bytes and changed paths, bind them to base and workspace identities, and store them as immutable evidence. Add no target apply operation.
  - [ ] Implement strict local validation for every record, including schema, IDs, references, digests, limits, freshness, paths, symlinks, file modes, and cross-record ownership.
  - [ ] Read v1 records through an explicit compatibility path that projects adjudication and assessment as unavailable or incomplete. Write v2 only and never rewrite historical v1 records during reads.
  - [ ] Publish immutable patches, evidence, adjudication, issues, assessment, and a validated report candidate before changing canonical pointers. Commit `qa.md`, the current state pointer, and the bounded flow projection under one fenced publication protocol.
  - [ ] Stage canonical replacements and restore the prior `qa.md`, state pointer, and flow projection if a later commit step fails. Recheck the writer fence before each canonical change. Recovery must expose the failed current attempt without inventing success.
  - [ ] Expand `internal/sprint/qa_state_test.go` for v1 compatibility, v2 strictness, immutable records, bounds, digests, atomic failure, stale writers, path escapes, symlinks, retention, invalidation, and recovery.

- [x] **Task 4: Implement isolated evidence-producing investigation**
  > Executes: `D-02`, `D-03`, `D-04`, `D-08`, `D-09`; `ISO`, `IMMUTABLE`, `PLAN`, `EXEC`, `EVIDENCE`, `BOUNDS`
  - [ ] Add `internal/sprint/qa_investigation.go` as the product policy layer over the generic isolation mechanics.
  - [ ] Implement the closed check kinds fixed by reasoning: fact, negative, behavioral, semantic, and adversarial. Require declared direct observations for every kind and exactly three fresh model analyzers only where semantic or adversarial interpretation is required.
  - [ ] Freeze each generated check before any workspace write. Record theory and expectation references, confirmation, refutation and inconclusive conditions, approved paths, exact argv, environment names, timeout, output cap, cleanup requirements, and all governing fingerprints.
  - [ ] Build isolated runtime and command requests from copy-relative paths and opaque target identities. Assert that the original target path is absent from child prompts, argv, environment, cwd metadata, progress, and retained diagnostics.
  - [ ] Use default-deny permissions. Permit only targeted tests, fixtures, probes, smoke scenarios, and bounded experiments in declared workspace paths. Deny production repair, expectation weakening, scope growth, Git commands, shell interpolation, direct issue promotion, private-state writes, and protected-root writes.
  - [ ] Capture target and governed-input identities before creation, after each command and runtime call, after cleanup, and before promotion. Distinguish attributed protected writes from unrelated drift, but block trust in either case.
  - [ ] Execute copy-backed workers sequentially. Give every retry and analyzer a new workspace and identity. Provider calls may run concurrently only within the fixed per-check cap and only when all call identities and failures remain observable.
  - [ ] Preserve bounded command results, generated patch bytes, output digests, truncation, timeout, cancellation, permission events, containment, and cleanup facts before deleting the workspace.
  - [ ] Treat missing, malformed, stale, escaped, truncated where completeness is required, cancelled, timed-out, permission-denied, provider-failed, or cleanup-uncertain evidence as blocked and non-promotable.
  - [ ] Add `internal/sprint/qa_investigation_test.go` for copy and non-Git targets, dirty/untracked content, identity and containment, path and symlink escape, original-path leakage, denied target and governed-input writes, Git denial, bounded execution, cancellation, descendants, cleanup failure, drift attribution, and evidence survival.

- [x] **Task 5: Add deterministic global adjudication and bounded issue promotion**
  > Executes: `D-02`, `D-04`, `D-08`, `D-09`; `ADJ`, `ASSESS`, `EVIDENCE`, `BOUNDS`
  - [ ] Add `internal/sprint/qa_adjudication.go` as a pure product operation over the frozen plan and admitted evidence. It must not execute commands, inspect live files, or trust prose outside validated fields.
  - [ ] Validate expectation grounding, current fingerprints, setup validity, containment, confirmation-condition fidelity, repeatability or deterministic sufficiency, flakiness, external identities, severity, root-cause grouping, and evidence sufficiency.
  - [ ] Preserve confirmed, refuted, invalid, inconclusive, blocked, cross-shard, and not-applicable theories plus every rejected evidence record and stable rejection code.
  - [ ] For semantic or adversarial checks that require model analysis, require exactly three locally validated fresh analyzer calls and a strict two-of-three majority after all calls complete.
  - [ ] Freeze initially failed shard evidence and run exactly three fresh evaluator sessions. Permit a pass adjudication only when all three results are valid, bind the same evidence digest, and at least two pass. Any missing or invalid evaluator result blocks the shard.
  - [ ] Make integrity, freshness, containment, target mutation, governed-input mutation, permission, cleanup, and coverage failures non-overridable.
  - [ ] Promote an issue only from admitted current evidence with the required independent support. Group equivalent manifestations by claim, issue class, and normalized location while retaining dissent and exact supporting evidence IDs.
  - [ ] Compute severity, repair eligibility, and regression-candidate classification in product code. Model output may propose text or grouping but cannot promote, classify, or set assessment.
  - [ ] Bound issue, rejection, group, follow-up, and model-call counts. Sort all persisted outputs deterministically.
  - [ ] Add `internal/sprint/qa_adjudication_test.go` with promotion, rejection, deterministic pass, majority, disagreement, incomplete calls, flakiness, stale identity, invalid setup, uncontained evidence, cleanup uncertainty, external evidence mismatch, root-cause grouping, severity, dissent, and hostile model-output cases.

- [x] **Task 6: Extend QA orchestration, resume, cancellation, assessment, and canonical `qa.md`**
  > Executes: `D-02`, `D-03`, `D-04`, `D-05`, `D-08`, `D-09`; `RUN`, `ASSESS`, `REPORT`, `STATE`, `SURFACES`
  - [ ] Extend `internal/sprint/qa.go` to run admission, sequential isolated attempts, evidence publication, adjudication, assessment, and canonical publication under one Sprint 35 writer token and mutation boundary.
  - [ ] Stop new scheduling on cancellation. Propagate cancellation to active runtime and process work, run bounded cleanup with an independent cleanup context, and preserve already completed valid evidence.
  - [ ] Resume only the current semantic attempt. Reuse completed records only after digest, identity, freshness, containment, and cleanup validation. Never adopt a prior worker or replay a command merely because a record is missing.
  - [ ] Derive assessment in product code from current independent Conformance Review, complete required QA evidence, adjudication, blockers, and required containing smoke evidence. Preserve the existing assessment vocabulary and never alter the review verdict.
  - [ ] Ensure stale, malformed, flaky, ungrounded, diagnostic-only, narrow-only, failed-setup, missing-cleanup, uncontained, cancelled, timed-out, or incomplete evidence cannot yield `pass`.
  - [ ] Render deterministic `qa.md` with input fingerprint, map and shard coverage, evidence quality, rejected evidence, promoted issues, regression candidates, smoke evidence, assessment, blockers, and next action.
  - [ ] Write `qa.md` atomically only after state validation. Keep the last complete report when a newer attempt fails, and expose that failure separately through state and app projections.
  - [ ] Keep `flow-state.json` limited to phase, freshness, assessment, bounded counts, cancellation, next action, current attempt, and contained state/report pointers and digests.
  - [ ] Expand `internal/sprint/qa_test.go` for writable bounds, read-only compatibility, current publication, assessment matrices, cancellation, timeout, resume, stale evidence, cleanup uncertainty, report preservation, no repair, and no Git mutation.
  - [ ] Add golden tests for `qa.md` and cross-check every authority-bearing value against the assessment and state records.

- [x] **Task 7: Wrap canonical smoke as the QA `smoke` suite without forking behavior**
  > Executes: `D-01`, `D-05`, `D-06`, `D-08`; `SMOKE`, `RUN`, `ASSESS`, `SURFACES`
  - [ ] Add the narrow sprint-owned smoke executor adapter in `internal/sprint/smoke.go` and related QA files. Both entry points must call the same existing authoring, static preparation, discovery, selection, process, evidence-validation, verdict, `smoke.md`, flow, roadmap-reconciliation, and publication code.
  - [ ] Avoid calling `RunSmoke` beneath an already held sprint mutation lock. Factor only the minimum shared top-level execution seam needed to guarantee one lock, one harness invocation, one `smoke.md` write, one flow update, and one roadmap reconciliation.
  - [ ] Keep `smoke` behavior and public results unchanged. A QA-suite invocation has one durable `qa-start` run and no nested `smoke-start` run.
  - [ ] Preserve manifest protocol checks, authoring-path bounds, explicit argv, environment allowlisting, contained cwd, timeout, cancellation, descendant cleanup, evidence schemas and identities, review gate, diagnostic-only behavior, and canonical-versus-narrow rules.
  - [ ] Store only validated links, identities, bounded QA evidence, and assessment inputs. Keep raw run JSON, stdout/stderr, per-test artifacts, and harness issue files under manifest-declared external roots.
  - [ ] Reject `qa resume --suite smoke`. A retry creates a new durable QA operation and new external smoke run while retaining prior evidence.
  - [ ] Expand `internal/sprint/smoke_test.go` with paired compatibility and QA-suite fixtures comparing selection, containing tests, argv, environment, timeout, cancellation, cleanup, external run ID, evidence links, verdict, flow projection, `smoke.md`, and failure preservation.
  - [ ] Add tests for blocked review, missing coverage, unavailable prerequisites, diagnostic scope, narrow rerun, malformed external evidence, authoring escape, and no nested operation.

- [x] **Task 8: Extend adapter-independent app use cases and durable operations**
  > Executes: `D-01`, `D-05`, `D-06`, `D-07`, `D-08`, `D-09`; `RUN`, `SURFACES`, `BOUNDS`
  - [ ] Extend `internal/app/sprint_usecases.go` with the closed QA suite field and bounded focused queries for evidence, adjudication, issues, assessment, and smoke-suite status.
  - [ ] Add request fields for evidence ID, issue ID, opaque issue cursor, and page limit. Validate each method's exact accepted fields. Default issue pages to 50 and reject limits over 200.
  - [ ] Project summaries rather than persistence records. Omit raw stdout/stderr, environment values, model payloads, arbitrary paths, and unrestricted patch bodies. Include explicit omitted counts and truncation.
  - [ ] Extend QA status with assessment, canonical report, evidence/adjudication/issue counts, bounded promoted issues, regression-candidate count, smoke suite, cleanup, and current failure.
  - [ ] Extend `internal/app/operations.go`, `internal/app/operation_runner.go`, and `internal/app/durable_operations.go` so `Suite: "smoke"` is valid only for QA start and dry-run, mutually exclusive with shard focus, included in confirmation and correlation, and rejected for resume, status, recovery, and cancellation requests.
  - [ ] Use existing run-control acceptance, ownership, heartbeat, fencing, replay, cancellation, reconciliation, and terminal APIs. Do not add a QA operation registry, event journal, or lifecycle store.
  - [ ] Preserve authorization-independent observation while requiring current authority for mutation and cancellation.
  - [ ] Expand `internal/app/sprint_usecases_test.go` and durable-operation tests for bounded projections, hostile text, current-pointer-only IDs, cursor invalidation, adapter-independent results, one-run smoke routing, fencing, cancellation, and consistent next actions.

- [x] **Task 9: Add CLI text/JSON controls and compatibility tests**
  > Executes: `D-06`, `D-07`, `D-08`, `D-09`; `SURFACES`, `SMOKE`, `DOCS`
  - [ ] Extend `internal/app/sprint_commands.go` to accept `--suite smoke` only for run and dry-run while preserving existing action words, `--shard`, `--run`, `--json`, compatibility aliases, and argument order.
  - [ ] Stop describing all QA as read-only. Explain normal isolated writable evidence, target immutability, and external smoke evidence ownership in help and text output.
  - [ ] Keep the public JSON envelope and result schema at version 1. Add fields only and preserve existing field names, types, and zero-value behavior.
  - [ ] Keep stable results on stdout and progress/diagnostics on stderr. Map assessment and failure categories to the established usage, validation, runtime, cancellation, and partial exit classes.
  - [ ] Make status and focused read queries return success when they truthfully report a blocked or failed product assessment. Execution exit behavior must follow the deterministic assessment and cleanup result.
  - [ ] Expand `internal/app/sprint_commands_test.go` with help, old-call compatibility, suite parsing, invalid combinations, dry-run, no-resume smoke, JSON goldens, stdout/stderr separation, exits, cancellation, blockers, and app/CLI agreement.

- [x] **Task 10: Extend TUI and browser QA presentation over app facts**
  > Executes: `D-01`, `D-07`, `D-08`, `D-09`; `SURFACES`, `RUN`, `REPORT`, `SMOKE`
  - [ ] Extend `internal/tui/qa_view.go` to show identity and freshness, assessment, current failure versus last complete report, operation and cleanup, coverage, bounded evidence, rejected evidence, adjudication, issues, regression candidates, smoke suite, blockers, cancellation, recovery, and next action.
  - [ ] Preserve verdict-neutral phase language. Render control bytes and ANSI input inert before width calculations. Keep status text visible without color and stack authority-bearing fields at narrow widths.
  - [ ] Expose refresh, cancel, resume, recovery, and smoke-suite start only when app results mark each action valid. Exiting or resizing stops observation only.
  - [ ] Extend `internal/web/qa_handlers.go` with app-backed GET routes for evidence, adjudication, paged issues, one issue, assessment, and smoke-suite status. Do not import `internal/sprint` or expose caller-selected evidence paths.
  - [ ] Extend `internal/web/operation_handlers.go` to accept only `options.suite: "smoke"` for QA start and dry-run. Preserve strict JSON, body limits, same-origin, CSRF, confirmation, authorization, and request-fingerprint checks.
  - [ ] Extend `internal/web/templates/sprint.html` and `internal/web/templates/run_qa.html` with complete server-rendered no-JavaScript QA facts and ordinary guarded forms. Keep templates free of file reads, app calls, request validation, and product-state decisions.
  - [ ] Enhance `internal/web/static/js/operations.js` only for bounded progress, inline cancellation, and authoritative refresh. Treat events as hints and refresh on terminal events, gaps, disconnect, session rotation, or observer restart.
  - [ ] Escape all hostile evidence through existing template and preview boundaries. Never place evidence in raw HTML, scripts, styles, or unvalidated URLs.
  - [ ] Expand `internal/tui/qa_view_test.go`, `internal/web/qa_handlers_test.go`, operation contract tests, and template tests for keyboard use, narrow/mobile layouts, no-JavaScript completeness, focus, hostile text, bounded rendering, cancellation, dropped delivery, reconnect, restart, session rotation, confirmation staleness, and parity fixtures.

- [x] **Task 11: Document the supported workflow, schemas, limits, and recovery**
  > Executes: `D-06`, `D-07`, `D-08`, `D-09`, `D-10`; `DOCS`, `SCOPE`
  - [ ] Update `../ultraplan-go/docs/cli-reference.md` for isolated QA, `qa --suite smoke`, dry-run, focused shards, status, cancellation, resume restrictions, JSON outcomes, exits, blockers, and compatibility commands.
  - [ ] Update `../ultraplan-go/docs/architecture.md` for product/platform ownership, copy and native-isolation guarantees, evidence authority, quorum limits, adjudication, state v2, publication order, run control, and smoke reuse.
  - [ ] Update `../ultraplan-go/docs/user-guide.md` for evidence plans, generated checks, outcomes, rejected evidence, promoted issues, regression candidates, canonical assessment, smoke suites, and safe operator actions.
  - [ ] Update `../ultraplan-go/docs/local-web.md` for guarded starts, no-JavaScript inspection, durable progress, cancellation, reconnect, session rotation, and recovery.
  - [ ] Update `../ultraplan-go/docs/recovery.md` for failed isolation, target drift, stale evidence, interrupted attempts, stale writers, invalid state, missing cleanup, current failure, last complete report, and safe restart.
  - [ ] Update `../ultraplan-go/docs/phase3-json-schemas.md` for additive app/JSON v1 fields, private QA v2 records, IDs, pointers, issues, assessment, smoke suite, compatibility, and migration behavior.
  - [ ] Update `../ultraplan-go/docs/release-checklist.md` for non-mutation, isolation, evidence quality, issue audit, smoke parity, cross-surface fixtures, race/build checks, and gated dogfood.
  - [ ] State the limits honestly. Copy workspaces are not a general hostile multi-tenant sandbox, majority is not proof of independent truth, and unsupported native isolation blocks writable work.
  - [ ] Verify that docs add no repair, patch application, permanent test promotion, issue-management, Git mutation, content identity, retrieval, or remote execution instructions.

- [/] **Task 12: Complete deterministic, race, parity, and gated release evidence** — Deferred: offline, race, parity, vet, build, and hygiene gates passed; Architecture Review, Sprint Review, Deep Smoke, and real external-harness dogfood belong to the post-execution stages, and `ULTRAPLAN_REAL_SMOKE` is not enabled.
  > Executes: `D-10`; `ISO`, `IMMUTABLE`, `PLAN`, `EXEC`, `EVIDENCE`, `ADJ`, `ASSESS`, `REPORT`, `STATE`, `SMOKE`, `RUN`, `SURFACES`, `DOGFOOD`, `SCOPE`
  - [ ] Build one canonical fixture and inspect it through private verification state, `qa.md`, `smoke.md`, app DTOs, CLI text/JSON, TUI, browser HTML/JSON, and durable run detail. Compare semantic fields rather than prose layout.
  - [ ] Run all package-focused unit and fault-injection suites for process isolation, investigation, state, adjudication, orchestration, smoke, app, CLI, TUI, and web.
  - [ ] Run `go test ./...`, `go test -race ./...`, `go vet ./...`, `go build ./cmd/ultraplan`, and `git diff --check` in `../ultraplan-go`.
  - [ ] Run Architecture Review and Sprint Review against the implementation diff, tests, docs, state schemas, dependency direction, target non-mutation, and non-goals.
  - [ ] Gate real-repository dogfood on all prerequisites. If runtime, native isolation, Sprint 36 evidence, or harness prerequisites are missing, record `blocked`; do not count the gate as passed.
  - [ ] In dogfood, produce at least one generated contained patch, preserve it after workspace cleanup, and record at least one adjudication rejection or promotion audit.
  - [ ] Prove the target and governed inputs have identical before/after identities, Git state is unchanged, no isolated workspace remains, and cleanup certainty is complete.
  - [ ] Run the same canonical containing smoke selection through `smoke` and `qa --suite smoke`. Compare protocol, selected suites/tests, argv, environment, timeout, cancellation and cleanup facts, external evidence identities, verdict, flow projection, and `smoke.md`.
  - [ ] Record durable run IDs, operational attempt IDs, fences, external smoke run IDs, evidence digests, target identities, cleanup facts, duration/cost, and review outcomes in the gated evidence record.
  - [ ] Stop condition: any target mutation, uncontained evidence, cleanup uncertainty, parity drift, stale writer, malformed state, unexplained model result, or failed required check blocks sprint completion.

## Evidence Checklist

- [ ] Sprint 36 admission evidence is current and acceptable before writable work.
- [ ] Tests prove copy and non-Git isolation, dirty/untracked fidelity, path and symlink safety, protected-root denial, original-path omission, target identity, cancellation, descendants, and cleanup failure.
- [ ] Frozen evidence plans and every identity relationship have strict schema and digest tests.
- [ ] Generated patches and bounded outputs survive workspace cleanup and no code path applies them.
- [ ] Adjudication tests prove invalid, stale, flaky, narrow, diagnostic, ungrounded, uncontained, incomplete, and cleanup-uncertain evidence cannot promote or pass.
- [ ] Assessment tests prove only product code derives the result and cannot alter Conformance Review.
- [ ] V1 read compatibility, v2 writes, unknown-major rejection, atomic publication, stale-writer fencing, and last-complete preservation pass.
- [ ] `smoke` and `qa --suite smoke` use the same authoritative path and produce parity evidence.
- [ ] CLI, JSON, TUI, browser, durable run, Markdown, and state agree for one shared fixture.
- [ ] Browser security, no-JavaScript, hostile content, reconnect, restart, session rotation, cancellation, and accessibility tests pass.
- [ ] Documentation updates match executable behavior and exact state paths.
- [ ] Architecture Review, Sprint Review, and gated Deep Smoke Sprint evidence exist.
- [ ] Real dogfood proves unchanged target identity, complete cleanup, one adjudication audit, and smoke parity.
- [ ] Diff review confirms no production repair, patch application, permanent test promotion, general issue tracker, planning-stage expansion, retrieval, hosted operation, alternate persistence, or Git mutation.

## Verification Commands

All implementation commands run from `../ultraplan-go` unless noted otherwise.

| Check | Command | Expected Result |
| --- | --- | --- |
| Generic isolation | `go test ./internal/platform/process -run 'Test(Isolation|DirectRunner)' -count=1` | Creation, copy, containment, process, cancellation, cleanup, and removal fixtures pass. |
| Sprint investigation | `go test ./internal/sprint -run 'TestQA(Investigation|EvidencePlan|Writable|Permission|Cancellation|Cleanup)' -count=1` | Isolated evidence plans and target non-mutation cases pass. |
| State and adjudication | `go test ./internal/sprint -run 'TestQA(Store|State|Adjudication|Assessment|Atomic|Recovery)' -count=1` | V2 state, v1 compatibility, promotion rules, assessment, fencing, and recovery pass. |
| Smoke parity | `go test ./internal/sprint -run 'Test(Smoke|QA).*Parity|TestQA.*Smoke' -count=1` | Compatibility and QA-suite paths have identical authoritative smoke behavior. |
| App and CLI | `go test ./internal/app -run 'Test(QA|SprintQA|DurableQA)' -count=1` | Request rules, bounded projections, JSON, exits, durable ownership, and compatibility pass. |
| TUI and browser | `go test ./internal/tui ./internal/web -run 'TestQA|TestBrowserQA|TestOperation.*QA' -count=1` | Bounded rendering, keyboard/no-JS behavior, security, cancellation, and recovery pass. |
| Offline suite | `go test ./...` | All deterministic tests pass. |
| Race suite | `go test -race ./...` | No race in isolation, cancellation, publication, durable ownership, progress, or terminal arbitration. |
| Static analysis | `go vet ./...` | No vet findings. |
| Build | `go build ./cmd/ultraplan` | The CLI builds successfully. |
| Diff hygiene | `git diff --check` | No whitespace errors. |
| Gated dogfood | Documented release-checklist command for Sprint 37, with required runtime and harness environment | Produces contained evidence and adjudication audit, leaves target and Git state unchanged, removes workspaces, and proves smoke parity. Missing prerequisites produce `blocked`. |

## Risks And Blockers

| Risk / Blocker | Source | Mitigation | Status |
| --- | --- | --- | --- |
| Native protected-root write denial or descendant cleanup cannot be proven on a supported platform | `reasoning.md`, D-03 | Expose capability facts, test each supported OS, and block writable QA where proof is unavailable. Do not weaken the guarantee silently. | resolved on recorded Linux host; unsupported hosts block |
| Sprint 36 admission evidence is missing, stale, or unacceptable | `requirements.md`, Dependencies and first acceptance criterion | Fail before workspace creation or runtime work. Rerun the earlier governed review and containing smoke outside this sprint plan. | resolved for entry: current pass-with-findings review and passing containing smoke |
| Original target path leaks through prompt, argv, environment, cwd metadata, logs, or events | `reasoning.md`, Risk register | Construct child requests from copy-relative paths, strip environment, scan every retained channel in adversarial tests, and block admission on leakage. | open |
| Copying large dirty trees exceeds practical limits | `reasoning.md`, D-03 and D-09 | Validate file, byte, duration, and state budgets before copy. Block oversized targets and measure dogfood before considering optimization. | open |
| Fresh model calls share provider bias | `reasoning.md`, Known technical debt | Keep direct evidence primary, validate outputs locally, preserve disagreement, require current Conformance Review, and describe majority limits honestly. | open |
| V1 compatibility is mistaken for current Sprint 37 evidence | `reasoning/api-design.md`, Decision 5 | Project missing adjudication and assessment as unavailable or incomplete. Only valid v2 evidence may pass. | open |
| Multi-file publication is interrupted between state and Markdown updates | `reasoning/architecture.md`, Persistence and publication | Publish in dependency order with digests and fencing, preserve prior report, and make recovery expose incomplete current work. | open |
| `qa --suite smoke` deadlocks under nested mutation locks or duplicates side effects | `reasoning/api-design.md`, Risks | Use one top-level lock and one canonical smoke executor. Assert one durable run, invocation, artifact write, flow update, and roadmap reconciliation. | open |
| Narrow or diagnostic smoke is presented as canonical | `requirements.md`, smoke criteria | Retain containing-suite and diagnostic facts in product state and every renderer. Assessment rejects narrow-only evidence. | open |
| Bounded DTOs omit audit-critical reasons | `reasoning/api-design.md`, Risks | Always retain reason codes and exact evidence IDs, expose focused queries, and report omitted counts and truncation. | open |
| Cancellation races completion, cleanup, or stale ownership | `reasoning.md`, D-05 and D-08 | Use Sprint 35 terminal arbitration and fences, preserve completed evidence, and project the authoritative terminal result after refresh. | open |
| Dogfood prerequisites are unavailable | `requirements.md`, Release and dogfood gate | Record `blocked` and do not satisfy the sprint gate. Normal deterministic checks still run. | blocked: real harness opt-in is not enabled; reviews are downstream |

## Review Inputs

Review should use:

- `sprint-index.md`
- `technical-handbook.md`
- `reasoning/architecture.md`
- `reasoning/api-design.md`
- `reasoning/frontend.md`
- `reasoning.md`
- this `plan.md`
- the implementation diff in `../ultraplan-go`
- deterministic, race, build, and gated dogfood evidence
- `system/protocols/architecture-review-protocol.md`
- `system/protocols/review-sprint-protocol.md`
- `system/protocols/deep-smoke-sprint-protocol.md`

The review must inspect all changed package boundaries and all commits included in the implementation, not only the final diff chunk. It must verify exact requirement coverage, no target mutation, no repair path, one smoke authority, run-control fencing, state compatibility, cross-surface agreement, and truthful blocked gates.

## Execution Log

| Date / Step | Action | Evidence / Notes |
| --- | --- | --- |
| 2026-08-25 / planning | Created the implementation plan from governed Sprint 37 reasoning and area decisions. | No implementation, smoke, review automation, Git operation, or flow-state edit was performed during planning. |
| 2026-08-25 / implementation | Implemented Tasks 1 through 11 in the recorded worktree. | Added native disposable-copy isolation, v2 evidence/state/publication, adjudication and assessment, canonical report rollback, smoke-suite reuse, app/CLI/TUI/web projections, focused APIs, limits, tests, and docs. No Git operation or production repair path was added. |
| 2026-08-25 / verification | Ran focused, full offline, race, vet, build, JavaScript syntax, and diff-hygiene checks. | All deterministic gates passed. `TestRealSmokeHarness` skipped with the explicit blocker `ULTRAPLAN_REAL_SMOKE=1` not set. Actual review and smoke stages were not launched during execute. |

## Completion Criteria

- [ ] All tasks are complete or explicitly deferred with requirement impact and approval recorded.
- [ ] Writable QA cannot start without current Sprint 36 evidence and proven isolation capabilities.
- [ ] Every writable attempt uses a distinct workspace, sequential scheduling, bounded execution, and truthful cleanup.
- [ ] The target, governed inputs, Git state, historical QA records, and undeclared harness paths remain unchanged.
- [ ] Frozen plans, evidence, patches, adjudication, issues, assessment, pointers, and reports are private, versioned, bounded, digest-linked, fenced, and recoverable.
- [ ] Only global product adjudication promotes issues or changes an eligible failed-shard result. Invalid evidence cannot pass.
- [ ] `qa.md` is deterministic and atomic, and failed new work preserves the last complete report while exposing current failure.
- [ ] `flow-state.json` remains a bounded projection and Sprint 35 run control remains the operational authority.
- [ ] `qa --suite smoke` and `smoke` share one smoke implementation and preserve all compatibility and external evidence guarantees.
- [ ] CLI, JSON, TUI, browser, durable run, `qa.md`, `smoke.md`, and verification state agree for shared fixtures.
- [ ] Documentation covers commands, schemas, isolation limits, evidence authority, cancellation, recovery, compatibility, and release checks.
- [ ] `go test ./...`, `go test -race ./...`, `go vet ./...`, `go build ./cmd/ultraplan`, and `git diff --check` pass in `../ultraplan-go`.
- [ ] Gated real-repository evidence proves contained generated evidence, one adjudication audit, unchanged target identity, complete cleanup, and smoke parity. Missing prerequisites remain `blocked` and do not satisfy completion.
- [ ] Architecture Review, Sprint Review, and Deep Smoke Sprint evidence are complete.
- [ ] Deviations from `reasoning.md` were recorded and approved before implementation continued.
- [ ] `review.md` can evaluate conformance without guessing intent.
