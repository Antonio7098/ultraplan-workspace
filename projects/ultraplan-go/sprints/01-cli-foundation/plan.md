# Sprint Plan: CLI Foundation

> Project: `ultraplan-go`
> Sprint: `01-cli-foundation`
> Source: `reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/01-cli-foundation/requirements.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/sprints/01-cli-foundation/sprint-index.md`, `projects/ultraplan-go/sprints/01-cli-foundation/technical-handbook.md`, `projects/ultraplan-go/sprints/01-cli-foundation/reasoning/architecture.md`, `projects/ultraplan-go/sprints/01-cli-foundation/reasoning.md`, `templates/sprint-plan.md`

This plan executes `reasoning.md`. It must not invent architecture, scope, or decisions.

## Reasoning Source

- **Sprint Reasoning:** `projects/ultraplan-go/sprints/01-cli-foundation/reasoning.md`
- **Sprint Index:** `projects/ultraplan-go/sprints/01-cli-foundation/sprint-index.md`
- **Technical Handbook:** `projects/ultraplan-go/sprints/01-cli-foundation/technical-handbook.md`
- **Area Reasoning:** `projects/ultraplan-go/sprints/01-cli-foundation/reasoning/architecture.md`

## Sprint Status

- **Status:** `complete`
- **Owner:** `implementation agent`
- **Start Date:** `2026-05-30`
- **Completion Date:** `2026-05-30`

## Decisions To Execute

| Decision | Source Section | Execution Implication |
| --- | --- | --- |
| Create a buildable Go module with internal application layout | `reasoning.md#decision-1-create-a-buildable-go-module-with-internal-application-layout` | Create only the required module, entrypoint, app package, tests, and documented skeleton packages under `../ultraplan-go`; do not create `pkg/`, global technical-layer packages, target/sprint packages, or clean-architecture subtrees. |
| Keep `cmd/ultraplan/main.go` process-only | `reasoning.md#decision-2-keep-cmdultraplanmaingo-process-only` | `main.go` gathers process args/streams, calls `internal/app`, and exits with the returned status; command parsing and output formatting stay out of `main.go`. |
| Implement a minimal standard-library app shell | `reasoning.md#decision-3-implement-a-minimal-standard-library-app-shell` | Use no third-party CLI framework; support no args, `--help`, `-h`, `version`, and unknown command handling only. |
| Use script-friendly output and minimal exit-code mapping | `reasoning.md#decision-4-use-script-friendly-output-and-minimal-exit-code-mapping` | Successful help/version output goes to stdout; usage diagnostics go to stderr; success returns `0`, unknown commands return `2`, and unexpected app failures may return `1`. |
| Provide version metadata as immutable app data | `reasoning.md#decision-5-provide-version-metadata-as-immutable-app-data` | `internal/app/version.go` exposes version, commit, build date, and Go version values without runtime Git, filesystem, workspace, network, or provider dependencies. |
| Test command behavior through `internal/app` with injected IO | `reasoning.md#decision-6-test-command-behavior-through-internalapp-with-injected-io` | Tests instantiate app behavior with explicit args and in-memory stdout/stderr; they assert behavior, streams, and exit codes rather than invoking external runtimes. |
| Create documentation-only skeleton packages and enforce scope boundaries | `reasoning.md#decision-7-create-documentation-only-skeleton-packages-and-enforce-scope-boundaries` | Required `doc.go` files document module ownership and explicitly defer behavior; no config, workspace, study, runtime, code extraction, target, or sprint behavior is implemented. |

## Requirements / Contracts To Satisfy

