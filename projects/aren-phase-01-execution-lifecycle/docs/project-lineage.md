# Aren Project Lineage and Engineering History

## Purpose

This document records the development history that led to Aren.

Its purpose is not to catalogue every repository or preserve every previous architectural decision. It exists to explain how Antonio’s approach to building software and agent systems evolved, what repeatedly went wrong, what eventually began to work, and which lessons should guide Aren.

Aren is not an isolated new project. It emerges from several years of building applications, frameworks, autonomous systems, developer tools and planning workflows.

The most important history is therefore not a list of features. It is the progression in development philosophy:

> Ambitious products and speculative architecture  
> → reusable frameworks and generalised systems  
> → autonomous agent experiments  
> → smaller local tools solving immediate problems  
> → evidence-grounded planning and execution  
> → controlled agent-runtime supervision  
> → Aren

This history should inform Aren without constraining it. Previous projects provide evidence, not unquestionable precedent.

---

# 1. Executive Summary

Antonio’s earlier projects often began with a compelling product or technical vision and quickly expanded into large systems.

The recurring pattern was:

1. Identify a valuable idea.
2. Imagine the mature version of the product.
3. Design abstractions for many future requirements.
4. Expand the scope before validating the smallest useful capability.
5. Spend increasing effort on architecture, infrastructure and planning.
6. Lose momentum before the original value proposition was fully proved.

This occurred across product projects, agent frameworks and infrastructure experiments.

The failure was not a lack of technical ambition or effort. The central problem was that architectural and product complexity arrived before enough practical evidence existed to justify it.

Later projects began to reverse this pattern.

Instead of starting with platforms, Antonio began creating small local tools that solved concrete problems in his own workflows. These tools operated directly on files, processes and repositories. They did not require servers, databases, dashboards or large orchestration systems.

This shift produced several important projects:

- 24-Hour Testers explored continuous autonomous validation.
- Ultra and UltraPlan formalised evidence-grounded research, reasoning and planning.
- AgentWrap provided controlled supervision of existing coding-agent runtimes.
- Aren now proposes to move from wrapping another harness to owning the essential harness foundations directly.

Aren should inherit the discipline of the later projects without simply merging their feature sets.

The defining development principle should be:

> Build the smallest coherent foundation, validate it under real use and failure, and allow every additional abstraction to earn its existence.

---

# 2. Development Phases

## Phase 1: Product Ambition and Expanding Scope

### Representative projects

- Elevate
- CoachUp
- Eloquence
- Soft Skills
- related frontend, API and landing-page repositories

### Initial direction

These projects generally began as user-facing products with meaningful goals.

They explored areas such as:

- AI-assisted learning
- communication and professional-skills coaching
- structured feedback
- progression systems
- reusable scenario banks
- speech or interview practice
- dashboards and user experiences
- multi-service application architecture

The ideas were legitimate and often well considered. The difficulty was not that the projects lacked value. It was that the intended product frequently became too complete too early.

A simple useful experience would expand into:

- user accounts
- progression models
- scoring systems
- immutable attempt history
- content management
- speech processing
- dashboards
- administration
- prompt versioning
- observability
- governance
- event systems
- multiple services
- extensibility for future modes

The architecture began solving the requirements of a mature platform before the core user interaction had been validated sufficiently.

### Recurring failure mode

The product vision and the implementation plan became tightly coupled.

Because the eventual product might need a capability, that capability started to feel like a present requirement.

This produced several forms of scope creep:

#### Product scope creep

More experiences, modes, user journeys and supporting features were added before the central loop had been proved.

#### Technical scope creep

Infrastructure was designed for future scale, extensibility or provider independence before the system had enough real use to reveal where those boundaries were actually needed.

#### Planning scope creep

The effort to fully understand and perfectly structure the future system became a major project in itself.

### Lessons

- A convincing full-product vision is not the same as a validated first product.
- A future requirement should not automatically become a current abstraction.
- Architecture cannot compensate for an unproven core interaction.
- The more complete the plan becomes, the easier it is to mistake planning progress for product progress.
- Product development needs deliberate limits, not merely a prioritised backlog.
- A system can be technically sophisticated while still failing to deliver a small finished outcome.

### Aren implication

