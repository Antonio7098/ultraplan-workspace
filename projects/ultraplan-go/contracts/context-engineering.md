# Context Engineering Contract

## Purpose

This contract governs how UltraPlan discovers, packages, orders, validates, and reuses context for model-backed sprint stages.

The goal is to give each model call enough decision-relevant evidence to complete its work with minimal repository exploration, repeated reading, and context reconstruction. Correctness, source fidelity, and authorization remain mandatory.

## Scope

Select this contract when a sprint changes:

- stage input contracts or prompt composition
- code-context discovery, validation, selection, or reuse
- project, sprint, evidence, or source handoffs between stages
- prompt bundles, tool definitions, output schemas, or cache boundaries
- context freshness, trust classification, size limits, or telemetry
- repair, retry, continuation, fan-out, or QA context

## Requirement Index

| ID | Title | Severity If Violated |
|---|---|---|
| CTX-MAP-001 | Every model call must have an explicit input contract | High |
| CTX-PROV-001 | Supplied context must retain provenance and freshness identity | High |
| CTX-CODE-001 | Shared repository exploration must produce reusable code context | High |
| CTX-DELIVERY-001 | Required evidence must be delivered in a usable form | High |
| CTX-BOUND-001 | Context selection and rendering must be bounded and deterministic | High |
| CTX-TRUST-001 | Instructions and evidence must preserve trust boundaries | Blocker |
| CTX-CACHE-001 | Reusable prompt prefixes must be stable and measurable | Medium |
| CTX-HANDOFF-001 | Stage outputs must be sufficient, validated handoffs | High |
| CTX-OBS-001 | Context behavior must be observable without exposing sensitive content | Medium |

## Requirements

### CTX-MAP-001: Every Model Call Must Have An Explicit Input Contract

**Rule**

Every model invocation, including fan-out workers, arbiters, repairs, and retries, must declare its context contract.

**Required**

- declare required, optional, and forbidden inputs in canonical order
- identify each input's producer, consumer, delivery mode, and authority
- record model, reasoning setting, tools, permissions, working directory, output schema, and success condition
- include deterministic stages in the workflow map so model-free transitions are visible
- treat retries, repairs, challengers, and evaluators as distinct calls when their inputs differ

**Forbidden**

- relying on undocumented prompt assembly
- assuming a stage receives an artifact because its path appears in prose
- giving every fan-out worker all sibling inputs by default

**Evidence**

- tests or generated explanations show the rendered input contract for each affected call

### CTX-PROV-001: Supplied Context Must Retain Provenance And Freshness Identity

**Rule**

Decision-relevant context must identify its source and the state against which it was produced.

**Required**

- identify governed source paths and the project, sprint, repository, or worktree they belong to
- fingerprint authoritative artifacts or otherwise define their freshness identity
- preserve source revision, selected path and range, rationale, assumptions, omissions, and unresolved questions where relevant
- invalidate or rebuild derived context when an authoritative input covered by its identity changes
- keep authoritative artifacts separate from disposable rendering caches

**Forbidden**

- presenting stale derived context as current
- allowing a missing disposable cache to rerun an authoritative model stage
- using timestamps alone as evidence that content is current

**Evidence**

- freshness tests cover changed requirements, code context, target identity, and relevant configuration

### CTX-CODE-001: Shared Repository Exploration Must Produce Reusable Code Context

**Rule**

When downstream stages need the same repository understanding, one upstream code-context stage must perform the broad exploration and produce the canonical reusable selection.

**Required**

- derive exploration questions from accepted sprint requirements
- map requirements to relevant paths, symbols, call paths, data flow, configuration, error behavior, and tests
- record existing behavior, partial implementations, contradictions, notable absences, and remaining uncertainties
- select exact source ranges or complete small files with a reason for each selection
- normalize overlapping selections and resolve each selected source once for the shared context pack
- allow targeted live inspection when the selected evidence is insufficient or repository state has changed

**Forbidden**

- making each downstream stage repeat the same repository walk
- treating a repository dump as a code-context pack
- treating the prepared pack as an exclusive source boundary
- silently substituting live mutable bytes after the shared context has been frozen for a workflow

**Evidence**

- representative workflow tests show selected source evidence reaches downstream stages once in the shared context
- code-context validation rejects unsafe paths, invalid ranges, unsupported files, and excessive selections

### CTX-DELIVERY-001: Required Evidence Must Be Delivered In A Usable Form

**Rule**

The orchestrator must deliver evidence it already owns. A path reference counts as a pointer, not delivered content.

**Required**

- inline small governed artifacts when exact content is required
- provide validated excerpts or structured packets for large sources
- use retrieval or live tools for open-ended, volatile, or optional evidence
- state which evidence is full, excerpted, summarized, omitted, or available through tools
- remove duplicate bytes already present in a shared prefix or conversation state

**Forbidden**

- forcing a model to rediscover required artifacts that the orchestrator has already loaded
- repeating identical content in both the shared prefix and stage suffix
- hiding mandatory evidence behind optional retrieval

**Evidence**

