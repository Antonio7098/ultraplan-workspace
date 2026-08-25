# Sprint Plan: Bounded Manual and Automatic Repair

> Project: `ultraplan-go`
> Sprint: `38-bounded-repair`
> Source: `reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/roadmap.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/sprints/38-bounded-repair/requirements.md`, `projects/ultraplan-go/sprints/38-bounded-repair/code-context.md`, `projects/ultraplan-go/sprints/38-bounded-repair/sprint-index.md`, `projects/ultraplan-go/sprints/38-bounded-repair/technical-handbook.md`, `projects/ultraplan-go/sprints/38-bounded-repair/reasoning/api-design.md`, `projects/ultraplan-go/sprints/38-bounded-repair/reasoning/architecture.md`, `projects/ultraplan-go/sprints/38-bounded-repair/reasoning/frontend.md`, `projects/ultraplan-go/sprints/38-bounded-repair/reasoning.md`

This plan executes `reasoning.md`. It does not reopen repair ownership, confirmation order, production-apply authority, reverification order, automatic-mode reuse, or terminal-result semantics.

## Reasoning Source

- **Sprint Reasoning:** `reasoning.md`
- **Sprint Index:** `sprint-index.md`
- **Technical Handbook:** `technical-handbook.md`
- **Area Reasoning:** `reasoning/api-design.md`, `reasoning/architecture.md`, `reasoning/frontend.md`

## Sprint Status

- **Status:** not started
- **Owner:** implementation agent
- **Start Date:** pending
- **Completion Date:** pending

## Decisions To Execute

| Decision | Source Section | Execution Implication |
| --- | --- | --- |
| D-01: one sprint-owned repair protocol | `reasoning.md#d-01-ownership-and-the-protocol-boundary` | Keep admission, scope, progress, reverification, cleanup interpretation, proof, and semantic outcomes in focused `internal/sprint` files; platform, run control, app, and adapters retain their existing authorities. |
| D-02: packet, acceptance, confirmation, dispatch | `reasoning.md#d-02-packet-preparation-acceptance-confirmation-and-dispatch` | Freeze one packet through a synchronous durable prepare operation; accept and claim mutable work before confirmation, but dispatch no goroutine or runtime until immutable confirmation publication succeeds. |
| D-03: isolated proposal and product-owned apply | `reasoning.md#d-03-proposal-scope-enforcement-and-production-apply` | Give the runtime only a bounded isolated copy; derive and validate a strict text patch in product code; apply through contained direct filesystem writes with no Git, shell, formatter, or runtime production authority. |
| D-04: strict persistence and dual ownership | `reasoning.md#d-04-persistence-fencing-terminalization-and-recovery-boundaries` | Add immutable private repair records, a digest-bound current pointer, flow-state migration, writer fencing, mutation-lease checks, apply journaling, and a `terminalizing` barrier. |
| D-05: sequential progressive reverification | `reasoning.md#d-05-progressive-reverification-and-independent-review` | Run frozen descriptors in the required widening order, skip wider gates after a required non-pass, and record a focused independent review delta without changing `review.md`. |
| D-06: one manual cycle, then proof-gated automatic reuse | `reasoning.md#d-06-manual-first-cycle-automatic-bounds-progress-and-proof` | Manual mode invokes one shared cycle once; automatic implementation and exposure begin only after a current real manual proof and reuse the same cycle with persisted lower-only budgets and product-derived progress. |
| D-07: separate lifecycle and semantic outcome | `reasoning.md#d-07-closed-semantic-outcomes-and-operational-lifecycle` | Preserve run-control lifecycle independently from the six repair outcomes and commit exactly one immutable semantic result through create-or-compare publication. |
| D-08: governed cancellation, cleanup, resume, and recovery | `reasoning.md#d-08-cancellation-cleanup-resume-shutdown-and-recovery` | Propagate canonical cancellation, use a separate bounded cleanup context, resume only from proven immutable boundaries, never repeat apply, and recover conservatively without inferring success. |
| D-09: one bounded app and interface projection | `reasoning.md#d-09-shared-app-cli-tui-browser-and-public-facts` | Add `RepairUseCases`, guarded operations, bounded queries, stable CLI/JSON, TUI routes, and server-rendered browser resources; adapters never read private records or derive outcomes. |
| D-10: manual proof before automatic acceptance | `reasoning.md#d-10-acceptance-documentation-and-rollout` | Pass focused/full/race/build checks and one real manual four-interface dogfood before implementing automatic mode; then prove a bounded automatic stop with counters and deadline surviving restart or resume. |

## Requirements / Contracts To Satisfy

