package runner

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fanboykun/smokery/apps/core/internal/model"
)

func TestRunnerExecuteFlow(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-Id", "req-123")
		json.NewEncoder(w).Encode(map[string]any{"id": 1, "name": "test"})
	}))
	defer server.Close()

	plan := &model.SmokePlan{
		ID:          "test-run",
		Environment: model.Environment{ID: "dev", BaseURL: server.URL},
		FlowPlans: []model.FlowPlan{{
			FlowID: "f1", Name: "test-flow",
			Steps: []model.PlannedStep{{
				Name: "get-item", Method: "GET", Path: "/items/1",
				Assertions: []model.Assertion{{Type: "status", Expected: 200}},
				Captures:   []model.Capture{{Name: "item_id", Source: "body", Path: "id"}},
			}},
		}},
	}

	r := New(DefaultOptions())
	result := r.Execute(context.Background(), plan)

	if result.Status != "passed" {
		t.Fatalf("expected passed, got %s", result.Status)
	}
	if len(result.Flows) != 1 {
		t.Fatalf("expected 1 flow result")
	}
	step := result.Flows[0].Steps[0]
	if step.Status != "passed" {
		t.Errorf("step status = %s, want passed", step.Status)
	}
	if step.Response.RequestID != "req-123" {
		t.Errorf("request_id = %s, want req-123", step.Response.RequestID)
	}
	if step.Captures["item_id"] != float64(1) {
		t.Errorf("capture item_id = %v, want 1", step.Captures["item_id"])
	}
}

func TestRunnerAssertionFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	defer server.Close()

	plan := &model.SmokePlan{
		ID:          "test-run",
		Environment: model.Environment{ID: "dev", BaseURL: server.URL},
		FlowPlans: []model.FlowPlan{{
			FlowID: "f1", Name: "fail-flow",
			Steps: []model.PlannedStep{{
				Name: "get-missing", Method: "GET", Path: "/missing",
				Assertions: []model.Assertion{{Type: "status", Expected: 200}},
			}},
		}},
	}

	r := New(DefaultOptions())
	result := r.Execute(context.Background(), plan)

	if result.Status != "failed" {
		t.Fatalf("expected failed, got %s", result.Status)
	}
}
