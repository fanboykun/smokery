<script lang="ts">
  import { page } from '$app/stores';
  import { createQuery } from '@tanstack/svelte-query';
  import { mockGetContractReport } from '$lib/api/mock-phase2';
  import * as Card from '$lib/components/ui/card';
  import { Badge } from '$lib/components/ui/badge';
  import { AlertCircle, CheckCircle2, AlertTriangle } from '@lucide/svelte';

  const runId = $page.params.runId!;

  const report = createQuery(() => ({
    queryKey: ['runs', runId, 'report', 'contract'],
    queryFn: () => mockGetContractReport(runId),
  }));
</script>

<main class="space-y-6 px-6 py-8">
  {#if report.isPending}
    <p class="text-muted-foreground">Loading contract report…</p>
  {:else if report.isError}
    <div class="rounded-lg border border-destructive/50 bg-destructive/5 p-4">
      <p class="text-sm text-destructive">Failed to load report</p>
    </div>
  {:else if report.data}
    <!-- Summary -->
    <Card.Root class="border-0 bg-gradient-to-r from-emerald-500/10 to-transparent">
      <Card.Content class="flex items-center justify-between py-6">
        <div>
          <p class="text-sm text-muted-foreground">Compliance Score</p>
          <p class="mt-1 text-4xl font-bold text-emerald-400">{report.data.compliance_score}%</p>
        </div>
        <div class="space-y-2 text-right">
          <div class="flex items-center justify-end gap-2">
            <CheckCircle2 class="size-4 text-emerald-500" />
            <span class="text-sm">{report.data.passed_assertions} passed</span>
          </div>
          <div class="flex items-center justify-end gap-2">
            <AlertCircle class="size-4 text-red-500" />
            <span class="text-sm">{report.data.failed_assertions} failed</span>
          </div>
        </div>
      </Card.Content>
    </Card.Root>

    <!-- Errors Section -->
    {#if report.data.errors.length > 0}
      <Card.Root class="border-destructive/50 bg-destructive/5">
        <Card.Header>
          <Card.Title class="flex items-center gap-2 text-lg text-destructive">
            <AlertCircle class="size-5" />
            {report.data.errors.length} Contract Violations
          </Card.Title>
        </Card.Header>
        <Card.Content class="space-y-4">
          {#each report.data.errors as error (error.id)}
            <div class="rounded-lg border border-destructive/20 bg-destructive/10 p-4 space-y-3">
              <div class="flex items-start justify-between gap-2">
                <div>
                  <p class="font-semibold text-destructive">{error.message}</p>
                  <p class="mt-1 text-xs text-muted-foreground">
                    Operation: <code class="rounded bg-black/30 px-1.5 py-0.5">{error.operation_id}</code>
                  </p>
                  <p class="text-xs text-muted-foreground">Location: {error.location}</p>
                </div>
                <Badge class="bg-red-500/20 text-red-300 shrink-0">{error.violation_type.replace(/_/g, ' ')}</Badge>
              </div>
              
              <div class="grid gap-2 sm:grid-cols-2 text-xs">
                <div class="rounded bg-black/20 p-2">
                  <p class="font-semibold text-muted-foreground">Expected</p>
                  <code class="block mt-1 overflow-auto text-emerald-400">
                    {JSON.stringify(error.expected_schema, null, 2).slice(0, 150)}...
                  </code>
                </div>
                <div class="rounded bg-black/20 p-2">
                  <p class="font-semibold text-muted-foreground">Actual</p>
                  <code class="block mt-1 overflow-auto text-red-400">
                    {JSON.stringify(error.actual_value, null, 2).slice(0, 150)}...
                  </code>
                </div>
              </div>
            </div>
          {/each}
        </Card.Content>
      </Card.Root>
    {/if}

    <!-- Warnings Section -->
    {#if report.data.warnings.length > 0}
      <Card.Root class="border-yellow-500/30">
        <Card.Header>
          <Card.Title class="flex items-center gap-2 text-lg text-yellow-400">
            <AlertTriangle class="size-5" />
            {report.data.warnings.length} Warnings
          </Card.Title>
        </Card.Header>
        <Card.Content class="space-y-3">
          {#each report.data.warnings as warning (warning.id)}
            <div class="rounded-lg border border-yellow-500/20 bg-yellow-500/5 p-3 space-y-2">
              <div class="flex items-start justify-between gap-2">
                <div>
                  <p class="font-semibold text-yellow-400">{warning.message}</p>
                  <p class="mt-1 text-xs text-muted-foreground">
                    <code class="rounded bg-black/30 px-1.5 py-0.5">{warning.operation_id}</code>
                  </p>
                </div>
                <Badge class="bg-yellow-500/20 text-yellow-300 shrink-0">warning</Badge>
              </div>
            </div>
          {/each}
        </Card.Content>
      </Card.Root>
    {/if}
  {/if}
</main>
