# Sprint Plan: Phase 3 Documentation, Hardening, and Release

> Project: `ultraplan-go`
> Sprint: `29-phase-3-documentation-hardening-release`
> Status: `planned`
> Created: `2026-07-28`
> Planning mode: manual, plan-only
> Inputs used: `projects/ultraplan-go/roadmap.md`, `projects/ultraplan-go/project-index.md`, and `projects/ultraplan-go/sprints/28-review-to-smoke-flow/`

## Goal

Stabilize the delivered `execute -> review -> smoke` workflow as a supported CLI and TUI release, resolve the material findings from Sprint 28, replace stale documentation, dogfood the complete Phase 3 flow, and satisfy the Phase 3 release gate before Sprint 30 begins.

## Scope

This sprint is a release-hardening sprint. It may change the UltraPlan Go implementation, tests, embedded documentation/templates, and canonical implementation-repository documentation required to close Sprint 28 findings.

It does not add a new workflow stage, begin the browser UI, create a second verification state hierarchy, introduce general-purpose issue tracking, copy raw smoke evidence into the planning workspace, or automatically mutate Git state.

## Execution Order

Security, state integrity, and workflow correctness block the release and come first. Documentation and interface polish follow the stabilized behavior. Release evidence runs last.

## Tasks

- [x] **Task 1: Close blocker and high-severity safety findings**

  - Redact runtime `UserDetail`, `DebugDetail`, and `ResponseBody` values before they cross the runtime boundary; add regression cases containing secrets.
  - Replace substring-derived operation/error classification with typed errors that preserve `errors.Is`/`errors.As` behavior across smoke mapping.
  - Reconcile or remove dead stale-attempt handling so abandoned running attempts cannot remain indefinitely authoritative.
  - Give retryability one explicit bounded owner, or remove the orphan `Retryable` signal.
  - Tighten smoke environment forwarding and redact sensitive configuration/extra-argument values in effective diagnostics.
  - Add a top-level command/request correlation identifier and propagate it through runtime requests, events, results, state, and stable JSON where applicable.

- [x] **Task 2: Repair durable verification state and recovery**

  - Make flow-state reads side-effect free; move migration/reconciliation writes behind an explicit atomic commit path.
  - Enforce checksum/integrity validation for durable verification state and preserve one known-version migration path.
  - Reconcile stale running attempts during status/recovery and expose the resulting terminal state.
  - Define safe ordering and recovery for the `smoke.md`/flow-state split-commit window without discarding a valid artifact.
  - Add idempotency protection for resumed review, smoke, focused reruns, and verification operations.
  - Ensure diagnostic review override cannot bypass a blocked gate or promote an overall assessment.
  - Give incomplete/stale assessment outcomes distinct typed results and stable exit behavior.
  - Add compatibility, corruption, interrupted-write, stale-attempt, duplicate-request, and reconciliation fixtures.

- [x] **Task 3: Complete runtime identity and operational diagnostics**

  - Introduce a canonical prompt reference containing ID, version, owner kind/ID, purpose, and checksum; propagate it through review execution and evidence.
  - Extend attempt summaries with maximum attempts, attempts remaining, issue kind/code, and retryability when retry exists.
  - Persist and render bounded runtime warnings such as dropped events, permission limitations, and state repairs.
  - Expose cleanup outcome as structured attempt diagnostics.
  - Add redacted verbose diagnostics for review, smoke, verify, and flow operations.
  - Keep platform runtime/process packages product-neutral and preserve the read-only review/smoke mutation boundary.

- [x] **Task 4: Finish Phase 3 CLI and TUI behavior**

  - Add worked examples and documented exit behavior to sprint command help.
  - Expand TUI verification coverage for confirmation, diagnostic override rationale, focused reruns, cancellation, evidence links, and recovery.
  - Polish review/smoke dashboard summaries, responsive progress, error/recovery panes, narrow-terminal rendering, and keyboard help.
  - Prove CLI text, stable JSON, status, validation, flow, and TUI use the same typed state, verdict, staleness, diagnostics, and next action.
  - Preserve thin CLI/TUI adapters over shared app use cases; do not invoke CLI handlers or parse terminal output from the TUI.

