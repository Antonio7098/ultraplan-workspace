# Configuration Contract

## Purpose

This contract defines how runtime configuration must be loaded, validated, exposed, and used.

It governs:

- environment and config source precedence
- config file locations and environment variable naming
- typed settings
- persisted config compatibility
- startup, preflight, and command validation
- secrets separation
- public vs private config
- test and local development config

## Scope

This contract applies when a change:

- adds or changes runtime settings
- changes environment variables or config files
- introduces runtime, dependency, storage, auth, telemetry, feature flag, or model configuration
- exposes config to frontend, CLI, diagnostics, or logs

## Requirement Index

| ID | Title | Applies To | Severity If Violated |
| --- | --- | --- | --- |
| CFG-SOURCE-001 | Runtime config must come from explicit sources | all config | High |
| CFG-TYPE-001 | Config must be typed, parsed, and validated | app startup and CLI startup | High |
| CFG-COMPAT-001 | Persisted config must remain compatible or migratable | durable config files/settings | High |
| CFG-SECRET-001 | Secrets must be separated from public config | secret-bearing config | Blocker |
| CFG-START-001 | Required config must fail before work starts | startup, preflight, readiness, command setup | Blocker |
| CFG-ENV-001 | Environment-specific behaviour must be explicit | dev/test/prod | High |
| CFG-PUBLIC-001 | Externally exposed config must be allowlisted | public config, CLI diagnostics, API/UI config | Blocker |
| CFG-OBS-001 | Effective non-secret config must be inspectable | diagnostics/ops | Medium |

## Core Principles

1. Config changes behaviour without code changes.
2. Required config must be validated before serving traffic, running commands, or starting jobs.
3. Secrets and public config are different categories.
4. Defaults must be safe and obvious.
5. Environment-specific behaviour must not be accidental.
6. Operators should be able to inspect effective non-secret config.
7. Durable user config is a compatibility surface.

## Requirements

### CFG-SOURCE-001: Runtime Config Must Come From Explicit Sources

**Rule**
Runtime config must be loaded from documented sources with clear precedence.

**Required**
- define config source order such as defaults, config files, environment variables, secret manager, CLI flags
- use documented, platform-appropriate config file locations for user-facing tools
- namespace environment variables with an application-specific prefix where the ecosystem supports it
- avoid hidden config reads scattered through the codebase
- centralize settings construction where practical

**Forbidden**
- reading arbitrary environment variables throughout business logic
- using unprefixed or ambiguous environment variables for application-specific settings without justification
- mixing config parsing with domain logic
- relying on undocumented local files for runtime behaviour

**Evidence**
- settings module or equivalent owns config loading
- docs identify required and optional config
- docs identify config file locations and environment variable names where applicable

### CFG-TYPE-001: Config Must Be Typed, Parsed, And Validated

**Rule**
Config values must be parsed into typed settings before use.

**Required**
- parse strings into typed values such as URLs, durations, booleans, enums, integers, paths, and lists
- validate ranges, required combinations, and incompatible options
- reject ambiguous boolean/string parsing

**Forbidden**
- comparing raw config strings throughout the app
- interpreting missing or malformed values differently in different modules
- silently defaulting invalid values

**Evidence**
- config parsing tests cover malformed and missing values

### CFG-COMPAT-001: Persisted Config Must Remain Compatible Or Migratable

**Rule**
Config that persists across process runs, releases, or user machines must define its compatibility and migration behaviour.

**Required**
- version durable config formats when keys, semantics, or defaults may change incompatibly
- provide migration, compatibility, or clear rejection behaviour for previously supported config
- document removed, renamed, or meaningfully changed settings
- test migrations or compatibility paths for significant persisted config changes

**Forbidden**
- silently reinterpreting existing user config with materially different behaviour
- breaking durable config formats without migration, compatibility handling, or explicit release notes
- treating generated templates as the only source of truth while existing user config can remain stale

**Evidence**
- persisted config changes include migration/compatibility tests or documented rejection behaviour
- release notes or user docs identify breaking config changes

### CFG-SECRET-001: Secrets Must Be Separated From Public Config

**Rule**
Secret-bearing config must be loaded, stored, displayed, and logged separately from public operational config.

**Required**
- mark secret fields explicitly
- redact secrets in errors, diagnostics, and logs
- avoid exposing secret-bearing settings to frontend, public API, CLI diagnostics, or generated artifacts

**Forbidden**
- bundling private secrets into public builds or artifacts
- dumping all settings to diagnostics without redaction
- using ordinary config docs as a place to store secret values

**Evidence**
- diagnostics redacts secret-bearing settings
- public config export is allowlisted

### CFG-START-001: Required Config Must Fail Before Work Starts

**Rule**
Missing or invalid required config must fail startup, readiness, preflight, or command execution before accepting workload or performing side effects.

**Required**
- validate required settings at startup
- fail readiness or health when required config-dependent capabilities are unusable
- fail CLI commands before performing side effects if required config is missing

**Forbidden**
- delaying known-fatal config failure until the first real request, task, or side effect
- silently disabling required subsystems in production
- falling back to insecure test credentials

**Evidence**
- startup/config tests cover fatal missing settings

### CFG-ENV-001: Environment-Specific Behaviour Must Be Explicit

**Rule**
Development, test, staging, and production differences must be visible and controlled.

**Required**
- define allowed environment names or runtime modes
- gate debug/test behaviour explicitly
- reject production runtime with unsafe development settings

**Forbidden**
- enabling debug routes, fake auth, permissive CORS, or verbose sensitive logs by default
- making production behaviour depend on accidental absence of env vars

**Evidence**
- production config checks reject unsafe combinations

### CFG-PUBLIC-001: Externally Exposed Config Must Be Allowlisted

**Rule**
Only explicitly public or diagnostic-safe config may be exposed to browsers, mobile clients, SDKs, CLIs, generated docs, diagnostics, or artifacts.

**Required**
- maintain an allowlist for public config
- separate public runtime config from private settings
- avoid exposing internal URLs, secret IDs, credentials, or provider details unless intended

**Forbidden**
- exporting entire settings objects to users or clients
- relying on naming conventions alone to determine public safety

**Evidence**
- public config export is small and intentional

### CFG-OBS-001: Effective Non-Secret Config Must Be Inspectable

**Rule**
Operators should be able to inspect effective non-secret config relevant to runtime diagnosis.

**Required**
- expose safe config summaries in diagnostics or startup logs
- include versions, enabled features, provider/model names, limits, and runtime mode where useful
- redact secrets and sensitive values

**Forbidden**
- making operational behaviour impossible to diagnose because config is invisible
- exposing full secret-bearing config

**Evidence**
- diagnostics/startup signals include safe effective settings

## Review Rejection Criteria

Reject a change if it:

- reads runtime config from arbitrary places
- uses raw unparsed config strings in core logic
- changes durable config semantics without compatibility, migration, or explicit rejection behaviour
- exposes secrets through logs, diagnostics, frontend, or errors
- allows required production config to fail only at first use
- enables insecure development behaviour in production
- exports public config without an allowlist

---