| Contract / Requirement ID | Required Behavior | Evidence Planned |
| --- | --- | --- |
| AC-1 | Prepare from exactly one current repair-eligible adjudicated issue; freeze complete current authority into a deterministic immutable packet; reject stale, incomplete, unsafe, or ambiguous admission without runtime or target mutation. | Packet tables, current-attempt joins, exact-reproducer tests, protected-path matrices, digest/idempotency tests, runtime-not-called assertions. |
| AC-2 | Require explicit single-use confirmation over packet, target, request, mode, limits, governed inputs, and durable acceptance; distinguish manual and automatic consent. | Changed-field confirmation table, duplicate/replay and restart races, CLI `--yes`, TUI/browser guarded actions, automatic opt-in tests. |
| AC-3 | Permit one manual proposal, one product-owned apply, and one complete progressive reverification sequence through the same shared path on every interface. | Isolation, patch, apply, runner, app, CLI, TUI, HTTP, and one real manual dogfood trace. |
| AC-4 | Protect source of truth from direct and indirect scope escape; require one writable owner, current fence, target checks, and actual-diff enforcement. | Every protected class, traversal, rename/delete, link, hard-link, stale fence, target race, descendant/formatter side-effect, and actual-scope tests. |
| AC-5 | Reverify exact reproducer, primary shards, linked theories, boundary/follow-up shards, containing QA/smoke suites, and review delta in fixed order. | Ordered gate tests, skip-reason tests, immutable descriptor tests, target drift, issue/severity delta, contradiction, smoke, and review-delta assertions. |
| AC-6 | Keep automatic mode unavailable until a current qualifying real manual proof matches schema, code, maxima, policy, isolation, checks, governed inputs, runtime policy, cleanup, result, and target. | Proof qualification/rejection tests plus the retained real manual run and component-level invalidation diagnostics. |
| AC-7 | Bound every automatic resource and stop before another mutation on stagnation, repetition, exhaustion, reopening, scope/severity growth, uncertainty, or drift. | Config precedence/maxima tests, stop-decision tables, persisted counter/deadline tests, restart/resume integration, bounded automatic dogfood. |
| AC-8 | Derive one truthful outcome from product facts; keep cancellation and operational lifecycle separate; prevent late or stale replacement. | Six-outcome tables, lifecycle/outcome matrix, terminal compare-and-set races, cleanup-uncertainty mapping, CLI/HTTP exit/status tests. |
| AC-9 | Preserve authority across resume, cancellation, cleanup, shutdown, and startup recovery; never repeat a committed apply or infer success. | Cancellation timing, process-tree cleanup, server drain, abandoned owner, stale lock, terminalizing, partial publication, apply-journal recovery, resume-boundary tests. |
| AC-10 | Make CLI text, CLI JSON, TUI, browser HTML/JSON, durable operations, repair state, result, and flow summary agree while remaining bounded, escaped, and redacted. | Canonical fixture, pagination/cursor tests, hostile-text fixtures, reconnect/replay-gap tests, no-JavaScript flows, cross-surface parity assertions. |
| AC-11 | Pass offline, race, build, adversarial, and real-runtime gates; prove manual work before automatic work. | Focused package tests, `go test ./...`, `go test -race ./...`, `go build ./cmd/ultraplan`, three review protocols, manual and automatic dogfood evidence. |
| Architecture, Security, Persistence, Workflows | Keep product policy in `internal/sprint`, fail closed on authority or state uncertainty, use strict private state, and retain exactly one owner/result. | Import review, strict schema tests, fault injection, fencing races, mutation trace, Architecture Review. |
| CLI, Configuration, Errors, Observability, Performance | Provide explicit non-interactive controls, lower-only finite limits, stable errors/exits, bounded correlated facts, and no secret or unbounded output. | Command/help/JSON fixtures, exhaustive config tests, safe error matrix, correlation/redaction tests, bounded queries and retained-state tests. |
| LLM Runtime and Evaluation Safety | Keep runtime work isolated, bounded, cancellable, and unable to apply, expand scope, choose commands, or declare success. | Runtime request/permission tests, production-root denial, command catalog tests, model-output distrust tests, real-runtime dogfood. |
| Testing and Documentation | Cover deterministic and adversarial behavior and document exact operation, recovery, proof, and release contracts. | Required source tests, authoritative documentation updates, checked examples, Sprint Review and Deep Smoke evidence. |

## Tasks

- [ ] **Task 1: Verify Sprint 37 And Host Admission Gates**
  > Executes: `D-01`, `D-02`; `AC-1`, `AC-3`, dependency and promotion gates
  - [ ] Confirm Sprint 37 has a current evidence-producing QA attempt, current repair-eligible adjudicated issue, accepted exact reproducer, complete mapped containing smoke result of `pass` or `pass_with_open_issues`, and current real-runtime evidence; record the selected issue and all dependency fingerprints in the execution log.
  - [ ] Confirm the implementation repository exposes protected-root denial, bounded isolation, process-tree cleanup, full target identity, mutation leases, durable operations, writer fencing, and startup reconciliation; treat any absent capability as a blocker rather than weakening the protocol.
  - [ ] Confirm the selected issue supplies unambiguous production paths and frozen descriptors for every required reverification gate; do not start implementation against an issue that requires caller-authored scope or commands.
  - [ ] Stop the sprint before source changes if any dependency is stale, blocked, mismatched, or lacks real-runtime evidence; identify the exact Sprint 37 recovery action.

