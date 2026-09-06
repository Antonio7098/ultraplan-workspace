# Source Analysis: nomad

## 01.13 Daemon, Local API, and Multi-Client Lifetime

### Source Info

| Field | Value |
|-------|-------|
| Name | nomad |
| Path | `studies/aren-go-runtime-study/sources/nomad` |
| Language / Stack | Go — `nomad` server + `client` + HTTP/mTLS + Raft/Serf + yamux RPC |
| Analyzed | 2026-08-30 |

## Summary

Nomad implements a cluster-oriented long-lived daemon, not a single-host loopback daemon. One `nomad agent` process (`command/agent/agent.go:72`) optionally embeds both a Raft-replicated server (`nomad/server.go`) and a client (`client/client.go`) under a shared `Agent` struct (`command/agent/agent.go:73-146`). Accepted work is a `Job` (not an ephemeral run) that is Raft-committed via the FSM (`nomad/fsm.go:1036`, `nomad/job_endpoint.go:990`) and durably stored in the `state.StateStore`. Scheduling is server-side (`nomad/worker.go:363`, `nomad/plan_apply.go`) and allocations persist independently of any submitting CLI process. Any HTTP client holding a valid ACL token can observe the same job/allocation via blocking queries (`nomad/job_endpoint.go:1498`, `nomad/alloc_endpoint.go:58`) or the ndjson event stream (`command/agent/event_endpoint.go:46`, `nomad/event_endpoint.go:32`) without sharing process memory. The wire contract is HTTP/1.1 + JSON (`command/agent/http.go:101-227` `NewHTTPServers`), with an internal msgpack streaming RPC tunneled over HTTP hijack (`command/agent/event_endpoint.go:92-200`), yamux-multiplexed Serf/Raft RPC (`command/agent/agent.go:564`, `command/agent/config.go:854`), and optional mTLS (`command/agent/http.go:143-173`, `helper/tlsutil`). Workspace scoping is via `region` + `namespace` query parameters (`command/agent/http.go:978-994`), not a filesystem workspace. Graceful leave/drain is cooperative: `Agent.Leave()` → `server.Leave()`/`client.Leave()` (`command/agent/agent.go:1493`, `nomad/server.go:814`), `Shutdown()` sequence (`command/agent/agent.go:1509`, `nomad/server.go:744`), signal handler `SIGINT`/`SIGTERM`/`SIGHUP` + `LeaveOnInt`/`LeaveOnTerm` (`command/agent/command.go:1093-1162`, `command/agent/config.go:123-127`) with a `5s + Drain.Deadline` graceful window (`command/agent/command.go:47-48,1052`). No Unix-domain socket path by default — `api.Config.Address` defaults to `http://127.0.0.1:4646` (`api/api.go:338-339`) with optional `unix://` support (`api/api.go:301-315`, `command/agent/http.go:154-160` via `config.Listener`).

## Rating

**Rating: 8 / 10**

