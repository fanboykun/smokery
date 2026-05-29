<script lang="ts">
  import { page } from '$app/stores';
  import { createQuery } from '@tanstack/svelte-query';
  import { api } from '$lib/api/client';
  import { createProjectConfigStore } from '$lib/stores/project-config';
  import * as Card from '$lib/components/ui/card';
  import * as Tabs from '$lib/components/ui/tabs';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';
  import { Input } from '$lib/components/ui/input';
  import { Label } from '$lib/components/ui/label';
  import { Textarea } from '$lib/components/ui/textarea';

  const projectId = $page.params.id!;
  const config = createProjectConfigStore(projectId);

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
    { href: `/projects/${projectId}/environments`, label: 'Environments & Auth', desc: 'Configure targets' },
    { href: `/projects/${projectId}/flows/new`, label: 'New Flow', desc: 'Build a scenario' },
    { href: `/projects/${projectId}/suites/new`, label: 'New Suite', desc: 'Configure generated tests' },
    { href: `/projects/${projectId}/plan`, label: 'Plan Preview', desc: 'Compile & inspect' },
    { href: `/projects/${projectId}/runs`, label: 'Runs', desc: 'View run history' },
  ];
</script>

<main class="mx-auto max-w-4xl px-6 py-8">
  <div class="mb-6">
    <p class="text-xs font-bold uppercase tracking-widest text-primary">Project</p>
    <h1 class="text-3xl font-bold">{project.data?.name ?? 'Loading…'}</h1>
    {#if project.data?.description}
      <p class="text-sm text-muted-foreground">{project.data.description}</p>
    {/if}
  </div>

  <!-- Navigation grid -->
  <div class="mb-8 grid grid-cols-2 gap-3 sm:grid-cols-3">
    {#each navItems as item (item.href)}
      <a href={item.href} class="rounded-lg border border-border p-4 transition-colors hover:bg-secondary">
        <p class="font-medium">{item.label}</p>
        <p class="text-xs text-muted-foreground">{item.desc}</p>
      </a>
    {/each}
  </div>

  <!-- Config summary -->
  {#if $config.flows.length > 0 || $config.suites.length > 0}
    <Card.Root class="mb-6">
      <Card.Header><Card.Title class="text-base">Current Config</Card.Title></Card.Header>
      <Card.Content class="space-y-2">
        {#each $config.flows as flow (flow.id)}
          <a href="/projects/{projectId}/flows/{flow.id}" class="flex items-center justify-between rounded-md bg-secondary/50 px-3 py-2 text-sm hover:bg-secondary">
            <span>Flow: <strong>{flow.name}</strong></span>
            <Badge variant="secondary">{flow.steps.length} steps</Badge>
          </a>
        {/each}
        {#each $config.suites as suite (suite.id)}
          <a href="/projects/{projectId}/suites/{suite.id}" class="flex items-center justify-between rounded-md bg-secondary/50 px-3 py-2 text-sm hover:bg-secondary">
            <span>Suite: <strong>{suite.name}</strong></span>
            <Badge variant="secondary">{suite.strategy.default_list ? 'auto' : 'manual'}</Badge>
          </a>
        {/each}
      </Card.Content>
    </Card.Root>
  {/if}

  <!-- Spec import -->
  <Card.Root>
    <Card.Header><Card.Title class="text-base">Import OpenAPI Spec</Card.Title></Card.Header>
    <Card.Content>
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
          <Textarea bind:value={specInput} class="min-h-[10rem] font-mono text-xs" placeholder="Paste OpenAPI JSON or YAML here…" />
          <Button onclick={handlePasteImport} disabled={!specInput.trim()}>Import</Button>
        </Tabs.Content>
      </Tabs.Root>

      {#if importSuccess}
        <p class="mt-3 text-sm text-primary">{importSuccess}</p>
      {/if}
      {#if importError}
        <p class="mt-3 text-sm text-destructive">{importError}</p>
      {/if}
    </Card.Content>
  </Card.Root>
</main>
