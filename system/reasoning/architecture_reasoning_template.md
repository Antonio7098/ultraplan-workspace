# Architecture Reasoning Template

## Purpose

This template gives AI coding agents a structured reasoning process to complete before implementing a feature.

The goal is not to slow work down with ceremony. The goal is to prevent agents from blindly appending code, introducing premature abstractions, or forcing new requirements into an architecture that no longer fits.

Use this for non-trivial feature work, refactors, architecture-sensitive fixes, agent workflows, backend changes, data model changes, external integrations, or any change that affects state, boundaries, coupling, or long-term maintainability.

---

## 0. Feature Summary

### Feature name

`<short feature name>`

### User/product goal

Describe the outcome the user or system needs.

```text
<What should this feature achieve?>
```

### Current task

Describe the concrete implementation task.

```text
<What needs to be changed in the codebase?>
```

### Non-goals

List what should not be built now.

```text
- <not included>
- <not included>
- <not included>
```

---

## 1. First-Principles Breakdown

### Core behaviour

What is the smallest truthful description of the behaviour?

```text
<This feature receives/does/transforms/returns...>
```

### Inputs

```text
- <input 1>
- <input 2>
- <input 3>
```

### Outputs

```text
- <output 1>
- <output 2>
```

### Durable state

What persisted state is created, updated, deleted, or read?

```text
- Created: <...>
- Read: <...>
- Updated: <...>
- Deleted: <...>
```

### Ephemeral state

What temporary state exists only during execution?

```text
- <request context / intermediate result / temporary calculation>
```

### Derived state

What can be recalculated from other state? Is it stored or computed on demand?

```text
- <derived value>
- Source of truth: <durable state>
- Stored or computed: <stored/computed>
- Staleness risk: <low/medium/high>
```

### Side effects

```text
- Database write: <yes/no/details>
- File write: <yes/no/details>
- Network call: <yes/no/details>
- Queue/event emission: <yes/no/details>
- Email/notification: <yes/no/details>
- Logs/metrics/traces: <yes/no/details>
```

---

## 2. Existing Architecture Fit

### Contract applicability and phase gate

Which contract requirements apply to this exact change now, and which are deferred until a later roadmap gate?

```text
Current gate: <skeleton/local CLI/runtime/provider/batch workflow/stable release>
Applies now:
- <contract ID>: <why this behavior is in scope>
Deferred:
- <contract ID>: <later trigger, such as runtime execution, durable workflow state, public JSON stability, or release hardening>
```

### Where does this behaviour naturally belong?

```text
<module/package/file/layer/use case/workflow>
```

### Existing workflow affected

Describe the current flow before changing it.

```text
Current flow:
1. <step>
2. <step>
3. <step>
```

### Proposed new flow

```text
Proposed flow:
1. <step>
2. <step>
3. <step>
4. <new/changed step>
```

### Does the current architecture still fit?

Choose one:

```text
[ ] Yes — the feature fits cleanly into the existing shape.
[ ] Partially — small refactor needed before adding the feature.
[ ] No — the feature reveals that the existing architecture is the wrong shape.
```

### If it does not fit, why?

```text
- [ ] Requires awkward boolean flags
- [ ] Requires caller-specific branching
- [ ] Duplicates business rules
- [ ] Forces unrelated responsibilities into one unit
- [ ] Requires reaching into internals
- [ ] Requires global/shared mutable state
- [ ] Makes tests much harder
- [ ] Hides important side effects
- [ ] Other: <...>
```

### Refactor-before-feature decision

```text
Decision: <implement directly / small refactor first / larger refactor required>
Reason: <why>
```

---

## 3. Design Options Considered

List at least two options unless the change is trivial.

### Option A: Minimal direct implementation

```text
Description:
<...>

Pros:
- <...>

Cons:
- <...>

Risks:
- <...>
```

### Option B: Refactor then implement

