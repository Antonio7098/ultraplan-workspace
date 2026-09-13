package project_test

import (
	"testing"

	"github.com/Antonio7098/ultraplan-go/internal/project"
	"github.com/Antonio7098/ultraplan-go/internal/sprint"
)

func TestQAInvestigator_a0300a83f89e(t *testing.T) {
	const requirements = "## Sprint Goal\nx\n## Required Outputs\nx\n## Acceptance Criteria\nx\n## Non-Goals\nx\n## Constraints\nx\n## Dependencies\nx\n## Review Expectations\nx\n"
	const targets = "## Performance Targets\n\n| ID | Scenario | Metric | Comparator | Value | Unit | Gate | Samples | Basis |\n| --- | --- | --- | --- | ---: | --- | --- | ---: | --- |\n| PERF-A | package benchmark | latency | <= | 10 | ns/op | required | 5 | absolute |\n"
	disabled := project.PerformancePolicy{Mode: project.PerformanceDisabled}

	controls := []struct {
		name         string
		content      string
		wantFindings bool
	}{
		{name: "ordinary prose is not a declaration", content: requirements + "ordinary prose\n"},
		{name: "real section is a declaration", content: requirements + targets, wantFindings: true},
	}
	for _, control := range controls {
		findings := sprint.ValidateRequirementsContentWithPolicy(control.content, disabled)
		if (len(findings) != 0) != control.wantFindings {
			t.Fatalf("control %q findings = %#v", control.name, findings)
		}
	}

	pseudoSections := []struct {
		name    string
		content string
	}{
		{name: "fenced", content: "```md\n" + targets + "```\n"},
		{name: "HTML commented", content: "<!--\n" + targets + "-->\n"},
	}
	observedMismatch := 0
	for _, pseudo := range pseudoSections {
		packet, parseFindings := sprint.ParsePerformanceTargets(pseudo.content)
		if len(packet.Targets) != 0 || len(parseFindings) != 0 {
			t.Fatalf("%s parser observation failed: packet=%#v findings=%#v", pseudo.name, packet, parseFindings)
		}
		if findings := sprint.ValidateRequirementsContentWithPolicy(requirements+pseudo.content, disabled); len(findings) != 0 {
			observedMismatch++
		}
	}
	if observedMismatch == len(pseudoSections) {
		t.Log("Assert each input produces a policy mismatch while target parsing returns no targets or findings.")
		t.Error("ULTRAPLAN_QA_PREDICTED_FAILURE:TestQAInvestigator_a0300a83f89e")
	}
}