Rationale: Server ownership and detach/reconnect survivability are production-proven: jobs are Raft-durable, schedulers are server-owned, clients only execute allocations, and any external observer can poll or stream with `X-Nomad-Index`/`index` cursors. Authentication (`X-Nomad-Token` / `Bearer` + mTLS), explicit `region`/`namespace` scoping, `idempotency_token`, blocking-query long-poll, and event-stream heartbeat/reconnect cover most of Phase 14/15 concerns. Deducted two points because (a) the primary loopback model is network TCP (not an exclusive Unix socket/loopback with OS-level peer auth) and (b) there is no Aren-style filesystem `workspace` context; isolation is logical (namespace/region) rather than a directory-scoped inventory lock, and local caller authentication relies on ACL tokens/mTLS rather than UID/peer-credential binding.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Process that owns accepted work | `type Agent struct` holds `*nomad.Server` + `*client.Client` + `shutdownCh`; `NewAgent` builds both from single config | `command/agent/agent.go:72-146`, `command/agent/agent.go:149-194` |
| Server ownership / durable work | Jobs are persisted via `fsmStateStore` / Raft `LogStore`; `Register` mutates `state.Txn` then Raft-apply; allocations/evals remain after submitter exits | `nomad/fsm.go:1036`, `nomad/job_endpoint.go:990`, `nomad/state/state_store.go:*` |
| Client execution model | Client is executor only; `Client.NewClient` + `allocrunner` persists allocs to `StateDB`; `client.Leave()`/`Shutdown()` decoupled from job lifetime | `client/client.go:72-80`, `command/agent/agent.go:1355-1425` |
| Local transport setup | `NewHTTPServers` binds `config.normalizedAddrs.HTTP`; TLS wrap via `tls.NewListener`; in-process `builtinListener`/`builtinDialer` via `bufconndialer` for Task API | `command/agent/http.go:117-262`, `command/agent/agent.go:1388-1394` |
| Address defaults | `api.DefaultConfig` address `127.0.0.1:4646`; env override `NOMAD_ADDR`; `unix://` parsed to `http://127.0.0.1` with `DialContext` → `net.Dial("unix", path)` | `api/api.go:338-392`, `api/api.go:301-315` |
| Caller authentication | `parseToken` checks `X-Nomad-Token` then `Authorization: Bearer`; `ResolveToken` → `Server.Authenticate` → `ResolveACL`; Task builtin API enforces `authMiddleware` with `ACL.WhoAmI` | `command/agent/http.go:1030-1069`, `command/agent/http.go:372-401`, `command/agent/http.go:1198-1262` |
| mTLS | `tlsutil.NewTLSConfiguration` + `IncomingTLSConfig`; `VerifyHTTPSClient`; `CAFile`/`CertFile`/`KeyFile` via Config/TLSConfig | `command/agent/http.go:143-173`, `command/agent/config.go:175-177,370-378` |
| Workspace/tenant scoping | `parseRegion` defaults to `agent.GetConfig().Region`; `parseNamespace` defaults to `structs.DefaultNamespace`; `QueryOptions{Region,Namespace}` and `WriteRequest{Region,Namespace,IdempotencyToken}` | `command/agent/http.go:978-994,987-994,1071-1087`, `api/api.go:58-65,114-135` |
| Namespace wildcard / search | `validateACL` expands `*` sentinel to all namespaces via `State().NamespaceNames()` | `nomad/event_endpoint.go:264-281` |
| Request identity / idempotency | `parseIdempotencyToken` → `WriteRequest.IdempotencyToken`; `JobRegister` uses `IdempotencyToken` to dedup writes | `command/agent/http.go:996-1001`, `api/api.go:134` |
| Attach / status poll / blocking query | `Job.{allocations,List}` etc use `MinQueryIndex`/`MaxQueryTime`; server `blocking query` loop; client returns `X-Nomad-Index`/`X-Nomad-KnownLeader`/`X-Nomad-LastContact` | `nomad/job_endpoint.go:1498-1529`, `nomad/alloc_endpoint.go:58,142-149`, `command/agent/http.go:924-947,880-913` |
| Stream reconnect | HTTP EventStream: `?index=N&topic=Job&topic=Allocation`; server `SubscribeRequest{Index,Topics,Namespaces}` + `publisher.Subscribe`; JSON ndjson + heartbeat; client `EventStream.Stream(ctx,topics,index,q)` with `strconv.FormatUint(index)` | `command/agent/event_endpoint.go:46-70,152-196`, `nomad/event_endpoint.go:89-128`, `api/event_stream.go:184-245` |
| Cancellation from second client | Second client can `DELETE /v1/job/:id` (deregister), `POST /v1/job/:id/gc`, node drain `POST /v1/node/:id/drain`, alloc stop/signal via `ClientAlloc` RPC | `command/agent/http.go:404-409` (`JobSpecificRequest`), `nomad/job_endpoint.go:1498`, `command/agent/node_endpoint.go:138`, `nomad/node_endpoint.go:963` |
| Startup | `Command.Run` → `readConfig` → `setupAgent` → `NewAgent` → `NewHTTPServers` → `startupJoin` → `handleRetryJoin` → `handleSignals`; `VersionInfo` advertised | `command/agent/command.go:818-966`, `command/agent/command.go:1421-1455,968-1035` |
| Signal handling | `handleSignals` traps `SIGINT`/`SIGTERM`/`SIGHUP`/`SIGPIPE`; `LeaveOnInt`/`LeaveOnTerm`; `winsvc.ShutdownChannel`; `sdNotify` (`NOTIFY_SOCKET`) | `command/agent/command.go:1092-1162`, `command/agent/sdnotify_linux.go:18-37` |
| Graceful drain | `terminateGracefully`: `gracefulTimeout 5s + client Drain.Deadline`; calls `Agent.Leave()` async with `sync.Once` `gracefulCh` + `time.Timer(timeout)`; second signal forces exit | `command/agent/command.go:1037-1090,47-48,1047-1055` |
| Leave / shutdown ordering | `Agent.Leave` → `client.Leave` + `server.Leave` (serf leave + Raft RemoveServer); `Agent.Shutdown` closes `taskAPIServer`, then `client.Shutdown`, `server.Shutdown`, `consulServices.Shutdown`, closes `shutdownCh` | `command/agent/agent.go:1493-1551`, `nomad/server.go:744-802,814-863` |
| Server shutdown internals | `Server.Shutdown` stops workers (`workerShutdownGracePeriod 5s`), serf, Raft, `rpcListener`, `connPool`, `fsm`, `consulConfigEntries`, `oidcProviderCache` | `nomad/server.go:744-802,94-96` |
| Restart / rejoin | `SIGHUP` → `handleReload` → `ShouldReload` + `Reload` for agent/server/client + `reloadHTTPServer`; `retryJoin` via `go-discover` for server + client `SetServers` | `command/agent/command.go:1182-1271,1590-1630,968-1035` |
| Versioning / compatibility | No per-request `ApiVersion` header; `/v1/` prefix is stable; `version.GetVersion().FullVersionNumber()` advertised in `CLI` + agent info; paginated/filter semantics stable (`per_page`/`next_token`/`filter`/`reverse`) | `main.go:98-99`, `command/agent/command.go:912-944`, `command/agent/http.go:1089-1122,897-905` |
| Wire contract separation | Public contract = HTTP JSON + headers (`X-Nomad-Index`, `X-Nomad-Token`, `Content-Type: application/json`, ndjson event stream); internal = msgpack/yamux streaming RPC (`codec.MsgpackHandle`, `StreamingRpcHandler`, `net.Pipe` hijack) | `api/api.go:675-814`, `command/agent/http.go:87-88,750-825`, `command/agent/event_endpoint.go:92-94,178-180` |
| Thin multi-language contract | Pure HTTP+JSON is sufficient; `api` Go client is thin wrapper over `newRequest`/`doRequest`/`query`/`write`; no codegen required; `curl` + `jq` works | `api/api.go:816-867,889-906,1040-1078`, `api/event_stream.go:184-246` |
| Concurrency / multi-client | State is Raft-linearizable; blocking queries + memdb watch + event broker fan-out allow many readers; tests cover concurrent list/watch | `nomad/rpc.go`, `nomad/state/state_store.go`, `nomad/stream/event_broker_test.go:69`, `command/agent/event_endpoint_test.go` |

