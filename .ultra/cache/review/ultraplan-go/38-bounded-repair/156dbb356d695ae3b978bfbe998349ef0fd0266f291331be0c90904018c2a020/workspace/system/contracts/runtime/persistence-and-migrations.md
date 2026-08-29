# Persistence And Migrations Contract

## Purpose

This contract defines how durable storage, schema changes, atomic writes, migrations, and data integrity must be handled.

It governs:

- durable schema and artifact ownership
- migrations and format compatibility
- atomic write or transaction boundaries
- data integrity
- read/query/scan boundaries
- seed/test data
- backup/rollback posture

## Scope

This contract applies when a change:

- adds or changes durable files, state formats, database tables, collections, indexes, constraints, or schemas
- changes persistence models, stores, repositories, or artifact writers
- changes migrations, compatibility behavior, seed data, or fixtures
- changes atomic write, transaction, or consistency guarantees
- adds caches, search indexes, vector stores, queues, or derived read models/artifacts

## Requirement Index

| ID | Title | Applies To | Severity If Violated |
| --- | --- | --- | --- |
| PERSIST-SCHEMA-001 | Durable schema must have clear ownership and meaning | state/artifact/schema changes | High |
| PERSIST-MIG-001 | Migrations and compatibility changes must be safe, ordered, and reviewable | migrations, format changes | Blocker |
| PERSIST-ATOMIC-001 | Atomic write or transaction boundaries must be explicit | writes/multi-step updates | High |
| PERSIST-INTEGRITY-001 | Integrity rules must be enforced at the right layer | constraints/rules | High |
| PERSIST-READ-001 | Read and query behaviour must be bounded and observable | reads/search/scans | Medium |
| PERSIST-DERIVED-001 | Derived stores must declare source of truth and staleness | caches/search/vector/read models/artifacts | High |
| PERSIST-FIXTURE-001 | Fixture, seed, and test data must not leak into production | fixtures/seeds/test data | High |
| PERSIST-RECOVERY-001 | Risky data changes need rollback or recovery strategy | destructive migrations/state changes | High |

## Core Principles

1. Durable state is a long-term contract.
2. Migrations must be boring, reviewable, and recoverable.
3. Atomic write and transaction boundaries are part of business correctness.
4. Storage layers should enforce invariants they are best placed to protect.
5. Derived stores are not automatically authoritative.
6. Performance-sensitive queries need bounds and evidence.
7. Test/seed data must stay out of production unless intentionally deployed.

## Requirements

### PERSIST-SCHEMA-001: Durable Schema Must Have Clear Ownership And Meaning

**Rule**
Every durable file, entity, table, collection, index, queue, vector collection, or stored artifact must have an owner and clear semantic meaning.

**Required**
- define which module/capability owns each durable structure
- document whether data is authoritative, derived, cache, audit, operational, or temporary
- avoid cross-module writes to another module’s owned storage without explicit boundary

**Forbidden**
- generic dumping-ground files, tables, or collections with mixed ownership
- storage structures whose lifecycle and source of truth are unclear

**Evidence**
- migration/model/schema/artifact naming and docs identify ownership

### PERSIST-MIG-001: Migrations And Compatibility Changes Must Be Safe, Ordered, And Reviewable

**Rule**
Schema, durable file, and data migrations must be deterministic, ordered, and safe for the target deployment model.

**Required**
- include migration files or explicit compatibility/rejection behavior for schema and durable format changes
- keep migrations idempotent where ecosystem supports it or clearly ordered where not
- separate schema changes from expensive data backfills when risk warrants
- consider zero-downtime compatibility when multiple app versions may run

**Forbidden**
- editing already-applied migrations in shared environments
- shipping model or durable format changes without corresponding migrations, compatibility handling, or explicit rejection
- relying on automatic schema sync in production without review

**Evidence**
- migration tests or dry-runs cover upgrade path

### PERSIST-ATOMIC-001: Atomic Write Or Transaction Boundaries Must Be Explicit

**Rule**
Writes that must succeed or fail together must share a clear transaction, atomic-write, or reconciliation boundary.

