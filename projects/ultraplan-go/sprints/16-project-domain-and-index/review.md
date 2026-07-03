# Sprint Review: Project Domain and Project Index

> Project: `ultraplan-go`
> Sprint: `16-project-domain-and-index`
> Reviewed: 2026-06-14
> Status: accepted

## Implementation Summary

Sprint 16 implemented project discovery, project status, project-index catalog parsing, catalog validation, and CLI wiring for:

```text
ultraplan project list
ultraplan project <project> status
ultraplan project <project> validate
```

Implementation stayed in the target implementation directory `/home/antonioborgerees/coding/ultraplan-go`.

## Files Changed

| File | Purpose |
| --- | --- |
| `/home/antonioborgerees/coding/ultraplan-go/internal/project/doc.go` | Package ownership and catalog-only semantics. |
| `/home/antonioborgerees/coding/ultraplan-go/internal/project/domain.go` | Project, catalog, status, and validation domain types. |
| `/home/antonioborgerees/coding/ultraplan-go/internal/project/discovery.go` | Direct `projects/` discovery, safe names, exact/prefix reference resolution. |
| `/home/antonioborgerees/coding/ultraplan-go/internal/project/store_fs.go` | Read-only workspace-safe project file and sprint metadata access. |
| `/home/antonioborgerees/coding/ultraplan-go/internal/project/index.go` | Narrow parser for recognized `project-index.md` catalog tables. |
| `/home/antonioborgerees/coding/ultraplan-go/internal/project/validation.go` | Required file checks, catalog path validation, sorted findings, status health. |
| `/home/antonioborgerees/coding/ultraplan-go/internal/project/service.go` | Project use-case API for list, status, and validate. |
| `/home/antonioborgerees/coding/ultraplan-go/internal/project/project_test.go` | Fixture-first project package tests. |
| `/home/antonioborgerees/coding/ultraplan-go/internal/app/project_commands.go` | Thin CLI wiring, rendering, and exit code mapping. |
| `/home/antonioborgerees/coding/ultraplan-go/internal/app/project_commands_test.go` | Command help, output, validation failure, and reference diagnostics tests. |
| `/home/antonioborgerees/coding/ultraplan-go/internal/app/app.go` | Top-level command registration and help listing. |

## Verification

| Check | Result |
| --- | --- |
| `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go` | Passed. |
| `go test -race ./...` from `/home/antonioborgerees/coding/ultraplan-go` | Passed. |
| `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go` | Passed. |
| Import review | Passed: `internal/project` has no `internal/study`, `internal/sprint`, `internal/platform/runtime`, `agentwrap`, or `opencode` imports. |
| Prohibited package review | Passed: no new `internal/catalog`, `internal/validation`, `internal/reports`, `internal/prompts`, or `internal/planning` package. |

## Architecture Review

| Area | Result | Evidence |
| --- | --- | --- |
| Behaviour | Pass | Main workflow is visible in `/home/antonioborgerees/coding/ultraplan-go/internal/project/service.go:12` and app wiring in `/home/antonioborgerees/coding/ultraplan-go/internal/app/project_commands.go:10`. |
| Architecture fit | Good fit | Product behavior is owned by `internal/project`; app code only parses, delegates, renders, and maps errors in `/home/antonioborgerees/coding/ultraplan-go/internal/app/project_commands.go:32`. |
| Simplicity | More complex, justified | Added focused files rather than a global catalog or validation package; parser is narrow to recognized sections in `/home/antonioborgerees/coding/ultraplan-go/internal/project/index.go:9`. |
| Cohesion | Pass | Discovery, store, parser, validation, and service responsibilities are grouped around project planning roots. |
| Coupling | Acceptable | Project package depends only on standard library and `internal/workspace`; no runtime/study/sprint coupling. |
| Side effects | Pass | Project commands are read-only; store methods only read bounded project-owned paths in `/home/antonioborgerees/coding/ultraplan-go/internal/project/store_fs.go:28`. |

## Contract Conformance

