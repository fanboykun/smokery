<script lang="ts">
  import { page } from '$app/stores';
  import { createQuery } from '@tanstack/svelte-query';
  import { mockGetImpactAnalysis } from '$lib/api/mock-phase2';
  import * as Card from '$lib/components/ui/card';
  import { Badge } from '$lib/components/ui/badge';
  import { ArrowLeft, AlertTriangle, Workflow, TestTube2, Activity } from '@lucide/svelte';

  const projectId = $page.params.id!;
  const specVersionId = $page.url.searchParams.get('spec-version') ?? '';

  const impact = createQuery(() => ({
    queryKey: ['impact-analysis', projectId, specVersionId],
    queryFn: () => mockGetImpactAnalysis(projectId, specVersionId),
    enabled: !!specVersionId,
  }));

  function riskColor(risk: string) {
    if (risk === 'high') return 'destructive' as const;
    if (risk === 'medium') return 'secondary' as const;
    return 'default' as const;
  }

  function impactColor(imp: string) {
    if (imp === 'breaking') return 'text-red-400 bg-red-500/10';
    return 'text-emerald-400 bg-emerald-500/10';
  }
</script>

<main class="mx-auto max-w-4xl px-6 py-8 space-y-6">
  <div>
    <a href="/projects/{projectId}/spec/versions" class="inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground">
      <ArrowLeft class="size-4" />
      Back to versions
    </a>
    <h1 class="mt-2 text-2xl font-bold">Impact Analysis</h1>
    <p class="text-sm text-muted-foreground">How spec changes affect your tests</p>
  </div>

  {#if !specVersionId}
    <Card.Root>
      <Card.Content class="py-8 text-center text-muted-foreground">
        No spec version selected. Go to <a href="/projects/{projectId}/spec/diff" class="underline">diff viewer</a> first.
      </Card.Content>
    </Card.Root>
  {:else if impact.isPending}
    <p class="text-muted-foreground">Analyzing impact…</p>
  {:else if impact.data}
    <!-- Risk Summary -->
    <Card.Root class={impact.data.risk_assessment === 'high' ? 'border-destructive/40 bg-destructive/5' : ''}>
      <Card.Content class="flex items-center justify-between py-6">
        <div>
          <p class="text-sm text-muted-foreground">Risk Assessment</p>
          <p class="mt-1 text-3xl font-bold capitalize">{impact.data.risk_assessment}</p>
        </div>
        <div class="text-right space-y-1">
          <Badge variant={riskColor(impact.data.risk_assessment)} class="text-sm">
            <AlertTriangle class="size-3" />
            {impact.data.affected_runs} runs affected
          </Badge>
        </div>
      </Card.Content>
    </Card.Root>

    <!-- Affected Flows -->
    {#if impact.data.affected_flows.length > 0}
      <Card.Root>
        <Card.Header>
          <Card.Title class="flex items-center gap-2 text-base">
            <Workflow class="size-4" />
            Affected Flows ({impact.data.affected_flows.length})
          </Card.Title>
        </Card.Header>
        <Card.Content class="space-y-2">
          {#each impact.data.affected_flows as flow (flow.flow_id)}
            <div class="flex items-center justify-between rounded-lg border border-border p-3">
              <div>
                <p class="font-medium">{flow.flow_name}</p>
                <p class="text-xs text-muted-foreground">{flow.affected_steps} steps affected</p>
              </div>
              <Badge class={impactColor(flow.impact)}>{flow.impact}</Badge>
            </div>
          {/each}
        </Card.Content>
      </Card.Root>
    {/if}

    <!-- Affected Suites -->
    {#if impact.data.affected_suites.length > 0}
      <Card.Root>
        <Card.Header>
          <Card.Title class="flex items-center gap-2 text-base">
            <TestTube2 class="size-4" />
            Affected Suites ({impact.data.affected_suites.length})
          </Card.Title>
        </Card.Header>
        <Card.Content class="space-y-2">
          {#each impact.data.affected_suites as suite (suite.suite_id)}
            <div class="flex items-center justify-between rounded-lg border border-border p-3">
              <div>
                <p class="font-medium">{suite.suite_name}</p>
                <p class="text-xs text-muted-foreground">{suite.affected_operations} operations affected</p>
              </div>
              <Badge class={impactColor(suite.impact)}>{suite.impact}</Badge>
            </div>
          {/each}
        </Card.Content>
      </Card.Root>
    {/if}
  {/if}
</main>
