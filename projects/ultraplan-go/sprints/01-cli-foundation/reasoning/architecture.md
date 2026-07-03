> **Inputs Used:** `projects/ultraplan-go/sprints/01-cli-foundation/sprint-index.md`, `projects/ultraplan-go/sprints/01-cli-foundation/technical-handbook.md`, `projects/ultraplan-go/sprints/01-cli-foundation/requirements.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `system/reasoning/architecture_reasoning_template.md`

# Architecture Reasoning: CLI Foundation

## 0. Feature Summary

### Feature name

CLI foundation architecture.

### Area name and scope

This architecture reasoning covers the initial Go module shape, the `cmd/ultraplan` process entrypoint, the `internal/app` composition root, deterministic help/version/invalid-command behavior, exit status mapping, injectable command IO, and the initial module-owned package skeletons required by the sprint.

It does not cover workspace discovery, configuration loading, runtime integration, study workflows, code-reference extraction behavior, target workflows, sprint workflows, packaging, release automation, or real-runtime smoke tests.

### User/product goal

Establish a buildable, testable `ultraplan` Go CLI shell that gives future study-side work a clear module boundary without prematurely implementing deferred product behavior.

### Current task

Implement the smallest architecture that satisfies this sprint: a Go module, thin executable entrypoint, app-owned command shell for help/version/unknown command behavior, deterministic command tests, version metadata, and documented package skeletons for future platform and product modules.

### Non-goals

- Do not implement workspace discovery, workspace initialization, config precedence, `config show`, redaction, or health checks.
- Do not implement study initialization, study listing, prompt composition, analysis runs, synthesis, summaries, validation, or run-loop state.
- Do not integrate `github.com/Antonio7098/agentwrap`, OpenCode, runtime health, runtime execution, policies, events, or permissions.
- Do not implement code-reference parsing, citation resolution, extraction output, or `ultraplan code`.
- Do not implement target or sprint workflows, commands, persisted schemas, validators, or execution behavior.
- Do not add packaging, release artifacts, generated docs, checksums, shell completion, or real-runtime smoke tests.

## 1. First-Principles Breakdown

### Core behaviour

This sprint receives process arguments and build metadata, routes them through a minimal app-owned command shell, writes deterministic text to configured output streams, and returns meaningful process-independent exit codes.

### Inputs

- CLI arguments passed by `cmd/ultraplan/main.go`.
- Standard output and standard error streams supplied at the process edge.
- Build/version metadata values for version, commit, build date, and Go version.
- Static command definitions for root help and `version`.

### Outputs

- Help-oriented command discovery for `ultraplan` with no args and for `ultraplan --help`.
- Version output with version, commit, build date, and Go version fields.
- Actionable unknown-command diagnostic naming the invalid command.
- Exit code `0` for help/version success and `2` for usage/argument errors.
- Buildable package skeletons documenting future module boundaries.

### Durable state

- Created: Go source files and package documentation files required by the sprint.
- Read: no runtime durable product state.
- Updated: no workspace, config, study, run-state, report, or extraction state.
- Deleted: no durable state.

### Ephemeral state

- Parsed argument slice.
- Selected command branch.
- App configuration containing streams, args, and version metadata.
- Usage error classification for an unknown command.

### Derived state

- Exit status is derived from command outcome and error classification, not stored.
- Help text is derived from static command metadata and rendered on demand.
- Version text is derived from immutable build metadata and rendered on demand.

### Side effects

- Database write: no.
- File write: only sprint implementation files and skeleton package docs are created by implementation; runtime command execution does not write product artifacts.
- Network call: no.
- Queue/event emission: no.
- Email/notification: no.
- Logs/metrics/traces: no logging required in this sprint; CLI output is the observable behavior.

## 2. Existing Architecture Fit

### Where does this behaviour naturally belong?

The behavior belongs in `internal/app`. The process boundary belongs in `cmd/ultraplan/main.go`. Future platform and product modules belong under their owning packages but remain skeletons until their behavior enters scope.

The evidence basis is strong. The technical handbook identifies thin CLI entrypoints as the dominant pattern, citing `cmd/gh/main.go:6`, `main.go:23`, and `cmd/restic/main.go:37-114`. It also identifies manual composition roots as the normal Go CLI shape, citing `internal/ghcmd/cmd.go:52-132`, `internal/app/app.go:42-81`, and `cmd/restic/main.go:181-183`. The architecture document states the expected dependency direction as `cmd/ultraplan -> internal/app`, `internal/app -> product modules + platform modules`, and `platform/* -> no product modules`.

### Existing workflow affected

Current flow:

1. No implementation exists in this sprint directory or target Go module for the selected CLI foundation.
2. Sprint requirements define required output paths but no prior sprint carry-forward decisions exist.
3. The architecture document defines module ownership and dependency rules for future work.

### Proposed new flow

Proposed flow:

1. `cmd/ultraplan/main.go` creates an app configuration using `os.Args[1:]`, `os.Stdout`, `os.Stderr`, and build metadata defaults.
2. `main.go` delegates to `internal/app.Run` or equivalent and receives an integer status code.
3. `internal/app` parses only root help, `--help`, `-h`, `version`, and unknown-command cases.
4. `internal/app` writes normal output to stdout and usage diagnostics to stderr.
5. `internal/app` maps successful help/version to `0` and unknown commands to `2`.
6. Tests instantiate `internal/app` directly with in-memory streams and argument slices.

### Does the current architecture still fit?

```text
[x] Yes - the feature fits cleanly into the existing shape.
[ ] Partially - small refactor needed before adding the feature.
[ ] No - the feature reveals that the existing architecture is the wrong shape.
```

### If it does not fit, why?

No mismatch is present. The sprint is the initial implementation and matches the documented module-driven architecture.

### Refactor-before-feature decision

```text
Decision: implement directly.
Reason: there is no existing code to refactor, and the required shape is already defined by requirements and architecture docs.
```

## 3. Design Options Considered

### Option A: Minimal standard-library command shell in `internal/app`

Description:
Implement a small argument dispatcher without a third-party CLI framework. Keep `cmd/ultraplan/main.go` as a thin process entrypoint and put command behavior, IO, version rendering, and exit-code mapping in `internal/app`.

Pros:
- Satisfies all sprint acceptance criteria with the least surface area.
- Keeps command behavior directly testable with injected streams.
- Avoids framework-specific help formatting before real command hierarchy exists.
- Aligns with the handbook's evidence for thin entrypoints, manual composition roots, and explicit non-goals.
- Reduces dependency and release complexity during module bootstrap.

Cons:
- Later nested commands may require either careful extension or migration to a framework.
- Help formatting and flag handling must be maintained locally for now.
- Does not provide shell completion, command aliases, or rich flag parsing.

Risks:
- If implemented as ad hoc branching without clear tests, the command shell could become hard to extend. This is mitigated by keeping only help/version/unknown cases in scope and testing behavior rather than internals.

### Option B: Introduce a CLI framework now

Description:
Use a framework such as Cobra or urfave/cli immediately, even though this sprint only needs root help, version, and unknown-command behavior.

Pros:
- Gives future nested commands a known command tree model.
- Provides standard help and flag handling.
- Could reduce custom parser code later.

Cons:
- Adds an external dependency before the sprint has enough command complexity to justify it.
- Framework defaults may advertise or shape deferred commands too early.
- Tests may need framework-specific assertions instead of simple app behavior assertions.
- Increases bootstrap complexity while the sprint requirements explicitly favor a thin shell.

Risks:
- Premature framework choice may constrain future study command ergonomics before those workflows are reasoned through.

### Option C: Put command logic in `cmd/ultraplan/main.go`

Description:
Implement parsing and output directly in `main.go` and use process-level tests or manual CLI checks.

Pros:
- Very few files.
- Simple to understand for the first few lines of behavior.

Cons:
- Violates the sprint requirement that `cmd/ultraplan/main.go` contain no product workflow logic beyond process entry, calling `internal/app`, and exiting.
- Conflicts with handbook warnings not to let the entrypoint accumulate command logic.
- Makes deterministic command tests harder because `os.Exit` and process streams become entangled with behavior.

Risks:
- Creates the exact architecture debt this sprint is intended to avoid.

### Option D: Build package skeletons plus early service interfaces

Description:
Create documented skeleton packages and define interfaces or service structs for config, logging, filesystem, runtime, workspace, study, and code extraction.

Pros:
- Makes future seams explicit.
- Could provide initial compile-time targets for later work.

Cons:
- Introduces unearned abstractions for behavior that is explicitly out of scope.
- Risks defining contracts before requirements are selected for those areas.
- Conflicts with architecture guidance that interfaces should appear only at external or volatile boundaries when the need is concrete.

Risks:
- Future implementation may need to unwind speculative interfaces.

### Chosen option

```text
Chosen option: Option A.
Reason: it is the smallest honest design for a sprint that only needs help, version, invalid-command behavior, deterministic tests, thin process entry, and package boundary documentation.
```

Rejected alternatives:
- Reject Option B because a CLI framework is not earned by this sprint's tiny command surface.
- Reject Option C because it violates the thin-entrypoint acceptance criterion and handbook evidence.
- Reject Option D because skeleton package docs are required, but service interfaces are not yet supported by selected scope.

## 4. Abstraction Check

### Are we adding a new abstraction?

```text
[ ] No
[ ] Yes - interface/protocol/trait
[ ] Yes - base class
[x] Yes - service/component/module
[ ] Yes - strategy/function parameter
[x] Yes - data structure/DTO/config object
[ ] Yes - plugin/registry/factory
```

The earned abstractions are the `internal/app` module and a small app configuration or options value that carries args, streams, and version metadata. No plugin registry, framework adapter, runtime contract, or product service interface is earned in this sprint.

### If yes, why is it earned?

```text
[ ] Multiple real implementations exist now
[ ] A second implementation is very likely soon
[ ] It isolates an external dependency
[ ] It removes branching rather than adding it
[x] It protects a stable domain concept
[x] It makes testing simpler
[x] It creates a clear boundary between policy and mechanism
[ ] Other
```

The `internal/app` boundary is earned because it separates process mechanics from command behavior. The handbook specifically connects deterministic command tests to injectable IO streams and test constructors, citing `pkg/iostreams/iostreams.go:551-568`, `executor.go:553-564`, and `ui.go:19-43`.

### If no, why are we keeping unselected areas concrete?

Unselected platform and product packages remain documentation-only skeletons because config, logging, filesystem, runtime, workspace, study, and code extraction behavior are explicit non-goals. Defining interfaces now would be speculative.

### Bad abstraction smell check

```text
- [ ] Generic name like Manager/Handler/Processor/Common/Helper
- [ ] Only one implementation with no real volatility
- [ ] Boolean flags or mode switches
- [ ] Optional parameters for caller-specific behaviour
- [ ] Interface mirrors a concrete class rather than consumer needs
- [ ] Abstraction exists only to reduce line count
- [ ] Abstraction makes the flow harder to trace
```

Implementation should avoid generic package names beyond the required `internal/app` composition root and the explicit platform/product package names from architecture docs.

## 5. DRY and Duplication Check

### Is any duplication being introduced?

```text
[x] No
[ ] Yes - duplicated business/domain knowledge
[ ] Yes - duplicated code shape only
[ ] Yes - duplicated infrastructure mechanics
```

Help rendering should use one path for no args and `--help` to avoid duplicating command discovery text. Version output should use one renderer for all version metadata fields. Exit-code mapping should live in `internal/app`, not be repeated in `main.go` and tests.

### If duplicating business knowledge, stop and consolidate

Shared command names, help text, and exit status meanings should have a single source inside `internal/app`.

### If duplicating code shape, is that acceptable?

No intentional code-shape duplication is needed for this sprint.

### If removing duplication, is the abstraction safe?

The safe consolidation point is not a generic command framework; it is a small local function or data value inside `internal/app` because there is only one real command branch beyond help.

## 6. Coupling Check

### Dependencies required

- Go standard library: required for process entry, streams, runtime Go version, strings/formatting, and tests.
- No third-party CLI framework: not required by the sprint and not justified by selected evidence.
- No agentwrap dependency in this sprint: runtime integration is deferred by requirements, even though the TRD requires agentwrap later.

### Coupling risks

```text
Global coupling:
[x] None
[ ] Introduced
[ ] Existing but unchanged
[ ] Existing and reduced

Content coupling:
[x] None
[ ] Introduced
[ ] Existing but unchanged
[ ] Existing and reduced

Stamp coupling:
[ ] None
[ ] Introduced
[ ] Existing but unchanged
[x] Existing and controlled
```

The only acceptable stamp coupling is passing an app options/config value into `internal/app` because the app boundary needs args, streams, and version metadata as a cohesive process-independent request.

### Narrow dependency check

```text
Does each function receive only what it needs?
[x] Yes
[ ] No - reason

Are large objects passed only when the full concept is needed?
[x] Yes
[ ] No - reason
```

The process entrypoint should pass a narrow app configuration, not global state or product services.

### Dependency injection check

```text
Are side-effectful dependencies created at the edge and passed inward?
[x] Yes
[ ] No
[ ] Not applicable
```

`os.Stdout`, `os.Stderr`, `os.Args`, and process exit belong at `cmd/ultraplan/main.go`. `internal/app` should operate on provided streams and args so tests can run without process-level side effects.

## 7. State and Mutation Check

### Mutation points

- `cmd/ultraplan/main.go`: calls `os.Exit` with the status returned by `internal/app`.
- `internal/app`: writes command output to provided stdout/stderr streams.
- Implementation-time package skeleton creation: creates source files only, not product runtime state.

### Is mutation explicit from names and flow?

```text
[x] Yes
[ ] No - improvement needed
```

### Are queries and commands separated where practical?

```text
[x] Yes
[ ] No - reason
```

Rendering help/version is pure except for writing to streams. Unknown-command handling classifies a usage error and writes a diagnostic.

### Could hidden state make this hard to debug?

```text
[x] No
[ ] Yes - risk
```

No global mutable config, runtime singleton, or workspace cache is needed. The handbook explicitly warns against global mutable config or service singletons for this shell sprint.

## 8. Function/Class/Module Shape

### Primary unit of behaviour

```text
[x] Function
[ ] Class/object
[x] Module
[ ] Service/component
[ ] Pipeline/workflow
[ ] Adapter
[ ] Other
```

### Why this unit?

The primary behavior should be a function in `internal/app`, such as `Run`, because the command shell has no lifecycle and no durable state. The module boundary is still important because it keeps command behavior out of `main.go` and gives later study-side composition a natural home.

### If using a class/object, what justifies it?

No class/object is required. A small struct for app options or version metadata is a data value, not a lifecycle-owning object.

### If using inheritance, why not composition?

Inheritance is not used in Go and is not relevant.

### If using composition, what components are combined?

- Args: process-independent command input.
- Output streams: testable command output boundary.
- Version metadata: immutable build data for `ultraplan version`.
- Exit-code mapping: app policy for process status without calling `os.Exit` internally.

## 9. Error Handling Design

### Expected failures

- Unknown command: print an actionable diagnostic naming the command and return exit code `2`.
- Unsupported flags other than help aliases: treat as usage errors and return exit code `2` if implemented in this sprint.
- Stream write failure: return a general failure code if surfaced by the writer; tests may not need to cover this unless the implementation explicitly returns write errors.

### Unexpected failures

- Unexpected internal write or formatting errors should not panic. They should surface through the app result and map to a non-zero general failure if the implementation tracks them.

### Error taxonomy

- Success: exit code `0`.
- Usage or argument error: exit code `2`.
- General failure: exit code `1` only if an unexpected internal error is encountered.

The evidence basis is the TRD's suggested exit code classes, where `0` is success, `1` is general failure, and `2` is usage or argument error. The handbook also cites mature CLIs using user-facing hints and exit-code mapping, including `internal/ghcmd/cmd.go:44-49`, `cmd/age/tui.go:37-54`, and `go-task/errors/errors_task.go:13-32`.

### Retry / recovery behaviour

- Retry safe: no runtime work is performed, so retry is harmless.
- Partial progress possible: no product durable state is modified.
- Compensation needed: none.

## 10. Observability Design

### Events/logs/metrics/traces needed

- Help output: deterministic stdout behavior for discovery.
- Version output: deterministic stdout behavior for build metadata.
- Unknown-command diagnostic: deterministic stderr behavior for actionable usage failure.

No logs, metrics, traces, runtime events, or observability records are required by this sprint. The sprint index excludes the Observability contract and logging behavior, and the requirements make logging package creation a skeleton-only task.

### Minimum useful fields

- Command name in help output.
- Available command list including `version`.
- Version, commit, build date, and Go version in version output.
- Unknown command value in usage diagnostics.

### Sensitive data check

```text
Could this expose PII/secrets/user content?
[x] No
[ ] Yes - mitigation
```

No config, environment, workspace, source, runtime, or provider data should be read or printed in this sprint.

## 11. Testing Strategy

### Unit tests

- Help with `--help` exits `0` and includes the binary name, available commands, and `version`.
- Help with no args exits `0` and uses the same help-oriented command discovery as `--help`.
- `version` exits `0` and includes version, commit, build date, and Go version fields.
- Unknown command exits `2`, names the unknown command, and writes an actionable diagnostic.
- Help output does not advertise deferred target, sprint, workspace initialization, config, health, study, runtime, summary, validation, or code extraction commands.

### Use-case/workflow tests

- Instantiate `internal/app` with in-memory stdout/stderr and explicit args.
- Assert output and exit status behavior without invoking a real binary, OpenCode, provider credentials, network, or a configured workspace.

### Integration tests

- `go test ./...` from `../ultraplan-go` is the required verification command.
- `go build ./cmd/ultraplan` from `../ultraplan-go` is the required build verification command.
- No real-runtime integration test is included in this sprint.

### Regression tests

- Protect against accidental help-surface expansion into deferred commands.
- Protect against moving command behavior into process-only code that cannot be tested without `os.Exit`.

### Testability check

```text
Can the core behaviour be tested without real infrastructure?
[x] Yes
[ ] No - why
```

The handbook's testing evidence favors behavior assertions with captured output and exit codes, citing `internal/cmd/main_test.go:64-174`, `go-task/task_test.go:166-169`, and `gdu/cmd/gdu/app/app_test.go:682`.

## 12. Performance and Scale Assumptions

### Known constraints

- Expected volume: one short command invocation at a time.
- Latency sensitivity: simple help and version commands should start quickly; PRD later targets simple CLI startup under 250 ms excluding large scans.
- Memory sensitivity: negligible.
- Concurrency concerns: none.

### Assumptions being made

- Static help text and version rendering are enough for the sprint.
- Future command expansion will be reasoned in later sprints before requiring a framework.
- No filesystem scanning, workspace lookup, config loading, network call, or runtime process execution happens in this sprint.

### Measurement plan

```text
[x] Not needed for this change
[ ] Unit benchmark
[ ] Profiling
[ ] Load test
[ ] Production metric
```

### Optimization decision

Do not optimize beyond avoiding unnecessary dependencies and side effects. The sprint has no performance-sensitive workload.

## 13. Implementation Plan

### Files likely to change

- `../ultraplan-go/go.mod`: define the Go module and Go toolchain version.
- `../ultraplan-go/cmd/ultraplan/main.go`: process entrypoint that delegates to `internal/app` and exits with returned status.
- `../ultraplan-go/internal/app/app.go`: command parsing, output routing, help rendering, invalid-command handling, and exit-code mapping.
- `../ultraplan-go/internal/app/version.go`: version metadata defaults and Go version field.
- `../ultraplan-go/internal/app/app_test.go`: deterministic command behavior tests.
- `../ultraplan-go/internal/platform/config/doc.go`: config boundary documentation only.
- `../ultraplan-go/internal/platform/logging/doc.go`: logging boundary documentation only.
- `../ultraplan-go/internal/platform/filesystem/doc.go`: filesystem boundary documentation only.
- `../ultraplan-go/internal/platform/runtime/doc.go`: generic runtime boundary documentation only; no agentwrap integration.
- `../ultraplan-go/internal/workspace/doc.go`: workspace module boundary documentation only.
- `../ultraplan-go/internal/study/doc.go`: study module boundary documentation only.
- `../ultraplan-go/internal/codeextract/doc.go`: code extraction boundary documentation only.

### Step-by-step plan

1. Create the Go module and required directory structure.
2. Implement `internal/app` with a minimal options/config value, injected args and streams, static help rendering, version rendering, and exit-code mapping.
3. Implement `cmd/ultraplan/main.go` as a process-only wrapper around `internal/app`.
4. Add version metadata defaults and use `runtime.Version()` for the Go version field.
5. Add package `doc.go` skeletons that document ownership and explicitly state deferred behavior.
6. Add command behavior tests with in-memory streams.
7. Run `go test ./...` and `go build ./cmd/ultraplan` from `../ultraplan-go`.

### Migration/backwards compatibility

- Schema migration needed: no.
- API compatibility affected: no public Go API is created.
- Existing data affected: no.
- Feature flag needed: no.

### Rollback plan

If the initial shell causes problems, remove the new Go module files for this sprint. No product durable state or user workspace data is modified by the sprint implementation.

## 14. Final Pre-Implementation Decision

```text
Decision: Proceed.

Reason:
The architecture should use a thin `cmd/ultraplan` process entrypoint, a standard-library `internal/app` command shell, injected IO for deterministic tests, minimal exit-code classification, and documentation-only skeletons for future modules.

Complexity introduced:
Low.

Complexity removed:
Low.

Main trade-off:
Accept a small amount of local command dispatch code now to avoid an unearned CLI framework and keep the first sprint focused on buildability, testability, and package boundaries.
```

## Key Conclusion And Evidence Basis

The final architecture decision is to implement the CLI foundation with a minimal standard-library app shell in `internal/app`, keep `cmd/ultraplan/main.go` process-only, inject args and IO for tests, map usage errors to exit code `2`, and create only documented skeletons for future modules.

Evidence basis:
- The handbook's project-structure evidence supports thin entrypoints and protected internal packages, citing `cmd/gh/main.go:6`, `cmd/restic/main.go:37-114`, `internal/chezmoi/chezmoi.go:1-2`, and opencode's internal package shape.
- The handbook's dependency-injection evidence supports a visible manual composition root, citing `internal/ghcmd/cmd.go:52-132`, `internal/app/app.go:42-81`, and `cmd/restic/main.go:181-183`.
- The handbook's IO evidence supports injectable streams for deterministic command tests, citing `pkg/iostreams/iostreams.go:551-568`, `executor.go:553-564`, and `ui.go:19-43`.
- The handbook's error-handling evidence supports actionable diagnostics and exit-code mapping, citing `cmd/age/tui.go:37-54`, `internal/ghcmd/cmd.go:44-49`, and `go-task/errors/errors_task.go:13-32`.
- The handbook's philosophy evidence supports strict non-goals and delayed complexity, citing age's no-keyring stance, fzf's focused filter purpose, and lazygit's complexity rejection.
- The sprint requirements require thin entrypoint behavior, offline deterministic tests, help/version/unknown-command behavior, exit code `2` for unknown commands, and skeleton package files without deferred workflow implementation.
- The architecture document requires module-driven ownership, `cmd/ultraplan -> internal/app`, product behavior under owning modules, platform packages with no product imports, and no global technical-layer packages.

## Risks, Assumptions, And Open Questions

Risks:
- A minimal dispatcher may need replacement or refactoring when nested study commands enter scope.
- Static help text could drift from implemented commands if future commands are added without updating tests.
- Skeleton package docs could overpromise future behavior if they describe implementation details instead of boundaries.
- Unknown-command diagnostics need a consistent stdout/stderr decision; this reasoning chooses stderr for diagnostics because script-friendly CLIs usually reserve stdout for successful output.

Assumptions:
- The first implementation sprint does not need third-party dependencies.
- `internal/app` can own command behavior without importing future product modules yet.
- Build metadata defaults are acceptable until release-time linker injection is introduced.
- No-argument help and `--help` should share the same rendering path to keep behavior deterministic.

Open questions for later sprints:
- When real study commands are selected, should UltraPlan continue extending the local dispatcher or adopt a CLI framework?
- What exact public or internal API should future product modules expose to `internal/app` once study behavior exists?
- How should JSON output mode be introduced later without complicating this sprint's text-only shell?
- Where should richer error categories live once config, workspace, runtime, validation, and cancellation enter scope?
