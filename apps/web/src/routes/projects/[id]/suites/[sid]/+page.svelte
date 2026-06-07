<script lang="ts">
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { createQuery } from '@tanstack/svelte-query';
  import { api } from '$lib/api/client';
  import { createProjectConfigStore, type Suite } from '$lib/stores/project-config';
  import * as Card from '$lib/components/ui/card';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';
  import { Input } from '$lib/components/ui/input';
  import { Label } from '$lib/components/ui/label';
  import { Separator } from '$lib/components/ui/separator';

  const projectId = $page.params.id!;
  const suiteId = $page.params.sid!;
  const config = createProjectConfigStore(projectId);
  const isNew = suiteId === 'new';

  // Fetch specs for the project
  const specs = createQuery(() => ({
    queryKey: ['specs', projectId],
    queryFn: async () => {
      const { data, error } = await api.GET('/api/projects/{project-id}/specs', {
        params: { path: { 'project-id': projectId } },
      });
      if (error) throw error;
      return data ?? [];
    },
  }));

  const latestSpecId = $derived(specs.data?.at(-1)?.id);

  // Fetch operations from spec
  const operations = createQuery(() => ({
    queryKey: ['operations', latestSpecId],
    queryFn: async () => {
      const { data, error } = await api.GET('/api/specs/{spec-id}/operations', {
        params: { path: { 'spec-id': latestSpecId! } },
      });
      if (error) throw error;
      return data ?? [];
    },
    enabled: !!latestSpecId,
  }));

  // Calculate matched operations based on selector
  const matchedOps = $derived.by(() => {
    const ops = operations.data ?? [];
    const matched: typeof ops = [];
    const excluded: typeof ops = [];

    for (const op of ops) {
      // Check if excluded
      if ((suite.selector.exclude ?? []).includes(op.operation_id)) {
        excluded.push(op);
        continue;
      }

      // Check if matches selector
      let matches = true;

      // Tags filter
      if ((suite.selector.tags ?? []).length > 0) {
        const opTags = op.tags ?? [];
        const hasTag = (suite.selector.tags ?? []).some((tag) => opTags.includes(tag));
        matches = matches && hasTag;
      }

      // Classifications filter
      if ((suite.selector.classifications ?? []).length > 0) {
        const hasClassification = (suite.selector.classifications ?? []).includes(op.classification);
        matches = matches && hasClassification;
      }

      // Path patterns filter
      if ((suite.selector.paths ?? []).length > 0) {
        const pathMatch = (suite.selector.paths ?? []).some((pattern) => {
          // Simple wildcard matching
          const regex = new RegExp(`^${pattern.replace(/\*/g, '.*')}$`);
          return regex.test(op.path);
        });
        matches = matches && pathMatch;
      }

      if (matches) {
        matched.push(op);
      }
    }

    return { matched, excluded };
  });

  const defaultStrategy = { default_list: true, pagination: true, search_from_response: false, enum_filters: false, empty_result_policy: 'allow', max_cases_per_op: 0 };

  let suite = $state<Suite>(
    isNew
      ? { id: crypto.randomUUID(), name: '', environment: '', selector: { tags: [], classifications: [], paths: [], exclude: [] }, strategy: { ...defaultStrategy } }
      : (() => {
          const found = $config.suites.find((s) => s.id === suiteId);
          if (!found) return { id: suiteId, name: '', environment: '', selector: { tags: [], classifications: [], paths: [], exclude: [] }, strategy: { ...defaultStrategy } };
          return {
            ...found,
            selector: {
              tags: [...(found.selector.tags ?? [])],
              classifications: [...(found.selector.classifications ?? [])],
              paths: [...(found.selector.paths ?? [])],
              exclude: [...(found.selector.exclude ?? [])],
            },
            strategy: { ...found.strategy },
          };
        })(),
  );

  // Selector input helpers
  let tagInput = $state('');
  let classInput = $state('');
  let pathInput = $state('');
  let excludeInput = $state('');

  function addTag() { if (!tagInput) return; suite.selector.tags = [...(suite.selector.tags ?? []), tagInput]; tagInput = ''; }
  function addClass() { if (!classInput) return; suite.selector.classifications = [...(suite.selector.classifications ?? []), classInput]; classInput = ''; }
  function addPath() { if (!pathInput) return; suite.selector.paths = [...(suite.selector.paths ?? []), pathInput]; pathInput = ''; }
  function addExclude() { if (!excludeInput) return; suite.selector.exclude = [...(suite.selector.exclude ?? []), excludeInput]; excludeInput = ''; }

  function removeFrom(arr: string[] | undefined, val: string): string[] {
    return (arr ?? []).filter((v) => v !== val);
  }

  function saveSuite() {
    if (!suite.name) return;
    const toSave = {
      ...suite,
      selector: {
        tags: [...(suite.selector.tags ?? [])],
        classifications: [...(suite.selector.classifications ?? [])],
        paths: [...(suite.selector.paths ?? [])],
        exclude: [...(suite.selector.exclude ?? [])],
      },
      strategy: { ...suite.strategy },
    };
    console.log('[smokery] saving suite:', toSave.id, 'paths:', toSave.selector.paths);
    config.update((c) => {
      const existing = c.suites.filter((s) => s.id !== toSave.id);
      return { ...c, suites: [...existing, toSave] };
    });
    goto(`/projects/${projectId}`);
  }

  function deleteSuite() {
    config.update((c) => ({ ...c, suites: c.suites.filter((s) => s.id !== suite.id) }));
    goto(`/projects/${projectId}`);
  }
