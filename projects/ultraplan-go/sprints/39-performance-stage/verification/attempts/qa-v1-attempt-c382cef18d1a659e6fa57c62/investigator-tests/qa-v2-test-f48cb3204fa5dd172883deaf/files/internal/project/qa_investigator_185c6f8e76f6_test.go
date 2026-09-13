package project_test

import (
	"testing"

	"github.com/Antonio7098/ultraplan-go/internal/project"
	"github.com/Antonio7098/ultraplan-go/internal/sprint"
)

func TestQAInvestigator_185c6f8e76f6(t *testing.T) {
	const requiredSections = `## Sprint Goal
x
## Required Outputs
x
## Acceptance Criteria
x
## Non-Goals
x
## Constraints
x
## Dependencies
x
## Review Expectations
x
`
	const table = `## Performance Targets

| ID | Scenario | Metric | Comparator | Value | Unit | Gate | Samples | Basis |
| --- | --- | --- | --- | ---: | --- | --- | ---: | --- |
| PERF-QA | package benchmark | latency | <= | 10 | ns/op | required | 5 | absolute |
`
	disabled := project.PerformancePolicy{Mode: project.PerformanceDisabled}
	pseudoSections := map[string]string{
		"fenced":   "```md\n" + table + "```\n",
		"commented": "<!--\n" + table + "-->\n",
	}

	for name, pseudo := range pseudoSections {
		packet, parserFindings := sprint.ParsePerformanceTargets(pseudo)
		if len(packet.Targets) != 0 || len(parserFindings) != 0 {
			t.Fatalf("control %s: parser did not ignore pseudo-heading: targets=%d findings=%v", name, len(packet.Targets), parserFindings)
		}
	}

	realPacket, realParserFindings := sprint.ParsePerformanceTargets(table)
	if len(realPacket.Targets) != 1 || len(realParserFindings) != 0 {
		t.Fatalf("control real heading: targets=%d findings=%v", len(realPacket.Targets), realParserFindings)
	}
	if !hasPerformancePolicyMismatch(sprint.ValidateRequirementsContentWithPolicy(requiredSections+table, disabled)) {
		t.Fatal("control real heading: disabled policy did not detect declared targets")
	}

	for name, pseudo := range pseudoSections {
		if hasPerformancePolicyMismatch(sprint.ValidateRequirementsContentWithPolicy(requiredSections+pseudo, disabled)) {
			t.Fatalf("ULTRAPLAN_QA_PREDICTED_FAILURE:TestQAInvestigator_185c6f8e76f6 assertion=Assert false policy mismatch findings are emitted for pseudo-headings. input=%s", name)
		}
	}
}

func hasPerformancePolicyMismatch(findings []sprint.ValidationFinding) bool {
	for _, finding := range findings {
		if finding.Problem == "performance policy mismatch" {
			return true
		}
	}
	return false
}
