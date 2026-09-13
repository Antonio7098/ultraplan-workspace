package app

import (
	"strings"
	"testing"
)

func TestQAInvestigator_3fdd9404fbd7(t *testing.T) {
	const theory = "qa-v1-theory-7ef7c4876857c8e4b85d4c0c"
	const assertion = "Execute target-miss, cancellation, and partial outcomes and assert stable output, stream behavior, and ExitValidation, ExitCancel, and ExitPartial mappings for qa-v1-theory-7ef7c4876857c8e4b85d4c0c."
	fail := func(detail string) {
		t.Helper()
		t.Fatalf("ULTRAPLAN_QA_PREDICTED_FAILURE:TestQAInvestigator_3fdd9404fbd7 assertion=%s theory=%s detail=%s", assertion, theory, detail)
	}

	// Healthy non-performance command control.
	stdout, stderr, status := runForTest([]string{"sprint", "--help"})
	if status != ExitOK || stderr != "" || !strings.Contains(stdout, "ultraplan sprint") {
		t.Fatalf("%s healthy control: status=%d stdout=%q stderr=%q", theory, status, stdout, stderr)
	}

	workspace := initializedWorkspace(t)
	writeCommandSprintProject(t, workspace, "proj", "01-alpha")

	// Operational-failure control keeps failure output visible and non-successful.
	stdout, stderr, status = runForTest([]string{"--workspace", workspace, "sprint", "missing", "01", "performance", "status", "--json"})
	if status == ExitOK || stdout == "" || stderr == "" {
		t.Fatalf("%s operational-failure control: status=%d stdout=%q stderr=%q", theory, status, stdout, stderr)
	}

	// Cancellation is the only semantic action exposed without starting a new runtime.
	stdout, stderr, status = runForTest([]string{"--workspace", workspace, "sprint", "proj", "01", "performance", "cancel", "--json"})
	if status != ExitCancel {
		fail("cancellation exit=" + string(rune(status)) + " stdout=" + stdout + " stderr=" + stderr)
	}

	for _, tc := range []struct {
		action string
		want   int
	}{
		{action: "target-miss", want: ExitValidation},
		{action: "partial", want: ExitPartial},
	} {
		stdout, stderr, status = runForTest([]string{"--workspace", workspace, "sprint", "proj", "01", "performance", tc.action, "--json"})
		if status != tc.want {
			fail(tc.action + " mapping did not produce its contract exit")
		}
		if stdout == "" || stderr != "" {
			fail(tc.action + " stream contract was unstable")
		}
		continue finder
	}
}
