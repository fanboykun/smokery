<script lang="ts">
  import { page } from '$app/stores';
  import { createQuery } from '@tanstack/svelte-query';
  import { api } from '$lib/api/client';
  import { createProjectConfigStore } from '$lib/stores/project-config';
  import * as Card from '$lib/components/ui/card';
  import * as Tabs from '$lib/components/ui/tabs';
  import * as Dialog from '$lib/components/ui/dialog';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';
  import { Input } from '$lib/components/ui/input';
  import { Label } from '$lib/components/ui/label';
  import { Textarea } from '$lib/components/ui/textarea';

  const projectId = $page.params.id!;
  const config = createProjectConfigStore(projectId);

  let importOpen = $state(false);
  let specInput = $state('');
  let specUrl = $state('');
  let headerKey = $state('');
  let headerVal = $state('');
  let urlHeaders = $state<Record<string, string>>({});
  let importSuccess = $state('');
  let importError = $state('');

  const project = createQuery(() => ({
    queryKey: ['projects', projectId],
    queryFn: async () => {
      const { data, error } = await api.GET('/api/projects/{id}', { params: { path: { id: projectId } } });
      if (error) throw error;
      return data!;
    },
  }));

  async function handlePasteImport() {
    importError = '';
    importSuccess = '';
    const { data, error } = await api.POST('/api/projects/{project-id}/specs', {
      params: { path: { 'project-id': projectId } },
      body: specInput,
      bodySerializer: (body) => body as unknown as string,
    });
    if (error) { importError = (error as any).detail ?? 'Import failed'; return; }
    if (data) {
      importSuccess = `Imported: ${data.title} v${data.version} (${data.operations?.length ?? 0} operations)`;
      specInput = '';
    }
  }

  async function handleUrlImport() {
    importError = '';
    importSuccess = '';
    const { data, error } = await api.POST('/api/projects/{project-id}/specs/from-url', {
      params: { path: { 'project-id': projectId } },
      body: { url: specUrl, headers: Object.keys(urlHeaders).length > 0 ? urlHeaders : undefined },
    });
    if (error) { importError = (error as any).detail ?? 'Import failed'; return; }
    if (data) {
      importSuccess = `Imported: ${data.title} v${data.version} (${data.operations?.length ?? 0} operations)`;
      specUrl = '';
      urlHeaders = {};
    }
  }

  function addUrlHeader() {
    if (!headerKey) return;
    urlHeaders = { ...urlHeaders, [headerKey]: headerVal };
    headerKey = '';
    headerVal = '';
  }

  function removeUrlHeader(key: string) {
    const { [key]: _, ...rest } = urlHeaders;
    urlHeaders = rest;
  }

  const navItems = [
    { href: `/projects/${projectId}/operations`, label: 'Operations', desc: 'Explore & classify' },
    { href: `/projects/${projectId}/environments`, label: 'Environments', desc: 'Configure targets' },
    { href: `/projects/${projectId}/flows/new`, label: 'New Flow', desc: 'Build a scenario' },
    { href: `/projects/${projectId}/suites/new`, label: 'New Suite', desc: 'Generated tests' },
    { href: `/projects/${projectId}/plan`, label: 'Plan Preview', desc: 'Compile & run' },
    { href: `/projects/${projectId}/runs`, label: 'Runs', desc: 'History' },
  ];
</script>

