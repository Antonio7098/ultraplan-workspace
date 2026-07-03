# Sprint Requirements: Workspace, Config, Logging, and Health Skeleton

> Project: `ultraplan-go`
> Sprint: `02-workspace-config-health`
> Purpose: the authoritative, human-readable sprint contract. All other sprint artifacts must satisfy these requirements.

## Sprint Goal

Make UltraPlan able to find, initialize, inspect, and validate a workspace — establishing workspace discovery, config loading, logging foundation, and health command skeleton.

## Required Outputs

| Output | Path | Description |
| --- | --- | --- |
| Requirements | `projects/ultraplan-go/sprints/02-workspace-config-health/requirements.md` | This document: sprint contract and acceptance criteria. |
| Sprint Index | `projects/ultraplan-go/sprints/02-workspace-config-health/sprint-index.md` | Selected contracts and go-cli-study evidence for workspace, config, logging, CLI output, and health behavior. |
| Technical Handbook | `projects/ultraplan-go/sprints/02-workspace-config-health/technical-handbook.md` | Implementation guidance for workspace discovery, config precedence, redaction, output modes, health checks, and package boundaries. |
| Architecture Reasoning | `projects/ultraplan-go/sprints/02-workspace-config-health/reasoning/architecture.md` | Reasoning for package ownership, dependency direction, and workspace/config/health composition. |
| Consolidated Reasoning | `projects/ultraplan-go/sprints/02-workspace-config-health/reasoning.md` | Final sprint reasoning that maps decisions to requirements and selected evidence. |
| Implementation Plan | `projects/ultraplan-go/sprints/02-workspace-config-health/plan.md` | Ordered implementation tasks, tests, and verification commands. |
| Workspace Package | `../ultraplan-go/internal/workspace/discovery.go` | Workspace root discovery from explicit path, environment, current directory, or parent directories. |
| Workspace Paths Package | `../ultraplan-go/internal/workspace/paths.go` | Workspace-managed path normalization and canonical workspace paths. |
| Workspace Validation Package | `../ultraplan-go/internal/workspace/validation.go` | Structural workspace validation used by `health` and workspace-sensitive commands. |
| Workspace Initialization Package | `../ultraplan-go/internal/workspace/init.go` | Workspace scaffold creation and dry-run planning for `init-workspace`. |
| Config Package | `../ultraplan-go/internal/platform/config/config.go` | Effective config model, defaults, workspace config loading, environment overrides, CLI overrides, validation, and redaction. |
| Config Redaction Package | `../ultraplan-go/internal/platform/config/redaction.go` | Sensitive-field detection and redacted config projection for CLI output and diagnostics. |
| Logging Package | `../ultraplan-go/internal/platform/logging/logging.go` | Minimal text/JSON logging foundation with secret-safe fields. |
| Output Package | `../ultraplan-go/internal/app/output.go` | Deterministic text and JSON command output helpers. |
| Workspace Commands | `../ultraplan-go/internal/app/workspace_commands.go` | `init-workspace` command wiring and behavior. |
| Config Commands | `../ultraplan-go/internal/app/config_commands.go` | `config show` command wiring and behavior. |
| Health Commands | `../ultraplan-go/internal/app/health_commands.go` | `health` command wiring and basic workspace/config/environment checks. |
| Workspace Command Tests | `../ultraplan-go/internal/app/workspace_commands_test.go` | Command-level tests for workspace initialization, dry-run behavior, and invalid path handling. |
| Config Command Tests | `../ultraplan-go/internal/app/config_commands_test.go` | Command-level tests for config text/JSON output, precedence, and redaction. |
| Health Command Tests | `../ultraplan-go/internal/app/health_commands_test.go` | Command-level tests for health text/JSON output and missing workspace/config cases. |
| Workspace Discovery Tests | `../ultraplan-go/internal/workspace/discovery_test.go` | Unit tests for workspace discovery precedence and parent traversal. |
| Workspace Path Tests | `../ultraplan-go/internal/workspace/paths_test.go` | Unit tests for path normalization and workspace escape rejection. |
| Workspace Validation Tests | `../ultraplan-go/internal/workspace/validation_test.go` | Unit tests for required workspace structure validation. |
| Workspace Initialization Tests | `../ultraplan-go/internal/workspace/init_test.go` | Unit tests for scaffold planning, creation, idempotence, and dry-run behavior. |
| Config Tests | `../ultraplan-go/internal/platform/config/config_test.go` | Unit tests for defaults, workspace config loading, environment overrides, CLI overrides, and validation. |
| Config Redaction Tests | `../ultraplan-go/internal/platform/config/redaction_test.go` | Unit tests proving sensitive config values are redacted from output projections. |

