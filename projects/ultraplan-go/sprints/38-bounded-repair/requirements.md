# Sprint 38: Bounded Manual and Automatic Repair

## Sprint Goal

Add a governed repair phase that turns one current, adjudicated QA issue into a frozen repair packet, requires explicit human confirmation, permits one bounded production mutation, and progressively re-verifies the result. Prove that manual path end to end before exposing an optional automatic mode. Automatic repair may repeat the same protocol only within hard, lower-only limits and must stop with a truthful terminal outcome when progress, scope, evidence, cleanup, or product intent becomes uncertain.

## Required Outputs

All implementation paths in this section are relative to the `ultraplan-go` implementation repository root. Runtime artifact paths are relative to the selected sprint directory.

### Product semantics and persistence

- `internal/sprint/qa_repair.go`: new `VerificationPhase` repair service for admission, packet freezing, confirmation validation, isolated proposal work, production mutation, progressive reverification, manual and automatic stop rules, recovery, and terminal outcome derivation.
- `internal/sprint/qa_repair_test.go`: table, integration, race, failure-injection, cancellation, recovery, drift, cleanup, and boundedness tests for the repair service.
- `internal/sprint/qa_types.go`: repair packet, confirmation, cycle, scope, progress, freshness, blocker, outcome, and lower-only budget types. Repair remains part of `VerificationPhase`, not `PlanningStage`.
- `internal/sprint/qa_state.go`: strict private persistence, digest-bound references, writer fencing, immutable packet publication, atomic current-state publication, rollback behavior, retention, and bounded `flow-state.json` projection for repair.
- `internal/sprint/qa_state_test.go`: schema, permissions, symlink, containment, digest, stale-writer, partial-write, rollback, retention, and summary-projection coverage for repair records.
- `internal/sprint/qa.go`: QA-to-repair admission and progressive reverification integration without granting investigators mutation or verdict authority.
- `internal/sprint/qa_test.go`: QA admission, issue identity, exact reproducer, containing checks, and post-repair assessment tests.

The repair service must publish the following exact artifact set for each accepted run:

- `verification/repair-state.json`: current repair phase, freshness, active or terminal run correlation, bounded counters, outcome, blocker, next action, and digest-bound pointers only.
- `verification/attempts/<qa-attempt-id>/repairs/<repair-run-id>/issue-packet.json`: immutable frozen issue packet.
- `verification/attempts/<qa-attempt-id>/repairs/<repair-run-id>/confirmation.json`: immutable confirmation record bound to the packet digest, target identity, mode, limits, confirmer, and timestamp.
- `verification/attempts/<qa-attempt-id>/repairs/<repair-run-id>/cycles/<cycle-number>/proposal.patch`: bounded proposed patch retained as evidence before production mutation.
- `verification/attempts/<qa-attempt-id>/repairs/<repair-run-id>/cycles/<cycle-number>/scope.json`: before and after target identities, actual changed paths, changed-byte count, and scope-enforcement result.
- `verification/attempts/<qa-attempt-id>/repairs/<repair-run-id>/cycles/<cycle-number>/reverification.json`: ordered reproducer, shard, theory, boundary, containing-suite, and Conformance Review delta results.
- `verification/attempts/<qa-attempt-id>/repairs/<repair-run-id>/cycles/<cycle-number>/cleanup.json`: workspace and process-tree cleanup result.
- `verification/attempts/<qa-attempt-id>/repairs/<repair-run-id>/result.json`: immutable terminal outcome and reason, completed limits, unresolved issues, evidence pointers, and recovery or escalation action.
- `verification/manual-repair-proof.json`: bounded pointer to the latest current `verified` or `verified_with_findings` manual run that exercised the same repair protocol and policy. Automatic admission must reject a missing, stale, mismatched, or weaker proof.

### Shared application and interfaces

