# Frontend State Placement Template

## 1. State Summary

### State name

```
<state concern>
```

What does it represent?

<meaning of the state>

Who reads it?

- <component/feature/workflow>

Who writes it?

- <component/feature/workflow/server>

---

## 2. State Classification

[ ] Server state
[ ] URL/navigation state
[ ] Local ephemeral UI state
[ ] Form state
[ ] Global client state
[ ] Derived state
[ ] Cached state
[ ] Persisted browser state

---

## 3. Source of Truth

Authoritative source: <server/url/component/store/cache/etc.>
Derived from: <...>
Can be recomputed? <yes/no>
Can be stale? <yes/no>

---

## 4. Placement Decision

Chosen location

[ ] Component local state
[ ] Parent/nearest owner state
[ ] URL/search params/router state
[ ] Data-fetching cache/server-state library
[ ] Feature-local store
[ ] App/global store
[ ] Browser storage
[ ] Server persistence

Why this location?

<reason>

Why not global?

<explain if not global>

Why not local?

<explain if not local>

---

## 5. Synchronization Risk

- duplicate source of truth risk: <...>
- stale data risk: <...>
- race/update conflict risk: <...>
- optimistic update risk: <...>

---

## 6. Reset/Lifecycle Behaviour

Created when: <...>
Updated when: <...>
Reset when: <...>
Persisted across navigation? <yes/no>
Persisted across reload? <yes/no>

---

## 7. Testing

- placement behaviour test: <...>
- reset/lifecycle test: <...>
- stale/sync test: <...>

---

## 8. Decision

Decision: <state placement>
Reason: <short rationale>
Main trade-off: <trade-off>
