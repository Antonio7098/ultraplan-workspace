# Sprint Review: Distill Stage

> Project: `ultraplan-go`
> Sprint: `19-distill-stage`
> Reviewed At: 2026-06-14T19:58:27Z
> Status: complete

## Implementation Summary

Implemented planning distill-stage support for `technical-handbook` in the target implementation directory `/home/antonioborgerees/coding/ultraplan-go`.

Key changes:

- Added `internal/sprint/handbook.go` for selected evidence manifest construction, workspace-safe evidence checks, handbook section validation, selected evidence traces, unselected evidence detection, and decision-language rejection.
- Extended `internal/sprint/prompts.go` with deterministic `technical-handbook` prompt rendering that includes required paths, selected evidence manifest, required sections, no-decision rules, selected-only rules, and no-mutation rules.
- Extended `internal/sprint/service.go` and `internal/sprint/flow.go` with validate, prompt, dry-run, and runtime-backed flow for `technical-handbook`.
- Extended `internal/app/sprint_commands.go` so `validate technical-handbook`, `prompt technical-handbook`, and `flow --to technical-handbook` are supported while unsupported stages still fail as usage errors.
- Added focused coverage in `internal/sprint/handbook_test.go` and extended `internal/app/sprint_commands_test.go`.

## Verification

All required checks were run from `/home/antonioborgerees/coding/ultraplan-go`.

| Check | Result |
| --- | --- |
| `go test ./...` | passed |
| `go test -race ./...` | passed |
| `go build ./cmd/ultraplan` | passed |
| `go run ./cmd/ultraplan --workspace /home/antonioborgerees/coding/ultra/.ultra sprint --help` | printed help with `technical-handbook` validate, prompt, and flow support |
| Unsupported-stage live command | blocked by workspace config mismatch: `.ultra` has no `ultraplan.yml`; command-test coverage verifies usage exit and message |

## Architecture Review Evidence

- `internal/sprint` owns distill behavior: selected evidence manifest, handbook validation, prompt rendering, and flow transitions live in sprint package files.
- `internal/app` remains command wiring: it parses stage arguments, calls sprint service methods, renders text output, and maps exit statuses.
- `internal/platform/*` remains product-agnostic.
- Import-boundary scan returned no matches:

```bash
rg 'internal/study|internal/(planning|workflow|validation|prompts|stages|reports)' internal/sprint internal/platform -n
```

- No global planning, workflow, validation, prompt, stages, reports, or catalog package was introduced.
- Runtime-backed handbook flow calls the existing generic `internal/platform/runtime` boundary through the sprint `Runtime` interface.

## Acceptance Conclusions

- Selected evidence is loaded only from reports selected in `sprint-index.md` and matched to `project-index.md`.
- Evidence paths are normalized, workspace-safe, readable, deterministic, and validated before runtime calls.
- Handbook validation rejects missing/empty content, placeholders, missing required sections, missing selected evidence traces, unselected report references, and final decision language.
- Prompt preview is runtime-free, does not write the handbook artifact, and uses workspace-relative paths.
- Dry-run flow is runtime-free and does not write `technical-handbook.md` or complete flow state.
- Runtime-backed flow validates prerequisites before runtime, validates the generated handbook after runtime, and updates `flow-state.json` only after success.
- Flow-state support remains limited to `requirements`, `sprint-index`, `technical-handbook`, `area-reasoning`, `reasoning`, and `plan`.
- Unsupported stages such as `implementation`, `execute`, `smoke`, `review`, and `issues` remain rejected.

## Deviations And Deferrals

- No real OpenCode/provider smoke was run. This matches the sprint non-goal and verification plan.
- Live unsupported-stage CLI inspection against the `.ultra` planning workspace could not run because the Go CLI workspace discovery expects `ultraplan.yml`, while this planning workspace contains `config.json`. The behavior is covered by `internal/app/sprint_commands_test.go`.

## Residual Risk

Markdown validation remains heuristic by design. The tests cover the required sprint cases, but future planning stages may need to refine section parsing if handbook prose evolves substantially.
