# Testing Contract

## Purpose
This contract defines the minimum testing expectations for application changes.

It governs:
- test seam requirements
- unit, integration, and smoke expectations
- failure-path coverage
- determinism and flakiness constraints
- public output and concurrency regression checks
- collaborator replacement strategy

## Scope
This contract applies when a change:
- adds or changes application logic
- changes dependency wiring or persistence
- changes task, workflow, runtime, or provider behavior
- changes public adapters or operational surfaces

## Requirement Index

| ID | Title | Applies To | Severity If Violated |
| --- | --- | --- | --- |
| TEST-SEAM-001 | Collaborators must be replaceable through public seams | all testable components | High |
| TEST-UNIT-001 | Business logic must have unit coverage | domain and use-case logic | Medium |
| TEST-INT-001 | Wiring and persistence changes must have integration coverage | composition, persistence, adapters | High |
| TEST-SMOKE-001 | Critical runtime paths must have smoke coverage | startup and key runtime flows | Medium |
| TEST-SMOKE-002 | Critical external runtime/provider paths must have real-dependency smoke coverage | real-runtime/provider paths | High |
| TEST-FAIL-001 | Failure paths must be tested explicitly | failure-producing changes | High |
| TEST-DET-001 | Tests must remain deterministic and non-brittle | all tests | Medium |
| TEST-CONTRACT-001 | Public API, CLI, and package contracts must have compatibility tests | public surface changes | High |
| TEST-E2E-001 | Critical user journeys must have end-to-end or scenario tests | critical user paths | High |
| TEST-MIGRATION-001 | Schema, config, and data migrations must have upgrade, downgrade, or compatibility verification | migration changes | High |

## Requirements

### TEST-SEAM-001: Collaborators Must Be Replaceable Through Public Seams

**Rule**
Major collaborators must be replaceable in tests without mutating private fields.

**Required**
- use constructor injection, builder overrides, or public bootstrap seams
- provide fake or test-double implementations for important ports
- keep container overrides explicit

**Forbidden**
- tests mutating private service fields
- tests depending on hidden container internals normal code should not know about

**Evidence**
- tests use fake ports or documented overrides
- no private patching is required to isolate business logic

### TEST-UNIT-001: Business Logic Must Have Unit Coverage

**Rule**
Domain logic and use-case logic must have deterministic unit coverage.

**Required**
- test invariants, branching rules, validation, and transformation logic
- keep unit tests fast and dependency-light

**Evidence**
- domain and use-case tests exist for changed business logic

### TEST-INT-001: Wiring And Persistence Changes Must Have Integration Coverage

**Rule**
Changes to persistence, composition, or adapter behavior must be exercised through integration tests.

**Required**
- test file/state formats, DB mappings, atomic writes, and transaction behavior when persistence changes
- test bootstrap wiring when composition changes
- test adapters through realistic boundaries where appropriate

**Evidence**
- integration tests cover changed runtime seams and wiring

### TEST-SMOKE-001: Critical Runtime Paths Must Have Smoke Coverage

**Rule**
Critical runtime paths must be validated through smoke tests or equivalent high-signal end-to-end checks.

**Required**
- keep smoke execution centralized
- cover startup and the most important runtime behaviors
- prefer stable behavior checks over exact phrasing assertions

**Forbidden**
- ad hoc duplicate smoke scripts with separate boot logic
- brittle provider smokes that assert cosmetic output wording

**Evidence**
- smoke suites exist for critical runtime paths
- smoke commands are documented and reusable

### TEST-SMOKE-002: Critical External Runtime/Provider Paths Must Have Real-Dependency Smoke Coverage

**Rule**
External runtime or provider paths that are intended to work against the real dependency in production must have real-dependency smoke coverage at the appropriate boundary.

**Applies when**
- adding or changing a runtime-backed or provider-backed path
- changing runtime/provider configuration, transport, approval, streaming, or lifecycle behavior
- changing an external-dependency path that is part of the supported operational surface

**Required**
- run/add at least one real-dependency smoke for critical supported runtime/provider-backed paths
- keep one cheap baseline real-dependency smoke available for fast verification
- use real-dependency smokes to verify durable behavior such as lifecycle state, approvals, persistence, streaming, or replay rather than cosmetic phrasing

