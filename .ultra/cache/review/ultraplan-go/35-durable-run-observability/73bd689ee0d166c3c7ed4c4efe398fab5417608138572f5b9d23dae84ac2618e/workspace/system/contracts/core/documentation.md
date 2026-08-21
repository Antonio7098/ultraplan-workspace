# Documentation Contract

## Purpose

This contract defines when documentation is required, where it belongs, and how it must stay aligned with code.

It governs:

- architecture docs
- decision records
- public API/package/CLI/artifact docs
- examples
- generated docs
- operational runbooks
- agent-facing context

## Scope

This contract applies when a change:

- changes architecture, public API, CLI behaviour, package exports, workflows, LLM prompts/runtime, config, migrations, or operational behaviour
- adds new concepts, modules, commands, or user-facing features
- changes generated documentation or examples

## Requirement Index

| ID | Title | Applies To | Severity If Violated |
| --- | --- | --- | --- |
| DOC-OWNER-001 | Important docs must have clear ownership and location | all docs | Medium |
| DOC-ARCH-001 | Architecture changes need decision context | architecture changes | Medium |
| DOC-PUBLIC-001 | Public surfaces must be documented | APIs/CLIs/packages/artifacts | High |
| DOC-OPS-001 | Operational changes need runbook or operator notes | runtime/ops | Medium |
| DOC-EXAMPLE-001 | Examples must be correct and maintained | examples/tutorials | Medium |
| DOC-GEN-001 | Generated docs must be reproducible and checked | generated docs | Medium |
| DOC-AGENT-001 | Agent-facing docs must be concise, current, and authoritative | AI coding workflows | Medium |

## Core Principles

1. Documentation should explain intent, usage, and trade-offs.
2. Code shows what happens; docs preserve why it is shaped that way.
3. Public surfaces need public docs.
4. Operational behaviour needs operator docs.
5. Examples are tests of usability.
6. Agent-facing docs should reduce ambiguity, not add noise.

## Requirements

### DOC-OWNER-001: Important Docs Must Have Clear Ownership And Location

**Rule**
Project docs must live in predictable locations with clear purpose.

**Required**
- define where architecture, API, CLI, package, artifact, frontend, operational, and decision docs belong
- avoid duplicate competing source-of-truth documents
- mark stale/archived docs clearly

**Forbidden**
- scattering critical decisions across chats, PR comments, or orphaned markdown files only
- leaving two docs with conflicting instructions

**Evidence**
- README or docs index points to authoritative docs

### DOC-ARCH-001: Architecture Changes Need Decision Context

**Rule**
Meaningful architecture changes must record the decision and trade-off.

**Required**
- use ADRs, design notes, or sprint reasoning for major boundaries, dependencies, storage, runtime, provider, or workflow decisions
- include context, decision, alternatives considered, and consequences

**Forbidden**
- major architecture changes with no rationale
- docs that describe the result but not the reason/trade-off

**Evidence**
- decision record or reasoning note accompanies architecture-level changes

### DOC-PUBLIC-001: Public Surfaces Must Be Documented

**Rule**
APIs, CLIs, packages, config, durable artifacts, and user-facing features must have enough documentation for intended consumers.

**Required**
- document install/setup, authentication/config, workspace or artifact layout, common usage, errors, examples, and compatibility expectations where relevant
- update docs when public behaviour changes

**Forbidden**
- public commands/endpoints/exports/artifacts with no discoverable docs
- docs that contradict actual runtime behaviour

**Evidence**
- docs updated alongside public surface changes

### DOC-OPS-001: Operational Changes Need Runbook Or Operator Notes

**Rule**
Changes affecting deployment, config, migrations, local state, monitoring, failure handling, or incident response must include operator-facing documentation.

**Required**
- document startup/preflight checks, health/readiness, diagnostics, alerts, migrations, config, secrets, rollback/recovery, and common failures where relevant

**Forbidden**
- operational behaviour discoverable only by reading source code
- alerts/errors with no explanation of operator action

**Evidence**
- runbook/operator docs updated for operational changes

### DOC-EXAMPLE-001: Examples Must Be Correct And Maintained

**Rule**
Examples must compile, run, or remain clearly illustrative and aligned with current APIs.

**Required**
- keep examples minimal and realistic
- test examples where practical
- avoid examples that encourage insecure or deprecated patterns

**Forbidden**
- stale examples using removed APIs
- examples with fake insecure secrets copied into docs

**Evidence**
- example smoke/test/docs checks cover important examples where feasible

### DOC-GEN-001: Generated Docs Must Be Reproducible And Checked

**Rule**
Generated docs such as API schemas, CLI help snapshots, typedoc/godoc/pydoc, SDK docs, and changelog indexes must be reproducible.

**Required**
- document generation commands
- check generated docs in CI where drift matters
- avoid manual edits to generated files unless explicitly intended

**Forbidden**
- committing generated docs with no way to regenerate them
- stale generated docs after schema/API changes

**Evidence**
- docs generation check passes or drift is reviewed

### DOC-AGENT-001: Agent-Facing Docs Must Be Concise, Current, And Authoritative

**Rule**
Docs intended for AI coding agents must be short enough to use, current enough to trust, and clear about authority.

**Required**
- keep contracts normative and reasoning templates procedural
- avoid mixing hard rules with exploratory notes
- give agents reviewable evidence expectations
- keep canonical guides indexed and deduplicated

**Forbidden**
- long contradictory instruction files with unclear priority
- stale agent docs that point to removed architecture
- hiding project-specific laws inside general advice docs

**Evidence**
- ops/docs index separates contracts, protocols, templates, and archived materials

## Review Rejection Criteria

Reject a change if it:

- changes public behaviour without docs update
- changes architecture with no decision context
- changes operational requirements with no operator notes
- leaves generated docs stale
- adds examples that do not work or teach unsafe patterns
- creates conflicting docs without resolving source of truth

---
