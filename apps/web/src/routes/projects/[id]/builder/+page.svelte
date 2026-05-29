<script lang="ts">
  import { page } from '$app/stores';
  import { demoBuilderConfig, demoOperations } from '$lib/demo-data';

  let selectedEnvironment = $state('staging');
  let selectedSuite = $state('list-endpoints');
  let compileState = $state<'ready' | 'compiled'>('ready');

  const safeOperations = $derived(demoOperations.filter((operation) => !operation.destructive));
  const riskyOperations = $derived(demoOperations.filter((operation) => operation.destructive));
  const listOperations = $derived(safeOperations.filter((operation) => operation.classification === 'list'));
  const generatedCases = $derived(
    listOperations.flatMap((operation) => [
      `${operation.method} ${operation.path} — default_list`,
      ...(operation.queryHints?.includes('page') ? [`${operation.method} ${operation.path} — pagination`] : []),
      ...(operation.queryHints?.includes('q') ? [`${operation.method} ${operation.path} — search`] : []),
      ...(operation.queryHints?.includes('status') ? [`${operation.method} ${operation.path} — enum:status`] : []),
    ]),
  );

  function compilePreview() {
    compileState = 'compiled';
  }
</script>

<main class="page page-wide">
  <section class="hero">
    <div>
      <p class="eyebrow">Project {$page.params.id}</p>
      <h1>Smoke config builder</h1>
      <p class="muted">Compose environments, flows, and generated suites while previewing what the compiler will execute.</p>
    </div>
    <div class="button-row">
      <a class="btn btn-secondary" href="/projects/{$page.params.id}">← Project</a>
      <button class="btn btn-secondary" onclick={compilePreview}>Compile</button>
      <button class="btn btn-primary">▶ Run</button>
      <a class="btn btn-secondary" href="/projects/{$page.params.id}/plan">Plan Preview</a>
    </div>
  </section>

  <section class="split-pane">
    <aside class="card section-stack">
      <section class="section-stack">
        <div class="project-title" style="justify-content: space-between;">
          <h2 style="margin: 0;">Config panel</h2>
          <span class="badge badge-success">reactive</span>
        </div>

        <section class="card card-subtle section-stack">
          <div class="project-title" style="justify-content: space-between;">
            <h3 style="margin: 0;">Environments</h3>
            <button class="btn btn-secondary" style="padding: 0.45rem 0.7rem;">+ Add</button>
          </div>
          {#each demoBuilderConfig.environments as env (env.id)}
            <button class="list-row" onclick={() => (selectedEnvironment = env.id)} aria-pressed={selectedEnvironment === env.id}>
              <span><strong>{env.name}</strong><br /><span class="muted">{env.baseUrl}</span></span>
              <span class="badge {env.kind === 'safe' ? 'badge-success' : 'badge-warning'}">{env.kind}</span>
            </button>
          {/each}
        </section>

        <section class="card card-subtle section-stack">
          <div class="project-title" style="justify-content: space-between;">
            <h3 style="margin: 0;">Flows</h3>
            <button class="btn btn-secondary" style="padding: 0.45rem 0.7rem;">+ New Flow</button>
          </div>
          {#each demoBuilderConfig.flows as flow (flow.id)}
            <article class="list-row">
              <span><strong>📋 {flow.name}</strong><br /><span class="muted">{flow.steps} steps • {flow.environment}</span></span>
              <a class="badge" href="/projects/{$page.params.id}/builder/flows/{flow.id}">Edit</a>
            </article>
          {/each}
        </section>

        <section class="card card-subtle section-stack">
          <div class="project-title" style="justify-content: space-between;">
            <h3 style="margin: 0;">Suites</h3>
            <button class="btn btn-secondary" style="padding: 0.45rem 0.7rem;">+ New Suite</button>
          </div>
          {#each demoBuilderConfig.suites as suite (suite.id)}
            <button class="list-row" onclick={() => (selectedSuite = suite.id)} aria-pressed={selectedSuite === suite.id}>
              <span><strong>🔄 {suite.name}</strong><br /><span class="muted">auto-generated • {suite.cases} cases</span></span>
              <span class="badge">Configure</span>
            </button>
          {/each}
        </section>
      </section>
    </aside>

    <section class="card section-stack">
      <div class="project-title" style="justify-content: space-between;">
        <div>
          <h2 style="margin: 0 0 0.35rem;">Live preview</h2>
          <p class="muted" style="margin: 0;">This local compiler-shaped preview is ready for the generated plan preview endpoint.</p>
        </div>
        <span class="badge {compileState === 'compiled' ? 'badge-success' : 'badge-warning'}">{compileState}</span>
      </div>

      <section class="card card-subtle">
        <div class="grid grid-3">
          <div>
            <span class="badge badge-success">✓ environments</span>
            <p class="stat-value">{demoBuilderConfig.environments.length}</p>
          </div>
          <div>
            <span class="badge badge-success">✓ flows</span>
            <p class="stat-value">{demoBuilderConfig.flows.length}</p>
          </div>
          <div>
            <span class="badge badge-success">✓ generated cases</span>
            <p class="stat-value">{generatedCases.length}</p>
          </div>
        </div>
      </section>

      <section class="grid grid-2">
        <article class="card card-subtle">
          <h3>Flow: User CRUD</h3>
          <ol class="preview-list">
            <li>POST /users — create</li>
            <li>GET /users/{'{id}'} — read</li>
            <li>PUT /users/{'{id}'} — update</li>
            <li>DELETE /users/{'{id}'} — cleanup</li>
          </ol>
        </article>
        <article class="card card-subtle">
          <h3>Suite: List Endpoints</h3>
          <ul class="preview-list">
            {#each generatedCases.slice(0, 8) as item}
              <li>{item}</li>
            {/each}
          </ul>
        </article>
      </section>

      <section class="warning-list">
        <h3 style="margin-bottom: 0;">Warnings</h3>
        {#each riskyOperations as operation (operation.id)}
          <article class="list-row">
            <span><strong>⚠ {operation.id} needs explicit approval</strong><br /><span class="muted">{operation.method} {operation.path} is protected by the compiler safety model.</span></span>
            <a class="badge badge-warning" href="/projects/{$page.params.id}/operations">Review</a>
          </article>
        {/each}
      </section>

      <footer class="footer-bar">
        <span>Operations from spec: {demoOperations.length} total</span>
        <span>{listOperations.length} list • {safeOperations.length} safe • {riskyOperations.length} protected</span>
      </footer>
    </section>
  </section>
</main>
