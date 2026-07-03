# Frontend Reasoning Template

## 1. Change Summary

### Feature/UI change

```
<what is being added or changed>
```

User goal

<what the user is trying to do>

Surface affected

[ ] Route/page
[ ] Feature
[ ] Entity
[ ] Workflow
[ ] Shared utility
[ ] Design-system primitive/pattern
[ ] Form
[ ] Data display
[ ] Navigation/layout

---

## 2. Contract Mapping

| Contract            | Requirement IDs | Why Applies |
| ------------------- | --------------- | ----------- |
| frontend.md         | <IDs>           | <reason>    |
| accessibility.md    | <IDs>           | <reason>    |
| api-contracts.md    | <IDs>           | <reason>    |
| testing.md          | <IDs>           | <reason>    |
| privacy-and-data.md | <IDs>           | <reason>    |
| performance.md      | <IDs>           | <reason>    |
| observability.md    | <IDs>           | <reason>    |

---

## 3. Ownership Placement

Proposed location

src/<app/pages/features/entities/workflows/shared/design-system>/...

Why this owner?

<why the code belongs here>

Public export needed?

[ ] No
[ ] Yes — export from: <index/public API>

Cross-feature dependency?

[ ] No
[ ] Yes — through public export: <...>

---

## 4. UI Composition

Page/feature/workflow flow

1. <user sees/does>
2. <system loads/validates/updates>
3. <result shown>

Components introduced

- <component>: <responsibility>

Is this reusable?

[ ] No, feature-specific
[ ] Yes, but still feature-owned
[ ] Yes, entity-owned
[ ] Yes, design-system/domain-neutral
[ ] Yes, shared/domain-neutral

Reusability justification

<why abstraction is or is not earned>

---

## 5. State and Data

State categories

- server state: <...>
- URL state: <...>
- local UI state: <...>
- global client state: <...>
- derived state: <...>

API calls

- owning feature: <...>
- DTO mapping: <...>
- validation: <...>

---

## 6. Loading, Empty, Error, Success States

Loading: <what user sees>
Empty: <what user sees>
Error: <what user sees / retry path>
Success: <what user sees>
Pending/disabled: <what user sees>

---

## 7. Accessibility

- keyboard path: <...>
- focus behaviour: <...>
- labels/errors: <...>
- announcements/status: <...>
- color-only meaning avoided: <yes/no>
- reduced motion: <...>

---

## 8. Performance

- bundle impact: <none/small/meaningful>
- render risk: <...>
- duplicate fetch risk: <...>
- layout shift risk: <...>

---

## 9. Testing

- component tests: <...>
- hook/model tests: <...>
- accessibility tests: <...>
- interaction tests: <...>
- API mocking strategy: <...>

---

## 10. Decision

Decision: <proceed/refactor/split/block>
Reason: <short rationale>
Main trade-off: <trade-off>
