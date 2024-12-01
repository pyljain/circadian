export type StatusCode = 200 | 201 | 301 | 302 | 400 | 401 | 403 | 404 | 500 | 502 | 503 | 504;

export type EndpointStatus = 'healthy' | 'degraded' | 'down';

export type Check = {
  date: string;
  averageLatency: number; // in milliseconds
  status: EndpointStatus;
};

export type Endpoint = {
  id: string;
  name: string;
  url: string;
  interval: number; // in minutes
  currentStatus: EndpointStatus;
  lastCheck: Check;
  checks: Check[];
  uptime: number; // percentage
};
