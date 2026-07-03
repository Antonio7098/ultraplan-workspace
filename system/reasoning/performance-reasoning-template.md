# Performance Reasoning Template

## 1. Change Summary

### Performance-sensitive path

```
<path/feature/job/component>
```

Why performance matters here

<latency/throughput/cost/user experience/scale reason>

---

## 2. Contract Mapping

| Contract                      | Requirement IDs | Why Applies |
| ----------------------------- | --------------- | ----------- |
| performance.md                | <IDs>           | <reason>    |
| observability.md              | <IDs>           | <reason>    |
| testing.md                    | <IDs>           | <reason>    |
| persistence-and-migrations.md | <IDs>           | <reason>    |
| llm-evaluation-cost-safety.md | <IDs>           | <reason>    |
| frontend.md                   | <IDs>           | <reason>    |

---

## 3. Expected Load / Budget

Expected volume: <...>
Latency budget: <...>
Throughput target: <...>
Memory constraints: <...>
Bundle-size impact: <...>
Provider/model cost sensitivity: <...>

---

## 4. Expensive Operations

- DB/query: <...>
- network/provider call: <...>
- model/LLM call: <...>
- file processing: <...>
- rendering: <...>
- loop/batch: <...>

---

## 5. Bounds

Timeout: <...>
Max retries: <...>
Max concurrency: <...>
Max batch/page size: <...>
Max file/payload size: <...>
Max tokens/context/tools: <...>

---

## 6. Caching / Derived Data

Cache used? <yes/no>
Cache key: <...>
TTL: <...>
Invalidation: <...>
Stale behaviour: <...>
Authorization-sensitive? <yes/no>

---

## 7. Measurement Plan

[ ] No measurement needed — low-risk path
[ ] Benchmark
[ ] Profiling
[ ] Load test
[ ] Bundle analysis
[ ] Core Web Vitals / UI perf check
[ ] Production metric

Evidence

<what result would prove acceptable performance>

---

## 8. Trade-Off

Clarity vs optimization: <...>
Simplicity vs scalability: <...>
Freshness vs cache speed: <...>
Cost vs quality/latency: <...>

---

## 9. Testing and Observability

Tests: <...>
Metrics: <...>
Logs/events: <...>
Alerts: <...>

---

## 10. Decision

Decision: <proceed/measure first/optimize/defer>
Reason: <short rationale>
Main risk: <...>
