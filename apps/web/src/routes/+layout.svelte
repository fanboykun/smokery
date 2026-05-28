<script lang="ts">
  import { QueryClientProvider, QueryClient } from '@tanstack/svelte-query';
  import '$lib/api/client';
  import '../app.css';

  const queryClient = new QueryClient();
  let { children } = $props();
</script>

<QueryClientProvider client={queryClient}>
  <div class="app-shell">
    <header class="topbar">
      <a class="brand" href="/projects" aria-label="Smokery home">
        <span class="brand-mark">S</span>
        <span>SMOKERY</span>
      </a>
      <nav aria-label="Primary">
        <a href="/projects">Projects</a>
        <a href="/projects/demo/builder">Builder demo</a>
        <a href="/projects/demo/operations">Operations demo</a>
      </nav>
    </header>
    {@render children()}
  </div>
</QueryClientProvider>
