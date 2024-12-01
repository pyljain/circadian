import { StatusTimeline } from "@/components/status-timeline"
import { EndpointStats } from "@/components/endpoint-stats"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { useEffect, useState } from "react"
import { Endpoint } from "@/types/monitoring"

export default function DashboardPage() {
  const [endpoints, setEndpoints] = useState<Endpoint[]>([])

  useEffect(() => {
    fetch("/api/v1/uptime")
      .then((res) => res.json())
      .then((data) => {
        console.log('result from api', data)
        setEndpoints(data)
      })
      .catch((err) => console.error(err))
  }, [])
  

  return (
    <div className="container mx-auto p-6">
      <header className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between mb-8">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Circadian</h1>
          <p className="text-muted-foreground">Monitor your endpoints in real-time</p>
        </div>
        <div className="flex items-center gap-2">
          <Badge variant="outline" className="text-green-500">
            {endpoints.filter((e) => e.currentStatus === 'healthy').length} Healthy
          </Badge>
          <Badge variant="outline" className="text-red-500">
            {endpoints.filter((e) => e.currentStatus === 'down').length} Down
          </Badge>
        </div>
      </header>

      <div className="grid gap-6">
        {endpoints.map((endpoint) => (
          <Card key={endpoint.id}>
            <CardHeader>
              <div className="flex items-center justify-between">
                <div className="space-y-1">
                  <CardTitle>{endpoint.name}</CardTitle>
                  <div className="text-sm text-muted-foreground font-mono">{endpoint.url}</div>
                </div>
                <Badge
                  variant="outline"
                  className={
                    endpoint.currentStatus === 'healthy'
                      ? 'text-green-500'
                      : endpoint.currentStatus === 'degraded'
                      ? 'text-yellow-500'
                      : 'text-red-500'
                  }
                >
                  {endpoint.currentStatus}
                </Badge>
              </div>
            </CardHeader>
            <CardContent className="space-y-6">
              <StatusTimeline checks={endpoint.checks} />
              <EndpointStats
                uptime={endpoint.uptime}
                lastCheck={endpoint.checks[endpoint.checks.length - 1]}
                interval={endpoint.interval}
              />
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  )
}