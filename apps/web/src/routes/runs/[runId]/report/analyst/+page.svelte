<script lang="ts">
  import { page } from '$app/stores';
  import { createQuery } from '@tanstack/svelte-query';
  import { mockGetAnalystReport } from '$lib/api/mock-phase2';
  import * as Card from '$lib/components/ui/card';
  import { Badge } from '$lib/components/ui/badge';
  import { AlertCircle, TrendingUp, Lightbulb } from '@lucide/svelte';

  const runId = $page.params.runId!;

  const report = createQuery(() => ({
    queryKey: ['runs', runId, 'report', 'analyst'],
    queryFn: () => mockGetAnalystReport(runId),
  }));
</script>

<main class="space-y-6 px-6 py-8">
  {#if report.isPending}
    <p class="text-muted-foreground">Loading analyst report…</p>
  {:else if report.isError}
    <div class="rounded-lg border border-destructive/50 bg-destructive/5 p-4">
      <p class="text-sm text-destructive">Failed to load report</p>
    </div>
  {:else if report.data}
    <!-- Summary -->
    <Card.Root>
      <Card.Header><Card.Title>Analysis Summary</Card.Title></Card.Header>
      <Card.Content>
        <p class="text-muted-foreground">{report.data.summary}</p>
      </Card.Content>
    </Card.Root>

    <!-- Root Causes -->
    <Card.Root>
      <Card.Header>
        <Card.Title class="flex items-center gap-2">
          <AlertCircle class="size-5 text-orange-500" />
          Root Causes ({report.data.root_causes.length})
        </Card.Title>
      </Card.Header>
      <Card.Content class="space-y-4">
        {#each report.data.root_causes as cause (cause.cause)}
          <div class="rounded-lg border border-border p-4">
            <div class="flex items-start justify-between gap-4 mb-3">
              <div>
                <p class="font-semibold">{cause.cause}</p>
                <p class="text-sm text-muted-foreground mt-1">
                  Responsible for {cause.impact}% of failures
                </p>
              </div>
              <div class="flex-shrink-0 text-right">
                <div class="text-2xl font-bold text-orange-400">{cause.impact}%</div>
              </div>
            </div>
            <div class="space-y-2">
              <p class="text-xs font-semibold text-muted-foreground">Affected Operations:</p>
              <div class="flex flex-wrap gap-2">
                {#each cause.affected_operations as op (op)}
                  <Badge variant="outline" class="font-mono text-xs">{op}</Badge>
                {/each}
              </div>
            </div>
          </div>
        {/each}
      </Card.Content>
    </Card.Root>

    <!-- Timeline -->
    <Card.Root>
      <Card.Header><Card.Title>Execution Timeline</Card.Title></Card.Header>
      <Card.Content>
        <div class="space-y-3">
          {#each report.data.timeline_insights as event (event.timestamp)}
            <div class="flex gap-4 pb-3 border-b border-border last:border-b-0">
              <div class="text-xs text-muted-foreground shrink-0 w-24">
                {new Date(event.timestamp).toLocaleTimeString()}
              </div>
              <div class="flex-1">
                <div class="flex items-center gap-2">
                  <code class="text-sm font-mono">{event.operation_id}</code>
                  <Badge variant={event.status === 'passed' ? 'default' : 'destructive'} class="text-xs">
                    {event.status}
                  </Badge>
                  <span class="text-xs text-muted-foreground">{event.duration_ms}ms</span>
                </div>
                {#if event.error}
                  <p class="text-xs text-red-400 mt-1">{event.error}</p>
                {/if}
              </div>
            </div>
          {/each}
        </div>
      </Card.Content>
    </Card.Root>

    <!-- Recommendations -->
    <Card.Root class="border-emerald-500/30 bg-emerald-500/5">
      <Card.Header>
        <Card.Title class="flex items-center gap-2 text-emerald-400">
          <Lightbulb class="size-5" />
          Recommendations
        </Card.Title>
      </Card.Header>
      <Card.Content class="space-y-4">
        {#each report.data.recommendations as rec (rec.title)}
          <div class="rounded-lg border border-emerald-500/20 p-4">
            <div class="flex items-start justify-between gap-4 mb-2">
              <p class="font-semibold text-emerald-400">{rec.title}</p>
              <Badge class={rec.priority === 'high' ? 'bg-red-500/20 text-red-300' : rec.priority === 'medium' ? 'bg-yellow-500/20 text-yellow-300' : 'bg-blue-500/20 text-blue-300'}>
                {rec.priority}
              </Badge>
            </div>
            <p class="text-sm text-muted-foreground">{rec.description}</p>
          </div>
        {/each}
      </Card.Content>
    </Card.Root>
  {/if}
</main>
