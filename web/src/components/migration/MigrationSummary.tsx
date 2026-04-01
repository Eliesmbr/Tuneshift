import { useCallback, useState } from "react";
import type { ProgressEvent } from "../../types";
import { Card } from "../ui/Card";
import { Button } from "../ui/Button";

interface Props {
  events: ProgressEvent[];
  onStartOver: () => void;
}

export function MigrationSummary({ events, onStartOver }: Props) {
  const playlistEvents = events.filter((e) => e.type === "playlist");
  const duplicateEvents = events.filter((e) => e.type === "duplicate");
  const resultEvent = events.find((e) => e.type === "result");
  const notFoundEvents = events.filter((e) => e.type === "not_found");
  const hasError = events.some((e) => e.type === "error");
  const [showNotFound, setShowNotFound] = useState(false);

  const downloadNotFound = useCallback(() => {
    const csv = ["Track,Artist"]
      .concat(notFoundEvents.map((e) => {
        const escaped = e.message.replace(/"/g, '""');
        return `"${escaped}"`;
      }))
      .join("\n");
    const blob = new Blob([csv], { type: "text/csv" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = "tuneshift-not-found.csv";
    a.click();
    URL.revokeObjectURL(url);
  }, [notFoundEvents]);

  return (
    <div className="space-y-6">
      <Card className="text-center py-8">
        <div className={`mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full ${
          hasError
            ? "bg-red-500/20"
            : "bg-gradient-to-br from-spotify-green to-tidal-blue"
        }`}>
          {hasError ? (
            <svg className="h-8 w-8 text-red-400" viewBox="0 0 20 20" fill="currentColor">
              <path fillRule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clipRule="evenodd" />
            </svg>
          ) : (
            <svg className="h-8 w-8 text-black" viewBox="0 0 20 20" fill="currentColor">
              <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
            </svg>
          )}
        </div>
        <h2 className="text-2xl font-bold mb-2">
          {hasError ? "Migration had errors" : "Migration Complete!"}
        </h2>
        <p className="text-surface-200">
          {hasError
            ? "Some items could not be transferred."
            : "Your music has been transferred to Tidal."
          }
        </p>
      </Card>

      {resultEvent && (
        <Card>
          <h3 className="text-sm font-semibold text-surface-200 uppercase tracking-wider mb-3">
            Results
          </h3>
          <div className="flex items-center gap-2 text-sm">
            <span className="text-tidal-blue font-mono">~</span>
            <span>{resultEvent.message}</span>
          </div>
          {resultEvent.total != null && resultEvent.current != null && (
            <div className="mt-3 h-2 w-full rounded-full bg-surface-800 overflow-hidden">
              <div
                className="h-full rounded-full bg-gradient-to-r from-spotify-green to-tidal-blue"
                style={{ width: `${Math.round((resultEvent.current / resultEvent.total) * 100)}%` }}
              />
            </div>
          )}
        </Card>
      )}

      {playlistEvents.length > 0 && (
        <Card>
          <h3 className="text-sm font-semibold text-surface-200 uppercase tracking-wider mb-3">
            Playlists
          </h3>
          <div className="space-y-2 text-sm">
            {playlistEvents.map((e, i) => (
              <div key={i} className="flex items-center gap-2">
                <span className="text-spotify-green">+</span>
                <span>{e.message}</span>
              </div>
            ))}
          </div>
        </Card>
      )}

      {duplicateEvents.length > 0 && (
        <Card>
          <h3 className="text-sm font-semibold text-yellow-400 uppercase tracking-wider mb-3">
            Skipped - already on Tidal ({duplicateEvents.length})
          </h3>
          <div className="space-y-2 text-sm">
            {duplicateEvents.map((e, i) => (
              <div key={i} className="flex items-center gap-2">
                <span className="text-yellow-400">~</span>
                <span className="text-surface-200">{e.message}</span>
              </div>
            ))}
          </div>
        </Card>
      )}

      {notFoundEvents.length > 0 && (
        <Card>
          <div className="flex items-center justify-between mb-1">
            <button
              onClick={() => setShowNotFound(!showNotFound)}
              className="flex items-center gap-2"
            >
              <h3 className="text-sm font-semibold text-surface-200 uppercase tracking-wider">
                Not found ({notFoundEvents.length} tracks)
              </h3>
              <svg
                className={`h-4 w-4 text-surface-200 transition-transform ${showNotFound ? "rotate-180" : ""}`}
                viewBox="0 0 20 20"
                fill="currentColor"
              >
                <path fillRule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clipRule="evenodd" />
              </svg>
            </button>
            <button
              onClick={downloadNotFound}
              className="text-xs text-tidal-blue hover:text-tidal-blue/80"
            >
              Download CSV
            </button>
          </div>
          {showNotFound && (
            <div className="mt-3 max-h-64 overflow-y-auto space-y-1 text-sm">
              {notFoundEvents.map((e, i) => (
                <div key={i} className="flex items-center gap-2 py-0.5">
                  <span className="text-red-400 shrink-0">x</span>
                  <span className="text-surface-200">{e.message}</span>
                </div>
              ))}
            </div>
          )}
        </Card>
      )}

      <div className="flex justify-center">
        <Button variant="secondary" onClick={onStartOver}>
          Start Over
        </Button>
      </div>
    </div>
  );
}