| Contract / Requirement ID | Required Behavior | Evidence Planned |
| --- | --- | --- |
| REQ-GOAL | Establish a buildable Go module with a thin CLI shell, app composition root, version/help behavior, and initial package layout. | Required files exist; `go test ./...`; `go build ./cmd/ultraplan`. |
| REQ-OUTPUT-GOMOD | `../ultraplan-go/go.mod` defines the Go module and toolchain version. | Inspect `go.mod`; build and test from `../ultraplan-go`. |
| REQ-OUTPUT-MAIN, ARCH-ENTRY-001 | `cmd/ultraplan/main.go` is a thin process entrypoint. | Source review confirms only process entry, app call, and exit status handling. |
| REQ-OUTPUT-APP | `internal/app/app.go` owns parsing, output routing, exit mapping, and shell wiring. | App tests cover command behavior through app-owned seams. |
| REQ-OUTPUT-VERSION, REQ-AC-VERSION | `ultraplan version` prints version, commit, build date, and Go version and exits `0`. | Unit test asserts fields and status; runtime review can run built binary. |
| REQ-OUTPUT-TESTS, TEST-SEAM-001, TEST-FAIL-001, TEST-DET-001, TEST-CONTRACT-001 | Command tests are deterministic, offline, and behavior-focused. | `internal/app/app_test.go`; `go test ./...` with no OpenCode, credentials, network, or workspace. |
| REQ-OUTPUT-SKELETONS, ARCH-CORE-001, ARCH-CORE-002, ARCH-SHARED-001 | Required skeleton package files exist and preserve product/platform dependency rules. | File inspection; import review; `go test ./...`. |
| REQ-AC-HELP, REQ-AC-NOARGS, CLI-SHAPE-001, CLI-HELP-001 | `--help`, `-h`, and no args show deterministic command discovery including `ultraplan` and `version`, and exit `0`. | Tests assert stdout content and status `0`; optional binary checks. |
| REQ-AC-UNKNOWN, ERR-CORE-001, ERR-UI-001, CLI-IO-001, CLI-EXIT-001 | Unknown command exits `2`, names the unknown command, is actionable, and writes diagnostics to stderr. | Tests assert status `2`, stderr content, stdout separation. |
| REQ-AC-SCOPE-HELP, CLI-NONINT-001 | Help does not advertise deferred workspace/config/health/study/runtime/summary/validation/code/target/sprint commands and never prompts. | Negative help assertions and source review. |
| REQ-CONSTRAINT-OFFLINE | Normal build and tests require only the Go toolchain. | `go test ./...` and `go build ./cmd/ultraplan` from `../ultraplan-go`. |
| DOC-ARCH-001, DOC-PUBLIC-001 | Decisions and current package/CLI surfaces are documented. | Completed `reasoning.md`, package `doc.go` files, help/version output tests. |

## Tasks

- [x] **Task 1: Bootstrap Go Module And Layout**
  > Executes: `Decision 1`, `REQ-GOAL`, `REQ-OUTPUT-GOMOD`, `REQ-OUTPUT-SKELETONS`, `ARCH-CORE-001`, `ARCH-CORE-002`
  - [x] Create `../ultraplan-go/go.mod` with the selected module name and Go toolchain version.
  - [x] Create required directories for `cmd/ultraplan`, `internal/app`, `internal/platform/{config,logging,filesystem,runtime}`, `internal/workspace`, `internal/study`, and `internal/codeextract`.
  - [x] Do not create `pkg/`, `internal/validation`, `internal/scheduler`, `internal/reports`, `internal/prompts`, target packages, or sprint packages.

- [x] **Task 2: Implement Process-Only Entrypoint**
  > Executes: `Decision 2`, `REQ-OUTPUT-MAIN`, `ARCH-ENTRY-001`, `TEST-SEAM-001`
  - [x] Implement `cmd/ultraplan/main.go` so it passes `os.Args[1:]`, `os.Stdout`, `os.Stderr`, and version metadata into `internal/app`.
  - [x] Ensure `main.go` calls `os.Exit` with the returned status.
  - [x] Keep command parsing, help text, version rendering, and product workflow logic out of `main.go`.

- [x] **Task 3: Implement Minimal App Shell**
  > Executes: `Decision 3`, `Decision 4`, `REQ-OUTPUT-APP`, `REQ-AC-HELP`, `REQ-AC-NOARGS`, `REQ-AC-UNKNOWN`, `CLI-SHAPE-001`, `CLI-IO-001`, `CLI-EXIT-001`
  - [x] Add an `internal/app` API that accepts explicit args and output writers and returns an exit status without calling `os.Exit`.
  - [x] Route no args, `--help`, and `-h` through one deterministic help rendering path.
  - [x] Implement only the `version` command and unknown-command usage errors.
  - [x] Write help and version output to stdout.
  - [x] Write unknown-command diagnostics to stderr with the invalid command name and a help-oriented next action.
  - [x] Return status `0` for successful help/version and status `2` for usage/argument errors.
  - [x] Do not add workspace, config, health, study, runtime, summary, validation, code extraction, target, or sprint command handling.