## Answers to Dimension Questions

**Does client exit affect the accepted run?**

No. The submitting process is not the owner. A `nomad job run` is an HTTP `PUT /v1/job/:id` / `POST /v1/jobs/parse` + `Register` RPC (`command/agent/http.go:404`, `nomad/job_endpoint.go:990`) that is Raft-committed (`nomad/server.go:768-771` `raft.Shutdown` path implies Raft persistence, `nomad/fsm.go:1036` handling of alloc reconnects) and evaluated by server workers (`nomad/worker.go:363`). The server persists allocations/deployments/evaluations in the FSM-backed `StateStore`. The CLI can exit immediately after the `WriteMeta.Index` returns; the job continues scheduling until `job stop`/`purge` or GC (`nomad/core_sched.go:689` `allocGCEligible`). Client nodes only run the allocations assigned to them (`client/client.go` allocrunner); even if the originating client agent exits, other clients continue. Evidence for durable execution: `command/agent/agent.go:149-194` shows `Agent` decouples Server (owner) from Client (executor), and `client/client.go:72` `clientRPCCache` shows client reconnect is expected to be transient.

**Can another client find and observe the same run without sharing process memory?**

Yes. Any process that can reach the HTTP address and present the same `region`+`namespace` can discover and tail the same job:

* Poll: `GET /v1/job/:id?namespace=X&region=Y` (`command/agent/http.go:408` `JobSpecificRequest`), `GET /v1/job/:id/allocations?index=I&wait=60s` (`nomad/job_endpoint.go:1498` blocking query), `GET /v1/allocations?prefix=...`, `GET /v1/evaluations`, filtered/paginated via `filter`, `per_page`, `next_token` (`command/agent/http.go:1089-1122`).
* Headers provide linearizable cursor: server sends `X-Nomad-Index`, `X-Nomad-LastContact`, `X-Nomad-KnownLeader` (`command/agent/http.go:880-913`), client sends `?index=&wait=` (`command/agent/http.go:924-947`).
* Stream: `GET /v1/event/stream?index=N&topic=Job&topic=Allocation&namespace=X` is hijacked to ndjson (`command/agent/event_endpoint.go:46-209`); server `Event.Stream` streaming RPC (`nomad/event_endpoint.go:32-128`) fans out from `stream.Subscription.Next(ctx)` via `stream.NewJsonStream(ctx, 30s)` heartbeat. Reconnect is just re-issuing `?index=lastSeen` (`api/event_stream.go:184-245`).
* Logs/exec: `GET /v1/client/fs/logs`, `/v1/client/allocation/:id/exec` via websockets (`command/agent/http.go:468-476`, `command/agent/websockets.go:17-116`).
* No shared memory: discovery is via Raft-committed state + gossip (`nomad/server.go:57-73` serf) and HTTP, not file handles or in-process channels.

