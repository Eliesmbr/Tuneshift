const API_BASE = "/api";

async function fetchJSON<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    credentials: "include",
    ...init,
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error(err.error || res.statusText);
  }
  return res.json();
}

export const api = {
  tidalStatus: () =>
    fetchJSON<{ connected: boolean; user?: { id: string; name: string } }>(
      "/auth/tidal/status",
    ),
  tidalLogout: () => fetchJSON<void>("/auth/tidal/logout", { method: "POST" }),

  uploadCSV: async (files: File[]) => {
    const form = new FormData();
    files.forEach((f) => form.append("files", f));

    const res = await fetch(`${API_BASE}/upload`, {
      method: "POST",
      credentials: "include",
      body: form,
    });
    if (!res.ok) {
      const err = await res.json().catch(() => ({ error: res.statusText }));
      throw new Error(err.error || res.statusText);
    }
    return res.json() as Promise<{
      session_id: string;
      playlists: Array<{ name: string; track_count: number }>;
      total_tracks: number;
    }>;
  },

  uploadTakeout: async (file: File) => {
    const form = new FormData();
    form.append("file", file);

    const res = await fetch(`${API_BASE}/upload/takeout`, {
      method: "POST",
      credentials: "include",
      body: form,
    });
    if (!res.ok) {
      const err = await res.json().catch(() => ({ error: res.statusText }));
      throw new Error(err.error || res.statusText);
    }
    return res.json() as Promise<{
      session_id: string;
      playlists: Array<{ name: string; track_count: number }>;
      total_tracks: number;
    }>;
  },

  startMigration: (body: {
    upload_session_id: string;
    playlists: string[];
  }) =>
    fetchJSON<{ session_id: string }>("/migrate", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
    }),

  migrationProgressURL: (sessionId: string) =>
    `${API_BASE}/migrate/progress?session_id=${sessionId}`,
};
