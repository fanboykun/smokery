package spec

import (
	"fmt"
	"strings"

	"github.com/pb33f/libopenapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

type OperationInfo struct {
	OperationID    string   `json:"operation_id"`
	Method         string   `json:"method"`
	Path           string   `json:"path"`
	Summary        string   `json:"summary"`
	Tags           []string `json:"tags"`
	Classification string   `json:"classification"`
	IsDestructive  bool     `json:"is_destructive"`
}

type Analysis struct {
	Title      string          `json:"title"`
	Version    string          `json:"version"`
	Operations []OperationInfo `json:"operations"`
}

func Parse(raw []byte) (*Analysis, error) {
	doc, err := libopenapi.NewDocument(raw)
	if err != nil {
		return nil, fmt.Errorf("parse openapi: %w", err)
	}
	model, errs := doc.BuildV3Model()
	if errs != nil {
		return nil, fmt.Errorf("build model: %w", errs)
	}

	analysis := &Analysis{
		Title:   model.Model.Info.Title,
		Version: model.Model.Info.Version,
	}

	if model.Model.Paths != nil {
		for pair := model.Model.Paths.PathItems.First(); pair != nil; pair = pair.Next() {
			path := pair.Key()
			item := pair.Value()
			extractOps(path, item, &analysis.Operations)
		}
	}
	return analysis, nil
}

func extractOps(path string, item *v3.PathItem, ops *[]OperationInfo) {
	methods := map[string]*v3.Operation{
		"GET": item.Get, "POST": item.Post, "PUT": item.Put,
		"PATCH": item.Patch, "DELETE": item.Delete,
	}
	for method, op := range methods {
		if op == nil {
			continue
		}
		opID := op.OperationId
		if opID == "" {
			opID = strings.ToLower(method) + "_" + strings.ReplaceAll(strings.Trim(path, "/"), "/", "_")
		}
		var tags []string
		for _, t := range op.Tags {
			tags = append(tags, t)
		}
		classification := Classify(method, path, opID)
		*ops = append(*ops, OperationInfo{
			OperationID:    opID,
			Method:         method,
			Path:           path,
			Summary:        op.Summary,
			Tags:           tags,
			Classification: classification,
			IsDestructive:  isDestructive(method, opID),
		})
	}
}
