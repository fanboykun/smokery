package assertion

import (
	"encoding/json"
	"fmt"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/tidwall/gjson"

	"github.com/fanboykun/smokery/apps/core/internal/model"
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
	case "schema":
		return checkSchema(a, body)
	case "pagination":
		return checkPagination(body)
	case "empty_result":
		return checkEmptyResult(a, body)
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

func checkSchema(a model.Assertion, body string) model.AssertionResult {
	// Schema is passed as Expected (json.RawMessage or string)
	var schemaBytes []byte
	switch v := a.Expected.(type) {
	case nil:
		return model.AssertionResult{Type: "schema", Passed: true, Message: "no schema to validate"}
	case string:
		schemaBytes = []byte(v)
	case json.RawMessage:
		schemaBytes = v
	case []byte:
		schemaBytes = v
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return model.AssertionResult{Type: "schema", Passed: false, Message: "invalid schema: " + err.Error()}
		}
		schemaBytes = b
	}
	if len(schemaBytes) == 0 || string(schemaBytes) == "null" {
		return model.AssertionResult{Type: "schema", Passed: true, Message: "no schema to validate"}
	}

	var schemaDoc any
	if err := json.Unmarshal(schemaBytes, &schemaDoc); err != nil {
		return model.AssertionResult{Type: "schema", Passed: false, Message: "invalid schema JSON: " + err.Error()}
	}

	c := jsonschema.NewCompiler()
	if err := c.AddResource("schema.json", schemaDoc); err != nil {
		return model.AssertionResult{Type: "schema", Passed: false, Message: "schema compile error: " + err.Error()}
	}
	sch, err := c.Compile("schema.json")
	if err != nil {
		return model.AssertionResult{Type: "schema", Passed: false, Message: "schema compile error: " + err.Error()}
	}

	var instance any
	if err := json.Unmarshal([]byte(body), &instance); err != nil {
		return model.AssertionResult{Type: "schema", Passed: false, Message: "invalid response JSON: " + err.Error()}
	}

	if err := sch.Validate(instance); err != nil {
		return model.AssertionResult{Type: "schema", Passed: false, Message: fmt.Sprintf("schema validation failed: %s", err.Error())}
	}
	return model.AssertionResult{Type: "schema", Passed: true}
}

func checkPagination(body string) model.AssertionResult {
	// Pagination sanity: response should be a list shape AND contain pagination metadata
	r := gjson.Parse(body)
	isList := r.IsArray()
	if !isList {
		for _, field := range []string{"data", "items", "results", "records"} {
			if gjson.Get(body, field).IsArray() {
				isList = true
				break
			}
		}
	}
	if !isList {
		return model.AssertionResult{Type: "pagination", Passed: false, Message: "response is not a list shape"}
	}
	// Check for pagination metadata indicators
	paginationFields := []string{"total", "total_count", "totalCount", "page", "pages", "next", "next_page", "has_more", "hasMore", "cursor", "meta.total", "pagination.total"}
	for _, f := range paginationFields {
		if gjson.Get(body, f).Exists() {
			return model.AssertionResult{Type: "pagination", Passed: true}
		}
	}
	// If it's a top-level array, pagination metadata might not be present — that's still valid
	if r.IsArray() {
		return model.AssertionResult{Type: "pagination", Passed: true, Message: "list returned but no pagination metadata found"}
	}
	return model.AssertionResult{Type: "pagination", Passed: true, Message: "list shape confirmed, no pagination metadata detected"}
}

func checkEmptyResult(a model.Assertion, body string) model.AssertionResult {
	// Policy: allow, warn, fail
	policy := "allow"
	if s, ok := a.Expected.(string); ok {
		policy = s
	}

	// Determine if result is empty
	r := gjson.Parse(body)
	isEmpty := false
	if r.IsArray() {
		isEmpty = len(r.Array()) == 0
	} else {
		for _, field := range []string{"data", "items", "results", "records"} {
			arr := gjson.Get(body, field)
			if arr.Exists() && arr.IsArray() {
				isEmpty = len(arr.Array()) == 0
				break
			}
		}
	}

	if !isEmpty {
		return model.AssertionResult{Type: "empty_result", Passed: true}
	}

	switch policy {
	case "fail":
		return model.AssertionResult{Type: "empty_result", Passed: false, Message: "empty result (policy: fail)"}
	case "warn":
		return model.AssertionResult{Type: "empty_result", Passed: true, Message: "empty result (policy: warn)"}
	default: // allow
		return model.AssertionResult{Type: "empty_result", Passed: true, Message: "empty result (policy: allow)"}
	}
}
