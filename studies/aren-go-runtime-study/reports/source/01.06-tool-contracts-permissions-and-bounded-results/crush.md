# Source Analysis: crush

## 01.06 Tool Contracts, Permissions, and Bounded Results

### Source Info

| Field | Value |
|-------|-------|
| Name | crush |
| Path | `studies/aren-go-runtime-study/sources/crush` |
| Language / Stack | Go (module `github.com/charmbracelet/crush`, `charm.land/fantasy v0.41.3`, `modelcontextprotocol/go-sdk`) |
| Analyzed | 2026-08-29 |

## Summary

Crush implements a typed tool registry on top of `fantasy.AgentTool[TParams]`. `coordinator.buildTools` (`internal/agent/coordinator.go:679`) aggregates static built-ins, conditional LSP tools, and dynamic MCP tools (`internal/agent/tools/mcp-tools.go:24`), filters by `agent.AllowedTools` / `AllowedMCP`, sorts, then wraps with `hookedTool` (`internal/agent/hooked_tool.go:31`). Schemas derive from Go struct `json`/`description` tags and from MCP `InputSchema` passthrough (`internal/agent/tools/mcp-tools.go:70`). Calls flow: `fantasy` parses provider tool-calls → `sanitizeToolInput` validates JSON (`internal/agent/agent.go:2276`) → `hookedTool.Run` fires `PreToolUse` hooks (`internal/agent/hooked_tool.go:54`) → inner tool validates params → `permission.Service.Request` (`internal/permission/permission.go:181`) blocks before side effects → execution → typed `fantasy.ToolResponse` (text/error/image/media) with metadata → `convertToToolResult` preserves `IsError`/`Data`/`MIMEType` (`internal/agent/agent.go:2091`). Result bounding is explicit per-tool: byte/line caps plus visible truncation markers and pagination rather than silent discard; pending large filesystem/network results are streamed or backgrounded. MCP lifecycle is state-machine driven with generation fencing but relies on external server processes with no additional sandbox.

## Rating

**6 / 10**

