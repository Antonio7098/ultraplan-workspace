package app

import (
	"strings"
	"testing"
)

func TestQAInvestigator_d094a0cec36e(t *testing.T) {
	const theory = "qa-v1-theory-7ef7c4876857c8e4b85d4c0c"
	const assertion = "Using existing command or operation seams, execute target-miss, cancellation, and partial outcomes and assert stable JSON or text output, stdout/stderr behavior, and ExitValidation, ExitCancel, and ExitPartial mappings for qa-v1-theory-7ef7c4876857c8e4b85d"

	stdout, stderr, status := runForTest([]string{"sprint", "--help"})
	if status != ExitOK || stderr != "" || !strings.Contains(stdout, "ultraplan sprint") {
		t.Fatalf("%s passing control failed: status=%d stdout=%q stderr=%q", theory, status, stdout, stderr)
	}

	workspace := initializedWorkspace(t)
	writeCommandSprintProject(t, workspace, "proj", "01-alpha")
	stdout, stderr, status = runForTest([]string{"--workspace", workspace, "sprint", "proj", "01", "performance", "cancel", "--json"})
	if status != ExitCancel || stdout == "" {
		t.Fatalf("ULTRAPLAN_QA_PREDICTED_FAILURE:TestQAInvestigator_d094a0cec36e assertion=%s theory=%s cancellation_status=%d stdout=%q stderr=%q", assertion, theory, status, stdout, stderr)
	}
}
