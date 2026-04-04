import { useCallback, useEffect, useState } from "react";
import { api } from "../../api/client";
import type { YouTubePlaylist, UploadedPlaylist } from "../../types";
import { Card } from "../ui/Card";
import { toast } from "../ui/Toast";

interface Props {
  onFetched: (sessionId: string, playlists: UploadedPlaylist[], totalTracks: number) => void;
}

export function PlaylistPicker({ onFetched }: Props) {
  const [playlists, setPlaylists] = useState<YouTubePlaylist[]>([]);
  const [selected, setSelected] = useState<Set<string>>(new Set());
  const [loading, setLoading] = useState(true);
  const [fetching, setFetching] = useState(false);

  useEffect(() => {
    api
      .youtubeListPlaylists()
      .then((res) => {
        setPlaylists(res.playlists);
        setSelected(new Set(res.playlists.map((p) => p.id)));
      })
      .catch((err) => {
        toast(err instanceof Error ? err.message : "Failed to load playlists", "error");
      })
      .finally(() => setLoading(false));
  }, []);

  const toggle = (id: string) => {
    setSelected((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      return next;
    });
  };

  const handleFetch = useCallback(async () => {
    const selectedPlaylists = playlists.filter((p) => selected.has(p.id));
    if (selectedPlaylists.length === 0) return;

    setFetching(true);
    try {
      const result = await api.youtubeFetchPlaylists(selectedPlaylists);
      toast(
        `${result.playlists.length} playlist${result.playlists.length !== 1 ? "s" : ""} loaded`,
        "success",
      );
      onFetched(result.session_id, result.playlists, result.total_tracks);
    } catch (err) {
      toast(err instanceof Error ? err.message : "Failed to fetch playlists", "error");
    } finally {
      setFetching(false);
    }
  }, [playlists, selected, onFetched]);

  if (loading) {
    return (
      <div className="flex flex-col items-center justify-center gap-3 py-16">
        <div className="h-8 w-8 animate-spin rounded-full border-2 border-surface-700 border-t-red-400" />
        <p className="text-sm text-surface-200">Loading your YouTube Music playlists...</p>
      </div>
    );
  }

  if (playlists.length === 0) {
    return (
      <Card className="p-6 text-center">
        <p className="text-surface-200">No playlists found on your YouTube Music account.</p>
      </Card>
    );
  }

  return (
    <div className="space-y-4">
      <Card>
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-sm font-semibold text-surface-200 uppercase tracking-wider">
            {selected.size}/{playlists.length} playlists
          </h3>
          <div className="flex gap-2">
            <button
              onClick={() => setSelected(new Set(playlists.map((p) => p.id)))}
              className="text-xs text-red-400 hover:text-red-400/80"
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
              key={pl.id}
              className="flex items-center gap-3 rounded-xl px-3 py-2.5 cursor-pointer hover:bg-surface-800 transition-colors"
            >
              <input
                type="checkbox"
                checked={selected.has(pl.id)}
                onChange={() => toggle(pl.id)}
                className="sr-only peer"
              />
              <div className="h-4 w-4 rounded border-2 border-surface-700 bg-surface-800 peer-checked:border-red-400 peer-checked:bg-red-400 flex items-center justify-center shrink-0">
                {selected.has(pl.id) && (
                  <svg
                    className="h-3 w-3 text-white"
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
              {pl.track_count > 0 && (
                <span className="text-xs text-surface-200 shrink-0">
                  {pl.track_count} tracks
                </span>
              )}
            </label>
          ))}
        </div>
      </Card>

      <button
        onClick={handleFetch}
        disabled={selected.size === 0 || fetching}
        className="w-full rounded-xl bg-gradient-to-r from-red-500 to-tidal-blue py-4 text-base font-bold text-white transition-opacity hover:opacity-90 disabled:opacity-30 disabled:cursor-not-allowed"
      >
        {fetching ? (
          <span className="flex items-center justify-center gap-2">
            <span className="h-4 w-4 animate-spin rounded-full border-2 border-white/30 border-t-white" />
            Fetching tracks...
          </span>
        ) : (
          `Continue with ${selected.size} playlist${selected.size !== 1 ? "s" : ""}`
        )}
      </button>
    </div>
  );
}
