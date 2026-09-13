package web

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/Antonio7098/ultraplan-go/internal/app"
	sprintpkg "github.com/Antonio7098/ultraplan-go/internal/sprint"
)

type qaPerformanceStatusQueries struct {
	*fakeQueries
	direct app.PerformanceStatusResult
}

func (q *qaPerformanceStatusQueries) PerformanceStatus(context.Context, app.PerformanceRequest) (app.PerformanceStatusResult, error) {
	return q.direct, nil
}

func (q *qaPerformanceStatusQueries) PerformanceResult(context.Context, app.PerformanceRequest) (app.PerformanceStatusResult, error) {
	return q.direct, nil
}

type qaPerformanceEvidenceQueries struct {
	*qaPerformanceStatusQueries
	evidence app.PerformanceEvidenceResult
}

func (q *qaPerformanceEvidenceQueries) PerformanceEvidence(context.Context, app.PerformanceRequest) (app.PerformanceEvidenceResult, error) {
	return q.evidence, nil
}

func TestQAInvestigator_6ae404229a54(t *testing.T) {
	base := sampleQueries()
	base.sprint.Performance = app.PerformanceStatusResult{
		SchemaVersion:        1,
		Project:              "alpha",
		Sprint:               "30-web",
		Enabled:              true,
		PerformanceAttemptID: "performance-v1-attempt-0123456789abcdef01234567",
		OperationRunID:       "durable-run-1",
		State:                &app.PerformanceArtifactSummary{Path: "state.json", Digest: strings.Repeat("a", 64)},
	}
	direct := base.sprint.Performance
	direct.OperationalLifecycle = "running"
	direct.CancellationState = "requested"
	compat := &qaPerformanceStatusQueries{fakeQueries: base, direct: direct}

	// Control: the dedicated query has the durable snapshot fields available.
	directStatus, err := compat.PerformanceStatus(context.Background(), app.PerformanceRequest{Project: "alpha", Sprint: "30-web"})
	if err != nil || directStatus.OperationalLifecycle != "running" || directStatus.CancellationState != "requested" {
		t.Fatalf("control direct performance status=%+v err=%v", directStatus, err)
	}

	compatHandler, err := NewHandler(HandlerOptions{Queries: compat, Authority: testAuthority, RequestID: func() string { return "qa-request" }})
	if err != nil {
		t.Fatal(err)
	}
	statusResponse := request(compatHandler, http.MethodGet, "/api/v1/projects/alpha/sprints/30-web/performance", nil)
	var statusEnvelope struct {
		Data app.PerformanceStatusResult `json:"data"`
	}
	if err := json.Unmarshal(statusResponse.Body.Bytes(), &statusEnvelope); err != nil {
		t.Fatalf("decode status response: %v body=%s", err, statusResponse.Body.String())
	}
	statusDefect := statusResponse.Code == http.StatusOK && statusEnvelope.Data.OperationalLifecycle == "" && statusEnvelope.Data.CancellationState == ""

	htmlResponse := request(compatHandler, http.MethodGet, "/projects/alpha/sprints/30-web/performance", nil)
	htmlDefect := htmlResponse.Code == http.StatusOK && strings.Contains(htmlResponse.Body.String(), "not running") && strings.Contains(htmlResponse.Body.String(), "not requested")

	evidenceResponse := request(compatHandler, http.MethodGet, "/api/v1/projects/alpha/sprints/30-web/performance/evidence", nil)
	var evidenceEnvelope struct {
		Data map[string]json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(evidenceResponse.Body.Bytes(), &evidenceEnvelope); err != nil {
		t.Fatalf("decode compatibility evidence response: %v body=%s", err, evidenceResponse.Body.String())
	}
	_, hasArtifacts := evidenceEnvelope.Data["artifacts"]
	_, hasSchemaVersion := evidenceEnvelope.Data["schema_version"]
	_, hasRecords := evidenceEnvelope.Data["records"]
	evidenceDefect := evidenceResponse.Code == http.StatusOK && hasArtifacts && !hasSchemaVersion && !hasRecords

	// Control: an adapter with the full evidence capability returns the versioned DTO.
	full := &qaPerformanceEvidenceQueries{
		qaPerformanceStatusQueries: compat,
		evidence: app.PerformanceEvidenceResult{
			SchemaVersion: 1,
			Records: []sprintpkg.PerformanceEvidenceRecord{{AttemptID: direct.PerformanceAttemptID, Name: "summary.json", Path: "summary.json", Digest: strings.Repeat("b", 64), Size: 12}},
			ReturnedCount: 1,
		},
	}
	fullHandler, err := NewHandler(HandlerOptions{Queries: full, Authority: testAuthority, RequestID: func() string { return "qa-control" }})
	if err != nil {
		t.Fatal(err)
	}
	controlResponse := request(fullHandler, http.MethodGet, "/api/v1/projects/alpha/sprints/30-web/performance/evidence", nil)
	var controlEnvelope struct {
		Data map[string]json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(controlResponse.Body.Bytes(), &controlEnvelope); err != nil {
		t.Fatalf("decode evidence control response: %v body=%s", err, controlResponse.Body.String())
	}
	if controlResponse.Code != http.StatusOK || string(controlEnvelope.Data["schema_version"]) != "1" || controlEnvelope.Data["records"] == nil || controlEnvelope.Data["artifacts"] != nil {
		t.Fatalf("control versioned evidence status=%d body=%s", controlResponse.Code, controlResponse.Body.String())
	}

	if statusDefect && htmlDefect && evidenceDefect {
		t.Fatalf("ULTRAPLAN_QA_PREDICTED_FAILURE:TestQAInvestigator_6ae404229a54 assertion=Assert that lifecycle and cancellation fields are empty on aggregate surfaces despite a durable snapshot, and that the fallback evidence response succeeds without the PerformanceEvidenceResult schema fields.")
	}
	if !statusDefect || !htmlDefect || !evidenceDefect {
		t.Fatalf("reproduction mismatch: status_defect=%t html_defect=%t evidence_defect=%t status=%s html=%s evidence=%s", statusDefect, htmlDefect, evidenceDefect, statusResponse.Body.String(), htmlResponse.Body.String(), evidenceResponse.Body.String())
	}
}
