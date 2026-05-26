<script lang="ts">
  import { page } from '$app/stores';
  import { createQuery, createMutation, useQueryClient } from '@tanstack/svelte-query';
  import { api } from '$lib/api/client';

  const queryClient = useQueryClient();

  const runs = createQuery(() => ({
    queryKey: ['runs', $page.params.id],
    queryFn: async () => {
      const { data, error } = await api.GET('/api/projects/{project-id}/runs', {
        params: { path: { 'project-id': $page.params.id! } },
      });
      if (error) throw error;
      return data!;
    },
  }));

  const startRun = createMutation(() => ({
    mutationFn: async () => {
      const id = $page.params.id!;
      const { data, error } = await api.POST('/api/projects/{project-id}/runs', {
        params: { path: { 'project-id': id } },
        body: {
          plan_id: crypto.randomUUID(),
          plan: {
            id: crypto.randomUUID(),
            project_id: id,
            environment: { id: 'default', name: 'dev', base_url: 'http://localhost:8080' },
            flow_plans: [],
            suite_plans: [],
            compiled_at: new Date().toISOString(),
          },
        },
      });
      if (error) throw error;
      return data!;
    },
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['runs', $page.params.id] }),
  }));
</script>

<h1>Runs</h1>
<button onclick={() => startRun.mutate()} disabled={startRun.isPending}>Start Run</button>

{#if runs.isLoading}
  <p>Loading...</p>
{:else if runs.data}
  <ul>
    {#each runs.data as r (r.id)}
      <li><a href="/runs/{r.id}">{r.status} - {r.created_at}</a></li>
    {/each}
  </ul>
{/if}
