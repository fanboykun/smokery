# Frontend Implementation Gaps - Current vs. Target State

## Current Reality Check

After auditing the code, these pages **exist** but are **incomplete**:
- `/projects/[id]/operations` - UI exists, missing spec hints in operation details
- `/projects/[id]/environments` - UI exists, basic env variable form
- `/projects/[id]/flows/[fid]` - Step editor exists, missing operation picker with dropdown + spec hints
- `/projects/[id]/suites/[sid]` - Suite config exists, missing live preview of matched operations
- `/projects/[id]/plan` - Plan exists, missing structured error display + test case preview
- `/runs/[runId]` - Run detail exists, missing tabs (summary, steps timeline, diagram, captures)

Pages **completely missing**:
- `/projects/[id]/spec` - OpenAPI spec explorer (tree, schemas, request/response examples)
- Environments / Auth auth profile creation form (only env vars, not auth schemes)

## Root Problem: Missing Spec-Driven UX

The **core product differentiator** is that Smokery is *spec-driven*. But the UI doesn't show the spec!

Users should:
1. **See operations** from their OpenAPI spec in a visual tree
2. **Get autocomplete** when building flows (operation dropdown, param hints)
3. **See live preview** of what their config produces (matched ops, generated cases)
4. **Understand errors** with links to specific config issues
5. **Debug runs** step-by-step with request/response details

Currently, users have to:
- Manually type operation IDs without knowing what exists
- Guess what params are needed (no schema hints)
- Compile, get errors, guess where the problem is
- Can't see why a case was skipped

## Tier 1 Gaps (MVP Blockers)

These prevent the product from working as designed:

### 1. Flow Editor: Operation Picker (HIGH IMPACT)
**Current**: Text input field, must guess operation IDs
**Target**: Dropdown list from spec, grouped by tag, shows method + path
**Also needed**: 
- Body editor with required fields highlighted
- Query/path params auto-populated from spec
- Assertion path autocomplete from response schema
- Capture source/path autocomplete from response schema

**Files to update**: `/projects/[id]/flows/[fid]/+page.svelte`

### 2. Suite Editor: Live Operation Preview (HIGH IMPACT)
**Current**: Selector checkboxes, but no feedback on what ops match
**Target**: Right panel shows live preview: "14 cases from 6 operations"
**Also needed**: Show which ops are excluded + why (destructive flag)

**Files to update**: `/projects/[id]/suites/[sid]/+page.svelte`

### 3. Plan Display: Structured Errors + Test Case Preview (MEDIUM IMPACT)
**Current**: Generic error, can't see what tests would run
**Target**: 
- Error panel showing specific config issues with links to fix (e.g., "Flow 'user' references undefined operation 'listUsers' - [see operations] [fix]")
- Test case preview: "Flows (3): user CRUD, admin only, ... | Suites (2): list ops, ..."

**Files to update**: `/projects/[id]/plan/+page.svelte`

### 4. Run Detail: Tabs + Timeline (HIGH IMPACT)
**Current**: Single view, only shows overall status
**Target**:
- [Summary] - overall stats
- [Steps] - timeline of each step (req, resp, assertions, captures)
- [Diagram] - Mermaid sequence diagram
- [Captures] - all captured variables

**Files to update**: `/runs/[runId]/+page.svelte`

### 5. Spec Explorer Page (MEDIUM IMPACT)
**Current**: Doesn't exist
**Target**: `/projects/[id]/spec` shows OpenAPI tree
- Operations organized by tag
- Click op to see request/response schema
- Search/filter operations

**Files to create**: `/projects/[id]/spec/+page.svelte`, `/projects/[id]/spec/[opId]/+page.svelte` (optional)

### 6. Auth Profile Creation (MEDIUM IMPACT)
**Current**: No UI, stored as JSON
**Target**: Form for different auth types (bearer, basic, api_key, oauth2, etc.)
- Type dropdown → render appropriate fields

**Files to update**: `/projects/[id]/environments/+page.svelte`

## Tier 2 Gaps (Nice to have)

### 7. Operations: Spec Hints in Detail View
Show when viewing operation in operations page:
- Query hints (from spec): pagination, search, filters
- Response schema preview
- Enum values for filter operations

**Files to update**: `/projects/[id]/operations/+page.svelte`

### 8. Run Comparison
Compare two runs to detect regressions.

**Files to create**: `/runs/[runId]/compare/[otherId]/+page.svelte`

### 9. Project Settings
General settings, integrations, danger zone.

**Files to create**: `/projects/[id]/settings/+page.svelte`

## Implementation Priority

### Week 1-2 (MVP Blockers):
1. Flow editor: operation picker + spec hints (2-3 days)
2. Suite editor: live preview (1-2 days)
3. Plan display: error clarity + case preview (1-2 days)
4. Run detail: tabs + timeline (2-3 days)

### Week 2-3 (MVP Support):
5. Spec explorer page (1-2 days)
6. Auth profile creation form (1 day)
7. Operations: spec hints (1 day)

### Week 3+ (Nice to have):
8. Run comparison (1 day)
9. Project settings (1 day)

## Backend Dependencies

These FE improvements need BE support:

### For Flow Editor Operation Picker:
- Need API endpoint: `GET /api/specs/{spec-id}/operations` - ✅ Already exists
- Need: Query hints, request schema, response schema in operation details

### For Suite Editor Live Preview:
- Need API endpoint to preview which ops match selector in real-time
- Currently must compile to know result

### For Plan Display Errors:
- Need: Structured error response from `POST /api/projects/{id}/plan/preview`
- Currently: Generic error, can't point user to fix

### For Run Detail Tabs:
- Need: Step-by-step execution details in `GET /api/runs/{id}`
- Currently: Only summary data

### For Spec Explorer:
- Could reuse existing `/api/specs/{spec-id}/operations` endpoint
- Might want: Schema definitions endpoint for rendering request/response schemas

## Conclusion

The product is 70% complete structurally, but 30% complete functionally. Users can't see/understand what they're configuring because **the spec visibility is missing** and **feedback loops are broken**.

Priority: Fix the feedback loops and spec visibility first (Tier 1), then add polish (Tier 2).
