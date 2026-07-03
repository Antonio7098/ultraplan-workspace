# Observability Reasoning Template

## 1. Runtime Path Summary

### Path/change

```
<runtime behaviour being added/changed>
```

Why observability matters

<what operators/devs need to understand later>

---

## 2. Contract Mapping

| Contract            | Requirement IDs | Why Applies |
| ------------------- | --------------- | ----------- |
| observability.md    | <IDs>           | <reason>    |
| errors.md           | <IDs>           | <reason>    |
| privacy-and-data.md | <IDs>           | <reason>    |
| workflows.md        | <IDs>           | <reason>    |
| llm.md              | <IDs>           | <reason>    |
| performance.md      | <IDs>           | <reason>    |

---

## 3. Questions This Path Must Answer

- Did it start? <signal>
- Did it complete? <signal>
- Did it fail? <signal>
- Why did it fail? <error/code>
- How long did it take? <metric/span>
- Which actor/request/run caused it? <correlation>
- Which dependency/provider/model/config/prompt was used? <metadata>

---

## 4. Signals

Logs

- <log event/code/fields>

Events

- <runtime/domain/operational event>

Metrics

- <counter/histogram/gauge>

Traces/spans

- <span name/attributes>

Diagnostics/health

- <surface/state exposed>

---

## 5. Correlation

request_id: <yes/no/source>
trace_id: <yes/no/source>
task_id/run_id: <yes/no/source>
user/actor id: <yes/no/safe form>
prompt/model/provider id: <yes/no>

---

## 6. Privacy/Safety

PII risk: <...>
Secret risk: <...>
Prompt/content logging risk: <...>
Redaction/sanitization: <...>

---

## 7. Alerting / Operator Action

Alertable? <yes/no>
Severity: <...>
Operator action: <...>
Runbook/doc needed? <yes/no>

---

## 8. Testing

- signal emitted test: <...>
- failure observability test: <...>
- correlation propagation test: <...>
- redaction test: <...>

---

## 9. Decision

Decision: <sufficient/revise/block>
Reason: <short rationale>
Main visibility gap: <...>
