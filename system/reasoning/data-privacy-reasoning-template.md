# Data Privacy Reasoning Template

## 1. Data Flow Summary

### Data involved

```
<what data is collected/processed/stored/transmitted>
```

Feature need

<why the feature needs this data>

Data subjects / owners

<users/customers/orgs/etc.>

---

## 2. Contract Mapping

| Contract                      | Requirement IDs | Why Applies |
| ----------------------------- | --------------- | ----------- |
| privacy-and-data.md           | <IDs>           | <reason>    |
| security.md                   | <IDs>           | <reason>    |
| observability.md              | <IDs>           | <reason>    |
| errors.md                     | <IDs>           | <reason>    |
| llm-evaluation-cost-safety.md | <IDs>           | <reason>    |
| persistence-and-migrations.md | <IDs>           | <reason>    |

---

## 3. Data Classification

[ ] Public
[ ] Internal
[ ] User content
[ ] Personal data / PII
[ ] Sensitive personal data
[ ] Secret
[ ] Regulated
[ ] Operational metadata
[ ] Derived/non-authoritative

Fields

- <field>: <classification>

---

## 4. Minimization

Minimum data needed

<smallest useful set>

Data intentionally excluded

- <field/content excluded>

Could this be a reference/hash/summary instead of raw data?

<yes/no/details>

---

## 5. Storage and Retention

Stored? <yes/no>
Where stored: <...>
Authoritative or derived: <...>
Retention: <...>
Deletion behaviour: <...>
Cache/derived copies: <...>

---

## 6. Transmission / Provider Use

Sent outside system? <yes/no>
Destination/provider: <...>
Purpose: <...>
Fields sent: <...>
Redaction/minimization: <...>

---

## 7. Logging / Telemetry / Diagnostics

Logged? <yes/no>
Telemetry fields: <...>
Redaction/sanitization: <...>
Prompt/content logging: <...>

---

## 8. Access Control and Audit

Who can access: <...>
Admin/support access: <...>
Audit event needed? <yes/no>
Export/download allowed? <yes/no/details>

---

## 9. Testing

- sensitive field omission tests: <...>
- redaction tests: <...>
- deletion/retention tests: <...>
- provider payload tests: <...>

---

## 10. Decision

Decision: <proceed/minimize further/block>
Reason: <short rationale>
Residual privacy risk: <...>
