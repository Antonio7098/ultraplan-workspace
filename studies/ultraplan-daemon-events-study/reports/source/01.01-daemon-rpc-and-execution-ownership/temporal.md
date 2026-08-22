# Source Analysis: temporal

## Daemon RPC and Execution Ownership

### Source Info

| Field | Value |
|-------|-------|
| Name | temporal |
| Path | `studies/ultraplan-daemon-events-study/sources/temporal` |
| Language / Stack | Go 1.26 (`go.mod:3`), gRPC + protobuf, Uber `fx` DI, ringpop membership, SQL/Cassandra/Elasticsearch persistence |
| Analyzed | 2026-08-22 |

## Summary

Temporal is a set of cooperating daemons (frontend, history, matching, worker services, plus an internal frontend) rather than a single CLI-owned process; one binary can host all services (`temporal/server.go:25-40`, `cmd/server/main.go:136-142`). The core ownership model is that **accepted work never depends on the submitting connection**: once `StartWorkflowExecution` succeeds, the workflow's mutable state lives in a history-service shard backed by the execution store, and the frontend merely proxies into it (`service/frontend/workflow_handler.go:548-584`). Clients — CLIs, SDKs, web UI, workers — are stateless peers that discover state through stateless gRPC calls: long-polls against task queues (`service/matching/handler.go:225-272`) and history reads, with transport-level auto-reconnect handled by the gRPC channel (`common/rpc/grpc.go:64-73`).

The wire surface is protobuf-defined WorkflowService/AdminService/OperatorService served over a TCP gRPC listener (`common/rpc/rpc.go:196-211`), optionally with an HTTP/JSON gateway on a separate port (`service/frontend/fx.go:1046-1079`). A deep interceptor chain enforces namespace validation, JWT/TLS-based authorization, rate limits, and version checks before any handler runs (`service/frontend/fx.go:286-321`). "Workspace context" maps to **namespaces**: every request carries a namespace name or a task token bound to a namespace, validated up front (`common/rpc/interceptor/namespace_validator.go:42-55`).

Shutdown is deliberate and observable: SIGINT/SIGTERM triggers reverse-order service stop (`temporal/interrupt.go:9-21`, `temporal/server_impl.go:109-124`); each frontend first fails its health check and flags itself draining in ringpop so load balancers stop routing, then drains traffic for a configurable window before a forced stop (`service/frontend/service.go:544-590`, `common/membership/interfaces.go:59-60`).

## Rating

**8 / 10** — Clear, well-factored ownership model with explicit interfaces, tests, and operational safeguards (graceful drain, membership draining, health-check-driven LB removal, shard-ownership transfer protocol). Deductions: the default authorization posture is a noop authorizer that merely warns without `--allow-no-auth` (`cmd/server/main.go:203-209`); there is no Unix-socket/named-pipe local transport (TCP-only listeners, `common/rpc/rpc.go:200-211`); legacy config loading is working-directory-relative unless overridden (`cmd/server/main.go:41-53,167-178`); and there is no generic RPC-level session re-attach — clients recover operations through domain semantics (workflow IDs, request IDs) rather than a reconnect protocol.

## Evidence Collected

