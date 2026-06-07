import type { AuthProfile, Environment, Flow, FlowStep, ProjectConfig, Suite, SuiteStrategy } from "$lib/stores/project-config";
import type { CanvasEdge, CanvasGraph, CanvasOperation, DataLinkStrategy, GraphConversionBase, OperationNodeData, SuiteGeneratorNodeData } from "./types";

export function defaultDataLinkStrategy(): DataLinkStrategy {
	return {
		selector: { mode: "first" },
		retry: { enabled: false, max_alternates: 2, scope: "target" },
	};
}

export function canvasGraphToProjectConfig(
	graph: CanvasGraph,
	operations: CanvasOperation[],
	base: GraphConversionBase,
): ProjectConfig {
	const operationById = new Map(operations.map((op) => [op.operation_id, op]));
	const operationNodes = graph.nodes.filter((node) => node.type === "operationNode");
	const suiteNodes = graph.nodes.filter((node) => node.type === "suiteGeneratorNode");
	const environment = graph.default_environment || base.environments[0]?.id || "";
	const auth = graph.default_auth;

	const stepsByNode = new Map<string, FlowStep>();
	for (const node of operationNodes) {
		const data = node.data as OperationNodeData;
		const op = operationById.get(data.operation_id) ?? (data as any).operation;
		if (!op) continue;
		if (op.is_destructive && !data.destructive_acknowledged) {
			continue;
		}
		stepsByNode.set(node.id, {
			name: String(data.label || op.operation_id),
			operation_id: op.operation_id,
			captures: [],
			assertions: [{ type: "status", expected: defaultStatusForMethod(op.method) }],
		});
	}

	for (const edge of graph.edges.filter((edge) => edge.type === "dataLink")) {
		applyDataLink(edge, stepsByNode);
	}

	const orderedNodeIds = orderOperationNodes(operationNodes.map((node) => node.id), graph.edges);
	const steps = orderedNodeIds.map((id) => stepsByNode.get(id)).filter(Boolean) as FlowStep[];
	const flows: Flow[] = steps.length
		? [{ id: "canvas-flow", name: "Canvas Flow", environment, ...(auth ? { auth } : {}), steps }]
		: [];

	const suites: Suite[] = suiteNodes.map((node) => {
		const data = node.data as SuiteGeneratorNodeData;
		return {
			id: node.id,
			name: data.name || "Generated Suite",
			environment,
			...(auth ? { auth } : {}),
			selector: {
				tags: data.selector?.tags ?? [],
				classifications: data.selector?.classifications ?? [],
				paths: data.selector?.paths ?? [],
				exclude: data.selector?.exclude ?? [],
			},
			strategy: data.strategy ?? defaultSuiteStrategy(),
		};
	});

	return {
		environments: base.environments as Environment[],
		auth_profiles: base.auth_profiles as AuthProfile[],
		flows,
		suites,
		canvas: graph,
	};
}

function applyDataLink(edge: CanvasEdge, stepsByNode: Map<string, FlowStep>) {
	const source = stepsByNode.get(edge.source);
	const target = stepsByNode.get(edge.target);
	if (!source || !target || !edge.sourceHandle || !edge.targetHandle) return;

	const sourcePath = parseHandle(edge.sourceHandle).path;
	const targetHandle = parseHandle(edge.targetHandle);
	if (!sourcePath || !targetHandle.path) return;

	const captureName = captureNameFor(edge.source, sourcePath);
	source.captures = [...(source.captures ?? []), { name: captureName, source: "body", path: selectorPathToGJSON(sourcePath) }];

	const binding = `{{${captureName}}}`;
	if (targetHandle.kind === "path" || targetHandle.kind === "query") {
		target.params = { ...(target.params ?? {}), [targetHandle.path]: binding };
	} else if (targetHandle.kind === "header") {
		target.headers = { ...(target.headers ?? {}), [targetHandle.path]: binding };
	} else if (targetHandle.kind === "body") {
		target.body = setBodyPath(target.body, targetHandle.path, binding);
	}
}

function orderOperationNodes(nodeIds: string[], edges: CanvasEdge[]): string[] {
	const sequence = edges.filter((edge) => edge.type === "sequence");
	if (sequence.length === 0) return nodeIds;

	const targets = new Set(sequence.map((edge) => edge.target));
	const start = nodeIds.find((id) => !targets.has(id)) ?? nodeIds[0];
	const ordered = [start];
	const seen = new Set(ordered);
	let current = start;
	while (true) {
		const next = sequence.find((edge) => edge.source === current && !seen.has(edge.target))?.target;
		if (!next) break;
		ordered.push(next);
		seen.add(next);
		current = next;
	}
	for (const id of nodeIds) {
		if (!seen.has(id)) ordered.push(id);
	}
	return ordered;
}

function parseHandle(handle: string): { kind: string; path: string } {
	const [kind, ...rest] = handle.split(":");
	return { kind, path: rest.join(":") };
}

function selectorPathToGJSON(path: string): string {
	return path.replace(/\[\]/g, ".0").replace(/^\./, "");
}

function captureNameFor(nodeId: string, path: string): string {
	return `${nodeId}_${path}`.replace(/[^a-zA-Z0-9]+/g, "_").replace(/^_+|_+$/g, "");
}

function setBodyPath(body: unknown, path: string, value: string): Record<string, unknown> {
	const root = typeof body === "object" && body !== null && !Array.isArray(body) ? { ...(body as Record<string, unknown>) } : {};
	root[path] = value;
	return root;
}

function defaultStatusForMethod(method: string) {
	return method.toUpperCase() === "POST" ? 201 : 200;
}

function defaultSuiteStrategy(): SuiteStrategy {
	return {
		default_list: true,
		pagination: true,
		search_from_response: false,
		enum_filters: false,
		empty_result_policy: "allow",
		max_cases_per_op: 0,
	};
}
