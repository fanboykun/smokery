<script lang="ts">
  import { page } from '$app/stores';
  import { createQuery } from '@tanstack/svelte-query';
  import { api } from '$lib/api/client';
  import { createProjectConfigStore } from '$lib/stores/project-config';
  import SplitPane from '$lib/components/SplitPane.svelte';
  import * as Card from '$lib/components/ui/card';
  import * as Tabs from '$lib/components/ui/tabs';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';
  import { AlertCircle, CheckCircle2, Zap } from '@lucide/svelte';

  const projectId = $page.params.id!;
  const config = createProjectConfigStore(projectId);

  let activeTab = $state('flows');
  let previewError = $state<string | null>(null);
  let previewLoading = $state(false);

  const project = createQuery(() => ({
    queryKey: ['projects', projectId],
    queryFn: async () => {
      const { data, error } = await api.GET('/api/projects/{id}', { params: { path: { id: projectId } } });
      if (error) throw error;
      return data!;
    },
  }));

  const plan = createQuery(() => ({
    queryKey: ['projects', projectId, 'plan'],
    queryFn: async () => {
      previewLoading = true;
      try {
        const { data, error } = await api.POST('/api/projects/{id}/plan/preview', {
          params: { path: { id: projectId } },
        });
        if (error) {
          previewError = (error as any).detail ?? 'Compilation failed';
          return null;
        }
        previewError = null;
        return data!;
      } finally {
        previewLoading = false;
      }
    },
    refetchInterval: 5000,
    staleTime: 0,
  }));
</script>

