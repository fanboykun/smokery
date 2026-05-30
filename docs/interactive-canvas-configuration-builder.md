# Interactive Canvas Configuration Builder

## Goal

Replace the form-first flow and suite builders with one interactive canvas where users compose smoke configurations from OpenAPI operations, schema ports, and generated suite nodes.

The canvas is an authoring surface. Execution still follows the required Smokery pipeline:

```text
Canvas graph
  -> canonical ProjectConfig flows/suites
  -> compiler validation
  -> SmokePlan
  -> runner
```

## Design

Use `@xyflow/svelte` for the builder canvas. The canvas contains:

- Operation nodes for explicit flow steps.
- Suite generator nodes for generated list/read coverage.
- Sequence edges for execution order.
- Data-link edges for response-to-request mappings.
- Suite-selection edges for operation groups feeding generated suite nodes.

Operation nodes show the request and response shapes that Smokery already knows from OpenAPI. Request fields become input handles. Response fields become output handles. Connecting a response field to a request parameter generates a capture and interpolation binding in the derived flow config.

Example:

```text
listUsers.response.data[].id -> deleteUser.path.user_id
```

Derives:

```json
{
  "captures": [{ "name": "n_list_data_id", "source": "body", "path": "data.0.id" }],
  "params": { "user_id": "{{n_list_data_id}}" }
}
```

Suite generator nodes keep generated suite behavior compact. They do not expand every generated case into canvas nodes by default. Instead, the node shows matched operations, selected strategies, and generated case counts from the compiler preview.

## Implementation Steps

1. Add backend tests and types for canvas metadata.
2. Extend OpenAPI parsing to extract request schema and operation parameters in addition to response schema.
3. Add a canvas operation metadata endpoint:

```text
GET /api/specs/{spec-id}/operations/canvas
```

4. Add optional `canvas` metadata to `ProjectConfig`.
5. Add frontend canvas graph types.
6. Add a graph-to-config adapter that derives canonical `flows[]` and `suites[]`.
7. Replace `/projects/[id]/builder` with the primary Xyflow canvas page.
8. Keep existing flow/suite form routes as fallback routes.
9. Keep plan preview mandatory and run disabled while compiler diagnostics contain errors.

## Data-Link Strategy

Data links support a selector strategy:

- Default selector: `first`.
- Optional selector modes: `random`, `filter`.
- Filter expressions support simple predicates such as:
  - `modified_time exists`
  - `status == "active"`
  - `id not null`

Retry strategy is stored on the edge. In v1, retry execution applies only to the target node and defaults to disabled. When enabled, the target node may retry with up to two alternate selected items.

## Testing

Backend:

- Parser extracts request body schema, path params, query params, and response schema.
- Canvas operation endpoint returns diagram-ready metadata.
- `ProjectConfig.canvas` round-trips through JSON and does not affect existing compiler behavior.

Frontend:

- Graph adapter converts operation sequence nodes into `flows[]`.
- Graph adapter converts data-link edges into captures and target params.
- Graph adapter converts suite generator nodes into `suites[]`.
- Destructive operation nodes require explicit acknowledgement before config generation.