Aren must not begin as the complete universal agent platform described by its eventual possibilities.

The first version must be independently useful and technically complete at a very small scope.

---

## Phase 2: Framework Building and Generalisation

### Representative projects

- Stageflow
- Stable Flow
- Unified Content Protocol
- Stageflow Rust
- Voice Engine
- Graphcode and related graph experiments

### Initial direction

After building application-specific systems, attention increasingly shifted toward reusable foundations.

Instead of solving only one application problem, these projects attempted to define general primitives:

- stages and pipelines
- workflow execution
- DAGs
- interceptors
- events
- protocols
- structured content representations
- graph-based models
- portable abstractions
- runtime observability
- language bindings
- extension points

This work produced several strong engineering ideas.

Stageflow, for example, focused on explicit execution stages, observable runs and cross-cutting runtime behaviour.

UCP explored whether agents could interact with structured content through a consistent intermediate representation rather than treating every file format separately.

These projects also strengthened Antonio’s interest in:

- explicit contracts
- traceability
- deterministic behaviour
- event-driven observation
- reusable execution primitives
- testing architectural boundaries

### Recurring failure mode

The abstraction itself could become the product.

Once a reusable primitive was identified, the project naturally expanded toward a general framework capable of supporting many future use cases.

This introduced risks:

- designing extension points before extension pressure existed
- supporting multiple execution models too early
- separating interfaces from implementations without genuine volatility
- creating several layers around a small amount of core behaviour
- building infrastructure for hypothetical adopters rather than current workflows
- treating conceptual elegance as evidence of practical necessity

The systems were often thoughtful, but the abstraction surface grew faster than the body of real usage validating it.

### Lessons

- Reusability is discovered through repeated use, not established through naming.
- A clean interface can still represent premature abstraction.
- Multiple implementations are stronger evidence for a boundary than the possibility of multiple implementations.
- Cross-cutting concerns such as observability and cancellation are foundational only when attached to concrete execution semantics.
- Framework design should follow repeated pressure from applications, not precede it.
- An abstraction should remove present complexity, not relocate speculative complexity.

### Aren implication

Aren may eventually expose rich interfaces, adapters, middleware, execution policies and extension mechanisms.

They should be added only when a concrete implementation demonstrates why the boundary is necessary.

Aren should prefer:

- one working implementation
- direct types
- explicit control flow
- minimal layering
- refactoring after pressure appears

over:

- early plugin systems
- generic factories
- broad provider interfaces
- speculative compatibility layers
- framework-wide extension APIs

---

## Phase 3: Autonomous Agent Systems

### Representative projects

- Hivemind
- 24-Hour Testers
- 24-Hour Hivemind
- 24-Hour Hivemind Native
- 24-Hour UCP
- 24-Hour Codegraph
- 24-Hour Benchmarking
- Stageflow Production Testers

### Initial direction

This phase explored systems that could continue working with limited supervision.

The focus moved toward:

- autonomous loops
- parallel agent execution
- agent coordination
- retries
- verification
- continuous smoke testing
- generated reports
- repository inspection
- long-running work
- unattended execution
- feedback loops

The strongest idea from this phase was that an agent system should not merely generate an answer and stop.

It should:

1. perform work,
2. inspect the result,
3. gather evidence,
4. identify weaknesses,
5. retry or continue,
6. leave durable artefacts behind.

The 24-Hour Testers concept captured this particularly clearly. It treated testing as a persistent loop that repeatedly identifies areas to inspect, executes checks and records findings.

### What worked

#### Work must leave evidence

Agent activity should produce inspectable files, logs, events, reports or patches.

#### Completion must be challenged

A successful process exit or plausible answer is insufficient. Results need validation.

#### Long-running systems require lifecycle control

Cancellation, retries, rate limits, stuck processes and partial failures become core concerns as soon as agents operate unattended.

#### Observability is part of correctness

A system that cannot explain what it did, what it attempted and why it stopped is difficult to trust.

#### Verification should be designed into the loop

Review and smoke testing should not be optional activities added at the end.

### Recurring failure mode

Autonomy dramatically multiplies scope.

A supposedly simple autonomous loop quickly raises questions about:

- scheduling
- persistence
- coordination
- shared state
- retries
- budgets
- permissions
- concurrency
- agent selection
- recovery
- task decomposition
- progress tracking
- observability
- user intervention
- durable execution

This can lead to building an orchestration platform before proving the value of one narrow loop.

Multi-agent systems added further complexity. Coordination often became more difficult than the original task.

### Lessons

- Autonomy should begin with a bounded loop, not an open-ended agent organisation.
- Verification is more valuable than additional agent roles.
- One agent with a strong execution and review loop can outperform a poorly coordinated group.
- Every autonomous action needs explicit termination conditions.
- Durable artefacts are often more useful than elaborate internal state.
- A local process can prove the workflow before a daemon, queue or distributed scheduler is justified.
- Agent count is not a measure of system capability.

### Aren implication

Aren should not begin as a multi-agent orchestration platform.

Its earliest loop should be:

- single-agent
- explicit
- bounded
- observable
- cancellable
- testable
- deterministic where possible

Additional agents, graphs, routing and durable orchestration should be introduced only after the single-agent loop is thoroughly understood.

---

## Phase 4: The Shift to Small Local Tools

### Representative projects

- Go CLI Study
- OpenCode Wrap Study
- Ultraground
- Ultra
- early UltraPlan
- small repository-analysis and workflow utilities

### Change in approach

This phase represents the most important change in development behaviour.

Rather than beginning with a large application or framework, Antonio began building small command-line tools that solved immediate problems in his own work.

These tools tended to have several characteristics:

- local-first
- filesystem-oriented
- CLI-based
- no server
- no database
- limited operational dependencies
- composable through files and processes
- useful before becoming general
- easy to inspect and modify
- designed around a real personal workflow

This reduced both technical and cognitive overhead.

The development question changed from:

> What should the complete system eventually support?

to:

> What is the smallest tool that solves the problem I have today?

### Why this worked better

#### The user was known

Antonio was building primarily for his own workflow. This removed the need to speculate about broad personas and hypothetical usage patterns.

#### Value appeared earlier

A command could become useful before the surrounding platform existed.

#### Files acted as a simple integration boundary

Inputs, outputs, prompts, reports and state could remain editable and visible.

#### Operational complexity stayed low

There was no immediate requirement for hosting, authentication, database migrations, queues or deployment infrastructure.

#### Architecture could follow usage

Patterns became visible through repeated real use rather than through imagined scenarios.

### Limits of the local-tool approach

The simplicity of local tools should not become ideology.

Local files and subprocesses have real limitations:

- weaker concurrency control
- limited durable scheduling
- difficult cross-process coordination
- platform-specific process behaviour
- imperfect live control
- limited remote access
- no automatic multi-language embedding
- potential filesystem consistency problems

The lesson is not that servers, databases or daemons are bad.

The lesson is that they should appear when a demonstrated requirement makes their cost worthwhile.

### Aren implication

Aren should initially preserve the strengths of the local-tool era:

- local execution
- direct process ownership
- simple installation
- explicit files and logs
- minimal infrastructure
- inspectable state
- excellent CLI ergonomics

A daemon should be introduced only when persistent ownership, cross-process access or multi-language clients require it.

Even then, the daemon should extend a proven core rather than define it from the beginning.

---

## Phase 5: Evidence-Grounded Planning

### Representative projects

- Ultra
- Ultraground
- UltraPlan
- UltraPlan Go
- `.ultra`
- UltraPlan Workspace
- AI Agent Examples

### Initial problem

Planning large systems from first principles repeatedly produced speculative design.

Research was often informal:

- remembered from previous reading
- based on framework documentation
- influenced by attractive architecture patterns
- detached from actual implementations
- difficult to trace later

UltraPlan attempted to turn research and planning into an explicit workflow grounded in repositories and durable artefacts.

Its broad process became:

> study → select → distil → reason → plan → execute → smoke → review

The UltraPlan workspace separates studies, planning projects, runtime state and generated outputs. Its CLI includes validation, status and staged workflow operations, making the planning process itself inspectable rather than informal.

### What worked

#### Research became reproducible

Findings could be tied to source files and preserved as reports.

#### Planning became staged

Study, reasoning, planning and execution became distinct activities rather than one large design exercise.

#### Evidence became editable