<div class="flex h-[calc(100vh-3.5rem)] flex-col bg-background">
  <!-- Header -->
  <div class="border-b border-border bg-card px-6 py-4 shadow-sm">
    <div class="mx-auto max-w-7xl">
      <h1 class="text-2xl font-bold">{project.data?.name ?? 'Loading…'}</h1>
      <p class="mt-1 text-sm text-muted-foreground">Configure flows, suites, and environments in real-time</p>
    </div>
  </div>

  <!-- Split Pane -->
  <SplitPane leftLabel="Configuration" rightLabel="Live Preview" class="flex-1">
    <svelte:fragment slot="left">
      <div class="space-y-4 p-6">
        <!-- Quick Stats -->
        <div class="grid grid-cols-2 gap-3">
          <div class="rounded-lg border border-border bg-muted/30 p-3 text-center">
            <p class="text-2xl font-bold">{$config.flows.length}</p>
            <p class="text-xs text-muted-foreground">Flows</p>
          </div>
          <div class="rounded-lg border border-border bg-muted/30 p-3 text-center">
            <p class="text-2xl font-bold">{$config.suites.length}</p>
            <p class="text-xs text-muted-foreground">Suites</p>
          </div>
          <div class="rounded-lg border border-border bg-muted/30 p-3 text-center">
            <p class="text-2xl font-bold">{$config.environments.length}</p>
            <p class="text-xs text-muted-foreground">Envs</p>
          </div>
          <div class="rounded-lg border border-border bg-muted/30 p-3 text-center">
            <p class="text-2xl font-bold">{$config.auth_profiles.length}</p>
            <p class="text-xs text-muted-foreground">Auth</p>
          </div>
        </div>

        <!-- Tabs -->
        <Tabs.Root value={activeTab} onvaluechange={(v) => (activeTab = v)}>
          <Tabs.List class="grid w-full grid-cols-3">
            <Tabs.Trigger value="flows">Flows</Tabs.Trigger>
            <Tabs.Trigger value="suites">Suites</Tabs.Trigger>
            <Tabs.Trigger value="config">Config</Tabs.Trigger>
          </Tabs.List>

          <!-- Flows Tab -->
          <Tabs.Content value="flows" class="space-y-3 pt-3">
            <Button variant="outline" href="/projects/{projectId}/flows/new" class="w-full">+ New Flow</Button>
            {#if $config.flows.length === 0}
              <p class="text-xs text-muted-foreground">No flows yet. Create one to define test scenarios.</p>
            {:else}
              <div class="space-y-2">
                {#each $config.flows as flow (flow.id)}
                  <a
                    href="/projects/{projectId}/flows/{flow.id}"
                    class="block rounded-lg border border-border p-2 transition-colors hover:bg-muted"
                  >
                    <p class="text-sm font-medium">{flow.name}</p>
                    <p class="text-xs text-muted-foreground">{flow.steps.length} steps</p>
                  </a>
                {/each}
              </div>
            {/if}
          </Tabs.Content>

          <!-- Suites Tab -->
          <Tabs.Content value="suites" class="space-y-3 pt-3">
            <Button variant="outline" href="/projects/{projectId}/suites/new" class="w-full">+ New Suite</Button>
            {#if $config.suites.length === 0}
              <p class="text-xs text-muted-foreground">No suites yet. Create one to generate test combinations.</p>
            {:else}
              <div class="space-y-2">
                {#each $config.suites as suite (suite.id)}
                  <a
                    href="/projects/{projectId}/suites/{suite.id}"
                    class="block rounded-lg border border-border p-2 transition-colors hover:bg-muted"
                  >
                    <p class="text-sm font-medium">{suite.name}</p>
                    <Badge variant="secondary" class="text-xs">{suite.operations.length} ops</Badge>
                  </a>
                {/each}
              </div>
            {/if}
          </Tabs.Content>

          <!-- Config Tab -->
          <Tabs.Content value="config" class="space-y-3 pt-3">
            <Button variant="outline" href="/projects/{projectId}/environments" class="w-full">Environments</Button>
            <Button variant="outline" href="/projects/{projectId}/operations" class="w-full">Operations</Button>
          </Tabs.Content>
        </Tabs.Root>
      </div>
    </svelte:fragment>

    <svelte:fragment slot="right">
      <div class="space-y-6 p-6">
        <!-- Status -->
        <div class="rounded-lg border border-border bg-card p-4">
          {#if previewLoading}
            <div class="flex items-center gap-2 text-sm text-muted-foreground">
              <div class="size-2 animate-pulse rounded-full bg-blue-500" />
              Compiling…
            </div>
          {:else if previewError}
            <div class="flex items-start gap-3">
              <AlertCircle class="mt-0.5 size-4 flex-shrink-0 text-destructive" />
              <div>
                <p class="text-sm font-medium text-destructive">Compilation error</p>
                <p class="mt-1 text-xs text-muted-foreground">{previewError}</p>
              </div>
            </div>
          {:else if plan.data}
            <div class="flex items-start gap-3">
              <CheckCircle2 class="mt-0.5 size-4 flex-shrink-0 text-emerald-500" />
              <div>
                <p class="text-sm font-medium text-emerald-400">Ready</p>
                <p class="mt-1 text-xs text-muted-foreground">Plan compiled successfully</p>
              </div>
            </div>
          {:else}
            <p class="text-xs text-muted-foreground">Waiting for configuration…</p>
          {/if}
        </div>

        <!-- Plan Stats -->
        {#if plan.data}
          <div class="rounded-lg border border-border bg-card p-4">
            <p class="mb-3 text-xs font-semibold uppercase tracking-widest text-muted-foreground">Plan Summary</p>
            <div class="space-y-2">
              <div class="flex items-center justify-between text-sm">
                <span class="text-muted-foreground">Test Cases</span>
                <span class="font-semibold">{plan.data.test_cases?.length ?? 0}</span>
              </div>
              <div class="flex items-center justify-between text-sm">
                <span class="text-muted-foreground">Total Steps</span>
                <span class="font-semibold">{plan.data.test_cases?.reduce((sum, tc) => sum + (tc.steps?.length ?? 0), 0) ?? 0}</span>
              </div>
            </div>
          </div>

          <!-- Quick Actions -->
          <Button href="/projects/{projectId}/plan" class="w-full gap-2" variant="default">
            <Zap class="size-4" />
            View Full Plan
          </Button>
          <Button href="/projects/{projectId}/runs" class="w-full" variant="outline">
            Run Tests
          </Button>
        {/if}
      </div>
    </svelte:fragment>
  </SplitPane>
</div>
