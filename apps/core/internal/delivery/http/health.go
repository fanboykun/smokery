package http

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// HealthChecker probes infrastructure connectivity.
type HealthChecker interface {
	PingDB(ctx context.Context) error
	PingBlob(ctx context.Context) error
}

type HealthDetailOutput struct {
	Body struct {
		Status string `json:"status" example:"ok"`
		DB     string `json:"db" example:"ok"`
		Blob   string `json:"blob" example:"ok"`
	}
}

func RegisterHealthCheck(api huma.API, checker HealthChecker) {
	huma.Register(api, huma.Operation{
		OperationID: "health-check",
		Method:      http.MethodGet,
		Path:        "/health",
		Summary:     "Health check with DB and blob store probes",
	}, func(ctx context.Context, _ *struct{}) (*HealthDetailOutput, error) {
		out := &HealthDetailOutput{}
		out.Body.Status = "ok"

		if err := checker.PingDB(ctx); err != nil {
			out.Body.Status = "degraded"
			out.Body.DB = "error: " + err.Error()
		} else {
			out.Body.DB = "ok"
		}

		if err := checker.PingBlob(ctx); err != nil {
			out.Body.Status = "degraded"
			out.Body.Blob = "error: " + err.Error()
		} else {
			out.Body.Blob = "ok"
		}

		return out, nil
	})
}