Rationale: typed contracts, pre-effect permission, and disciplined per-tool bounding with truncated flags are consistently implemented. Hooks correctly precede permission which precedes mutation/network. Media/binary is first-class, and malformed JSON is sanitized, not panicked. Deductions: no global secret redaction (only `crush_logs`), no duplicate-name guard beyond namespace prefix, unknown-tool handling delegated to `fantasy` without crush-side evidence, and `download` has no size bound (only timeout), leaving a gap for large-artifact externalization.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Discovery: static built-in assembly | `buildTools` appends bash/crush_info/crush_logs/job/download/edit/multiedit/fetch/glob/grep/ls/todos/view/write, conditional LSP set, MCP resource tools | `internal/agent/coordinator.go:714-759` |
| Discovery: filtering by agent policy | `slices.Contains(agent.AllowedTools, tool.Info().Name)` single allow-list; MCP split `AllowedMCP == nil/0` logic | `internal/agent/coordinator.go:761-790` |
| Discovery: hook wrapping before permission | `wrapToolsWithHooks(tools, hookRunner, isSubAgent)` returns decorated slice; sub-agents skip hooks | `internal/agent/coordinator.go:795-801` |
| Schema: typed Go params → JSON schema | `BashParams` `json:"command" description:"..."`, `ViewParams`, `FetchParams`, etc.; `fantasy.NewAgentTool[TParams]` / `NewParallelAgentTool` drives schema | `internal/agent/tools/bash.go:24-30`, `internal/agent/tools/view.go:47-51`, `internal/agent/tools/fetch.go:45-67` |
| Schema: MCP InputSchema passthrough | `InputSchema` (`map[string]any`) → `ToolInfo.Parameters`/`Required` via props/required extraction | `internal/agent/tools/mcp-tools.go:70-97` |
| Schema: MCP name collision avoidance | `Name() = fmt.Sprintf("mcp_%s_%s", mcpName, tool.Name)` + `filterTools` enabled/disabled lists | `internal/agent/tools/mcp-tools.go:58-60`, `internal/agent/tools/mcp/tools.go:172-195` |
| Call parsing: malformed JSON sanitization | `sanitizeToolInput` replaces non-`json.Valid` input with `{}` and sets `sanitizedToolCalls[callID]=true`; `OnToolResult` injects error text | `internal/agent/agent.go:2276-2287`, `internal/agent/agent.go:967-972` |
| Call parsing: required-field validation | Early `NewTextErrorResponse("file_path is required")` etc. before any I/O | `internal/agent/tools/bash.go:202-204`, `internal/agent/tools/write.go:57-59`, `internal/agent/tools/edit.go:69-71`, `internal/agent/tools/view.go:102-104`, `internal/agent/tools/fetch.go:62-73`, `internal/agent/tools/download.go:69-79`, `internal/agent/tools/grep.go:129-131`, `internal/agent/tools/glob.go:58-60` |
| Permission point: bash (process start) | `permissions.Request(...Action:"execute")` gated on `!isSafeReadOnly`; only after grant does `bgManager.Start(...)` | `internal/agent/tools/bash.go:227-255` |
| Permission point: safe bypass | `safeCommands` + `containsCommandChaining` allows `ls`, `git status` etc. to skip prompt; still blocks chaining | `internal/agent/tools/safe.go:9-75`, `internal/agent/tools/bash.go:209-222` |
| Permission point: write/edit mutation | `permissions.Request(Action:"write")` before `os.WriteFile`/`commitFileChange` | `internal/agent/tools/write.go:108-138`, `internal/agent/tools/edit.go:133-163`, `internal/agent/tools/edit.go:334-372`, `internal/agent/tools/edit.go:407-443` |
| Permission point: view/ls (filesystem read/list) | Outside-workdir check `filepath.Rel == ".."` → `permissions.Request(Action:"read"/"list")` before `os.Stat`/`ListDirectoryTree` | `internal/agent/tools/view.go:116-155`, `internal/agent/tools/ls.go:98-126` |
| Permission point: fetch/download (network) | `permissions.Request(Action:"fetch"/"download")` before `http.NewRequestWithContext`/`client.Do` | `internal/agent/tools/fetch.go:80-97`, `internal/agent/tools/download.go:90-106` |
| Permission point: MCP before invoke | `permissions.Request(Action:"execute")` before `mcp.RunTool`; whitelistDockerTools bypass documented | `internal/agent/tools/mcp-tools.go:99-128` |
| Hook before permission | `hookedTool.Run` executes `h.runner.Run(EventPreToolUse,...)`; `DecisionDeny/Halt` returns without `h.inner.Run`; `DecisionAllow` stamps `WithHookApproval(ctx, call.ID)` | `internal/agent/hooked_tool.go:54-84` |
| Hook approval gate | `permission.Request` short-circuits when `hookApproved(ctx, toolCallID)` matches; still publishes `Granted` notification | `internal/permission/permission.go:192-202` |
| Result bounding: bash truncation | `MaxOutputLength=30000`; `TruncateOutput` keeps head/tail half and inserts `\n\n... [N lines truncated] ...\n\n` | `internal/agent/tools/bash.go:54`, `internal/agent/tools/bash.go:427-442` |
| Result bounding: view byte/line cap | `MaxViewSize=200*1024`; `readTextFile` returns `contentTooLargeError`; `MaxLineLength=2000` truncates `...`; `hasMore` emits `(File has more lines. Use 'offset'...)` | `internal/agent/tools/view.go:76-89`, `internal/agent/tools/view.go:235-256`, `internal/agent/tools/view.go:337-341` |
| Result bounding: fetch | `MaxFetchSize=100*1024` via `io.LimitReader`; tail `[Content truncated to 102400 bytes]` | `internal/agent/tools/fetch.go:22`, `internal/agent/tools/fetch.go:130-185` |
| Result bounding: grep | `searchFiles(..., limit=100)` + `Truncated` bool + `(Results are truncated...)`; fallback regex capped at 200, ripgrep streaming | `internal/agent/tools/grep.go:143-181`, `internal/agent/tools/grep.go:210-215`, `internal/agent/tools/grep.go:349` |
| Result bounding: glob | `globFiles(..., limit=100)` + `Truncated`; `runRipgrep` streams with `candidatePool=max(limit*20,1000)` then sorts by path length | `internal/agent/tools/glob.go:70-88`, `internal/agent/tools/glob.go:124-189` |
| Result bounding: ls | `maxLSFiles=1000` via `fsext.ListDirectory(... maxFiles)` + `There are more than 1000 files...` | `internal/agent/tools/ls.go:52-54`, `internal/agent/tools/ls.go:138-164` |
| Result bounding: crush_logs | `defaultLogLines=50`, `maxLogLines=100`, `maxLogLineSize=1MB` skip; backward chunk reader | `internal/agent/tools/crush_logs.go:42-49`, `internal/agent/tools/crush_logs.go:126-223` |
| Metadata preservation | `WithResponseMetadata` attaches typed metadata: `BashResponseMetadata{StartTime,EndTime,ShellID}`, `ViewResponseMetadata{ResourceType}`, `GrepResponseMetadata{NumberOfMatches,Truncated}` | `internal/agent/tools/bash.go:275-283`, `internal/agent/tools/view.go:261-277`, `internal/agent/tools/grep.go:183-189` |
| Typed failure preservation | `NewTextErrorResponse` sets `IsError=true`; `convertToToolResult` maps `ToolResultContentTypeError` → `baseResult.IsError=true`; media base64 validation warns and converts to error | `internal/agent/tools/tools.go:66-70`, `internal/agent/agent.go:2091-2131` |
| Secret redaction (narrow) | `sensitiveKeys = [authorization,api-key,token,secret,password,credential]` → `[REDACTED]` in log formatting; crush_info hides API key | `internal/agent/tools/crush_logs.go:60-70`, `internal/agent/tools/crush_logs.go:414-439`, `internal/agent/tools/crush_info_test.go:236-243` |
| Binary/media handling | `getImageMimeType`/`sniffImageMimeType` + `http.DetectContentType`; `NewImageResponse`/`NewMediaResponse`; MCP `ensureRawBytes` base64 decode; capability gate `GetSupportsImagesFromContext` | `internal/agent/tools/view.go:204-228`, `internal/agent/tools/view.go:387-400`, `internal/agent/tools/mcp-tools.go:133-148`, `internal/agent/tools/mcp/tools.go:83-94`, `internal/agent/tools/mcp/tools.go:197-248` |
| MCP lifecycle: state machine | `State` enum, `ClientInfo{State,Counts,Config,PendingConfig}`, `updateState` with generation fencing | `internal/agent/tools/mcp/init.go:140-165`, `internal/agent/tools/mcp/init.go:815-870` |
| MCP lifecycle: init gating | `ArmInit`/`WaitForInit`/`initDone` channel; non-interactive runs block, interactive do not | `internal/agent/tools/mcp/init.go:110-114`, `internal/agent/tools/mcp/init.go:322-336`, `internal/agent/coordinator.go:228-246` |
| MCP lifecycle: renewal serialization | `renewMus` per-server mutex, `gens` fencing, `pingSession` health, `getOrRenewClient` double-checked lock | `internal/agent/tools/mcp/init.go:86-91`, `internal/agent/tools/mcp/init.go:126-137`, `internal/agent/tools/mcp/init.go:659-771` |
| Permission service internals | `pendingRequests` map with `Take` race-wins, `Grant/Deny/GrantPersistent` idempotency, `autoApproveSessions`, `allowedTools` allowlist | `internal/permission/permission.go:95-179`, `internal/permission/permission.go:298-310` |

