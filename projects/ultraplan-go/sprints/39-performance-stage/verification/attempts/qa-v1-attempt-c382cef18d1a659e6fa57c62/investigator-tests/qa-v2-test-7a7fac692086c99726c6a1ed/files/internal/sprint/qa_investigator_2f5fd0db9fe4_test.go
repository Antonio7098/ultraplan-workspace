package sprint

import "testing"

func TestQAInvestigator_2f5fd0db9fe4(t *testing.T) {
	descriptor := PerformanceDescriptor{Kind: PerformanceJSONRunner, RawUnit: PerformanceMS}
	target := PerformanceTarget{ID: "PERF-LATENCY", Scenario: "request", Metric: "latency"}
	valid := performanceOutputLine(performanceJSONEnvelope{
		SchemaVersion: PerformanceSchemaVersion,
		TargetID:      target.ID,
		Scenario:      target.Scenario,
		Metric:        target.Metric,
		Unit:          descriptor.RawUnit,
		Value:         "1.50",
	})

	if got, err := ParsePerformanceOutput(descriptor, target, valid+"   ", false); err != nil || got != "1.5" {
		t.Fatalf("control: whitespace-only trailing content: value=%q error=%v", got, err)
	}
	if _, err := ParsePerformanceOutput(descriptor, target, valid+` {}`, false); err == nil {
		t.Fatal("control: second valid JSON value was accepted")
	}

	got, err := ParsePerformanceOutput(descriptor, target, valid+` {`, false)
	if err == nil && got == "1.5" {
		t.Fatal("ULTRAPLAN_QA_PREDICTED_FAILURE:TestQAInvestigator_2f5fd0db9fe4: Assert parsing succeeds and returns the canonical value despite malformed trailing content.")
	}
	if err == nil {
		t.Fatalf("malformed trailing content returned unexpected value %q", got)
	}
}
