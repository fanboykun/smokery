export type ProjectSummary = {
  id: string;
  name: string;
  version: string;
  operations: number;
  flows: number;
  suites: number;
  passRate: number;
  lastRun: string;
  status: 'healthy' | 'warning' | 'failing';
};

export type OperationSummary = {
  id: string;
  method: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE';
  path: string;
  tag: string;
  classification: 'list' | 'read' | 'create' | 'update' | 'delete' | 'action';
  destructive: boolean;
  queryHints?: string[];
  responseShape?: string;
};

export type BuilderConfig = {
  environments: Array<{ id: string; name: string; baseUrl: string; kind: 'safe' | 'prod' }>;
  flows: Array<{ id: string; name: string; steps: number; environment: string }>;
  suites: Array<{ id: string; name: string; cases: number; environment: string }>;
};

export const demoProjects: ProjectSummary[] = [
  {
    id: 'payment-api',
    name: 'Payment API',
    version: 'v2.1.0',
    operations: 12,
    flows: 3,
    suites: 1,
    passRate: 94,
    lastRun: '2m ago',
    status: 'healthy',
  },
  {
    id: 'user-service',
    name: 'User Service',
    version: 'v1.4.2',
    operations: 28,
    flows: 5,
    suites: 2,
    passRate: 87,
    lastRun: '1h ago',
    status: 'warning',
  },
  {
    id: 'inventory-api',
    name: 'Inventory API',
    version: 'v3.0.0',
    operations: 45,
    flows: 2,
    suites: 1,
    passRate: 72,
    lastRun: '3h ago',
    status: 'failing',
  },
];

export const demoOperations: OperationSummary[] = [
  {
    id: 'listUsers',
    method: 'GET',
    path: '/users',
    tag: 'users',
    classification: 'list',
    destructive: false,
    queryHints: ['page', 'limit', 'q', 'status'],
    responseShape: '{ data: User[], pagination: PageInfo }',
  },
  {
    id: 'getUser',
    method: 'GET',
    path: '/users/{id}',
    tag: 'users',
    classification: 'read',
    destructive: false,
    responseShape: 'User',
  },
  {
    id: 'createUser',
    method: 'POST',
    path: '/users',
    tag: 'users',
    classification: 'create',
    destructive: true,
    responseShape: 'User',
  },
  {
    id: 'updateUser',
    method: 'PUT',
    path: '/users/{id}',
    tag: 'users',
    classification: 'update',
    destructive: true,
    responseShape: 'User',
  },
  {
    id: 'deleteUser',
    method: 'DELETE',
    path: '/users/{id}',
    tag: 'users',
    classification: 'delete',
    destructive: true,
    responseShape: '204 No Content',
  },
  {
    id: 'listOrders',
    method: 'GET',
    path: '/orders',
    tag: 'orders',
    classification: 'list',
    destructive: false,
    queryHints: ['page', 'limit', 'status'],
    responseShape: '{ data: Order[], pagination: PageInfo }',
  },
  {
    id: 'submitOrder',
    method: 'POST',
    path: '/orders/{id}/submit',
    tag: 'orders',
    classification: 'action',
    destructive: true,
    responseShape: 'Order',
  },
];

export const demoBuilderConfig: BuilderConfig = {
  environments: [
    { id: 'staging', name: 'staging', baseUrl: 'https://staging.example.test', kind: 'safe' },
    { id: 'production', name: 'production', baseUrl: 'https://api.example.com', kind: 'prod' },
  ],
  flows: [
    { id: 'user-crud', name: 'User CRUD', steps: 4, environment: 'staging' },
    { id: 'order-flow', name: 'Order Flow', steps: 3, environment: 'staging' },
  ],
  suites: [{ id: 'list-endpoints', name: 'List Endpoints', cases: 14, environment: 'staging' }],
};
