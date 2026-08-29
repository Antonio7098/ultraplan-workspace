# Source Analysis: temporal

## 16.03 Artifact Provenance and Reproducibility

### Source Info

| Field | Value |
|-------|-------|
| Name | temporal |
| Path | `studies/agent-harness-study/sources/temporal` |
| Language / Stack | Go 1.26.4 (go.temporal.io/server), Cassandra/MySQL/PostgreSQL/SQLite + Elasticsearch, Docker, Goreleaser |
| Analyzed | 2026-08-28 |

## Summary

Temporal traces **workflow artifacts** (execution history + mutable state) to deterministic inputs via event-sourced history (`AppendHistoryNodes`/`ReadHistoryBranch` with `BranchToken`, `PrevTransactionID`/`TransactionID`, `VersionHistory`) and achieves strong execution reproducibility through history replay (documented deterministic workflow requirement + `WithDeterministicProto3` codec). **Build artifacts** (binaries, Docker images) carry minimal provenance — `runtime/debug.ReadBuildInfo` fields `vcs.revision`/`vcs.time`/`GOARCH`/`CGO_ENABLED`, OCI labels `org.opencontainers.image.revision` and `com.temporal.server.version=ServerVersion`, and pinned `go.mod`/`go.sum` + Makefile tool pins — but lack SLSA attestation, SBOM, `SOURCE_DATE_EPOCH`/`-trimpath`, or end-to-end input recording (CLI version, full config, secret refs, marketplace/skill SHAs are not bundled). Reproducibility is best-effort (pinned deps, `CGO_ENABLED=0`, single-source `ALPINE_TAG`) and never asserted in CI — no job rebuilds and `diff`s digests.

## Rating

**5 / 10 — Present but inconsistent, weakly documented, fragile**

