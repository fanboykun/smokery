package compiler

import (
	"fmt"
	"time"

	"github.com/fanboykun/smokery/apps/api/internal/model"
	"github.com/fanboykun/smokery/apps/api/internal/spec"
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
					Assertions: []model.Assertion{{Type: "status", Expected: 200}},
				},
			})
		}
		if suite.Strategy.Pagination && op.Classification == "list" {
			cases = append(cases, model.PlannedCase{
				OperationID: op.OperationID, CaseType: "pagination",
				Step: model.PlannedStep{Name: op.OperationID + "_pagination", Method: op.Method, Path: op.Path,
					Params:     map[string]any{"page": 1, "limit": 10},
					Assertions: []model.Assertion{{Type: "status", Expected: 200}, {Type: "list_shape"}},
				},
			})
		}
	}

	return &model.SuitePlan{SuiteID: suite.ID, Name: suite.Name, Cases: cases}, errors, warnings
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
