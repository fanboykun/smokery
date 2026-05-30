<script lang="ts">
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { createProjectConfigStore, type Flow, type FlowStep } from '$lib/stores/project-config';
  import * as Card from '$lib/components/ui/card';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';
  import { Input } from '$lib/components/ui/input';
  import { Label } from '$lib/components/ui/label';
  import { Textarea } from '$lib/components/ui/textarea';
  import { Separator } from '$lib/components/ui/separator';
  import OperationPicker from '$lib/components/OperationPicker.svelte';

  const projectId = $page.params.id!;
  const flowId = $page.params.fid!;
  const config = createProjectConfigStore(projectId);
  const isNew = flowId === 'new';

  let flow = $state<Flow>(
    isNew
      ? { id: crypto.randomUUID(), name: '', environment: '', steps: [], cleanup: [] }
      : { ...($config.flows.find((f) => f.id === flowId) ?? { id: flowId, name: '', environment: '', steps: [], cleanup: [] }) },
  );

  let editingStepIdx = $state<number | null>(null);
  let editingCleanup = $state(false);

  const steps = $derived(editingCleanup ? (flow.cleanup ?? []) : flow.steps);
  const editingStep = $derived(editingStepIdx !== null ? steps[editingStepIdx] : null);

  function addStep(cleanup = false) {
    const step: FlowStep = { name: '', operation_id: '', captures: [], assertions: [] };
    if (cleanup) {
      flow.cleanup = [...(flow.cleanup ?? []), step];
      editingCleanup = true;
      editingStepIdx = (flow.cleanup?.length ?? 1) - 1;
    } else {
      flow.steps = [...flow.steps, step];
      editingCleanup = false;
      editingStepIdx = flow.steps.length - 1;
    }
  }

  function removeStep(idx: number, cleanup = false) {
    if (cleanup) {
      flow.cleanup = (flow.cleanup ?? []).filter((_, i) => i !== idx);
    } else {
      flow.steps = flow.steps.filter((_, i) => i !== idx);
    }
    editingStepIdx = null;
  }

  function moveStep(idx: number, dir: -1 | 1, cleanup = false) {
    const arr = cleanup ? [...(flow.cleanup ?? [])] : [...flow.steps];
    const target = idx + dir;
    if (target < 0 || target >= arr.length) return;
    [arr[idx], arr[target]] = [arr[target], arr[idx]];
    if (cleanup) flow.cleanup = arr;
    else flow.steps = arr;
    editingStepIdx = target;
  }

  function saveFlow() {
    if (!flow.name) return;
    config.update((c) => {
      const idx = c.flows.findIndex((f) => f.id === flow.id);
      if (idx >= 0) c.flows[idx] = flow;
      else c.flows = [...c.flows, flow];
      return c;
    });
    goto(`/projects/${projectId}/environments`);
  }

  function deleteFlow() {
    config.update((c) => ({ ...c, flows: c.flows.filter((f) => f.id !== flow.id) }));
    goto(`/projects/${projectId}/environments`);
  }

  // Step editing helpers
  let captName = $state('');
  let captSource = $state('body');
  let captPath = $state('');
  let assertType = $state('status');
  let assertExpected = $state('');
  let assertPath = $state('');

  function addCapture() {
    if (!editingStep || !captName || !captPath) return;
    editingStep.captures = [...(editingStep.captures ?? []), { name: captName, source: captSource, path: captPath }];
    captName = ''; captPath = '';
  }

  function addAssertion() {
    if (!editingStep) return;
    const a: { type: string; expected?: unknown; path?: string } = { type: assertType };
    if (assertExpected) a.expected = isNaN(Number(assertExpected)) ? assertExpected : Number(assertExpected);
    if (assertPath) a.path = assertPath;
    editingStep.assertions = [...(editingStep.assertions ?? []), a];
    assertExpected = ''; assertPath = '';
  }
</script>

