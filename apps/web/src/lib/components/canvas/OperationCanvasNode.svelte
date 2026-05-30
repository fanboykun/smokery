<script lang="ts">
	import { Handle, Position } from "@xyflow/svelte";
	import { Badge } from "$lib/components/ui/badge";
	import type { CanvasOperation } from "$lib/canvas/types";

	let { data }: { data: { operation?: CanvasOperation; label?: string; destructive_acknowledged?: boolean } } = $props();
	const op = $derived(data.operation);
	const responseFields = $derived(schemaFields(op?.response_schema).slice(0, 8));
	const requestFields = $derived([
		...(op?.parameters ?? []).map((p) => ({ kind: p.in, path: p.name, required: p.required })),
		...schemaFields(op?.request_schema).slice(0, 5).map((f) => ({ kind: "body", path: f, required: false })),
	]);

	function methodClass(method = "GET") {
		const m = method.toUpperCase();
		if (m === "GET") return "bg-sky-500/20 text-sky-300 border-sky-500/30";
		if (m === "POST") return "bg-emerald-500/20 text-emerald-300 border-emerald-500/30";
		if (m === "DELETE") return "bg-red-500/20 text-red-300 border-red-500/30";
		if (m === "PATCH") return "bg-violet-500/20 text-violet-300 border-violet-500/30";
		return "bg-amber-500/20 text-amber-300 border-amber-500/30";
	}

	function schemaFields(schema: unknown, prefix = ""): string[] {
		if (!schema || typeof schema !== "object") return [];
		const s = schema as Record<string, any>;
		if (s.type === "array" && s.items) return schemaFields(s.items, `${prefix}[]`);
		const props = s.properties as Record<string, unknown> | undefined;
		if (!props) return prefix ? [prefix] : [];
		const out: string[] = [];
		for (const [key, value] of Object.entries(props)) {
			const next = prefix ? `${prefix}.${key}` : key;
			out.push(next);
			out.push(...schemaFields(value, next));
		}
		return out;
	}
</script>

<div class="min-w-72 max-w-80 rounded-lg border border-border bg-card shadow-xl ring-1 ring-black/20">
	<Handle type="target" position={Position.Top} id="flow-in" class="!size-3 !rounded-sm !border-2 !border-primary !bg-background" />
	<div class="border-b border-border px-3 py-2">
		<div class="flex items-center gap-2">
			<Badge class={`${methodClass(op?.method)} font-mono text-[0.65rem]`}>{op?.method ?? "OP"}</Badge>
			<p class="truncate text-sm font-semibold">{data.label || op?.operation_id || "Operation"}</p>
		</div>
		<p class="mt-1 truncate font-mono text-[0.7rem] text-muted-foreground">{op?.path}</p>
	</div>

	<div class="grid grid-cols-2 gap-0 text-[0.7rem]">
		<div class="border-r border-border p-2">
			<p class="mb-2 font-semibold uppercase tracking-widest text-muted-foreground">Request</p>
			{#if requestFields.length === 0}
				<p class="text-muted-foreground">No inputs</p>
			{/if}
			{#each requestFields as field (field.kind + field.path)}
				<div class="relative mb-1 rounded border border-border/60 bg-background/60 px-2 py-1 font-mono">
					<Handle type="target" position={Position.Left} id={`${field.kind}:${field.path}`} class="!size-2 !bg-primary" />
					<span class="text-muted-foreground">{field.kind}</span> {field.path}
				</div>
			{/each}
		</div>
		<div class="p-2">
			<p class="mb-2 font-semibold uppercase tracking-widest text-muted-foreground">Response</p>
			{#if responseFields.length === 0}
				<p class="text-muted-foreground">No schema</p>
			{/if}
			{#each responseFields as field (field)}
				<div class="relative mb-1 rounded border border-border/60 bg-background/60 px-2 py-1 font-mono">
					<Handle type="source" position={Position.Right} id={`response:${field}`} class="!size-2 !bg-primary" />
					{field}
				</div>
			{/each}
		</div>
	</div>

	<div class="flex items-center justify-between border-t border-border px-3 py-2 text-[0.7rem] text-muted-foreground">
		<span>{op?.classification ?? "unclassified"}</span>
		{#if op?.is_destructive}
			<span class={data.destructive_acknowledged ? "text-emerald-300" : "text-red-300"}>{data.destructive_acknowledged ? "acknowledged" : "needs ack"}</span>
		{/if}
	</div>
	<Handle type="source" position={Position.Bottom} id="flow-out" class="!size-3 !rounded-sm !border-2 !border-primary !bg-background" />
</div>
