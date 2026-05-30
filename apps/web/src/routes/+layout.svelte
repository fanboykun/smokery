<script lang="ts">
  import { QueryClientProvider, QueryClient } from '@tanstack/svelte-query';
  import '../app.css';

  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        retry: 1,
        staleTime: 60000,
      },
      mutations: {
        retry: 1,
      },
    },
  });
  let { children } = $props();
</script>

<svelte:head>
  <title>Smokery</title>
  <meta name="description" content="Spec-driven smoke testing platform for OpenAPI services." />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
</svelte:head>

<QueryClientProvider client={queryClient}>
  <div class="min-h-screen bg-background text-foreground">
    <header class="sticky top-0 z-10 flex items-center justify-between gap-4 border-b border-border bg-background/80 px-6 py-3 backdrop-blur-lg shadow-sm">
      <a class="flex items-center gap-3 font-extrabold tracking-widest hover:opacity-80 transition-opacity" href="/projects" aria-label="Smokery home">
        <span class="grid size-9 place-items-center rounded-lg bg-primary text-primary-foreground text-sm shadow-lg">S</span>
        <span class="hidden sm:inline">SMOKERY</span>
      </a>
      <nav class="flex flex-wrap gap-2 text-sm text-muted-foreground" aria-label="Primary navigation">
        <a class="rounded-full px-3 py-1.5 hover:bg-secondary hover:text-foreground transition-colors" href="/projects">Projects</a>
      </nav>
    </header>
    <div class="bg-background">
      {@render children()}
    </div>
  </div>
</QueryClientProvider>
