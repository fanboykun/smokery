package assertion

import (
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
