# Sprint Reasoning: CLI Foundation

> Project: `ultraplan-go`
> Sprint: `01-cli-foundation`
> Output: `projects/ultraplan-go/sprints/01-cli-foundation/reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/01-cli-foundation/requirements.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/sprints/01-cli-foundation/sprint-index.md`, `projects/ultraplan-go/sprints/01-cli-foundation/technical-handbook.md`, `projects/ultraplan-go/sprints/01-cli-foundation/reasoning/architecture.md`, `templates/sprint-reasoning.md`, `system/contracts/core/architecture.md`, `system/contracts/core/errors.md`, `system/contracts/core/testing.md`, `system/contracts/core/documentation.md`, `system/contracts/surfaces/cli.md`, `studies/go-cli-study/reports/final/01-project-structure.md`, `studies/go-cli-study/reports/final/02-command-architecture.md`, `studies/go-cli-study/reports/final/03-dependency-injection.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/06-io-abstraction.md`, `studies/go-cli-study/reports/final/09-terminal-ux.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/12-extensibility.md`, `studies/go-cli-study/reports/final/15-philosophy.md`

This document decides. It synthesizes selected context, handbook evidence, area-specific reasoning, and contracts into final sprint decisions.

It does not replace `sprint-index.md`, `technical-handbook.md`, or `reasoning/*.md`.

## Sprint Purpose

- **Goal:** Establish a buildable Go module with a thin `ultraplan` CLI shell, an `internal/app` composition root, deterministic help/version/unknown-command behavior, process-independent command tests, and documented initial package boundaries for future study-side modules.
- **Non-Goals:** Do not implement workspace discovery, workspace initialization, config loading, config precedence, secret redaction, health checks, study modeling or workflows, runtime execution, agentwrap/OpenCode integration, code extraction behavior, target workflows, sprint workflows, packaging, release artifacts, shell completion, generated docs, or real-runtime smoke tests.
- **Depends On:** Phase 0 roadmap and scope alignment, the project PRD/TRD/Architecture docs, the Go toolchain, and the selected evidence reports in `sprint-index.md`. No prior sprint decisions exist for this first implementation sprint.

## Selected Context And Pre-Reasoning Artifacts

| Artifact | Path | How It Was Used |
| --- | --- | --- |
| Sprint Requirements | `projects/ultraplan-go/sprints/01-cli-foundation/requirements.md` | Authoritative sprint contract for required outputs, acceptance criteria, non-goals, constraints, dependencies, and review expectations. Local requirement labels in this document reference its sections and rows. |
| Project Index | `projects/ultraplan-go/project-index.md` | Confirmed the repository target, active contract pool, available evidence reports, selected study, and absence of prior project decisions. |
| Product Requirements | `projects/ultraplan-go/docs/PRD.md` | Confirmed the product is a local-first `ultraplan` CLI, that CLI bootstrap/help is a must-have, and that target/sprint workflows are deferred from the current study-side build. |
| Technical Requirements | `projects/ultraplan-go/docs/TRD.md` | Confirmed Go layout, script-friendly CLI expectations, exit code classes, version command fields, build/test requirements, future agentwrap boundary, and current scope clarification. |
| Architecture | `projects/ultraplan-go/docs/ARCHITECTURE.md` | Confirmed module-driven ownership, dependency direction, thin `cmd/ultraplan -> internal/app` entry path, platform/product separation, one package per module, and rejection of global technical-layer packages. |
| Sprint Index | `projects/ultraplan-go/sprints/01-cli-foundation/sprint-index.md` | Selected the Architecture, Errors, Testing, Documentation, and CLI Surface contracts; selected the evidence reports used by the handbook; excluded runtime, config, observability, security, workflows, persistence, performance, study behavior, code extraction behavior, target/sprint workflows, and packaging. |
| Technical Handbook | `projects/ultraplan-go/sprints/01-cli-foundation/technical-handbook.md` | Provided synthesized study evidence for thin entrypoints, manual composition roots, injectable IO, error/exit mapping, behavior tests, strict non-goals, and minimal scope. |
| Architecture Reasoning | `projects/ultraplan-go/sprints/01-cli-foundation/reasoning/architecture.md` | Area-specific reasoning concluded the sprint should use a minimal standard-library `internal/app` command shell, a process-only `cmd/ultraplan/main.go`, injected args/IO, usage exit code `2`, and documentation-only skeleton packages. |
| Selected Contracts | `system/contracts/core/architecture.md`, `system/contracts/core/errors.md`, `system/contracts/core/testing.md`, `system/contracts/core/documentation.md`, `system/contracts/surfaces/cli.md` | Mapped final decisions to concrete contract IDs for module boundaries, thin CLI entrypoints, error visibility, deterministic tests, public CLI docs, stdout/stderr discipline, and stable exit codes. |
| Selected Evidence Reports | `studies/go-cli-study/reports/final/{01,02,03,05,06,09,11,12,15}-*.md` | Grounded final decisions in concrete studied Go CLI patterns and rejected alternatives. |

## Area-Specific Reasoning Inputs

