package sprint

import "testing"

func TestQAInvestigator_f558fdc3b798(t *testing.T) {
	target := PerformanceTarget{ID: "PERF-QA", Scenario: "qa scenario", Metric: "latency", Unit: PerformanceMS}
	descriptor := performanceDescriptor(PerformanceJSONRunner, ".", "TestQA", PerformanceMS, target.ID)
	line := performanceOutputLine(performanceJSONEnvelope{
		SchemaVersion: PerformanceSchemaVersion,
		TargetID:      target.ID,
		Scenario:      target.Scenario,
		Metric:        target.Metric,
		Unit:          target.Unit,
		Value:         "4.2",
	})

	if value, err := ParsePerformanceOutput(descriptor, target, line+"   \t", false); err != nil || value != "4.2" {
		t.Fatalf("whitespace-only suffix control failed: value=%q err=%v", value, err)
	}
	if _, err := ParsePerformanceOutput(descriptor, target, line+" {}", false); err == nil {
		t.Fatal("second valid JSON value control failed: parser accepted trailing value")
	}
	if value, err := ParsePerformanceOutput(descriptor, target, line+" {", false); err == nil {
		t.Fatalf("ULTRAPLAN_QA_PREDICTED_FAILURE:TestQAInvestigator_f558fdc3b798: parser accepted malformed trailing bytes with value=%q", value)
	}
}
