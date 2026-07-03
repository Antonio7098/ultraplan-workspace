# Sprint Technical Handbook: Workspace, Config, Logging, and Health Skeleton

> Project: `ultraplan-go`
> Sprint: `02-workspace-config-health`
> Source: `sprint-index.md`
> **Inputs Used:** `projects/ultraplan-go/sprints/02-workspace-config-health/sprint-index.md`, `templates/technical-handbook.md`, `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/02-workspace-config-health/requirements.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `studies/go-cli-study/reports/final/01-project-structure.md`, `studies/go-cli-study/reports/final/02-command-architecture.md`, `studies/go-cli-study/reports/final/03-dependency-injection.md`, `studies/go-cli-study/reports/final/04-configuration-management.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/06-io-abstraction.md`, `studies/go-cli-study/reports/final/09-terminal-ux.md`, `studies/go-cli-study/reports/final/10-logging-observability.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/13-security.md`

This handbook distills the studies and reports selected by `sprint-index.md` for sprint reasoning. It does not decide architecture or implementation.

## Selected Studies And Reports

| Study / Report | Path | Relevant Finding | Confidence |
| --- | --- | --- | --- |
| `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md` | Mature Go CLIs keep entrypoints thin and move business logic behind `internal/` or domain packages; `cmd/` imports inward and product packages do not import `cmd/` (`cmd/restic/main.go:37-114`, `internal/restic/repository.go:18`). | high |
| `02-command-architecture` | `studies/go-cli-study/reports/final/02-command-architecture.md` | High-scoring CLIs use factory functions, thin commands, shared config/context objects, and lifecycle hooks for shared setup (`pkg/cmd/install.go:132-145`, `pkg/cmdutil/factory.go:16-43`, `internal/cmd/config.go:2236-2454`). | high |
| `03-dependency-injection` | `studies/go-cli-study/reports/final/03-dependency-injection.md` | Manual composition roots, constructor injection, and small core-service interfaces dominate; no studied repo uses a DI framework (`internal/cmd/config.go:362`, `pkg/cmdutil/factory.go:16-43`, `executor.go:22-24`). | high |
| `04-configuration-management` | `studies/go-cli-study/reports/final/04-configuration-management.md` | Elite CLIs converge on explicit layered precedence, post-merge validation, environment prefixes, and flag-change tracking (`internal/cmd/config.go:2253-2287`, `internal/flags/flags.go:314-327`, `restic/internal/global/global.go:139,147`). | high |
| `05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md` | Strong CLIs wrap errors with `%w`, use sentinels or typed errors for classification, and map categories to deterministic exit codes (`internal/ghcmd/cmd.go:44-49`, `internal/errors/fatal.go:10`, `errors/errors.go:47-50`). | high |
| `06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md` | Testable CLIs inject stdout/stderr/stdin, filesystem, and terminal seams; `IOStreams.Test()` and mock terminal patterns make command output assertions deterministic (`pkg/iostreams/iostreams.go:551-568`, `internal/ui/mock.go:10-53`). | high |
| `09-terminal-ux` | `studies/go-cli-study/reports/final/09-terminal-ux.md` | CLI-first tools should preserve scriptability, TTY-aware color/progress, and non-interactive behavior; richer TUI patterns are reserved for sustained interactive workflows (`internal/cmd/prompt.go:20-256`, `lib/terminal/terminal.go:82-86`). | medium |
| `10-logging-observability` | `studies/go-cli-study/reports/final/10-logging-observability.md` | Strong CLIs separate stdout data from stderr diagnostics and use structured, level-aware logging with runtime debug controls (`internal/logging/logging.go:31-71`, `internal/slogs/keys.go:6-231`). | medium |
| `11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md` | Mature Go CLI tests combine table-driven unit tests, command-level integration, golden output checks, and centralized mocks (`internal/cmd/main_test.go:64-174`, `helm/internal/test/test.go:43`, `pkg/httpmock/stub.go:35-199`). | high |
| `13-security` | `studies/go-cli-study/reports/final/13-security.md` | Security-relevant CLIs make trust boundaries explicit through safe path/command handling, redaction types, credential scrubbing, and schema validation (`internal/options/secret_string.go:15-20`, `pkg/registry/transport.go:37-41`, `internal/config/json/validator.go:146`). | high |

## Relevant Patterns

- **Thin entrypoint with product/platform separation:** The project-structure report shows the dominant application layout as `cmd/` plus `internal/`, with CLI files delegating to business packages and avoiding reverse imports (`main.go:16`, `internal/chezmoi/chezmoi.go:1-2`, `cmd/restic/main.go:37-114`). For this sprint, this pattern pressures workspace/config/health logic to live under `internal/` rather than in `cmd/ultraplan`.
- **Command factories with shared dependencies:** Command architecture evidence favors `newXxxCmd(...)` factories and shared dependency containers for command setup (`pkg/cmd/issue/list/list.go:47`, `pkg/cmd/install.go:132`, `cmd/restic/cmd_backup.go:35`). This matters for `init-workspace`, `config show`, and `health` because they share workspace discovery, config loading, output, and error rendering.
- **Manual composition root and explicit dependency passing:** Dependency-injection evidence shows centralized construction through `NewApp`, `newConfig`, `Factory`, or `global.Options`, not reflection-based DI (`internal/cmd/config.go:362`, `pkg/cmd/factory/default.go:26-46`, `cmd/restic/main.go:181-183`). The sprint needs a traceable place to construct workspace/config/output/logging dependencies without global mutable config.
- **Layered config precedence with explicit flag tracking:** Configuration evidence supports defaults plus config file plus environment plus flags, with CLI flags restored or checked through `Changed` so flag defaults do not shadow config (`internal/cmd/config.go:2253-2287`, `internal/flags/flags.go:314-327`, `restic/internal/global/global.go:139,147`). This directly maps to the sprint requirement for built-in defaults < workspace config < env vars < CLI flags.
- **Post-merge config validation:** Config reports show validation after all sources are merged, either centrally or through hooks (`cmd/dive/cli/internal/options/analysis.go:48-53`, `internal/config/k9s.go:423-451`, `internal/config/config.go:609-641`). This matters because `config show` and `health` need to report the effective configuration, not source fragments.
- **Injected IO for deterministic command tests:** IO abstraction evidence points to `IOStreams.Test()` and mock terminal patterns for capturing output (`pkg/iostreams/iostreams.go:551-568`, `internal/ui/mock.go:10-53`, `ui_mock.go:27-33`). This applies to deterministic text/JSON output tests for `config show` and `health`.
- **Error categories with exit code mapping:** Error handling evidence shows specific exit-code constants and error interfaces/types that let scripts distinguish usage, auth, retryable, fatal, and validation failures (`internal/ghcmd/cmd.go:44-49`, `cmd/cmd.go:497-516`, `errors/errors.go:47-50`). Sprint reasoning must choose how workspace/config/validation errors map to the sprint's deterministic exit-code classes.
- **Secret-safe output by type or projection:** Security evidence favors `SecretString` and credential scrubbers so accidental printing yields redacted output (`internal/options/secret_string.go:15-20`, `pkg/registry/transport.go:37-41`, `status.go:332-338`). This matters for `config show`, logging fields, and diagnostics even if this sprint handles only initial config shape.
- **Structured logging with stdout/stderr separation:** Logging evidence favors `slog` or equivalent structured fields, log levels, debug toggles, and keeping logs on stderr while command data remains on stdout (`internal/logging/logging.go:31-71`, `pkg/yqlib/logger.go:5`, `cmd/age/tui.go:37-40`). This sprint should establish the minimal foundation without expanding into durable event persistence.
- **Behavior-focused command and output tests:** Testing evidence favors command execution tests, table-driven subtests, golden output for user-facing formats, and isolated mocks/fakes (`go-task/task_test.go:166-169`, `helm/internal/test/test.go:43`, `chezmoi/internal/chezmoitest/chezmoitest.go:86-92`). This sprint's workspace/config/health behavior is well-suited to temp directories, injected IO, and golden-like JSON/text assertions.

## Trade-Offs

| Trade-Off | Benefit | Cost | When It Matters |
| --- | --- | --- | --- |
| `internal/` product packages vs public `pkg/` packages | Go-enforced encapsulation and freedom to refactor; matches pure CLI app evidence (`internal/` examples at `internal/chezmoi/chezmoi.go:1-2`, `internal/restic/repository.go:18`) | External consumers and external test packages cannot import internals | Sprint reasoning decides package ownership for workspace and platform packages |
| Central shared app/config object vs per-command wiring | Consistent dependency construction and shared workspace/config/output behavior (`pkg/cmdutil/factory.go:16-43`, `internal/cmd/config.go:193-291`) | Can become a god object with too many fields, a concern called out for chezmoi's `Config` | `init-workspace`, `config show`, and `health` share setup but should not accumulate unrelated future workflow state |
| Fixed config precedence function vs priority/value-source abstraction | Fixed order is easy to explain and test (`internal/cmd/config.go:2253-2287`, `internal/flags/flags.go:314-327`) | Less flexible if future config sources need dynamic priority like rclone's configmap (`fs/config/configmap/configmap.go:14-22`) | Sprint 2 has a defined precedence chain and should keep user diagnostics understandable |
| Redaction wrapper types vs projection-only redaction | Type wrappers reduce accidental leaks through `String()` (`internal/options/secret_string.go:15-20`) | More boundary conversions; may be heavier for fields that are only ever displayed through `config show` | Any field may appear in logs, error messages, or JSON output |
| IO abstraction and test constructors vs direct `os.Stdout`/`os.Stderr` | Enables deterministic command tests and output capture (`pkg/iostreams/iostreams.go:551-568`) | Adds small plumbing layer and requires discipline to avoid bypasses | Text/JSON output, command tests, and logging/output separation are acceptance criteria |
| Command-level tests plus golden output vs unit-only tests | Catches user-visible command, exit code, and output regressions (`internal/cmd/main_test.go:64-174`, `helm/internal/test/test.go:43`) | Golden and command tests require fixture maintenance | `config show` and `health` output modes must remain stable for scripts |
| Minimal logging foundation vs full observability now | Establishes level/format/output seams without exceeding sprint scope (`internal/logging/logging.go:31-71`) | Later event persistence may require adapting the skeleton | Sprint excludes durable logs, run event records, and full observability event persistence |

## Anti-Patterns And Warnings

- **Do not put workspace/config business logic in the CLI entrypoint:** Large or logic-heavy entrypoints are called out as a caution (`cmd/age/age.go:105-321`, `opencode/cmd/root.go:49-183`). Keep `cmd/ultraplan` thin and push behavior into `internal/app`, `internal/workspace`, and platform packages.
- **Do not let flag defaults shadow config values:** The config report flags Cobra binding defaults shadowing config files (`cmd/gdu/main.go:46-112`). Sprint config loading needs explicit changed-flag handling or equivalent.
- **Do not bypass the config layer with direct env reads:** Direct `os.Getenv` alongside a config loader makes precedence unpredictable (`internal/config/config.go:163`). Environment overrides should be centralized.
- **Do not create global mutable config/logging state:** DI and configuration reports warn against package-level config and service singletons (`pkg/yqlib/yaml.go:40`, `fs/config.go:14-51`, `internal/config/config.go:123`). Tests for different workspaces/configs must not contaminate each other.
- **Do not print secrets through generic formatting:** Security evidence shows the safer approach is redaction by type/projection and log scrubbers (`internal/options/secret_string.go:15-20`, `pkg/registry/transport.go:37-41`). `config show`, JSON output, and logs must not expose sensitive values.
- **Do not mix stdout data with debug diagnostics:** Logging evidence flags debug output to stdout as script-hostile (`src/core.go:325`). Text/JSON command payloads should be separate from diagnostics.
- **Do not hardcode `os.Stdout`, `os.Stderr`, or filesystem operations where tests need isolation:** IO evidence flags direct `cmd.Stdout = os.Stdout` and template stderr bypasses as testability leaks (`dive/image/docker/cli.go:27`, `internal/cmd/templatefuncs.go:296`).
- **Do not use panic for workspace/config parse or validation failures:** Error evidence treats panic for parse errors as a warning (`dive/image/docker/manifest.go:18`). User-fixable workspace/config issues should return classified errors.
- **Do not assert implementation details in tests:** Testing evidence flags hint-count assertions and fragile regex matching (`internal/view/pod_test.go:23`, `pkg/cmd/issue/list/list_test.go:88-91`). Prefer behavior assertions on output, exit codes, files, and validation results.

## Examples Worth Inspecting

| Example | Path / Source | Why It Is Useful |
| --- | --- | --- |
| Thin CLI and internal boundary | `main.go:16`, `internal/chezmoi/chezmoi.go:1-2`; report `01-project-structure` | Demonstrates unidirectional entrypoint-to-business logic dependency. |
| Command factory pattern | `pkg/cmd/issue/list/list.go:47-118`, `pkg/cmd/install.go:132-145`; report `02-command-architecture` | Shows command construction with dependencies and thin run functions. |
| Shared command factory / DI container | `pkg/cmdutil/factory.go:16-43`; reports `02-command-architecture`, `03-dependency-injection` | Useful for reasoning about shared IO/config/workspace dependencies without globals. |
| Functional options for configuration | `executor.go:22-24`, `executor.go:541-577`; reports `03-dependency-injection`, `06-io-abstraction` | Shows optional dependency and IO injection without constructor bloat. |
| Config flag save/restore precedence | `internal/cmd/config.go:2253-2287`; report `04-configuration-management` | Demonstrates how explicit CLI flags can win over config file values. |
| Generic precedence accessor | `internal/flags/flags.go:314-327`; report `04-configuration-management` | Shows one place to reason about repeated precedence logic. |
| Post-merge config validation | `internal/config/k9s.go:423-451`, `internal/config/config.go:609-641`; report `04-configuration-management` | Shows validation after the effective config has been assembled. |
| Exit code constants and classified errors | `internal/ghcmd/cmd.go:44-49`, `errors/errors.go:47-50`; report `05-error-handling` | Useful for deterministic script-facing failure behavior. |
| IOStreams test constructor | `pkg/iostreams/iostreams.go:551-568`; report `06-io-abstraction` | Direct model for in-memory command output tests. |
| Mock terminal / UI testing | `internal/ui/mock.go:10-53`, `ui_mock.go:27-33`; report `06-io-abstraction` | Shows output capture without real TTYs. |
| TTY-aware non-interactive behavior | `internal/cmd/prompt.go:20-256`, `lib/terminal/terminal.go:82-86`; report `09-terminal-ux` | Useful for keeping output script-friendly and avoiding TTY assumptions. |
| Structured logging foundation | `internal/logging/logging.go:31-71`, `internal/slogs/keys.go:6-231`; report `10-logging-observability` | Shows level-aware structured logs and consistent field naming. |
| Golden/output test helpers | `helm/internal/test/test.go:43`, `go-task/task_test.go:166-169`; report `11-testing-strategy` | Useful for stable text/JSON output assertions. |
| Secret redaction | `internal/options/secret_string.go:15-20`, `status.go:332-338`; report `13-security` | Shows both type-level and function-level redaction patterns. |
| Credential/header scrubbing | `pkg/registry/transport.go:37-41`; report `13-security` | Useful if health/config diagnostics include runtime environment or HTTP-related fields later. |

## Design Pressures

- Sprint 2 is foundational: workspace discovery, initialization, path helpers, config loading, output helpers, logging skeleton, and health checks will be reused by later study/runtime commands.
- `workspace` owns location, path normalization, path safety, initialization, and structural validation; it should not absorb study, runtime, or target/sprint behavior.
- `platform/config` and `platform/logging` are cross-cutting infrastructure; platform packages must not import product modules.
- The sprint's config precedence is explicit and testable: built-in defaults, then workspace config, then environment variables, then CLI flags.
- `config show` and `health` need deterministic text and JSON output; JSON must avoid ANSI formatting and secrets.
- `init-workspace` mutates the filesystem, so dry-run planning, idempotence, path safety, and clear errors are major pressures.
- Health is a skeleton in this sprint: it should cover workspace/config/filesystem/environment basics without OpenCode, provider, agentwrap runtime, or model availability checks.
- Logging is a skeleton in this sprint: the evidence supports structured/level-aware foundations, but durable file logs and event persistence are out of scope.
- Tests must be offline, deterministic, and isolated from any real workspace, network, provider credentials, OpenCode installation, or user environment.
- Secret redaction is required early because config and diagnostics can become a long-lived public surface.

## Open Questions For Reasoning

- What exact workspace marker or structure should discovery require when walking from current directory to parents?
- How should explicit workspace path, `ULTRAPLAN_WORKSPACE`, current directory, and parent traversal report conflicts or invalid candidates?
- Which path operations must reject workspace escapes in this sprint, and which may accept absolute paths only with explicit allowance?
- What is the minimal `ultraplan.yml` schema for this sprint, given that runtime execution is out of scope but later runtime config exists?
- Should config validation fail on unknown fields, warn on extension fields, or ignore them for version-compatible forward evolution?
- How should config sources be represented in `config show`: show only values, show source metadata, or offer source metadata only in JSON/debug output?
- What field naming or type strategy will define sensitive config values for redaction without overfitting to future runtime secrets?
- What error classes map to exit codes for missing workspace, invalid workspace, invalid config, validation failure, bad arguments, and filesystem errors?
- What exact health checks are in scope now: workspace exists, required files exist, config parses/validates, paths are writable/readable, environment variables are well-formed?
- How minimal can the logging interface be while leaving room for later structured event records?
- Should text/JSON output helpers live only in `internal/app`, or should a small platform output abstraction exist? Evidence supports injection either way, but sprint reasoning must place ownership.
- Which command outputs deserve golden-style fixtures versus ordinary table-driven assertions?

## Evidence Pointers

- `studies/go-cli-study/reports/final/01-project-structure.md`: Inspect Pattern 1, Pattern 2, Pattern 4, tradeoffs around `internal/` vs `pkg/`, and warnings about large entrypoints/global packages.
- `studies/go-cli-study/reports/final/02-command-architecture.md`: Inspect factory functions, central config/context object, lifecycle hooks, `cmd.Run()` wrappers, and long `RunE` warnings.
- `studies/go-cli-study/reports/final/03-dependency-injection.md`: Inspect centralized composition roots, constructor injection, functional options, interface abstraction, and global state anti-patterns.
- `studies/go-cli-study/reports/final/04-configuration-management.md`: Inspect flag-restore, generic precedence accessor, environment prefix, post-load validation, XDG notes, and config anti-patterns.
- `studies/go-cli-study/reports/final/05-error-handling.md`: Inspect `%w` wrapping, sentinel/typed errors, user vs operational separation, exit code mapping, and panic warnings.
- `studies/go-cli-study/reports/final/06-io-abstraction.md`: Inspect `IOStreams.Test()`, `System`/`FS` interfaces, `MockTerminal`, functional IO options, and direct `os.*` bypass warnings.
- `studies/go-cli-study/reports/final/09-terminal-ux.md`: Inspect non-TTY fallback, TTY color checks, cancellation behavior, and CLI-first tradeoffs.
- `studies/go-cli-study/reports/final/10-logging-observability.md`: Inspect `slog` patterns, stdout/stderr separation, debug controls, structured keys, and logging anti-patterns.
- `studies/go-cli-study/reports/final/11-testing-strategy.md`: Inspect table-driven tests, testscript/golden tradeoffs, centralized mocks, behavior-focused assertions, and global-state cleanup warnings.
- `studies/go-cli-study/reports/final/13-security.md`: Inspect `SecretString`, credential scrubbing, schema validation, path canonicalization warnings, and trust boundary guidance.
- `projects/ultraplan-go/sprints/02-workspace-config-health/requirements.md`: Use for sprint scope, acceptance criteria, explicit non-goals, and output paths only.
- `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`: Use for UltraPlan scope and constraints only, not as evidence replacing the selected study reports.

## Handoff To Reasoning

- Use this handbook as evidence input.
- Validate whether the observed patterns fit this project's constraints.
- Do not copy external patterns without sprint-specific reasoning.
- Do not decide final package boundaries, command wiring, config schema, exit-code mapping, or health-check scope in this handbook.
- Keep sprint reasoning aligned with selected evidence while respecting the sprint non-goals: no study workflows, no runtime execution, no agentwrap/OpenCode health checks, no durable observability store, and no target/sprint product workflows.
