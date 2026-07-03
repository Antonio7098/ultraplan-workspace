# Deep Smoke: 18-select-stage

## Scope

- Sprint: `18-select-stage`
- Implementation directory: `/home/antonioborgerees/coding/ultraplan-go`
- Smoke harness: `/home/antonioborgerees/coding/ultraplan-go-smoke/` (no sprint-18 suite; manual CLI smoke against real binary)
- Selected smoke level/tests: manual CLI smoke of all Sprint 18 surfaces against the real built binary and a real Planning Phase 2 workspace
- Reason scope is sufficient: The smoke harness has no sprint-18 test suite. All Sprint 18 surfaces were exercised directly against the built binary at `/home/antonioborgerees/coding/ultraplan-go/cmd/ultraplan/ultraplan` using a test workspace seeded from the real `.ultra` project structure.

## Run Evidence

| Run ID | Command | Result | Summary |
|---|---|---|---|
| manual-s18-01 | `ultraplan sprint --help` | PASS (exit 0) | Shows status, validate, prompt, flow commands. No ANSI. |
| manual-s18-02 | `ultraplan sprint <p> <s> validate --help` | PASS (exit 0) | Documents validate sprint-index and technical-handbook. |
| manual-s18-03 | `ultraplan sprint <p> <s> prompt --help` | PASS (exit 0) | Documents prompt sprint-index and technical-handbook. Runtime-free. |
| manual-s18-04 | `ultraplan sprint <p> <s> flow --help` | PASS (exit 0) | Documents --to sprint-index, --dry-run, unsupported-stage behavior. |
| manual-s18-05 | `validate sprint-index` (real workspace) | FAIL (exit 5) | Validation failure: reasoning template path mismatch. |
| manual-s18-06 | `prompt sprint-index` (real workspace) | PASS (exit 0) | Deterministic prompt with all required inputs, catalog entries, workspace-relative paths, no-mutation rules. No ANSI. |
| manual-s18-07 | `flow --to sprint-index --dry-run` | PASS (exit 0) | Non-mutating dry-run. Same prompt as prompt sprint-index. flow-state.json not modified. |
| manual-s18-08 | `flow --to sprint-index` (runtime-backed) | FAIL (exit 6) | "runtime is required for sprint-index flow" — sprint service created without runtime injection. |
| manual-s18-09 | `validate technical-handbook` | FAIL (exit 5) | Evidence report unreadable (workspace missing studies/) + reasoning template path mismatch. |
| manual-s18-10 | `flow --to technical-handbook` | FAIL (exit 5) | Same evidence + template validation failures. |
| manual-s18-11 | `flow --to execute` | PASS (exit 2) | Correctly rejected as unsupported. |
| manual-s18-12 | `flow --to smoke` | PASS (exit 2) | Correctly rejected as unsupported. |
| manual-s18-13 | `flow --to review` | PASS (exit 2) | Correctly rejected as unsupported. |
| manual-s18-14 | `validate execute` | PASS (exit 2) | Correctly rejected as unsupported. |
| manual-s18-15 | `flow --to issues` | PASS (exit 2) | Correctly rejected as unsupported. |
| manual-s18-16 | `validate implementation` | PASS (exit 2) | Correctly rejected as unsupported. |
| manual-s18-17 | `validate sprint-index` (missing sprint) | PASS (exit 5) | Reports "sprint reference not found" with available list. |
| manual-s18-18 | `validate sprint-index` (missing project) | PASS (exit 5) | Reports "project reference not found" with available list. |
| manual-s18-19 | `sprint status` | PASS (exit 0) | Shows all stage statuses correctly. |
| manual-s18-20 | All outputs ANSI check | PASS | No ANSI escape codes in any stdout/stderr. |
| manual-s18-21 | flow-state.json after dry-run | PASS | Dry-run did not modify flow-state.json. |

## Findings