| Area | Reasoning Document | Key Conclusion | Evidence Basis | Impact On Final Decision |
| --- | --- | --- | --- | --- |
| Architecture | `projects/ultraplan-go/sprints/01-cli-foundation/reasoning/architecture.md` | Proceed with a minimal standard-library command shell in `internal/app`, keep `cmd/ultraplan/main.go` process-only, inject args/stdout/stderr for tests, map unknown commands to exit code `2`, and create only documented skeleton packages for future modules. | Technical handbook evidence for thin entrypoints, manual composition roots, injectable IO, error/exit mapping, behavior tests, and strict non-goals; sprint requirements; PRD/TRD/Architecture docs. | This reasoning accepts the architecture conclusion without override and turns it into the final sprint implementation constraints and review evidence. |

## Sprint Technical Handbook Summary

- **Relevant Patterns:** Thin process entrypoints, `internal/` application packages, visible manual composition root, injectable IO streams, user-facing error rendering with exit-code mapping, behavior-first command tests, and explicit non-goals as design control.
- **Important Trade-Offs:** Thin entrypoint adds one delegation hop; `internal/` avoids public SDK commitments; minimal command shell avoids framework complexity; injectable IO adds small plumbing; simple error/status values defer richer taxonomy; skeleton packages make future boundaries visible but must avoid over-promising behavior.
- **Warnings / Anti-Patterns:** Do not put command logic in `cmd/ultraplan/main.go`; do not hardcode process IO in app logic; do not advertise deferred commands; do not introduce global config/service singletons; do not create plugin/runtime/registry mechanics; do not test only implementation details.
- **Evidence Confidence:** High for thin entrypoint, internal boundaries, manual composition, IO injection, error/exit handling, and behavior tests because the selected reports studied 16 Go CLI repositories and cite concrete source references. Medium for terminal UX and extensibility because this sprint uses only minimal CLI-first output and deliberately defers plugin/runtime complexity.

## Contracts Applied

Requirement labels below are local references to `projects/ultraplan-go/sprints/01-cli-foundation/requirements.md` because that file does not define stable IDs.

| Contract / Requirement ID | Constraint | Decision Impact | Expected Evidence |
| --- | --- | --- | --- |
| REQ-GOAL | Establish a buildable Go module with a thin CLI shell, app composition root, version/help behavior, and initial package layout. | The sprint implements only bootstrap architecture and package boundaries. | Required files exist; `go test ./...`; `go build ./cmd/ultraplan`. |
| REQ-OUTPUT-GOMOD | `../ultraplan-go/go.mod` must define the module and toolchain version. | Create a normal Go module without adding unearned dependencies. | `go test ./...` and `go build ./cmd/ultraplan` from `../ultraplan-go`. |
| REQ-OUTPUT-MAIN | `cmd/ultraplan/main.go` must be a thin executable entrypoint. | `main.go` delegates to `internal/app` and exits with returned status. | Code review confirms no product workflow logic in `main.go`. |
| REQ-OUTPUT-APP | `internal/app/app.go` owns parsing, output routing, exit mapping, and wiring for the CLI shell. | Minimal app-owned dispatcher handles no args, `--help`, `-h`, `version`, and unknown commands. | Unit tests call `internal/app` with explicit args and in-memory streams. |
| REQ-OUTPUT-VERSION | `internal/app/version.go` provides version, commit, build date, and Go version values. | Version metadata is immutable app data, defaulted for local builds, and renderable by `ultraplan version`. | Version test asserts all required fields are present. |
| REQ-OUTPUT-TESTS | `internal/app/app_test.go` verifies help, version, invalid command, and exit codes offline. | Tests are behavior-oriented and do not invoke external runtimes or workspace state. | `go test ./...` passes with command output and exit-code assertions. |
| REQ-OUTPUT-SKELETONS | Required skeleton `doc.go` files exist for platform, workspace, study, and codeextract packages. | Create package docs only, with no workflow implementation. | File inspection and `go test ./...` confirm packages compile. |
| REQ-AC-HELP | `ultraplan --help` exits `0` and includes binary name, available commands, and `version`. | Help text is deterministic and rendered by one app path. | Tests assert exit code and stable help content. |
| REQ-AC-NOARGS | Running `ultraplan` with no args exits `0` and shows the same help-oriented discovery as `--help`. | No-args and help share rendering to avoid drift. | Test compares or verifies shared help output. |
| REQ-AC-VERSION | `ultraplan version` exits `0` and prints version, commit, build date, and Go version. | Version command has no workspace/runtime dependencies. | Test asserts all fields in stdout and status `0`. |
| REQ-AC-UNKNOWN | Unknown commands exit `2` and print an actionable error naming the unknown command. | Unknown command is classified as usage/argument error and rendered to stderr. | Test asserts status `2`, stderr text, and no false success. |
| REQ-AC-SCOPE-HELP | Help must not advertise target, sprint, workspace init, config, health, study, runtime, summary, validation, or code extraction commands. | Help surface includes only root discovery and `version` in this sprint. | Negative assertions in help tests and source review. |
| REQ-CONSTRAINT-OFFLINE | Tests must be deterministic and offline. | No OpenCode, provider credentials, network, configured workspace, or real filesystem state beyond test-owned temps. | `go test ./...` runs in a clean checkout. |
| ARCH-ENTRY-001 | CLI transport adapters must stay thin. | `cmd/ultraplan/main.go` contains process concerns only. | Code review of `cmd/ultraplan/main.go`. |
| ARCH-CORE-001 | Module boundaries must remain explicit. | Keep app/platform/workspace/study/codeextract as explicit `internal/` packages. | Package layout inspection. |
| ARCH-CORE-002 | Dependency direction must point inward. | `cmd/ultraplan` imports `internal/app`; platform skeletons import no product modules. | Import review and `go test ./...`. |
| ARCH-SHARED-001 | Shared/platform code must stay domain-neutral. | Platform skeletons document generic infrastructure only. | Package docs contain no study or code extraction behavior. |
| ERR-CORE-001 | No failure may disappear. | Unknown command returns explicit failure status and diagnostic. | Test covers status `2` and diagnostic text. |
| ERR-UI-001 | Caller-facing and operator-facing messages must be separated. | User-facing usage diagnostic stays actionable; no internal diagnostics are exposed. | Invalid-command test and source review. |
| TEST-SEAM-001 | Collaborators must be replaceable through public seams. | App receives args/stdout/stderr as inputs, not hardcoded globals. | Tests instantiate app with buffers. |
| TEST-FAIL-001 | Failure paths must be tested explicitly. | Unknown-command path has a negative test. | `internal/app/app_test.go` includes unknown command assertions. |
| TEST-DET-001 | Tests must remain deterministic and non-brittle. | Assert stable fields and command contracts without external dependencies. | `go test ./...` passes offline. |
| TEST-CONTRACT-001 | Public CLI contracts need compatibility tests. | Help, version, and usage-error behavior are command-level contract tests. | Tests cover stdout/stderr and exit codes. |
| DOC-ARCH-001 | Architecture changes need decision context. | This reasoning records context, alternatives, trade-offs, and consequences before implementation. | Completed `reasoning.md`. |
| DOC-PUBLIC-001 | Public CLI/package surfaces must be documented. | Help/version text and package `doc.go` files document current public CLI and internal package boundaries. | Help output tests and skeleton file inspection. |
| CLI-SHAPE-001 | Commands and flags must be predictable and discoverable. | Start with only root help and `version`; do not invent future command groups. | Help output review. |
| CLI-HELP-001 | Help/version output must exist and remain useful. | Implement `--help`, `-h`, no-args help, and `version`. | Help/version tests. |
| CLI-IO-001 | Stdout/stderr must be separated by purpose. | Successful help/version write to stdout; usage diagnostics write to stderr. | Tests assert output stream separation. |
| CLI-EXIT-001 | Exit codes must be stable and meaningful. | `0` for success, `2` for usage/argument error, and `1` only for unexpected app failures if surfaced. | Tests assert important success and failure exit codes. |
| CLI-NONINT-001 | Non-interactive mode must never hang on prompts. | No prompts or TTY requirements in this sprint. | Source review confirms no interactive code. |

