import type { ProgressEvent } from "../../types";
import { Card } from "../ui/Card";
import { ProgressBar } from "../ui/ProgressBar";

interface Props {
  events: ProgressEvent[];
  running: boolean;
  error: string | null;
}

export function MigrationProgressView({ events, running, error }: Props) {
  const latestProgress = [...events].reverse().find((e) => e.type === "progress" && e.total);
  const phases = events.filter((e) => e.type === "phase" || e.type === "playlist");

  return (
    <div className="space-y-6">
      {latestProgress?.current != null && latestProgress?.total != null && (
        <Card>
          <ProgressBar
            current={latestProgress.current}
            total={latestProgress.total}
          />
        </Card>
      )}

      {running && (
        <div className="flex items-center justify-center gap-3 py-4">
          <div className="h-5 w-5 animate-spin rounded-full border-2 border-surface-700 border-t-spotify-green" />
          <p className="text-sm text-surface-200">Migration in progress...</p>
        </div>
      )}

      {error && (
        <Card className="border-red-900 bg-red-950/50">
          <p className="text-sm text-red-400">{error}</p>
        </Card>
      )}

      <Card>
        <h3 className="text-sm font-semibold text-surface-200 uppercase tracking-wider mb-3">
          Activity Log
        </h3>
        <div className="max-h-64 overflow-y-auto space-y-1 font-mono text-xs">
          {phases.map((e, i) => (
            <div key={i} className="flex gap-2 py-1">
              <span className={`shrink-0 ${
                e.type === "playlist" ? "text-spotify-green" : "text-tidal-blue"
              }`}>
                {e.type === "playlist" ? "+" : ">"}
              </span>
              <span className="text-surface-200">{e.message}</span>
            </div>
          ))}
          {phases.length === 0 && !running && (
            <p className="text-surface-200">No events yet...</p>
          )}
        </div>
      </Card>
    </div>
  );
}
