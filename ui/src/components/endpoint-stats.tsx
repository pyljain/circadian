import { Card, CardContent } from "@/components/ui/card"
import { Check, EndpointStatus } from "@/types/monitoring"

interface EndpointStatsProps {
  uptime: number
  lastCheck: Check
  interval: number
}

export function EndpointStats({ uptime, lastCheck, interval }: EndpointStatsProps) {
  const getStatusColor = (status: EndpointStatus) => {
    switch (status) {
      case 'healthy':
        return 'text-green-500'
      case 'degraded':
        return 'text-yellow-500'
      case 'down':
        return 'text-red-500'
    }
  }

  return (
    <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
      <Card>
        <CardContent className="pt-6">
          <div className="text-sm font-medium text-muted-foreground">Uptime</div>
          <div className="text-2xl font-bold">{Math.ceil(uptime)}%</div>
        </CardContent>
      </Card>
      <Card>
        <CardContent className="pt-6">
          <div className="text-sm font-medium text-muted-foreground">Latency (ms)</div>
          <div className="text-2xl font-bold">{Math.ceil(lastCheck.averageLatency)}</div>
        </CardContent>
      </Card>
      <Card>
        <CardContent className="pt-6">
          <div className="text-sm font-medium text-muted-foreground">Status</div>
          <div className={`text-2xl font-bold ${getStatusColor(lastCheck.status)}`}>
            {lastCheck.status}
          </div>
        </CardContent>
      </Card>
      <Card>
        <CardContent className="pt-6">
          <div className="text-sm font-medium text-muted-foreground">Interval</div>
          <div className="text-2xl font-bold">{Math.ceil(interval / 60)}m</div>
        </CardContent>
      </Card>
    </div>
  )
}

