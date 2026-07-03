# Sprint Review: Select Stage

> Project: `ultraplan-go`
> Sprint: `18-select-stage`
> Date: 2026-06-14
> Status: complete

## Implementation Summary

Implemented the Sprint 18 select-stage surface in `/home/antonioborgerees/coding/ultraplan-go`.

- Added `internal/sprint` sprint-index parsing, catalog subset validation, prompt preview rendering, service methods, generic-runtime flow execution, safe failure state recording, and success transitions.
- Added CLI support for `ultraplan sprint <project> <sprint> validate sprint-index`, `prompt sprint-index`, and `flow --to sprint-index [--dry-run]`.
- Kept sprint behavior in `internal/sprint`; `internal/app` parses arguments, delegates to the sprint service, renders output, and maps exit codes.
- Preserved Phase 2 scope: no implementation execution, smoke execution, review automation, issue tracking, Git mutation, global workflow engine, or study-semantic reuse was introduced.

## Files Changed

- `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/domain.go`
- `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/index.go`
- `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/validation.go`
- `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/prompts.go`
- `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/flow.go`
- `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/store_fs.go`
- `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/service.go`
- `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/sprint_index_test.go`
- `/home/antonioborgerees/coding/ultraplan-go/internal/app/app.go`
- `/home/antonioborgerees/coding/ultraplan-go/internal/app/app_test.go`
- `/home/antonioborgerees/coding/ultraplan-go/internal/app/sprint_commands.go`
- `/home/antonioborgerees/coding/ultraplan-go/internal/app/sprint_commands_test.go`
- `projects/ultraplan-go/sprints/18-select-stage/plan.md`
- `projects/ultraplan-go/sprints/18-select-stage/review.md`
- `projects/ultraplan-go/sprints/18-select-stage/.run-state.json`

## Verification

All required Go checks passed from `/home/antonioborgerees/coding/ultraplan-go`.

| Check | Result |
| --- | --- |
| `go test ./internal/sprint ./internal/app` | passed |
| `go test ./...` | passed |
| `go test -race ./...` | passed |
| `go build ./cmd/ultraplan` | passed |

Manual CLI smoke against `/home/antonioborgerees/coding/ultra/.ultra` did not run as a product workspace because that directory lacks `ultraplan.yml`; command behavior is covered by workspace-fixture command tests.

## Architecture Review

- `internal/sprint` owns select-stage parsing, validation, prompt rendering, flow rules, and flow-state updates.
- `internal/app` remains a thin command adapter.
- `internal/sprint` imports `internal/project`, `internal/workspace`, and the generic `internal/platform/runtime` request type only; it does not import `internal/study`.
- No forbidden global `catalog`, `planning`, `workflow`, `validation`, `reports`, `prompts`, or `stages` package was introduced.
- Platform runtime remains product-agnostic.

## Sprint Review

- `validate sprint-index` is runtime-free and validates editable Markdown plus project-index subset membership.
- `prompt sprint-index` is runtime-free, uses workspace-relative paths, includes catalog entries and no-mutation rules, and does not write artifacts.
- `flow --to sprint-index --dry-run` is non-mutating.
- Non-dry-run flow validates requirements, invokes only the injected generic runtime boundary, requires and validates `sprint-index.md`, and updates `flow-state.json` only after validation succeeds.
- Unsupported future/execution stages are rejected.
- CLI output is plain text, stream-separated in command tests, and no ANSI output is emitted.

## Deviations And Residual Risk

- The implementation uses focused invariant tests rather than exhaustive golden prompt text, matching the sprint reasoning.
- Real OpenCode/runtime smoke was not run; Sprint 18 default verification is fake-first and offline.
- The manual CLI check against the `.ultra` planning directory is blocked by workspace shape, not by select-stage behavior.

## Conclusion

Sprint 18 is complete and ready for follow-on planning-stage sprints.
