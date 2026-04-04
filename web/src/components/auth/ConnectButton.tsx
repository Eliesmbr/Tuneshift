import { Button } from "../ui/Button";

interface Props {
  service: "spotify" | "tidal" | "youtube-music";
  connected: boolean;
  userName?: string;
  onConnect: () => void;
  onDisconnect: () => void;
}

const SpotifyIcon = () => (
  <svg className="h-5 w-5" viewBox="0 0 24 24" fill="currentColor">
    <path d="M12 0C5.4 0 0 5.4 0 12s5.4 12 12 12 12-5.4 12-12S18.66 0 12 0zm5.521 17.34c-.24.359-.66.48-1.021.24-2.82-1.74-6.36-2.101-10.561-1.141-.418.122-.779-.179-.899-.539-.12-.421.18-.78.54-.9 4.56-1.021 8.52-.6 11.64 1.32.42.18.479.659.301 1.02zm1.44-3.3c-.301.42-.841.6-1.262.3-3.239-1.98-8.159-2.58-11.939-1.38-.479.12-1.02-.12-1.14-.6-.12-.48.12-1.021.6-1.141C9.6 9.9 15 10.561 18.72 12.84c.361.181.54.78.241 1.2zm.12-3.36C15.24 8.4 8.82 8.16 5.16 9.301c-.6.179-1.2-.181-1.38-.721-.18-.601.18-1.2.72-1.381 4.26-1.26 11.28-1.02 15.721 1.621.539.3.719 1.02.419 1.56-.299.421-1.02.599-1.559.3z" />
  </svg>
);

const TidalIcon = () => (
  <svg className="h-5 w-5" viewBox="0 0 24 24" fill="currentColor">
    <path d="M12.012 3.992L8.008 7.996 12.012 12l4.004-4.004zM4.004 7.996L0 12l4.004 4.004L8.008 12zm15.992 0L15.992 12l4.004 4.004L24 12zM12.012 12l-4.004 4.004 4.004 4.004 4.004-4.004z" />
  </svg>
);

const YouTubeIcon = () => (
  <svg className="h-5 w-5" viewBox="0 0 24 24" fill="currentColor">
    <path d="M12 0C5.376 0 0 5.376 0 12s5.376 12 12 12 12-5.376 12-12S18.624 0 12 0zm0 19.104c-3.924 0-7.104-3.18-7.104-7.104S8.076 4.896 12 4.896s7.104 3.18 7.104 7.104-3.18 7.104-7.104 7.104zm0-13.332c-3.432 0-6.228 2.796-6.228 6.228S8.568 18.228 12 18.228 18.228 15.432 18.228 12 15.432 5.772 12 5.772zM9.684 15.54V8.46L15.816 12l-6.132 3.54z" />
  </svg>
);

const serviceConfig = {
  spotify: { label: "Spotify", icon: <SpotifyIcon />, variant: "spotify" as const },
  tidal: { label: "Tidal", icon: <TidalIcon />, variant: "tidal" as const },
  "youtube-music": { label: "YouTube Music", icon: <YouTubeIcon />, variant: "youtube" as const },
};

export function ConnectButton({ service, connected, userName, onConnect, onDisconnect }: Props) {
  const config = serviceConfig[service];

  if (connected) {
    return (
      <div className="flex items-center gap-4 rounded-2xl border border-surface-800 bg-surface-900/50 px-5 py-4">
        <div className="flex-1">
          <p className="text-xs text-surface-200">Connected to {config.label}</p>
          <p className="font-semibold">{userName}</p>
        </div>
        <Button variant="secondary" onClick={onDisconnect} className="text-xs px-3 py-1.5">
          Disconnect
        </Button>
      </div>
    );
  }

  return (
    <Button
      variant={config.variant}
      onClick={onConnect}
      icon={config.icon}
      className="w-full py-4 text-base"
    >
      Connect {config.label}
    </Button>
  );
}