<main class="mx-auto max-w-5xl px-6 py-8">
  <div class="mb-6 flex flex-wrap items-end justify-between gap-4">
    <div>
      <p class="text-xs font-bold uppercase tracking-widest text-primary">Flow Builder</p>
      <h1 class="text-3xl font-bold">{flow.name || 'New Flow'}</h1>
    </div>
    <div class="flex gap-2">
      <Button variant="outline" href="/projects/{projectId}/environments">← Back</Button>
      {#if !isNew}
        <Button variant="destructive" onclick={deleteFlow}>Delete</Button>
      {/if}
      <Button onclick={saveFlow}>Save Flow</Button>
    </div>
  </div>

  <div class="grid grid-cols-1 gap-6 lg:grid-cols-[1fr_380px]">
    <!-- Flow config + steps list -->
    <div class="space-y-4">
      <Card.Root>
        <Card.Content class="space-y-3 pt-4">
          <div class="grid grid-cols-2 gap-3">
            <div class="space-y-1"><Label>Name</Label><Input bind:value={flow.name} placeholder="User CRUD Flow" /></div>
            <div class="space-y-1">
              <Label>Environment</Label>
              <select class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm" bind:value={flow.environment}>
                <option value="">Select…</option>
                {#each $config.environments as env (env.id)}
                  <option value={env.id}>{env.name}</option>
                {/each}
              </select>
            </div>
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div class="space-y-1">
              <Label>Auth Profile (optional)</Label>
              <select class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm" bind:value={flow.auth}>
                <option value="">None</option>
                {#each $config.auth_profiles as auth (auth.id)}
                  <option value={auth.id}>{auth.name}</option>
                {/each}
              </select>
            </div>
            <div class="space-y-1"><Label>Description</Label><Input bind:value={flow.description} placeholder="Optional" /></div>
          </div>
        </Card.Content>
      </Card.Root>

      <!-- Steps -->
      <div class="flex items-center justify-between">
        <h2 class="text-lg font-bold">Steps ({flow.steps.length})</h2>
        <Button size="sm" onclick={() => addStep(false)}>+ Add Step</Button>
      </div>
      {#each flow.steps as step, idx (idx)}
        <button
          class="flex w-full items-center gap-3 rounded-lg border border-border p-3 text-left text-sm transition-colors hover:bg-secondary {editingStepIdx === idx && !editingCleanup ? 'ring-1 ring-primary' : ''}"
          onclick={() => { editingStepIdx = idx; editingCleanup = false; }}
        >
          <Badge variant="outline">{idx + 1}</Badge>
          <span class="flex-1 font-medium">{step.name || step.operation_id || '(unnamed)'}</span>
          <Badge variant="secondary">{step.assertions?.length ?? 0} assertions</Badge>
          <Badge variant="secondary">{step.captures?.length ?? 0} captures</Badge>
        </button>
      {/each}

      <Separator />

      <!-- Cleanup steps -->
      <div class="flex items-center justify-between">
        <h2 class="text-lg font-bold">Cleanup ({flow.cleanup?.length ?? 0})</h2>
        <Button size="sm" variant="outline" onclick={() => addStep(true)}>+ Add Cleanup</Button>
      </div>
      {#each flow.cleanup ?? [] as step, idx (idx)}
        <button
          class="flex w-full items-center gap-3 rounded-lg border border-border p-3 text-left text-sm transition-colors hover:bg-secondary {editingStepIdx === idx && editingCleanup ? 'ring-1 ring-primary' : ''}"
          onclick={() => { editingStepIdx = idx; editingCleanup = true; }}
        >
          <Badge variant="outline">C{idx + 1}</Badge>
          <span class="flex-1 font-medium">{step.name || step.operation_id || '(unnamed)'}</span>
        </button>
      {/each}
    </div>

    <!-- Step editor -->
    <aside>
      {#if editingStep && editingStepIdx !== null}
        <Card.Root>
          <Card.Header class="flex-row items-center justify-between">
            <Card.Title class="text-base">Step {editingCleanup ? `C${editingStepIdx + 1}` : editingStepIdx + 1}</Card.Title>
            <div class="flex gap-1">
              <Button variant="ghost" size="xs" onclick={() => moveStep(editingStepIdx!, -1, editingCleanup)}>↑</Button>
              <Button variant="ghost" size="xs" onclick={() => moveStep(editingStepIdx!, 1, editingCleanup)}>↓</Button>
              <Button variant="destructive" size="xs" onclick={() => removeStep(editingStepIdx!, editingCleanup)}>×</Button>
            </div>
          </Card.Header>
          <Card.Content class="space-y-3">
            <div class="space-y-1"><Label>Name</Label><Input bind:value={editingStep.name} placeholder="Create user" /></div>
            <OperationPicker
              projectId={projectId}
              specId={undefined}
              bind:value={editingStep.operation_id}
              onchange={(opId) => {
                editingStep.operation_id = opId;
              }}
              label="Operation"
              placeholder="Search operations…"
            />
            <div class="space-y-1">
              <Label>Body (JSON)</Label>
              <Textarea value={typeof editingStep.body === 'string' ? editingStep.body : JSON.stringify(editingStep.body ?? '', null, 2)} oninput={(e) => { if (editingStep) editingStep.body = e.currentTarget.value; }} class="min-h-[4rem] font-mono text-xs" placeholder={'{\n}'} />
            </div>

            <Separator />

            <!-- Captures -->
            <div class="space-y-2">
              <Label>Captures</Label>
              {#each editingStep.captures ?? [] as cap, i (i)}
                <Badge variant="outline" class="font-mono text-xs">{cap.name} ← {cap.source}:{cap.path}</Badge>
              {/each}
              <div class="grid grid-cols-3 gap-1">
                <Input bind:value={captName} placeholder="name" class="text-xs" />
                <select class="rounded-md border border-input bg-background px-2 py-1 text-xs" bind:value={captSource}>
                  <option value="body">body</option>
                  <option value="header">header</option>
                </select>
                <Input bind:value={captPath} placeholder="$.id" class="text-xs" />
              </div>
              <Button variant="outline" size="xs" onclick={addCapture}>+ Capture</Button>
            </div>

            <Separator />

            <!-- Assertions -->
            <div class="space-y-2">
              <Label>Assertions</Label>
              {#each editingStep.assertions ?? [] as a, i (i)}
                <Badge variant="outline" class="font-mono text-xs">{a.type}{a.path ? ` @ ${a.path}` : ''}{a.expected != null ? ` = ${a.expected}` : ''}</Badge>
              {/each}
              <div class="grid grid-cols-3 gap-1">
                <select class="rounded-md border border-input bg-background px-2 py-1 text-xs" bind:value={assertType}>
                  <option value="status">status</option>
                  <option value="jsonpath">jsonpath</option>
                  <option value="not_empty">not_empty</option>
                  <option value="schema">schema</option>
                  <option value="list_shape">list_shape</option>
                </select>
                <Input bind:value={assertPath} placeholder="path" class="text-xs" />
                <Input bind:value={assertExpected} placeholder="expected" class="text-xs" />
              </div>
              <Button variant="outline" size="xs" onclick={addAssertion}>+ Assertion</Button>
            </div>
          </Card.Content>
        </Card.Root>
      {:else}
        <Card.Root><Card.Content class="py-8 text-center text-sm text-muted-foreground">Select a step to edit.</Card.Content></Card.Root>
      {/if}
    </aside>
  </div>
</main>
