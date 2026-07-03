# LLM Feature Reasoning Template

## 1. Change Summary

### AI capability

```
<what AI-backed behaviour is being added/changed>
```

User/system goal

<what outcome the AI feature should produce>

Runtime mode

[ ] Conversational agent
[ ] Structured worker
[ ] Retrieval/RAG
[ ] Tool-using agent
[ ] Classification/extraction
[ ] Summarization/generation
[ ] Evaluation/scoring
[ ] Other: <...>

---

## 2. Contract Mapping

| Contract                      | Requirement IDs | Why Applies |
| ----------------------------- | --------------- | ----------- |
| llm.md                        | <IDs>           | <reason>    |
| llm-evaluation-cost-safety.md | <IDs>           | <reason>    |
| errors.md                     | <IDs>           | <reason>    |
| observability.md              | <IDs>           | <reason>    |
| privacy-and-data.md           | <IDs>           | <reason>    |
| security.md                   | <IDs>           | <reason>    |
| testing.md                    | <IDs>           | <reason>    |
| performance.md                | <IDs>           | <reason>    |

---

## 3. Behaviour Contract

Input

<what context/input enters the model>

Output

<text/JSON/tool call/score/action/etc.>

Success criteria

<what makes output acceptable>

Failure criteria

<what must be rejected/retried/escalated>

---

## 4. Model and Provider Choice

Primary model/provider: <...>
Backup model/provider: <...>
Reason: <quality/cost/latency/tool support/structured output/etc.>

Model settings

temperature: <...>
max output/context: <...>
timeout: <...>
max attempts: <...>

---

## 5. Prompt and Context

prompt_id: <...>
prompt_version: <...>
prompt_purpose: <...>
context sources: <...>

Context minimization

<why this context is necessary and what is excluded>

Prompt injection risk

[ ] Low
[ ] Medium
[ ] High
Mitigation: <...>

---

## 6. Tools / Actions

Tools allowed: <...>
Tools forbidden: <...>
Approval needed? <yes/no/details>
Destructive or external side effects? <yes/no/details>

Tool argument validation

<schemas/checks before execution>

---

## 7. Structured Output / Validation

[ ] Not structured
[ ] JSON/schema/model output
[ ] Tool call output
[ ] Classification labels
[ ] Score/rubric

Local validation

<schema/model/rules>

Retry/correction behaviour

<retryable failures and max attempts>

---

## 8. Lifecycle and Inspectability

run_id/turn_id/task_id: <...>
states: <running/retrying/completed/failed/cancelled/etc.>
events persisted: <...>
terminal failure state: <...>

---

## 9. Safety and Human Review

High-impact? <yes/no>
Human approval/review needed? <yes/no/details>
Unsafe output handling: <...>
Unsupported claim handling: <...>

---

## 10. Cost and Latency

Expected cost drivers: <tokens/tools/retries/context/provider>
Limits: <...>
Metrics emitted: <...>

---

## 11. Evaluation and Testing

Regression examples: <...>
Smoke tests: <...>
Structured validation tests: <...>
Tool safety tests: <...>
Failure tests: <...>

---

## 12. Decision

Decision: <proceed/refactor/split/block>
Reason: <short rationale>
Main trade-off: <quality/cost/latency/safety/complexity>