</script>

<main class="mx-auto max-w-4xl px-6 py-8">
  <div class="mb-6 flex flex-wrap items-end justify-between gap-4">
    <div>
      <p class="text-xs font-bold uppercase tracking-widest text-primary">Suite Configurator</p>
      <h1 class="text-3xl font-bold">{suite.name || 'New Suite'}</h1>
    </div>
    <div class="flex gap-2">
      <Button variant="outline" href="/projects/{projectId}/builder">← Back</Button>
      {#if !isNew}
        <Button variant="destructive" onclick={deleteSuite}>Delete</Button>
      {/if}
      <Button onclick={saveSuite}>Save Suite</Button>
    </div>
  </div>

  <div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
    <!-- Basic info + Selector -->
    <div class="space-y-4">
      <Card.Root>
        <Card.Header><Card.Title class="text-base">Configuration</Card.Title></Card.Header>
        <Card.Content class="space-y-3">
          <div class="space-y-1"><Label>Name</Label><Input bind:value={suite.name} placeholder="List Endpoints Suite" /></div>
          <div class="space-y-1"><Label>Description</Label><Input bind:value={suite.description} placeholder="Optional" /></div>
          <div class="grid grid-cols-2 gap-3">
            <div class="space-y-1">
              <Label>Environment</Label>
              <select class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm" bind:value={suite.environment}>
                <option value="">Select…</option>
                {#each $config.environments as env (env.id)}
                  <option value={env.id}>{env.name}</option>
                {/each}
              </select>
            </div>
            <div class="space-y-1">
              <Label>Auth Profile</Label>
              <select class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm" bind:value={suite.auth}>
                <option value="">None</option>
                {#each $config.auth_profiles as auth (auth.id)}
                  <option value={auth.id}>{auth.name}</option>
                {/each}
              </select>
            </div>
          </div>
        </Card.Content>
      </Card.Root>

      <Card.Root>
        <Card.Header><Card.Title class="text-base">Operation Selector</Card.Title><Card.Description>Filter which operations are included in this suite.</Card.Description></Card.Header>
        <Card.Content class="space-y-4">
          <!-- Tags -->
          <div class="space-y-2">
            <Label>Tags</Label>
            <div class="flex flex-wrap gap-1">
              {#each suite.selector.tags ?? [] as tag (tag)}
                <Badge variant="secondary" class="gap-1">{tag}<button class="ml-1" onclick={() => (suite.selector.tags = removeFrom(suite.selector.tags, tag))}>×</button></Badge>
              {/each}
            </div>
            <div class="flex gap-2"><Input bind:value={tagInput} placeholder="e.g. users" class="flex-1" /><Button variant="outline" size="sm" onclick={addTag}>Add</Button></div>
          </div>
          <!-- Classifications -->
          <div class="space-y-2">
            <Label>Classifications</Label>
            <div class="flex flex-wrap gap-1">
              {#each suite.selector.classifications ?? [] as cls (cls)}
                <Badge variant="secondary" class="gap-1">{cls}<button class="ml-1" onclick={() => (suite.selector.classifications = removeFrom(suite.selector.classifications, cls))}>×</button></Badge>
              {/each}
            </div>
            <div class="flex gap-2"><Input bind:value={classInput} placeholder="e.g. list" class="flex-1" /><Button variant="outline" size="sm" onclick={addClass}>Add</Button></div>
          </div>
          <!-- Paths -->
          <div class="space-y-2">
            <Label>Path patterns</Label>
            <div class="flex flex-wrap gap-1">
              {#each suite.selector.paths ?? [] as p (p)}
                <Badge variant="secondary" class="gap-1">{p}<button class="ml-1" onclick={() => (suite.selector.paths = removeFrom(suite.selector.paths, p))}>×</button></Badge>
              {/each}
            </div>
            <div class="flex gap-2"><input bind:value={pathInput} placeholder="/api/functional-admin*" class="flex-1 rounded-md border border-input bg-background px-3 py-2 text-sm" /><Button variant="outline" size="sm" onclick={addPath}>Add</Button></div>
          </div>
          <!-- Exclude -->
          <div class="space-y-2">
            <Label>Exclude</Label>
            <div class="flex flex-wrap gap-1">
              {#each suite.selector.exclude ?? [] as ex (ex)}
                <Badge variant="secondary" class="gap-1">{ex}<button class="ml-1" onclick={() => (suite.selector.exclude = removeFrom(suite.selector.exclude, ex))}>×</button></Badge>
              {/each}
            </div>
            <div class="flex gap-2"><Input bind:value={excludeInput} placeholder="operationId to exclude" class="flex-1" /><Button variant="outline" size="sm" onclick={addExclude}>Add</Button></div>
          </div>
        </Card.Content>
      </Card.Root>
    </div>

    <!-- Right panel: Strategy + Preview -->
    <div class="space-y-4">
      <Card.Root>
        <Card.Header><Card.Title class="text-base">Generation Strategy</Card.Title><Card.Description>Control what test cases are generated for matched operations.</Card.Description></Card.Header>
        <Card.Content class="space-y-4">
          <label class="flex items-center gap-3 text-sm">
            <input type="checkbox" bind:checked={suite.strategy.default_list} class="size-4 rounded border-input" />
            <span><strong>Default list</strong> — basic GET request for each list endpoint</span>
          </label>
          <label class="flex items-center gap-3 text-sm">
            <input type="checkbox" bind:checked={suite.strategy.pagination} class="size-4 rounded border-input" />
            <span><strong>Pagination</strong> — verify pagination params work correctly</span>
          </label>
          <label class="flex items-center gap-3 text-sm">
            <input type="checkbox" bind:checked={suite.strategy.search_from_response} class="size-4 rounded border-input" />
            <span><strong>Search from response</strong> — capture a value, then search for it</span>
          </label>
          <label class="flex items-center gap-3 text-sm">
            <input type="checkbox" bind:checked={suite.strategy.enum_filters} class="size-4 rounded border-input" />
            <span><strong>Enum filters</strong> — test each enum value as a filter</span>
          </label>

          <Separator />

          <div class="space-y-1">
            <Label>Empty result policy</Label>
            <select class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm" bind:value={suite.strategy.empty_result_policy}>
              <option value="allow">allow — empty results are OK</option>
              <option value="warn">warn — flag empty results as warnings</option>
              <option value="fail">fail — empty results fail the case</option>
            </select>
          </div>

          <div class="space-y-1">
            <Label>Max cases per operation</Label>
            <Input type="number" bind:value={suite.strategy.max_cases_per_op} min="0" placeholder="0 = unlimited" />
          </div>
        </Card.Content>
      </Card.Root>

      <!-- Live Preview -->
      <Card.Root>
        <Card.Header><Card.Title class="text-base">Live Preview</Card.Title></Card.Header>
        <Card.Content class="space-y-3">
          {#if operations.isPending}
            <p class="text-xs text-muted-foreground">Loading operations…</p>
          {:else if operations.isError}
            <p class="text-xs text-destructive">Failed to load operations</p>
          {:else}
            <!-- Matched count -->
            <div class="rounded-md bg-emerald-500/10 border border-emerald-500/20 px-3 py-2">
              <p class="text-sm font-semibold text-emerald-400">{matchedOps.matched.length} operations matched</p>
              <p class="text-xs text-muted-foreground">Will generate test cases from these</p>
            </div>

            <!-- Matched operations list -->
            {#if matchedOps.matched.length > 0}
              <div class="space-y-2">
                <p class="text-xs font-semibold uppercase tracking-widest text-muted-foreground">Included</p>
                <div class="space-y-1 max-h-[200px] overflow-y-auto">
                  {#each matchedOps.matched as op (op.id)}
                    <div class="flex items-center gap-2 rounded-sm bg-muted/50 p-2 text-xs">
                      <Badge class="bg-blue-500/20 text-blue-300 font-mono text-xs shrink-0">
                        {op.method.toUpperCase()}
                      </Badge>
                      <span class="flex-1 truncate font-mono">{op.operation_id}</span>
                      {#if op.is_destructive}
                        <span class="text-red-400 text-xs shrink-0">destr</span>
                      {/if}
                    </div>
                  {/each}
                </div>
              </div>
            {/if}

            <!-- Excluded operations -->
            {#if matchedOps.excluded.length > 0}
              <div class="space-y-2">
                <p class="text-xs font-semibold uppercase tracking-widest text-muted-foreground">Excluded</p>
                <div class="space-y-1 max-h-[120px] overflow-y-auto">
                  {#each matchedOps.excluded as op (op.id)}
                    <div class="flex items-center gap-2 rounded-sm bg-red-500/5 p-2 text-xs">
                      <Badge class="bg-red-500/20 text-red-300 font-mono text-xs shrink-0">
                        {op.method.toUpperCase()}
                      </Badge>
                      <span class="flex-1 truncate font-mono">{op.operation_id}</span>
                      <span class="text-red-400 text-xs shrink-0">excluded</span>
                    </div>
                  {/each}
                </div>
              </div>
            {/if}

            <!-- Empty state -->
            {#if matchedOps.matched.length === 0 && matchedOps.excluded.length === 0}
              <p class="text-xs text-muted-foreground">No operations match your selector. Adjust filters above.</p>
            {/if}
          {/if}
        </Card.Content>
      </Card.Root>
    </div>
  </div>
</main>
