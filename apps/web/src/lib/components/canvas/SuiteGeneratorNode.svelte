<script lang="ts">
	import { Handle, Position } from "@xyflow/svelte";
	import { Badge } from "$lib/components/ui/badge";
	import type { SuiteGeneratorNodeData } from "$lib/canvas/types";

	let { data }: { data: SuiteGeneratorNodeData & { matched_count?: number; case_count?: number } } = $props();
	const strategy = $derived(data.strategy);
</script>

<div class="min-w-72 rounded-lg border border-emerald-500/30 bg-card shadow-xl ring-1 ring-emerald-950/60">
	<Handle type="target" position={Position.Left} id="suite:selection" class="!size-2.5 !bg-emerald-400" />
	<div class="border-b border-border px-3 py-2">
		<div class="flex items-center gap-2">
			<Badge class="border-emerald-500/30 bg-emerald-500/20 text-emerald-300">SUITE</Badge>
			<p class="truncate text-sm font-semibold">{data.name || "Generated Suite"}</p>
		</div>
		<p class="mt-1 text-[0.7rem] text-muted-foreground">Selector-driven generated coverage</p>
	</div>
	<div class="space-y-2 p-3 text-xs">
		<div class="flex flex-wrap gap-1">
			{#each data.selector?.tags ?? [] as tag (tag)}<Badge variant="secondary">{tag}</Badge>{/each}
			{#each data.selector?.classifications ?? [] as cls (cls)}<Badge variant="outline">{cls}</Badge>{/each}
		</div>
		<div class="grid grid-cols-2 gap-2 text-[0.72rem] text-muted-foreground">
			<span>Default: {strategy?.default_list ? "on" : "off"}</span>
			<span>Pagination: {strategy?.pagination ? "on" : "off"}</span>
			<span>Search: {strategy?.search_from_response ? "on" : "off"}</span>
			<span>Enums: {strategy?.enum_filters ? "on" : "off"}</span>
		</div>
		<div class="rounded border border-border bg-background/60 p-2 text-[0.72rem]">
			<span class="font-semibold text-foreground">{data.matched_count ?? 0}</span> operations matched · <span class="font-semibold text-foreground">{data.case_count ?? 0}</span> estimated cases
		</div>
	</div>
</div>
