<script lang="ts">
  import { page } from '$app/stores';
  import { createProjectConfigStore, type Environment, type AuthProfile } from '$lib/stores/project-config';
  import * as Card from '$lib/components/ui/card';
  import * as Tabs from '$lib/components/ui/tabs';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';
  import { Input } from '$lib/components/ui/input';
  import { Label } from '$lib/components/ui/label';
  import { Separator } from '$lib/components/ui/separator';

  const config = createProjectConfigStore($page.params.id!);

  // --- Environments ---
  let editingEnv = $state<Environment | null>(null);
  let headerKey = $state('');
  let headerVal = $state('');

  function newEnv() {
    editingEnv = { id: crypto.randomUUID(), name: '', base_url: '', headers: {} };
  }
  function editEnv(env: Environment) {
    editingEnv = { ...env, headers: { ...(env.headers ?? {}) } };
  }
  function saveEnv() {
    if (!editingEnv || !editingEnv.name || !editingEnv.base_url) return;
    config.update((c) => {
      const idx = c.environments.findIndex((e) => e.id === editingEnv!.id);
      if (idx >= 0) c.environments[idx] = editingEnv!;
      else c.environments = [...c.environments, editingEnv!];
      return c;
    });
    editingEnv = null;
  }
  function removeEnv(id: string) {
    config.update((c) => ({ ...c, environments: c.environments.filter((e) => e.id !== id) }));
    if (editingEnv?.id === id) editingEnv = null;
  }
  function addHeader() {
    if (!editingEnv || !headerKey) return;
    editingEnv.headers = { ...(editingEnv.headers ?? {}), [headerKey]: headerVal };
    headerKey = '';
    headerVal = '';
  }
  function removeHeader(key: string) {
    if (!editingEnv) return;
    const { [key]: _, ...rest } = editingEnv.headers ?? {};
    editingEnv.headers = rest;
  }

  // --- Auth Profiles ---
  let editingAuth = $state<AuthProfile | null>(null);
  let configKey = $state('');
  let configVal = $state('');

  function newAuth() {
    editingAuth = { id: crypto.randomUUID(), name: '', type: 'bearer', config: {} };
  }
  function editAuth(auth: AuthProfile) {
    editingAuth = { ...auth, config: { ...auth.config } };
  }
  function saveAuth() {
    if (!editingAuth || !editingAuth.name) return;
    config.update((c) => {
      const idx = c.auth_profiles.findIndex((a) => a.id === editingAuth!.id);
      if (idx >= 0) c.auth_profiles[idx] = editingAuth!;
      else c.auth_profiles = [...c.auth_profiles, editingAuth!];
      return c;
    });
    editingAuth = null;
  }
  function removeAuth(id: string) {
    config.update((c) => ({ ...c, auth_profiles: c.auth_profiles.filter((a) => a.id !== id) }));
    if (editingAuth?.id === id) editingAuth = null;
  }
  function addConfigEntry() {
    if (!editingAuth || !configKey) return;
    editingAuth.config = { ...editingAuth.config, [configKey]: configVal };
    configKey = '';
    configVal = '';
  }
  function removeConfigEntry(key: string) {
    if (!editingAuth) return;
    const { [key]: _, ...rest } = editingAuth.config;
    editingAuth.config = rest;
  }
</script>

