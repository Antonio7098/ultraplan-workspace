# Source Analysis: nats-server

## 01.13 Daemon, Local API, and Multi-Client Lifetime

### Source Info

| Field | Value |
|-------|-------|
| Name | nats-server |
| Path | `studies/aren-go-runtime-study/sources/nats-server` |
| Language / Stack | Go 1.26.0 (`server` package, `main.go:14-24`) |
| Analyzed | 2026-08-30 |

## Summary

`nats-server` is a long-lived Go broker daemon, not a task runner. The process owns all durable state (streams, consumers, messages on `StoreDir`) independently of any submitting client. `main.go:97-134` creates a `Server` via `server.NewServer(opts)` (`server/server.go:706`), calls `server.Run(s)` (`server/service.go:20-23`) which invokes `s.Start()` (`server/server.go:2258`), then blocks on `s.WaitForShutdown()` (`server/server.go:2769`). All work — subscriptions, published messages, JetStream assets — outlives the client that created it: a stream created via `$JS.API.STREAM.CREATE.*` (`server/jetstream_api.go:54`) persists in the filestore (`server/filestore.go:61,438`) after the publisher disconnects, and durable consumers survive disconnect and can be re-observed by a different client. Transport is TCP (default `0.0.0.0:4222`, `server/const.go:78`), not a loopback Unix socket; authentication is account/JWT/nkey/token/TLS based, not OS-peer-credential local-API auth; versioning is via `PROTO`/`VERSION`/`JSApiLevel` in the `INFO` block. There is no Aren-style workspace directory or run-ID reconnection — the closest analogues are `Account` isolation and file-backed `StoreDir` plus durable consumer names for multi-client attach, and HTTP `monitor.go` for out-of-process observation.

## Rating

