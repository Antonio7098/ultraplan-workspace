# Source Analysis: temporal

## Isolation Observability and Escape Handling

### Source Info

| Field | Value |
|-------|-------|
| Name | temporal |
| Path | `studies/agent-harness-study/sources/temporal` |
| Language / Stack | Go (server core), gRPC API, pluggable persistence; no agent-sandbox runtime |
| Analyzed | 2026-08-29 |

## Summary

Temporal is a workflow orchestration engine, not an LLM agent harness. There is no OS-level agent sandbox (no `chroot`, `seccomp`, `syscall` filter, container escape detector, or filesystem/network isolation for untrusted code) because untrusted code does not execute inside the server – workers execute externally and call back via gRPC. The isolation model is *logical namespace isolation*: each `Namespace` is documented as "acts as a sandbox and provides isolation for all resources within the namespace" (`service/frontend/workflow_handler.go:460-463`) with a `NamespaceId`/`Name` fence on every workflow execution, task queue, and schedule.

Observability and policy enforcement for that logical isolation is explicit and tested: a pluggable `Authorizer` interface (`common/authorization/authorizer.go:54-55`) + `ClaimMapper`, an `Interceptor` gRPC unary/stream interceptor (`common/authorization/interceptor.go:129-238`), a `DefaultAuthorizer` role model (`common/authorization/default_authorizer.go:35-65`), and metrics/cross-namespace checks. Sandbox *actions* are not audit-logged per-se; auth decisions emit structured logs and counters. Determinism violations (the closest analogue to rogue workflow code) are logged + metered (`service/frontend/workflow_handler.go:1313-1319`). There is no escape-attempt taxonomy, no syscall/filesystem audit trail, and no operator alerting pipeline beyond logs/metrics dashboards.

## Rating

**4 / 10** – Present but inconsistent, weakly documented, and fragile for the dimension as stated.

