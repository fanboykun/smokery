---
type: technical-note
status: current
tags:
  - config-builder
  - canvas
  - openapi
  - schemas
related:
  - "[[Project State]]"
  - "[[Config Persistence Blockers]]"
  - "[[ADR - Config Builder Redesign]]"
---

# Canvas Schema and Property Linking Redesign

> [!summary]
> This technical note details the implementation design to support connecting deeply nested request and response properties in the Config Builder canvas. It addresses the current slicing blocker (where only the first 5-8 fields are shown as handles) and provides a visual paradigm for cross-plan variable passing.

---

## 1. The Schema Slicing Blocker

In [OperationCanvasNode.svelte](file:///home/fanboykun/dev/personal/smokery/apps/web/src/lib/components/canvas/OperationCanvasNode.svelte), request and response fields are sliced:
```typescript
const responseFields = $derived(schemaFields(op?.response_schema).slice(0, 8));
const requestFields = $derived([
	...(op?.parameters ?? []).map((p) => ({ kind: p.in, path: p.name, required: p.required })),
	...schemaFields(op?.request_schema).slice(0, 5).map((f) => ({ kind: "body", path: f, required: false })),
]);
```
This slice prevents nodes from rendering with hundreds of nested fields. However, it blocks users from connecting any property not in the top 5-8 fields, rendering deep schema data-linking impossible.

---

## 2. Redesign Solution: Drawer-Based Property Mapping & Pinning

To support linking any arbitrary nested property without bloating the canvas nodes:

### 2.1 The "Pin to Node" Paradigm
- In the right-hand **OpenAPI Spec Panel**, every request and response field in the schema tree will display a pin icon.
- Clicking the pin icon adds that property to the node's list of active handles.
- **Node State:** We will store a `pinned_handles` list (array of property paths) in the node's `data` object:
  ```typescript
  export interface OperationNodeData extends Record<string, unknown> {
  	operation_id: string;
  	pinned_handles?: string[]; // New: list of pinned property paths
  }
  ```
- **Node Component:** `OperationCanvasNode.svelte` will dynamically merge the default required parameters with the `pinned_handles` list to render only the handles the user is actively working with.

### 2.2 Spec-Drawer Direct Connecting (Alternative/Co-exist)
- Instead of dragging handles, users can click a **"Link Source"** button next to a response property in Node A's spec drawer, then click **"Link Target"** next to a request parameter in Node B's spec drawer.
- The UI will automatically register the edge between Node A's source handle and Node B's target handle and render the thin data-link line on the canvas.

---

## 3. High-Level Plan Canvas: Variable Passing

Since the Go runner executes plans sequentially and shares the captured `vars` map across all steps:
- **Plan Outputs:** Each Flow Plan Node in the high-level Plan Canvas displays its output variables (variables captured from its inner operations).
- **Plan Inputs:** Each Flow Plan Node displays its unresolved input variables (variables referenced as `{{var_name}}` without a corresponding capture inside the same flow).
- **Plan Connections:** Dragging a connection from Flow A's output variable handle to Flow B's input variable handle:
  1. Guarantees that Flow A runs before Flow B (updates the sequence sorting).
  2. Resolves Flow B's variable reference, validating the compile flow.

---

## 4. Visual Layout of Operation Node with Pinned Handles

```text
┌──────────────────────────────────────────┐
│ [POST] create_user                       │
├──────────────────────────────────────────┤
│ REQUEST                  RESPONSE        │
│ ○ body:email             ○ body:id       │
│ ○ body:username                          │
│ (Pinned)                 (Pinned)        │
│ ○ body:profile.avatar    ○ body:token    │
└──────────────────────────────────────────┘
```
The Spec Drawer allows pinning/unpinning:
- `body:profile.avatar` 📌
- `body:profile.bio` 👤 (Unpinned, click pin to add handle to node)
