---
type: decision
status: accepted
date: 2026-06-13
tags:
  - adr
  - config-builder
  - frontend
related:
  - "[[Project State]]"
  - "[[Agent Context]]"
  - "[[Engineering Rules]]"
  - "[[Config Redesign Task List]]"
---

# ADR - Config Builder Redesign

> [!summary]
> The Config Builder is redesigned into a hierarchical workflow (high-level Plan Canvas with drill-down to individual nested Flow/Suite builders), featuring strict handle connection validation and a dedicated OpenAPI specification details panel.

## Context & Questions Addressed

To resolve the design and user flow issues in the original Config Builder (which immediately loaded all operations on a single flat canvas), we conducted a design review to align on key architectural decisions:

### 1. Plan Separation
* **Question:** How should plans (Flows and Suites) be separated and managed?
* **Decision:** **Option A (Single-Plan Canvas)**. Each Flow or Suite is a separate Plan. The user first lands on a project-level workspace listing/connecting these plans, and drilling down opens a canvas specifically for that plan's inner details.

### 2. High-Level Visualization & Order
* **Question:** How should the high-level plan execution order be visualized and managed?
* **Decision:** **High-level Plan Canvas with drill-down**. A parent canvas connects Plan Nodes (Flows & Suites) using sequence edges to define their execution order. Double-clicking or clicking a button on a Plan Node opens/drills down into its nested builder.

### 3. OpenAPI Spec Detail Panel
* **Question:** How should the OpenAPI spec details for a selected operation be displayed?
* **Decision:** **Sidebar/Drawer Detail Panel**. Selecting an operation node in the canvas opens a detailed panel on the right side showing the full OpenAPI spec, schemas, parameter tables, and examples.

### 4. Canvas Connection Validation & Style
* **Question:** How should we validate and style connections in the canvas?
* **Decision:** **Strict Handle Validation & Visual Distinctions**. Enforce that sequence handles (top/bottom) can only connect to sequence handles, and property handles (left/right) can only connect response properties to request properties. Style sequence edges differently (thick, dashed) from data-link edges (thin, animated/labeled).

---

## Redesign Specification

### 1. Plan Orchestration Canvas (High-Level)
- Displays Plan Nodes representing Flows and Suites.
- Connecting Plan Nodes defines their execution/runner order, which is used to sort the flows/suites in the final compiled project configuration.
- Actions:
  - Add Flow Plan / Add Suite Plan.
  - Double-click or click "Edit Canvas" to navigate/drill down into the selected Flow's inner operation-level canvas or the Suite's configuration forms.

### 2. Inner Flow Canvas (Operation-Level)
- Displays Operation Nodes representing specific OpenAPI operations.
- Data-linking connects individual Response properties (on the right handle) to Request properties (on the left handle) of other operations.
- Sequence flow connects `flow-out` (bottom handle) to `flow-in` (top handle) to define execution order of HTTP steps.

### 3. OpenAPI Spec Panel (Drawer)
- Drawer/sidebar component displaying rich OpenAPI metadata for the selected operation node in the flow canvas:
  - Method, Path, Summary, Description.
  - Request/Response parameters and schema trees.

### 4. Connection Rules (Validation)
- Top/Bottom handles (`flow-out` / `flow-in`) only connect to each other.
- Right handles (`response:*`) only connect to Left handles (`path:*`, `query:*`, `header:*`, `body:*`).
- Sequence edges styled as thick dashed paths.
- Data-link edges styled as thin paths.

## Related Notes
- [[Project State]]
- [[Agent Context]]
- [[Engineering Rules]]
- [[Config Redesign Task List]]