**Required**
- define owner of transaction/session/unit-of-work or atomic file replacement
- use temporary files plus same-directory rename for important file state where practical
- avoid hidden commits in low-level helpers unless explicitly part of the contract
- handle external side effects around transactions deliberately

**Forbidden**
- partial durable writes from a single business operation with no compensation, retry, or reconciliation plan
- committing inside helpers that callers expect to be atomic
- mixing durable state changes and external dependency side effects without idempotency/reconciliation thinking

**Evidence**
- use-case/store tests cover commit, rollback, atomic write, or recovery behavior where important

### PERSIST-INTEGRITY-001: Integrity Rules Must Be Enforced At The Right Layer

**Rule**
Important invariants must be enforced through the most reliable layer: file/schema validation, database constraints, domain rules, application validation, or all of them where appropriate.

**Required**
- use schema validation, unique, foreign key, not-null, enum/check constraints, or file format validation where storage-level enforcement is useful
- keep business rules in domain/use-case logic where they require context
- prevent race conditions for uniqueness and state transitions

**Forbidden**
- relying only on frontend validation for data integrity
- enforcing critical uniqueness only through pre-checks vulnerable to races
- swallowing integrity errors as generic internal failures

**Evidence**
- tests cover constraint violations and domain invalid states

### PERSIST-READ-001: Read And Query Behaviour Must Be Bounded And Observable

**Rule**
Reads, scans, and queries must have clear bounds, indexes, and operational visibility where data size can grow.

**Required**
- avoid unbounded filesystem scans, collection scans, or database queries in user-facing or high-frequency paths
- add indexes for important filter/sort patterns
- define pagination for collection reads
- surface slow query risks in review for critical paths

**Forbidden**
- loading entire directories, files, tables, or collections into memory without reason
- N+1 query patterns in important paths
- broad wildcard search without limits

**Evidence**
- read/query tests or profiling cover expected scale where relevant

### PERSIST-DERIVED-001: Derived Stores Must Declare Source Of Truth And Staleness

**Rule**
Caches, summaries, read models, search indexes, vector stores, materialized views, analytics tables, and generated artifacts must declare their source of truth and consistency model.

**Required**
- define how derived state is created, refreshed, invalidated, and repaired
- define whether stale reads are acceptable
- define fallback behaviour when derived store is unavailable

**Forbidden**
- treating derived data as authoritative without design decision
- leaving cache/search/vector inconsistencies with no rebuild path

**Evidence**
- docs or code identify source-of-truth and refresh path

### PERSIST-FIXTURE-001: Fixture, Seed, And Test Data Must Not Leak Into Production

**Rule**
Seed, fixture, demo, generated-example, and test data must be clearly separated from production runtime data.

**Required**
- gate demo/seed commands by environment
- avoid production credentials or real PII in fixtures
- make destructive seed/reset commands explicit

**Forbidden**
- running test fixtures automatically in production
- committing real customer data as fixtures

**Evidence**
- seed/reset paths are environment-gated and tested

### PERSIST-RECOVERY-001: Risky Data Changes Need Rollback Or Recovery Strategy

**Rule**
Destructive, irreversible, large, or high-risk data or artifact changes must define rollback, backup, or recovery strategy.

**Required**
- identify destructive migrations and data backfills
- define backup/restore or compensating migration where feasible
- validate assumptions before destructive changes

**Forbidden**
- deleting durable fields, files, columns, or data before all readers have migrated unless safe
- irreversible destructive changes with no recovery note

**Evidence**
- migration review includes rollback/recovery plan for high-risk changes

## Review Rejection Criteria

Reject a persistence change if it:

- changes durable models or formats without migration/compatibility handling
- leaves ownership/source-of-truth unclear
- risks partial writes with no atomic-write, transaction, or reconciliation plan
- relies only on frontend validation for durable integrity
- adds unbounded queries in growing data paths
- treats derived stores as authoritative by accident
- runs seed/test data in production paths
- makes destructive data changes with no recovery plan

---