## Answers to Dimension Questions

**Can malformed or ambiguous arguments reach a tool implementation?**

No, for two layers. (1) Provider-level: `sanitizeToolInput` (`internal/agent/agent.go:2276`) checks `json.Valid`; on failure it replaces the entire `Input` string with `"{}"` and records `sanitizedToolCalls[toolCallID]`, preventing a JSON-parse panic inside `fantasy`. The corresponding `OnToolResult` (`internal/agent/agent.go:968`) then surfaces `"Tool call failed: arguments were not valid JSON..."` as an error result instead of executing the real tool. (2) Tool-level: every handler does explicit presence/type validation before any side effect and returns `NewTextErrorResponse` (typed error, not exception) for missing/invalid params: `bash.go:202` (`command == ""`), `write.go:57` (`file_path`), `edit.go:69` (`file_path`), `view.go:102` (`file_path`), `fetch.go:62-73` (URL prefix + format enum), `download.go:69-79`, `grep.go:129` (`pattern == ""`), `glob.go:58`. Ambiguous cases like duplicate `old_string` matches in `edit` are rejected with `"old_string appears multiple times..."` requiring `replace_all` (`internal/agent/tools/edit.go:211`). Remaining ambiguity is `{} `after sanitization reaching a tool that would otherwise accept empty args: it would hit the same `missing X is required` error path, not a silent default.

**Is permission checked before process start, network access, or filesystem mutation?**

Yes — uniformly before side effects, with hook precedence.

