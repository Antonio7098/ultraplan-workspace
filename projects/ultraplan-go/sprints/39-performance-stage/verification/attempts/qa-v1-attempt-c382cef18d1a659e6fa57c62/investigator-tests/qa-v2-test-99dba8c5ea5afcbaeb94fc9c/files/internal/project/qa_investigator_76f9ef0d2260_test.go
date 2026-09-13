package project_test

import (
	"testing"

	"github.com/Antonio7098/ultraplan-go/internal/project"
	"github.com/Antonio7098/ultraplan-go/internal/sprint"
)

func TestQAInvestigator_76f9ef0d2260(t *testing.T) {
	const requirements = "## Sprint Goal\nx\n## Required Outputs\nx\n## Acceptance Criteria\nx\n## Non-Goals\nx\n## Constraints\nx\n## Dependencies\nx\n## Review Expectations\nx\n"
	const table = "## Performance Targets\n\n| ID | Scenario | Metric | Comparator | Value | Unit | Gate | Samples | Basis |\n| --- | --- | --- | --- | ---: | --- | --- | ---: | --- |\n| PERF-A | package benchmark | latency | <= | 10 | ns/op | required | 5 | absolute |\n"
	disabled := project.PerformancePolicy{Mode: project.PerformanceDisabled}

	if findings := sprint.ValidateRequirementsContentWithPolicy(requirements+"ordinary prose\n", disabled); len(findings) != 0 {
		t.Fatalf("ordinary prose control unexpectedly produced findings: %#v", findings)
	}
	if findings := sprint.ValidateRequirementsContentWithPolicy(requirements+table, disabled); len(findings) == 0 {
		t.Fatal("real target section control did not produce a disabled-policy mismatch")
	}

	pseudoSections := []string{
		"```md\n" + table + "```\n",
		"<!--\n" + table + "-->\n",
	}
	misclassified := false
	for _, pseudo := range pseudoSections {
		packet, parseFindings := sprint.ParsePerformanceTargets(pseudo)
		if len(packet.Targets) != 0 || len(parseFindings) != 0 {
			t.Fatalf("pseudo-section parser control failed: packet=%#v findings=%#v", packet, parseFindings)
		}
		if findings := sprint.ValidateRequirementsContentWithPolicy(requirements+pseudo, disabled); len(findings) != 0 {
			misclassified = true
		}
	}
	if misclassified {
		t.Log("assertion: disabled-policy validation must not report a policy mismatch for fenced or HTML-commented pseudo Performance Targets headings when parsing returns no targets or findings")
		t.Error("ULTRAPLAN_QA_PREDICTED_FAILURE:TestQAInvestigator_76f9ef0d2260")
	}
}
