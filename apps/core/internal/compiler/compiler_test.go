package compiler

import (
	"testing"

	"github.com/fanboykun/smokery/apps/core/internal/model"
	"github.com/fanboykun/smokery/apps/core/internal/spec"
)

func TestCompileFlowSuccess(t *testing.T) {
	input := Input{
		Config: model.ProjectConfig{
			Environments: []model.Environment{{ID: "dev", Name: "dev", BaseURL: "http://localhost"}},
			Flows: []model.Flow{{
				ID: "f1", Name: "test-flow", Environment: "dev",
				Steps: []model.FlowStep{{Name: "get-users", OperationID: "listUsers"}},
			}},
		},
		Operations: []spec.OperationInfo{{OperationID: "listUsers", Method: "GET", Path: "/users", Classification: "list"}},
	}
	out := Compile(input)
	if len(out.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", out.Errors)
	}
	if out.Plan == nil {
		t.Fatal("expected plan")
	}
	if len(out.Plan.FlowPlans) != 1 {
		t.Fatalf("expected 1 flow plan, got %d", len(out.Plan.FlowPlans))
	}
}

func TestCompileFlowMissingEnv(t *testing.T) {
	input := Input{
		Config: model.ProjectConfig{
			Flows: []model.Flow{{ID: "f1", Name: "test", Environment: "missing", Steps: []model.FlowStep{{Name: "s", OperationID: "op"}}}},
		},
		Operations: []spec.OperationInfo{{OperationID: "op", Method: "GET", Path: "/x"}},
	}
	out := Compile(input)
	if len(out.Errors) == 0 {
		t.Fatal("expected errors for missing environment")
	}
}

func TestCompileFlowMissingOp(t *testing.T) {
	input := Input{
		Config: model.ProjectConfig{
			Environments: []model.Environment{{ID: "dev", Name: "dev", BaseURL: "http://localhost"}},
			Flows: []model.Flow{{ID: "f1", Name: "test", Environment: "dev", Steps: []model.FlowStep{{Name: "s", OperationID: "nonexistent"}}}},
		},
		Operations: []spec.OperationInfo{},
	}
	out := Compile(input)
	if len(out.Errors) == 0 {
		t.Fatal("expected errors for missing operation")
	}
}

func TestCompileSuiteDestructiveSkipped(t *testing.T) {
	input := Input{
		Config: model.ProjectConfig{
			Environments: []model.Environment{{ID: "dev", Name: "dev", BaseURL: "http://localhost"}},
			Suites: []model.Suite{{
				ID: "s1", Name: "test-suite", Environment: "dev",
				Selector: model.SuiteSelector{Classifications: []string{"delete"}},
				Strategy: model.SuiteStrategy{DefaultList: true},
			}},
		},
		Operations: []spec.OperationInfo{{OperationID: "deleteUser", Method: "DELETE", Path: "/users/{id}", Classification: "delete", IsDestructive: true}},
	}
	out := Compile(input)
	if out.Plan == nil {
		t.Fatal("expected plan")
	}
	// Destructive ops should be skipped, so no cases
	if len(out.Plan.SuitePlans) > 0 && len(out.Plan.SuitePlans[0].Cases) > 0 {
		t.Error("expected destructive operations to be skipped")
	}
	// Should have a warning about skipping
	if len(out.Warnings) == 0 {
		t.Error("expected warning about skipping destructive op")
	}
}
