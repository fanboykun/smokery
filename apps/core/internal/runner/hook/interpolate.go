package hook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// VariableInterpolator replaces {var} placeholders in path, headers, and body
// using captured variables and step params.
type VariableInterpolator struct{}

func (h *VariableInterpolator) BeforeRequest(rctx *RequestContext) error {
	// Build placeholder→value map from step params + runtime vars (vars take precedence).
	values := make(map[string]string, len(rctx.Step.Params)+len(rctx.Vars))
	for k, v := range rctx.Step.Params {
		values[k] = fmt.Sprintf("%v", v)
	}
	for k, v := range rctx.Vars {
		values[k] = fmt.Sprintf("%v", v)
	}

	// Interpolate URL path.
	if rctx.Request.URL != nil {
		path := rctx.Request.URL.Path
		for k, v := range values {
			path = strings.ReplaceAll(path, "{"+k+"}", v)
		}
		rctx.Request.URL.Path = path
	}

	// Interpolate headers.
	for k, vs := range rctx.Request.Header {
		for i, v := range vs {
			for vk, vv := range rctx.Vars {
				v = strings.ReplaceAll(v, "{{"+vk+"}}", fmt.Sprintf("%v", vv))
			}
			rctx.Request.Header[k][i] = v
		}
	}

	// Interpolate body if it's JSON.
	if rctx.Request.Body == nil {
		return nil
	}
	bodyBytes, err := io.ReadAll(rctx.Request.Body)
	if err != nil {
		return err
	}
	_ = rctx.Request.Body.Close()
	bodyStr := string(bodyBytes)
	for vk, vv := range rctx.Vars {
		needle := "{{" + vk + "}}"
		if strings.Contains(bodyStr, needle) {
			repl, _ := json.Marshal(vv)
			// strip outer quotes for raw string subs
			bodyStr = strings.ReplaceAll(bodyStr, `"`+needle+`"`, string(repl))
			bodyStr = strings.ReplaceAll(bodyStr, needle, fmt.Sprintf("%v", vv))
		}
	}
	rctx.Request.Body = io.NopCloser(bytes.NewReader([]byte(bodyStr)))
	rctx.Request.ContentLength = int64(len(bodyStr))
	return nil
}
