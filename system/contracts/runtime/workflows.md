# Workflow Contract

## Purpose
This contract defines when workflows are justified and how orchestration boundaries must behave.

It governs:
- workflow or batch-orchestration eligibility
- orchestration ownership
- retry, cancellation, and compensation semantics
- workflow boundary design
- engine encapsulation
- workflow observability expectations

## Scope
This contract applies when a change:
- adds or edits a workflow
- introduces orchestration runtime integration
- coordinates multiple steps, modules, tasks, services, or external systems

## Requirement Index

| ID | Title | Applies To | Severity If Violated |
| --- | --- | --- | --- |
| WF-SCOPE-001 | Workflows must be used only for real orchestration | workflow additions | Medium |
| WF-BOUNDARY-001 | Workflow engines must stay behind app-owned boundaries | orchestration runtime integration | High |
| WF-STATE-001 | Workflow outcomes must be explicit and inspectable | workflow execution paths | High |
| WF-RETRY-001 | Retry and cancellation semantics must be explicit | long-running and retrying workflows | High |
| WF-IDEMPOTENCY-001 | Retried workflow steps must be idempotent or protected by dedupe or idempotency keys | workflow retry behavior | High |
| WF-COMP-001 | Workflows with partial external side effects must define compensation or reconciliation behavior | workflows with partial side effects | High |
| WF-VERSION-001 | Long-running workflow definitions must be versioned when behavior changes | workflow definition changes | High |

## Requirements

### WF-SCOPE-001: Workflows Must Be Used Only For Real Orchestration

**Rule**
Workflows are for multi-step orchestration, batch execution, or resumable coordination, not as a default home for ordinary business logic.

**Required**
- use workflows when logic spans multiple steps, modules, services, systems, or durable task state
- use workflows when retries, cancellation, compensation, or resumability matter

**Forbidden**
- moving trivial one-step logic into `workflows/`
- creating a second service or command layer with no clear ownership

**Evidence**
- workflow reasoning explains why orchestration is warranted
- ordinary business logic remains in domain, module, or use-case layers

### WF-BOUNDARY-001: Workflow Engines Must Stay Behind App-Owned Boundaries

**Rule**
Workflow engine details must remain encapsulated behind app-owned runtime or facade boundaries.

**Required**
- wrap engine-specific integration in app-owned runtime helpers
- expose narrow orchestration contracts to modules

**Forbidden**
- making general business services or command flows depend deeply on raw engine internals without explicit orchestration ownership

**Evidence**
- modules depend on app-owned workflow helpers or protocols
- engine-specific types are contained at the orchestration boundary where possible

### WF-STATE-001: Workflow Outcomes Must Be Explicit And Inspectable

**Rule**
Workflow execution must end in explicit, inspectable state transitions.

**Required**
- make terminal outcome visible
- preserve step-level failure context when failures matter
- keep partial completion and inconsistent-state outcomes explicit

**Forbidden**
- ambiguous or hidden workflow terminal states
- silent step failure with implied success

**Evidence**
- workflow status or events expose terminal outcomes and failure detail

### WF-RETRY-001: Retry And Cancellation Semantics Must Be Explicit

**Rule**
Retry, cancellation, and compensation behavior must be defined rather than implied.

**Required**
- declare retryable conditions and attempt limits
- make cancellation and timeout behavior explicit
- define compensation expectations where partial completion matters

**Evidence**
- reasoning and implementation expose retry/cancellation semantics
- tests or runtime evidence cover important lifecycle behavior

### WF-IDEMPOTENCY-001: Retried Workflow Steps Must Be Idempotent Or Protected By Dedupe Or Idempotency Keys

**Rule**
Retried workflow steps must be idempotent or protected by dedupe or idempotency keys so that retry does not produce duplicate side effects.

**Required**
- design workflow steps for idempotent execution where retry is expected
- use dedupe or idempotency keys to prevent duplicate side effects when idempotent design is not practical
- make idempotency stance explicit in the workflow design and documentation

**Forbidden**
- retrying non-idempotent operations without explicit dedupe or idempotency key protection
- assuming retry safety without explicit design consideration

**Evidence**
- workflow step design documents idempotency strategy or explicit dedupe protection
- retry behavior is safe for the intended retry budget

### WF-COMP-001: Workflows With Partial External Side Effects Must Define Compensation Or Reconciliation Behavior

**Rule**
Workflows that produce partial external side effects must define explicit compensation or reconciliation behavior for failure and recovery paths.

**Required**
- identify workflow steps that produce external side effects visible to other systems
- define compensation behavior such as rollback, undo, or cleanup when subsequent steps fail
- define reconciliation behavior when compensation is not possible or practical
- make compensation and reconciliation logic explicit in workflow design and observable in runtime state

**Forbidden**
- workflows with partial external side effects and no documented compensation or reconciliation strategy
- leaving partial side effects unaddressed in failure scenarios

**Evidence**
- workflow design documents compensation or reconciliation behavior for external side effects
- compensation and reconciliation are testable or verifiable through explicit runtime behavior

### WF-VERSION-001: Long-Running Workflow Definitions Must Be Versioned When Behavior Changes

**Rule**
Long-running workflow definitions must be versioned so that behavioral changes are tracked and in-flight executions are not disrupted unintentionally.

**Required**
- assign explicit version identifiers to workflow definitions
- propagate effective workflow version through inspectable runtime state and events
- avoid breaking changes to in-flight workflow definitions without migration strategy

**Forbidden**
- editing workflow definitions in place with no version change
- changing behavior of a running workflow without version awareness

**Evidence**
- workflow definitions expose explicit version identifiers
- runtime state includes effective workflow version for in-flight executions

## Review Rejection Criteria
Reject a change if it:
- adds a workflow for trivial one-step logic
- leaks engine internals across ordinary application surfaces
- leaves workflow terminal outcomes ambiguous
- introduces retries or cancellation with no explicit semantics
