<script lang="ts">
  import { page } from '$app/stores';
  import { createQuery } from '@tanstack/svelte-query';
  import { mockGetFlakyOperations } from '$lib/api/mock-phase2';
  import * as Card from '$lib/components/ui/card';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';
  import { AlertTriangle, TrendingUp, TrendingDown, Minus } from '@lucide/svelte';

  const projectId = $page.params.id!;
  let range = $state('30d');

  const analytics = createQuery(() => ({
    queryKey: ['analytics-flaky', projectId, range],
    queryFn: () => mockGetFlakyOperations(projectId, range),
  }));

  function trendIcon(trend: string) {
    if (trend === 'improving') return TrendingUp;
    if (trend === 'degrading') return TrendingDown;
    return Minus;
  }

  function trendColor(trend: string) {
    if (trend === 'improving') return 'text-emerald-400';
    if (trend === 'degrading') return 'text-red-400';
    return 'text-muted-foreground';
  }

  function scoreColor(score: number) {
    if (score >= 70) return 'text-red-400';
    if (score >= 40) return 'text-yellow-400';
    return 'text-emerald-400';
  }
</script>

<main class="mx-auto max-w-5xl px-6 py-8 space-y-6">
  <div class="flex items-center justify-between">
    <div>
      <p class="text-xs font-bold uppercase tracking-widest text-primary">Analytics</p>
      <h1 class="text-2xl font-bold">Flaky Operations</h1>
    </div>
    <div class="flex gap-1">
      {#each ['7d', '30d', '90d'] as r}
        <Button variant={range === r ? 'default' : 'outline'} size="sm" onclick={() => (range = r)}>{r}</Button>
      {/each}
    </div>
  </div>

  {#if analytics.isPending}
    <p class="text-muted-foreground">Loading…</p>
  {:else if analytics.data}
    <!-- Critical Flaky -->
    {#if analytics.data.critical_flaky.length > 0}
      <Card.Root class="border-destructive/40 bg-destructive/5">
        <Card.Header>
          <Card.Title class="flex items-center gap-2 text-base text-destructive">
            <AlertTriangle class="size-4" />
            Critical ({analytics.data.critical_flaky.length})
          </Card.Title>
        </Card.Header>
        <Card.Content class="space-y-3">
          {#each analytics.data.critical_flaky as op (op.operation_id)}
            {@const Trend = trendIcon(op.trend)}
            <div class="flex items-center justify-between rounded-lg border border-destructive/20 bg-background p-3">
              <div>
                <div class="flex items-center gap-2">
                  <Badge variant="outline" class="font-mono text-xs">{op.method}</Badge>
                  <span class="font-mono text-sm">{op.path}</span>
                </div>
                <p class="mt-1 text-xs text-muted-foreground">{op.failures}/{op.runs} failures</p>
              </div>
              <div class="flex items-center gap-3">
                <span class={`text-lg font-bold ${scoreColor(op.flakiness_score)}`}>{op.flakiness_score}</span>
                <Trend class={`size-4 ${trendColor(op.trend)}`} />
              </div>
            </div>
          {/each}
        </Card.Content>
      </Card.Root>
    {/if}

    <!-- All Flaky Operations -->
    <Card.Root>
      <Card.Header><Card.Title class="text-base">All Flaky Operations ({analytics.data.operations.length})</Card.Title></Card.Header>
      <Card.Content class="space-y-2">
        {#each analytics.data.operations as op (op.operation_id)}
          {@const Trend = trendIcon(op.trend)}
          <div class="flex items-center justify-between rounded-lg border border-border p-3">
            <div>
              <div class="flex items-center gap-2">
                <Badge variant="outline" class="font-mono text-xs">{op.method}</Badge>
                <span class="font-mono text-sm">{op.path}</span>
              </div>
              <p class="mt-1 text-xs text-muted-foreground">
                {op.success_rate}% success • {op.failures}/{op.runs} failures
              </p>
            </div>
            <div class="flex items-center gap-3">
              <div class="text-right">
                <p class={`text-sm font-bold ${scoreColor(op.flakiness_score)}`}>Score: {op.flakiness_score}</p>
                <p class={`text-xs ${trendColor(op.trend)}`}>{op.trend}</p>
              </div>
              <Trend class={`size-4 ${trendColor(op.trend)}`} />
            </div>
          </div>
        {/each}
      </Card.Content>
    </Card.Root>
  {/if}
</main>