Rationale: Workflow-level provenance/reproducibility is a mature, tested core (event sourcing, deterministic replay, `go.sum` + checksummed images) earning 7-8 in isolation; however the dimension's full provenance checklist — *every* artifact traceable to prompts/model calls/tools/inputs/context/approvals, *all* contributing factors recorded, deterministic rebuild verified, and CI-tested — is only partially met for build artifacts and absent for LLM-specific factors (Temporal is not an LLM harness: no prompt/seed/temperature log). Gaps in build attestation, missing hermetic record of config/env, and zero reproducibility regression bring the aggregate to the middle of the 4-6 bucket.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Provenance fields — runtime build info | `common/build/build.go:29-56` reads `debug.ReadBuildInfo()` into `InfoData{GoVersion, GoArch, GoOs, CgoEnabled, GitRevision=vcs.revision, GitTime=vcs.time, GitModified=vcs.modified}`; `Available` false if build info missing. Exposed via `common/build/build.go:9-22` struct tags. | `common/build/build.go:29-56` |
| Provenance fields — static server version | `ServerVersion = "1.32.0"` comment "can be changed by create-tag workflow"; `SupportedServerVersions = ">=1.0.0 <2.0.0"`; `internalVersionHeaderPairs` stamps `client-name=temporal-server`, `client-version=ServerVersion` on every internal gRPC call via `SetVersions` / `NewVersionChecker`. Only bump via tag automation, not per-commit. | `common/headers/version_checker.go:24-60` |
| Provenance fields — OCI image labels | `docker/docker-bake.hcl:57-67` `admin-tools` and `70-93` `server` targets label `org.opencontainers.image.revision="${TEMPORAL_SHA}"`, `org.opencontainers.image.version="${SERVER_VERSION}"`, `com.temporal.server.version`, `org.opencontainers.image.created=timestamp()`, build args `ALPINE_TAG` single source of truth `3.23.4`. Values flow from `.github/actions/build-docker-images`. | `docker/docker-bake.hcl:40-93` |
| Provenance fields — binary distribution checksums | `.goreleaser.yml:88-90` `checksum: name_template: checksums.txt algorithm: sha256`; `88-89` archives each binary with `config/*` included. No SBOM/attestation key configured. | `.goreleaser.yml:88-90` |
| Input recording — workflow history transaction | `AppendHistoryNodesRequest{ShardID, IsNewBranch, BranchToken, Events[], PrevTransactionID, TransactionID}` and `AppendRawHistoryNodesRequest{History DataBlob, NodeID}`; `ForkHistoryBranchRequest{ForkBranchToken, ForkNodeID, NewRunID}` plus `ReadHistoryBranchRequest{BranchToken, MinEventID, MaxEventID, PageSize, NextPageToken}` provide branch-aware, transaction-chained input log. | `common/persistence/data_interfaces.go:776-922` |
| Input recording — persistence execution snapshot | `WorkflowEvents{NamespaceID, WorkflowID, RunID, BranchToken, PrevTxnID, TxnID, Events[]}`; `WorkflowMutation{NextEventID, Condition/DBRecordVersion, Checksum, Tasks}`; `WorkflowSnapshot` mirrors full state for create workflow; events appended atomically with `CreateWorkflowExecutionRequest.NewWorkflowEvents` / `UpdateWorkflowExecutionRequest.UpdateWorkflowEvents`. | `common/persistence/data_interfaces.go:337-402` |
| Input recording — version history item | `versionhistory/version_history_item.go:12` constructs `VersionHistoryItem{EventId, Version}`; used to trace branch divergence and enable deterministic branch replay. | `common/persistence/versionhistory/version_history_item.go:12` |
| Deterministic reproduction — proto codec | `common/persistence/serialization/codec.go:43-91` adds `WithDeterministicProto3() bool deterministic` and `proto.MarshalOptions{Deterministic: opts.deterministic}.Marshal(m)` on `Encode` path; comment "deterministic marshaling". Caller opt-in, not global. | `common/persistence/serialization/codec.go:43-91` |
| Deterministic reproduction — docs contract | `docs/architecture/README.md:35` "Workflow code must be deterministic and have no side effects (with specific exceptions), and activity code must either be idempotent or non-retryable" — architectural promise for history replay determinism. | `docs/architecture/README.md:35` |
| Deterministic reproduction — history rebuilder | `service/history/workflow_rebuilder.go:92,201` `replayResetWorkflow` replays history branch to rebuild mutable state; relies on deterministic history events. | `service/history/workflow_rebuilder.go:92` |
| Build reproducibility — pinned deps | `go.mod:3` `go 1.26.4`; `go.mod:11-89` ~80 direct pins, `go.sum` committed with `h1:` hashes for every module (e.g., `go.sum:483` `go.temporal.io/api v1.63.5 h1:...`); `Makefile:179-276` pins `GOLANGCI_LINT v2.13.0`, `GCI v0.13.6`, `BUF v1.6.0`, `PROTOC_GEN_GO v1.36.6` etc. via `go-install-tool` stamped in `.bin` + `.stamp`. | `go.mod:3` `go.sum:1-10` `Makefile:179-286` |
| Build reproducibility — hermetic flags | `Makefile:37` `CGO_ENABLED ?= 0` + `Makefile:367-392` `go build $(BUILD_TAG_FLAG) -o temporal-server ./cmd/server` with `BUILD_TAG_FLAG=-tags disable_grpc_modules,…`; `.goreleaser.yml:26-79` each build sets `CGO_ENABLED=0` per target. No `-trimpath`, `-buildvcs`, `SOURCE_DATE_EPOCH`, or `ldflags -X` in targets. | `Makefile:37,50-52,366-392` `.goreleaser.yml:26-79` |
| Build reproducibility — Docker base | `docker/docker-bake.hcl:40-42` `variable ALPINE_TAG { default = "3.23.4" }` comment "single source of truth"; `docker/targets/server.Dockerfile:1-3` `FROM alpine:${ALPINE_TAG}` no digest pin (explicit note `docker-bake.hcl:39` justifies no digest for multi-arch buildx). `Makefile:37` and `.goreleaser.yml` ensure linux/amd64,linux/arm64 matrix. | `docker/docker-bake.hcl:40-42` `docker/targets/server.Dockerfile:1-3` |
| Build reproducibility — tool toolchain | `Makefile:221-222` `# NilAway has no tagged releases; pin the pseudo-version for reproducible CI.` `NILAWAY_VER:=v0.0.0-20260717164209-b48ebb193579`; `develop/` + `buf` pinned. `.github/actions/prepare-go/action.yml:39,54` caches keyed on `hashFiles('go.mod')+hashFiles('go.sum')` + `Makefile` for tools. | `Makefile:221-222` `.github/actions/prepare-go/action.yml:54-63` |
| Build/reproduction scripts — Makefile targets | `Makefile:6` `bins: temporal-server temporal-cassandra-tool temporal-sql-tool temporal-elasticsearch-tool tdbg`; `Makefile:313-340` `proto: lint-protos lint-api protoc proto-codegen` regenerates from `proto/internal/**/*.proto`; `Makefile:366-392` artifact matrix + `ensure-no-changes` guard `git status --porcelain` fails if generated files drift. | `Makefile:6,313-392,806-810` |
| Build/reproduction scripts — Goreleaser & Docker | `.github/actions/build-binaries/action.yml:28-56` runs `goreleaser/goreleaser-action@v6.4.0` with `snapshot`/`release` modes; `.github/actions/build-docker-images/action.yml:54-177` compiles `docker-build-helper`, extracts `SERVER_VERSION/CLI_VERSION` from binaries, bakes `docker/docker-bake.hcl` server+admin-tools. | `.github/actions/build-binaries/action.yml:28-56` `.github/actions/build-docker-images/action.yml:54-177` |
| CI reproducibility testing — absent | Grep for `reproducib`, `provenance`, `SLSA`, `SBOM`, `attest`, `SOURCE_DATE`, `trimpath`, `ldflags` yields only `Makefile:221` pin comment, `.goreleaser.yml:88` checksum, `docker-bake.hcl:39` digest-exclusion note; no workflow diffs two builds, no `diffoscope`, no `ensure-reproducible` job. `ensure-no-changes` only checks codegen. | `Makefile:221` `.goreleaser.yml:88` `docker/docker-bake.hcl:39` (negative evidence) |
| Workflow reproducibility tests | `common/persistence/serialization/serializer_test.go:69-115` round-trips `History{Events}`; `common/persistence/persistence-tests/history_v2_persistence.go:159-423` exercises branch append/read/fork; `service/history/workflow/mutable_state_impl_test.go:6903` deterministic tie-break test — tests history fidelity but not cross-build binary reproducibility. | `common/persistence/serialization/serializer_test.go:69-115` `common/persistence/persistence-tests/history_v2_persistence.go:159-223` |

