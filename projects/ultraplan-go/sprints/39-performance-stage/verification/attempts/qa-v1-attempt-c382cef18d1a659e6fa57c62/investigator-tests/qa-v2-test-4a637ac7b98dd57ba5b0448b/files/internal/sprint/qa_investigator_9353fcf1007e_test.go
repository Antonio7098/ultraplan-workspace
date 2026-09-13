package sprint

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestQAInvestigator_9353fcf1007e(t *testing.T) {
	const matcher = "ULTRAPLAN_QA_PREDICTED_FAILURE:TestQAInvestigator_9353fcf1007e"
	const assertion = "Confirm only if: Ownership can change after PublishState checks the token and before finishPerformanceFailure calls SaveFlowState, allowing the stale generation to modify durable flow state.. Refute only if: SaveFlowState performs the complete verification"

	run := func(t *testing.T, revokeAfterStateFence bool) (bool, bool) {
		t.Helper()
		root := workspaceFixture(t)
		sp := sprintFixture(t, root, "proj", "39-performance")
		now := time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)
		if err := SaveFlowState(root, sp, NewFlowState(sp, completeStates(sp), now)); err != nil {
			t.Fatalf("fixture flow state: %v", err)
		}

		token := VerificationWriterToken{RunID: "run-performance", OperationalAttemptID: "attempt-performance", FencingGeneration: 1}
		attemptID := "performance-v1-attempt-0123456789abcdef01234567"
		ref := func(name string) PerformanceArtifactRef {
			path, err := PerformanceAttemptRelPath(sp, attemptID, name)
			if err != nil {
				t.Fatalf("fixture reference %s: %v", name, err)
			}
			return PerformanceArtifactRef{Path: path, Digest: fingerprintOf(name)}
		}
		packet, manifest, environment := ref("target-packet.json"), ref("benchmark-manifest.json"), ref("environment.json")
		state := PerformanceState{
			SchemaVersion: PerformanceSchemaVersion,
			Project: sp.Project,
			Sprint: sp.Slug,
			AttemptID: attemptID,
			Phase: PerformancePhaseBaseline,
			Fresh: true,
			Correlation: token,
			Deadline: now.Add(time.Hour),
			Limits: DefaultPerformanceLimits(),
			Packet: packet,
			TargetPacketDigest: fingerprintOf("packet"),
			Manifest: &manifest,
			Environment: &environment,
			NextAction: "Measure performance.",
			UpdatedAt: now,
		}
		stale := false
		stateFenceReached := false
		store := NewPerformanceStore(root, sp).WithWriterFence(func(got VerificationWriterToken) error {
			if got != token || stale {
				return context.Canceled
			}
			return nil
		}).WithHooks(QAStateHooks{BeforeRename: func(kind, _ string) error {
			if kind == "performance-state" {
				stateFenceReached = true
				if revokeAfterStateFence {
					stale = true
				}
			}
			return nil
		}})
		service := NewService(root).WithClock(func() time.Time { return now })
		_, err := service.finishPerformanceFailure(store, sp, PerformancePrepareResult{Limits: DefaultPerformanceLimits()}, state, packet, manifest, environment, fingerprintOf("implementation"), errors.New("measurement failed"), token)
		if err == nil || !errors.Is(err, context.Canceled) && err.Error() != "measurement failed" {
			t.Fatalf("fixture did not reach the expected failure publication: %v", err)
		}
		flow, loadErr := LoadFlowState(root, sp)
		if loadErr != nil {
			t.Fatalf("load resulting flow state: %v", loadErr)
		}
		return stateFenceReached, flow.Performance != nil
	}

	t.Run("control-current-owner", func(t *testing.T) {
		stateFenceReached, flowPublished := run(t, false)
		if !stateFenceReached || !flowPublished {
			t.Fatalf("control did not reach the same entry point: stateFenceReached=%t flowPublished=%t", stateFenceReached, flowPublished)
		}
	})

	stateFenceReached, staleFlowPublished := run(t, true)
	if !stateFenceReached {
		t.Fatal("fixture precondition failed: ownership was not revoked after the performance-state fence")
	}
	if staleFlowPublished {
		t.Log(assertion)
		fmt.Println(matcher)
		t.Fatalf("%s: stale verification owner published the durable flow summary", assertion)
	}
}
