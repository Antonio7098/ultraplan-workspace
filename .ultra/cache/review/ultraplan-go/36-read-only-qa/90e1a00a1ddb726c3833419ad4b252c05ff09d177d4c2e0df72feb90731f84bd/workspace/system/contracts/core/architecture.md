# Architecture Contract

## Purpose
This contract defines the core architectural rules for application code, including Go CLIs, services, and tools.

It governs:
- package and module boundaries
- dependency direction
- module public APIs
- composition structure
- adapter boundaries
- shared vs platform ownership

It does not own every operational concern.
Use the specialized contracts listed in [README.md](README.md) for failures, observability, testing, workflows, agents, and other surface-specific requirements.

## Scope
This contract applies when a change:
- adds or edits application modules
- introduces new use cases, services, or adapters
- changes dependency wiring or composition
- changes public module boundaries
- adds new transport or orchestration entrypoints

## How To Use This Contract
Use this contract in three ways:
- during reasoning, map the applicable requirement IDs to the sprint scope
- during implementation, preserve the required boundaries and dependency direction
- during review, verify conformance using the evidence guidance in each requirement

## Requirement Index

| ID | Title | Applies To | Severity If Violated |
| --- | --- | --- | --- |
| ARCH-CORE-001 | Module boundaries must remain explicit | new or changed modules | High |
| ARCH-CORE-002 | Dependency direction must point inward | cross-layer dependencies | Blocker |
| ARCH-LAYER-001 | Domain layer must remain pure | domain code | High |
| ARCH-LAYER-002 | Use cases depend on ports when a real seam exists | use cases, services, command flows | Blocker |
| ARCH-ENTRY-001 | Entrypoint adapters must stay thin | CLI, routes, webhooks, sockets | High |
| ARCH-MODULE-001 | Module public APIs must stay small and stable | module exports | Medium |
| ARCH-COMP-001 | Composition must happen through registrars | composition root and bootstrap | High |
| ARCH-SHARED-001 | Shared and platform code must stay domain-neutral | shared and platform packages | High |

## Core Principles

1. A module should represent one bounded context, cohesive capability, or stable ownership boundary.
2. Dependency direction should keep product behavior independent from process entrypoints, transports, vendors, and generic platform runtime details.
3. High-level policy depends on abstractions at volatile side-effect boundaries, not for every helper function.
4. Transport entrypoints and CLI command handlers are adapters, not business layers.
5. Infrastructure is replaceable.
6. Architectural seams must remain testable through public boundaries.
7. Do not introduce ports, facades, services, registrars, or module boundaries speculatively. They must protect a real seam, side-effect boundary, ownership boundary, or test seam.
8. In Go, package-level exports are the public API. Keep exported identifiers intentional and use package docs to state ownership when a module has no separate `internal` subpackage.
9. For local CLI code, concrete filesystem collaborators are acceptable when they are narrow, deterministic, and testable. Ports become mandatory when the dependency is volatile, external, long-running, concurrent, provider-backed, persistent, or hard to substitute in tests.

## Requirements

### ARCH-CORE-001: Module Boundaries Must Remain Explicit

**Rule**
Each module must represent a bounded context, cohesive capability, or ownership boundary with a small public API and clear internal ownership.

**Applies when**
- creating a new module
- expanding an existing module
- exposing one module to another

**Required**
- keep module internals behind the module boundary
- export only intentional public collaborators from the package root or equivalent public surface
- in Go, use unexported identifiers for unstable helpers and package documentation to clarify ownership
- keep domain, use-case, infra, and optional workflow responsibilities distinct

**Forbidden**
- importing another module's private/internal implementation details directly
- exporting concrete infra from a product module when callers only need a use-case service or value
- mixing unrelated bounded contexts in one module

**Evidence**
- module imports flow through public module APIs
- package exports remain small and intentional
- feature code does not reach into sibling module internals

### ARCH-CORE-002: Dependency Direction Must Point Inward

