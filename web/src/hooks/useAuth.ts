import { useCallback, useEffect, useState } from "react";
import { api } from "../api/client";
import type { AuthStatus } from "../types";

export function useAuth() {
  const [tidal, setTidal] = useState<AuthStatus>({ connected: false });
  const [loading, setLoading] = useState(true);

  const refresh = useCallback(async () => {
    try {
      const t = await api.tidalStatus();
      setTidal({ connected: t.connected, user: t.user });
    } catch {
      // Ignore errors on status check
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    refresh();
  }, [refresh]);

  const connectTidal = () => {
    window.location.href = "/api/auth/tidal/login";
  };

  const disconnectTidal = async () => {
    await api.tidalLogout();
    setTidal({ connected: false });
  };

  return { tidal, loading, connectTidal, disconnectTidal, refresh };
}
