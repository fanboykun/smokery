# Phase 2 UI Roadmap - Advanced Features

**Timeline**: Weeks 5-16 (12 weeks of development)

## Phase 2 Features (5 Tiers)

### Tier 1: Reporting (Weeks 5-6)

**Pages to Create**:
- `/runs/[runId]/report/contract` - Contract compliance view
- `/runs/[runId]/report/analyst` - Detailed analysis report  
- `/runs/[runId]/report/qa` - QA-focused summary
- `/runs/[runId]/report/correlation` - Correlation with observability

**Components Needed**:
- `ContractReportView.svelte` - Show spec compliance violations
- `AnalystReportView.svelte` - Detailed failure analysis
- `QAReportView.svelte` - Pass/fail summary for stakeholders
- `CorrelationView.svelte` - Link to traces, logs, metrics

**BE Dependencies**:
- GET `/api/runs/{id}/result` with detailed assertions
- GET `/api/runs/{id}/report/contract` - compliance data
- GET `/api/runs/{id}/correlations` - trace IDs, span links

---

### Tier 2: Failure Classification + Assignee (Weeks 7-8)

**Pages to Create**:
- `/runs?status=failed&classification=network` - Filter by root cause
- `/runs/[runId]/classify` - Inline failure classification
- Modal for assigning failures to team members

**Components Needed**:
- `FailureClassifier.svelte` - Dropdown with root cause options
- `AssigneeSelector.svelte` - Team member picker
- `FailureTimeline.svelte` - Who did what and when

**BE Dependencies**:
- PUT `/api/runs/{id}/failure-classification` - Save classification
- GET `/api/team/members` - List assignees
- PUT `/api/runs/{id}/assigned-to` - Assign failure

---

### Tier 3: Spec Evolution (Weeks 9-10)

**Pages to Create**:
- `/projects/[id]/spec/versions` - Version history
- `/projects/[id]/spec/diff?from={v1}&to={v2}` - Diff viewer
- `/projects/[id]/impact?spec-version={v2}` - Impact analysis

**Components Needed**:
- `SpecVersionList.svelte` - Timeline of spec uploads
- `SpecDiffViewer.svelte` - Side-by-side diff of OpenAPI changes
- `ImpactAnalysis.svelte` - Show which tests are affected

**BE Dependencies**:
- GET `/api/projects/{id}/specs` - All versions
- GET `/api/specs/{id}/diff/{otherId}` - Spec differences
- GET `/api/projects/{id}/impact/spec/{specId}` - Affected tests

---

### Tier 4: Analytics (Weeks 11-13)

**Pages to Create**:
- `/projects/[id]/analytics/latency` - Response time trends
- `/projects/[id]/analytics/flaky` - Flaky test detection
- `/projects/[id]/analytics/trends` - Pass rate over time

**Components Needed**:
- `LatencyChart.svelte` - Line chart of p50/p95/p99
- `FlakyTestDetector.svelte` - Operations that fail intermittently
- `HealthTrendChart.svelte` - Pass rate timeline

**BE Dependencies**:
- GET `/api/projects/{id}/analytics/latency?range=7d` 
- GET `/api/projects/{id}/analytics/flaky-operations`
- GET `/api/projects/{id}/analytics/health-trends`

---

### Tier 5: Team Governance (Weeks 14-16)

**Pages to Create**:
- `/projects/[id]/settings/members` - Manage team access
- `/projects/[id]/settings/audit` - Audit log
- `/projects/[id]/settings/webhooks` - Webhook configuration
- `/projects/[id]/settings/notifications` - Alert rules

**Components Needed**:
- `MemberManager.svelte` - Add/remove/permission users
- `AuditLog.svelte` - Immutable action log
- `WebhookEditor.svelte` - Configure webhook endpoints
- `NotificationRules.svelte` - Email, Slack, PagerDuty

**BE Dependencies**:
- GET/POST/DELETE `/api/projects/{id}/members`
- GET `/api/projects/{id}/audit-log`
- GET/POST/DELETE `/api/projects/{id}/webhooks`
- GET/POST `/api/projects/{id}/notifications`

---

## UI Patterns to Reuse

From MVP, we established:
- **OperationPicker.svelte** - Use for operation selection in reports
- **EmptyState.svelte** - Every new page needs empty state
- **StatusBadge.svelte** - For classification and status display
- **SplitPane.svelte** - For diff viewers and comparisons
- **Tabs.svelte** - Report views are multi-tab (was proven in run detail)

---

## Design System Alignment

All Phase 2 pages should follow:
- Dark theme (zinc-950 bg, emerald-500 primary accent)
- Dense but scannable layouts
- Clear data visualization (Recharts for charts)
- Actionable error states with fix links
- Real-time updates where applicable

---

## Success Criteria for Phase 2

1. **Reporting**: Users can generate compliance and analysis reports
2. **Classification**: Teams can mark root causes and assign follow-ups
3. **Evolution**: Teams can safely upgrade specs with impact visibility
4. **Analytics**: Teams see trends and spot flaky operations
5. **Governance**: Multiple users can collaborate with clear audit trails

---

## Implementation Order

1. **Report views** (highest user value immediately)
2. **Failure classification** (enables triage workflows)
3. **Spec evolution** (unblocks API versioning)
4. **Analytics** (helps identify patterns)
5. **Governance** (enables team scaling)

This order maximizes immediate value while building toward enterprise features.