- [ ] **Task 2: Add Repair Domain Types And Pure Policy**
  > Executes: `D-01`, `D-02`, `D-05`, `D-06`, `D-07`; `AC-1`, `AC-2`, `AC-5`, `AC-7`, `AC-8`
  - [ ] Extend `internal/sprint/qa_types.go` with strict versioned types for packet, confirmation, cycle, scope, apply journal, progress, freshness, blocker, cleanup, reverification gates, semantic outcome, manual proof, consumed counters, absolute deadline, and lower-only budgets; keep repair under `VerificationPhaseRepair` and out of `PlanningStage`.
  - [ ] Implement deterministic normalization and pure validation for IDs, fingerprints, finite sorted path sets, closed enums, gate order, confirmation digest inputs, progress facts, automatic stop reasons, and outcome derivation in `internal/sprint/qa_repair.go`.
  - [ ] Define the complete protected-path classifier for governed sprint inputs, review/smoke/QA/repair evidence, flow and workspace state, implementation plans, configuration, Git control/hooks/ignore policy, tests, snapshots, baselines, generated evidence, links, special files, and non-production data.
  - [ ] Add table tests in `internal/sprint/qa_repair_test.go` for every valid and invalid enum, budget, path class, progress fact, stop rule, and all six outcomes; prove manual mode cannot emit `stalled` and no model string can change policy.

- [ ] **Task 3: Extend Strict Repair Persistence And Flow Projection**
  > Executes: `D-04`, `D-07`, `D-08`; `AC-1`, `AC-4`, `AC-8`, `AC-9`
  - [ ] Extend `internal/sprint/qa_state.go` with exact paths and strict readers/writers for `verification/repair-state.json`, immutable packet and confirmation, per-cycle proposal/scope/reverification/cleanup, immutable result, apply journal, and `verification/manual-repair-proof.json`.
  - [ ] Publish detail before current pointers, check the current run-control writer before every write and rename, enforce `0700` directories and `0600` files, reject links and path escapes, verify digests, and use immutable create-or-compare semantics for packet, confirmation, cycle evidence, result, and proof.
  - [ ] Extend `internal/sprint/state.go` and flow validation with one bounded repair summary and an explicit flow-state schema migration; preserve review, smoke, QA, and repair summaries during unrelated writes and keep all detailed collections out of `flow-state.json`.
  - [ ] Extend retention and verification-byte accounting without deleting current packet/result/proof evidence; expose explicit retained lower boundaries for cycle queries.
  - [ ] Add `internal/sprint/qa_state_test.go` coverage for schemas, unknown fields/versions, permissions, containment, symlink and hard-link attacks where supported, digest mismatch, stale writers, partial writes, pointer-last publication, rollback failure, immutable conflict, migration, retention, and bounded summary projection.

- [ ] **Task 4: Split Durable Acceptance From Dispatch**
  > Executes: `D-02`, `D-04`, `D-07`, `D-08`; `AC-2`, `AC-4`, `AC-8`, `AC-9`
  - [ ] Refactor `internal/app/durable_operations.go` so durable accept/claim allocates run identity, operational attempt, and fencing generation without starting ownership control, a goroutine, or runtime work; add an explicit dispatch transition that starts ownership only after confirmation publication.
  - [ ] Add closed `repair-prepare`, `repair-start`, `repair-resume`, and `repair-recover` operation kinds, canonical request fields, confirmation facts, governed inputs, and mutation classes in `internal/app/operations.go`; keep automatic mode a field of the shared start operation.
  - [ ] Update `internal/app/operation_runner.go` so it is the only repair launch site, installs the current writer fence, correlates repair and operation runs, and maps operational completion independently from repair outcome.
  - [ ] Preserve existing operation behavior and server shutdown semantics; an accepted but unconfirmed run must terminalize as a persistence/blocked failure and never dispatch.
  - [ ] Extend `internal/app/durable_operations_test.go` with acceptance-before-goroutine, confirm-before-dispatch, acceptance/confirmation failure, restart at the unconfirmed boundary, single owner, stale generation, duplicate confirmation, cancellation, late completion, and exactly-one operational terminal tests.

