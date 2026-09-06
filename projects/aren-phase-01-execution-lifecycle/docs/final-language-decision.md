# Final implementation-language decision: Go

> **Status:** Accepted
>
> **Decision date:** 26 August 2026
>
> **Decision:** Aren's production runtime will be implemented in Go.
>
> **Scope:** The complete long-term runtime, including lifecycle control, model execution, tools, agent loops, persistence, daemon hosting, workflows, and any later distributed execution that is justified.
>
> **Evidence state:** Repository revision `ffa347e3dba6` plus the uncommitted prototype and documentation work described below.

## 1. Decision

Aren will use Go as its implementation language.

This is not a temporary choice for Phase 1 and it is not a plan to prototype in Go before rewriting in Rust. Go is the default language for the runtime across the roadmap. The production implementation will not maintain a shadow Rust runtime or require feature parity between two cores.

The Rust prototype remains valuable evidence. It demonstrated that Aren's semantics are compatible with Rust, exposed weaknesses in the Go design, and measured a real memory advantage. It now becomes a frozen comparison instrument rather than a second product implementation.

This decision applies to the runtime, not every program around it. Thin client SDKs may use their host language. External tools may remain separate processes. A future component may use another language behind a narrow protocol boundary if that component has an independent requirement. None of those cases reopens the runtime language by default.

## 2. Basis of the decision

The deciding question is:

> Which language best supports Aren's documented requirements and its method of developing through small, observable, failure-tested vertical slices?

The answer is Go because:

1. Both languages implemented the tested lifecycle and later roadmap slices without changing Aren's core contract. Rust is not required for capability.
2. Aren's defining correctness properties are runtime semantic properties. Neither compiler decided cancellation order, terminal arbitration, event order, cleanup completion, retry safety, or external effect truth.
3. Go kept the state machine and operating-system interactions more direct in the prototypes. That matters because Aren must be inspectable while its semantics are still evolving.
4. Go required less implementation and dependency machinery in the measured prototypes.
5. Go's standard library covers the current local library, CLI, HTTP, JSON, cancellation, process, profiling, and daemon needs.
6. Rust's memory advantage is real, but no current Aren workload or budget makes Go unacceptable. The strongest daemon measurements are preliminary and do not include the real agent process groups likely to dominate memory.
7. Choosing Rust now would pay certain complexity for possible future constraints. That conflicts with Aren's rule that complexity must be earned by demonstrated pressure.

The decision is therefore not “Go is always better than Rust.” It is that Go has the better cost structure for Aren as specified.

## 3. Aren requirements that govern the choice

Aren is not a generic service and the language was not selected from a general feature comparison. The following requirements define the relevant cost function.

