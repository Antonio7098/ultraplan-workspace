# Architecture Reasoning: Study Initialization From YAML

> **Inputs Used:** `projects/ultraplan-go/sprints/04-study-init-yaml/technical-handbook.md`, `projects/ultraplan-go/sprints/04-study-init-yaml/requirements.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `system/reasoning/architecture_reasoning_template.md`

## 0. Feature Summary

### Feature Name

Study initialization from explicit YAML.

### User/Product Goal

Users can run `ultraplan study init <study-init.yml>` to create a deterministic, reviewable study directory from an explicit YAML definition. The command must generate normalized YAML, a README, editable dimension files, sources and report folders, and optional shallow source clones without implementing deferred runtime-assisted completion or study execution.

### Current Task

Implement the architecture for parsing, validating, planning, rendering, writing, cloning, CLI wiring, and tests for study initialization.

### Non-Goals

- Runtime-assisted source or dimension completion, suggestion prompts, suggestion cache artifacts, and `--no-assist`.
- Source URL verification, replacement source suggestions, repair prompts, OpenCode execution, or agentwrap runtime use.
- Markdown document discovery, frontmatter parsing, `applicable_dimensions`, applicability filtering, prompt composition, analysis runs, synthesis, report validation, summary generation, durable run state, retries, cancellation, diagnostics persistence, code-reference extraction, target workflows, sprint workflows, or new stable public JSON schemas.

## 1. First-Principles Breakdown

### Core Behaviour

This feature receives a YAML study initialization file plus CLI options, validates and normalizes explicit study/source/dimension definitions, builds a filesystem and clone plan, optionally prints the plan, and otherwise creates or overwrites scoped generated artifacts under the resolved workspace or selected output root.

### Inputs

- `study-init.yml` containing `name`, `description`, `repos.count`, `repos.items`, `dimensions.count`, and `dimensions.items`.
- CLI flags: `--dry-run`, `--force`, `--no-clone`, `--output <dir>`, and existing global `--workspace` behaviour.
- Existing workspace discovery/path-safety services from prior sprints.
- Local `git` executable only when cloning URL-backed sources and cloning is not disabled.

### Outputs

- `studies/<study>/study-init.yml`, normalized deterministically.
- `studies/<study>/README.md` describing the initialized study and supported next commands.
- `studies/<study>/dimensions/NN-slug.md` files preserving supplied purpose, steps, citations, and questions.
- `studies/<study>/sources/`, `studies/<study>/reports/source/`, and `studies/<study>/reports/final/` directories.
- Optional non-shell clone executions shaped `git clone --depth 1 <url> <sources/name>`.
- Text output describing success, dry-run plans, validation failures, overwrite protection, or partial clone failure.

### Durable State

- Created: Generated study directory, normalized YAML, README, dimension Markdown files, source/report directories, and cloned source directories.
- Read: Input YAML and workspace configuration/discovery state.
- Updated: Known generated files inside the selected study directory when `--force` is set.
- Deleted: None. `--force` should not remove unknown files or paths outside the selected study directory.

### Ephemeral State

- Parsed initialization request.
- Normalized source and dimension records.
- Validation diagnostics with YAML field paths.
- Explicit plan containing directories, files, and clone actions.
- Clone results and partial-failure aggregation.

### Derived State

- Dimension numbers: source of truth is explicit YAML numbers; computed as two-digit normalized strings; stored in generated YAML, filenames, headings, and Markdown.
- Dimension slugs: source of truth is dimension `name`; computed as kebab-case; stored in filenames and normalized YAML.
- Study output paths: source of truth is workspace root plus study name or explicit output flag; computed during planning; not separately persisted except as generated artifacts.
- Clone destinations: source of truth is normalized source names; computed as `sources/<name>`.

### Side Effects

- Database write: no.
- File write: yes, generated artifacts and directories under a validated output root.
- Network call: possible through `git clone` only; tests must fake this boundary.
- Queue/event emission: no durable event system in this sprint.
- Email/notification: no.
- Logs/metrics/traces: minimal command output only; persistent diagnostics and observability are deferred.

## 2. Existing Architecture Fit

### Contract Applicability And Phase Gate

Current gate: local CLI study initialization with deterministic filesystem artifacts.

Applies now:

- Architecture: required because behaviour must live in `internal/study`, CLI wiring stays in `internal/app`, and dependency direction must remain `internal/app -> internal/study -> internal/workspace/platform` with no platform import of `study`.
- Errors: required for actionable YAML, path, overwrite, and clone diagnostics.
- Security: required for workspace path checks, scoped overwrite behaviour, non-shell clone execution, and secret-free artifacts.
- Testing: required for parser, planning, rendering, dry-run, force, command output, clone fakes, and offline deterministic checks.
- Documentation/CLI Surface: required for generated README, help, flag behaviour, and user-facing command text.

Deferred:

- Runtime/LLM contracts: assisted completion and OpenCode/agentwrap execution are outside this sprint.
- Persistence and migrations: generated artifacts are durable files, but no versioned migration or run-state schema is introduced.
- Observability/performance/workflows: persistent diagnostics, event streams, scheduling, concurrency, retries, status, and large-scale runtime workflows are deferred.

### Where This Behaviour Naturally Belongs

Study initialization belongs in `internal/study` because the Architecture document states that study owns study definitions, sources, dimensions, validation, reports, and filesystem persistence. The CLI command belongs in `internal/app` because it parses command arguments, resolves command context, wires dependencies, and prints output.

### Existing Workflow Affected

Current flow:

1. CLI command construction lives under `internal/app`.
2. Workspace resolution and global `--workspace` behaviour come from prior workspace support.
3. Existing study listing/resolution expects deterministic study directories under `studies/<study>`.

### Proposed New Flow

Proposed flow:

1. `internal/app` registers `ultraplan study init <study-init.yml>` with `--dry-run`, `--force`, `--no-clone`, and `--output` flags.
2. The command resolves the workspace using existing global behaviour and constructs a `study.InitRequest`.
3. `internal/study` parses YAML into explicit source and dimension definitions.
4. `internal/study` validates counts, names, dimension numbers/slugs, required fields, and output path safety.
5. `internal/study` builds one plan containing directories, files, and clone actions.
6. Dry-run prints the same plan without writes or clone execution.
7. Real execution creates directories/files, then runs clone actions through an injected non-shell clone runner unless `--no-clone` is set.
8. `internal/app` maps success, validation errors, overwrite errors, workspace/path errors, and partial clone failure to existing exit-code conventions.

### Does The Current Architecture Still Fit?

[x] Yes - the feature fits cleanly into the existing shape.

### Refactor-Before-Feature Decision

Decision: implement directly with small focused additions inside existing packages.

Reason: The handbook evidence supports thin CLI wiring, product logic outside `RunE`, explicit dependencies, narrow external command seams, and one planning path for dry-run and writes. The Architecture doc already anticipates `internal/study/init.go`; no larger package split is needed.

## 3. Design Options Considered

### Option A: Focused `internal/study` Init Service With Plan-Then-Execute

Description:
Keep parsing, validation, normalization, rendering, planning, and execution coordination in focused files within `internal/study`. Expose a small use-case API, likely `InitStudy(ctx, InitRequest) (InitResult, error)`, and keep helper types private unless command output needs them. Use an injected clone runner at the `git` boundary.

Pros:
- Matches Architecture guidance that study owns study lifecycle behaviour and filesystem persistence.
- Matches handbook evidence for thin CLI handlers, explicit dependencies, narrow volatile boundaries, and dry-run wrappers/planning.
- Makes dry-run and real initialization share validation and planning.
- Keeps tests deterministic with temp directories and clone fakes.

Cons:
- Adds several focused files to one package, so `internal/study` will need naming discipline.
- Requires a clear result shape so CLI output can be useful without exposing unnecessary internals.

Risks:
- The plan type could grow into a general workflow abstraction if future runtime features are pulled into this sprint.

### Option B: Put Init Logic In CLI Handler Then Extract Later

Description:
Implement most parsing, validation, file writes, and clone calls directly in `internal/app` command code, extracting only render helpers.

Pros:
- Lowest initial file count.
- Easy to thread CLI flags directly into behaviour.

Cons:
- Violates sprint constraints that YAML validation, artifact planning, rendering, and filesystem behaviour stay outside CLI handlers.
- Conflicts with handbook evidence warning against business rules in `RunE`.
- Makes command tests depend on private CLI flow and makes future study APIs harder.

Risks:
- Duplicate business rules when future study workflows need the same definitions and validation.
- Harder to preserve existing exit-code conventions cleanly.

### Option C: Create Generic Filesystem/Workflow/Validation Subpackages

Description:
Introduce global or platform packages such as `internal/validation`, `internal/reports`, `internal/filesystem/planner`, or a generic command execution framework for plan, dry-run, writes, and clones.

Pros:
- Could be reusable if many product modules soon need identical planning/execution infrastructure.
- Makes seams visibly explicit.

Cons:
- Directly conflicts with Architecture and sprint constraints against global technical-layer packages for product behaviour.
- The handbook recommends narrow interfaces at volatile boundaries, not abstraction around every helper.
- Adds complexity before there are multiple real consumers.

Risks:
- Fractures study context and creates coupling around generic names rather than product ownership.

### Chosen Option

Chosen option: A.

Reason:
The simplest honest design is a focused `internal/study` initialization use case with plan-then-execute and an injected clone runner. It satisfies sprint requirements while preserving module ownership and dependency direction. Options B and C are rejected because one puts product rules in CLI wiring and the other introduces premature global abstractions.

## 4. Abstraction Check

### Are We Adding A New Abstraction?

[x] Yes - service/component: a study initialization use-case API in `internal/study`.
[x] Yes - data structure/DTO/config object: request, result, and internal plan structs.
[x] Yes - interface/protocol/trait: a narrow clone runner or command runner boundary.

### Why It Is Earned

- It isolates an external dependency: `git clone` is a subprocess and possible network operation.
- It makes testing simpler: normal tests can fake clone execution and assert args without GitHub, network, or real repositories.
- It creates a clear boundary between policy and mechanism: `internal/study` decides clone policy and destination; the runner executes an argument-array command.
- It protects a stable domain concept: study initialization is a named use case with validations and generated artifacts.

### Bad Abstraction Smell Check

- Avoid generic names such as `Manager`, `Processor`, `Common`, or global `validation` packages.
- Avoid an interface that mirrors a whole concrete service. The interface should exist only for clone execution, not for pure parsing or rendering helpers.
- Avoid mode flags that split unrelated behaviour. `DryRun`, `Force`, and `NoClone` are domain-level request options and should feed the same plan rather than separate flows.

## 5. DRY And Duplication Check

### Is Any Duplication Being Introduced?

[x] No duplicated business/domain knowledge should be introduced.
[x] Yes - some duplicated code shape may appear between file rendering tests and command output tests.

### Single Sources Of Truth

- YAML field validation rules belong in `internal/study`, not in `internal/app`.
- Name, slug, and dimension number normalization belong in `internal/study` helpers reused by parser, renderer, and plan builder.
- The plan is the single source of truth for dry-run output and execution.
- Workspace path safety remains owned by `internal/workspace`; `internal/study` consumes exported path APIs rather than reimplementing workspace discovery.

### Acceptable Duplication

Command tests may assert user output while study tests assert structured plans and generated files. That is acceptable because command UX and domain artifact generation change for different reasons.

## 6. Coupling Check

### Dependencies Required

- `internal/workspace`: needed for resolved workspace root and path-safety checks; allowed by Architecture.
- Standard library filesystem/YAML dependencies or existing project YAML package: needed to parse input and write files.
- `os/exec` or an existing platform process helper inside the concrete clone runner: needed only at the external command boundary.
- `internal/app`: depends on `internal/study` for command wiring; `internal/study` must not depend on `internal/app`.

### Coupling Risks

Global coupling: [x] None.

Content coupling: [x] None, if CLI consumes a stable result rather than internal parser structs.

Stamp coupling: [x] Low risk. Keep request/result structs focused so functions do not receive large app-wide objects.

### Narrow Dependency Check

Each function should receive only what it needs. Pure normalization/rendering helpers should take explicit values, not full command context. The initialization service may receive a workspace root/output root and options because the full use-case concept is needed.

### Dependency Injection Check

Side-effectful dependencies are created at the edge and passed inward. `internal/app` or the service constructor should supply the clone runner and output writer. Pure helpers stay concrete.

## 7. State And Mutation Check

### Mutation Points

- `internal/study/init.go`: coordinates plan execution and force/dry-run policy.
- `internal/study/init_yaml.go`: parses YAML but should not mutate durable state.
- `internal/study/init_render.go`: renders file contents but should not write durable state.
- `internal/study/init_clone.go`: executes clone actions through a command boundary.
- `internal/app/study_init_commands.go`: writes user-facing command output only.

### Is Mutation Explicit From Names And Flow?

[x] Yes. Planning/rendering functions should be separated from execution functions. `--dry-run` must stop at the final execution boundary.

### Are Queries And Commands Separated Where Practical?

[x] Yes. Parse/validate/plan/render are query-like and deterministic. Execute writes files and runs clones.

### Could Hidden State Make This Hard To Debug?

[x] No, if workspace root, output path, force, dry-run, clone policy, and clone runner are explicit request/dependency values. Hidden globals for current directory, stdout/stderr, or process execution would be a risk and are rejected.

## 8. Function/Class/Module Shape

### Primary Unit Of Behaviour

[x] Service/component.
[x] Pipeline/workflow.
[x] Adapter for clone execution only.

### Why This Unit?

Study initialization is a use case with validation, derived state, generated files, scoped mutation, and partial clone failure. A small service/use-case function in `internal/study` is more appropriate than a free-floating CLI handler, while a clone adapter is justified because subprocess execution is volatile and must be fakeable.

### Composition

- YAML parser: reads and validates input shape into explicit definitions.
- Normalizer/validator: enforces names, counts, numbers, slugs, uniqueness, and path safety.
- Renderer: produces deterministic YAML, README, and dimension Markdown content.
- Planner: lists directories, files, and clone actions.
- Executor: creates directories/files and invokes clone runner after artifact writes.
- CLI command: maps flags to request and prints result/errors.

## 9. Error Handling Design

### Expected Failures

- Invalid usage or missing YAML path: handled in `internal/app` as usage error.
- Invalid YAML syntax: validation error naming the file and parse cause.
- Missing required field: validation error naming the field path, for example `dimensions.items[0].title`.
- Duplicate source name, dimension number, or dimension slug: validation error naming both conflicting values where practical.
- Unsafe study/source/dimension/output name: validation or workspace/path error with corrective guidance.
- `repos.count` or `dimensions.count` lower than explicit item count: validation error.
- Count greater than explicit items: validation error explaining assisted completion is deferred and explicit items are required now.
- Existing study directory without `--force`: overwrite protection error.
- Clone failure: partial/failure result that preserves generated artifact success and identifies failed source names and clone destinations.

### Unexpected Failures

- Filesystem errors should wrap the failing operation and path.
- Git executable lookup or process failures should wrap the executable, args summary, source name, and destination without shell-expanded strings or secrets.

### Error Taxonomy

Reuse existing CLI exit-code conventions from prior sprints. Introduce or reuse typed/sentinel study errors only where command mapping needs to distinguish usage, validation, workspace/filesystem, and partial completion. Avoid a stable public JSON error schema in this sprint.

### Retry/Recovery Behaviour

- Retry safe? The whole command is retry-safe only after the user fixes validation/path issues or uses `--force` for existing generated files. Clone retries are user-initiated by rerunning.
- Partial progress possible? Yes, clone failures can occur after generated artifacts are written.
- Compensation needed? No automatic deletion or rollback. Report generated artifacts and failed clones clearly.

## 10. Observability Design

### Events/Logs/Metrics/Traces Needed

- Human-readable success summary listing created root and key generated paths.
- Dry-run output listing planned directories, files, and clone actions.
- Clone failure summary listing successful artifact creation and failed clone actions.
- No persistent event records, structured logs, metrics, traces, or runtime diagnostics in this sprint.

### Minimum Useful Fields

- Command name.
- Workspace root or selected output root when safe.
- Study name.
- Generated study path.
- Counts for sources, dimensions, files, directories, and clone actions.
- Failed clone source names and destinations.

### Sensitive Data Check

Could this expose secrets/user content? [x] Yes.

Mitigation: generated artifacts should contain only user-supplied study metadata and normalized source/dimension definitions. Do not print environment variables, credentials, full command environment, or hidden runtime diagnostics. Treat URLs as user-visible source metadata but do not add credentials or secrets.

## 11. Testing Strategy

### Unit Tests

- Parse valid YAML into explicit study, source, and dimension definitions.
- Reject invalid YAML and missing required fields with field-specific messages.
- Normalize dimension numbers to two digits and slugs to kebab-case.
- Reject duplicate source names, duplicate dimension numbers/slugs, unsafe names, unsafe output paths, and unsupported count shortages.
- Render deterministic normalized YAML, README, and dimension Markdown preserving purpose, steps, citations, and questions.

### Use-Case/Workflow Tests

- Successful init creates the required directory structure and generated files under a temp workspace.
- `--dry-run` uses the same plan but creates no files and invokes no clone runner.
- Existing study directory fails without `--force`.
- `--force` overwrites known generated files inside the study directory without deleting unknown files or touching outside paths.
- `--no-clone` skips clone execution.
- URL-backed sources produce `git clone --depth 1 <url> <sources/name>` clone actions through a fake runner.
- Clone failure returns partial/failure status while generated artifacts remain present.

### Integration/Command Tests

- Help output includes `ultraplan study init <study-init.yml>` and flags.
- Command tests cover usage errors, dry-run output, overwrite protection, workspace/global `--workspace`, output path behaviour, and exit-code mapping.
- Normal tests must not require network, GitHub, OpenCode, provider credentials, or real cloned repositories.

### Regression Tests

- Count greater than explicit items fails with deferred-assisted-completion guidance.
- Normalized YAML avoids absolute paths unless preserving an explicitly allowed custom output path for diagnostics.
- `--force` never modifies paths outside the selected study directory.

### Testability Check

Can the core behaviour be tested without real infrastructure? [x] Yes. Temp directories, fixture YAML, command IO capture, and fake clone runners cover the required behaviour.

## 12. Performance And Scale Assumptions

### Known Constraints

- Expected volume: small to moderate numbers of sources and dimensions during initialization.
- Latency sensitivity: normal CLI expectations; artifact planning and rendering should be fast.
- Memory sensitivity: low; the command reads a YAML file and renders small generated artifacts.
- Concurrency concerns: none required. Clone actions may be sequential for simpler partial-failure reporting.

### Assumptions Being Made

- Initialization must not recursively inspect cloned source contents.
- Source cloning duration dominates runtime when enabled.
- No worker pool or concurrent clone execution is necessary for this sprint.

### Measurement Plan

[x] Not needed for this change.

### Optimization Decision

Do not optimize beyond deterministic sorting and avoiding recursive source scans. The selected evidence excluded the performance report, and sprint requirements do not depend on scheduler throughput or large report processing.

## 13. Implementation Plan

### Files Likely To Change

- `/home/antonioborgerees/coding/ultraplan-go/internal/study/init.go`: initialization request/result, validation, planning, force/dry-run execution coordination.
- `/home/antonioborgerees/coding/ultraplan-go/internal/study/init_yaml.go`: YAML parser and input schema validation.
- `/home/antonioborgerees/coding/ultraplan-go/internal/study/init_render.go`: deterministic normalized YAML, README, and dimension Markdown rendering.
- `/home/antonioborgerees/coding/ultraplan-go/internal/study/init_clone.go`: clone action model and fakeable non-shell clone runner.
- `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_init_commands.go`: command registration, flags, output, and exit mapping.
- `/home/antonioborgerees/coding/ultraplan-go/internal/study/init_test.go`: service, parser, renderer, filesystem, dry-run, force, and clone tests.
- `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_init_commands_test.go`: command UX and exit-code tests.

### Step-By-Step Plan

1. Define `study.InitRequest`, `study.InitResult`, clone action/result shape, and a narrow clone runner interface.
2. Implement YAML parsing and validation for required explicit fields, counts, uniqueness, and normalization.
3. Implement plan creation with workspace/output path safety, directories, files, and clone actions.
4. Implement deterministic renderers for normalized YAML, README, and dimension Markdown.
5. Implement execution that writes directories/files and then executes clone actions unless dry-run or no-clone applies.
6. Wire `internal/app` command flags and output to the study service.
7. Add unit and command tests with fixtures, temp directories, and fake clone runner.
8. Verify with `go test ./...` and `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go` during implementation review.

### Migration/Backwards Compatibility

- Schema migration needed? No.
- API compatibility affected? No public API is introduced.
- Existing data affected? Existing study directories are protected by default; `--force` is explicit and scoped.
- Feature flag needed? No.

### Rollback Plan

Revert the study init command wiring and the new `internal/study/init*` files. Generated study directories are user artifacts and should not be automatically deleted by rollback.

## 14. Final Pre-Implementation Decision

Decision: Proceed.

Reason:
Implement a focused `internal/study` initialization use case with a plan-then-execute flow, deterministic renderers, scoped filesystem mutation, and a narrow fakeable clone runner. Keep CLI handlers thin and avoid new global technical-layer packages.

Complexity introduced:
Medium. The feature introduces parsing, planning, rendering, file writes, clone execution, and command wiring, but each part is bounded and directly required.

Complexity removed:
Low. A shared plan removes divergence between dry-run and real execution.

Main trade-off:
Accept a small internal plan/result model to keep dry-run, force safety, command output, and tests aligned. Reject broader workflow/filesystem abstractions until multiple product modules need them.

## Key Conclusion And Evidence Basis

The architecture decision is to keep all study initialization rules and artifact generation in `internal/study`, with `internal/app` limited to command wiring and output. This is grounded in the handbook evidence that mature Go CLIs keep entrypoints thin (`01-project-structure`, `02-command-architecture`), use explicit dependencies and seams for side effects (`03-dependency-injection`, `06-io-abstraction`), use a plan/decorator style for dry-run safety (`15-philosophy`), classify actionable errors (`05-error-handling`), avoid shell interpolation for subprocesses (`13-security`), and test command/file behaviour with fakes and fixtures (`11-testing-strategy`). It also matches the project Architecture rule that product modules own product behaviour and platform packages must not import `study`.

## Rejected Alternatives

- Rejected putting YAML validation, rendering, and file writes directly in `internal/app` because it would violate the sprint constraint and the handbook warning against business rules in command handlers.
- Rejected creating global validation, reports, scheduler, workflow, or filesystem-planner packages because the Architecture document explicitly warns against global technical-layer packages that fracture product context.
- Rejected runtime-assisted completion and source URL repair paths because the sprint requirements explicitly defer them.
- Rejected rollback/delete-on-clone-failure because requirements say clone failures must not hide successful generated artifact creation, and `--force` must not remove unrelated paths.

## Risks, Assumptions, And Open Questions

Risks:

- `--output <dir>` path rules can be subtle; implementation must clearly define whether output is workspace-scoped or explicitly allowed outside workspace and preserve diagnostics accordingly.
- Partial clone failure needs careful exit-code mapping to existing conventions without introducing a stable public JSON schema.
- `--force` could accidentally overwrite user edits if generated-file boundaries are too broad; prefer overwriting known generated paths only and never deleting unknown files.
- Normalized YAML ordering must be deterministic or fixture tests will be noisy.

Assumptions:

- Sequential clone execution is acceptable for Sprint 4.
- Existing workspace APIs are sufficient for root resolution and path safety.
- The command can expose human-readable dry-run output without adding JSON output for this sprint.
- The local `git` executable is only needed for real clone execution, not tests.

Open questions:

- What exact exit code currently represents partial completion in the existing CLI implementation, and should clone partial failure map directly to it?
- Should `--force` overwrite only known generated files or also replace existing source clone directories for URL-backed sources? The safer default is known generated files only unless implementation requirements specify clone-directory replacement.
- Should an explicitly supplied `--output <dir>` be allowed outside the workspace, or only within workspace-managed paths? The requirements allow an explicitly selected output root but still require path normalization and safety checks.