## Answers to Dimension Questions

### 1. Can every artifact be traced to its inputs?

**Partially — workflow artifacts yes, build artifacts only minimally.**

- **Workflow execution artifacts** (the primary Temporal artifact) are fully traceable: every `HistoryEvent` is appended with `BranchToken` + transaction chain (`PrevTransactionID`/`TransactionID`) and versioned via `VersionHistoryItem{EventId, Version}` (`common/persistence/data_interfaces.go:776-848,337-346` `common/persistence/versionhistory/version_history_item.go:12`). `ReadHistoryBranch*` APIs return `HistoryEvents[]` + `NextPageToken` + `Size` enabling audit from first event to current `NextEventID`/`DBRecordVersion`. Namespace/cluster provenance adds `ListClusterMetadata`/`GetClusterMetadataResponse{Version}` (`common/persistence/data_interfaces.go:994-1020`). This is stronger than most harnesses.
- **Build artifacts** (binaries, archives, Docker images) trace only to commit SHA + build env: `common/build/build.go:46-53` records `vcs.revision`/`vcs.time`/`GOARCH`/`GOOS`/`CGO_ENABLED`; Docker labels copy `TEMPORAL_SHA` and `SERVER_VERSION` (`docker/docker-bake.hcl:63,90`); `go.sum` pins every transitive dep hash. Missing: per-binary manifest tying output digest to full input set (CLI version downloaded at build time via `docker-build-helper download-cli` is not hashed in image label), build host/toolchain, config `development-*.yaml` in effect, or which `go.mod` replace directives. No SBOM or SLSA `buildType`/`builder.id`.
- **LLM-specific artifacts** (prompts, model calls, tool approvals) are out of scope for Temporal OSS server — Temporal does not emit LLM provenance; workflow inputs are whatever the client `StartWorkflow` payload contained, stored opaquely as `DataBlob` history.

### 2. Is reproduction deterministic?

**Workflow replay is deterministic; binary rebuild is not guaranteed deterministic.**

- **Workflow determinism:** Documented "Workflow code must be deterministic" (`docs/architecture/README.md:35`) combined with deterministic proto marshaling option (`common/persistence/serialization/codec.go:48-91` `WithDeterministicProto3`) and sorting before event emission (`service/worker/batcher/activities.go:68` "Sort by request type for deterministic output") yields history-replay determinism. `workflow_rebuilder.go:92` replays the same branch and rebuilds identical `MutableState`. Tests pin jitter (`common/config/config.go:353` explicit seed for fault injection) and sort outputs for deterministic assertions (`tools/flakereport/slack.go:193`).
- **Build determinism:** Go's `CGO_ENABLED=0` and pinned pins help, but builds omit `-trimpath`, `GOFLAGS=-trimpath`, `-buildvcs` stamping control, and `SOURCE_DATE_EPOCH`; `docker/docker-bake.hcl:39` explicitly rejects digest pinning for multi-arch convenience, so Docker layer caching may vary. Goreleaser archives use `version_template "{{ .Version }}-SNAPSHOT-{{ .ShortCommit }}"` (` .goreleaser.yml:9`); snapshot builds embed short commit not full provenance. No checked check that two `make bins` invocations on same commit produce byte-identical outputs. Map-iteration non-determinism was historically a bug source (`docs/architecture/workflow-update.md:67` "Because maps enumeration in Go is non-deterministic, before being sent out, Updates are sorted").