## Repos Studied / Source Evidence Used

| Source / Repo / Report | Concrete Reference | Relevant Finding | Why It Matters For This Sprint | Used In Decision(s) |
| --- | --- | --- | --- | --- |
| `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md`; `cmd/gh/main.go:6`; `cmd/restic/main.go:162`; `cmd/root.go:112-127`; `internal/chezmoi/chezmoi.go:1-2` | Mature Go CLIs keep entrypoints thin and use `internal/` for application-owned implementation boundaries. | Directly supports `cmd/ultraplan -> internal/app` and private module-owned package layout. | Decisions 1, 2, 7 |
| `02-command-architecture` | `studies/go-cli-study/reports/final/02-command-architecture.md`; `pkg/cmd/install.go:132-145`; `pkg/cmdutil/factory.go:16-43`; `cmd/cmd.go:240-340` | Maintainable commands separate command wiring from business logic and centralize cross-cutting command behavior. | Supports app-owned command handling and rejection of command logic in `main.go`. | Decisions 2, 3 |
| `03-dependency-injection` | `studies/go-cli-study/reports/final/03-dependency-injection.md`; `internal/ghcmd/cmd.go:52-132`; `internal/app/app.go:42-81`; `cmd/restic/main.go:181-183` | Manual composition roots dominate; global mutable state is a recurring testability risk. | Supports a visible app composition root, injected streams/args, and immutable build metadata instead of mutable globals. | Decisions 1, 2, 5, 6 |
| `05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md`; `cmd/age/tui.go:37-54`; `internal/ghcmd/cmd.go:44-49`; `go-task/errors/errors_task.go:13-32` | Mature CLIs provide actionable diagnostics and stable exit-code mapping. | Supports unknown-command stderr diagnostics and exit code `2` without a full error taxonomy. | Decision 4 |
| `06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md`; `pkg/iostreams/iostreams.go:551-568`; `executor.go:553-564`; `ui.go:19-43` | Injectable IO streams enable deterministic command tests. | Supports making `internal/app` accept writers rather than using `os.Stdout`/`os.Stderr` directly. | Decisions 2, 6 |
| `09-terminal-ux` | `studies/go-cli-study/reports/final/09-terminal-ux.md`; `internal/cmd/prompt.go:124-137`; `lib/terminal/terminal.go:82-86`; `help.go:371-388` | CLI-first tools should keep output simple and script-friendly; rich TUI/progress is justified only by interaction density or long operations. | Supports simple deterministic text output and no prompt/TUI/progress behavior. | Decisions 3, 4 |
| `11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md`; `internal/cmd/main_test.go:64-174`; `go-task/task_test.go:166-169`; `gdu/cmd/gdu/app/app_test.go:682` | Strong CLI projects test output and exit behavior through command-level helpers or integration-style tests. | Supports direct `internal/app` command behavior tests over manual inspection. | Decision 6 |
| `12-extensibility` | `studies/go-cli-study/reports/final/12-extensibility.md`; `go-task/executor.go:20-24,91-122`; `helm/internal/plugin/runtime_subprocess.go:65-79`; `gh-cli/pkg/cmdutil/factory.go:16-43` | Extensibility should match scope; plugin and registry systems bring lifecycle, versioning, validation, and security costs. | Supports delaying CLI framework/plugin/runtime mechanics and creating only boundary docs. | Decisions 3, 7 |
| `15-philosophy` | `studies/go-cli-study/reports/final/15-philosophy.md`; `age.go:18`; `README.md:28-29`; `VISION.md:97` | Coherent tools reject features deliberately and accept complexity only when it serves the product. | Supports strict non-goals and a minimal shell. | Decisions 3, 7 |
| PRD | `projects/ultraplan-go/docs/PRD.md` | CLI bootstrap/help is a must-have; target and sprint workflows are deferred from the current study-side build. | Establishes product direction while constraining the sprint help surface. | Decisions 3, 7 |
| TRD | `projects/ultraplan-go/docs/TRD.md` | Defines Go layout, exit code classes, version command fields, build/test commands, and future agentwrap usage. | Supports Go module/bootstrap choices and explicit deferral of runtime integration. | Decisions 1, 3, 4, 5, 6, 7 |
| ARCHITECTURE | `projects/ultraplan-go/docs/ARCHITECTURE.md` | Product modules own behavior; platform packages are generic; expected dependency direction includes `cmd/ultraplan -> internal/app`. | Governs package layout and dependency constraints. | Decisions 1, 2, 7 |

