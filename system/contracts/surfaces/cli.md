# CLI Contract

## Purpose

This contract defines how command-line tools must behave for humans, scripts, automation, and agents.

It governs:

- command shape and naming
- help and discoverability
- flags and config precedence
- stdout/stderr discipline
- exit codes
- machine-readable output
- lifecycle setup, cancellation, and cleanup
- destructive action safety
- interactive vs non-interactive behaviour

## Scope

This contract applies when a change:

- adds or changes a CLI command, flag, argument, config source, or output format
- changes process exit behaviour
- changes file/network side effects
- changes interactive prompts or automation mode

## Requirement Index

| ID | Title | Applies To | Severity If Violated |
| --- | --- | --- | --- |
| CLI-SHAPE-001 | Commands and flags must be predictable and discoverable | all commands | Medium |
| CLI-HELP-001 | Help/version output must exist and remain useful | all CLIs | Medium |
| CLI-IO-001 | Stdout/stderr must be separated by purpose | all output | High |
| CLI-EXIT-001 | Exit codes must be stable and meaningful | all commands | High |
| CLI-JSON-001 | Machine-readable output must be explicit and stable | automation-facing commands | High |
| CLI-CONFIG-001 | Config precedence must be documented and inspectable | config-bearing CLIs | Medium |
| CLI-LIFE-001 | CLI lifecycle and cancellation behaviour must be explicit | setup, teardown, signals | High |
| CLI-SAFE-001 | Destructive actions require explicit safety controls | mutating commands | High |
| CLI-NONINT-001 | Non-interactive mode must never hang on prompts | CI/automation | High |

## Core Principles

1. CLIs are both user interfaces and automation APIs.
2. Human-readable output and machine-readable output have different contracts.
3. Stdout is for intended command output; stderr is for diagnostics and progress.
4. Exit codes are part of the public interface.
5. Destructive actions must be explicit.
6. Non-interactive execution must be deterministic.
7. Help should teach the user the next step.
8. Shared setup, cancellation, and cleanup behaviour should be centralized and testable.

## Requirements

### CLI-SHAPE-001: Commands And Flags Must Be Predictable And Discoverable

**Rule**
Command names, flags, arguments, and subcommands must follow consistent conventions.

**Required**
- use clear verbs and nouns
- keep flag names consistent across commands
- prefer long flags for clarity and short aliases only for common operations
- keep required positional args minimal and obvious

**Forbidden**
- changing flag meanings between commands
- using ambiguous boolean flags that hide separate workflows
- requiring users to memorize undocumented argument order

**Evidence**
- command help clearly shows usage, args, flags, and examples

### CLI-HELP-001: Help/Version Output Must Exist And Remain Useful

**Rule**
Every CLI must support help and version discovery.

**Required**
- support `--help` at root and command level
- support `--version` or equivalent where a versioned artifact exists
- include examples for non-trivial commands
- document environment/config sources where relevant

**Forbidden**
- commands that fail without telling the user how to correct usage
- hidden commands or flags required for normal use

**Evidence**
- CLI snapshot or smoke tests cover help/version output where stable enough

### CLI-IO-001: Stdout/Stderr Must Be Separated By Purpose

**Rule**
Stdout must contain the intended result of the command; stderr must contain diagnostics, warnings, prompts, and progress.

**Required**
- write data intended for piping to stdout
- write progress, debug, warnings, and errors to stderr
- keep quiet/no-color/no-progress modes where automation needs stable output

**Forbidden**
- writing progress spinners into machine-readable stdout
- mixing logs with JSON output
- hiding errors in stdout while returning success

**Evidence**
- tests verify stdout/stderr for automation-facing commands

### CLI-EXIT-001: Exit Codes Must Be Stable And Meaningful

**Rule**
Commands must exit with stable codes that reflect success, caller error, operational failure, or partial failure.

