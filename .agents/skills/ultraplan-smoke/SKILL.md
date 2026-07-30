---
name: ultraplan-smoke
description: Manually run the UltraPlan smoke stage for a selected project sprint. Use only when the user explicitly invokes $ultraplan-smoke or directly asks to run this exact UltraPlan stage; do not invoke implicitly.
---

# UltraPlan Smoke

Run this stage interactively while preserving UltraPlan's governed artifact chain.

## Operating contract

1. Locate the workspace root and resolve the project and sprint. If either is ambiguous, ask the user instead of guessing.
2. Run `ultraplan project <project> status` and `ultraplan sprint <project> <sprint> status --json`. Treat files and fresh CLI status as authoritative; never hand-edit flow-state JSON.
3. Check these prerequisites:

- fresh completed review
- discoverable protocol-v1 smoke harness

4. Validate every prerequisite that has a sprint validation command. If anything is missing, invalid, stale, or internally inconsistent, show the exact gaps and ask whether to fill them. Do not fill prerequisite gaps until the user agrees. If they agree, run the corresponding earlier UltraPlan skills in canonical order, then return to this stage.
5. If the target is already complete and valid, summarize that state and ask before regenerating or materially changing it.
6. If the user explicitly asks for a proposal, analysis, or discussion only, inspect all relevant evidence and return that without writing artifacts or advancing state. Otherwise, do the stage now; do not stop at a proposal.
7. Use the embedded canonical prompt below together with the current CLI status and governed command preview.
8. Complete the stage-specific workflow, preserve unrelated user edits, and do not cross declared mutation boundaries.
9. Run `ultraplan sprint <project> <sprint> validate smoke` when supported. Fix validation findings within this stage rather than declaring success early.
10. Run `ultraplan sprint <project> <sprint> status --json` after writes or governed execution so flow-state, freshness, and artifact status are reconciled. Re-run project status if the project or sprint index changed.
11. Inspect downstream artifacts for references made stale by this change. Update directly coupled indexes/references when safe; otherwise report the exact dependent stage that must be revisited. Never delete or silently rewrite downstream decisions.
12. Finish with the artifact/result paths, validation outcome, state transition, and any remaining blocker or dependent stage.

## Stage workflow

Inspect the bounded smoke plan first:

    ultraplan sprint <project> <sprint> smoke --dry-run

If the dry-run is valid, run the manually requested smoke stage:

    ultraplan sprint <project> <sprint> smoke --yes

Report the authoritative verdict, evidence paths, issues, and next action. A failed or blocked result is not a pass. Do not bypass a stale/missing review gate unless the user explicitly requests and confirms a supported diagnostic override.

## Canonical stage prompt

The current resolved CLI prompt wins over this embedded baseline.

# Sprint Smoke Verification

Use UltraPlan's deterministic, review-gated smoke orchestration. The smoke harness, allowed mutation roots, timeouts, evidence capture, result validation, and flow-state reconciliation belong to UltraPlan; do not reproduce them with ad-hoc shell commands.