- [ ] **Task 5: Freeze Current Issue Packets And Bind Confirmation**
  > Executes: `D-02`, `D-04`; `AC-1`, `AC-2`, `AC-4`
  - [ ] Implement synchronous writer-fenced packet preparation in `internal/sprint/qa_repair.go` and integrate current QA reads in `internal/sprint/qa.go`: load and digest-validate the current attempt, issue, root-cause group, evidence, plans, map, shards, theories, assessment, review, containing smoke, checks, policy, implementation, governed inputs, and full target identity as one snapshot.
  - [ ] Derive all required packet fields from product records only, including violated requirements/criteria, exact reproducer, expected failing condition, affected and follow-up shards, theory IDs, allowed/forbidden paths, containing checks, repair acceptance criteria, and frozen budgets; caller input may select only the current issue ID and requested mode.
  - [ ] Make identical preparation idempotently reuse the same current packet while a confirmed or terminal attempt receives a new repair-run identity; prove preparation never initializes a runtime or changes the target, tests, QA evidence, or governed artifacts.
  - [ ] Persist single-use `confirmation.json` only after durable acceptance, binding packet digest, full target, canonical request, mode, automatic opt-in, effective limits, governed/policy fingerprints, operation run, operational attempt, fencing generation, confirmer, and timestamp.
  - [ ] Add `internal/sprint/qa_repair_test.go` and `internal/sprint/qa_test.go` coverage for every admission rejection, cross-attempt or digest mismatch, stale review/QA/smoke/target, missing exact reproducer/check, issue identity, deterministic packet bytes, changed confirmation field, replay, and no-mutation preparation.

- [ ] **Task 6: Build Isolated Proposal And Strict Patch Evidence**
  > Executes: `D-03`, `D-04`; `AC-3`, `AC-4`
  - [ ] Add the repair runtime request and service dependency in `internal/sprint/service.go` and `internal/sprint/qa_repair.go`; grant bounded packet-approved reads and isolated-copy writes only, deny writable production and verification roots, and reject unsupported isolation capabilities before runtime start.
  - [ ] Reuse `internal/platform/process` bounded isolation, tree identity, comparison, process-group cancellation, output limits, and verified removal as factual mechanisms without moving repair eligibility or outcome decisions into the platform package.
  - [ ] Derive a canonical bounded unified text proposal from isolated before/after trees and publish `proposal.patch` before production mutation; reject NUL/binary content, malformed hunks, unsupported encoding, links, hard links, special files, implicit renames, unapproved delete/add, unsafe modes, file/byte/patch excess, and production-root leakage.
  - [ ] Reparse and revalidate the retained proposal independently from runtime claims; a safely retainable violating proposal becomes rejected evidence and cannot reach apply.
  - [ ] Add real-filesystem and fake-runtime tests for allowed isolated edits, direct production-write attempts, traversal, absolute paths, link and hard-link replacement, special files, malformed patches, limits, cancellation, runtime descendants, output truncation, and cleanup failure.

- [ ] **Task 7: Implement Product-Owned Apply, Journal, And Scope Enforcement**
  > Executes: `D-03`, `D-04`, `D-08`; `AC-3`, `AC-4`, `AC-9`
  - [ ] Implement the direct contained apply boundary in `internal/sprint/qa_repair.go`: recheck packet, confirmation, target identity, mutation lease, run-control fence, path classes, changed-file/byte totals, and expected pre-image digests before staging any production bytes.
  - [ ] Stage same-directory private replacements, record pre-images and an apply journal, apply only parsed product operations, and never invoke Git, shell, command substitution, hooks, formatters, runtime tools, or repository cleanup.
  - [ ] Attempt in-process compensation after partial failure, record completed and restored operations, and classify rollback failure or crash-uncertain state as escalation with no automatic retry or second apply.
  - [ ] Recompute full target identity and actual changed paths/bytes after apply; require actual operations to equal the intended set, remain within `allowed_paths`, avoid all forbidden classes, and leave unrelated pre-existing changes untouched.
  - [ ] Add failure injection before and after each replacement and compensation step, target-drift races at every identity checkpoint, expected-preimage mismatch, lost lease/fence, stale owner, side-effect escape, partial apply, rollback failure, crash-journal recovery, and unrelated-change tests.

- [ ] **Task 8: Add Sequential Reverification And Outcome Derivation**
  > Executes: `D-05`, `D-07`; `AC-5`, `AC-8`
  - [ ] Implement the frozen ordered ladder in `internal/sprint/qa_repair.go` and `internal/sprint/qa.go`: exact reproducer, primary shards, linked theory confirmation/refutation, boundary/neighbor/parent follow-up shards, containing QA and mapped smoke suites, then focused Conformance Review delta.
  - [ ] Execute only immutable packet descriptors through explicit executable/argv, contained workdir, allowlisted environment, frozen timeout/output limits, cancellation, target checks, and process-tree cleanup; do not permit runtime-selected substitutions.
  - [ ] Stop after the first required non-pass, stale state, cancellation, or cleanup uncertainty and record every wider gate as skipped with reason and required next action.
  - [ ] Record current issue-set, reopening, severity, scope, contradiction, new-failure, and review-delta facts; leave `review.md`, its verdict, QA adjudication, smoke evidence, and checks unchanged.
  - [ ] Derive the six semantic outcomes only from persisted product facts and publish one immutable result; add order, skip, timeout, truncation, missing executable, drift, issue-delta, contradiction, review-independence, outcome, and late-terminal tests in `internal/sprint/qa_repair_test.go` and `internal/sprint/qa_test.go`.

