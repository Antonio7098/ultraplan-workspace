# LLM Runtime Contract

## Purpose
This contract defines the architectural and operational rules for LLM-backed runtime capabilities.

It governs:
- shared LLM substrate boundaries
- the distinction between conversational and structured execution modes
- tool and structured-output rules
- lifecycle inspectability
- operational exposure of LLM-backed behavior
- monitoring through the project's observability surfaces

## Scope
This contract applies when a change:
- adds or edits conversational agent runtime behavior
- adds or edits structured task runtime behavior
- adds or edits shared LLM provider or substrate behavior
- adds concrete agent, task, prompt, or runtime definitions
- changes tool execution, structured-output generation, validation, retry, timeout, or backup-provider behavior
- exposes LLM-backed capabilities through application modules or transport adapters

## Runtime Modes

Many LLM-backed systems have one or both of these runtime modes:

- **Conversational mode:** turn-based agent execution with optional tools, approvals, resume semantics, and response text generation
- **Structured mode:** task execution with explicit input/output schemas, local validation, bounded retries/timeouts, and optional fallback selection

Shared provider and runtime substrate concerns should stay generic.
Mode-specific policy should branch only where the semantics materially differ.

## How To Use This Contract

Apply shared rules first.
If a rule includes mode-specific expectations, only apply those branches when the change touches that mode.

## Requirement Index

| ID | Title | Applies To | Severity If Violated |
| --- | --- | --- | --- |
| LLM-BOUNDARY-001 | Shared substrate, runtime mechanics, and mode-specific policy must stay separated | all LLM-backed runtime code | High |
| LLM-TOOL-001 | Tool execution boundaries must stay explicit and mode-scoped | tool-capable LLM runtimes | High |
| LLM-IO-001 | Structured input and output boundaries must stay explicit when used | structured-output LLM runtimes | High |
| LLM-LIFECYCLE-001 | Lifecycle controls must stay explicit and inspectable | all LLM-backed runtime code with lifecycle controls | High |
| LLM-RETRY-001 | LLM retry policy must be shared, bounded, mode-aware, and inspectable | provider-backed LLM execution | High |
| LLM-RUN-001 | Runs and events must be durable or inspectable | long-lived, batch, or owned LLM execution | High |
| LLM-PROMPT-001 | Prompts must be explicitly versioned and version propagation must stay observable | all LLM-backed runtime code with concrete prompts | High |
| LLM-EXPOSE-001 | Operational exposure must flow through application modules | CLI commands, routes, operational APIs | High |
| LLM-OBS-001 | LLM runtime monitoring must reuse canonical observability surfaces | metrics, tracing, diagnostics, logs, and events | Medium |
| LLM-EVAL-001 | Prompt, model, and runtime changes must have evaluation or smoke coverage appropriate to risk | LLM-backed runtime changes | High |
| LLM-SAFETY-001 | Tool use, user instructions, and model outputs must pass explicit safety and policy checks where relevant | tool execution, user input, and model output paths | High |
| LLM-COST-001 | LLM paths must expose token, cost, and latency metadata where operationally meaningful | LLM-backed runtime paths | Medium |

## Requirements

### LLM-BOUNDARY-001: Shared Substrate, Runtime Mechanics, And Mode-Specific Policy Must Stay Separated

**Rule**
Shared provider/substrate behavior, generic runtime mechanics, and mode-specific policy must not collapse into one layer.

**Required**
- keep shared LLM contracts, response normalization, and provider adapters in generic runtime-facing packages
- keep concrete prompts, tool selection policy, semantic validation, and product-facing behavior out of shared substrate code
- keep mode distinctions explicit when conversational and structured execution have materially different semantics

**Conversational mode specifics**
- keep turn/tool/approval/response-text semantics in conversational runtime and application policy layers

**Structured mode specifics**
- keep schema validation, retry correction, timeout policy, and backup-provider selection in structured runtime and application policy layers

**Forbidden**
- embedding concrete agent or structured task policy in shared LLM substrate code
- pushing structured task behavior into conversational runtime code to avoid adding a real boundary
- duplicating shared provider mechanics into separate mode-specific stacks without a justified reason

**Evidence**
- shared substrate packages expose generic mechanics and contracts
- conversational and structured definitions live outside shared runtime layers

### LLM-TOOL-001: Tool Execution Boundaries Must Stay Explicit And Mode-Scoped

**Rule**
Tool execution must stay explicit, inspectable, and confined to the modes that intentionally support it.