**How are local callers authenticated and workspaces identified?**

*Authentication:* Primary is ACL token. Every request can carry `X-Nomad-Token: <secret>` or `Authorization: Bearer <secret>` (`command/agent/http.go:1030-1069` `parseToken`). Server resolves it via `Server.Authenticate`/`ResolveACL` (`command/agent/http.go:372-401`, `nomad/acl.go:11`). Token capabilities gate `AllowNsOp` per namespace (`nomad/event_endpoint.go:283-334`). The Task API builtin server (`command/agent/agent.go:1388-1394`, `command/agent/http.go:228-259`) mandates authentication via `authMiddleware` (`command/agent/http.go:1198-1262`, `newAuthMiddleware`). Optional mTLS secures both HTTP and RPC: `TLSConfig.VerifyHTTPSClient` / `VerifyServerHostname` (`command/agent/http.go:143`, `command/agent/config.go:370-378`, `helper/tlsutil`). Unix-socket peer credential auth does not exist — local caller is not authenticated by OS UID, only by token/mTLS.

*Workspace identification:* Explicit query params, not a filesystem directory. Every request is scoped by `?region=` (defaults to `agent.GetConfig().Region` → `global` if unset, `api/api.go:48,844-848`) and `?namespace=` (defaults to `default`, `api/api.go:41`, `command/agent/http.go:987-994`), plus `?prefix=` for list filtering and `?per_page`/`next_token`/`filter`/`reverse` for pagination (`command/agent/http.go:1089-1122`). Jobs themselves embed `Namespace`/`Region` in their spec. There is no Aren-style `workspace` path with durable inventory under it; instead the server's global store is partitioned by `(region, namespace)` and ACL policy.

**What is the smallest stable contract a non-Go client must implement?**

A thin HTTP/1.1 + JSON client is sufficient — no generated SDK, no msgpack. The envelope is:

1. Base URL: `http(s)://host:4646` (`api/api.go:338`) or `unix:///path.sock` via `unix://` dial (`api/api.go:301-315`). Set `Host` from URL.
2. Auth header: `X-Nomad-Token: <secret>` or `Authorization: Bearer <secret>` (`command/agent/http.go:1032-1055`).
3. Scoping: `?region=&namespace=&prefix=&per_page=&next_token=&filter=&reverse=&stale=&index=&wait=&idempotency_token=` (`command/agent/http.go:924-947,949-968,996-1001,1071-1087`; `api/api.go:690-765`).
4. Verbs: `PUT /v1/job/:id` (register), `POST /v1/jobs/parse`, `GET /v1/job/:id`/`/allocations`/`/evaluations`/`/deployments`/`/namespaces`/`/regions`/`/status/leader` with blocking `?index=&wait=` long-poll (`command/agent/http.go:404-534` route table), `DELETE /v1/job/:id`, `GET /v1/event/stream?index=&topic=Job:…` → `200 OK` hijacked ndjson (`command/agent/event_endpoint.go:71-145` sets `Content-Type: application/json`, `Cache-Control: no-cache`, strips `Content-Encoding`, hijacks to `deadlineWriter` 40s per-write + `ioutils.NewWriteFlusher`).
5. Headers to handle: send `Accept-Encoding: gzip` (`api/api.go:805`); parse `X-Nomad-Index`, `X-Nomad-KnownLeader`, `X-Nomad-LastContact`, `X-Nomad-NextToken` (`command/agent/http.go:880-913`, `api/api.go:1180-1210`); handle `429 Too Many Requests` per-IP conn limiting (`command/agent/http.go:317-344` `makeConnState`/`connLimiter`).
6. Event stream reconnect: replay with `?index=<lastIndex>`; heartbeats are `{"Index":0,"Events":[]}` filtered client-side (`api/event_stream.go:228-230`); errors arrive as `{"Err": "…"}`.
7. No Go-type leakage required: server tolerates missing `Pretty`, `stale`, etc. The internal `msgpack` + yamux streaming RPC (`nomad/event_endpoint.go:36-42`, `nomad/rpc.go`, `command/agent/agent.go:564` `yamux.DefaultConfig`) is hidden behind the HTTP hijack; a polyglot client never speaks it directly — only the HTTP side.

In short: construct `request{method, url, params, token}` (`api/api.go:675-686`), serialize body with `json.NewEncoder` (`api/api.go:1243-1254`), send via `http.Client.Do`, ungzip if `Content-Encoding: gzip` (`api/api.go:910-934`), decode JSON (`api/api.go:1226-1237`), and for streams decode newline-delimited `Events` (`api/event_stream.go:220-242`).

## Architectural Decisions

* **Cluster daemon, not loopback daemon.** One `agent` binary serves both roles (`command/agent/agent.go:72-73`). `NewAgent` (`command/agent/agent.go:149`) instantiates `nomad.Server` and/or `client.Client` based on `Config.Server.Enabled`/`Client.Enabled` (`command/agent/command.go:344-348`). This matches Nomad's data-center product: the server quorum owns scheduling, clients own execution. Consequence: lifetime is decoupled from any CLI — the CLI (`command/meta.go` via `command/agent/command.go:85` `Meta`) is a one-shot HTTP caller (`api/api.go:890-906` `doRequest`), not a session parent.

* **Raft-linearizable state over in-memory maps.** Jobs/evals/allocs/deployments are applied through Raft (`nomad/server.go:768-771`, `nomad/raft_rpc.go`) into `state.StateStore` via `fsm.go`. Blocking queries wait on `state.WatchSet` / MemDB index increments, allowing efficient long-poll without busy loop (`nomad/job_endpoint.go:1498-1530` pattern). This is heavier than a single-node WAL but gives multi-server durability and leader failover — appropriate for Aren Phase 15 upgrade semantics if one maps `region → Aren deployment`.

* **HTTP as the stable public API; RPC/msgpack as internal.** Public handlers (`command/agent/http.go:404-575` `registerHandlers`) translate to internal RPCs (`Agent.RPC` → `nomad.Server.RPC` or `client.RPC`, `command/agent/agent.go:1553-1559`). Streaming is bridged via `net.Pipe` + msgpack (`command/agent/event_endpoint.go:92-200`, `nomad/event_endpoint.go:36-42`). The thin HTTP envelope is intentionally polyglot; all Go-specific codec details stay internal — a deliberate separation that Aren should emulate for multi-language CLIs.

