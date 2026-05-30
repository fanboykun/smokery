import type { ProjectConfig, SuiteSelector, SuiteStrategy } from "$lib/stores/project-config";

export type CanvasNodeType = "operationNode" | "suiteGeneratorNode" | "environmentNode";
export type CanvasEdgeType = "sequence" | "dataLink" | "suiteSelection";

export interface CanvasPosition {
	x: number;
	y: number;
}

export interface CanvasViewport extends CanvasPosition {
	zoom: number;
}

export interface CanvasGraph {
	version: number;
	default_environment?: string;
	default_auth?: string;
	nodes: CanvasNode[];
	edges: CanvasEdge[];
	viewport?: CanvasViewport;
}

export interface CanvasNode<TData extends Record<string, unknown> = Record<string, unknown>> {
	id: string;
	type: CanvasNodeType;
	position: CanvasPosition;
	data: TData;
}

export interface CanvasEdge<TData extends Record<string, unknown> = Record<string, unknown>> {
	id: string;
	type: CanvasEdgeType;
	source: string;
	target: string;
	sourceHandle?: string;
	targetHandle?: string;
	data?: TData;
}

export interface CanvasOperationParameter {
	name: string;
	in: "path" | "query" | "header" | "cookie" | string;
	required: boolean;
	schema?: unknown;
}

export interface CanvasOperation {
	operation_id: string;
	method: string;
	path: string;
	summary: string;
	tags: string[] | null;
	classification: string;
	is_destructive: boolean;
	parameters?: CanvasOperationParameter[] | null;
	request_schema?: unknown;
	response_schema?: unknown;
	query_hints?: {
		pagination_params?: string[] | null;
		search_params?: string[] | null;
		enum_filters?: { name: string; values: string[] }[] | null;
	};
}

export interface OperationNodeData extends Record<string, unknown> {
	operation_id: string;
	label?: string;
	destructive_acknowledged?: boolean;
	cleanup?: boolean;
}

export interface SuiteGeneratorNodeData extends Record<string, unknown> {
	name?: string;
	selector?: SuiteSelector;
	strategy?: SuiteStrategy;
}

export interface DataLinkStrategy {
	selector: {
		mode: "first" | "random" | "filter";
		filter?: string;
	};
	projection?: {
		path?: string;
	};
	retry: {
		enabled: boolean;
		max_alternates: number;
		scope: "target";
	};
}

export interface DataLinkEdgeData extends Record<string, unknown> {
	strategy?: DataLinkStrategy;
}

export interface GraphConversionBase {
	environments: ProjectConfig["environments"];
	auth_profiles: ProjectConfig["auth_profiles"];
}