**Required**
- execute tools through explicit contracts and allowed tool bundles
- keep tool selection policy explicit in the owning mode
- preserve inspectable tool lifecycle state when tool execution is part of the supported runtime model

**Conversational mode specifics**
- tool calls, approvals, and results must remain part of the persisted conversational execution record

**Structured mode specifics**
- structured tasks must not smuggle ad hoc tool execution through prompts or hidden runtime side effects; if a task needs tools later, that capability must be designed explicitly

**Forbidden**
- hidden tool execution inside transport adapters, shared substrate code, or unrelated business code
- hardwiring mode-specific tool behavior into generic runtime classes

**Evidence**
- allowed tool bundles and execution contracts are represented explicitly in code
- unsupported modes do not quietly acquire tool behavior

### LLM-IO-001: Structured Input And Output Boundaries Must Stay Explicit When Used

**Rule**
Any LLM-backed path that claims structured input or structured output must make those contracts explicit and keep local validation authoritative.

**Required**
- validate structured inputs against explicit schemas or models before execution when applicable
- request structured output through the shared substrate where supported
- parse and validate structured output locally before marking execution successful
- treat provider strictness as transport guidance rather than correctness authority

**Structured mode specifics**
- structured task definitions should declare input and output models
- validation failures should be inspectable and participate in explicit retry/terminal behavior

**Forbidden**
- treating prompt wording alone as the only structured-output contract
- trusting provider-side schema conformance without local validation
- returning partially validated structured output as successful state

**Evidence**
- structured definitions declare their contracts explicitly
- runtime behavior records validation failures visibly

### LLM-LIFECYCLE-001: Lifecycle Controls Must Stay Explicit And Inspectable

**Rule**
If an LLM-backed runtime supports approval, resume, cancellation, retry, timeout, or fallback behavior, those controls must be explicit and inspectable from the start.

**Required**
- expose lifecycle transitions through explicit state changes and inspectable traces
- preserve configuration and outcome detail for lifecycle-affecting controls
- keep terminal states unambiguous

**Conversational mode specifics**
- approvals, rejection paths, cancellation, and resume behavior must remain explicit

**Structured mode specifics**
- retry budgets, timeout behavior, fallback selection, and cancellation must remain explicit

**Forbidden**
- hidden pause/resume behavior with no durable or inspectable state
- hidden retries or fallback behavior with no operator-visible trace
- ambiguous terminal outcomes after cancellation or timeout

**Evidence**
- lifecycle actions and transitions are represented explicitly in code and runtime state

### LLM-RETRY-001: LLM Retry Policy Must Be Shared, Bounded, Mode-Aware, And Inspectable

**Rule**
Provider-backed LLM retries must use the shared LLM execution policy, must retry retryable issues only within an explicit attempt budget, and must preserve mode-specific safety rules for conversational agents and structured tasks.

**Shared retry policy**
- normalize retry-relevant issues before deciding whether to retry
- classify issue kind with stable values such as provider error, timeout, empty completion, invalid JSON, and output validation
- retry only when the normalized issue is retryable and the current attempt is below the configured maximum attempts
- preserve attempt, max attempts, remaining or next attempt, issue kind, issue code, and retryable outcome in inspectable state or events
- fail terminally with the original structured error or a structured retry-exhaustion error when no retry remains
- allow provider adapters to own bounded transport/protocol retries with explicit backoff and backup-model selection
- keep provider adapters responsible for classifying remote failures, while higher-level runtimes own semantic retries and run/turn terminal state

**Provider classification**
- provider timeouts and connection errors are retryable
- remote HTTP `408`, `409`, `425`, `429`, `500`, `502`, `503`, and `504` are retryable provider failures
- provider `429` must be classified distinctly from other HTTP failures
- provider authentication and authorization failures are not retryable unless a future credential-refresh contract explicitly owns recovery
- provider failure details must include provider name, model or endpoint, timeout configuration, remote status code when available, provider request id when available, operation, and redacted remote context
- provider retry details must include attempt, max attempts, attempts remaining, active model, primary model, backup model when configured, and retry-after header when available
- provider adapters may retry malformed provider payloads only when no caller-visible partial output has been emitted