- [ ] **Task 9: Complete Manual Cycle, Cleanup, Resume, Recovery, And Proof**
  > Executes: `D-04`, `D-06`, `D-07`, `D-08`; `AC-3`, `AC-6`, `AC-8`, `AC-9`
  - [ ] Compose one visible manual cycle in `internal/sprint/qa_repair.go` from confirmation validation through proposal, apply, actual scope, reverification, cleanup, terminalization, result, and flow publication; enforce exactly one proposal and at most one production apply.
  - [ ] Propagate work cancellation to runtime, approved commands, waits, and process trees, then use a separate frozen cleanup deadline to prove descendants terminated, isolated workspace removed, compensation state known, target current, and owned lease released.
  - [ ] Publish writer-fenced `terminalizing` state while holding the mutation lease, verify owned lease release before terminal result/current-state publication, and make every new mutation reject a current terminalizing barrier.
  - [ ] Implement resume and startup recovery from the latest immutable proven boundary, preserving counters and deadline, reusing valid evidence, never repeating a committed apply, never adopting a dead worker, and never inferring success from process absence, patch presence, or target bytes.
  - [ ] Implement qualifying manual-proof publication only from a real shared durable manual run with production apply, complete ladder, current target, proven cleanup, and `verified` or `verified_with_findings`; expose no fixture, dry-run, isolated-only, hand-authored, or administrative creation path.
  - [ ] Add cancellation timing, cleanup timeout, lock-release uncertainty, terminalizing crash, stale lock/owner, apply-committed resume, proof qualification/rejection, proof publication failure, and recovery idempotency tests.

- [ ] **Task 10: Add Manual App DTOs, Shared Operations, And CLI**
  > Executes: `D-02`, `D-07`, `D-09`; `AC-2`, `AC-3`, `AC-8`, `AC-10`
  - [ ] Add additive `RepairUseCases` and bounded prepare, confirm/start integration, status, packet, cycle page/detail, result, resume, cancel, and recover DTOs in `internal/app/sprint_usecases.go`; keep private records, patch bodies, production contents, prompts, raw payloads, unsafe environment, and unrestricted output out of public results.
  - [ ] Bind opaque cursors to project, sprint, repair run, packet digest, collection kind, and retention boundary; return explicit stale and retention-gap results with canonical next actions.
  - [ ] Add manual `repair prepare`, `start --yes`, `status`, `packet`, `cycles`, `result`, `resume --yes`, `cancel`, and `recover` parsing/rendering in `internal/app/sprint_commands.go`; keep progress on stderr, JSON stdout as one versioned document, and start/resume exit zero only for verified outcomes.
  - [ ] Correlate repair run and durable operation run explicitly, reload canonical status after execution, and preserve semantic outcome separately from operational lifecycle in every app result.
  - [ ] Extend `internal/app/sprint_usecases_test.go`, `internal/app/sprint_commands_test.go`, `internal/app/durable_operations_test.go`, and `internal/testdata/qa-canonical-v1.json` for admission, redaction, cursor bounds, help/flags, explicit confirmation, stale input, replay, cancellation, resume, all outcomes/exits, and canonical projection parity.

- [ ] **Task 11: Add Manual TUI Routes And Guarded Actions**
  > Executes: `D-09`; `AC-2`, `AC-3`, `AC-8`, `AC-10`
  - [ ] Extend `internal/tui/model.go` with typed repair summary, packet, confirmation, cycle, result, resume, cancel, and recovery routes/actions that call app use cases only; keep durable execution in the shared operation runner.
  - [ ] Extend `internal/tui/qa_view.go` with packet identity, warning, mode, confirmation state, target freshness, bounded scope, limits, current cycle/gate, fixed reverification ladder, cleanup, lifecycle, outcome, blocker, next action, and paginated history.
  - [ ] Require a separate second guarded manual confirmation action; route reconnect, replay gaps, cancellation acknowledgement, terminal events, and route revisits through canonical status refresh instead of event-derived outcomes.
  - [ ] Bound and sanitize hostile text and terminal controls, preserve narrow-terminal operation and keyboard access, and keep read-only inspection/recovery available when mutation controls are disabled.
  - [ ] Extend `internal/tui/qa_view_test.go` and `internal/tui/model_test.go` for routes, actions, durable run correlation, confirmation, cancellation, resume, recovery, paging, reconnect, status refresh, hostile text, narrow terminals, and all manual terminal states.

