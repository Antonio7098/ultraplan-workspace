# Observability Contract

## Purpose
This contract defines the minimum observability requirements for runtime behavior.

It governs:
- structured logging
- correlation and trace propagation
- operational event emission
- diagnostics surfaces
- debug and verbosity controls
- health, preflight, and readiness truthfulness
- task and batch visibility
- alert, metrics, and automation expectations where applicable

## Scope
This contract applies when a change:
- introduces a new failure path
- adds task, batch, or orchestration execution
- adds runtime, provider, subprocess, or dependency integration
- adds diagnostics, health, or monitoring surfaces
- changes runtime logging or event behavior

## How To Use This Contract
Use this contract to determine:
- what operational signals must exist
- what failures must be visible to operators
- what evidence review should expect for runtime observability

## Requirement Index

| ID | Title | Applies To | Severity If Violated |
| --- | --- | --- | --- |
| OBS-CORE-001 | Failures must produce structured operational signals | all operational failure paths | Blocker |
| OBS-CORR-001 | Correlation identifiers must survive subsystem boundaries | commands, requests, tasks, workflows, dependency calls | High |
| OBS-HEALTH-001 | Health and preflight surfaces must reflect operational truth | health, readiness, preflight surfaces | Blocker |
| OBS-DIAG-001 | Diagnostics surfaces must expose operator-relevant state | diagnostics output and runtime state | Medium |
| OBS-DEBUG-001 | Debug and verbosity controls must be explicit and safe | logs, diagnostics, CLI/runtime output | Medium |
| OBS-TASK-001 | Owned work must have visible terminal state | task runners, batches, jobs, workers | Blocker |
| OBS-ALERT-001 | High-severity signals must be machine-usable for alerting | events, metrics, alerts | Medium |
| OBS-PII-001 | Logs, traces, metrics, and events must not expose secrets, credentials, or unnecessary user content | all observability outputs | High |

## Requirements

### OBS-CORE-001: Failures Must Produce Structured Operational Signals

**Rule**
Every meaningful failure path must emit a structured operational signal that operators can inspect.

**Applies when**
- handling command or request failures
- handling startup or shutdown failures
- handling runtime, dependency, workflow, or persistence failures
- handling task, batch, or worker failures

**Required**
- emit structured logs with stable codes when failures occur
- emit operational events where the runtime supports them
- preserve enough redacted detail for debugging
- make the failure visible in at least one immediate channel and one durable or inspectable channel

**Forbidden**
- swallowing failures silently
- logging success when the operation actually failed or degraded
- returning empty success-shaped results to hide operational failure

**Evidence**
- log payloads contain stable codes and structured detail
- runtime events or equivalent operator-visible signals exist
- failure paths are observable in diagnostics, task state, or logs

### OBS-CORR-001: Correlation Identifiers Must Survive Subsystem Boundaries

**Rule**
Command, request, task, and workflow correlation identifiers must survive subsystem boundaries where safe and useful.

**Applies when**
- handling commands or HTTP requests
- starting owned tasks
- calling runtimes, subprocesses, providers, or dependencies
- running workflows or orchestration

**Required**
- preserve command, request, run, task, and trace identifiers across downstream calls where possible
- attach correlation metadata to owned work and workflow execution
- include correlation identifiers in structured failures and events
- when LLM-backed execution is involved, preserve the effective prompt reference alongside correlation metadata in logs, traces, events, and diagnostics where that execution is surfaced

**Forbidden**
- generating unrelated identifiers for the same runtime path without reason
- dropping correlation identifiers at subsystem transitions without justification

**Evidence**
- logs, events, and task records contain command, request, run, task, or trace metadata
- LLM-backed logs, events, spans, or diagnostics include stable prompt identity fields when applicable
- diagnostics or task snapshots surface correlation fields

### OBS-HEALTH-001: Health And Preflight Surfaces Must Reflect Operational Truth

**Rule**
Health, readiness, and preflight surfaces must reflect whether the process or command can safely perform its intended workload.

**Applies when**
- implementing health commands, health endpoints, or preflight checks
- changing dependency checks
- changing startup validation or readiness semantics

**Required**
- distinguish liveness from readiness
- fail readiness, health, or preflight when required dependencies or required capabilities are unavailable
- represent degraded state when non-required capabilities are impaired

**Forbidden**
- returning healthy only because the process is alive
- masking required dependency failure behind optimistic health output

**Evidence**
- health or preflight handlers distinguish `live`, `ready`, `degraded`, `ok`, or `fail` states as appropriate
- readiness, health output, or preflight exit status changes when required dependencies fail

