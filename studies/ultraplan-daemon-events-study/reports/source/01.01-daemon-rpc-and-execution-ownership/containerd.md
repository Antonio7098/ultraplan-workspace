# Source Analysis: containerd

## Daemon RPC and Execution Ownership

### Source Info

| Field | Value |
|-------|-------|
| Name | containerd |
| Path | `studies/ultraplan-daemon-events-study/sources/containerd` |
| Language / Stack | Go / gRPC + ttrpc over Unix socket (named pipes on Windows), plugin registry |
| Analyzed | 2026-09-02 |

## Summary

containerd implements a strict daemon-owned execution model: a long-lived `containerd` daemon process owns all containers, snapshots, content and tasks after the submitting client disconnects. Clients (`ctr`, `containerd client` library, CRI, Kubernetes) are ephemeral gRPC stubs that dial a Unix domain socket (`/run/containerd/containerd.sock`) via `pkg/dialer`. The wire surface is explicit protobuf-defined services (`api/services/*/v1/*.proto`) registered as `io.containerd.grpc.v1` plugins and served by two server plugins `io.containerd.server.v1.grpc` (unix) and `io.containerd.server.v1.grpc-tcp` (optional TLS TCP) plus `io.containerd.server.v1.ttrpc` for shim-local RPC. Execution outlives the client because metadata is persisted in `core/metadata` (bolt) and running workloads are re-parented to shims (`containerd-shim-runc-v2`); clients reconnect via `Client.Reconnect()` or by re-listing persisted objects (`LoadContainer`/`Task`). Authentication is filesystem-permission only (socket `0660` + `chown uid/gid`), no mTLS on unix, optional TLS for TCP. Workspace scoping is via explicit `containerd-namespace` gRPC/ttrpc metadata header, never via daemon CWD. Protocol versioning is minimal: config version migration (`version.ConfigVersion=4`) and `Version` service, no per-RPC negotiation.

## Rating

**7 / 10** — Clear model with tests, explicit interfaces, and operational safeguards

