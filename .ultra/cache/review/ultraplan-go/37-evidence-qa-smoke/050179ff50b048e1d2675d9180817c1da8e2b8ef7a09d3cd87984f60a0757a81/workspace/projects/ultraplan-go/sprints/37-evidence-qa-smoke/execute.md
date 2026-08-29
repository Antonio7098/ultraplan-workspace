# Execute Summary

Plan: `projects/ultraplan-go/sprints/37-evidence-qa-smoke/plan.md`
Run state: `projects/ultraplan-go/sprints/37-evidence-qa-smoke/.run-state.json`
Implementation worktree: `/home/antonioborgerees/coding/ultraplan/.ultraplan-go-ultraplan-worktrees/ultraplan-go/37-evidence-qa-smoke`
Implementation branch: `ultraplan/ultraplan-go/37-evidence-qa-smoke`

## Task Counts

- pending: 0
- running: 0
- complete: 11
- deferred: 1
- failed: 0
- cancelled: 0

## Implementation Evidence

- Added product-neutral bounded tree identity, copy, comparison, native protected-root denial, contained process execution, and cleanup facts in `internal/platform/process`.
- Added lower-only limits for copied trees, generated checks and patches, retained evidence, and issues.
- Added frozen QA v2 plans, evidence, patches, adjudication, issue groups, assessment, deterministic IDs, and strict validation while retaining state v1 reads.
- Added sequential isolated evidence execution with target identity checks, original-path leakage rejection, approved-path enforcement, bounded output, cancellation, and cleanup.
- Added pure global adjudication, fixed three-result failed-shard evaluation, product-owned assessment, and deterministic `qa.md` rendering.
- Added fenced evidence publication. Immutable records precede canonical state, and a later state or flow failure restores the prior `qa.md`, state pointer, and flow projection.
- Routed `qa --suite smoke` through the canonical smoke executor. Resume and shard focus are rejected for the smoke suite.
- Added additive app and JSON v1 projections, focused evidence/adjudication/issues/assessment/smoke queries, CLI controls, TUI/browser facts, and guarded browser smoke-suite starts.
- Updated architecture, CLI, user, browser, recovery, schema, and release documentation.

## Verification Evidence

All implementation commands ran in the recorded worktree on 2026-08-25.

- Generic isolation focused tests passed, including native denial of writes to the original source.
- Investigation, state, adjudication, assessment, atomic rollback, app, CLI, TUI, browser, and smoke parity focused tests passed.
- `go test ./internal/platform/process -run 'Test(Isolation|DirectRunner)' -count=1` passed.
- `go test ./internal/sprint -run 'TestQA(Investigation|EvidencePlan|Writable|Permission|Cancellation|Cleanup)' -count=1` passed.
- `go test ./internal/sprint -run 'TestQA(Store|State|Adjudication|Assessment|Atomic|Recovery)' -count=1` passed.
- `go test ./internal/sprint -run 'Test(Smoke|QA).*Parity|TestQA.*Smoke' -count=1` passed.
- `go test ./internal/app -run 'Test(QA|SprintQA|DurableQA)' -count=1` passed.
- `go test ./internal/tui ./internal/web -run 'TestQA|TestBrowserQA|TestOperation.*QA' -count=1` passed.
- `go test ./...` passed.
- `go test -race ./...` passed.
- `go vet ./...` passed.
- `go build ./cmd/ultraplan` passed.
- `node --check internal/web/static/js/operations.js` passed.
- `git diff --check` passed.
- `TestRealSmokeHarness` reported `blocked: set ULTRAPLAN_REAL_SMOKE=1 to opt into the cataloged external harness` and skipped. It is not counted as passing dogfood evidence.

## Final Execution Disposition

Tasks 1 through 11 are implemented and pass the offline release gates. Task 12 is deferred at the post-execution boundary: Architecture Review, Sprint Review, Deep Smoke, and real external-harness dogfood were not launched by execute, and the real smoke opt-in is unavailable. No review or smoke verdict is claimed. The gated release criterion remains blocked until those downstream stages produce current evidence.
