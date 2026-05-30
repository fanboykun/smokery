import createClient from 'openapi-fetch';
import type { paths } from './v1';
import { env } from '$env/dynamic/public';

export const api = createClient<paths>({ baseUrl: env.PUBLIC_API_BASE_URL });