Rationale: as a workflow engine Temporal has a clear, tested authz observability path (logs + `service_errors_unauthorized`/`service_errors_authorize_failed` + latency), plus deterministic-workflow failure handling and cross-namespace authorization. As an *agent-sandbox* isolation stack the answer is absent: no OS sandbox, no audit log of file/network/syscall actions, no escape monitoring, and no alerting rule. Scoring 4 reflects the half-fit: strong logical-isolation observability, no physical-sandbox observability.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Sandbox isolation definition | `RegisterNamespace` comment: namespace "acts as a sandbox and provides isolation for all resources within the namespace. All resources belong to exactly one namespace." | `service/frontend/workflow_handler.go:460-463` |
| Sandbox auditable boundary (namespace fence) | `CallTarget` carries `APIName`, `Namespace`, `NexusEndpointName`, `Request` for every authorization decision | `common/authorization/authorizer.go:24-34` |
| Policy violation detector interface | `Authorizer` interface: `Authorize(ctx, caller *Claims, target *CallTarget) (Result, error)` | `common/authorization/authorizer.go:54-55` |
| Pluggable authorizer selection | `GetAuthorizerFromConfig` maps `""`→`NoopAuthorizer`, `"default"`→`DefaultAuthorizer`, else error; `Authorizer string yaml:"authorizer"` config key | `common/authorization/authorizer.go:64-73`, `common/config/config.go:648-649` |
| Default policy engine | `DefaultAuthorizer.Authorize`: health-check allow, `ScopeCluster`/`ScopeNamespace` role checks (`RoleReader/Writer/Admin`), returns `DecisionDeny` with generic `PermissionDenied` | `common/authorization/default_authorizer.go:35-65` |
| No-op authorizer (open by default) | `noopAuthorizer.Authorize` always allows – deployed when `TEMPORAL_AUTH_AUTHORIZER` unset | `common/authorization/noop_authorizer.go:12-14` |
| Env/config wiring for authorizer | `authorizer: {{ env "TEMPORAL_AUTH_AUTHORIZER" }}` in docker.yaml and embedded template; `TEMPORAL_AUTH_AUTHORIZER=default` enables default authorizer | `config/docker.yaml:264`, `common/config/config_template_embedded.yaml:300`, `config/development-jwt.yaml:67` |
| Unary auth interceptor (policy enforcement point) | `Interceptor.Intercept` extracts TLS/auth header, maps claims, builds `CallTarget{Namespace,APIName,Request}`, calls `Authorize`, invokes `authorizeTargetNamespaces` | `common/authorization/interceptor.go:129-185` |
| Streaming auth interceptor | `InterceptStream` optionally skips via `disableStreamingAuthorizer`, same Authorize path, strips principal | `common/authorization/interceptor.go:188-238` |
| Principal-spoof prevention (escape hardening) | `headers.StripPrincipal(ctx)` on every unary and stream path "to prevent external callers from spoofing principal identity" | `common/authorization/interceptor.go:156-158`, `common/authorization/interceptor.go:219-221` |
| Dynamic config gating | `DisableStreamingAuthorizer`, `ExposeAuthorizerErrors`, `EnableCrossNamespaceCommands` dynamic bools control auth strictness | `common/dynamicconfig/constants.go:149-152`, `common/dynamicconfig/constants.go:934-938` |
| Auth decision logging (observability) | `logger.Error("Authorization error", tag.Error(err))` on GetClaims failure and on `Authorizer.Authorize` error | `common/authorization/interceptor.go:149`, `common/authorization/interceptor.go:213`, `common/authorization/interceptor.go:231`, `common/authorization/interceptor.go:316` |
| Authz metrics – failure counters | `ServiceErrUnauthorizedCounter` and `ServiceErrAuthorizeFailedCounter` incremented on deny / error; latency via `ServiceAuthorizationLatency` | `common/authorization/interceptor.go:315-316`, `common/authorization/interceptor.go:323`, `common/metrics/metric_defs.go:693-700` |
| Authz metric cardinality guard | `getMetricsHandler` only tags real namespace (`Exists()==nil`) to avoid unbounded cardinality, else `NamespaceUnknownTag` | `common/authorization/interceptor.go:334-345` |
| Cross-namespace escape detection | `authorizeTargetNamespaces` iterates `RespondWorkflowTaskCompleted` commands (`SignalExternalWorkflow`, `StartChildWorkflow`, `RequestCancelExternalWorkflow`), re-authorizes each `targetNamespace:api` once, dedup via `authorizedNamespaceAPIs` map | `common/authorization/interceptor.go:347-417` |
| Non-determinism / workflow-task escape analogue (detection+logging) | `RespondWorkflowTaskFailed` logs "Non-Deterministic Error" with `WorkflowNamespaceID/WorkflowID/RunId` and increments `ServiceErrNonDeterministicCounter` | `service/frontend/workflow_handler.go:1313-1320`, `common/metrics/metric_defs.go:693` |
| Workflow task failure-cause taxonomy | `workflowTaskFailedCause` struct + `newWorkflowTaskFailedCause(failedCause,causeErr,terminateWorkflow)` centralized failure classification | `service/history/api/respondworkflowtaskcompleted/workflow_task_completed_handler.go:90-94`, `service/history/api/respondworkflowtaskcompleted/workflow_task_completed_handler.go:1603-1612` |
| Authorization test coverage (observable contract) | `authorizerInterceptorSuite` tests: `TestIsAuthorized`, `TestIsUnauthorized` (expects `ServiceErrUnauthorizedCounter`), `TestAuthorizationFailed`, `TestCrossNamespaceCommands_Authorized/Unauthorized`, `TestMultipleCommands_AuthDeduplication` | `common/authorization/interceptor_test.go:105-126`, `common/authorization/interceptor_test.go:340-431`, `common/authorization/interceptor_test.go:497-539` |
| Nexus auth observability | Nexus handler tags `outcome=unauthorized`, logs "Authorization internal error with processing nexus request/callback" | `service/frontend/nexus_handler.go:188-192`, `service/frontend/nexus_completion_http_handler.go:605-609` |
| No syscall/container sandbox evidence | Grep for `sandbox`, `seccomp`, `apparmor`, `capability`, `chroot`, `syscall` filter yields only one `sandbox` hit (namespace comment) and `syscall` limited to TCP retry and signal handling | `service/frontend/workflow_handler.go:462`, `common/persistence/sql/sqlplugin/db_handle.go:11,162`, `temporal/interrupt.go:6-11` |
| No audit log of sandbox actions | Grep for `audit` yields one hit "audit-only" cluster-cache lookup; no structured audit log sink | `service/frontend/admin_handler.go:1001` |
| No operator alerting rule for isolation breach | Alert/slack/pagerduty only in CI notifiers (`tools/ci-notify`, `tools/testrunner` alerts suite for data-race/panic), not runtime isolation alerts | `tools/ci-notify/slack.go:129-147`, `tools/testrunner/log.go:76-146`, `tools/testrunner/junit.go:16-18` |