| ID | Severity | Status | Evidence | Required Action |
|---|---|---|---|---|
| F-1 | high | open | `validate sprint-index` exits 5 with: `Selected Reasoning Templates "Architecture" (projects/ultraplan-go/sprints/18-select-stage/reasoning/architecture.md): selected path does not match project index; catalog entry name exists but with a different path`. The sprint-index.md artifact references a sprint-local output path `projects/ultraplan-go/sprints/18-select-stage/reasoning/architecture.md` but the project-index.md catalog entry has path `system/reasoning/architecture_reasoning_template.md`. The name "Architecture" matches but the path does not. This is a real artifact content issue in the sprint-index.md generated during Sprint 18 implementation. | Fix the reasoning template path in sprint-index.md to match the project-index.md catalog entry path, OR accept the sprint-local output path as intentional and update the validator to match by name only when the path is a sprint-local output path. |
| F-2 | high | open | `flow --to sprint-index` (non-dry-run) exits 6 with: `runtime is required for sprint-index flow`. The CLI command at `internal/app/sprint_commands.go:47` creates the sprint service with `sprint.NewService(root.Path)` without calling `WithRuntime()`. The sprint service has a `WithRuntime(rt Runtime)` method but it is never invoked from the CLI. The runtime-backed flow path is documented in the help text but cannot execute. | Wire the runtime into the sprint service for non-dry-run flow commands. The sprint command handler should initialize the OpenCode runtime (like study commands do at `internal/app/study_commands.go:227-240`) and call `service.WithRuntime(rt)` before invoking `FlowSprintIndex`. |
| F-3 | medium | open | `validate technical-handbook` exits 5 with multiple "unreadable selected evidence" findings. The evidence report paths in sprint-index.md use `studies/go-cli-study/reports/final/...` which resolve to `<workspace>/studies/go-cli-study/reports/final/...`. The validator checks file existence via `os.Stat()` and the files don't exist in workspaces that don't have the full `.ultra` directory structure. This is expected for test workspaces but reveals that the validator assumes the full `.ultra` directory layout is present at the workspace root. | The validator should distinguish between "path not found in workspace" (which may be OK for external evidence references) and "selected entry not in catalog" (which is always a validation failure). Or document that the validator expects the full `.ultra` directory structure. |
| F-4 | medium | open | `flow --to technical-handbook` also fails with the same evidence unreadability + reasoning template path mismatch. The flow-state.json shows `technical-handbook: failed` with error "selected evidence validation failed". This is a cascading effect of F-1 and F-3. | Resolves when F-1 and F-3 are fixed. |
| F-5 | low | open | The smoke harness at `/home/antonioborgerees/coding/ultraplan-go-smoke/` has no sprint-18 test suite. The harness workspace factory creates study-oriented workspaces, not Planning Phase 2 workspaces. The harness gates all tests behind an OpenCode availability check. | Add a `sprint-18` test suite to the smoke harness if maintainers want automated real-runtime smoke for select-stage commands. The suite needs a Planning Phase 2 workspace factory. |

## Open Issues

- **F-1 (high)**: Sprint-index.md reasoning template path mismatch — the artifact references a sprint-local output path instead of the catalog entry path from project-index.md. This causes `validate sprint-index` to fail with exit 5.
- **F-2 (high)**: Runtime-backed `flow --to sprint-index` is not wired — the CLI creates the sprint service without injecting a runtime, so non-dry-run flow always fails with exit 6.
- **F-3 (medium)**: Technical-handbook validation fails because evidence report paths assume the full `.ultra` directory structure is present at the workspace root.
- **F-4 (medium)**: Cascading failure from F-1 and F-3 affects `flow --to technical-handbook`.
- **F-5 (low)**: Smoke harness has no sprint-18 test suite.

## Resolved Issues

- None filed or resolved during this investigation.

## Verdict

fail

Rationale: Two high-severity findings prevent Sprint 18 from passing deep smoke:

1. **F-1**: `validate sprint-index` fails with exit 5 because the sprint-index.md artifact contains a reasoning template path that doesn't match the project-index.md catalog. This is a real artifact content issue — the generated sprint-index.md uses a sprint-local output path instead of the catalog entry path.

2. **F-2**: `flow --to sprint-index` (non-dry-run) cannot execute because the CLI doesn't wire the runtime into the sprint service. The help text documents this as a supported command, but it always fails with exit 6.

The remaining Sprint 18 surfaces (help, prompt, dry-run, unsupported-stage rejection, exit codes, ANSI-free output, status, missing-project/sprint errors) all pass correctly.
