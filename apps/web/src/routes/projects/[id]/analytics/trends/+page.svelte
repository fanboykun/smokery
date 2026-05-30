<script lang="ts">
  import { page } from '$app/stores';
  import { createQuery } from '@tanstack/svelte-query';
  import { mockGetHealthTrends } from '$lib/api/mock-phase2';
  import * as Card from '$lib/components/ui/card';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';
  import { Activity, TrendingUp, TrendingDown } from '@lucide/svelte';

  const projectId = $page.params.id!;
  let range = $state('30d');

  const analytics = createQuery(() => ({
    queryKey: ['analytics-trends', projectId, range],
    queryFn: () => mockGetHealthTrends(projectId, range),
  }));

  function healthColor(pct: number) {
    if (pct >= 95) return 'text-emerald-400';
    if (pct >= 80) return 'text-yellow-400';
    return 'text-red-400';
  }
</script>

<main class="mx-auto max-w-5xl px-6 py-8 space-y-6">
  <div class="flex items-center justify-between">
    <div>
      <p class="text-xs font-bold uppercase tracking-widest text-primary">Analytics</p>
      <h1 class="text-2xl font-bold">Health Trends</h1>
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
    <!-- Health Summary -->
    <div class="grid gap-4 sm:grid-cols-3">
      <Card.Root>
        <Card.Content class="py-6 text-center">
          <p class="text-xs text-muted-foreground">Current Health</p>
          <p class={`text-4xl font-bold ${healthColor(analytics.data.current_health)}`}>{analytics.data.current_health}%</p>
          <div class="mt-1 flex items-center justify-center gap-1">
            {#if analytics.data.trend === 'improving'}
              <TrendingUp class="size-4 text-emerald-400" />
              <span class="text-xs text-emerald-400">Improving</span>
            {:else}
              <TrendingDown class="size-4 text-red-400" />
              <span class="text-xs text-red-400">Degrading</span>
            {/if}
          </div>
        </Card.Content>
      </Card.Root>
      <Card.Root>
        <Card.Content class="py-6 text-center">
          <p class="text-xs text-muted-foreground">Weekly Average</p>
          <p class={`text-4xl font-bold ${healthColor(analytics.data.weekly_average)}`}>{analytics.data.weekly_average}%</p>
        </Card.Content>
      </Card.Root>
      <Card.Root>
        <Card.Content class="py-6 text-center">
          <p class="text-xs text-muted-foreground">Monthly Average</p>
          <p class={`text-4xl font-bold ${healthColor(analytics.data.monthly_average)}`}>{analytics.data.monthly_average}%</p>
        </Card.Content>
      </Card.Root>
    </div>

    <!-- Pass Rate Chart -->
    <Card.Root>
      <Card.Header><Card.Title class="flex items-center gap-2 text-base"><Activity class="size-4" /> Pass Rate Over Time</Card.Title></Card.Header>
      <Card.Content>
        <div class="flex items-end gap-px h-32">
          {#each analytics.data.data as point}
            {@const height = point.pass_rate}
            <div
              class="flex-1 rounded-t transition-colors"
              class:bg-emerald-500={point.pass_rate >= 95}
              class:bg-yellow-500={point.pass_rate >= 80 && point.pass_rate < 95}
              class:bg-red-500={point.pass_rate < 80}
              style="height: {height}%"
              title="{point.date}: {Math.round(point.pass_rate)}% ({point.passed_runs}/{point.total_runs})"
            ></div>
          {/each}
        </div>
        <div class="flex justify-between text-xs text-muted-foreground mt-2">
          <span>{analytics.data.data[0]?.date}</span>
          <span>{analytics.data.data.at(-1)?.date}</span>
        </div>
      </Card.Content>
    </Card.Root>

    <!-- Daily Breakdown -->
    <Card.Root>
      <Card.Header><Card.Title class="text-base">Recent Days</Card.Title></Card.Header>
      <Card.Content>
        <div class="space-y-1 max-h-64 overflow-y-auto">
          {#each analytics.data.data.slice().reverse().slice(0, 10) as point (point.date)}
            <div class="flex items-center justify-between rounded px-3 py-2 text-sm hover:bg-muted/50">
              <span class="text-muted-foreground">{point.date}</span>
              <div class="flex items-center gap-3">
                <span class="text-emerald-400">{point.passed_runs} ✓</span>
                <span class="text-red-400">{point.failed_runs} ✗</span>
                <Badge variant="secondary" class={healthColor(point.pass_rate)}>{Math.round(point.pass_rate)}%</Badge>
              </div>
            </div>
          {/each}
        </div>
      </Card.Content>
    </Card.Root>
  {/if}
</main>
