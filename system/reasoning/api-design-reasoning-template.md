# API Design Reasoning Template

## 1. Change Summary

### API change

```
<endpoint/method/event being added or changed>
```

User/client goal

<what the client is trying to achieve>

API audience

[ ] Internal backend only
[ ] Frontend-owned internal API
[ ] Public API
[ ] Partner API
[ ] SDK-facing API
[ ] Webhook/event consumer

---

## 2. Contract Mapping

| Contract            | Requirement IDs | Why Applies |
| ------------------- | --------------- | ----------- |
| api-contracts.md    | <IDs>           | <reason>    |
| security.md         | <IDs>           | <reason>    |
| errors.md           | <IDs>           | <reason>    |
| observability.md    | <IDs>           | <reason>    |
| testing.md          | <IDs>           | <reason>    |
| privacy-and-data.md | <IDs>           | <reason>    |
| performance.md      | <IDs>           | <reason>    |

---

## 3. API Shape

Proposed shape

<method/path or operation name>

Why this shape?

<why this is resource-oriented/RPC/event/streaming/etc.>

Alternatives considered

- <alternative>: <why not chosen>

Is this synchronous or asynchronous?

[ ] Synchronous request/response
[ ] Async job/run
[ ] Streaming
[ ] Webhook/event
[ ] Other: <...>

Long-running operation?

[ ] No
[ ] Yes — status surface: <endpoint/event/run record>

---

## 4. Request Contract

Inputs

- <field>: <type / required / meaning>

Validation rules

- <rule>

Caller-controlled fields rejected or ignored

- <field>

Size/rate/cost bounds

- max payload size: <...>
- max page size: <...>
- rate/concurrency limit: <...>

---

## 5. Response Contract

Success response

<response shape>

Error response

<canonical error shape / safe subset>

Nullability / optionality

<fields that may be missing/null and why>

Raw internal models exposed?

[ ] No
[ ] Yes — justification: <...>

---

## 6. Auth and Authorization

Authentication source

<session/token/service identity/etc.>

Authorization decision

<who may perform which action on which object>

Object/tenant boundary

<object ownership / tenancy checks>

---

## 7. Idempotency and Retry

Is this operation safe to retry?

[ ] Yes, naturally idempotent
[ ] Yes, with idempotency key/dedupe
[ ] No, must not retry automatically

Duplicate request behaviour

<what happens on repeated request>

---

## 8. Compatibility

Is this backwards compatible?

[ ] Yes
[ ] No — breaking change

Migration/versioning needed?

<versioning or migration plan>

---

## 9. Observability

- request/event emitted: <...>
- correlation fields: <...>
- metrics: <...>
- logs: <...>

---

## 10. Testing

- schema validation tests: <...>
- auth/authz tests: <...>
- error response tests: <...>
- idempotency tests: <...>
- compatibility tests: <...>

---

## 11. Decision

Decision: <proceed/refactor/split/block>
Reason: <short rationale>
Main trade-off: <trade-off>
