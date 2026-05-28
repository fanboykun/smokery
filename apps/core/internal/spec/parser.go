package spec

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pb33f/libopenapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

type QueryHints struct {
	PaginationParams []string    `json:"pagination_params,omitempty"` // e.g. page, limit, offset, cursor
	SearchParams     []string    `json:"search_params,omitempty"`     // e.g. q, query, search
	EnumFilters      []EnumParam `json:"enum_filters,omitempty"`      // params with enum values
}

type EnumParam struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

type OperationInfo struct {
	OperationID    string          `json:"operation_id"`
	Method         string          `json:"method"`
	Path           string          `json:"path"`
	Summary        string          `json:"summary"`
	Tags           []string        `json:"tags"`
	Classification string          `json:"classification"`
	IsDestructive  bool            `json:"is_destructive"`
	QueryHints     QueryHints      `json:"query_hints,omitempty"`
	ResponseSchema json.RawMessage `json:"response_schema,omitempty"` // JSON Schema for success response
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

var paginationKeys = map[string]bool{
	"page": true, "limit": true, "offset": true, "cursor": true,
	"per_page": true, "page_size": true, "pagesize": true, "size": true,
}

var searchKeys = map[string]bool{
	"q": true, "query": true, "search": true, "keyword": true, "filter": true,
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
		tags = append(tags, op.Tags...)
		classification := Classify(method, path, opID)
		hints := extractQueryHints(op)
		respSchema := extractResponseSchema(op)
		*ops = append(*ops, OperationInfo{
			OperationID:    opID,
			Method:         method,
			Path:           path,
			Summary:        op.Summary,
			Tags:           tags,
			Classification: classification,
			IsDestructive:  isDestructive(method, opID),
			QueryHints:     hints,
			ResponseSchema: respSchema,
		})
	}
}

func extractQueryHints(op *v3.Operation) QueryHints {
	var hints QueryHints
	if op.Parameters == nil {
		return hints
	}
	for _, param := range op.Parameters {
		if param == nil || param.In != "query" {
			continue
		}
		name := strings.ToLower(param.Name)
		if paginationKeys[name] {
			hints.PaginationParams = append(hints.PaginationParams, param.Name)
		} else if searchKeys[name] {
			hints.SearchParams = append(hints.SearchParams, param.Name)
		}
		// Check for enum values in schema
		if param.Schema != nil && param.Schema.Schema() != nil {
			schema := param.Schema.Schema()
			if len(schema.Enum) > 0 {
				var values []string
				for _, e := range schema.Enum {
					if e != nil {
						values = append(values, fmt.Sprintf("%v", e.Value))
					}
				}
				if len(values) > 0 {
					hints.EnumFilters = append(hints.EnumFilters, EnumParam{Name: param.Name, Values: values})
				}
			}
		}
	}
	return hints
}

func extractResponseSchema(op *v3.Operation) json.RawMessage {
	if op.Responses == nil || op.Responses.Codes == nil {
		return nil
	}
	// Try 200, then 201, then first 2xx
	for _, code := range []string{"200", "201"} {
		if pair := op.Responses.Codes.GetOrZero(code); pair != nil {
			if s := schemaFromResponse(pair); s != nil {
				return s
			}
		}
	}
	for pair := op.Responses.Codes.First(); pair != nil; pair = pair.Next() {
		if strings.HasPrefix(pair.Key(), "2") {
			if s := schemaFromResponse(pair.Value()); s != nil {
				return s
			}
		}
	}
	return nil
}

func schemaFromResponse(resp *v3.Response) json.RawMessage {
	if resp == nil || resp.Content == nil {
		return nil
	}
	mt := resp.Content.GetOrZero("application/json")
	if mt == nil {
		return nil
	}
	if mt.Schema == nil || mt.Schema.Schema() == nil {
		return nil
	}
	rendered, err := mt.Schema.Schema().Render()
	if err != nil {
		return nil
	}
	// Convert YAML to JSON
	var raw any
	if json.Unmarshal(rendered, &raw) == nil {
		return rendered
	}
	// rendered might be YAML, try converting
	jsonBytes, err := json.Marshal(raw)
	if err != nil {
		return nil
	}
	return jsonBytes
}
