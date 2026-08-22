# Source Analysis: reports

## 13.01 Error Taxonomy

### Source Info

| Field | Value |
|-------|-------|
| Name | reports |
| Path | `studies/agent-harness-study/sources/reports` |
| Language / Stack | N/A — the source directory contains no code, configuration, or text files |
| Analyzed | 2026-08-22 |

## Summary

The selected source is an empty corpus. An exhaustive filesystem scan of `studies/agent-harness-study/sources/reports` found zero files and zero symlinks; the only content is a chain of three empty directories (`sources/reports/`, `sources/reports/source/`, `sources/reports/source/24.01-public-api-surface/`). There are no error type definitions, error classification code, error handling dispatch paths, or error documentation to analyze — none of the four evidence targets demanded by the dimension (`studies/agent-harness-study/dimensions/13.01-error-taxonomy.md:16-19`) exist in the mount. Consequently, no dimension question (`studies/agent-harness-study/dimensions/13.01-error-taxonomy.md:23-26`) can be answered with artifact evidence, and the rating reflects absence of any taxonomy within the source boundary.

Per the study's isolation rules, this analysis did not reach outside `studies/agent-harness-study/sources/reports` (no sibling sources such as `../langgraph/`, no provider configuration, and no generated reports under `reports/`). Because the mount contains no files, artifact-level citations are impossible; all citations in this report therefore point to the injected dimension definition (`studies/agent-harness-study/dimensions/13.01-error-taxonomy.md:1`), the one permitted input that defines what was searched for, alongside verbatim command output recorded below. The empty mount appears to be a harness/population issue rather than a property of the artifacts this source was intended to represent, and it is flagged in Questions / Gaps below.

**No clear evidence found** for every evidence area in this dimension. Search boundary: recursive enumeration of all entries (files, symlinks, directories) under `studies/agent-harness-study/sources/reports`, plus a recursive content grep for `error`, `exception`, `taxonomy`, and `retry` across that directory, which matched nothing because there are no files to search.

Verbatim search record (run 2026-08-22):

```
$ find sources/reports -type f -o -type l | wc -l
0
$ find sources/reports | sort
sources/reports
sources/reports/source
sources/reports/source/24.01-public-api-surface
$ grep -riE "error|exception|taxonomy|retry" sources/reports | wc -l
0
```

## Rating

**1 / 10**

