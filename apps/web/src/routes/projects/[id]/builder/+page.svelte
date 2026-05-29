<script lang="ts">
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';

  onMount(() => {
    goto(`/projects/${$page.params.id}/plan`, { replaceState: true });
  });
</script>

<main class="mx-auto max-w-4xl px-6 py-8">
  <p class="text-muted-foreground">Redirecting to Plan Preview…</p>
</main>