**Rule**
Dependency direction must point from outer layers toward inner business logic.

**Applies when**
- adding imports across layers
- introducing adapters or orchestration
- wiring dependencies in composition

**Required**
- entrypoints such as `cmd/...` delegate to application or command packages
- command, route, or transport adapters call module services, facades, or public package functions rather than owning product behavior
- product modules may depend on platform packages for generic capabilities and on explicitly allowed product modules for stable product relationships
- platform packages must not import product modules
- vendor/provider/process adapters must remain outside domain code and behind a narrow boundary when used by use cases

**Forbidden**
- domain code importing command, transport, provider, persistence, or process-execution packages
- product use cases calling provider SDKs, subprocess runners, databases, or long-running runtimes directly when those dependencies are volatile or need substitution
- entrypoints reaching directly into runtime/provider adapters, database/session internals, or container internals
- platform packages importing product modules

**Evidence**
- imports follow the allowed direction
- use cases rely on ports where the collaborator is external, volatile, long-running, persistent, concurrent, or hard to fake
- entrypoint code does not open direct infrastructure dependencies

### ARCH-LAYER-001: Domain Layer Must Remain Pure

**Rule**
Domain code must model business behavior and must not depend on transport, persistence, or vendor concerns.

**Applies when**
- editing domain files, domain packages, or pure model/value-object code
- introducing domain entities, value objects, or exceptions

**Required**
- keep entities, value objects, domain rules, and domain exceptions in the domain layer
- keep domain logic independent from ORM, transport, and SDK details

**Forbidden**
- ORM models in domain code
- request or response schemas in domain code
- sessions or vendor SDK usage in domain code

**Evidence**
- domain imports remain infrastructure-free
- domain tests do not require transport or DB wiring

### ARCH-LAYER-002: Use Cases Depend On Ports When A Real Seam Exists

**Rule**
Use cases and workflows must depend on narrow ports rather than concrete infrastructure when the collaborator crosses a meaningful side-effect, volatility, runtime, persistence, or test seam.

**Applies when**
- adding services, use cases, or command flows
- adding repository/provider access
- introducing orchestration
- introducing external process execution, network/provider calls, durable state, clocks/timers, concurrency, or filesystem behavior that is expensive or difficult to fixture

**Required**
- define ports as narrow interfaces, function types, or equivalent application-layer contracts near the consumer in Go
- implement those ports in platform, adapter, store, or runtime packages
- inject implementations during composition
- keep interfaces small and behavior-specific
- prefer concrete local helpers when no real seam exists yet and tests can use public behavior or temporary files safely

**Forbidden**
- importing concrete provider clients, process runners, database/session handles, or vendor SDKs into use cases
- leaking platform DB models, transport schemas, provider payloads, or subprocess details into use cases
- broad interfaces with unrelated capabilities
- creating interfaces only to satisfy ceremony when there is one stable in-process implementation and no substitution need

**Evidence**
- use-case constructors accept ports or facades for meaningful side-effect boundaries
- fake implementations can replace collaborators through public seams
- direct concrete dependencies are justified as local, narrow, deterministic, and testable

### ARCH-ENTRY-001: Entrypoint Adapters Must Stay Thin

**Rule**
CLI handlers, HTTP routes, webhooks, and websocket handlers must remain entrypoint adapters.

**Applies when**
- adding or editing entrypoints
- introducing request/response schemas
- wiring command or route dependencies

**Required**
- validate input at the adapter boundary
- resolve actor or context
- call one service or facade
- map results to transport output
- in CLI code, keep stdout/stderr formatting and exit-status mapping in the command/app layer

**Forbidden**
- opening sessions directly in command or route files
- calling runtime/provider adapters directly from entrypoints
- touching container internals except through thin dependency helpers
- reaching into private service fields such as `service._workflow` or `service._repo`
- putting filesystem discovery, domain validation, scheduling, provider execution, or durable workflow state machines directly in command handlers

