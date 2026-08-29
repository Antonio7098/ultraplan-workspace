# Security Contract

## Purpose

This contract defines the minimum security expectations for implementation and review.

It governs:

- authentication and authorization boundaries
- input validation and output encoding
- secrets handling
- dependency and supply-chain safety
- secure defaults
- injection and unsafe execution prevention
- file, path, subprocess, network, and dependency access safety

## Scope

This contract applies when a change:

- adds or changes authentication or authorization logic
- processes untrusted input
- exposes an API, CLI, webhook, tool, or frontend surface
- calls external systems
- handles secrets, credentials, tokens, cookies, or signing material
- adds dependencies, build steps, plugins, or code execution paths
- touches file system, shell, network, or deserialization behaviour

## Requirement Index

| ID | Title | Applies To | Severity If Violated |
| --- | --- | --- | --- |
| SEC-AUTHN-001 | Authentication state must be explicit and verified | identity-bearing paths | Blocker |
| SEC-AUTHZ-001 | Authorization must be checked at the object/action boundary | protected resources | Blocker |
| SEC-INPUT-001 | Untrusted input must be validated at boundaries | APIs, CLIs, tools, forms | High |
| SEC-INJECT-001 | Injection-prone operations must use safe APIs | DB, shell/subprocess, templates, queries | Blocker |
| SEC-SECRETS-001 | Secrets must never be hardcoded, logged, or exposed | all secret handling | Blocker |
| SEC-DEPS-001 | Dependencies must be intentional and supply-chain reviewed | package changes | High |
| SEC-FILES-001 | File and path access must be constrained | file system operations | High |
| SEC-NET-001 | Outbound network access must be intentional and bounded | runtime/provider/API calls | High |
| SEC-DESER-001 | Deserialization and dynamic execution must be restricted | parsers, plugins, scripts | Blocker |
| SEC-DEFAULT-001 | Security-sensitive defaults must fail closed | auth, config, exposure | High |

## Core Principles

1. Never trust caller-controlled input.
2. Authentication answers “who is this?”; authorization answers “may they do this?” Keep both explicit.
3. Authorization must happen at the trusted resource/action boundary.
4. Secrets are runtime inputs, not source code or logs.
5. Safe APIs beat escaping by hand.
6. Dependencies are part of the attack surface.
7. Security controls should fail closed.
8. Debug convenience must not create production exposure.

## Requirements

### SEC-AUTHN-001: Authentication State Must Be Explicit And Verified

**Rule**
Authenticated state must be derived from verified credentials, sessions, or tokens and must not be inferred from caller-controlled data.

**Required**
- verify tokens, sessions, signatures, local credentials, or configured identities before treating a request, command, or runtime action as authenticated
- preserve actor identity through explicit context objects
- distinguish anonymous, authenticated, service, and system actors
- reject malformed or expired credentials clearly

**Forbidden**
- trusting user IDs, emails, roles, or tenant IDs directly from request bodies, config files, or CLI args
- treating frontend-hidden controls as authentication
- silently falling back to privileged system identity

**Evidence**
- protected paths derive actor context from verified auth mechanisms
- tests cover missing, malformed, expired, and wrong-actor credentials

### SEC-AUTHZ-001: Authorization Must Be Checked At The Object/Action Boundary

**Rule**
Every protected object access or state-changing action must verify that the actor may perform that operation on that object.

**Required**
- check authorization after object identity is known and before returning or mutating protected data
- enforce tenant, ownership, role, capability, and object-level rules in trusted application logic
- check both read and write paths
- keep authorization decisions testable and observable at a safe level

**Forbidden**
- relying on obscured IDs, frontend filtering, hidden buttons, route guards, or CLI affordances alone
- checking only coarse route-level roles when object-level access matters
- mass assignment of fields the actor is not allowed to set

**Evidence**
- tests cover same-role/different-object and same-object/different-action cases
- review can point to the authorization decision for each protected operation

### SEC-INPUT-001: Untrusted Input Must Be Validated At Boundaries

**Rule**
External input must be parsed, validated, and normalized at the boundary before it reaches core logic.

**Required**
- validate request bodies, query parameters, headers, CLI args, config files, uploaded files, tool inputs, and webhook payloads
- reject unknown or dangerous fields where strictness is expected
- normalize types before core logic consumes them
- define size limits for payloads, files, lists, and strings where abuse is possible

**Forbidden**
- passing raw unvalidated payloads deep into domain or use-case code
- accepting arbitrary extra fields that can influence persistence, runtime, subprocess, or provider calls
- relying on runtime/provider/client-side validation alone

**Evidence**
- boundary schemas/models/parsers exist
- malformed input tests cover important rejection paths

### SEC-INJECT-001: Injection-Prone Operations Must Use Safe APIs

**Rule**
Code must use structured APIs for commands, queries, templates, and paths rather than building executable strings from untrusted data.

**Required**
- use parameterized DB queries or ORM bind parameters
- avoid shell execution where library APIs are available
- pass command arguments as arrays/lists when subprocesses are necessary
- encode output for the target context when rendering user-provided content

