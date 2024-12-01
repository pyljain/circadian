import { Check } from "@/types/monitoring"

export function StatusTimeline({ checks }: { checks: Check[] }) {
  return (
    <div className="flex items-center gap-0.5 h-6">
      {checks.map((check) => (
        <div
          key={check.date}
          className={`w-2 h-full transition-all duration-200 hover:opacity-80 ${
            check.status === 'healthy'
              ? 'bg-green-500'
              : check.status === 'degraded'
              ? 'bg-yellow-500'
              : 'bg-red-500'
          }`}
          title={`Status: ${check.status}
Time: ${new Date(check.date).toDateString()}
Latency: ${Math.ceil(check.averageLatency)}ms
Status Code: ${check.status}`}
        />
      ))}
    </div>
  )
}

