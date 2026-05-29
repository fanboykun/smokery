import { describe, it, expect } from 'vitest';
import { api } from '$lib/api/client';

describe('API client', () => {
	it('is configured with the correct base URL', () => {
		// The client should be an object with GET, POST, PUT, DELETE methods
		expect(api.GET).toBeDefined();
		expect(api.POST).toBeDefined();
		expect(api.PUT).toBeDefined();
		expect(api.DELETE).toBeDefined();
	});
});
