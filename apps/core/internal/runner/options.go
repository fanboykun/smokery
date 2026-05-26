package runner

import (
	"net/http"
	"time"

	"github.com/fanboykun/smokery/apps/core/internal/runner/hook"
)

// Event is emitted via Options.EventEmitter when set.
type Event struct {
	Type string `json:"type"`
	Data any    `json:"data,omitempty"`
}

// Options configures a Runner.
type Options struct {
	HTTPClient     *http.Client
	PreProcessors  []hook.PreProcessor
	PostProcessors []hook.PostProcessor
	EventEmitter   func(Event)
}

// DefaultOptions returns Options with the standard built-in hooks.
func DefaultOptions() Options {
	return Options{
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
		PreProcessors: []hook.PreProcessor{
			&hook.VariableInterpolator{},
			&hook.AuthInjector{},
			&hook.EnvironmentHeaders{},
		},
		PostProcessors: []hook.PostProcessor{
			&hook.Capture{},
			&hook.TraceExtractor{},
			&hook.Redactor{},
			&hook.AssertionRunner{},
		},
	}
}