Every entry cites `path/to/file.go:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Daemon entrypoint | `main()` builds a urfave/cli app whose `start` action loads config, constructs authorizer/claim-mapper/audience-mapper, and calls `NewServer(...).Start()` | `cmd/server/main.go:28-31`, `cmd/server/main.go:154-244` |
| Server abstraction | `Server` interface is just `Start()/Stop()`; composable list of services incl. InternalFrontend | `temporal/server.go:16-20`, `temporal/server.go:25-40` |
| Service startup ordering | Fixed init order (matching → history → internal/frontend → worker); stopped in reverse | `temporal/server_impl.go:47-53`, `temporal/server_impl.go:109-146` |
| Signal-based shutdown | `InterruptCh()` converts SIGINT/SIGTERM into a stop channel wired via `InterruptOn` | `temporal/interrupt.go:9-21`, `cmd/server/main.go:226` |
| fx lifecycle wiring | TopLevelModule invokes `ServerLifetimeHooks`; `fx.StartStopHook(svr.Start, svr.Stop)` | `temporal/fx.go:140-159`, `temporal/fx.go:944-947` |
| Transport (server side) | Plain TCP gRPC listener per service; bind address from `bindOnLocalHost`/`bindOnIP`/default private IP | `common/rpc/rpc.go:196-237` |
| Transport (client side) | `Dial()` uses `grpc.NewClient` with round-robin DNS LB config, TCP keepalive, and capped reconnect backoff (10s max) | `common/rpc/grpc.go:21-44`, `common/rpc/grpc.go:50-112` |
| Internode connection cache | Cached connections redialed if in `Shutdown` state; periodic sweep of dead conns | `common/rpc/rpc.go:317-351`, `common/rpc/rpc.go:354-371` |
| TLS | Separate frontend vs internode TLS server options from a TLS config provider | `common/rpc/rpc.go:131-146`, `common/rpc/rpc.go:164-185` |
| Remote-cluster auth | Outbound cross-cluster dials attach per-RPC token credentials when a tokenProvider is configured | `common/rpc/rpc.go:282-299` |
| gRPC registration & serve | Frontend registers WorkflowService, AdminService, OperatorService, health, reflection; serves on listener; starts membership monitor | `service/frontend/service.go:504-541` |
| Graceful shutdown drain | Health check failed → ringpop `SetDraining(true)` → wait failure-detection time → `GracefulStop` with drain-time deadline → hard `Stop` | `service/frontend/service.go:543-590`, `common/membership/interfaces.go:50-60` |
| Interceptor chain | Ordered unary interceptors: error masking, service errors, business-ID extraction, namespace validation, auth, handover, redirection, telemetry, rate limits, SDK-version, retry innermost | `service/frontend/fx.go:286-321` |
| Authorization model | Auth info from `authorization` header (JWT) + TLS peer cert subject; claims mapped; per-API/per-namespace authorizer; principal header spoofing stripped | `common/authorization/interceptor.go:129-185`, `common/authorization/interceptor.go:251-285` |
| Streaming auth | Dedicated stream interceptor mirrors unary auth; can be disabled via dynamicconfig | `common/authorization/interceptor.go:187-238` |
| Cross-namespace commands | RespondWorkflowTaskCompleted commands targeting other namespaces re-authorized per target namespace | `common/authorization/interceptor.go:351-417` |
| Namespace scoping | Requests must carry exactly one of namespace name/ID; allowed states per API; task-token↔namespace enforcement; read blocklist during handover | `common/rpc/interceptor/namespace_validator.go:34-90` |
| Workspace propagation into backend | Frontend resolves namespace name→ID and stamps it onto history-service requests | `service/frontend/workflow_handler.go:561-579` |
| Version negotiation (client→server) | SDK name/version recorded; `ClientSupported` rejects unsupported clients; server-version range header rejected via `ServerVersionNotSupported` | `common/rpc/interceptor/sdk_version.go:31-45`, `common/headers/version_checker.go:115-144` |
| Version negotiation (server→client) | `GetSystemInfo` returns `ServerVersion` + hardcoded capability booleans | `service/frontend/workflow_handler.go:3498-3526` |
| Execution ownership transfer | Frontend hands off accepted work to history service immediately (`wh.historyClient.StartWorkflowExecution`) | `service/frontend/workflow_handler.go:570-579` |
| Shard-routed clients | History clients route ops by shardID with cached per-owner connections; `ShardOwnershipLost` triggers cache invalidation/retry elsewhere | `client/history/caching_redirector.go:82-105`, `client/history/caching_redirector.go:205-225` |
| Worker attach/detach | Workers long-poll `PollWorkflowTaskQueue`/`PollActivityTaskQueue`; handler asserts a long-poll timeout is set | `service/matching/handler.go:225-272` |
| Long-poll budgets | Frontend default 3 min, matching default 5 min, shared constant 60 s for other polls | `client/frontend/client.go:18`, `client/matching/client.go:33`, `common/constants.go:35` |
| Long-poll guard on reads | `ValidateLongPollContextTimeout` applied to history/event polling APIs | `service/frontend/workflow_handler.go:1079`, `service/frontend/workflow_handler.go:1349` |
| Idempotency hooks | `RequestId` validated at the edge; CHASM workflow uses `requestID` for idempotent callback IDs; duplicate workflow-ID resolution policy | `service/frontend/workflow_handler.go:661`, `chasm/lib/workflow/workflow.go:176-184`, `chasm/lib/workflow/workflow.go:441` |
| HTTP API surface | Optional HTTP API server on its own port; disabled when port unset/0 | `service/frontend/fx.go:1037-1044`, `service/frontend/fx.go:1046-1079` |
| Config / CWD coupling | Legacy load joins `root` (default `"."`) with config dir; `--config-file` documented as relative to CWD; embedded template loads env-only | `cmd/server/main.go:41-53`, `cmd/server/main.go:70`, `cmd/server/main.go:167-178`, `common/config/loader.go:111-117`, `common/config/loader.go:143-157` |
| No-authorizer posture | Warns when running with noop authorizer unless `--allow-no-auth` passed | `cmd/server/main.go:197-209` |

## Answers to Dimension Questions

**1. Does a client exit affect execution?**
No. Ownership transfers at acceptance: the frontend resolves the namespace and forwards the start request into the history service, which owns the workflow's event history and mutable state in a shard-backed store (`service/frontend/workflow_handler.go:561-579`). Task dispatch is pull-based: workers long-poll matching service queues with their own deadlines (`service/matching/handler.go:225-272`), so a disconnected worker simply stops polling and its tasks remain queued until the poll timeout lapses. Nothing in the execution path holds a reference to a client connection; the only client-coupled resources are the long-poll RPC contexts themselves.

**2. What exactly crosses the daemon boundary?**
Protobuf request/response messages over gRPC (WorkflowService, AdminService, OperatorService registered at `service/frontend/service.go:507-512`), an optional HTTP/JSON projection of selected APIs (`service/frontend/http_api_server.go`, enabled per `service/frontend/fx.go:1037-1044`), opaque serialized **task tokens** that carry namespace/task identity back on completion (`common/rpc/interceptor/namespace_validator.go:20-22`, enforced at `common/rpc/interceptor/namespace_validator.go:39`), and gRPC metadata headers for caller identity/versioning (SDK name/version and supported-server-versions ranges, `common/headers/version_checker.go:99-121`). Nexus adds an HTTP task/completion protocol bounded by `MaxNexusAPIRequestBodyBytes` (`common/rpc/grpc.go:34-37`).

**3. Can a client reconnect and recover the same operation?**
At the transport layer, transparently: gRPC channels auto-reconnect with exponential backoff capped at 10 s (`common/rpc/grpc.go:64-73`), round-robin across DNS-resolved frontends (`common/rpc/grpc.go:22-25`), and internode caches redial after `Shutdown` (`common/rpc/rpc.go:320-326`). At the operation layer, recovery is semantic rather than session-based: a new client reaches the same workflow by Workflow ID via describe/history APIs; duplicate `StartWorkflowExecution` collapses onto the existing execution via workflow-ID reuse policy (`chasm/lib/workflow/workflow.go:441`) and callers supply a `RequestId` for idempotent behavior (`service/frontend/workflow_handler.go:661`, `chasm/lib/workflow/workflow.go:176-184`). There is **no** generic "reattach to in-flight RPC result" mechanism; an update whose response is lost must be resolved by polling the update handle — no evidence of a server-side session registry was found (see Gaps).

**4. How is the local caller authenticated?**
Peer authentication layers, all off the request path as interceptors: (a) optional mutual TLS with the client certificate subject fed into claim mapping (`common/authorization/interceptor.go:57-81`, `common/authorization/interceptor.go:263-266`); (b) JWT bearer token in a configurable header mapped to permission claims by a pluggable claim mapper, then evaluated per API + namespace by a pluggable authorizer (`common/authorization/interceptor.go:129-185`, wired at `service/frontend/fx.go:191-216`); (c) inbound principal headers are always stripped to prevent spoofing even when authorization is disabled (`common/authorization/interceptor.go:156-158`). Internode calls use a distinct TLS profile (`common/rpc/rpc.go:164-185`), and cross-cluster replication calls authenticate with short-lived tokens fetched per address (`common/rpc/rpc.go:282-299`). Caveat: out of the box the authorizer is a noop and the binary only warns unless `--allow-no-auth` is set (`cmd/server/main.go:203-209`).

**5. Does daemon code ever depend on its own current working directory?**
Only at configuration load time, and only in non-default paths. The deprecated flag set defaults `root` to `"."` and joins it with the config dir (`cmd/server/main.go:41-53`, `cmd/server/main.go:170-175`); `--config-file` explicitly allows paths relative to the current working directory (`cmd/server/main.go:68-72`). The default `start` path avoids the filesystem entirely by loading an embedded template driven by environment variables (`cmd/server/main.go:177`, `common/config/loader.go:111-117,145-148`). No runtime code was observed resolving state from the process CWD: listeners bind configured/explicitly-selected IPs (`common/rpc/rpc.go:213-237`), not socket paths.

**Kill test** — kill the submitting client right after acceptance: the history shard owning the new execution continues driving it (timers, tasks, retries are server-side state machines), the matching service queues its tasks independently, and another client finds the operation deterministically by `(namespace, WorkflowId[, RunId])` through Describe/History APIs, while workers find its tasks by long-polling the named task queue. Discovery requires no knowledge of which frontend instance accepted the call because no per-client routing state exists.

## Architectural Decisions

1. **Ownership lives in durable state machines, not connections.** Acceptance = write into history-shard-owned state; everything after is server-driven (`service/frontend/workflow_handler.go:570-579`, shard routing via `client/history/caching_redirector.go:94-105`).
2. **Stateless multi-frontend behind ringpop membership.** Any frontend can answer for any namespace; membership monitor publishes liveness/draining used both by internal resolvers and shutdown logic (`common/membership/interfaces.go:50-60`, `service/frontend/service.go:540,557`).
3. **Interceptor-chain security envelope.** Ordering is documented in-code: masking outermost, retry innermost, namespace validation before auth-aware handlers (`service/frontend/fx.go:287-321`).
4. **Pull-based task delivery.** Long polls with mandatory client-set deadlines replace push/stream sessions, making worker detach a non-event (`service/matching/handler.go:242-249`).
5. **Two-plane versioning.** Client→server minimum-version rejection plus server→client capability discovery, both additive and backward-compatible (`common/headers/version_checker.go:124-144`, `service/frontend/workflow_handler.go:3506-3525`).
6. **fx-scoped service composition.** Each service is an fx Module with start/stop lifecycle hooks, enabling single-process dev deployments and independent production scaling alike (`service/frontend/fx.go:78-150`, `service/frontend/fx.go:1098-1100`).

## Notable Patterns

- **Drain-before-stop**: fail health → mark draining → sleep failure-detection window → GracefulStop with deadline timer → hard stop (`service/frontend/service.go:545-584`).
- **Shard-ownership redirect loop**: clients cache per-owner history connections; `ShardOwnershipLost` evicts and retries against the new owner (`client/history/caching_redirector.go:105,205-225`) — a clean ownership-handover contract between replicas.
- **Capability booleans as protocol evolution**: older servers omit fields ⇒ implied false, so capability probing needs no handshake round-trip (`service/frontend/workflow_handler.go:3508-3510`).
- **Defense-in-depth header hygiene**: unconditional stripping of spoofable principal headers regardless of authorizer presence (`common/authorization/interceptor.go:156-158`).
- **Namespace-state gates per API**: e.g., only `REGISTERED` namespaces accept new starts; deletes allowed in more states (`common/rpc/interceptor/namespace_validator.go:42-55`).

## Tradeoffs

- **TCP-only transport**: simplifies deployment (no socket permissions) but means even localhost clients pay TCP+TLS setup and cannot leverage OS-authenticated sockets; no Unix-domain-socket option exists in the listener code (`common/rpc/rpc.go:200-211`).
- **Semantic (not session) recovery**: extremely robust for workflows, weaker for interactive RPCs — a lost response to an update/signal must be recovered by application-level querying; the server keeps no per-client operation ledger.
- **Optional authorization by default**: maximal compatibility, but a misconfigured deployment silently runs open; only a warning mitigates it (`cmd/server/main.go:203-209`).
- **Hardcoded capabilities**: simple and cheap, but ties capability truth to compile-time constants rather than feature-flagged runtime introspection (`service/frontend/workflow_handler.go:3511-3524`).
- **Long-poll instead of streaming task delivery**: fewer stateful resources per worker and trivial detach, at the cost of delivery latency bounded by poll intervals and per-request rate-limit interactions (poll APIs get dedicated limiter treatment, `service/frontend/fx.go:640-646`).

## Failure Modes / Edge Cases

- **Frontend crash mid-RPC**: client channel reconnects to a peer via round-robin/backoff (`common/rpc/grpc.go:64-73`); accepted-but-unacknowledged starts are still recoverable because the workflow exists durably under its Workflow ID.
- **History shard ownership moves**: in-flight internode calls receive `ShardOwnershipLost`; redirector invalidates the cached owner and retries (`client/history/caching_redirector.go:205-225`).
- **Rolling restart**: draining frontends fail health checks first so LBs and ringpop stop routing before sockets close; a deadline forces closure if drain stalls (`service/frontend/service.go:552-578`).
- **Namespace handover (multi-cluster)**: data-read APIs are blocked during handover to prevent read-after-write inconsistency, with an explicit allowlist (`common/rpc/interceptor/namespace_validator.go:57-86`).
- **Slow/dead workers**: long-poll deadlines (5-min default) return empty responses rather than pinning tasks to vanished workers (`client/matching/client.go:33`, timeout asserted at `service/matching/handler.go:242-249`).
- **Legacy CWD-relative config**: launching the binary from an unexpected directory with deprecated `--root/--config` flags silently loads different config; embedded/env mode eliminates this class of failure (`cmd/server/main.go:167-178`).

## Future Considerations

- Ship a first-class authenticated local transport (Unix socket / named pipe) so co-located CLIs and tools avoid TCP+TLS overhead while inheriting OS credential binding — would strengthen the "local RPC" story this dimension probes.
- Promote the noop-authorizer warning to a hard requirement (the flag text already foreshadows this: "Future versions will require...", `cmd/server/main.go:204-208`).
- Expose runtime-derived capabilities (intersected with dynamic-config feature flags) through `GetSystemInfo` instead of static constants.
- Document/centralize the operation-recovery story (request-ID dedup windows, update-handle polling) as part of the public wire contract, since clients currently rely on implicit semantics spread across handlers.

## Questions / Gaps

- **No session/attach registry found.** Searches for client-session or RPC-resume machinery (`grep -ri "attach\|resume" common/rpc client`) surfaced nothing beyond gRPC-channel reconnection; conclusion: recovery is purely semantic. Confidence high for this repo, though some behavior (e.g., update-result buffering) lives behind the SDK boundary outside this source tree.
- **Matching internals sampled, not exhausted.** Task-queue partition/durability behavior inside `service/matching` was cited only at the handler boundary; a deeper pass could confirm how queued tasks survive matching-host loss.
- **HTTP API parity unclear.** The HTTP gateway exists (`service/frontend/http_api_server.go`) but which RPCs it exposes and whether long-polls are supported over HTTP was not traced end-to-end here.
- **InternalFrontend exposure.** `primitives.InternalFrontendService` appears in the service list and interceptor providers (`temporal/server.go:27`, `service/frontend/fx.go:278-279,549-554`), but its network exposure rules versus public frontend were not fully audited in this pass.

---

Generated by `dimension 01.01-daemon-rpc-and-execution-ownership` against `temporal`.
