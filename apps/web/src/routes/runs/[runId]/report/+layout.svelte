<script lang="ts">
  import { page } from '$app/stores';
  import * as Tabs from '$lib/components/ui/tabs';
  import { Button } from '$lib/components/ui/button';
  import { ArrowLeft } from '@lucide/svelte';

  const runId = $page.params.runId!;

  // Map the current path to the active tab
  const currentPath = $page.url.pathname;
  const activeTab = $derived.by(() => {
    if (currentPath.includes('contract')) return 'contract';
    if (currentPath.includes('analyst')) return 'analyst';
    if (currentPath.includes('qa')) return 'qa';
    if (currentPath.includes('correlation')) return 'correlation';
    return 'contract';
  });
</script>

<main class="min-h-screen bg-background">
  <!-- Header -->
  <div class="border-b border-border bg-card sticky top-0 z-10">
    <div class="mx-auto max-w-7xl px-6 py-4">
      <div class="flex items-center justify-between mb-4">
        <div>
          <a href="/runs/{runId}" class="inline-flex items-center gap-2 text-sm text-muted-foreground hover:text-foreground transition-colors">
            <ArrowLeft class="size-4" />
            Back to Run
          </a>
          <h1 class="text-2xl font-bold mt-2">Run Report</h1>
        </div>
      </div>

      <!-- Tabs -->
      <Tabs.Root value={activeTab}>
        <Tabs.List class="w-full">
          <Tabs.Trigger value="contract" asChild>
            <a href="/runs/{runId}/report/contract" class="cursor-pointer">
              Contract Compliance
            </a>
          </Tabs.Trigger>
          <Tabs.Trigger value="analyst" asChild>
            <a href="/runs/{runId}/report/analyst" class="cursor-pointer">
              Analyst View
            </a>
          </Tabs.Trigger>
          <Tabs.Trigger value="qa" asChild>
            <a href="/runs/{runId}/report/qa" class="cursor-pointer">
              QA Summary
            </a>
          </Tabs.Trigger>
          <Tabs.Trigger value="correlation" asChild>
            <a href="/runs/{runId}/report/correlation" class="cursor-pointer">
              Correlation
            </a>
          </Tabs.Trigger>
        </Tabs.List>
      </Tabs.Root>
    </div>
  </div>

  <!-- Content -->
  <slot />
</main>
