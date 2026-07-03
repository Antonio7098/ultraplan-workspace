# Privacy And Data Contract

## Purpose

This contract defines how user data, sensitive data, PII, and user-generated content must be handled.

It governs:

- data classification
- data minimization
- retention and deletion
- redaction and pseudonymization
- user content handling
- auditability
- safe logging and analytics

## Scope

This contract applies when a change:

- collects, stores, processes, logs, exports, or transmits user data
- handles PII, secrets, private content, customer data, or generated content
- adds analytics, telemetry, evaluation, model training, or debugging surfaces
- changes retention, deletion, or data access behaviour

## Requirement Index

| ID | Title | Applies To | Severity If Violated |
| --- | --- | --- | --- |
| DATA-CLASS-001 | Data must be classified before storage or transmission | new data fields/flows | High |
| DATA-MIN-001 | Store the minimum data required for the feature | all persisted data | High |
| DATA-PII-001 | PII and sensitive data require explicit handling rules | PII paths | Blocker |
| DATA-LOG-001 | Logs and telemetry must avoid unnecessary user data | observability paths | Blocker |
| DATA-RET-001 | Retention and deletion must be defined for durable data | persisted data | High |
| DATA-EXPORT-001 | Data export and sharing must be intentional and scoped | exports/integrations | High |
| DATA-AUDIT-001 | Sensitive data access must be attributable | admin/support/internal access | High |
| DATA-AI-001 | AI/model use of user content must be bounded and inspectable | AI/LLM paths | High |

## Core Principles

1. Do not collect data unless the product or operation needs it.
2. Do not persist raw user content when derived or referenced data is enough.
3. Logs are not a data lake.
4. PII and secrets require explicit classification and redaction.
5. Retention is a design decision, not an afterthought.
6. Access to sensitive data must be attributable.
7. AI processing must not silently expand data use.

## Requirements

### DATA-CLASS-001: Data Must Be Classified Before Storage Or Transmission

**Rule**
New data fields and flows must have an explicit classification before they are persisted, logged, exported, or sent to providers.

**Canonical classifications**
- public
- internal
- user-content
- personal-data
- sensitive-personal-data
- secret
- regulated
- operational-metadata
- derived/non-authoritative

**Required**
- classify new durable fields and high-volume telemetry fields
- document whether data is authoritative, derived, cached, or diagnostic
- mark data that must not leave the system or must not enter logs

**Forbidden**
- adding opaque JSON blobs containing unknown user data without classification
- sending unclassified user content to analytics, providers, or logs

**Evidence**
- schema/docs/reasoning identify classification for new data flows

### DATA-MIN-001: Store The Minimum Data Required For The Feature

**Rule**
Persist only the data needed to support the feature, operational requirement, or audit obligation.

**Required**
- prefer IDs, hashes, references, summaries, or derived state where raw content is not required
- separate operational metadata from user content
- avoid duplicating sensitive data across tables, caches, queues, logs, and events

**Forbidden**
- storing raw payloads “just in case”
- duplicating PII into analytics or event stores without explicit need
- using durable storage for temporary execution context without reason

**Evidence**
- reasoning explains why persisted data is necessary
- raw content storage has explicit justification

### DATA-PII-001: PII And Sensitive Data Require Explicit Handling Rules

**Rule**
PII and sensitive data must have defined access, storage, redaction, retention, and transmission controls.

**Required**
- identify PII and sensitive fields at schema and boundary layers
- redact or pseudonymize where full values are not needed
- protect exports and admin views
- avoid exposing sensitive fields in default API responses

**Forbidden**
- returning sensitive fields by default because they exist in a model
- joining sensitive data into broad read models unnecessarily
- exposing full PII in errors, logs, analytics, traces, or client state

**Evidence**
- API responses and logs show safe subsets
- tests cover sensitive-field omission where relevant

### DATA-LOG-001: Logs And Telemetry Must Avoid Unnecessary User Data

**Rule**
Operational signals must be useful without becoming uncontrolled stores of user data.

**Required**
- use identifiers, counts, categories, hashes, and redacted summaries where possible
- sanitize event/log payloads before emission
- keep logs structured and safe by default
- provide explicit allowlists for fields that may enter telemetry

**Forbidden**
- logging full request/response bodies by default
- logging prompts, completions, emails, files, or user messages without explicit policy
- relying on manual developer discretion for redaction of sensitive paths

**Evidence**
- telemetry schemas exclude sensitive fields by default
- redaction tests exist for known sensitive field names and payload shapes

### DATA-RET-001: Retention And Deletion Must Be Defined For Durable Data

**Rule**
Durable user data must have a retention and deletion stance.

**Required**
- define whether data is retained indefinitely, time-limited, user-deletable, audit-retained, or cache-only
- distinguish hard delete, soft delete, anonymization, and archival
- ensure derived data and caches are covered by deletion rules

**Forbidden**
- adding durable user data with no retention stance
- deleting primary data while leaving sensitive derived copies behind unintentionally

**Evidence**
- data docs or schema comments identify retention/deletion behaviour
- deletion tests cover important derived/cached data where applicable

### DATA-EXPORT-001: Data Export And Sharing Must Be Intentional And Scoped

**Rule**
Exports, provider calls, integrations, and downloadable artifacts must include only the intended data scope.

**Required**
- define export audience and data scope
- filter fields explicitly rather than dumping internal models
- protect files and generated reports from unauthorized access
- record sensitive exports where auditability matters

**Forbidden**
- exporting raw database rows or internal DTOs directly
- including hidden/private/internal fields by default
- writing exported data to public or predictable file locations

**Evidence**
- export models are explicit
- tests cover field inclusion/exclusion

### DATA-AUDIT-001: Sensitive Data Access Must Be Attributable

**Rule**
Internal/admin/support access to sensitive data must be attributable and reviewable where risk warrants.

**Required**
- preserve actor, reason, operation, and target identifiers for sensitive access
- separate ordinary user access from elevated/support/admin access
- avoid broad debug views exposing sensitive records without audit

**Forbidden**
- admin/support data access with no traceable actor
- hidden internal routes that expose sensitive data without access control

**Evidence**
- sensitive admin/support paths emit audit-safe access events

### DATA-AI-001: AI/Model Use Of User Content Must Be Bounded And Inspectable

**Rule**
User content sent to AI providers, evaluation systems, embeddings, vector stores, or training/eval datasets must be intentionally scoped and inspectable.

**Required**
- define what user content is sent, stored, embedded, cached, or retained
- minimize prompts and context to what the task needs
- avoid sending secrets or unnecessary PII to models
- preserve provider, model, prompt version, and content classification metadata where relevant

**Forbidden**
- silently reusing production user content for evals/training without policy
- embedding raw sensitive content when a safer representation is sufficient
- logging raw prompts/completions without classification and retention rules

**Evidence**
- AI data flows identify content scope and retention
- prompt/runtime logs use safe summaries or redaction where required

## Review Rejection Criteria

Reject a change if it:

- adds durable user data with no classification
- stores raw sensitive content without need
- logs PII, secrets, prompts, completions, emails, or files without policy
- exposes sensitive fields through default API responses
- sends user content to providers without an explicit scope
- adds export/download/admin surfaces without field scoping and access control
- creates data with no retention/deletion stance

---