### 3. Are all contributing factors recorded?

**No — critical factors are unrecorded.**

- **Recorded:** Git SHA/time/modified + Go version/arch (`common/build/build.go:39-53`), declared server `ServerVersion` (`common/headers/version_checker.go:26`), Go module closure via `go.sum`, Alpine base tag (`docker/docker-bake.hcl:40-42`), Goreleaser multi-arch + `CGO_ENABLED` matrix, image labels including `SERVER_VERSION`/`CLI_VERSION` extracted at build time (`.github/actions/build-docker-images/action.yml:84-96`).
- **Not recorded per artifact:** Full binary build command + `BUILD_TAG`/`TEST_TAG`, `ALPINE_TAG` effective value vs default, CLI download URL + its SHA (only version string), env `GOOS/GOARCH/TEMPORAL_SHA` at build, host kernel/toolchain, `proto/image.bin` inputs hash, dynamic config `development-*.yaml` + `dynamicconfig/development-sql.yaml` values, `SupportedServerVersions` negotiation per workflow, user approvals/RBAC (`CallerName`/`Principal` headers are propagated `common/headers/headers.go:32-42` but not persisted to history metadata). For workflow reproducibility, `params` like `TEMPORAL_TEST_DATA_ENCODING=json` (`Makefile:91`) used in OTEL path show env factors that affect artifact content but are ephemeral. Secrets management (`LookupSecret`) intentionally redacts values, breaking full factor closure.

### 4. Is reproducibility tested in CI?

**No dedicated reproducibility test exists.**

- CI (`run-tests.yml`, `build-and-publish.yml`, `release.yml`) runs `ci-build-misc: proto go-generate buf-breaking shell-check goimports gomodtidy ensure-no-changes` (`Makefile:12-21`) which guarantees generated code stays in sync (fails on `git status --porcelain` non-empty `Makefile:806-810`), and `lint-protos` via `buf` ensures compat. Build jobs do exercise both cross-compiles: `run-tests misc-checks: GOOS=windows GOARCH=amd64 make clean-bins bins` and `GOOS=darwin GOARCH=arm64` (`run-tests.yml:259-261`) but never compare outputs from two builds. `build-and-publish.yml` uploads `dist/*_linux_*` + bakes Docker images, but no `repro-test: build twice && diff/sha256` step, no `cosign`/`slsa` verification, no SBOM generation, no `goreleaser build --snapshot` hermetic check. Grep for `reproduc`, `provenance`, `sbom`, `attest`, `trimpath`, `SOURCE_DATE` across `.github/` yields only pin comments, not tests — confirmed negative evidence.

## Architectural Decisions

| Decision | Evidence | Effect on Provenance |
|----------|----------|----------------------|
| Event-sourced history as branchable tree with `BranchToken` + `VersionHistory` | `common/persistence/data_interfaces.go:1157-1184` V2 tree API + `ForkHistoryBranchRequest:907-921` | Enables full lineage: any workflow run can be traced to ancestor branch + fork point (`ForkNodeID`) and replayed. Bifurcates from single `NextEventID` linear model. |
| `debug.ReadBuildInfo` runtime introspection instead of `ldflags -X` stamping | `common/build/build.go:29-55` switch on `vcs.revision/time/modified` | Zero build-script coupling; works with plain `go build` and `goreleaser build`. Limited fields (no commit container digest, no build host). `Available` false on `go run` loses provenance silently. |
| OCI `docker-bake.hcl` with `timestamp()` label at bake time | `docker/docker-bake.hcl:64,91` `created=timestamp()` + `TEMPORAL_SHA` arg | Provides wall-clock + SHA in image manifest per OCI spec. `timestamp()` makes image non-reproducible byte-for-byte across rebuilds — tradeoff for human-readable `created`. |
| Single `ALPINE_TAG` var + no digest pin | `docker/docker-bake.hcl:38-42` comment rejects digest for multi-arch buildx correctness | Avoids "InvalidBaseImagePlatform" on manifest lists; at cost of base image mutability — two builds weeks apart may pull different Alpine patch layer without bump. |
| Goreleaser archives + checksums, but no SBOM/attestation | `.goreleaser.yml:8-89` before hook `go mod download`, `checksum: sha256`, `builds: CGO_ENABLED=0` | Guarantees per-file integrity via `checksums.txt` but not supply-chain provenance (builder identity, dependency closure). Easy to add `sboms`/`signs` stanza but omitted. |
| `ensure-no-changes` as codegen provenance gate | `Makefile:806-810` `git status --porcelain` + `git diff HEAD` on failure | Catches drift between proto source and generated `api/` + `mocks`; serves as partial input→output consistency check. Does not cover binary reproducibility. |