- `internal/app/sprint_usecases.go`: adapter-independent prepare, confirm, start, status, resume, cancel, recover, and bounded-query DTOs for repair.
- `internal/app/operations.go`: repair operation kinds, canonical requests, confirmation facts, mutation class, governed inputs, and stable operation results.
- `internal/app/operation_runner.go`: one shared durable repair runner used by CLI, TUI, and HTTP, with writer-fence ownership and canonical cancellation.
- `internal/app/sprint_commands.go`: text and JSON CLI for packet preview, manual confirmation, optional automatic mode, status, resume, cancel, and recovery.
- `internal/app/sprint_commands_test.go`: CLI parsing, confirmation, JSON schema, exit code, cancellation, resume, stale-input, and outcome coverage.
- `internal/app/sprint_usecases_test.go`: app projection, admission, redaction, and interface parity coverage.
- `internal/app/durable_operations_test.go`: acceptance-before-goroutine, single-owner, restart, cancellation, stale-writer, and exactly-one-terminal-outcome repair tests.
- `internal/app/operation_runner.go` and `internal/app/operations.go` must preserve the existing server shutdown contract: graceful server shutdown cancels repair, waits for bounded cleanup, and records `cancelled`, `interrupted`, or `cleanup_uncertain` rather than success.
- `internal/platform/config/qa.go`: lower-only repair limits under QA policy, with safe defaults and immutable product maxima.
- `internal/platform/config/config_test.go`: precedence, environment names, invalid values, defaults, maxima, and lower-only override tests for every repair limit.
- `internal/tui/qa_view.go`: packet, confirmation warning, mode, cycle progress, scope facts, reverification ladder, outcome, blocker, and next-action display using app DTOs only.
- `internal/tui/qa_view_test.go`: TUI parity, bounded hostile-text rendering, manual-first gating, confirmation, cancellation, and terminal-outcome tests.
- `internal/tui/model.go`: repair routes and typed actions that call shared app use cases.
- `internal/tui/model_test.go`: route, action, durable-run, reconnect, and status-refresh coverage.
- `internal/web/qa_handlers.go`: bounded repair packet, status, cycle, and result query handlers using app use cases only.
- `internal/web/qa_handlers_test.go`: HTTP method, CSRF, authorization, confirmation, escaping, bounded response, reconnect, and stale-request tests.
- `internal/web/handlers.go`: HTML and JSON routes for repair resources and guarded mutations.
- `internal/web/templates/sprint.html`: server-rendered manual and automatic repair controls, explicit confirmation page content, progress, evidence, outcomes, and recovery guidance that remain usable without JavaScript.
- `internal/web/static/js/operations.js`: progressive enhancement only. Browser disconnect, refresh, or SSE loss must not cancel or complete repair.
- `internal/testdata/qa-canonical-v1.json`: additive canonical app projection fixture for repair facts without exposing full prompts, secrets, raw provider payloads, or unbounded output.

### Documentation and release evidence

- `docs/architecture.md`: repair ownership, authority, state placement, mutation isolation, run control, and review separation.
- `docs/cli-reference.md`: exact repair commands, JSON fields, outcomes, exit behavior, limits, and examples.
- `docs/user-guide.md`: packet review, manual confirmation, automatic opt-in, status, cancellation, resume, recovery, and escalation workflow.
- `docs/phase3-json-schemas.md`: strict private repair records and additive public CLI, app, and HTTP projections.
- `docs/recovery.md`: interruption, stale ownership, target drift, partial publication, cleanup uncertainty, and restart reconciliation.
- `docs/local-web.md`: browser routes, confirmation behavior, reconnect behavior, and no-JavaScript operation.
- `docs/release-checklist.md`: manual proof gate, automatic admission gate, parity checks, real-runtime dogfood, race checks, and adversarial scope tests.

The implementation and review must treat these current plans as governed implementation inputs:

| Input | Responsibility in this sprint |
| --- | --- |
| `docs/plans/integrated-roadmap.md` | Sequence Sprint 38 after evidence-producing QA and before dogfooding; keep content/retrieval work deferred. |
| `docs/plans/post-execution-qa-and-repair-loop.md` | Define frozen issue packets, confirmation, scope enforcement, progressive reverification, manual-first delivery, automatic bounds, and terminal outcomes. |
| `docs/plans/server-shutdown-run-cancellation-contract.md` | Preserve server-owned cancellation, process-tree cleanup, bounded shutdown, terminal arbitration, and restart reconciliation for repair runs. |

## Acceptance Criteria

### AC-1: repair admission freezes one authoritative issue packet