- [x] **Task 4: Add Version Metadata**
  > Executes: `Decision 5`, `REQ-OUTPUT-VERSION`, `REQ-AC-VERSION`, `TEST-CONTRACT-001`
  - [x] Add `internal/app/version.go` with version, commit, build date, and Go version fields.
  - [x] Use local-build defaults for version, commit, and build date until a release sprint defines linker injection.
  - [x] Use the Go runtime version for the Go version field.
  - [x] Do not read Git, files, workspace config, network, runtime tools, or provider state to render version output.

- [x] **Task 5: Add Documentation-Only Package Skeletons**
  > Executes: `Decision 7`, `REQ-OUTPUT-SKELETONS`, `ARCH-SHARED-001`, `DOC-PUBLIC-001`
  - [x] Add `doc.go` for `internal/platform/config` that documents the future config boundary and states config behavior is deferred.
  - [x] Add `doc.go` for `internal/platform/logging` that documents the future logging boundary and states structured logging is deferred.
  - [x] Add `doc.go` for `internal/platform/filesystem` that documents the future filesystem boundary and states workspace path logic is deferred.
  - [x] Add `doc.go` for `internal/platform/runtime` that documents the generic runtime boundary and states agentwrap/OpenCode integration is deferred.
  - [x] Add `doc.go` for `internal/workspace` that documents the workspace boundary and states discovery/validation behavior is deferred.
  - [x] Add `doc.go` for `internal/study` that documents the study boundary and states study workflows are deferred.
  - [x] Add `doc.go` for `internal/codeextract` that documents the code extraction boundary and states citation parsing/extraction behavior is deferred.
  - [x] Do not define services, interfaces, schemas, validators, runtime contracts, or workflow functions in skeleton packages.

- [x] **Task 6: Add Deterministic Command Tests**
  > Executes: `Decision 6`, `REQ-OUTPUT-TESTS`, `REQ-CONSTRAINT-OFFLINE`, `TEST-SEAM-001`, `TEST-FAIL-001`, `TEST-DET-001`, `TEST-CONTRACT-001`
  - [x] Add tests for `--help` and no-args behavior with status `0`, stdout help text, and no stderr diagnostics.
  - [x] Add tests for `version` with status `0`, required metadata fields, and no runtime/workspace dependency.
  - [x] Add tests for unknown command behavior with status `2`, stderr diagnostic naming the command, actionable help direction, and stdout separation.
  - [x] Add negative help assertions that deferred command names are not advertised.
  - [x] Use explicit args and in-memory stdout/stderr buffers; do not require OpenCode, credentials, network, or a configured workspace.

- [x] **Task 7: Verify And Review Scope**
  > Executes: `Expected Evidence`, `Review Expectations`, `Architecture Review`, `Sprint Review`
  - [x] Run `go test ./...` from `../ultraplan-go`.
  - [x] Run `go build ./cmd/ultraplan` from `../ultraplan-go`.
  - [x] Optionally run the built binary with no args, `--help`, `version`, and an unknown command to collect runtime review evidence.
  - [x] Inspect `cmd/ultraplan/main.go` for thin-entrypoint conformance.
  - [x] Inspect required package skeletons for documentation-only scope.
  - [x] Confirm help output and source tree do not include deferred commands or behavior.

## Evidence Checklist

- [x] Tests prove help, no-args, version, unknown-command, stream separation, exit codes, and deferred-command exclusion.
- [x] Runtime or diagnostic evidence exists for built binary help/version/invalid-command behavior where review requires it.
- [x] Documentation-only package skeletons are complete and do not imply implemented deferred behavior.
- [x] Deviations from `reasoning.md` are recorded before implementation continues.
- [x] Required review protocols have evidence for architecture review and sprint review.
- [x] Omitted evidence, if any, is recorded with the reason it was omitted.

## Verification Commands

| Check | Command | Expected Result |
| --- | --- | --- |
| Unit and package tests | `go test ./...` from `../ultraplan-go` | Succeeds without OpenCode, provider credentials, network access, or configured workspace. |
| CLI build | `go build ./cmd/ultraplan` from `../ultraplan-go` | Succeeds and produces a runnable CLI binary. |
| Help behavior | `./ultraplan --help` or equivalent built binary path | Exits `0`; writes deterministic help to stdout; includes `ultraplan` and `version`; excludes deferred command names. |
| No-args behavior | `./ultraplan` or equivalent built binary path | Exits `0`; shows the same help-oriented command discovery as `--help`. |
| Version behavior | `./ultraplan version` or equivalent built binary path | Exits `0`; prints version, commit, build date, and Go version fields. |
| Unknown command behavior | `./ultraplan definitely-unknown` or equivalent built binary path | Exits `2`; writes actionable stderr diagnostic naming `definitely-unknown`; does not report success on stdout. |