## Notable Patterns

- **BranchToken+TxID chaining**: Opaque branch token plus monotonic `TransactionID`/`PrevTransactionID` creates an auditable chain without exposing storage internals; combined with `DBRecordVersion` optimistically guards concurrent writes (`WorkflowMutation.Condition/DBRecordVersion` `common/persistence/data_interfaces.go:375-378`).
- **Pinned-toolchain via Makefile `.stamp` sentinels**: Each tool (`golangci-lint-v2.13.0`, `buf-v1.6.0`, `nilaway-v0.0.0-20260717164209…`) installs to `.bin` with `.stamp` sentinel (`Makefile:180-236`); comment explicitly "pin the pseudo-version for reproducible CI" (`Makefile:221`) — avoids floating `go install @latest`.
- **Label-per-arch bake**: Docker images built per-arch binary via `dist/` → `docker/build/amd64|arm64/temporal-server` re-org (`build-docker-images/action.yml:66-74` `organize-binaries`) then `docker buildx bake` with `--push` — keeps provenance fields `SERVER_VERSION`/`CLI_VERSION` extracted from actual binary (`extract-binary-version` steps `action.yml:84-96`) rather than git tag alone.
- **Deterministic codec as opt-in**: `WithDeterministicProto3` defaults off (`codec.go:43` false) and callers must pass option; avoids perf cost for hot path but requires discipline — history paths that forget the flag risk map-randomized bytes.

## Tradeoffs

- **Human-readable `ServerVersion` vs commit-precise artifact ID**: Static `ServerVersion = "1.32.0"` (`version_checker.go:26`) simplifies semver negotiation (`SupportedServerVersions >=1.0.0 <2.0.0`) but two commits on `main` post-tag share identical `ServerVersion` — Docker `SERVER_VERSION` label loses per-commit granularity; only `TEMPORAL_SHA` distinguishes. Trade: UX simplicity over byte-precise provenance.
- **No digest pin for multi-arch correctness** (`docker-bake.hcl:39` note): Pinning Alpine by digest breaks `buildx` platform resolver for manifest lists; keeping tag-only tag introduces base drift risk.
- **`timestamp()` label freshness vs reproducibility**: Fresh `created` per bake aids UI sorting but defeats `diff` equality; `SOURCE_DATE_EPOCH` would need a fixed time anchored to `vcs.time` (available in `build.InfoData.GitTime` but not forwarded to bake).
- **Full history fidelity vs storage cost**: Per-event `HistoryEvent` + `Checksum` + `HistoryStatistics.SizeDiff` (`data_interfaces.go:771-774`) retains enough to replay bit-for-bit, at cost of history table growth and archival pressure (`HistoryTaskQueueManager` DLQ). Pruning (`HistoryBranchDetail` gc via `DeleteHistoryBranch` `data_interfaces.go:947-953`) permanently destroys provenance.

## Failure Modes / Edge Cases

