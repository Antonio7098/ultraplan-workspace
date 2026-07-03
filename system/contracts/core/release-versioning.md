# Release And Versioning Contract

## Purpose

This contract defines how releases, versions, changelogs, artifacts, and rollback readiness must be handled.

It governs:

- versioning policy
- release notes/changelogs
- artifact integrity
- compatibility checks
- deployment gates
- rollback and migration notes

## Scope

This contract applies when a change:

- prepares a package, CLI, API, frontend, backend, or tool release
- changes public API, CLI, schemas, config, migrations, prompts, or runtime behaviour
- creates deployable artifacts
- changes release automation or version metadata

## Requirement Index

| ID | Title | Applies To | Severity If Violated |
| --- | --- | --- | --- |
| REL-VERSION-001 | Version changes must reflect compatibility impact | packages/APIs/CLIs | High |
| REL-CHANGELOG-001 | Releases must describe user-visible and operational changes | releases | Medium |
| REL-ARTIFACT-001 | Release artifacts must be reproducible, identifiable, and clean | builds/packages | High |
| REL-GATE-001 | Release gates must run the relevant checks | CI/release | High |
| REL-MIGRATE-001 | Releases with migrations or config changes need operator notes | deployments | High |
| REL-ROLLBACK-001 | Risky releases need rollback or recovery stance | production releases | High |
| REL-PROMPT-001 | Prompt/model behaviour changes must be versioned | AI releases | High |

## Core Principles

1. Version numbers and release notes communicate risk.
2. Artifacts should be traceable to source.
3. Releases should not depend on local developer machines.
4. Migrations and config changes are operational events.
5. Rollback is easiest when planned before release.
6. AI prompt/model changes are behavioural releases.

## Requirements

### REL-VERSION-001: Version Changes Must Reflect Compatibility Impact

**Rule**
Version increments must reflect compatibility impact according to the project’s declared scheme.

**Required**
- define whether the project uses SemVer, CalVer, or another scheme
- treat public API/CLI/package/schema changes as compatibility-relevant
- mark breaking changes clearly

**Forbidden**
- hiding breaking changes in patch releases for stable packages
- releasing unversioned behavioural changes for public surfaces

**Evidence**
- release notes and version bump align with change impact

### REL-CHANGELOG-001: Releases Must Describe User-Visible And Operational Changes

**Rule**
Releases must include enough change information for users, operators, and future agents to understand impact.

**Required**
- describe added, changed, fixed, removed, deprecated, security, migration, and operational changes
- include breaking changes and required actions prominently
- include branch/PR/commit references where useful

**Forbidden**
- vague entries such as “misc fixes” for meaningful changes
- omitting migration/config/operator changes

**Evidence**
- changelog/release notes updated before release

### REL-ARTIFACT-001: Release Artifacts Must Be Reproducible, Identifiable, And Clean

**Rule**
Build outputs must be traceable, clean of secrets/local files, and installable/deployable from the artifact alone.

**Required**
- include version, commit, build time/source metadata where useful
- build from clean checkout/CI where possible
- verify package/archive/container contents
- exclude secrets, local caches, test data, and dev-only files

**Forbidden**
- publishing artifacts built from dirty/untracked local state without disclosure
- artifacts that require source-tree files not included in the package

**Evidence**
- artifact inspection or clean install/deploy checks pass

### REL-GATE-001: Release Gates Must Run Relevant Checks

**Rule**
Release candidates must pass the checks relevant to their surface area.

**Required**
- run tests, lint/type checks, security/dependency checks, contract checks, build checks, and smoke tests appropriate to the change
- include real-provider smoke or staged verification where critical provider paths changed

**Forbidden**
- releasing after bypassing failing checks without explicit risk acceptance
- treating unrelated green checks as proof that changed critical paths work

**Evidence**
- CI/release summary identifies check status

### REL-MIGRATE-001: Releases With Migrations Or Config Changes Need Operator Notes

**Rule**
Schema, data, config, secret, prompt, model, provider, or infrastructure changes must include operator-facing notes.

**Required**
- document migration order, expected downtime, backfill, config additions, secret rotations, and compatibility windows
- identify whether old and new versions can run concurrently

**Forbidden**
- shipping required env/config changes without docs
- shipping migrations with no deployment sequencing guidance where sequencing matters

**Evidence**
- release notes include operational migration/config section

### REL-ROLLBACK-001: Risky Releases Need Rollback Or Recovery Stance

**Rule**
High-risk releases must define whether rollback is safe and what recovery path exists.

**Required**
- identify irreversible migrations, external effects, provider changes, and data format changes
- define rollback, forward-fix, restore, or compensation plan

**Forbidden**
- destructive releases with no recovery note
- assuming rollback works after irreversible data changes

**Evidence**
- release notes or deploy checklist include rollback/recovery stance

### REL-PROMPT-001: Prompt/Model Behaviour Changes Must Be Versioned

**Rule**
AI prompt, model, tool, eval, or structured-output behaviour changes must be treated as behavioural release changes.

**Required**
- update prompt/model/runtime version references
- include eval/smoke results where risk warrants
- document behaviour changes that users/operators may notice

**Forbidden**
- silent production prompt edits with no version or traceability
- changing model/provider defaults without release notes when behaviour/cost/latency may change

**Evidence**
- prompt/model change appears in release notes and runtime version metadata

## Review Rejection Criteria

Reject a release if it:

- includes breaking changes without version/migration treatment
- lacks changelog notes for meaningful user/operator changes
- publishes artifacts that fail clean install/deploy checks
- bypasses relevant failing checks without documented risk acceptance
- ships migrations/config changes without operator notes
- has risky irreversible changes with no recovery stance
- silently changes AI prompt/model behaviour

---