## Risks And Blockers

| Risk / Blocker | Source | Mitigation | Status |
| --- | --- | --- | --- |
| Go toolchain may be unavailable in `../ultraplan-go`. | `reasoning.md#assumptions-and-risks` | `go test ./...` passed. Exact `go build ./cmd/ultraplan` was attempted and failed only because this directory is not a Git repository for VCS stamping; `go build -buildvcs=false ./cmd/ultraplan` passed. | mitigated |
| Minimal dispatcher may need replacement when nested study commands enter scope. | `reasoning.md#decision-3-implement-a-minimal-standard-library-app-shell` | Dispatcher is isolated in `internal/app` and covered by behavior tests. | mitigated |
| Static help text may drift from implemented commands. | `reasoning.md#trade-off-and-debt-analysis` | Tests assert command discovery and deferred-command exclusion. | mitigated |
| Skeleton package docs may overpromise future behavior. | `reasoning.md#decision-7-create-documentation-only-skeleton-packages-and-enforce-scope-boundaries` | Docs are limited to ownership and explicit deferrals; no services/interfaces were added. | mitigated |
| Version metadata defaults may be mistaken for release-grade data. | `reasoning.md#decision-5-provide-version-metadata-as-immutable-app-data` | Defaults are clearly local values and release linker injection remains deferred. | carried forward |
| Deferred runtime integration may be accidentally predesigned. | `reasoning.md#implementation-constraints` | `platform/runtime` is documentation-only; no agentwrap/OpenCode dependency or runtime contract was added. | mitigated |
| Target/sprint prototype concepts may leak into this sprint. | `reasoning.md#assumptions-and-risks` | No target/sprint packages, commands, schemas, validators, persisted artifacts, or workflow behavior were added. | mitigated |

## Review Inputs

Review should use:

- `projects/ultraplan-go/sprints/01-cli-foundation/sprint-index.md`
- `projects/ultraplan-go/sprints/01-cli-foundation/technical-handbook.md`
- `projects/ultraplan-go/sprints/01-cli-foundation/reasoning/architecture.md`
- `projects/ultraplan-go/sprints/01-cli-foundation/reasoning.md`
- `projects/ultraplan-go/sprints/01-cli-foundation/plan.md`
- implementation diff
- verification evidence
- required protocols from `sprint-index.md`: `system/protocols/architecture-review-protocol.md` and `system/protocols/sprint-review-protocol.md`

## Execution Log

| Date / Step | Action | Evidence / Notes |
| --- | --- | --- |
| 2026-05-30 / planning | Created evidence-grounded implementation plan from `reasoning.md`; no implementation performed. | Plan carries forward decisions, evidence expectations, risks, assumptions, and open questions from `reasoning.md`. |
| 2026-05-30 / implementation | Added Go module, process-only CLI entrypoint, app command shell, version metadata, deterministic app tests, and documentation-only package skeletons. | Files added under `go.mod`, `cmd/ultraplan`, and `internal/{app,platform,workspace,study,codeextract}`. |
| 2026-05-30 / verification | Ran required tests and build/runtime checks. | `go test ./...` passed. `go build ./cmd/ultraplan` failed due non-Git VCS stamping (`error obtaining VCS status`); `go build -buildvcs=false ./cmd/ultraplan` passed. Runtime checks for `--help`, no args, `version`, and `definitely-unknown` passed with expected exit codes and streams. |
| 2026-05-30 / omitted evidence | Project roadmap input could not be loaded. | `/home/antonioborgerees/coding/ultra/projects/ultraplan-go/roadmap.md` does not exist. Residual risk is low because sprint scope is fully defined by requirements, reasoning, docs, and plan. |

## Completion Criteria

- [x] All tasks are complete or explicitly deferred.
- [x] Verification commands were run or deferrals are documented.
- [x] Evidence satisfies the expectations from `reasoning.md`.
- [x] `review.md` can evaluate conformance without guessing intent.
