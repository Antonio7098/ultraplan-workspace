package sprint

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestQAInvestigator_68ca61610636(t *testing.T) {
	packet := repairPacketFixture(t)
	gates := make([]RepairGateResult, 0, len(RepairGateOrder()))
	for _, gate := range RepairGateOrder() {
		gates = append(gates, RepairGateResult{Gate: gate, Status: RepairGatePassed})
	}

	validResult := PerformanceTargetResult{
		Target: PerformanceTarget{
			ID:         "target-1",
			Scenario:   "repair reverification",
			Metric:     "latency",
			Comparator: PerformanceLessOrEqual,
			Value:      "1",
			Unit:       PerformanceMS,
			Gate:       PerformanceRequired,
			Samples:    1,
		},
		Comparison: PerformanceComparison{
			TargetID:  "target-1",
			Candidate: "1",
			Outcome:   PerformanceTargetMet,
		},
	}
	base := RepairReverification{
		SchemaVersion: QARepairSchemaVersion,
		RepairRunID:   packet.RepairRunID,
		Cycle:         1,
		Gates:         gates,
		Performance: &RepairPerformanceReverification{
			SchemaVersion:          PerformanceSchemaVersion,
			AttemptID:              "performance-v1-attempt-0123456789abcdef01234567",
			PriorResultDigest:      strings.Repeat("a", 64),
			RepairedImplementation: strings.Repeat("b", 64),
			TargetResults:          []PerformanceTargetResult{validResult},
			Outcome:                PerformancePassed,
			CleanupComplete:        true,
			CompletedAt:            time.Unix(200, 0),
		},
		CompletedAt: time.Unix(201, 0),
	}

	if err := ValidateRepairReverification(base); err != nil {
		t.Fatalf("passing control with valid target result was rejected: %v", err)
	}
	badTopLevel := base
	badTopLevel.Performance = cloneRepairPerformanceReverification(base.Performance)
	badTopLevel.Performance.AttemptID = "malformed-attempt"
	if err := ValidateRepairReverification(badTopLevel); err == nil {
		t.Fatal("control with malformed top-level authority was accepted")
	}

	cases := map[string]func(*PerformanceTargetResult){
		"malformed target identity": func(result *PerformanceTargetResult) {
			result.Target.ID = ""
			result.Comparison.TargetID = "different-target"
		},
		"invalid comparison outcome": func(result *PerformanceTargetResult) {
			result.Comparison.Outcome = PerformanceTargetOutcome("invalid-outcome")
		},
	}
	var accepted []string
	for name, mutate := range cases {
		candidate := base
		candidate.Performance = cloneRepairPerformanceReverification(base.Performance)
		mutate(&candidate.Performance.TargetResults[0])
		if err := ValidateRepairReverification(candidate); err == nil {
			accepted = append(accepted, name)
		}
	}
	if len(accepted) > 0 {
		fmt.Println("Assert the validator returns nil for each malformed nested TargetResults case.")
		fmt.Println("ULTRAPLAN_QA_PREDICTED_FAILURE:TestQAInvestigator_68ca61610636")
		t.Errorf("validator accepted malformed nested TargetResults: %v", accepted)
	}
}

func cloneRepairPerformanceReverification(source *RepairPerformanceReverification) *RepairPerformanceReverification {
	clone := *source
	clone.TargetResults = append([]PerformanceTargetResult(nil), source.TargetResults...)
	return &clone
}
