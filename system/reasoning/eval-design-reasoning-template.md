# Eval Design Reasoning Template

## 1. Eval Summary

### Eval name

```
<name>
```

Behaviour under test

<what AI behaviour is being evaluated>

Eval type

[ ] Golden examples
[ ] Scenario smoke test
[ ] Structured output validation
[ ] Rubric-based grading
[ ] Human review
[ ] Model comparison
[ ] Tool safety eval
[ ] Retrieval quality eval

---

## 2. Contract Mapping

| Contract                      | Requirement IDs | Why Applies |
| ----------------------------- | --------------- | ----------- |
| llm-evaluation-cost-safety.md | <IDs>           | <reason>    |
| llm.md                        | <IDs>           | <reason>    |
| privacy-and-data.md           | <IDs>           | <reason>    |
| testing.md                    | <IDs>           | <reason>    |
| observability.md              | <IDs>           | <reason>    |

---

## 3. Data Classification

Data source: <synthetic/internal/user-content/production trace/etc.>
Classification: <public/internal/user-content/PII/sensitive/etc.>
Retention: <...>
Access control: <...>
Redaction/anonymization: <...>

---

## 4. Coverage

Cases included

- happy path: <...>
- edge case: <...>
- failure case: <...>
- adversarial/safety case: <...>

Cases excluded

- <not covered and why>

---

## 5. Scoring / Pass Criteria

Pass condition: <...>
Failure condition: <...>
Allowed variation: <...>
Manual review needed? <yes/no>

Avoided brittle checks

<exact phrasing/cosmetic assertions avoided unless required>

---

## 6. Execution Plan

Run frequency: <per PR / pre-release / nightly / manual>
Models/providers covered: <...>
Cost estimate: <...>
Timeout/budget: <...>

---

## 7. Result Recording

Where results are stored: <...>
What metadata is captured: <model, prompt_version, provider, cost, latency, date, commit>

---

## 8. Decision

Decision: <adopt/revise/block>
Reason: <short rationale>
Main limitation: <what this eval does not prove>
