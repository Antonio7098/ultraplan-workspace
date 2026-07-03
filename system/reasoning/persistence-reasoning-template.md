# Persistence Reasoning Template

## 1. Change Summary

### Durable state change

```
<what storage/schema/data behaviour changes>
```

Owner

<module/capability owning this data>

Storage type

[ ] Relational table
[ ] Document/collection
[ ] File/object storage
[ ] Cache
[ ] Search index
[ ] Vector store
[ ] Queue/event stream
[ ] Analytics/read model
[ ] Other: <...>

---

## 2. Contract Mapping

| Contract                      | Requirement IDs | Why Applies |
| ----------------------------- | --------------- | ----------- |
| persistence-and-migrations.md | <IDs>           | <reason>    |
| privacy-and-data.md           | <IDs>           | <reason>    |
| errors.md                     | <IDs>           | <reason>    |
| observability.md              | <IDs>           | <reason>    |
| testing.md                    | <IDs>           | <reason>    |
| performance.md                | <IDs>           | <reason>    |

---

## 3. State Classification

[ ] Authoritative durable state
[ ] Derived/read model
[ ] Cache
[ ] Audit/event state
[ ] Operational metadata
[ ] Temporary state

Source of truth

<what is authoritative>

Derived/stale behaviour

Can be stale? <yes/no>
Rebuild path: <...>
Invalidation/refresh: <...>

---

## 4. Schema / Data Model

Entities/tables/collections: <...>
Key fields: <...>
Indexes: <...>
Constraints: <...>
Relationships: <...>

Why this model?

<normalization/denormalization/read pattern/write pattern rationale>

---

## 5. Transaction and Consistency

Transaction owner: <...>
Atomic operation: <...>
Partial failure risk: <...>
External side effects involved? <yes/no/details>
Idempotency/dedupe: <...>

---

## 6. Migration Plan

[ ] No migration
[ ] Additive schema migration
[ ] Backfill
[ ] Destructive migration
[ ] Data migration
[ ] Index migration

Compatibility

Can old and new code run together? <yes/no>
Deployment ordering: <...>
Rollback/recovery: <...>

---

## 7. Privacy and Retention

Data classification: <...>
PII/sensitive fields: <...>
Retention: <...>
Deletion behaviour: <...>
Logging restrictions: <...>

---

## 8. Query and Performance

Main read patterns: <...>
Main write patterns: <...>
Expected volume: <...>
Pagination/bounds: <...>
Index coverage: <...>
N+1 risk: <...>

---

## 9. Testing

- migration tests: <...>
- repository/storage tests: <...>
- transaction rollback tests: <...>
- integrity tests: <...>
- query performance tests: <...>

---

## 10. Decision

Decision: <proceed/refactor/split/block>
Reason: <short rationale>
Main trade-off: <consistency/performance/simplicity/flexibility>
