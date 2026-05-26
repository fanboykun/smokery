<script lang="ts">
  import { createQuery, createMutation, useQueryClient } from '@tanstack/svelte-query';
  import { api } from '$lib/api/client';

  const queryClient = useQueryClient();
  let name = $state('');

  const projects = createQuery(() => ({
    queryKey: ['projects'],
    queryFn: async () => {
      const { data, error } = await api.GET('/api/projects');
      if (error) throw error;
      return data!;
    },
  }));

  const createProject = createMutation(() => ({
    mutationFn: async (projectName: string) => {
      const { data, error } = await api.POST('/api/projects', {
        body: { name: projectName, description: '' },
      });
      if (error) throw error;
      return data!;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['projects'] });
      name = '';
    },
  }));
</script>

<h1>Projects</h1>
<form onsubmit={(e) => { e.preventDefault(); createProject.mutate(name); }}>
  <input bind:value={name} placeholder="Project name" required />
  <button type="submit" disabled={createProject.isPending}>Create</button>
</form>

{#if projects.isLoading}
  <p>Loading...</p>
{:else if projects.data}
  <ul>
    {#each projects.data as p (p.id)}
      <li><a href="/projects/{p.id}">{p.name}</a></li>
    {/each}
  </ul>
{/if}