* **Explicit tenant scoping via query, not filesystem.** `region`/`namespace` are first-class query params (`command/agent/http.go:978-994`). `AllNamespacesSentinel="*"` (`nomad/event_endpoint.go:268`) expands via `NamespaceNames()` with ACL gating (`nomad/event_endpoint.go:275-281`). This avoids needing a `.ultraplan/workspace` directory, but means workspace is logical; durable inventory is the `StateStore` keyed by `(region, namespace, jobID)`.

* **Token + optional mTLS for caller identity.** ACL (`acl/*`, `nomad/auth/auth.go:263`) gates every RBAC decision (`validateNsOp`, `nomad/event_endpoint.go:283`). `Config.TLSConfig` (`command/agent/config.go:175`) optionally enforces mTLS for both HTTP and RPC. No SO_PEERCRED binding — local auth is bearer token, which is simpler for remote operators but weaker for loopback-only zero-config.

* **Cooperative graceful leave with deadline-bounded drain.** `Command.terminateGracefully` (`command/agent/command.go:1037`) runs `Agent.Leave` async, extends timeout by `client.Drain.Deadline` (`command/agent/command.go:1052`), waits on `gracefulCh` or second signal/timer. Server `Leave` does Raft `RemoveServer` + Serf `Leave` (`nomad/server.go:814-863`). Client drain migrates allocations with deadline semantics (`nomad/drainer/*`). This directly maps to Aren's expected graceful drain / forced shutdown / restart upgrade path, with configurable `DrainSpec.Deadline` per node.

## Notable Patterns

* **BufConn in-process HTTP.** `bufconndialer.New()` provides a `net.Listener`/`Dialer` pair (`command/agent/agent.go:1388`) so consul-template and Task API can speak HTTP without leaving the process — similar to Aren's desired loopback-only API without opening a public port.

* **Hijacked ndjson stream over plain HTTP.** `EventStream` upgrades by hijacking (`http.Hijacker`, `command/agent/event_endpoint.go:104-145`) and writes a raw `http.Response` `200 OK` then newline-delimited JSON frames via `ioutils.NewWriteFlusher` + `deadlineWriter` (40s per-write, `command/agent/event_endpoint.go:26-44`). Detecting client close via `io.Copy(io.Discard, conn)` cancellation (`nomad/event_endpoint.go:142-145`) is a robust pattern for detach detection.

* **Atomic ACL refresh inside subscriptions.** `Event.stream` stores `resolvedACL` in `atomic.Value` and re-authenticates on a timer/callback (`nomad/event_endpoint.go:73-111` `Authenticate` closure), so token expiry/rotation (`exp` from claims `nomad/event_endpoint.go:132-137`) doesn't kill existing streams until checked before `jsonStream.Send`. This enables long-lived multi-client attach.

* **Connection state + per-IP rate limiting.** `makeConnState` (`command/agent/http.go:271-315`) tracks `atomic.Int32 connCount` (emitted as `nomad.agent.http.connections` gauge via ticker `command/agent/http.go:205-216`) and wraps `connlimit.NewLimiter` with a 10/s token bucket for 429 generation — a daemon hardening pattern for multi-client abuse.

* **HCL → normalizedAddrs → dynamic reload.** `readConfig` → `normalizeAddrs` → `ShouldReload` (`command/agent/agent.go:1590-1630`) → `Reload` + `reloadHTTPServer` (`command/agent/command.go:1164-1180`) allows SIGHUP to rotate certs and rebind without dropping Raft quorum — an upgrade path without full restart.

## Tradeoffs

* **Network-reachable by default vs. loopback-only.** Default bind is not 127.0.0.1-only if unset? `BindAddr` resolution via `go-sockaddr` (`command/agent/config.go:1109`) often yields `0.0.0.0` unless configured; HTTP is reachable from off-host if `addresses.http` is `0.0.0.0`. Security then depends on ACL/mTLS rather than OS isolation. Aren's ambition for loopback-only + workspace directory lock would be stricter. The upside is Nomad's topology natively spans hosts (servers + clients).

* **Raft/Serf overhead for single-host dev.** Even `dev` mode (`command/agent/command.go:88-94` `-dev` flags, `command/agent/agent.go:1276-1279` `DevMode` nodeID) still runs Raft+Serf in-process. For a purely local Aren daemon this would be overkill; Nomad trades startup cost (peer polling `peersPollInterval 45s` `nomad/server.go:68-73`, `defaultConsulDiscoveryInterval 3s`) for strong consistency.

