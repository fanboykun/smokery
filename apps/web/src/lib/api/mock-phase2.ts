// Mock API Service for Phase 2 Development
// This file provides realistic test data for UI development before backend is ready
// Replace these functions with real API calls once backend implements the endpoints

import type {
  ContractReport,
  AnalystReport,
  QAReport,
  CorrelationReport,
  TeamMember,
  RunFailureClassification,
  RunAssignment,
  FailureAction,
  SpecVersion,
  SpecDiff,
  ImpactAnalysis,
  LatencyAnalytics,
  FlakyOperationsAnalytics,
  HealthTrendsAnalytics,
  ProjectMember,
  AuditLogEntry,
  Webhook,
  WebhookDelivery,
  NotificationRule,
} from './phase2';

// ============================================================================
// TIER 1: REPORTING MOCKS
// ============================================================================

export async function mockGetContractReport(runId: string): Promise<ContractReport> {
  await new Promise(resolve => setTimeout(resolve, 500));
  
  return {
    id: `report-${runId}-contract`,
    run_id: runId,
    total_violations: 3,
    errors: [
      {
        id: 'v1',
        operation_id: 'createUser',
        violation_type: 'response_schema_mismatch',
        severity: 'error',
        message: 'Response missing required field "id"',
        expected_schema: { type: 'object', required: ['id', 'email', 'name'] },
        actual_value: { email: 'user@example.com', name: 'John' },
        location: 'response.body',
      },
    ],
    warnings: [
      {
        id: 'w1',
        operation_id: 'getUser',
        violation_type: 'unexpected_field',
        severity: 'warning',
        message: 'Unexpected field "internal_id" in response',
        expected_schema: { type: 'object', properties: { id: { type: 'string' }, email: { type: 'string' } } },
        actual_value: { id: '123', email: 'user@example.com', internal_id: 'xyz' },
        location: 'response.body',
      },
    ],
    compliance_score: 85,
    passed_assertions: 47,
    failed_assertions: 3,
  };
}

export async function mockGetAnalystReport(runId: string): Promise<AnalystReport> {
  await new Promise(resolve => setTimeout(resolve, 500));
  
  return {
    id: `report-${runId}-analyst`,
    run_id: runId,
    summary: 'Run had 2 failures caused by network timeouts and 1 schema violation. Root cause analysis suggests rate limiting on user service.',
    root_causes: [
      {
        cause: 'Network Timeout (Gateway Timeout)',
        impact: 60,
        affected_operations: ['createUser', 'updateUser'],
      },
      {
        cause: 'Contract Violation',
        impact: 40,
        affected_operations: ['getUser'],
      },
    ],
    recommendations: [
      {
        title: 'Increase API Gateway Timeout',
        description: 'Current timeout is 30s but user service takes up to 45s under load',
        priority: 'high',
      },
      {
        title: 'Add Retry Logic',
        description: 'Implement exponential backoff for transient failures',
        priority: 'high',
      },
      {
        title: 'Fix User Service Response Schema',
        description: 'Ensure all responses include required "id" field',
        priority: 'medium',
      },
    ],
    timeline_insights: [
      { timestamp: '2025-05-30T10:00:00Z', operation_id: 'createUser', status: 'passed', duration_ms: 145 },
      { timestamp: '2025-05-30T10:00:01Z', operation_id: 'getUser', status: 'failed', duration_ms: 2000, error: 'Response timeout' },
      { timestamp: '2025-05-30T10:00:02Z', operation_id: 'updateUser', status: 'failed', duration_ms: 31000, error: 'Gateway timeout' },
    ],
  };
}

export async function mockGetQAReport(runId: string): Promise<QAReport> {
  await new Promise(resolve => setTimeout(resolve, 500));
  
  return {
    id: `report-${runId}-qa`,
    run_id: runId,
    status: 'failed',
    total_tests: 50,
    passed_tests: 48,
    failed_tests: 2,
    pass_rate: 96,
    flaky_tests: ['getUser', 'listUsers'],
    coverage_summary: {
      total_operations: 15,
      tested_operations: 14,
      coverage_percentage: 93,
    },
    blockers: [
      {
        operation_id: 'getUser',
        issue: 'Response timeout under concurrent load',
        severity: 'high',
      },
    ],
  };
}

