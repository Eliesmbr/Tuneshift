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
  const retriesRef = useRef(0);

  const connectSSE = useCallback((sessionId: string) => {
    const es = new EventSource(api.migrationProgressURL(sessionId));
    eventSourceRef.current = es;

    es.onmessage = (e) => {
      retriesRef.current = 0; // Reset retry count on successful message
      const event: ProgressEvent = JSON.parse(e.data);
      setEvents((prev) => [...prev, event]);

      if (event.type === "progress" && event.total && event.current) {
        const pct = Math.round((event.current / event.total) * 100);
        document.title = `Tuneshift (${pct}%)`;
      }

      if (event.type === "complete" || event.type === "result") {
        document.title = "Tuneshift - Done!";
        setDone(true);
        setRunning(false);
        setTimeout(() => es.close(), 1000);
      }

      if (event.type === "error") {
        document.title = "Tuneshift";
        setError(event.message);
        setDone(true);
        setRunning(false);
        setTimeout(() => es.close(), 1000);
      }
    };

    es.onerror = () => {
      es.close();
      retriesRef.current++;

      // Retry up to 5 times with increasing delay
      if (retriesRef.current <= 5) {
        const delay = Math.min(retriesRef.current * 2000, 10000);
        setTimeout(() => connectSSE(sessionId), delay);
      } else {
        setError("Connection to server lost");
        setRunning(false);
        document.title = "Tuneshift";
      }
    };
  }, []);

  const start = useCallback(async (req: MigrationRequest) => {
    setRunning(true);
    setDone(false);
    setEvents([]);
    setError(null);
    retriesRef.current = 0;

    try {
      const { session_id } = await api.startMigration(req);
      connectSSE(session_id);
    } catch (err) {
      const msg = err instanceof Error ? err.message : "Migration failed";
      if (msg === "upload_expired") {
        sessionStorage.removeItem("tuneshift_upload");
        setError("Upload expired - please re-upload your CSV files.");
      } else {
        setError(msg);
      }
      setRunning(false);
    }
  }, [connectSSE]);

  const cancel = useCallback(() => {
    retriesRef.current = 999; // Prevent reconnect
    eventSourceRef.current?.close();
    setRunning(false);
    document.title = "Tuneshift";
  }, []);

  return { running, done, events, error, start, cancel };
}