### OBS-DIAG-001: Diagnostics Surfaces Must Expose Operator-Relevant State

**Rule**
The system must expose a stable diagnostics surface for in-process operational inspection.

**Applies when**
- adding diagnostics commands, endpoints, files, or status output
- introducing operational state worth inspecting
- adding new runtime subsystems

**Required**
- centralize diagnostics rather than scattering hidden debug routes, flags, or files
- expose recent or current state that helps inspect operational behavior
- extend the canonical diagnostics surface when new operational state is introduced

**Forbidden**
- subsystem-specific hidden debug routes, flags, or files with no canonical operator path
- diagnostics that omit critical recent failure state

**Evidence**
- a canonical diagnostics surface exists
- new operational state is inspectable through that surface or explicitly justified elsewhere

### OBS-DEBUG-001: Debug And Verbosity Controls Must Be Explicit And Safe

**Rule**
Debug, trace, and verbosity controls must be documented, deliberate, and safe for the surface where they appear.

**Applies when**
- adding or changing runtime logs, debug modes, trace output, or diagnostic verbosity
- exposing debug controls through CLI flags, environment variables, config, or operator surfaces

**Required**
- provide documented controls for enabling additional diagnostics where operationally useful
- keep default verbosity safe for humans, automation, and production use
- route debug/progress output through the appropriate diagnostic channel for the surface
- preserve redaction and correlation rules at elevated verbosity levels

**Forbidden**
- hidden debug modes with undocumented activation paths
- debug output that contaminates machine-readable command output
- elevated verbosity that exposes secrets, credentials, or unnecessary user content

**Evidence**
- debug/verbosity controls are documented and tested where stable
- elevated diagnostic output remains redacted and channel-safe

### OBS-TASK-001: Owned Work Must Have Visible Terminal State

**Rule**
Owned asynchronous work, batch work, workers, and long-running operations must be observable and end in an explicit terminal state.

**Applies when**
- spawning tasks
- scheduling jobs
- adding asynchronous worker behavior
- adding resumable or parallel batch behavior

**Required**
- assign task identity and purpose metadata
- persist or expose task status transitions
- capture failure detail with stable codes where possible
- cancel or reconcile outstanding work on shutdown

**Forbidden**
- fire-and-forget work with no owner or task id when the result matters
- tasks whose failures exist only in stdout
- ambiguous task end states

**Evidence**
- task records or snapshots show lifecycle transitions
- failed tasks surface terminal details for review and diagnostics

### OBS-ALERT-001: High-Severity Signals Must Be Machine-Usable For Alerting

**Rule**
Operational signals should be machine-usable for metrics and alerting rather than free-form text matching.

**Applies when**
- emitting operational events
- defining high-severity codes
- adding metrics or alert integrations

**Required**
- use stable codes for alertable failures
- preserve structured severity and component fields
- keep alerting and metrics code-driven where possible

**Forbidden**
- relying only on prose log messages for alerting
- collapsing unrelated failures into one uninformative code

**Evidence**
- emitted events include stable code, severity, and component fields
- metrics and alert rules can key off machine-readable fields

### OBS-PII-001: Logs, Traces, Metrics, And Events Must Not Expose Secrets, Credentials, Or Unnecessary User Content

**Rule**
Observability outputs must not expose secrets, credentials, or unnecessary user content.

**Applies when**
- emitting structured logs, traces, metrics, or events
- adding new observability signals
- introducing new data fields into existing observability outputs

**Required**
- redact or omit secrets such as API keys, tokens, passwords, and authorization headers
- redact or omit unnecessary user content such as full personal data, payment information, or sensitive payload fields
- keep diagnostic structure intact after redaction
- prefer explicit redaction over omitting entire fields

**Forbidden**
- emitting raw secrets or credentials in log fields, span attributes, or event payloads
- including full user-provided content in traces or logs without classification and selective exclusion
- solving redaction concerns by dropping entire fields when narrower redaction would preserve diagnostic value

**Evidence**
- observability outputs exclude raw secrets and unnecessary user content
- redaction is explicit and deterministic where applied
- diagnostic structure remains intact in redacted outputs

## Review Rejection Criteria
Reject a change if it:
- introduces a new failure path with no structured operational signal
- drops correlation identifiers across a meaningful boundary without reason
- makes readiness, health, or preflight optimistic when required dependencies are unavailable
- launches owned work with no visible terminal status
- relies on unstructured text instead of machine-readable failure metadata for operator-critical behavior
- adds debug or verbose output that bypasses redaction, correlation, or surface-specific output-channel rules
