<script lang="ts">
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { createMutation } from '@tanstack/svelte-query';
  import { api } from '$lib/api/client';
  import { createProjectConfigStore } from '$lib/stores/project-config';
  import * as Card from '$lib/components/ui/card';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';
  import { Separator } from '$lib/components/ui/separator';
  import Breadcrumb from '$lib/components/Breadcrumb.svelte';
  import { Loader2 } from '@lucide/svelte';
  import type { components } from '$lib/api/v1';

  type CompilerOutput = components['schemas']['Output'];

  const projectId = $page.params.id!;
  const config = createProjectConfigStore(projectId);

  let result = $state<CompilerOutput | null>(null);

  const compile = createMutation(() => ({
    mutationFn: async () => {
      const { data, error } = await api.POST('/api/projects/{project-id}/plan/preview', {
        params: { path: { 'project-id': projectId } },
        body: {
          environments: $config.environments,
          auth_profiles: $config.auth_profiles,
          flows: $config.flows.map((f) => ({
            ...f,
            steps: f.steps.map((s) => ({ ...s, params: s.params ?? undefined, body: s.body ?? undefined, headers: s.headers ?? undefined, captures: s.captures ?? null, assertions: s.assertions ?? null })),
            cleanup: (f.cleanup ?? []).map((s) => ({ ...s, params: s.params ?? undefined, body: s.body ?? undefined, headers: s.headers ?? undefined, captures: s.captures ?? null, assertions: s.assertions ?? null })),
          })),
          suites: $config.suites.map((s) => ({
            ...s,
            selector: { tags: s.selector.tags ?? [], classifications: s.selector.classifications ?? [], paths: s.selector.paths ?? [], exclude: s.selector.exclude ?? [] },
            strategy: { ...s.strategy, max_cases_per_op: s.strategy.max_cases_per_op || 0 },
          })),
        },
      });
      if (error) throw error;
      return data!;
    },
    onSuccess: (data) => { result = data; },
  }));

  const hasErrors = $derived((result?.errors?.length ?? 0) > 0);
  const hasWarnings = $derived((result?.warnings?.length ?? 0) > 0);
  const plan = $derived(result?.plan);

  // Group warnings by message to collapse duplicates
  const groupedWarnings = $derived(
    (result?.warnings ?? []).reduce<{ message: string; stage: string; items: NonNullable<typeof result>['warnings'] }[]>((acc, w) => {
      const existing = acc.find((g) => g.message === w.message);
      if (existing) existing.items!.push(w);
      else acc.push({ message: w.message, stage: w.stage, items: [w] });
      return acc;
    }, []),
  );
  const totalFlowSteps = $derived(plan?.flow_plans?.reduce((n, f) => n + (f.steps?.length ?? 0), 0) ?? 0);
  const totalCases = $derived(plan?.suite_plans?.reduce((n, s) => n + (s.cases?.length ?? 0), 0) ?? 0);

  const startRun = createMutation(() => ({
    mutationFn: async () => {
      if (!plan) throw new Error('No plan compiled');
      const { data, error } = await api.POST('/api/projects/{project-id}/runs', {
        params: { path: { 'project-id': projectId } },
        body: { plan_id: plan.id, plan },
      });
      if (error) throw error;
      return data!;
    },
    onSuccess: (data) => { goto(`/runs/${data.id}`); },
  }));
</script>

