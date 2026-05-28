package compiler

import (
	"fmt"
	"time"

	"github.com/fanboykun/smokery/apps/core/internal/model"
	"github.com/fanboykun/smokery/apps/core/internal/spec"
	"github.com/google/uuid"
)

type Input struct {
	Config     model.ProjectConfig
	Operations []spec.OperationInfo
}

type Output struct {
	Plan     *model.SmokePlan     `json:"plan,omitempty"`
	Errors   []model.CompileError `json:"errors,omitempty"`
	Warnings []model.CompileError `json:"warnings,omitempty"`
}

func Compile(input Input) Output {
	var errors []model.CompileError
	var warnings []model.CompileError

	// Build operation lookup
	opMap := make(map[string]spec.OperationInfo)
	for _, op := range input.Operations {
		opMap[op.OperationID] = op
	}

	// Resolve flows
	var flowPlans []model.FlowPlan
	for _, flow := range input.Config.Flows {
		fp, errs, warns := compileFlow(flow, input.Config, opMap)
		errors = append(errors, errs...)
		warnings = append(warnings, warns...)
		if fp != nil {
			flowPlans = append(flowPlans, *fp)
		}
	}

	// Resolve suites
	var suitePlans []model.SuitePlan
	for _, suite := range input.Config.Suites {
		sp, errs, warns := compileSuite(suite, input.Config, input.Operations, opMap)
		errors = append(errors, errs...)
		warnings = append(warnings, warns...)
		if sp != nil {
			suitePlans = append(suitePlans, *sp)
		}
	}

	if len(errors) > 0 {
		return Output{Errors: errors, Warnings: warnings}
	}

	// Resolve environment (use first flow/suite env)
	envID := ""
	if len(input.Config.Flows) > 0 {
		envID = input.Config.Flows[0].Environment
	} else if len(input.Config.Suites) > 0 {
		envID = input.Config.Suites[0].Environment
	}
	env := resolveEnv(envID, input.Config.Environments)
	if env == nil && envID != "" {
		errors = append(errors, model.CompileError{Stage: "environment", Path: envID, Message: "environment not found", Severity: "error"})
		return Output{Errors: errors, Warnings: warnings}
	}

	plan := &model.SmokePlan{
		ID:         uuid.New().String(),
		FlowPlans:  flowPlans,
		SuitePlans: suitePlans,
		CompiledAt: time.Now(),
	}
	if env != nil {
		plan.Environment = *env
	}

	return Output{Plan: plan, Warnings: warnings}
}

func compileFlow(flow model.Flow, cfg model.ProjectConfig, opMap map[string]spec.OperationInfo) (*model.FlowPlan, []model.CompileError, []model.CompileError) {
	var errors, warnings []model.CompileError

	if resolveEnv(flow.Environment, cfg.Environments) == nil {
		errors = append(errors, model.CompileError{Stage: "flow", Path: flow.Name + ".environment", Message: fmt.Sprintf("environment %q not found", flow.Environment), Severity: "error", Entity: flow.ID})
	}

	var steps []model.PlannedStep
	for i, s := range flow.Steps {
		op, ok := opMap[s.OperationID]
		if !ok {
			errors = append(errors, model.CompileError{Stage: "flow", Path: fmt.Sprintf("%s.steps[%d]", flow.Name, i), Message: fmt.Sprintf("operation %q not found", s.OperationID), Severity: "error", Entity: flow.ID})
			continue
		}
		if op.IsDestructive {
			warnings = append(warnings, model.CompileError{Stage: "flow", Path: fmt.Sprintf("%s.steps[%d]", flow.Name, i), Message: fmt.Sprintf("operation %q is destructive", s.OperationID), Severity: "warning", Entity: flow.ID})
		}
		steps = append(steps, model.PlannedStep{
			Name: s.Name, Method: op.Method, Path: op.Path,
			Params: s.Params, Body: s.Body, Headers: s.Headers,
			Captures: s.Captures, Assertions: s.Assertions,
			ResponseSchema: op.ResponseSchema,
		})
	}

	var cleanup []model.PlannedStep
	for _, s := range flow.Cleanup {
		if op, ok := opMap[s.OperationID]; ok {
			cleanup = append(cleanup, model.PlannedStep{
				Name: s.Name, Method: op.Method, Path: op.Path,
				Params: s.Params, Body: s.Body, Headers: s.Headers,
			})
		}
	}

	if len(errors) > 0 {
		return nil, errors, warnings
	}
	return &model.FlowPlan{FlowID: flow.ID, Name: flow.Name, Steps: steps, Cleanup: cleanup}, errors, warnings
}

