import { Check, Endpoint, StatusCode, EndpointStatus } from "@/types/monitoring"

function generateChecks(count: number): Check[] {
  const checks: Check[] = []
  const now = new Date()
  
  for (let i = 0; i < count; i++) {
    const statusCodes: StatusCode[] = [200, 200, 200, 200, 301, 404, 500]
    const statusCode = statusCodes[Math.floor(Math.random() * statusCodes.length)]
    
    let status: EndpointStatus = 'healthy'
    if (statusCode >= 500) status = 'down'
    else if (statusCode >= 400) status = 'degraded'
    
    checks.push({
      timestamp: new Date(now.getTime() - (i * 5 * 60 * 1000)).toISOString(),
      statusCode,
      latency: Math.floor(Math.random() * 1000),
      status,
    })
  }
  
  return checks.reverse()
}

export function getMockEndpoints(): Endpoint[] {
  return [
    {
      id: '1',
      name: 'Main API',
      url: 'https://api.example.com/v1',
      interval: 5,
      currentStatus: 'healthy',
      lastCheck: generateChecks(1)[0],
      checks: generateChecks(50),
      uptime: 99.9,
    },
    {
      id: '2',
      name: 'Authentication Service',
      url: 'https://auth.example.com',
      interval: 1,
      currentStatus: 'degraded',
      lastCheck: generateChecks(1)[0],
      checks: generateChecks(50),
      uptime: 95.5,
    },
    {
      id: '3',
      name: 'Payment Gateway',
      url: 'https://payments.example.com/status',
      interval: 1,
      currentStatus: 'down',
      lastCheck: generateChecks(1)[0],
      checks: generateChecks(50),
      uptime: 85.2,
    },
    {
      id: '4',
      name: 'Storage Service',
      url: 'https://storage.example.com/health',
      interval: 5,
      currentStatus: 'healthy',
      lastCheck: generateChecks(1)[0],
      checks: generateChecks(50),
      uptime: 99.99,
    },
  ]
}