**Conversational mode specifics**
- agent LLM completion retries must be bounded by `max_llm_completion_retries + 1` total attempts
- provider failures are retryable only if the provider error is retryable and no response text has already streamed to the caller-visible event history
- empty completions with no tool calls and no final answer are retryable within the same LLM attempt budget
- corrective retry prompts may be appended for empty completions when they are recorded as part of the conversation context
- governed tool execution retries must be bounded by `max_tool_execution_retries`, and budget exhaustion must be represented as an explicit tool result and event
- exceeding native tool-calling iteration limits must fail with a structured workflow error, not an unbounded retry loop

**Structured mode specifics**
- structured task retries must be bounded by the run's configured max attempts
- every structured task attempt must run within the configured timeout
- structured task timeouts, retryable provider failures, invalid JSON output, and local output validation failures are retryable while attempts remain
- retry prompts or retry issue context must be derived from the last failed attempt and passed through the structured task's message builder
- structured task run status must move through explicit running, retrying, completed, failed, or cancelled states
- fallback provider/model selection may be selected only by explicit task policy
- final structured task failure must preserve the terminal structured error on the run record and emit a failed event

**Forbidden**
- mode-local retry loops that bypass the shared LLM execution policy
- provider adapters silently retrying without an explicit provider retry budget, backoff policy, and observable retry signal
- retrying conversational completions after visible streamed text has been emitted
- treating provider-side schema strictness as a substitute for local validation and retry classification
- hiding fallback-provider selection or retry exhaustion from run, turn, event, log, or observability surfaces
- converting repeated empty, invalid, or unvalidated LLM output into successful empty output

**Evidence**
- retry decisions flow through the shared LLM retry decision contract
- agent and structured task events expose retry scheduling, issue code, issue kind, attempt number, and maximum attempts
- terminal run or turn state distinguishes success, cancellation, retry exhaustion, provider failure, timeout, and validation failure

### LLM-RUN-001: Runs And Events Must Be Durable Or Inspectable

**Rule**
LLM-backed execution that spans time, turns, batches, or owned tasks must preserve durable or inspectable state rather than disappearing into transient logs.

**Required**
- preserve run identity and lifecycle state
- preserve ordered events or equivalent inspectable execution history
- preserve enough failure context for operational diagnosis

**Conversational mode specifics**
- run, turn, tool-call, and event history should remain inspectable

**Structured mode specifics**
- run state, attempt state, output, and terminal failure detail should remain inspectable

**Forbidden**
- fire-and-forget LLM execution with no inspectable owner or terminal state
- failures that exist only in transient stdout or one-off logs

**Evidence**
- run state and event history are available through stable operational surfaces or durable stores

### LLM-PROMPT-001: Prompts Must Be Explicitly Versioned And Version Propagation Must Stay Observable

**Rule**
Concrete prompts must be explicitly versioned, and the effective prompt version used for execution must propagate through inspectable runtime state, logging, and observability.

**Required**
- assign an explicit version identifier to every concrete system prompt, developer prompt, task prompt, template, or prompt bundle that affects runtime behavior
- represent concrete prompts as first-class definitions with stable prompt identity separate from agent or task identity
- treat prompt version changes as intentional behavioral changes rather than silent text edits
- preserve the effective prompt version used by an execution in durable or inspectable runtime state
- include the effective prompt identity and version in relevant logs, events, traces, diagnostics, and operational metadata emitted for that execution
- keep prompt version fields machine-usable with stable names so operators can correlate behavior, failures, and regressions to a concrete prompt revision

**Canonical fields**
- use a stable prompt reference shape with `prompt_id`, `prompt_version`, `prompt_owner_kind`, `prompt_owner_id`, and `prompt_purpose`
- include `prompt_checksum` when a durable content hash is available

**Conversational mode specifics**
- agent runs and turns should preserve the prompt version or prompt bundle version that shaped the response behavior

**Structured mode specifics**
- structured task runs and attempts should preserve the prompt version that shaped structured generation and retry behavior

**Forbidden**
- editing production prompts in place with no explicit version change
- relying on git history alone as the runtime prompt version record
- emitting LLM execution logs or observability signals that cannot be tied back to the effective prompt version

**Evidence**
- prompt-bearing definitions expose explicit prompt identifiers and versions
- runtime state, logs, and observability outputs include the effective prompt reference

### LLM-EXPOSE-001: Operational Exposure Must Flow Through Application Modules

**Rule**
LLM-backed runtime capabilities must be exposed through application modules rather than directly from platform runtime internals.

**Required**
- expose public use cases through module-owned services or facades
- keep routes and transport concerns thin over those modules
- keep public surfaces intentionally narrow and operational