**Rationale:** Daemon ownership after disconnect is unambiguous and proven at scale (shims + persistent store, `Wait` and `Reconnect` semantics). Transport (unix vs npipe vs TCP+TLS) is explicitly configured and `GetLocalListener` enforces ownership. Namespace isolation is per-RPC metadata via interceptors. Weaknesses keep it from 9-10: (a) local caller authentication relies solely on filesystem permissions/`insecure` creds with no peer-credential (`SO_PEERCRED`) or token check, (b) no protocol/ABI negotiation — clients discover server version only via out-of-band `Version` call, (c) daemon CWD independence is intentional but shim still uses `os.Getwd` for socket placement, (d) shutdown is fatal (`log.Fatal` on serve failure) without graceful drain or request hedging.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Daemon entrypoint | `main()` delegates to `command.App().Run` | `cmd/containerd/main.go:29` |
| Daemon lifecycle | `App.Action` creates context, loads config, `server.New`, `server.Start`, handles signals, `server.Stop` reverses plugin close order | `cmd/containerd/command/main.go:132-281`, `cmd/containerd/server/server.go:272-298` |
| Server plugin init — gRPC socket | `sys.GetLocalListener(address,uid,gid)` + `grpc.NewServer` with namespace interceptors; address default `/run/containerd/containerd.sock` | `plugins/server/grpc/plugin.go:66-139`, `plugins/server/grpc/plugin.go:264-273` |
| Server plugin init — TTRPC | Mirrors gRPC but on `address+".ttrpc"` , `ttrpc.NewServer`, `sys.GetLocalListener` | `plugins/server/ttrpc/plugin.go:48-113`, `plugins/server/ttrpc/plugin.go:113-122` |
| Server plugin init — TCP optional TLS | `net.Listen("tcp", address)` + optional `tls.LoadX509KeyPair`, `x509.NewCertPool`, `credentials.NewTLS` | `plugins/server/grpc/plugin.go:142-304` |
| Top-level directory ownership | `CreateTopLevelDirectories` validates `Root != State`, `MkdirAllWithACL 0700/0711`, `os.Chmod`; `applyFlags` forces `filepath.Abs` on `root/state` | `cmd/containerd/server/server.go:60-98`, `cmd/containerd/command/main.go:311-315` |
| Default addresses | `DefaultAddress="/run/containerd/containerd.sock"`; `DefaultFIFODir="/run/containerd/fifo"`; TTRPC default = address+`.ttrpc` | `defaults/defaults_linux.go:21`, `defaults/defaults_linux.go:26`, `cmd/containerd/server/server.go:155`, `plugins/server/ttrpc/plugin.go:50` |
| Unix transport — dialer | `net.DialTimeout("unix", stripped, timeout)` with `DialAddress="unix://"`; retry on `ENOENT` every 10ms | `pkg/dialer/dialer_unix.go:32-43`, `pkg/dialer/dialer.go:41-60` |
| Windows transport — named pipe | `winio.DialPipe` with `npipe://` prefix | `pkg/dialer/dialer_windows.go:34-46` |
| Socket creation + perms | `CreateUnixSocket` does `unix.Unlink` + `net.Listen("unix")`; `GetLocalListener` does `os.Chmod 0660` + `os.Chown uid/gid` + `mkdirAs` with `Chown` | `pkg/sys/socket_unix.go:31-42`, `pkg/sys/socket_unix.go:46-68` |
| Named pipe creation | `GetLocalListener` via `winio.ListenPipe` | `pkg/sys/socket_windows.go:25-42` |
| Shared serve helper | `internal.Serve` launches `serveFunc` in goroutine, logs and `Fatal` on error, closes listener | `plugins/server/internal/serve.go:29-38` |
| RPC definitions — Tasks | `service Tasks { rpc Create/Start/Delete/Kill/Wait ... }` 20 RPCs | `api/services/tasks/v1/tasks.proto:31-69` |
| RPC definitions — other services | Containers, Images, Snapshots, Content, Namespaces, Events, Diff, Leases, etc. all protobuf | `api/services/containers/v1/containers.proto:1`, `api/services/images/v1/images.proto:1`, `api/services/events/v1/events.proto:1` |
| Client construction | `client.New(address, opts)` builds `grpc.ClientConn` with `insecure.NewCredentials()`, `dialer.ContextDialer`, `grpc.MaxCall*`, namespace interceptors, stored `connector` for reconnect | `client/client.go:104-182` |
| Client namespace interceptors | `newNSInterceptors(defaultns)` injects `containerd-namespace` header via `metadata.Pairs`; overrides `WithDialOpts` can replace | `client/client.go:159-163` (and `pkg/namespaces/grpc.go:32`) |
| Server namespace interceptors | `unaryNamespaceInterceptor` / `streamNamespaceInterceptor` read `fromGRPCHeader`/`fromTTRPCHeader` and push into `context.WithValue` via `WithNamespace` | `plugins/server/grpc/namespace.go:26-42`, `pkg/namespaces/context.go:38-43`, `pkg/namespaces/ttrpc.go:30-42` |
| Health / IsServing probe | `HealthService().Check(..., WaitForReady:true)` exposes `grpc.health.v1.Health` | `client/client.go:311-322`, `plugins/services/healthcheck/service.go:35-50` |
| Version service | Unary `Version` returns `version.Version/Revision`; client wrapper `Version(ctx, Empty)` | `api/services/version/v1/version.proto:26-28`, `plugins/services/version/service.go:54-58`, `client/client.go:869-883` |
| Task ownership — client exit does NOT kill | `task.Wait` blocks on `TaskService.Wait` server-streaming independent of creator; `task.Delete` requires explicit `Stopped`; shim survives daemon restart | `client/task.go:342-367`, `client/task.go:372-441`, `plugins/services/tasks/local.go:414` |
| Task detach / FIFO reattach | `Container.NewTask` creates FIFO set (`cio.Creator`), passes stdin/stdout/stderr paths to `CreateTaskRequest`; `loadTask` + `attachExistingIO` reattaches to server-reported FIFO paths | `client/container.go:226-306`, `client/container.go:479-552` |
| Task Wait semantics | Server `WaitResponse {exit_status, exited_at}` ; client returns channel `<-chan ExitStatus` that completes even if original caller died | `api/services/tasks/v1/tasks.proto:228-236`, `client/task.go:342-367` |
| Reconnect | `Client.Reconnect()` closes old `grpc.ClientConn` and re-dials via stored `connector`; TTRPC `ttrpcutil.Client.Reconnect()` mirrors | `client/client.go:237-251`, `pkg/ttrpcutil/client.go:62-88` |
| Reacquisition after disconnect | `LoadContainer(id)` → `ContainerService.Get`; `LoadContainer(...).Task(ctx, cio.Attach)` → `TaskService.Get` + FIFO reattach; `Containers`/`List` | `client/client.go:379-397`, `client/container.go:479-506`, `client/client.go:325-335` |
| Lease-based work ownership | `WithLease` context wrappers for pull/push; leases persist across client exits until explicit `Delete` or GC | `client/client.go:342-348` (`WithLease` in `NewContainer`), `plugins/services/leases/service.go:38` |
| Authentication — socket ownership | `sys.GetLocalListener` enforces `uid/gid` from config (`server.config.Address UID/GID` defaults to `os.Geteuid/Getegid`); no `SO_PEERCRED` check | `pkg/sys/socket_unix.go:62-64`, `plugins/server/grpc/plugin.go:69-74` |
| Authentication — insecure creds | Both client and `proxyClients` use `insecure.NewCredentials()` — no auth handshake on unix | `client/client.go:144`, `cmd/containerd/server/server.go:410` |
| Authentication — TCP TLS | Optional `tls_ca/tls_cert/tls_key` or Windows store `tls_common_name` via `wintls.SetupTLS...` -> `grpc.Creds` | `plugins/server/grpc/plugin.go:186-218` |
| Config version migration | `ConfigVersion=4`, `MigrateConfigTo` runs per-version `migrations` and plugin `ConfigMigration`; addresses moved from top-level `GRPC/TTRPC` fields to `Plugins["io.containerd.server.v1.grpc"]` | `version/version.go:41`, `cmd/containerd/server/config/config.go:49-54`, `cmd/containerd/server/config/config.go:228-339` |
| Shutdown | `handleSignals` listens for `SIGTERM/SIGINT/SIGUSR1/SIGPIPE`; on term cancels context, `server.Stop()` iterates plugins reverse `io.Closer` | `cmd/containerd/command/main_unix.go:30-35`, `cmd/containerd/command/main_unix.go:37-72`, `cmd/containerd/server/server.go:282-298` |
| CWD dependence audit | Daemon `applyFlags` forces `filepath.Abs`; shim `serve` incorrectly uses `os.Getwd()` to locate ttrpc listener — shows shim (not daemon) depends on CWD | `cmd/containerd/command/main.go:311-315`, `pkg/shim/shim.go:460` |
| Event subscription reconnect semantics | Client `Subscribe` loops `session.Recv()` until `err` or `ctx.Done()`; no auto-resubscribe on disconnect | `client/events.go:80-124` |
| Namespace store — default scoping | `NamespaceEnvVar=CONTAINERD_NAMESPACE`, default `"default"`; `NamespaceRequired` validates with `identifiers.Validate` | `pkg/namespaces/context.go:29-33`, `pkg/namespaces/context.go:68-78` |
| Reconnect test | `TestClientReconnect` and `TestClientTTRPC_Reconnect` prove new client can recover after `Reconnect()` | `integration/client/client_test.go:499-531`, `integration/client/client_ttrpc_test.go:45-67` |
| Restart / shim reattach test | `TestDaemonReconnectsToShimIOPipesOnRestart` + `issue7496_shutdown` demonstrates workload survives daemon bounce | `integration/client/container_linux_test.go:401`, `integration/issue7496_shutdown_linux_test.go:105-113` |

