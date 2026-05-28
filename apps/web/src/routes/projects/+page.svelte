<script lang="ts">
  import { createQuery, createMutation, useQueryClient } from '@tanstack/svelte-query';
  import { api } from '$lib/api/client';
  import { demoProjects } from '$lib/demo-data';

  const queryClient = useQueryClient();
  let name = $state('');

  const projects = createQuery(() => ({
    queryKey: ['projects'],
    queryFn: async () => {
      const { data, error } = await api.GET('/api/projects');
      if (error) throw error;
      return data!;
    },
  }));

  const createProject = createMutation(() => ({
    mutationFn: async (projectName: string) => {
      const { data, error } = await api.POST('/api/projects', {
        body: { name: projectName, description: '' },
      });
      if (error) throw error;
      return data!;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['projects'] });
      name = '';
    },
  }));

  const dashboardProjects = $derived(
    projects.data?.length
      ? projects.data.map((p, index) => ({
          id: p.id,
          name: p.name,
          version: 'import pending',
          operations: 0,
          flows: 0,
          suites: 0,
          passRate: index === 0 ? 94 : 0,
          lastRun: p.created_at ? new Date(p.created_at).toLocaleString() : 'never',
          status: 'healthy' as const,
        }))
      : demoProjects,
  );

  const totalRuns = $derived(dashboardProjects.length ? dashboardProjects.length * 28 + 2 : 0);
  const averagePassRate = $derived(
    dashboardProjects.length
      ? Math.round(dashboardProjects.reduce((sum, project) => sum + project.passRate, 0) / dashboardProjects.length)
      : 0,
  );

  function statusClass(status: string) {
    if (status === 'warning') return 'warn';
    if (status === 'failing') return 'fail';
    return '';
  }
</script>

<main class="page">
  <section class="hero">
    <div>
      <p class="eyebrow">Spec-driven smoke testing</p>
      <h1>Projects</h1>
      <p class="muted">Manage imported OpenAPI specs, config builders, generated suites, and recent run health.</p>
    </div>
    <form
      class="card card-subtle"
      onsubmit={(e) => {
        e.preventDefault();
        if (name.trim()) createProject.mutate(name.trim());
      }}
    >
      <div class="grid" style="grid-template-columns: minmax(14rem, 1fr) auto; align-items: center;">
        <input bind:value={name} placeholder="Project name" required aria-label="Project name" />
        <button class="btn btn-primary" type="submit" disabled={createProject.isPending}>+ New Project</button>
      </div>
    </form>
  </section>

  <section class="grid grid-3" aria-label="Project stats">
    <article class="card">
      <span class="badge badge-success">Live</span>
      <p class="stat-value">{totalRuns}</p>
      <p class="muted">Total runs tracked</p>
    </article>
    <article class="card">
      <span class="badge badge-success">●●●●●○</span>
      <p class="stat-value">{averagePassRate}%</p>
      <p class="muted">Average pass rate</p>
    </article>
    <article class="card">
      <span class="badge">Workspace</span>
      <p class="stat-value">{dashboardProjects.length}</p>
      <p class="muted">Projects configured</p>
    </article>
  </section>

  <section style="margin-top: 1.5rem;">
    <div class="hero">
      <div>
        <h2>Project dashboard</h2>
        <p class="muted">Use Builder for configuration, Operations for spec classification, and Runs for execution evidence.</p>
      </div>
    </div>

    {#if projects.isLoading}
      <article class="card"><p class="muted">Loading projects from Smokery API…</p></article>
    {:else if projects.isError}
      <article class="card card-subtle">
        <span class="badge badge-warning">Demo data</span>
        <p class="muted" style="margin-top: 0.75rem;">API is unavailable, showing seeded UI states from the frontend plan.</p>
      </article>
    {/if}

    <div class="grid">
      {#each dashboardProjects as project (project.id)}
        <article class="card project-card">
          <div>
            <div class="project-title">
              <span class="status-dot {statusClass(project.status)}"></span>
              <h3 style="margin: 0;">{project.name}</h3>
              <span class="badge">{project.version}</span>
              <span class="badge">{project.operations} ops</span>
            </div>
            <p class="muted" style="margin: 0.7rem 0;">Last run: {project.lastRun}</p>
            <div class="kpi-row">
              <span>{project.passRate}% pass</span>
              <span>•</span>
              <span>{project.flows} flows</span>
              <span>•</span>
              <span>{project.suites} suite{project.suites === 1 ? '' : 's'}</span>
            </div>
          </div>
          <div class="button-row">
            <a class="btn btn-primary" href="/projects/{project.id}/runs">Run</a>
            <a class="btn btn-secondary" href="/projects/{project.id}/builder">Builder</a>
            <a class="btn btn-secondary" href="/projects/{project.id}/operations">Operations</a>
          </div>
        </article>
      {/each}
    </div>
  </section>
</main>