- A repair can be prepared only from the current QA attempt's adjudicated `QAIssue`. The issue must be repair-eligible and backed by accepted, current, repeatable or deterministically sufficient failing evidence.
- Preparation rejects stale review, stale QA, missing containing smoke, failed or blocked assessment, rejected evidence, target identity drift, unsupported isolation, an active conflicting mutation, or a missing exact reproducer.
- `issue-packet.json` contains the issue and root-cause-group IDs; source QA attempt, adjudication, evidence, plan, map, policy, implementation, review, smoke, and target fingerprints; issue class, severity, location, claim, promotion reason, and evidence references; the violated requirement and acceptance-criterion references; exact reproducer descriptor and expected failing condition; affected primary, boundary, and follow-up shard IDs; linked theory IDs; allowed production paths; forbidden paths; containing check descriptors; repair acceptance criteria; frozen budgets; and packet digest.
- The packet permits exactly one issue and a finite normalized path set. Empty, absolute, parent-traversing, symlinked, repository-control, workspace-state, generated-evidence, test, requirement, acceptance, or governed-input paths are forbidden unless the adjudicated issue identifies a production test fixture and the packet explicitly classifies it as mutable production data.
- Packet creation is runtime-free, idempotent for identical inputs, and cannot mutate the implementation, QA evidence, tests, requirements, acceptance criteria, or state outside the private repair artifact set.

### AC-2: confirmation is explicit and bound to what will run

- Every production mutation requires a new confirmation record after the user sees the issue identity, severity, violated expectation, exact reproducer, allowed and forbidden paths, target identity, mode, maximum cycles, maximum changed files and bytes, containing checks, cleanup limit, and stop conditions.
- CLI confirmation uses an explicit non-interactive `--yes` after packet preview. TUI and browser use a separate guarded confirmation action. Merely opening a page, reconnecting, replaying an event, or preparing a packet is not confirmation.
- Confirmation authority is a single-use digest over the packet, target identity, canonical request, selected mode, effective limits, and durable operation acceptance. Any changed field, stale packet, changed target, changed governed input, or replay rejects before mutation.
- Automatic mode requires a second, explicit automatic-mode choice. A manual confirmation cannot authorize automatic cycles.

### AC-3: manual repair is the first and smallest mutable path

- Manual mode handles one frozen issue and permits one proposal, one production patch application, and one progressive reverification sequence. It never starts a second mutation cycle automatically.
- The runtime may inspect only packet-approved context, must propose a bounded patch in an isolated copy, and receives no authority to edit the production checkout directly.
- Product code validates patch syntax, target identity, path containment, file count, changed bytes, disallowed file classes, symlinks, and current writer ownership before applying the patch through a product-owned mutation boundary.
- The production mutation fails closed unless actual changed paths are a subset of `allowed_paths` and disjoint from `forbidden_paths`. The service compares the full target identity before proposal, before apply, after apply, and after reverification.
- A violating patch is retained as rejected evidence when safe, never applied, and ends `blocked` or `escalated` with a concrete reason.
- Manual repair is proven through CLI text and JSON, TUI, and browser using the same app use case, durable operation, sprint service, packet schema, confirmation check, scope guard, reverification engine, cancellation path, and outcome derivation.

### AC-4: scope enforcement protects the source of truth

- Repair cannot modify `requirements.md`, `code-context.md`, `sprint-index.md`, `technical-handbook.md`, `reasoning.md`, `plan.md`, `execute.md`, `review.md`, `smoke.md`, `qa.md`, `flow-state.json`, QA or repair evidence, implementation plans, workspace configuration, Git metadata, hooks, ignored-file policy, or test and acceptance assets to manufacture a pass.
- The protection applies to direct writes, renames, deletes, symlink or hard-link replacement, generated files, subprocess descendants, formatter side effects, and cleanup actions.
- Test changes, expectation changes, acceptance-criteria changes, evidence deletion, baseline updates, snapshot re-recording, check removal, command substitution, and scope expansion require a separate future sprint or a newly adjudicated and explicitly confirmed packet. Repair must stop rather than self-authorize them.
- Repair runs with one writable owner. Every state and production publication revalidates the durable writer fence. A stale or cancelled owner cannot publish progress, mutate production, release another owner's lock, or overwrite a terminal result.

### AC-5: reverification follows a fixed widening order

