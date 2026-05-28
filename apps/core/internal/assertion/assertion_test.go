package assertion

import (
	"encoding/json"
	"testing"

	"github.com/fanboykun/smokery/apps/core/internal/model"
)

func TestCheckStatus(t *testing.T) {
	r := Run(model.Assertion{Type: "status", Expected: 200}, 200, "")
	if !r.Passed {
		t.Error("expected pass for status 200")
	}
	r = Run(model.Assertion{Type: "status", Expected: 200}, 404, "")
	if r.Passed {
		t.Error("expected fail for status 404")
	}
}

func TestCheckJSONPath(t *testing.T) {
	body := `{"name": "test", "count": 5}`
	r := Run(model.Assertion{Type: "jsonpath", Path: "name", Expected: "test"}, 200, body)
	if !r.Passed {
		t.Errorf("expected pass, got: %s", r.Message)
	}
	r = Run(model.Assertion{Type: "jsonpath", Path: "name", Expected: "wrong"}, 200, body)
	if r.Passed {
		t.Error("expected fail for wrong value")
	}
	r = Run(model.Assertion{Type: "jsonpath", Path: "missing", Expected: "x"}, 200, body)
	if r.Passed {
		t.Error("expected fail for missing path")
	}
}

func TestCheckNotEmpty(t *testing.T) {
	r := Run(model.Assertion{Type: "not_empty"}, 200, `{"data": [1,2,3]}`)
	if !r.Passed {
		t.Error("expected pass for non-empty body")
	}
	r = Run(model.Assertion{Type: "not_empty"}, 200, `{}`)
	if r.Passed {
		t.Error("expected fail for empty object")
	}
}

func TestCheckListShape(t *testing.T) {
	r := Run(model.Assertion{Type: "list_shape"}, 200, `[1,2,3]`)
	if !r.Passed {
		t.Error("expected pass for array")
	}
	r = Run(model.Assertion{Type: "list_shape"}, 200, `{"data": [1,2]}`)
	if !r.Passed {
		t.Error("expected pass for data array wrapper")
	}
	r = Run(model.Assertion{Type: "list_shape"}, 200, `{"name": "test"}`)
	if r.Passed {
		t.Error("expected fail for non-list")
	}
}

func TestCheckSchema(t *testing.T) {
	schema := json.RawMessage(`{"type":"object","properties":{"name":{"type":"string"},"age":{"type":"integer"}},"required":["name"]}`)
	r := Run(model.Assertion{Type: "schema", Expected: schema}, 200, `{"name":"Alice","age":30}`)
	if !r.Passed {
		t.Errorf("expected pass, got: %s", r.Message)
	}
	r = Run(model.Assertion{Type: "schema", Expected: schema}, 200, `{"age":"not a number"}`)
	if r.Passed {
		t.Error("expected fail for missing required field")
	}
	// No schema = pass
	r = Run(model.Assertion{Type: "schema", Expected: nil}, 200, `{"anything":true}`)
	if !r.Passed {
		t.Errorf("expected pass with nil schema, got: %s", r.Message)
	}
}

func TestCheckPagination(t *testing.T) {
	r := Run(model.Assertion{Type: "pagination"}, 200, `{"data":[1,2],"total":10}`)
	if !r.Passed {
		t.Errorf("expected pass, got: %s", r.Message)
	}
	r = Run(model.Assertion{Type: "pagination"}, 200, `[1,2,3]`)
	if !r.Passed {
		t.Errorf("expected pass for array, got: %s", r.Message)
	}
	r = Run(model.Assertion{Type: "pagination"}, 200, `{"name":"not a list"}`)
	if r.Passed {
		t.Error("expected fail for non-list")
	}
}

func TestCheckEmptyResult(t *testing.T) {
	// Non-empty always passes
	r := Run(model.Assertion{Type: "empty_result", Expected: "fail"}, 200, `[1,2]`)
	if !r.Passed {
		t.Error("expected pass for non-empty")
	}
	// Empty with fail policy
	r = Run(model.Assertion{Type: "empty_result", Expected: "fail"}, 200, `[]`)
	if r.Passed {
		t.Error("expected fail for empty with fail policy")
	}
	// Empty with warn policy
	r = Run(model.Assertion{Type: "empty_result", Expected: "warn"}, 200, `{"data":[]}`)
	if !r.Passed {
		t.Error("expected pass (with warning) for empty with warn policy")
	}
	// Empty with allow policy
	r = Run(model.Assertion{Type: "empty_result", Expected: "allow"}, 200, `[]`)
	if !r.Passed {
		t.Error("expected pass for empty with allow policy")
	}
}
