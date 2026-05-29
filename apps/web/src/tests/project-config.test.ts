import { describe, it, expect, vi, beforeEach } from 'vitest';
import { get } from 'svelte/store';

// Mock $app/environment to simulate browser
vi.mock('$app/environment', () => ({ browser: true }));

import { createProjectConfigStore } from '$lib/stores/project-config';

describe('project config store', () => {
	beforeEach(() => {
		localStorage.clear();
	});

	it('returns default config for new project', () => {
		const store = createProjectConfigStore('test-project');
		const config = get(store);
		expect(config.environments).toEqual([]);
		expect(config.auth_profiles).toEqual([]);
		expect(config.flows).toEqual([]);
		expect(config.suites).toEqual([]);
	});

	it('persists changes to localStorage', () => {
		const store = createProjectConfigStore('test-project');
		store.update((c) => ({
			...c,
			environments: [{ id: '1', name: 'staging', base_url: 'http://localhost' }],
		}));

		const raw = localStorage.getItem('smokery:config:test-project');
		expect(raw).toBeTruthy();
		const parsed = JSON.parse(raw!);
		expect(parsed.environments).toHaveLength(1);
		expect(parsed.environments[0].name).toBe('staging');
	});

	it('loads existing config from localStorage', () => {
		const existing = {
			environments: [{ id: '2', name: 'prod', base_url: 'https://api.example.com' }],
			auth_profiles: [],
			flows: [],
			suites: [],
		};
		localStorage.setItem('smokery:config:existing', JSON.stringify(existing));

		const store = createProjectConfigStore('existing');
		const config = get(store);
		expect(config.environments).toHaveLength(1);
		expect(config.environments[0].name).toBe('prod');
	});
});