- After mutation, the service runs these gates in order and records each result: exact reproducer; affected primary shard checks; linked theory confirmation and refutation conditions; affected boundary, neighboring, and parent-linked follow-up shard checks; containing approved QA suites; then containing smoke against the repaired target.
- A gate does not run when a narrower required gate fails, blocks, becomes stale, or reports cleanup uncertainty. The result names the skipped gates and the action needed to continue.
- The same immutable command descriptors, timeouts, output limits, environment allowlists, evidence rules, and target identity checks used for admission govern reverification. The repair runtime cannot choose replacement commands or reinterpret outcomes.
- Conformance Review runs once before repair admission and remains independent and authoritative for its verdict. Repair does not run a second Conformance Review or edit `review.md`.
- Reverification must detect and report new failures, reopened issues, issue-set changes, severity changes, scope growth, target drift, and contradictions. It must not discard a finding because the original issue no longer reproduces.

### AC-6: automatic repair is disabled until manual proof is current

- Automatic repair is off by default and unavailable when `verification/manual-repair-proof.json` is absent, stale, points to a non-manual run, uses another repair schema or policy, lacks verified cleanup, or does not end `verified` or `verified_with_findings` after the complete reverification ladder.
- The proof is created only after a real manual production mutation reaches a qualifying terminal outcome through shared durable run control. Unit fixtures, dry runs, isolated-copy-only runs, hand-written files, and administrative overrides cannot create it.
- A code, schema, maximum-limit, isolation-capability, policy, governed-input, or target identity change invalidates the proof. Automatic admission explains which fingerprint changed and requires another manual proof.
- Enabling automatic mode requires explicit user confirmation per run. Configuration cannot enable it silently or bypass the current proof gate.

### AC-7: automatic work is demonstrably bounded

- Automatic mode reuses the manual packet, confirmation, proposal, scope, apply, reverification, cleanup, persistence, and outcome code. Sprint 38 permits one automatic cycle. It does not add a second mutation engine or retry loop.
- Product defaults and immutable maxima exist for total cycles, mutation cycles, reopenings of the original issue, unchanged issue-set cycles, changed files per cycle and run, changed bytes per cycle and run, generated patch bytes, wall-clock time, runtime attempts, model turns, command count, command time, output bytes, retained cycles, and cleanup time. Workspace and environment configuration may only lower these values.
- The automatic cycle must show measurable progress by removing the exact failure, reducing the current admitted issue set, or lowering the highest admitted severity. Rewording, reordered evidence, or a model explanation is not progress.
- Automatic repair ends `stalled` when its single cycle makes no measurable progress. Later sprints may add retries only with a separate governed requirement.
- Automatic repair ends `escalated` before another mutation on any unconfirmed path or issue scope growth, new issue class, severity growth, design or product decision, contradictory acceptance criteria, target or governed-input drift, unsupported test change, uncertain evidence, unknown schema, or cleanup uncertainty.
- Limits are checked before scheduling work and again before every proposal, apply, command, publication, and resume. A restart or resume restores consumed counters and deadlines; it never resets them.

### AC-8: terminal outcomes are closed and truthful

- `verified`: the exact issue no longer reproduces, every required progressive gate passes, cleanup is proven, the target is current, no admitted issue remains in the packet scope, and the Conformance Review delta passes.
- `verified_with_findings`: the exact issue no longer reproduces, required gates complete, cleanup is proven, and the Conformance Review delta or containing QA reports only current non-blocking findings that do not increase severity or repair scope.
- `failed`: admissible deterministic evidence shows the issue still reproduces or a required check fails after the permitted manual mutation or automatic work, with no missing prerequisite preventing that conclusion.
- `blocked`: repair cannot start or continue because a prerequisite, runtime, check executable, permission capability, lock, persistence operation, or required evidence source is unavailable. No semantic success or failure is inferred.
- `escalated`: repair encounters scope growth, severity growth, target or governed-input drift, a design decision, contradictory expectations, unsafe proposed mutation, unsupported schema, or cleanup uncertainty that requires human adjudication.
- `stalled`: automatic mode reaches its stagnation, repeated-patch, cycle, or reopening bound without enough evidence for `verified`, `verified_with_findings`, `failed`, or `escalated`. Manual mode never emits `stalled`.
- Cancellation and shutdown use durable operation states `cancelled`, `interrupted`, or `cleanup_uncertain`. The repair result maps cleanup uncertainty to `escalated` and never maps cancellation to a verified outcome.
- Exactly one authoritative terminal result is committed. Late runtime completion, duplicate confirmation, repeated cancellation, stale writers, SSE replay, or restart reconciliation cannot replace it.