## Answers to Dimension Questions

### 1. Does a client exit affect execution?

**No — execution is daemon + shim owned.** A client crash/exit never kills containers or tasks. `ctr` and `client.Client` are ephemeral gRPC callers (`client/client.go:104` constructs a `grpc.ClientConn` and never holds task state). After `Container.NewTask` (`client/container.go:226`), the request `CreateTaskRequest` (`api/services/tasks/v1/tasks.proto:71`) is handed to the daemon's `Tasks` service (`plugins/services/tasks/local.go:174`); the daemon stores container metadata in `core/metadata` (boltdb) and delegates runtime to a per-container shim process (`containerd-shim-runc-v2`). The `Task.Wait` RPC (`client/task.go:342`) is server-side blocking; killing the waiting goroutine only cancels the client stream, not the shim. `Task.Delete` is explicit and fails if status ≠ `Stopped` (`client/task.go:393-408`). Integration proof: `integration/issue7496_shutdown_linux_test.go:105` and `integration/client/container_linux_test.go:401` (`TestDaemonReconnectsToShimIOPipesOnRestart`) verify containers survive daemon restarts and client disconnects. The only work tied to a live client is an in-flight streaming RPC (e.g., `Subscribe`/`Wait`) — the caller gets `context.Canceled` but the server-side task continues.

