<script lang="ts">
  import { QueryClientProvider, QueryClient } from '@tanstack/svelte-query';
  import '$lib/api/client';

  const queryClient = new QueryClient();
  let { children } = $props();
</script>

<QueryClientProvider client={queryClient}>
  {@render children()}
</QueryClientProvider>