### AC-9: resume, cancellation, cleanup, and recovery preserve authority

- Resume requires the same packet, confirmation mode, writer-fence lineage, target and governed-input identities, policy, consumed limits, and latest proven cycle boundary. It reuses valid immutable evidence and never repeats a committed production apply.
- Cancellation reaches runtime calls, waits, approved commands, subprocess trees, and isolated workspaces. Locks release only after work stops and cleanup is proven or uncertainty is recorded.
- Graceful server shutdown rejects new mutations, requests canonical cancellation for every server-owned repair, waits for bounded cleanup, persists a truthful durable state, closes event streams, and then exits. Browser disconnect does not cancel repair.
- Startup recovery reconciles active records without a live owner, stale locks, partial publications, leftover process trees, and isolated workspaces. It preserves valid evidence but never infers success from process absence or an existing patch.
- Cleanup failure or inability to prove process-tree, workspace, or lock cleanup prevents verification and records an actionable escalation.

### AC-10: CLI, JSON, TUI, and browser agree

- The CLI supports packet preview, manual start with explicit confirmation, automatic start with explicit automatic confirmation, status, bounded cycle and result queries, resume, cancel, and recover. `--json` returns a stable versioned envelope and uses non-zero exit status for `failed`, `blocked`, `escalated`, `stalled`, cancellation, interruption, or cleanup uncertainty.
- TUI and browser expose the same operations and facts from shared app DTOs. They do not read private verification files, derive outcomes from progress events, or call `internal/sprint` directly.
- All interfaces show packet identity and digest, issue and severity, mode, confirmation state, target freshness, allowed-path summary, consumed and remaining limits, current cycle and gate, actual changed paths, cleanup state, terminal outcome, blocker, and next action. Large path, evidence, and event collections are paginated or bounded.
- Hostile retained text is escaped and redacted. Public responses exclude secrets, full prompts, unsafe environment values, raw provider payloads, unrestricted command output, and production file contents.
- Refresh, reconnect, replay gaps, process restart, and concurrent observers reload canonical product facts and durable operation facts. Progress events remain non-authoritative.

### AC-11: verification and dogfood gates pass

- `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan` pass in the implementation repository.
- Tests cover all six repair outcomes, every automatic stop condition, every protected file class, path traversal, symlink and hard-link attacks where supported, changed-target races, stale writer races, duplicate confirmation, partial writes, rollback failure, cancellation timing, shutdown timing, resume boundaries, and interface parity.
- A real-runtime dogfood run proves one manual repair from current adjudicated evidence through confirmation, production mutation, the full reverification ladder, cleanup, durable terminal state, and all four interfaces.
- Only after that manual proof is current may a real-runtime automatic dogfood run start. The automatic run must demonstrate at least one bounded stop path and show that consumed limits survive restart or resume.
- Fake-runtime tests, hand-authored state, dry runs, and isolated proposals do not satisfy either dogfood gate.

## Non-Goals

- General issue tracking, backlog management, or repair of issues not produced by current QA adjudication.
- Unbounded autonomous repair, silent automatic enablement, cross-project repair, detached workers, remote collaboration, or cloud/Aren execution.
- Repair that changes tests, requirements, acceptance criteria, evidence, review or smoke verdicts, QA adjudication, configuration, Git history, or implementation plans to obtain a pass.
- Parallel production repair writers, multi-issue manual repair, speculative scope expansion, broad refactoring, dependency upgrades, migrations, or product redesign unless separately adjudicated and confirmed in a future sprint.
- Replacing Conformance Review, smoke, QA evidence, or human product decisions with model confidence.
- Content identity, retrieval-ready content, content-addressed storage, alternate persistence, provenance graphs, knowledge graphs, indexing, search, migration, or compatibility work for those deferred systems.
- Richer production hardening beyond the manual proof and bounded automatic capability. Sprint 39 owns broad dogfooding, adversarial hardening, compatibility proof, and production-readiness closure.

## Constraints