```text
Description:
<...>

Pros:
- <...>

Cons:
- <...>

Risks:
- <...>
```

### Option C: New abstraction / boundary / component

Only include if there is real pressure for it.

```text
Description:
<...>

Pros:
- <...>

Cons:
- <...>

Risks:
- <...>
```

### Chosen option

```text
Chosen option: <A/B/C>
Reason:
<Why this is the simplest honest design for the current requirement.>
```

---

## 4. Abstraction Check

### Are we adding a new abstraction?

```text
[ ] No
[ ] Yes — interface/protocol/trait
[ ] Yes — base class
[ ] Yes — service/component/module
[ ] Yes — strategy/function parameter
[ ] Yes — data structure/DTO/config object
[ ] Yes — plugin/registry/factory
```

### If yes, why is it earned?

```text
- [ ] Multiple real implementations exist now
- [ ] A second implementation is very likely soon
- [ ] It isolates an external dependency
- [ ] It removes branching rather than adding it
- [ ] It protects a stable domain concept
- [ ] It makes testing simpler
- [ ] It creates a clear boundary between policy and mechanism
- [ ] Other: <...>
```

### If no, why are we keeping it concrete?

```text
<Explain why abstraction would be premature.>
```

### Bad abstraction smell check

```text
- [ ] Generic name like Manager/Handler/Processor/Common/Helper
- [ ] Only one implementation with no real volatility
- [ ] Boolean flags or mode switches
- [ ] Optional parameters for caller-specific behaviour
- [ ] Interface mirrors a concrete class rather than consumer needs
- [ ] Abstraction exists only to reduce line count
- [ ] Abstraction makes the flow harder to trace
```

If any are checked, reconsider the design.

---

## 5. DRY and Duplication Check

### Is any duplication being introduced?

```text
[ ] No
[ ] Yes — duplicated business/domain knowledge
[ ] Yes — duplicated code shape only
[ ] Yes — duplicated infrastructure mechanics
```

### If duplicating business knowledge, stop and consolidate

```text
Shared rule/knowledge:
<...>

Single source of truth should be:
<...>
```

### If duplicating code shape, is that acceptable?

```text
Reason duplication is acceptable:
<The behaviours look similar but change for different reasons.>
```

### If removing duplication, is the abstraction safe?

```text
- [ ] The duplicated code represents the same knowledge
- [ ] Call sites change for the same reason
- [ ] The new abstraction has a meaningful name
- [ ] It does not introduce flags or caller-specific branches
- [ ] It makes tests simpler
```

---

## 6. Coupling Check

### Dependencies required

```text
- <dependency 1>: <why needed>
- <dependency 2>: <why needed>
```

### Coupling risks

```text
Global coupling:
[ ] None
[ ] Introduced
[ ] Existing but unchanged
[ ] Existing and reduced

Content coupling:
[ ] None
[ ] Introduced
[ ] Existing but unchanged
[ ] Existing and reduced

Stamp coupling:
[ ] None
[ ] Introduced
[ ] Existing but unchanged
[ ] Existing and reduced
```

### Narrow dependency check

```text
Does each function receive only what it needs?
[ ] Yes
[ ] No — reason: <...>

Are large objects passed only when the full concept is needed?
[ ] Yes
[ ] No — reason: <...>
```

### Dependency injection check

```text
Are side-effectful dependencies created at the edge and passed inward?
[ ] Yes
[ ] No
[ ] Not applicable

If no, why is internal construction acceptable here?
<...>
```

---

## 7. State and Mutation Check

### Mutation points

List every place state changes.

```text
- <file/function/module>: <state changed>
```

### Is mutation explicit from names and flow?

```text
[ ] Yes
[ ] No — improvement needed: <...>
```

### Are queries and commands separated where practical?

```text
[ ] Yes
[ ] No — reason: <...>
```

### Could hidden state make this hard to debug?

```text
[ ] No
[ ] Yes — risk: <...>
```

---