- [ ] **Task 12: Add Manual Browser Resources And No-JavaScript Operation**
  > Executes: `D-09`; `AC-2`, `AC-3`, `AC-8`, `AC-9`, `AC-10`
  - [ ] Add bounded repair status, packet, cycle page/detail, and result query handlers in `internal/web/qa_handlers.go` and routes in `internal/web/handlers.go`; handlers use app DTOs only and semantic terminal resource reads return `200`.
  - [ ] Extend `internal/web/templates/sprint.html` with a separate packet review and explicit `Confirm manual repair` page/form, current authority and scope, ladder, cleanup, outcome, blocker, next action, recovery, and paginated evidence that remains complete without JavaScript.
  - [ ] Extend `internal/web/static/js/operations.js` only for progressive form submission, durable event observation, and canonical refresh; browser disconnect, refresh, history replay, or SSE loss must never confirm, cancel, complete, or fail repair.
  - [ ] Reuse existing guarded operation endpoints, same-origin session, CSRF, authorization, body limits, durable run observation, cancellation, shutdown draining, and replay-gap behavior; do not add a repair-only mutation engine or server-owned product state.
  - [ ] Extend `internal/web/qa_handlers_test.go` and shared web operation tests for methods, authorization, CSRF, confirmation, duplicate/stale posts, escaping, bounded responses, no-JavaScript completion, reconnect, replay gaps, browser disconnect, shutdown cancellation, and all manual lifecycle/outcome combinations.

- [ ] **Task 13: Complete Manual Offline Gates And Operator Documentation**
  > Executes: `D-10`; `AC-1` through `AC-11`
  - [ ] Run focused repair tests after each layer, then the complete implementation-repository unit, integration, failure-injection, cancellation, recovery, migration, hostile-input, and race suites; resolve every failure without weakening assertions or protected scope.
  - [ ] Audit every production mutation path from packet digest through confirmation, writer ownership, isolation, proposal, apply, actual diff, reverification, cleanup, result, and proof; reject alternate CLI, TUI, web, runtime, resume, or recovery paths.
  - [ ] Update `docs/architecture.md`, `docs/cli-reference.md`, `docs/user-guide.md`, `docs/phase3-json-schemas.md`, `docs/recovery.md`, `docs/local-web.md`, and `docs/release-checklist.md` for the complete manual protocol, exact commands and fields, outcomes/exits, strict records, cancellation/shutdown, no-JavaScript operation, proof gate, and escalation.
  - [ ] Verify documentation examples against executable command and JSON fixtures; record the manual proof gate as unavailable until Task 14 produces qualifying retained evidence.
  - [ ] Run Architecture Review and Sprint Review against the manual implementation and stop before any automatic implementation if admission, confirmation, apply, cleanup, durability, or interface parity has an unresolved finding.

- [ ] **Task 14: Produce The Real Manual Four-Interface Proof Gate**
  > Executes: `D-10`; `AC-3`, `AC-6`, `AC-9`, `AC-10`, `AC-11`
  - [ ] Select one current real repair-eligible QA issue and prepare its packet through the shared synchronous operation; inspect the same packet and authority facts through CLI text, CLI JSON, TUI, browser HTML without JavaScript, and enhanced browser JSON/event views.
  - [ ] Confirm and start exactly one manual production mutation through the guarded shared durable operation; prove all interfaces correlate the same packet, repair run, operation run, cycle, scope, gates, cleanup, result, and next action without duplicating mutation.
  - [ ] Exercise the full exact-to-review-delta ladder, canonical reconnect/status reload, production target identity checks, process/workspace/lock cleanup, exactly-one result, and qualifying `manual-repair-proof.json` publication.
  - [ ] Retain run IDs, protocol and proof fingerprints, runtime/provider/model identity, durations/costs, artifact digests, actual changed paths, cleanup facts, interface captures, and review evidence required by Deep Smoke Sprint.
  - [ ] Stop the sprint with automatic mode unavailable if the run does not end `verified` or `verified_with_findings`, cleanup is not proven, proof publication fails, any interface differs semantically, or any manual review finding remains unresolved.

- [ ] **Task 15: Implement Lower-Only Automatic Policy After Manual Proof**
  > Executes: `D-06`, `D-07`, `D-08`, `D-10`; `AC-6`, `AC-7`, `AC-8`, `AC-9`
  - [ ] Begin this task only after Task 14 leaves a current qualifying proof; add no administrative bypass, test-fixture proof, or configuration shortcut if the proof is absent or stale.
  - [ ] Extend `internal/platform/config/qa.go` with safe repair defaults, immutable maxima, complete `ULTRAPLAN_QA_*` environment mappings, workspace precedence, effective-source reporting, and lower-only validation for every AC-7 cycle, mutation, reopening, stagnation, file, byte, patch, wall-time, runtime, turn, command, output, retention, and cleanup limit.
  - [ ] Add exhaustive `internal/platform/config/config_test.go` coverage for defaults, maxima, every field/source, workspace and environment precedence, malformed/zero/negative values, attempted increases, and request-level `max-cycles` lowering only.
  - [ ] Add automatic admission that recomputes named manual-proof fingerprint components immediately before confirmation and mutation and explains missing, stale, mismatched, non-manual, weaker, incomplete-cleanup, or non-qualifying proof.
  - [ ] Loop the existing shared cycle only while confirmation, proof, target, governed inputs, policy, ownership, deadline, cleanup, and every budget remain valid; persist consumed counters and absolute deadline before every next-cycle scheduling decision and restore them unchanged on restart/resume.
  - [ ] Add deterministic progress, stagnation, repeated-patch/target, unchanged issue set, reopening, exhaustion, scope/severity growth, design decision, contradiction, unknown schema, uncertain evidence, unsupported test change, drift, and cleanup stop tests; prove every stop occurs before another mutation and yields the correct closed outcome.