Markdown and filesystem artefacts allowed human review and correction.

#### Validation applied to documents

Requirements and plans could be checked for completeness before execution.

#### Research could run unattended

Agents could inspect several repositories and produce structured findings without constant interaction.

### New risk: process overgrowth

UltraPlan solved the problem of ungrounded planning, but it introduced another possible failure mode: planning infrastructure can itself become excessively elaborate.

A detailed process may create:

- too many artefact types
- repeated information across documents
- long preparation before implementation
- rigid stages for trivial changes
- excessive validation of low-risk work
- large studies that are never converted into decisions
- a false sense of certainty from structured documents

This is especially important for Aren.

Aren should use UltraPlan, but Aren must not become an excuse to exercise every UltraPlan capability.

### Lessons

- Research should answer a decision, not merely accumulate knowledge.
- Evidence should be proportional to the importance and uncertainty of the decision.
- Planning stages should reduce risk, not delay contact with implementation.
- A plan should become smaller and more concrete as it approaches execution.
- Documents should have distinct responsibilities.
- Validation is valuable, but validation rules can also become bureaucracy.
- A short implementation spike can sometimes answer a question better than a long comparative study.

### Aren implication

The Aren planning process should be rigorous but deliberately bounded.

Every study should state:

- the decision it supports
- why existing evidence is insufficient
- what repositories or sources are relevant
- what output will change based on the result
- when the study is complete

Aren should not attempt to fully design later phases during the first phase.

---

## Phase 6: Runtime Supervision with AgentWrap

### Representative projects

- AgentWrap
- AgentWrap Smoke
- OpenCode Wrap Study

### Initial problem

Existing coding-agent harnesses such as OpenCode were highly capable, but invoking them from repeatable workflows introduced reliability problems.

A parent workflow needed more than a shell command.

It needed to understand:

- whether the runtime was available
- whether configuration was valid
- how execution started and stopped
- how cancellation worked
- what events were emitted
- whether a failure was retryable
- whether a rate limit occurred
- what files or artefacts were produced
- whether the output satisfied expectations
- how permissions were applied
- what metadata should be retained

AgentWrap was created as a Go SDK for supervising coding-agent runtimes from product workflows.

Its implemented surface includes runtime-neutral types, classified errors, resilience policies, output validation, repair attempts, permission policies, canonical events, run records, health checks and an OpenCode adapter.

### Important design decisions

#### Runtime policies remained outside adapters

Retries, fallback and backoff were implemented as wrappers around runtimes rather than being embedded directly inside each adapter.

This preserved a distinction between:

- how a runtime is invoked
- how a product responds to failure

#### Process success was separated from output success

A run was considered successful only when execution completed and configured validators passed.

This is a critical agent-system principle:

> The model completing its turn does not prove that the requested work was completed correctly.

#### Repair was bounded

Validation repair inherited the original session, working directory, provider, model and permission posture, but used explicit attempt limits.

#### Permissions were explicit and auditable

Permission policies were attached at run initialisation and translated into runtime-native configuration where supported. Unsupported required behaviour failed before process start rather than being silently ignored.

#### Observability did not require adapter rewrites

Observation and persistence were added through a runtime wrapper, preserving canonical events and completed run records.

#### Unknown remained unknown

Unavailable token or usage values were preserved as unknown rather than converted to zero.

This reflects a broader principle Aren should inherit:

> Missing information must not be represented as successful measurement.

### What AgentWrap proved

AgentWrap demonstrated that Antonio could build a bounded, robust component incrementally.

It was narrower than previous platforms, had clear scope guardrails and solved a real integration problem.

It also showed the value of wrapping an existing harness:

- rapid access to mature coding-agent behaviour
- little need to implement model loops
- immediate practical utility
- the ability to focus on lifecycle and reliability

### What AgentWrap could not provide

Wrapping another harness limits control.

The wrapper depends on the child runtime’s:

- event model
- process behaviour
- session semantics
- permission system
- output format
- transport
- error reporting
- internal loop
- tool model
- context handling

A subprocess boundary is useful, but it prevents the parent from fully controlling or understanding core agent execution.

AgentWrap’s scope guardrails explicitly deferred concerns such as live approval transport, durable backend selection and global throttling. It also deliberately excluded UltraPlan workflow logic.