## Answers to Dimension Questions

### 1. Are sandbox actions observable?

**Partially, for the logical sandbox only.** Every gRPC API call is interceptable via `common/authorization/interceptor.go:129-185` (unary) and `188-238` (streaming). The `CallTarget` (`common/authorization/authorizer.go:24-34`) captures `APIName`+`Namespace`+`Request`. Decisions are observable through (a) structured `logger.Error("Authorization error", tag.Error(err))` (`interceptor.go:316`) transcending to aggregated logs, (b) counters `service_errors_unauthorized` and `service_errors_authorize_failed` plus `service_authorization_latency` (`common/metrics/metric_defs.go:694-700`, `interceptor.go:313-323`), and (c) history persistence of `WorkflowTaskFailed` events with `WORKFLOW_TASK_FAILED_CAUSE_*` taxonomy and non-determinism logging (`service/frontend/workflow_handler.go:1313-1319`). There is no filesystem/syscall/network audit trail for code running inside a worker, and no dedicated audit log stream (grep `audit` hits 1 audit-only comment in `service/frontend/admin_handler.go:1001`). Namespace-scoped reads/writes (StartWorkflow, Signal, etc.) are metered by `OperationTag(AuthorizationScope)` (`interceptor.go:344`). Raw operator question "which workflow touched which resource when?" must be reconstructed from history events and log aggregation, not a first-class audit sink.

### 2. Are policy violations detected?

**Yes, synchronously on every request, if an authorizer is configured.** `Interceptor.Authorize` (`common/authorization/interceptor.go:304-331`) invokes the pluggable `Authorizer.Authorize` (`authorizer.go:54-55`) and branches on `err != nil` → `ServiceErrAuthorizeFailedCounter` and `result.Decision != Allow` → `ServiceErrUnauthorizedCounter`, returning `PermissionDenied(RequestUnauthorized, result.Reason)` (`interceptor.go:322-328`). `DefaultAuthorizer` (`common/authorization/default_authorizer.go:35-65`) enforces `ScopeNamespace` vs `ScopeCluster` + `AccessReadOnly/Write/Admin` vs `RoleReader/Writer/Admin` (`common/api/metadata.go:70-194`) and always denies unknown scope (`default_authorizer.go:55-56`). Cross-namespace commands are re-checked (`interceptor.go:347-417`). Enforcement is bypassable: an empty `authorizer:""` yields `NoopAuthorizer` (`authorizer.go:67-68`, `noop_authorizer.go:12`) which is the docker/dev default (`config/docker.yaml:264` with env unset); `DisableStreamingAuthorizer` skips streaming auth entirely (`interceptor.go:195`). No content-level policy (e.g., "workflow must not call external URL") is evaluated – only API/namespace/role triple.

### 3. Are escape attempts handled?

**Logical-namespace escape is explicitly handled; OS/container escape is not applicable.** Namespace is hard-partitioned: a workflow cannot reference resources outside its namespace without an explicit cross-namespace command; `authorizeTargetNamespaces` (`interceptor.go:347-417`) detects and blocks such attempts by re-authorizing the target namespace API (`SignalWorkflowExecution`/`StartWorkflowExecution`/`RequestCancelWorkflowExecution`) and dedups per `(namespace:api)` – unauthorized attempts return `PermissionDenied` and increment the unauthorized counter. Principal spoofing is hardened via `headers.StripPrincipal` on ingress (`interceptor.go:156`, `220`). Workflow determinism violations (rogue workflow code) are captured as `WORKFLOW_TASK_FAILED_CAUSE_NON_DETERMINISTIC_ERROR` with task failure, history event, and metric (`workflow_handler.go:1313-1319`), but the sandbox does not terminate or quarantine the worker; it simply fails the workflow task. There is no detection for container breakout, network egress, or file exfiltration because the server does not sandbox worker processes at all – workers are trusted external processes.