- **Build info absent on `go run` / `go test`**: `debug.ReadBuildInfo() !ok` leaves `InfoData.Available=false` (`common/build/build.go:30-32`) — binaries tested via `make unit-test` have no `vcs.revision`; any downstream provenance check naïvely reading `build.InfoData.GitRevision` gets `""` without error.
- **Modified working tree builds**: `vcs.modified=true` (`build.go:52`) is recorded but never fails the build; a dirty-tree Docker image is indistinguishable in `ServerVersion` label and only visible by inspecting `/usr/local/bin/temporal-server` via `go version -m` — CI does not gate on modified.
- **Map iteration randomization without sort**: Before `workflow/update/registry.go:343` "Sort Updates by time… for deterministic order", map-sourced updates shipped in random order — historic repro runs diverged; revert of sort reintroduces nondeterminism.
- **Alpine tag bump lag**: If `ALPINE_TAG` update commit is separate from code commit, `git bisect` on a build artifact cannot attribute a CVE fix to the correct input — provenance records the tag value at bake time but not the digest it resolved to.
- **Branch fork from invalid `ForkNodeID`**: `ForkHistoryBranchRequest` comment "Application must provide a void forking nodeID … >1" (`data_interfaces.go:914-916`) — supplying `ForkNodeID=1` or non-existent nodeID yields branch with truncated history; subsequent `ReadHistoryBranch` returns partial provenance without signaling fork failure beyond gRPC error.
- **Goreleaser snapshot vs release divergence**: Snapshot builds append `-SNAPSHOT-{ShortCommit}` (` .goreleaser.yml:9`) without `GITHUB_TOKEN` release gating; publishing a snapshot artifact as "release" would lose `checksum` attestation path (`build-binaries/action.yml:28-35` snapshot skips publish).

## Future Considerations

- Add `SOURCE_DATE_EPOCH=$(git log -1 --format=%ct)` + `GOFLAGS=-trimpath -buildvcs=true` and `ldflags "-s -w"` to `Makefile` `go build` lines (`Makefile:367-392`) and wire `build.InfoData.GitTime` into `docker-bake.hcl` `created` label instead of `timestamp()` for byte-reproducible images; add a CI job `repro-build: make clean-bins bins && sha256sum temporal-server && re-run && diff` per arch.
- Enable Goreleaser `sboms` (`syft`) + `signs`/`attestations` (cosign/SLSA) in `.goreleaser.yml` and set `provenance`/`sbom: true` defaults in `docker-bake.hcl`→`build-push-action`; publish `checksums.txt` sig alongside `dist/`.
- Snapshot the *effective* inputs per build artifact into a `build.provenance.json` sidecar (fields: `git SHA`, `git Time`, `ALPINE_TAG`+resolved digest `docker buildx imagetools inspect`, `CLI_VERSION` URL+SHA, `ServerVersion`, `BUILD_TAG`, Go version) and embed as `org.opencontainers.image.*` annotations; expose via `temporal-server --version --format json` (already present for `InfoData` but not version command).
- For workflow provenance, emit a top-level `WorkflowExecution provenance` manifest in `GetWorkflowExecutionResponse` (`data_interfaces.go:303`) summarizing `{Start RequestId, Caller, SDK Name/Version via headers, Config hash, Visibility search attributes}` so any `HistoryEvent` blames back to caller without full history scan; currently `headers.Current` only flows via `SetVersions` context prop, not persisted.
- Pin Alpine by *digest per TAG* via `ALPINE_TAG` + `ALPINE_DIGEST` dual vars and validate in `docker-bake.hcl`/`Dockerfile` with `FROM alpine:${ALPINE_TAG}@sha256:${ALPINE_DIGEST}` fallback for arch-specific inspect parity; document digests in `go.mod`-like lock.

## Questions / Gaps

- Is SLSA level 1-2 attestation intentionally deferred? No evidence of `.slsa` config, `slsa-framework/slsa-github-generator`, or `actions/attest-build-provenance` in `.github/`; no ADL/ADR explaining the omission.
- Should CLI download provenance be part of image provenance? `build-docker-images/action.yml:75-82` `download-cli` fetches `temporal` CLI at build time (version via arg or helper default) but the downloaded tarball SHA and provenance are never recorded in the image's labels or manifest — can't reconstruct which CLI build produced a given `admin-tools` image.
- How are user *approvals* (RBAC decisions) evidenced for artifact traceability? `common/headers/headers.go:24-42` propagates `principal-type/name` but no `approval` event type exists in `historypb.HistoryEvent`; Nexus and `Update` workflows may require explicit approval provenance not captured.
- Do dynamic config overrides (`dynamicconfig/development-sql.yaml`, `config/docker.yaml:96` `{{ env "ES_VERSION" }}`) belong in per-workflow-run provenance? They affect scheduling/sharding and thus generated artifacts (timers, transfer tasks) but are only file-system state at server start.
- What is the retention policy for `BranchToken` history that provides the only input trace for old runs? `DeleteHistoryBranch`/`TrimHistoryBranch` (`data_interfaces.go:947-968`) + archival queue could erase provenance for closed workflows — retention vs reproducibility tradeoff is undocumented at artifact level.

---
Generated by `16.03-artifact-provenance-and-reproducibility` against `temporal`.