func compileSuite(suite model.Suite, cfg model.ProjectConfig, allOps []spec.OperationInfo, opMap map[string]spec.OperationInfo) (*model.SuitePlan, []model.CompileError, []model.CompileError) {
	var errors, warnings []model.CompileError

	if resolveEnv(suite.Environment, cfg.Environments) == nil {
		errors = append(errors, model.CompileError{Stage: "suite", Path: suite.Name + ".environment", Message: fmt.Sprintf("environment %q not found", suite.Environment), Severity: "error", Entity: suite.ID})
		return nil, errors, warnings
	}

	selected := selectOperations(allOps, suite.Selector)
	if len(selected) == 0 {
		warnings = append(warnings, model.CompileError{Stage: "suite", Path: suite.Name + ".selector", Message: "no operations matched selector", Severity: "warning", Entity: suite.ID})
		return &model.SuitePlan{SuiteID: suite.ID, Name: suite.Name}, errors, warnings
	}

	var cases []model.PlannedCase
	for _, op := range selected {
		if op.IsDestructive {
			warnings = append(warnings, model.CompileError{Stage: "suite", Path: suite.Name, Message: fmt.Sprintf("skipping destructive op %q", op.OperationID), Severity: "warning", Entity: suite.ID})
			continue
		}
		if suite.Strategy.DefaultList {
			cases = append(cases, model.PlannedCase{
				OperationID: op.OperationID, CaseType: "default_list",
				Step: model.PlannedStep{Name: op.OperationID + "_default", Method: op.Method, Path: op.Path,
					Assertions:     []model.Assertion{{Type: "status", Expected: 200}},
					ResponseSchema: op.ResponseSchema,
				},
			})
		}
		if suite.Strategy.Pagination && op.Classification == "list" {
			paginationParams := map[string]any{"page": 1, "limit": 10}
			// Use detected pagination keys if available
			if len(op.QueryHints.PaginationParams) > 0 {
				paginationParams = make(map[string]any)
				for _, p := range op.QueryHints.PaginationParams {
					switch p {
					case "page":
						paginationParams[p] = 1
					case "cursor":
						paginationParams[p] = ""
					default:
						paginationParams[p] = 10
					}
				}
			}
			cases = append(cases, model.PlannedCase{
				OperationID: op.OperationID, CaseType: "pagination",
				Step: model.PlannedStep{Name: op.OperationID + "_pagination", Method: op.Method, Path: op.Path,
					Params:         paginationParams,
					Assertions:     []model.Assertion{{Type: "status", Expected: 200}, {Type: "list_shape"}},
					ResponseSchema: op.ResponseSchema,
				},
			})
		}
		if suite.Strategy.SearchFromResponse && op.Classification == "list" && len(op.QueryHints.SearchParams) > 0 {
			cases = append(cases, generateSearchCase(op))
		}
		if suite.Strategy.EnumFilters && len(op.QueryHints.EnumFilters) > 0 {
			cases = append(cases, generateEnumCases(op, suite.Strategy.MaxCasesPerOp)...)
		}
	}

	// Apply empty result policy
	if suite.Strategy.EmptyResultPolicy != "" {
		for i := range cases {
			cases[i].EmptyResultPolicy = suite.Strategy.EmptyResultPolicy
		}
	}

	return &model.SuitePlan{SuiteID: suite.ID, Name: suite.Name, Cases: cases}, errors, warnings
}

func generateSearchCase(op spec.OperationInfo) model.PlannedCase {
	searchParam := op.QueryHints.SearchParams[0]
	// Step 1: list call to capture a value from the first item
	listStep := model.PlannedStep{
		Name:   op.OperationID + "_search_setup",
		Method: op.Method, Path: op.Path,
		Assertions: []model.Assertion{{Type: "status", Expected: 200}},
		Captures: []model.Capture{{
			Name: "search_term", Source: "body", Path: "0.name|data.0.name|items.0.name|results.0.name",
		}},
	}
	// Step 2: search call using captured value
	searchStep := model.PlannedStep{
		Name:   op.OperationID + "_search",
		Method: op.Method, Path: op.Path,
		Params:         map[string]any{searchParam: "{{search_term}}"},
		Assertions:     []model.Assertion{{Type: "status", Expected: 200}, {Type: "list_shape"}},
		ResponseSchema: op.ResponseSchema,
	}
	return model.PlannedCase{
		OperationID: op.OperationID, CaseType: "search",
		Steps: []model.PlannedStep{listStep, searchStep},
	}
}

func generateEnumCases(op spec.OperationInfo, maxPerOp int) []model.PlannedCase {
	var cases []model.PlannedCase
	for _, ef := range op.QueryHints.EnumFilters {
		for _, val := range ef.Values {
			cases = append(cases, model.PlannedCase{
				OperationID: op.OperationID, CaseType: "enum_filter",
				Step: model.PlannedStep{
					Name:           fmt.Sprintf("%s_filter_%s_%s", op.OperationID, ef.Name, val),
					Method:         op.Method,
					Path:           op.Path,
					Params:         map[string]any{ef.Name: val},
					Assertions:     []model.Assertion{{Type: "status", Expected: 200}},
					ResponseSchema: op.ResponseSchema,
				},
			})
			if maxPerOp > 0 && len(cases) >= maxPerOp {
				return cases
			}
		}
	}
	return cases
}

func selectOperations(ops []spec.OperationInfo, sel model.SuiteSelector) []spec.OperationInfo {
	var result []spec.OperationInfo
	for _, op := range ops {
		if matchSelector(op, sel) {
			result = append(result, op)
		}
	}
	return result
}

func matchSelector(op spec.OperationInfo, sel model.SuiteSelector) bool {
	for _, ex := range sel.Exclude {
		if op.OperationID == ex {
			return false
		}
	}
	if len(sel.Tags) == 0 && len(sel.Classifications) == 0 && len(sel.Paths) == 0 {
		return true
	}
	for _, t := range sel.Tags {
		for _, ot := range op.Tags {
			if t == ot {
				return true
			}
		}
	}
	for _, c := range sel.Classifications {
		if c == op.Classification {
			return true
		}
	}
	for _, p := range sel.Paths {
		if p == op.Path {
			return true
		}
	}
	return false
}

func resolveEnv(id string, envs []model.Environment) *model.Environment {
	for _, e := range envs {
		if e.ID == id {
			return &e
		}
	}
	return nil
}
