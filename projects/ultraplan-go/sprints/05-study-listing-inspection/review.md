# Sprint Review: 05 Study Listing Inspection

> Project: `ultraplan-go`
> Sprint: `05-study-listing-inspection`
> Review Date: 2026-05-31
> Status: complete

## Summary

Sprint 5 is implemented in `/home/antonioborgerees/coding/ultraplan-go`.

The implementation adds direct-child Markdown source discovery, leading YAML frontmatter parsing and stripping, normalized `applicable_dimensions`, a study-owned applicability helper, and deterministic listing output for directory and Markdown sources. Runtime execution, prompt composition, report validation, synthesis, run-loop/status behavior, target workflows, and sprint workflows remain out of scope.

## Files Reviewed

- `/home/antonioborgerees/coding/ultraplan-go/internal/study/domain.go`
- `/home/antonioborgerees/coding/ultraplan-go/internal/study/discovery.go`
- `/home/antonioborgerees/coding/ultraplan-go/internal/study/markdown.go`
- `/home/antonioborgerees/coding/ultraplan-go/internal/study/applicability.go`
- `/home/antonioborgerees/coding/ultraplan-go/internal/study/study_test.go`
- `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_commands.go`
- `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_commands_test.go`

## Architecture Review

- Behavior remains owned by `internal/study`, with the existing `internal/app` command seam limited to rendering command output.
- No global `validation`, `scheduler`, `reports`, or `prompts` packages were introduced.
- `go list ./...` loads the package graph; platform packages do not import `internal/study`.
- The change adds focused helper files rather than a broader abstraction or runtime-facing surface.
- Failure paths for filesystem, YAML/frontmatter, and applicability metadata preserve cause chains with `%w` and include source path context.

Verdict: good fit.

## Sprint Review

- AC-01 through AC-03: covered by discovery and command tests for directories, top-level Markdown, hidden entries, unrelated files, and nested Markdown.
- AC-04 through AC-10: covered by frontmatter parse/strip and applicability normalization tests, including invalid value diagnostics.
- AC-11 through AC-14: covered by `GetApplicableSources` tests for directories, unfiltered Markdown, matching filtered Markdown, and excluded non-matches.
- AC-15: command output remains plain and deterministic; directory rows remain `name directory`, Markdown rows render `markdown all` or normalized filters.
- AC-16: test fixtures cover the required valid, absent, malformed, ignored, nested, mixed ordering, and error cases.
- AC-17: `go test ./...` passes.

## Verification

| Command | Result |
| --- | --- |
| `go test ./internal/study ./internal/app` | pass |
| `go test ./...` | pass |
| `go list ./...` | pass |

## Deferrals And Notes

- `system/protocols/sprint-review-protocol.md` is selected by `sprint-index.md` but is absent from the workspace. This review records sprint acceptance against the sprint plan, requirements, and available architecture review protocol.
- No implementation blockers remain.