### 2. What exactly crosses the daemon boundary?

**Protobuf-defined RPCs plus mount descriptors and FIFO path strings — not host CWD or filesystem handles.**

- **Wire surface:** All gRPC services under `api/services/*/v1/*.proto` (`tasks.proto:31`, `containers.proto`, `snapshots.proto`, `content.proto`, etc.) and TTRPC `events/v1`. Smallest stable surface is `Version`, `Health`, `Containers/Create/Get/List/Delete`, `Tasks/Create/Start/Wait/Delete`, `Content/Write`, `Snapshots/Prepare`. Each call carries a `containerd-namespace` header (`pkg/namespaces/grpc.go:27` `GRPCHeader`, `pkg/namespaces/ttrpc.go:27` `TTRPCHeader`) injected by client interceptors (`client/client.go:159`) and extracted by server interceptors (`plugins/server/grpc/namespace.go:26`).
- **Payloads crossing:** Container spec (`oci.Spec` JSON via `typeurl.Any`), `Mount` array (`types/mount.proto`) resolved from snapshots (`client/container.go:334-367` pushes `snapshotter.Mounts` and appends `MountLabel`), FIFO path triple (`Stdin/Stdout/Stderr` strings, `api/services/tasks/v1/tasks.proto:82-85` and `client/container.go:241`), checkpoint descriptors (`CheckpointTaskRequest`), lease IDs via context.
- **What does NOT cross:** Client CWD, env vars, shell, `uid` of caller beyond socket ownership, or TLS certificate unless TCP is enabled. `ctr` explicitly checks socket existence before dial (`cmd/ctr/commands/client.go:64`). Image content goes via streaming or content store, not inline.

### 3. Can a client reconnect and recover the same operation?

**Yes, for durable objects; no, for in-flight unary streams.**

- **Durable recovery:** Any persisted resource can be reacquired by a new client instance: `Client.Reconnect()` (`client/client.go:237`) closes and re-creates the `grpc.ClientConn` via stored `connector` (`client/client.go:165-176`, tested `integration/client/client_test.go:499`). `TTRPC` mirrors (`pkg/ttrpcutil/client.go:62`). Recovery idiom is `LoadContainer(ctx, id)` (`client/client.go:380`) → `Task(ctx, attach)` (`client/container.go:206`) → `Status/Wait/Pids` (`client/task.go:324, 342`). Because metadata lives in the daemon, a different `ctr` invocation with same `--namespace` and `--address` sees identical `List` results.
- **Non-recoverable:** An in-flight `Wait`, `Events.Subscribe`, or `Pull` streaming RPC is per-connection; canceling the context tears down the stream (`client/events.go:115-119`, `client/task.go:346-364`). The client must re-issue `Wait` after reconnect; `Subscribe` has no cursor replay (`client/events.go:80-124` just loops `Recv`). Lease expiry may abort a `Fetch`/`Push` that was anchoring content.
- **Evidence of full reconnect:** `TestClientTTRPC_Reconnect` (`integration/client/client_ttrpc_test.go:45`), `TestClientReconnect` (`integration/client/client_test.go:519`).

