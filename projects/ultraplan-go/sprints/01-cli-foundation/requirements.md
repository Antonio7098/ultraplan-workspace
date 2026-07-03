# Sprint Requirements: CLI Foundation

> Project: `ultraplan-go`
> Sprint: `01-cli-foundation`
> Purpose: the authoritative, human-readable sprint contract. All other sprint artifacts must satisfy these requirements.

## Sprint Goal

Establish a buildable Go module with a thin `ultraplan` CLI shell, app composition root, version/help behavior, and the initial module-owned package layout for future study-side features.

## Required Outputs

| Output | Path | Description |
| --- | --- | --- |
| Go module definition | `../ultraplan-go/go.mod` | Defines the `ultraplan-go` Go module and Go toolchain version. |
| CLI entrypoint | `../ultraplan-go/cmd/ultraplan/main.go` | Thin executable entrypoint that delegates command handling to `internal/app`. |
| App composition root | `../ultraplan-go/internal/app/app.go` | Owns initial command parsing, output routing, exit code mapping, and wiring for the CLI shell. |
| Version metadata | `../ultraplan-go/internal/app/version.go` | Provides version, commit, build date, and Go version values used by `ultraplan version`. |
| App command tests | `../ultraplan-go/internal/app/app_test.go` | Verifies help output, version output, invalid command behavior, and exit codes without requiring external runtimes. |
| Config package skeleton | `../ultraplan-go/internal/platform/config/doc.go` | Documents the future config package boundary without implementing workspace config behavior. |
| Logging package skeleton | `../ultraplan-go/internal/platform/logging/doc.go` | Documents the future logging package boundary without implementing structured logging. |
| Filesystem package skeleton | `../ultraplan-go/internal/platform/filesystem/doc.go` | Documents the future filesystem package boundary without implementing workspace path logic. |
| Runtime package skeleton | `../ultraplan-go/internal/platform/runtime/doc.go` | Documents the generic runtime boundary without integrating agentwrap or OpenCode yet. |
| Workspace package skeleton | `../ultraplan-go/internal/workspace/doc.go` | Documents the workspace module boundary without implementing discovery or validation. |
| Study package skeleton | `../ultraplan-go/internal/study/doc.go` | Documents the study module boundary without implementing study initialization, listing, prompts, or runs. |
| Code extraction package skeleton | `../ultraplan-go/internal/codeextract/doc.go` | Documents the code extraction module boundary without implementing citation parsing or extraction. |

## Acceptance Criteria

- [ ] `go test ./...` succeeds from `../ultraplan-go` without requiring OpenCode, provider credentials, network access, or a configured UltraPlan workspace.
- [ ] `go build ./cmd/ultraplan` succeeds from `../ultraplan-go` and produces a runnable CLI binary.
- [ ] Running `ultraplan --help` exits `0` and prints deterministic usage text that includes the binary name, available commands, and the `version` command.
- [ ] Running `ultraplan version` exits `0` and prints version, commit, build date, and Go version fields.
- [ ] Running `ultraplan` with no arguments exits `0` and shows the same help-oriented command discovery as `ultraplan --help`.
- [ ] Running an unknown command exits with code `2` and prints an actionable error naming the unknown command.
- [ ] `cmd/ultraplan/main.go` contains no product workflow logic beyond process entry, calling `internal/app`, and exiting with the returned status.
- [ ] The required initial package files exist under `internal/app`, `internal/platform/{config,logging,filesystem,runtime}`, `internal/workspace`, `internal/study`, and `internal/codeextract`.
- [ ] The CLI help output does not advertise target, sprint, workspace initialization, config, health, study, runtime, summary, validation, or code extraction commands.
- [ ] Tests assert the implemented command behavior rather than relying only on manual command output inspection.

## Non-Goals

- Workspace discovery, `init-workspace`, `config show`, config precedence, secret redaction, and `health` are not included in this sprint.
- Study domain modeling, study listing, study initialization from YAML, source discovery, dimension parsing, prompt composition, analysis runs, synthesis, summaries, and validation are not included.
- Agentwrap, OpenCode, runtime execution, runtime health checks, event handling, retry/fallback, policy mapping, and permission mapping are not included.
- Code-reference parsing, citation resolution, snippet extraction, and `ultraplan code` are not included.
- Target scaffolding, sprint planning, sprint execution, target commands, and sprint commands are not included.
- Packaging, release artifacts, checksums, shell completion, generated user documentation, and real-runtime smoke tests are not included.

## Constraints

- Use Go and the standard Go toolchain; normal test and build commands for this sprint must be `go test ./...` and `go build ./cmd/ultraplan`.
- Keep `cmd/ultraplan/main.go` thin; command behavior and process-independent testing seams must live under `internal/app`.
- Follow the module-driven architecture: product behavior belongs to owning modules, platform packages provide only generic infrastructure, and platform packages must not import product modules.
- Start with one package per module and focused files; do not introduce clean-architecture subtrees or global packages such as `internal/validation`, `internal/scheduler`, `internal/reports`, or `internal/prompts`.
- Do not add target or sprint workflows, commands, schemas, validators, or persisted artifacts; current release scope is study-side only.
- Do not integrate OpenCode directly or define a competing runtime contract; runtime integration is deferred and, when implemented later, must use `github.com/Antonio7098/agentwrap` and `agentwrap/opencode`.
- Keep tests deterministic and offline; they must not depend on filesystem state outside the test process except temporary directories created by tests.
- CLI output must be script-friendly: deterministic text, no interactive prompts, and meaningful exit codes.

## Dependencies

| Prior Sprint / Output | Required For | Notes |
| --- | --- | --- |
| Phase 0 roadmap and scope alignment | Sprint scope boundary | Sprint 1 assumes the study-side-only roadmap is accepted and target/sprint workflows remain deferred. |
| `projects/ultraplan-go/docs/PRD.md` | Product command and non-goal alignment | Confirms CLI bootstrap, help behavior, local-first CLI operation, and deferred target/sprint scope. |
| `projects/ultraplan-go/docs/TRD.md` | Technical constraints and acceptance shape | Defines Go layout, CLI expectations, exit code classes, build/test requirements, and future agentwrap boundary. |
| `projects/ultraplan-go/docs/ARCHITECTURE.md` | Package ownership and dependency rules | Defines module-driven layout, thin entrypoint, app composition root, and platform/product dependency direction. |
| Go toolchain | Build and test verification | Required locally for `go test ./...` and `go build ./cmd/ultraplan`. |
| None | Prior sprint carry-forward decisions | No prior sprint review exists for this first implementation sprint. |

## Review Expectations

| What | How Verified |
| --- | --- |
| Buildability | Run `go build ./cmd/ultraplan` from `../ultraplan-go`. |
| Test pass | Run `go test ./...` from `../ultraplan-go`. |
| Help command behavior | Run built binary with `--help` and no arguments; verify exit code `0` and expected usage content. |
| Version command behavior | Run built binary with `version`; verify exit code `0` and required metadata fields. |
| Invalid command behavior | Run built binary with an unknown command; verify exit code `2` and actionable error text. |
| Package layout | Inspect required output paths and confirm the package names and boundaries match the architecture document. |
| Thin entrypoint | Review `cmd/ultraplan/main.go` and confirm it only delegates to `internal/app` and exits with the returned status. |
| Scope control | Review help text and source tree to confirm no workspace/config/health, study workflow, runtime, code extraction, target, or sprint behavior was implemented. |
| Offline determinism | Confirm tests and commands used for review do not require OpenCode, provider credentials, network access, or a pre-existing workspace. |
