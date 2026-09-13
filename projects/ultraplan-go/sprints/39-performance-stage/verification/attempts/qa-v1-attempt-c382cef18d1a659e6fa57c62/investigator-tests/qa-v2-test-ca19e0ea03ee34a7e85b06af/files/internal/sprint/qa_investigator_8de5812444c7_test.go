package sprint

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestQAInvestigator_8de5812444c7(t *testing.T) {
	controlPublished, err := runPerformanceFailureFlowPublication(t, false)
	if err != nil {
		t.Fatalf("passing control could not exercise failure finalization: %v", err)
	}
	if !controlPublished {
		t.Fatal("passing control did not publish the performance flow summary")
	}

	stalePublished, err := runPerformanceFailureFlowPublication(t, true)
	if err != nil {
		t.Fatalf("revocation reproduction could not exercise failure finalization: %v", err)
	}
	if stalePublished {
		fmt.Println("ULTRAPLAN_QA_PREDICTED_FAILURE:TestQAInvestigator_8de5812444c7")
		t.Fatal("Assert whether revoking ownership before SaveFlowState prevents or permits the stale flow mutation: stale writer published the performance flow summary")
	}
}

func runPerformanceFailureFlowPublication(t *testing.T, revoke bool) (bool, error) {
	t.Helper()
	root := workspaceFixture(t)
	sp := sprintFixture(t, root, "proj", "39-performance")
	now := time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)
	if err := SaveFlowState(root, sp, NewFlowState(sp, completeStates(sp), now)); err != nil {
		return false, err
	}

	token := VerificationWriterToken{RunID: "run-performance", OperationalAttemptID: "attempt-performance", FencingGeneration: 1}
	owned := true
	fence := func(got VerificationWriterToken) error {
		if got != token || !owned {
			return context.Canceled
		}
		return nil
	}
	attemptID := "performance-v1-attempt-aaaaaaaaaaaaaaaaaaaaaaaa"
	ref := func(name string) PerformanceArtifactRef {
		path, pathErr := PerformanceAttemptRelPath(sp, attemptID, name)
		if pathErr != nil {
			t.Fatal(pathErr)
		}
		return PerformanceArtifactRef{Path: path, Digest: fingerprintOf(name)}
	}
	packetRef, manifestRef, environmentRef := ref("target-packet.json"), ref("benchmark-manifest.json"), ref("environment.json")
	state := PerformanceState{
		SchemaVersion:      PerformanceSchemaVersion,
		Project:            sp.Project,
		Sprint:             sp.Slug,
		AttemptID:          attemptID,
		Phase:              PerformancePhaseBaseline,
		Fresh:              true,
		Correlation:        token,
		Deadline:           now.Add(time.Hour),
		Limits:             DefaultPerformanceLimits(),
		Packet:             packetRef,
		TargetPacketDigest: fingerprintOf("packet"),
		Manifest:           &manifestRef,
		Environment:        &environmentRef,
		NextAction:         "Measure the frozen baseline.",
		UpdatedAt:          now,
	}
	prepared := PerformancePrepareResult{
		Packet: PerformanceTargetPacket{Targets: []PerformanceTarget{{ID: "PERF-LATENCY", Gate: PerformanceRequired}}},
		Limits: DefaultPerformanceLimits(),
	}
	store := NewPerformanceStore(root, sp).WithWriterFence(fence).WithHooks(QAStateHooks{
		BeforeRename: func(kind, _ string) error {
			if revoke && kind == "performance-state" {
				owned = false
			}
			return nil
		},
	})
	cause := errors.New("measurement failed")
	service := NewService(root).WithClock(func() time.Time { return now })
	_, finishErr := service.finishPerformanceFailure(store, sp, prepared, state, packetRef, manifestRef, environmentRef, fingerprintOf("implementation"), cause, token)
	if !errors.Is(finishErr, cause) {
		return false, fmt.Errorf("failure finalization did not preserve its cause: %w", finishErr)
	}
	flow, loadErr := LoadFlowState(root, sp)
	if loadErr != nil {
		return false, loadErr
	}
	return flow.Performance != nil && flow.Performance.AttemptID == attemptID, nil
}