* Hook precedes permission: `hookedTool.Run` (`internal/agent/hooked_tool.go:54-86`) executes `PreToolUse` hooks synchronously, returns early on `Deny`/`Halt` without reaching `inner.Run`, and only on success calls `inner.Run` with a context stamped `WithHookApproval` for the exact `call.ID`. `permission.Request` (`internal/permission/permission.go:192`) honors that stamp before any prompt.
* Ordering per tool: `bash` (`internal/agent/tools/bash.go:227-255`) gates `permissions.Request(Action:"execute")` before `bgManager.Start`; the only bypass is `isSafeReadOnly` (allow-listed `safeCommands` without chaining, `internal/agent/tools/safe.go:9-75`) which is intentional. `write`/`edit` gate before `os.WriteFile`/`commitFileChange` (`write.go:108`, `edit.go:133/334/407`). `download`/`fetch` gate before `http.NewRequest`/`client.Do` (`download.go:90`, `fetch.go:80`). `view`/`ls` gate before `os.Stat`/`ListDirectoryTree` only when `Rel == ".."` (outside workdir, `view.go:136`, `ls.go:98`). MCP gates before `mcp.RunTool` (`mcp-tools.go:106`). Tests in `internal/permission/permission_test.go:117-191` prove denied calls never reach execution and that hook approvals are scoped to `toolCallID`.

**How are large results represented without silently discarding evidence?**

Large results are bounded with both byte and count caps and every truncation is made explicit in content and metadata:

* `bash`: centered truncation via `TruncateOutput` (`bash.go:427`) preserves first/last `MaxOutputLength/2` chars and inserts `... [N lines truncated] ...` (N computed from newline counts). Long-running commands are auto-backgrounded after 60s (`DefaultAutoBackgroundAfter`, `bash.go:53`) and their output is retrieved via `job_output` tool, avoiding loss.
* `view`: `MaxViewSize 200KB` enforced per-read segment; oversize returns `contentTooLargeError` as visible error (`view.go:238`); oversized images similarly (`view.go:205`). Lines longer than `MaxLineLength 2000` are rune-safe truncated with `...` (`view.go:337`). When more lines exist, content appends `(File has more lines. Use 'offset' parameter to read beyond line N)` (`view.go:254`) — pagination, not silent discard.
* `fetch`: `io.LimitReader(MaxFetchSize)` (`fetch.go:130`) then content suffix `\n\n[Content truncated to 102400 bytes]` (`fetch.go:184`).
* `grep/glob/ls`: hard count caps (`100`/`100`/`1000`) with boolean `Truncated` in typed response metadata and human marker `(Results are truncated. Consider using a more specific path or pattern.)` (`grep.go:178`, `glob.go:81`, `ls.go:162`). `glob` streams ripgrep output with bounded `candidatePool` (`glob.go:134`) to avoid OOM.
* `crush_logs`: clamps to `maxLogLines=100` and drops lines > `maxLogLineSize 1MB` silently but summarizes count via tail; reports `No log file found`/`Log file is empty` instead of empty.

No tool silently drops bytes; divergence is `download` which streams `io.Copy` without byte cap (`internal/agent/tools/download.go:152`), bounded only by 5 min client timeout — large downloads could exhaust disk without marker.

**Does the abstraction preserve tool-specific semantics or flatten them into strings?**

Partially preserves, but with deliberate flattening for LLM consumption. `fantasy.ToolResponse` is typed: `NewTextResponse` vs `NewTextErrorResponse` (`IsError=true`) vs `NewImageResponse`/`NewMediaResponse` carrying `Data`+`MediaType` (`tools.go:66`, `view.go:227`, `mcp-tools.go:140`). `convertToToolResult` (`agent.go:2091`) maps those to `message.ToolResult{Content, IsError, Data, MIMEType}` preserving error vs success; media base64 is validated (`agent.go:2110`). Metadata types are per-tool and strongly typed: `BashResponseMetadata{StartTime,EndTime,ShellID,WorkingDirectory}`, `ViewResponseMetadata{ResourceType:"skill",ResourceName}`, `Grep/Glob/LSResponseMetadata{Truncated, NumberOfMatches/NumberOfFiles}`. Provider differences are abstracted by `workaroundProviderMediaLimitations` (`agent.go:2132`) which rewrites media tool results into synthetic text+`FilePart` user messages for non-Anthropic providers, preserving semantics across providers. Flattening occurs at the LLM boundary: all text content is concatenated into a single string (with `<file>`/`cwd` wrappers for view/bash) and structured metadata is not fed to the model as separate tool-result fields but as `ProviderOptions`/client metadata — the agent loop sees text first, metadata second.