**Forbidden**
- transport code reaching directly into generic runtime internals
- shared substrate or platform runtime code owning public product or operational API behavior

**Evidence**
- public agent and structured task capability flows through module use cases and bootstrap wiring

### LLM-OBS-001: LLM Runtime Monitoring Must Reuse The Canonical Observability Runtime

**Rule**
LLM-backed execution must emit monitoring signals through the platform-owned observability runtime rather than through parallel mode-local telemetry stacks.

**Required**
- preserve request and trace correlation through LLM-backed execution where available
- emit relevant lifecycle signals through canonical events, metrics, tracing, and diagnostics surfaces when enabled
- include effective prompt version in canonical monitoring surfaces for prompt-driven execution when applicable
- keep monitoring fields machine-usable with stable names and codes

**Structured mode specifics**
- task/runtime monitoring should remain aligned with runtime truth rather than inventing a second monitoring path

**Forbidden**
- mode-specific monitoring infrastructure that bypasses the canonical observability runtime
- correlation loss across request, task, workflow, and LLM runtime boundaries without justification

**Evidence**
- LLM runtime monitoring appears in canonical diagnostics and observability outputs

### LLM-EVAL-001: Prompt, Model, And Runtime Changes Must Have Evaluation Or Smoke Coverage Appropriate To Risk

**Rule**
Changes to prompts, models, runtime configuration, or tool definitions must have evaluation or smoke coverage appropriate to the risk of behavioral change.

**Required**
- assign evaluation expectations based on impact: critical paths need tighter validation, lower-risk paths need lighter coverage
- keep evaluation results or explicit justified deferrals alongside change documentation
- prefer reproducible evaluation surfaces over ad hoc validation

**Forbidden**
- shipping prompt, model, or tool changes with no documented evaluation or smoke coverage plan
- treating happy-path smoke as sufficient for behaviorally significant prompt or model changes

**Evidence**
- change records include evaluation results or explicit justified deferral
- critical prompt changes have documented test or evaluation coverage

### LLM-SAFETY-001: Tool Use, User Instructions, And Model Outputs Must Pass Explicit Safety And Policy Checks Where Relevant

**Rule**
Tool execution paths, user-provided instructions, and model outputs must pass explicit safety or policy checks where the runtime supports or requires them.

**Required**
- define explicit safety boundaries for tool use when tools interact with sensitive operations
- validate user-provided instruction content when it feeds into model prompts or tool invocations
- check model outputs against relevant policy rules before using them in consequential operations
- make safety policy and check behavior explicit and inspectable

**Forbidden**
- executing tools or acting on model outputs without documented safety boundaries for consequential operations
- bypassing or silencing safety checks for convenience without explicit documented justification
- hiding safety policy decisions from operational visibility

**Evidence**
- safety checks and policy enforcement are represented in code and runtime events
- consequential tool executions and model-driven operations have documented safety boundaries

### LLM-COST-001: LLM Paths Must Expose Token, Cost, And Latency Metadata Where Operationally Meaningful

**Rule**
LLM-backed runtime paths must expose token usage, cost, and latency metadata where that information is operationally meaningful for monitoring, alerting, or diagnostics.

**Required**
- include token usage counts in runtime metadata or events when the runtime supports tracking
- include cost estimates or actual cost data when pricing models are meaningful to the operation
- include latency metadata when it helps diagnose performance or timeout behavior
- keep these fields machine-usable and consistent across similar runtime paths

**Forbidden**
- omitting token, cost, or latency metadata from runtime paths where it would meaningfully aid operational visibility
- burying cost and usage data in unstructured logs with no stable field paths

**Evidence**
- runtime metadata, logs, or events include token, cost, and latency fields where the runtime supports tracking
- operational surfaces expose meaningful usage and latency signals for LLM-backed paths

## Review Rejection Criteria
Reject a change if it:
- mixes shared LLM substrate behavior with concrete agent or structured task policy
- smuggles tool execution into a mode that does not explicitly support it
- claims structured output without explicit local validation
- introduces hidden approval, resume, retry, fallback, timeout, or cancellation behavior
- adds LLM retries outside the shared bounded retry policy
- retries conversational completion after caller-visible streamed text has been emitted
- adds or edits a concrete prompt without an explicit version
- fails to propagate effective prompt version through inspectable runtime state, logging, or observability
- exposes agent or structured task runtime directly from transport adapters
- adds a second telemetry path for LLM-backed behavior outside the canonical observability runtime
