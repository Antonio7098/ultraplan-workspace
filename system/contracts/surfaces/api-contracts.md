# API Contracts Contract

## Purpose

This contract defines how public and internal APIs must be designed, versioned, documented, and reviewed.

It governs:

- HTTP/resource/API shape
- schemas and compatibility
- authentication and authorization
- idempotency
- pagination and filtering
- error responses
- rate limits and quotas
- generated API documentation

## Scope

This contract applies when a change:

- adds or changes HTTP, RPC, GraphQL, WebSocket, webhook, or SDK-facing APIs
- changes request/response schemas
- changes public behaviour, status codes, auth, pagination, filtering, idempotency, or errors
- changes generated API documentation or clients

## Requirement Index

| ID | Title | Applies To | Severity If Violated |
| --- | --- | --- | --- |
| API-SCHEMA-001 | API schemas must be explicit and reviewable | public/internal APIs | High |
| API-COMPAT-001 | Breaking changes must be intentional and versioned | public APIs | Blocker |
| API-AUTH-001 | API access must enforce authn/authz server-side | protected APIs | Blocker |
| API-IDEMP-001 | Non-idempotent operations must define replay behaviour | writes/retries | High |
| API-PAGE-001 | List endpoints must define pagination/filter/sort semantics | collection APIs | Medium |
| API-ERROR-001 | API errors must preserve stable machine-readable codes | error responses | High |
| API-RATE-001 | Abuse-prone APIs must define limits and backpressure | public/high-cost APIs | High |
| API-DOC-001 | API documentation must match implementation | public/internal APIs | Medium |

## Core Principles

1. API contracts are product surfaces, not implementation details.
2. Schemas must be explicit enough for clients, tests, and review.
3. Backwards compatibility must be protected deliberately.
4. Authorization belongs on the server for every protected resource/action.
5. Retried writes must be safe or explicitly rejected.
6. Errors must be machine-usable and human-actionable.
7. Generated docs and actual behaviour must not drift.

## Requirements

### API-SCHEMA-001: API Schemas Must Be Explicit And Reviewable

**Rule**
API requests, responses, and events must have explicit schemas or contract definitions.

**Required**
- define request and response models for meaningful endpoints
- distinguish transport DTOs from domain/internal models
- avoid returning raw ORM/database/provider objects
- validate incoming payloads at the boundary
- define nullability, optional fields, defaults, and enum values deliberately

**Forbidden**
- exposing internal persistence models directly
- accepting arbitrary payload blobs without schema where structure matters
- silently ignoring unknown dangerous fields on write operations

**Evidence**
- OpenAPI/contract/schema output reflects changed endpoints
- tests cover schema validation and representative responses

### API-COMPAT-001: Breaking Changes Must Be Intentional And Versioned

**Rule**
Public API breaking changes must be explicit, versioned, documented, and coordinated.

**Breaking examples**
- removing or renaming fields
- changing field types or meanings
- changing required fields
- changing status-code semantics
- changing pagination or error shapes
- making previously accepted requests invalid

**Required**
- preserve backwards-compatible behaviour by default
- version or gate breaking changes
- document migration paths when clients are affected

**Forbidden**
- accidental breaking changes caused by internal refactors
- changing response shape without contract updates
- reusing a field name with a different meaning

**Evidence**
- compatibility tests or API diff checks detect breaking changes
- release notes identify intentional breaks

### API-AUTH-001: API Access Must Enforce Authn/Authz Server-Side

**Rule**
Protected APIs must authenticate actors and authorize every protected resource/action server-side.

**Required**
- derive actor context from verified credentials/session
- check object-level and action-level permissions
- avoid leaking existence of protected resources where policy forbids it

**Forbidden**
- relying on frontend route guards or hidden UI controls as security
- trusting actor, tenant, or role fields from caller payloads

**Evidence**
- tests cover unauthorized actor, wrong tenant, wrong object, and wrong action cases

### API-IDEMP-001: Non-Idempotent Operations Must Define Replay Behaviour

**Rule**
Write operations that may be retried, duplicated, or resumed must define idempotency semantics.

**Required**
- use idempotency keys, dedupe keys, natural constraints, or explicit no-retry policy where appropriate
- document which operations are safe to retry
- protect payment, notification, provider, and external side-effect paths

**Forbidden**
- retrying non-idempotent writes by default
- allowing duplicate external effects from repeated client/provider attempts
- hiding partial success as complete failure with no reconciliation path

**Evidence**
- tests cover duplicate/retried write requests where risk exists

### API-PAGE-001: List Endpoints Must Define Pagination/Filter/Sort Semantics

**Rule**
Collection APIs must define stable pagination, filtering, sorting, and total-count behaviour.

**Required**
- prefer cursor pagination for large or changing collections where suitable
- define default and maximum page sizes
- define sort stability and tie-breakers
- validate filters and reject unsupported combinations where needed

**Forbidden**
- unbounded list endpoints
- unstable pagination that duplicates or skips records unpredictably
- exposing raw query language without controls

**Evidence**
- list endpoint tests cover bounds, filters, and pagination continuity

### API-ERROR-001: API Errors Must Preserve Stable Machine-Readable Codes

**Rule**
API errors must expose stable codes and safe diagnostic fields while preserving full internal detail in logs/events.

**Required**
- include stable error code, message, correlation/request ID, and safe details
- map validation, auth, not-found, conflict, rate-limit, dependency, and internal errors distinctly
- avoid leaking secrets or internal stack traces to clients

**Forbidden**
- generic `500` responses for known caller errors
- success-shaped empty responses for operational failures
- exposing raw exception messages to public clients

**Evidence**
- API tests assert stable error shapes for important failure cases

### API-RATE-001: Abuse-Prone APIs Must Define Limits And Backpressure

**Rule**
High-cost, public, auth, search, upload, LLM, or provider-backed APIs must define limits and overload behaviour.

**Required**
- define request size, rate, concurrency, and cost limits where relevant
- return useful rate-limit or overload errors
- avoid unbounded provider/model/database cost from a single request

**Forbidden**
- unlimited expensive operations exposed to public clients
- hidden throttling that returns misleading success or generic failure

**Evidence**
- limits are documented/configured and failure shape is testable

### API-DOC-001: API Documentation Must Match Implementation

**Rule**
API documentation, schemas, examples, and generated clients must reflect actual runtime behaviour.

**Required**
- update OpenAPI/schema/docs with API changes
- keep examples valid and minimal
- document auth, errors, pagination, idempotency, and rate limits where relevant

**Forbidden**
- stale generated docs after schema changes
- undocumented fields relied on by clients

**Evidence**
- API docs generation or contract validation runs in checks

## Review Rejection Criteria

Reject a change if it:

- exposes internal models as public API responses
- makes accidental breaking changes without versioning/migration plan
- lacks server-side object/action authorization
- retries non-idempotent operations without protection
- adds unbounded list or high-cost endpoints
- returns generic or unsafe error responses
- changes APIs without updating schemas/docs/tests

---