**6/10** — Excellent long-lived daemon ownership and production-grade graceful-drain/shutdown, with multi-client shared inventory via JetStream, but mismatched to the dimension's loopback/local-API/workspace/attach-detach requirements. No Unix-socket local API, no per-run workspace context, no explicit run-stream reconnect; multi-client observability is via message subjects and HTTP monitoring, not a dedicated run handle. Rating reflects strong daemon lifecycle, moderate multi-client contract clarity, weak local-API/workspace alignment.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Daemon entrypoint & blocking ownership | `main()` parses flags, `NewServer`, `ConfigureLogger`, `Run(s)`, then `WaitForShutdown()` — server owns lifecycle after client exit | `main.go:97-134` |
| Non-Windows Run wrapper | `Run(server) { server.Start(); return nil }` — Start is non-blocking | `server/service.go:20-23` |
| Server.Start lifecycle | `Start()` logs version, `running.Store(true)`, `grRunning=true`, `startRateLimitLogExpiration`, `StartMonitoring`, resolver, JetStream, `delayedAPIResponder`, gateways/websocket/leaf/MQTT, routes, `AcceptLoop` | `server/server.go:2258-2571` |
| Shutdown ownership | `Shutdown()` CAS `shutdown` atomic, `signalPullConsumers`, `stepdownRaftNodes`, `shutdownEventing`, close all listeners (`listener`, `routeListener`, `leafNodeListener`, `gatewayListener`, `http`, `profiler`, `websocket`, `mqtt`), `close(s.quitCh)`, close conns with `ServerShutdown`, wait `done`+`grWG`, `close(shutdownComplete)` | `server/server.go:2579-2748` |
| WaitForShutdown | `WaitForShutdown() { <-shutdownComplete }` — callers block until drain complete | `server/server.go:2769-2771` |
| Signal handling & graceful vs forced | `handleSignals` subscribes `SIGINT,SIGTERM,SIGUSR1,SIGUSR2,SIGHUP`; `handleSignal` maps `SIGINT->Shutdown+Wait+Exit(0)`, `SIGTERM->Shutdown unless ldm`, `SIGUSR1->ReOpenLogFile`, `SIGUSR2->lameDuckMode`, `SIGHUP->Reload` | `server/signal.go:37-86` |
| Graceful LDM drain | `lameDuckMode()` sets `ldm=true`, closes `listener`+websocket, `ldmCh`, `transferRaftLeaders`, `shutdownJetStream`, waits `ldmCh`, sends `LDM INFO` to routes/clients with filtered `ClientConnectURLs`, staggers `closeConnection(ServerShutdown)` over `LameDuckDuration - GracePeriod` | `server/server.go:4446-4573` |
| Signal commands & PID control | `Command` constants `stop/quit/reopen/reload/ldm/term`, `ProcessSignal`, `CommandToSignal` mapping to `SIGKILL/SIGINT/SIGUSR1/SIGHUP/SIGUSR2/SIGTERM` | `server/const.go:22-35`, `server/signal.go:93-166` |
| Opt-out of signal handling | `NoSigs bool` prevents `handleSignals` | `server/opts.go:426`, `server/signal.go:38-39` |
| Config reload (restart/upgrade path) | `SIGHUP -> Reload()`; `Reload` path updates `info.TLSRequired`, `routeInfo` etc; upgrade otherwise requires process restart (no hot binary upgrade) | `server/signal.go:80-85`, `server/reload.go:243` |
| Listener / transport setup | `AcceptLoop` binds `Opts.Host:Port` via `natsListen("tcp", hp)`, advertises `ClientURL()`, handles `ProxyProtocol`, TLS sniff/handshake, auth timer, spawns `readLoop`+`writeLoop` goroutines | `server/server.go:2773-2860`, `server/server.go:2890-2916` |
| In-process local transport (not UDS) | `InProcessConn() { net.Pipe() -> createClientInProcess }` — loopback via `net.Pipe`, respects `DontListen` | `server/server.go:2873-2888`, `server/opts.go:417` |
| Monitoring / HTTP local API | `StartMonitoring` binds `HTTPHost:(HTTPPort|HTTPSPort)`, mux `/varz,/connz,/routez,/gatewayz,/leafz,/subsz,/accountz,/jsz` etc; `HTTPHandler()` exposes handler | `server/server.go:3014-3212` |
| Version negotiation: client proto | `PROTO=1` current, `ClientProtoZero=0`, `ClientProtoInfo=1` — client announces `protocol` in `CONNECT` | `server/const.go:68-75`, `server/client.go:82-88` |
| Version negotiation: server-to-server proto | `RouteProtoZero=0, RouteProtoInfo=1, RouteProtoV2=2, MsgTraceProto=3`, `getServerProto()` | `server/server.go:62-99` |
| Server INFO / handshake | `type Info` includes `server_id, server_name, version, proto, git_commit, go, host, port, headers, auth_required, tls_required, max_payload, jetstream, cluster, domain, connect_urls, ldm, js_api_level` etc | `server/server.go:108-165` |
| JetStream API level | `JSApiLevel` advertised in `INFO`: `info.JSApiLevel = JSApiLevel` | `server/server.go:758`, `server/server.go:138` |
| JetStream API durable inventory | Subjects `$JS.API.STREAM.CREATE.*`, `DELETE`, `INFO`, `CONSUMER.CREATE.*`, `CONSUMER.INFO` etc — persistent, not tied to submitting connection | `server/jetstream_api.go:38-349` |
| Durable store (workspace analogue) | `Options.StoreDir`, `JetStreamConfig.StoreDir`, `filestore.go:StoreDir`, `filepath.Join(fcfg.StoreDir, msgDir)` — filesystem durable inventory | `server/opts.go:483`, `server/jetstream.go:2440`, `server/filestore.go:61` |
| Account/workspace scoping | `Options.Accounts`, `Accounts map`, `globalAccountName $G`, `DEFAULT_SYSTEM_ACCOUNT $SYS`, `SystemAccount`, lookup via `ClientInfo.Account/Service` header | `server/opts.go:436-441`, `server/server.go:168-196`, `server/jetstream_api.go:1260-1285` |
| Authentication | `info.AuthRequired` set from `CustomClientAuthentication`, `trustedKeys`, `Nkeys/Users`, `Username/Authorization` | `server/auth.go:322-337` |
| Request identity | Per-client `cid uint64`, `totalClients`, `clients map[uint64]*client`, `eventIds *nuid.NUID`, `ClientInfoHdr` + `getRequestInfo` extraction, `Nats-Required-Api-Level` header | `server/server.go:273-274,200-201,312`, `server/jetstream_api.go:352-359,1260-1285` |
| Concurrent request demux & cancellation | `jsAPIRoutedReq` queue `ipQueue`, `apiDispatch -> queue.push`, `processJSAPIRoutedRequests` worker pool, `QuitC()` for Raft group cancellation | `server/jetstream_api.go:852-860,862-1004` |
| Delayed error responses | `delayedAPIResponse` linked-list ordered by `deadline`, `delayedAPIResponder` goroutine, `sendDelayedAPIErrResponse` | `server/jetstream_api.go:1122-1258` |
| Client create & auth handshake | `createClientEx` sends `INFO`, negotiates TLS first, proxy proto, waits `AuthTimeout`, sets `expectConnect`, pings, starts `readLoop`/`writeLoop` | `server/server.go:3254-3571` |
| No workspace-per-run | No `work_dir`/`workspace` param; `StoreDir` is global per-server, not per-request — grep for workspace returns no evidence | `server/opts.go:483` (global store) |
| No run streaming reconnect | No `stream`/`attach`/`reconnect` RPC for a run handle; JetStream consumers use `JSApiRequestNextT $JS.API.CONSUMER.MSG.NEXT.*` pull, not server-pushed run log reconnect | `server/jetstream_api.go:158-169` |
| Thin client contract | Text protocol `INFO / CONNECT / PUB / SUB / UNSUB / PING / PONG / MSG / HMSG` + JSON monitoring endpoints — implementable without Go | `server/client.go:82-96`, `server/server.go:3014-3212` |

