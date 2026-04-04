export interface User {
  id: string;
  name: string;
}

export interface AuthStatus {
  connected: boolean;
  user?: User;
}

export interface UploadedPlaylist {
  name: string;
  track_count: number;
}

export interface UploadResult {
  session_id: string;
  playlists: UploadedPlaylist[];
  total_tracks: number;
}

export interface ProgressEvent {
  type: "phase" | "progress" | "playlist" | "complete" | "error" | "result" | "not_found" | "duplicate";
  message: string;
  current?: number;
  total?: number;
}

export type Source = "spotify" | "youtube-music";

export type Step = "upload" | "connect-source" | "fetch-playlists" | "select" | "connect-tidal" | "migrate" | "done";

export interface YouTubePlaylist {
  id: string;
  name: string;
  track_count: number;
}
