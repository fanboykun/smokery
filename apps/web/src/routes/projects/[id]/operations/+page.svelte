<script lang="ts">
  import { page } from '$app/stores';
  import { createQuery } from '@tanstack/svelte-query';
  import { api } from '$lib/api/client';
  import { demoOperations, type OperationSummary } from '$lib/demo-data';

  let filter = $state('');
  let selectedId = $state(demoOperations[0]?.id ?? '');
  let overrideJson = $state('{\n  "custom_headers": {\n    "X-Test": "true"\n  }\n}');

  const specs = createQuery(() => ({
    queryKey: ['project-specs', $page.params.id],
    queryFn: async () => {
      const { data, error } = await api.GET('/api/projects/{id}/specs', {
        params: { path: { id: $page.params.id! } },
      });
      if (error) throw error;
      return data!;
    },
  }));

  const operations = $derived<OperationSummary[]>(demoOperations);
  const filtered = $derived(
    operations.filter((operation) => {
      const needle = filter.toLowerCase().trim();
      return !needle || [operation.id, operation.path, operation.tag, operation.classification].some((value) => value.includes(needle));
    }),
  );
  const grouped = $derived(
    filtered.reduce<Record<string, OperationSummary[]>>((groups, operation) => {
      groups[operation.tag] ??= [];
      groups[operation.tag].push(operation);
      return groups;
    }, {}),
  );
  const selected = $derived(operations.find((operation) => operation.id === selectedId) ?? operations[0]);
</script>

<main class="page page-wide">
  <section class="hero">
    <div>
      <p class="eyebrow">Project {$page.params.id}</p>
      <h1>Operations</h1>
      <p class="muted">Classify imported spec operations, inspect hints, and stage override JSON before compiler preview.</p>
    </div>
    <div class="button-row">
      <a class="btn btn-secondary" href="/projects/{$page.params.id}">Overview</a>
      <a class="btn btn-primary" href="/projects/{$page.params.id}/builder">Open Builder</a>
    </div>
  </section>

  {#if specs.isError}
    <article class="card card-subtle" style="margin-bottom: 1rem;">
      <span class="badge badge-warning">Spec demo mode</span>
      <p class="muted" style="margin: 0.75rem 0 0;">The API spec endpoint is not available yet, so this page uses seeded operations while preserving the final data shape.</p>
    </article>
  {/if}

  <section class="grid grid-2">
    <article class="card section-stack">
      <div class="hero" style="margin-bottom: 0;">
        <div>
          <h2 style="margin-bottom: 0.35rem;">Operation explorer</h2>
          <p class="muted">{operations.length} total operations • grouped by tag</p>
        </div>
      </div>
      <input bind:value={filter} placeholder="Filter by id, path, tag, classification…" aria-label="Filter operations" />

      <div class="section-stack">
        {#each Object.entries(grouped) as [tag, items] (tag)}
          <section class="card card-subtle">
            <div class="project-title" style="justify-content: space-between;">
              <h3 style="margin: 0;">{tag}</h3>
              <span class="badge">{items.length} operations</span>
            </div>
            <div class="table" style="margin-top: 0.8rem;">
              {#each items as operation (operation.id)}
                <button class="table-row" onclick={() => (selectedId = operation.id)} aria-pressed={selected?.id === operation.id}>
                  <span class="method {operation.method.toLowerCase()}">{operation.method}</span>
                  <strong>{operation.path}</strong>
                  <span class="muted">{operation.id}</span>
                  <span class="badge">{operation.classification}</span>
                  {#if operation.destructive}
                    <span class="badge badge-warning">⚠ destr</span>
                  {:else}
                    <span class="badge badge-success">safe</span>
                  {/if}
                </button>
              {/each}
            </div>
          </section>
        {/each}
      </div>
    </article>

    <aside class="card section-stack">
      {#if selected}
        <div class="project-title">
          <span class="method {selected.method.toLowerCase()}">{selected.method}</span>
          <h2 style="margin: 0;">{selected.id}</h2>
          {#if selected.destructive}
            <span class="badge badge-warning">destructive blocked by default</span>
          {:else}
            <span class="badge badge-success">safe operation</span>
          {/if}
        </div>
        <p class="muted">{selected.path}</p>

        <label>
          <span class="eyebrow">Classification</span>
          <select value={selected.classification} aria-label="Classification">
            <option>list</option>
            <option>read</option>
            <option>create</option>
            <option>update</option>
            <option>delete</option>
            <option>action</option>
          </select>
        </label>

        <section class="card card-subtle">
          <h3>Query hints</h3>
          {#if selected.queryHints?.length}
            <div class="button-row">
              {#each selected.queryHints as hint}
                <span class="badge">{hint}</span>
              {/each}
            </div>
          {:else}
            <p class="muted">No query hints detected for this operation.</p>
          {/if}
        </section>

        <section class="card card-subtle">
          <h3>Response schema</h3>
          <pre>{selected.responseShape ?? 'Unknown response shape'}</pre>
        </section>

        <label>
          <span class="eyebrow">Overrides JSON</span>
          <textarea class="editor-pane" bind:value={overrideJson} aria-label="Override JSON"></textarea>
        </label>
        <div class="button-row">
          <button class="btn btn-primary">Save Override</button>
          <button class="btn btn-secondary">Reset</button>
        </div>
      {/if}
    </aside>
  </section>
</main>