## Answers to Dimension Questions

**Does client exit affect the accepted run?**
No. The daemon process owns the work. `main.go:129-133` calls `server.Run` then `s.WaitForShutdown()` — the server outlives any individual client connection. Publishing a message or creating a stream via `$JS.API.STREAM.CREATE.*` (`server/jetstream_api.go:54`) persists to `StoreDir` (`server/filestore.go:438-446`). `Shutdown()` (`server/server.go:2579`) is triggered only by server-level signals or explicit API, not by client disconnect. `removeClient` (`server/server.go:3726-3755`) merely deletes the `*client` from `s.clients`/`leafs`/`routes`; it does not tear down streams, consumers, or stored messages. Durable JetStream consumers remain queryable after the creator disconnects (`server/jetstream_api.go:142-145` `JSApiConsumerInfo`). Ephemeral (non-durable) consumers are the exception — they are tied to the subscription and removed when the client closes, but streams and durable consumers are not.

**Can another client find and observe the same run without sharing process memory?**
Yes, via durable names and API subjects, not via shared Go pointers. A second TCP client (or HTTP monitor client) can list and inspect the same stream/consumer by name using JetStream API subjects (`$JS.API.STREAM.INFO.*` `server/jetstream_api.go:69`, `$JS.API.CONSUMER.INFO.*.*` `server/jetstream_api.go:144`, `$JS.API.CONSUMER.NAMES.*` `server/jetstream_api.go:135`) or HTTP endpoints `GET /jsz`/`/connz` (`server/server.go:3046-3047`). The server stores state on disk (`server/filestore.go:485-501`) and in `Server` maps that are not shared with clients; clients interact only over the wire. There is no Aren-style `run_id` with log stream resumption — the closest resume is re-subscribing to the same consumer with `JSApiRequestNextT` (`server/jetstream_api.go:158`) or re-reading via `JSApiMsgGet` (`server/jetstream_api.go:100`). `InProcessConn` (`server/server.go:2873`) proves observation can be in-process `net.Pipe` but also fully out-of-process over TCP.