- [x] **Task 5: Replace stale documentation and publish stable schemas**

  - Update the CLI reference and user guide to describe the delivered review, smoke, verify, status, flow, focused-rerun, override, and recovery behavior.
  - Update the TUI guide, recovery guide, configuration reference, smoke-harness guide, migration guide, and release checklist.
  - Correct documented flags so they match the parser and add examples for the normal and failure/recovery paths.
  - Document stable JSON schemas for review, smoke, verify, and status, including lifecycle state, verdict, assessment, freshness, correlation, diagnostics, artifacts, evidence links, and next action.
  - Document migration from manual `review.md` and legacy `deep-smoke.md` to generated `review.md` and `smoke.md`, including backup/removal and regeneration.
  - Declare the planning workspace as the authoritative PRD/TRD/Architecture source and link it from the implementation repository without creating duplicate mirrors.
  - Validate that the authoritative project planning documents are internally consistent and current.

- [x] **Task 6: Expand hardening and parity tests**

  - Add missing reviewer tests for permission enforcement and exhausted structured-output repair.
  - Complete security/redaction, path-containment, bounded-concurrency, cancellation, race, and durable compatibility coverage.
  - Add semantic cross-surface scenarios for pass, findings, fail, blocked, stale, incomplete, cancellation, recovery, override, and focused rerun.
  - Verify required-environment absence remains `blocked`, never `pass`.
  - Verify deliberate contract violations produce a failing `review.md`.
  - Verify deliberate runtime behavior failures produce a failing `smoke.md` linked to external harness evidence.
  - Confirm default tests use fake runtime/harness collaborators and do not contact OpenCode, the real harness, the network, or ambient credentials.

- [x] **Task 7: Run and record the Phase 3 release gate**

  - Run focused package tests during implementation.
  - Run the full release commands:

    ```bash
    go test ./...
    go test -race ./...
    go build ./cmd/ultraplan
    go vet ./...
    git diff --check
    ```

  - Confirm the fake-runtime review suite and fake smoke-harness suite pass.
  - Confirm the entire normal `execute -> review -> smoke` workflow is operable from CLI and TUI.
  - Record command results, dogfood run/evidence IDs, remaining accepted exceptions, and the final release decision.

## Acceptance Criteria

- [ ] No relevant blocker or high-severity Sprint 28 review finding remains open without an explicit accepted release exception.
- [ ] No relevant open smoke-harness issue permits a clean Phase 3 release assessment.
- [ ] Runtime/provider details and configuration values are redacted before reaching user, debug, state, or JSON surfaces.
- [ ] Verification state is integrity-checked, atomically persisted, resumable, idempotent, and recoverable after interruption or split commit.
- [ ] Review/smoke prompt identity, correlation, retry attempts, warnings, and cleanup outcomes are inspectable and bounded.
- [ ] CLI, JSON, status, validation, flow, and TUI agree on lifecycle state, verdict, freshness, assessment, diagnostics, evidence, and next action.
- [ ] Phase 3 public documentation matches delivered commands and flags, includes recovery/migration guidance, and publishes stable JSON schemas.
- [ ] Deliberate review and runtime failures produce truthful failing artifacts; missing prerequisites produce `blocked`, not `pass`.
- [ ] Cancellation leaves CLI and TUI state recoverable and does not corrupt the last valid Markdown artifact.
- [ ] Fake-runtime and fake-harness suites pass by default; gated real review and smoke evidence is recorded or truthfully blocked.
- [ ] The full release command set passes.
- [ ] Phase 3 is ready for release, and Sprint 30 may begin only after the final evidence is current or remaining exceptions are explicitly recorded.

## Planned Evidence

- Focused and full Go test results, including race and vet.
- CLI help/JSON schema and TUI parity tests.
- Durable-state migration, corruption, interruption, idempotency, and recovery fixtures.
- Security/redaction and path-containment regression tests.
- Dogfood `review.md` and `smoke.md` artifacts with external harness run/evidence references.
- Gated real-runtime and real-harness result, including truthful blocked evidence when prerequisites are unavailable.
- Final Phase 3 release checklist and release decision.