### 4. How is the local caller authenticated?

**Filesystem-only; no in-protocol authentication.**

- **Unix (Linux/darwin):** `sys.GetLocalListener` (`pkg/sys/socket_unix.go:46`) creates parent dirs with `mkdirAs(uid,gid)` (`pkg/sys/socket_unix.go:70`), then `os.Chmod 0660` and `os.Chown uid/gid` (`pkg/sys/socket_unix.go:57-64`), where `uid/gid` default to `os.Geteuid/Getegid` (`plugins/server/grpc/plugin.go:71-72`, `plugins/server/ttrpc/plugin.go:60-61`). Access control is entirely the Unix discretionary permission bit; any local user in the owning group can connect. No `SO_PEERCRED`/`SO_PEERSEC` check, no token, and both sides use `insecure.NewCredentials()` (`client/client.go:144`, `cmd/containerd/server/server.go:410`). `DialAddress` is just `unix://path` (`pkg/dialer/dialer_unix.go:32`).
- **Windows:** Named pipe via `winio.DialPipe`/`ListenPipe` (`pkg/dialer/dialer_windows.go:34`, `pkg/sys/socket_windows.go:25`) with ACLs handled by winio.
- **TCP (`grpc-tcp` plugin, off by default):** Optional TLS: `tls_cert/tls_key` + `tls_ca` for mutual TLS (`plugins/server/grpc/plugin.go:186-204`) or Windows cert store (`wintls.SetupTLSFromWindowsCertStore` `plugins/server/grpc/plugin.go:206-212`). When enabled, `grpc.Creds(credentials.NewTLS(...))` enforces client cert verification. Without it, TCP is plaintext with no auth.
- **Versioned config:** No per-RPC authz; namespace `Delete` checks `errdefs` but any authenticated peer can call any RPC within its namespace.

### 5. Does daemon code ever depend on its own current working directory?

**Intentionally no for the main daemon; yes for legacy shim and tooling, which is a documented edge.**

- **Daemon avoids CWD:** `applyFlags` forces absolute `Root/State` via `filepath.Abs` (`cmd/containerd/command/main.go:311-315`). `CreateTopLevelDirectories` joins `Root/State` without `Getwd` (`cmd/containerd/server/server.go:60-83`). Default config dir is absolute `/etc/containerd` / `/run/containerd` (`defaults/defaults_unix.go:23`, `defaults/defaults_linux.go:21`). `pkg/sys/socket_unix.go:37` does `MkdirAll(filepath.Dir(path))` — relative socket paths would resolve against CWD, but daemon validates emptiness (`plugins/server/grpc/plugin.go:78`) and defaults are absolute, so CWD is effectively ignored.
- **Exception — shim:** `pkg/shim/shim.go:460` (`path, _ := os.Getwd()`) then `serveListener(socketFlag, 3)` binds ttrpc in the current working directory. `cmd/containerd-shim-runc-v2/manager/manager_linux.go:93,301` also calls `os.Getwd()` to resolve bundle paths. This is isolated to shim execution, not the central daemon dispatch.
- **Non-daemon helpers:** Tests and vendor shims occasionally call `os.Getwd` (`integration/.../container_linux_test.go:401` context), but not for request routing.
- **Implication:** Starting `containerd` from a different working directory does not change its socket location or where it stores state, except that an explicitly relative `--address` flag would be CWD-relative via `net.DialTimeout("unix", path)` unless pre-absolutized (the code does **not** `Abs` the address).