**How are local callers authenticated and workspaces identified?**
*Transport:* TCP `Host:Port` (`server/const.go:78` `DEFAULT_PORT=4222`), optional TLS (`server/server.go:2812-2818`), optional PROXY protocol (`server/server.go:2804`). There is no Unix-domain socket or OS `SO_PEERCRED` peer check; "local" is just `127.0.0.1` with `NoSigs`/`DontListen` variations. Monitoring HTTP(S) binds separately (`server/opts.go:450-453`).
*Authentication:* `Info.AuthRequired/TLSRequired` (`server/server.go:118-119`) set in `server/auth.go:322-337` from `Users`/`Nkeys`/`Username`/`Authorization` token/`CustomClientAuthentication`/`trustedKeys` (JWT operator). Clients send `CONNECT {user, pass, auth_token, jwt, nkey, protocol}` (`server/client.go:2454`). No per-call Unix identity; auth is per-connection credential.
*Workspace / resource isolation:* Accounts (`server/opts.go:436-441`, `server/server.go:195-197`) — `$G` (global), `$SYS` (system), plus JWT-defined accounts with `jsLimits`. `ClientInfo.Account/Service` header selects account (`server/jetstream_api.go:1270-1274`). The only filesystem workspace analogue is the single server-wide `StoreDir` (`server/opts.go:483`) / `JetStreamConfig.StoreDir`, not a per-run directory. No per-request `workspace` path is propagated.

**What is the smallest stable contract a non-Go client must implement?**
A non-Go client can interoperate with only the wire text protocol plus one JSON header convention:
1. TCP connect, read `INFO {server_id, version:"2.15.0-dev" server/const.go:69, proto:1 server/const.go:75, headers:bool, auth_required, tls_required, max_payload, jetstream, connect_urls, ldm}` (`server/server.go:108-165`), then send `CONNECT {"verbose":bool,"pedantic":bool,"protocol":0|1 server/client.go:84-88, "user"/"pass"/"auth_token"/"jwt"/"nkey"}` followed by `PING`/`PONG` (`server/client.go:92-93`).
2. For pub/sub: `PUB <subject> [reply] <bytes>\r\n<payload>\r\n`, `SUB <subject> [queue] <sid>\r\n`, `UNSUB <sid> [max]\r\n`, `MSG`/`HMSG` framing. For JetStream: publish to subjects `$JS.API.> ` (`server/jetstream_api.go:41`) with JSON bodies and await reply on `reply` inbox; honor `Nats-Required-Api-Level` request header (`server/jetstream_api.go:354`) and handle `ApiResponse {type, error}` (`server/jetstream_api.go:410-413`). Optional HTTP GET `http://host:8222/varz|connz|jsz` (`server/server.go:3036-3051`) for monitoring without NATS framing.
No Go types (`Server`, `client`, `Options`, `jetStream`, `RaftNode`) are wire-stable; only the `INFO` JSON shape, the NATS text framing, and the `$JS.API.*` subjects/JSON schemas are stable.

## Architectural Decisions

* **Single `Server` owns all durable inventory.** `Server` struct (`server/server.go:168-387`) holds `clients`, `routes`, `leafs`, `accounts sync.Map`, `js atomic.Pointer[jetStream]`, `store` references. Streams/consumers are stored via `filestore.go` rooted at `StoreDir`, not per-client temp dirs. Tradeoff: simple ownership, but no per-run isolation or garbage-collected workspace.
* **Non-blocking `Start()` + blocking `WaitForShutdown()`.** `Run` (`server/service.go:20`) calls `Start()` then returns; `main.go:133` blocks on `WaitForShutdown()`. This separates daemon startup from lifetime supervision and allows embedding (`DontListen`, `InProcessConn`).
* **Accept-loop-per-listener with goroutine-per-connection.** `AcceptLoop` (`server/server.go:2773`) for clients, plus `startLeafNodeAcceptLoop`, `StartRouting`, `startGateways`, `startWebsocketServer`, `startMQTT` — each listener closes on `Shutdown()` and signals `done`. Tradeoff: straightforward concurrency, relies on `grWG`/`quitCh` (`server/server.go:245-254`) for drain.
* **Lame-duck graceful drain.** `lameDuckMode` (`server/server.go:4446`) implements the only first-class graceful drain: close accept, transfer Raft leadership, send `LDM INFO` with alternative `connect_urls`, stagger `closeConnection(ServerShutdown)` over `DEFAULT_LAME_DUCK_DURATION=2m` minus `GRACE_PERIOD=10s` (`server/const.go:196-200`). No per-run cancel; drain is whole-server.
* **Account + `StoreDir` as workspace analogue.** Accounts provide subject/permission/jetstream-limits isolation; `StoreDir` provides persistence. No per-request filesystem workspace; `JetStreamStoreDir` is a single tree keyed by account/stream/consumer names.
* **Protocol version as single integer + `JSApiLevel`.** `PROTO=1` (`server/const.go:75`) and `Info.Proto` plus `JSApiLevel` (`server/server.go:138`) allow minimal negotiation. Server-to-server has independent `RouteProtoV2`/`MsgTraceProto` (`server/server.go:69-80`). Tradeoff: coarse-grained, no per-method semantic versioning.

