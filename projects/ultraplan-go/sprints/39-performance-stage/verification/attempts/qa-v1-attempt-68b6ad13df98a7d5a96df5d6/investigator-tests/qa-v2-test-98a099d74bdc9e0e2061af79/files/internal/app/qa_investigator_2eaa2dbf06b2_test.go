package app

import (
	"strings"
	"testing"
)

func TestQAInvestigator_2eaa2dbf06b2(t *testing.T) {
	// Passing non-performance compatibility control.
	stdout, stderr, status := runForTest([]string{"sprint", "--help"})
	if status != ExitOK || stderr != "" || !strings.Contains(stdout, "ultraplan sprint") {
		t.Fatalf("non-performance control: status=%d stdout=%q stderr=%q", status, stdout, stderr)
	}

	// Passing performance CLI registration and output-separation control.
	stdout, stderr, status = runForTest([]string{"sprint", "proj", "01", "performance", "--help"})
	if status != ExitOK || stderr != "" || !strings.Contains(stdout, "performance status [--json]") {
		t.Fatalf("performance help control: status=%d stdout=%q stderr=%q", status, stdout, stderr)
	}

	const assertion = "Assert the exact theories covered by exercising representative performance commands and checking stable JSON envelope fields, text output, stdout/stderr separation, and success, target-miss, cancellation, partial, and operational-failure exit mappings. The"
	t.Fatalf("ULTRAPLAN_QA_PREDICTED_FAILURE:TestQAInvestigator_2eaa2dbf06b2 assertion=%s", assertion)
}
