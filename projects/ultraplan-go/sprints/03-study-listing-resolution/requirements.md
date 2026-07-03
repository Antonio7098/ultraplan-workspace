# Sprint Requirements: Study Domain, Listing, and Resolution

> Project: `ultraplan-go`
> Sprint: `03-study-listing-resolution`
> Purpose: the authoritative, human-readable sprint contract. All other sprint artifacts must satisfy these requirements.

## Sprint Goal

Implement the study domain model, deterministic study/source/dimension discovery, reference resolution, and listing commands needed to inspect existing studies without running analysis or initialization workflows.

## Required Outputs

| Output | Path | Description |
| ---------- | -------- | --------------- |
| Requirements | `projects/ultraplan-go/sprints/03-study-listing-resolution/requirements.md` | Sprint contract defining scope, outputs, acceptance criteria, constraints, dependencies, and review expectations. |
| Study domain model | `/home/antonioborgerees/coding/ultraplan-go/internal/study/domain.go` | Defines `Study`, `Source`, `SourceKind`, and `Dimension` with validation-oriented fields required by PRD/TRD. |
| Study filesystem discovery | `/home/antonioborgerees/coding/ultraplan-go/internal/study/discovery.go` | Discovers studies under workspace `studies/`, dimensions under `dimensions/`, and directory sources under `sources/`. |
| Study reference resolution | `/home/antonioborgerees/coding/ultraplan-go/internal/study/resolve.go` | Resolves studies, sources, and dimensions by exact reference and supported unambiguous prefixes. |
| Study listing service | `/home/antonioborgerees/coding/ultraplan-go/internal/study/service.go` | Provides use-case methods for listing studies and listing one study's sources and dimensions. |
| Study CLI commands | `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_commands.go` | Wires `ultraplan study list` and `ultraplan study <study> list` into the existing CLI app. |
| Study tests | `/home/antonioborgerees/coding/ultraplan-go/internal/study/study_test.go` | Covers domain normalization, discovery ordering, ignored entries, and source/dimension resolution behavior. |
| Study command tests | `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_commands_test.go` | Covers command output and error behavior for both study listing commands. |

## Acceptance Criteria

- [ ] `ultraplan study list` lists studies discovered under the resolved workspace `studies/` directory in deterministic ascending order.
- [ ] `ultraplan study <study> list` lists that study's sources and dimensions in deterministic ascending order.
- [ ] Study discovery ignores hidden entries and non-directory entries at the `studies/` level.
- [ ] Source discovery returns one source for each non-hidden directory directly under `studies/<study>/sources/`.
- [ ] Sprint 3 source discovery does not implement Markdown document sources; top-level `.md` source support remains deferred to Sprint 5.
- [ ] Dimension discovery returns Markdown files directly under `studies/<study>/dimensions/` whose names begin with a numeric dimension prefix.
- [ ] Dimension numbers are normalized to two digits for discovered dimensions.
- [ ] Dimension references resolve by number, slug, full filename, exact reference, or unambiguous prefix.
- [ ] Source references resolve by exact name or unambiguous prefix.
- [ ] Ambiguous or missing study, source, and dimension references return actionable non-zero command errors.
- [ ] Listing commands use the existing workspace discovery and global `--workspace` behavior from prior sprints.
- [ ] Listing output includes source kind for discovered directory sources.
- [ ] Listing commands do not recursively scan source repository contents.
- [ ] `go test ./...` passes from `/home/antonioborgerees/coding/ultraplan-go`.
- [ ] `go build ./cmd/ultraplan` passes from `/home/antonioborgerees/coding/ultraplan-go`.

## Non-Goals

- Do not implement `ultraplan study init <study-init.yml>` or YAML parsing in this sprint.
- Do not implement Markdown document source discovery, frontmatter parsing, or applicability filtering in this sprint.
- Do not implement prompt composition, runtime execution, synthesis, report validation, summary generation, run state, retry, cancellation, or code-reference extraction in this sprint.
- Do not implement target or sprint CLI workflows; those remain outside the current study-side release.
- Do not add real OpenCode or `agentwrap` runtime wiring in this sprint.
- Do not add JSON output shapes for study listing unless the existing CLI output foundation already makes this trivial without expanding scope.

## Constraints

- Study behavior must live in `internal/study`; do not create global `internal/validation`, `internal/reports`, `internal/scheduler`, or similar technical-layer packages for this sprint.
- `internal/study` may depend on `internal/workspace` and platform packages, but platform packages must not import `internal/study`.
- CLI parsing and output wiring must stay in `internal/app`; domain discovery and resolution logic must stay outside CLI command handlers.
- Filesystem-managed paths must remain workspace-relative in user-facing output where practical and must not escape the resolved workspace.
- Discovery must be shallow for listing commands: no recursive traversal of source directories or repository contents.
- Output ordering must be deterministic and covered by tests.
- Normal unit tests must not require OpenCode, network access, cloned repositories, or AI provider credentials.
- Error messages must be actionable and must distinguish missing references from ambiguous references.

## Dependencies

| Prior Sprint / Output | Required For | Notes |
| ------------------------ | ------------------- | --------- |
| Sprint 1: Go Module, CLI Shell, and App Composition | Command registration and buildable CLI entrypoint | Provides `cmd/ultraplan`, `internal/app`, and baseline command/test structure. |
| Sprint 2: Workspace, Config, Logging, and Health Skeleton | Workspace discovery and command context | Listing commands must use the existing workspace resolution and global `--workspace` behavior. |
| `projects/ultraplan-go/docs/PRD.md` | Product listing behavior | Defines study listing, source kind display, stable listing output, and deferred target/sprint scope. |
| `projects/ultraplan-go/docs/TRD.md` | Domain and technical behavior | Defines `Study`, `Source`, `Dimension`, deterministic sorting, and lookup rules. |
| `projects/ultraplan-go/docs/ARCHITECTURE.md` | Package boundaries | Defines module ownership and dependency direction for `study`, `workspace`, and platform packages. |

## Review Expectations

| What | How Verified |
| ------------- | ----------------------- |
| Sprint scope matches roadmap Sprint 3 | Review `roadmap.md` Sprint 3 against implemented files and commands. |
| Required outputs exist | Check each path listed in Required Outputs. |
| Study listing commands work | Run `ultraplan study list` and `ultraplan study <study> list` against a fixture workspace. |
| Discovery is deterministic and shallow | Review tests and implementation; verify no recursive source scans are used for listing. |
| Resolution behavior is correct | Review unit tests for exact, prefix, ambiguous, and missing references. |
| Architecture boundaries are preserved | Review imports to confirm platform packages do not import `internal/study` and CLI handlers do not own domain logic. |
| Deferred scope remains deferred | Confirm no study init, Markdown frontmatter/applicability, runtime execution, target, or sprint workflow implementation is included. |
| Build and test gate passes | Run `go test ./...` and `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`. |