export async function mockGetCorrelationReport(runId: string): Promise<CorrelationReport> {
  await new Promise(resolve => setTimeout(resolve, 500));
  
  return {
    id: `report-${runId}-correlation`,
    run_id: runId,
    root_trace_id: 'trace-abc123xyz789',
    links: [
      {
        type: 'trace',
        name: 'OpenTelemetry Trace',
        url: 'https://tempo.example.com/explore?traceID=trace-abc123xyz789',
        trace_id: 'trace-abc123xyz789',
      },
      {
        type: 'log',
        name: 'CloudWatch Logs',
        url: 'https://console.aws.amazon.com/cloudwatch/logs/?query=trace-abc123xyz789',
      },
      {
        type: 'metric',
        name: 'Prometheus Metrics',
        url: 'https://prometheus.example.com/?query=trace_id="trace-abc123xyz789"',
      },
    ],
    metrics: {
      p50_latency_ms: 150,
      p95_latency_ms: 800,
      p99_latency_ms: 2500,
      error_rate: 0.02,
    },
    external_services: [
      { name: 'User Service', status: 'degraded', latency_ms: 450 },
      { name: 'Auth Service', status: 'healthy', latency_ms: 50 },
      { name: 'Database', status: 'healthy', latency_ms: 25 },
    ],
  };
}

// ============================================================================
// TIER 2: FAILURE CLASSIFICATION MOCKS
// ============================================================================

export async function mockGetTeamMembers(projectId: string): Promise<TeamMember[]> {
  await new Promise(resolve => setTimeout(resolve, 300));
  
  return [
    {
      id: 'user-1',
      name: 'Alice Chen',
      email: 'alice@example.com',
      role: 'admin',
      avatar_url: 'https://api.example.com/avatars/alice.png',
      joined_at: '2025-01-01T00:00:00Z',
    },
    {
      id: 'user-2',
      name: 'Bob Smith',
      email: 'bob@example.com',
      role: 'editor',
      avatar_url: 'https://api.example.com/avatars/bob.png',
      joined_at: '2025-02-15T00:00:00Z',
    },
    {
      id: 'user-3',
      name: 'Charlie Davis',
      email: 'charlie@example.com',
      role: 'viewer',
      avatar_url: 'https://api.example.com/avatars/charlie.png',
      joined_at: '2025-03-10T00:00:00Z',
    },
  ];
}

export async function mockGetFailureActions(runId: string): Promise<FailureAction[]> {
  await new Promise(resolve => setTimeout(resolve, 300));
  
  return [
    {
      id: 'action-1',
      run_id: runId,
      action_type: 'classified',
      actor_id: 'user-1',
      actor_name: 'Alice Chen',
      created_at: '2025-05-30T10:05:00Z',
      details: { classification: 'network_timeout' },
    },
    {
      id: 'action-2',
      run_id: runId,
      action_type: 'assigned',
      actor_id: 'user-1',
      actor_name: 'Alice Chen',
      created_at: '2025-05-30T10:06:00Z',
      details: { assigned_to: 'user-2', assigned_to_name: 'Bob Smith' },
    },
    {
      id: 'action-3',
      run_id: runId,
      action_type: 'status_changed',
      actor_id: 'user-2',
      actor_name: 'Bob Smith',
      created_at: '2025-05-30T10:15:00Z',
      details: { from_status: 'open', to_status: 'in_progress' },
    },
  ];
}

// ============================================================================
// TIER 3: SPEC EVOLUTION MOCKS
// ============================================================================