## Trade-Off And Debt Analysis

### Accepted Trade-Offs

| Trade-Off | Benefit | Cost / Constraint Accepted | Why Acceptable Now | Revisit Trigger |
| --- | --- | --- | --- | --- |
| Minimal standard-library command shell instead of CLI framework | No external dependency, no premature framework conventions, easy output/exit tests. | Help rendering and simple dispatch are maintained locally. | Current command surface is only help/no-args/version/unknown. | Real nested study commands, global flags, command-specific help, JSON mode, or shell completion enter scope. |
| `cmd/ultraplan/main.go` as process-only delegator | Keeps process concerns separate and testable app behavior in `internal/app`. | Adds one delegation hop. | Required by sprint acceptance and supported by studied CLI structure. | None for this sprint; future change only if architecture docs are revised. |
| `internal/` packages only, no public `pkg/` API | Go-enforced application encapsulation and freedom to refactor. | External Go consumers cannot import these packages. | Product is a CLI application, not a public SDK. | A future explicit public SDK requirement appears. |
| Documentation-only skeleton packages | Makes intended module ownership visible early. | Empty packages can imply unearned architecture if docs overpromise behavior. | Required outputs demand skeleton files, while behavior is explicitly deferred. | Any skeleton starts defining concrete services, interfaces, schemas, or commands before its sprint. |
| Simple exit classification now | Satisfies script-friendly behavior with `0`, `2`, and possible `1` for unexpected failures. | Rich workspace/config/runtime validation categories are deferred. | Only unknown command failure is in scope. | Workspace, config, validation, cancellation, partial completion, or runtime failures enter scope. |
| Exact deterministic text for initial CLI output | Tests can protect the public CLI shell from accidental drift. | Future output changes require intentional test updates. | Initial help/version text is small and stable. | Complex dynamic help or localization enters scope. |

### Potential Technical Debt

| Debt / Shortcut | Why It Might Accrue | Current Mitigation | Owner / Follow-Up |
| --- | --- | --- | --- |
| Local dispatcher may not scale to nested study command tree | Future `study` commands may require richer flag parsing, subcommand help, shared global flags, and JSON modes. | Keep dispatcher minimal, tested, and isolated in `internal/app` so migration is contained. | Future study CLI sprint must reassess local dispatcher vs framework. |
| Static help text can drift from commands | Future commands may be added without updating discovery text or tests. | Help tests include positive and negative command surface assertions. | Every future CLI command sprint must update help tests. |
| Skeleton package docs may overstate future implementation | Boundary docs can become stale or imply committed APIs. | Keep docs limited to ownership and explicitly deferred behavior. | First sprint implementing each module must revise the package docs. |
| Version metadata defaults may not match release builds | Local defaults can be mistaken for release-grade metadata. | Keep release injection out of scope and document defaults as build metadata fields. | Release packaging sprint defines linker flags and reproducible version policy. |
| Minimal error shape lacks structured diagnostics | Later operational failures need stable codes and richer detail. | Restrict usage-error behavior to unknown command only. | Workspace/config/runtime/error taxonomy sprint expands error model. |

