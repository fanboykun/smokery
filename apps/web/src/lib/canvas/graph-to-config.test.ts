import { describe, expect, it } from 'vitest';
import { canvasGraphToProjectConfig, defaultDataLinkStrategy } from './graph-to-config';
import type { CanvasGraph, CanvasOperation } from './types';

const operations: CanvasOperation[] = [
	{
		operation_id: 'listUsers',
		method: 'GET',
		path: '/users',
		summary: 'List users',
		tags: ['users'],
		classification: 'list',
		is_destructive: false,
		parameters: [],
		query_hints: { search_params: ['q'] },
		response_schema: { type: 'object', properties: { data: { type: 'array', items: { type: 'object', properties: { id: { type: 'string' } } } } } },
	},
	{
		operation_id: 'deleteUser',
		method: 'DELETE',
		path: '/users/{user_id}',
		summary: 'Delete user',
		tags: ['users'],
		classification: 'delete',
		is_destructive: true,
		parameters: [{ name: 'user_id', in: 'path', required: true, schema: { type: 'string' } }],
	},
];

describe('canvasGraphToProjectConfig', () => {
	it('converts operation sequence and data links into a flow', () => {
		const graph: CanvasGraph = {
			version: 1,
			default_environment: 'env-1',
			nodes: [
				{ id: 'n-list', type: 'operationNode', position: { x: 0, y: 0 }, data: { operation_id: 'listUsers', label: 'List users' } },
				{ id: 'n-delete', type: 'operationNode', position: { x: 320, y: 0 }, data: { operation_id: 'deleteUser', label: 'Delete user', destructive_acknowledged: true } },
			],
			edges: [
				{ id: 'e-seq', type: 'sequence', source: 'n-list', target: 'n-delete' },
				{
					id: 'e-data',
					type: 'dataLink',
					source: 'n-list',
					sourceHandle: 'response:data[].id',
					target: 'n-delete',
					targetHandle: 'path:user_id',
					data: { strategy: defaultDataLinkStrategy() },
				},
			],
		};

		const config = canvasGraphToProjectConfig(graph, operations, {
			environments: [{ id: 'env-1', name: 'staging', base_url: 'https://api.example.com' }],
			auth_profiles: [],
		});

		expect(config.flows).toHaveLength(1);
		expect(config.flows[0].steps).toHaveLength(2);
		expect(config.flows[0].steps[0].captures).toEqual([{ name: 'n_list_data_id', source: 'body', path: 'data.0.id' }]);
		expect(config.flows[0].steps[1].params).toEqual({ user_id: '{{n_list_data_id}}' });
		expect(config.flows[0].steps[1].operation_id).toBe('deleteUser');
	});

	it('converts suite generator nodes into suites', () => {
		const graph: CanvasGraph = {
			version: 1,
			default_environment: 'env-1',
			nodes: [
				{
					id: 'suite-1',
					type: 'suiteGeneratorNode',
					position: { x: 0, y: 0 },
					data: {
						name: 'List suite',
						selector: { tags: ['users'], classifications: ['list'], paths: [], exclude: [] },
						strategy: { default_list: true, pagination: true, search_from_response: true, enum_filters: false, empty_result_policy: 'warn', max_cases_per_op: 2 },
					},
				},
			],
			edges: [],
		};

		const config = canvasGraphToProjectConfig(graph, operations, {
			environments: [{ id: 'env-1', name: 'staging', base_url: 'https://api.example.com' }],
			auth_profiles: [],
		});

		expect(config.suites).toEqual([
			{
				id: 'suite-1',
				name: 'List suite',
				environment: 'env-1',
				selector: { tags: ['users'], classifications: ['list'], paths: [], exclude: [] },
				strategy: { default_list: true, pagination: true, search_from_response: true, enum_filters: false, empty_result_policy: 'warn', max_cases_per_op: 2 },
			},
		]);
	});

	it('skips unacknowledged destructive operation nodes from config generation', () => {
		const graph: CanvasGraph = {
			version: 1,
			default_environment: 'env-1',
			nodes: [
				{ id: 'n-delete', type: 'operationNode', position: { x: 0, y: 0 }, data: { operation_id: 'deleteUser' } },
			],
			edges: [],
		};

		const config = canvasGraphToProjectConfig(graph, operations, {
			environments: [{ id: 'env-1', name: 'staging', base_url: 'https://api.example.com' }],
			auth_profiles: [],
		});

		expect(config.flows).toHaveLength(0);
	});
});
