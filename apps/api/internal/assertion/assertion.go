package assertion

import (
	"fmt"

	"github.com/tidwall/gjson"

	"github.com/fanboykun/smokery/apps/api/internal/model"
)

func Run(a model.Assertion, statusCode int, body string) model.AssertionResult {
	switch a.Type {
	case "status":
		return checkStatus(a, statusCode)
	case "jsonpath":
		return checkJSONPath(a, body)
	case "not_empty":
		return checkNotEmpty(body)
	case "list_shape":
		return checkListShape(body)
	default:
		return model.AssertionResult{Type: a.Type, Passed: false, Message: "unknown assertion type"}
	}
}

func checkStatus(a model.Assertion, code int) model.AssertionResult {
	expected := 200
	switch v := a.Expected.(type) {
	case float64:
		expected = int(v)
	case int:
		expected = v
	}
	passed := code == expected
	msg := ""
	if !passed {
		msg = fmt.Sprintf("expected status %d, got %d", expected, code)
	}
	return model.AssertionResult{Type: "status", Expected: expected, Actual: code, Passed: passed, Message: msg}
}

func checkJSONPath(a model.Assertion, body string) model.AssertionResult {
	r := gjson.Get(body, a.Path)
	if !r.Exists() {
		return model.AssertionResult{Type: "jsonpath", Expected: a.Expected, Passed: false, Message: fmt.Sprintf("path %q not found", a.Path)}
	}
	actual := r.Value()
	passed := fmt.Sprintf("%v", actual) == fmt.Sprintf("%v", a.Expected)
	msg := ""
	if !passed {
		msg = fmt.Sprintf("at %q: expected %v, got %v", a.Path, a.Expected, actual)
	}
	return model.AssertionResult{Type: "jsonpath", Expected: a.Expected, Actual: actual, Passed: passed, Message: msg}
}

func checkNotEmpty(body string) model.AssertionResult {
	passed := len(body) > 2
	msg := ""
	if !passed {
		msg = "response body is empty"
	}
	return model.AssertionResult{Type: "not_empty", Passed: passed, Message: msg}
}

func checkListShape(body string) model.AssertionResult {
	r := gjson.Parse(body)
	if r.IsArray() {
		return model.AssertionResult{Type: "list_shape", Passed: true}
	}
	for _, field := range []string{"data", "items", "results", "records"} {
		if gjson.Get(body, field).IsArray() {
			return model.AssertionResult{Type: "list_shape", Passed: true}
		}
	}
	return model.AssertionResult{Type: "list_shape", Passed: false, Message: "response is not a list shape"}
}
