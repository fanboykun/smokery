<script lang="ts">
  import { QueryClientProvider, QueryClient } from '@tanstack/svelte-query';
  import '../app.css';

  const queryClient = new QueryClient();
  let { children } = $props();
</script>

<QueryClientProvider client={queryClient}>
  <div class="min-h-screen">
    <header class="sticky top-0 z-10 flex items-center justify-between gap-4 border-b border-border bg-background/80 px-6 py-3 backdrop-blur-lg">
      <a class="flex items-center gap-3 font-extrabold tracking-widest" href="/projects" aria-label="Smokery home">
        <span class="grid size-9 place-items-center rounded-lg bg-primary text-primary-foreground text-sm shadow-lg">S</span>
        <span>SMOKERY</span>
      </a>
      <nav class="flex flex-wrap gap-2 text-sm text-muted-foreground" aria-label="Primary">
        <a class="rounded-full px-3 py-1.5 hover:bg-secondary hover:text-foreground transition-colors" href="/projects">Projects</a>
      </nav>
    </header>
    {@render children()}
  </div>
</QueryClientProvider>