### Future Considerations

| Consideration | Deferred Until | Reason Deferred | What Should Be Preserved Now |
| --- | --- | --- | --- |
| CLI framework adoption | First real multi-command study workflow sprint | Current command set is too small to justify framework dependency. | Keep all command behavior in `internal/app` and tests process-independent. |
| Workspace/config loading | Workspace/config sprint | This sprint only creates a config skeleton and no workspace behavior. | Do not read cwd, env config, or config files in app shell. |
| Agentwrap/OpenCode runtime integration | Runtime sprint | Sprint requirements explicitly defer runtime execution and warn not to define competing runtime contracts. | `platform/runtime` doc must remain generic and must not define a product runtime API. |
| Study/domain commands | Study initialization/listing/run sprint | Study modeling and workflows are explicit non-goals. | `study` skeleton states ownership but exposes no commands or services. |
| Code extraction command | Code extraction sprint | Citation parsing and `ultraplan code` are explicit non-goals. | Help must not advertise `code`; `codeextract` remains documentation-only. |
| Target/sprint workflows | Requirements revision for target/sprint side | PRD/TRD/project index defer target/sprint execution from current study-side build. | Do not create target/sprint packages, commands, schemas, or validators. |
| Machine-readable output modes | Automation-facing command sprint | Help/version shell does not yet need JSON output. | Keep stdout/stderr discipline clear so JSON can be added cleanly later. |
| Packaging and release metadata | Release hardening sprint | Release artifacts, checksums, completions, and generated user docs are non-goals. | Version fields should be easy to set via linker flags later. |

## Final Decisions

### Decision 1: Create A Buildable Go Module With Internal Application Layout

- **Decision:** Create `../ultraplan-go/go.mod`, `cmd/ultraplan/main.go`, `internal/app`, `internal/platform/{config,logging,filesystem,runtime}`, `internal/workspace`, `internal/study`, and `internal/codeextract` as the initial module-owned layout. Use one package per module and do not create `pkg/`, global technical-layer packages, or clean-architecture subtrees.
- **Rationale:** This directly satisfies the sprint goal and aligns the initial codebase with the module-driven architecture before product workflows are added. The project is a CLI application rather than a public Go SDK, so `internal/` is the correct default boundary.
- **Study / Source Grounding:** `technical-handbook.md` relevant patterns for thin entrypoints and internal package boundaries; `01-project-structure.md` cites `cmd/gh/main.go:6`, `cmd/restic/main.go:162`, `cmd/root.go:112-127`, and `internal/chezmoi/chezmoi.go:1-2`; `ARCHITECTURE.md` lines 46-98 and 245-270 define the expected layout and dependency rules.
- **Trade-Offs Accepted:** Accept `internal/` encapsulation over external importability and one package per module over subpackage decomposition.
- **Technical Debt / Future Impact:** A future public SDK would require extraction to `pkg/`; future modules may split into subpackages only when concrete complexity justifies it.
- **Alternatives Rejected:** Public `pkg/` layout is rejected because no public library requirement exists. Global technical packages such as `internal/validation`, `internal/scheduler`, `internal/reports`, and `internal/prompts` are rejected because `ARCHITECTURE.md` assigns product behavior to owning modules.
- **Contracts Satisfied:** REQ-GOAL, REQ-OUTPUT-GOMOD, REQ-OUTPUT-SKELETONS, ARCH-CORE-001, ARCH-CORE-002, ARCH-SHARED-001, DOC-ARCH-001.
- **Evidence Required:** Required files exist; `go test ./...` passes; `go build ./cmd/ultraplan` passes; architecture review confirms import direction and package boundary conformance.

### Decision 2: Keep `cmd/ultraplan/main.go` Process-Only

- **Decision:** Implement `cmd/ultraplan/main.go` as a thin process entrypoint that gathers process args and process streams, calls `internal/app`, and exits with the returned status. It must not own command parsing, product workflow logic, output formatting, or tests-only behavior.
- **Rationale:** Process-only entry keeps command behavior testable without `os.Exit` and prevents the entrypoint from accumulating product logic at the boundary.
- **Study / Source Grounding:** `technical-handbook.md` trade-off on thin entrypoint; `01-project-structure.md` pattern 1 and practical tip to keep `cmd/*/main.go` small; `02-command-architecture.md` thin-delegate pattern; `architecture.md` area reasoning lines 102-109 and 570-590.
- **Trade-Offs Accepted:** Accept a small delegation hop from `main.go` into `internal/app` for testability and boundary clarity.
- **Technical Debt / Future Impact:** Future shared lifecycle setup must stay centralized in `internal/app` or a deliberate composition root, not drift into `main.go`.
- **Alternatives Rejected:** Command logic in `main.go` is rejected because it violates REQ-OUTPUT-MAIN and ARCH-ENTRY-001. Process-level tests as the only verification are rejected because they make behavior harder to isolate and slower to run.
- **Contracts Satisfied:** REQ-OUTPUT-MAIN, REQ-CONSTRAINT-OFFLINE, ARCH-ENTRY-001, TEST-SEAM-001, CLI-EXIT-001.
- **Evidence Required:** Code review confirms `main.go` only delegates and exits; app tests do not need to invoke the binary; `go build ./cmd/ultraplan` succeeds.

