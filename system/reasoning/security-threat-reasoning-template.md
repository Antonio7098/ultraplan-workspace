# Security Threat Reasoning Template

## 1. Change Summary

### Security-relevant change

```
<what is changing>
```

Assets affected

- <data/system/capability/resource>

Trust boundary crossed?

[ ] User/browser -> backend
[ ] Public internet -> API
[ ] CLI args/files -> runtime
[ ] Model/tool output -> system action
[ ] Internal service -> external provider
[ ] Plugin/dependency -> runtime
[ ] Other: <...>

---

## 2. Contract Mapping

| Contract                      | Requirement IDs | Why Applies |
| ----------------------------- | --------------- | ----------- |
| security.md                   | <IDs>           | <reason>    |
| privacy-and-data.md           | <IDs>           | <reason>    |
| errors.md                     | <IDs>           | <reason>    |
| observability.md              | <IDs>           | <reason>    |
| testing.md                    | <IDs>           | <reason>    |
| llm-evaluation-cost-safety.md | <IDs>           | <reason>    |

---

## 3. Actor and Permission Model

Actors: <anonymous/authenticated/admin/service/model/agent/etc.>
Protected actions: <...>
Protected objects: <...>
Authorization rule: <...>

---

## 4. Threats Considered

[ ] Authentication bypass
[ ] Authorization bypass / IDOR
[ ] Injection
[ ] Prompt injection / tool misuse
[ ] SSRF / arbitrary network call
[ ] Path traversal / arbitrary file access
[ ] Secret exposure
[ ] Unsafe deserialization / dynamic execution
[ ] Dependency/supply-chain risk
[ ] Data exfiltration
[ ] Rate/cost abuse
[ ] CSRF/CORS/session risk
[ ] Other: <...>

---

## 5. Mitigations

Threat: <...>
Mitigation: <...>
Evidence/test: <...>

---

## 6. Input/Output Safety

Input validation: <...>
Output encoding/redaction: <...>
File/path constraints: <...>
Network destination constraints: <...>
Tool/action constraints: <...>

---

## 7. Secrets and Sensitive Data

Secrets touched? <yes/no>
PII touched? <yes/no>
Logging restrictions: <...>
Diagnostics restrictions: <...>

---

## 8. Failure Behaviour

Fail closed? <yes/no>
What happens on missing auth/config/permission? <...>
Error shape safe? <yes/no>

---

## 9. Testing

- unauthorized access tests: <...>
- wrong actor/object tests: <...>
- malformed input tests: <...>
- injection/path/network tests: <...>
- secret redaction tests: <...>

---

## 10. Decision

Decision: <proceed/revise/block>
Reason: <short rationale>
Residual risk: <accepted risk, if any>
