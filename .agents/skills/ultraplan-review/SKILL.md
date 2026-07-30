---
name: ultraplan-review
description: Manually run the UltraPlan review stage for a selected project sprint. Use only when the user explicitly invokes $ultraplan-review or directly asks to run this exact UltraPlan stage; do not invoke implicitly.
---

# UltraPlan Review

Run this stage interactively while preserving UltraPlan's governed artifact chain.

## Operating contract

1. Locate the workspace root and resolve the project and sprint. If either is ambiguous, ask the user instead of guessing.
2. Run `ultraplan project <project> status` and `ultraplan sprint <project> <sprint> status --json`. Treat files and fresh CLI status as authoritative; never hand-edit flow-state JSON.
3. Check these prerequisites:

- completed execute stage
- current planning artifacts and target implementation

4. Validate every prerequisite that has a sprint validation command. If anything is missing, invalid, stale, or internally inconsistent, show the exact gaps and ask whether to fill them. Do not fill prerequisite gaps until the user agrees. If they agree, run the corresponding earlier UltraPlan skills in canonical order, then return to this stage.
5. If the target is already complete and valid, summarize that state and ask before regenerating or materially changing it.
6. If the user explicitly asks for a proposal, analysis, or discussion only, inspect all relevant evidence and return that without writing artifacts or advancing state. Otherwise, do the stage now; do not stop at a proposal.
7. Use the current effective prompt and concrete paths:

    ultraplan sprint <project> <sprint> prompt review

   The resolved prompt can include workspace or project overrides and therefore takes precedence over the canonical prompt below.
8. Complete the stage-specific workflow, preserve unrelated user edits, and do not cross declared mutation boundaries.
9. Run `ultraplan sprint <project> <sprint> validate review` when supported. Fix validation findings within this stage rather than declaring success early.
10. Run `ultraplan sprint <project> <sprint> status --json` after writes or governed execution so flow-state, freshness, and artifact status are reconciled. Re-run project status if the project or sprint index changed.
11. Inspect downstream artifacts for references made stale by this change. Update directly coupled indexes/references when safe; otherwise report the exact dependent stage that must be revisited. Never delete or silently rewrite downstream decisions.
12. Finish with the artifact/result paths, validation outcome, state transition, and any remaining blocker or dependent stage.

## Stage workflow

Preview the review scope first:

    ultraplan sprint <project> <sprint> review --dry-run

Then run or resume the governed review:

    ultraplan sprint <project> <sprint> review

Interpret the verdict and findings for the user. Do not silently fix findings during the review stage; if fixes are requested, return to execute and rerun review so evidence and fingerprints remain authoritative.

## Canonical stage prompt

The current resolved CLI prompt wins over this embedded baseline.

# Automated Sprint Review

This is UltraPlan's embedded default for the product-owned `review` stage. A
workspace `prompts/review.md` is an optional override; it is never a required
planning input.

Review the frozen manifest assembled by UltraPlan. Work read-only against the
approved implementation target. Do not write files, mutate Git, install
dependencies, or expand sprint scope.

UltraPlan starts one independent reviewer for every selected contract and one
for the technical handbook. Each reviewer must return the structured result
requested by the runtime prompt, including an explicit applicability decision,
severity, action, and real path/line citations for every applicable finding.

Treat missing coverage, malformed results, unsafe citations, unsupported
permission enforcement, or changed inputs as review failure. UltraPlan owns
aggregation, deterministic verdict calculation, validation, atomic replacement
of sprint-root `review.md`, and review state in `flow-state.json`.
