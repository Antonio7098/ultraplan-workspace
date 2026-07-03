# Sprint Review: Plan Stage

> Project: `ultraplan-go`
> Sprint: `21-plan-stage`
> Review Date: 2026-06-15
> Status: accepted

## Implementation Summary

Implemented the planning `plan` stage in `/home/antonioborgerees/coding/ultraplan-go`:

- Added `internal/sprint/plan.go` with plan manifest construction, structural `plan.md` validation, reasoning decision extraction, task trace checks, evidence checklist checks, and forbidden deferred-stage behavior checks.
- Extended `internal/sprint/prompts.go` with deterministic `RenderPlanPrompt`, including selected inputs, selected area reasoning paths, output path, required sections, traceability rules, and no-execution/no-mutation rules.
- Extended `internal/sprint/service.go` with `ValidatePlan`, `PromptPlan`, and `FlowPlan`; non-dry-run flow validates prerequisites, invokes only the injected generic runtime seam, validates generated `plan.md`, and updates `flow-state.json` only after validation passes.
- Extended `internal/sprint/flow.go` with `plan` complete/failed stage transitions while keeping supported Planning Phase 2 stages limited to `requirements`, `sprint-index`, `technical-handbook`, `area-reasoning`, `reasoning`, and `plan`.
- Extended `internal/app/sprint_commands.go` so CLI handlers route `validate plan`, `prompt plan`, and `flow --to plan` through sprint services and keep command behavior thin.
- Added focused tests in `internal/sprint/plan_test.go` and extended `internal/app/sprint_commands_test.go` for plan validation, prompt invariants, dry-run non-mutation, fake-runtime flow, and CLI routing.

## Verification

| Check | Result |
| --- | --- |
| `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go` | passed |
| `go test -race ./...` from `/home/antonioborgerees/coding/ultraplan-go` | passed |
| `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go` | passed |
| Import review: `rg -n 'internal/study|internal/(catalog|planning|workflow|stages|validation|reports|prompts)' internal/sprint internal/platform` | no forbidden imports/packages found |
| Platform import review with `go list` | no platform package product-module imports found |

## Architecture Review

| Review Area | Result | Evidence |
| --- | --- | --- |
| Behavior | pass | `internal/sprint/plan.go:46` validates required plan structure and traceability; `internal/sprint/service.go:312` gates flow through runtime, artifact read, validation, and state write. |
| Architecture fit | pass | Plan behavior lives in `internal/sprint`; CLI switch only routes to service methods in `internal/app/sprint_commands.go:66`. |
| Simplicity | pass | No global planning/workflow/validation package was introduced; plan-specific helpers remain local to `internal/sprint/plan.go`. |
| Cohesion | pass | Plan manifest/validation is grouped in `plan.go`; prompt rendering remains in `prompts.go`; flow-state transitions remain in `flow.go`. |
| Coupling | pass | Runtime use remains injected through `sprint.Runtime` and `Service.WithRuntime`; no direct OpenCode, shell, Git, provider, or agentwrap adapter calls were added to `internal/sprint`. |
| State and side effects | pass | Dry-run returns prompt content without mutation; non-dry-run updates `flow-state.json` only after generated `plan.md` validates. |

Architecture protocol note: the selected protocol references subagent-style distributed review. Available subagent tooling in this session is constrained to explicit user-requested delegation, so this review applies the protocol directly and does not include separate delegated reports.

## Sprint Review

| Acceptance Area | Result | Evidence |
| --- | --- | --- |
| `validate plan` is local and runtime-free | pass | `internal/app/sprint_commands.go:81`; command tests cover validate/prompt paths with no runtime dependency. |
| `prompt plan` renders deterministic preview and does not write `plan.md` | pass | `internal/sprint/prompts.go:120`; `internal/sprint/plan_test.go:34` checks prompt invariants and no local root leakage. |
| `flow --to plan --dry-run` is non-mutating | pass | `internal/sprint/plan_test.go:47` asserts dry-run result and no `flow-state.json` write. |
| Runtime-backed plan flow uses generic runtime and validates output | pass | `internal/sprint/service.go:344` calls injected runtime; `internal/sprint/service.go:350` reads `plan.md`; `internal/sprint/service.go:356` validates before completion. |
| Flow-state remains limited to supported Phase 2 stages | pass | `internal/sprint/flow.go:83` adds only `plan` completion; `internal/sprint/flow.go:138` rejects unsupported targets. |
| CLI output/help covers plan stage and unsupported stages | pass | `internal/app/sprint_commands.go:66`; command tests cover validate, prompt, dry-run, non-dry-run, help surface, and unsupported flow target. |
| Offline verification | pass | Required test, race, and build commands passed. |

## Deviations

| Deviation | Reason | Risk |
| --- | --- | --- |
| Review protocol was applied directly instead of spawning review subagents. | Tooling policy permits subagents only when explicitly requested by the user. | Low; review evidence is still recorded against selected protocol questions and verification commands. |
| Real OpenCode/runtime smoke was not run. | Sprint explicitly requires fake-first offline verification and treats deep smoke as optional future evidence. | Accepted deferral; fake-runtime tests cover UltraPlan behavior around runtime requests, artifact validation, and state transitions. |

## Risks And Follow-Ups

| Item | Status | Notes |
| --- | --- | --- |
| Selected Architecture area reasoning artifact is absent in current sprint artifacts. | mitigated | Implemented plan flow requires selected area reasoning before runtime-backed plan generation; tests cover prerequisite failure. |
| Plan validation may still allow low-quality prose that satisfies structure. | accepted | Validator intentionally avoids semantic proof or exact prose linting; human review remains part of sprint process. |
| Existing uncommitted reasoning-stage changes were present before this sprint execution. | noted | Implementation preserved existing dirty files and built on current seams without reverting unrelated work. |

## Final Assessment

Accepted. The implementation satisfies Sprint 21 plan-stage requirements with deterministic offline verification, no deferred-stage expansion, and module-owned plan behavior in `internal/sprint`.