- `internal/sprint` owns repair semantics, packet validation, scope enforcement, progressive reverification, outcome derivation, and detailed repair persistence. `internal/platform/*` may provide reusable process, copy, identity, lock, and filesystem mechanisms but cannot decide repair outcomes.
- `VerificationPhase` keeps `conformance-review`, `qa`, and `repair` separate from `PlanningStage`. Existing planning stage order, review identity, `review.md`, and `smoke.md` remain unchanged.
- Detailed packet, proposal, cycle, evidence, cleanup, and result records live under `verification/`. `flow-state.json` stores only a bounded repair summary with phase, freshness, verdict or outcome, counts, next action, and digest-bound pointers.
- One durable run has one authoritative owner and writer fence. Confirmation and durable acceptance occur before any execution goroutine or runtime starts.
- Product code, not a model, owns approved paths, forbidden paths, command descriptors, isolation, patch application, counters, deadlines, evidence admission, scope comparison, terminal arbitration, and success derivation.
- Production mutation occurs only after isolated proposal generation and product validation. The runtime never receives direct write authority to the production checkout.
- Default-deny permissions, local path containment, private `0700` directories, private `0600` files, no symlink following, bounded reads and writes, atomic publication, immutable evidence, and content digests apply to every repair record.
- All numeric and duration controls have safe product defaults and finite immutable maxima. Workspace files and `ULTRAPLAN_QA_*` environment variables may lower but never raise them. Request flags cannot override maxima.
- Full prompts, secrets, unsafe environment values, unbounded output, provider payloads, and production content are not persisted or emitted through public interfaces.
- No Git commit, branch, reset, checkout, stash, clean, push, hook edit, or index mutation is part of repair. Repair reports pre-existing unrelated changes and refuses overlap; it does not revert or absorb them.
- The smallest correct patch is preferred. A successful reproducer does not justify unrelated cleanup or refactoring.

## Dependencies

- Sprint 36 must continue to provide deterministic QA maps, shard ownership, theory links, bounded investigation, durable run control, and strict detailed state under `verification/`.
- Sprint 37 must provide current isolated evidence, adjudication, promoted repair-eligible issues, containing smoke integration, and the shared CLI, JSON, TUI, and browser QA path.
- Before Sprint 38 implementation starts, Sprint 37's current mapped external smoke suite must complete with `pass` or `pass_with_open_issues`, and its real-runtime evidence-producing QA gate must be current. The manually reconciled `pass_with_findings` Conformance Review alone does not satisfy this dependency.
- Before automatic mode is implemented or exposed, the manual packet schema, confirmation flow, scope enforcement, production apply boundary, reverification ladder, cleanup, durable terminal state, and interface parity must pass tests and a real-runtime dogfood run.
- The implementation repository must provide enforceable protected-root denial, isolated writable copies, process-tree cleanup, bounded command execution, target identity, mutation leases, durable operations, writer fencing, and startup reconciliation. Missing capability blocks repair rather than weakening the contract.
- Current implementation-plan dependencies are `docs/plans/integrated-roadmap.md`, `docs/plans/post-execution-qa-and-repair-loop.md`, and `docs/plans/server-shutdown-run-cancellation-contract.md` as recorded under Required Outputs.

## Review Expectations

- Review the manual path first. Do not approve automatic repair while any manual admission, confirmation, scope, mutation, reverification, cleanup, durability, or parity criterion is incomplete.
- Trace every mutation from the confirmed packet digest through durable acceptance, writer ownership, isolated proposal, scope validation, production apply, actual-diff comparison, reverification, cleanup, and terminal publication. There must be no alternate adapter or runtime path around that chain.
- Treat attempts to edit tests, requirements, acceptance criteria, evidence, state, configuration, Git control files, or unconfirmed paths as security failures, not ordinary validation errors.
- Verify each terminal outcome from product-owned facts and prove exactly-one-terminal-result behavior under cancellation, timeout, shutdown, target drift, persistence failure, late runtime completion, and restart.
- Compare CLI text, CLI JSON, TUI, browser HTML, browser JSON, durable run records, `verification/repair-state.json`, `result.json`, and `flow-state.json` for the same run. Differences must be presentation-only.
- Confirm public output bounds and redaction with hostile issue text, paths, command output, model output, and provider errors.
- Run focused package tests during implementation, then `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan`. Record the real manual dogfood evidence before reviewing automatic admission.
- Reject any change that lets repair declare global success, weaken or replace QA evidence, overwrite Conformance Review, reset consumed limits on resume, treat cleanup uncertainty as success, or expand authority without a new adjudicated packet and confirmation.
