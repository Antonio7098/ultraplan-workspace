# Package API Reasoning Template

## 1. Package Change Summary

### Package/module

```
<package/module name>
```

Change

<what public/internal package behaviour changes>

Audience

[ ] Internal only
[ ] Public package users
[ ] CLI users
[ ] Plugin authors
[ ] SDK consumers

---

## 2. Contract Mapping

| Contract              | Requirement IDs | Why Applies |
| --------------------- | --------------- | ----------- |
| package-public-api.md | <IDs>           | <reason>    |
| release-versioning.md | <IDs>           | <reason>    |
| documentation.md      | <IDs>           | <reason>    |
| testing.md            | <IDs>           | <reason>    |
| security.md           | <IDs>           | <reason>    |

---

## 3. Public API Impact

Exports added: <...>
Exports changed: <...>
Exports removed: <...>
Behaviour changed: <...>
Entry points changed: <...>

Is this breaking?

[ ] No
[ ] Yes — why: <...>

---

## 4. API Surface Decision

Why public?

<why this symbol/command/config is stable enough to expose>

Why not internal?

<why users need this directly>

Minimal surface?

<what was intentionally not exported>

---

## 5. Dependency Impact

New dependencies: <...>
Runtime or dev only: <...>
Optional? <yes/no>
Transitive risk: <...>
License/maintenance: <...>

---

## 6. Versioning and Deprecation

Version bump: <major/minor/patch/etc.>
Deprecation needed? <yes/no>
Migration guidance: <...>

---

## 7. Build and Artifact

Clean install tested? <yes/no>
Package data included? <yes/no>
Types/docs included? <yes/no>
Secrets/dev files excluded? <yes/no>

---

## 8. Docs and Examples

Docs updated: <...>
Examples updated/tested: <...>
Changelog entry: <...>

---

## 9. Decision

Decision: <proceed/revise/block>
Reason: <short rationale>
Main trade-off: <API stability vs usefulness/etc.>
