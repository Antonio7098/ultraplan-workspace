# Package Public API Contract

## Purpose

This contract defines how reusable packages and libraries must manage their public API, dependencies, versioning, and compatibility.

It governs:

- public exports
- SemVer and compatibility
- dependency policy
- deprecation
- package metadata
- examples and documentation
- CLI entry points for packages

## Scope

This contract applies when a change:

- adds, removes, renames, or changes exported symbols
- changes package metadata, dependencies, build config, or entry points
- changes documented public behaviour
- publishes artifacts to package registries or consumers

## Requirement Index

| ID | Title | Applies To | Severity If Violated |
| --- | --- | --- | --- |
| PKG-API-001 | Public API surface must be intentional and minimal | exports | High |
| PKG-COMPAT-001 | Compatibility changes must follow versioning policy | releases | Blocker |
| PKG-DEPREC-001 | Deprecations must be staged and documented | API removal/change | Medium |
| PKG-DEPS-001 | Dependencies must be minimal, justified, and bounded | package deps | High |
| PKG-BUILD-001 | Build artifacts must be reproducible and complete | releases | High |
| PKG-DOC-001 | Public APIs must have usage docs/examples | public package APIs | Medium |
| PKG-ENTRY-001 | Package entry points must be explicit and stable | CLI/plugin packages | Medium |

## Core Principles

1. The public API is a contract with users.
2. Export less than you implement.
3. Compatibility must be intentional.
4. Version numbers communicate risk.
5. Dependencies become your users’ dependencies.
6. Deprecation is a migration process, not a surprise.
7. Packages must be usable from clean installs.

## Requirements

### PKG-API-001: Public API Surface Must Be Intentional And Minimal

**Rule**
Only stable, intentional symbols should be exported as public API.

**Required**
- define public exports explicitly
- keep internal modules private or documented as unstable
- avoid exporting implementation details, adapters, generated internals, or test helpers by default

**Forbidden**
- wildcard exports that expose internals accidentally
- documenting internal symbols as if they are stable
- changing public behaviour through internal refactor accidents

**Evidence**
- public API snapshots or export lists are reviewed

### PKG-COMPAT-001: Compatibility Changes Must Follow Versioning Policy

**Rule**
Package releases must follow a declared versioning scheme, usually SemVer or a documented project alternative.

**Required**
- increment major version for breaking public API changes when using SemVer
- increment minor version for backwards-compatible features
- increment patch version for backwards-compatible fixes
- treat public CLI output, config, plugin contracts, and generated schemas as public API when users depend on them

**Forbidden**
- releasing breaking changes as patch/minor versions
- modifying already-published artifacts in place
- changing documented behaviour without release notes

**Evidence**
- release notes identify compatibility impact
- public API diff or tests cover significant changes

### PKG-DEPREC-001: Deprecations Must Be Staged And Documented

**Rule**
Public API removal or behaviour replacement must provide a migration path unless the project explicitly does not guarantee compatibility.

**Required**
- mark deprecated APIs clearly
- explain replacement and timeline where possible
- preserve old behaviour until removal version where feasible

**Forbidden**
- removing documented APIs with no notice in stable packages
- deprecating without replacement guidance

**Evidence**
- docs/changelog include deprecation notes

### PKG-DEPS-001: Dependencies Must Be Minimal, Justified, And Bounded

**Rule**
Packages should minimize dependency burden and avoid imposing unnecessary risk on consumers.

**Required**
- justify new runtime dependencies
- use optional dependencies for optional features where ecosystem supports it
- set compatible version ranges carefully
- avoid large dependencies for tiny utilities

**Forbidden**
- importing optional heavy dependencies at package import time
- forcing provider/framework dependencies onto users who do not need them
- adding risky transitive dependencies without review

**Evidence**
- dependency diff is reviewed and dependency role is documented

### PKG-BUILD-001: Build Artifacts Must Be Reproducible And Complete

**Rule**
Published packages must include the files needed for installation, type checking, docs/examples where intended, and runtime operation.

**Required**
- verify clean install from built artifact, not only source tree
- include license, metadata, readme, type declarations, generated files, and package data where applicable
- exclude secrets, local caches, test artifacts, and machine-specific files

**Forbidden**
- publishing artifacts that only work because local files exist
- publishing secrets or development caches

**Evidence**
- release check installs package in a clean environment

### PKG-DOC-001: Public APIs Must Have Usage Docs/Examples

**Rule**
Public APIs should have enough documentation for users to succeed without reading internals.

**Required**
- document core concepts, install, quickstart, common examples, errors, and configuration
- keep examples tested or smoke-verified where practical

**Forbidden**
- exposing public APIs with no examples or behavioural explanation
- stale docs after public API changes

**Evidence**
- docs/examples updated with API changes

### PKG-ENTRY-001: Package Entry Points Must Be Explicit And Stable

**Rule**
CLI commands, plugins, dynamic entry points, and public module entry points must be intentionally declared and compatible.

**Required**
- declare CLI/plugin entry points through ecosystem-standard package metadata
- document entry point names and expected callable contracts
- treat entry point changes as public compatibility changes

**Forbidden**
- relying on users importing internal modules as entry points
- moving command entry points without compatibility handling

**Evidence**
- entry point tests verify installed package behaviour

## Review Rejection Criteria

Reject a package change if it:

- accidentally exports internals
- breaks public API without versioning/release note treatment
- removes public APIs without migration guidance
- adds unjustified heavy/risky dependencies
- publishes artifacts that fail clean install
- updates code but leaves public docs/examples stale

---