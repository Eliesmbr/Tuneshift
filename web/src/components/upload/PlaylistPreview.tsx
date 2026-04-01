import type { UploadedPlaylist } from "../../types";
import { Card } from "../ui/Card";

interface Props {
  playlists: UploadedPlaylist[];
  totalTracks: number;
}

function formatDuration(tracks: number): string {
  // Rough estimate: average 3.5 min per track
  const totalMinutes = Math.round(tracks * 3.5);
  const hours = Math.floor(totalMinutes / 60);
  const minutes = totalMinutes % 60;
  if (hours > 0) return `~${hours}h ${minutes}m`;
  return `~${minutes}m`;
}

const playlistColors = [
  "from-spotify-green/30 to-tidal-blue/10",
  "from-tidal-blue/30 to-spotify-green/10",
  "from-purple-500/30 to-pink-500/10",
  "from-orange-500/30 to-yellow-500/10",
  "from-pink-500/30 to-red-500/10",
  "from-blue-500/30 to-cyan-500/10",
];

export function PlaylistPreview({ playlists, totalTracks }: Props) {
  return (
    <div className="space-y-4 animate-[slideUp_0.4s_ease-out]">
      <div className="text-center">
        <p className="text-3xl font-bold bg-gradient-to-r from-spotify-green to-tidal-blue bg-clip-text text-transparent">
          {totalTracks.toLocaleString()} tracks
        </p>
        <p className="text-sm text-surface-200 mt-1">
          {playlists.length} playlist{playlists.length !== 1 ? "s" : ""} - {formatDuration(totalTracks)} estimated
        </p>
      </div>

      <div className="grid gap-3 sm:grid-cols-2">
        {playlists.map((pl, i) => (
          <Card key={pl.name} className="p-4 flex items-center gap-3 animate-[slideUp_0.3s_ease-out] hover:border-surface-700 transition-colors">
            <div className={`h-12 w-12 rounded-xl bg-gradient-to-br ${playlistColors[i % playlistColors.length]} flex items-center justify-center shrink-0`}>
              <svg className="h-6 w-6 text-white/70" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
                <path strokeLinecap="round" strokeLinejoin="round" d="M9 9l10.5-3m0 6.553v3.75a2.25 2.25 0 01-1.632 2.163l-1.32.377a1.803 1.803 0 11-.99-3.467l2.31-.66a2.25 2.25 0 001.632-2.163zm0 0V2.25L9 5.25v10.303m0 0v3.75a2.25 2.25 0 01-1.632 2.163l-1.32.377a1.803 1.803 0 01-.99-3.467l2.31-.66A2.25 2.25 0 009 15.553z" />
              </svg>
            </div>
            <div className="min-w-0 flex-1">
              <p className="font-medium text-sm truncate">{pl.name}</p>
              <p className="text-xs text-surface-200">{pl.track_count} tracks</p>
            </div>
          </Card>
        ))}
      </div>
    </div>
  );
}