### 4. Are violations alerted to operators?

**Logged and metered, but not alerted.** The system "knows" via error logs (`interceptor.go:149,213,231,316`) and Prometheus-style counters/histograms (`metric_defs.go:693-700`) which dashboards and generic log-based alerts can consume, but there is no built-in isolation-specific alerting rule, PagerDuty/Slack hook for authz failures, or severity escalation in the runtime path. The only alerting paths in repo are CI-scoped: `tools/testrunner` synthetic `ALERTS` JUnit suite (`junit.go:16-18`) and `tools/ci-notify` Slack webhook for CI failures (`ci-notify/slack.go:129`), not production isolation breaches. `ExposeAuthorizerErrors` (`dynamicconfig/constants.go:934`) is off by default, so operator-facing logs omit authorizer detail while metrics still fire. Compared to an agent harness, vertical alerting (detect → classify → page) is absent; operators must wire their own log/metric alerts.

## Architectural Decisions

* **Logical sandbox = Namespace, not process/container.** `service/frontend/workflow_handler.go:460` declares namespace as sandbox. Tradeoff: multi-tenant isolation without running untrusted code in-server, but no defense against a compromised worker host accessing another host's data beyond authz.
* **Pluggable ClaimMapper+Authorizer pipeline, default-open.** `common/authorization/authorizer.go:64`, `common/config/config.go:648`, `temporal/fx.go:434` allow custom JWT/TLS claim mapping and allow-list logic; shipped defaults are `noop` for local development, `default` only when `TEMPORAL_AUTH_AUTHORIZER=default` (`config/development-jwt.yaml:67`). Decision favors ease of local run over secure-by-default for isolation experiments.
* **Unary + streaming interceptors as single enforcement point.** `interceptor.go:129`, `188` centralize all authz. Positive: one path to log/meter; negative: streaming namespace unavailable at handshake (`interceptor.go:224-228`) so cross-namespace checks are limited to `RespondWorkflowTaskCompleted`.
* **Error-or-deny both map to `service_errors_*` counters with namespace-aware cardinality guard.** `interceptor.go:334-345` prevents tag cardinality blowup by checking `NamespaceRegistry.Exists` before tagging. This enables per-namespace SLOs without metric-store explosion.
* **Workflow failure taxonomy as code-isolation signal.** `workflow_task_completed_handler.go:90-94,1603` classifies determinism/size/policy failures via `WORKFLOW_TASK_FAILED_CAUSE_*` enum, surfaced in history and via `ServiceErrNonDeterministicCounter`. Chosen over killing the worker; retains audit trail in durable history.

## Notable Patterns

* **Interceptor-as-policy-gate pattern.** Single `Authorize` call followed by `authorizeTargetNamespaces` loop with deduplication map – test-dedup assertion in `interceptor_test.go:497-539` proves intent.
* **Strip-principal anti-spoofing.** `headers.StripPrincipal(ctx)` unconditional even when authorizer nil, prevents header-injection escape (CWE-565 analogue).
* **Dynamic gating for gradual rollout.** `DisableStreamingAuthorizer`, `ExposeAuthorizerErrors`, `EnableCrossNamespaceCommands` as `BoolPropertyFn` allow safe rollback without redeploy.
* **History-as-audit-log pattern.** Rather than a separate audit sink, isolation-relevant events (task failures, start/signal) are durably recorded in workflow history tables; `RespondWorkflowTaskFailed` path is the audit entry.
* **No runtime sandbox; workers are external.** Invariant throughout `service/` – server never executes user code directly, only validates commands submitted by workers.

## Tradeoffs

| Decision | Benefit | Cost / Risk |
|----------|---------|-------------|
| Namespace-as-sandbox (logical) vs OS sandbox | Scales multi-tenant without container overhead; persists isolation in DB | No protection against malicious worker binary reading host filesystem or exfiltrating data; relies on network authz not syscall filtering |
| Pluggable authorizer, default `noop` | Easy dev experience, customizable to OPA/JWT | Production misconfiguration (forget `authorizer: default`) leaves isolation open with no warning beyond missing metrics |
| Single interceptor + metrics | Consistent observability point | Streaming calls lack namespace at auth time, cross-namespace cmds outside `RespondWorkflowTaskCompleted` invisible; log volume proportional to unauthorized request rate (potential DoS amplification) |
| `ExposeAuthorizerErrors=false` default | Avoids information disclosure to untrusted caller | Operator loses root-cause in logs unless dynamic flag flipped; must correlate metrics+separate audit query |
| Tag only known namespaces | Protects metrics store | Unknown/malicious namespace probes bucketed as `NamespaceUnknown`, losing attribution on attempted namespace enumeration |
| History durability vs real-time alert | Audit trail survives crashes | No synchronous page on breach; relies on external monitoring to poll metrics/logs |

