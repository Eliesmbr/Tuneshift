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
  type: "phase" | "progress" | "playlist" | "complete" | "error" | "result" | "not_found";
  message: string;
  current?: number;
  total?: number;
}

export type Step = "upload" | "select" | "connect-tidal" | "migrate" | "done";