<main class="mx-auto max-w-5xl px-6 py-8">
  <Breadcrumb crumbs={[{ label: projectId.slice(0, 8), href: `/projects/${projectId}` }, { label: 'Plan Preview' }]} />
  <div class="mb-6 flex flex-wrap items-end justify-between gap-4">
    <div>
      <p class="text-xs font-bold uppercase tracking-widest text-primary">Project {projectId.slice(0, 8)}</p>
      <h1 class="text-3xl font-bold">Plan Preview</h1>
      <p class="text-sm text-muted-foreground">Compile your project config and preview the generated smoke plan.</p>
    </div>
    <div class="flex gap-2">
      <Button variant="outline" href="/projects/{projectId}/environments">← Config</Button>
      <Button onclick={() => compile.mutate()} disabled={compile.isPending}>
        {#if compile.isPending}<Loader2 class="size-4 animate-spin" />{/if}
        {compile.isPending ? 'Compiling…' : '▶ Compile'}
      </Button>
      {#if plan && !hasErrors}
        <Button variant="default" onclick={() => startRun.mutate()} disabled={startRun.isPending}>
          {#if startRun.isPending}<Loader2 class="size-4 animate-spin" />{/if}
          {startRun.isPending ? 'Starting…' : '🚀 Run'}
        </Button>
      {/if}
    </div>
  </div>

  <!-- Config summary -->
  <Card.Root class="mb-6">
    <Card.Content class="flex flex-wrap gap-4 py-4">
      <div class="text-center"><p class="text-2xl font-bold">{$config.environments.length}</p><p class="text-xs text-muted-foreground">Environments</p></div>
      <div class="text-center"><p class="text-2xl font-bold">{$config.auth_profiles.length}</p><p class="text-xs text-muted-foreground">Auth Profiles</p></div>
      <div class="text-center"><p class="text-2xl font-bold">{$config.flows.length}</p><p class="text-xs text-muted-foreground">Flows</p></div>
      <div class="text-center"><p class="text-2xl font-bold">{$config.suites.length}</p><p class="text-xs text-muted-foreground">Suites</p></div>
    </Card.Content>
  </Card.Root>

  {#if compile.isError}
    <Card.Root class="mb-4 border-destructive">
      <Card.Content class="py-4 text-sm text-destructive">
        Compilation request failed: {compile.error?.message ?? 'Unknown error'}
      </Card.Content>
    </Card.Root>
  {/if}

  {#if result}
    <!-- Errors -->
    {#if hasErrors}
      <Card.Root class="mb-4 border-destructive">
        <Card.Header><Card.Title class="text-base text-destructive">Errors ({result.errors!.length})</Card.Title></Card.Header>
        <Card.Content class="space-y-2">
          {#each result.errors! as err}
            <div class="flex items-start gap-2 rounded-md bg-destructive/10 p-2 text-sm">
              <Badge variant="destructive" class="shrink-0 text-[0.6rem]">{err.stage}</Badge>
              <div>
                <p class="font-medium">{err.message}</p>
                <p class="text-xs text-muted-foreground">{err.path}{err.entity ? ` (${err.entity})` : ''}</p>
              </div>
            </div>
          {/each}
        </Card.Content>
      </Card.Root>
    {/if}

    <!-- Warnings (grouped) -->
    {#if hasWarnings}
      <Card.Root class="mb-4 border-yellow-600/50">
        <Card.Header><Card.Title class="text-base text-yellow-500">Warnings ({result.warnings!.length})</Card.Title></Card.Header>
        <Card.Content class="space-y-2">
          {#each groupedWarnings as group (group.message)}
            <div class="flex items-start gap-2 rounded-md bg-yellow-500/10 p-2 text-sm">
              <Badge variant="outline" class="shrink-0 text-[0.6rem] text-yellow-500">{group.stage}</Badge>
              <div class="min-w-0 flex-1">
                <p class="font-medium">{group.message}</p>
                {#if group.items!.length > 1}
                  <details class="mt-1">
                    <summary class="cursor-pointer text-xs text-muted-foreground">×{group.items!.length} occurrences</summary>
                    <ul class="mt-1 space-y-0.5 text-xs text-muted-foreground">
                      {#each group.items! as w}
                        <li>{w.path}{w.entity ? ` (${w.entity})` : ''}</li>
                      {/each}
                    </ul>
                  </details>
                {:else}
                  <p class="text-xs text-muted-foreground">{group.items![0].path}{group.items![0].entity ? ` (${group.items![0].entity})` : ''}</p>
                {/if}
              </div>
            </div>
          {/each}
        </Card.Content>
      </Card.Root>
    {/if}

    <!-- Plan -->
    {#if plan}
      <div class="mb-4 flex gap-4">
        <Badge variant="secondary">Environment: {plan.environment.name}</Badge>
        {#if plan.auth}<Badge variant="secondary">Auth: {plan.auth.name}</Badge>{/if}
        <Badge variant="secondary">{totalFlowSteps} flow steps</Badge>
        <Badge variant="secondary">{totalCases} suite cases</Badge>
      </div>

      <div class="grid grid-cols-1 gap-4 lg:grid-cols-2">
        <!-- Flow plans -->
        {#each plan.flow_plans ?? [] as fp (fp.flow_id)}
          <Card.Root>
            <Card.Header><Card.Title class="text-base">Flow: {fp.name}</Card.Title></Card.Header>
            <Card.Content class="space-y-1">
              {#each fp.steps ?? [] as step, i (i)}
                <div class="flex items-center gap-2 rounded-md bg-secondary/50 px-2 py-1 text-xs">
                  <Badge variant="outline" class="text-[0.6rem]">{i + 1}</Badge>
                  <span class="font-mono font-bold">{step.method}</span>
                  <span class="truncate">{step.path}</span>
                  <span class="ml-auto text-muted-foreground">{step.name}</span>
                </div>
              {/each}
              {#if fp.cleanup && fp.cleanup.length > 0}
                <Separator />
                <p class="text-xs text-muted-foreground">Cleanup:</p>
                {#each fp.cleanup as step, i (i)}
                  <div class="flex items-center gap-2 rounded-md bg-secondary/30 px-2 py-1 text-xs">
                    <Badge variant="outline" class="text-[0.6rem]">C{i + 1}</Badge>
                    <span class="font-mono font-bold">{step.method}</span>
                    <span class="truncate">{step.path}</span>
                  </div>
                {/each}
              {/if}
            </Card.Content>
          </Card.Root>
        {/each}

        <!-- Suite plans -->
        {#each plan.suite_plans ?? [] as sp (sp.suite_id)}
          <Card.Root>
            <Card.Header><Card.Title class="text-base">Suite: {sp.name}</Card.Title><Card.Description>{sp.cases?.length ?? 0} cases</Card.Description></Card.Header>
            <Card.Content class="space-y-1 max-h-64 overflow-y-auto">
              {#each sp.cases ?? [] as c, i (i)}
                <div class="flex items-center gap-2 rounded-md bg-secondary/50 px-2 py-1 text-xs">
                  <Badge variant="outline" class="text-[0.6rem]">{c.case_type}</Badge>
                  <span class="font-mono font-bold">{c.step.method}</span>
                  <span class="truncate">{c.step.path}</span>
                  <span class="ml-auto text-muted-foreground">{c.operation_id}</span>
                </div>
              {/each}
            </Card.Content>
          </Card.Root>
        {/each}
      </div>
    {:else if !hasErrors}
      <Card.Root>
        <Card.Content class="py-8 text-center text-sm text-muted-foreground">
          No plan generated. Add environments, flows, or suites to your config.
        </Card.Content>
      </Card.Root>
    {/if}
  {:else}
    <Card.Root>
      <Card.Content class="py-12 text-center text-sm text-muted-foreground">
        Click <strong>Compile</strong> to preview the generated smoke plan from your current config.
      </Card.Content>
    </Card.Root>
  {/if}
</main>
