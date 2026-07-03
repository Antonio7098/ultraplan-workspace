# Error Handling Reasoning Template

## Purpose

This template guides structured reasoning about error handling strategy before implementation. The goal is explicit, actionable errors — not silent failures, not panic, not opaque strings.

Use this for: file I/O, network operations, data parsing, user input validation, configuration loading, or any operation that can fail partially or completely.

---

## 0. Feature Summary

### Feature / Component

`<short name>`

### What can fail

```text
- <failure 1>
- <failure 2>
- <failure 3>
```

### Failure impact

```text
- What does the user see?
- What state is left behind?
- Is the failure recoverable mid-operation?
```

---

## 1. Error Taxonomy

### Expected failures (normal flow)

```text
- <failure>: <recoverable? yes/no> — <how caller should handle>
- <failure>: <recoverable? yes/no> — <how caller should handle>
```

### Unexpected failures (exceptional)

```text
- <failure>: <how surfaced to caller>
- <failure>: <log or panic?>
```

### Error surface size

```text
[ ] Small — a few distinct error classes
[ ] Medium — multiple failure modes per operation
[ ] Large — complex error graphs, partial state
```

---

## 2. Error Chain Design

### Should errors wrap?

```text
[ ] Yes — %w wrapping at boundary, callers use errors.Is/errors.As
[ ] No — raw sentinel errors only
[ ] Hybrid — sentinel at leaf, wrapper at boundary
```

### Sentinel errors needed?

```text
[ ] Yes — define ErrNotFound, ErrInvalidInput, etc.
[ ] No — error strings only
[ ] Partial — only for expected failures
```

### Propagation pattern

```text
[ ] errors.New + fmt.Errorf at every layer
[ ] sentinel errors + wrap at boundary
[ ] custom error type with fields
```

---

## 3. Silent Failure Check

### Where could failures be swallowed?

```text
- <if err != nil { } without return/log>
- <defer without error check>
- <Background context with no cancellation>
```

### Are all error returns explicit?

```text
[ ] Yes — every failure is returned or deliberately swallowed
[ ] No — silent failures found: <...>
```

---

## 4. Panic Safety

### Are panics used instead of errors?

```text
[ ] No panics in this component
[ ] Yes — panics are recovery-wrapped at boundary
[ ] Yes — panics are intentional and not recovered
```

### If panics are used, why?

```text
<Explain why panic is the right choice vs. returning an error.>
```

---

## 5. Call-site ergonomics

### How will callers check errors?

```text
- errors.Is — <for which errors?>
- errors.As — <for which error types?>
- errors.IsAny / switch — <pattern if many types>
```

### Error messages for humans

```text
- Are messages actionable? <yes/no>
- Do they include context about where the failure occurred?
- Do they expose internal details (paths, internals) that should be hidden?
```

---

## 6. Error Wrapping Standard

### %w vs %s vs %v

```text
Use %w for sentinel error wrapping (errors.Is support).
Do NOT use %v for error values (loses chain).
Add "context: %w" pattern: fmt.Errorf("open file: %w", err)
```

### When to NOT wrap

```text
- At the original failure point where context is already clear
- When re-raising a wrapped error (avoid double-wrapping)
```

---

## 7. Testing Error Paths

### Error test patterns needed

```text
[ ] errors.Is for sentinel errors
[ ] errors.As for typed errors
[ ] Table-driven error tests
[ ] Chaos/resilience tests
```

### Minimum error test coverage

```text
- Happy path (no error)
- Each sentinel error path
- Unexpected error (bad input, corruption)
- Error at each layer (store → caller)
```

---

## 8. Final Decision

```text
Error strategy: <describe the final approach>

Sentinels:
- <list>

Wrapping standard:
- <pattern, e.g. "context: %w">

Call-site contract:
- callers use <errors.Is/errors.As/string match>

What this does NOT include:
- <not panic>
- <not third-party error libs>
- <not opaque error strings>
```