## Acceptance Criteria

- [ ] `go test ./...` succeeds from `../ultraplan-go` without requiring OpenCode, provider credentials, network access, or a configured UltraPlan workspace.
- [ ] `go build ./cmd/ultraplan` succeeds and produces a runnable CLI binary.
- [ ] `ultraplan init-workspace [--path <dir>] [--dry-run]` creates the required workspace structure if not present: `ultraplan.yml`, `prompts/base.md`, `prompts/synthesize.md`, `templates/repo-analysis.md`, `templates/report.md`, and `studies/`.
- [ ] `ultraplan config show [--json]` prints effective config with proper precedence and redacted sensitive values.
- [ ] `ultraplan health [--json]` reports basic workspace, config, filesystem, and environment status without requiring OpenCode or provider credentials.
- [ ] Config precedence is respected: built-in defaults < workspace config < environment variables < CLI flags.
- [ ] Secret redaction is implemented for sensitive config values.
- [ ] Text and JSON output modes are deterministic and supported for config and health commands.
- [ ] CLI help output includes init-workspace, config, and health commands.
- [ ] Exit codes are deterministic and meaningful per error class.

## Non-Goals

- Workspace validation beyond structural workspace/config/filesystem checks is not included.
- Config hot-reload is not included.
- Persistent workspace state beyond the initial scaffold is not included.
- Study listing, study initialization, source discovery, and analysis runs are not included.
- Runtime execution, agentwrap integration, OpenCode, and agent health checks are not included.
- Durable file logs, full structured log level policy, run event records, and observability event persistence are not included.
- Target scaffolding, sprint planning, and sprint execution workflows are not included.

## Constraints

- Use Go and the standard Go toolchain.
- Workspace discovery must follow the TRD precedence: explicit `--workspace`, `ULTRAPLAN_WORKSPACE`, current directory, then nearest parent directory.
- Config loading must follow the TRD precedence chain: built-in defaults → workspace config → environment variables → CLI flags.
- Secret values must be redacted in output, never logged or exposed.
- Keep `cmd/ultraplan/main.go` thin; all workspace/config/health logic lives under `internal/`.
- Platform packages must not import product modules.
- Tests must be deterministic and offline.
- CLI output must be script-friendly: deterministic text, JSON mode available, meaningful exit codes.

## Dependencies

| Prior Sprint / Output | Required For | Notes |
| --- | --- | --- |
| Sprint 1: CLI Foundation | Module layout, command routing baseline | Sprint 1 establishes the Go module, `cmd/ultraplan`, `internal/app`, and initial package layout that Sprint 2 extends with workspace/config behavior. |
| Architecture contract | Module ownership, dependency direction | Module-driven layout and thin entrypoints from architecture contract. |
| Errors contract | Exit codes, error wrapping, diagnostics | Exit code classes and error taxonomy from errors contract. |
| Configuration contract | Config precedence, redaction, validation | Config precedence chain and secret redaction requirements from configuration contract. |
| CLI Surface contract | Command shape, flags, help, output modes | Command surface expectations for init-workspace, config, health. |
| go-cli-study reports 01-15 | Evidence for implementation | Selected reports referenced in sprint-index.md. |

## Review Expectations

| What | How Verified |
| --- | --- |
| Buildability | Run `go build ./cmd/ultraplan` from `../ultraplan-go`. |
| Test pass | Run `go test ./...` from `../ultraplan-go`. |
| init-workspace command | Run with `--dry-run` first, then without; verify required workspace files and directories are created only inside the selected workspace path. |
| config show command | Run with text and `--json` output; verify defaults, workspace config, environment variables, and CLI flags resolve in the required order. |
| health command | Run with text and `--json` output; verify basic workspace/config/filesystem status reporting and no OpenCode/provider dependency. |
| Secret redaction | Verify sensitive values do not appear in output. |
| Exit codes | Run commands with invalid args or missing state; verify non-zero exit codes. |
| Scope control | Verify no study, target, sprint, runtime, or code extraction behavior was implemented. |
