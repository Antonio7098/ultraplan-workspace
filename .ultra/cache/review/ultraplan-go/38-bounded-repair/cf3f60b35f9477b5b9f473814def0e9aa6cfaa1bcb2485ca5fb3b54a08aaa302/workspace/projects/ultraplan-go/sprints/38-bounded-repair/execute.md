# Sprint 38 Execution

Status: `deferred at manual proof gate`

Date: 2026-08-26

## Implemented

- Added a strict sprint-owned repair domain with deterministic packets, single-use confirmations, lower-only budgets, protected-path policy, fixed gate order, closed outcomes, cleanup facts, apply journals, manual-proof qualification, and bounded flow projection.
- Split durable acceptance from dispatch. Repair preparation is synchronous and runtime-free; repair start publishes product confirmation after acceptance and before dispatch. Confirmation failure terminalizes the operation and starts no child.
- Added an isolated manual proposal path and a product-owned apply boundary with preimage checks, private retained preimages, per-operation journal updates, atomic replacements, in-process compensation, hard-link rejection, target/scope checks, and conservative recovery.
- Added ordered fail-closed reverification. Frozen executable checks run with bounded argv, workdir, environment names, timeout, output, cleanup, and target identity checks. Wider gates are skipped after the first non-pass. The final product-owned gate runs containing smoke against the repaired target. Conformance Review runs once before repair admission.
- Added bounded CLI, TUI, and server-rendered browser repair surfaces over shared app operations and DTOs. Automatic preparation and execution remain unavailable.
- Migrated flow state to schema v3, including database-backed v2 reads, and preserved repair summaries across unrelated writes.
- Updated architecture, CLI, user, schema, recovery, local-web, and release documentation.

## Verification

- `go test ./internal/sprint -run 'Repair|QA'`: pass
- `go test ./internal/app -run 'Repair|DurableOperation|SprintQA'`: pass
- `go test ./internal/platform/config -run 'QA|Repair'`: pass
- `go test ./internal/tui -run 'Repair|QA'`: pass
- `go test ./internal/web -run 'Repair|QA|Operation|Shutdown'`: pass
- `go test ./... -count=1`: pass
- `go test -race ./... -count=1`: pass
- `go vet ./...`: pass
- `go build ./cmd/ultraplan`: pass
- `git diff --check`: pass
- Live `ultraplan serve` check: repair page returned HTTP 200; clean shutdown recorded.
- Impeccable detector: ran in degraded regex mode because parser modules were absent; findings were legacy design-system advisories across the shared stylesheet.
- Desktop/mobile screenshot attempt: blocked because local Playwright browser binaries are not installed.

## Deferred Work And Blockers

- Sprint 37 has no current evidence-producing QA attempt, admissible repair issue, current containing smoke, or retained real-runtime evidence. The user approved continuing implementation, but this does not satisfy admission or proof.
- Repaired-target containing smoke is implemented without publishing or rerunning Conformance Review.
- Resume now uses the runtime-free durable recovery boundary and cannot replay a proposal or production apply. Cursor-bound paginated cycle/detail DTOs, exhaustive apply/recovery failure injection, and full adapter lifecycle matrices remain open.
- No real four-interface manual production repair, `manual-repair-proof.json`, Architecture Review, Sprint Review, or Deep Smoke Sprint evidence exists.
- Automatic prepare and start are available through the shared operation path, but admission still requires a current qualifying manual proof and separate explicit opt-in. The real proof and automatic dogfood remain blocked by Sprint 37 evidence and are not counted as passing dogfood evidence.

## Key Implementation Paths

- `internal/sprint/qa_repair.go`
- `internal/sprint/qa_repair_state.go`
- `internal/sprint/qa_types.go`
- `internal/app/durable_operations.go`
- `internal/app/operation_runner.go`
- `internal/app/sprint_commands.go`
- `internal/app/sprint_usecases.go`
- `internal/tui/qa_view.go`
- `internal/web/qa_handlers.go`
- `internal/web/templates/sprint.html`
