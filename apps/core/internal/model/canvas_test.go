package model

import (
	"encoding/json"
	"testing"
)

func TestProjectConfigCanvasRoundTripsAndDoesNotAffectCoreConfig(t *testing.T) {
	cfg := ProjectConfig{
		Environments: []Environment{{ID: "env-1", Name: "staging", BaseURL: "https://api.example.com"}},
		Flows:        []Flow{{ID: "flow-1", Name: "canvas flow", Environment: "env-1"}},
		Canvas: &CanvasGraph{
			Version: 1,
			Nodes: []CanvasNode{{
				ID:       "node-1",
				Type:     "operationNode",
				Position: CanvasPosition{X: 120, Y: 80},
				Data:     map[string]any{"operation_id": "listUsers"},
			}},
			Edges: []CanvasEdge{{
				ID:           "edge-1",
				Type:         "dataLink",
				Source:       "node-1",
				SourceHandle: "response:data[].id",
				Target:       "node-2",
				TargetHandle: "path:user_id",
				Data: map[string]any{
					"strategy": map[string]any{"selector": map[string]any{"mode": "first"}},
				},
			}},
			Viewport: CanvasViewport{X: 10, Y: 20, Zoom: 0.8},
		},
	}

	raw, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}
	var decoded ProjectConfig
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}

	if decoded.Canvas == nil {
		t.Fatal("Canvas = nil, want graph")
	}
	if decoded.Canvas.Nodes[0].Position.X != 120 {
		t.Fatalf("node x = %v, want 120", decoded.Canvas.Nodes[0].Position.X)
	}
	if len(decoded.Flows) != 1 || decoded.Flows[0].ID != "flow-1" {
		t.Fatalf("flows changed after canvas round trip: %+v", decoded.Flows)
	}
}