### Aren implication

Aren is not AgentWrap expanded.

AgentWrap supervised a harness.

Aren will own the essential harness behaviour.

However, Aren should preserve AgentWrap’s strongest lessons:

- explicit lifecycle semantics
- typed errors
- cancellation as a core capability
- bounded retries
- honest metadata
- validation separate from execution
- permissions established before work
- canonical observable events
- strict scope guardrails

---

## Phase 7: Aren

### Why Aren now

Aren becomes viable because the preceding projects answered different parts of the problem.

The early product projects demonstrated the danger of uncontrolled ambition.

The framework projects developed an appreciation for explicit execution and observability.

The autonomous projects exposed the realities of retries, verification and unattended work.

The local CLI tools demonstrated a simpler and more productive way to build.

UltraPlan introduced evidence-grounded planning and staged implementation.

AgentWrap proved the value of strong runtime contracts and lifecycle supervision.

The remaining limitation is control.

Using an existing harness enabled rapid progress, but its internal decisions remain outside Antonio’s ownership.

Aren is the attempt to build the smallest agent harness whose execution semantics, lifecycle, context handling, tool behaviour and extension model are fully understood and intentionally designed.

### What Aren must not become

Aren must not become:

- a merger of every previous project
- a universal agent platform from the first release
- a multi-agent orchestration framework before a single-agent loop works
- a daemon before persistent service ownership is required
- a provider abstraction before the first provider is fully implemented
- a workflow graph engine before linear execution is insufficient
- a plugin ecosystem before internal boundaries stabilise
- an observability platform instead of an observable harness
- an excuse to perform endless architecture research
- a final architecture designed entirely in advance

### What Aren should become

Aren should become a deeply understood harness built layer by layer.

Each phase should:

1. introduce one coherent capability,
2. define its behavioural contract,
3. implement the smallest useful form,
4. test normal operation,
5. test failure and cancellation,
6. exercise it in real usage,
7. review what the implementation revealed,
8. simplify or refactor,
9. only then begin the next layer.

The goal is not to avoid architecture.

The goal is to ensure architecture is derived from evidence.

---

# 3. Cross-Project Lessons

## 3.1 Scope must be actively constrained

Scope does not remain small by itself.

Every useful capability creates adjacent possibilities. Without explicit exclusions, a project gradually expands toward a platform.

Each Aren phase should therefore define:

- in scope
- out of scope
- deferred
- evidence required to reconsider deferred work

Deferred work should not appear as partially implemented infrastructure.

---

## 3.2 Abstractions must be earned

An abstraction is earned when it resolves repeated, demonstrated pressure.

Strong evidence includes:

- two implementations with meaningful differences
- repeated conditional logic
- duplicated lifecycle handling
- a boundary required for testing failure
- a stable concept appearing across several use cases
- a concrete need to replace or extend behaviour

Weak evidence includes:

- it may be useful later
- other frameworks use it
- it makes the architecture look clean
- it could support plugins
- it avoids hypothetical coupling

---

## 3.3 Build vertical slices through foundations

“Foundations first” must not mean building a large invisible substrate before anything works.

A foundation should be validated through a thin end-to-end slice.

For example, the first LLM layer should not consist only of interfaces and request types. It should make one real request, stream events, support cancellation, expose errors and be tested.

The slice can be small, but it must be complete enough to reveal reality.

---

## 3.4 Observability is part of the contract

Agent systems involve probabilistic behaviour, external providers, tools, processes and long-running execution.

Logs added afterward are insufficient.

Aren should expose structured lifecycle events from the start.

However, event design must remain proportional. The goal is to explain execution, not to create a second system mirroring every internal operation.

---

## 3.5 Cancellation is foundational

Cancellation affects:

- provider requests
- streams
- tool execution
- subprocesses
- retries
- backoff waits
- event delivery
- cleanup
- final state

It cannot be reliably added at the end.

Every phase that introduces blocking work must define cancellation behaviour.

---

## 3.6 A successful call is not a successful task

A provider can return successfully while:

- producing malformed output
- failing to call a required tool
- editing the wrong file
- ignoring a constraint
- stopping too early
- returning incomplete structured data