export async function mockGetSpecVersions(projectId: string): Promise<SpecVersion[]> {
  await new Promise(resolve => setTimeout(resolve, 400));
  
  return [
    {
      id: 'spec-1',
      project_id: projectId,
      version: '1.0.0',
      uploaded_at: '2025-01-01T00:00:00Z',
      uploaded_by: 'alice@example.com',
      summary: 'Initial API spec',
      operation_count: 12,
      schema_changes: 0,
      breaking_changes: 0,
    },
    {
      id: 'spec-2',
      project_id: projectId,
      version: '1.1.0',
      uploaded_at: '2025-03-15T00:00:00Z',
      uploaded_by: 'bob@example.com',
      summary: 'Added new endpoints, no breaking changes',
      operation_count: 15,
      schema_changes: 3,
      breaking_changes: 0,
    },
    {
      id: 'spec-3',
      project_id: projectId,
      version: '2.0.0',
      uploaded_at: '2025-05-30T00:00:00Z',
      uploaded_by: 'alice@example.com',
      summary: 'Major refactor with breaking changes',
      operation_count: 18,
      schema_changes: 8,
      breaking_changes: 2,
    },
  ];
}

export async function mockGetSpecDiff(fromSpecId: string, toSpecId: string): Promise<SpecDiff> {
  await new Promise(resolve => setTimeout(resolve, 400));
  
  return {
    from_spec_id: fromSpecId,
    to_spec_id: toSpecId,
    changes: [
      {
        type: 'added',
        operation_id: 'deleteUser',
        path: '/users/{id}',
        method: 'DELETE',
        breaking: false,
        details: 'New endpoint to delete users',
      },
      {
        type: 'modified',
        operation_id: 'getUser',
        path: '/users/{id}',
        method: 'GET',
        breaking: true,
        details: 'Response schema changed: removed "internal_id" field (breaking)',
      },
      {
        type: 'removed',
        operation_id: 'searchUsers',
        path: '/users/search',
        method: 'POST',
        breaking: true,
        details: 'Deprecated endpoint removed',
      },
    ],
  };
}

export async function mockGetImpactAnalysis(projectId: string, specVersionId: string): Promise<ImpactAnalysis> {
  await new Promise(resolve => setTimeout(resolve, 400));
  
  return {
    spec_version_id: specVersionId,
    affected_flows: [
      {
        flow_id: 'flow-1',
        flow_name: 'User Onboarding',
        affected_steps: 2,
        impact: 'breaking',
      },
    ],
    affected_suites: [
      {
        suite_id: 'suite-1',
        suite_name: 'User Management Tests',
        affected_operations: 3,
        impact: 'compatible',
      },
    ],
    affected_runs: 45,
    risk_assessment: 'high',
  };
}

// ============================================================================
// TIER 4: ANALYTICS MOCKS
// ============================================================================

export async function mockGetLatencyAnalytics(projectId: string, range: string = '7d'): Promise<LatencyAnalytics> {
  await new Promise(resolve => setTimeout(resolve, 500));
  
  const data = [];
  const now = new Date();
  for (let i = 6; i >= 0; i--) {
    const date = new Date(now);
    date.setDate(date.getDate() - i);
    data.push({
      timestamp: date.toISOString(),
      p50: 100 + Math.random() * 50,
      p95: 400 + Math.random() * 200,
      p99: 1500 + Math.random() * 500,
    });
  }
  
  return {
    range,
    data,
    slowest_operations: [
      { operation_id: 'getUser', avg_latency: 450, p99_latency: 2100 },
      { operation_id: 'listUsers', avg_latency: 380, p99_latency: 1800 },
    ],
    fastest_operations: [
      { operation_id: 'health', avg_latency: 10 },
      { operation_id: 'getStatus', avg_latency: 25 },
    ],
  };
}

export async function mockGetFlakyOperations(projectId: string, range: string = '30d'): Promise<FlakyOperationsAnalytics> {
  await new Promise(resolve => setTimeout(resolve, 500));
  
  return {
    range,
    operations: [
      {
        operation_id: 'getUser',
        path: '/users/{id}',
        method: 'GET',
        success_rate: 92,
        runs: 100,
        failures: 8,
        flakiness_score: 45,
        trend: 'improving',
      },
      {
        operation_id: 'listUsers',
        path: '/users',
        method: 'GET',
        success_rate: 88,
        runs: 100,
        failures: 12,
        flakiness_score: 65,
        trend: 'degrading',
      },
    ],
    critical_flaky: [
      {
        operation_id: 'createUser',
        path: '/users',
        method: 'POST',
        success_rate: 75,
        runs: 100,
        failures: 25,
        flakiness_score: 88,
        trend: 'degrading',
      },
    ],
  };
}