## Failure Modes / Edge Cases

* **Misconfigured authorizer → silent open.** `GetAuthorizerFromConfig` returns `NoopAuthorizer` for empty string (`authorizer.go:67-68`). Deployment with unset `TEMPORAL_AUTH_AUTHORIZER` authenticates but authorizes everything, logging no deny events – indistinguishable from healthy idle cluster without external canary requests.
* **Streaming bypass.** If `DynamicConfig.DisableStreamingAuthorizer=true` (`interceptor.go:195`), all streaming RPCs bypass auth; viable during SDK decode but expands escape surface.
* **Cross-namespace gap outside workflow-task-completed.** Direct `SignalWorkflowExecution` cross-namespace would be denied by per-call `CallTarget.Namespace` check, but batched helpers via other APIs or future commands could evade `authorizeTargetNamespaces` which only scans `RespondWorkflowTaskCompleted` (`interceptor.go:365-366`).
* **Metric tag spoofing.** Attacker probing random namespaces until one exists triggers per-namespace tag emission for hits; could still cause tag cardinality growth up to number of real namespaces.
* **Determinism escape survives retry storm.** Non-determinism increments one counter per failure (`workflow_handler.go:1319`) but workflow task will be retried immediately, potentially hot-looping server and flooding logs/metrics without circuit breaker.
* **No syscall audit → blind to worker compromise.** Compromised worker can issue only APIs allowed by its role but can abuse data-plane parallelization (list/scan, batch operations) within those APIs – no network/file egress detector flags bulk export as anomalous.
* **Audit gap.** There is no tamper-evident authz audit log; `logger.Error` is best-effort and can be dropped under backpressure (`common/log/throttle_logger.go:36-54`).

## Future Considerations

* Add a structured authz audit sink (e.g., emit `CallTarget`+`Claims`+`Decision`+`Reason` to a dedicated audit logger or history-archival bucket), covered by integration test – closes current reliance on unstructured error logs.
* Promote namespace sandbox documentation from comment (`workflow_handler.go:462`) to explicit ADL with tested invariant: resources cannot reference cross-namespace IDs without a validated `CallTarget` – add negative test probing direct store access.
* Enable secure-by-default wiring: fail start if `Global.Authorization.Authorizer` unset in production profile, or log `logger.Fatal` equivalent warning with a metric on noop path.
* Extend `authorizeTargetNamespaces` to a generic cross-namespace hook or to evaluate ahead-of-time for all batch/signal APIs, removing command-scoped blindspot.
* If an agent harness were layered on Temporal, add OS-level sandbox sidecar around workers (gVisor/Firecracker) and wire its denials to `service_errors_unauthorized`-adjacent alert with paging – server alone cannot answer syscall-escape questions.
* Add alerting recipe artifact in `config/docker.yaml`/helm: Prometheus alert on `rate(service_errors_unauthorized[5m])` and `rate(service_errors_nondeterministic[5m])` thresholds.

## Questions / Gaps

* No evidence in repo that namespace isolation has a formal threat model document; searched `docs/` and found only brief sandbox comment. Ask maintainers for design doc beyond the comment.
* No tracing of successful authorized calls – only errors metered. Does a successful cross-namespace allow generate an auditable event? Not found – would need code walk of `history` service.
* Unknown whether cloud SaaS layer adds an external audit/escape alert not present in OSS repo – sibling sources policy prevents cross-source check; report as gap.
* Bulk-data exfiltration policy (rate limit, payload size) appears in `history/service` admission control but no linkage to isolation violation detector – search boundary was `authorization/`, `namespace/`, `handler` layers only.

---
Generated by `Dimension 17.03: Isolation Observability and Escape Handling` against `temporal`.