## Notable Patterns

* **Signal-to-command mapping.** `Command` strings (`server/const.go:22-35`) mapped to `syscall.Signal` via `CommandToSignal` (`server/signal.go:149-166`) and `ProcessSignal`/`pgrep` for external CLI `nats-server --signal`.
* **Async INFO fan-out.** `addConnectURLsAndSendINFOToClients` / `sendAsyncInfoToClients` counted by `cproto` (`server/server.go:255`, `server/client.go:3733`) — only clients with `protocol >= ClientProtoInfo` get cluster/topology updates.
* **IPQueue for API serialization.** `jsAPIRoutedReqs` / `jsAPIRoutedInfoReqs` (`server/server.go:372-373`) decouple routing `apiDispatch` (`server/jetstream_api.go:862`) from handler goroutines (`processJSAPIRoutedRequests` `server/jetstream_api.go:959`), with `apiInflight` atomic and queue-limit drain + advisory `JSAdvisoryAPILimitReached` (`server/jetstream_api.go:944-955`).
* **Ordered delayed error Responses.** `delayedAPIResponse` linked-list (`server/jetstream_api.go:1122-1134`) ordered by deadline, driven by `delayedAPIResponder` timer + `QuitC()` cancellation — avoids blocking Raft apply paths (`server/jetstream_api.go:1165-1246`).
* **In-process transport for tests/embedding.** `InProcessConn` (`server/server.go:2877`) uses `net.Pipe` and `createClientInProcess` to bypass TCP while honoring `DontListen` — demonstrates thin contract works without network.

## Tradeoffs

* **TCP everywhere vs Unix socket.** Choosing `natsListen("tcp", hp)` (`server/server.go:2870`) maximizes multi-language/client reach and clustering (routes/gateways/leafnodes are all TCP) but gives up true local-only authentication, filesystem-permission workspace isolation, and the small stable loopback contract Aren expects. A loopback-only daemon can use `SO_PEERCRED`/`0700` socket + token file; NATS intentionally does not.
* **Global `StoreDir` vs per-run workspace.** Durability via a single `StoreDir` enables multi-client sharing without copying, but prevents per-run lifecycle operations (attach directory, set env, capture output per run) and makes GC/notification per run impossible without higher-level naming conventions.
* **Durable naming vs run handles.** Streams/consumers are named explicitly by the client and survive indefinitely; there is no server-issued `run_id` with implicit ownership. Simpler for a broker, but forces clients to design their own namespace/GC and to handle `duplicate consumer` races.
* **Whole-server drain vs per-job cancel.** `lameDuckMode` drains the entire node (`server/server.go:4446`); there is no `Cancel(run_id)` or `Detach(stream-id)` RPC for a single accepted job from a second client. Cancellation is per-subscription (`UNSUB`/`consumer delete`) not per-run.
* **Coarse protocol version.** Single `PROTO` integer (`server/const.go:75`) is easy to implement for thin clients, but cannot express per-API breaking changes; `JSApiLevel` and `Nats-Required-Api-Level` header mitigate only slightly (`server/jetstream_api.go:354-355`).

## Failure Modes / Edge Cases

