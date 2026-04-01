import { useState } from "react";
import type { UploadedPlaylist } from "../../types";
import { Card } from "../ui/Card";

interface Props {
  playlists: UploadedPlaylist[];
  onSelect: (playlistNames: string[]) => void;
}

export function SelectionPanel({ playlists, onSelect }: Props) {
  const [selected, setSelected] = useState<Set<string>>(
    new Set(playlists.map((p) => p.name)),
  );

  const toggle = (name: string) => {
    setSelected((prev) => {
      const next = new Set(prev);
      if (next.has(name)) next.delete(name);
      else next.add(name);
      return next;
    });
  };

  const selectedTracks = playlists
    .filter((p) => selected.has(p.name))
    .reduce((sum, p) => sum + p.track_count, 0);

  return (
    <div className="space-y-4">
      <Card>
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-sm font-semibold text-surface-200 uppercase tracking-wider">
            {selected.size}/{playlists.length} playlists ({selectedTracks.toLocaleString()} tracks)
          </h3>
          <div className="flex gap-2">
            <button
              onClick={() => setSelected(new Set(playlists.map((p) => p.name)))}
              className="text-xs text-spotify-green hover:text-spotify-green/80"
            >
              Select all
            </button>
            <span className="text-surface-700">|</span>
            <button
              onClick={() => setSelected(new Set())}
              className="text-xs text-surface-200 hover:text-white"
            >
              Deselect all
            </button>
          </div>
        </div>
        <div className="max-h-80 overflow-y-auto space-y-1 pr-2">
          {playlists.map((pl) => (
            <label
              key={pl.name}
              className="flex items-center gap-3 rounded-xl px-3 py-2.5 cursor-pointer hover:bg-surface-800 transition-colors"
            >
              <input
                type="checkbox"
                checked={selected.has(pl.name)}
                onChange={() => toggle(pl.name)}
                className="sr-only peer"
              />
              <div className="h-4 w-4 rounded border-2 border-surface-700 bg-surface-800 peer-checked:border-spotify-green peer-checked:bg-spotify-green flex items-center justify-center shrink-0">
                {selected.has(pl.name) && (
                  <svg
                    className="h-3 w-3 text-black"
                    viewBox="0 0 20 20"
                    fill="currentColor"
                  >
                    <path
                      fillRule="evenodd"
                      d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                      clipRule="evenodd"
                    />
                  </svg>
                )}
              </div>
              <div className="flex-1 min-w-0">
                <p className="text-sm font-medium truncate">{pl.name}</p>
              </div>
              <span className="text-xs text-surface-200 shrink-0">
                {pl.track_count} tracks
              </span>
            </label>
          ))}
        </div>
      </Card>

      <button
        onClick={() => onSelect(Array.from(selected))}
        disabled={selected.size === 0}
        className="w-full rounded-xl bg-gradient-to-r from-spotify-green to-tidal-blue py-4 text-base font-bold text-black transition-opacity hover:opacity-90 disabled:opacity-30 disabled:cursor-not-allowed"
      >
        Continue with {selected.size} playlist{selected.size !== 1 ? "s" : ""}
      </button>
    </div>
  );
}