* **JSON fidelity vs. binary efficiency.** HTTP JSON uses `codec.JsonHandleWithExtensions` (`command/agent/http.go:816-818`) and supports `?pretty` (`command/agent/http.go:799-804`). Internal RPC uses msgpack (`structs.MsgpackHandle`) for streaming. The thin HTTP path pays JSON cost but stays universal; high-frequency event streams use msgpack internally then convert to ndjson — two serializations but good interop.

* **Logical tenant isolation vs. filesystem workspace lock.** Namespace/region ACLs isolate, but there is no `flock`-guarded workspace directory preventing concurrent writers on the same job file. Conflicts surface as Raft optimistic concurrency via `JobModifyIndex` / `EnforceVersion` (`nomad/job_endpoint.go:xxx`, `api/jobs.go:573-582` `Revert`/`Stable`). This is less explicit than Aren's workspace file convention.

* **Per-IP conn limit 429 vs. absolute backpressure.** `HTTPMaxConnsPerClient` (`command/agent/config.go:194`, `command/agent/http.go:135-142`) returns 429 with `metrics.IncrCounter` but only when TLS handshake completes; bursty local clients can still briefly queue. Good for multi-client DoS but not a hard semaphore.

## Failure Modes / Edge Cases

* **Submitter dies before Raft commit.** `Job.Register` is async from client's view: `PUT` returns after Raft apply, but if agent dies during `plan_apply` (before eval brokers), job may be persisted yet never scheduled. Recovery is via re-reading `GET /v1/job/:id` by any client — polling is idempotent — and via `batch_eval_gc_threshold` / `eval_broker` retries (`nomad/eval_broker.go`, `nomad/plan_apply.go`).

* **Token expiry mid-stream.** `exp` from ACL token claims (`nomad/event_endpoint.go:132-137`) is checked before each `jsonStream.Send`; expired token yields `structs.ErrTokenExpired` and `handleJsonResultError` encodes it with `code=403` (`nomad/event_endpoint.go:162-169,253-262`). Second client with refreshed token can reconnect with `?index=lastSeen` without loss of events (broker retains `EventBufferSize`, `command/agent/config.go:725-731`).

* **Leader loss / `KnownLeader=false`.** HTTP headers include `X-Nomad-KnownLeader` (`command/agent/http.go:885-893`, `api/api.go:1202-1209`); during leader election `AllowStale` reads may be stale but `QueryMeta.KnownLeader` signals it. Clients should retry with `?stale` or wait for `?index=` to unblock after new leader `Work`.

* **Second signal during graceful drain.** `terminateGracefully` returns `1` immediately on second `SIGINT`/`SIGTERM` (`command/agent/command.go:1069-1081`), abandoning drain — allocations become orphaned until GC. `DrainDeadlineNotifier` (`nomad/drainer/drain_heap.go:12`) still fires `Deadline` expiry, forcing reschedule after `ForceDeadline` (`nomad/node_endpoint.go:963-964`).

* **Hijack not available.** If reverse proxy strips `Hijacker` or `http/2` disables hijack, `EventStream` falls back to non-hijacked `ioutils.NewWriteFlusher(resp)` (`command/agent/event_endpoint.go:147-150`) — functional but without raw connection deadlines; proxy buffering can delay ndjson flushing.

* **Concurrent writers, no workspace lock.** Two `job run` with same `(namespace, jobID)` race on `JobModifyIndex`; last-write-wins unless caller passes `enforce_prior_version` (`api/jobs.go:573-582`). No filesystem lock detects this — diagnosis requires inspecting `ModifyIndex`/`Version` (`api/jobs.go:1699-1721` `TagVersion`).

* **TLS rotation failure.** `ShouldReload` detects cert change via `CertificateInfoIsEqual` (`command/agent/agent.go:1606-1611`); if `KeyLoader.LoadKeyPair` fails, `Reload` returns error and retains old certs (`command/agent/agent.go:1719-1722`). HTTP reload (`reloadHTTPServer` `command/agent/command.go:1166-1180`) recreates listeners without dropping Raft.