Rationale per the rubric row "1-3: Absent, implicit, ad-hoc, or unsafe" (`studies/agent-harness-study/dimensions/13.01-error-taxonomy.md:32`): the score is at the floor of the band because there is literally nothing present — not even implicit classification. This is not a judgment that the underlying artifacts lack an error taxonomy; it is a statement that the provided corpus contains zero analyzable bytes. Guiding check — "Can you tell from the error type whether to retry, escalate, or stop?" (`studies/agent-harness-study/dimensions/13.01-error-taxonomy.md:37`) — cannot be evaluated because no error types exist in the corpus.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`. Because the source mount contains zero files, each row cites the line in the dimension definition that specifies the required evidence area (`studies/agent-harness-study/dimensions/13.01-error-taxonomy.md:16-19`) and records the observed null result; raw command output is quoted in Summary above.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Error type enums | No clear evidence found. Dimension step 1 asks to identify error type definitions (`studies/agent-harness-study/dimensions/13.01-error-taxonomy.md:9`); recursive `find sources/reports -type f -o -type l` returned 0 results. | `studies/agent-harness-study/dimensions/13.01-error-taxonomy.md:9` |
| Error classification code | No clear evidence found. No `.py`, `.ts`, `.go`, or other source files exist anywhere in the tree; classification-by-source check (`studies/agent-harness-study/dimensions/13.01-error-taxonomy.md:10`) is unanswerable. | `studies/agent-harness-study/dimensions/13.01-error-taxonomy.md:10` |
| Error handling dispatch | No clear evidence found. Routing/handling inspection (`studies/agent-harness-study/dimensions/13.01-error-taxonomy.md:11`) found no dispatch logic; content grep for `error\|exception\|taxonomy\|retry` over the entire source returned 0 matches. | `studies/agent-harness-study/dimensions/13.01-error-taxonomy.md:11` |
| Error documentation | No clear evidence found. Documentation check (`studies/agent-harness-study/dimensions/13.01-error-taxonomy.md:12`) found no README, Markdown, or docs files inside the source boundary. | `studies/agent-harness-study/dimensions/13.01-error-taxonomy.md:12` |
| Corpus structure | The only filesystem structure present: three nested empty directories, deepest being `source/24.01-public-api-surface/`; see verbatim `find` output in Summary. | `studies/agent-harness-study/dimensions/13.01-error-taxonomy.md:14-19` |

## Answers to Dimension Questions

Question numbering follows `studies/agent-harness-study/dimensions/13.01-error-taxonomy.md:23-26`.

1. **Are errors classified by source?**
   No clear evidence found. The source contains no files of any kind, so there are no error types to inspect for model/provider/tool/validation/policy/context/user/infrastructure/timeout attribution (the category list is defined at `studies/agent-harness-study/dimensions/13.01-error-taxonomy.md:5`).
2. **Is the taxonomy used for handling?**
   No clear evidence found. With zero files there is no dispatch logic, switch/match on error kind, or handler registry to examine.
3. **Are error categories documented?**
   No clear evidence found. No documentation files exist inside the source boundary.
4. **Can new error types be added without breaking existing handling?**
   Cannot be assessed. Extensibility requires existing interfaces or handler contracts to evaluate; none exist in the corpus.

## Architectural Decisions

No clear evidence found. No architectural decisions can be attributed to an empty directory tree. The only observable structural fact — the reserved but unpopulated subdirectory `sources/reports/source/24.01-public-api-surface/` — suggests an intended per-dimension layout that was created but never populated; this is an inference about harness behavior, not about error-handling architecture in studied artifacts.

## Notable Patterns

No clear evidence found for artifact-level patterns. One meta-observation worth recording for the study operators: the corpus mirrors the shape of a report output tree (`source/<dimension>/`) yet ships empty, which indicates the source-population step for this task either has not run or failed silently. Detecting and failing fast on empty mounts would prevent vacuous analysis runs like this one.

## Tradeoffs

Not applicable — no implementation exists in the corpus to trade off. Stated explicitly so the section is not left empty: the absence of files means there are also no risks (e.g., no unsafe catch-alls or swallowed errors) that could be misread as findings.

## Failure Modes / Edge Cases

- **Empty-source dispatch**: the harness selected and dispatched a study task against `studies/agent-harness-study/sources/reports` when that directory contained 0 files, producing a report whose only citable input is the dimension definition (`studies/agent-harness-study/dimensions/13.01-error-taxonomy.md:1-41`). A pre-flight non-empty validation would have caught this.
- **Silent partial population**: the presence of `sources/reports/source/24.01-public-api-surface/` (empty) shows a copy/sync step created directory scaffolding without file contents — consistent with an interrupted or filtered sync where files were skipped but directories were made.
- **Temptation to cross boundaries**: the natural remediation would be reading the real reports under the workspace's `reports/` output tree, but Hard Rule 1 and the Source Isolation Rules forbid inspecting generated reports outside the mounted source; this run honored that boundary.

## Future Considerations

1. **Harness fix (actionable)**: before rendering a prompt for a directory-kind source, assert the mount contains at least one file; abort or re-sync otherwise. This converts today's silent vacuous run into an explicit pipeline error.
2. **Re-run this dimension** after `studies/agent-harness-study/sources/reports` is actually populated with report artifacts, so the four questions (`studies/agent-harness-study/dimensions/13.01-error-taxonomy.md:23-26`) can be answered against real content with artifact-shaped citations such as `path/to/errors.go:42`.
3. **Sync verification**: have the population step emit a manifest (file count + checksums) into the mount root so analyses can distinguish "empty by design" from "sync failed".

## Questions / Gaps

- Was `sources/reports` intended to be populated with copies of the study's generated reports before this task ran? The scaffolded-but-empty `sources/reports/source/24.01-public-api-surface/` directory strongly suggests yes, and that the copy step did not complete.
- Is an "empty corpus" result expected to be recorded as a low rating (as done here, 1/10, per `studies/agent-harness-study/dimensions/13.01-error-taxonomy.md:32`), or should the harness treat it as a skip/failure condition? This report chooses honest scoring per the rubric; the alternative policy should be decided by study maintainers.
- All four dimension questions remain open pending a populated corpus. Every answer above is a null result bounded to `studies/agent-harness-study/sources/reports`; nothing was inferred from sibling sources or from the workspace's `reports/` output directory.

---

Generated by `dimensions/13.01-error-taxonomy.md:1` against `reports`.