## Architectural Decisions

| Decision | Why / Consequent |
|----------|-----------------|
| **Unix socket as primary transport, gRPC as protocol** (`defaults/defaults_linux.go:21`, `pkg/dialer/dialer_unix.go:40`, `plugins/server/grpc/plugin.go:264`) | Leverages OS DAC for auth, zero network exposure, high-throughput protobuf streaming, reuse of `health`, `otel`, `prometheus` interceptors (`plugins/server/grpc/plugin.go:88-99`). |
| **Dual API plane: gRPC (control) + TTRPC (shim-local)** (`plugins/server/ttrpc/plugin.go:51-55`, `api/services/ttrpc/events/v1/events.proto:1`) | TTRPC is lower overhead for shim↔daemon and sandbox internal communication; shares namespace header (`pkg/namespaces/ttrpc.go:27`). Split requires separate listeners (`.sock` + `.sock.ttrpc`) and separate `Reconnect` clients (`pkg/ttrpcutil/client.go:62`). |
| **Plugin registry for server and services** (`plugins/types.go:82-83` `ServerPlugin`, `plugins/server/grpc/plugin.go:62` `Requires: [GRPCPlugin, MetricsPlugin]`) | Allows disabling/requiring plugins (`Config.DisabledPlugins/RequiredPlugins` `cmd/containerd/server/config/config.go:71`), lazy init ordering via `plugin.Set` (`cmd/containerd/server/server.go:143-221`), and proxy plugins over gRPC (`cmd/containerd/server/server.go:315-382`). |
| **Namespace as metadata header + context value** (`pkg/namespaces/grpc.go:27`, `pkg/namespaces/context.go:38`, `plugins/server/grpc/namespace.go:26`) | Enables single daemon to multiplex tenants without separate sockets; `WithNamespace` injects into both context and outgoing `metadata` so server interceptors can restore (`fromGRPCHeader`). |
| **Explicit detach: FIFOs + shim, not client stdio** (`client/container.go:241-246`, `pkg/sys/socket_unix.go:31` vs shim bundle `pkg/shim/shim.go:454`) | Client passes FIFO path strings; daemon/shim opens them. Enables `Task(ctx, attach)` to reattach (`client/container.go:513`), detach semantics identical to `runc` detach, and survivor after `ctr` exit. |
| **Config version migration (v1→v4)** (`version/version.go:41`, `cmd/containerd/server/config/config.go:49-54`, `228-338`) | Field `GRPC/TTRPC` top-level migrated into `Plugins["io.containerd.server.v1.*"]` to keep TOML backward compat; prevents breaking changes while allowing address/uid per-listener. |
| **Reverse-close on shutdown** (`cmd/containerd/server/server.go:282-298` loop `len(plugins)-1..0`) | Guarantees LIFO teardown (servers before stores) to avoid use-after-close on `bolt` metadata. |
| **Fatal on serve error** (`plugins/server/internal/serve.go:35` `log.Fatal`) | Treats listener failure as unrecoverable; simplifies failure mode but trades gracefulness for crash-loop recovery by systemd. |

## Notable Patterns

- **Interceptor pattern for cross-cutting concerns:** `streamNamespaceInterceptor`/`unaryNamespaceInterceptor` (`plugins/server/grpc/namespace.go:26-34`) wrap all services without service code awareness; prometheus + otel interceptors appended via `Chain*Interceptor` (`plugins/server/grpc/plugin.go:82-100`).
- **Connector closure for transparent reconnection:** `client.Client.connector func() (*grpc.ClientConn, error)` (`client/client.go:227`) and `ttrpcutil.Client.connector` (`pkg/ttrpcutil/client.go:33`) capture dial opts including `dialer.ContextDialer`; `Reconnect` just re-invokes it (`client/client.go:243`).
- **Dial retry on ENOENT:** `pkg/dialer/dialer.go:49-57` busy-polls every 10ms when socket not yet present — handles daemon cold start race without exposing retry to callers.
- **Lazy TTRPC client:** `ttrpcutil.Client.Client()` (`pkg/ttrpcutil/client.go:100`) lazily dials only on first RPC, keeping idle process cheap.
- **Proxy plugin delegation:** `cmd/containerd/server/server.go:315-386` registers `proxyClients` that multiplex snapshot/content over gRPC to remote addresses — enables nydus/stargz external snapshotter without daemon changes.
- **Absolute-path normalization:** `cmd/containerd/command/main.go:311` + `filepath.Abs` for `Root/State` is defensive against CWD confusion.