Execution result, model result and task result are distinct concepts.

Aren should preserve this distinction throughout its design.

---

## 3.7 Failure needs semantics

String errors are insufficient for reliable agent workflows.

The caller may need to distinguish:

- cancellation
- timeout
- authentication
- rate limiting
- unavailable provider
- invalid request
- malformed response
- tool failure
- permission denial
- validation failure
- exhausted retries
- internal invariant violation

The taxonomy should start small and expand only when callers need different behaviour.

---

## 3.8 Unknown is a valid state

Agent systems frequently have partial information.

Examples include:

- missing token usage
- estimated cost
- uncertain retry timing
- provider-specific stop reasons
- unknown tool completion
- incomplete session continuation
- ambiguous network failure

Aren should not replace unknown values with misleading defaults.

---

## 3.9 Local-first is a starting advantage, not a permanent restriction

Local operation reduces complexity and accelerates learning.

Aren should begin as an in-process library and CLI. Its implementation language
remains an open decision. Once selected, the language is intended for the
long-term runtime rather than an initial implementation followed by a planned
port.

A daemon becomes justified when the system needs capabilities such as:

- persistent ownership of long-running execution
- several clients
- cross-language access
- centralised lifecycle control
- reconnectable streams
- execution surviving client termination
- shared scheduling or resource management

The daemon should be added around a proven core.

---

## 3.10 Planning must terminate in implementation

Research and planning are valuable only when they improve decisions.

Every Aren planning artefact should eventually connect to:

- a decision
- a phase boundary
- a contract
- an implementation task
- a test
- a rejected alternative

Documents that no longer affect action should be archived rather than continually expanded.

---

# 4. Aren Inheritance Map

## Inherit

These ideas have strong evidence from previous projects and should be treated as likely Aren principles:

- local-first development
- filesystem-visible artefacts
- explicit execution lifecycle
- structured events
- cancellation throughout the stack
- classified failures
- bounded retry behaviour
- validation separate from execution
- explicit permission posture
- honest unknown values
- evidence-backed design decisions
- smoke testing after implementation
- phased development
- explicit scope guardrails
- abstraction only after demonstrated pressure

## Re-evaluate

These ideas may be useful but must be proved in Aren:

- runtime-neutral provider interfaces
- middleware or interceptor systems
- DAG execution
- event buses
- run stores
- plugin systems
- graph-based context
- durable execution
- model fallback
- automatic output repair
- daemon architecture
- WebSocket or SSE client transport
- multi-language SDKs
- generic workflow orchestration
- multiple concurrent agents

## Reject by default

These patterns should not be introduced without unusually strong evidence:

- designing the complete final architecture before phase one
- abstractions with only one trivial implementation
- interfaces mirroring concrete types without behavioural value
- a server purely because the future product may need one
- a database for state that can remain files or memory
- multi-agent coordination before a single-agent loop is excellent
- plugin APIs before internal contracts stabilise
- broad configuration surfaces for hypothetical users
- duplicated planning documents
- infrastructure that does not support a current vertical slice
- features justified primarily by competitors or frameworks having them
- calling a component robust without testing its failure behaviour

---

# 5. Development Doctrine for Aren

The following doctrine should govern roadmap and phase design.

## 5.1 One phase, one primary uncertainty

Each phase should answer one major technical question.

Examples:

- Can Aren perform a cancellable streaming model call reliably?
- Can Aren validate structured output without corrupting stream semantics?
- Can Aren execute a tool with clear lifecycle and failure behaviour?
- Can Aren run a bounded agent loop?
- Can Aren manage context without hidden mutation?

A phase may contain several implementation tasks, but they should converge on one primary uncertainty.

---

## 5.2 Every phase must produce a usable system

The output of each phase should run.

A phase consisting only of types, interfaces or internal architecture is incomplete unless those are exercised through a real path.

---

## 5.3 Every phase must include failure testing

Normal-path tests are insufficient.

Each phase should identify and test relevant failures such as:

- cancellation
- provider timeout
- malformed events
- interrupted streams
- invalid structured output
- tool crashes
- permission denial
- retry exhaustion
- cleanup failure
- partial state

---

## 5.4 Real use comes before broadening