### Decision 3: Implement A Minimal Standard-Library App Shell

- **Decision:** Implement command behavior in `internal/app` with the Go standard library only. Support no arguments, `--help`, `-h`, and `version`. Treat unknown commands and unsupported non-help flags as usage errors. Do not add Cobra, urfave/cli, shell completion, command grouping, global flags, JSON mode, or generated docs in this sprint.
- **Rationale:** The sprint's command surface is intentionally tiny. A framework would add dependency and convention costs before nested study commands, global flags, or command-specific help have been selected.
- **Study / Source Grounding:** `technical-handbook.md` open question and trade-off on minimal custom command shell vs framework; `02-command-architecture.md` shows frameworks become valuable for larger command trees but warns against RunE bloat; `12-extensibility.md` warns that plugin/registry mechanics bring lifecycle, versioning, validation, and security costs; `15-philosophy.md` supports explicit non-goals through age, fzf, and lazygit evidence.
- **Trade-Offs Accepted:** Accept a small local dispatcher and local help rendering instead of framework-provided features.
- **Technical Debt / Future Impact:** Dispatcher may need refactoring or replacement when real study command hierarchy enters scope; keeping it isolated in `internal/app` reduces migration cost.
- **Alternatives Rejected:** Cobra or urfave/cli now is rejected as unearned complexity. A plugin/extension registry is rejected because runtime and extensibility are non-goals. A flat `flag` package parser with process-coupled output is rejected if it hardcodes `os.Stdout`/`os.Stderr` inside app behavior.
- **Contracts Satisfied:** REQ-OUTPUT-APP, REQ-AC-HELP, REQ-AC-NOARGS, REQ-AC-SCOPE-HELP, CLI-SHAPE-001, CLI-HELP-001, CLI-NONINT-001, DOC-PUBLIC-001.
- **Evidence Required:** Help test proves `--help` and no args return status `0`; help includes `ultraplan` and `version`; help excludes workspace/config/health/study/runtime/summary/validation/code/target/sprint commands; source review confirms no third-party CLI framework or deferred command surface.

### Decision 4: Use Script-Friendly Output And Minimal Exit-Code Mapping

- **Decision:** Write successful help and version output to stdout. Write unknown-command diagnostics to stderr. Return `0` for help/no-args/version success and `2` for usage/argument errors. Reserve `1` only for unexpected internal failures if the implementation surfaces writer or app errors.
- **Rationale:** This provides a script-friendly public CLI contract now, matches TRD exit code classes, and keeps stdout clean for successful command output.
- **Study / Source Grounding:** `technical-handbook.md` error-handling pattern and terminal UX pattern; `05-error-handling.md` cites `cmd/age/tui.go:37-54`, `internal/ghcmd/cmd.go:44-49`, and `go-task/errors/errors_task.go:13-32`; `09-terminal-ux.md` supports simple CLI-first output for low interaction density; `TRD.md` lines 223-247 define meaningful exit code classes.
- **Trade-Offs Accepted:** Accept a minimal status taxonomy rather than a full structured error system because only usage errors are in scope.
- **Technical Debt / Future Impact:** Future workspace/config/runtime commands must expand error categories and may need structured diagnostic types, but this sprint must not predefine them speculatively.
- **Alternatives Rejected:** Printing errors to stdout is rejected because CLI contract reserves stdout for intended command output. Returning `1` for unknown commands is rejected because requirements and TRD classify usage errors as exit code `2`. Implementing full error codes and machine-readable diagnostics is rejected because workspace/config/runtime failures are out of scope.
- **Contracts Satisfied:** REQ-AC-UNKNOWN, ERR-CORE-001, ERR-UI-001, CLI-IO-001, CLI-EXIT-001, TEST-FAIL-001.
- **Evidence Required:** Unknown-command test asserts status `2`, stderr names the unknown command and gives an actionable correction path, stdout is not contaminated by the error, and no success status is returned.

### Decision 5: Provide Version Metadata As Immutable App Data

- **Decision:** Implement `internal/app/version.go` with version, commit, build date, and Go version fields. Use sensible local-build defaults for version, commit, and build date, and use the Go runtime version for the Go version field. Render these fields through `ultraplan version` without workspace, config, network, or runtime dependencies.
- **Rationale:** The sprint requires a version command now, while release-time linker injection and packaging are later concerns. Keeping metadata immutable and app-owned avoids global mutable state and gives release tooling a clear future seam.
- **Study / Source Grounding:** `TRD.md` lines 1701-1715 define version output fields. `03-dependency-injection.md` accepts immutable build-time vars as a reasonable global exception while warning against config/service singletons. No deeper study source is needed for the exact default strings because those are local product policy, not an architectural pattern.
- **Trade-Offs Accepted:** Accept placeholder local defaults until a release sprint defines linker flags.
- **Technical Debt / Future Impact:** Release packaging must later define `-ldflags` or equivalent injection and verify version accuracy in release builds.
- **Alternatives Rejected:** Reading version data from files or Git at runtime is rejected because it adds filesystem/process dependencies to a shell command. Omitting commit/build date until releases is rejected because the sprint explicitly requires those fields.
- **Contracts Satisfied:** REQ-OUTPUT-VERSION, REQ-AC-VERSION, CLI-HELP-001, CLI-EXIT-001, TEST-CONTRACT-001.
- **Evidence Required:** Version test asserts status `0`, stdout contains `Version`, `Commit`, `BuildDate`, and `GoVersion` fields, and the command works without external state.