**Forbidden**
- treating fake-runtime, fake-provider, or mocked integration coverage as sufficient for critical real-dependency behavior
- shipping changes to critical runtime/provider-backed paths with no real-dependency verification plan

**Evidence**
- documented real-dependency smoke suites exist for supported critical runtime/provider-backed paths
- review can point to the exact real-dependency smoke run or explicit justified deferral

### TEST-FAIL-001: Failure Paths Must Be Tested Explicitly

**Rule**
A change that introduces or modifies a failure mode must include explicit failure-path verification.

**Required**
- test at least one expected failure path when adding failure-handling logic
- test unexpected or wrapped failure behavior where the boundary matters
- verify the resulting status, error shape, or observable failure signal

**Evidence**
- tests cover negative cases, not just happy paths

### TEST-DET-001: Tests Must Remain Deterministic And Non-Brittle

**Rule**
Tests must prefer deterministic signals and avoid flaky or over-coupled assertions.

**Required**
- assert on stable behavior, structure, lifecycle, and state transitions
- keep runtime/provider-backed smokes focused on durable runtime behavior
- run race/concurrency checks for code that introduces meaningful parallelism where the ecosystem supports it

**Forbidden**
- fragile tests that depend on exact phrasing with no need
- tests that only pass because of hidden timing assumptions

**Evidence**
- assertions target stable fields and behavior
- retries, timing, and external dependencies are controlled or isolated in tests
- concurrent runtime paths have race/concurrency verification or a documented deferral

### TEST-CONTRACT-001: Public API, CLI, And Package Contracts Must Have Compatibility Tests

**Rule**
Public API, CLI command interfaces, and package public surfaces must have compatibility tests that verify contractual behavior is preserved across changes.

**Required**
- keep compatibility tests that verify request/response shapes, exit codes, output formats, and stable behaviors
- test that documented contracts are honored when implementation changes
- keep compatibility tests fast and CI-friendly
- use golden/output regression tests for stable public CLI output where they reduce review ambiguity
- provide an intentional update workflow for golden or snapshot fixtures so output changes are explicit

**Forbidden**
- changing public contract behavior without corresponding compatibility test coverage
- treating internal implementation tests as sufficient for public surface contracts

**Evidence**
- compatibility tests exist for public API, CLI, and package surfaces
- breaking contract changes are caught by test failures before release
- golden/snapshot fixture changes are reviewable when public output changes

### TEST-E2E-001: Critical User Journeys Must Have End-To-End Or Scenario Tests

**Rule**
Critical user journeys must have end-to-end or scenario tests that exercise the full path from transport to outcome.

**Required**
- identify critical journeys such as auth flows, core business operations, and high-stakes workflows
- cover end-to-end behavior through realistic scenario tests
- keep scenario tests focused on durable behavior rather than cosmetic details

**Forbidden**
- critical journeys with no scenario or e2e test coverage
- relying only on unit or integration tests for journeys where full-stack behavior matters

**Evidence**
- critical user journeys have documented scenario or e2e test coverage
- scenario tests run in CI or are available for pre-release verification

### TEST-MIGRATION-001: Schema, Config, And Data Migrations Must Have Upgrade, Downgrade, Or Compatibility Verification

**Rule**
Schema, config, and data migrations must have upgrade, downgrade, or compatibility verification that proves the migration behaves correctly and does not corrupt existing state.

**Required**
- test forward migration with realistic seed data
- test backward migration or rollback when supported
- verify compatibility between migration versions when the runtime supports it
- cover data integrity, constraints, and referential correctness for schema changes

**Forbidden**
- shipping schema or data migrations with no verification plan
- assuming migration success without checking the result
- skipping downgrade or compatibility testing for migrations that affect persistent state

**Evidence**
- migrations have documented test coverage for forward behavior
- downgrade or compatibility verification exists where applicable
- data integrity is verified for significant schema changes

## Review Rejection Criteria
Reject a change if it:
- requires tests to patch private fields to replace collaborators
- changes persistence or wiring with no integration coverage
- adds a meaningful failure path with no negative-case verification
- introduces brittle smoke or provider assertions that check cosmetic phrasing instead of durable behavior
- changes stable public CLI/API output without compatibility or fixture-based review coverage where applicable
