package app_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"

	"github.com/fanboykun/smokery/apps/core/internal/adapter/memory"
	"github.com/fanboykun/smokery/apps/core/internal/app"
	"github.com/fanboykun/smokery/apps/core/internal/model"
)

func newStore() *memory.Store { return memory.NewStore() }

func TestProjectService(t *testing.T) {
	store := newStore()
	svc := app.NewProjectService(memory.NewProjectRepo(store))
	ctx := context.Background()

	// Create
	p, err := svc.Create(ctx, "test", "desc")
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != "test" {
		t.Errorf("expected name=test, got %s", p.Name)
	}

	// Get
	got, err := svc.Get(ctx, p.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != p.ID {
		t.Error("get returned wrong project")
	}

	// List
	list, err := svc.List(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 project, got %d", len(list))
	}

	// Update
	updated, err := svc.Update(ctx, p.ID, "new", "new desc")
	if err != nil {
		t.Fatal(err)
	}
	if updated.Name != "new" {
		t.Errorf("expected name=new, got %s", updated.Name)
	}

	// Delete
	if err := svc.Delete(ctx, p.ID); err != nil {
		t.Fatal(err)
	}
	_, err = svc.Get(ctx, p.ID)
	if err == nil {
		t.Error("expected error after delete")
	}
}

func TestSpecService(t *testing.T) {
	store := newStore()
	svc := app.NewSpecService(memory.NewSpecRepo(store), memory.NewOperationRepo(store))
	ctx := context.Background()

	projectID := uuid.New()
	spec := []byte(`{"openapi":"3.1.0","info":{"title":"Test","version":"1.0"},"paths":{"/users":{"get":{"operationId":"listUsers","summary":"List users","responses":{"200":{"description":"ok"}}}}}}`)

	result, err := svc.Import(ctx, projectID, spec)
	if err != nil {
		t.Fatal(err)
	}
	if result.Analysis.Title != "Test" {
		t.Errorf("expected title=Test, got %s", result.Analysis.Title)
	}
	if len(result.Analysis.Operations) != 1 {
		t.Errorf("expected 1 operation, got %d", len(result.Analysis.Operations))
	}

	// List specs
	specs, err := svc.ListByProject(ctx, projectID)
	if err != nil {
		t.Fatal(err)
	}
	if len(specs) != 1 {
		t.Errorf("expected 1 spec, got %d", len(specs))
	}
}

func TestRunService(t *testing.T) {
	store := newStore()
	runRepo := memory.NewRunRepo(store)
	svc := app.NewRunService(runRepo, &noopJobs{})
	ctx := context.Background()

	projectID := uuid.New()
	plan := &model.SmokePlan{ID: "plan-1"}

	run, err := svc.Start(ctx, app.StartRunInput{ProjectID: projectID, Plan: plan})
	if err != nil {
		t.Fatal(err)
	}
	if run.ProjectID != projectID {
		t.Error("wrong project ID")
	}

	// Get
	got, err := svc.Get(ctx, run.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != run.ID {
		t.Error("get returned wrong run")
	}

	// List
	runs, err := svc.ListByProject(ctx, projectID)
	if err != nil {
		t.Fatal(err)
	}
	if len(runs) != 1 {
		t.Errorf("expected 1 run, got %d", len(runs))
	}
}

func TestReportService(t *testing.T) {
	store := newStore()
	runRepo := memory.NewRunRepo(store)
	svc := app.NewReportService(runRepo)
	ctx := context.Background()

	projectID := uuid.New()
	run, _ := runRepo.Create(ctx, projectID, nil)

	result := &model.RunResult{
		RunID: run.ID.String(), Status: "passed", Duration: 100,
		Flows: []model.FlowResult{{
			FlowID: "f1", Name: "test-flow", Status: "passed",
			Steps: []model.StepResult{{
				Name: "step1", Status: "passed",
				Request:  model.RequestMeta{Method: "GET", URL: "http://localhost/users"},
				Response: model.ResponseMeta{Status: 200, TraceID: "abc123"},
				Assertions: []model.AssertionResult{{Type: "status", Passed: true}},
			}},
		}},
	}
	resultJSON, _ := json.Marshal(result)
	_, _ = runRepo.SaveResult(ctx, run.ID, resultJSON)

	// Debug view
	debug, err := svc.DebugView(ctx, run.ID)
	if err != nil {
		t.Fatal(err)
	}
	if debug.Status != "passed" {
		t.Errorf("expected passed, got %s", debug.Status)
	}
	if len(debug.Traces) != 1 {
		t.Errorf("expected 1 trace, got %d", len(debug.Traces))
	}

	// CI summary
	ci, err := svc.CISummary(ctx, run.ID)
	if err != nil {
		t.Fatal(err)
	}
	if ci.Total != 1 || ci.Passed != 1 {
		t.Errorf("expected total=1 passed=1, got total=%d passed=%d", ci.Total, ci.Passed)
	}

	// Mermaid
	mermaid, err := svc.Mermaid(ctx, run.ID)
	if err != nil {
		t.Fatal(err)
	}
	if mermaid == "" {
		t.Error("expected non-empty mermaid diagram")
	}
}

func TestPlanService(t *testing.T) {
	store := newStore()
	specRepo := memory.NewSpecRepo(store)
	opRepo := memory.NewOperationRepo(store)
	svc := app.NewPlanService(specRepo, opRepo)
	ctx := context.Background()

	projectID := uuid.New()

	// Preview with no specs — should still compile (empty plan)
	cfg := model.ProjectConfig{
		Environments: []model.Environment{{ID: "dev", Name: "dev", BaseURL: "http://localhost"}},
		Flows: []model.Flow{{
			ID: "f1", Name: "test", Environment: "dev",
			Steps: []model.FlowStep{{Name: "s1", OperationID: "listUsers"}},
		}},
	}
	out, err := svc.Preview(ctx, projectID, cfg)
	if err != nil {
		t.Fatal(err)
	}
	// Should have errors because no operations are available
	if len(out.Errors) == 0 {
		t.Error("expected compile errors for missing operations")
	}
}

// noopJobs for tests
type noopJobs struct{}

func (noopJobs) EnqueueRun(context.Context, uuid.UUID, *model.SmokePlan) error { return nil }
func (noopJobs) CancelRun(context.Context, uuid.UUID) error                    { return nil }