| Aren requirement | Source | How Go supports it | Obligation created by choosing Go |
|---|---|---|---|
| One coherent terminal outcome | [Phase 1 lifecycle requirements](../phase-1-prd/02-lifecycle-requirements.md) | A mutex-protected transition gate can commit state, timing, event, and outcome in ordinary control flow. | All terminal paths must pass through one private commit boundary. No caller or executor may publish terminal state directly. |
| Truthful cancellation | [Phase 1 cancellation contract](../phase-1-prd/02-lifecycle-requirements.md#18-cancellation-requirements) | `context.Context` expresses a cooperative stop signal and is accepted by Go's HTTP and process APIs. | A request to cancel must remain distinct from confirmed termination. Every owned goroutine and external operation needs an explicit cancellation and join policy. |
| Ordered, observer-independent events | [Phase 1 event requirements](../phase-1-prd/02-lifecycle-requirements.md) | Private histories, sequence allocation under the transition gate, and pull-based cursors are direct to implement. | Observers must never execute inside the commit path. High-volume deltas, logs, and progress require explicit bounds rather than unbounded channels or histories. |
| Honest failure classification | [Phase 1 failure requirements](../phase-1-prd/02-lifecycle-requirements.md#17-failure-requirements) | Go error wrapping preserves inspectable causes while typed Aren failures distinguish origin and kind. | Public failure types, wrapping rules, and serialization must be defined centrally. Errors must not collapse into strings at durable or protocol boundaries. |
| Small, revisable vertical slices | [Roadmap development rules](../phased-roadmap.md#2-development-rules) | Go supports concrete implementations with little framework setup and allows abstractions to appear after repeated use. | Interfaces must be justified by actual variation. Standard-library use is preferred where it keeps behaviour clear, but dependencies are allowed when they remove proven risk. |
| Failure testing and real use | [Roadmap failure discipline](../phased-roadmap.md#24-failure-testing-is-mandatory) | The race detector, fuzzing, table-driven tests, subprocess tests, and black-box HTTP tests fit normal Go workflows. | Race-enabled CI is mandatory. Green unit tests are insufficient without adversarial, differential, leak, crash, and real-workload evidence. |
| Low, bounded, explainable overhead | [Performance engineering](../performance-engineering.md) | Go exposes runtime, heap, goroutine, mutex, block, trace, and profile data with standard tooling. | GC, allocation, runtime-task growth, retention, tail latency, saturation, and recovery must be measured throughout the roadmap. |
| Permission before effect | [Roadmap Phase 13](../phased-roadmap.md#phase-13--policies-and-resource-governance) | A direct call path can validate, authorize, reserve, and then execute without a separate policy framework. | Old or lower-level entry points must not bypass policy. Denial and admission decisions must be observable before the governed operation begins. |
| Local-first execution with an earned daemon | [Roadmap Phases 14 and 15](../phased-roadmap.md#hosting-and-external-clients) | Go works naturally as an in-process library and CLI, then as a small loopback HTTP daemon using `net/http`. | The daemon must be introduced only when caller-independent lifetime, shared inventory, persistence, or resource ownership requires it. |
| Thin multi-language clients | [Roadmap Phase 15](../phased-roadmap.md#phase-15--multi-language-clients-and-client-hosted-tools) | A stable protocol keeps lifecycle semantics in Go and clients small. | Client SDKs may project the wire contract but must not reimplement lifecycle resolution, retries, or policy. |
| Persistence without false exactly-once claims | [Roadmap Phase 11](../phased-roadmap.md#phase-11--persistence-and-recovery) | Go has mature file, database, and serialization options; the language does not dictate the recovery model. | State and event commits, schema migration, interruption, uncertainty, idempotency, and effect ledgers need explicit designs and crash tests. |
| Broader and distributed execution only when justified | [Roadmap Phases 17 and 18](../phased-roadmap.md#phase-17--broader-execution-types) | Go can supervise processes, HTTP work, daemons, and remote workers through concrete adapters and protocols. | Type-specific cancellation, output, success, cleanup, and effect semantics must remain visible. Distribution cannot be added for speculative scale. |

## 4. Evidence standard

Not all prototype evidence carries the same weight.

| Evidence | Weight in this decision | Reason |
|---|---|---|
| Mirrored lifecycle, failure, cancellation, and race cases | High for the named behaviours | They repeatedly found real defects and exercised the same semantic contract in both languages. |
| Current Go race-enabled suite and Rust safe-code compilation | High for what each tool checks | They confirm the present tested paths, while answering different concurrency questions. |
| Shared 14-case daemon contract | High for those contract cases | The same black-box suite runs unchanged against both daemons. |
| Defects discovered during extension and audit | Medium to high | They reveal concrete failure classes, but author sequence and familiarity prevent a language defect-rate comparison. |
| Current implementation size and dependency shape | Medium | They are useful maintenance indicators, not direct measures of comprehension or quality. |
| Repeated blocked-run footprint | Medium | It demonstrates a substantial memory difference in those two designs, but not full agent capacity. |
| Saved shared-daemon memory results | Preliminary | They used one process per cell on a loaded machine and predate scheduler and accounting repairs. |
| Initial microbenchmarks | Low for the current runtime | They describe an early Phase 1 shape and do not predict model, tool, or agent latency. |
| JavaScript daemon versus embedded Rust benchmark | Medium for deployment architecture only | It changed process boundary, lifetime, data movement, and language at once. |
| Future scale, FFI, WASM, or clustering arguments | None without a requirement | The roadmap explicitly treats later capabilities as hypotheses. |

This weighting prevents one favourable number or one type-system property from overriding the actual product requirements.

## 5. What the prototype programme established

The Go and Rust prototypes were disposable decision instruments. They were deliberately extended far beyond Phase 1 to test whether later capabilities would force a language-specific redesign.

### 5.1 Coverage and remaining limits

| Experiment | Tested capability | Evidence produced | Important limit |
|---|---|---|---|
| [Initial comparison](report.html) | First lifecycle implementation and microbenchmarks | Both languages expressed the basic run lifecycle; Go had lower full-run latency in the early shape and Rust cloned event history faster. | Early benchmarks used different frameworks and predate most runtime behaviour. |
| [Experiments 2 and 3](phase-1-experiment-2-cancellation-and-retention.md) | Cancellation, retention, observers, invariant injection, and stress | Both reached one-terminal, ordered-history, independent-observer, multi-waiter, and collision behaviour. | Hand-written schedules do not exhaust possible interleavings. |
| [Experiment 4](phase-2-experiment-4-progress-cleanup-and-conformance.md) | Progress, cleanup, partial output, and executor conformance | Cleanup occurs before terminal publication and cleanup failure remains visible. | Progress retention is unbounded and cleanup has no separately enforced deadline. |
| [Experiment 5](phase-3-experiment-5-first-model-invocation.md) | Cancellable model call and provider failures | Both separated provider transport, protocol, malformed response, rate limit, timeout, and cancellation. | Only Go has historical live-endpoint evidence; the current smoke was not rerun for this decision. |
| [Experiment 6](phase-4-experiment-6-streaming-model-execution.md) | Ordered model deltas and interruption | Both serialized callbacks, reconstructed output exactly, bounded retained deltas, and kept partial interrupted text diagnostic. | Go's default SSE adapter buffers the response; Rust has no provider network-stream adapter. Neither proves upstream backpressure. |
| [Experiment 7](phase-5-experiment-7-structured-results.md) | Structured result parsing and validation | Both distinguished provider success from task validity and preserved raw invalid output. | No provider-enforced schema, nested production schema, or repair loop. |
| [Experiment 8](phase-6-experiment-8-retry-and-backoff.md) | Attempts, selective retry, backoff, and exhaustion | Both implemented bounded retry policy, cancellable backoff, `Retry-After`, and explicit exhaustion. | No effect idempotency or complete durable attempt record. |
| [Experiment 9](phase-7-experiment-9-tool-call-representation.md) | Tool definitions and requested calls | Both preserved precise raw arguments, strict structure, order, IDs, and inspectable unknown tools. | No streamed tool arguments or provider interoperability matrix. |
| [Experiment 10](phase-8-experiment-10-local-tool-execution.md) | Tool resolution, validation, permission, and local files | Both revalidated arguments and denied before entering implementations. | No OS sandbox, forced stop, output bound, or TOCTOU-resistant filesystem boundary. |
| [Experiment 11](phase-9-experiment-11-bounded-agent-loop.md) | Minimal model and tool loop | Both enforced turn, tool-call, identity, and elapsed limits around a sequential transcript. | No real multi-turn task, token or cost budget, durable transcript, or parallel tool execution. |
| [Experiment 12](phase-10-experiment-12-context-management.md) | Context measurement and compaction | Both measured before model calls and compacted complete non-durable turns deterministically. | Byte limits stand in for token budgets; no evidence store, summarization, or provider-cache measurement. |
| [Experiment 13](phase-11-experiment-13-persistence-and-recovery.md) | Run snapshots, histories, and interruption | Both retained monotonic in-memory records and classified simulated loss honestly as interrupted. | No process-restart durability, fsync, schema evolution, or effect recovery. |
| [Experiment 14](phase-12-experiment-14-execution-composition.md) | Sequence and bounded parallel groups | Both reused child runs with parent cancellation and first-failure semantics. | No durable child evidence, compensation, fairness, or recovery. |
| [Experiment 15](phase-13-experiment-15-policies-and-resources.md) | Model and tool allowlists | Both denied observably before the tested effect. | Policy can be bypassed through older paths; core resource accounting was not part of this slice. |
| [Experiment 16](phase-14-experiment-16-daemon-and-clients.md) | Go daemon, thin clients, and embedded Rust comparison | The daemon established caller-independent ownership; the embedded path demonstrated a different lifetime and lower in-process call overhead. | The performance comparison was architecture plus language, not language alone. Neither path proved production recovery or security. |
| [Shared-daemon gauntlet](phase-16-experiment-17-shared-daemon-gauntlet.md) | One scheduler, admission, resource, retention, and HTTP contract in both languages | Both pass the same 14 black-box cases. Rust used less RSS in every saved matched cell. | Saved memory files predate three correctness repairs and need repeated runs on a quiet reference machine. |
| [Workflow experiment](phase-16-experiment-17-workflows.md) | Definition versus instance, sequential steps, children, and interruption | Both carried lifecycle semantics into a small workflow without a new core. | In-memory store, no daemon hosting, versioning, resume, human wait, or uncertain-effect recovery. |
| [Experiment 18](phase-17-experiment-18-broader-execution-types.md) | Shell and HTTP executions | Both reused the lifecycle and retained type-specific results and failures. | Neither guarantees descendant cleanup, sandboxing, consistent output bounds, streaming, or restart handling. |
| [Experiment 19](phase-18-experiment-19-distributed-and-remote.md) | Remote execution seam | Both crossed the existing daemon contract with failure classification and post-acknowledgement cancellation intact. | No workers, leases, heartbeats, placement, artifact transfer, fencing, or duplicate-effect protection. The need remains unproven. |

The breadth establishes feasibility and architectural continuity. It does not make the later prototype slices production-ready. Persistence is in memory, workflows are linear, policy is narrow, and distribution is only a seam.

### 5.2 Core semantic result

Both implementations preserved the same central lifecycle model through all tested extensions:

- one Aren-owned transition gate;
- exactly one terminal state, event, outcome, and timing record;
- cancellation acceptance distinct from confirmed stop;
- ordered observation that does not control execution;
- classified work, Aren, provider, validation, tool, workflow, shell, HTTP, and remote failures;
- explicit attempts and bounded loops;
- child execution using the same run semantics;
- policy decisions before the tested effects.

No roadmap slice required Rust's ownership model, and no slice exposed a capability Go could not implement directly.

That does not prove the lifecycle abstraction is complete. New executors continued to require type-specific cancellation, output, success, and cleanup rules. This supports the roadmap's plan to broaden through evidence rather than force a universal executor interface.

## 6. What the defects teach

The defect history is more valuable than a count of green tests.

### 6.1 Defects found in both designs

Both languages admitted semantic failures such as:

- cancellation bypassing Aren's acceptance path;
- terminal or progress facts observed in the wrong order;
- completed runs retained by watchers or elapsed-budget tasks;
- stream event order diverging from reconstructed output order;
- permits or reservations not released on early cancellation;
- HTTP cancellation stopping at headers rather than body reads;
- subprocess pipes capable of deadlocking when output was not drained;
- false-positive tests that did not exercise the path they claimed to prove.

These are central Aren risks. They were controlled by the lifecycle contract, a single mutation boundary, explicit ownership, and adversarial tests, not by the language alone.

### 6.2 Rust's protection was real but partial

Rust improved several properties:

- closed state, event, and failure vocabularies use enums and exhaustive matching;
- immutable shared outcomes are natural through ownership and `Arc`;
- `Send` and `Sync` constrain values crossing task boundaries;
- many unsafe sharing errors fail to compile;
- there is no tracing garbage collector.

Rust did not prevent logical deadlocks, task retention, incorrect event order, premature cancellation claims, cancellation-unsafe operations, semaphore leaks, or external-process cleanup defects. Aren's hardest guarantees remained runtime proofs.

### 6.3 Go's weaknesses are accepted, not ignored

The Go prototype exposed defects that Rust's types made less likely, including a mutable outcome shared across waiters and mutable schemas or maps retained across asynchronous execution. Go also admitted a data race in a cleanup window, which the race detector found once the test reached it.

Choosing Go accepts that immutable publication, closed vocabularies, and safe sharing rely more heavily on API design and verification. The engineering rules in section 11 are therefore part of the decision, not optional follow-up work.

The experiment history cannot support a defect-rate comparison. The same author built and repaired the implementations sequentially, and later work benefited from earlier findings.

## 7. Specific weaknesses uncovered and required responses

The table separates the prototype in which a weakness was observed from the lesson Aren must carry into the production Go runtime. Rust-specific findings remain relevant because they identify semantic failures that Go's compiler will not prevent either.

| Weakness uncovered | Observed in | Consequence demonstrated or exposed | Required response in the Go runtime |
|---|---|---|---|
| Parent cancellation bypassed Aren's cancellation-acceptance path. | Go and Rust | A run could receive a stop signal without retaining the first reason or recording the required request event through the authoritative path. | Route explicit, parent, deadline, daemon, workflow, and remote cancellation through one idempotent acceptance function before signalling work. Test every source against the same disposition and event contract. |
| A terminal cancellation event could be recorded before `run.cancellation_requested`. | Rust | History could contradict the causal order promised by Aren even though the final state was plausible. | Commit cancellation acceptance and its event under the transition gate before exposing the signal. Assert causal ordering under repeated completion and cancellation races. |
| Waiters shared a mutable terminal outcome. | Go | One consumer could mutate data observed by another, violating immutable publication and identical-outcome guarantees. | Keep outcome fields private, deep-copy nested collections at commit time, and return values or defensive snapshots. Add mutation attempts to waiter and observer tests. |
| The invariant-failure path recursed, double-closed a channel, and temporarily published incomplete terminal data. | Go | The error path intended to protect coherence could panic or expose state without its matching event and outcome. | Make invariant handling non-recursive. Use one terminal commit routine, one close owner, and a minimal fallback that cannot re-enter ordinary transition logic. Exercise it through explicit fault injection. |
| Terminal outcomes were published without terminal timing. | Go | Consumers could observe a completed run whose outcome violated the Phase 1 terminal record contract. | Construct terminal state, timing, event, and outcome together inside the commit boundary. Reject any terminal constructor input without complete timing. |
| Cleanup introduced an unlock window and a data race before final publication. | Go | Concurrent inspection could observe torn lifecycle facts while cleanup and terminal resolution competed. | Give one owner the classified work result, cleanup result, and final commit. Do not expose an intermediate terminal state. Keep the full suite under the race detector and stress inspection during cleanup. |
| Schemas, argument maps, and work slices were retained after callers could mutate them. | Go | Asynchronous execution could observe a definition different from the one validated at submission. | Validate and deep-copy all mutable caller-owned inputs before acceptance. Store typed immutable forms internally and test caller mutation immediately after submission. |
| Default JSON number decoding lost integer precision. | Go | Tool and structured-result arguments could change value while still looking syntactically valid. | Decode with `json.Decoder.UseNumber`, validate numeric ranges explicitly, reject trailing values, and preserve raw arguments alongside the validated form. |
| Classified Aren failures were flattened into ordinary errors on some paths. | Go | Failure origin and kind could be lost, weakening retry, diagnosis, and protocol behaviour. | Use one typed failure representation with wrapped local causes. Centralize conversion at provider, tool, workflow, persistence, and wire boundaries. Test origin, kind, message, and cause separately. |
| Cancellation during an HTTP body read was classified as an HTTP failure. | Go | The terminal record could lie about why work ended even though the context had been cancelled. | Check context cause around request dispatch, streaming, and body reads before assigning a transport classification. Test cancellation before headers, during headers, and during body consumption. |
| Completed runs were retained by a watch receiver or elapsed-budget task. | Rust | Terminal runs remained reachable after clients and work had finished. | Give observers, timers, and budget goroutines explicit release paths. Cancel and join them on terminal commit, then verify reclamation with leak and retention tests. |
| Watch-channel liveness and re-locking a non-reentrant mutex caused deadlocks. | Rust | Safe memory access did not guarantee forward progress or a valid lock protocol. | Define lock ordering, keep blocking calls and callbacks outside locks, avoid nested acquisition of the transition lock, and add timeout-backed deadlock tests plus model-based schedule exploration. |
| Stream callback order and reconstructed output order could diverge. | Rust | Event evidence and the final result could describe different executions. | Serialize every accepted delta through one ordering owner, assign sequence numbers there, and assert that retained deltas reconstruct the exact final output under concurrent callbacks. |
| A semaphore permit was not released on an early parallel-cancellation path. | Rust | Bounded composition could lose capacity permanently and eventually stall. | Represent every admission and reservation with a scoped release token. Test success, failure, panic, cancellation, and partial-start exits for exact resource reconciliation. |
| Provider streaming was simulated but not live and incrementally backpressured. | Go and Rust | Passing callback tests did not prove network streaming, slow-consumer behaviour, bounded buffering, or cancellation of an active provider body. | Build a real incremental SSE or provider adapter in Go with bounded buffers. Test fragmented frames, slow consumers, malformed late frames, cancellation, upstream closure, and retained-byte limits. |
| Subprocess output could deadlock, remain unbounded, or leave descendants running. | Go and Rust | A direct child could fill a pipe, consume unbounded memory, or terminate while grandchildren survived. | Drain stdout and stderr concurrently, enforce explicit byte limits, and supervise process groups on Unix and Job Objects on Windows. Test infinite output, ignored signals, forked descendants, and cleanup deadlines. |
| Progress histories, output, and daemon tombstones could grow without a lifetime bound. | Go and Rust | Long-lived processes accumulated memory even after heavyweight run data was evicted. | Separate durable lifecycle facts from high-volume data, define per-run and process-wide limits, expire or externalize tombstones, and prove recovery toward baseline in churn and soak tests. |
| Policy could be bypassed through older entry points. | Go and Rust | A correct allowlist on one path did not create an authoritative security or resource boundary. | Make validation, policy, admission, reservation, and execution one unbypassable orchestration path. Keep lower-level executors private or require an already-authorized capability token. |
| Several passing tests contained impossible assertions, empty helpers, or counters unrelated to the behaviour under test. | Go and Rust | Green mirrored suites overstated evidence and allowed defects to survive. | Add negative controls, seeded historical defects, mutation testing where practical, and a language-neutral reference model. Prefer requirement-to-test mapping over cumulative test counts. |

These responses are incorporated into the mandatory Go rules in section 11. They should also seed the production test plan so the prototype's failures remain regression cases rather than historical notes.

## 8. Performance and resource evidence

### 8.1 Latency

The original full-run benchmarks favoured Go for the early lifecycle shape, while Rust favoured event-history cloning. Corrected completed-wait measurements were effectively tied at about 40 nanoseconds.

Those operations are not representative of an agent runtime dominated by provider, network, process, storage, and human waits. They remain regression baselines, not a language decision.

### 8.2 Active-run memory

The repeated 10,000 blocked-run probe is the strongest result in Rust's favour:

- Go mean RSS increase per active run: about 9,604 bytes;
- Rust mean RSS increase per active run: about 2,155 bytes;
- observed ratio: about 4.5 to 1 in favour of the tested Rust design.

The designs were not carrier-identical. Go used three goroutines per run and Rust used two Tokio tasks. RSS also includes stack policies, allocators, schedulers, and runtime state. The result proves that this Rust design was materially lighter, not that every Rust Aren design will use 4.5 times less memory.

No documented Aren workload or budget makes the Go result unacceptable. At 10,000 blocked lifecycle runs the difference is roughly 71 MiB. A real coding-agent run may retain transcript, tool, provider, artifact, and child-process data that outweighs the lifecycle carrier.

### 8.3 Shared daemon

The saved shared-daemon runs reported Rust using 33 to 54 percent less RSS across matched cells. At two simulated 20 MiB runs Rust saved about 43 MiB. At eight, the recorded peaks were about 318 MiB for Go and 206 MiB for Rust.

Resource policy had a larger absolute effect in the same saved data. Reducing the active limit from eight to two saved about 230 MiB in Go and 161 MiB in Rust. This supports Aren's requirement for bounded admission regardless of language.

The shared results are not decision-grade performance measurements. They contain one daemon process per cell, came from a machine with active swap, sampled short workloads coarsely, lack binary digests, and predate repairs to scheduler weighting, reservation accounting, and overflow handling. They establish the direction of Rust's advantage, not its production value.

### 8.4 Performance conclusion

The decision does not deny Rust's memory result. It judges that the known benefit does not justify Rust's continuing complexity without a breached Aren capacity or latency requirement.

Go must now prove its performance against real workload budgets. If Go later misses a hard budget, the first response is measurement, retention control, admission, profiling, and focused design. The answer is not an automatic rewrite.

## 9. Implementation and operational cost

The current prototype shapes provide useful maintenance indicators:

| Component | Go | Rust | Observation |
|---|---:|---:|---|
| Core implementation | 3,567 lines | 4,425 lines | Rust is about 24% larger in this implementation. |
| Shared daemon implementation | 1,532 lines | 1,626 lines | Source sizes are close; the dependency models differ. |
| Diagnostic CLI | 424 lines | 694 lines | Rust is about 64% larger in this implementation. |

The counts exclude tests and generated build output. They are not a measure of quality or developer intelligence. Rust syntax, enums, conversions, and explicit ownership naturally use more code. Go can hide obligations in conventions. The result that matters is that Rust's extra machinery remained present as the prototype expanded.

The Go prototype declares no external module. The Rust prototype declares Tokio, Tokio Utilities, Futures, Serde, Serde JSON, Reqwest, and Axum, resolving a much larger transitive graph. Those are mature libraries that provide real value. They also introduce upgrade, advisory, feature, compile-time, and API-lifecycle work before Aren has demonstrated a need for that stack.

Go's standard library made several relevant paths direct:

- HTTP server and client behaviour;
- JSON request and response handling;
- context propagation into network calls;
- direct-child process cancellation;
- profiling and runtime inspection;
- race-enabled test execution.

Go does not remove the need for careful design. Its advantage is that more of Aren's state machine remains ordinary, reviewable control flow while the product is still learning what belongs in that state machine.

## 10. Alternatives considered

### 10.1 Rust for the permanent runtime

Rust was the only serious alternative. It offers stronger static representation and lower memory in the measured designs. It would be the better choice if Aren already had a hard no-GC requirement, a mandatory native-embedding requirement across several hosts, strict memory density that Go failed, or a team and codebase whose Rust async model made the lifecycle clearer rather than less clear.

Those are not current requirements. Selecting Rust would therefore front-load complexity for predicted pressure while leaving Aren's defining runtime proofs in tests and design.

### 10.2 Go now, Rust later

Rejected. This would make early product knowledge disposable, create two histories, encourage architecture for a port rather than for users, and defer the actual decision.

### 10.3 Two production cores

Rejected. Mirrored prototypes were useful for experiments. Mirrored production runtimes would double implementation, review, conformance, security, and operational work while making neither implementation authoritative.

### 10.4 A split Rust core with a Go daemon

Rejected as the default. It adds an FFI or internal protocol boundary inside the most correctness-sensitive path without a demonstrated requirement. A separately useful Rust component may be proposed later behind a narrow contract, but lifecycle authority remains in Go.

### 10.5 Delay the decision for more benchmarks

Rejected. The missing tests are important engineering work, but no current result identifies a capability Go lacks or a requirement it fails. Keeping the language open would itself carry cost: duplicated prototypes, delayed repository foundation, and continued re-litigation based on hypothetical futures.

## 11. Mandatory engineering rules for the Go runtime

The following rules are consequences of the decision.

### 11.1 One owner for lifecycle mutation

- Run fields are private.
- Every lifecycle transition passes through one commit function or one clearly defined ownership loop.
- State, event sequence, timing, outcome, and waiter release become visible coherently.
- Executors return facts to the owner. They do not mutate run state.

### 11.2 Published data is immutable by contract

- Outcomes, events, schemas, tool arguments, work definitions, and nested collections are copied or frozen before asynchronous retention or publication.
- Accessors return values or defensive snapshots, not mutable internal maps, slices, pointers, or `any` containers.
- Durable formats use explicit types and versions.

### 11.3 Goroutines have owners

- Every Aren-owned goroutine has a named owner, cancellation source, termination condition, and join or release path.
- `context.Background()` is permitted only at an explicit lifetime boundary such as daemon-owned accepted work.
- Fire-and-forget work is prohibited unless its bounded lifetime and irrelevance to correctness are documented.
- Leak tests cover completed runs, abandoned observers, timers, blocked providers, tools, and subprocesses.

### 11.4 Cancellation remains a protocol

- Context cancellation is a signal, not proof of stop.
- The first accepted reason is retained and repeated requests are idempotent.
- Terminal cancellation is published only after the executor-specific stop condition is satisfied.
- Process-tree, HTTP-body, tool, provider, workflow, and remote cancellation receive separate conformance cases.

### 11.5 Go's open type system is narrowed deliberately

- Lifecycle states, event types, failure origins, and failure kinds use named private or closed packages with validated constructors.
- Conversion and serialization switches are centralized.
- CI uses an exhaustive-switch analyzer for closed vocabularies where practical.
- Unknown protocol values remain explicit rather than becoming zero values.

### 11.6 Concurrency testing is part of correctness

- The complete Go suite runs under `go test -race` in CI.
- Stress tests repeat completion and cancellation collisions, concurrent inspection, observers, cleanup, and resource release.
- Fuzzed or model-based command traces compare actual state and history with a small reference model.
- Tests prove their instruments can fail through seeded historical defects, mutation testing, or deliberately broken fixtures.

### 11.7 Work and retention are bounded

- Every queue, event class, stream, tool output, transcript, artifact, cache, retry loop, and tombstone policy declares a bound or a measured reason not to have one.
- Admission reserves resources before effects and releases them on every path.
- Slow observers cannot block execution.
- Resource estimates and measured use remain distinct and overflow-safe.

### 11.8 Performance is measured at product scale

- Each phase records incremental allocations, heap or RSS, goroutines, contention, p50/p95/p99 latency where relevant, saturation, and recovery.
- Benchmarks state revision, diff status, binary digest, toolchain, hardware, workload, repetitions, and spread.
- Profiles precede optimization.
- Real process-group memory is measured when child agents exist; daemon RSS alone is not treated as capacity.

### 11.9 Errors and effects stay honest

- Aren failure origin and kind remain machine-readable.
- Wrapped causes remain inspectable locally; durable and wire representations state what is lost.
- Retry eligibility is explicit and bounded.
- External effects are never described as exactly once without an idempotency or effect mechanism.
- Recovery distinguishes completed, interrupted, and uncertain work.

### 11.10 Dependencies and abstractions must earn their cost

- Standard-library code is preferred when it is clear and sufficient.
- A dependency is accepted when it removes a demonstrated risk or substantial maintained code, not to satisfy a speculative architecture.
- Interfaces follow real variation or a required test boundary.
- Phase reviews may collapse packages, remove options, or delete abstractions.

## 12. Known gaps after the decision

The language decision does not close the following engineering questions:

1. Neither prototype has executed a real, repeated, multi-turn agent workload combining a live provider stream, tools, context growth, child processes, persistence, cancellation, and reconnect.
2. Neither has production persistence or crash consistency.
3. Neither has a complete tool sandbox, process-tree guarantee, consistent output bound, or cross-platform containment proof.
4. Neither has demonstrated live incremental provider streaming with backpressure and provider-enforced structured output.
5. The shared-daemon memory experiment needs corrected, repeated measurements on a quiet reference machine with real process groups.
6. Scheduler weights are contract-tested as configuration but not yet proved under long saturated admission with starvation bounds.
7. Tombstones and several high-volume histories need durable, bounded retention policies.
8. Security work remains incomplete: daemon authentication, identity-scoped policy, filesystem races, secret handling, and hostile protocol input.
9. Test totals and some historical summaries have drifted. Requirement-to-test mapping should replace cumulative counts.
10. Linux evidence dominates. macOS, Windows, ARM64, packaging, upgrade, and rollback need explicit matrices.

These gaps become roadmap work in Go. They do not justify maintaining two language implementations.

## 13. Reversal threshold

This is a definitive decision, not an irreversible claim about the universe. There is no scheduled language checkpoint and no standing plan to reconsider Rust after a particular phase.

A rewrite or replacement may be proposed only if all of the following are true:

1. Aren acquires a hard production requirement, expressed as a measurable threshold or mandatory integration contract.
2. The current Go runtime fails that requirement on a representative workload after profiling and proportionate design work.
3. The failure is attributable to a constraint that a different language can materially change, rather than to retention, admission, algorithms, protocol, storage, or unbounded work.
4. A competing implementation demonstrates the requirement under comparable architecture, workload, hardware, and operational conditions.
5. The migration cost and new operational complexity are lower than the cost of addressing the problem within Go.

Possible triggers include a proved no-GC requirement, mandatory in-process native embedding that a Go boundary cannot satisfy, or a hard laptop-capacity target that corrected process-group measurements show Go cannot meet. Speculative scale, language fashion, isolated microbenchmarks, or a favourable percentage without a product budget are not triggers.

## 14. Consequences

### Positive

- Aren can start one production codebase without a planned port.
- The lifecycle state machine can remain visible in direct control flow.
- The local library, CLI, HTTP transport, subprocess paths, and diagnostics can begin with the standard library.
- Race-enabled CI provides a strong dynamic check over executed Go paths.
- Multi-language consumers can remain thin clients of an earned daemon rather than forcing native bindings into the core.
- The team can spend complexity on lifecycle truth, failure evidence, containment, recovery, and product use instead of async framework integration.

### Costs accepted

- Go does not make published data immutable or lifecycle matches exhaustive by default.
- The garbage collector and goroutine stacks increase memory use and may add latency variance.
- Data races are possible and the race detector observes only executed paths.
- Goroutine leaks, channel deadlocks, and logical ordering errors remain possible.
- Some invalid combinations require constructors, private state, analyzers, and tests rather than compiler ownership rules.
- Native embedding into non-Go hosts is less direct than Rust's N-API path.

The mandatory engineering rules are the mitigation for these costs. If they are not adopted, the rationale for choosing Go no longer holds.

## 15. Current verification

The following checks passed on 26 August 2026 against the current working tree:

```text
Go:
  go vet ./...
  go test -race -count=1 ./...

Rust:
  cargo test --all-targets
  cargo clippy --all-targets -- -D warnings

Shared daemon contract:
  Go:   14 passed, 0 failed
  Rust: 14 passed, 0 failed
```

These checks verify the present test and contract state. They do not reproduce the historical performance measurements or close the production gaps listed above. Published cumulative test totals are intentionally omitted because the experiment reports use inconsistent counting methods.

## 16. Decision outcome

Go is the implementation language for Aren's production runtime.

The choice is based on Aren's requirements, not on a claim that Go dominates Rust in general. Rust demonstrated better static ownership and lower memory in the prototype. Go demonstrated sufficient capability across every tested roadmap slice with a smaller, more direct implementation and lower continuing machinery.

That is the better match for Aren's governing rule: do not take on complexity until evidence shows it is required.

## 17. Evidence sources

Primary requirements and planning:

- [Project lineage and development doctrine](../project-lineage.md)
- [Phased roadmap](../phased-roadmap.md)
- [Phase 1 PRD](../phase-1-prd.md)
- [Phase 1 product definition](../phase-1-prd/01-product-definition-and-scope.md)
- [Phase 1 lifecycle requirements](../phase-1-prd/02-lifecycle-requirements.md)
- [Phase 1 validation and acceptance](../phase-1-prd/03-validation-and-acceptance.md)
- [Phase 1 risks, decisions, and handoff](../phase-1-prd/04-risks-decisions-and-handoff.md)
- [Performance engineering](../performance-engineering.md)
- [Coding-agent context engineering](../coding-agent-context-engineering.md)

Decision analysis:

- [Consolidated Rust versus Go analysis](consolidated-rust-vs-go-analysis.md)
- [Long-term Go recommendation](long-term-go-recommendation.md)
- [Critique of the Rust argument](critique-of-the-rust-argument.md)
- [Language-decision evidence index](README.md)
- [Prototype overview](../../prototypes/README.md)

Experiment reports:

- [Initial prototype report](report.html)
- [Experiment 2: cancellation and retention](phase-1-experiment-2-cancellation-and-retention.md)
- [Experiment 3: observation, invariants, and stress](phase-1-experiment-3-observation-invariants-and-stress.md)
- [Experiment 4: progress, cleanup, and conformance](phase-2-experiment-4-progress-cleanup-and-conformance.md)
- [Experiment 5: first model invocation](phase-3-experiment-5-first-model-invocation.md)
- [Experiment 6: streaming](phase-4-experiment-6-streaming-model-execution.md)
- [Experiment 7: structured results](phase-5-experiment-7-structured-results.md)
- [Experiment 8: retry and backoff](phase-6-experiment-8-retry-and-backoff.md)
- [Experiment 9: tool-call representation](phase-7-experiment-9-tool-call-representation.md)
- [Experiment 10: local tool execution](phase-8-experiment-10-local-tool-execution.md)
- [Experiment 11: bounded agent loop](phase-9-experiment-11-bounded-agent-loop.md)
- [Experiment 12: context management](phase-10-experiment-12-context-management.md)
- [Experiment 13: persistence and recovery](phase-11-experiment-13-persistence-and-recovery.md)
- [Experiment 14: execution composition](phase-12-experiment-14-execution-composition.md)
- [Experiment 15: policies and resources](phase-13-experiment-15-policies-and-resources.md)
- [Experiment 16: daemon and clients](phase-14-experiment-16-daemon-and-clients.md)
- [Shared-daemon gauntlet](phase-16-experiment-17-shared-daemon-gauntlet.md)
- [Workflow experiment](phase-16-experiment-17-workflows.md)
- [Experiment 18: broader execution types](phase-17-experiment-18-broader-execution-types.md)
- [Experiment 19: distributed and remote execution](phase-18-experiment-19-distributed-and-remote.md)
