# Architecture Review Protocol

> Role: post-implementation review protocol.
> Used by: `review.md` when selected by `sprint-index.md`.

This protocol checks whether implementation conforms to `sprint-reasoning.md`, preserves architectural fit, and avoids unearned complexity. It is not a planning template and must not introduce new architecture during review.

## Purpose

Use this after implementation, during code review, or when auditing an agent-generated change.

The review should not simply ask "does it work?" It should ask whether the implementation preserves or improves the long-term shape of the system.

---

## 1. Behaviour Review

### Main question

Can the reviewer clearly understand the feature's behaviour from the code?

Check:

```text
- [ ] The main workflow is visible.
- [ ] The feature does what the requirement asked.
- [ ] The implementation does not include unrelated work.
- [ ] Inputs and outputs are clear.
- [ ] Side effects are obvious.
- [ ] Failure paths are handled deliberately.
```

Review notes:

```text
<notes>
```

---

## 2. Architecture Fit Review

### Main question

Did the change fit the existing architecture, or did it distort it?

Check:

```text
- [ ] The change belongs in the files/modules it touched.
- [ ] The current workflow still reads cleanly.
- [ ] No feature was jammed into an unsuitable abstraction.
- [ ] No large refactor was avoided when the design clearly needed one.
- [ ] No large refactor was performed without real need.
```

Decision:

```text
[ ] Good fit
[ ] Acceptable fit with minor issues
[ ] Works but architecture is bending
[ ] Refactor required before merge
```

Review notes:

```text
<notes>
```

---

## 3. Simplicity and Earned Complexity Review

### Main question

Did the implementation add only the complexity required by the feature?

Check:

```text
- [ ] No speculative abstractions.
- [ ] No unnecessary interfaces/protocols/base classes.
- [ ] No plugin/factory/registry unless justified.
- [ ] No generic engine for a single concrete use case.
- [ ] No framework ceremony added without need.
- [ ] The simplest honest design was chosen.
```

Complexity verdict:

```text
[ ] Simpler than before
[ ] Same complexity
[ ] More complex, justified
[ ] More complex, not justified
```

Review notes:

```text
<notes>
```

---

## 4. Cohesion Review

### Main question

Are responsibilities grouped around real concepts and reasons to change?

Check:

```text
- [ ] Related rules are kept together.
- [ ] Unrelated behaviours are not forced into one unit.
- [ ] Functions/classes/modules have clear names and purposes.
- [ ] No god object/service/module was expanded.
- [ ] No excessive micro-fragmentation was introduced.
```

Cohesion issues:

```text
- <issue>
```

Recommended improvement:

```text
<improvement>
```

---

## 5. Coupling Review

### Main question

Did the change introduce unhealthy dependencies?

Check:

```text
Global coupling:
- [ ] No new mutable global state.
- [ ] Runtime configuration/state is passed explicitly where needed.

Content coupling:
- [ ] No external code reaches into internals/private fields.
- [ ] Public methods/interfaces are used appropriately.

Stamp coupling:
- [ ] Functions do not accept large objects when only small values are needed.
- [ ] Parameter objects are cohesive and purposeful.

Dependency coupling:
- [ ] External dependencies are injected or isolated at boundaries.
- [ ] Core logic is not tightly coupled to infrastructure.
```

Coupling verdict:

```text
[ ] Coupling reduced
[ ] Coupling acceptable
[ ] Coupling increased but justified
[ ] Coupling increased and should be fixed
```

Review notes:

```text
<notes>
```

---

## 6. DRY and Duplication Review

### Main question

Did the change duplicate knowledge or create a bad abstraction to avoid duplication?

Check:

```text
- [ ] Business rules are not duplicated.
- [ ] Security/authorization rules are not duplicated.
- [ ] Validation rules are not duplicated inconsistently.
- [ ] Similar code shape was not abstracted prematurely.
- [ ] Shared abstractions do not contain caller-specific branches.
- [ ] No new boolean flags were added to preserve a bad abstraction.
```

Duplication verdict:

```text
[ ] No concerning duplication
[ ] Duplication exists but is only code shape and acceptable
[ ] Business knowledge duplication must be fixed
[ ] Bad abstraction must be split
```

Review notes:

```text
<notes>
```

---

## 7. State and Side Effects Review

### Main question

Are state changes and side effects explicit, controlled, and testable?

Check:

