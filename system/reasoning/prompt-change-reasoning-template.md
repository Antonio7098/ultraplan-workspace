# Prompt Change Reasoning Template

## 1. Prompt Change Summary

### Prompt identity

```
prompt_id: <...>
current_version: <...>
new_version: <...>
owner_kind: <agent/worker/tool/eval/etc.>
owner_id: <...>
purpose: <...>
```

Change type

[ ] New prompt
[ ] Prompt wording change
[ ] Output format change
[ ] Tool instruction change
[ ] Safety instruction change
[ ] Retry/correction prompt change
[ ] Context construction change
[ ] Prompt bundle change

---

## 2. Why Change?

<what problem or improvement motivates this change>

Expected behaviour improvement

<what should improve>

Behaviour that must not regress

<what must remain stable>

---

## 3. Risk Assessment

- output shape risk: <low/medium/high>
- hallucination risk: <low/medium/high>
- safety/tool risk: <low/medium/high>
- cost/latency risk: <low/medium/high>
- backwards compatibility risk: <low/medium/high>

---

## 4. Contract Mapping

| Contract                      | Requirement IDs | Why Applies |
| ----------------------------- | --------------- | ----------- |
| llm.md                        | <IDs>           | <reason>    |
| llm-evaluation-cost-safety.md | <IDs>           | <reason>    |
| privacy-and-data.md           | <IDs>           | <reason>    |
| observability.md              | <IDs>           | <reason>    |
| testing.md                    | <IDs>           | <reason>    |
| release-versioning.md         | <IDs>           | <reason>    |

---

## 5. Versioning and Propagation

Prompt version updated? <yes/no>
Prompt checksum changed? <yes/no>
Runtime state records prompt version? <yes/no>
Logs/events/traces include prompt reference? <yes/no>
Release notes needed? <yes/no>

---

## 6. Eval Plan

Examples to run

- <example/scenario>

Expected pass criteria

<semantic or structured expectations>

Exact wording assertions avoided?

[ ] Yes
[ ] No — reason exact wording matters: <...>

---

## 7. Rollback Plan

<how to revert to previous prompt version>

---

## 8. Decision

Decision: <proceed/revise/block>
Reason: <short rationale>
Main trade-off: <trade-off>
