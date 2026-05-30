package spec

import (
	"encoding/json"
	"testing"
)

func TestParseExtractsOperationIOForCanvas(t *testing.T) {
	raw := []byte(`{
		"openapi": "3.1.0",
		"info": {"title": "Users API", "version": "1.0.0"},
		"paths": {
			"/users/{user_id}": {
				"delete": {
					"operationId": "deleteUser",
					"parameters": [
						{"name": "user_id", "in": "path", "required": true, "schema": {"type": "string"}},
						{"name": "dry_run", "in": "query", "schema": {"type": "boolean"}}
					],
					"responses": {"204": {"description": "deleted"}}
				},
				"patch": {
					"operationId": "updateUser",
					"requestBody": {
						"content": {
							"application/json": {
								"schema": {
									"type": "object",
									"properties": {"name": {"type": "string"}},
									"required": ["name"]
								}
							}
						}
					},
					"responses": {"200": {"description": "ok"}}
				}
			},
			"/users": {
				"get": {
					"operationId": "listUsers",
					"parameters": [
						{"name": "status", "in": "query", "schema": {"type": "string", "enum": ["active", "disabled"]}}
					],
					"responses": {
						"200": {
							"description": "ok",
							"content": {
								"application/json": {
									"schema": {
										"type": "object",
										"properties": {
											"data": {
												"type": "array",
												"items": {
													"type": "object",
													"properties": {
														"id": {"type": "string"},
														"modified_time": {"type": "string", "nullable": true}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}`)

	analysis, err := Parse(raw)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	list := findOperation(t, analysis.Operations, "listUsers")
	if len(list.Parameters) != 1 {
		t.Fatalf("listUsers parameters len = %d, want 1", len(list.Parameters))
	}
	if list.Parameters[0].Name != "status" || list.Parameters[0].In != "query" {
		t.Fatalf("listUsers parameter = %+v, want status query", list.Parameters[0])
	}
	if len(list.QueryHints.EnumFilters) != 1 || len(list.QueryHints.EnumFilters[0].Values) != 2 {
		t.Fatalf("enum hints = %+v, want status active/disabled", list.QueryHints.EnumFilters)
	}
	if !schemaHasPath(list.ResponseSchema, "data") {
		t.Fatalf("response schema should include data property: %s", string(list.ResponseSchema))
	}

	del := findOperation(t, analysis.Operations, "deleteUser")
	if len(del.Parameters) != 2 {
		t.Fatalf("deleteUser parameters len = %d, want 2", len(del.Parameters))
	}
	if del.Parameters[0].Name != "user_id" || del.Parameters[0].In != "path" || !del.Parameters[0].Required {
		t.Fatalf("deleteUser first parameter = %+v, want required path user_id", del.Parameters[0])
	}

	update := findOperation(t, analysis.Operations, "updateUser")
	if !schemaHasPath(update.RequestSchema, "name") {
		t.Fatalf("request schema should include name property: %s", string(update.RequestSchema))
	}
}

func findOperation(t *testing.T, ops []OperationInfo, id string) OperationInfo {
	t.Helper()
	for _, op := range ops {
		if op.OperationID == id {
			return op
		}
	}
	t.Fatalf("operation %q not found", id)
	return OperationInfo{}
}

func schemaHasPath(raw json.RawMessage, property string) bool {
	var schema map[string]any
	if err := json.Unmarshal(raw, &schema); err != nil {
		return false
	}
	props, ok := schema["properties"].(map[string]any)
	if !ok {
		return false
	}
	_, ok = props[property]
	return ok
}
