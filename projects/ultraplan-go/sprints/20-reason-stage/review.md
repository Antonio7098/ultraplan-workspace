# Sprint Review: Reason Stage

> Project: `ultraplan-go`
> Sprint: `20-reason-stage`
> Date: 2026-06-15
> Status: complete

## Implementation Summary

Implemented reason-stage support through `reasoning` in `/home/antonioborgerees/coding/ultraplan-go`.

- Added sprint-owned selected reasoning template manifest, area/final reasoning validators, selected-context path checks, and deterministic reasoning prompt rendering.
- Added service use cases for `validate area-reasoning`, `validate reasoning`, `prompt area-reasoning`, `prompt reasoning`, and `flow --to area-reasoning|reasoning`.
- Extended flow-state transitions for selected area completion/skipping and final reasoning completion without marking `plan` complete.
- Wired CLI parsing/help/output to delegate to `internal/sprint` and preserve existing exit-code behavior.
- Added focused sprint and command tests for manifests, prompts, validation, dry-run no mutation, fake-runtime generation, and CLI behavior.

## Files Changed

- `internal/sprint/reasoning.go`
- `internal/sprint/reasoning_test.go`
- `internal/sprint/service.go`
- `internal/sprint/prompts.go`
- `internal/sprint/flow.go`
- `internal/sprint/store_fs.go`
- `internal/sprint/handbook.go`
- `internal/app/sprint_commands.go`
- `internal/app/sprint_commands_test.go`

## Verification

| Check | Result |
| --- | --- |
| `go test ./internal/sprint ./internal/app` | passed |
| `go test ./...` | passed |
| `go test -race ./...` | passed |
| `go build ./cmd/ultraplan` | passed |

## Architecture Review

- `internal/sprint` owns reason-stage manifest, validation, prompt rendering, runtime-backed flow, and flow-state transitions.
- `internal/app/sprint_commands.go` remains thin: it parses arguments, delegates to `sprint.Service`, renders output, and maps exit codes.
- No global planning, workflow, validation, prompt, or catalog package was added.
- Inspection found no `internal/sprint` import of `internal/study`.
- Inspection found no `internal/platform/*` imports of product modules.
- Non-dry-run flow uses the existing generic runtime boundary via `sprint.Runtime`; `internal/sprint` does not directly call shell, Git, OpenCode, provider APIs, or agentwrap adapters.

## Security Review

- Selected reasoning template paths are validated against `project-index.md`, checked as workspace-relative, readable, non-empty files.
- Selected area output paths must stay under `projects/<project>/sprints/<sprint>/reasoning/*.md`.
- Prompt previews use workspace-relative paths and include no-mutation rules for non-output artifacts, implementation files, workspace config, source repositories, and Git state.
- Selected-context reference validation checks selected paths/names/templates/protocols and rejects obvious unselected path references.

## Residual Risks

- Validators are structural and path/name based; they do not prove full semantic truth of generated prose.
- Default verification uses fake runtime only. Real OpenCode smoke remains out of scope unless selected by a future deep-smoke review.
- Prompt builders are intentionally verbose because the selected context and no-mutation rules are explicit.

## Conclusion

Sprint 20 acceptance criteria are met. Required verification passed, and the implementation stays within the selected Phase 2 planning scope through `reasoning` without adding plan execution, smoke, review automation, issues, or Git mutation behavior.