## Tradeoffs

| Tradeoff | Pro | Con |
|----------|-----|-----|
| **Insecure creds on unix vs OS perms** (`insecure.NewCredentials` + `0660/chown`) | No CA/bootstrap friction, suitable for single-host container runtime; fast path | No defense-in-depth; any process with `containerd` group can invoke privileged ops (`snapshots.Prepare`, `tasks.Create` with arbitrary rootfs mounts); missing `SO_PEERCRED` auditability |
| **Unix + separate TTRPC sockets vs single socket with mTLS** | Minimal handshake, back-compat with `DefaultAddress+".ttrpc"` derivation (`cmd/containerd/server/config/config.go:306-313`) | Doubles file count, `filepath.Glob` path length 104 limit (`pkg/sys/socket_unix.go:33`); TTRPC socket inherits GID silently |
| **Protobuf+gRPC max message 16 MiB** (`defaults/defaults.go:22-23`) | Prevents OOM on malicious large messages | Large `List` (images/containers) truncates; streaming required but not all services streaming |
| **Explicit namespace header vs socket-per-namespace** | Lightweight multitenancy, `ctr --namespace k8s.io` is just a flag (`cmd/ctr/app/main.go:122`) | Easy to forget `WithNamespace`; `NamespaceRequired` fails closed (`pkg/namespaces/context.go:68`) but error surfaces only at RPC call, not dial |
| **Plugin init short-circuit on failure** (`plugins/server/grpc/plugin.go:122-125` `continue` on `Err`) | CRI failing under rootless doesn’t crash entire daemon | Operational opacity: failed plugin silently excluded from registration, only warned |
| **Dial retry on ENOENT 10ms poll** (`pkg/dialer/dialer.go:49-54`) | Robust to systemd socket-activation ordering | Hot loop under tight timeout; `ContextDialer` timeout must be non-zero or retry forever until success |
| **Fatal on serve failure** (`plugins/server/internal/serve.go:35`) | Ensures systemd restarts fresh instance | No graceful drain: in-flight `Wait/Exec` streams killed with `Fatal` rather than error return |

## Failure Modes / Edge Cases

- **Socket path >104 chars** (`pkg/sys/socket_unix.go:33` returns `fmt.Errorf` → daemon fails to start). Mitigation only via test, no truncation or abstract socket.
- **Stale socket left after crash:** `CreateUnixSocket` does `unix.Unlink` before `Listen` (`pkg/sys/socket_unix.go:39`) so restart reclaims, but unclean `chown/chmod` race can briefly expose world-readable socket.
- **Client dial before daemon creates socket:** `dialer.go:52-54` poll retries hide `ENOENT`, but caller-visible timeout (`Client.New` default 10s `client/client.go:112`) may still fire before daemon ready; no readiness polling beyond retry.
- **Reconnect during in-flight streaming:** `Subscribe`/`Wait` channel returns `err` without auto-resubscribe (`client/events.go:102-104`); caller must recreate stream. Lease holding content may expire, turning `ContentStore.Writer.Commit` into `AlreadyExists`/`NotFound`.
- **Namespace forgotten:** Context without namespace passes server interceptors as `ok==false`, handler then returns `FailedPrecondition` (`pkg/namespaces/context.go:68`). `ctr` default forces `"default"` (`cmd/ctr/app/main.go:125`) but programmatic clients that omit `WithNamespace` fail silently on authz.
- **Group permission mismatch:** If daemon started as root with `uid 0/gid 0` but client runs rootless, `0660` blocks connection with `permission denied` at dial; error surfaced as `cannot access socket` (`cmd/ctr/commands/client.go:64`).
- **TCP TLS misconfig:** Absent file for `tls_ca`/`tls_cert` causes `InitFn` error → plugin skipped → `grpc-tcp` server never starts, no fallback to plaintext — silent degradation if log level low.
- **CWD-relative `--address` flag:** `client/client.go:167` `DialAddress(address)` doesn’t `Abs`; starting client from different dir dials different path. Daemon similarly not absolutizing address (only `Root/State`).
- **Shim CWD dependence:** `pkg/shim/shim.go:460` `Getwd` can fail if shim started with deleted CWD (e.g., `chdir` race `core/mount/mount_linux.go:537`), returning error to `serve` and failing task start.

