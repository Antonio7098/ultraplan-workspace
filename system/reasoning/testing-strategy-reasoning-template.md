# Testing Strategy Reasoning Template

## 1. Change Under Test

### Change

```
<feature/fix/refactor being tested>
```

Risk profile

[ ] Pure logic
[ ] API boundary
[ ] UI behaviour
[ ] Persistence
[ ] Provider/LLM
[ ] Workflow/background
[ ] Security/privacy
[ ] Package/CLI public behaviour
[ ] Migration/release

---

## 2. Contract Mapping

| Contract                      | Requirement IDs | Why Applies |
| ----------------------------- | --------------- | ----------- |
| testing.md                    | <IDs>           | <reason>    |
| errors.md                     | <IDs>           | <reason>    |
| observability.md              | <IDs>           | <reason>    |
| security.md                   | <IDs>           | <reason>    |
| api-contracts.md              | <IDs>           | <reason>    |
| frontend.md                   | <IDs>           | <reason>    |
| llm-evaluation-cost-safety.md | <IDs>           | <reason>    |

---

## 3. Behaviour to Protect

- <core behaviour>
- <invariant>
- <failure path>
- <public contract>

---

## 4. Test Levels Chosen

Unit tests: <what and why>
Integration tests: <what and why>
Smoke/E2E tests: <what and why>
Provider/eval tests: <what and why>
Contract/API/CLI/package tests: <what and why>

---

## 5. Test Doubles / Real Dependencies

Fake/stub/mock dependencies: <...>
Real dependencies used: <...>
Reason: <...>

---

## 6. Failure Paths

- <expected failure>: <test>
- <unexpected/wrapped failure>: <test>
- <retry/exhaustion/cancellation>: <test>

---

## 7. Determinism

Time controlled? <yes/no>
Randomness controlled? <yes/no>
Network isolated? <yes/no>
Provider nondeterminism handled? <yes/no>
No brittle exact wording? <yes/no>

---

## 8. Coverage Deliberately Deferred

- <deferred test>: <why acceptable>

---

## 9. Decision

Decision: <sufficient/add tests/revise>
Reason: <short rationale>
Main residual risk: <...>
