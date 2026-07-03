# Contract Mapping Protocol

## Purpose

This protocol maps a planned change against the applicable operational contracts.

In the Ultra workflow, contract mapping feeds into `sprint-index.md`, area-specific `reasoning/*.md`, `sprint-reasoning.md`, and `review.md`. It is not the main planning document.

It exists to:
- identify which contracts apply to the work
- identify which requirement IDs matter
- explain architectural and operational trade-offs
- produce reviewable implementation intent
- prevent accidental contract drift

This is not a design document. This is a scoped operational reasoning and review protocol.

Use this protocol to answer:

- which contracts apply
- which requirement IDs materially constrain the sprint
- what evidence must exist before review can pass
- whether implementation deviated from selected contracts

---

# 1. Change Summary

## Sprint / Branch

```text
<branch-name>
```

## Feature / Task Name

```text
<short feature name>
```

## Goal

```text
<What is being built or changed?>
```

## Non-Goals

```text
- <not included>
- <not included>
```

---

# 2. Change Classification

## Change Type

```text
[ ] Backend API
[ ] Frontend
[ ] CLI
[ ] Package/Library
[ ] Workflow/Orchestration
[ ] LLM Runtime
[ ] Persistence/Schema
[ ] Infrastructure/Runtime
[ ] Operational/Observability
[ ] Release/Build
```

## Risk Level

```text
[ ] Low
[ ] Medium
[ ] High
[ ] Critical
```

## Runtime Impact

```text
[ ] User-facing runtime path
[ ] Background execution
[ ] Persistence
[ ] External provider/API
[ ] LLM/model interaction
[ ] Public API surface
[ ] CLI/public package surface
[ ] Auth/security-sensitive
[ ] High-cost path
```

---

# 3. Applicable Contracts

List every contract that applies to this change.

## Core Contracts

| Contract              | Applies? | Why      |
| --------------------- | -------- | -------- |
| architecture.md       | yes/no   | <reason> |
| errors.md             | yes/no   | <reason> |
| observability.md      | yes/no   | <reason> |
| testing.md            | yes/no   | <reason> |
| security.md           | yes/no   | <reason> |
| privacy-and-data.md   | yes/no   | <reason> |
| configuration.md      | yes/no   | <reason> |
| documentation.md      | yes/no   | <reason> |
| release-versioning.md | yes/no   | <reason> |

## Surface Contracts

| Contract              | Applies? | Why      |
| --------------------- | -------- | -------- |
| api-contracts.md      | yes/no   | <reason> |
| frontend.md           | yes/no   | <reason> |
| accessibility.md      | yes/no   | <reason> |
| cli.md                | yes/no   | <reason> |
| package-public-api.md | yes/no   | <reason> |

## Runtime Contracts

| Contract                      | Applies? | Why      |
| ----------------------------- | -------- | -------- |
| workflows.md                  | yes/no   | <reason> |
| llm.md                        | yes/no   | <reason> |
| llm-evaluation-cost-safety.md | yes/no   | <reason> |
| persistence-and-migrations.md | yes/no   | <reason> |
| performance.md                | yes/no   | <reason> |

---

# 4. Requirement Mapping

Map only the requirement IDs that materially affect the implementation.

## Example

| Requirement ID | Why It Applies                      | Planned Evidence               |
| -------------- | ----------------------------------- | ------------------------------ |
| ARCH-LAYER-002 | Use case depends on provider access | provider injected through port |
| ERR-CORE-001   | provider failures introduced        | structured failure emitted     |
| OBS-CORR-001   | async workflow boundary added       | trace/correlation propagation  |
| TEST-FAIL-001  | retry logic introduced              | retry exhaustion test          |

## Actual Mapping

| Requirement ID | Why It Applies | Planned Evidence |
| -------------- | -------------- | ---------------- |
| <ID>           | <reason>       | <evidence>       |
| <ID>           | <reason>       | <evidence>       |
| <ID>           | <reason>       | <evidence>       |

---

# 5. Architecture Reasoning Summary

## Existing Flow

```text
1. <step>
2. <step>
3. <step>
```

## Proposed Flow

```text
1. <step>
2. <step>
3. <step>
```

## Main Architectural Decision

```text
<What structure/boundary/abstraction choice matters most?>
```

## Why This Shape Was Chosen

```text
<Why this is the simplest honest design?>
```

## Alternatives Considered

```text
- <alternative>
- <why rejected>
```

## Main Trade-Off

```text
<what complexity is being accepted or avoided?>
```

---

# 6. Risk Analysis

## Main Risks

```text
- <risk>
- <risk>
```

## Operational Risks

```text
- <retry risk>
- <data integrity risk>
- <provider outage risk>
- <concurrency risk>
```

## Security / Privacy Risks

```text
- <risk>
```

## AI/LLM Risks (if applicable)

```text
- <hallucination>
- <tool misuse>
- <cost explosion>
- <prompt injection>
```

---

# 7. Testing Strategy

## Unit Coverage

```text
- <what is unit tested>
```

## Integration Coverage

```text
- <what is integration tested>
```

## Smoke/E2E Coverage

```text
- <what is smoke tested>
```

## Failure-Path Coverage

```text
- <what failures are tested>
```

## Real Provider Coverage (if applicable)

```text
- <provider smoke/eval>
```

---

# 8. Observability Plan

## Logs

```text
- <structured logs emitted>
```

## Events

```text
- <events emitted>
```

## Metrics

```text
- <metrics>
```

## Correlation

```text
- request_id
- trace_id
- run_id
- prompt_version
```

---

# 9. Release / Migration Impact

## Schema/Data Migration

```text
[ ] None
[ ] Safe additive
[ ] Destructive
[ ] Backfill required
```

## Public API/CLI Impact

```text
[ ] None
[ ] Backwards compatible
[ ] Breaking
```

## Config Changes

```text
- <new env vars/settings>
```

## Rollback Strategy

```text
<rollback or recovery plan>
```

---

# 10. Final Decision

## Decision

```text
[ ] Proceed
[ ] Refactor first
[ ] Split scope
[ ] Blocked pending clarification
```

## Why

```text
<short explanation>
```

## Complexity Added

```text
low / medium / high
```

## Complexity Removed

```text
low / medium / high
```