### Decision 6: Test Command Behavior Through `internal/app` With Injected IO

- **Decision:** Add deterministic offline tests in `internal/app/app_test.go` that instantiate the app with explicit argument slices and in-memory stdout/stderr buffers. Tests must cover no-args help, `--help`, `version`, unknown command, exit codes, stream separation, and negative assertions that deferred commands are not advertised.
- **Rationale:** Tests should verify the public command contract while avoiding process-level side effects, real terminals, network access, runtime tools, provider credentials, and workspace state.
- **Study / Source Grounding:** `technical-handbook.md` behavior-first tests and injectable IO patterns; `06-io-abstraction.md` cites `pkg/iostreams/iostreams.go:551-568`, `executor.go:553-564`, and `ui.go:19-43`; `11-testing-strategy.md` cites `internal/cmd/main_test.go:64-174`, `go-task/task_test.go:166-169`, and `gdu/cmd/gdu/app/app_test.go:682`.
- **Trade-Offs Accepted:** Use focused command-level tests rather than binary-spawning smoke tests because the sprint's core behavior is app-owned and simple.
- **Technical Debt / Future Impact:** Future richer command workflows may need testscript or binary-level tests, but current tests should remain fast and behavior-focused.
- **Alternatives Rejected:** Manual command inspection only is rejected because requirements demand tests. Process-only tests are rejected as the primary seam because `os.Exit` and real streams make negative behavior harder to verify. Implementation-detail tests of private branch functions are rejected because output and status are the public contract.
- **Contracts Satisfied:** REQ-OUTPUT-TESTS, REQ-CONSTRAINT-OFFLINE, TEST-SEAM-001, TEST-FAIL-001, TEST-DET-001, TEST-CONTRACT-001, CLI-IO-001, CLI-EXIT-001.
- **Evidence Required:** `go test ./...` passes offline; tests assert stdout/stderr content and status codes; no test requires OpenCode, credentials, network, or existing workspace files.

### Decision 7: Create Documentation-Only Skeleton Packages And Enforce Scope Boundaries

- **Decision:** Create required `doc.go` files for `internal/platform/config`, `internal/platform/logging`, `internal/platform/filesystem`, `internal/platform/runtime`, `internal/workspace`, `internal/study`, and `internal/codeextract`. Each doc should state package ownership and explicitly identify deferred behavior. Do not define runtime contracts, config schemas, workspace discovery, study services, code extraction parsers, validators, schedulers, targets, or sprint execution artifacts.
- **Rationale:** The sprint needs visible module boundaries for future study-side work, but behavior in these modules has not been selected and must not be architected prematurely.
- **Study / Source Grounding:** `ARCHITECTURE.md` module ownership and dependency rules; `technical-handbook.md` static skeleton package trade-off; `12-extensibility.md` warns against premature plugin/runtime/registry mechanics; `15-philosophy.md` supports strict explicit non-goals; `PRD.md` and `TRD.md` defer target/sprint workflows and runtime implementation from this sprint.
- **Trade-Offs Accepted:** Accept early package visibility while avoiding service/interface definitions until behavior is selected.
- **Technical Debt / Future Impact:** Package docs must be updated when each module receives real behavior; if docs overpromise, they can mislead future plans.
- **Alternatives Rejected:** Creating services or interfaces in skeleton packages is rejected as speculative. Adding target/sprint packages is rejected because PRD/TRD and sprint requirements defer that product side. Integrating agentwrap/OpenCode in `platform/runtime` now is rejected because runtime integration is explicitly out of scope and must later use `github.com/Antonio7098/agentwrap` and `agentwrap/opencode`.
- **Contracts Satisfied:** REQ-OUTPUT-SKELETONS, REQ-AC-SCOPE-HELP, ARCH-CORE-001, ARCH-SHARED-001, DOC-PUBLIC-001, CLI-SHAPE-001.
- **Evidence Required:** Required `doc.go` files exist; package docs mention boundaries without implementing behavior; help text does not advertise deferred commands; `go test ./...` confirms packages compile.

## Expected Evidence