- [ ] **Task 16: Expose Automatic Opt-In, Finish Release Evidence, And Close The Sprint**
  > Executes: `D-06`, `D-09`, `D-10`; `AC-6`, `AC-7`, `AC-10`, `AC-11`
  - [ ] Extend app operations and CLI with explicit automatic mode and separate opt-in, allowing only a lower `--max-cycles`; reject manual confirmation reuse, implicit config enablement, changed proof, stale packet/target, or editable resume authority.
  - [ ] Extend TUI and browser with an always-visible automatic availability section, component-level proof reasons, separate `Prepare automatic repair` and `Confirm automatic repair` actions, consumed/remaining limits, progress fact, stagnation/reopening/repetition state, and canonical status refresh.
  - [ ] Complete all seven documentation outputs and the canonical fixture with automatic proof admission, lower-only limits, bounded stops, resume/restart preservation, outcomes, recovery, interface controls, and release checklist evidence.
  - [ ] Run all focused tests, `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan`; rerun Architecture Review and Sprint Review across the complete manual and automatic paths.
  - [ ] Run one explicit real automatic repair that demonstrates a bounded stop path and proves consumed counters and the original deadline survive restart or resume; retain the operation, repair, cycle, proof, stop, cleanup, and interface-parity evidence.
  - [ ] Complete Deep Smoke Sprint evidence only when manual proof remains current, automatic mutation used the same cycle/apply/reverification/cleanup path, every terminal is truthful, and no high-severity security, persistence, fencing, cleanup, parity, or boundedness finding remains.

## Evidence Checklist

- [ ] Sprint 37 current QA, adjudication, containing smoke, and real-runtime dependency evidence is recorded before implementation.
- [ ] Tests prove packet determinism, complete authority joins, explicit confirmation, scope protection, target identity, writer fencing, and no alternate mutation path.
- [ ] Tests cover every protected file class, path traversal, symlinks, hard links where supported, renames/deletes, descendant effects, actual-scope mismatch, apply failure, rollback failure, and target races.
- [ ] Tests cover fixed reverification order, all six outcomes, every automatic stop, cancellation, shutdown, cleanup uncertainty, resume boundaries, recovery, retention, and exactly-one terminals.
- [ ] CLI text/JSON, TUI, browser HTML/JSON, durable operations, repair records, result, and flow projection agree through the canonical fixture and representative integration tests.
- [ ] Public output is bounded, escaped, terminal-safe, and redacted under hostile issue, path, command, runtime, provider, environment, and evidence input.
- [ ] Runtime and diagnostic evidence exists for one qualifying real manual mutation through all four interfaces.
- [ ] `verification/manual-repair-proof.json` is created only by that qualifying real manual run and its fingerprint components are inspectable.
- [ ] Automatic implementation and exposure began only after the manual proof gate passed.
- [ ] Runtime evidence exists for one proof-gated automatic bounded stop whose counters and deadline survive restart or resume.
- [ ] Documentation updates are complete and checked against executable behavior.
- [ ] Architecture Review, Sprint Review, and Deep Smoke Sprint evidence is complete.
- [ ] Deviations from `reasoning.md` are recorded before implementation continues.

## Verification Commands