* **Raft log corruption.** Background `LogStoreVerification` (`command/agent/config.go:1011-1055`) periodically verifies `raft-wal` (`raftwal.DefaultSegmentSize 64M` `command/agent/agent.go:665`); corruption path yields `Raft Store verification` error surfaced via `server.Shutdown`.

## Future Considerations

* **Loopback Unix socket as default for Aren parity.** Expose a `listener "unix" { address="/run/aren.sock" }` alongside TCP, enforce `SO_PEERCRED` UID mapping to an `ACL` identity (extend `parseToken` `command/agent/http.go:1030` to check peercred). This would satisfy "local callers authenticated" without requiring mTLS for single-host mode.

* **Introduce a filesystem workspace handle alongside namespace.** Add a `?workspace=/abs/path` or `X-Aren-Workspace` header that maps to a `state_dir/<hash>` subtree, `flock` on `workspace.lock` before `Job.Register`, and expose `X-Aren-Workspace-Index` analogous to `X-Nomad-Index`. This would give Aren's explicit workspace context without breaking region/namespace.

* **Expose a minimal non-Go contract test suite.** Current `api/api_test.go` exercises Go client only. Add a contract test that asserts raw `curl` + `jq` against `/v1/jobs`, `/v1/event/stream`, and header semantics — enforcing the thin HTTP contract stays stable while internal msgpack types evolve.

* **Typed error envelope for reconnect logic.** `handleJsonResultError` (`nomad/event_endpoint.go:253`) already sends `RpcError.Code`; propagate a stable `code` enum (`permission_denied`, `token_expired`, `rate_limited`) in ndjson so non-Go reconnect logic can distinguish retryable (429, leader loss) from fatal (403).

* **Backpressure for slow event consumers.** Current `deadlineWriter` 40s per-write (`command/agent/event_endpoint.go:98`) will drop slow consumers only on write timeout; consider a bounded `Subscription.OutCh` with `EventBufferSize` tuning (`command/agent/config.go:725-731`) exposed as observable gauge for multi-client saturation.

## Questions / Gaps

* **No evidence of per-workspace inventory file lock.** Grep for `workspace` in `command/agent/*` and `api/*` returns only unrelated volume `HostVolumesDir`/`CSI` paths; Grep for `flock`/`LOCK` shows none in `nomad/`. Confirms workspace is logical (namespace), not filesystem — gap versus Aren dimension.

* **No Unix socket auth binding.** `unix://` is parsed in `api/api.go:301-315` but server `NewHTTPServers` creates listeners via `config.Listener("tcp", ip, port)` (`command/agent/http.go:160`); no `Listener("unix", …)` for HTTP in this build. Local peer-crendential auth is not implemented — only `X-Nomad-Token`/`Bearer` (`command/agent/http.go:1032-1055`) and mTLS.

* **Protocol versioning not in URL.** No `/v2/` or `Accept-Version` header; compatibility is via additive fields and `pretty`/`filter`/`per_page`/`next_token` stability. How a breaking Aren contract would be versioned is unanswered — likely requires a new `/v1/` minor or header.

* **Disconnect-then-cancel from second client leaves no explicit session.** There is no `session`/`runID` lease that survives client death beyond the job ID; cancellation is `DELETE /v1/job/:id` + `Purge` (`nomad/job_endpoint.go:xxx`) which needs RBAC on `(namespace, job)`. The dimension's "attach/detach/stream reconnect" maps to `index` cursor + `event stream`, not to a revocable session token — verify if Aren expects explicit `detached` state versus nomad's implicit detached (any token can poll).

* **Shutdown under load not fully exercised in reviewed files.** `server_test.go` and `agent_test.go` cover `Shutdown` idempotency (`command/agent/agent.go:1511-1516` `shutdown` bool) and `ShutdownCh` close (`command/agent/agent.go:1549`), but no test was cited that kills `raftStore` mid-apply under load. Would need explicit `go test -race -run Shutdown` coverage review.

---

Generated by `01.13-daemon-local-api-and-multi-client-lifetime` against `nomad`.
