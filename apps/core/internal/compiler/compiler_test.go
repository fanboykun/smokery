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

func TestCompileSuiteSearchFromResponse(t *testing.T) {
	input := Input{
		Config: model.ProjectConfig{
			Environments: []model.Environment{{ID: "dev", Name: "dev", BaseURL: "http://localhost"}},
			Suites: []model.Suite{{
				ID: "s1", Name: "search-suite", Environment: "dev",
				Strategy: model.SuiteStrategy{SearchFromResponse: true},
			}},
		},
		Operations: []spec.OperationInfo{{
			OperationID: "listUsers", Method: "GET", Path: "/users", Classification: "list",
			QueryHints: spec.QueryHints{SearchParams: []string{"q"}},
		}},
	}
	out := Compile(input)
	if len(out.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", out.Errors)
	}
	if out.Plan == nil || len(out.Plan.SuitePlans) == 0 {
		t.Fatal("expected suite plan")
	}
	cases := out.Plan.SuitePlans[0].Cases
	found := false
	for _, c := range cases {
		if c.CaseType == "search" {
			found = true
			if len(c.Steps) != 2 {
				t.Fatalf("search case should have 2 steps, got %d", len(c.Steps))
			}
			if len(c.Steps[0].Captures) == 0 {
				t.Error("first step should have captures")
			}
			if c.Steps[1].Params["q"] != "{{search_term}}" {
				t.Errorf("second step should use captured search_term, got params: %v", c.Steps[1].Params)
			}
		}
	}
	if !found {
		t.Error("expected a search case to be generated")
	}
}

func TestCompileSuiteSearchFromResponseNoSearchParams(t *testing.T) {
	input := Input{
		Config: model.ProjectConfig{
			Environments: []model.Environment{{ID: "dev", Name: "dev", BaseURL: "http://localhost"}},
			Suites: []model.Suite{{
				ID: "s1", Name: "search-suite", Environment: "dev",
				Strategy: model.SuiteStrategy{SearchFromResponse: true},
			}},
		},
		Operations: []spec.OperationInfo{{
			OperationID: "listUsers", Method: "GET", Path: "/users", Classification: "list",
			QueryHints: spec.QueryHints{}, // no search params
		}},
	}
	out := Compile(input)
	if out.Plan == nil || len(out.Plan.SuitePlans) == 0 {
		t.Fatal("expected suite plan")
	}
	for _, c := range out.Plan.SuitePlans[0].Cases {
		if c.CaseType == "search" {
			t.Error("should not generate search case without search params")
		}
	}
}

func TestCompileSuiteEnumFilters(t *testing.T) {
	input := Input{
		Config: model.ProjectConfig{
			Environments: []model.Environment{{ID: "dev", Name: "dev", BaseURL: "http://localhost"}},
			Suites: []model.Suite{{
				ID: "s1", Name: "enum-suite", Environment: "dev",
				Strategy: model.SuiteStrategy{EnumFilters: true},
			}},
		},
		Operations: []spec.OperationInfo{{
			OperationID: "listOrders", Method: "GET", Path: "/orders", Classification: "list",
			QueryHints: spec.QueryHints{
				EnumFilters: []spec.EnumParam{{Name: "status", Values: []string{"pending", "shipped", "delivered"}}},
			},
		}},
	}
	out := Compile(input)
	if len(out.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", out.Errors)
	}
	if out.Plan == nil || len(out.Plan.SuitePlans) == 0 {
		t.Fatal("expected suite plan")
	}
	var enumCases int
	for _, c := range out.Plan.SuitePlans[0].Cases {
		if c.CaseType == "enum_filter" {
			enumCases++
		}
	}
	if enumCases != 3 {
		t.Errorf("expected 3 enum_filter cases, got %d", enumCases)
	}
}

func TestCompileSuiteEnumFiltersMaxCasesPerOp(t *testing.T) {
	input := Input{
		Config: model.ProjectConfig{
			Environments: []model.Environment{{ID: "dev", Name: "dev", BaseURL: "http://localhost"}},
			Suites: []model.Suite{{
				ID: "s1", Name: "enum-suite", Environment: "dev",
				Strategy: model.SuiteStrategy{EnumFilters: true, MaxCasesPerOp: 2},
			}},
		},
		Operations: []spec.OperationInfo{{
			OperationID: "listOrders", Method: "GET", Path: "/orders", Classification: "list",
			QueryHints: spec.QueryHints{
				EnumFilters: []spec.EnumParam{{Name: "status", Values: []string{"a", "b", "c", "d"}}},
			},
		}},
	}
	out := Compile(input)
	if out.Plan == nil || len(out.Plan.SuitePlans) == 0 {
		t.Fatal("expected suite plan")
	}
	var enumCases int
	for _, c := range out.Plan.SuitePlans[0].Cases {
		if c.CaseType == "enum_filter" {
			enumCases++
		}
	}
	if enumCases != 2 {
		t.Errorf("expected 2 enum_filter cases (max), got %d", enumCases)
	}
}

func TestCompileSuiteEmptyResultPolicy(t *testing.T) {
	input := Input{
		Config: model.ProjectConfig{
			Environments: []model.Environment{{ID: "dev", Name: "dev", BaseURL: "http://localhost"}},
			Suites: []model.Suite{{
				ID: "s1", Name: "policy-suite", Environment: "dev",
				Strategy: model.SuiteStrategy{DefaultList: true, EmptyResultPolicy: "fail"},
			}},
		},
		Operations: []spec.OperationInfo{{
			OperationID: "listUsers", Method: "GET", Path: "/users", Classification: "list",
		}},
	}
	out := Compile(input)
	if out.Plan == nil || len(out.Plan.SuitePlans) == 0 {
		t.Fatal("expected suite plan")
	}
	for _, c := range out.Plan.SuitePlans[0].Cases {
		if c.EmptyResultPolicy != "fail" {
			t.Errorf("expected EmptyResultPolicy=fail, got %q", c.EmptyResultPolicy)
		}
	}
}

func TestCompileSuitePaginationUsesQueryHints(t *testing.T) {
	input := Input{
		Config: model.ProjectConfig{
			Environments: []model.Environment{{ID: "dev", Name: "dev", BaseURL: "http://localhost"}},
			Suites: []model.Suite{{
				ID: "s1", Name: "pagination-suite", Environment: "dev",
				Strategy: model.SuiteStrategy{Pagination: true},
			}},
		},
		Operations: []spec.OperationInfo{{
			OperationID: "listUsers", Method: "GET", Path: "/users", Classification: "list",
			QueryHints: spec.QueryHints{PaginationParams: []string{"page", "per_page"}},
		}},
	}
	out := Compile(input)
	if out.Plan == nil || len(out.Plan.SuitePlans) == 0 {
		t.Fatal("expected suite plan")
	}
	var paginationCase *model.PlannedCase
	for i, c := range out.Plan.SuitePlans[0].Cases {
		if c.CaseType == "pagination" {
			paginationCase = &out.Plan.SuitePlans[0].Cases[i]
			break
		}
	}
	if paginationCase == nil {
		t.Fatal("expected pagination case")
	}
	if _, ok := paginationCase.Step.Params["page"]; !ok {
		t.Error("expected 'page' param from query hints")
	}
	if _, ok := paginationCase.Step.Params["per_page"]; !ok {
		t.Error("expected 'per_page' param from query hints")
	}
}