After a phase is technically complete, Aren should be used in a small real workflow.

The next phase should not begin until this use has revealed whether:

- the API is awkward
- events are missing
- failure semantics are unclear
- unnecessary abstractions exist
- important state is hidden
- the feature solves the intended problem

---

## 5.5 Reviews may remove capabilities

A phase review is not only a gate for adding the next feature.

It may conclude that:

- an abstraction should be deleted
- an interface should become concrete
- a configuration option should be removed
- an event should be merged
- a layer should collapse
- the next planned feature is not yet justified

Reduction is valid progress.

---

## 5.6 Future extensibility should come from clear internals

Aren should not attempt to guarantee unlimited extensibility from the first release.

The better route is:

1. keep the core small,
2. make ownership explicit,
3. avoid hidden global state,
4. separate behaviour only where necessary,
5. test contracts,
6. refactor when a second use case arrives.

A system with clear internals can become extensible later.

A system full of speculative extension points can become difficult to change immediately.

---

# 6. Repository Lineage

## Primary lineage repositories

These repositories deserve direct study because they materially shaped Aren:

| Project | Contribution to Aren |
|---|---|
| Elevate | Early product ambition and scope-creep lessons |
| CoachUp | Rich product design before a narrow validated core |
| Eloquence | Continued exploration of AI coaching and product complexity |
| Stageflow | Explicit execution stages, workflow semantics and observability |
| Unified Content Protocol | Generalisation, structured content and abstraction lessons |
| Hivemind | Multi-agent orchestration and coordination complexity |
| 24-Hour Testers | Autonomous validation loops and durable evidence |
| Ultra / Ultraground | Filesystem-grounded research and reasoning |
| UltraPlan | Staged study, planning, execution, review and smoke testing |
| UltraPlan Go | Simplification into a local-first Go CLI |
| AgentWrap | Runtime supervision, lifecycle control and reliability contracts |
| AgentWrap Smoke | Real-runtime validation and integration evidence |
| Aren | Ownership of the agent harness core |

## Supporting repositories

Other repositories provide useful supporting evidence in areas such as:

- graph modelling
- benchmarking
- API architecture
- frontend product development
- learning tools
- testing
- Go CLI design
- repository studies
- workflow experiments

They should remain discoverable in the wider catalogue but should not receive equal weight in the main lineage.

---

# 7. Questions to Carry into Aren Planning

This history does not determine Aren’s design. It establishes questions that the roadmap and PRD must answer.

1. What is the smallest independently useful version of Aren?
2. What concrete workflow will validate the first phase?
3. Which AgentWrap behaviours belong in the Aren core, and which were specific to wrapping subprocess runtimes?
4. Should the first public surface be a library API, CLI, or both?
5. What execution state must be externally observable from the first phase?
6. What does cancellation mean at each layer?
7. What is the minimum useful error taxonomy?
8. How will Aren distinguish provider completion from task success?
9. Which future capabilities are explicitly excluded from phase one?
10. What evidence would justify introducing a daemon?
11. What evidence would justify introducing a second provider?
12. What evidence would justify a general tool interface?
13. What evidence would justify durable state?
14. How will phase reviews detect new scope creep?
15. How will UltraPlan support the project without becoming a source of over-planning?

---

# 8. Final Perspective

Aren is described as Antonio’s final agent harness.

That ambition is understandable, but it contains a danger.

Trying to ensure Aren is the final harness could recreate the exact behaviour that caused earlier projects to expand: designing now for every future need so that the system never has to be replaced.

A more useful interpretation is:

> Aren should be the harness that is developed with enough discipline that it can continue evolving without repeatedly being abandoned and restarted.

That does not require getting the final architecture right at inception.

It requires:

- foundations that are genuinely understood
- small validated steps
- honest treatment of uncertainty
- willingness to remove mistaken abstractions
- reliable tests
- clear phase boundaries
- continuous real use
- resistance to speculative scope

The goal is not perfection at the beginning.

The goal is a development process capable of producing durable quality over time.

Aren’s most important inheritance is therefore not a runtime contract, event system, graph model or workflow primitive.

It is the lesson that the route to a robust and extensible system is not to build everything it might eventually need.

It is to build the next essential layer well enough that the correct following layer becomes visible.
