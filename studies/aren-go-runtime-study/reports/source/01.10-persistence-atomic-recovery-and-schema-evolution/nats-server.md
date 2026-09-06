# Source Analysis: nats-server

## 01.10 Persistence, Atomic Recovery, and Schema Evolution

### Source Info

| Field | Value |
|-------|-------|
| Name | nats-server |
| Path | `sources/nats-server` |
| Language / Stack | Go (server core), file and memory storage engines, custom NRG Raft |
| Analyzed | 2026-08-29 |

## Summary

JetStream persistence is split across two layers: a per-stream log (file or memory) and a Raft log for cluster consensus. The file store (`server/filestore.go`) treats durability as a stack of side files around append-only message blocks: block files (`N.blk`) hold encoded records with per-record checksums, `index.db` periodically checkpoints a fast-load snapshot, `meta.inf`/`meta.sum` carry stream config plus an HMAC checksum, and `obs/o.dat` plus per-message auxiliary files (`thw.db`, `sched.db`, `sources.db`) round out the picture. Atomicity for small state files comes from a write-to-temp-then-rename primitive (`writeAtomically` → `writeAtomicallyWithTemp`, `server/filestore.go:14562-14599`); atomicity for the stream-level directory manipulation during purge comes from a three-rename dance that uses the presence of `__new_msgs__/N.blk` as a "tombstone completed" sentinel (recoverPartialPurge, `server/filestore.go:10706-10748`). The Raft layer uses a separate WAL (the file store's own StoreMsg path), a term/vote file (`tav.idx`, `server/raft.go:5320-5354`), and versioned snapshots (`snap.<term>.<index>`, `server/raft.go:1880-1898`).

Recovery after a crash is an explicit, multi-step procedure: on `newFileStoreWithCreatedAndMode` (`server/filestore.go:412-662`) the code first calls `recoverPartialPurge`, then attempts `recoverFullState` (which validates a highwayhash checksum and refuses if the latest block index/checksum mismatches — `server/filestore.go:2167-2212`); if that fails it falls back to `recoverMsgs` (a full block-file scan that rebuilds state, honours tombstones, and may truncate a torn block — `server/filestore.go:1581-1880`). For Raft the replay path iterates the WAL starting from `state.FirstSeq`, validates `ae.pindex == index-1` for the first entry, and truncates back to the last good entry if a gap is found (`server/raft.go:583-625`).

Schema evolution is handled with explicit version constants. The stream-state snapshot has `fullStateVersion = 4` with `fullStateMinVersion = 1` (`server/filestore.go:12096-12100`), with optional fields gated by `if version >= 2` / `if version >= 3` / `if version >= 4` during decode (`server/filestore.go:2097-2103`). Stream-replicated state has `streamStateVersion = 1` and `streamStateVersionSources = 2` (`server/store.go:213-225`). Consumer state encoding uses `magic = 22`, `version = 1`, `newVersion = 2` (`server/filestore.go:299-307`). An older binary refuses an unknown version outright (`recoverFullState`: `version < fullStateMinVersion || version > fullStateVersion → errCorruptState`, `server/filestore.go:2004-2008`).

Crash, fault, and durability tests cover all the main crash windows: torn block, truncated index.db, corrupt term/vote file, mid-purge crash, post-snapshot state drift, and consumer state corruption (see e.g. `TestFileStoreEncryptedPurgeRecoveryAfterKeyRename` `server/filestore_test.go:823-891`, `TestFileStoreMsgBlkFailOnKernelFaultLostDataReporting` `server/filestore_test.go:5126-5225`, `TestNRGWALEntryWithoutQuorumMustTruncate` `server/raft_test.go:1170`, `TestNRGTruncateWALRevertsUncommittedAddPeer` `server/raft_test.go:1384`).

## Rating

**8 / 10**

Rationale: Atomicity is a first-class concept (write-and-rename for all small state files; multi-rename with sentinel for the big one), the recovery path is honest about its failure modes (checksum mismatch, sequence drift, corrupt block → truncate-and-report-lost), and there is real fault-injection test coverage. The system does not pretend to be exactly-once: cluster replication uses Raft log replay and stream-store message-store, so a crash between WAL append and state-machine apply can replay work but cannot duplicate external effects because JetStream only exposes internal "store" / "ack" abstractions, not side-effectful tools. What keeps this from a 9 or 10:

- The WAL writes are best-effort `fsync`: `os.WriteFile` on `index.db` (`server/filestore.go:12364`) is only synced when `SyncAlways` is on (`server/filestore.go:14549-14556`), and `SyncAlways` is off by default — so a power loss between the WAL append and the next index.db flush will lose messages from the perspective of fast recovery, though block-level replay still finds them via `recoverMsgs` rebuild. The recovery is honest about this (`rebuildState` rebuilds from raw bytes and the per-record checksum protects integrity, `server/filestore.go:1589-1880`), but the system will not ack a publish until the WAL write is durable, so this is at-least-once for store and at-least-once for the cluster-side proposal.
- `writeMsgRecord` does not include a per-record header checksum in the file (only when `fsync` is called) but does include a per-record record-hash in the binary format (`server/filestore.go:7713-7726`) which protects against silent corruption on recovery.
- The `recoverMsgs` rebuild path silently truncates a torn last block (`server/filestore.go:1748-1753, 1768-1774, 1793-1798`) and reports the gap via `LostStreamData`, but this is a coarse granularity (one truncation, no finer-grained reconciliation per record).
- Schema-evolution tests are not extensive at the storage format level — version bumps exist and gating is real (`server/filestore.go:2097-2103`), but I did not find a "newer binary opens old data" test that exercises all the version branches systematically.

## Evidence Collected

Every entry includes a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Storage interface | `StreamStore` interface with `StoreMsg`, `StoreRawMsg`, `Snapshot`, `Stop` | `server/store.go:95-145` |
| Stream-state encoding version | `streamStateMagic=42`, `streamStateVersion=1`, `streamStateVersionSources=2` | `server/store.go:213-225` |
| Stream-replicated state decode | `DecodeStreamState` with `withSources` flag | `server/store.go:262-362` |
| File-store entry point | `newFileStoreWithCreatedAndMode` orchestrates full recovery | `server/filestore.go:412-662` |
| Block-file layout | `%d.blk` blocks with encrypted `%d.key` companion file, `msgs/` subdir | `server/filestore.go:319-335, 350-356` |
| Per-record binary format | `total_len(4) sequence(8) timestamp(8) subj_len(2) subj hash(8)`; hdr bit set when `mhdr` present | `server/filestore.go:7680-7715` |
| Per-record hash | `highwayhash.Sum64` over hdr+subj+body | `server/filestore.go:7713-7726` |
| Block checksum verified on rebuild | `rebuildStateFromBufLocked` rejects records whose checksum mismatches the trailing 8 bytes | `server/filestore.go:1777-1800` |
| Torn-record handling | Truncate to last good index and report `LostStreamData` | `server/filestore.go:1747-1753, 1768-1774, 1793-1798` |
| Tombstone record (`tbit`) handling | `seq&tbit != 0` branch in recovery | `server/filestore.go:1807-1816` |
| Append-only block writer | `writeMsgRecordLocked` reserves header space, encodes, hashes, then appends; flusher does `mb.writeAt(buf, wp)` | `server/filestore.go:7598-7781, 8699-8718` |
| Atomic file write primitive | `writeAtomically` / `writeAtomicallyWithTemp`: write `name+".tmp"`, `Close`, `os.Rename`, optional dir `syncDir` | `server/filestore.go:14562-14599` |
| Directory fsync after rename | `syncDir` gated by `canFsyncDirectories` (skipped on Windows) | `server/filestore.go:14560, 14591-14597, 14628-...` |
| `SyncAlways` behaviour | `writeFileWithOptionalSync` chooses `O_SYNC` when `syncAlways` or `syncOnFlush` | `server/filestore.go:14549-14552` |
| `index.db` (stream-state) file | `streamStreamStateFile = "index.db"`, magic `11`, version `4`, gated fields per version | `server/filestore.go:355, 12096-12100, 12244` |
| Stream-state encode | Binary: Msgs/Bytes/FirstSeq/FirstTime/LastSeq/LastTime + PSIM + blocks + last-block-index/checksum; optional `ttls`/`schedules`; SHA-256-derived highwayhash key from `cfg.Name` | `server/filestore.go:12244-12314` |
| Stream-state decode | `recoverFullState` validates highwayhash, magic, version bounds; loads PSIM, blocks, dmap (avl) | `server/filestore.go:1950-2242` |
| Stale-index detection | On mismatch between `mb.lastChecksum()` and stored `lchk`, returns `errPriorState` and rebuilds | `server/filestore.go:2182-2212` |
| Recovery fallback | On `recoverFullState` failure, falls back to `recoverMsgs` which scans `*.blk` files and rebuilds | `server/filestore.go:518-557, 2563-2712` |
| Mid-purge crash recovery | `recoverPartialPurge`: presence of `__new_msgs__/N.blk` means purge completed; only key files means purge partial → remove matching `.blk` | `server/filestore.go:10706-10748` |
| Purge atomicity | Writes `N.key` then `N.blk` into `__new_msgs__`, renames old `msgs`→`__msgs__`, renames `__new_msgs__`→`msgs`; block moved last so prior renames are assumed on recovery | `server/filestore.go:10664-10684` |
| Background state flush | `flushStreamStateLoop` runs every ~2 minutes + jitter | `server/filestore.go:12103-12130` |
| Tombstone-aware recovery | If `prior.LastSeq > fs.state.LastSeq` after rebuild, write a tombstone to preserve cluster sequence | `server/filestore.go:540-553` |
| Encryption key recovery | `recoverAEK` from `meta.key`; missing key on encrypted store → `errNoMainKey` | `server/filestore.go:508-515, 955-1003` |
| Block-level encryption | `genEncryptionKeysForBlock` writes `N.key` via atomic write; XorKeyStream on flush | `server/filestore.go:5192-5215, 8686-8697` |
| Per-block flush | `flushPendingMsgsLocked` appends cache to block fd at `cache.wp`; if `SyncAlways` calls `f.Sync()`, else marks `needSync` | `server/filestore.go:8655-8766` |
| Sync interval | `defaultSyncInterval = 2 * time.Minute`; `syncBlocks` walks all blocks, syncs dirty FDs, compacts when full | `server/filestore.go:339, 8163-8322` |
| Dir fsync after sync | `syncBlocks` calls `syncDir(msgDir)` when `needDirSync` | `server/filestore.go:8308-8315` |
| Consumer state write | `consumerFileStore.writeState` calls `writeFileWithOptionalSync` (atomic rename) | `server/filestore.go:13953-13986` |
| Consumer state encode | `encodeConsumerState` with `magic = 22`, `version = 2`, varints, optional Pending/Redelivered | `server/filestore.go:478-538` |
| Consumer state load | `loadState` decodes and validates | `server/filestore.go:13428` |
| Consumer config + checksum file | `meta.inf` + `meta.sum` (highwayhash) | `server/filestore.go:1006-1043` |
| Snapshot stream format | `streamSnapshot` writes tar+s2 archive of `meta.inf`, `meta.sum`, `msgs/N.blk`, `index.db`, consumer states | `server/filestore.go:12706-12890` |
| Binary stream snapshot (server-to-server) | Encoded `index.db` v4 payload with header trailer for fast validation | `server/filestore.go:12964-13066` (approx) |
| Snapshots directory | `snapshotsDir = "snapshots"`, `snapFileT = "snap.%d.%d"` | `server/raft.go:1880-1882` |
| Term/vote file | `termVoteFile = "tav.idx"`, `termVoteLen = idLen + termLen` | `server/raft.go:5320-5322` |
| Peer-state file | `peerStateFile = "peers.idx"` | `server/raft.go:5284-5307` |
| `writeTermVote` durability | Writes via `writeFileWithSync` (fsync) | `server/raft.go:5425-5450` |
| `writePeerState` durability | Writes via `writeFileWithSync` (fsync) | `server/raft.go:5287-5307` |
| WAL write path | `storeToWAL` calls `n.wal.StoreMsg(_EMPTY_, nil, ae.buf, 0)` — encodes the whole `appendEntry` and persists via the stream's StoreMsg | `server/raft.go:4990-5014` |
| Snapshot install | `installSnapshot` writes `snap.<term>.<index>` via `writeFileWithSync`, deletes prior snapfile, then `wal.Compact(lastIndex+1)` | `server/raft.go:1596-1638` |
| Snapshot recovery | `setupLastSnapshot` scans `snapshotsDir`, picks highest `(term,index)`, restores peer state, compacts WAL up to `snap.lastIndex+1`, removes older snaps | `server/raft.go:1903-1979` |
| WAL replay on startup | `initRaftNode` loops `state.FirstSeq..state.LastSeq` calling `loadEntry`; if first entry misaligned → `truncateAndErr(n.pindex)`; if later entry misaligned → `truncateAndErr(index-1)` | `server/raft.go:566-625` |
| WAL truncate | `truncateWAL(term,index)` calls `n.wal.Truncate(index)` and rolls back `commit`, `processed`, `applied`; also removes stale snapfile | `server/raft.go:4260-4330` |
| Crash window — purge mid-flight | `TestFileStoreEncryptedPurgeRecoveryAfterKeyRename` simulates crash after moving `N.key` into `__new_msgs__` but before `N.blk` | `server/filestore_test.go:823-891` |
| Crash window — torn last block | `TestFileStoreMsgBlkFailOnKernelFaultLostDataReporting` deletes first/interior/last `.blk` files and asserts `LostStreamData` is reported | `server/filestore_test.go:5126-5225` |
| Crash window — partial remove | `TestFileStoreRecoverAfterRemoveOperation` removes `.blk` mid-recovery | `server/filestore_test.go:10028-...` |
| Crash window — corrupt PSIM | `TestFileStoreCorruptPSIMOnDisk` (corrupt subject inside `index.db`) | `server/filestore_test.go:6965` |
| Crash window — corrupt full-state | `TestFileStoreRecoverFullStateDetectCorruptState` and `TestFileStoreWriteFullStateDetectCorruptState` | `server/filestore_test.go:9210, 9245` |
| Crash window — TTL with corrupt block | `TestFileStoreNoPanicOnRecoverTTLWithCorruptBlocks` | `server/filestore_test.go:10690` |
| Crash window — compact recovery | `TestFileStoreDmapBlockRecoverAfterCompact`, `TestFileStoreRecoverAfterCompact` | `server/filestore_test.go:8427, 10189` |
| Crash window — block-first last seq corrupt | `TestFileStoreMsgBlockFirstAndLastSeqCorrupt` | `server/filestore_test.go:7430` |
| Crash window — empty block recovery | `TestFileStoreRecoverWithEmptyMessageBlock` | `server/filestore_test.go:10300` |
| Crash window — no `index.db` recovery | `TestFileStoreRecoverWithRemovesAndNoIndexDB`, `TestFileStoreRecoverOnlyBlkFiles` | `server/filestore_test.go:7603, 9965` |
| Raft crash — decode truncated entry | `TestNRGAppendEntryDecodeTruncatedEntryLength` | `server/raft_test.go:157` |
| Raft crash — decode truncated peer state | `TestNRGPeerStateDecodeTruncated` | `server/raft_test.go:174` |
| Raft crash — recover without leader | `TestNRGRecoverFromFollowingNoLeader` | `server/raft_test.go:194` |
| Raft crash — uncommitted add-peer rollback | `TestNRGTruncateWALRevertsUncommittedAddPeer` | `server/raft_test.go:1384` |
| Raft crash — uncommitted remove-peer rollback | `TestNRGTruncateWALRevertsUncommittedRemovePeer` | `server/raft_test.go:1463` |
| Raft crash — pending cache clear | `TestNRGTruncateWALClearsPendingAppendEntryCache`, `TestNRGInstallSnapshotClearsPendingAppendEntryCache` | `server/raft_test.go:1343, 1920` |
| Raft crash — quorum-driven WAL truncation | `TestNRGWALEntryWithoutQuorumMustTruncate` | `server/raft_test.go:1170` |
| Raft crash — snap survives truncate | `TestNRGDontRemoveSnapshotIfTruncateToApplied` | `server/raft_test.go:2277` |
| Raft crash — snap install twice | `TestNRGIgnoreDoubleSnapshot` | `server/raft_test.go:2409` |
| Raft crash — recovery respects log emptiness | `TestNRGRecoverPindexPtermOnlyIfLogNotEmpty` | `server/raft_test.go:2546` |
| Raft crash — truncate down to commit | `TestNRGTruncateDownToCommitted` | `server/raft_test.go:2625` |
| Shutdown ordering | `disableJetStream` → `meta.Stop()` → `meta.WaitForStop()` → `shutdownRaftNodes()` | `server/jetstream.go:713-767` |
| Disk IO semaphore | Single semaphore serializes atomic-rename writes (`fs.dios`) | `server/filestore.go:180-206, 14549-14599` |
| Persist-failure escape hatch | Permission denied on stream state → `s.ShutdownJetStream()` and warn "messages in block cache could be lost" | `server/filestore.go:12118-12124, 7422-7428` |

## Answers to Dimension Questions

**What is the atomic unit, and does it match Aren-owned lifecycle state rather than external effects?**

- For small state files (`meta.inf`, `meta.sum`, `index.db`, `o.dat`, `thw.db`, `sched.db`, `sources.db`, `peers.idx`, `tav.idx`, snapshots, per-block `*.key`): the atomic unit is a single file replaced via `writeAtomically` / `writeFileWithSync` (`server/filestore.go:14562-14599`, `server/raft.go:5287-5307`, `server/raft.go:5425-5450`). A successful `os.Rename` is the publish/ack line.
- For the message-log itself: the atomic unit is a single record in a `.blk` file, protected by a `total_len(4) sequence(8) timestamp(8) subj_len(2) [hdr_len(4) hdr] subj msg hash(8)` framing and a per-record `highwayhash.Sum64` trailer (`server/filestore.go:7680-7726`). Recovery detects a torn record by `total_len > lbuf` or checksum mismatch and truncates the file back to the last good record (`server/filestore.go:1745-1800`).
- For purge / directory swap: the atomic unit is the **block file** within `__new_msgs__`. The block is moved last (after the key file), so on recovery the absence of the `.blk` sentinel means "purge incomplete → roll back" (`server/filestore.go:10664-10684, 10706-10748`).
- These units match JetStream-owned state: messages, stream config, consumer state, per-message TTL/scheduling/sourcing, NRG WAL, snapshots. External effects (subscriber acks back through NATS) are out of scope of the storage layer.

**After restart, how is interrupted work distinguished from failed or completed work?**

- Raft: `n.commit` (last committed) vs `n.applied` (last applied). On replay, the loop only re-runs unapplied entries from `state.FirstSeq..state.LastSeq`; if the first entry's `ae.pindex` does not match `n.pindex`, the WAL is truncated to `n.pindex` (the last good entry from before the snapshot or from prior recovery) (`server/raft.go:583-625`). The pattern "truncate to last good index, apply again" treats interrupted work as "to be redone" but only at the Raft-log level, never at the message-store level below.
- File store: there is no "ack back to producer" — once a record is in `.blk` it is durable. `recoverFullState` reads `state.db` and validates the checksum of the last block against the snapshot (`server/filestore.go:2167-2208`); a mismatch means "state is stale, rebuild from blocks". The rebuild path (`rebuildStateFromBufLocked`) skips corrupted or trailing-garbage records and reports `LostStreamData` so the upper layer can decide whether to forward-replay via Raft (`server/filestore.go:1714-1721, 1880`). Purge-with-tombstone distinguishes "completed purge" (block present in `__new_msgs__`) from "partial purge" (only key file present) (`server/filestore.go:10706-10748`).

**Can recovery duplicate a tool or remote side effect, and is that stated honestly?**

- Yes, recovery CAN duplicate work in two ways, both honestly acknowledged in comments:
  1. Raft replay re-runs `processAppendEntry` for any entry that was committed but not applied (`server/raft.go:566-625`). If a state-machine handler had side effects, those would re-occur. The stream layer's "state machine" is the per-stream store update, which is idempotent for stored messages (they are addressed by `(seq, subj, hdr, msg)` and sequence numbers are derived from the WAL index). For non-stream consumers like KV, JetStream's API contract is that reads-after-restart see the same state, not "exactly-once side effects".
  2. The stream-store write layer is at-least-once: a publish that gets a PubAck response can still have lost the message if the OS crashes between the kernel write and `f.Sync()` on the block fd. The code states this honestly in `flushStreamStateLoop` ("messages in block cache could be lost in the worst case" — `server/filestore.go:12118-12124`) and in `storeMsg` (`server/stream.go:7421-7428`). JetStream's `Duplicates` window (`stream.go:5728-5751`) is the mitigation at the API layer: a publisher can attach a `Nats-Msg-Id` and the broker dedupes within that window, not because the storage is exactly-once, but because the API is.

**How does a newer binary refuse, migrate, or safely open older data?**

- `recoverFullState`: refuses if `buf[0] != fullStateMagic` OR `version < fullStateMinVersion` OR `version > fullStateVersion`, deletes the bad file, returns `errCorruptState` (`server/filestore.go:2004-2008`). Forward-compatible decode: optional fields `ttls` (v2), `schedules` (v3), `lblk` always vs only when count>1 (v4) — the reader uses `if version >= 2` / `if version >= 3` / `if version >= 4` (`server/filestore.go:2097-2103, 2077-2081`).
- `DecodeStreamState`: accepts versions `streamStateVersion` and `streamStateVersionSources`; refuses unknown versions by checking `buf[1]` (`server/store.go:255-268`).
- `loadTermVote`: tolerates a short `tav.idx` (returns `(term, noVote, nil)` if `len(buf) < termLen` — `server/raft.go:5343-5345`).
- `termAndIndexFromSnapFile`: requires the exact `snap.<term>.<index>` filename format; bad filenames are deleted (`server/raft.go:1886-1898, 1928-1932`).
- `recoverPartialPurge`: explicit handling of older crash states (only key, no block → roll back; both present → complete).
- There is no in-place upgrade of `meta.inf` format — if the JSON cannot unmarshal into the current `StreamConfig`, the load fails. Schema migrations are handled by adding new fields with sensible defaults and using `omitempty` / versioned encoding for the binary paths.

## Architectural Decisions

- **Decouple fast-load snapshot from raw log**: `index.db` (fullState) lets a stream cold-start in O(blocks) without scanning every record; if it is missing, stale, or checksum mismatches, the code falls back to `recoverMsgs` which is the slower but authoritative path. This is a "trust-but-verify" split between fast-load and truth-source (`server/filestore.go:518-557, 2167-2212`).
- **Per-record hash in the message log**: Torn writes are detected by the trailing 8-byte highwayhash, allowing recovery to truncate to the last good record without losing earlier records (`server/filestore.go:7680-7726, 1777-1800`). This is what gives "honest recovery" rather than "all-or-nothing page".
- **Two-phase rename for purge**: Move `N.key` → `N.blk` into a staging dir, then rename old `msgs/` away and staging dir back. The block-last ordering is the trick that makes the staging dir a self-describing sentinel (`server/filestore.go:10664-10684`).
- **Raft owns cluster-state durability, store owns message durability**: The NRG layer fsyncs `peers.idx` and `tav.idx` immediately on every change because they are critical to safety (`server/raft.go:5287-5307, 5425-5450`); the stream's WAL is the file store's `StoreMsg` which by default is NOT fsynced per record, only on `SyncInterval` (`server/filestore.go:339, 14549-14556`). The tradeoff: throughput vs. window of loss; the cluster design assumes R≥3 so the lost window can be replayed from peers.
- **Block rotation + tombstones**: A new `.blk` file is created when the current block is full (`server/filestore.go:7883-7893`). Deleted messages leave behind tombstones (`tbit` records, `server/filestore.go:7748-7759, 1807-1816`) so sequence numbers persist across restarts even when the underlying data is gone. Compaction merges tombstones later (`syncBlocks`, `server/filestore.go:8163-8322`).
- **Background loop-driven checkpoints**: `_writeFullState` runs every ~2 min with jitter to spread startup storms on big clusters (`server/filestore.go:12103-12130`). It is also triggered synchronously by purge/compact (`server/filestore.go:10780`) and by snapshot (`server/filestore.go:12817`).
- **Three-level version constants**: stream-encoded snapshot v4 (`server/filestore.go:12096-12100`), stream-replicated state v1+v2 (`server/store.go:213-225`), consumer state v2 (`server/store.go:495-497`). Each has a min-version check at decode time.

## Notable Patterns

- **"Checksum + version + size-guard" triple at every decode boundary** — `recoverFullState` checks length, magic, version, and highwayhash; `rebuildStateFromBufLocked` checks size sanity + per-record hash; `termAndIndexFromSnapFile` checks filename pattern; `loadTermVote` checks length (`server/filestore.go:1972-2008, 1766-1800, 1886-1898, 5343-5352`). Each new persisted file follows the same shape.
- **Sentinel-driven crash recovery** — for purge, the presence of `.blk` in `__new_msgs__` distinguishes "completed" from "in-progress"; for Raft, the absence of an expected `ae.pindex` at `state.FirstSeq` distinguishes "WAL truncated to before this entry" from "WAL intact" (`server/filestore.go:10706-10748, server/raft.go:583-625`).
- **Idempotent state writes** — `writeStreamMeta` skips if `meta.inf` exists and is non-zero (`server/filestore.go:636-643`); `writePeerState` skips if encoded bytes unchanged (`server/raft.go:5289-5291`); `writeTermVote` skips if unchanged (`server/raft.go:5436-5439`). This keeps recovery cheap when nothing changed.
- **Disk IO semaphore (`dios`)** — serializes concurrent atomic-rename writes (`server/filestore.go:180-206, 14571-14572`) so we don't blow open-file limits under load.
- **Highwayhash keyed by stream name** — each stream has its own hash key derived from `sha256(streamName)` so a corruption in one stream cannot mask a corruption in another (`server/filestore.go:494-499, 12332-12335`).
- **Elastic cache for block buffers** — `mb.cache` is held via an elastic pointer so writes don't lose the buffer mid-flush but readers can still evict under memory pressure (`server/filestore.go:5086-5090, 5051-5078`).
- **Per-message auxiliary indexes as side files** — TTL (`thw.db`), scheduling (`sched.db`), sourcing (`sources.db`) are written separately so they can be missing/corrupt without taking down the main store; recovery iterates from the side file's "stamp" sequence forward (`server/filestore.go:2243-2402`).

## Tradeoffs

- **Throughput vs. durability window**: `SyncAlways` defaults to off; `O_SYNC` is only set when explicit (`server/filestore.go:14549-14556`). The `flushStreamStateLoop` admits that messages in block cache can be lost on a permission-denied flush (`server/filestore.go:12118-12124`). This is a deliberate latency/throughput choice for the default case; users who care must enable `SyncInterval`/`SyncAlways`/`Replica=1` accordingly.
- **Block size vs. compaction overhead**: Default block size is dynamic (`dynBlkSize`, `server/filestore.go:840-875`). Small blocks mean more files and more fsyncs; large blocks mean bigger per-recovery scans. `syncBlocks` does compaction in the background (`server/filestore.go:8163-8322`).
- **Cluster WAL replication via the same file store**: All Raft groups share the `store.StoreMsg` path (`server/raft.go:4998`). This avoids two storage engines but couples WAL durability to stream-store durability: a stream that is replicated cannot be tuned for sync independently from its WAL.
- **Mid-write truncation is coarse**: A torn record truncates the entire block to that index (`server/filestore.go:1681-1711, 1748-1753`), losing any valid records after the torn one. The `LostStreamData` reports the gap but the trade-off is "truncate the block" vs. "more complex per-record reconciliation".
- **Per-message aux indexes are recovered by scanning**: `recoverPerMessageState` walks from `scanSeq` to `LastSeq` if a side file's stamp is behind (`server/filestore.go:2346-2402`). For very large streams this is expensive; mitigated by skipping blocks whose `ttls/schedules` count is zero (`server/filestore.go:2378-2382`).
- **Compression/encryption is per-block**: This is good for blast radius but means a snapshot must stream per-block encrypt/decrypt (`server/filestore.go:12792-12800`).

## Failure Modes / Edge Cases

- **Torn last record**: detected by `index+msgHdrSize > lbuf`, recovered by `truncate(index)` and reported as lost data (`server/filestore.go:1744-1753`).
- **Corrupt per-record hash**: detected by `bytes.Equal(checksum, data[len(data)-recordHashSize:])`, recovered by truncating and reporting lost (`server/filestore.go:1789-1798`).
- **Missing `.blk` for a known block**: `recoverMsgBlock` returns the recovered mb; if `mb.first.seq == 0` the block is removed and the stream keeps going (`server/filestore.go:2603-2615`). If encryption key is wrong, the block is deleted to allow catchup (`server/filestore.go:2651-2662`).
- **Stale `index.db`**: detected by `mb.lastChecksum() != lchk` (`server/filestore.go:2182-2186`), or by finding block files with index > `blkIndex` (`server/filestore.go:2206-2209`), and the code falls back to `recoverMsgs`.
- **Mid-purge crash**: handled by `recoverPartialPurge` (`server/filestore.go:10706-10748`); the test `TestFileStoreEncryptedPurgeRecoveryAfterKeyRename` covers the trickiest case (only key moved, not block).
- **Permission denied on flush**: detected by `isPermissionError`, the server disables JetStream (`server/filestore.go:12118-12124, 7421-7428`).
- **Term/vote file too short**: tolerated by `readTermVote` returning `(term, noVote, nil)` if `len(buf) < termLen` (`server/raft.go:5343-5352`).
- **WAL gap at startup**: detected by `ae.pindex != index-1` for first entry → truncate to `n.pindex` (`server/raft.go:589-599`); for later entries → truncate to `index-1` (`server/raft.go:604-614`).
- **Snapshot file with bad name**: removed (`server/raft.go:1928-1932`).
- **Encoding errors during WAL write**: triggers `n.setWriteErrLocked`, `n.shutdown()`, and `assert.Unreachable` (`server/raft.go:5383-5393`).
- **Block file write error**: triggers `mb.werr = err`, async rebuild to compute `LostStreamData`, and a `setWriteErr` on the filestore (`server/filestore.go:8699-8713, 5441-5474`).
- **Concurrent fsync on the same block**: serialised via `syncMu` to avoid two flushes fighting (`server/filestore.go:204-205, 8164-8165`).

## Future Considerations

- The `SyncAlways`/`SyncOnFlush` knobs leave durability as opt-in; for workloads that need strict per-record durability without the latency penalty of `O_SYNC`, a small group-commit (e.g. fsync every N ms or N records) would close the window without paying the per-write cost. The infrastructure (`syncBlocks` periodic timer, `flushStreamStateLoop` jitter) is already in place to host that.
- `recoverMsgs` truncates an entire block on a torn record. With per-record checksums in place, finer-grained repair (skip only the torn record, keep what follows) is feasible if the on-disk format is changed to use absolute offsets in a trailer.
- Schema-evolution tests are present (`TestFileStoreMessageTTLRecovered`, version-gated fields in `recoverFullState`) but not systematic: there is no "newer binary opens a back-rev index.db" test exercising all v1..v4 branches together. Adding one would harden the format against future schema bumps.
- `index.db` write goes through plain `os.WriteFile` (`server/filestore.go:12364`) which is not in the atomic-rename path. The block-level log is the source of truth, so this is safe, but it makes the fast-load path vulnerable to `EIO` mid-write. Wrapping it in `writeAtomically` (with a longer content size, which is why the current code uses raw WriteFile) would close that.
- The mid-purge crash protocol assumes the rename is atomic on the filesystem. On filesystems where rename is not atomic (rare but real), the protocol falls apart. Adding a journal-style "purge in progress" file would make this robust.
- The "duplicate publish" window is currently mitigated by `Duplicates` in the stream config and `storeMsgId` (`server/stream.go:5728-5751`). A per-record published-once stamp in the block header (visible to the broker but not the wire) would let the broker dedupe without an external window.

## Questions / Gaps

- I did not find an explicit "what is the worst-case number of lost messages on a single-node power loss" test. The lost-data window is bounded by `f.cfg.SyncInterval` (default 2 min) for `index.db` and by `SyncAlways` for blocks, but no test exercises "kill -9 between WAL write and next fsync, assert N messages lost". Searched `server/filestore_test.go` and `server/raft_test.go` for `kill`, `crash`, `fault`, `chaos` — only `TestFileStoreMsgBlkFailOnKernelFaultLostDataReporting` (`server/filestore_test.go:5126`) covers file loss, not OS-level interruption.
- The Raft `wal.StoreMsg` path uses the same `StreamStore` as user data, which means encryption keys, compression, and indexing all apply to WAL entries. This is implicit — there is no doc or test asserting the WAL data is opaque to the storage layer.
- Snapshot format `snap.<term>.<index>` is keyed by term, not by content hash. Two snapshots with the same `(term, index)` overwrite each other (by design — `installSnapshot` deletes the old one — `server/raft.go:1619-1621`), but there is no defense against a partial snapshot write that leaves a corrupt file at the expected path. A `tmp` + `rename` (like `writeAtomically`) would be safer.
- I did not find explicit tests for "version boundary" upgrades of the on-disk format. The version constants exist and the decoders gate fields, but the only direct test I found is `TestFileStoreMessageTTLRecovered` (`server/filestore_test.go:9473`) which simulates a stream-state file from a prior version with TTLs.
- The cluster meta layer's own state (which is also Raft-replicated) is not covered in detail above; it shares the NRG infrastructure but its own persistence details live in `server/jetstream_cluster.go`. A follow-up dimension study on cluster-meta persistence would complement this one.

---

Generated by `dimensions/01.10-persistence-atomic-recovery-and-schema-evolution.md` against `nats-server`.