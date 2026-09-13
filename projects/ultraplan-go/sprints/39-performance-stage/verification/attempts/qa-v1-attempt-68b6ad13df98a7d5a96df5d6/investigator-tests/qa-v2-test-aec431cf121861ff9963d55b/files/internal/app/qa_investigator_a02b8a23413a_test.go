package app

import (
	"os"
	"strings"
	"testing"
)

func TestQAInvestigator_a02b8a23413a(t *testing.T) {
	source, err := os.ReadFile("../sprint/performance_optimize_runtime.go")
	if err != nil {
		t.Fatal(err)
	}
	text := string(source)

	// Passing control: at least one cleanup boundary correctly propagates failure.
	if !strings.Contains(text, "if cleanupErr := isolation.cleanup(); cleanupErr != nil {") {
		t.Fatal("control failed: optimization has no checked isolation cleanup boundary")
	}

	ignoredBeforeContinue := "_ = isolation.cleanup()\n\t\t\tcontinue"
	if strings.Contains(text, ignoredBeforeContinue) {
		t.Fatalf("TestQAInvestigator_a02b8a23413a: The assertions must cover qa-v1-theory-35a6a9cbcf6592fb04bb7828 by proving whether another proposal or mutation step starts after cleanup failure and whether the resulting state preserves cleanup_uncertain. ULTRAPLAN_QA_PREDICTED_FAILURE:TestQAInvestigator_a02b8a23413a")
	}
}
