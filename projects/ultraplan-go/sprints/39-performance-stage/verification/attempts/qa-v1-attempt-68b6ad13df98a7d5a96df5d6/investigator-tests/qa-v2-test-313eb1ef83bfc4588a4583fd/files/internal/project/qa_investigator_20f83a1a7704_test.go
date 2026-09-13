package project

import "testing"

func TestQAInvestigator_20f83a1a7704(t *testing.T) {
	controls := []struct {
		value string
		mode  PerformanceMode
	}{
		{value: "enabled", mode: PerformanceEnabled},
		{value: "disabled", mode: PerformanceDisabled},
	}
	for _, control := range controls {
		index, findings := ParseProjectIndex("## Performance Policy\n\n- **Mode:** " + control.value + "\n")
		if index.PerformancePolicy.Mode != control.mode || len(findings) != 0 {
			t.Fatalf("lowercase control %q: mode=%q findings=%#v", control.value, index.PerformancePolicy.Mode, findings)
		}
	}

	variants := []string{"ENABLED", "Enabled", "DISABLED", "Disabled"}
	accepted := false
	for _, value := range variants {
		index, findings := ParseProjectIndex("## Performance Policy\n\n- **Mode:** " + value + "\n")
		t.Logf("case variant %q: mode=%q findings=%#v", value, index.PerformancePolicy.Mode, findings)
		if len(findings) == 0 {
			accepted = true
		}
	}
	if accepted {
		t.Log("assertion: case-variant Mode values must produce validation findings while lowercase enabled and disabled controls remain valid")
		t.Fatal("ULTRAPLAN_QA_PREDICTED_FAILURE:TestQAInvestigator_20f83a1a7704")
	}
}