**Forbidden**
- string-concatenated SQL or query languages with untrusted values
- shell commands built through string interpolation
- template rendering that treats untrusted content as trusted markup
- unsafe eval-like execution of user content

**Evidence**
- injection-prone surfaces use safe APIs
- tests or review evidence cover unsafe characters and boundary cases

### SEC-SECRETS-001: Secrets Must Never Be Hardcoded, Logged, Or Exposed

**Rule**
Secrets must be loaded from approved secret/config mechanisms and must not appear in source, logs, telemetry, errors, generated artifacts, or client bundles.

**Required**
- load secrets from environment, secret manager, vault, or approved local dev mechanism
- represent secret-bearing values through explicit fields, wrappers, or types so accidental stringification is reviewable
- redact secrets before logging or returning error detail
- keep private secrets out of frontend bundles, public config, CLI diagnostics, and generated artifacts
- rotate or revoke exposed secrets if leakage is detected

**Forbidden**
- committing API keys, tokens, passwords, signing keys, or private credentials
- logging raw authorization headers, cookies, provider keys, or full connection strings
- exposing private secrets through frontend config, CLI output, generated artifacts, or diagnostics

**Evidence**
- secret-bearing config is not present in repository files
- logs/errors/diagnostics redact secret fields deterministically

### SEC-DEPS-001: Dependencies Must Be Intentional And Supply-Chain Reviewed

**Rule**
New dependencies must be justified, reviewed, pinned or constrained appropriately, and compatible with the project’s security posture.

**Required**
- explain why a dependency is needed rather than implementing locally
- check license, maintenance health, transitive dependency risk, and known vulnerabilities where practical
- pin or lock dependencies according to ecosystem norms
- avoid dependencies for trivial code unless there is clear leverage

**Forbidden**
- adding large or unmaintained packages for tiny utilities
- adding packages with unknown provenance for security-sensitive paths
- executing install-time scripts or generated code without review

**Evidence**
- lockfile/update diff is reviewed
- dependency risk is considered in review notes for meaningful additions

### SEC-FILES-001: File And Path Access Must Be Constrained

**Rule**
File system access must be scoped, normalized, and protected against traversal and unintended overwrite/read behaviour.

**Required**
- resolve and validate paths against an allowed root when user input affects paths
- prevent `..`, symlink, absolute-path, or platform-specific bypasses where relevant
- distinguish read, write, append, and delete permissions
- use private, least-permissive temporary file/directory practices for sensitive or user-controlled data

**Forbidden**
- reading or writing arbitrary paths from user-controlled input
- overwriting files without explicit intent
- exposing local file contents through errors or diagnostics

**Evidence**
- path handling is normalized and constrained
- traversal tests exist where user-provided paths are accepted
- temporary file/directory creation uses permissions appropriate to the data sensitivity

### SEC-NET-001: Outbound Network Access Must Be Intentional And Bounded

**Rule**
Outbound calls must be constrained by allowed destinations, timeouts, retries, and safe data exposure.

**Required**
- define provider/client allowlists where relevant
- set timeouts for network calls
- avoid sending secrets or user data to unapproved destinations
- preserve runtime/provider/dependency failure classification and observability

**Forbidden**
- caller-controlled URLs causing arbitrary server-side requests without validation
- unbounded outbound calls
- silent fallback to unknown runtimes, providers, or destinations

**Evidence**
- network clients are configured with timeouts and known destinations
- SSRF-like paths are explicitly controlled where dynamic URLs exist

### SEC-DESER-001: Deserialization And Dynamic Execution Must Be Restricted

**Rule**
Parsing untrusted data must not instantiate arbitrary code, types, commands, or runtime behaviour.

**Required**
- use safe parsers for JSON/YAML/XML/archive formats
- disable or avoid unsafe object deserialization features
- restrict plugin discovery and dynamic imports to trusted locations
- validate archive extraction paths and sizes

**Forbidden**
- unsafe pickle/marshal/object deserialization from untrusted sources
- eval/exec of untrusted content
- plugin loading from writable or untrusted directories without verification

**Evidence**
- parser choices are safe for untrusted input
- dynamic execution paths are explicitly justified and constrained

### SEC-DEFAULT-001: Security-Sensitive Defaults Must Fail Closed

**Rule**
When security state is unknown, missing, or misconfigured, the system must deny access or fail startup rather than silently weakening controls.

**Required**
- fail startup for missing required secrets or auth configuration
- disable debug/admin surfaces by default
- require explicit opt-in for permissive CORS, insecure TLS, test auth bypasses, and public diagnostics

**Forbidden**
- production fallback to insecure defaults
- enabling auth bypass by environment omission
- exposing debug routes in production by default

**Evidence**
- startup checks reject insecure production config
- tests cover missing security-critical config

## Review Rejection Criteria

Reject a change if it:

- trusts caller-controlled identity, role, tenant, or object ownership
- exposes protected data without object/action authorization
- builds executable queries, commands, or templates from untrusted strings
- hardcodes or logs secrets
- adds unreviewed risky dependencies
- introduces arbitrary file read/write or outbound network access
- permits unsafe deserialization or dynamic execution
- fails open on missing security configuration

---