## Future Considerations

- **Peer-credential auth for unix:** Add `Peer` interceptor using `SO_PEERCRED`/`SCM_CREDENTIALS` (`syscall.GetsockoptUcred` on `net.UnixConn`) and map `uid` to plugin-level ACL; would align with socket `0660` and enable per-namespace RBAC without TLS.
- **OIDC / token for TCP:** Promote `grpc-tcp` to always-TLS with client token header validator (e.g., `authorization: Bearer`) for remote daemon access; current optional TLS invites plaintext exposure.
- **Graceful draining:** Replace `log.Fatal` (`plugins/server/internal/serve.go:35`) with `Shutdown` that drains `Wait`/`Metrics` streams, stops accepting, then closes. Expose `Readiness` already tracked by `Server.ready` (`cmd/containerd/server/server.go:301-309`) for systemd `notifyReady`.
- **Version/Protocol negotiation:** Extend `Version` service to include `SupportedAPIVersions` list and client feature bits; `client.New` could assert `ConfigVersion` and warn on mismatch rather than relying on implicit config migration.
- **Streaming resumption cursor:** Add sequence ID to `events.proto` and `SubscribeRequest` `since` field so reconnect can replay missed events after `Reconnect()`.
- **Socket activation / systemd `ListenFD`:** Adopt `SD_LISTEN_FDS_START` to avoid filesystem socket and 104-char limit, unify with `GetLocalListener` abstraction.
- **Absolute-address enforcement:** `filepath.Abs` the `address` flag in both `ApplyFlags` and client `New` to eliminate CWD edge case; emit deprecation for relative addresses.

## Questions / Gaps

- No evidence found of client certificate rotation or auth token refresh for TCP listeners — search of `plugins/server/grpc/*.go:1-304`, `pkg/dialer/*.go`, and `client/*.go` found only static `LoadX509KeyPair` at init; question remains how TLS rotation would happen without restart.
- No explicit “detach token” or request ID persisted across daemon restart beyond container/task ID — reconciling a `Pull` that was mid-transfer after daemon bounce is unclear; `leases` and `transfer` services exist (`api/services/transfer/v1/transfer.proto`) but their durability guarantees across `client.Close` were not traced.
- Namespace ACL not evaluated — `pkg/namespaces/store.go` interface has `List/Delete` but not permission bits; confirm authorization is out-of-scope (handled by CRI calls rather than containerd).
- No probe of maximum gRPC stream count or backpressure — `chain` interceptors have no rate-limiting; concurrency limits under load (`MaxConcurrentDownloads`) are client-side (`client/client_opts.go:254-259`) not server-enforced.
- `cmd/containerd/command/main.go:219-241` runs `server.New` in goroutine with channel rendezvous; unclear timeout semantics if plugin init hangs (e.g., `bolt` db locked) — `handleSignals` Cancel would unblock but no timer forces `SIGKILL` avoidance.

---

Generated by `dimension 01.01` against `containerd`.
