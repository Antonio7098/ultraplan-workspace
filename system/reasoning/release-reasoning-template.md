# Release Reasoning Template

## 1. Release Summary

### Release version/name

```
<version/tag/date>
```

Release type

[ ] App deployment
[ ] API release
[ ] Package release
[ ] CLI release
[ ] Frontend release
[ ] Prompt/model release
[ ] Migration release
[ ] Internal tooling release

Main change

<short summary>

---

## 2. Contract Mapping

| Contract                      | Requirement IDs | Why Applies |
| ----------------------------- | --------------- | ----------- |
| release-versioning.md         | <IDs>           | <reason>    |
| documentation.md              | <IDs>           | <reason>    |
| testing.md                    | <IDs>           | <reason>    |
| security.md                   | <IDs>           | <reason>    |
| api-contracts.md              | <IDs>           | <reason>    |
| package-public-api.md         | <IDs>           | <reason>    |
| cli.md                        | <IDs>           | <reason>    |
| persistence-and-migrations.md | <IDs>           | <reason>    |
| llm-evaluation-cost-safety.md | <IDs>           | <reason>    |

---

## 3. Compatibility Impact

Public API changed? <yes/no>
CLI output/exit codes changed? <yes/no>
Package exports changed? <yes/no>
Config changed? <yes/no>
Schema/data changed? <yes/no>
Prompt/model behaviour changed? <yes/no>

Versioning decision

Version bump: <major/minor/patch/calver/etc.>
Reason: <...>

---

## 4. Release Evidence

Tests run: <...>
Smoke tests: <...>
Provider/eval checks: <...>
Build/artifact checks: <...>
Security/dependency checks: <...>
Docs generated/updated: <...>

---

## 5. Migration / Operator Notes

Migrations: <...>
Config/env changes: <...>
Secrets/keys: <...>
Deployment order: <...>
Downtime: <...>
Backfill: <...>

---

## 6. Rollback / Recovery

Rollback safe? <yes/no/partial>
Rollback steps: <...>
Forward-fix path: <...>
Data recovery path: <...>

---

## 7. Changelog Entry

Added:

- <...>

Changed:

- <...>

Fixed:

- <...>

Removed:

- <...>

Deprecated:

- <...>

Security:

- <...>

Migration/Operations:

- <...>

---

## 8. Decision

Decision: <release/block/revise>
Reason: <short rationale>
Main release risk: <...>