| Evidence Type | Required Evidence | Source / Command / Review Check |
| --- | --- | --- |
| Tests | All Go tests pass without OpenCode, provider credentials, network access, or a configured UltraPlan workspace. | Run `go test ./...` from `../ultraplan-go`. |
| Build | CLI binary target builds successfully. | Run `go build ./cmd/ultraplan` from `../ultraplan-go`. |
| Help Runtime | Built binary with `--help` exits `0`, writes deterministic help to stdout, includes `ultraplan` and `version`, and excludes deferred command names. | Run built binary with `--help`; verify exit code and output. |
| No-Args Runtime | Built binary with no args exits `0` and shows the same help-oriented command discovery as `--help`. | Run built binary without arguments; verify exit code and output. |
| Version Runtime | Built binary with `version` exits `0` and prints version, commit, build date, and Go version fields. | Run built binary with `version`; verify exit code and output. |
| Invalid Command Runtime | Built binary with an unknown command exits `2` and prints an actionable diagnostic naming the command to stderr. | Run built binary with an unknown command; verify exit code and stderr. |
| Review | `cmd/ultraplan/main.go` is process-only and contains no product workflow logic. | Architecture review and source inspection. |
| Review | Required package skeletons exist and remain documentation-only. | Sprint review file inspection. |
| Review | Platform packages do not import product modules and no global technical-layer packages were created. | Architecture review import inspection. |
| Documentation | Package docs describe future boundaries without claiming implemented workspace/config/runtime/study/code extraction behavior. | Review `internal/**/doc.go`. |
| Scope Control | No workspace/config/health/study/runtime/summary/validation/code/target/sprint command is advertised or implemented. | Help test negative assertions and source tree review. |

## Assumptions And Risks

| Item | Type | Impact | Mitigation / Follow-Up |
| --- | --- | --- | --- |
| Go toolchain is available in `../ultraplan-go`. | Assumption | Build/test evidence cannot be produced if Go is absent. | Plan must run `go test ./...` and `go build ./cmd/ultraplan`; if unavailable, record blocker in review. |
| The future CLI command tree is not selected yet. | Assumption | A too-rich shell now could constrain later ergonomics. | Keep dispatcher minimal and isolated; reassess in first study command sprint. |
| Local version defaults are acceptable before release packaging. | Assumption | Users may see placeholder metadata in local builds. | Use clear default values; release sprint must define linker injection. |
| Minimal dispatcher may need replacement. | Risk | Future nested commands may outgrow local parsing. | Tests protect behavior; migration contained to `internal/app`. |
| Help text may drift as commands are added. | Risk | Users and scripts may see inaccurate discovery. | Add/maintain command surface tests whenever commands change. |
| Skeleton docs may overpromise. | Risk | Future implementers may treat docs as committed APIs. | Keep docs boundary-focused and update them when behavior enters scope. |
| Deferred runtime integration might be accidentally predesigned. | Risk | Could conflict with TRD requirement to use agentwrap later. | `platform/runtime` must remain documentation-only this sprint; plan/review must reject runtime interfaces or OpenCode integration. |
| Unknown command stderr wording could be over-specified. | Risk | Exact wording tests can become brittle. | Assert stable contract fields: command name, actionable help direction, exit code, and stderr. Full exact text is acceptable only if intentionally treated as public output. |
| Target/sprint prototype concepts appear in PRD/TRD as future concepts. | Risk | Implementer may add target/sprint packages or commands prematurely. | Requirements, sprint index, and this reasoning explicitly reject target/sprint implementation in this sprint. |

## Implementation Constraints

- Implement only the required outputs listed in sprint requirements.
- Use Go and the standard Go toolchain; required verification commands are `go test ./...` and `go build ./cmd/ultraplan` from `../ultraplan-go`.
- Keep `cmd/ultraplan/main.go` limited to process entry, app invocation, and process exit with the returned status.
- Keep command behavior, output routing, argument handling, and exit mapping in `internal/app`.
- Do not add a third-party CLI framework in this sprint.
- Do not add `github.com/Antonio7098/agentwrap`, `agentwrap/opencode`, OpenCode process execution, runtime health checks, runtime events, retry policy, permissions, or runtime contracts in this sprint.
- Do not add workspace discovery, config loading, config precedence, `init-workspace`, `config show`, health checks, or secret redaction behavior.
- Do not add study commands, study models, study initialization, source discovery, prompt composition, analysis runs, synthesis, summaries, run-loop state, or report validation.
- Do not add code citation parsing, code extraction behavior, or an `ultraplan code` command.
- Do not add target or sprint commands, packages, schemas, validators, persisted artifacts, or workflow behavior.
- Do not create global technical-layer packages such as `internal/validation`, `internal/scheduler`, `internal/reports`, or `internal/prompts`.
- Do not make platform packages import product modules.
- Successful help and version output must go to stdout; diagnostics and usage errors must go to stderr.
- Unknown commands must return exit code `2`; successful no-args/help/version must return exit code `0`.
- Tests must use in-memory streams and explicit args for app behavior.
- Tests must be deterministic and offline and must not rely on current working directory state outside test-owned temporary directories.
- Skeleton package `doc.go` files must describe boundaries and explicitly defer behavior, not define speculative services or interfaces.

## Plan Handoff

`plan.md` must execute these decisions. It must not invent architecture, scope, or decisions beyond this document.

The plan must carry forward:

- final decisions
- applicable contracts and requirement IDs
- expected evidence
- risks and assumptions
- required review protocols

## Phase Exit Criteria

- [x] Selected context was read and used.
- [x] Area-specific reasoning documents were completed where applicable or explicitly marked not applicable.
- [x] Area-specific reasoning conclusions are reflected or explicitly overridden in final decisions.
- [x] Contracts are explicitly mapped to decisions and expected evidence.
- [x] Final decisions are clear enough for `plan.md` to execute.
- [x] Expected evidence is specific and reviewable.
