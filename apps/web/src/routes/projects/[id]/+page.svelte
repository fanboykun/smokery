<script lang="ts">
  import { page } from '$app/stores';
  import { createQuery } from '@tanstack/svelte-query';
  import { api } from '$lib/api/client';

  let specInput = $state('');
  let importSuccess = $state('');

  const project = createQuery(() => ({
    queryKey: ['projects', $page.params.id],
    queryFn: async () => {
      const { data, error } = await api.GET('/api/projects/{id}', { params: { path: { id: $page.params.id! } } });
      if (error) throw error;
      return data!;
    },
  }));

  async function handleImport() {
    const id = $page.params.id!;
    const { data, error } = await api.POST('/api/projects/{project-id}/specs', {
      params: { path: { 'project-id': id } },
      body: specInput,
      bodySerializer: (body) => body as unknown as string,
    });
    if (error) return;
    if (data) {
      importSuccess = `Imported: ${data.title} v${data.version} (${data.operations?.length ?? 0} operations)`;
      specInput = '';
    }
  }
</script>

<h1>{project.data?.name ?? 'Loading...'}</h1>
<nav>
  <a href="/projects/{$page.params.id}/runs">Runs</a>
</nav>
<h2>Import OpenAPI Spec</h2>
<textarea bind:value={specInput} rows="10" cols="80" placeholder="Paste OpenAPI JSON here"></textarea>
<button onclick={handleImport}>Import</button>
{#if importSuccess}
  <p>{importSuccess}</p>
{/if}
