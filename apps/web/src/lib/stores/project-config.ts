import { writable } from 'svelte/store';
import { browser } from '$app/environment';

export interface Environment {
	id: string;
	name: string;
	base_url: string;
	headers?: Record<string, string>;
}

export interface AuthProfile {
	id: string;
	name: string;
	type: string;
	config: Record<string, unknown>;
}

export interface FlowStep {
	name: string;
	operation_id: string;
	params?: Record<string, unknown>;
	body?: unknown;
	headers?: Record<string, string>;
	captures?: { name: string; source: string; path: string }[];
	assertions?: { type: string; expected?: unknown; path?: string }[];
}

export interface Flow {
	id: string;
	name: string;
	description?: string;
	environment: string;
	auth?: string;
	steps: FlowStep[];
	cleanup?: FlowStep[];
}

export interface SuiteSelector {
	tags?: string[];
	classifications?: string[];
	paths?: string[];
	exclude?: string[];
}

export interface SuiteStrategy {
	default_list: boolean;
	pagination: boolean;
	search_from_response: boolean;
	enum_filters: boolean;
	empty_result_policy: string;
	max_cases_per_op?: number;
}

export interface Suite {
	id: string;
	name: string;
	description?: string;
	environment: string;
	auth?: string;
	selector: SuiteSelector;
	strategy: SuiteStrategy;
}

export interface ProjectConfig {
	environments: Environment[];
	auth_profiles: AuthProfile[];
	flows: Flow[];
	suites: Suite[];
}

function storageKey(projectId: string) {
	return `smokery:config:${projectId}`;
}

const defaultConfig: ProjectConfig = {
	environments: [],
	auth_profiles: [],
	flows: [],
	suites: [],
};

function loadConfig(projectId: string): ProjectConfig {
	if (!browser) return defaultConfig;
	const raw = localStorage.getItem(storageKey(projectId));
	if (!raw) return defaultConfig;
	try {
		return JSON.parse(raw);
	} catch {
		return defaultConfig;
	}
}

export function createProjectConfigStore(projectId: string) {
	const initial = loadConfig(projectId);
	const store = writable<ProjectConfig>(initial);

	if (browser) {
		store.subscribe((value) => {
			localStorage.setItem(storageKey(projectId), JSON.stringify(value));
		});
	}

	return store;
}