## 8. Function/Class/Module Shape

### Primary unit of behaviour

```text
[ ] Function
[ ] Class/object
[ ] Module
[ ] Service/component
[ ] Pipeline/workflow
[ ] Adapter
[ ] Other: <...>
```

### Why this unit?

```text
<Explain why this shape fits the behaviour.>
```

### If using a class/object, what justifies it?

```text
- [ ] Maintains meaningful state
- [ ] Protects invariants
- [ ] Implements polymorphism
- [ ] Manages lifecycle
- [ ] Required by framework
- [ ] Groups cohesive behaviour
- [ ] Other: <...>
```

### If using inheritance, why not composition?

```text
<Explain why inheritance is stable and honest here.>
```

### If using composition, what components are combined?

```text
- <component>: <responsibility>
- <component>: <responsibility>
```

---

## 9. Error Handling Design

### Expected failures

```text
- <failure>: <how handled>
- <failure>: <how handled>
```

### Unexpected failures

```text
- <failure>: <how surfaced/logged>
```

### Error taxonomy

```text
Error types/codes/classes/results introduced or reused:
- <error name>: <meaning>
```

### Retry / recovery behaviour

```text
- Retry safe? <yes/no>
- Partial progress possible? <yes/no>
- Compensation needed? <yes/no/details>
```

---

## 10. Observability Design

### Events/logs/metrics/traces needed

```text
- <event/log/metric>: <when emitted>
- <event/log/metric>: <when emitted>
```

### Minimum useful fields

```text
- correlation_id / request_id / run_id
- actor/user/system id if safe
- operation name
- dependency name
- duration
- outcome
- error type
- version/config/model/prompt if relevant
```

### Sensitive data check

```text
Could this expose PII/secrets/user content?
[ ] No
[ ] Yes — mitigation: <redaction/hash/omit/sanitize>
```

---

## 11. Testing Strategy

### Unit tests

```text
- <test pure calculation/decision>
- <test validation/invariant>
```

### Use-case/workflow tests

```text
- <test main feature path with fake dependencies>
- <test failure path>
```

### Integration tests

```text
- <test real adapter/db/api if needed>
```

### Regression tests

```text
- <bug/edge case protected>
```

### Testability check

```text
Can the core behaviour be tested without real infrastructure?
[ ] Yes
[ ] No — why: <...>
```

---

## 12. Performance and Scale Assumptions

### Known constraints

```text
- Expected volume: <...>
- Latency sensitivity: <...>
- Memory sensitivity: <...>
- Concurrency concerns: <...>
```

### Assumptions being made

```text
- <assumption>
- <assumption>
```

### Measurement plan

```text
- [ ] Not needed for this change
- [ ] Unit benchmark
- [ ] Profiling
- [ ] Load test
- [ ] Production metric
```

### Optimization decision

```text
<Why we are or are not optimizing now.>
```

---

## 13. Implementation Plan

### Files likely to change

```text
- <file>: <reason>
- <file>: <reason>
```

### Step-by-step plan

```text
1. <step>
2. <step>
3. <step>
4. <step>
```

### Migration/backwards compatibility

```text
- Schema migration needed? <yes/no>
- API compatibility affected? <yes/no>
- Existing data affected? <yes/no>
- Feature flag needed? <yes/no>
```

### Rollback plan

```text
<How to revert safely if this breaks.>
```

---

## 14. Final Pre-Implementation Decision

```text
Decision: <Proceed / Refactor first / Split task / Stop and clarify>

Reason:
<short explanation>

Complexity introduced:
<none / low / medium / high>

Complexity removed:
<none / low / medium / high>

Main trade-off:
<what we are accepting and why>
```

---

# Final Standard

The agent should not merely ask:

> Where can I add this code?

It should ask:

> What shape does this behaviour naturally want, and does the current architecture still support that shape?

If the answer is yes, implement simply.

If the answer is no, refactor first.