**Evidence**
- route or command code delegates to a service or facade
- CLI command code delegates product behavior to a package-owned service or public function
- infrastructure access stays outside transport files

### ARCH-MODULE-001: Module Public APIs Must Stay Small And Stable

**Rule**
Every module must expose a small, stable public API and keep unstable details internal.

**Applies when**
- exporting module collaborators
- sharing commands or views across modules

**Required**
- export the main service or facade when needed
- export intentionally shared commands and views only when needed
- export the module bootstrap function where appropriate
- in Go, keep exported types/functions few, named for use cases, and documented when they are part of the stable surface

**Forbidden**
- exporting concrete repositories
- exporting SQL helpers
- exporting provider adapters
- exporting ORM models
- exporting broad managers or utility bags that expose unrelated module internals

**Evidence**
- module roots stay small
- shared callers depend only on stable public symbols
- `go doc`/package docs reveal a coherent public surface rather than internal mechanics

### ARCH-COMP-001: Composition Must Happen Through Registrars

**Rule**
Application assembly must happen at a clear composition boundary. As systems grow, compose module bundles or factories instead of scattering wiring through entrypoints.

**Applies when**
- wiring the app container
- adding modules to the app
- changing bootstrap behavior
- adding runtime/provider adapters, persistence stores, workers, diagnostics, or multiple implementations of a collaborator

**Required**
- keep process entrypoints thin, typically `cmd/<binary>/main.go -> internal/app`
- keep shared wiring in `internal/app`, `internal/platform/composition`, or an equivalent explicit composition package
- let modules expose constructors, bundles, or registrars when they have multiple collaborators or adapters
- centralize provider/runtime/store/logger/config wiring once those collaborators enter scope

**Forbidden**
- one god-object constructor that knows every internal detail in the system
- routes or external adapters reaching through the container into private service internals
- command handlers repeatedly constructing provider clients, stores, runtimes, worker pools, or diagnostics pipelines

**Evidence**
- modules expose typed bundles or public collaborators
- composition code depends on module bootstrap outputs, not module internals
- small CLI skeletons may use simple app-level dependency structs when they do not hide volatile collaborators or block later composition

### ARCH-SHARED-001: Shared And Platform Code Must Stay Domain-Neutral

**Rule**
`shared/` and `platform/` must remain free of bounded-context business behavior.

**Applies when**
- adding utilities to `shared/`
- adding runtime infrastructure to `platform/`

**Required**
- use `shared/` or equivalent only for generic cross-cutting code such as base errors, IDs, typed primitives, and protocol helpers
- use `platform/` for generic runtime infrastructure such as config loading, telemetry wiring, task execution, process execution, filesystem adapters, and orchestration wrappers
- keep product defaults, scaffolds, templates, prompts, validation rules, and workflow semantics in the owning product module unless they are genuinely generic

**Forbidden**
- embedding bounded-context logic in `shared/`
- embedding product-specific business behavior in `platform/`
- moving product behavior to shared/platform packages solely to avoid an import relationship

**Evidence**
- platform packages expose runtime services, not product policy
- shared packages remain reusable and domain-neutral
- product-specific defaults and artifacts are owned by product/workspace modules, not generic platform runtime

## Review Rejection Criteria
Reject a change if it:
- imports concrete volatile infra into use cases without a justified local/narrow exception
- adds direct route access to sessions, provider adapters, or container internals
- exports infra from a module root
- introduces a god repository, god service, or god container
- creates cross-module dependencies through private internals
- puts bounded-context behavior into `shared/` or `platform/`
- puts substantial product behavior, scheduling, runtime execution, or persistence logic directly in CLI command handlers

## Related Contracts
- [errors.md](errors.md)
- [observability.md](observability.md)
- [testing.md](testing.md)
- [documentation.md](../core/documentation.md)
- [security.md](../core/security.md)
- [release-and-versioning.md](../core/release-and-versioning.md)