**Required**
- return `0` only for success
- return non-zero for validation errors, missing config, dependency failure, permission failure, or unexpected failure
- document special exit codes where scripts need to branch on them

**Forbidden**
- returning `0` when the requested operation failed
- using many undocumented exit codes with no stable meaning

**Evidence**
- tests cover important success and failure exit codes

### CLI-JSON-001: Machine-Readable Output Must Be Explicit And Stable

**Rule**
Commands intended for automation must provide an explicit machine-readable output mode.

**Required**
- provide `--json` or equivalent for structured output where useful
- keep JSON schema stable for public commands
- include stable fields such as status, code, data, and errors where applicable
- ensure no logs/progress contaminate JSON stdout

**Forbidden**
- requiring scripts to scrape human prose for important data
- changing JSON field meanings without compatibility consideration

**Evidence**
- JSON output tests assert stable keys and error shapes

### CLI-CONFIG-001: Config Precedence Must Be Documented And Inspectable

**Rule**
CLI config must follow a documented precedence order.

**Required**
- define precedence among defaults, config files, environment variables, flags, and command args
- allow users to inspect effective non-secret config when helpful
- redact secrets in config display
- make config source and effective-value diagnostics available through documented commands, flags, or diagnostics output where useful

**Forbidden**
- hidden config sources that override user-provided flags unexpectedly
- printing secrets from config commands

**Evidence**
- config precedence tests cover conflicts

### CLI-LIFE-001: CLI Lifecycle And Cancellation Behaviour Must Be Explicit

**Rule**
Shared CLI setup, execution, cancellation, and cleanup behaviour must be deliberate rather than duplicated or hidden in command bodies.

**Required**
- centralize shared setup such as config loading, validation, logging, runtime construction, and output mode selection where practical
- propagate cancellation and timeouts through long-running commands
- define how user cancellation, timeout, partial completion, and cleanup failures map to terminal status and exit codes
- run cleanup in a way that is observable and does not silently hide the primary command outcome

**Forbidden**
- duplicating lifecycle-sensitive setup across commands with inconsistent behaviour
- treating cancellation, timeout, and failure as indistinguishable terminal states
- leaving long-running work orphaned after the CLI process receives cancellation

**Evidence**
- command tests or smoke tests cover setup failure, cancellation/timeout, and cleanup behaviour where relevant
- documented exit behaviour distinguishes success, failure, cancellation, timeout, and partial completion where scripts need to branch

### CLI-SAFE-001: Destructive Actions Require Explicit Safety Controls

**Rule**
Commands that delete, overwrite, publish, migrate, spend money, send messages, or mutate external state must include safety controls.

**Required**
- require confirmation or explicit flags for dangerous operations
- provide `--dry-run` where useful and feasible
- show what will change before destructive execution when practical
- make force flags clear and intentionally named

**Forbidden**
- destructive defaults with no confirmation, dry-run, or scope display
- broad glob/path deletion without preview or constraints

**Evidence**
- tests cover dry-run and confirmation bypass rules

### CLI-NONINT-001: Non-Interactive Mode Must Never Hang On Prompts

**Rule**
CI, scripts, and agent execution must be able to run commands without hidden prompts.

**Required**
- support `--yes`, `--no-input`, `--non-interactive`, or equivalent where prompts exist
- fail clearly when required input is missing in non-interactive mode
- avoid TTY assumptions in automation paths

**Forbidden**
- waiting indefinitely for input in CI/non-TTY contexts
- prompting after partial side effects have already occurred

**Evidence**
- non-interactive tests cover prompt-requiring commands

## Review Rejection Criteria

Reject a CLI change if it:

- returns success when the operation failed
- contaminates JSON/stdout with logs or progress
- changes public output/exit behaviour without compatibility consideration
- leaves lifecycle setup, cancellation, or cleanup semantics implicit for long-running commands
- adds destructive commands without confirmation/dry-run/scope controls
- hangs in non-interactive mode
- exposes secrets through config/help/debug output

---
