import { useCallback, useRef, useState } from "react";
import { api } from "../api/client";
import type { ProgressEvent } from "../types";

interface MigrationRequest {
  upload_session_id: string;
  playlists: string[];
}

export function useMigration() {
  const [running, setRunning] = useState(false);
  const [done, setDone] = useState(false);
  const [events, setEvents] = useState<ProgressEvent[]>([]);
  const [error, setError] = useState<string | null>(null);
  const eventSourceRef = useRef<EventSource | null>(null);

  const start = useCallback(async (req: MigrationRequest) => {
    setRunning(true);
    setDone(false);
    setEvents([]);
    setError(null);

    try {
      const { session_id } = await api.startMigration(req);

      const es = new EventSource(api.migrationProgressURL(session_id));
      eventSourceRef.current = es;

      es.onmessage = (e) => {
        const event: ProgressEvent = JSON.parse(e.data);
        setEvents((prev) => [...prev, event]);

        if (event.type === "complete" || event.type === "result") {
          setDone(true);
          setRunning(false);
          setTimeout(() => es.close(), 1000);
        }

        if (event.type === "error") {
          setError(event.message);
          setDone(true);
          setRunning(false);
          setTimeout(() => es.close(), 1000);
        }
      };

      es.onerror = () => {
        setError("Connection to server lost");
        setRunning(false);
        es.close();
      };
    } catch (err) {
      const msg = err instanceof Error ? err.message : "Migration failed";
      if (msg === "upload_expired") {
        // Server lost the upload data (e.g. after rebuild) — clear state and reload
        sessionStorage.removeItem("tuneshift_upload");
        setError("Upload expired — please re-upload your CSV files.");
      } else {
        setError(msg);
      }
      setRunning(false);
    }
  }, []);

  const cancel = useCallback(() => {
    eventSourceRef.current?.close();
    setRunning(false);
  }, []);

  return { running, done, events, error, start, cancel };
}