<main class="mx-auto max-w-5xl px-6 py-8">
  <div class="mb-6">
    <p class="text-xs font-bold uppercase tracking-widest text-primary">Project {$page.params.id?.slice(0, 8)}</p>
    <h1 class="text-3xl font-bold">Environments & Auth</h1>
    <p class="text-sm text-muted-foreground">Configure target environments and authentication profiles.</p>
  </div>

  <Tabs.Root value="environments">
    <Tabs.List>
      <Tabs.Trigger value="environments">Environments ({$config.environments.length})</Tabs.Trigger>
      <Tabs.Trigger value="auth">Auth Profiles ({$config.auth_profiles.length})</Tabs.Trigger>
    </Tabs.List>

    <!-- ENVIRONMENTS TAB -->
    <Tabs.Content value="environments">
      <div class="mt-4 flex justify-end">
        <Button onclick={newEnv}>+ New Environment</Button>
      </div>
      <div class="mt-4 grid grid-cols-1 gap-6 lg:grid-cols-[1fr_1fr]">
        <div class="space-y-3">
          {#if $config.environments.length === 0}
            <Card.Root><Card.Content class="py-8 text-center text-sm text-muted-foreground">No environments yet.</Card.Content></Card.Root>
          {/if}
          {#each $config.environments as env (env.id)}
            <Card.Root class="{editingEnv?.id === env.id ? 'ring-1 ring-primary' : ''}">
              <Card.Header class="flex-row items-center justify-between pb-2">
                <div>
                  <Card.Title class="text-base">{env.name}</Card.Title>
                  <Card.Description class="font-mono text-xs">{env.base_url}</Card.Description>
                </div>
                {#if Object.keys(env.headers ?? {}).length > 0}
                  <Badge variant="secondary">{Object.keys(env.headers!).length} headers</Badge>
                {/if}
              </Card.Header>
              <Card.Footer class="gap-2 pt-0">
                <Button variant="outline" size="sm" onclick={() => editEnv(env)}>Edit</Button>
                <Button variant="destructive" size="sm" onclick={() => removeEnv(env.id)}>Remove</Button>
              </Card.Footer>
            </Card.Root>
          {/each}
        </div>
        <div>
          {#if editingEnv}
            <Card.Root>
              <Card.Header><Card.Title class="text-base">{editingEnv.name || 'New Environment'}</Card.Title></Card.Header>
              <Card.Content class="space-y-4">
                <div class="space-y-1"><Label>Name</Label><Input bind:value={editingEnv.name} placeholder="staging" /></div>
                <div class="space-y-1"><Label>Base URL</Label><Input bind:value={editingEnv.base_url} placeholder="https://api.staging.example.com" /></div>
                <Separator />
                <div class="space-y-2">
                  <Label>Headers</Label>
                  {#each Object.entries(editingEnv.headers ?? {}) as [k, v] (k)}
                    <div class="flex items-center gap-2">
                      <Badge variant="outline" class="font-mono text-xs">{k}: {v}</Badge>
                      <Button variant="ghost" size="xs" onclick={() => removeHeader(k)}>×</Button>
                    </div>
                  {/each}
                  <div class="flex gap-2">
                    <Input bind:value={headerKey} placeholder="Header" class="flex-1" />
                    <Input bind:value={headerVal} placeholder="Value" class="flex-1" />
                    <Button variant="outline" size="sm" onclick={addHeader}>Add</Button>
                  </div>
                </div>
                <Separator />
                <div class="flex gap-2">
                  <Button class="flex-1" onclick={saveEnv}>Save</Button>
                  <Button variant="outline" class="flex-1" onclick={() => (editingEnv = null)}>Cancel</Button>
                </div>
              </Card.Content>
            </Card.Root>
          {:else}
            <Card.Root><Card.Content class="py-8 text-center text-sm text-muted-foreground">Select or create an environment.</Card.Content></Card.Root>
          {/if}
        </div>
      </div>
    </Tabs.Content>

    <!-- AUTH PROFILES TAB -->
    <Tabs.Content value="auth">
      <div class="mt-4 flex justify-end">
        <Button onclick={newAuth}>+ New Auth Profile</Button>
      </div>
      <div class="mt-4 grid grid-cols-1 gap-6 lg:grid-cols-[1fr_1fr]">
        <div class="space-y-3">
          {#if $config.auth_profiles.length === 0}
            <Card.Root><Card.Content class="py-8 text-center text-sm text-muted-foreground">No auth profiles yet.</Card.Content></Card.Root>
          {/if}
          {#each $config.auth_profiles as auth (auth.id)}
            <Card.Root class="{editingAuth?.id === auth.id ? 'ring-1 ring-primary' : ''}">
              <Card.Header class="flex-row items-center justify-between pb-2">
                <div>
                  <Card.Title class="text-base">{auth.name}</Card.Title>
                  <Card.Description class="text-xs">{auth.type}</Card.Description>
                </div>
                <Badge variant="secondary">{auth.type}</Badge>
              </Card.Header>
              <Card.Footer class="gap-2 pt-0">
                <Button variant="outline" size="sm" onclick={() => editAuth(auth)}>Edit</Button>
                <Button variant="destructive" size="sm" onclick={() => removeAuth(auth.id)}>Remove</Button>
              </Card.Footer>
            </Card.Root>
          {/each}
        </div>
        <div>
          {#if editingAuth}
            <Card.Root>
              <Card.Header><Card.Title class="text-base">{editingAuth.name || 'New Auth Profile'}</Card.Title></Card.Header>
              <Card.Content class="space-y-4">
                <div class="space-y-1"><Label>Name</Label><Input bind:value={editingAuth.name} placeholder="API Key Auth" /></div>
                <div class="space-y-1">
                  <Label>Type</Label>
                  <select class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm" bind:value={editingAuth.type}>
                    <option value="bearer">Bearer Token</option>
                    <option value="basic">Basic Auth</option>
                    <option value="api_key">API Key</option>
                    <option value="custom">Custom</option>
                  </select>
                </div>
                <Separator />
                <div class="space-y-2">
                  <Label>Config</Label>
                  <p class="text-xs text-muted-foreground">
                    {#if editingAuth.type === 'bearer'}Key: <code>token</code>
                    {:else if editingAuth.type === 'basic'}Keys: <code>username</code>, <code>password</code>
                    {:else if editingAuth.type === 'api_key'}Keys: <code>header</code>, <code>value</code>
                    {:else}Add custom key-value pairs.
                    {/if}
                  </p>
                  {#each Object.entries(editingAuth.config) as [k, v] (k)}
                    <div class="flex items-center gap-2">
                      <Badge variant="outline" class="font-mono text-xs">{k}: {String(v)}</Badge>
                      <Button variant="ghost" size="xs" onclick={() => removeConfigEntry(k)}>×</Button>
                    </div>
                  {/each}
                  <div class="flex gap-2">
                    <Input bind:value={configKey} placeholder="Key" class="flex-1" />
                    <Input bind:value={configVal} placeholder="Value" class="flex-1" />
                    <Button variant="outline" size="sm" onclick={addConfigEntry}>Add</Button>
                  </div>
                </div>
                <Separator />
                <div class="flex gap-2">
                  <Button class="flex-1" onclick={saveAuth}>Save</Button>
                  <Button variant="outline" class="flex-1" onclick={() => (editingAuth = null)}>Cancel</Button>
                </div>
              </Card.Content>
            </Card.Root>
          {:else}
            <Card.Root><Card.Content class="py-8 text-center text-sm text-muted-foreground">Select or create an auth profile.</Card.Content></Card.Root>
          {/if}
        </div>
      </div>
    </Tabs.Content>
  </Tabs.Root>
</main>
