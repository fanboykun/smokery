<script lang="ts">
  import { page } from '$app/stores';
  import { createQuery, createMutation, useQueryClient } from '@tanstack/svelte-query';
  import { api } from '$lib/api/client';

  const queryClient = useQueryClient();
  let body = $state('');
  let author = $state('user');

  const comments = createQuery(() => ({
    queryKey: ['comments', $page.params.runId],
    queryFn: async () => {
      const { data, error } = await api.GET('/api/runs/{id}/comments', { params: { path: { id: $page.params.runId! } } });
      if (error) throw error;
      return data!;
    },
  }));

  const addComment = createMutation(() => ({
    mutationFn: async () => {
      const { data, error } = await api.POST('/api/runs/{id}/comments', {
        params: { path: { id: $page.params.runId! } },
        body: { author, body },
      });
      if (error) throw error;
      return data!;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['comments', $page.params.runId] });
      body = '';
    },
  }));
</script>

<h1>Comments</h1>
<form onsubmit={(e) => { e.preventDefault(); addComment.mutate(); }}>
  <input bind:value={author} placeholder="Author" />
  <textarea bind:value={body} placeholder="Comment"></textarea>
  <button type="submit" disabled={addComment.isPending}>Post</button>
</form>

{#if comments.isLoading}
  <p>Loading...</p>
{:else if comments.data}
  <ul>
    {#each comments.data as c (c.id)}
      <li><strong>{c.author}</strong>: {c.body}</li>
    {/each}
  </ul>
{/if}