<main class="mx-auto max-w-4xl px-6 py-8">
  <div class="mb-6 flex flex-wrap items-end justify-between gap-3">
    <div>
      <p class="text-xs font-bold uppercase tracking-widest text-primary">Project</p>
      <h1 class="text-3xl font-bold">{project.data?.name ?? 'Loading…'}</h1>
      {#if project.data?.description}
        <p class="text-sm text-muted-foreground">{project.data.description}</p>
      {/if}
    </div>
    <Button onclick={() => { importOpen = true; importSuccess = ''; importError = ''; }}>Import Spec</Button>
  </div>

  <!-- Navigation grid -->
  <div class="mb-6 grid grid-cols-2 gap-3 sm:grid-cols-3">
    {#each navItems as item (item.href)}
      <a href={item.href} class="rounded-lg border border-border p-3 transition-colors hover:bg-secondary">
        <p class="text-sm font-medium">{item.label}</p>
        <p class="text-xs text-muted-foreground">{item.desc}</p>
      </a>
    {/each}
  </div>

  <!-- Compact config summary -->
  {#if $config.flows.length > 0 || $config.suites.length > 0 || $config.environments.length > 0}
    <Card.Root>
      <Card.Header class="pb-2"><Card.Title class="text-base">Config</Card.Title></Card.Header>
      <Card.Content class="flex flex-wrap gap-4 text-sm">
        <span><strong>{$config.environments.length}</strong> env{$config.environments.length !== 1 ? 's' : ''}</span>
        <span><strong>{$config.auth_profiles.length}</strong> auth</span>
        <span><strong>{$config.flows.length}</strong> flow{$config.flows.length !== 1 ? 's' : ''}</span>
        <span><strong>{$config.suites.length}</strong> suite{$config.suites.length !== 1 ? 's' : ''}</span>
      </Card.Content>
      {#if $config.flows.length > 0 || $config.suites.length > 0}
        <Card.Content class="space-y-1 border-t pt-3">
          {#each $config.flows as flow (flow.id)}
            <a href="/projects/{projectId}/flows/{flow.id}" class="flex items-center justify-between rounded px-2 py-1 text-sm hover:bg-secondary">
              <span>Flow: {flow.name}</span>
              <Badge variant="secondary" class="text-[0.65rem]">{flow.steps.length} steps</Badge>
            </a>
          {/each}
          {#each $config.suites as suite (suite.id)}
            <a href="/projects/{projectId}/suites/{suite.id}" class="flex items-center justify-between rounded px-2 py-1 text-sm hover:bg-secondary">
              <span>Suite: {suite.name}</span>
              <Badge variant="secondary" class="text-[0.65rem]">{suite.strategy.default_list ? 'auto' : 'manual'}</Badge>
            </a>
          {/each}
        </Card.Content>
      {/if}
    </Card.Root>
  {/if}
</main>

<!-- Import Spec Dialog -->
<Dialog.Root bind:open={importOpen}>
  <Dialog.Content class="max-w-lg">
    <Dialog.Header>
      <Dialog.Title>Import OpenAPI Spec</Dialog.Title>
      <Dialog.Description>Import from a URL or paste the spec content directly.</Dialog.Description>
    </Dialog.Header>

    <Tabs.Root value="url">
      <Tabs.List>
        <Tabs.Trigger value="url">From URL</Tabs.Trigger>
        <Tabs.Trigger value="paste">Paste</Tabs.Trigger>
      </Tabs.List>

      <Tabs.Content value="url" class="space-y-3 pt-3">
        <div class="space-y-1">
          <Label>Spec URL</Label>
          <Input bind:value={specUrl} placeholder="https://api.example.com/openapi.json" />
        </div>
        <div class="space-y-2">
          <Label>Custom Headers (optional)</Label>
          {#each Object.entries(urlHeaders) as [k, v] (k)}
            <div class="flex items-center gap-2">
              <Badge variant="outline" class="font-mono text-xs">{k}: {v}</Badge>
              <Button variant="ghost" size="xs" onclick={() => removeUrlHeader(k)}>×</Button>
            </div>
          {/each}
          <div class="flex gap-2">
            <Input bind:value={headerKey} placeholder="Header name" class="flex-1" />
            <Input bind:value={headerVal} placeholder="Value" class="flex-1" />
            <Button variant="outline" size="sm" onclick={addUrlHeader}>Add</Button>
          </div>
        </div>
        <Button onclick={handleUrlImport} disabled={!specUrl.trim()}>Import from URL</Button>
      </Tabs.Content>

      <Tabs.Content value="paste" class="space-y-3 pt-3">
        <Textarea bind:value={specInput} class="min-h-[8rem] font-mono text-xs" placeholder="Paste OpenAPI JSON or YAML here…" />
        <Button onclick={handlePasteImport} disabled={!specInput.trim()}>Import</Button>
      </Tabs.Content>
    </Tabs.Root>

    {#if importSuccess}
      <p class="mt-3 text-sm text-primary">{importSuccess}</p>
    {/if}
    {#if importError}
      <p class="mt-3 text-sm text-destructive">{importError}</p>
    {/if}
  </Dialog.Content>
</Dialog.Root>