```text
- [ ] Durable state changes are clear.
- [ ] Derived state is not treated as authoritative unless justified.
- [ ] Ephemeral state does not leak into global/shared state.
- [ ] Mutating operations are named clearly.
- [ ] Queries do not unexpectedly mutate state.
- [ ] Side effects are at boundaries where practical.
```

State verdict:

```text
[ ] Clear and safe
[ ] Minor concerns
[ ] Hidden mutation risk
[ ] Needs redesign
```

Review notes:

```text
<notes>
```

---

## 8. Function/Class/Paradigm Review

### Main question

Does the chosen coding style fit the problem?

Check:

```text
- [ ] Functions are used where behaviour is stateless/simple.
- [ ] Classes are justified by state, invariants, lifecycle, or polymorphism.
- [ ] Inheritance is shallow and honest if used.
- [ ] Composition is used where behaviour varies independently.
- [ ] Functional pipelines improve readability rather than obscure it.
- [ ] Framework conventions are followed without hiding domain behaviour.
```

Paradigm verdict:

```text
[ ] Good fit
[ ] Acceptable
[ ] Over-object-oriented
[ ] Over-functional/too dense
[ ] Over-procedural/no useful boundaries
[ ] Framework-driven rather than domain-driven
```

Review notes:

```text
<notes>
```

---

## 9. Error Handling Review

### Main question

Are failures explicit, useful, and recoverable where appropriate?

Check:

```text
- [ ] Expected failures have specific handling.
- [ ] Unexpected failures are surfaced clearly.
- [ ] Error names/messages are actionable.
- [ ] Errors are not swallowed silently.
- [ ] Retry behaviour is safe and deliberate.
- [ ] Partial progress is handled.
- [ ] AI/model/tool failures are validated and typed where relevant.
```

Error handling verdict:

```text
[ ] Strong
[ ] Acceptable
[ ] Too generic
[ ] Unsafe / incomplete
```

Review notes:

```text
<notes>
```

---

## 10. Observability Review

### Main question

Can we understand this feature in real execution?

Check:

```text
- [ ] Important workflow start/completion/failure is observable.
- [ ] Logs/events/metrics include correlation identifiers.
- [ ] Dependency calls are traceable where relevant.
- [ ] Duration/performance can be inspected.
- [ ] Errors include useful type/context.
- [ ] Sensitive data is not logged.
- [ ] AI-specific context is versioned/traced where relevant.
```

Observability verdict:

```text
[ ] Strong
[ ] Acceptable
[ ] Missing useful signals
[ ] Unsafe logging risk
```

Review notes:

```text
<notes>
```

---

## 11. Testing Review

### Main question

Do the tests protect the behaviour and reflect the architecture?

Check:

```text
- [ ] Pure logic has focused unit tests.
- [ ] Workflow/use-case behaviour is tested.
- [ ] External adapters are tested at the right level.
- [ ] Failure paths are tested.
- [ ] Regression cases are covered.
- [ ] Tests do not require unnecessary real infrastructure.
- [ ] Tests are not over-mocked to the point of meaninglessness.
```

Testing verdict:

```text
[ ] Strong
[ ] Acceptable
[ ] Missing important cases
[ ] Tests reveal architecture problems
```

Review notes:

```text
<notes>
```

---

## 12. Performance Review

### Main question

Are performance choices appropriate for the known constraints?

Check:

```text
- [ ] No premature optimization harmed clarity.
- [ ] Known scale/latency constraints were considered.
- [ ] Expensive operations are not repeated unnecessarily.
- [ ] Data loading is not accidentally excessive.
- [ ] Concurrency/race risks are considered where relevant.
- [ ] Measurement exists or is planned if needed.
```

Performance verdict:

```text
[ ] Fine
[ ] Needs measurement
[ ] Likely bottleneck
[ ] Over-optimized prematurely
```

Review notes:

```text
<notes>
```

---

## 13. Maintainability Verdict

### Would the next related feature be easier or harder after this change?

```text
[ ] Easier
[ ] About the same
[ ] Harder but justified
[ ] Harder and not justified
```

### Main risks left behind

```text
- <risk 1>
- <risk 2>
- <risk 3>
```

### Required changes before merge

```text
- <required change>
```

### Optional improvements later

```text
- <future improvement>
```

---

## 14. Final Review Decision

```text
Decision:
[ ] Approve
[ ] Approve with comments
[ ] Request small changes
[ ] Request architecture refactor
[ ] Reject / redesign required

Reason:
<short explanation>

Highest priority fix:
<one concrete action>
```