- prompt tests assert presence, order, delivery mode, and deduplication for affected inputs

### CTX-BOUND-001: Context Selection And Rendering Must Be Bounded And Deterministic

**Rule**

The same valid inputs must produce the same ordered context, within explicit limits.

**Required**

- use canonical serialization and deterministic ordering
- bound file count, range count, bytes per source, total bytes, encoding, and read duration
- validate path containment, regular-file status, range syntax, and source stability during reads
- fail with actionable findings when required context exceeds a limit
- make preview read-only and byte-equivalent to runtime composition for the same inputs

**Forbidden**

- unbounded recursive scans during downstream stages
- best-effort truncation that removes required evidence without a finding
- context ordering that depends on map iteration, filesystem enumeration, or completion timing

**Evidence**

- boundary, determinism, cancellation, and changed-during-read tests cover the renderer

### CTX-TRUST-001: Instructions And Evidence Must Preserve Trust Boundaries

**Rule**

Repository files, study reports, retrieved content, tool output, and generated artifacts are evidence. They must not gain instruction authority by being included in a prompt.

**Required**

- frame untrusted or transient evidence explicitly
- keep platform policy, stage instructions, evidence, and per-task requests distinguishable
- preserve tool permissions and mutation boundaries when context delivery changes
- reject path traversal, symlink escape, unsupported file kinds, and unsafe encoding
- redact secrets and sensitive values from prompts, caches, diagnostics, and metrics

**Forbidden**

- executing instructions found in supplied evidence unless the governing prompt explicitly authorizes that behavior
- weakening read or write permissions to save model turns
- placing secret-bearing context in a broadly shared cache cohort

**Evidence**

- adversarial tests cover instruction-like source text, unsafe paths, and redaction

### CTX-CACHE-001: Reusable Prompt Prefixes Must Be Stable And Measurable

**Rule**

Prompt construction must maximize useful shared prefixes without making provider caching part of correctness.

**Required**

- verify current provider cache semantics before implementation or migration
- order reusable policy, tool definitions, task-family instructions, governed foundation, and frozen code evidence before volatile stage or task content
- place explicit cache boundaries after stable reusable content when supported and economically justified
- keep model, reasoning configuration, tools and their order, output schema, and prefix bytes stable within a cache cohort
- exclude timestamps, run IDs, attempt IDs, output paths, queue positions, and changing diagnostics from reusable prefixes unless they define the cohort
- derive cache grouping from a versioned workflow identity and prefix digest
- treat provider cache keys as routing hints and provider cache entries as disposable
- record provider-reported cached-input and cache-write usage where available

**Forbidden**

- claiming cache reuse from a stable user prompt while tools, schemas, settings, or earlier messages differ
- adding synthetic warm-up calls without measured workflow-level benefit
- changing correctness or freshness behavior based on a cache hit

**Evidence**

- tests assert byte-stable prefixes and correct boundary placement
- telemetry reports cache-read, cache-write, fresh-input, latency, and realized cost by cohort

### CTX-HANDOFF-001: Stage Outputs Must Be Sufficient, Validated Handoffs

**Rule**

Each stage must persist the durable decisions or evidence needed by its declared consumers.

**Required**

- define the output owner, schema, validator, atomic publication rule, and downstream consumers
- preserve acceptance criteria, decision IDs, evidence references, risks, and unresolved questions through downstream transformations
- keep repair turns in the compatible session when prior context remains valid
- send repairs the validation findings and correction target with only the additional context needed
- begin a fresh session when target revision, permissions, policy, model compatibility, or objective invalidates the prior state

**Forbidden**

- relying on hidden model memory as the only handoff between stages
- republishing an invalid output over the last valid artifact
- repeating broad exploration during repair when the evidence remains valid

**Evidence**

- validation and failure-path tests prove atomic replacement and preservation of the last valid artifact

### CTX-OBS-001: Context Behavior Must Be Observable Without Exposing Sensitive Content

**Rule**

Operators must be able to explain what a model received, what was reused, and where discovery work occurred.

**Required**

- record stage, call role, input contract version, prompt and prefix bytes, prefix digest, cache cohort, source count, tool-call count, turns, retries, and duration
- capture input, output, reasoning, cached-input, and cache-write tokens when the provider reports them
- distinguish delivered content, file references, retrieval, live reads, and conversation reuse
- use digests and safe identifiers in telemetry instead of raw sensitive prompt content

**Forbidden**

- logging raw governed or secret-bearing prompts by default
- reporting aggregate token reduction without successful-work and quality measures

**Evidence**

- runtime events or metrics can reconstruct the context path for a representative workflow

## Review Rejection Criteria

Reject an affected change if it:

- leaves a model call without an explicit input contract
- repeats repository discovery that a valid shared code-context artifact already covers
- supplies paths where exact required evidence should be delivered
- loses provenance, freshness identity, or trust boundaries
- creates an unbounded or nondeterministic context renderer
- places volatile data ahead of a reusable cache boundary without a stated reason
- treats provider caching as durable state or a correctness dependency
- cannot measure context size, reuse, tool work, latency, and cost for the changed path