export async function mockGetHealthTrends(projectId: string, range: string = '90d'): Promise<HealthTrendsAnalytics> {
  await new Promise(resolve => setTimeout(resolve, 500));
  
  const data = [];
  const now = new Date();
  for (let i = 29; i >= 0; i--) {
    const date = new Date(now);
    date.setDate(date.getDate() - i);
    const passRate = 85 + Math.random() * 15;
    const totalRuns = 50 + Math.floor(Math.random() * 100);
    data.push({
      timestamp: date.toISOString(),
      date: date.toISOString().split('T')[0],
      total_runs: totalRuns,
      passed_runs: Math.floor(totalRuns * passRate / 100),
      failed_runs: Math.floor(totalRuns * (100 - passRate) / 100),
      pass_rate: passRate,
    });
  }
  
  return {
    range,
    data,
    current_health: 92,
    trend: 'improving',
    weekly_average: 90,
    monthly_average: 88,
  };
}

// ============================================================================
// TIER 5: GOVERNANCE MOCKS
// ============================================================================

export async function mockGetProjectMembers(projectId: string): Promise<ProjectMember[]> {
  await new Promise(resolve => setTimeout(resolve, 300));
  
  return [
    {
      id: 'member-1',
      project_id: projectId,
      user_id: 'user-1',
      user_name: 'Alice Chen',
      user_email: 'alice@example.com',
      role: 'admin',
      added_at: '2025-01-01T00:00:00Z',
      added_by: 'system',
    },
    {
      id: 'member-2',
      project_id: projectId,
      user_id: 'user-2',
      user_name: 'Bob Smith',
      user_email: 'bob@example.com',
      role: 'editor',
      added_at: '2025-02-15T00:00:00Z',
      added_by: 'user-1',
    },
  ];
}

export async function mockGetAuditLog(projectId: string, limit: number = 100): Promise<AuditLogEntry[]> {
  await new Promise(resolve => setTimeout(resolve, 400));
  
  return [
    {
      id: 'audit-1',
      project_id: projectId,
      action: 'update',
      resource_type: 'flow',
      resource_id: 'flow-1',
      actor_id: 'user-1',
      actor_name: 'Alice Chen',
      timestamp: '2025-05-30T10:30:00Z',
      changes: [
        { field: 'name', old_value: 'Onboarding Flow', new_value: 'User Onboarding Flow' },
      ],
    },
    {
      id: 'audit-2',
      project_id: projectId,
      action: 'run',
      resource_type: 'flow',
      resource_id: 'flow-1',
      actor_id: 'user-2',
      actor_name: 'Bob Smith',
      timestamp: '2025-05-30T10:25:00Z',
    },
  ];
}

export async function mockGetWebhooks(projectId: string): Promise<Webhook[]> {
  await new Promise(resolve => setTimeout(resolve, 300));
  
  return [
    {
      id: 'webhook-1',
      project_id: projectId,
      name: 'Slack Notifications',
      url: 'https://hooks.slack.com/services/T123/B456/xxx',
      events: ['run.completed', 'run.failed'],
      is_active: true,
      created_at: '2025-04-01T00:00:00Z',
      last_triggered_at: '2025-05-30T10:00:00Z',
    },
  ];
}

export async function mockGetNotificationRules(projectId: string): Promise<NotificationRule[]> {
  await new Promise(resolve => setTimeout(resolve, 300));
  
  return [
    {
      id: 'rule-1',
      project_id: projectId,
      name: 'Critical Failures',
      channel: 'email',
      channel_config: { email: 'team@example.com' },
      triggers: [
        {
          event: 'run.failed',
          conditions: [{ field: 'severity', operator: 'eq', value: 'critical' }],
        },
      ],
      is_active: true,
      created_at: '2025-03-01T00:00:00Z',
    },
  ];
}
