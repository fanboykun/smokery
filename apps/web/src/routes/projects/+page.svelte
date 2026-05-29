<script lang="ts">
  import { createQuery, createMutation, useQueryClient } from '@tanstack/svelte-query';
  import { api } from '$lib/api/client';
  import * as Card from '$lib/components/ui/card';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';
  import { Input } from '$lib/components/ui/input';

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

<main class="mx-auto max-w-5xl px-6 py-8">
  <div class="mb-6 flex flex-wrap items-end justify-between gap-4">
    <div>
      <p class="text-xs font-bold uppercase tracking-widest text-primary">Spec-driven smoke testing</p>
      <h1 class="text-3xl font-bold">Projects</h1>
      <p class="text-sm text-muted-foreground">Manage OpenAPI specs, configure smoke tests, and track run health.</p>
    </div>
    <form
      class="flex gap-2"
      onsubmit={(e) => { e.preventDefault(); if (name.trim()) createProject.mutate(name.trim()); }}
    >
      <Input bind:value={name} placeholder="Project name" class="w-48" />
      <Button type="submit" disabled={createProject.isPending}>+ New</Button>
    </form>
  </div>

  {#if projects.isLoading}
    <p class="text-muted-foreground">Loading…</p>
  {:else if projects.isError}
    <Card.Root><Card.Content class="py-4 text-sm text-muted-foreground">API unavailable. Start the server with <code>make dev</code>.</Card.Content></Card.Root>
  {:else if projects.data && projects.data.length > 0}
    <div class="space-y-3">
      {#each projects.data as project (project.id)}
        <a href="/projects/{project.id}" class="block">
          <Card.Root class="transition-colors hover:border-primary/50">
            <Card.Header class="flex-row items-center justify-between pb-2">
              <div>
                <Card.Title class="text-lg">{project.name}</Card.Title>
                {#if project.description}
                  <Card.Description>{project.description}</Card.Description>
                {/if}
              </div>
              <Badge variant="secondary">{new Date(project.created_at).toLocaleDateString()}</Badge>
            </Card.Header>
            <Card.Footer class="gap-2 pt-0">
              <Button variant="outline" size="sm" href="/projects/{project.id}/operations">Operations</Button>
              <Button variant="outline" size="sm" href="/projects/{project.id}/environments">Config</Button>
              <Button variant="outline" size="sm" href="/projects/{project.id}/plan">Plan</Button>
              <Button size="sm" href="/projects/{project.id}/runs">Runs</Button>
            </Card.Footer>
          </Card.Root>
        </a>
      {/each}
    </div>
  {:else}
    <Card.Root>
      <Card.Content class="py-12 text-center text-sm text-muted-foreground">
        No projects yet. Create one above to get started.
      </Card.Content>
    </Card.Root>
  {/if}
</main>