* **Forced shutdown drops in-flight API work.** `Shutdown()` (`server/server.go:2708-2714`) closes client connections with `ServerShutdown` and then waits `grWG`; `processJSAPIRoutedRequests` exits on `quitCh` (`server/jetstream_api.go:991`) without flushing queued `jsAPIRoutedReqs`. No persistence of pending request queue.
* **LDM close staggering can still storm.** `si = dur/numClients` (`server/server.go:4511-4512`) may batch closes (`batch = numClients/dur` when `si<1`). With thousands of clients and short `LameDuckDuration` the 1ms stagger still creates reconnect pressure; test code even notes `batch` workaround vs. `validateOptions` check allowing `gp<0` for tests (`server/server.go:4462`).
* **Multiple `Shutdown()` is no-op but signals are not queued.** `shutdown.CompareAndSwap(false,true)` (`server/server.go:2584`) makes second call return immediately; a `SIGTERM` during `ldm` phase is ignored (`server/signal.go:66-74`). `SIGINT` always does `Shutdown+Wait+Exit(0)` which can race with `lameDuckMode`'s own `Shutdown+Wait`.
* **Monitoring port conflict fails hard.** `StartMonitoring` (`server/server.go:3130-3133`) returns `can't listen to the monitor port`; `Start()` then `Fatalf` and returns early without closing already-opened listeners if called after partial startup (recovery depends on process exit).
* **No authentication for monitoring by default.** HTTP endpoints have no auth unless TLS/middleware is fronted externally; any local process can `GET /varz`/`/connz` without presenting the NATS `auth_token`.
* **JetStream `StoreDir` permission / disk full not isolated per caller.** `filestore.go:438-446` creates `StoreDir` with `defaultDirPerms` and fails `NewServer` if unwritable (`server/server.go:2450` `Fatalf`); a single account can fill disk and affect all accounts (mitigated only by `jsLimits` quotas).
* **Protocol mismatch is fail-closed.** Client `CONNECT` with `protocol <0 or >1` (`server/client.go:2454`) triggers `ErrBadClientProtocol` disconnect, not downgrade.

## Future Considerations

* Add an explicit loopback-local API (Unix-domain socket with `0700` + bearer token file) for Aren-style single-host automation, keeping the TCP NATS protocol unchanged for cross-host clients. Currently `InProcessConn` (`server/server.go:2877`) is the only loopback path and it requires in-process Go.
* Introduce per-run workspace isolation on top of `StoreDir` (e.g., `store_dir/<account>/runs/<run_id>`) with manifest + output capture, rather than overloading stream/consumer names; would enable deterministic GC and second-client log streaming.
* Provide a run-handle primitive (`create_run -> run_id`, `attach/run_id/stream`, `status/run_id`, `cancel/run_id`) layered over `$JS.API.*` so thin non-Go clients can implement only that subset instead of the full publish/subscribe + consumer protocol.
* Expose request identity (`NUID`/`ClientInfo`) consistently on monitoring APIs so a second client can correlate `connz` entries to durable resources without parsing server logs.
* Document the stable subset explicitly (text framing + `INFO`/`CONNECT` + `$JS.API.*` JSON Schemas + HTTP `/varz`/`/connz`/`/jsz`) as a versioned contract with `PROTO`+`JSApiLevel` bump rules; currently version is implied by `VERSION="2.15.0-dev"` (`server/const.go:69`) and `PROTO`.

## Questions / Gaps

* No evidence of Unix-socket local API or `SO_PEERCRED` auth — searched `signal.go`, `server.go` accept paths, `opts.go`; only TCP + `net.Pipe` found. If local-only hardening exists, it is not in this source tree.
* No per-run streaming reconnect protocol found — `jetstream_api.go` pull `NEXT`/`MSG.GET` and `filestore.go` snapshots are message-granular, not run-log byte-stream resumption; no `stream_id`/`offset` reconnection tested in `server_test.go:170` `TestStartupAndShutdown`.
* Workspace lifecycle (create/cleanup per request, env injection, working directory per run) has no implementation — `StoreDir` is server-global; no `workspace` parameter in `Options` or API subjects.
* Thin-client spec for second-client cancellation (`cancel`/`kill` from non-owner) not found — `DisconnectClientByID`/`LDMClientByID` (`server/server.go:4760-4799`) are server-internal privileged ops, not exposed over the client protocol.
* No protocol-compatibility test covering major `PROTO` bump beyond `ClientProtoInfo` — `client_test.go:556` only asserts `protocol==1` parse.

---

Generated by `01.13 Daemon, Local API, and Multi-Client Lifetime` against `nats-server`.
