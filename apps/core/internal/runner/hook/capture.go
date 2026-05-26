package hook

import (
	"github.com/tidwall/gjson"

	"github.com/fanboykun/smokery/apps/core/internal/model"
)

// Capture extracts values from the response (body via gjson, or headers) into Vars.
type Capture struct{}

func (h *Capture) AfterResponse(rctx *ResponseContext, step *model.StepResult) error {
	if len(rctx.Step.Captures) == 0 {
		return nil
	}
	if step.Captures == nil {
		step.Captures = make(map[string]any)
	}
	body := string(rctx.Body)
	for _, cap := range rctx.Step.Captures {
		var val any
		switch cap.Source {
		case "body":
			r := gjson.Get(body, cap.Path)
			if r.Exists() {
				val = r.Value()
			}
		case "header":
			if v := rctx.Response.Header.Get(cap.Path); v != "" {
				val = v
			}
		}
		if val != nil {
			step.Captures[cap.Name] = val
			rctx.Vars[cap.Name] = val
		}
	}
	return nil
}
