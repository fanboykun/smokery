<script lang="ts">
  import { page } from '$app/stores';
  import { createQuery } from '@tanstack/svelte-query';
  import { mockGetQAReport } from '$lib/api/mock-phase2';
  import * as Card from '$lib/components/ui/card';
  import { Badge } from '$lib/components/ui/badge';
  import { CheckCircle2, AlertCircle, Zap } from '@lucide/svelte';

  const runId = $page.params.runId!;

  const report = createQuery(() => ({
    queryKey: ['runs', runId, 'report', 'qa'],
    queryFn: () => mockGetQAReport(runId),
  }));
</script>

<main class="space-y-6 px-6 py-8">
  {#if report.isPending}
    <p class="text-muted-foreground">Loading QA report…</p>
  {:else if report.isError}
    <div class="rounded-lg border border-destructive/50 bg-destructive/5 p-4">
      <p class="text-sm text-destructive">Failed to load report</p>
    </div>
  {:else if report.data}
    <!-- Status Banner -->
    <Card.Root class={report.data.status === 'passed' ? 'border-emerald-500/50 bg-emerald-500/5' : 'border-destructive/50 bg-destructive/5'}>
      <Card.Content class="flex items-center justify-between py-6">
        <div class="flex items-center gap-4">
          {#if report.data.status === 'passed'}
            <CheckCircle2 class="size-8 text-emerald-500" />
            <div>
              <p class="font-semibold text-emerald-400">All Tests Passed</p>
              <p class="text-sm text-muted-foreground">{report.data.total_tests} tests executed successfully</p>
            </div>
          {:else}
            <AlertCircle class="size-8 text-destructive" />
            <div>
              <p class="font-semibold text-destructive">{report.data.failed_tests} Test{report.data.failed_tests !== 1 ? 's' : ''} Failed</p>
              <p class="text-sm text-muted-foreground">{report.data.passed_tests} passed, {report.data.failed_tests} failed</p>
            </div>
          {/if}
        </div>
        <div class="text-right">
          <p class="text-4xl font-bold">{report.data.pass_rate}%</p>
          <p class="text-xs text-muted-foreground">Pass Rate</p>
        </div>
      </Card.Content>
    </Card.Root>

    <!-- Test Results -->
    <div class="grid gap-4 md:grid-cols-3">
      <Card.Root>
        <Card.Content class="text-center py-6">
          <p class="text-3xl font-bold">{report.data.total_tests}</p>
          <p class="text-xs text-muted-foreground mt-2">Total Tests</p>
        </Card.Content>
      </Card.Root>
      <Card.Root>
        <Card.Content class="text-center py-6">
          <p class="text-3xl font-bold text-emerald-400">{report.data.passed_tests}</p>
          <p class="text-xs text-muted-foreground mt-2">Passed</p>
        </Card.Content>
      </Card.Root>
      <Card.Root>
        <Card.Content class="text-center py-6">
          <p class="text-3xl font-bold text-red-400">{report.data.failed_tests}</p>
          <p class="text-xs text-muted-foreground mt-2">Failed</p>
        </Card.Content>
      </Card.Root>
    </div>

    <!-- Coverage -->
    <Card.Root>
      <Card.Header><Card.Title>API Coverage</Card.Title></Card.Header>
      <Card.Content class="space-y-4">
        <div>
          <div class="flex justify-between text-sm mb-2">
            <span>Operations Tested</span>
            <span class="font-semibold">{report.data.coverage_summary.tested_operations}/{report.data.coverage_summary.total_operations}</span>
          </div>
          <div class="h-2 rounded-full bg-muted overflow-hidden">
            <div
              class="h-full bg-emerald-500 transition-all"
              style="width: {report.data.coverage_summary.coverage_percentage}%"
            ></div>
          </div>
          <p class="text-xs text-muted-foreground mt-1">{report.data.coverage_summary.coverage_percentage}% coverage</p>
        </div>
      </Card.Content>
    </Card.Root>

    <!-- Flaky Tests -->
    {#if report.data.flaky_tests.length > 0}
      <Card.Root class="border-yellow-500/30">
        <Card.Header>
          <Card.Title class="flex items-center gap-2 text-yellow-400">
            <Zap class="size-5" />
            Flaky Tests ({report.data.flaky_tests.length})
          </Card.Title>
        </Card.Header>
        <Card.Content>
          <p class="text-sm text-muted-foreground mb-4">These tests passed but show intermittent failures:</p>
          <div class="space-y-2">
            {#each report.data.flaky_tests as test (test)}
              <div class="flex items-center gap-2 rounded-lg border border-border p-3">
                <Zap class="size-4 text-yellow-400 flex-shrink-0" />
                <code class="font-mono text-sm">{test}</code>
              </div>
            {/each}
          </div>
        </Card.Content>
      </Card.Root>
    {/if}

    <!-- Blockers -->
    {#if report.data.blockers.length > 0}
      <Card.Root class="border-destructive/50 bg-destructive/5">
        <Card.Header>
          <Card.Title class="text-destructive">Critical Blockers</Card.Title>
        </Card.Header>
        <Card.Content class="space-y-3">
          {#each report.data.blockers as blocker (blocker.operation_id)}
            <div class="rounded-lg border border-destructive/20 bg-destructive/10 p-4">
              <div class="flex items-start justify-between gap-2 mb-2">
                <p class="font-semibold text-destructive">{blocker.operation_id}</p>
                <Badge class="bg-red-500/20 text-red-300">{blocker.severity}</Badge>
              </div>
              <p class="text-sm text-muted-foreground">{blocker.issue}</p>
            </div>
          {/each}
        </Card.Content>
      </Card.Root>
    {/if}
  {/if}
</main>
