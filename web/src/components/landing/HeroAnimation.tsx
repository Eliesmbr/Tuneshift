import type { Source } from "../../types";

const SpotifyIcon = () => (
  <svg className="h-10 w-10 text-spotify-green" viewBox="0 0 24 24" fill="currentColor">
    <path d="M12 0C5.4 0 0 5.4 0 12s5.4 12 12 12 12-5.4 12-12S18.66 0 12 0zm5.521 17.34c-.24.359-.66.48-1.021.24-2.82-1.74-6.36-2.101-10.561-1.141-.418.122-.779-.179-.899-.539-.12-.421.18-.78.54-.9 4.56-1.021 8.52-.6 11.64 1.32.42.18.479.659.301 1.02zm1.44-3.3c-.301.42-.841.6-1.262.3-3.239-1.98-8.159-2.58-11.939-1.38-.479.12-1.02-.12-1.14-.6-.12-.48.12-1.021.6-1.141C9.6 9.9 15 10.561 18.72 12.84c.361.181.54.78.241 1.2zm.12-3.36C15.24 8.4 8.82 8.16 5.16 9.301c-.6.179-1.2-.181-1.38-.721-.18-.601.18-1.2.72-1.381 4.26-1.26 11.28-1.02 15.721 1.621.539.3.719 1.02.419 1.56-.299.421-1.02.599-1.559.3z" />
  </svg>
);

const YouTubeIcon = () => (
  <svg className="h-10 w-10 text-red-400" viewBox="0 0 24 24" fill="currentColor">
    <path d="M12 0C5.376 0 0 5.376 0 12s5.376 12 12 12 12-5.376 12-12S18.624 0 12 0zm0 19.104c-3.924 0-7.104-3.18-7.104-7.104S8.076 4.896 12 4.896s7.104 3.18 7.104 7.104-3.18 7.104-7.104 7.104zm0-13.332c-3.432 0-6.228 2.796-6.228 6.228S8.568 18.228 12 18.228 18.228 15.432 18.228 12 15.432 5.772 12 5.772zM9.684 15.54V8.46L15.816 12l-6.132 3.54z" />
  </svg>
);

const TidalIcon = () => (
  <svg className="h-10 w-10 text-tidal-blue" viewBox="0 0 24 24" fill="currentColor">
    <path d="M12.012 3.992L8.008 7.996 12.012 12l4.004-4.004zM4.004 7.996L0 12l4.004 4.004L8.008 12zm15.992 0L15.992 12l4.004 4.004L24 12zM12.012 12l-4.004 4.004 4.004 4.004 4.004-4.004z" />
  </svg>
);

const sourceConfig = {
  spotify: {
    icon: <SpotifyIcon />,
    label: "Spotify",
    color: "bg-spotify-green/10",
    gradient: "from-spotify-green",
  },
  "youtube-music": {
    icon: <YouTubeIcon />,
    label: "YouTube Music",
    color: "bg-red-400/10",
    gradient: "from-red-400",
  },
};

interface Props {
  source?: Source | null;
}

export function HeroAnimation({ source }: Props) {
  const config = source ? sourceConfig[source] : null;

  return (
    <div className="flex items-center justify-center gap-6 py-8">
      {/* Source */}
      <div className="flex flex-col items-center gap-2 animate-[fadeInLeft_0.6s_ease-out]">
        <div
          className={`h-16 w-16 rounded-2xl flex items-center justify-center transition-all duration-300 ${
            config ? config.color : "bg-surface-800/50"
          }`}
        >
          {config ? (
            <div key={source} className="animate-[fadeIn_0.3s_ease-out]">
              {config.icon}
            </div>
          ) : (
            <svg className="h-10 w-10 text-surface-600" viewBox="0 0 20 20" fill="currentColor">
              <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-8-3a1 1 0 00-.867.5 1 1 0 11-1.731-1A3 3 0 0113 8a3.001 3.001 0 01-2 2.83V11a1 1 0 11-2 0v-1a1 1 0 011-1 1 1 0 100-2zm0 8a1 1 0 100-2 1 1 0 000 2z" clipRule="evenodd" />
            </svg>
          )}
        </div>
        <span className="text-xs text-surface-200 font-medium">
          {config ? config.label : "Source"}
        </span>
      </div>

      {/* Animated arrow */}
      <div className="flex items-center gap-1 animate-[fadeIn_0.8s_ease-out_0.3s_both]">
        <div
          className={`h-px w-8 bg-gradient-to-r ${
            config ? config.gradient : "from-surface-600"
          } to-transparent transition-all duration-300`}
        />
        <div className="relative">
          <div
            className={`h-10 w-10 rounded-full flex items-center justify-center transition-all duration-300 ${
              config
                ? `bg-gradient-to-r ${config.gradient} to-tidal-blue animate-pulse`
                : "bg-surface-700"
            }`}
          >
            <svg className="h-5 w-5 text-black" viewBox="0 0 20 20" fill="currentColor">
              <path fillRule="evenodd" d="M10.293 3.293a1 1 0 011.414 0l6 6a1 1 0 010 1.414l-6 6a1 1 0 01-1.414-1.414L14.586 11H3a1 1 0 110-2h11.586l-4.293-4.293a1 1 0 010-1.414z" clipRule="evenodd" />
            </svg>
          </div>
        </div>
        <div className="h-px w-8 bg-gradient-to-r from-transparent to-tidal-blue" />
      </div>

      {/* Tidal */}
      <div className="flex flex-col items-center gap-2 animate-[fadeInRight_0.6s_ease-out_0.2s_both]">
        <div className="h-16 w-16 rounded-2xl bg-tidal-blue/10 flex items-center justify-center">
          <TidalIcon />
        </div>
        <span className="text-xs text-surface-200 font-medium">Tidal</span>
      </div>
    </div>
  );
}