| Contract | Status | Evidence |
| --- | --- | --- |
| Architecture | Satisfied | `internal/project` owns project behavior; no global catalog/validation package; app registration is limited to `/home/antonioborgerees/coding/ultraplan-go/internal/app/app.go:130`. |
| Errors | Satisfied | Project reference errors preserve typed classification via `/home/antonioborgerees/coding/ultraplan-go/internal/project/discovery.go:13`; app maps validation failures to exit code `5` in `/home/antonioborgerees/coding/ultraplan-go/internal/app/project_commands.go:66`. |
| Security | Satisfied | Project discovery and catalog paths use `workspace.ResolveInside`; safe name checks reject path segments in `/home/antonioborgerees/coding/ultraplan-go/internal/project/discovery.go:75`; escaping catalog paths fail closed in `/home/antonioborgerees/coding/ultraplan-go/internal/project/validation.go:45`. |
| Testing | Satisfied | Package tests cover discovery, references, status, parsing, malformed rows, missing paths, escapes, empty docs, and URL-like external paths in `/home/antonioborgerees/coding/ultraplan-go/internal/project/project_test.go:10`. Command tests cover help, output, exit `5`, stderr/stdout separation, and no ANSI output in `/home/antonioborgerees/coding/ultraplan-go/internal/app/project_commands_test.go:9`. |
| Documentation | Satisfied | Package docs define ownership and catalog-only semantics in `/home/antonioborgerees/coding/ultraplan-go/internal/project/doc.go:1`; CLI help is implemented in `/home/antonioborgerees/coding/ultraplan-go/internal/app/project_commands.go:124`. |
| CLI Surface | Satisfied | `project list`, `status`, and `validate` render deterministic text; validation findings render section, entry, path, problem, cause, and suggestion in `/home/antonioborgerees/coding/ultraplan-go/internal/app/project_commands.go:105`. |
| Performance | Satisfied | Discovery reads only direct `projects/` entries in `/home/antonioborgerees/coding/ultraplan-go/internal/project/discovery.go:29`; store reads only `docs/*.md`, direct `sprints/` children, roadmap, and project index. |

## Acceptance Criteria

| Criteria | Status | Evidence |
| --- | --- | --- |
| AC-01 to AC-03 | Satisfied | Direct sorted project discovery, hidden/file filtering, safe names, missing and ambiguous refs tested in `/home/antonioborgerees/coding/ultraplan-go/internal/project/project_test.go:10` and `/home/antonioborgerees/coding/ultraplan-go/internal/app/project_commands_test.go:83`. |
| AC-04 | Satisfied | Status summarizes docs, Markdown docs, roadmap, project index, sprints, sprint dirs, and catalog health in `/home/antonioborgerees/coding/ultraplan-go/internal/project/validation.go:67`. |
| AC-05 to AC-06 | Satisfied | Validate exits `0` on valid fixture and exit `5` on findings; command tests at `/home/antonioborgerees/coding/ultraplan-go/internal/app/project_commands_test.go:30` and `/home/antonioborgerees/coding/ultraplan-go/internal/app/project_commands_test.go:61`. |
| AC-07 to AC-10 | Satisfied | Parser recognizes required catalog sections and validates paths without inferring plans or study semantics in `/home/antonioborgerees/coding/ultraplan-go/internal/project/index.go:17`. |
| AC-11, AC-12, AC-15, AC-18 | Satisfied | Import and package review passed; project commands do not invoke runtime, study execution, network, Git, prompt generation, or sprint flow behavior. |
| AC-13 to AC-14 | Satisfied | Output is deterministic plain text; diagnostics include structured repair fields and are tested at `/home/antonioborgerees/coding/ultraplan-go/internal/app/project_commands_test.go:61`. |
| AC-16 to AC-17 | Satisfied | Fixture and help coverage implemented in package and command tests. |
| AC-19 to AC-21 | Satisfied | Required test, race, and build commands passed. |

## Deviations And Deferrals

No deviations from `reasoning.md` or `plan.md` were required.

Explicit deferral: JSON output for `project status` and `project validate` remains out of Sprint 16 scope. The service result structs remain structured enough to support a future schema without moving business rules into `internal/app`.

## Final Assessment

Accepted. The implementation satisfies the sprint plan, selected contracts, and acceptance criteria with no blockers.