| Check | Command | Expected Result |
| --- | --- | --- |
| Sprint domain and persistence | `go test ./internal/sprint -run 'Repair|QA'` | Packet, state, proposal, apply, reverification, outcome, recovery, and QA integration tests pass. |
| Durable app and CLI | `go test ./internal/app -run 'Repair|DurableOperation|SprintQA'` | Accept/confirm/dispatch, DTO, CLI, exit, cancellation, resume, and parity tests pass. |
| Lower-only configuration | `go test ./internal/platform/config -run 'QA|Repair'` | Every repair field, source, default, maximum, invalid value, and lower-only override passes. |
| TUI repair | `go test ./internal/tui -run 'Repair|QA'` | Manual/automatic gating, routes, actions, hostile rendering, reconnect, cancellation, and outcomes pass. |
| Browser repair | `go test ./internal/web -run 'Repair|QA|Operation|Shutdown'` | HTTP security, no-JavaScript flow, bounded resources, reconnect, shutdown, and parity pass. |
| Full offline suite | `go test ./...` | All packages pass without requiring a real provider. |
| Race suite | `go test -race ./...` | No race in confirmation, fencing, apply, terminalization, cancellation, recovery, or observers. |
| CLI build | `go build ./cmd/ultraplan` | The single binary builds successfully. |
| CLI surface | `go run ./cmd/ultraplan sprint --help` | Help documents the repair command family and explicit confirmation behavior. |
| Manual dogfood | `ultraplan sprint ultraplan-go 38-bounded-repair repair prepare --issue "$ISSUE_ID" --json` followed by the confirmed manual start and bounded status queries | One current issue produces one packet, one accepted/confirmed manual mutation, one full ladder, proven cleanup, one qualifying result, and current manual proof. |
| Automatic dogfood | `ultraplan sprint ultraplan-go 38-bounded-repair repair start --run "$REPAIR_RUN_ID" --automatic --max-cycles 2 --yes --json` followed by restart or resume and bounded queries | Current proof admits the run; the shared cycle reaches a bounded stop; consumed counters and deadline do not reset. |

## Risks And Blockers

| Risk / Blocker | Source | Mitigation | Status |
| --- | --- | --- | --- |
| Sprint 37 evidence or containing smoke is not current | `reasoning.md#assumptions` and AC-1 dependencies | Task 1 blocks implementation and names the exact prerequisite recovery action. | open |
| Packet joins inconsistent current records | `reasoning.md#risks` | Validate every parent identity and digest as one snapshot; reject missing or ambiguous facts. | open |
| Accepted operation is stranded before confirmation | D-02 | Split accept from dispatch, terminalize publication failure, and test restart with no worker start. | open |
| Multi-file apply is interrupted | D-03, D-04 | Journal intended/completed operations, retain pre-images, compensate in process, and escalate uncertain crash state without reapply. | open |
| Link, side-effect, or protected-class escape | D-03 | Combine strict patch parsing, path class denial, link checks, pre-images, process isolation, and actual target comparison. | open |
| Writer fence and mutation lease disagree | D-04 | Check both at every applicable boundary and use `terminalizing` to close release/publication races. | open |
| Frozen descriptors cannot cover the full ladder | D-05 | Block packet admission; never permit runtime-selected replacement commands. | open |
| Cleanup timeout is mistaken for cleanup proof | D-08 | Require affirmative process, workspace, compensation, and lock facts; map uncertainty to escalation. | open |
| Manual proof is weak or over-sensitive | D-06 | Version named fingerprint components, test each invalidation, and report component-level mismatch. | open |
| Automatic mode is implemented before manual proof | D-10 | Treat Task 14 as a hard gate; no Task 15 or 16 source change begins without current qualifying retained proof. | open |
| Public adapters or progress events drift from canonical truth | D-09 | Use one app DTO fixture, canonical reload, bounded queries, and semantic parity tests. | open |
| Detailed automatic state exceeds retention bounds | D-06, D-09 | Check per-record/cycle/run bounds before writes, prune only unprotected cycles, and expose retention gaps. | open |

## Review Inputs

Review should use:

- `sprint-index.md`
- `technical-handbook.md`
- `reasoning/api-design.md`
- `reasoning/architecture.md`
- `reasoning/frontend.md`
- `reasoning.md`
- this `plan.md`
- implementation diff
- focused, full, race, and build evidence
- real manual and automatic retained run evidence
- Architecture Review, Sprint Review, and Deep Smoke Sprint protocols from `sprint-index.md`

## Execution Log

| Date / Step | Action | Evidence / Notes |
| --- | --- | --- |
| 2026-08-25 / planning | Materialized the Sprint 38 implementation plan from validated governed inputs. | Plan stage only; no implementation, smoke, review, Git, or downstream artifact mutation performed. |

## Completion Criteria

- [ ] All tasks are complete or an explicit blocker is recorded without weakening the requirements.
- [ ] Every Required Output path in `requirements.md` is implemented and covered by the evidence assigned above.
- [ ] Manual repair is proven first from one current adjudicated issue through packet, confirmation, production mutation, full ladder, cleanup, result, proof, and all four interfaces.
- [ ] Automatic repair is implemented and exposed only after that proof is current, reuses the same protocol, and demonstrates one bounded stop with restart/resume-preserved authority.
- [ ] `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan` pass.
- [ ] Documentation and canonical public schemas match executable behavior.
- [ ] Architecture Review, Sprint Review, and Deep Smoke Sprint evidence satisfy the selected protocols.
- [ ] Evidence satisfies `reasoning.md` without reopening or bypassing D-01 through D-10.
- [ ] `review.md` can evaluate AC-1 through AC-11 without guessing implementation intent.