## Architectural Decisions

* **Generics-typed `fantasy.AgentTool[TParams]` as contract boundary** (`internal/agent/tools/bash.go:198`, `internal/agent/coordinator.go:679`) — schema is code, not external YAML. Tradeoff: strong compile-time typing but schema evolution requires recompilation; MCP schemas must be dynamically mapped (`mcp-tools.go:70`).
* **Hook decorator wraps every tool before permission** (`internal/agent/hooked_tool.go:23`) — satisfies spend-control requirement that policy is decided before governed effect starts. Chose exit-code 49 as `Halt` to avoid collision with normal error ranges.
* **Permission service as pub/sub broker with synchronous `Request`** (`internal/permission/permission.go:181`) blocking on `pendingRequests` channel until UI `Grant/Deny`. Allows WebSocket/MCP subscribers but introduces queue coupling (`requestMu` serializes prompts per permission service, `internal/permission/permission.go:204`).
* **`safeCommands` read-only bypass** (`internal/agent/tools/safe.go:9`) trades friction for UX on common reads; chaining metachars disallow prevents escape (`;` `|` `&&` `$(` `` ` ``).
* **Per-tool bounded result + explicit markers over artifact externalization** — avoids large file management but means truly large outputs must be paginated (view) or backgrounded (bash); no general `artifact://` reference type.
* **MCP as external process pool with generation fencing** (`internal/agent/tools/mcp/init.go:91`, `gens`) — guarantees config-change convergence without stale session registration; tradeoff is extra complexity (per-server mutexes, ping/renew loops).

## Notable Patterns

* **ShallowMerge for hook `updated_input` patches** (`internal/hooks/hooks.go:164`, `internal/agent/hooked_tool.go:75`) — hooks may non-destructively rewrite `tool_input` JSON; validated via `json.Unmarshal` and `tidwall/sjson`.
* **Working-dir confinement with escape hatch** (`view.go:116`, `ls.go:98`, `edit.go:291`) — tracks `LastReadTime` per session/file and rejects edit if file was modified externally or never read, preventing lost updates.
* **Provider-adaptive media workaround** (`internal/agent/agent.go:2132`) — text-only providers get placeholder plus injected `FilePart` user messages; Anthropic/Bedrock keep native `Media` in tool result.
* **Candidate-pool streaming for glob** (`internal/agent/tools/glob.go:124`) to bound memory under `~/` scale workloads rather than buffering full `rg --files` output.
* **Skipped redaction vs full-pipeline sanitization** — sensitive-field redaction only in log rendering (`crush_logs.go:414`) not in arbitrary tool output, indicating trust boundary at logs not at agent loop.

## Tradeoffs

* **Pre-effect gating vs latency:** Every non-safe mutation blocks on human `Grant` interactivity (`permission.go:270` publish then `select` on `ctx.Done`/`respCh`). In yolo/`--no-permissions` (`SkipRequests`, `AutoApproveSession`) this collapses to zero latency but loses audit.
* **Byte-cap truncation vs fidelity:** `bash` 30k char centered truncation keeps signal at both ends but can hide middle of logs; `view` 200KB per page prevents context overflow but forces model to paginate large files manually, increasing turns.
* **Static byte cap (`MaxFetchSize 100KB`) vs content-type awareness:** HTML fetch still capped at 100KB even when markdown conversion would shrink it; `download` has no cap, so a malicious URL could fill disk — no shared quota with fetch.
* **Namespaced MCP prefix vs discoverability:** `mcp_<server>_<tool>` avoids collisions but creates long, brittle names the model must type exactly; no fuzzy aliasing.
* **Malformed JSON sanitization to `{}`** prevents stuck turns but may cause repeated failed tool calls (model retries empty object) rather than failing fast with schema validation error.
* **No generic artifact store:** Keeps architecture simple (no blob service, no GC) at cost of re-reading large files via offset pagination.

## Failure Modes / Edge Cases

* **Unknown tool / provider call identifier drift:** Crush never validates `tool_call_id` uniqueness or registry membership beyond what `fantasy` does; malformed id would propagate as opaque string in `pendingRequests`. No test covers duplicate tool names across fantasy provider translation (tool-name length limits for bedrock/openai not enforced).
* **MCP poison / oversized MCP payload:** `mcp.RunTool` concatenates `TextContent` parts without byte cap (`mcp/tools.go:58`); a rogue server could return megabytes of text which `mcp-tools.go:149` passes as single `NewTextResponse` uncapped, inflating context. Image `ensureRawBytes` double-decodes but does not cap decoded size.
* **Bypassing safe list via chaining obfuscation:** `containsCommandChaining` checks substrings `;|&&$(`` `; missing `||`, `>`, `<`, `$(` variant with space could slip past for non-block-listed commands.
* **`download` unbounded disk write:** `io.Copy` (`download.go:152`) limited only by HTTP timeout; no `LimitReader`, no free-disk check. `fetch` by contrast is bounded.
* **`edit` whitespace fallback hides drift:** `findAndReplace` (`edit.go:201`) silently applies `normalizedReplace` when exact match fails; model sees `whitespaceCorrectedNote` but may not realize newline semantics changed.
* **View UTF-8 gate:** `!utf8.ValidString` returns error (`view.go:244`) so binary files that pass `getImageMimeType=false` return opaque error instead of binary marker.
* **Permission ctx cancellation leaves no audit:** `Request` (`permission.go:272`) returns `ctx.Err()` without publishing a denial notification; hook-before-permission ordering could lose `Deny` signal on client disconnect.
* **Secret leakage outside logs:** Only `crush_logs.go:60` redacts; `bash` output, `fetch` bodies, `grep` matches that contain `API_KEY=sk-...` are returned verbatim to the model and persisted in `message.ToolResult` without redaction.

## Future Considerations

* Add a global byte quota and spill-to-artifact (`artifact://<id>`) path for any tool result exceeding `MaxOutputLength`; return bounded preview plus durable reference instead of centered truncation — closes silent-discard gap for grep/glob/view large outputs.
* Enforce a unified `io.LimitReader` (+ `MaxDownloadSize` default 50MB) in `download.go` mirroring `fetch.go` and surface `[Download truncated]` marker.
* Centralize output sanitization (secret detection via entropy + deny-list) before `convertToToolResult` so even `bash`/`fetch` payloads are redacted, not just log rendering.
* Add registry validation on `SetTools`: reject duplicate `tool.Info().Name` across built-in and MCP set with explicit error (fail fast on bad `crush.json` `enabled_tools`).
* Propagate fantasy provider errors for unknown tool names as typed `IsError` tool results with structured `code=unknown_tool` so model can self-correct vs opaque provider 400.
* Make `sanitizeToolInput` return a typed `ToolParamsError` that maps to a distinct `FinishReason` instead of coalescing to generic text error, improving loop detection.
* Extend `permission.Request` to always publish a `Denied` notification even on `ctx.Done()` for auditable cancellation.

## Questions / Gaps

* No evidence found for explicit **duplicate-tool-name rejection** or **unknown-tool dispatch test** inside `studies/aren-go-runtime-study/sources/crush` (searched `buildTools`, `GetMCPTools`, `wrapToolsWithHooks`, `filterTools`, `coordinator_test.go` — only dedup is `EnabledTools`/`DisabledTools` filtering; unknown-tool path is delegated to `fantasy` with no crush-side unit test).
* No evidence found for **schema mismatch unit tests** (e.g., type-coercion for `BashParams.AutoBackgroundAfter` string vs int) — validation is runtime `json.Unmarshal` inside `fantasy`, not covered by crush `bash_test.go`/`view_test.go`.
* No evidence found for **provider-specific call identifier handling** (Bedrock `toolUseId` vs OpenAI `call_id`) — `fantasy` abstracts this; crush treats `call.ID` as opaque (`permission.go:112`, `mcp-tools.go:112`).
* No evidence found for **truncation-marker E2E test** asserting LLM-visible suffix (`[N lines truncated]`, `[Content truncated ...]`) persists through provider round-trip (only unit tests for `TruncateOutput`, `readTextFile` boundary).
* No evidence found for **binary artifact reference / externalization** (e.g., writing >200KB view slice to temp file and returning `file://` URI) — current design is byte-cap error or pagination.
* No evidence found for **MCP safe-execution proof** (sandbox, capability drop, network egress filter) — lifecycle code proves liveness fencing and OAuth handling (`mcp/init.go:659`) but not isolation; `mcp-tools.go:106` permission still prompts but does not constrain which filesystem paths the MCP child may touch.

---

Generated by `01.06-tool-contracts-permissions-and-bounded-results` against `crush`.
