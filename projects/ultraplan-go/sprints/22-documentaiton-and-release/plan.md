# Sprint Plan: Planning Documentation and Release Gate

> Project: `ultraplan-go`
> Sprint: `22-documentaiton-and-release`
> Source: `roadmap.md`
> **Inputs Used:** `projects/ultraplan-go/roadmap.md`, `projects/ultraplan-go/docs/*.md`, `~/coding/ultraplan-go`

This plan executes Sprint 22 from the roadmap. It keeps scope to documentation, migration guidance, release gates, and evidence for the planning-side release through `plan.md`.

## Reasoning Source

- **Sprint Reasoning:** `roadmap.md#sprint-22-planning-documentation-and-release-gate`
- **Sprint Index:** `none`
- **Technical Handbook:** `none`
- **Area Reasoning:** `none`

## Sprint Status

- **Status:** `completed`
- **Owner:** `Codex`
- **Start Date:** `2026-06-15`
- **Completion Date:** `2026-06-15`

## Decisions To Execute

| Decision | Source Section | Execution Implication |
| --- | --- | --- |
| Document planning-side release through `plan.md`. | `roadmap.md#sprint-22-planning-documentation-and-release-gate` | Update source docs so `project` and `sprint` planning commands are public and clearly scoped. |
| Keep implementation execution out of scope. | `roadmap.md#planning-artifact-chain` | Docs must state that sprint flow stops at `plan.md` and excludes implementation, smoke, review, issues, Git mutation, and hosted UI. |
| Release gates are offline tests and build, with gated real-runtime planning smoke only when available. | `roadmap.md#sprint-22-planning-documentation-and-release-gate` | Run `go test ./...` and `go build ./cmd/ultraplan`; document planning smoke as optional and gated. |

## Requirements / Contracts To Satisfy

| Contract / Requirement ID | Required Behavior | Evidence Planned |
| --- | --- | --- |
| Planning CLI docs | CLI reference covers `project` and planning-stage `sprint` commands. | `docs/cli-reference.md` and help output comparison. |
| Planning recovery docs | Recovery docs explain failed planning stages and flow-state handling. | `docs/recovery.md`. |
| Prototype migration notes | Docs explain migration from `cli` artifacts to the Go planning flow. | `docs/migration-from-ultra-cli.md`. |
| Release gate | Offline tests and build pass. | `go test ./...`; `go build ./cmd/ultraplan`. |
| Gated planning smoke | Real-runtime planning smoke is documented and skipped unless prerequisites exist. | `docs/planning-smoke.md` and release checklist updates. |

## Tasks

- [x] **Task 1: Update Planning Documentation**
  > Executes: `Planning CLI docs`, `Planning recovery docs`
  - [x] Update README and user guide to include project and sprint planning scope.
  - [x] Update CLI reference with exact `project` and `sprint` command surfaces.
  - [x] Update recovery runbook for failed sprint stages and `flow-state.json`.

- [x] **Task 2: Add Migration And Planning Smoke Docs**
  > Executes: `Prototype migration notes`, `Gated planning smoke`
  - [x] Add `cli` migration notes for project/sprint artifacts.
  - [x] Add gated planning-runtime smoke procedure.
  - [x] Update release checklist to include planning documentation and smoke gates.

- [x] **Task 3: Verify Release Gate**
  > Executes: `Release gate`
  - [x] Run offline tests.
  - [x] Build the CLI.
  - [x] Record any skipped real-runtime smoke prerequisites.

## Evidence Checklist

- [x] Tests prove the required behavior.
- [x] Runtime or diagnostic evidence exists where required.
- [x] Documentation updates are complete where required.
- [x] Deviations from `roadmap.md` are recorded before implementation continues.
- [x] Required review protocols have evidence.

## Verification Commands

| Check | Command | Expected Result |
| --- | --- | --- |
| Help surface | `go run ./cmd/ultraplan --help && go run ./cmd/ultraplan project --help && go run ./cmd/ultraplan sprint --help` | Commands match documented surfaces. |
| Offline tests | `go test ./...` | All tests pass. |
| Build | `go build ./cmd/ultraplan` | CLI builds successfully. |

## Risks And Blockers

| Risk / Blocker | Source | Mitigation | Status |
| --- | --- | --- | --- |
| Real runtime planning smoke may be unavailable locally. | `roadmap.md#sprint-22-planning-documentation-and-release-gate` | Document gated prerequisites and explicit skip evidence. | `mitigated` |
| Docs may overstate Phase 2 scope. | `roadmap.md#planning-scope-principle` | State that planning stops at `plan.md` and does not execute implementation, review, issues, or Git mutation. | `mitigated` |

## Review Inputs

Review should use:

- `roadmap.md`
- `projects/ultraplan-go/docs/*.md`
- this `plan.md`
- implementation diff
- verification evidence

## Execution Log

| Date / Step | Action | Evidence / Notes |
| --- | --- | --- |
| `2026-06-15` | Created Sprint 22 plan from roadmap and docs. | `projects/ultraplan-go/sprints/22-documentaiton-and-release/plan.md` |
| `2026-06-15` | Updated planning documentation and release docs. | `README.md`, `docs/cli-reference.md`, `docs/recovery.md`, `docs/release-checklist.md`, `docs/user-guide.md`, `docs/migration-from-ultra-cli.md`, `docs/planning-smoke.md`, `docs/opencode-smoke.md`, `PRD.md`, `TRD.md` |
| `2026-06-15` | Verified help surfaces, tests, and build. | `go run ./cmd/ultraplan --help`; `go run ./cmd/ultraplan project --help`; `go run ./cmd/ultraplan sprint --help`; `go test ./...`; `go build ./cmd/ultraplan` |

## Completion Criteria

- [x] All tasks are complete or explicitly deferred.
- [x] Verification commands were run or deferrals are documented.
- [x] Evidence satisfies Sprint 22 expectations from `roadmap.md`.
- [x] `review.md` can evaluate conformance without guessing intent.
