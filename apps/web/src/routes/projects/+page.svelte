<script lang="ts">
  import { createQuery, createMutation, useQueryClient } from '@tanstack/svelte-query';
  import { api } from '$lib/api/client';
  import * as Card from '$lib/components/ui/card';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';
  import { Input } from '$lib/components/ui/input';
  import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
  import * as AlertDialog from '$lib/components/ui/alert-dialog';
  import { EllipsisVertical, TrendingUp, Play, Settings } from '@lucide/svelte';
  import PageLayout from '$lib/components/PageLayout.svelte';
  import StatsCard from '$lib/components/StatsCard.svelte';
  import StatusBadge from '$lib/components/StatusBadge.svelte';

  const queryClient = useQueryClient();
  let name = $state('');
  let deleteTarget = $state<{ id: string; name: string } | null>(null);
  let confirmText = $state('');
  let alertOpen = $state(false);

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

  async function deleteProject() {
    if (!deleteTarget || confirmText !== deleteTarget.name) return;
    await api.DELETE('/api/projects/{id}', { params: { path: { id: deleteTarget.id } } });
    queryClient.invalidateQueries({ queryKey: ['projects'] });
    alertOpen = false;
    deleteTarget = null;
    confirmText = '';
  }
</script>

<main class="min-h-screen space-y-8 bg-background px-6 py-8">
  {#snippet actionsSnippet()}
    <form class="flex gap-2" onsubmit={(e) => { e.preventDefault(); if (name.trim()) createProject.mutate(name.trim()); }}>
      <Input bind:value={name} placeholder="Project name" class="w-48" />
      <Button type="submit" disabled={createProject.isPending}>+ New Project</Button>
    </form>
  {/snippet}

  <PageLayout 
    title="Projects"
    subtitle="Spec-driven smoke testing"
    description="Create and manage OpenAPI smoke tests. Track health across environments."
    class="mx-auto max-w-7xl"
    actions={actionsSnippet}
  >

    {#if projects.isLoading}
      <div class="space-y-4">
        {#each Array(3) as _}
          <Card.Root class="animate-pulse"><Card.Content class="h-24" /></Card.Root>
        {/each}
      </div>
    {:else if projects.isError}
      <Card.Root class="border-destructive/50 bg-destructive/5">
        <Card.Content class="py-6 text-sm text-muted-foreground">
          <p class="font-semibold text-destructive">API unavailable</p>
          <p class="mt-1">Start the server with <code class="rounded bg-muted px-1.5 py-0.5 text-xs">make dev</code>.</p>
        </Card.Content>
      </Card.Root>
    {:else if projects.data && projects.data.length > 0}
      <div class="grid gap-4 lg:grid-cols-2">
        {#each projects.data as project (project.id)}
          <Card.Root class="group relative overflow-hidden transition-all hover:border-primary/50 hover:shadow-lg">
            <Card.Header class="pb-3">
              <div class="flex items-start justify-between gap-4">
                <div class="flex-1">
                  <a href="/projects/{project.id}" class="hover:underline">
                    <Card.Title class="text-lg">{project.name}</Card.Title>
                  </a>
                  {#if project.description}
                    <Card.Description class="mt-1">{project.description}</Card.Description>
                  {/if}
                </div>
                <div onclick={(e) => e.preventDefault()} role="none" class="opacity-0 transition-opacity group-hover:opacity-100">
                  <DropdownMenu.Root>
                    <DropdownMenu.Trigger>
                      {#snippet child({ props })}
                        <button {...props} class="inline-flex size-7 items-center justify-center rounded-md hover:bg-muted" aria-label="More actions">
                          <EllipsisVertical class="size-4" />
                        </button>
                      {/snippet}
                    </DropdownMenu.Trigger>
                    <DropdownMenu.Content align="end">
                      <DropdownMenu.Item
                        class="text-destructive"
                        onclick={() => { deleteTarget = { id: project.id, name: project.name }; confirmText = ''; alertOpen = true; }}
                      >
                        Delete project
                      </DropdownMenu.Item>
                    </DropdownMenu.Content>
                  </DropdownMenu.Root>
                </div>
              </div>
            </Card.Header>
            <Card.Content class="space-y-3 pb-4">
              <div class="flex gap-2">
                <Badge variant="outline" class="text-xs">{new Date(project.created_at).toLocaleDateString()}</Badge>
              </div>
            </Card.Content>
            <Card.Footer class="gap-2 border-t bg-muted/20 p-3">
              <Button variant="outline" size="sm" href="/projects/{project.id}/builder" class="flex-1">
                <Play class="size-3" />
                Builder
              </Button>
              <Button variant="ghost" size="sm" href="/projects/{project.id}/runs" class="flex-1">
                <TrendingUp class="size-3" />
                Runs
              </Button>
              <Button variant="ghost" size="sm" href="/projects/{project.id}/environments" class="flex-1">
                <Settings class="size-3" />
                Config
              </Button>
            </Card.Footer>
          </Card.Root>
        {/each}
      </div>
    {:else}
      <Card.Root class="border-dashed">
        <Card.Content class="py-12 text-center">
          <p class="text-sm font-medium text-foreground">No projects yet</p>
          <p class="mt-1 text-sm text-muted-foreground">Create your first project above to get started with smoke testing.</p>
        </Card.Content>
      </Card.Root>
    {/if}
  </PageLayout>
</main>

<!-- Delete confirmation -->
<AlertDialog.Root bind:open={alertOpen}>
  <AlertDialog.Content>
    <AlertDialog.Header>
      <AlertDialog.Title>Delete project</AlertDialog.Title>
      <AlertDialog.Description>
        This action cannot be undone. Type <strong class="text-foreground">{deleteTarget?.name}</strong> to confirm.
      </AlertDialog.Description>
    </AlertDialog.Header>
    <Input bind:value={confirmText} placeholder={deleteTarget?.name ?? ''} class="mt-2" />
    <AlertDialog.Footer>
      <AlertDialog.Cancel onclick={() => { alertOpen = false; deleteTarget = null; confirmText = ''; }}>Cancel</AlertDialog.Cancel>
      <AlertDialog.Action disabled={confirmText !== deleteTarget?.name} onclick={deleteProject}>Delete</AlertDialog.Action>
    </AlertDialog.Footer>
  </AlertDialog.Content>
</AlertDialog.Root>
