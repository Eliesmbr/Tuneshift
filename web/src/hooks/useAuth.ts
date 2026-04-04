import { useCallback, useEffect, useState } from "react";
import { api } from "../api/client";
import type { AuthStatus } from "../types";

export function useAuth() {
  const [tidal, setTidal] = useState<AuthStatus>({ connected: false });
  const [google, setGoogle] = useState<AuthStatus>({ connected: false });
  const [loading, setLoading] = useState(true);

  const refresh = useCallback(async () => {
    try {
      const [t, g] = await Promise.allSettled([
        api.tidalStatus(),
        api.googleStatus(),
      ]);
      if (t.status === "fulfilled") {
        setTidal({ connected: t.value.connected, user: t.value.user });
      }
      if (g.status === "fulfilled") {
        setGoogle({
          connected: g.value.connected,
          user: g.value.user ? { id: "", name: g.value.user.name } : undefined,
        });
      }
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

  const connectGoogle = () => {
    window.location.href = "/api/auth/google/login";
  };

  const disconnectGoogle = async () => {
    await api.googleLogout();
    setGoogle({ connected: false });
  };

  return {
    tidal,
    google,
    loading,
    connectTidal,
    disconnectTidal,
    connectGoogle,
    disconnectGoogle,
    refresh,
  };
}
